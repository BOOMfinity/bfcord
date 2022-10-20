package api

import (
	"fmt"
	"github.com/BOOMfinity/bfcord/api/builders"

	"github.com/BOOMfinity/bfcord/discord"
	"github.com/andersfylling/snowflake/v5"
	"github.com/segmentio/encoding/json"
	"github.com/valyala/fasthttp"
)

var _ = (discord.ChannelQuery)(&ChannelQuery{})

type ChannelQuery struct {
	client *Client
	emptyOptions[discord.ChannelQuery]
	id snowflake.ID
}

func (c ChannelQuery) Edit() discord.UpdateChannelTypeSelector {
	return builders.UpdateChannelTypeSelector{ID: c.id}
}

func (c ChannelQuery) Messages() discord.ChannelMessagesQuery {
	return NewChannelMessagesQuery(c.client, c.id)
}

func (c ChannelQuery) Get() (ch discord.Channel, err error) {
	req := c.client.New(true)
	req.SetRequestURI(fmt.Sprintf("%v/channels/%v", FullApiUrl, c.id))
	err = c.client.DoResult(req, &ch)
	return
}

func (c ChannelQuery) Invites() (invites []discord.InviteWithMeta, err error) {
	req := c.client.New(true)
	req.SetRequestURI(fmt.Sprintf("%v/channels/%v/invites", FullApiUrl, c.id))
	err = c.client.DoResult(req, &invites)
	return
}

func (c ChannelQuery) Delete() error {
	req := c.client.New(true)
	req.SetRequestURI(fmt.Sprintf("%v/channels/%v", FullApiUrl, c.id))
	req.Header.SetMethod(fasthttp.MethodDelete)
	if c.reason != "" {
		req.Header.Set("X-Audit-Log-Reason", c.reason)
	}
	return c.client.DoNoResp(req)
}

func (c ChannelQuery) Bulk(ids []snowflake.ID) error {
	req := c.client.New(true)
	req.SetRequestURI(fmt.Sprintf("%v/channels/%v/messages/bulk-delete", FullApiUrl, c.id))
	data := struct {
		Messages []snowflake.ID `json:"messages"`
	}{ids}
	raw, err := json.Marshal(data)
	if err != nil {
		return err
	}
	req.SetBody(raw)
	req.Header.SetMethod(fasthttp.MethodPost)
	if c.reason != "" {
		req.Header.Set("X-Audit-Log-Reason", c.reason)
	}
	return c.client.DoNoResp(req)
}

func (c ChannelQuery) Follow(target snowflake.ID) error {
	req := c.client.New(true)
	req.SetRequestURI(fmt.Sprintf("%v/channels/%v", FullApiUrl, c.id))
	data := struct {
		WebhookChannelID snowflake.ID `json:"webhook_channel_id"`
	}{target}
	raw, err := json.Marshal(data)
	if err != nil {
		return err
	}
	req.SetBody(raw)
	req.Header.SetMethod(fasthttp.MethodDelete)
	return c.client.DoNoResp(req)
}

func (c ChannelQuery) Pinned() (msg []discord.Message, err error) {
	req := c.client.New(true)
	req.SetRequestURI(fmt.Sprintf("%v/channels/%v/pins", FullApiUrl, c.id))
	err = c.client.DoResult(req, &msg)
	return
}

func (c ChannelQuery) StartThread(name string) discord.CreateThreadTypeSelector {
	return builders.CreateThreadTypeSelector{Channel: c.id, Data: discord.ThreadCreate{Name: &name}}
}

func (c ChannelQuery) StartForumThread(name string) discord.CreateForumMessageBuilder {
	return builders.NewCreateForumMessageBuilder(c.id, name)
}

func (c ChannelQuery) Join() error {
	req := c.client.New(true)
	req.SetRequestURI(fmt.Sprintf("%v/channels/%v/thread-members/@me", FullApiUrl, c.id))
	req.Header.SetMethod(fasthttp.MethodPut)
	return c.client.DoNoResp(req)
}

