curl -X GET 'http://localhost:8080/payment?reqid=abcde'
curl -d '{"reqid":"value8", "callback":":9090/"}' -H "Content-Type: application/json" 127.0.0.1:8080/payment

curl -X POST -H "Content-Type: application/json" 127.0.0.1:8080/moneygohome
