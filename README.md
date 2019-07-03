<p align="center">
<a href="README_CN.md"><img src="https://img.shields.io/badge/language-中文文档-red.svg?longCache=true&style=flat-square"></a>
</p>

# Cryptocurrency payment plugin

Web developer can accept cryptocurrency without understand Bitcoin, EOS full API. No need to setup full node.

The standalone program is a battery included solution. Developer just need to call it's http api, show payment information to client, program will visit callback url when your client paid cryptocurrency. Cryptocurrency is automatically transferred to your account when your client paid.

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
#### Accept cryptocurrency payment
To accept bitcoin or eos payment, developer need to call localhost:8080/payment by http POST,  with parameter in body. An unique string, a callback url and an expired timer should be in body. The unique string can be anything like uuid. Callback url will be visited by the program when your client paid to you. The callback mechanism will be expired if 60 minutes if you give expiredafter 60, the callback will always work if you give it a ZERO.

Following curl is an example:
```shell
curl -d '{"reqid":"value8", "callback":":9090/", "expiredafter":60}' -H "Content-Type: application/json" 127.0.0.1:8080/payment
```
The command just tell the program to create a payment address for an unique string "value8", visit localhost:9090/ when user paid in 60 minutes.

the result of the command will be 
```json
{"Reqid":"value8","Payment_methods":[{"Name":"XLM","PaymentAddress":"","PaymentAccount":"GD77JOIFC622O5HXU446VIKGR5A5HMSTAUKO2FSN5CIVWPHXDBGIAG7Y","PaymentMemo":"3f8db42022b5bc32","Priceinusd":"0.10472789","Priceinbtc":"0.00000925"},{"Name":"EOS","PaymentAddress":"","PaymentAccount":"eoswithmixin","PaymentMemo":"302c37ebff05ccf09dd7296053d1924a","Priceinusd":"5.9436916","Priceinbtc":"0.00052505"},{"Name":"ETH","PaymentAddress":"0x365DA43BC7B22CD4334c3f35eD189C8357D4bEd6","PaymentAccount":"","PaymentMemo":"","Priceinusd":"295.86024062","Priceinbtc":"0.02613571"}],"Payment_records":null,"Balance":null}
```
Your client need value in Payment_methods. There three payment methods in the example.

There are two types of payment method:
1. Bitcoin/Ethereum style: PaymentAddress is not empty, PaymentAccount and PaymentMemo are all empty. You just  show Ethererum Name and PaymentAddress to your clients, they just need to transfer token to the address.
2. EOS/Stellar style: PaymentAddress is empty, PaymentAccount and PaymentMemo are not empty. You need to show Asset Name and both of PaymentAccount and PaymentMemo to user, and remind user need to input BOTH of PaymentAccount and PaymentMemo. Transfer asset to PaymentAccount without memo is a common mistake, and it can not be reverted because current Mixin Network limitation.

Asset current price in USD and Bitcoin is inside payment record.
```json
{"Priceinusd":"0.10472789","Priceinbtc":"0.00000925"}
```


#### Query payment status
fetch the payment status by visit localhost:8080/payment with parameter reqid

Example:
```shell
curl -X GET 'http://localhost:8080/payment?reqid=value8'
```

Response will be similar to following if payment is not yet confirmed
```json
{"Reqid":"value8","Payment_methods":[{"Name":"XLM","PaymentAddress":"","PaymentAccount":"GD77JOIFC622O5HXU446VIKGR5A5HMSTAUKO2FSN5CIVWPHXDBGIAG7Y","PaymentMemo":"3f8db42022b5bc32","Priceinusd":"0.10472789","Priceinbtc":"0.00000925"},{"Name":"EOS","PaymentAddress":"","PaymentAccount":"eoswithmixin","PaymentMemo":"302c37ebff05ccf09dd7296053d1924a","Priceinusd":"5.9436916","Priceinbtc":"0.00052505"},{"Name":"ETH","PaymentAddress":"0x365DA43BC7B22CD4334c3f35eD189C8357D4bEd6","PaymentAccount":"","PaymentMemo":"","Priceinusd":"295.86024062","Priceinbtc":"0.02613571"}],"Payment_records":null,"Balance":null}
```
The paymnet_records is empty here.

