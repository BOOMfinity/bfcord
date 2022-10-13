package permissions

import (
	"github.com/andersfylling/snowflake/v5"
)

type Overwrite struct {
	ID    snowflake.ID  `json:"id"`
	Type  OverwriteType `json:"type"`
	Allow Permission    `json:"allow"`
	Deny  Permission    `json:"deny"`
}

type OverwriteType uint8

const (
	RoleOverwrite OverwriteType = iota
	MemberOverwrite
)
