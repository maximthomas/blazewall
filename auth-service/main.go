package main

import (
	"flag"
	"os"

	"github.com/maximthomas/blazewall/auth-service/repo"

	"github.com/maximthomas/blazewall/auth-service/server"

	"github.com/maximthomas/blazewall/auth-service/config"
)

var yamlConfigFile = flag.String("yc", "./config/auth-config.yaml", "Yaml config file path")
var port = flag.String("p", "8080", "Gateway service port")

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	flag.Parse()

	configReader, err := os.Open(*yamlConfigFile)
	check(err)
	config.Init(configReader)
	repo.InitSessionRepo()
	router := server.GetRouter()
	router.Run(":" + *port)
}
