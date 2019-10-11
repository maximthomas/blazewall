package main

import (
	"flag"
	"os"

	"github.com/maximthomas/blazewall/user-service/server"

	"github.com/maximthomas/blazewall/user-service/config"
)

var port = flag.String("p", "8080", "User service port")
var yamlConfigFile = flag.String("yc", "./config/user-service.yaml", "Yaml config file path")

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
	router := server.SetupRouter()
	router.Run(":" + *port)
}
