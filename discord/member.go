package discord

import (
	"github.com/BOOMfinity/bfcord/discord/permissions"
	"github.com/BOOMfinity/bfcord/internal/slices"
	"github.com/BOOMfinity/bfcord/internal/timeconv"
	"github.com/andersfylling/snowflake/v5"
)

type MemberWithUser struct {
	Member
	User User `json:"user"`
}

// Member
//
// Reference: https://discord.com/developers/docs/resources/guild#guild-member-object
type Member struct {
	JoinedAt                   timeconv.Timestamp `json:"joined_at"`
	PremiumSince               timeconv.Timestamp `json:"premium_since"`
	Nick                       string             `json:"nick"`
	Avatar                     string             `json:"avatar"`
	CommunicationDisabledUntil string             `json:"communication_disabled_until"`
	Roles                      []snowflake.ID     `json:"roles"`
	UserID                     snowflake.ID       `json:"user_id"`
	GuildID                    snowflake.ID       `json:"guild_id"`
	Deaf                       bool               `json:"deaf"`
	Mute                       bool               `json:"mute"`
	Pending                    bool               `json:"pending"`
}

func (x Member) API(client ClientQuery) GuildMemberQuery {
	return x.Guild(client).Member(x.UserID)
}

func (x Member) User(client ClientQuery) UserQuery {
	return client.User(x.UserID)
}

func (x Member) Guild(client ClientQuery) GuildQuery {
	return client.Guild(x.GuildID)
}

func (x Member) Kick(client ClientQuery) error {
	return x.Guild(client).Member(x.UserID).Kick()
}

func (x Member) Ban(client ClientQuery, days uint8) error {
	return x.Guild(client).Member(x.UserID).Ban(days)
}

func (x Member) PermissionsIn(bot ClientQuery, channel snowflake.ID) (resperm permissions.Permission, err error) {
	return bot.Guild(x.GuildID).Member(x.UserID).PermissionsIn(channel)
}

func (x Member) Permissions(bot ClientQuery) (perm permissions.Permission, err error) {
	return bot.Guild(x.GuildID).Member(x.UserID).Permissions()
}

func BasePermissions(member Member, guild BaseGuild) (perm permissions.Permission) {
	if guild.OwnerID == member.UserID {
		return permissions.All
	}
	for i := range guild.Roles {
		for i2 := range member.Roles {
			if member.Roles[i2] == guild.Roles[i].ID {
				perm.Add(guild.Roles[i].Permissions)
			}
		}
	}
	if perm.Administrator() {
		return permissions.All
	}
	return
}

func ChannelPermissions(guild BaseGuild, member Member, overwrites []permissions.Overwrite) (perm permissions.Permission) {
	base := BasePermissions(member, guild)
	if base.Administrator() {
		return permissions.All
	}
	data, ok := slices.FindCopy(overwrites, func(item permissions.Overwrite) bool {
		return item.ID == guild.ID
	})
	if ok {
		perm.Add(data.Allow)
		perm.Clear(data.Deny)
	}
	for i := range member.Roles {
		data, ok = slices.FindCopy(overwrites, func(item permissions.Overwrite) bool {
			return item.ID == member.Roles[i]
		})
		if !ok {
			continue
		}
		for i2 := range guild.Roles {
			if guild.Roles[i2].ID == member.Roles[i] {
				perm.Add(data.Allow)
				perm.Clear(data.Deny)
			}
		}
	}
	data, ok = slices.FindCopy(overwrites, func(item permissions.Overwrite) bool {
		return item.ID == member.UserID
	})
	if ok {
		perm.Add(data.Allow)
		perm.Clear(data.Deny)
	}
	return
}
