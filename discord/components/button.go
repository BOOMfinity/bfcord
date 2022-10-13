package components

type ButtonBuilder struct {
	baseBuilder
	//id       string
	label string
	url   string
	emoji PartialEmoji
	//disabled bool
	style ButtonStyle
}

func (v *ButtonBuilder) SetCustomID(id string) *ButtonBuilder {
	v.id = id
	return v
}

func (v *ButtonBuilder) SetDisabled(disabled bool) *ButtonBuilder {
	v.disabled = disabled
	return v
}

func (v *ButtonBuilder) SetLabel(label string) *ButtonBuilder {
	v.label = label
	return v
}

func (v *ButtonBuilder) SetUrl(url string) *ButtonBuilder {
	v.url = url
	return v
}

func (v *ButtonBuilder) SetEmoji(emoji PartialEmoji) *ButtonBuilder {
	v.emoji = emoji
	return v
}

func (v *ButtonBuilder) ToComponent() (x Component) {
	x.Type = TypeButton
	if v.style != ButtonStyleLink {
		x.CustomID = v.id
	} else {
		x.Url = v.url
	}
	x.Label = v.label
	x.Disabled = v.disabled
	if !v.emoji.ID.IsZero() && v.emoji.Name == "" {
		x.Emoji = &v.emoji
	}
	x.Style = v.style
	return
}

func NewButton(style ButtonStyle) *ButtonBuilder {
	return &ButtonBuilder{
		baseBuilder: baseBuilder{t: TypeButton},
		style:       style,
	}
}
