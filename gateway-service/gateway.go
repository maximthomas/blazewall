package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"path/filepath"
)

const sessionKey = "blazewall-session"

type Gateway struct {
	ProtectedSitesConfig map[string]ProtectedSiteConfig
	SessionRepository    SessionRepository
}

func NewGateway(protectedSitesConfig []ProtectedSiteConfig, sessionRepository SessionRepository) *Gateway {
	gateway := new(Gateway)
	gateway.ProtectedSitesConfig = make(map[string]ProtectedSiteConfig)
	gateway.SessionRepository = sessionRepository
	for _, protectedSiteConfig := range protectedSitesConfig {
		targetURL, err := url.Parse(protectedSiteConfig.TargetHost)
		if err != nil {
			log.Fatal(err)
			panic(err)
		}
		protectedSiteConfig = ProtectedSiteConfig{
			RequestHost:          protectedSiteConfig.RequestHost,
			TargetHost:           protectedSiteConfig.TargetHost,
			ProtectedPathsConfig: protectedSiteConfig.ProtectedPathsConfig,
			proxy:                httputil.NewSingleHostReverseProxy(targetURL),
		}
		gateway.ProtectedSitesConfig[protectedSiteConfig.RequestHost] = protectedSiteConfig
	}

	return gateway
}

func (g *Gateway) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path

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
			sessionID, _ := g.getSessionID(r)
			session, ok := g.SessionRepository.GetSession(sessionID)
			var sessionPtr *Session
			if ok {
				sessionPtr = &session
			}

			valid := protected.PolicyValidator.ValidatePolicy(r, sessionPtr)
			if valid {
				r.Header.Add("X-Forwarded-For", r.Host)
				protectedSiteConfig.proxy.ServeHTTP(w, r)
			} else {
				w.WriteHeader(http.StatusUnauthorized)
				fmt.Fprintf(w, "Unauthorized")
			}
			return
		}
	}
	protectedSiteConfig.proxy.ServeHTTP(w, r)
}

// ErrNoSessionID is returned from request
var ErrNoSessionID = errors.New("no session id presented in the request")

func (g *Gateway) getSessionID(r *http.Request) (string, error) {
	sessionCookie, err := r.Cookie(sessionKey)
	if err == nil {
		return sessionCookie.Value, nil
	}

	if len(r.Header[sessionKey]) > 0 {
		return r.Header[sessionKey][0], nil
	}
	return "", ErrNoSessionID
}
