package components

import "github.com/segmentio/encoding/json"

type baseBuilder struct {
	id       string
	disabled bool
	t        Type
}

func (v *baseBuilder) Type() Type {
	return v.t
}

func (v *baseBuilder) SetCustomID(id string) *baseBuilder {
	v.id = id
	return v
}

func (v *baseBuilder) SetDisabled(disabled bool) *baseBuilder {
	v.disabled = disabled
	return v
}

type ActionRowBuilder struct {
	baseBuilder
	components []Component
}

func (a *ActionRowBuilder) Size() int {
	return len(a.components)
}

func (a *ActionRowBuilder) Add(components ...Component) *ActionRowBuilder {
	a.components = components
	return a
}

func (a *ActionRowBuilder) ToComponent() (x Component) {
	x.Type = TypeActionRow
	x.Components = a.components
	return
}

func NewActionRow() *ActionRowBuilder {
	return &ActionRowBuilder{}
}

type ActionRowItem interface {
	ToComponent() Component
	Type() Type
}

type components []Component

type List [][]ActionRowItem

var emptyArray = []byte("[]")
var jsonNull = []byte("null")

func (v components) MarshalJSON() ([]byte, error) {
	if v == nil {
		return jsonNull, nil
	}
	if len(v) == 0 {
		return emptyArray, nil
	}
	type ComponentListCopy components
	return json.Marshal(ComponentListCopy(v))
}

func (v *components) UnmarshalJSON(bytes []byte) error {
	type ComponentListCopy components
	var list ComponentListCopy
	err := json.Unmarshal(bytes, &list)
	if err != nil {
		return err
	}
	*v = components(list)
	return nil
}

func (rows List) Build() (components []Component) {
	for i := range rows {
		row := Component{Type: TypeActionRow}
		for i2 := range rows[i] {
			row.Components = append(row.Components, rows[i][i2].ToComponent())
		}
		components = append(components, row)
	}
	return
}
