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

type MessageResolver struct {
	Channel snowflake.ID
	Message snowflake.ID
	client  *client
}

func (v MessageResolver) Answer(id uint, params MessagePollVotersParams) ([]discord.User, error) {
	values := url.Values{}
	users := make([]discord.User, 0, inlineif.IfElse(params.Limit == 0, 100, params.Limit))
	fetchLimit := inlineif.IfElse(params.Limit == 0 || params.Limit > 100, 100, params.Limit)
	values.Set("limit", fmt.Sprint(fetchLimit))

	for {
		if params.After.Valid() {
			values.Set("after", params.After.String())
		}
		fetched, err := httpc.NewJSONRequest[MessagePollVoters](v.client.http, func(b httpc.RequestBuilder) error {
			return b.Execute("channels", v.Channel.String(), "polls", v.Message.String(), "answers", fmt.Sprint(id)+"?"+values.Encode())
		})
		if err != nil {
			return nil, err
		}
		users = append(users, fetched.Users...)
		if len(fetched.Users) < int(fetchLimit) || len(users) > int(params.Limit) {
			break
		}
		if len(fetched.Users) > 0 {
			params.After = fetched.Users[len(fetched.Users)-1].ID
		}
	}

	return users, nil
}

func (v MessageResolver) EndPoll() (discord.Message, error) {
	return httpc.NewJSONRequest[discord.Message](v.client.http, func(b httpc.RequestBuilder) error {
		b.Method(fasthttp.MethodPost)
		return b.Execute("channels", v.Channel.String(), "polls", v.Message.String(), "expire")
	})
}

func (v MessageResolver) Reaction(emoji string) ReactionClient {
	o := ReactionResolver{
		Channel: v.Channel,
		Message: v.Message,
		Emoji:   emoji,
		client:  v.client,
	}
	return o
}

func (v MessageResolver) StartThread(data StartThreadParams, reason ...string) (discord.Channel, error) {
	return httpc.NewJSONRequest[discord.Channel](v.client.http, func(b httpc.RequestBuilder) error {
		b.Method(fasthttp.MethodPost)
		b.Body(data)
		b.Reason(reason...)
		return b.Execute("channels", v.Channel.String(), "messages", v.Message.String(), "threads")
	})
}

func (v MessageResolver) Get() (dst discord.Message, _ error) {
	return httpc.NewJSONRequest[discord.Message](v.client.http, func(b httpc.RequestBuilder) error {
		return b.Execute("channels", v.Channel.String(), "messages", v.Message.String())
	})
}

func (v MessageResolver) Delete(reason ...string) error {
	return httpc.NewRequest(v.client.http, func(req httpc.RequestBuilder) error {
		req.Method(fasthttp.MethodDelete)
		req.Reason(reason...)
		return req.Execute("channels", v.Channel.String(), "messages", v.Message.String())
	})
}

func (v MessageResolver) Pin(reason ...string) error {
	return httpc.NewRequest(v.client.http, func(req httpc.RequestBuilder) error {
		req.Method(fasthttp.MethodPut)
		req.Reason(reason...)
		return req.Execute("channels", v.Channel.String(), "pins", v.Message.String())
	})
}

func (v MessageResolver) Unpin(reason ...string) error {
	return httpc.NewRequest(v.client.http, func(req httpc.RequestBuilder) error {
		req.Method(fasthttp.MethodDelete)
		req.Reason(reason...)
		return req.Execute("channels", v.Channel.String(), "pins", v.Message.String())
	})
}

func (v MessageResolver) Update(params EditMessageParams) (discord.Message, error) {
	return httpc.NewJSONRequest[discord.Message](v.client.http, func(b httpc.RequestBuilder) error {
		uploadFiles(b, params, params.Attachments)
		return b.Execute("channels", v.Channel.String(), "messages", v.Message.String())
	})
}

func (v MessageResolver) CrossPost() error {
	return httpc.NewRequest(v.client.http, func(b httpc.RequestBuilder) error {
		b.Method(fasthttp.MethodPost)
		return b.Execute("channels", v.Channel.String(), "messages", v.Message.String(), "crosspost")
	})
}

func (v MessageResolver) DeleteAllReactions() error {
	return httpc.NewRequest(v.client.http, func(b httpc.RequestBuilder) error {
		b.Method(fasthttp.MethodDelete)
		return b.Execute("channels", v.Channel.String(), "messages", v.Message.String(), "reactions")
	})
}
