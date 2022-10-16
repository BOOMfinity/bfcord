[![Go Reference](https://pkg.go.dev/badge/github.com/BOOMfinity/bfcord.svg)](https://pkg.go.dev/github.com/BOOMfinity/bfcord)
# bfcord
### Discord API library written in golang with focus on simplicity and performance

# EARLY ACCESS / BETA
For this moment we CAN'T guarantee full stability for longer time perspective. The library is actively tested on our multiple bots and by our testers.

## Getting started
### Installation
```
go get github.com/BOOMfinity/bfcord
```
### Basic client example
Bfcord uses options pattern in more complex constructors
```go
// create client
bot, err := client.New("token", client.WithIntents(intents.X|intents.Y) /*other options: client.WithX()*/)

// register events
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
```

# Packages and features
## Client ([godoc]())
Client is an intuitive wrapper over gateway, api and cache packages. It provides convenient methods for interacting with discord, listening for events, and fetching data from API and cache.
## Slash ([godoc]()) and Interactive ([godoc]())
Slash package provides basic API client for managing slash commands

Interactive is a utility package giving easier ways to use and interact with message components (buttons, lists, etc.)
## Cache ([godoc]())
Cache is accessible under client.Store() method. However, you can (should?) use methods built into client whenever possible (ex: bot.Guild(), Channel(), User()), as they automatically fall back to API if object is not in cache yet.

bfcord uses modern and performant stores built with SliceMap approach and generics. This ensures low load on GC while maintaining high performance.
### Examples
#### Pull a guild (using client's method)
```go
guild, err := bot.Guild(613425648685547541).Get()
```
#### Pull a guild (cache only)
```go
guild, err := bot.Store().Guilds().Get(613425648685547541)
```
#### Find guilds with more than 1000 members
```go
guilds := bot.Store().Guilds().Filter(func(guild discord.BaseGuild) bool {
	return guild.MemberCount > 1000
})
```
## API ([godoc]())
> **Note**
> Not all endpoints are implemented at this moment - if something is missing and you need it - notify us 

API is accessible under client.API() method, and is split into logical categories.
For more complex endpoints builder pattern is used.
### Examples
#### Delete a channel
```go
bot.API().Channel(41514364364367).Delete()
```
#### Send a message with embed
```go
msg, err := bot.API().Channel(1654564363426634).SendMessage().Embed(discord.MessageEmbed{
	Description: "example",
	Title:       "example",
}).Execute(bot)
```
#### Edit a guild
```go
bot.API().Guild(3213232312312).Edit().Name("new name").Description("new description").AFKChannelID(1232312312414).Execute(bot, "test reason")
```

## Gateway ([godoc]())
Gateway methods are accessible within client:
* [`client.ChangeVoiceState()`]()
* [`clientczy .FetchMembers()`]()

Refer to docs for more advanced usage.
## Voice
See: [readme.md]() and [godoc]() of voice package

# Modular usage
Each of the packages can be used independently (excl. client). For example, you can grab just the gateway and api package, and write your own logic around it.

# Why another library? Is it better than \<insert any library here\> ?
We don't want to compare ourselves. This library was written because other libraries weren't suitable for our needs (and we wanted to get some experience in golang). 

However, we can say that it's probably the most performant and memory efficient discord library right now.
