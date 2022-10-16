/*
Package client provides a user-friendly wrapper over gateway, cache and API.

# Getting started

Create a client and register event listeners

	bot, err := client.New("token") // options: client.WithX()

	bot.Sub().Interaction(func(bot client.Client, shard *gateway.Shard, ev *interactions.Interaction) {
		HandleInteraction(bot, cfg, ev) // your implementation
	})

	bot.Sub().MessageCreate(func(bot client.Client, shard *gateway.Shard, ev discord.Message) {
		HandleMessage(bot, ev) // your implementation
	})

	// start the client
	err = bot.Start(context.Background())
	if err != nil {
		panic(err)
	}
	bot.Wait()

# Accessing API

API is accessible using [Client].API() method. Refer to [api] package for more information.

# Accessing data

There are 3 ways of accessing various discord data.
  - With client helper methods, which use both cache and API (defined in [discord.ClientQuery])
  - Using cache ([Client].Store(), [cache.Store])
  - Using API (above)

# Accessing gateway

Current shard is transmitted in all events; To access any shard from client, use [Client].Get()

There are also two gateway methods accessible from client:
  - FetchMembers()
  - ChangeVoiceState()
*/
package client
