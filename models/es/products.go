package es

type Products struct {
	Id            string  `json:"id" bson:"id"`
	Name          string  `json:"name" bson:"name"`
	Description   string  `json:"description" bson:"description"`
	SellPrice     float64 `json:"sellPrice" bson:"sell_price"`
	OriginalPrice float64 `json:"originalPrice" bson:"original_price"`
	Vendor        string  `json:"vendor" bson:"vendor"`
	CreatedAt     string  `json:"createdAt" bson:"created_at"`
	UpdateAt      string  `json:"updatedAt" bson:"update_at"`
}
