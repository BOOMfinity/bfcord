package client

import (
	"fmt"
	"github.com/BOOMfinity/bfcord/discord"
	"github.com/andersfylling/snowflake/v5"
)

func (s *sessionImpl) PermissionsIn(guildID, channelID, memberID snowflake.ID) (discord.Permission, error) {
	guild, err := s.Guild(guildID).Get()
	if err != nil {
		return 0, fmt.Errorf("failed to get the guild: %w", err)
	}
	channel, err := s.Channel(channelID).Get()
	if err != nil {
		return 0, fmt.Errorf("failed to get the channel: %w", err)
	}
	member, err := s.Guild(guildID).Member(memberID).Get()
	if err != nil {
		return 0, fmt.Errorf("failed to get the member: %w", err)
	}

	return ComputePermissions(guild, member, channel), nil
}

func ComputeBasePermissions(guild discord.Guild, member discord.MemberWithUser) (p discord.Permission) {
	if guild.OwnerID == member.User.ID {
		return discord.PermissionAdministrator
	}

	everyone, _ := guild.Role(guild.ID)

	p = everyone.Permissions

	for _, rid := range member.Roles {
		if rid == guild.ID {
			continue
		}
		role, _ := guild.Role(rid)
		p = p.Add(role.Permissions)
	}
	return
}

func ComputeOverwrites(base discord.Permission, guild discord.Guild, member discord.MemberWithUser, overwrites discord.PermissionOverwrites) (p discord.Permission) {
	if base.Has(discord.PermissionAdministrator) {
		return base
	}

	p = base

	{
		everyone, ok := overwrites.Get(guild.ID)
		if ok {
			p = p.Remove(everyone.Deny)
			p = p.Add(everyone.Allow)
		}
	}

	{
		allow := discord.Permission(0)
		deny := discord.Permission(0)

		for _, rid := range member.Roles {
			if rid == guild.ID {
				continue
			}
			role, ok := overwrites.Get(rid)
			if ok {
				allow = allow.Add(role.Allow)
				deny = deny.Add(role.Deny)
			}
		}

		p = p.Remove(deny)
		p = p.Add(allow)
	}
	{
		memberOverwrite, ok := overwrites.Get(member.User.ID)
		if ok {
			p = p.Remove(memberOverwrite.Deny)
			p = p.Add(memberOverwrite.Allow)
		}
	}
	return
}

func ComputePermissions(guild discord.Guild, member discord.MemberWithUser, channel discord.Channel) (p discord.Permission) {
	base := ComputeBasePermissions(guild, member)
	return ComputeOverwrites(base, guild, member, channel.PermissionOverwrites)
}
