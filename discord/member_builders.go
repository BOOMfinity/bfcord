package discord

import (
	"github.com/andersfylling/snowflake/v5"
	"time"
)

type UpdateGuildMemberBuilder interface {
	Nick(name string) UpdateGuildMemberBuilder
	Roles(roles []snowflake.ID) UpdateGuildMemberBuilder
	Mute(isMuted bool) UpdateGuildMemberBuilder
	Deaf(isDeafened bool) UpdateGuildMemberBuilder
	VoiceChannel(channel snowflake.ID) UpdateGuildMemberBuilder
	DisableCommunicationUntil(t time.Time) UpdateGuildMemberBuilder
	Execute(api ClientQuery, reason ...string) (member *MemberWithUser, err error)
}
