package interactions

import (
	"fmt"
	"math"

	"github.com/BOOMfinity/bfcord/api"
	"github.com/BOOMfinity/bfcord/discord/components"
	"github.com/segmentio/encoding/json"
	"github.com/valyala/fasthttp"
)

type ModalBuilder struct {
	i     *Interaction
	id    string
	title string
	items []*TextFieldData
}

func (x *ModalBuilder) AddTextField(id string, label string) *TextFieldBuilder {
	field := NewTextFieldBuilder(id, label)
	x.items = append(x.items, field.data)
	return field
}

func (x *ModalBuilder) Execute() error {
	req := http.New(false)
	req.SetRequestURI(fmt.Sprintf(api.FullApiUrl+"/interactions/%v/%v/callback", x.i.ID.String(), x.i.Token))
	req.Header.SetMethod(fasthttp.MethodPost)
	data := make(map[string]any)
	index := 0
	rows := make([]map[string]any, int(math.Ceil(float64(len(x.items))/5)))
	row := make([]*TextFieldData, 0, 5)
	for i := 0; i < len(x.items); i++ {
		if rows[index] == nil {
			rows[index] = make(map[string]any, 2)
		}
		rows[index]["type"] = components.TypeActionRow
		field := x.items[i]
		if len(row) < 4 {
			row = append(row, field)
			continue
		}
		rows[index]["components"] = row
		row = row[:0]
		index++
		i--
	}
	if len(row) != 0 {
		rows[index]["components"] = row
	}
	modalData := make(map[string]any)
	data["type"] = ModalCallback
	modalData["title"] = x.title
	modalData["custom_id"] = x.id
	modalData["components"] = rows
	data["data"] = modalData
	raw, err := json.Marshal(data)
	if err != nil {
		return err
	}
	req.SetBody(raw)
	return http.DoNoResp(req)
}

type TextFieldStyle uint8

const (
	SingleLineStyle TextFieldStyle = iota + 1
	MultiLineStyle
)

type TextFieldBuilder struct {
	data *TextFieldData
}

func (x *TextFieldBuilder) NotRequired() *TextFieldBuilder {
	x.data.Required = false
	return x
}

func (x *TextFieldBuilder) MultiLine() *TextFieldBuilder {
	x.data.Style = 2
	return x
}

func NewModalBuilder(i *Interaction, id string, title string) *ModalBuilder {
	return &ModalBuilder{
		i:     i,
		id:    id,
		title: title,
	}
}

func NewTextFieldBuilder(id string, label string) *TextFieldBuilder {
	return &TextFieldBuilder{
		data: &TextFieldData{
			Type:     4,
			ID:       id,
			Style:    1,
			Label:    label,
			Required: true,
			Max:      4000,
		},
	}
}

type TextFieldData struct {
	Value       *string        `json:"value,omitempty"`
	Placeholder *string        `json:"placeholder,omitempty"`
	Label       string         `json:"label,omitempty"`
	ID          string         `json:"custom_id,omitempty"`
	Min         uint16         `json:"min_length,omitempty"`
	Max         uint16         `json:"max_length,omitempty"`
	Type        uint8          `json:"type"`
	Required    bool           `json:"required,omitempty"`
	Style       TextFieldStyle `json:"style,omitempty"`
}
