package api

import (
	"github.com/BOOMfinity/bfcord/discord"
	"github.com/BOOMfinity/bfcord/internal/httpc"
	"github.com/BOOMfinity/golog/v2"
	"github.com/andersfylling/snowflake/v5"
	"github.com/valyala/fasthttp"
)

type Client interface {
	Guild(id snowflake.ID) GuildClient
	Channel(id snowflake.ID) ChannelClient
	Slash() SlashClient
	User(id snowflake.ID) UserClient
	Stage(id snowflake.ID) StageClient
	ModifyCurrentUser(params ModifyCurrentUserParams) (discord.User, error)
	GatewayInfo() (BotGateway, error)
	GetCurrentUser() (discord.User, error)
	Interaction(id snowflake.ID, token string) InteractionClient
}

type client struct {
	user  discord.User
	http  *httpc.Client
	proxy CacheProxy
}

func (c *client) Interaction(id snowflake.ID, token string) InteractionClient {
	return InteractionResolver{
		client: c,
		ID:     id,
		Token:  token,
	}
}

func (c *client) Guild(id snowflake.ID) GuildClient {
	return GuildResolver{
		ID:     id,
		client: c,
	}
}

func (c *client) Channel(id snowflake.ID) ChannelClient {
	return ChannelResolver{
		ID:     id,
		client: c,
	}
}

func (c *client) Slash() SlashClient {
	user, err := c.GetCurrentUser()
	if err != nil {
		panic("failed to get current user")
	}
	return SlashResolver{
		client:        c,
		CurrentUserID: user.ID,
	}
}

func (c *client) Stage(id snowflake.ID) StageClient {
	return StageResolver{
		ID:     id,
		client: c,
	}
}

func (c *client) User(id snowflake.ID) UserClient {
	return UserResolver{
		ID:     id,
		client: c,
	}
}

func (c *client) ModifyCurrentUser(params ModifyCurrentUserParams) (discord.User, error) {
	return httpc.NewJSONRequest[discord.User](c.http, func(b httpc.RequestBuilder) error {
		b.Method(fasthttp.MethodPatch)
		b.Body(params)
		return b.Execute("users", "@me")
	})
}

func (c *client) GatewayInfo() (BotGateway, error) {
	return httpc.NewJSONRequest[BotGateway](c.http, func(b httpc.RequestBuilder) error {
		return b.Execute("gateway", "bot")
	})
}

func (c *client) GetCurrentUser() (u discord.User, _ error) {
	if c.user.ID.Valid() {
		return c.user, nil
	}
	defer func() {
		c.user = u
	}()
	return httpc.NewJSONRequest[discord.User](c.http, func(b httpc.RequestBuilder) error {
		return b.Execute("users", "@me")
	})
}

func NewClient(log golog.Logger, token string, opts ...ClientOption) Client {
	c := &client{
		proxy: noopProxy{},
		http:  httpc.NewClient(token, log),
	}
	for _, opt := range opts {
		opt(c)
	}
	return c
}

type ClientOption func(c *client)

func WithCacheProxy(proxy CacheProxy) ClientOption {
	return func(c *client) {
		c.proxy = proxy
	}
}

func WithUserID(id snowflake.ID) ClientOption {
	return func(c *client) {
		c.user.ID = id
	}
}
