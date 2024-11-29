package api

import (
	"fmt"
	"github.com/BOOMfinity/bfcord/internal/httpc"
	"github.com/BOOMfinity/go-utils/inlineif"
	"github.com/andersfylling/snowflake/v5"
	"net/url"
)

type GuildEventQueryResolver struct {
	client *client
	Guild  snowflake.ID
	Event  snowflake.ID
}

func (g GuildEventQueryResolver) fetch(after, before snowflake.ID, withMember bool, limit uint) ([]ScheduledEventUser, error) {
	values := url.Values{}
	fetchAll := limit == 0
	if fetchAll {
		limit = 100
	}
	fetchLimit := inlineif.IfElse(limit > 100, 100, limit)
	values.Set("limit", fmt.Sprint(fetchLimit))
	if withMember {
		values.Set("with_member", "true")
	}
	data := make([]ScheduledEventUser, 0, inlineif.IfElse(limit > 100, limit, fetchLimit))

	for {
		if before.Valid() {
			values.Set("before", before.String())
		}
		if after.Valid() {
			values.Set("after", after.String())
		}
		users, err := httpc.NewJSONRequest[[]ScheduledEventUser](g.client.http, func(b httpc.RequestBuilder) error {
			return b.Execute("guilds", g.Guild.String(), "scheduled-events", g.Event.String(), "users"+values.Encode())
		})
		if err != nil {
			return data, err
		}
		data = append(data, users...)
		if (len(users) < int(fetchLimit)) || ((len(data) >= int(limit)) && !fetchAll) {
			break
		}
		if before.Valid() {
			before = users[0].User.ID
		} else {
			after = users[len(users)-1].User.ID
		}
	}

	return data, nil
}

func (g GuildEventQueryResolver) Before(id snowflake.ID, withMember bool, limit uint) ([]ScheduledEventUser, error) {
	return g.fetch(0, id, withMember, limit)
}

func (g GuildEventQueryResolver) After(id snowflake.ID, withMember bool, limit uint) ([]ScheduledEventUser, error) {
	return g.fetch(id, 0, withMember, limit)
}

func (g GuildEventQueryResolver) Latest(withMember bool, limit uint) ([]ScheduledEventUser, error) {
	return g.fetch(0, 0, withMember, limit)
}
