# crypto currency payment module for web store
Web developer can accept cryptocurrency without pain.
Steps:
### 1. Create a Mixin Messenger account.
Visit https://mixin.one/messenger to download App from AppStore, Google Play.

中国大陆用户可以访问 https://a.app.qq.com/o/simple.jsp?pkgname=one.mixin.messenger  下载

### 2. Active developer account and create an app
Log in to https://developer.mixin.one with your mixin messenger account

There is a [tutorial](https://mixin-network.gitbook.io/mixin-network/mixin-messenger-app/create-bot-account)

### Clone, build, run
```shell
git clone https://github.com/myrual/mixin-network-snapshot-golang
cd mixin-network-snapshot-golang
```

2. Edit some of code
```go
const (
	userid      = "3c5fd587-5ac3-4fb6-b294-423ba3473f7d"
	sessionid   = "42848ded-0ffd-45eb-9b46-094d5542ee01"
	private_key = `-----BEGIN RSA PRIVATE KEY-----
MIICXAIBAAKBgQDACTrT4uaB9el9qe0MUOsFrm8kpaDI9PowauMB1Ha25mpfL+5h
MFqISLS5z2P89nAsXBg+KyQ2gAVA6rBwW/ZqOc1PKiJhhLBS80nzo3ayfv7OlzNG
IxMyqD5izCCtPixnqpuTePoPWq4CNZlxop0VYklsWEfU0U0qqMBgmtqYfQIDAQAB
AoGAR8crZed5oTn5fC73m5LjRcxdXqVJ49MtcMuC7jwr41FckRepUkpwjGAgrRMH
nJXAd9Q0e4hEkNppHEqciGLXR1dQfZnaM1Gnv7mD3oSgHaH+4qAMnNOCpvwW4Eu3
yp9b1UGj9SvM3D2BrpA+MGf0E/yEJzpRcT956W6SPYYSegECQQDm4uTK+teoxr1Z
agJZuCta+IhMzpxIWMob+JN/Huf7OnRcIa9JpXngg4tHOUWmZCDQdqeJMpaQc8SQ
44hba015AkEA1OyJswNIhdmvVp5P1zgREVVRK6JloYwmAtj+Qo4pWJ117LqH4w+b
491r4AeLEGh8VrZ4k6Hp+Cm783S2jTAWJQJARbWdlHdV45xVkQiDuyjy1h2RsXb0
EpfUNcvAZLIlImIMvcBh1x+CA7pTs+Zj1BAJJEee37qJYQXDBGfeRJPKKQJAVG+c
x42Ew/eoTZwoIzvLoOkJcFlNHjwaksSER9ZiVQ7URdVOr99vvXQAJG45Wn9k12oy
9LCfvNan/wqIngK0tQJBAL1Wc02seEbMeWyt5jycJEhn6G8F18s9S1v0GXb4U/7/
6Y87P3TmDLcEuCXkrbZQaCX7jVLu0BkDw8To58TWjh0=	
-----END RSA PRIVATE KEY-----`
)
```
3. Build
```shell
go build mixin_snap.go
```
4. Run
```shell
./mixin_snap
```

## How to 
#### Create payment request

Create payment with request id value8, let the program visit http://127.0.0.1:9090/ with http POST 

```shell
curl -d '{"reqid":"value8", "callback":":9090/"}' -H "Content-Type: application/json" 127.0.0.1:8080/payment
```


#### Get payment status
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

#### callback url 
The program will visit following url when user pay to payment address
```json
"http://127.0.0.1"+callbackurl
```
the method is POST, http body will be
```json
{"Reqid":"value8","Callbackurl":":9090/","Paymentrecord":{"Amount":"0.01","AssetId":"56e63c06-b506-4ec5-885a-4a5ac17b83c1","created_at":"2019-06-20T07:33:06.445471337Z","snapshot_id":"a6603374-509b-4015-a192-c63bfa8def5f"}}
```


### Get all payment asset
1. All income payment will be automatically sent to your own Mixin Messenger account.
### Manually send all money to your Mixin Messenger account
```shell
curl -X POST -H "Content-Type: application/json" 127.0.0.1:8080/moneygohome
```

response will be similar to follow
```json
total 20 account will send all balance to admin
```
