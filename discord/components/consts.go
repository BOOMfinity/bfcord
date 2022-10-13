package components

type Type uint8

const (
	TypeActionRow Type = iota + 1
	TypeButton
	TypeSelectMenu
)

type ButtonStyle uint8

const (
	ButtonStylePrimary ButtonStyle = iota + 1
	ButtonStyleSecondary
	ButtonStyleSuccess
	ButtonStyleDanger
	ButtonStyleLink
)
