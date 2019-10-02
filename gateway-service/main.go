package main

import (
	"flag"
	"log"
	"net/http"
	"os"
)

var yamlConfigFile = flag.String("yc", "", "Yaml config file path")
var port = flag.String("p", "8080", "Gateway service port")
var sessionServiceEndpoint = flag.String("sess", "http://session-service:8080/session-service/v1/sessions", "Session service endpoint")
var authSessionID = flag.String("sID", "BlazewallSession", "Session service cookie name")

func check(err error) {
	if err != nil {
		panic(err)
	}
}

var sessionRepo *SessionRepository

func main() {

	flag.Parse()

	sr := RestSessionRepository{endpoint: *sessionServiceEndpoint}

	var sitesConfig []ProtectedSiteConfig
	if *yamlConfigFile != "" {
		configReader, err := os.Open(*yamlConfigFile)
		check(err)

		sitesConfig, err = NewProtectedSitesConfigYaml(configReader)
		check(err)
	}

	log.Printf("sites config: %#v", sitesConfig)

	gateway := NewGateway(sitesConfig, &sr)

	http.ListenAndServe(":"+*port, gateway)

}
