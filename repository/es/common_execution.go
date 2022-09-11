package es

import (
	"encoding/json"

	"github.com/elastic/go-elasticsearch/v7/esapi"
	esModel "github.com/razorpay-go/models/es"
)

func CommonExecution(res *esapi.Response, err error) (esModel.CommonEsFields, error) {
	var commonEsFields esModel.CommonEsFields
	if err != nil {
		return commonEsFields, err
	}
	defer res.Body.Close()
	if res.IsError() {
		var e map[string]interface{}
		if err := json.NewDecoder(res.Body).Decode(&e); err != nil {
			return commonEsFields, err
		}
	}
	var esResponse map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&esResponse); err != nil {
		return commonEsFields, err
	}
	responseByte, _ := json.Marshal(esResponse)
	json.Unmarshal([]byte(responseByte), &commonEsFields)
	return commonEsFields, nil
}
