package es

type Orders struct {
	Id                     string  `json:"id" bson:"id"`
	UserId                 string  `json:"userId" bson:"user_id"`
	ProductId              string  `json:"productId" bson:"product_id"`
	SellPrice              float64 `json:"sellPrice" bson:"sell_price"`
	OriginalPrice          float64 `json:"originalPrice" bson:"original_price"`
	Status                 string  `json:"status" bson:"status"`
	PaymentId              string  `json:"paymentId" bson:"paymentId"`
	PaymentMerchant        string  `json:"paymentMerchant" bson:"payment_merchant"`
	PaymentMerchantOrderId string  `json:"paymentMerchantOrderId" bson:"payment_merchant_order_id"`
	CreatedAt              string  `json:"createdAt" bson:"created_at"`
	UpdateAt               string  `json:"updatedAt" bson:"update_at"`
}
