package api

import (
	"net/url"
	"strings"
	"time"

	"github.com/BOOMfinity/bfcord/discord"
	"github.com/BOOMfinity/bfcord/internal/httpc"
	"github.com/BOOMfinity/golog/v2"
	"github.com/andersfylling/snowflake/v5"
	"github.com/valyala/fasthttp"
)

var oauthClient = httpc.NewClient("", golog.New("oauth").Level(golog.LevelError))

const (
	oauthAuthorizeUrl = "https://discord.com/oauth2/authorize"
)

type OAuthScope = string

const (
	OAuthScopeIdentify OAuthScope = "identify"
	OAuthScopeGuilds   OAuthScope = "guilds"
)

type OAuthPrompt = string

const (
	OAuthPromptNone    OAuthPrompt = "none"
	OAuthPromptConsent OAuthPrompt = "consent"
)

type OAuthConfig struct {
	RedirectURL  string
	ClientID     snowflake.ID
	ClientSecret string
	Scopes       []OAuthScope
	Prompt       OAuthPrompt
}

func (cfg OAuthConfig) AuthorizeURL(state string) string {
	values := url.Values{}
	values.Set("client_id", cfg.ClientID.String())
	values.Set("response_type", "code")
	values.Set("state", state)
	if cfg.Prompt != "" {
		values.Set("prompt", cfg.Prompt)
	}
	values.Set("redirect_uri", cfg.RedirectURL)
	values.Set("scope", strings.Join(cfg.Scopes, "+"))
	x, _ := url.QueryUnescape(values.Encode())
	return oauthAuthorizeUrl + "?" + x
}

type OAuthToken struct {
	AccessToken  string    `json:"access_token"`
	TokenType    string    `json:"token_type"`
	ExpiresIn    uint      `json:"expires_in"`
	RefreshToken string    `json:"refresh_token"`
	Scope        string    `json:"scope"`
	CreatedAt    time.Time `json:"created_at"`
}

func (t OAuthToken) ExpiresAt() time.Time {
	return t.CreatedAt.Add(time.Duration(t.ExpiresIn) * time.Second)
}

func (t *OAuthToken) Refresh(cfg OAuthConfig) error {
	data, err := httpc.NewJSONRequest[OAuthToken](oauthClient, func(b httpc.RequestBuilder) error {
		b.NoAuth()
		b.Method(fasthttp.MethodPost)
		b.OnRequest(func(req *fasthttp.Request) error {
			req.Header.SetContentType("application/x-www-form-urlencoded")
			values := url.Values{}
			values.Set("grant_type", "refresh_token")
			values.Set("refresh_token", t.RefreshToken)
			values.Set("client_id", cfg.ClientID.String())
			values.Set("client_secret", cfg.ClientSecret)
			req.SetBodyString(values.Encode())
			return nil
		})
		return b.Execute("oauth2", "token")
	})
	if err != nil {
		return err
	}
	data.CreatedAt = time.Now()
	*t = data
	return nil
}

func (t *OAuthToken) FetchUser() (discord.User, error) {
	return httpc.NewJSONRequest[discord.User](oauthClient, func(b httpc.RequestBuilder) error {
		b.NoAuth()
		b.Header("Authorization", t.TokenType+" "+t.AccessToken)
		return b.Execute("users", "@me")
	})
}

func (t *OAuthToken) FetchGuilds() ([]OAuthGuild, error) {
	return httpc.NewJSONRequest[[]OAuthGuild](oauthClient, func(b httpc.RequestBuilder) error {
		b.NoAuth()
		b.Header("Authorization", t.TokenType+" "+t.AccessToken)
		return b.Execute("users", "@me", "guilds"+"?with_counts=true")
	})
}

func (cfg OAuthConfig) OAuthExchange(code string) (OAuthToken, error) {
	data, err := httpc.NewJSONRequest[OAuthToken](oauthClient, func(b httpc.RequestBuilder) error {
		b.NoAuth()
		b.Method(fasthttp.MethodPost)
		b.OnRequest(func(req *fasthttp.Request) error {
			req.Header.SetContentType("application/x-www-form-urlencoded")
			values := url.Values{}
			values.Set("grant_type", "authorization_code")
			values.Set("code", code)
			values.Set("redirect_uri", cfg.RedirectURL)
			values.Set("client_id", cfg.ClientID.String())
			values.Set("client_secret", cfg.ClientSecret)
			req.SetBodyString(values.Encode())
			return nil
		})
		return b.Execute("oauth2", "token")
	})
	if err != nil {
		return data, err
	}
	data.CreatedAt = time.Now()
	return data, nil
}

type OAuthGuild struct {
	ID                       snowflake.ID       `json:"id"`
	Name                     string             `json:"name"`
	Icon                     string             `json:"icon,omitempty"`
	Banner                   string             `json:"banner,omitempty"`
	Owner                    bool               `json:"owner"`
	Permissions              discord.Permission `json:"permissions"`
	Features                 []string           `json:"features"`
	ApproximateMemberCount   uint               `json:"approximate_member_count,omitempty"`
	ApproximatePresenceCount uint               `json:"approximate_presence_count,omitempty"`
}
