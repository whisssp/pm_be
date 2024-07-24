package main

import (
	"fmt"
	"pm/infrastructure/config"
	"pm/infrastructure/persistences/base"
	"pm/infrastructure/routes"
)

//	@title			Product Management API
//	@version		1.0
//	@description	This is a sample server celler server.
func main() {
	appConfig, err := config.LoadConfig()
	if err != nil {
		fmt.Println("error loading config")
	}

	persistence := base.InitPersistence(appConfig)

	server := routes.InitServer(appConfig.Server.Port, persistence)
	server.Run()
}