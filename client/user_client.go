package client

import (
	"github.com/BOOMfinity/bfcord/api"
	"github.com/BOOMfinity/bfcord/discord"
	"github.com/andersfylling/snowflake/v5"
)

func getOrSet[T any](sess Session, getFn func() (T, error), fetchFn func() (T, error), save func(data T) error) (T, error) {
	if sess.Cache() != nil {
		data, err := getFn()
		if err == nil {
			return data, nil
		}
	}
	data, err := fetchFn()
	if err != nil {
		return data, err
	}
	if sess.Cache() != nil {
		if err = save(data); err != nil {
			return data, err
		}
	}
	return data, nil
}

type userClient struct {
	api.UserClient
	id   snowflake.ID
	sess Session
}

func (c userClient) Get() (discord.User, error) {
	return getOrSet(c.sess, func() (discord.User, error) {
		return c.sess.Cache().Users().Get(c.id)
	}, func() (discord.User, error) {
		return c.UserClient.Get()
	}, func(data discord.User) error {
		return c.sess.Cache().Users().Set(c.id, data)
	})
}
