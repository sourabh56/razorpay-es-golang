package controllers

import (
	"io/ioutil"

	"github.com/gin-gonic/gin"
	"github.com/razorpay-go/services/orders"
)

type OrdersController struct{}

var OrderService = orders.OrdersService{}

func (OrdersController) PlaceProductOrder(c *gin.Context) {
	var apiRes = ApiResponse{
		Status: "error",
	}
	var stringBody = ""
	if body, err := ioutil.ReadAll(c.Request.Body); err != nil {
		if err != nil {
			apiRes.Message = err.Error()
			c.AbortWithStatusJSON(500, apiRes)
		}
	} else {
		stringBody = string(body)
	}
	response, err := OrderService.PlaceOrder(stringBody)
	if err != nil {
		apiRes.Message = err.Error()
		c.AbortWithStatusJSON(500, apiRes)
		return
	}
	apiRes.Data = response
	apiRes.Status = "success"
	c.JSON(200, apiRes)
	return
}
