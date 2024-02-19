package main

import (
	// "fmt"

	"github.com/meta-node-blockchain/meta-node/cmd/chiabai/api/routers"

	// c_config "github.com/meta-node-blockchain/meta-node/cmd/chiabai/config"
	// c_config "github.com/meta-node-blockchain/meta-node/cmd/client/pkg/config"
	// cc_config "github.com/meta-node-blockchain/meta-node/cmd/chiabai/config"

	"github.com/meta-node-blockchain/meta-node/cmd/chiabai/database"
	// "github.com/meta-node-blockchain/meta-node/pkg/logger"
)

func main() {
	// load config
	config, err := cc_config.LoadConfig(cc_config.CONFIG_FILE_PATH)
	if err != nil {
		logger.Error(fmt.Sprintf("error when loading config %v", err))
		panic(fmt.Sprintf("error when loading config %v", err))
	}
	cConfig := config.(*c_config.ClientConfig)
	// Initialize the database connection
	database.InitDatabase()
	// Code to initialize the database
	// Initialize the Gin router
	
	router := routers.InitRouter()
	// Run the server
	if err := router.Run(":2000"); err != nil {
		panic(err)
	}

}

