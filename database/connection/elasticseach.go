package connection

import (
	"fmt"

	elasticsearch7 "github.com/elastic/go-elasticsearch/v7"
	"github.com/razorpay-go/config"
)

var es *elasticsearch7.Client

func CreateEsConnection() {
	cfg := elasticsearch7.Config{
		Addresses: []string{
			config.EsUrl,
		},
	}
	esc, err := elasticsearch7.NewClient(cfg)
	if err != nil {
		panic(err.Error())
	}
	fmt.Printf("Es %s connected", config.EsUrl)
	es = esc
}

func GetEsDB() *elasticsearch7.Client {
	return es
}
