package api

import (
	"fmt"
	"github.com/BOOMfinity/bfcord/discord"
	"github.com/BOOMfinity/bfcord/internal/httpc"
	"github.com/andersfylling/snowflake/v5"
	"github.com/valyala/fasthttp"
)

type InteractionResolver struct {
	Token  string
	client *client
	ID     snowflake.ID
}

func (i InteractionResolver) Pong() error {
	return httpc.NewRequest(i.client.http, func(b httpc.RequestBuilder) error {
		b.Body(map[string]any{
			"type": discord.InteractionCallbackPong,
		})
		b.Method(fasthttp.MethodPost)
		return b.Execute("interactions", i.ID.String(), i.Token, "callback")
	})
}

func (i InteractionResolver) Response() InteractionResponseClient {
	return InteractionResponseResolver{
		client: i.client,
		Token:  i.Token,
	}
}

func (i InteractionResolver) Update(params EditMessageParams) error {
	return httpc.NewRequest(i.client.http, func(b httpc.RequestBuilder) error {
		b.Body(map[string]any{
			"type": discord.InteractionCallbackUpdateMessage,
			"data": params,
		})
		b.Method(fasthttp.MethodPost)
		return b.Execute("interactions", i.ID.String(), i.Token, "callback")
	})
}

func (i InteractionResolver) DeferredUpdate() error {
	return httpc.NewRequest(i.client.http, func(b httpc.RequestBuilder) error {
		b.Body(map[string]any{
			"type": discord.InteractionCallbackDeferredChannelMessage,
		})
		b.Method(fasthttp.MethodPost)
		return b.Execute("interactions", i.ID.String(), i.Token, "callback")
	})
}

func (i InteractionResolver) DeferredComponentUpdate() error {
	return httpc.NewRequest(i.client.http, func(b httpc.RequestBuilder) error {
		b.Body(map[string]any{
			"type": discord.InteractionCallbackDeferredUpdateMessage,
		})
		b.Method(fasthttp.MethodPost)
		return b.Execute("interactions", i.ID.String(), i.Token, "callback")
	})
}

func (i InteractionResolver) Reply(params InteractionMessageParams) error {

	return httpc.NewRequest(i.client.http, func(b httpc.RequestBuilder) error {
		data := map[string]any{
			"type": discord.InteractionCallbackChannelMessage,
			"data": params,
		}
		uploadFiles(b, data, params.Attachments)
		b.Method(fasthttp.MethodPost)
		return b.Execute("interactions", i.ID.String(), i.Token, "callback")
	})
}

func (i InteractionResolver) SendFollowUp(params FollowUpParams) (discord.Message, error) {
	user, err := i.client.GetCurrentUser()
	if err != nil {
		return discord.Message{}, fmt.Errorf("failed to get current user: %w", err)
	}
	return httpc.NewJSONRequest[discord.Message](i.client.http, func(b httpc.RequestBuilder) error {
		b.Method(fasthttp.MethodPost)
		uploadFiles(b, params, params.Attachments)
		return b.Execute("webhooks", user.ID.String(), i.Token)
	})
}

func (i InteractionResolver) FollowUp(id snowflake.ID) FollowUpClient {
	return FollowUpResolver{
		client: i.client,
		ID:     id,
		Token:  i.Token,
	}
}

func (i InteractionResolver) AutoComplete(choices []discord.CommandChoice) error {
	return httpc.NewRequest(i.client.http, func(b httpc.RequestBuilder) error {
		b.Body(map[string]any{
			"type": discord.InteractionCallbackAutoCompleteResult,
			"data": map[string]any{
				"choices": choices,
			},
		})
		b.Method(fasthttp.MethodPost)
		return b.Execute("interactions", i.ID.String(), i.Token, "callback")
	})
}

func (i InteractionResolver) TextInput(params TextInputParams) error {
	return httpc.NewRequest(i.client.http, func(b httpc.RequestBuilder) error {
		b.Body(map[string]any{
			"type": discord.InteractionCallbackModal,
			"data": params,
		})
		b.Method(fasthttp.MethodPost)
		return b.Execute("interactions", i.ID.String(), i.Token, "callback")
	})
}

type InteractionResponseResolver struct {
	client *client
	Token  string
}

func (i InteractionResponseResolver) Get() (discord.Message, error) {
	user, err := i.client.GetCurrentUser()
	if err != nil {
		return discord.Message{}, fmt.Errorf("failed to get current user: %w", err)
	}
	return httpc.NewJSONRequest[discord.Message](i.client.http, func(b httpc.RequestBuilder) error {
		return b.Execute("webhooks", user.ID.String(), i.Token, "messages", "@original")
	})
}

func (i InteractionResponseResolver) Delete() error {
	user, err := i.client.GetCurrentUser()
	if err != nil {
		return fmt.Errorf("failed to get current user: %w", err)
	}
	return httpc.NewRequest(i.client.http, func(b httpc.RequestBuilder) error {
		b.Method(fasthttp.MethodDelete)
		return b.Execute("webhooks", user.ID.String(), i.Token, "messages", "@original")
	})
}

func (i InteractionResponseResolver) Edit(params EditMessageParams) (discord.Message, error) {
	user, err := i.client.GetCurrentUser()
	if err != nil {
		return discord.Message{}, fmt.Errorf("failed to get current user: %w", err)
	}
	return httpc.NewJSONRequest[discord.Message](i.client.http, func(b httpc.RequestBuilder) error {
		b.Method(fasthttp.MethodPatch)
		uploadFiles(b, params, params.Attachments)
		return b.Execute("webhooks", user.ID.String(), i.Token, "messages", "@original")
	})
}

type FollowUpResolver struct {
	client *client
	ID     snowflake.ID
	Token  string
}

func (f FollowUpResolver) Get() (discord.Message, error) {
	user, err := f.client.GetCurrentUser()
	if err != nil {
		return discord.Message{}, fmt.Errorf("failed to get current user: %w", err)
	}
	return httpc.NewJSONRequest[discord.Message](f.client.http, func(b httpc.RequestBuilder) error {
		return b.Execute("webhooks", user.ID.String(), f.Token, "messages", f.ID.String())
	})
}

func (f FollowUpResolver) Update(params EditMessageParams) (discord.Message, error) {
	user, err := f.client.GetCurrentUser()
	if err != nil {
		return discord.Message{}, fmt.Errorf("failed to get current user: %w", err)
	}
	return httpc.NewJSONRequest[discord.Message](f.client.http, func(b httpc.RequestBuilder) error {
		b.Method(fasthttp.MethodPatch)
		uploadFiles(b, params, params.Attachments)
		return b.Execute("webhooks", user.ID.String(), f.Token, "messages", f.ID.String())
	})
}

func (f FollowUpResolver) Delete() error {
	user, err := f.client.GetCurrentUser()
	if err != nil {
		return fmt.Errorf("failed to get current user: %w", err)
	}
	return httpc.NewRequest(f.client.http, func(b httpc.RequestBuilder) error {
		b.Method(fasthttp.MethodDelete)
		return b.Execute("webhooks", user.ID.String(), f.Token, "messages", f.ID.String())
	})
}
