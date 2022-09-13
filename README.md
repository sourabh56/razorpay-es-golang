# README

## About

POC to setup **razorpay golang gin** integartion.
We have used elastiserarch to keep data
products > keep product details
orders > keep product order status

## Mapping
mapping is available in **./mapping/es**

## APIs
**GET** http://localhost:5001/api/v1/products?query=Downshifter


**POST** http://localhost:5001/api/v1/order


**POST**  http://localhost:5001/api/v1/payment/verify


## CURL
``` sh
curl --location --request GET 'localhost:5001/api/v1/products?query=Downshifter' \
--header 'Authorization: a'
```

``` sh
curl --location --request POST 'localhost:5001/api/v1/order' \
--header 'Authorization;' \
--header 'Content-Type: application/json' \
--data-raw '{
    "productId": "631c69b9a05e47120e00d792"
}'
```


``` sh
curl --location --request POST 'localhost:5001/api/v1/payment/verify' \
--header 'Authorization;' \
--header 'Content-Type: application/json' \
--data-raw '{
    "orderId": "order_KGaSnjgysnmqDV",
    "merchantPaymentId": "pay_KGaWeoPJkCQnnX",
    "merchantSignature": "9f93f4bcb6ff4ac59ce42e736811e392483d9526f1bf8c5b8798b5adb26eb280"
}'
```
