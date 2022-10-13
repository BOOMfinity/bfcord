package components

type SelectMenuBuilder struct {
	baseBuilder
	//id          string
	placeholder string
	options     SelectMenuOptions
	//disabled    bool
	minValues uint8
	maxValues uint8
}

func (v *SelectMenuBuilder) SetCustomID(id string) *SelectMenuBuilder {
	v.id = id
	return v
}

func (v *SelectMenuBuilder) SetDisabled(disabled bool) *SelectMenuBuilder {
	v.disabled = disabled
	return v
}

func (v *SelectMenuBuilder) ToComponent() (x Component) {
	x.CustomID = v.id
	x.Type = TypeSelectMenu
	x.Disabled = v.disabled
	x.Options = v.options
	x.MaxValues = v.maxValues
	x.MinValues = v.minValues
	return
}

func (v *SelectMenuBuilder) SetPlaceholder(text string) *SelectMenuBuilder {
	v.placeholder = text
	return v
}

func (v *SelectMenuBuilder) SetMinValues(min uint8) *SelectMenuBuilder {
	v.minValues = min
	return v
}

func (v *SelectMenuBuilder) SetMaxValues(max uint8) *SelectMenuBuilder {
	v.maxValues = max
	return v
}

func (v *SelectMenuBuilder) SetOptions(options SelectMenuOptions) *SelectMenuBuilder {
	v.options = options
	return v
}

func NewSelectMenu() *SelectMenuBuilder {
	return &SelectMenuBuilder{
		baseBuilder: baseBuilder{t: TypeSelectMenu},
		minValues:   1,
		maxValues:   1,
	}
}
