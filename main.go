package main

import (
	"github.com/razorpay-go/database/connection"
	"github.com/razorpay-go/routes"
)

func main() {
	connection.CreateEsConnection()
	routes.Init()
}