Response will be similar to following if payment is already confirmed
```json
{"Reqid":"value8","Payment_methods":[{"Name":"XLM","PaymentAddress":"","PaymentAccount":"GD77JOIFC622O5HXU446VIKGR5A5HMSTAUKO2FSN5CIVWPHXDBGIAG7Y","PaymentMemo":"3f8db42022b5bc32","Priceinusd":"0.10472789","Priceinbtc":"0.00000925"},{"Name":"EOS","PaymentAddress":"","PaymentAccount":"eoswithmixin","PaymentMemo":"302c37ebff05ccf09dd7296053d1924a","Priceinusd":"5.9436916","Priceinbtc":"0.00052505"},{"Name":"ETH","PaymentAddress":"0x365DA43BC7B22CD4334c3f35eD189C8357D4bEd6","PaymentAccount":"","PaymentMemo":"","Priceinusd":"295.86024062","Priceinbtc":"0.02613571"}],"Payment_records":[{"Amount":"0.1","AssetId":"","created_at":"2019-06-20T02:00:39.650472961Z","snapshot_id":"570233aa-3c91-45cd-a6ec-0e9724165300"},{"Amount":"0.01","AssetId":"6cfe566e-4aad-470b-8c9a-2fd35b49c68d","created_at":"2019-06-20T02:33:50.152539755Z","snapshot_id":"88859d4d-5bee-4fb5-aef6-ac01dc3a43c6"},{"Amount":"0.01","AssetId":"6cfe566e-4aad-470b-8c9a-2fd35b49c68d","created_at":"2019-06-20T02:37:05.870885973Z","snapshot_id":"6530f455-3238-491a-a9c5-bbcb52bcc306"},{"Amount":"0.001","AssetId":"6cfe566e-4aad-470b-8c9a-2fd35b49c68d","created_at":"2019-06-20T02:40:53.251365044Z","snapshot_id":"f2c8a751-3d30-472e-bf76-924787f341b9"},{"Amount":"0.001","AssetId":"6cfe566e-4aad-470b-8c9a-2fd35b49c68d","created_at":"2019-06-20T02:59:28.854380284Z","snapshot_id":"3ebfd5a3-bd29-4e32-bd06-2506bee3da99"},{"Amount":"-0.122","AssetId":"6cfe566e-4aad-470b-8c9a-2fd35b49c68d","created_at":"2019-06-20T03:00:17.249302744Z","snapshot_id":"0bfe6f6b-1ff8-4144-9786-52d6a6459b19"}],"Balance":null}
```
the payment_records are filled by transaction information. one
```json
{"Amount":"0.01","AssetId":"6cfe566e-4aad-470b-8c9a-2fd35b49c68d",
"created_at":"2019-06-20T02:37:05.870885973Z","snapshot_id":"6530f455-3238-491a-a9c5-bbcb52bcc306"}
```
This is a payment from user: 
* amount 0.01
* asset id 6cfe566e-4aad-470b-8c9a-2fd35b49c68d, EOS main chain token
* create at UTC 2019-06-20T02:37:05.870885973
* unique snapshot id in Mixin Network is 6530f455-3238-491a-a9c5-bbcb52bcc306, you can verify the payment by https://mixin.one/snapshots/6530f455-3238-491a-a9c5-bbcb52bcc306

#### callback url 
The program will visit following url when payment is confirmed.
```json
"http://127.0.0.1"+callbackurl
```
The http visit method is POST, json body parameter is following
```json
{"Reqid":"value8","Callbackurl":":9090/","Paymentrecord":{"Amount":"0.01","AssetId":"56e63c06-b506-4ec5-885a-4a5ac17b83c1","created_at":"2019-06-20T07:33:06.445471337Z","snapshot_id":"a6603374-509b-4015-a192-c63bfa8def5f"}}
```


### Did all asset belongs to developer?
1. All income payment will be AUTOMATICALLY sent to your own Mixin Messenger account with ZERO transaction fee in 1 seconds. 
2. You can also ask the program send all money to your Mixin Messenger account if the program exit accidently.

```shell
curl -X POST -H "Content-Type: application/json" 127.0.0.1:8080/moneygohome
```

response will be similar to follow
```json
total 20 account will send all balance to admin
```

### payment confirmation time
1. EOS: 3 minutes
2. Stellar: 2 minutes
3. Bitcoin/USDT: 60 minutes
4. Litecoin/Ethererum/DOGE: 120 minutes

Why so long? Because it is Mixin Network confirmation time for different asset. You can not change now.

### What kind of currency can be supported
All asset supported by Mixin Network:
BTC, USDT, BCH, ETH and ERC20, ETC, EOS and token issue on EOS, DASH, Litecoin, Doge, Horizen, MGD, NEM, XRP, XLM, TRON and TRC10, Zcash. 

### Recommend Currency
Current default cryptocurrency is EOS and XLM because transaction can be confirmed in 3 minutes.

To support more currency, just add more asset into the default_asset_id_group.
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


TO BE DONE:
1. All asset can be withdrawed to developer's cold wallet.
2. One type asset can be exchanged to USDT or Bitcoin automatically through DEX.
3. Support Mixin Messenger User to pay.
4. ~~Latest USD price for every asset~~
