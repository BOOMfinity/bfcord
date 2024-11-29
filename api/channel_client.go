package api

import (
	"fmt"
	"github.com/BOOMfinity/bfcord/discord"
	"github.com/BOOMfinity/bfcord/internal/httpc"
	"github.com/BOOMfinity/go-utils/inlineif"
	"github.com/andersfylling/snowflake/v5"
	"github.com/valyala/fasthttp"
)

type ChannelResolver struct {
	ID     snowflake.ID
	client *client
}

func (c ChannelResolver) JoinThread() error {
	return httpc.NewRequest(c.client.http, func(b httpc.RequestBuilder) error {
		b.Method(fasthttp.MethodPut)
		return b.Execute("channels", c.ID.String(), "thread-members", "@me")
	})
}

func (c ChannelResolver) AddThreadMember(id snowflake.ID) error {
	return httpc.NewRequest(c.client.http, func(b httpc.RequestBuilder) error {
		b.Method(fasthttp.MethodPut)
		return b.Execute("channels", c.ID.String(), "thread-members", id.String())
	})
}

func (c ChannelResolver) LeaveThread() error {
	return httpc.NewRequest(c.client.http, func(b httpc.RequestBuilder) error {
		b.Method(fasthttp.MethodDelete)
		return b.Execute("channels", c.ID.String(), "thread-members", "@me")
	})
}

func (c ChannelResolver) RemoveThreadMember(id snowflake.ID) error {
	return httpc.NewRequest(c.client.http, func(b httpc.RequestBuilder) error {
		b.Method(fasthttp.MethodDelete)
		return b.Execute("channels", c.ID.String(), "thread-members", id.String())
	})
}

func (c ChannelResolver) ThreadMember(id snowflake.ID, withMember bool) (discord.ThreadMember, error) {
	return httpc.NewJSONRequest[discord.ThreadMember](c.client.http, func(b httpc.RequestBuilder) error {
		return b.Execute("channels", c.ID.String(), "thread-members", id.String()+fmt.Sprintf("?with_member=%v", withMember))
	})
}

func (c ChannelResolver) AddRecipient(id snowflake.ID, userToken, userNickname string) error {
	return httpc.NewRequest(c.client.http, func(b httpc.RequestBuilder) error {
		b.Method(fasthttp.MethodPut)
		b.Body(map[string]any{
			"access_token": userToken,
			"nick":         userNickname,
		})
		return b.Execute("channels", c.ID.String(), "recipients", id.String())
	})
}

func (c ChannelResolver) RemoveRecipient(id snowflake.ID) error {
	return httpc.NewRequest(c.client.http, func(b httpc.RequestBuilder) error {
		b.Method(fasthttp.MethodDelete)
		return b.Execute("channels", c.ID.String(), "recipients", id.String())
	})
}

func (c ChannelResolver) Modify(params ModifyChannelParams, reason ...string) (discord.Channel, error) {
	return httpc.NewJSONRequest[discord.Channel](c.client.http, func(b httpc.RequestBuilder) error {
		b.Method(fasthttp.MethodPatch)
		if params.GroupDM != nil {
			b.Body(params.GroupDM)
		} else if params.Guild != nil {
			b.Body(params.Guild)
		} else if params.Thread != nil {
			b.Body(params.Thread)
		}
		b.Reason(reason...)
		return b.Execute("channels", c.ID.String())
	})
}

func (c ChannelResolver) StartThread(data StartThreadWithoutMessageParams, reason ...string) (discord.Channel, error) {
	return httpc.NewJSONRequest[discord.Channel](c.client.http, func(b httpc.RequestBuilder) error {
		b.Method(fasthttp.MethodPost)
		b.Body(data)
		b.Reason(reason...)
		return b.Execute("channels", c.ID.String(), "threads")
	})
}

func (c ChannelResolver) StartForumMediaThread(data StartForumOrMediaThreadParams, reason ...string) (discord.Channel, error) {
	return httpc.NewJSONRequest[discord.Channel](c.client.http, func(b httpc.RequestBuilder) error {
		b.Method(fasthttp.MethodPost)
		b.Body(data)
		b.Reason(reason...)
		return b.Execute("channels", c.ID.String(), "threads")
	})
}

func (c ChannelResolver) UpdateChannelPermissions(id snowflake.ID, data UpdateChannelPermissionsParams, reason ...string) error {
	return httpc.NewRequest(c.client.http, func(b httpc.RequestBuilder) error {
		b.Method(fasthttp.MethodPut)
		b.Reason(reason...)
		b.Body(data)
		return b.Execute("channels", c.ID.String(), "permissions", id.String())
	})
}

