package discord

import (
	"fmt"
	"github.com/BOOMfinity/bfcord/api/images"
	"github.com/BOOMfinity/bfcord/discord/permissions"
	"github.com/BOOMfinity/bfcord/errs"
	"github.com/BOOMfinity/bfcord/internal/slices"
	"github.com/andersfylling/snowflake/v5"
	"sort"
)

type RoleBuilder interface {
	Name(str string) RoleBuilder
	Permissions(perms permissions.Permission) RoleBuilder
	Color(c int64) RoleBuilder
	ShowSeparately(bool) RoleBuilder
	Icon(i *images.MediaBuilder) RoleBuilder
	UnicodeEmoji(str string) RoleBuilder
	Mentionable(bool) RoleBuilder
	Execute(api ClientQuery, reasons ...string) (role Role, err error)
}

type RolePositions []Role

func (r *RolePositions) Map() (d []map[string]any) {
	for i := range *r {
		d = append(d, map[string]any{
			"id":       (*r)[i].ID,
			"position": i + 1,
		})
	}
	return
}

func (r *RolePositions) Set(role snowflake.ID, pos uint8) error {
	sort.SliceStable(*r, func(a, z int) bool {
		return (*r)[a].Position < (*r)[z].Position
	})
	for i := range *r {
		_x := (*r)[i]
		fmt.Println(_x.ID, _x.Position)
	}
	_currIndex := slices.FindIndex(*r, func(item Role) bool {
		return item.ID == role
	})
	if _currIndex == -1 {
		return errs.ItemNotFound
	}
	if pos == 0 {
		pos = 1
	}
	if int(pos) > len(*r) {
		pos = uint8(len(*r))
	}
	_roleCpy := (*r)[_currIndex]
	*r = append((*r)[:_currIndex], (*r)[_currIndex+1:]...)
	*r = append((*r)[:pos-1], append([]Role{_roleCpy}, (*r)[pos-1:]...)...)
	for i := range *r {
		(*r)[i].Position = i + 1
	}
	println()
	for i := range *r {
		_x := (*r)[i]
		fmt.Println(_x.ID, _x.Position)
	}
	return nil
}
