package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"path/filepath"
)

//Gateway main structure, holds protected sites config and session repository
type Gateway struct {
	ProtectedSitesConfig map[string]ProtectedSiteConfig
	SessionRepository    SessionRepository
}

//NewGateway constructs new Gateway Instance
func NewGateway(protectedSitesConfig []ProtectedSiteConfig, sessionRepository SessionRepository) *Gateway {
	gateway := new(Gateway)
	gateway.ProtectedSitesConfig = make(map[string]ProtectedSiteConfig)
	gateway.SessionRepository = sessionRepository
	for _, protectedSiteConfig := range protectedSitesConfig {
		gateway.ProtectedSitesConfig[protectedSiteConfig.RequestHost] = protectedSiteConfig
	}
	return gateway
}

func (g *Gateway) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path

	dump, err := httputil.DumpRequest(r, false)
	if err != nil {
		log.Fatal(err)
		return
	}
	log.Printf("%s", dump)
	protectedSiteConfig, exists := g.ProtectedSitesConfig[r.Host]
	if !exists {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Invalid host")
		return
	}

	for _, protected := range protectedSiteConfig.ProtectedPathsConfig {
		match, err := filepath.Match(protected.URLPattern, path)
		if err != nil {
			log.Fatal(err)
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, "Unauthorized")
			return
		}
		if match {
			sessionID, err := g.getSessionID(r)
			var sessionPtr *Session
			if err == nil {
				session, ok := g.SessionRepository.GetSession(sessionID)
				if ok {
					sessionPtr = &session
				}
			}

			valid := protected.PolicyValidator.ValidatePolicy(r, sessionPtr)
			if valid {
				r.Header.Add("X-Forwarded-For", r.Host)
				if sessionPtr != nil {
					sessionJSONBytes, err := json.Marshal(*sessionPtr)
					if err != nil {
						log.Fatalf("error occurred %v", err)
						panic(err)
					}
					r.Header.Add("X-Blazewall-Session", string(sessionJSONBytes))
				}

				protectedSiteConfig.proxy.ServeHTTP(w, r)
			} else {
				http.Redirect(w, r, protected.AuthURL, http.StatusFound)
			}
			return
		}
	}
	protectedSiteConfig.proxy.ServeHTTP(w, r)
}

// ErrNoSessionID is returned from request
var ErrNoSessionID = errors.New("no session id presented in the request")

func (g *Gateway) getSessionID(r *http.Request) (string, error) {
	sessionCookie, err := r.Cookie(*authSessionID)
	if err == nil {
		return sessionCookie.Value, nil
	}

	if len(r.Header[*authSessionID]) > 0 {
		return r.Header[*authSessionID][0], nil
	}
	return "", ErrNoSessionID
}