func (c ChannelResolver) DeleteChannelPermission(id snowflake.ID, reason ...string) error {
	return httpc.NewRequest(c.client.http, func(b httpc.RequestBuilder) error {
		b.Method(fasthttp.MethodDelete)
		b.Reason(reason...)
		return b.Execute("channels", c.ID.String(), "permissions", id.String())
	})
}

func (c ChannelResolver) FollowAnnouncementChannel(webhook snowflake.ID, reason ...string) (discord.FollowedChannel, error) {
	return httpc.NewJSONRequest[discord.FollowedChannel](c.client.http, func(b httpc.RequestBuilder) error {
		b.Method(fasthttp.MethodPost)
		b.Reason(reason...)
		b.Body(webhook)
		return b.Execute("channels", c.ID.String(), "followers")
	})
}

func (c ChannelResolver) Pins() ([]discord.Message, error) {
	return httpc.NewJSONRequest[[]discord.Message](c.client.http, func(b httpc.RequestBuilder) error {
		return b.Execute("channels", c.ID.String(), "pins")
	})
}

func (c ChannelResolver) Invites() ([]discord.Invite, error) {
	return httpc.NewJSONRequest[[]discord.Invite](c.client.http, func(b httpc.RequestBuilder) error {
		return b.Execute("channels", c.ID.String(), "invites")
	})
}

func (c ChannelResolver) CreateInvite(data CreateChannelInviteParams, reason ...string) (discord.Invite, error) {
	return httpc.NewJSONRequest[discord.Invite](c.client.http, func(b httpc.RequestBuilder) error {
		b.Method(fasthttp.MethodPost)
		b.Reason(reason...)
		b.Body(data)
		return b.Execute("channels", c.ID.String(), "invites")
	})
}

func (c ChannelResolver) Delete(reason ...string) error {
	return httpc.NewRequest(c.client.http, func(b httpc.RequestBuilder) error {
		b.Method(fasthttp.MethodDelete)
		b.Reason(reason...)
		return b.Execute("channels", c.ID.String())
	})
}

func (c ChannelResolver) Webhooks() ([]discord.Webhook, error) {
	//TODO implement me
	panic("implement me")
}

func (c ChannelResolver) CreateWebhook(params CreateWebhookParams, reason ...string) (discord.Webhook, error) {
	return httpc.NewJSONRequest[discord.Webhook](c.client.http, func(b httpc.RequestBuilder) error {
		b.Method(fasthttp.MethodPost)
		b.Body(params)
		b.Reason(reason...)
		return b.Execute("channels", c.ID.String(), "webhooks")
	})
}

func (c ChannelResolver) Get() (dst discord.Channel, err error) {
	defer c.client.proxy.AddChannel(dst)
	return httpc.NewJSONRequest[discord.Channel](c.client.http, func(b httpc.RequestBuilder) error {
		return b.Execute("channels", c.ID.String())
	})
}

func (c ChannelResolver) SendMessage(params CreateMessageParams) (discord.Message, error) {
	return httpc.NewJSONRequest[discord.Message](c.client.http, func(b httpc.RequestBuilder) error {
		b.Method(fasthttp.MethodPost)
		uploadFiles(b, params, params.Attachments)
		return b.Execute("channels", c.ID.String(), "messages")
	})
}

func (c ChannelResolver) Messages() MessagesQuery {
	return MessagesResolver{
		client:  c.client,
		Channel: c.ID,
	}
}

func (c ChannelResolver) Message(id snowflake.ID) MessageClient {
	return MessageResolver{
		client:  c.client,
		Message: id,
		Channel: c.ID,
	}
}

func (c ChannelResolver) BulkDelete(messages []snowflake.ID, reason ...string) error {
	for len(messages) > 0 {
		ids := messages[:inlineif.IfElse(len(messages) > 100, 100, len(messages))]
		messages = messages[len(ids):]
		if err := httpc.NewRequest(c.client.http, func(b httpc.RequestBuilder) error {
			b.Method(fasthttp.MethodPost)
			b.Body(map[string]any{
				"messages": ids,
			})
			b.Reason(reason...)
			return b.Execute("channels", c.ID.String(), "messages", "bulk-delete")
		}); err != nil {
			return err
		}
	}
	return nil
}
