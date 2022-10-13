package discord

import (
	"github.com/andersfylling/snowflake/v5"
)

type MessageReaction struct {
	Emoji     Emoji        `json:"emoji"`
	ChannelID snowflake.ID `json:"-"`
	MessageID snowflake.ID `json:"-"`
	Count     int          `json:"count"`
	Me        bool         `json:"me"`
}

func (v MessageReaction) API(api ClientQuery) MessageReactionQuery {
	return v.Message(api).Reaction(v.Emoji.ToString())
}

func (v MessageReaction) RemoveOwn(api ClientQuery) error {
	return v.API(api).RemoveOwn()
}

func (v MessageReaction) Remove(api ClientQuery, userID snowflake.ID) error {
	return v.API(api).Remove(userID)
}

func (v MessageReaction) Channel(api ClientQuery) ChannelQuery {
	return api.Channel(v.ChannelID)
}

func (v MessageReaction) Message(api ClientQuery) MessageQuery {
	return api.Channel(v.ChannelID).Message(v.MessageID)
}

func (v MessageReaction) Users(api ClientQuery, limit uint64) ([]User, error) {
	return v.API(api).All(limit)
}
