package slash

import "github.com/andersfylling/snowflake/v5"

// CommandPermissionType https://discord.com/developers/docs/interactions/application-commands#application-command-permissions-object-application-command-permission-type
type CommandPermissionType int

const (
	PermissionTypeRole CommandPermissionType = iota + 1
	PermissionTypeUser
)

// CommandPermission https://discord.com/developers/docs/interactions/application-commands#application-command-permissions-object-application-command-permissions-structure
type CommandPermission struct {
	ID         snowflake.ID          `json:"id"`
	Type       CommandPermissionType `json:"type"`
	GuildID    snowflake.ID          `json:"guild_id"`
	Permission bool                  `json:"permission"`
}

type PermissionList []CommandPermission

func (x PermissionList) Get(id snowflake.ID) *CommandPermission {
	for i := range x {
		if x[i].ID == id {
			return &x[i]
		}
	}
	return nil
}

func (x *PermissionList) Allow(t CommandPermissionType, id snowflake.ID) {
	perm := x.Get(id)
	if perm != nil {
		perm.Permission = true
		return
	}
	*x = append(*x, CommandPermission{ID: id, Type: t, Permission: true})
}

func (x *PermissionList) Disallow(t CommandPermissionType, id snowflake.ID) {
	perm := x.Get(id)
	if perm != nil {
		perm.Permission = false
		return
	}
	*x = append(*x, CommandPermission{ID: id, Type: t, Permission: false})
}

func (x PermissionList) findIndex(id snowflake.ID) int {
	for i := range x {
		if x[i].ID == id {
			return i
		}
	}
	return -1
}

func (x *PermissionList) Remove(id snowflake.ID) bool {
	if index := x.findIndex(id); index != -1 {
		*x = append((*x)[:index], (*x)[index+1:]...)
		return true
	}
	return false
}
