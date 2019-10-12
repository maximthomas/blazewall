package models

type Session struct {
	ID string `json:"id,omitempty"`

	UserID string `json:"userId,omitempty"`

	Realm string `json:"realm,omitempty"`

	Expired int64 `json:"expired,omitempty"`

	Properties map[string]string `json:"properties,omitempty"`
}