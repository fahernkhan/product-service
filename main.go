package main

import (
	"product-service/cmd/product/handler"
	"product-service/cmd/product/repository"
	"product-service/cmd/product/resource"
	"product-service/cmd/product/service"
	"product-service/cmd/product/usecase"
	"product-service/config"
	"product-service/infrastructure/log"
	"product-service/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.LoadConfig()
	redis := resource.InitRedis(&cfg)
	db := resource.InitDB(&cfg)

	log.SetupLogger()
	// routes
	// prepare each layer

	productRepository := repository.NewProductRepository(db, redis)
	productService := service.NewProductService(*productRepository)
	productUsecase := usecase.NewProductUsecase(*productService)
	productHandler := handler.NewProductHandler(*productUsecase)

	port := cfg.App.Port
	router := gin.Default()
	routes.SetupRoutes(router, *productHandler)
	router.Run(":" + port)

	log.Logger.Printf("Server running on port: %s", port)
}
