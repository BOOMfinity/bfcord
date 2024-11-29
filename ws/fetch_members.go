package ws

import (
	"context"
	"fmt"
	"time"

	"github.com/andersfylling/snowflake/v5"
	"github.com/segmentio/encoding/json"

	"github.com/BOOMfinity/bfcord"
	"github.com/BOOMfinity/bfcord/discord"
)

func (g *gatewayImpl) FetchMembers(ctx context.Context, params RequestGuildMembersParams) (_ []discord.MemberWithUser, _ []discord.Presence, err error) {
	if g.Status() != StatusConnected {
		return nil, nil, ErrGatewayNotConnected
	}
	g.sendEvent(sendEvent[RequestGuildMembersParams]{
		OpCode: 8,
		Data:   params,
	})
	events, close := g.Listen()
	defer close()
	var (
		members   []discord.MemberWithUser
		presences []discord.Presence
		notFound  []snowflake.ID
	)
	data := new(GuildMembersChunkEvent)
	timer := time.NewTimer(15 * time.Second)
	defer func() {
		if !timer.Stop() {
			<-timer.C
		}
		if err == nil && len(notFound) > 0 {
			err = bfcord.PointerOf(ErrNotFound(notFound))
		}
	}()
loop:
	for {
	sel:
		select {
		case msg := <-events:
			if ev, ok := msg.(InternalDispatchEvent); ok {
				if !timer.Stop() {
					<-timer.C
				}
				defer timer.Reset(15 * time.Second)
				defer ev.Dereference()
				if ev.Event != "GUILD_MEMBERS_CHUNK" {
					break sel
				}
				if err = json.Unmarshal(ev.Data, &data); err != nil {
					return members, presences, fmt.Errorf("failed to unmarshal chunk: %w", err)
				}
				if data.Nonce != params.Nonce {
					break sel
				}
				if data.ChunkIndex == 0 && data.ChunkCount > 1 {
					members = make([]discord.MemberWithUser, 0, len(data.Members)*int(data.ChunkCount))
					presences = make([]discord.Presence, 0, len(data.Presences)*int(data.ChunkCount))
				}
				members = append(members, data.Members...)
				presences = append(presences, data.Presences...)
				notFound = append(notFound, data.NotFound...)
				if data.ChunkIndex-1 == data.ChunkCount {
					break loop
				}
			}
		case <-timer.C:
			return members, presences, ErrFetchingMembersTimedOut
		case <-ctx.Done():
			return members, presences, ctx.Err()
		}
	}
	return members, presences, nil
}
