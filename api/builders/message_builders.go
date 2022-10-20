package builders

import (
	"fmt"
	"github.com/BOOMfinity/bfcord/discord"
	"github.com/BOOMfinity/bfcord/discord/components"
	"github.com/andersfylling/snowflake/v5"
)

var _ = (discord.CreateMessageBuilder)(&MessageBuilder[discord.CreateMessageBuilder]{})
var _ = (discord.WebhookExecuteBuilder)(&WebhookExecuteBuilder[discord.WebhookExecuteBuilder]{})
var _ = (discord.WebhookUpdateMessageBuilder)(&WebhookUpdateMessageBuilder[discord.WebhookUpdateMessageBuilder]{})

func NewWebhookUpdateMessageBuilder(id snowflake.ID, token string, message snowflake.ID, api discord.ClientQuery) *WebhookUpdateMessageBuilder[discord.WebhookUpdateMessageBuilder] {
	bl := &WebhookUpdateMessageBuilder[discord.WebhookUpdateMessageBuilder]{}
	bl.ID = id
	bl.B = bl
	bl.Token = token
	bl.Message = message
	bl.api = api
	return bl
}

type WebhookUpdateMessageBuilder[B any] struct {
	MessageBuilder[B]
	api      discord.ClientQuery
	Token    string
	ID       snowflake.ID
	Message  snowflake.ID
	ThreadID snowflake.ID
}

func (b *WebhookUpdateMessageBuilder[B]) Thread(id snowflake.ID) B {
	b.ThreadID = id
	return b.B
}

func (b *WebhookUpdateMessageBuilder[B]) Execute() (msg discord.Message, err error) {
	return b.api.LowLevel().UpdateWebhookMessage(b.ID, b.Token, b.Message, b.Create.MessageCreate, b.ThreadID)
}

func NewWebhookExecuteBuilder(id snowflake.ID, token string, api discord.ClientQuery) *WebhookExecuteBuilder[discord.WebhookExecuteBuilder] {
	bl := &WebhookExecuteBuilder[discord.WebhookExecuteBuilder]{}
	bl.B = bl
	bl.ID = id
	bl.Token = token
	bl.api = api
	return bl
}

type WebhookExecuteBuilder[B any] struct {
	MessageBuilder[B]
	api      discord.ClientQuery
	Token    string
	ID       snowflake.ID
	ThreadID snowflake.ID
	Nowait   bool
}

func (b *WebhookExecuteBuilder[B]) Thread(id snowflake.ID) B {
	b.ThreadID = id
	return b.B
}

func (b *WebhookExecuteBuilder[B]) AvatarURL(url string) B {
	b.Create.AvatarURL = &url
	return b.B
}

func (b *WebhookExecuteBuilder[B]) Username(name string) B {
	b.Create.Username = &name
	return b.B
}

func (b *WebhookExecuteBuilder[B]) NoWait() B {
	b.Nowait = true
	return b.B
}

func (b *WebhookExecuteBuilder[B]) Execute() (msg discord.Message, err error) {
	return b.api.LowLevel().ExecuteWebhook(b.ID, b.Token, b.Create, !b.Nowait, b.ThreadID)
}

func NewUpdateMessageBuilder(channel, message snowflake.ID) *MessageBuilder[discord.MessageBuilder] {
	bl := &MessageBuilder[discord.MessageBuilder]{}
	bl.B = bl
	bl.MessageID = message
	bl.ChannelID = channel
	return bl
}

func NewCreateMessageBuilder(channel snowflake.ID) *MessageBuilder[discord.CreateMessageBuilder] {
	bl := &MessageBuilder[discord.CreateMessageBuilder]{}
	bl.B = bl
	bl.ChannelID = channel
	return bl
}

type MessageBuilder[B any] struct {
	Create              discord.WebhookExecute
	B                   B
	ChannelID           snowflake.ID
	MessageID           snowflake.ID
	attachmentsDisabled bool
}

func (m *MessageBuilder[B]) ActionRow(items ...components.ActionRowItem) B {
	if m.Create.Components == nil {
		m.Create.Components = new([]components.Component)
	}
	row := components.NewActionRow()
	for i := range items {
		row.Add(items[i].ToComponent())
	}
	*m.Create.Components = append(*m.Create.Components, row.ToComponent())
	return m.B
}

