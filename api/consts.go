package api

const (
	ApiUrl     = "https://discord.com/api"
	ApiVersion = "v10"
	FullApiUrl = ApiUrl + "/" + ApiVersion
)

type emptyOptions[V any] struct {
	data   V
	reason string
}

func (e *emptyOptions[V]) NoCache() V {
	return e.data
}

func (e *emptyOptions[V]) NoAPI() V {
	return e.data
}

func (e *emptyOptions[V]) Reason(str string) V {
	e.reason = str
	return e.data
}
