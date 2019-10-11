package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"path/filepath"

	"github.com/maximthomas/blazewall/gateway-service/models"

	"github.com/maximthomas/blazewall/gateway-service/repo"

	"github.com/maximthomas/blazewall/gateway-service/config"
)

//Gateway main structure, holds protected sites config and session repository
type Gateway struct {
	ProtectedSitesConfig map[string]config.ProtectedSiteConfig
	SessionRepository    repo.SessionRepository
}

//NewGateway constructs new Gateway Instance
func NewGateway(protectedSitesConfig []config.ProtectedSiteConfig, sessionRepository repo.SessionRepository) *Gateway {
	gateway := new(Gateway)
	gateway.ProtectedSitesConfig = make(map[string]config.ProtectedSiteConfig)
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
			var sessionPtr *models.Session
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

				protectedSiteConfig.Proxy.ServeHTTP(w, r)
			} else {
				http.Redirect(w, r, protected.AuthURL, http.StatusFound)
			}
			return
		}
	}
	protectedSiteConfig.Proxy.ServeHTTP(w, r)
}

// ErrNoSessionID is returned from request
var ErrNoSessionID = errors.New("no session id presented in the request")

func (g *Gateway) getSessionID(r *http.Request) (string, error) {

	gc := config.GetConfig()

	sessionCookie, err := r.Cookie(gc.SessionID)
	if err == nil {
		return sessionCookie.Value, nil
	}

	if len(r.Header[gc.SessionID]) > 0 {
		return r.Header[gc.SessionID][0], nil
	}
	return "", ErrNoSessionID
}
