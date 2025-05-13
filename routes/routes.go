package routes

import (
	"product-service/cmd/product/handler"
	"product-service/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine, productHandler handler.ProductHandler) {
	router.Use(middleware.RequestLogger())
	router.POST("/v1/product", productHandler.ProductManagement)
	router.POST("/v1/product_category", productHandler.ProductCategoryManagement)

	router.GET("/v1/product/:id", productHandler.GetProductInfo)
	router.GET("/v1/product_category/:id", productHandler.GetProductCategoryInfo)
}
