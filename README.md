# A plugin enable web app accept crypto currency payment
Web developer can accept cryptocurrency without pain.
Steps:
### 1. Create a Mixin Messenger account.
Visit https://mixin.one/messenger to download App from AppStore, Google Play.

中国大陆用户可以访问 https://a.app.qq.com/o/simple.jsp?pkgname=one.mixin.messenger  下载

### 2. Active developer account and create an app
Log in to https://developer.mixin.one with your mixin messenger account

This [tutorial](https://mixin-network.gitbook.io/mixin-network/mixin-messenger-app/create-bot-account) is very useful for new developer to create app.

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
	ADMIN_MessengerID = ""//this is your mixin messenger id, you can find your id in contact page.
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

5. Database

A sqlite3 file with name test.db will be generated in same folder.

## How to 
#### Create payment request

Make a unique string, and setup the callback url. The program will visit the callback URL if client pay to the deposit address

Example:
```shell
curl -d '{"reqid":"value8", "callback":":9090/"}' -H "Content-Type: application/json" 127.0.0.1:8080/payment
```
the result from the curl is following, you can see the deposit address for EOS and XLM.
```json
{"Reqid":"value9","Payment_methods":[{"Name":"XLM","PublicKey":"","AccountName":"GD77JOIFC622O5HXU446VIKGR5A5HMSTAUKO2FSN5CIVWPHXDBGIAG7Y","AccountTag":"39819a44ac87dd2c"},{"Name":"EOS","PublicKey":"","AccountName":"eoswithmixin","AccountTag":"7648a59ae0eaee11d5d7f90c0f334eb1"}],"Payment_records":null,"Balance":null}
```
#### Get payment status
fetch the payment status

Example:
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
http method is POST, http body will be
```json
{"Reqid":"value8","Callbackurl":":9090/","Paymentrecord":{"Amount":"0.01","AssetId":"56e63c06-b506-4ec5-885a-4a5ac17b83c1","created_at":"2019-06-20T07:33:06.445471337Z","snapshot_id":"a6603374-509b-4015-a192-c63bfa8def5f"}}
```


### Get all payment asset
1. All income payment will be AUTOMATICALLY sent to your own Mixin Messenger account.
2. You can also manually send all money to your Mixin Messenger account
```shell
curl -X POST -H "Content-Type: application/json" 127.0.0.1:8080/moneygohome
```

response will be similar to follow
```json
total 20 account will send all balance to admin
```

### Confirmation times
1. EOS: 3 minutes

### Support cryptocurrency
BTC, USDT, BCH, ETH and ERC20,ETC, EOS and token issue on, DASH, Litecoin, Doge,Horizen, MGD, NEM, XRP, XLM, TRON and TRC10, Zcash. 

### Recommend Currency

Default Crypty currency is EOS and XLM because transaction can be confirmed in 3 minutes.

You can change the 	default_asset_id_group to support more ASSET.
```go
const (
	BTC_ASSET_ID  = "c6d0c728-2624-429b-8e0d-d9d19b6592fa"
	EOS_ASSET_ID  = "6cfe566e-4aad-470b-8c9a-2fd35b49c68d"
	USDT_ASSET_ID = "815b0b1a-2764-3736-8faa-42d694fa620a"
	ETC_ASSET_ID  = "2204c1ee-0ea2-4add-bb9a-b3719cfff93a"
	XRP_ASSET_ID  = "23dfb5a5-5d7b-48b6-905f-3970e3176e27"
	XEM_ASSET_ID  = "27921032-f73e-434e-955f-43d55672ee31"
	ETH_ASSET_ID  = "43d61dcd-e413-450d-80b8-101d5e903357"
	DASH_ASSET_ID = "6472e7e3-75fd-48b6-b1dc-28d294ee1476"
	DOGE_ASSET_ID = "6770a1e5-6086-44d5-b60f-545f9d9e8ffd"
	LTC_ASSET_ID  = "76c802a2-7c88-447f-a93e-c29c9e5dd9c8"
	SIA_ASSET_ID  = "990c4c29-57e9-48f6-9819-7d986ea44985"
	ZEN_ASSET_ID  = "a2c5d22b-62a2-4c13-b3f0-013290dbac60"
	ZEC_ASSET_ID  = "c996abc9-d94e-4494-b1cf-2a3fd3ac5714"
	BCH_ASSET_ID  = "fd11b6e3-0b87-41f1-a41f-f0e9b49e5bf0"
	XIN_ASSET_ID  = "c94ac88f-4671-3976-b60a-09064f1811e8"
	CNB_ASSET_ID  = "965e5c6e-434c-3fa9-b780-c50f43cd955c"
	XLM_ASSET_ID  = "56e63c06-b506-4ec5-885a-4a5ac17b83c1"
	TRON_ASSET_ID = "25dabac5-056a-48ff-b9f9-f67395dc407c"
	........
)
.......
.......

	// to support more asset, just add them in the following array
	default_asset_id_group := []string{XLM_ASSET_ID, EOS_ASSET_ID}
```
|Asset|estimate confirmation duration|
|-|-|
|TRON| 5 minutes|
|BTC| 1 hour|
|USDT Omni|  1 hour|
