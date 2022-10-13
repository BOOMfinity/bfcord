package interactions

import (
	"github.com/BOOMfinity/bfcord/discord/components"
)

type ComponentList []components.Component

func (x ComponentList) find(slice []components.Component, id string) (c components.Component, found bool) {
	for i := range slice {
		item := slice[i]
		if item.CustomID == id {
			return item, true
		}
		if len(item.Components) > 0 {
			return x.find(item.Components, id)
		}
	}
	return
}

func (x ComponentList) Get(id string) (c components.Component, found bool) {
	return x.find(x, id)
}

func (x ComponentList) Has(id string) bool {
	if _, found := x.Get(id); found {
		return true
	}
	return false
}
