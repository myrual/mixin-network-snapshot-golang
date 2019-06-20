# crypto currency payment module for web store
Web developer can accept cryptocurrency without pain.
Steps:
### 1. Create a Mixin Messenger account.
### 2. Active developer account and create an app
Log in to https://developer.mixin.one with your mixin messenger account

### Clone, build, run
1. Clone 
2. Edit config
3. Build
4. Run


### Create payment by http post
Create payment with request id value8, let the program visit http://127.0.0.1:9090/ with http POST 

```shell
curl -d '{"reqid":"value8", "callback":":9090/"}' -H "Content-Type: application/json" 127.0.0.1:8080/payment
```
The http post body will be 
```json
{"Reqid":"value8","Callbackurl":":9090/","Paymentrecord":{"Amount":"0.01","AssetId":"56e63c06-b506-4ec5-885a-4a5ac17b83c1","created_at":"2019-06-20T07:33:06.445471337Z","snapshot_id":"a6603374-509b-4015-a192-c63bfa8def5f"}}
```

### Get payment status by http get
fetch the payment status
```shell
curl -X GET 'http://localhost:8080/payment?reqid=value7'
```

Response will be similar to following if payment is not confirmed
```json
{"Reqid":"value6","Payment_methods":[{"Name":"XLM","PublicKey":"","AccountName":"GD77JOIFC622O5HXU446VIKGR5A5HMSTAUKO2FSN5CIVWPHXDBGIAG7Y","AccountTag":"dfc6af4e022c3a11"},{"Name":"EOS","PublicKey":"","AccountName":"eoswithmixin","AccountTag":"d457cab41245ca0531f64947d1bb958a"}],"Payment_records":null,"Balance":null}
```
Response will be similar to following if payment is already confirmed
```json
{"Reqid":"value8","Payment_methods":[{"Name":"XLM","PublicKey":"","AccountName":"GD77JOIFC622O5HXU446VIKGR5A5HMSTAUKO2FSN5CIVWPHXDBGIAG7Y","AccountTag":"62d0d256dcf15608"},{"Name":"EOS","PublicKey":"","AccountName":"eoswithmixin","AccountTag":"7481cd36f77953f129c194d3444ae2ff"}],"Payment_records":[{"Amount":"0.1","AssetId":"","created_at":"2019-06-20T02:00:39.650472961Z","snapshot_id":"570233aa-3c91-45cd-a6ec-0e9724165300"},{"Amount":"0.01","AssetId":"6cfe566e-4aad-470b-8c9a-2fd35b49c68d","created_at":"2019-06-20T02:33:50.152539755Z","snapshot_id":"88859d4d-5bee-4fb5-aef6-ac01dc3a43c6"},{"Amount":"0.01","AssetId":"6cfe566e-4aad-470b-8c9a-2fd35b49c68d","created_at":"2019-06-20T02:37:05.870885973Z","snapshot_id":"6530f455-3238-491a-a9c5-bbcb52bcc306"},{"Amount":"0.001","AssetId":"6cfe566e-4aad-470b-8c9a-2fd35b49c68d","created_at":"2019-06-20T02:40:53.251365044Z","snapshot_id":"f2c8a751-3d30-472e-bf76-924787f341b9"},{"Amount":"0.001","AssetId":"6cfe566e-4aad-470b-8c9a-2fd35b49c68d","created_at":"2019-06-20T02:59:28.854380284Z","snapshot_id":"3ebfd5a3-bd29-4e32-bd06-2506bee3da99"},{"Amount":"-0.122","AssetId":"6cfe566e-4aad-470b-8c9a-2fd35b49c68d","created_at":"2019-06-20T03:00:17.249302744Z","snapshot_id":"0bfe6f6b-1ff8-4144-9786-52d6a6459b19"}],"Balance":null}
```
### All income will be automatically sent to your own Mixin Messenger account
### Manually send all money to your Mixin Messenger account
```shell
curl -X POST -H "Content-Type: application/json" 127.0.0.1:8080/moneygohome
```

response will be similar to follow
```json
total 20 account will send all balance to admin
```
