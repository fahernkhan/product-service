package handler

import (
	// golang package
	"fmt"
	"net/http"
	"product-service/cmd/product/usecase"
	"product-service/infrastructure/log"
	"product-service/models"
	"strconv"

	// external package
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type ProductHandler struct {
	ProductUsecase usecase.ProductUsecase
}

func NewProductHandler(orderUsecase usecase.ProductUsecase) *ProductHandler {
	return &ProductHandler{
		ProductUsecase: orderUsecase,
	}
}

// GetProductInfo get product info by given c pointer of gin.Context.
func (h *ProductHandler) GetProductInfo(c *gin.Context) {
	productIDstr := c.Param("id")

	productID, err := strconv.ParseInt(productIDstr, 10, 64)
	if err != nil {
		log.Logger.WithFields(logrus.Fields{
			"productID": productIDstr,
		}).Errorf("strconv.ParseInt got error %v", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error_message": "Invalid Product ID",
		})

		return
	}

	product, err := h.ProductUsecase.GetProductByID(c.Request.Context(), productID)
	if err != nil {
		log.Logger.WithFields(logrus.Fields{
			"productID": productID,
		}).Errorf("h.ProductUsecase.GetProductByID() got error %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error_message": err,
		})

		return
	}

	if product.ID == 0 {
		log.Logger.WithFields(logrus.Fields{
			"productID": productID,
		}).Info("Product ID not found")
		c.JSON(http.StatusBadRequest, gin.H{
			"error_message": "Product Not Exists",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"product": product,
	})
}

// GetProductCategoryInfo get product category info by given c pointer of gin.Context.
func (h *ProductHandler) GetProductCategoryInfo(c *gin.Context) {
	productCategoryIDstr := c.Param("id")

	productCategoryID, err := strconv.Atoi(productCategoryIDstr)
	if err != nil {
		log.Logger.WithFields(logrus.Fields{
			"productID": productCategoryIDstr,
		}).Errorf("strconv.Atoi got error %v", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error_message": "Invalid Product ID",
		})

		return
	}

	productCategory, err := h.ProductUsecase.GetProductCategoryByID(c.Request.Context(), productCategoryID)
	if err != nil {
		log.Logger.WithFields(logrus.Fields{
			"productCategoryID": productCategoryID,
		}).Errorf("h.ProductUsecase.GetProductCategoryByID() got error %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error_message": err,
		})

		return
	}

	if productCategory.ID == 0 {
		log.Logger.WithFields(logrus.Fields{
			"productCategoryID": productCategoryID,
		}).Info("Product Category Not Found")
		c.JSON(http.StatusBadRequest, gin.H{
			"error_message": "Product Category Not Exists",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"productCategory": productCategory,
	})
}

