package es

import (
	"bytes"
	"context"
	"encoding/json"
	"strings"

	"github.com/elastic/go-elasticsearch/v7/esapi"
	"github.com/razorpay-go/database/connection"
	esModel "github.com/razorpay-go/models/es"
)

type OrderRepo struct{}

func (OrderRepo) CreateOne(body esModel.Orders) (esModel.Orders, error) {
	var response esModel.Orders
	es := connection.GetEsDB()
	orderId := body.Id
	bodyString, _ := json.Marshal(body)

	resInit := esapi.IndexRequest{
		Index:        "orders",
		DocumentID:   orderId,
		DocumentType: "_doc",
		Refresh:      "true",
		Body:         strings.NewReader(string(bodyString)),
	}
	res, err := resInit.Do(context.Background(), es)
	commonEsFields, err := CommonExecution(res, err)
	if err != nil {
		return response, err
	}
	orderDetails, orderFindError := OrderRepo{}.FindOne(commonEsFields.ESId)
	if orderFindError != nil {
		return orderDetails, orderFindError
	}
	return orderDetails, nil
}

func (OrderRepo) FindOne(id string) (esModel.Orders, error) {
	var response esModel.Orders
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
		es.Search.WithIndex("orders"),
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

func (OrderRepo) FindOnePaymentOrder(paymentOrderId string) (esModel.Orders, error) {
	var response esModel.Orders
	var buf bytes.Buffer
	query := map[string]interface{}{
		"query": map[string]interface{}{
			"term": map[string]interface{}{
				"paymentMerchantOrderId": paymentOrderId,
			},
		},
	}

	if err := json.NewEncoder(&buf).Encode(query); err != nil {
		return response, err
	}
	es := connection.GetEsDB()
	res, err := es.Search(
		es.Search.WithContext(context.Background()),
		es.Search.WithIndex("orders"),
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
