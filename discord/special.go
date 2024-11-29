package discord

type MemberWithUser struct {
	Member
	User User `json:"user,omitempty"`
}
