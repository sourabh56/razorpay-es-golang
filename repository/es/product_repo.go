package es

import (
	"bytes"
	"context"
	"encoding/json"

	"github.com/razorpay-go/database/connection"
	esModel "github.com/razorpay-go/models/es"
)

type ProductRepo struct{}

func (ProductRepo) FindAll(productQuery string) ([]esModel.Products, error) {
	var response []esModel.Products
	var buf bytes.Buffer
	query := map[string]interface{}{
		"query": map[string]interface{}{
			"query_string": map[string]interface{}{
				"fields": []string{"description"},
				"query":  productQuery,
			},
		},
	}

	if err := json.NewEncoder(&buf).Encode(query); err != nil {
		return response, err
	}
	es := connection.GetEsDB()
	res, err := es.Search(
		es.Search.WithContext(context.Background()),
		es.Search.WithIndex("products"),
		es.Search.WithBody(&buf),
		es.Search.WithTrackTotalHits(true),
		es.Search.WithPretty(),
	)
	commonEsFields, err := CommonExecution(res, err)
	if err != nil {
		return response, err
	}
	for _, hit := range commonEsFields.Hits.Hits {
		temp, _ := json.Marshal(hit.Source)
		var tempProducts esModel.Products
		json.Unmarshal([]byte(temp), &tempProducts)
		response = append(response, tempProducts)
	}
	return response, nil
}

func (ProductRepo) FindOne(id string) (esModel.Products, error) {
	var response esModel.Products
	var buf bytes.Buffer
	query := map[string]interface{}{
		"query": map[string]interface{}{
			"term": map[string]interface{}{
				"_id": id,
			},
		},
	}

	if err := json.NewEncoder(&buf).Encode(query); err != nil {
		return response, err
	}
	es := connection.GetEsDB()
	res, err := es.Search(
		es.Search.WithContext(context.Background()),
		es.Search.WithIndex("products"),
		es.Search.WithBody(&buf),
		es.Search.WithTrackTotalHits(true),
		es.Search.WithPretty(),
	)
	commonEsFields, err := CommonExecution(res, err)
	if err != nil {
		return response, err
	}
	for _, hit := range commonEsFields.Hits.Hits {
		temp, _ := json.Marshal(hit.Source)
		json.Unmarshal([]byte(temp), &response)
	}
	return response, nil
}
