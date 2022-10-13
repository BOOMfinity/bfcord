package interactions

import (
	"errors"
	"github.com/BOOMfinity/bfcord/discord"
	"github.com/BOOMfinity/bfcord/slash"
	"github.com/andersfylling/snowflake/v5"
	"github.com/segmentio/encoding/json"
)

var (
	ValueParseFailedErr = errors.New("value parse failed")
)

type Option struct {
	Value   interface{}      `json:"value,omitempty"`
	Name    string           `json:"name"`
	Options OptionList       `json:"options,omitempty"`
	Type    slash.OptionType `json:"type"`
	Focused bool             `json:"focused,omitempty"`
}

func (o *Option) Default(x interface{}) *Option {
	if o.Value == nil {
		o.Value = x
	}
	return o
}

func (o *Option) Nil() bool {
	if o.Value == nil {
		return true
	}
	return false
}

func (o *Option) Snowflake() (snowflake.ID, error) {
	if x, ok := o.Value.(string); ok {
		return snowflake.ParseSnowflakeString(x), nil
	}
	return 0, ValueParseFailedErr
}

func (o *Option) Attachment() (att discord.Attachment, err error) {
	data, err := json.Marshal(o.Value)
	if err != nil {
		return
	}
	err = json.Unmarshal(data, &att)
	return
}

func (o *Option) Int() (int, error) {
	float, err := o.Float()
	if err != nil {
		return 0, err
	}
	return int(float), nil
}

func (o *Option) String() (string, error) {
	if x, ok := o.Value.(string); ok {
		return x, nil
	}
	return "", ValueParseFailedErr
}

func (o *Option) Float() (float64, error) {
	if x, ok := o.Value.(float64); ok {
		return x, nil
	}
	return 0, ValueParseFailedErr
}

func (o *Option) Bool() (bool, error) {
	if x, ok := o.Value.(bool); ok {
		return x, nil
	}
	return false, ValueParseFailedErr
}

type OptionList []*Option

func (x OptionList) Get(name string) *Option {
	if index := x.findIndex(name); index != -1 {
		return x[index]
	}
	return &Option{}
}

func (x OptionList) findIndex(name string) int {
	for i := range x {
		if x[i].Name == name {
			return i
		}
	}
	return -1
}

func (x OptionList) Has(name string) bool {
	if x.findIndex(name) != -1 {
		return true
	}
	return false
}
