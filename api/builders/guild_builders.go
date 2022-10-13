package builders

import (
	"fmt"
	"github.com/BOOMfinity/bfcord/api/images"
	"github.com/BOOMfinity/bfcord/discord"
	"github.com/BOOMfinity/go-utils/nullable"
	"github.com/andersfylling/snowflake/v5"
	"strings"
)

func NewGuildBuilder(id snowflake.ID) *GuildBuilder {
	return &GuildBuilder{id: id}
}

type GuildBuilder struct {
	data                   discord.GuildUpdate
	iconBuilder            *images.MediaBuilder
	splashBuilder          *images.MediaBuilder
	discoverySplashBuilder *images.MediaBuilder
	bannerBuilder          *images.MediaBuilder
	id                     snowflake.ID
}

func (g *GuildBuilder) Name(str string) discord.GuildBuilder {
	g.data.Name = &str
	return g
}

func (g *GuildBuilder) Description(desc string) discord.GuildBuilder {
	g.data.Description = &desc
	return g
}

func (g *GuildBuilder) PremiumProgressBar(enabled bool) discord.GuildBuilder {
	g.data.PremiumProgressBarEnabled = &enabled
	return g
}

func (g *GuildBuilder) VerificationLevel(lvl discord.GuildVerificationLevel) discord.GuildBuilder {
	g.data.VerificationLevel = &lvl
	return g
}

func (g *GuildBuilder) DefaultMessageNotifications(n discord.GuildDefaultNotifications) discord.GuildBuilder {
	g.data.DefaultMessageNotifications = &n
	return g
}

func (g *GuildBuilder) ExplicitContentFilter(f discord.GuildExplicitContentFilter) discord.GuildBuilder {
	g.data.ExplicitContentFilter = &f
	return g
}

func (g *GuildBuilder) AFKChannelID(id snowflake.ID) discord.GuildBuilder {
	g.data.AFKChannelID = nullable.New[snowflake.ID]()
	if id == 0 {
		g.data.AFKChannelID.Clear()
	} else {
		g.data.AFKChannelID.Set(id)
	}
	return g
}

func (g *GuildBuilder) AFKTimeout(timeout uint32) discord.GuildBuilder {
	g.data.AFKTimeout = &timeout
	return g
}

func (g *GuildBuilder) Icon(image *images.MediaBuilder) discord.GuildBuilder {
	g.data.Icon = nullable.New[string]()
	if image == nil {
		g.iconBuilder = images.NewMediaBuilder()
	} else {
		g.iconBuilder = image
	}
	return g
}

func (g *GuildBuilder) TransferOwner(id snowflake.ID) discord.GuildBuilder {
	g.data.OwnerID = &id
	return g
}

func (g *GuildBuilder) Splash(image *images.MediaBuilder) discord.GuildBuilder {
	g.data.Splash = nullable.New[string]()
	if image == nil {
		g.splashBuilder = images.NewMediaBuilder()
	} else {
		g.splashBuilder = image
	}
	return g
}

func (g *GuildBuilder) DiscoverySplash(image *images.MediaBuilder) discord.GuildBuilder {
	g.data.DiscoverySplash = nullable.New[string]()
	if image == nil {
		g.discoverySplashBuilder = images.NewMediaBuilder()
	} else {
		g.discoverySplashBuilder = image
	}
	return g
}

func (g *GuildBuilder) Banner(image *images.MediaBuilder) discord.GuildBuilder {
	g.data.Banner = nullable.New[string]()
	if image == nil {
		g.iconBuilder = images.NewMediaBuilder()
	} else {
		g.iconBuilder = image
	}
	return g
}

func (g *GuildBuilder) SystemChannel(id snowflake.ID, flags ...discord.SystemChannelFlag) discord.GuildBuilder {
	g.data.SystemChannelID = nullable.New[snowflake.ID]()
	if id == 0 {
		g.data.SystemChannelID.Clear()
	} else {
		g.data.SystemChannelID.Set(id)
	}
	var _t discord.SystemChannelFlag
	if flags != nil {
		g.data.SystemChannelFlags = &_t
		for i := range flags {
			*g.data.SystemChannelFlags |= flags[i]
		}
	} else {
		g.data.SystemChannelFlags = &_t
	}
	return g
}

func (g *GuildBuilder) RulesChannel(id snowflake.ID) discord.GuildBuilder {
	g.data.RulesChannelID = &id
	return g
}

func (g *GuildBuilder) PublicUpdatesChannel(id snowflake.ID) discord.GuildBuilder {
	g.data.PublicUpdatesChannelID = &id
	return g
}

func (g *GuildBuilder) Locale(l string) discord.GuildBuilder {
	g.data.PreferredLocale = &l
	return g
}

func (g *GuildBuilder) Execute(api discord.ClientQuery, reason ...string) (guild *discord.Guild, err error) {
	if g.bannerBuilder != nil {
		b, _err := g.bannerBuilder.ToBase64()
		if _err != nil {
			err = fmt.Errorf("failed to parse image to base64: %w", _err)
			return
		}
		g.data.Banner.Set(b)
	}
	if g.splashBuilder != nil {
		b, _err := g.splashBuilder.ToBase64()
		if _err != nil {
			err = fmt.Errorf("failed to parse image to base64: %w", _err)
			return
		}
		g.data.Splash.Set(b)
	}
	if g.discoverySplashBuilder != nil {
		b, _err := g.discoverySplashBuilder.ToBase64()
		if _err != nil {
			err = fmt.Errorf("failed to parse image to base64: %w", _err)
			return
		}
		g.data.DiscoverySplash.Set(b)
	}
	if g.iconBuilder != nil {
		b, _err := g.iconBuilder.ToBase64()
		if _err != nil {
			err = fmt.Errorf("failed to parse image to base64: %w", _err)
			return
		}
		g.data.Icon.Set(b)
	}
	return api.LowLevel().Reason(strings.Join(reason, " ")).UpdateGuild(g.id, g.data)
}
