package discord

import (
	"github.com/BOOMfinity/bfcord/api/images"
	"github.com/andersfylling/snowflake/v5"
)

type GuildBuilder interface {
	Name(str string) GuildBuilder
	Description(desc string) GuildBuilder
	PremiumProgressBar(enabled bool) GuildBuilder
	VerificationLevel(lvl GuildVerificationLevel) GuildBuilder
	DefaultMessageNotifications(n GuildDefaultNotifications) GuildBuilder
	ExplicitContentFilter(f GuildExplicitContentFilter) GuildBuilder
	AFKChannelID(id snowflake.ID) GuildBuilder
	AFKTimeout(timeout uint32) GuildBuilder
	Icon(image *images.MediaBuilder) GuildBuilder
	TransferOwner(id snowflake.ID) GuildBuilder
	Splash(image *images.MediaBuilder) GuildBuilder
	DiscoverySplash(image *images.MediaBuilder) GuildBuilder
	Banner(image *images.MediaBuilder) GuildBuilder
	SystemChannel(id snowflake.ID, flags ...SystemChannelFlag) GuildBuilder
	RulesChannel(id snowflake.ID) GuildBuilder
	PublicUpdatesChannel(id snowflake.ID) GuildBuilder
	Locale(l string) GuildBuilder
	Execute(api ClientQuery, reason ...string) (guild *Guild, err error)
}
