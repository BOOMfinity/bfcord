package builders

import (
	"github.com/BOOMfinity/bfcord/discord"
	"github.com/andersfylling/snowflake/v5"
)

var _ = (discord.CreateForumMessageBuilder)(&ForumMessageBuilder[discord.CreateForumMessageBuilder]{})

type ForumMessageBuilder[B any] struct {
	*MessageBuilder[B]
	forum discord.ForumMessageCreate
}

func (f *ForumMessageBuilder[B]) Execute(api discord.ClientQuery) (discord.ChannelWithMessage, error) {
	f.forum.Message = f.Create.MessageCreate
	return api.LowLevel().CreateForumMessage(f.ChannelID, f.forum)
}

func (f *ForumMessageBuilder[B]) RawForum() discord.ForumMessageCreate {
	f.forum.Message = f.Create.MessageCreate
	return f.forum
}

func (f *ForumMessageBuilder[B]) Builder() B {
	return f.B
}

func (f *ForumMessageBuilder[B]) AutoArchiveDuration(dur discord.ThreadArchiveDuration) B {
	f.forum.AutoArchiveDuration = &dur
	return f.B
}

func (f *ForumMessageBuilder[B]) RateLimitPerUser(limit uint32) B {
	f.forum.RateLimitPerUser = limit
	return f.B
}

func (f *ForumMessageBuilder[B]) Tags(t []snowflake.ID) B {
	f.forum.AppliedTags = t
	return f.B
}

func NewCreateForumMessageBuilder(id snowflake.ID, name string) *ForumMessageBuilder[discord.CreateForumMessageBuilder] {
	bl := &ForumMessageBuilder[discord.CreateForumMessageBuilder]{}
	bl.MessageBuilder = &MessageBuilder[discord.CreateForumMessageBuilder]{}
	bl.B = bl
	bl.forum.Name = &name
	bl.ChannelID = id
	return bl
}
