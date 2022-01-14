<p align="center">
<a href="README_CN.md"><img src="https://img.shields.io/badge/language-中文文档-red.svg?longCache=true&style=flat-square"></a>
</p>

# Cryptocurrency payment plugin(No maintainer NOW, No response if you ask question)

Accept cryptocurrency payment can be painless, free and secure. 
* No need to setup full Bitcoin/Ethereum/EOS full node(it cost your hundreds gigabytes).
* No need to pay expensive service fee, your program, your money.
* All payment is automatically transfer to your personal account on the fly. No money to lose even database is stolen.

Developer call to localhost http api, show payment information to client, program will visit webhook when client paid cryptocurrency.

ATTENTION: If you ever use code on or before tag v0.0.1, the current master branch is not backward compatible.

Steps:
### 1. Create a Mixin Messenger account.
Visit https://mixin.one/messenger to download App from AppStore, Google Play.

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

#### Query current cryptocurrency price, so developer know how many asset client need to transfer.
```shell
curl -X GET 'http://localhost:8080/assetsprice'
```
Result is following.
```json
[
	{"Fullname":"Stellar","Symbol":"XLM","USDPrice":0.10357796,"BTCPrice":0.00000889,"Assetid":"56e63c06-b506-4ec5-885a-4a5ac17b83c1"},
	{"Fullname":"EOS","Symbol":"EOS","USDPrice":5.96024263,"BTCPrice":0.00051165,"Assetid":"6cfe566e-4aad-470b-8c9a-2fd35b49c68d"},
	{"Fullname":"Ether","Symbol":"ETH","USDPrice":294.61322131,"BTCPrice":0.02529107,"Assetid":"43d61dcd-e413-450d-80b8-101d5e903357"}
]
```

If your order is valued about 1 USD, that means client need to deposit about 10 XLM, or 0.17 EOS.

#### Create charge
To accept bitcoin or eos payment, developer need to call localhost:8080/charges by http POST,  with parameter in body. 

POST /charges

|Attributes| type | description|
|--|--|--|
|currency| String | Currency code associated with the amount.  Only EOS/XLM/ETH is supported currently|
|amount| Float64 | Positive float|
|customerid| String | This field is optional and can be used to attach an identifier of your choice to the charge. Must not exceed 64 characters|
|webhookurl| String | program will visit localhost+webhook when user pay enough currency before charge is expired |
|expiredafter| uint | the webhook will be expired after xx minutes. User can pay to an expired charge , program keep income record and will transfer asset to admin account|


Example: let client "client1245" pay 0.001 ETH, notify developer's app by POST localhost:9090/123 when user pay enough currency in 60 minutes.
```shell
curl -d '{"currency":"ETH", "amount":0.001, "customerid":"client1245", "webhookurl":":9090/123", "expiredafter":60}' -H "Content-Type: application/json" 127.0.0.1:8080/charges
```
The command just tell the program to create a ETH charge address for customer id "client1245", visit localhost:9090/123 when user paid enough asset to the address in 60 minutes.

the result of the command will be 
```json
{
	"Id":3,
	"Currency":"ETH",
	"Amount":0.001,
	"Customerid":"client1245",
	"Webhookurl":":9090/123",
	"Expired_after":60,
	"Paymentmethod":{
		"Name":"ETH",
		"PaymentAddress":"0x130D3e6655f073e33235e567E7A1e1E1f59ddD79",
		"PaymentAccount":"",
		"PaymentMemo":"",
		"Priceinusd":"310.40105841",
		"Priceinbtc":"0.02374051"
		},
	"Receivedamount":0,
	"Paidstatus":0}
```
Client need to tranfser 0.001 ETH to address 0x130D3e6655f073e33235e567E7A1e1E1f59ddD79 to finish the payment. 

If you want to accept EOS

