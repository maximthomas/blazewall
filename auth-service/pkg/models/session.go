package models

//Session struct represents session object from session service
type Session struct {
	ID         string            `json:"id,omitempty"`
	Properties map[string]string `json:"properties,omitempty"`
}
