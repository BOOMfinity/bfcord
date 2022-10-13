package api

import (
	"fmt"
	"net/url"

	"github.com/BOOMfinity/bfcord/discord"
	"github.com/andersfylling/snowflake/v5"
	"github.com/valyala/fasthttp"
)

type ReactionQuery struct {
	emptyOptions[discord.MessageReactionQuery]
	client  *Client
	emoji   string
	channel snowflake.ID
	message snowflake.ID
}

func (r ReactionQuery) RemoveAll() (err error) {
	req := r.client.New(true)
	req.SetRequestURI(fmt.Sprintf(FullApiUrl+"/channels/%v/messages/%v/reactions/%v", r.channel, r.message, r.emoji))
	req.Header.SetMethod(fasthttp.MethodDelete)
	return r.client.DoNoResp(req)
}

func (r ReactionQuery) RemoveOwn() (err error) {
	req := r.client.New(true)
	req.SetRequestURI(fmt.Sprintf(FullApiUrl+"/channels/%v/messages/%v/reactions/%v/@me", r.channel, r.message, r.emoji))
	req.Header.SetMethod(fasthttp.MethodDelete)
	return r.client.DoNoResp(req)
}

func (r ReactionQuery) After(limit uint64, after snowflake.ID) (users []discord.User, err error) {
	if limit != 0 {
		users = make([]discord.User, 0, limit)
	}
	fetch := uint64(100)
	unlimited := limit == 0
	for limit > 0 || unlimited {
		if limit > 0 {
			if fetch > limit {
				fetch = limit
			}
			limit -= fetch
		}

		reactions, err := r.Range(fetch, after)
		if err != nil {
			return nil, err
		}
		users = append(reactions, users...)
		if uint64(len(reactions)) < fetch {
			break
		}
		after = users[len(users)-1].ID
	}

	return
}

func (r ReactionQuery) Range(limit uint64, after snowflake.ID) (users []discord.User, err error) {
	switch {
	case limit == 0:
		limit = 25
		break
	case limit > 100:
		limit = 100
		break
	}
	params := url.Values{}
	params.Add("limit", fmt.Sprint(limit))
	if after.Valid() {
		params.Add("after", after.String())
	}
	users = make([]discord.User, 0, limit)
	req := r.client.New(true)
	req.SetRequestURI(fmt.Sprintf(FullApiUrl+"/channels/%v/messages/%v/reactions/%v?%v", r.channel, r.message, r.emoji, params.Encode()))
	err = r.client.DoResult(req, &users)
	if err != nil {
		return
	}
	return
}

func (r ReactionQuery) All(limit uint64) (users []discord.User, err error) {
	if limit != 0 {
		users = make([]discord.User, 0, limit)
	}

	return r.After(limit, 0)
}

func (r ReactionQuery) Remove(userID snowflake.ID) (err error) {
	req := r.client.New(true)
	req.SetRequestURI(fmt.Sprintf(FullApiUrl+"/channels/%v/messages/%v/reactions/%v/%v", r.channel, r.message, r.emoji, userID))
	req.Header.SetMethod(fasthttp.MethodDelete)
	return r.client.DoNoResp(req)
}

func (r ReactionQuery) Emoji() string {
	return r.emoji
}

func (r ReactionQuery) Channel() snowflake.ID {
	return r.channel
}

func (r ReactionQuery) Message() snowflake.ID {
	return r.message
}

func NewReactionQuery(client *Client, channel, message snowflake.ID, emoji string) *ReactionQuery {
	data := &ReactionQuery{
		message: message,
		emoji:   emoji,
		channel: channel,
		client:  client,
	}
	data.emptyOptions = emptyOptions[discord.MessageReactionQuery]{data: data}
	return data
}
