package api

import (
	"fmt"
	"github.com/BOOMfinity/bfcord/api/builders"
	"net/url"

	"github.com/BOOMfinity/bfcord/discord"
	"github.com/andersfylling/snowflake/v5"
	"github.com/valyala/fasthttp"
)

type MessageQuery struct {
	client *Client
	emptyOptions[discord.MessageQuery]
	channel snowflake.ID
	id      snowflake.ID
}

func (m MessageQuery) CrossPost() error {
	req := m.client.New(true)
	req.SetRequestURI(fmt.Sprintf(FullApiUrl+"/channels/%v/messages/%v/crosspost", m.channel, m.id))
	req.Header.SetMethod(fasthttp.MethodPost)
	return m.client.DoNoResp(req)
}

func (m MessageQuery) Pin() error {
	req := m.client.New(true)
	req.SetRequestURI(fmt.Sprintf(FullApiUrl+"/channels/%v/pins/%v", m.channel, m.id))
	req.Header.SetMethod(fasthttp.MethodPut)
	if m.reason != "" {
		req.Header.Set("X-Audit-Log-Reason", m.reason)
	}
	return m.client.DoNoResp(req)
}

func (m MessageQuery) UnPin() error {
	req := m.client.New(true)
	req.SetRequestURI(fmt.Sprintf(FullApiUrl+"/channels/%v/pins/%v", m.channel, m.id))
	req.Header.SetMethod(fasthttp.MethodDelete)
	if m.reason != "" {
		req.Header.Set("X-Audit-Log-Reason", m.reason)
	}
	return m.client.DoNoResp(req)
}

func (m MessageQuery) StartThread(name string) discord.CreateThreadChannelBuilder {
	return builders.NewCreateThreadChannelBuilder(m.channel, m.id, name)
}

func (m MessageQuery) Get() (msg *discord.Message, err error) {
	req := m.client.New(true)
	req.SetRequestURI(fmt.Sprintf(FullApiUrl+"/channels/%v/messages/%v", m.channel, m.id))
	err = m.client.DoResult(req, &msg)
	if err != nil {
		return
	}
	msg.Patch()
	return
}

func (m MessageQuery) Delete() (err error) {
	req := m.client.New(true)
	req.SetRequestURI(fmt.Sprintf(FullApiUrl+"/channels/%v/messages/%v", m.channel, m.id))
	req.Header.SetMethod(fasthttp.MethodDelete)
	if m.reason != "" {
		req.Header.Set("X-Audit-Log-Reason", m.reason)
	}
	return m.client.DoNoResp(req)
}

func (m MessageQuery) React(emoji string) error {
	req := m.client.New(true)
	req.SetRequestURI(fmt.Sprintf(FullApiUrl+"/channels/%v/messages/%v/reactions/%v/@me", m.channel, m.id, url.QueryEscape(emoji)))
	req.Header.SetMethod(fasthttp.MethodPut)
	return m.client.DoNoResp(req)
}

func (m MessageQuery) Reaction(emoji string) discord.MessageReactionQuery {
	return NewReactionQuery(m.client, m.channel, m.id, url.QueryEscape(emoji))
}

func (m MessageQuery) RemoveAllReactions() (err error) {
	req := m.client.New(true)
	req.SetRequestURI(fmt.Sprintf(FullApiUrl+"/channels/%v/messages/%v/reactions", m.channel, m.id))
	req.Header.SetMethod(fasthttp.MethodDelete)
	return m.client.DoNoResp(req)
}

func (m MessageQuery) Edit() discord.MessageBuilder {
	return builders.NewUpdateMessageBuilder(m.channel, m.id)
}

func (m MessageQuery) ID() snowflake.ID {
	return m.id
}

func (m MessageQuery) ChannelID() snowflake.ID {
	return m.channel
}

func NewMessageQuery(client *Client, channel, id snowflake.ID) *MessageQuery {
	data := &MessageQuery{
		id:      id,
		channel: channel,
		client:  client,
	}
	data.emptyOptions = emptyOptions[discord.MessageQuery]{data: data}
	return data
}
