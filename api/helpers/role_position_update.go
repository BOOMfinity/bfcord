package helpers

import (
	"github.com/BOOMfinity/bfcord/api"
	"slices"

	"github.com/andersfylling/snowflake/v5"

	"github.com/BOOMfinity/bfcord/discord"
)

func UpdateRolePosition(roles []discord.Role, positions map[snowflake.ID]uint) []api.GuildRolePosition {
	data := make([]api.GuildRolePosition, 0, len(roles))
	slices.SortStableFunc(roles, func(a, b discord.Role) int {
		return a.Position - b.Position
	})
	for roleID, pos := range positions {
		index := slices.IndexFunc(roles, func(obj discord.Role) bool {
			return obj.ID == roleID
		})
		if index == -1 {
			continue
		}
		if pos > uint(len(roles)) {
			pos = uint(len(roles))
		} else if pos == 0 {
			pos = 1
		}
		role := roles[index]
		roles = slices.DeleteFunc(roles, func(obj discord.Role) bool {
			return obj.ID == roleID
		})
		roles = slices.Insert(roles, int(pos-1), role)
	}
	for i, r := range roles {
		data = append(data, api.GuildRolePosition{
			ID:       r.ID,
			Position: uint(i + 1),
		})
	}
	return data
}
