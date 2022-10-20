package builders

import (
	"github.com/BOOMfinity/bfcord/discord"
	"github.com/BOOMfinity/bfcord/internal/timeconv"
	"github.com/BOOMfinity/go-utils/nullable"
	"github.com/andersfylling/snowflake/v5"
	"strings"
	"time"
)

func NewUpdateGuildMemberBuilder(guild, member snowflake.ID) *UpdateGuildMemberBuilder {
	return &UpdateGuildMemberBuilder{ID: member, Guild: guild}
}

type UpdateGuildMemberBuilder struct {
	Data  discord.MemberUpdate
	Guild snowflake.ID
	ID    snowflake.ID
}

func (u *UpdateGuildMemberBuilder) Nick(name string) discord.UpdateGuildMemberBuilder {
	u.Data.Nick = &name
	return u
}

func (u *UpdateGuildMemberBuilder) Roles(roles []snowflake.ID) discord.UpdateGuildMemberBuilder {
	u.Data.Roles = &roles
	return u
}

func (u *UpdateGuildMemberBuilder) Mute(isMuted bool) discord.UpdateGuildMemberBuilder {
	u.Data.Mute = &isMuted
	return u
}

func (u *UpdateGuildMemberBuilder) Deaf(isDeafened bool) discord.UpdateGuildMemberBuilder {
	u.Data.Deaf = &isDeafened
	return u
}

func (u *UpdateGuildMemberBuilder) VoiceChannel(channel snowflake.ID) discord.UpdateGuildMemberBuilder {
	if u.Data.ChannelID == nil {
		u.Data.ChannelID = &nullable.Nullable[snowflake.ID]{}
	}
	if channel.IsZero() {
		u.Data.ChannelID.Clear()
	} else {
		u.Data.ChannelID.Set(channel)
	}
	return u
}

func (u *UpdateGuildMemberBuilder) DisableCommunicationUntil(t time.Time) discord.UpdateGuildMemberBuilder {
	u.Data.CommunicationDisabledUntil = &timeconv.Timestamp{Time: t}
	return u
}

func (u *UpdateGuildMemberBuilder) Execute(api discord.ClientQuery, reason ...string) (member discord.MemberWithUser, err error) {
	return api.LowLevel().Reason(strings.Join(reason, " ")).UpdateGuildMember(u.Guild, u.ID, u.Data)
}
