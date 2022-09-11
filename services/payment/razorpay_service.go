package payment

import (
	"crypto/hmac"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/hex"
	"encoding/json"
	"errors"
	"strings"

	"github.com/razorpay-go/config"
	"github.com/razorpay-go/constants"
	esModel "github.com/razorpay-go/models/es"
	"github.com/razorpay-go/repository/es"
	razorpay "github.com/razorpay/razorpay-go"
)

var razorpayClient = razorpay.NewClient(config.RazorpayKey, config.RazorpaySecret)
var orderRepo = es.OrderRepo{}

type RazorPayService struct{}

func (RazorPayService) CreatePaymentOrder(amount float64, currency string, orderId string) (string, error) {
	amountInBaseUnit := amount * constants.PAYMENT_BASE_UNIT
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
		return status, errors.New(constants.REQUIRED_KEY_VAL)
	}
	orderDetail, err := getOrder(paymentRequest.OrderId)
	if err != nil {
		return status, err
	}
	orderDetail.PaymentId = paymentRequest.MerchantPaymentId
	if generateSignature(paymentRequest.OrderId, orderDetail.PaymentId, paymentRequest.MerchantSignature) == false {
		return status, errors.New(constants.INVD_SIG)
	}
	// getting payment status at razor pay side
	currentPaymentStatus, err := getPaymentCurrentStatus(orderDetail.PaymentId)
	if err != nil {
		return status, err
	}
	// update db status
	orderDetail.Status = currentPaymentStatus
	_, orderUpdateErr := orderRepo.CreateOne(orderDetail)
	if orderUpdateErr != nil {
		return status, orderUpdateErr
	}
	if orderDetail.Status == constants.PAYMENT_FAILED_STATUS {
		return status, errors.New(constants.VERIFICATION_FAILED)
	}
	return true, nil
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
		return o, errors.New(constants.ORDER_ID_REQ)
	}
	order, err := orderRepo.FindOnePaymentOrder(id)
	if err != nil {
		return order, err
	}
	if order.Id == "" {
		return order, errors.New(constants.INVD_ORD_ID)
	}
	return order, nil
}

func getPaymentCurrentStatus(paymentId string) (string, error) {
	status, err := razorpayClient.Payment.Fetch(paymentId, nil, nil)
	if err != nil {
		return "", err
	}
	currentStatus := status["status"].(string)
	// possible CAPTURED, AUTHORIZED, FAILED
	returnStatus := strings.ToUpper(currentStatus)
	if returnStatus == constants.PAYMENT_CAPTURED_STATUS || returnStatus == constants.PAYMENT_AUTHORIZED_STATUS {
		returnStatus = constants.PAYMENT_COMPLETED_STATUS
	} else if returnStatus == constants.PAYMENT_FAILED_STATUS {
		returnStatus = constants.PAYMENT_FAILED_STATUS
	}
	return returnStatus, nil
}
