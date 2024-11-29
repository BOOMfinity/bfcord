package api

import (
	"fmt"
	"github.com/BOOMfinity/bfcord/discord"
	"github.com/BOOMfinity/bfcord/internal/httpc"
	"github.com/BOOMfinity/go-utils/inlineif"
	"github.com/andersfylling/snowflake/v5"
	"github.com/valyala/fasthttp"
	"net/url"
)

type ReactionResolver struct {
	client  *client
	Channel snowflake.ID
	Message snowflake.ID
	Emoji   string
}

func (r ReactionResolver) Reactions(opts MessageReactionsParams) (_ []discord.User, err error) {
	values := url.Values{}
	users := make([]discord.User, 0, inlineif.IfElse(opts.Limit == 0, 100, opts.Limit))
	fetched := make([]discord.User, 0, inlineif.IfElse(opts.Limit == 0 || opts.Limit > 100, 100, opts.Limit))
	fetchLimit := inlineif.IfElse(opts.Limit == 0 || opts.Limit > 100, 100, opts.Limit)
	values.Set("limit", fmt.Sprint(fetchLimit))

	for {
		if opts.After.Valid() {
			values.Set("after", opts.After.String())
		}
		fetched, err = httpc.NewJSONRequest[[]discord.User](r.client.http, func(b httpc.RequestBuilder) error {
			return b.Execute("channels", r.Channel.String(), "messages", r.Message.String(), "reactions", url.PathEscape(r.Emoji)+"?"+values.Encode())
		})
		if err != nil {
			return nil, fmt.Errorf("failed to fetch reaction users: %w", err)
		}
		users = append(users, fetched...)
		if opts.Limit != 0 && len(users) >= int(opts.Limit) {
			break
		}
		if len(fetched) < int(fetchLimit) {
			break
		}
		if len(fetched) > 0 {
			opts.After = fetched[len(fetched)-1].ID
		}
	}
	return users, nil
}

func (r ReactionResolver) React() error {
	return httpc.NewRequest(r.client.http, func(b httpc.RequestBuilder) error {
		b.Method(fasthttp.MethodPut)
		return b.Execute("channels", r.Channel.String(), "messages", r.Message.String(), "reactions", url.PathEscape(r.Emoji), "@me")
	})
}

func (r ReactionResolver) DeleteOwn() error {
	return httpc.NewRequest(r.client.http, func(b httpc.RequestBuilder) error {
		b.Method(fasthttp.MethodDelete)
		return b.Execute("channels", r.Channel.String(), "messages", r.Message.String(), "reactions", url.PathEscape(r.Emoji), "@me")
	})
}

func (r ReactionResolver) Delete(user snowflake.ID) error {
	return httpc.NewRequest(r.client.http, func(b httpc.RequestBuilder) error {
		b.Method(fasthttp.MethodDelete)
		return b.Execute("channels", r.Channel.String(), "messages", r.Message.String(), "reactions", url.PathEscape(r.Emoji), user.String())
	})
}

func (r ReactionResolver) DeleteAll() error {
	return httpc.NewRequest(r.client.http, func(b httpc.RequestBuilder) error {
		b.Method(fasthttp.MethodDelete)
		return b.Execute("channels", r.Channel.String(), "messages", r.Message.String(), "reactions", url.PathEscape(r.Emoji))
	})
}
