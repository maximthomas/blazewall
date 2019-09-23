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

	var sitesConfigPtr *[]ProtectedSiteConfig
	if *yamlConfigFile != "" {
		configReader, err := os.Open(*yamlConfigFile)
		check(err)

		sitesConfigPtr, err = NewProtectedSitesConfigYaml(configReader)
		check(err)
	}

	if sitesConfigPtr == nil {
		log.Fatal("error reading sites config")
	}

	sitesConfig := *sitesConfigPtr

	log.Printf("sites config: %#v", sitesConfig)

	gateway := NewGateway(sitesConfig, sessionRepo)

	http.ListenAndServe(":"+*port, gateway)

}
