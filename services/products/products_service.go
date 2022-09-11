package products

import (
	esModel "github.com/razorpay-go/models/es"
	"github.com/razorpay-go/repository/es"
)

type ProductService struct{}

var productRepo = es.ProductRepo{}

func (ProductService) SearchProducts(query string) ([]esModel.Products, error) {
	query = query + "*"
	products, err := productRepo.FindAll(query)
	if err != nil {
		return products, err
	}
	if len(products) == 0 {
		products = []esModel.Products{}
	}
	return products, nil
}

func (ProductService) FindProduct(id string) (esModel.Products, error) {
	productDetail, err := productRepo.FindOne(id)
	if err != nil {
		return productDetail, err
	}
	return productDetail, nil
}
