package api

import (
	"fmt"
	"github.com/BOOMfinity/bfcord/api/builders"
	"github.com/valyala/fasthttp"

	"github.com/BOOMfinity/bfcord/discord"
	"github.com/andersfylling/snowflake/v5"
)

type UserQuery struct {
	client *Client
	emptyOptions[discord.UserQuery]
	id snowflake.ID
}

func (u UserQuery) Get() (user *discord.User, err error) {
	req := u.client.New(true)
	req.SetRequestURI(fmt.Sprintf("%v/users/%v", FullApiUrl, u.id))
	err = u.client.DoResult(req, &user)
	return
}

func (u UserQuery) Send() (bld discord.CreateMessageBuilder, err error) {
	channel, err := u.CreateDM()
	if err != nil {
		return nil, err
	}
	return builders.NewCreateMessageBuilder(channel.ID), nil
}

func (u UserQuery) ID() snowflake.ID {
	return u.id
}

func (u UserQuery) CreateDM() (channel *discord.Channel, err error) {
	req := u.client.New(true)
	req.SetRequestURI(fmt.Sprintf("%v/users/@me/channels", FullApiUrl))
	req.SetBodyString(fmt.Sprintf(`{"recipient_id": %v}`, u.id))
	req.Header.SetMethod(fasthttp.MethodPost)
	err = u.client.DoResult(req, &channel)
	return
}

func NewUserQuery(client *Client, id snowflake.ID) *UserQuery {
	data := &UserQuery{
		id:     id,
		client: client,
	}
	data.emptyOptions = emptyOptions[discord.UserQuery]{data: data}
	return data
}
