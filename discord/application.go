package discord

import "github.com/andersfylling/snowflake/v5"

type Application struct {
	ID                             snowflake.ID `json:"id,omitempty"`
	Name                           string       `json:"name,omitempty"`
	Icon                           string       `json:"icon,omitempty"`
	Description                    string       `json:"description,omitempty"`
	RPCOrigins                     []string     `json:"rpc_origins,omitempty"`
	BotPublic                      bool         `json:"bot_public,omitempty"`
	BotRequireCodeGrant            bool         `json:"bot_require_code_grant,omitempty"`
	Bot                            User         `json:"bot"`
	TermsOfServiceURL              string       `json:"terms_of_service_url,omitempty"`
	PrivacyPolicyURL               string       `json:"privacy_policy_url,omitempty"`
	Owner                          User         `json:"owner"`
	VerifyKey                      string       `json:"verify_key,omitempty"`
	GuildID                        snowflake.ID `json:"guild_id,omitempty"`
	Guild                          Guild        `json:"guild"`
	PrimarySkuID                   snowflake.ID `json:"primary_sku_id,omitempty"`
	Slug                           string       `json:"slug,omitempty"`
	CoverImage                     string       `json:"cover_image,omitempty"`
	Flags                          uint         `json:"flags,omitempty"`
	ApproximateGuildCount          uint         `json:"approximate_guild_count,omitempty"`
	RedirectUris                   []string     `json:"redirect_uris,omitempty"`
	InteractionsEndpointURL        string       `json:"interactions_endpoint_url,omitempty"`
	RoleConnectionsVerificationURL string       `json:"role_connections_verification_url,omitempty"`
	Tags                           []string     `json:"tags,omitempty"`
	CustomInstallURL               string       `json:"custom_install_url,omitempty"`
}
