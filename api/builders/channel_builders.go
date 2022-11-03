package builders

import (
	"fmt"
	"github.com/BOOMfinity/bfcord/discord"
	"github.com/BOOMfinity/bfcord/discord/permissions"
	"github.com/andersfylling/snowflake/v5"
	"strings"
)

var _ = (discord.CreateThreadChannelBuilder)(&ChannelBuilder[discord.CreateThreadChannelBuilder]{})
var _ = (discord.GuildChannelBuilder)(&ChannelBuilder[discord.GuildChannelBuilder]{})
var _ = (discord.UpdateThreadChannelBuilder)(&ChannelBuilder[discord.UpdateThreadChannelBuilder]{})
var _ = (discord.UpdateGuildChannelBuilder)(&ChannelBuilder[discord.UpdateGuildChannelBuilder]{})

type CreateThreadTypeSelector struct {
	Data    discord.ThreadCreate
	Channel snowflake.ID
	Message snowflake.ID
}

func (c CreateThreadTypeSelector) Public() discord.CreateThreadChannelBuilder {
	bl := &CreateThreadBuilder[discord.CreateThreadChannelBuilder]{}
	public := discord.ChannelTypePublicThread
	bl.B = bl
	bl.Data = c.Data
	bl.Data.Type = &public
	bl.ID = c.Channel
	bl.Message = c.Message
	return bl
}

func (c CreateThreadTypeSelector) Private() discord.CreateThreadChannelBuilder {
	bl := &CreateThreadBuilder[discord.CreateThreadChannelBuilder]{}
	public := discord.ChannelTypePrivateThread
	bl.B = bl
	bl.Data = c.Data
	bl.Data.Type = &public
	bl.ID = c.Channel
	bl.Message = c.Message
	return bl
}

func (c CreateThreadTypeSelector) News() discord.CreateThreadChannelBuilder {
	bl := &CreateThreadBuilder[discord.CreateThreadChannelBuilder]{}
	public := discord.ChannelTypeNewsThread
	bl.B = bl
	bl.Data = c.Data
	bl.Data.Type = &public
	bl.ID = c.Channel
	bl.Message = c.Message
	return bl
}

type UpdateChannelTypeSelector struct {
	ID snowflake.ID
}

func (u UpdateChannelTypeSelector) Thread() discord.UpdateThreadChannelBuilder {
	return NewUpdateThreadChannelBuilder(u.ID)
}

func (u UpdateChannelTypeSelector) Guild() discord.UpdateGuildChannelBuilder {
	return NewUpdateChannelBuilder(u.ID)
}

type CreateThreadBuilder[B any] struct {
	ChannelBuilder[B]
	Data    discord.ThreadCreate
	ID      snowflake.ID
	Message snowflake.ID
}

func (c *CreateThreadBuilder[B]) Execute(api discord.ClientQuery, reason ...string) (ch *discord.Channel, err error) {
	return api.LowLevel().Reason(strings.Join(reason, " ")).StartThread(c.ID, c.Message, c.Data)
}

func NewCreateThreadChannelBuilder(channel, message snowflake.ID, name string) *CreateThreadBuilder[discord.CreateThreadChannelBuilder] {
	bl := &CreateThreadBuilder[discord.CreateThreadChannelBuilder]{}
	bl.B = bl
	bl.Name(name)
	bl.Message = message
	bl.ID = channel
	return bl
}

func NewUpdateThreadChannelBuilder(channel snowflake.ID) *CreateThreadBuilder[discord.UpdateThreadChannelBuilder] {
	bl := &CreateThreadBuilder[discord.UpdateThreadChannelBuilder]{}
	bl.B = bl
	bl.ID = channel
	return bl
}

func NewUpdateChannelBuilder(channel snowflake.ID) *ChannelBuilder[discord.UpdateGuildChannelBuilder] {
	bl := &ChannelBuilder[discord.UpdateGuildChannelBuilder]{}
	bl.B = bl
	bl.ID = channel
	return bl
}

func NewCreateChannelBuilder(guild snowflake.ID, name string) *ChannelBuilder[discord.GuildChannelBuilder] {
	bl := &ChannelBuilder[discord.GuildChannelBuilder]{}
	bl.Name(name)
	bl.B = bl
	bl.Guild = guild
	return bl
}

type ChannelBuilder[B any] struct {
	Data  discord.ChannelUpdate
	B     B
	ID    snowflake.ID
	Guild snowflake.ID
}

func (c *ChannelBuilder[B]) Builder() B {
	return c.B
}

func (c *ChannelBuilder[B]) Name(name string) B {
	c.Data.Name = &name
	return c.B
}

func (c *ChannelBuilder[B]) Archived(archived bool) B {
	c.Data.Archived = &archived
	return c.B
}

func (c *ChannelBuilder[B]) AutoArchiveDuration(dur discord.ThreadArchiveDuration) B {
	c.Data.AutoArchiveDuration = &dur
	return c.B
}

func (c *ChannelBuilder[B]) Locked(locked bool) B {
	c.Data.Locked = &locked
	return c.B
}

func (c *ChannelBuilder[B]) Invitable(invitable bool) B {
	c.Data.Invitable = &invitable
	return c.B
}

func (c *ChannelBuilder[B]) Execute(api discord.ClientQuery, reason ...string) (ch *discord.Channel, err error) {
	ll := api.LowLevel()
	if len(reason) > 0 {
		ll = ll.Reason(strings.Join(reason, " "))
	}
	if c.ID.IsZero() {
		// create
		ch, err = ll.CreateGuildChannel(c.Guild, c.Data)
		if err != nil {
			return
		}
	} else {
		// edit
		ch, err = ll.UpdateChannel(c.Guild, c.Data)
		if err != nil {
			return
		}
	}
	if c.Data.Position != nil {
		channel, _err := api.Channel(c.ID).Get()
		if _err != nil {
			err = fmt.Errorf("could not fetch channel: %w", _err)
			return
		}
		channels, _err := api.Guild(channel.GuildID).NoCache().Channels()
		if _err != nil {
			err = fmt.Errorf("could not fetch guild channels: %w", _err)
			return
		}
		bl := discord.NewGuildChannelPositionsBuilder(channels)
		bl.Pos(c.ID, *c.Data.Position)
		err = api.Guild(channel.GuildID).UpdateChannelPositions(bl)
	}
	return
}

func (c *ChannelBuilder[B]) Type(t discord.ChannelType) B {
	c.Data.Type = &t
	return c.B
}

func (c *ChannelBuilder[B]) Topic(topic string) B {
	c.Data.Topic = &topic
	return c.B
}

func (c *ChannelBuilder[B]) Bitrate(bitrate uint64) B {
	c.Data.Bitrate = &bitrate
	return c.B
}

func (c *ChannelBuilder[B]) UserLimit(limit uint16) B {
	c.Data.UserLimit = &limit
	return c.B
}

func (c *ChannelBuilder[B]) RateLimitPerUser(limit uint32) B {
	c.Data.RateLimitPerUser = &limit
	return c.B
}

func (c *ChannelBuilder[B]) Position(pos int) B {
	c.Data.Position = &pos
	return c.B
}

func (c *ChannelBuilder[B]) Parent(id snowflake.ID) B {
	c.Data.ParentID = &id
	return c.B
}

func (c *ChannelBuilder[B]) NSFW(isNSFW bool) B {
	c.Data.Nsfw = &isNSFW
	return c.B
}

func (c *ChannelBuilder[B]) Overwrites(perms []permissions.Overwrite) B {
	c.Data.Overwrites = &perms
	return c.B
}
