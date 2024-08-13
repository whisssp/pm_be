package main

import (
	"pm/infrastructure/config"
	"pm/infrastructure/persistences/base"
	"pm/infrastructure/persistences/base/logger"
	"pm/infrastructure/routes"
)

// @title			Product Management API
// @version		1.0
// @description	This is a sample server celler server.
func main() {

	persistence := base.InitPersistence(config.Configs)
	logger.NewLogger()
	server := routes.InitServer(config.Configs.Server.Port, persistence)
	server.Run()
}
