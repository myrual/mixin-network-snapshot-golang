# 数字货币收款插件
不需要了解比特币，EOS的全部API，不需要搭建全节点就可以接收数字货币付款.

这个程序是一个全集成方案，开发者只需要通过http api就可以调用接口，把付款方式展示给消费者，程序会自动访问回掉URL。

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
#### 如何接受数字货币付款
为了接受比特币，EOS支付，开发者需要用http POST方法访问 localhost:8080/payment，参数放在body里面。 body里面应该有 唯一标示字符串，回掉URL，以及回掉过期时间。 唯一标示字符串可以是任意字符，uuid也可以。 程序收到用户的付款之后会用http post方法访问回掉本机Callback url。
回掉机制有有效期，过了有效期，回掉机制会停止。如果回掉有效期参数为60， 那么回掉机制会在60分钟后过期。
curl例子
```shell
curl -d '{"reqid":"value8", "callback":":9090/", "expiredafter":60}' -H "Content-Type: application/json" 127.0.0.1:8080/payment
```
这条指令指示程序为value8创建一个支付地址，如果在60分钟内收到用户支付，那么程序会http post方式访问 localhost:9090/，同时带有参数。

这条指令的返回结果是
```json
{"Reqid":"value8","Payment_methods":[{"Name":"XLM","PaymentAddress":"","PaymentAccount":"GD77JOIFC622O5HXU446VIKGR5A5HMSTAUKO2FSN5CIVWPHXDBGIAG7Y","PaymentMemo":"3f8db42022b5bc32","Priceinusd":"0.10472789","Priceinbtc":"0.00000925"},{"Name":"EOS","PaymentAddress":"","PaymentAccount":"eoswithmixin","PaymentMemo":"302c37ebff05ccf09dd7296053d1924a","Priceinusd":"5.9436916","Priceinbtc":"0.00052505"},{"Name":"ETH","PaymentAddress":"0x365DA43BC7B22CD4334c3f35eD189C8357D4bEd6","PaymentAccount":"","PaymentMemo":"","Priceinusd":"295.86024062","Priceinbtc":"0.02613571"}],"Payment_records":null,"Balance":null}
```
Payment_methods里面的结果是给客户看的，这个例子有三个支付方法。

有两种风格的支付：
1. 比特币/以太坊: PaymentAddress 不是空，PaymentAccount 和 PaymentMemo是空。这种情况下，你只需要给用户展示资产名字 以太坊和PaymentAddress，客户只需要向以太坊地址付款。
2. EOS/行星 : PaymentAddress 是空, PaymentAccount 和 PaymentMemo 都有内容。这种情况下，你需要给用户展示资产名字，收款账户和收款备注，并且严肃的提醒用户同时填写收款账户和收款备注，客户如果忘记填写备注，会导致不能到账，而且无法退款。

Payment_methods的记录内容里面有该资产当前的美元价格和比特币价格。
```json
{"Priceinusd":"0.10472789","Priceinbtc":"0.00000925"}
```

#### 检查收款状态
通过参数 reqid 访问 localhost:8080/payment 可以查询收款状态和记录。

例子:
```shell
curl -X GET 'http://localhost:8080/payment?reqid=value8'
```

如果客户还没有支付，那么结果是这样的
```json
{"Reqid":"value8","Payment_methods":[{"Name":"XLM","PaymentAddress":"","PaymentAccount":"GD77JOIFC622O5HXU446VIKGR5A5HMSTAUKO2FSN5CIVWPHXDBGIAG7Y","PaymentMemo":"3f8db42022b5bc32","Priceinusd":"0.10472789","Priceinbtc":"0.00000925"},{"Name":"EOS","PaymentAddress":"","PaymentAccount":"eoswithmixin","PaymentMemo":"302c37ebff05ccf09dd7296053d1924a","Priceinusd":"5.9436916","Priceinbtc":"0.00052505"},{"Name":"ETH","PaymentAddress":"0x365DA43BC7B22CD4334c3f35eD189C8357D4bEd6","PaymentAccount":"","PaymentMemo":"","Priceinusd":"295.86024062","Priceinbtc":"0.02613571"}],"Payment_records":null,"Balance":null}
```
paymnet_records 是空

