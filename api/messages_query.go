package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/BOOMfinity/bfcord/discord"
	"github.com/andersfylling/snowflake/v5"
)

type _data_ struct {
	Around snowflake.ID `json:"around,omitempty"`
	After  snowflake.ID `json:"after,omitempty"`
	Before snowflake.ID `json:"before,omitempty"`
	Limit  uint16       `json:"limit,omitempty"`
}

type ChannelMessagesQuery struct {
	api     *Client
	channel snowflake.ID
}

func (v ChannelMessagesQuery) Around(id snowflake.ID, limit uint16) (msgs []discord.Message, err error) {
	if !id.Valid() {
		return nil, nil
	}
	if limit > 100 {
		limit = 100
	}
	if limit < 1 {
		limit = 50
	}
	raw, err := json.Marshal(_data_{
		Limit:  limit,
		Around: id,
	})
	if err != nil {
		return
	}
	req := v.api.New(true)
	req.SetHost(fmt.Sprintf(FullApiUrl+"/channels/%v/messages", v.channel))
	req.SetBody(raw)
	err = v.api.DoResult(req, &msgs)
	for i := range msgs {
		msgs[i].Patch()
	}
	return
}

func (v ChannelMessagesQuery) After(ctx context.Context, id snowflake.ID, limit uint16) (msgs []discord.Message, err error) {
	if !id.Valid() {
		return nil, nil
	}
	if limit < 1 {
		limit = 1000
	}
	fetchLimit := limit
	if fetchLimit > 100 {
		fetchLimit = 100
	}
	if fetchLimit < 1 {
		fetchLimit = 100
	}
	params := url.Values{}
	params.Set("limit", fmt.Sprint(fetchLimit))
	temp := make([]discord.Message, 0, 100)
	for uint16(len(msgs)) < limit {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
			params.Set("after", id.String())
			req := v.api.New(true)
			req.SetRequestURI(fmt.Sprintf("%v/channels/%v/messages?%v", FullApiUrl, v.channel.String(), params.Encode()))
			err = v.api.DoResult(req, &temp)
			if err != nil {
				return nil, err
			}
			id = temp[len(temp)-1].ID
			for i := range temp {
				temp[i].Patch()
			}
			msgs = append(msgs, temp...)
			if len(temp) < 100 {
				break
			}
			temp = make([]discord.Message, 0, 100)
		}
	}
	return
}

func (v ChannelMessagesQuery) Before(ctx context.Context, id snowflake.ID, limit uint16) (msgs []discord.Message, err error) {
	if !id.Valid() {
		return nil, nil
	}
	if limit < 1 {
		limit = 1000
	}
	fetchLimit := limit
	if fetchLimit > 100 {
		fetchLimit = 100
	}
	if fetchLimit < 1 {
		fetchLimit = 100
	}
	params := url.Values{}
	params.Set("limit", fmt.Sprint(fetchLimit))
	temp := make([]discord.Message, 0, 100)
	for uint16(len(msgs)) < limit {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
			params.Set("before", id.String())
			req := v.api.New(true)
			req.SetRequestURI(fmt.Sprintf("%v/channels/%v/messages?%v", FullApiUrl, v.channel.String(), params.Encode()))
			err = v.api.DoResult(req, &temp)
			if err != nil {
				return nil, err
			}
			id = temp[len(temp)-1].ID
			for i := range temp {
				temp[i].Patch()
			}
			msgs = append(msgs, temp...)
			if len(temp) < 100 {
				break
			}
			temp = make([]discord.Message, 0, 100)
		}
	}
	return
}

func (v ChannelMessagesQuery) Latest(limit uint16) (msgs []discord.Message, err error) {
	if limit == 0 {
		limit = 50
	}
	req := v.api.New(true)
	req.SetRequestURI(fmt.Sprintf("%v/channels/%v/messages?limit=%v", FullApiUrl, v.channel, limit))
	err = v.api.DoResult(req, &msgs)
	for i := range msgs {
		msgs[i].Patch()
	}
	return
}

func (v ChannelMessagesQuery) ID() snowflake.ID {
	return v.channel
}

func NewChannelMessagesQuery(client *Client, channel snowflake.ID) ChannelMessagesQuery {
	data := ChannelMessagesQuery{
		channel: channel,
		api:     client,
	}
	//data.emptyOptions = emptyOptions[discord.MessageQuery]{data: data}
	return data
}
