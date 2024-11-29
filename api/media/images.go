package media

import (
	"fmt"
	"strings"

	"github.com/andersfylling/snowflake/v5"
)

const CDNUrl = "https://cdn.discordapp.com/"

func JoinPath(seg ...string) string {
	return CDNUrl + strings.Join(seg, "/")
}

func CustomEmoji(id snowflake.ID, opts ...CDNOption) string {
	options := parseOptions(opts)
	return JoinPath("emojis", fmt.Sprintf("%s.%s", id, options.Format))
}

func GuildIcon(id snowflake.ID, hash string, opts ...CDNOption) string {
	if hash == "" {
		return ""
	}
	options := parseOptions(opts)
	if strings.HasPrefix(hash, "a_") && options.Animated {
		options.Format = "gif"
	}
	return JoinPath("icons", id.String(), fmt.Sprintf("%s.%s", hash, options.Format))
}

func GuildSplash(id snowflake.ID, hash string, opts ...CDNOption) string {
	if hash == "" {
		return ""
	}
	options := parseOptions(opts)
	return JoinPath("splashes", id.String(), fmt.Sprintf("%s.%s", hash, options.Format))
}

func GuildDiscoverySplash(id snowflake.ID, hash string, opts ...CDNOption) string {
	if hash == "" {
		return ""
	}
	options := parseOptions(opts)
	return JoinPath("discovery-splashes", id.String(), fmt.Sprintf("%s.%s", hash, options.Format))
}

func GuildBanner(id snowflake.ID, hash string, opts ...CDNOption) string {
	if hash == "" {
		return ""
	}
	options := parseOptions(opts)
	if strings.HasPrefix(hash, "a_") && options.Animated {
		options.Format = "gif"
	}
	return JoinPath("banners", id.String(), fmt.Sprintf("%s.%s", hash, options.Format))
}

func UserBanner(id snowflake.ID, hash string, opts ...CDNOption) string {
	if hash == "" {
		return ""
	}
	options := parseOptions(opts)
	if strings.HasPrefix(hash, "a_") && options.Animated {
		options.Format = "gif"
	}
	return JoinPath("banners", id.String(), fmt.Sprintf("%s.%s", hash, options.Format))
}

func DefaultUserAvatar(id snowflake.ID) string {
	return JoinPath("embed", "avatars", fmt.Sprintf("%d.png", (uint64(id)>>uint64(22))%uint64(6)))
}

func UserAvatar(id snowflake.ID, hash string, opts ...CDNOption) string {
	options := parseOptions(opts)
	if hash == "" && options.FallbackToDefault {
		return DefaultUserAvatar(id)
	} else if hash == "" {
		return ""
	}
	if strings.HasPrefix(hash, "a_") && options.Animated {
		options.Format = "gif"
	}
	return JoinPath("avatars", id.String(), fmt.Sprintf("%s.%s", hash, options.Format))
}

func GuildMemberAvatar(guild, user snowflake.ID, hash string, opts ...CDNOption) string {
	options := parseOptions(opts)
	if hash == "" {
		return ""
	}
	if strings.HasPrefix(hash, "a_") && options.Animated {
		options.Format = "gif"
	}
	return JoinPath("guilds", guild.String(), "users", user.String(), "avatars", fmt.Sprintf("%s.%s", hash, options.Format))
}

func AvatarDecoration(hash string) string {
	return JoinPath("avatar-decoration-presets", fmt.Sprintf("%s.png", hash))
}

func ApplicationIcon(id snowflake.ID, hash string, opts ...CDNOption) string {
	if hash == "" {
		return ""
	}
	options := parseOptions(opts)
	return JoinPath("app-icons", id.String(), fmt.Sprintf("%s.%s", hash, options.Format))
}

func ApplicationAsset(id snowflake.ID, hash string, opts ...CDNOption) string {
	if hash == "" {
		return ""
	}
	options := parseOptions(opts)
	return JoinPath("app-assets", id.String(), fmt.Sprintf("%s.%s", hash, options.Format))
}

func AchievementIcon(app, achievement snowflake.ID, hash string, opts ...CDNOption) string {
	if hash == "" {
		return ""
	}
	options := parseOptions(opts)
	return JoinPath("app-assets", app.String(), "achievements", achievement.String(), fmt.Sprintf("%s.%s", hash, options.Format))
}

func StorePageAsset(app, asset snowflake.ID, hash string, opts ...CDNOption) string {
	if hash == "" {
		return ""
	}
	options := parseOptions(opts)
	return JoinPath("app-assets", app.String(), "store", asset.String(), fmt.Sprintf("%s.%s", hash, options.Format))
}

func StickerPackBanner(app, asset snowflake.ID, hash string, opts ...CDNOption) string {
	if hash == "" {
		return ""
	}
	options := parseOptions(opts)
	return JoinPath("app-assets", app.String(), "store", asset.String(), fmt.Sprintf("%s.%s", hash, options.Format))
}

func TeamIcon(team snowflake.ID, hash string, opts ...CDNOption) string {
	if hash == "" {
		return ""
	}
	options := parseOptions(opts)
	return JoinPath("team-icons", team.String(), fmt.Sprintf("%s.%s", hash, options.Format))
}

func Sticker(sticker snowflake.ID, opts ...CDNOption) string {
	options := parseOptions(opts)
	return JoinPath("stickers", fmt.Sprintf("%d.%s", sticker, options.Format))
}

func RoleIcon(role snowflake.ID, hash string, opts ...CDNOption) string {
	if hash == "" {
		return ""
	}
	options := parseOptions(opts)
	return JoinPath("role-icons", role.String(), fmt.Sprintf("%s.%s", hash, options.Format))
}

func GuildScheduledEventCover(event snowflake.ID, hash string, opts ...CDNOption) string {
	if hash == "" {
		return ""
	}
	options := parseOptions(opts)
	return JoinPath("guild-events", event.String(), fmt.Sprintf("%s.%s", hash, options.Format))
}

func GuildMemberBanner(guild, user snowflake.ID, hash string, opts ...CDNOption) string {
	options := parseOptions(opts)
	if hash == "" {
		return ""
	}
	if strings.HasPrefix(hash, "a_") && options.Animated {
		options.Format = "gif"
	}
	return JoinPath("guilds", guild.String(), "users", user.String(), "banners", fmt.Sprintf("%s.%s", hash, options.Format))
}

func parseOptions(opts []CDNOption) (cdn CDNOptions) {
	cdn.Format = "png"
	cdn.Animated = true
	cdn.FallbackToDefault = true
	for _, opt := range opts {
		cdn = opt(cdn)
	}
	return
}

type CDNOptions struct {
	Format            string
	Animated          bool
	FallbackToDefault bool
}

type CDNOption func(opts CDNOptions) CDNOptions

func Format(f string) CDNOption {
	return func(opts CDNOptions) CDNOptions {
		opts.Format = f
		return opts
	}
}

func DoNotAnimate() CDNOption {
	return func(opts CDNOptions) CDNOptions {
		opts.Animated = false
		return opts
	}
}

func DoNotUseDefault() CDNOption {
	return func(opts CDNOptions) CDNOptions {
		opts.FallbackToDefault = false
		return opts
	}
}
