package api

import (
	"github.com/BOOMfinity/bfcord/discord"
	"github.com/BOOMfinity/bfcord/internal/httpc"
	"github.com/andersfylling/snowflake/v5"
	"github.com/valyala/fasthttp"
)

type StageResolver struct {
	client *client
	ID     snowflake.ID
}

func (s StageResolver) Get() (discord.StageInstance, error) {
	return httpc.NewJSONRequest[discord.StageInstance](s.client.http, func(b httpc.RequestBuilder) error {
		return b.Execute("stage-instances", s.ID.String())
	})
}

func (s StageResolver) Create(params CreateStageInstanceParams, reason ...string) (discord.StageInstance, error) {
	return httpc.NewJSONRequest[discord.StageInstance](s.client.http, func(b httpc.RequestBuilder) error {
		b.Method(fasthttp.MethodPost)
		b.Body(params)
		b.Reason(reason...)
		return b.Execute("stage-instances")
	})
}

func (s StageResolver) Modify(params ModifyStageInstanceParams, reason ...string) (discord.StageInstance, error) {
	return httpc.NewJSONRequest[discord.StageInstance](s.client.http, func(b httpc.RequestBuilder) error {
		b.Method(fasthttp.MethodPatch)
		b.Body(params)
		b.Reason(reason...)
		return b.Execute("stage-instances", s.ID.String())
	})
}

func (s StageResolver) Delete(reason ...string) error {
	return httpc.NewRequest(s.client.http, func(b httpc.RequestBuilder) error {
		b.Method(fasthttp.MethodDelete)
		b.Reason(reason...)
		return b.Execute("stage-instances", s.ID.String())
	})
}
