curl -X GET 'http://localhost:8080/payment?reqid=abcde'
curl -d '{"reqid":"value7", "callback":"callbackurl"}' -H "Content-Type: application/json" 127.0.0.1:8080/payment
