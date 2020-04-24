package models

//Session struct represents session object from session service
type Session struct {
	ID         string            `json:"id,omitempty"`
	UserID     string            `json:"userId,omitempty"`
	Realm      string            `json:"realm,omitempty"`
	Properties map[string]string `json:"properties,omitempty"`
}
