package client

import (
	"slices"

	"github.com/BOOMfinity/bfcord/discord"
)

func SortedGuildRoles(roles []discord.Role) []discord.Role {
	slices.SortStableFunc(roles, func(a, b discord.Role) int {
		if a.Position == b.Position {
			return int(a.ID - b.ID)
		} else {
			return b.Position - a.Position
		}
	})
	return roles
}
