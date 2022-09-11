package payment

import (
	"crypto/hmac"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/hex"
	"encoding/json"
	"errors"

	"github.com/razorpay-go/config"
	esModel "github.com/razorpay-go/models/es"
	"github.com/razorpay-go/repository/es"
	razorpay "github.com/razorpay/razorpay-go"
)

var razorpayClient = razorpay.NewClient(config.RazorpayKey, config.RazorpaySecret)
var orderRepo = es.OrderRepo{}

type RazorPayService struct{}

func (RazorPayService) CreatePaymentOrder(amount float64, currency string, orderId string) (string, error) {
	amountInBaseUnit := amount * 100
	createRequest := map[string]interface{}{
		"amount":   amountInBaseUnit,
		"currency": currency,
	}
	body, err := razorpayClient.Order.Create(createRequest, nil)
	if err != nil {
		return "", err
	}
	return body["id"].(string), nil
}

func (RazorPayService) VerifyPayment(body string) (bool, error) {
	status := false
	type paymentRequestS struct {
		OrderId           string `json:"orderId"`
		MerchantPaymentId string `json:"merchantPaymentId"`
		MerchantSignature string `json:"merchantSignature"`
	}
	paymentRequest := paymentRequestS{}
	er := json.Unmarshal([]byte(body), &paymentRequest)
	if er != nil {
		return status, er
	}
	if paymentRequest.OrderId == "" || paymentRequest.MerchantPaymentId == "" || paymentRequest.MerchantSignature == "" {
		return status, errors.New("Validation error, required keys are missing.")
	}
	orderDetail, err := getOrder(paymentRequest.OrderId)
	if err != nil {
		return status, err
	}
	orderDetail.PaymentId = paymentRequest.MerchantPaymentId
	if generateSignature(paymentRequest.OrderId, paymentRequest.MerchantPaymentId, paymentRequest.MerchantSignature) {
		orderDetail.Status = "COMPLETED"
		status = true
	} else {
		orderDetail.Status = "FAILED"
	}
	_, orderUpdateErr := orderRepo.CreateOne(orderDetail)
	if orderUpdateErr != nil {
		return false, orderUpdateErr
	}
	if orderDetail.Status == "FAILED" {
		return status, errors.New("Verification failed.")
	}
	return status, nil
}

func generateSignature(orderId string, paymentId string, signature string) bool {
	data := orderId + "|" + paymentId
	h := hmac.New(sha256.New, []byte(config.RazorpaySecret))
	_, err := h.Write([]byte(data))

	if err != nil {
		return false
	}
	sha := hex.EncodeToString(h.Sum(nil))

	if subtle.ConstantTimeCompare([]byte(sha), []byte(signature)) == 1 {
		return true
	}
	return false
}

func getOrder(id string) (esModel.Orders, error) {
	o := esModel.Orders{}
	if id == "" {
		return o, errors.New("order id is required.")
	}
	order, err := orderRepo.FindOnePaymentOrder(id)
	if err != nil {
		return order, err
	}
	if order.Id == "" {
		return order, errors.New("Invalid order id.")
	}
	return order, nil
}