func (m *MessageBuilder[B]) AutoActionRows(items ...components.ActionRowItem) B {
	m.Create.Components = new([]components.Component)
	row := components.NewActionRow()
	for i := range items {
		item := items[i]
		if row.Size() == 5 || (item.Type() == components.TypeSelectMenu && row.Size() > 0) {
			*m.Create.Components = append(*m.Create.Components, row.ToComponent())
			row = components.NewActionRow()
		}
		row.Add(item.ToComponent())
	}
	if row.Size() > 0 {
		*m.Create.Components = append(*m.Create.Components, row.ToComponent())
		row = nil
	}
	return m.B
}

func (m *MessageBuilder[B]) Builder() B {
	return m.B
}

func (m *MessageBuilder[B]) DoNotKeepFiles() B {
	m.attachmentsDisabled = true
	return m.B
}

func (m *MessageBuilder[B]) KeepFiles(files []discord.Attachment) B {
	m.Create.Attachments = &files
	m.attachmentsDisabled = true
	return m.B
}

func (m *MessageBuilder[B]) Content(str string) B {
	m.Create.Content = &str
	return m.B
}

func (m *MessageBuilder[B]) Embed(embed discord.MessageEmbed) B {
	if m.Create.Embeds == nil {
		m.Create.Embeds = new([]discord.MessageEmbed)
	}
	*m.Create.Embeds = append(*m.Create.Embeds, embed)
	return m.B
}

func (m *MessageBuilder[B]) Embeds(embeds []discord.MessageEmbed) B {
	if m.Create.Embeds == nil {
		m.Create.Embeds = new([]discord.MessageEmbed)
	}
	*m.Create.Embeds = embeds
	return m.B
}

func (m *MessageBuilder[B]) Components(list []components.Component) B {
	m.Create.Components = &list
	return m.B
}

func (m *MessageBuilder[B]) File(f discord.MessageFile) B {
	if m.Create.Files == nil {
		m.Create.Files = new([]discord.MessageFile)
	}
	*m.Create.Files = append(*m.Create.Files, f)
	return m.B
}

func (m *MessageBuilder[B]) Files(f []discord.MessageFile) B {
	if m.Create.Files == nil {
		m.Create.Files = new([]discord.MessageFile)
	}
	*m.Create.Files = f
	return m.B
}

func (m *MessageBuilder[B]) ClearEmbeds() B {
	m.Create.Embeds = new([]discord.MessageEmbed)
	*m.Create.Embeds = []discord.MessageEmbed{}
	return m.B
}

func (m *MessageBuilder[B]) ClearFiles() B {
	m.Create.Files = new([]discord.MessageFile)
	*m.Create.Files = []discord.MessageFile{}
	return m.B
}

func (m *MessageBuilder[B]) ClearContent() B {
	m.Create.Content = new(string)
	*m.Create.Content = ""
	return m.B
}

func (m *MessageBuilder[B]) ClearComponents() B {
	m.Create.Components = new([]components.Component)
	*m.Create.Components = []components.Component{}
	return m.B
}

func (m *MessageBuilder[B]) Raw() discord.MessageCreate {
	return m.Create.MessageCreate
}

func (m *MessageBuilder[B]) Reference(ref discord.MessageReference) B {
	m.Create.MessageReference = &ref
	return m.B
}

func (m *MessageBuilder[B]) TTS() B {
	var x = true
	m.Create.TTS = &x
	return m.B
}

func (m *MessageBuilder[B]) SuppressEmbeds() B {
	// TODO: flags
	panic("implement me")
}

func (m *MessageBuilder[B]) Execute(api discord.ClientQuery) (discord.Message, error) {
	if m.MessageID.Valid() {
		if !m.attachmentsDisabled {
			msg, err := api.Channel(m.ChannelID).Message(m.MessageID).Get()
			if err != nil {
				return discord.Message{}, fmt.Errorf("couldn't fetch attachments: %w", err)
			}
			m.KeepFiles(msg.Attachments)
		} else {
			m.KeepFiles([]discord.Attachment{})
		}
		return api.LowLevel().UpdateMessage(m.ChannelID, m.MessageID, m.Create.MessageCreate)
	} else {
		return api.LowLevel().CreateMessage(m.ChannelID, m.Create.MessageCreate)
	}
}
