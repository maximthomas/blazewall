package models

type Callback struct {
	Name       string            `json:"name,omitempty"`
	Type       string            `json:"type"`
	Value      string            `json:"value"`
	Validation string            `json:"validation,omitempty"`
	Required   bool              `json:"required,omitempty"`
	Properties map[string]string `json:"properties,omitempty"`
	Error      string            `json:"error,omitempty"`
}

type CallbackRequest struct {
	Module    string     `json:"module,omitempty"`
	Callbacks []Callback `json:"callbacks,omitempty"`
}
