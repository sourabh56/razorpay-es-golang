package routes

import (
	"time"

	"github.com/razorpay-go/constants"
	"github.com/razorpay-go/http/controllers"
	"github.com/razorpay-go/http/middlewares"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func NewRouter() *gin.Engine {
	router := gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	router.NoRoute(func(c *gin.Context) {
		c.JSON(404, gin.H{"code": 404, "message": constants.MSG_NOT_FOUND_ERR})
	})
	router.GET("/healthCheck", func(c *gin.Context) {
		c.JSON(200, gin.H{"code": 200, "message": ""})
	})

	router.Use(cors.New(cors.Config{
		AllowMethods:     []string{"POST", "GET", "DELETE"},
		AllowHeaders:     []string{"Origin", "x-api-key", "apikey"},
		AllowCredentials: true,
		AllowAllOrigins:  true,
		MaxAge:           10 * time.Minute,
	}))
	apiv1 := router.Group("/api/v1")
	OrdersController := new(controllers.OrdersController)
	ProductsController := new(controllers.ProductsController)
	PaymentController := new(controllers.PaymentController)

	apiv1.GET("/products", middlewares.TokenMiddleware(), ProductsController.GetProducts)
	apiv1.POST("/order", middlewares.TokenMiddleware(), OrdersController.PlaceProductOrder)
	apiv1.POST("/payment/verify", middlewares.TokenMiddleware(), PaymentController.VerifyPayment)

	return router
}
