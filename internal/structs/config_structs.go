package structs

type ElementType int

const (
	Text ElementType = iota + 1
	Select
	Radio
	Number
	Password
	Disabled
	Info
	FileUpload
	TextArea
)

type ModuleConfig struct {
	Fields []Element `json:"fields"`
}

type Element struct {
	Label            string      `json:"label"`
	Type             ElementType `json:"type"`
	ExpectedJsonName string      `json:"expected_json_name"`
	Rationale        string      `json:"rationale"`
	Value            string      `json:"value"`
	PossibleValues   []string    `json:"possible_values"`
	Required         bool        `json:"required"`
}
