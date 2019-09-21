package main

import (
	"net/http"
)

type PolicyValidator interface {
	ValidatePolicy(r *http.Request, s *Session) bool
}

type AllowedPolicyValidator struct{}

func (a AllowedPolicyValidator) ValidatePolicy(r *http.Request, s *Session) bool {
	return true
}

type DeniedPolicyValidator struct{}

func (d DeniedPolicyValidator) ValidatePolicy(r *http.Request, s *Session) bool {
	return false
}

type AuthenticatedUserPolicyValidator struct{}

func (a AuthenticatedUserPolicyValidator) ValidatePolicy(r *http.Request, s *Session) bool {
	if s != nil {
		return true
	}
	return false
}

type RealmsPolicyValidator struct {
	Realms []string
}

func (rp RealmsPolicyValidator) ValidatePolicy(r *http.Request, s *Session) bool {
	if s == nil {
		return false
	}

	if len(rp.Realms) > 0 && arrayContains(s.Realm, rp.Realms) {
		return true
	}
	return false
}

func arrayContains(val string, array []string) bool {
	for _, el := range array {
		if el == val {
			return true
		}
	}
	return false

}
