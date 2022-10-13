package discord

type ReactionStore []MessageReaction

func (v ReactionStore) Get(emoji string) (reaction MessageReaction, ok bool) {
	for i := range v {
		e := v[i].Emoji
		if e.ToString() == emoji || e.Name == emoji || e.ID.String() == emoji {
			return v[i], true
		}
	}
	return
}

type Slice[V any] []V

func (s Slice[V]) Contains(fn func(item V) bool) bool {
	for i := range s {
		if fn(s[i]) {
			return true
		}
	}
	return false
}

func (s Slice[V]) Find(fn func(item V) bool) (item V, found bool) {
	for i := range s {
		if fn(s[i]) {
			return s[i], true
		}
	}
	return
}

func (s Slice[V]) Filter(fn func(item V) bool) (items Slice[V]) {
	for i := range s {
		if fn(s[i]) {
			items = append(items, s[i])
		}
	}
	return
}
