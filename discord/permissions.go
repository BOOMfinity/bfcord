package discord

import (
	"bytes"
	"fmt"
	"strconv"

	"github.com/BOOMfinity/go-utils/ubytes"
	"github.com/andersfylling/snowflake/v5"
)

type Permission uint64

func (p *Permission) UnmarshalJSON(b []byte) error {
	if bytes.Equal(b, []byte("null")) {
		return nil
	}
	if val, err := strconv.ParseUint(ubytes.ToString(b[1:len(b)-1]), 10, 64); err != nil {
		return fmt.Errorf("failed to parse permission: %w", err)
	} else {
		*p = Permission(val)
		return nil
	}
}

func (p Permission) MarshalJSON() ([]byte, error) {
	return ubytes.ToBytes(strconv.FormatUint(uint64(p), 10)), nil
}

func (p Permission) Admin() bool {
	return p&PermissionAdministrator == PermissionAdministrator
}

func (p Permission) Has(perm Permission) bool {
	if p.Admin() {
		return true
	}
	return (p & perm) == perm
}

func (p Permission) Add(perms ...Permission) (ret Permission) {
	for _, perm := range perms {
		ret |= perm
	}
	return
}

func (p Permission) Remove(perms ...Permission) (ret Permission) {
	for _, perm := range perms {
		ret &= ^perm
	}
	return
}

func (p Permission) Toggle(perms ...Permission) (ret Permission) {
	for _, perm := range perms {
		ret ^= perm
	}
	return
}

const (
	PermissionCreateInstantInvite              Permission = 1 << 0
	PermissionKickMembers                      Permission = 1 << 1
	PermissionBanMembers                       Permission = 1 << 2
	PermissionAdministrator                    Permission = 1 << 3
	PermissionManageChannels                   Permission = 1 << 4
	PermissionManageGuild                      Permission = 1 << 5
	PermissionAddReactions                     Permission = 1 << 6
	PermissionViewAuditLog                     Permission = 1 << 7
	PermissionPrioritySpeaker                  Permission = 1 << 8
	PermissionStream                           Permission = 1 << 9
	PermissionViewChannel                      Permission = 1 << 10
	PermissionSendMessages                     Permission = 1 << 11
	PermissionSendTTSMessages                  Permission = 1 << 12
	PermissionManageMessages                   Permission = 1 << 13
	PermissionEmbedLinks                       Permission = 1 << 14
	PermissionAttachFiles                      Permission = 1 << 15
	PermissionReadMessageHistory               Permission = 1 << 16
	PermissionMentionEveryone                  Permission = 1 << 17
	PermissionUseExternalEmojis                Permission = 1 << 18
	PermissionViewGuildInsights                Permission = 1 << 19
	PermissionConnect                          Permission = 1 << 20
	PermissionSpeak                            Permission = 1 << 21
	PermissionMuteMembers                      Permission = 1 << 22
	PermissionDeafenMembers                    Permission = 1 << 23
	PermissionMoveMembers                      Permission = 1 << 24
	PermissionUseVAD                           Permission = 1 << 25
	PermissionChangeNickname                   Permission = 1 << 26
	PermissionManageNicknames                  Permission = 1 << 27
	PermissionManageRoles                      Permission = 1 << 28
	PermissionManageWebhooks                   Permission = 1 << 29
	PermissionManageEmojisAndStickers          Permission = 1 << 30
	PermissionUseApplicationCommands           Permission = 1 << 31
	PermissionRequestToSpeak                   Permission = 1 << 32
	PermissionManageEvents                     Permission = 1 << 33
	PermissionManageThreads                    Permission = 1 << 34
	PermissionCreatePublicThreads              Permission = 1 << 35
	PermissionCreatePrivateThreads             Permission = 1 << 36
	PermissionUseExternalStickers              Permission = 1 << 37
	PermissionSendMessagesInThreads            Permission = 1 << 38
	PermissionUseEmbeddedActivities            Permission = 1 << 39
	PermissionModerateMembers                  Permission = 1 << 40
	PermissionViewCreatorMonetizationAnalytics Permission = 1 << 41
	PermissionUseSoundboard                    Permission = 1 << 42
	PermissionCreateGuildExpressions           Permission = 1 << 43
	PermissionCreateEvents                     Permission = 1 << 44
	PermissionUseExternalSounds                Permission = 1 << 45
	PermissionSendVoiceMessages                Permission = 1 << 46
	PermissionSendPolls                        Permission = 1 << 49
)

type PermissionOverwrite struct {
	ID    snowflake.ID            `json:"id,omitempty"`
	Type  PermissionOverwriteType `json:"type,omitempty"`
	Allow Permission              `json:"allow,omitempty"`
	Deny  Permission              `json:"deny,omitempty"`
}

func (o PermissionOverwrite) Denied(perm Permission) bool {
	return o.Deny.Has(perm)
}

func (o PermissionOverwrite) Allowed(perm Permission) bool {
	return o.Allow.Has(perm)
}

func (o PermissionOverwrite) Neutral(perm Permission) bool {
	return !(o.Denied(perm) || o.Allowed(perm))
}

type PermissionOverwriteType uint

const (
	PermissionOverwriteRole PermissionOverwriteType = iota
	PermissionOverwriteMember
)

type PermissionOverwrites []PermissionOverwrite

func (list PermissionOverwrites) Get(id snowflake.ID) (v PermissionOverwrite, _ bool) {
	for _, v = range list {
		if v.ID == id {
			return v, true
		}
	}
	return
}
