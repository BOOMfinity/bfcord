package slash

type Choice struct {
	Value             any               `json:"value,omitempty"`
	NameLocalizations map[string]string `json:"name_localizations,omitempty"`
	Name              string            `json:"name"`
}
