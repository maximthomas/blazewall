package main

import (
	"flag"
	"log"
	"net/http"
	"os"

	"github.com/maximthomas/blazewall/gateway-service/server"

	"github.com/maximthomas/blazewall/gateway-service/config"
	"github.com/maximthomas/blazewall/gateway-service/repo"
)

var yamlConfigFile = flag.String("yc", "./config/gateway-config.yaml", "Yaml config file path")
var port = flag.String("p", "8080", "Gateway service port")

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	log.SetFlags(log.LstdFlags | log.Llongfile)
	flag.Parse()
	configReader, err := os.Open(*yamlConfigFile)
	check(err)
	config.Init(configReader)
	gc := config.GetConfig()
	repo.Init()
	sr := repo.GetSessionRepository()

	log.Printf("gateway config: %#v", gc)

	gateway := server.NewGateway(gc.ProtectedSitesConfig, sr)

	http.ListenAndServe(":"+*port, gateway)

}