func (c ChannelQuery) AddMember(id snowflake.ID) error {
	req := c.client.New(true)
	req.SetRequestURI(fmt.Sprintf("%v/channels/%v/thread-members/%v", FullApiUrl, c.id, id))
	req.Header.SetMethod(fasthttp.MethodPut)
	return c.client.DoNoResp(req)
}

func (c ChannelQuery) Leave() error {
	req := c.client.New(true)
	req.SetRequestURI(fmt.Sprintf("%v/channels/%v/thread-members/@me", FullApiUrl, c.id))
	req.Header.SetMethod(fasthttp.MethodDelete)
	return c.client.DoNoResp(req)
}

func (c ChannelQuery) RemoveMember(id snowflake.ID) error {
	req := c.client.New(true)
	req.SetRequestURI(fmt.Sprintf("%v/channels/%v/thread-members/%v", FullApiUrl, c.id, id))
	req.Header.SetMethod(fasthttp.MethodDelete)
	return c.client.DoNoResp(req)
}

func (c ChannelQuery) GetThreadMember(id snowflake.ID) (tm discord.ThreadMember, err error) {
	req := c.client.New(true)
	req.SetRequestURI(fmt.Sprintf("%v/channels/%v/thread-members/%v", FullApiUrl, c.id, id))
	err = c.client.DoResult(req, &tm)
	return
}

func (c ChannelQuery) SendMessage() discord.CreateMessageBuilder {
	return builders.NewCreateMessageBuilder(c.id)
}

func (c ChannelQuery) Message(id snowflake.ID) discord.MessageQuery {
	return NewMessageQuery(c.client, c.id, id)
}

func (c ChannelQuery) Stage() discord.StageQuery {
	return NewStageQuery(c.client, c.id)
}

func (c ChannelQuery) ID() snowflake.ID {
	return c.id
}

func NewChannelQuery(client *Client, id snowflake.ID) *ChannelQuery {
	d := &ChannelQuery{
		id:     id,
		client: client,
	}
	d.emptyOptions = emptyOptions[discord.ChannelQuery]{data: d}
	return d
}

var _ = (discord.StageQuery)(&StageQuery{})

type StageQuery struct {
	client *Client
	emptyOptions[discord.StageQuery]
	id snowflake.ID
}

func (s StageQuery) Create(topic string, notify bool) (stage discord.StageInstance, err error) {
	req := s.client.New(true)
	req.Header.SetMethod(fasthttp.MethodPost)
	req.SetRequestURI(FullApiUrl + "/stage-instances")
	if topic != "" {
		json, _err := json.Marshal(map[string]any{
			"topic":                   topic,
			"channel_id":              s.id,
			"send_start_notification": notify,
		})
		if _err != nil {
			err = fmt.Errorf("failed to marshal json body: %w", _err)
			return
		}
		req.SetBody(json)
	}
	err = s.client.DoResult(req, &stage)
	return
}

func (s StageQuery) Get() (stage discord.StageInstance, err error) {
	req := s.client.New(true)
	req.SetRequestURI(FullApiUrl + "/stage-instances")
	err = s.client.DoResult(req, &stage)
	return
}

func (s StageQuery) Modify(topic string) (stage discord.StageInstance, err error) {
	req := s.client.New(true)
	req.Header.SetMethod(fasthttp.MethodPatch)
	req.SetRequestURI(fmt.Sprintf("%v/stage-instances/%v", FullApiUrl, s.id))
	json, _err := json.Marshal(map[string]any{
		"topic": topic,
	})
	if _err != nil {
		err = fmt.Errorf("failed to marshal json body: %w", _err)
		return
	}
	req.SetBody(json)
	err = s.client.DoResult(req, &stage)
	return
}

func (s StageQuery) Delete() error {
	req := s.client.New(true)
	req.Header.SetMethod(fasthttp.MethodDelete)
	req.SetRequestURI(fmt.Sprintf("%v/stage-instances/%v", FullApiUrl, s.id))
	return s.client.DoNoResp(req)
}

func NewStageQuery(client *Client, id snowflake.ID) *StageQuery {
	d := &StageQuery{
		id:     id,
		client: client,
	}
	d.emptyOptions = emptyOptions[discord.StageQuery]{data: d}
	return d
}
