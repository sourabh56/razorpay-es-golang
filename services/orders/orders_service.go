package orders

import (
	"encoding/json"
	"errors"
	"time"

	"github.com/razorpay-go/constants"
	esModel "github.com/razorpay-go/models/es"
	"github.com/razorpay-go/repository/es"
	"github.com/razorpay-go/services/payment"
	"github.com/razorpay-go/services/products"
	"github.com/razorpay-go/utils"
)

type OrdersService struct{}

var productService = products.ProductService{}
var orderRepo = es.OrderRepo{}
var paymentService = payment.RazorPayService{}

func (OrdersService) PlaceOrder(body string) (esModel.Orders, error) {
	r := esModel.Orders{}
	type orderRequestS struct {
		ProductId string `json:"productId"`
	}
	orderRequest := orderRequestS{}
	er := json.Unmarshal([]byte(body), &orderRequest)
	if er != nil {
		return r, er
	}
	// validate if product id is correct
	productDetail, productValidationError := getProductDetail(orderRequest.ProductId)
	if productValidationError != nil {
		return r, productValidationError
	}
	// place request with status pending
	orderResponse, err := initOrderRequest(productDetail)
	if err != nil {
		return orderResponse, err
	}
	return orderResponse, nil
}

func getProductDetail(id string) (esModel.Products, error) {
	var productDetail = esModel.Products{}
	if id == "" {
		return productDetail, errors.New("Product id is required.")
	}
	productExists, err := productService.FindProduct(id)
	if err != nil {
		return productDetail, err
	} else if productExists.Id == "" {
		return productDetail, errors.New("Product not found.")
	}
	return productExists, nil
}

func initOrderRequest(productDetail esModel.Products) (esModel.Orders, error) {
	var orderRequest = esModel.Orders{
		CreatedAt:       time.Now().Format(constants.MONGO_TIME_FORMAT),
		UpdateAt:        time.Now().Format(constants.MONGO_TIME_FORMAT),
		Id:              utils.RandomString(24),
		UserId:          utils.RandomString(4),
		ProductId:       productDetail.Id,
		SellPrice:       productDetail.SellPrice,
		OriginalPrice:   productDetail.OriginalPrice,
		Status:          "PENDING",
		PaymentId:       "",
		PaymentMerchant: "RAZORPAY",
	}
	orderInitResponse, orderError := orderRepo.CreateOne(orderRequest)
	if orderError != nil {
		return orderInitResponse, orderError
	}
	paymentOrderId, orderCreationError := paymentService.CreatePaymentOrder(orderInitResponse.SellPrice, "INR", orderInitResponse.Id)
	if orderCreationError != nil {
		return orderInitResponse, orderCreationError
	}
	orderInitResponse.PaymentMerchantOrderId = paymentOrderId

	orderupdateResponse, updateOrderError := orderRepo.CreateOne(orderInitResponse)
	if updateOrderError != nil {
		return orderupdateResponse, updateOrderError
	}
	return orderupdateResponse, orderError
}

func (OrdersService) UpdateOrderStatus(data esModel.Orders) {

}
