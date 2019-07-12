# 数字货币收款插件
轻松，免费，安全的数字货币收款方案
* 无需搭建比特币/以太坊/EOS全节点(每一个都需要几百G空间)
* 无手续费，你的程序你做主
* 所有收到的钱实时自动转移到开发者个人账户，即使被拖库也没钱可盗。

开发者访问本地 http 接口，向用户展示付款方法，用户付款后程序会访问本地回调URL

步骤:
### 1. 创建一个Mixin Messenger账户.
访问 https://mixin.one/messenger 下载对应手机端App。

中国大陆用户可以访问 https://a.app.qq.com/o/simple.jsp?pkgname=one.mixin.messenger  下载

### 2. 激活开发者账号
登陆 https://developer.mixin.one ，用App扫码登录

这个 [教程](https://mixin-network.gitbook.io/mixin-network/mixin-messenger-app/create-bot-account)对于新开发者很有用。

### Clone, build, run
```shell
git clone https://github.com/myrual/mixin-network-snapshot-golang
cd mixin-network-snapshot-golang
```

2. 编辑一部分配置信息
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
3. 编译
```shell
go build mixin_snap.go
```
4. 运行
```shell
./mixin_snap
```

5. 数据库
同一目录下会生成一个test.db 的sqlite3文件。


## 如何使用
#### 获取数字资产当前价格信息，因此可以计算客户应该付多少数字资产
```shell
curl -X GET 'http://localhost:8080/assetsprice'
```

价格结果如下，其中Full Name是该币种全名，Symbol是在交易所和钱包的缩写符号，USDPrice是当前资产美元价格，BTCPrice同理。
```json
[
	{"Fullname":"Stellar","Symbol":"XLM","USDPrice":0.10357796,"BTCPrice":0.00000889,"Assetid":"56e63c06-b506-4ec5-885a-4a5ac17b83c1"},
	{"Fullname":"EOS","Symbol":"EOS","USDPrice":5.96024263,"BTCPrice":0.00051165,"Assetid":"6cfe566e-4aad-470b-8c9a-2fd35b49c68d"},
	{"Fullname":"Ether","Symbol":"ETH","USDPrice":294.61322131,"BTCPrice":0.02529107,"Assetid":"43d61dcd-e413-450d-80b8-101d5e903357"}
]
```
如果订单价值1美金，那么根据资产价格可以知道客户需要 10 XLM, 或者 0.17 EOS。
#### 创建支付请求
用POST方法访问 localhost:8080/charges，参数如下 

POST /charges

|Attributes| type | description|
|--|--|--|
|currency| String | Currency code associated with the amount.  Only EOS/XLM/ETH is supported currently|
|amount| Float64 | Positive float|
|customerid| String | This field is optional and can be used to attach an identifier of your choice to the charge. Must not exceed 64 characters|
|webhookurl| String | program will visit localhost+webhook when user pay enough currency before charge is expired |
|expiredafter| uint | the webhook will be expired after xx minutes. User can pay to an expired charge , program keep income record and will transfer asset to admin account|

举例: 需要让客户 "client1245" 支付 0.001 ETH, 60分钟内支付完成之后访问用 POST 访问 localhost:9090/123。
```shell
curl -d '{"currency":"ETH", "amount":0.001, "customerid":"client1245", "webhookurl":":9090/123", "expiredafter":60}' -H "Content-Type: application/json" 127.0.0.1:8080/charges
```

这条指令返回结果
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
客户需要向以太坊地址 0x130D3e6655f073e33235e567E7A1e1E1f59ddD79 支付0.001 ETH来完成支付。 

如果想收EOS
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
客户需要向EOS账户 eoswithmixin 支付 0.001 EOS, 并且必须填写支付备注 a01a148f234ea8be0229a4422d21e7f3。 
![](https://github.com/myrual/mixin-network-snapshot-golang/raw/master/EOS_pay.jpg)

如果想收XLM
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
客户需要向Stellar账户 GD77JOIFC622O5HXU446VIKGR5A5HMSTAUKO2FSN5CIVWPHXDBGIAG7Y 支付 0.001 XLM, 并且必须填写支付备注 45da67ad857c907a 
![](https://github.com/myrual/mixin-network-snapshot-golang/raw/master/XLM_pay.jpg)


Payment_method里面有两种类型的支付：
1. 比特币/以太坊: PaymentAddress 不是空，PaymentAccount 和 PaymentMemo是空。这种情况下，你只需要给用户展示资产名字 以太坊和PaymentAddress，客户只需要向以太坊地址付款。在这个例子里面，向用户展示资产名称 ETH，以及收款地址 0x365DA43BC7B22CD4334c3f35eD189C8357D4bEd6，以及你期望的以太坊数量。
2. EOS/行星 : PaymentAddress 是空, PaymentAccount 和 PaymentMemo 都有内容。这种情况下，你需要给用户展示资产名字，收款账户和收款备注，并且严肃的提醒用户同时填写收款账户和收款备注，客户如果忘记填写备注，会导致不能到账，而且无法退款。

Payment_method的记录内容里面有该资产当前的美元价格和比特币价格，开发者可以根据订单的美元价格来计算客户应该支付多少数字货币。
```json
{"Priceinusd":"0.10472789","Priceinbtc":"0.00000925"}
```

支持的货币列表

|Currency| 说明 | 介绍|
|-| - | - |
|EOS|EOS.io 主网token|-|
|XLM|Stellar 主网token|-|
|BTC|比特币|-|
|UDT|Tether USD|基于比特币的USDT，不是ERC20的代币|
|XRP|锐波币|-|
|LTC|莱特币|-|

#### 检查收款状态
访问 localhost:8080/charges, 带有参数charge_id

例子
```shell
 curl -X GET 'http://localhost:8080/charges?charge_id=3'

```

如果客户还没有支付，结果如下
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

如果客户已经支付了，那么结果如下
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

支付状态 Paidstatus的解释

|值 | 解释|
|--|--|
|0| 还没有支付|
|1| 支付不足|
|2| 支付完毕|
|3| 支付超过需求|

payment_records 有支付信息. 其中一个支付信息如下

#### 回掉URL
用户支付完毕后，程序会访问本地+webook url
```json
"http://127.0.0.1"+webhookurl
```
http 方法为POST，body 参数如下
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


### 所有资产都属于开发者自己么？
1. 所有资产都会自动被转移到你指定的账户, 免手续费，1秒到账。
2. 在该程序重启，或者意外退出之后，你可以手动发指令要求程序把所有资产都转移到你指定的账户。

例子：
```shell
curl -X POST -H "Content-Type: application/json" 127.0.0.1:8080/moneygohome
```

结果如下
```json
total 20 account will send all balance to admin
```

### 支付的确认时间
1. EOS: 3 分钟
2. Stellar: 2 分钟
3. Bitcoin/USDT: 60 分钟
4. Litecoin/Ethererum/DOGE: 120 分钟

#### 什么是确认时间？为什么要关心确认时间？
数字货币从用户发起转账请求，到收款方确认这笔付款不能回滚需要一点时间，比特币需要的时间长，其他需要的时间短一点。

### 支持哪些资产
理论上Mixin Network支持的都可以接受。现在支持
BTC, USDT, BCH, 以太坊和 ERC20, ETC, EOS 以及EOS上发的token, DASH, Litecoin, Doge, Horizen, MGD, NEM, XRP, XLM, 波场和波场上发的TRC10, Zcash. 

### 目前的代码库默认支持的资产
现在代码里面默认支持的资产是EOS和恒星，因为他们都可以3分钟完成支付确认。

想要支持更多的币，把对应资产的变量放到 default_asset_id_group 里面就可以.
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


下一步的开发任务:
1. 所有的资产可以自动提取到开发者自己的冷钱包，而不是只能转移到Mixin Messenger账户。
2. 可以把收到的资产通过去中心化交易所自动转换成USDT或者比特币。
3. 支持Mixin Messenger用户付款。
4. ~~可以提供资产对应的美元价格~~ 在commit 8a634e23254e4841c2a9c3114b3eb847d46f55fc 中已经完成。
