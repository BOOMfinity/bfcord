package interactions

import (
	"github.com/BOOMfinity/bfcord/api/builders"
	"github.com/BOOMfinity/bfcord/discord"
	"github.com/andersfylling/snowflake/v5"
)

type FollowUpMessage struct {
	*discord.Message
	i *Interaction
}

func (x FollowUpMessage) Edit() FollowUpUpdateBuilder {
	bl := followUpUpdateBuilder{id: x.ID, i: x.i}
	sbl := &builders.WebhookUpdateMessageBuilder[FollowUpUpdateBuilder]{}
	sbl.B = bl
	bl.ExpandableWebhookUpdateMessageBuilder = sbl
	return bl
}

func (x FollowUpMessage) Delete() error {
	return http.Webhook(x.i.ApplicationID, x.i.Token).DeleteMessage(x.ID)
}

type followUpCreateBuilder struct {
	discord.ExpandableWebhookExecuteBuilder[FollowUpCreateBuilder]

	i *Interaction
}

func (v followUpCreateBuilder) Execute() (msg *FollowUpMessage, err error) {
	m, err := http.LowLevel().ExecuteWebhook(v.i.ApplicationID, v.i.Token, discord.WebhookExecute{MessageCreate: v.Raw()}, true, 0)
	if err != nil {
		return
	}
	return &FollowUpMessage{i: v.i, Message: m}, nil
}

type followUpUpdateBuilder struct {
	discord.ExpandableWebhookUpdateMessageBuilder[FollowUpUpdateBuilder]

	i  *Interaction
	id snowflake.ID
}

func (v followUpUpdateBuilder) Execute() (msg *FollowUpMessage, err error) {
	m, err := http.LowLevel().UpdateWebhookMessage(v.i.ApplicationID, v.i.Token, v.id, v.Raw(), 0)
	if err != nil {
		return
	}
	return &FollowUpMessage{i: v.i, Message: m}, nil
}
