package routes

import "github.com/razorpay-go/config"

func Init() {
	r := NewRouter()
	port := config.Port
	r.Run(port)
}
