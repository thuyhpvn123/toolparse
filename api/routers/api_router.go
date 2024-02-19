package routers

import (
	"fmt"
	"net/http"
	"github.com/gin-gonic/gin"
	"github.com/meta-node-blockchain/meta-node/cmd/chiabai/api/controllers"
	c_config "github.com/meta-node-blockchain/meta-node/cmd/client/pkg/config"
	cc_config "github.com/meta-node-blockchain/meta-node/cmd/chiabai/config"
	"github.com/meta-node-blockchain/meta-node/pkg/logger"
)

// SetupRouter sets up the API routes and returns the Gin router.
func InitRouter() *gin.Engine {
	server := controllers.Server{}
	config, err := c_config.LoadConfig(cc_config.CONFIG_FILE_PATH)
	if err != nil {
		logger.Error(fmt.Sprintf("error when loading config %v", err))
		panic(fmt.Sprintf("error when loading config %v", err))
	}
	cConfig := config.(*c_config.ClientConfig)

	server.Init(cConfig)
	router := gin.Default()
	// Set a lower memory limit for multipart forms (default is 32 MiB)
	router.MaxMultipartMemory = 8 << 20 // 8 MiB

	v1 := router.Group("/api/v1")
	{
		v1.StaticFS("", http.Dir("frontend/public"))
	}
	router.GET("/ws", func(c *gin.Context) {
		server.WebsocketHandler(c.Writer, c.Request)
	})

	return router
}

