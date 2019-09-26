package main

import (
	"flag"
	"log"
	"net/http"
	"os"
)

var yamlConfigFile = flag.String("yc", "", "Yaml config file path")
var port = flag.String("p", "8080", "Gateway service port")

func check(err error) {
	if err != nil {
		panic(err)
	}
}

var sessionRepo *SessionRepository

func main() {

	flag.Parse()

	sessionRepo := &InMemorySessionRepository{}

	var sitesConfig []ProtectedSiteConfig
	if *yamlConfigFile != "" {
		configReader, err := os.Open(*yamlConfigFile)
		check(err)

		sitesConfig, err = NewProtectedSitesConfigYaml(configReader)
		check(err)
	}

	log.Printf("sites config: %#v", sitesConfig)

	gateway := NewGateway(sitesConfig, sessionRepo)

	http.ListenAndServe(":"+*port, gateway)

}
