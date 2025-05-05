package main

import (
	"product-service/cmd/product/resource"
	"product-service/config"
	"product-service/infrastructure/log"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.LoadConfig()
	resource.InitRedis(&cfg)
	resource.InitDB(&cfg)

	log.SetupLogger()
	// routes
	// prepare each layer

	port := cfg.App.Port
	router := gin.Default()

	router.Run(":" + port)
	log.Logger.Printf("Server running on port: %s", port)
}
