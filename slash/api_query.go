package slash

import "github.com/andersfylling/snowflake/v5"

type Query interface {
	Global() GlobalQuery
	Guild(id snowflake.ID) GuildQuery
}

type GuildQuery interface {
	Get(id snowflake.ID) (Command, error)
	Create(name, desc string) CommandBuilder
	Edit(id snowflake.ID) EditCommandBuilder
	Delete(id snowflake.ID) error
	Permissions() (PermissionList, error)
	CommandPermissions(id snowflake.ID) (PermissionList, error)
	EditPermissions(id snowflake.ID, perms PermissionList) error
}

type GlobalQuery interface {
	List() ([]Command, error)
	Create(name, desc string) CommandBuilder
	Get(id snowflake.ID) (Command, error)
	Edit(id snowflake.ID) EditCommandBuilder
	Delete(id snowflake.ID) error
}