如果客户已经支付了，结果是这样的。

```json
{"Reqid":"value8","Payment_methods":[{"Name":"XLM","PaymentAddress":"","PaymentAccount":"GD77JOIFC622O5HXU446VIKGR5A5HMSTAUKO2FSN5CIVWPHXDBGIAG7Y","PaymentMemo":"3f8db42022b5bc32","Priceinusd":"0.10472789","Priceinbtc":"0.00000925"},{"Name":"EOS","PaymentAddress":"","PaymentAccount":"eoswithmixin","PaymentMemo":"302c37ebff05ccf09dd7296053d1924a","Priceinusd":"5.9436916","Priceinbtc":"0.00052505"},{"Name":"ETH","PaymentAddress":"0x365DA43BC7B22CD4334c3f35eD189C8357D4bEd6","PaymentAccount":"","PaymentMemo":"","Priceinusd":"295.86024062","Priceinbtc":"0.02613571"}],"Payment_records":[{"Amount":"0.1","AssetId":"","created_at":"2019-06-20T02:00:39.650472961Z","snapshot_id":"570233aa-3c91-45cd-a6ec-0e9724165300"},{"Amount":"0.01","AssetId":"6cfe566e-4aad-470b-8c9a-2fd35b49c68d","created_at":"2019-06-20T02:33:50.152539755Z","snapshot_id":"88859d4d-5bee-4fb5-aef6-ac01dc3a43c6"},{"Amount":"0.01","AssetId":"6cfe566e-4aad-470b-8c9a-2fd35b49c68d","created_at":"2019-06-20T02:37:05.870885973Z","snapshot_id":"6530f455-3238-491a-a9c5-bbcb52bcc306"},{"Amount":"0.001","AssetId":"6cfe566e-4aad-470b-8c9a-2fd35b49c68d","created_at":"2019-06-20T02:40:53.251365044Z","snapshot_id":"f2c8a751-3d30-472e-bf76-924787f341b9"},{"Amount":"0.001","AssetId":"6cfe566e-4aad-470b-8c9a-2fd35b49c68d","created_at":"2019-06-20T02:59:28.854380284Z","snapshot_id":"3ebfd5a3-bd29-4e32-bd06-2506bee3da99"},{"Amount":"-0.122","AssetId":"6cfe566e-4aad-470b-8c9a-2fd35b49c68d","created_at":"2019-06-20T03:00:17.249302744Z","snapshot_id":"0bfe6f6b-1ff8-4144-9786-52d6a6459b19"}],"Balance":null}
```
payment_records 有支付信息. 其中一个支付信息如下
```json
{"Amount":"0.01","AssetId":"6cfe566e-4aad-470b-8c9a-2fd35b49c68d",
"created_at":"2019-06-20T02:37:05.870885973Z","snapshot_id":"6530f455-3238-491a-a9c5-bbcb52bcc306"}
```
这是一条来自客户的支付: 
* 数量 0.01
* 资产id 6cfe566e-4aad-470b-8c9a-2fd35b49c68d，是 EOS 主网token
* 支付生成于 UTC 2019-06-20T02:37:05.870885973
* 该支付在Mixin Network内的唯一标示号 6530f455-3238-491a-a9c5-bbcb52bcc306，你可以在浏览器里面验证这笔交易 https://mixin.one/snapshots/6530f455-3238-491a-a9c5-bbcb52bcc306

#### 回掉URL
在有效期内收到用户付款，程序会访问本地的回掉URL。
```json
"http://127.0.0.1"+callbackurl
```
http访问方法是POST，参数在body里面，例子如下
```json
{"Reqid":"value8","Callbackurl":":9090/","Paymentrecord":{"Amount":"0.01","AssetId":"56e63c06-b506-4ec5-885a-4a5ac17b83c1","created_at":"2019-06-20T07:33:06.445471337Z","snapshot_id":"a6603374-509b-4015-a192-c63bfa8def5f"}}
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

什么是确认时间？大部分数字货币从用户发起转账请求，到收款方确认这笔付款不能回滚需要一点时间。

为什么这么长？这是Mixin Network 本身的设定，你现在改不了。

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
