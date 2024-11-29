package client

import (
	"github.com/BOOMfinity/bfcord/discord"
	"github.com/andersfylling/snowflake/v5"
	"slices"
)

func MemberHighestRole(roles []discord.Role, member []snowflake.ID) discord.Role {
	slices.SortStableFunc(roles, func(a, b discord.Role) int {
		return b.Position - a.Position
	})
	for _, role := range roles {
		if slices.ContainsFunc(member, func(id snowflake.ID) bool {
			return id == role.ID
		}) {
			return role
		}
	}
	return roles[len(roles)-1]
}

func SortedMemberRoles(guild []discord.Role, member []snowflake.ID) (result []discord.Role) {
	slices.SortStableFunc(guild, func(a, b discord.Role) int {
		return b.Position - a.Position
	})
	for _, role := range guild {
		if slices.ContainsFunc(member, func(id snowflake.ID) bool {
			return id == role.ID
		}) {
			result = append(result, role)
		}
	}
	return
}
