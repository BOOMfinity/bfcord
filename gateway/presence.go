package gateway

import (
	"context"
	"github.com/BOOMfinity/bfcord/discord"
)

type PresenceSet interface {
	Set(status discord.PresenceStatus, ac discord.Activity)
	SetCustom(data PresenceUpdate)
}

type presenceSet struct {
	g *Gateway
}

func (v *presenceSet) Set(status discord.PresenceStatus, ac discord.Activity) {
	v.SetCustom(PresenceUpdate{
		Status:     status,
		Activities: []discord.Activity{ac},
	})
}

func (v *presenceSet) SetCustom(data PresenceUpdate) {
	_ = v.g.ws.WriteJSON(context.Background(), map[string]any{
		"op": PresenceUpdateOp,
		"d":  data,
	})
}