// ProductManagement product management by given c pointer of gin.Context.
func (h *ProductHandler) ProductManagement(c *gin.Context) {
	var param models.ProductManagementParameter
	if err := c.ShouldBindJSON(&param); err != nil {
		log.Logger.Error(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{
			"error_message": "Invalid Input",
		})

		return
	}

	// validateProductManagementParameter
	if param.Action == "" {
		log.Logger.Error("missing parameter action")
		c.JSON(http.StatusBadRequest, gin.H{
			"error_message": "Missing required parameter",
		})

		return
	}

	switch param.Action {
	case "add":
		if param.ID != 0 {
			log.Logger.WithFields(logrus.Fields{
				"param": param,
			}).Error("invalid request - product id is not empty")
			c.JSON(http.StatusBadRequest, gin.H{
				"error_message": "Invalid Request",
			})

			return
		}

		productID, err := h.ProductUsecase.CreateNewProduct(c.Request.Context(), &param.Product)
		if err != nil {
			log.Logger.WithFields(logrus.Fields{
				"param": param,
			}).Errorf("h.ProductUsecase.CreateNewProduct() got error %v", err)

			c.JSON(http.StatusInternalServerError, gin.H{
				"error_message": err,
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": fmt.Sprintf("Sucessfully create new product: %d", productID),
		})

		return
	case "edit":
		if param.ID == 0 {
			log.Logger.WithFields(logrus.Fields{
				"param": param,
			}).Error("invalid request - product id is empty")
			c.JSON(http.StatusBadRequest, gin.H{
				"error_message": "Invalid Request",
			})

			return
		}

		product, err := h.ProductUsecase.EditProduct(c.Request.Context(), &param.Product)
		if err != nil {
			log.Logger.WithFields(logrus.Fields{
				"param": param,
			}).Errorf("h.ProductUsecase.EditProduct() got error %v", err)

			c.JSON(http.StatusInternalServerError, gin.H{
				"error_message": err,
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "Success edit product!",
			"product": product,
		})

		return
	case "delete":
		if param.ID == 0 {
			log.Logger.WithFields(logrus.Fields{
				"param": param,
			}).Error("invalid request - product id is empty")
			c.JSON(http.StatusBadRequest, gin.H{
				"error_message": "Invalid Request",
			})

			return
		}

		err := h.ProductUsecase.DeleteProduct(c.Request.Context(), param.ID)
		if err != nil {
			log.Logger.WithFields(logrus.Fields{
				"param": param,
			}).Errorf("h.ProductUsecase.DeleteProduct() got error %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"error_message": err,
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": fmt.Sprintf("Product %d successfully deleted!", param.ID),
		})
	default:
		log.Logger.Errorf("Invalid action: %s", param.Action)
		c.JSON(http.StatusBadRequest, gin.H{
			"error_message": "Invalid Action",
		})

		return
	}

}

// ProductCategoryManagement product category management by given c pointer of gin.Context.
func (h *ProductHandler) ProductCategoryManagement(c *gin.Context) {
	var param models.ProductCategoryManagementParameter
	if err := c.ShouldBindJSON(&param); err != nil {
		log.Logger.Error(err.Error()) // utk debugging
		c.JSON(http.StatusBadRequest, gin.H{
			"error_message": "Invalid Input",
		})
		return
	}

	if param.Action == "" {
		log.Logger.Error("missing parameter action")
		c.JSON(http.StatusBadRequest, gin.H{
			"error_message": "Missing required parameter",
		})

		return
	}

	switch param.Action {
	case "add":
		if param.ID != 0 {
			log.Logger.WithFields(logrus.Fields{
				"param": param,
			}).Error("invalid request - product category id is not empty")
			c.JSON(http.StatusBadRequest, gin.H{
				"error_message": "Invalid Request",
			})

			return
		}

		productCategoryID, err := h.ProductUsecase.CreateNewProductCategory(c.Request.Context(), &param.ProductCategory)
		if err != nil {
			log.Logger.WithFields(logrus.Fields{
				"param": param,
			}).Errorf("h.ProductUsecase.CreateNewProductCategory got error %v", err)

			c.JSON(http.StatusInternalServerError, gin.H{
				"error_message": err,
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": fmt.Sprintf("Successfully create new product category: %d", productCategoryID),
		})

		return
	case "edit":
		if param.ID == 0 {
			log.Logger.WithFields(logrus.Fields{
				"param": param,
			}).Error("invalid request - product id is empty")
			c.JSON(http.StatusBadRequest, gin.H{
				"error_message": "Invalid Request",
			})

			return
		}

		productCategory, err := h.ProductUsecase.EditProductCategory(c.Request.Context(), &param.ProductCategory)
		if err != nil {
			log.Logger.WithFields(logrus.Fields{
				"param": param,
			}).Errorf("h.ProductUsecase.EditProductCategory got error %v", err)

			c.JSON(http.StatusInternalServerError, gin.H{
				"error_message": err,
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message":         "Success Edit Product",
			"productCategory": productCategory,
		})

		return
	case "delete":
		if param.ID == 0 {
			log.Logger.WithFields(logrus.Fields{
				"param": param,
			}).Error("invalid request - product id is empty")
			c.JSON(http.StatusBadRequest, gin.H{
				"error_message": "Invalid Request",
			})
			return
		}

		err := h.ProductUsecase.DeleteProductCategory(c.Request.Context(), param.ID)
		if err != nil {
			log.Logger.WithFields(logrus.Fields{
				"param": param, // notes: kalau ada PII --> prevent print log PII data
			}).Errorf("h.ProductUsecase.DeleteProductCategory() got error %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"error_message": err,
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": fmt.Sprintf("Product Category ID %d successfully deleted!", param.ID),
		})

		return
	default:
		log.Logger.Errorf("Invalid action: %s", param.Action)
		c.JSON(http.StatusBadRequest, gin.H{
			"error_message": "Invalid Action",
		})

		return
	}
}
