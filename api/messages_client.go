package api

import (
	"fmt"
	"github.com/BOOMfinity/bfcord/discord"
	"github.com/BOOMfinity/bfcord/internal/httpc"
	"github.com/BOOMfinity/go-utils/inlineif"
	"github.com/andersfylling/snowflake/v5"
	"net/url"
)

type MessagesResolver struct {
	Channel snowflake.ID
	client  *client
}

func (m MessagesResolver) fetch(around, before, after snowflake.ID, limit uint) ([]discord.Message, error) {
	buff := make([]discord.Message, 0, limit)
	values := url.Values{}
	if limit != 0 {
		values.Set("limit", fmt.Sprint(inlineif.IfElse(limit > 100, 100, limit)))
	}
	for uint(len(buff)) < limit {
		if around.Valid() {
			values.Set("around", around.String())
		}
		if before.Valid() {
			values.Set("before", before.String())
		}
		if after.Valid() {
			values.Set("after", after.String())
		}
		messages, err := httpc.NewJSONRequest[[]discord.Message](m.client.http, func(b httpc.RequestBuilder) error {
			return b.Execute("channels", m.Channel.String(), "messages?"+values.Encode())
		})
		if err != nil {
			return buff, fmt.Errorf("failed to make request: %w", err)
		}
		buff = append(buff, messages...)
		if uint(len(messages)) < inlineif.IfElse(limit > 100, 100, limit) {
			break
		}
		if len(messages) > 0 {
			before = messages[len(messages)-1].ID
		}
		around = 0
		after = 0

	}
	return buff, nil
}

func (m MessagesResolver) Latest(limit uint) ([]discord.Message, error) {
	return m.fetch(0, 0, 0, limit)
}

func (m MessagesResolver) Before(id snowflake.ID, limit uint) ([]discord.Message, error) {
	return m.fetch(0, id, 0, limit)
}

func (m MessagesResolver) After(id snowflake.ID, limit uint) ([]discord.Message, error) {
	return m.fetch(0, 0, id, limit)
}

func (m MessagesResolver) Around(id snowflake.ID, limit uint) ([]discord.Message, error) {
	return m.fetch(id, 0, 0, limit)
}
