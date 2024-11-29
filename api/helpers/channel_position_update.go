package helpers

import (
	"github.com/BOOMfinity/bfcord/api"
	"slices"

	"github.com/BOOMfinity/go-utils/inlineif"
	"github.com/andersfylling/snowflake/v5"

	"github.com/BOOMfinity/bfcord/discord"
)

func UpdateChannelPosition(channels []discord.Channel, positions map[snowflake.ID]api.GuildChannelPosition) []api.GuildChannelPosition {
	data := make([]api.GuildChannelPosition, 0, len(channels))
	sameType := make([]discord.Channel, 0, len(channels))
	for id, params := range positions {
		sameType = make([]discord.Channel, 0, len(channels))
		// make channel slice copy to prevent working on original slice
		copied := make([]discord.Channel, len(channels))
		copy(copied, channels)

		index := slices.IndexFunc(copied, func(obj discord.Channel) bool {
			return obj.ID == id
		})

		// this should not happen but just in case - just skip :P
		if index == -1 {
			continue
		}

		if params.ParentID != nil {
			copied[index].ParentID = *params.ParentID
		}

		orig := copied[index]

		// generate slice of channels with the same "sorting" type + skip channels that don't have the same parent id.
		// for example, all text-based channels are sorted together (news, forum, etc...)
		for _, obj := range copied {
			if obj.ID == id {
				continue
			}
			if !SameChannelSortingType(orig, obj) {
				continue
			}
			if orig.ParentID != obj.ParentID {
				continue
			}
			sameType = append(sameType, obj)
		}

		slices.SortStableFunc(sameType, func(a, b discord.Channel) int {
			if a.Position == b.Position {
				return int(a.ID - b.ID)
			}
			return int(a.Position) - int(b.Position)
		})

		// check if position to set is not too high or too low.
		if params.Position > uint(len(sameType)) {
			params.Position = uint(len(sameType)) + 1
		}

		if params.Position == 0 {
			params.Position = 1
		}

		sameType = slices.Insert(sameType, int(params.Position)-1, orig)

		// update positions and remove duplications
		for i, obj := range sameType {
			sameType[i].Position = uint(i + 1)
			for ii, v := range copied {
				if v.ID == obj.ID {
					copied[ii].Position = uint(i + 1)
				}
			}
			// make sure there is no duplication when sorting multiple channels at the same time.
			data = slices.DeleteFunc(data, func(v api.GuildChannelPosition) bool {
				return v.ID == obj.ID
			})
		}

		// after updating positions and checking for duplications, we are pushing all changed channels to the data slice that will be sent to API.
		for _, obj := range sameType {
			data = append(data, api.GuildChannelPosition{
				Position:        obj.Position,
				ParentID:        inlineif.IfElse(obj.ID == id, params.ParentID, nil),
				LockPermissions: inlineif.IfElse(obj.ID == id, params.LockPermissions, false),
				ID:              obj.ID,
			})
		}

		// update original slice for next position update
		channels = make([]discord.Channel, len(channels))
		copy(channels, copied)
	}
	return data
}

func check(a, b discord.Channel) bool {
	if a.Thread() {
		return b.Thread()
	}
	if a.Type == discord.ChannelTypeVoice {
		switch b.Type {
		case discord.ChannelTypeVoice, discord.ChannelTypeStageVoice:
			return true
		}
		return false
	}
	if a.Type == discord.ChannelTypeCategory {
		return b.Type == discord.ChannelTypeCategory
	}
	return true
}

func SameChannelSortingType(orig, ch discord.Channel) bool {
	return check(orig, ch) && check(ch, orig)
}
