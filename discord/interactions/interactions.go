package interactions

import (
	"fmt"
	"github.com/BOOMfinity/bfcord/api/builders"

	"github.com/BOOMfinity/bfcord/api"
	"github.com/BOOMfinity/bfcord/slash"
	"github.com/segmentio/encoding/json"
	"github.com/valyala/fasthttp"

	"github.com/BOOMfinity/bfcord/discord"
	"github.com/BOOMfinity/bfcord/discord/components"
	"github.com/andersfylling/snowflake/v5"
)

type ActionType uint8

const (
	PingAction ActionType = iota + 1
	CommandAction
	MessageComponentAction
	AutocompleteAction
	ModalSubmitAction
)

var (
	http = api.NewClient("", api.WithRetries(0))
)

type ResolvedData struct {
	Users map[snowflake.ID]discord.User `json:"users,omitempty"`
	// Always partial - missing: deaf and mute
	Members map[snowflake.ID]discord.Member `json:"members,omitempty"`
	Roles   map[snowflake.ID]discord.Role   `json:"roles,omitempty"`
	// Always partial - only id, name, type, permissions, thread_metadata and parent_id fields are included
	Channels map[snowflake.ID]discord.Channel     `json:"channels,omitempty"`
	Messages map[snowflake.ID]discord.BaseMessage `json:"messages,omitempty"`
}

type ActionData struct {
	Resolved      ResolvedData        `json:"resolved,omitempty"`
	Name          string              `json:"name"`
	TargetID      string              `json:"target_id,omitempty"`
	CustomID      string              `json:"custom_id,omitempty"`
	Options       OptionList          `json:"options,omitempty"`
	Values        []string            `json:"values,omitempty"`
	Components    ComponentList       `json:"components,omitempty"`
	ID            snowflake.Snowflake `json:"id"`
	Type          uint                `json:"type"`
	ComponentType components.Type     `json:"component_type,omitempty"`
}

type Interaction struct {
	Token         string                 `json:"token"`
	Locale        string                 `json:"locale,omitempty"`
	Member        discord.MemberWithUser `json:"member,omitempty"`
	Data          ActionData             `json:"data,omitempty"`
	User          discord.User           `json:"user,omitempty"`
	Message       discord.Message        `json:"message"`
	GuildID       snowflake.Snowflake    `json:"guild_id,omitempty"`
	ChannelID     snowflake.Snowflake    `json:"channel_id,omitempty"`
	ID            snowflake.Snowflake    `json:"id,omitempty"`
	ApplicationID snowflake.Snowflake    `json:"application_id,omitempty"`
	Type          ActionType             `json:"type,omitempty"`
}

func (i *Interaction) AutocompleteReply(choices []slash.Choice) error {
	custom := i.CustomReply()
	custom.Type(AutocompleteResultCallback)
	custom.Data(ResponseData{Choices: &choices})
	return custom.Execute()
}

func (i *Interaction) CustomReply() OriginalCustomBuilder {
	return &originalCustomBuilder{i: i}
}

func (i *Interaction) ComponentMessageUpdate() InteractionUpdateBuilder {
	bl := &interactionBuilder[InteractionUpdateBuilder]{
		t: UpdateMessageCallback,
		i: i,
	}
	bl.BaseMessageBuilder = &builders.MessageBuilder[InteractionUpdateBuilder]{B: bl}
	return bl
}

func (i *Interaction) FetchOriginalResponse() error {
	req := http.New(false)
	req.SetRequestURI(fmt.Sprintf(api.FullApiUrl+"/webhooks/%v/%v/messages/@original", i.ApplicationID, i.Token))
	defer i.Message.Patch()
	return http.DoResult(req, &i.Message)
}

func (i *Interaction) EditOriginalReply() OriginalUpdateBuilder {
	bl := &originalUpdateBuilder{
		i: i,
	}
	bl.BaseMessageBuilder = &builders.MessageBuilder[OriginalUpdateBuilder]{B: bl}
	return bl
}

func (i *Interaction) Delete() error {
	req := http.New(false)
	req.SetRequestURI(fmt.Sprintf(api.FullApiUrl+"/webhooks/%v/%v/messages/@original", i.ApplicationID, i.Token))
	req.Header.SetMethod(fasthttp.MethodDelete)
	return http.DoNoResp(req)
}

func (i *Interaction) SendMessageReply() ReplyBuilder {
	bl := &interactionBuilder[ReplyBuilder]{t: ChannelMessageWithSourceCallback, i: i}
	bl.BaseMessageBuilder = &builders.MessageBuilder[ReplyBuilder]{B: bl}
	return bl
}

