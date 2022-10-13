package gateway

import (
	"context"
	"errors"
	"math/rand"
	"time"

	"github.com/BOOMfinity/bfcord/discord"
	"github.com/andersfylling/snowflake/v5"
)

func (g *Gateway) RequestMembers(ctx context.Context, opts RequestMembers) (members []discord.MemberWithUser, presences []discord.BasePresence, err error) {
	read := make(chan RequestMembersResponse)
	g.mut.Lock()
generate:
	key := randomString(32)
	_, ok := g.membersChannel[key]
	if ok {
		goto generate
	}
	opts.Nonce = key
	g.membersChannel[key] = read
	g.mut.Unlock()
	err = g.ws.WriteJSON(ctx, map[string]any{
		"op": RequestGuildMembersOp,
		"d":  opts,
	})
	if err != nil {
		return
	}
	for {
		select {
		case data := <-read:
			for i := range data.Members {
				data.Members[i].UserID = data.Members[i].User.ID
				data.Members[i].GuildID = opts.GuildID
			}
			members = append(members, data.Members...)
			for i := range data.Presences {
				data.Presences[i].UserID = data.Presences[i].User.ID
				presences = append(presences, data.Presences[i].BasePresence)
			}
			if data.ChunkIndex+1 == data.ChunkCount {
				g.mut.Lock()
				delete(g.membersChannel, key)
				g.mut.Unlock()
				return
			}
		case <-ctx.Done():
			return nil, nil, err
		case <-time.After(15 * time.Second):
			return nil, nil, errors.New("timed out")
		}
	}
}

type RequestMembers struct {
	Query     string         `json:"query"`
	Nonce     string         `json:"nonce,omitempty"`
	UserIDs   []snowflake.ID `json:"user_ids,omitempty"`
	GuildID   snowflake.ID   `json:"guild_id"`
	Limit     int            `json:"limit"`
	Presences bool           `json:"presences,omitempty"`
}

type RequestMembersResponse struct {
	Nonce      string                   `json:"nonce"`
	Members    []discord.MemberWithUser `json:"members"`
	NotFound   []snowflake.ID           `json:"not_found"`
	Presences  []discord.Presence       `json:"presences"`
	ChunkIndex int                      `json:"chunk_index"`
	ChunkCount int                      `json:"chunk_count"`
	GuildID    snowflake.ID             `json:"guild_id"`
}

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func randomString(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Int63()%int64(len(letterBytes))]
	}
	return string(b)
}