```shell
 $ curl -d '{"currency":"EOS", "amount":0.001, "customerid":"client1245", "webhookurl":":9090/123", "expiredafter":5}' -H "Content-Type: application/json" 127.0.0.1:8080/charges
```
```json
{
	"Id":2,
	"Currency":"EOS",
	"Amount":0.001,
	"Customerid":"client1245",
	"Webhookurl":":9090/123",
	"Expired_after":5,
	"Paymentmethod":{
		"Name":"EOS",
		"PaymentAddress":"",
		"PaymentAccount":"eoswithmixin",
		"PaymentMemo":"a01a148f234ea8be0229a4422d21e7f3",
		"Priceinusd":"4.63264861",
		"Priceinbtc":"0.00040277"
	},
	"Receivedamount":0,
	"Paidstatus":0
}
```
Client need to tranfser 0.001 EOS to account eoswithmixin, and MUST fill memo a01a148f234ea8be0229a4422d21e7f3 to finish the payment. 
![](https://github.com/myrual/mixin-network-snapshot-golang/raw/master/EOS_pay.jpg)

If you want to accept Stellar XLM
```shell
curl -d '{"currency":"XLM", "amount":0.001, "customerid":"client1245", "webhookurl":":9090/123", "expiredafter":5}' -H "Content-Type: application/json" 127.0.0.1:8080/charges
```
```json
{
	"Id":3,
	"Currency":"XLM",
	"Amount":0.001,
	"Customerid":"client1245",
	"Webhookurl":":9090/123",
	"Expired_after":5,
	"Paymentmethod":{
		"Name":"XLM",
		"PaymentAddress":"",
		"PaymentAccount":"GD77JOIFC622O5HXU446VIKGR5A5HMSTAUKO2FSN5CIVWPHXDBGIAG7Y",
		"PaymentMemo":"45da67ad857c907a",
		"Priceinusd":"0.08866487",
		"Priceinbtc":"0.00000769"
	},
	"Receivedamount":0,
	"Paidstatus":0
}
```
Client need to tranfser 0.001 XLM to account GD77JOIFC622O5HXU446VIKGR5A5HMSTAUKO2FSN5CIVWPHXDBGIAG7Y, and MUST fill memo 45da67ad857c907a to finish the payment. 
![](https://github.com/myrual/mixin-network-snapshot-golang/raw/master/XLM_pay.jpg)

There are two types of payment method:
1. Bitcoin/Ethereum style: PaymentAddress is not empty, PaymentAccount and PaymentMemo are all empty. You just  show Ethererum Name and PaymentAddress to your clients, they just need to transfer token to the address. In this example, show asset name ETH, payment address 0x365DA43BC7B22CD4334c3f35eD189C8357D4bEd6 and payment amount to your client.
2. EOS/Stellar style: PaymentAddress is empty, PaymentAccount and PaymentMemo are not empty. You need to show Asset Name and both of PaymentAccount and PaymentMemo to user, and remind user need to input BOTH of PaymentAccount and PaymentMemo. Transfer asset to PaymentAccount without memo is a common mistake, and it can not be reverted because current Mixin Network limitation. In this example, show asset name EOS, payment account eoswithmixin , payment memo 302c37ebff05ccf09dd7296053d1924a.

Asset current price in USD and Bitcoin is inside payment record, so developer can calculate how many asset client should transfer to the address or account.

```json
{
	"Priceinusd":"310.40105841",
	"Priceinbtc":"0.02374051"
}	
```

Currency list

|Currency| Explain | introduction|
|-| - | - |
|EOS|EOS.io main chain token|-|
|XLM|Stellar main chain token|-|
|BTC|Bitcoin|-|
|UDT|Tether USD|Running on Bitcoin instead of Ethereum|
|XRP|Ripple|-|
|LTC|Litecoin|-|

#### Query payment status
fetch the payment status by visit localhost:8080/charges with parameter charge_id

Example:
```shell
 curl -X GET 'http://localhost:8080/charges?charge_id=3'

```

Response will be similar to following if payment is not yet confirmed
```json
{
	"Id":3,
	"Currency":"ETH",
	"Amount":0.001,
	"Customerid":"client1245",
	"Webhookurl":":9090/123",
	"Expired_after":60,
	"Paymentmethod":{
		"Name":"ETH",
		"PaymentAddress":"0x130D3e6655f073e33235e567E7A1e1E1f59ddD79",
		"PaymentAccount":"",
		"PaymentMemo":"",
		"Priceinusd":"310.40105841",
		"Priceinbtc":"0.02374051"
		},
	"Receivedamount":0,
	"Paidstatus":0}
}
```

Response will be similar to following if payment is already confirmed
```json
{
	"Id":3,
	"Currency":"ETH",
	"Amount":0.001,
	"Customerid":"client1245",
	"Webhookurl":":9090/123",
	"Expired_after":60,
	"Paymentmethod":{
		"Name":"ETH",
		"PaymentAddress":"0x130D3e6655f073e33235e567E7A1e1E1f59ddD79",
		"PaymentAccount":"",
		"PaymentMemo":"",
		"Priceinusd":"309.75108846",
		"Priceinbtc":"0.02369282"
	},
	"Receivedamount":0.002,
	"Paidstatus":2
}
```

Paid status

|value | description|
|--|--|
|0| not yet paid|
|1| partial paid|
|2| paid|
|3| over paid|

The payment address is a deposit address in cryptocurrency world, so user can deposit any amount. 


#### payment notification webhook
The program will visit webhook url when user paid and confirmed by program. 
```json
"http://127.0.0.1"+webhookurl
```
The http visit method is POST, json body parameter is following
```json
{
	"Id":3,
	"Currency":"ETH",
	"Amount":0.001,
	"Customerid":"client1245",
	"Webhookurl":":9090/123",
	"Expired_after":60,
	"Paymentmethod":{
		"Name":"ETH",
		"PaymentAddress":"0x130D3e6655f073e33235e567E7A1e1E1f59ddD79",
		"PaymentAccount":"",
		"PaymentMemo":"",
		"Priceinusd":"309.75108846",
		"Priceinbtc":"0.02369282"
	},
	"Receivedamount":0.0021,
	"Paidstatus":2
}
```
Developer can know when, which asset is paid by client, and what's the payment value in USD and Bitcoin.


### How did developer get all asset?
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

#### What is confirmation time, why need to care about it?
A cryptocurrency transaction created by your client need to be confirmed by network, Bitcoin network need long time to confirm, other blockchain need less time.

### What kind of currency can be supported
All asset supported by Mixin Network:
BTC, USDT, BCH, ETH and ERC20, ETC, EOS and token issue on EOS, DASH, Litecoin, Doge, Horizen, MGD, NEM, XRP, XLM, TRON and TRC10, Zcash. 

### Recommend Currency
Three kind of currency : ETH, EOS, XLM are accepted in code.

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
	default_asset_id_group := []string{XLM_ASSET_ID, EOS_ASSET_ID, ETH_ASSET_ID}
```


TO BE DONE:
1. All asset can be withdrawed to developer's cold wallet.
2. One type asset can be exchanged to USDT or Bitcoin automatically through DEX.
3. Support Mixin Messenger User to pay.
4. ~~Latest USD price for every asset~~ Implemented in commit 8a634e23254e4841c2a9c3114b3eb847d46f55fc