// Defer loading false works only for component-based interactions!
func (i *Interaction) Defer(loading bool, ephemeral bool) error {
	var flag InteractionCallbackFlags
	if ephemeral {
		flag = EphemeralFlag
	}
	if !i.IsMessageComponent() {
		loading = true
	}
	custom := i.CustomReply()
	if loading {
		custom.Type(DeferredChannelMessageWithSourceCallback)
	} else {
		custom.Type(DeferredUpdateMessageCallback)
	}
	if flag != 0 {
		custom.Data(ResponseData{Flags: flag})
	}
	return custom.Execute()
}

func (i *Interaction) CreateFollowUp() FollowUpCreateBuilder {
	bl := followUpCreateBuilder{i: i}
	bl.ExpandableWebhookExecuteBuilder = &builders.WebhookExecuteBuilder[FollowUpCreateBuilder]{}
	return bl
}

func (i *Interaction) CreateModal(id string, title string) *ModalBuilder {
	return NewModalBuilder(i, id, title)
}

func (i *Interaction) IsCommand() bool {
	if i.Type == CommandAction {
		return true
	}
	return false
}

func (i *Interaction) IsMessageComponent() bool {
	if i.Type == MessageComponentAction {
		return true
	}
	return false
}

func (i *Interaction) IsModalSubmit() bool {
	if i.Type == ModalSubmitAction {
		return true
	}
	return false
}

func (i *Interaction) IsAutocomplete() bool {
	if i.Type == AutocompleteAction {
		return true
	}
	return false
}

func (i *Interaction) IsButton() bool {
	if i.Data.ComponentType == components.TypeButton {
		return true
	}
	return false
}

func (i *Interaction) IsSelectMenu() bool {
	if i.Data.ComponentType == components.TypeSelectMenu {
		return true
	}
	return false
}

// UserID might be present (or not) in multiple places - this method gets it wherever it exists.
func (i *Interaction) UserID() snowflake.Snowflake {
	if !i.User.ID.IsZero() {
		return i.User.ID
	}
	if !i.Member.UserID.IsZero() {
		return i.Member.UserID
	}
	return i.Member.User.ID
}

func (i *Interaction) InGuild() bool {
	if i.GuildID.Valid() {
		return true
	}
	return false
}

type ReplyBuilder ExpandableReplyBuilder[ReplyBuilder]

type ExpandableReplyBuilder[B any] interface {
	discord.BaseMessageBuilder[B]
	Ephemeral() B
	Execute() (err error)
}

type InteractionUpdateBuilder interface {
	discord.BaseMessageBuilder[InteractionUpdateBuilder]
	Execute() (err error)
}

type interactionBuilder[B any] struct {
	discord.BaseMessageBuilder[B]
	i    *Interaction
	data ResponseData
	t    CallbackType
}

func (v *interactionBuilder[B]) Ephemeral() B {
	v.data.Flags = EphemeralFlag
	return v.BaseMessageBuilder.Builder()
}

func (v *interactionBuilder[B]) Execute() (err error) {
	v.data.MessageCreate = v.Raw()
	custom := v.i.CustomReply()
	custom.Type(v.t)
	custom.Data(v.data)
	return custom.Execute()
}

type FollowUpCreateBuilder interface {
	discord.ExpandableWebhookExecuteBuilder[FollowUpCreateBuilder]
	Execute() (msg FollowUpMessage, err error)
}

type FollowUpUpdateBuilder interface {
	discord.BaseMessageBuilder[FollowUpUpdateBuilder]
	Execute() (msg FollowUpMessage, err error)
}

type OriginalUpdateBuilder interface {
	discord.BaseMessageBuilder[OriginalUpdateBuilder]
	Execute() (err error)
}

type originalUpdateBuilder struct {
	discord.BaseMessageBuilder[OriginalUpdateBuilder]
	i *Interaction
}

func (v *originalUpdateBuilder) Execute() (err error) {
	msg, err := http.LowLevel().Message(fasthttp.MethodPatch, fmt.Sprintf(api.FullApiUrl+"/webhooks/%v/%v/messages/@original", v.i.ApplicationID, v.i.Token), v.Raw())
	if err != nil {
		return
	}
	msg.Patch()
	v.i.Message = *msg
	return
}

type OriginalCustomBuilder interface {
	Type(t CallbackType)
	Data(x ResponseData)
	Execute() (err error)
}

type originalCustomBuilder struct {
	i *Interaction
	d *ResponseData
	t CallbackType
}

func (o *originalCustomBuilder) Type(t CallbackType) {
	o.t = t
}

func (o *originalCustomBuilder) Data(x ResponseData) {
	o.d = &x
}

func (o *originalCustomBuilder) Execute() (err error) {
	req := http.New(false)
	req.SetRequestURI(fmt.Sprintf(api.FullApiUrl+"/interactions/%v/%v/callback", o.i.ID.String(), o.i.Token))
	req.Header.SetMethod(fasthttp.MethodPost)
	raw, err := json.Marshal(InteractionResponse{Type: o.t, Data: o.d})
	if err != nil {
		return err
	}
	req.SetBody(raw)
	return http.DoNoResp(req)
}
