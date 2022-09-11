package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/razorpay-go/services/products"
)

type ProductsController struct{}

type ApiResponse struct {
	Data    interface{} `json:"data"`
	Status  string      `json:"status"`
	Message string      `json:"message"`
}

var ProductService = products.ProductService{}

func (ProductsController) GetProducts(c *gin.Context) {
	urlQuery := c.Request.URL.Query()
	query := urlQuery.Get("query")
	response, err := ProductService.SearchProducts(query)
	var apiRes = ApiResponse{
		Data:   response,
		Status: "error",
	}
	if err != nil {
		apiRes.Message = err.Error()
		c.AbortWithStatusJSON(500, apiRes)
		return
	}
	apiRes.Status = "success"
	c.JSON(200, apiRes)
	return
}
