package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/http/httputil"
	"net/url"
	"testing"
)

const proxyHost = "gateway-service"

const wantResponse = "valid response"

func TestPassThrough(t *testing.T) {

	protectedServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, wantResponse)
	}))

	protectedURL, _ := url.Parse(protectedServer.URL)

	gateway := NewGateway([]ProtectedSiteConfig{
		{
			RequestHost: proxyHost,
			TargetHost:  protectedServer.URL,
			ProtectedPathsConfig: []ProtectedPathConfig{
				{
					URLPattern:      "/",
					PolicyValidator: AllowedPolicyValidator{},
				},
				{
					URLPattern:      "/protected",
					PolicyValidator: DeniedPolicyValidator{},
				},
			},
			proxy: httputil.NewSingleHostReverseProxy(protectedURL),
		},
	}, NewInMemorySessionRepository(nil))

	t.Run("test proxy ok", func(t *testing.T) {
		allowedRequest := newGetRequest("/")
		response := httptest.NewRecorder()

		gateway.ServeHTTP(response, allowedRequest)

		assertStatus(t, response.Code, http.StatusOK)
		assertBody(t, response.Body.String(), wantResponse)
	})

	t.Run("test proxy 401", func(t *testing.T) {
		allowedRequest := newGetRequest("/protected")
		response := httptest.NewRecorder()

		gateway.ServeHTTP(response, allowedRequest)

		assertStatus(t, response.Code, http.StatusUnauthorized)
		assertBody(t, response.Body.String(), "Unauthorized")
	})
}

func BenchmarkTestPass(b *testing.B) {

	protectedServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, wantResponse)
	}))

	gateway := NewGateway([]ProtectedSiteConfig{
		{
			RequestHost: proxyHost,
			TargetHost:  protectedServer.URL,
			ProtectedPathsConfig: []ProtectedPathConfig{
				{
					URLPattern:      "/",
					PolicyValidator: AllowedPolicyValidator{},
				},
				{
					URLPattern:      "/protected",
					PolicyValidator: DeniedPolicyValidator{},
				},
			},
		},
	},
		NewInMemorySessionRepository(nil))

	for i := 0; i < b.N; i++ {
		b.Run("test proxy ok", func(b *testing.B) {
			deniedRequest := newGetRequest("/protected")
			response := httptest.NewRecorder()
			gateway.ServeHTTP(response, deniedRequest)
		})
	}
}

func assertStatus(t *testing.T, got, want int) {
	t.Helper()
	if got != want {
		t.Errorf("did not get correct status, got %d, want %d", got, want)
	}
}

func assertBody(t *testing.T, got, want string) {
	t.Helper()
	if got != want {
		t.Errorf("did not get correct body, got %s, want %s", got, want)
	}

}
func newGetRequest(path string) *http.Request {
	req, _ := http.NewRequest(http.MethodGet, path, nil)
	req.Host = proxyHost
	return req
}
