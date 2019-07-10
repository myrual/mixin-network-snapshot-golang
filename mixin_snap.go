package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"time"

	messenger "github.com/MooooonStar/mixin-sdk-go/messenger"
	mixin "github.com/MooooonStar/mixin-sdk-go/network"
	"github.com/gofrs/uuid"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

const (
	userid      = "3c5fd587-5ac3-4fb6-b294-423ba3473f7d"
	sessionid   = "82365995-4d86-4676-a632-4354f1a04954"
	private_key = `-----BEGIN RSA PRIVATE KEY-----
	MIICXgIBAAKBgQCIJJ+PfATaR3UDegG13fs0Rgx21+jpEB/aVOqrDNQUmuY5k/MC
	GRicNeIa+J2gT6lzPfWDL8z5IABnlc6EgI4BKM1ok0UgeQy/DWvDAmik47vEMpjU
	JjKwgwyarQA/yWl1R/9djt/cVOwxekOXJqNAjsfGzjArgUAzBZ3OYNr15wIDAQAB
	AoGAL5kmVCMXCz3KcmG4sV4f0qHe/7nzC3EAwfPIa+87Qsz5Sw4n+wbNLOhF2gos
	Cf1wEAOMj8YpkrwWiCC/KGJNwyGvCMWs166oA0KFPeMhpKMlqXE+DpNJQzit3mb1
	nZf18MXSMHw2u8LFfMDjJb2QZTOMVvWjYKnVg1AchfcPbeECQQDtW26N9lMTlzXO
	MxrtheenHuSEgm55wGRhKquxDpXbJNcalHAnIvfd5tGbuciMbVlEk01MF7ggEWID
	IhSHXhvRAkEAktYRV2XetVvF13orMGMPag9akdCjCpnA8Yxn11S0HSGsM6nkrMKl
	i7t+qDTjNUODsnOMzb34iJCNtaaHGL48NwJBAKiQ6Yfaaw+bsLObKcGL+oN+dg4B
	T5IZ52/2TO62jAiRNk6DIs84j03BUhVFML9CHUaNUjT7F2F21uOgvXGRjTECQQCG
	Uq2qddY1sa5aX7gCm5wOOd1wZpu/pseKMBcONL5Pp+4PlOtL3wPxv6Mt3LO8lfZz
	2KCF1bL1usbn1V7gk6YhAkEAhegwxueVDfayjuOSX0yHKY6VAECELpCqpNy1YEon
	9oKozCDe1Kg/ZLWGMP9cIayb/eyUjJeyUPP8dR2rfhFBkA==
-----END RSA PRIVATE KEY-----`

	ADMIN_MessengerID = "31367"
)

type Snapshot struct {
	Amount string `json:"amount"`
	Asset  struct {
		AssetId string `json:"asset_id"`
	} `json:"asset"`
	AssetId    string    `json:"asset_id"`
	CreatedAt  time.Time `json:"created_at"`
	SnapshotId string    `json:"snapshot_id"`
	Source     string    `json:"source"`
	Type       string    `json:"type"`
	//only available when http request include correct token
	UserId          string `json:"user_id"`
	TraceId         string `json:"trace_id"`
	OpponentId      string `json:"opponent_id"`
	Sender          string `json:"sender"`
	Data            string `json:"data"`
	Transactionhash string `json:"transaction_hash"`
}
type Payment_Record struct {
	Amount     string
	AssetId    string
	CreatedAt  time.Time `json:"created_at"`
	SnapshotId string    `json:"snapshot_id"`
}

type Profile struct {
	CreatedAt time.Time `json:"created_at"`
}

type DepositAddressResonse struct {
	PublicKey    string `json:"public_key"`
	AccountName  string `json:"account_name"`
	AccountTag   string `json:"account_tag"`
	IconURL      string `json:"icon_url"`
	Confirmblock uint   `json:"confirmations"`
	Symbol       string `json:"symbol"`
	Name         string `json:"name"`
	Chainid      string `json:"chain_id"`
	Assetkey     string `json:"asset_key"`
	Assetid      string `json:"asset_id"`
	Priceusd     string `json:"price_usd"`
	Pricebtc     string `json:"price_btc"`
}
type DepositNetResponse struct {
	Error         error
	Accountid     string
	Assetid       string
	MixinResponse MixinDepositResponse
}
type MixinDepositResponse struct {
	Data  *DepositAddressResonse `json:"data"`
	Error string                 `json:"error"`
}
type Snapshotindb struct {
	gorm.Model
	SnapshotId    string `gorm:"primary_key"`
	Amount        string
	AssetId       string `gorm:"index"`
	Source        string `gorm:"index"`
	SnapCreatedAt time.Time
	UserId        string `gorm:"index"`
	TraceId       string
	OpponentId    string
	Data          string
	OPType        string
}

type MixinAccountindb struct {
	gorm.Model
	Userid        string `gorm:"primary_key"`
	Sessionid     string
	Pintoken      string
	Privatekey    string
	Pin           string
	ClientReqid   uint
	Utccreated_at time.Time
}

type DepositAddressindb struct {
	gorm.Model
	Accountrecord_id uint
	Assetid          string
	Publicaddress    string
	Accountname      string
	Accounttag       string
	Iconurl          string
	Confirmblock     uint
}

type AssetInformationindb struct {
	gorm.Model
	Assetid  string
	Chainid  string
	Symbol   string
	Name     string
	Assetkey string
}

type MessengerUserindb struct {
	gorm.Model
	Messengerid string `gorm:"primary_key"`
	Uuid        string
}

type Assetpriceindb struct {
	gorm.Model
	Assetid    string `gorm:"primary_key"`
	Priceinusd string
	Priceinbtc string
}
type ClientReq struct {
	gorm.Model
	Reqid          string
	Callbackurl    string
	MixinAccountid uint
	Callbackfired  bool
}
type CallbackRespone struct {
	Reqid         string
	Callbackurl   string
	Paymentrecord Payment_Record
	Valueinusd    float64
	Valueinbtc    float64
}

type Searchtaskindb struct {
	gorm.Model
	Starttime         time.Time
	Endtime           time.Time
	Yesterday2today   bool
	Assetid           string
	Ongoing           bool
	Userid            string
	Includesubaccount bool
	Taskexpired_at    time.Time
}
type Searchtaskinram struct {
	Starttime         time.Time
	Endtime           time.Time
	Taskexpired_at    time.Time
	Yesterday2today   bool
	Assetid           string
	Ongoing           bool
	Userid            string
	Sessionid         string
	Privatekey        string
	Includesubaccount bool
}

type BotConfig struct {
	user_id     string
	session_id  string
	private_key string
}

type SnapNetResponse struct {
	Error        error
	MixinRespone MixinResponse
}

type MixinResponse struct {
	Data  []*Snapshot `json:"data"`
	Error string      `json:"error"`
}

type TransferNetRespone struct {
	TransferRes TransferResponse
	Error       error
}
type TransferResponse struct {
	Data  Transfer `json:"data"`
	Error string   `json:"error"`
}
type Transfer struct {
	Optype         string `json:"type"`
	Snapshotid     string `json:"snapshot_id"`
	OpponentId     string `json:"opponent_id"`
	Assetid        string `json:"asset_id"`
	Amount         string `json:"amount"`
	Memo           string `json:"memo"`
	Snap_createdat string `json:"created_at"`
}

type BalanceNetResponse struct {
	Balance BalanceResponse
	Error   error
}
type BalanceResponse struct {
	Data  []*Asset `json:"data"`
	Error string   `json:"error"`
}

type Asset struct {
	Optype  string `json:"type"`
	Assetid string `json:"asset_id"`
	Balance string `json:"balance"`
}
type ProfileResponse struct {
	Data  *Profile `json:"data"`
	Error string   `json:"error"`
}
type Searchtask struct {
	start_t            time.Time
	end_t              time.Time
	task_expired_after time.Time
	yesterday2today    bool
	max_len            int
	asset_id           string
	ongoing            bool
	userconfig         BotConfig
	includesubaccount  bool
}

type Searchprogress struct {
	search_task Searchtask
	Error       error
}

type PaymentReqhttp struct {
	Reqid         string `json:"reqid"`
	Callback      string `json:"callback"`
	Expired_after uint32 `json:"expiredafter"`
}

type PaymentReq struct {
	Method        string
	Reqid         string
	Callback      string
	Expired_after uint32
	Res_c         chan PaymentRes
}
type OPReq struct {
	op_code string
	Res_c   chan []byte
}
type PaymentMethod struct {
	Name           string
	PaymentAddress string
	PaymentAccount string
	PaymentMemo    string
	Priceinusd     string
	Priceinbtc     string
}
type PaymentRes struct {
	Reqid           string
	Payment_methods []PaymentMethod
	Payment_records []Payment_Record
	Balance         []Asset
	ReceivedinUSD   float64
	ReceivedinBTC   float64
}
type ChargeReqhttp struct {
	Currency      string  `json:"currency"`
	Amount        float64 `json:"amount"`
	Customerid    string  `json:"customerid"`
	Webhookurl    string  `json:"webhookurl"`
	Expired_after uint32  `json:"expiredafter"`
}
type ChargeResponse struct {
	Currency       string
	Amount         float64
	Customerid     string
	Webhookurl     string
	Expired_after  uint32
	Paymentmethod  PaymentMethod
	Receivedamount float64
	Paidstatus     uint32
}

type ChargeReq struct {
	Method        string
	Currency      string
	Amount        float64
	Customerid    string
	Webhookurl    string
	Expired_after uint32
	Res_c         chan ChargeResponse
}
type ChargeRecordindb struct {
	gorm.Model
	Currency       string
	Amount         float64
	Customerid     string
	Webhookurl     string
	Expiredafter   uint32
	MixinAccountid uint
}
type ChargeRelationWithMixinAccountindb struct {
	gorm.Model
	Chargerecordid uint
	Mixinaccountid uint
}

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

	PREDEFINE_PIN             = "198435"
	PREDEFINE_NAME            = "tom"
	scan_interval_in_seconds  = 5
	op_all_money_go_home      = "allmoneygohome"
	op_all_snaps              = "allsnaps"
	assets_price              = "assetsprice"
	scan_stop_after_n_minutes = 240
	local_web_port            = ":8080"
)

func read_asset_deposit_address(asset_id string, user_id string, session_id string, private_key string, deposit_c chan DepositNetResponse) {
	result, err := mixin.Deposit(asset_id, user_id, session_id, private_key)

	if err != nil {
		deposit_c <- DepositNetResponse{
			Error:     err,
			Accountid: user_id,
			Assetid:   asset_id,
		}
		return
	}

	var resp MixinDepositResponse
	err = json.Unmarshal(result, &resp)

	if err != nil {
		deposit_c <- DepositNetResponse{
			Error:     err,
			Accountid: user_id,
			Assetid:   asset_id,
		}
	}
	if resp.Error != "" {
		log.Println("Server return error", resp.Error, " for req:")
		return
	}

	deposit_c <- DepositNetResponse{
		Accountid:     user_id,
		Assetid:       asset_id,
		MixinResponse: resp,
	}
}

//read snapshot related to the account or account created by the account
//given asset id and kick off time:
//    the routine will read and filter snapshot endless,
//    push snap result into channel
//    and progress to another channel
//given asset id and kick off time and end time:
//    the routine will read and filter snapshot between the kick off and end time,
//    filter snapshot and push data to channel, and progress to another channel

func read_useruuid_from(user_id string, session_id string, private_key string, messengerid string) string {
	botUser := messenger.NewMessenger(user_id, session_id, private_key)
	ctx := context.Background()
	user, err := botUser.SearchUser(ctx, ADMIN_MessengerID)

	if err != nil {
		log.Println(err)
		return ""
	}

	return user.UserId
}

func read_bot_created_time(user_id string, session_id string, private_key string) time.Time {
	botUser := mixin.NewUser(user_id, session_id, private_key, "")
	profile, err := botUser.ReadProfile()

	if err != nil {
		return time.Date(0, 0, 0, 0, 0, 0, 0, time.UTC)
	}

	var resp ProfileResponse
	err = json.Unmarshal(profile, &resp)

	if err != nil {
		return time.Date(0, 0, 0, 0, 0, 0, 0, time.UTC)
	}
	if resp.Error != "" {
		return time.Date(0, 0, 0, 0, 0, 0, 0, time.UTC)
	}
	return resp.Data.CreatedAt
}

func read_snap_to_future(req_task Searchtask, result_chan chan *Snapshot, in_progress_c chan Searchprogress) {
	for {

		if req_task.end_t.IsZero() == false && time.Now().After(req_task.end_t) {
			log.Println("Exit task because user set end time and it is passed now ")
			p := Searchprogress{
				search_task: req_task,
			}
			p.search_task.ongoing = false
			in_progress_c <- p
			return
		}
		var snaps []byte
		var err error

		snaps, err = mixin.NetworkSnapshots(req_task.asset_id, req_task.start_t, "ASC", req_task.max_len, req_task.userconfig.user_id, req_task.userconfig.session_id, req_task.userconfig.private_key)

		if err != nil {
			in_progress_c <- Searchprogress{
				Error: err,
			}
			continue
		}

		var resp MixinResponse
		err = json.Unmarshal(snaps, &resp)

		if err != nil {
			in_progress_c <- Searchprogress{
				Error: err,
			}
			continue
		}
		if resp.Error != "" {
			log.Fatal("Server return error", resp.Error, " for req:", req_task.asset_id, " start ", req_task.start_t)
			return
		}
		len_of_snap := len(resp.Data)
		for _, v := range resp.Data {

			if v.UserId != "" {
				result_chan <- v
			}
		}
		if len_of_snap == 0 {
			p := Searchprogress{
				search_task: req_task,
			}
			in_progress_c <- p
			//nothing is searched, wait
			time.Sleep(scan_interval_in_seconds * time.Second)
			continue
		} else {
			last_element := resp.Data[len(resp.Data)-1]
			req_task.start_t = last_element.CreatedAt
			p := Searchprogress{
				search_task: req_task,
			}
			p.search_task.start_t = last_element.CreatedAt
			p.search_task.ongoing = true
			in_progress_c <- p
			if len_of_snap < req_task.max_len {
				time.Sleep(scan_interval_in_seconds * time.Second)
			}
		}
	}
}
func makePaymentHandle(input chan PaymentReq) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			keys, ok := r.URL.Query()["reqid"]
			if ok != true || len(keys[0]) < 1 {
				io.WriteString(w, "Missing parameter reqid!\n")
				return
			}
			payment_res_c := make(chan PaymentRes, 1)
			req := PaymentReq{
				Method: "GET",
				Reqid:  keys[0],
				Res_c:  payment_res_c,
			}
			input <- req
			v := <-payment_res_c
			b, jserr := json.Marshal(v)
			if jserr != nil {
				log.Println(jserr)
			} else {
				w.Write(b)
			}
		case "POST":
			d := json.NewDecoder(r.Body)
			var p PaymentReqhttp
			errjs := d.Decode(&p)
			if errjs != nil {
				http.Error(w, errjs.Error(), http.StatusInternalServerError)
			}
			payment_res_c := make(chan PaymentRes, 1)
			req := PaymentReq{
				Reqid:         p.Reqid,
				Callback:      p.Callback,
				Res_c:         payment_res_c,
				Expired_after: p.Expired_after,
			}
			input <- req
			v := <-payment_res_c
			b, jserr := json.Marshal(v)
			if jserr != nil {
				log.Println(jserr)
			} else {
				w.Write(b)
			}
		default:
			io.WriteString(w, "Wrong!\n")
		}
	}
}
func makeChargeHandle(input chan ChargeReq) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			keys, ok := r.URL.Query()["customerid"]
			if ok != true || len(keys[0]) < 1 {
				io.WriteString(w, "Missing parameter customerid!\n")
				return
			}
			charge_res_c := make(chan ChargeResponse, 1)
			req := ChargeReq{
				Method:     "GET",
				Customerid: keys[0],
				Res_c:      charge_res_c,
			}
			input <- req
			v := <-charge_res_c
			b, jserr := json.Marshal(v)
			if jserr != nil {
				log.Println(jserr)
			} else {
				w.Write(b)
			}
		case "POST":
			d := json.NewDecoder(r.Body)
			var p ChargeReqhttp
			errjs := d.Decode(&p)
			if errjs != nil {
				http.Error(w, errjs.Error(), http.StatusInternalServerError)
			}
			charge_res_c := make(chan ChargeResponse, 1)
			req := ChargeReq{
				Currency:      p.Currency,
				Amount:        p.Amount,
				Customerid:    p.Customerid,
				Webhookurl:    p.Webhookurl,
				Expired_after: p.Expired_after,
				Res_c:         charge_res_c,
			}
			input <- req
			v := <-charge_res_c
			b, jserr := json.Marshal(v)
			if jserr != nil {
				log.Println(jserr)
			} else {
				w.Write(b)
			}
		default:
			io.WriteString(w, "Wrong!\n")
		}
	}
}
func moneyGoHomeHandle(input chan OPReq) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "POST":
			response_c := make(chan []byte, 2)
			input <- OPReq{
				op_code: op_all_money_go_home,
				Res_c:   response_c,
			}
			result := <-response_c
			w.Write(result)
		default:
			io.WriteString(w, "Wrong!\n")
		}
	}
}

func allsnapsHandle(input chan OPReq) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			response_c := make(chan []byte, 2)
			input <- OPReq{
				op_code: op_all_snaps,
				Res_c:   response_c,
			}
			result := <-response_c
			w.Write(result)
		default:
			io.WriteString(w, "Wrong!\n")
		}
	}
}
func assetspriceHandle(input chan OPReq) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			response_c := make(chan []byte, 2)
			input <- OPReq{
				op_code: assets_price,
				Res_c:   response_c,
			}
			result := <-response_c
			w.Write(result)
		default:
			io.WriteString(w, "Wrong!\n")
		}
	}
}
func paymentHandle(w http.ResponseWriter, req *http.Request) {

}

func user_interact(cmd_c chan PaymentReq, op_c chan OPReq, charge_c chan ChargeReq) {

	http.HandleFunc("/payment", makePaymentHandle(cmd_c))
	http.HandleFunc("/charges", makeChargeHandle(charge_c))
	http.HandleFunc("/moneygohome", moneyGoHomeHandle(op_c))
	http.HandleFunc("/snaps", allsnapsHandle(op_c))
	http.HandleFunc("/assetsprice", assetspriceHandle(op_c))
	log.Fatal(http.ListenAndServe(local_web_port, nil))
	log.Println("after web")
}

func fire_callback_url(v CallbackRespone) {
	jsonValue, jserr := json.Marshal(v)
	if jserr != nil {
		return
	}
	localURL := "http://127.0.0.1" + v.Callbackurl
	_, err := http.Post(localURL, "application/json", bytes.NewBuffer(jsonValue))
	if err != nil {
		log.Println(err)
	}
}
func all_money_gomyhome(userid string, sessionid string, privatekey string, pin string, pintoken string, admin_uuid string) {
	this_user := mixin.NewUser(userid, sessionid, privatekey, pin, pintoken)
	balance, err := this_user.ReadAssets()
	if err != nil {
		log.Println(err)
		return
	} else {
		var resp BalanceResponse
		err = json.Unmarshal(balance, &resp)
		if err != nil {
			log.Println(err)
			return
		}
		if resp.Error != "" {
			log.Println(resp.Error)
			return
		}
		for _, v := range resp.Data {
			if v.Balance == "0" {
				continue
			} else {
				trans_result, trans_err := this_user.Transfer(admin_uuid, v.Balance, v.Assetid, "allmoneygomyhome", uuid.Must(uuid.NewV4()).String())
				if trans_err != nil {
					log.Println(trans_err)
				} else {
					var resp TransferNetRespone
					err = json.Unmarshal(trans_result, &resp)

					if err != nil {
						log.Println(err)
					} else {
						if resp.TransferRes.Error != "" {
							log.Println(resp.TransferRes.Error)
						} else {
							log.Println(resp.TransferRes.Data.Snapshotid)
						}
					}

				}

			}
		}
	}

}
func create_mixin_account(account_name string, predefine_pin string, user_id string, session_id string, private_key string, result_chan chan MixinAccountindb) {
	user, err := mixin.CreateAppUser(account_name, predefine_pin, user_id, session_id, private_key)
	if err != nil {
		log.Println(err)
	} else {
		created_time, err := time.Parse(time.RFC3339Nano, user.CreatedAt)
		if err != nil {
			log.Println(err)
		} else {
			new_user := MixinAccountindb{
				Userid:        user.UserId,
				Sessionid:     user.SessionId,
				Pintoken:      user.PinToken,
				Privatekey:    user.PrivateKey,
				Pin:           predefine_pin,
				ClientReqid:   0,
				Utccreated_at: created_time,
			}
			result_chan <- new_user
		}

	}
}

func search_userincome(asset_id string, userid string, sessionid string, privatekey string, in_result_chan chan *Snapshot, in_progress_c chan Searchprogress, use_created_at time.Time, end_at time.Time, search_expired_after time.Time) {
	req_task := Searchtask{
		start_t:         end_at,
		end_t:           use_created_at,
		max_len:         500,
		yesterday2today: false,
		asset_id:        asset_id,
		userconfig: BotConfig{
			user_id:     userid,
			session_id:  sessionid,
			private_key: privatekey,
		},
		ongoing:            true,
		includesubaccount:  false,
		task_expired_after: search_expired_after,
	}
	for {
		if req_task.task_expired_after.IsZero() == false && time.Now().After(req_task.task_expired_after) {
			p := Searchprogress{
				search_task: req_task,
			}
			p.search_task.ongoing = false
			in_progress_c <- p
			log.Println("task is expired")
			return
		}
		var snaps []byte
		var err error
		snaps, err = mixin.MyNetworkSnapshots(req_task.asset_id, req_task.start_t, req_task.max_len, req_task.userconfig.user_id, req_task.userconfig.session_id, req_task.userconfig.private_key)

		if err != nil {
			in_progress_c <- Searchprogress{
				Error: err,
			}
			continue
		}
		var resp MixinResponse
		err = json.Unmarshal(snaps, &resp)

		if err != nil {
			in_progress_c <- Searchprogress{
				Error: err,
			}
			continue
		}
		if resp.Error != "" {
			log.Fatal("Server return error", resp.Error, " for req:", req_task.asset_id, " start ", req_task.start_t)
			return
		}
		len_of_snap := len(resp.Data)
		for _, v := range resp.Data {
			v.UserId = req_task.userconfig.user_id
			in_result_chan <- v
		}
		if len_of_snap == 0 {
			req_task.start_t = time.Now()
			p := Searchprogress{
				search_task: req_task,
			}
			in_progress_c <- p
			time.Sleep(30 * time.Second)
		} else {
			last_element := resp.Data[len(resp.Data)-1]
			req_task.start_t = last_element.CreatedAt
			p := Searchprogress{
				search_task: req_task,
			}
			p.search_task.start_t = last_element.CreatedAt

			if req_task.end_t.IsZero() == false && last_element.CreatedAt.Before(req_task.end_t) {
				p.search_task.ongoing = false
				in_progress_c <- p
				return
			}
			in_progress_c <- p
		}
	}
}
func restore_searchsnap(bot_config BotConfig, in_result_chan chan *Snapshot, in_progress_c chan Searchprogress, default_asset_id_group []string, searchtasks_array_inram []Searchtaskinram) {
	if len(searchtasks_array_inram) > 0 {
		for _, v := range searchtasks_array_inram {
			if v.Ongoing == true {
				log.Println(v.Ongoing, v.Starttime, v.Endtime, v.Userid, v.Assetid)
				unfinished_req_task := Searchtask{
					start_t:         v.Starttime,
					end_t:           v.Endtime,
					max_len:         500,
					yesterday2today: v.Yesterday2today,
					asset_id:        v.Assetid,
					ongoing:         v.Ongoing,
					userconfig: BotConfig{
						user_id:     v.Userid,
						session_id:  v.Sessionid,
						private_key: v.Privatekey,
					},
					includesubaccount:  v.Includesubaccount,
					task_expired_after: v.Taskexpired_at,
				}
				if v.Includesubaccount == false {
					go search_userincome(v.Assetid, v.Userid, v.Sessionid, v.Privatekey, in_result_chan, in_progress_c, v.Endtime, time.Now(), v.Taskexpired_at)
				} else {
					if v.Yesterday2today {
						go read_snap_to_future(unfinished_req_task, in_result_chan, in_progress_c)
					}
				}

			}
		}
	} else {
		botCreateAt := read_bot_created_time(bot_config.user_id, bot_config.session_id, bot_config.private_key)
		if botCreateAt.IsZero() {
			panic("Read bot profile failed")
		} else {
			log.Println("I am created at ", botCreateAt)
			for _, v := range default_asset_id_group {
				search_asset_task := Searchtask{
					start_t:           botCreateAt,
					max_len:           500,
					yesterday2today:   false,
					asset_id:          v,
					userconfig:        bot_config,
					includesubaccount: true,
				}
				go read_snap_to_future(search_asset_task, in_result_chan, in_progress_c)
			}
		}
	}
}

func main() {
	var my_snapshot_chan = make(chan *Snapshot, 1000)
	var asset_received_snap_chan = make(chan *Snapshot, 1000)
	var global_progress_c = make(chan Searchprogress, 1000)
	var quit_chan = make(chan int, 2)
	var req_cmd_chan = make(chan PaymentReq, 2)
	var charge_cmd_chan = make(chan ChargeReq, 2)
	var single_direction_op_cmd_chan = make(chan OPReq, 2)
	var new_account_received_chan = make(chan MixinAccountindb, 100)
	var payment_received_asset_chan = make(chan CallbackRespone, 100)
	var account_deposit_address_receive_chan = make(chan DepositNetResponse, 100)
	var should_create_more_account_c = make(chan uint, 10)
	var update_asset_price_c = make(chan uint, 10)
	// to support more asset, just add them in the following array
	default_asset_id_group := []string{XLM_ASSET_ID, EOS_ASSET_ID, ETH_ASSET_ID}
	timer1 := time.NewTimer(1 * time.Minute)

	db, err := gorm.Open("sqlite3", "test.db")
	if err != nil {
		panic("failed to connect database")
	}
	defer db.Close()

	db.AutoMigrate(&Snapshotindb{})
	db.AutoMigrate(&Searchtaskindb{})
	db.AutoMigrate(&MixinAccountindb{})
	db.AutoMigrate(&ClientReq{})
	db.AutoMigrate(&DepositAddressindb{})
	db.AutoMigrate(&AssetInformationindb{})
	db.AutoMigrate(&MessengerUserindb{})
	db.AutoMigrate(&Assetpriceindb{})
	db.AutoMigrate(&ChargeRecordindb{})
	var bot_config_instance = BotConfig{
		user_id:     userid,
		session_id:  sessionid,
		private_key: private_key,
	}

	//startup
	var admin_uuid_record MessengerUserindb
	db.Find(&MessengerUserindb{Messengerid: ADMIN_MessengerID}).First(&admin_uuid_record)
	if admin_uuid_record.ID == 0 {
		result := read_useruuid_from(bot_config_instance.user_id, bot_config_instance.session_id, bot_config_instance.private_key, ADMIN_MessengerID)
		if result != "" {
			log.Println(result)
			db.Create(&MessengerUserindb{Messengerid: ADMIN_MessengerID, Uuid: result})
		} else {
			log.Fatal("Failed to read admin uuid by it's messenger id")
		}
	}
	db.Find(&MessengerUserindb{Messengerid: ADMIN_MessengerID}).First(&admin_uuid_record)
	var ongoing_searchtasks_indb []Searchtaskindb
	var ongoing_searchtasks_inram []Searchtaskinram
	db.Find(&ongoing_searchtasks_indb)
	for _, v := range ongoing_searchtasks_indb {
		var this_user_record MixinAccountindb
		if db.Where(&MixinAccountindb{Userid: v.Userid}).First(&this_user_record).RecordNotFound() == false {
			var this_search_task_ram Searchtaskinram
			this_search_task_ram.Starttime = v.Starttime
			this_search_task_ram.Endtime = v.Endtime
			this_search_task_ram.Taskexpired_at = v.Taskexpired_at
			this_search_task_ram.Yesterday2today = v.Yesterday2today
			this_search_task_ram.Assetid = v.Assetid
			this_search_task_ram.Ongoing = v.Ongoing

			this_search_task_ram.Userid = v.Userid
			this_search_task_ram.Sessionid = this_user_record.Sessionid
			this_search_task_ram.Privatekey = this_user_record.Privatekey

			this_search_task_ram.Includesubaccount = v.Includesubaccount

			ongoing_searchtasks_inram = append(ongoing_searchtasks_inram, this_search_task_ram)
		} else {
			if v.Userid == bot_config_instance.user_id {
				var this_search_task_ram Searchtaskinram
				this_search_task_ram.Starttime = v.Starttime
				this_search_task_ram.Endtime = v.Endtime
				this_search_task_ram.Taskexpired_at = v.Taskexpired_at
				this_search_task_ram.Yesterday2today = v.Yesterday2today
				this_search_task_ram.Assetid = v.Assetid
				this_search_task_ram.Ongoing = v.Ongoing
				this_search_task_ram.Userid = v.Userid
				this_search_task_ram.Sessionid = bot_config_instance.session_id
				this_search_task_ram.Privatekey = bot_config_instance.private_key
				this_search_task_ram.Includesubaccount = v.Includesubaccount
				ongoing_searchtasks_inram = append(ongoing_searchtasks_inram, this_search_task_ram)
			}
		}
	}

	restore_searchsnap(bot_config_instance, my_snapshot_chan, global_progress_c, default_asset_id_group, ongoing_searchtasks_inram)
	go user_interact(req_cmd_chan, single_direction_op_cmd_chan, charge_cmd_chan)

	should_create_more_account_c <- 1
	update_asset_price_c <- 1
	for {
		select {
		case pv := <-global_progress_c:
			if pv.Error != nil {
				log.Println(pv.Error)
				continue
			}
			searchtaskindb := Searchtaskindb{}
			query_task := Searchtaskindb{
				Endtime:           pv.search_task.end_t,
				Assetid:           pv.search_task.asset_id,
				Userid:            pv.search_task.userconfig.user_id,
				Includesubaccount: pv.search_task.includesubaccount,
			}
			if db.Where(&query_task).First(&searchtaskindb).RecordNotFound() {
				var this_record = Searchtaskindb{
					Starttime:         pv.search_task.start_t,
					Endtime:           pv.search_task.end_t,
					Yesterday2today:   pv.search_task.yesterday2today,
					Assetid:           pv.search_task.asset_id,
					Ongoing:           pv.search_task.ongoing,
					Userid:            pv.search_task.userconfig.user_id,
					Includesubaccount: pv.search_task.includesubaccount,
					Taskexpired_at:    pv.search_task.task_expired_after,
				}
				db.Create(&this_record)
			} else {
				searchtaskindb.Starttime = pv.search_task.start_t
				searchtaskindb.Ongoing = pv.search_task.ongoing
				db.Save(&searchtaskindb)
			}
		case v := <-my_snapshot_chan:
			snapInDb := Snapshotindb{
				SnapshotId: v.SnapshotId,
			}
			db.First(&snapInDb, "snapshot_id = ?", v.SnapshotId)
			if snapInDb.CreatedAt.IsZero() {
				var thisrecord = Snapshotindb{
					SnapshotId:    v.SnapshotId,
					Amount:        v.Amount,
					AssetId:       v.Asset.AssetId,
					Source:        v.Source,
					SnapCreatedAt: v.CreatedAt,
					UserId:        v.UserId,
					TraceId:       v.TraceId,
					OpponentId:    v.OpponentId,
					Data:          v.Data,
					OPType:        v.Type,
				}
				if v.AssetId != "" {
					thisrecord.AssetId = v.AssetId
				}
				db.Create(&thisrecord)
				f, err := strconv.ParseFloat(v.Amount, 64)
				if err != nil {
					log.Println(err)
				} else {
					if f > 0 {
						asset_received_snap_chan <- v
					}
				}
			}
		case v := <-asset_received_snap_chan:
			var matched_account MixinAccountindb
			if db.Where(&MixinAccountindb{Userid: v.UserId}).First(&matched_account).RecordNotFound() == false {
				go all_money_gomyhome(matched_account.Userid, matched_account.Sessionid, matched_account.Privatekey, matched_account.Pin, matched_account.Pintoken, admin_uuid_record.Uuid)
				if matched_account.ClientReqid != 0 {
					var matched_req ClientReq
					db.First(&matched_req, matched_account.ClientReqid)
					if matched_req.ID != 0 {
						var callback_response CallbackRespone
						callback_response.Reqid = matched_req.Reqid
						callback_response.Callbackurl = matched_req.Callbackurl
						callback_response.Paymentrecord = Payment_Record{
							CreatedAt:  v.CreatedAt,
							Amount:     v.Amount,
							AssetId:    v.AssetId,
							SnapshotId: v.SnapshotId,
						}
						var asset_price Assetpriceindb
						if db.Where(&Assetpriceindb{Assetid: v.AssetId}).First(&asset_price).RecordNotFound() == false {
							amount_float, _ := strconv.ParseFloat(v.Amount, 64)

							float_price_usd, _ := strconv.ParseFloat(asset_price.Priceinusd, 64)
							callback_response.Valueinusd = amount_float * float_price_usd

							float_price_btc, _ := strconv.ParseFloat(asset_price.Priceinbtc, 64)
							callback_response.Valueinbtc += amount_float * float_price_btc
						}
						payment_received_asset_chan <- callback_response
					}
				}
			}

		case v := <-payment_received_asset_chan:
			go fire_callback_url(v)
		case <-quit_chan:
			log.Println("finished")
			return

		case new_user := <-new_account_received_chan:
			db.Create(&new_user)
			for _, v := range default_asset_id_group {
				depositRecord := DepositAddressindb{
					Accountrecord_id: new_user.ID,
					Assetid:          v,
				}
				db.Create(&depositRecord)
				go read_asset_deposit_address(v, new_user.Userid, new_user.Sessionid, new_user.Privatekey, account_deposit_address_receive_chan)
			}
		case asset_deposit_address_result := <-account_deposit_address_receive_chan:
			if asset_deposit_address_result.Error == nil {
				var matched_user MixinAccountindb
				db.Where(&MixinAccountindb{Userid: asset_deposit_address_result.Accountid}).First(&matched_user)
				var depositRecord DepositAddressindb
				db.Where("accountrecord_id = ?", matched_user.ID).Where("assetid = ?", asset_deposit_address_result.Assetid).First(&depositRecord)
				if depositRecord.CreatedAt.IsZero() {
					panic("The record should has been created when the user is created")
				} else {
					depositRecord.Publicaddress = asset_deposit_address_result.MixinResponse.Data.PublicKey
					depositRecord.Accountname = asset_deposit_address_result.MixinResponse.Data.AccountName
					depositRecord.Accounttag = asset_deposit_address_result.MixinResponse.Data.AccountTag
					depositRecord.Confirmblock = asset_deposit_address_result.MixinResponse.Data.Confirmblock
					depositRecord.Iconurl = asset_deposit_address_result.MixinResponse.Data.IconURL
					db.Save(&depositRecord)

					var asset_record AssetInformationindb
					if db.Where(&AssetInformationindb{Assetid: depositRecord.Assetid}).First(&asset_record).RecordNotFound() {
						//first found asset
						asset_record.Symbol = asset_deposit_address_result.MixinResponse.Data.Symbol
						asset_record.Name = asset_deposit_address_result.MixinResponse.Data.Name
						asset_record.Assetkey = asset_deposit_address_result.MixinResponse.Data.Assetkey
						asset_record.Chainid = asset_deposit_address_result.MixinResponse.Data.Chainid
						asset_record.Assetid = asset_deposit_address_result.Assetid
						db.Create(&asset_record)
					}
				}
			}

		case <-timer1.C:
			should_create_more_account_c <- 1
			update_asset_price_c <- 1

		case <-update_asset_price_c:
			for _, asset_id := range default_asset_id_group {
				result, err := mixin.Deposit(asset_id, bot_config_instance.user_id, bot_config_instance.session_id, bot_config_instance.private_key)

				if err != nil {
					log.Fatal(err)
					continue
				}

				var resp MixinDepositResponse
				err = json.Unmarshal(result, &resp)

				if err != nil {
					log.Fatal(err)
					continue
				}
				if resp.Error != "" {
					log.Fatal(resp.Error)
					continue
				}
				var this_asset_price Assetpriceindb
				this_asset_price.Assetid = asset_id

				if db.Where(&Assetpriceindb{Assetid: asset_id}).First(&this_asset_price).RecordNotFound() == true {
					var new_asset_price Assetpriceindb
					new_asset_price.Assetid = asset_id
					new_asset_price.Priceinusd = resp.Data.Priceusd
					new_asset_price.Priceinbtc = resp.Data.Pricebtc
					db.Create(&new_asset_price)
				} else {
					this_asset_price.Priceinusd = resp.Data.Priceusd
					this_asset_price.Priceinbtc = resp.Data.Pricebtc
					db.Save(&this_asset_price)
				}
			}
		case <-should_create_more_account_c:
			var free_mixinaccounts []MixinAccountindb
			db.Model(&MixinAccountindb{}).Where("client_reqid = ?", "0").Find(&free_mixinaccounts)
			available_mixin_account := len(free_mixinaccounts)
			if available_mixin_account < 10 {
				for i := 20; i > available_mixin_account; i-- {
					go create_mixin_account(PREDEFINE_NAME, PREDEFINE_PIN, bot_config_instance.user_id, bot_config_instance.session_id, bot_config_instance.private_key, new_account_received_chan)
				}
			}

			//read all free account, and check all deposit address is ready
			for _, account := range free_mixinaccounts {
				var payment_addresses []DepositAddressindb
				db.Where(&DepositAddressindb{Accountrecord_id: account.ID}).Find(&payment_addresses)
				for _, payment_address := range payment_addresses {
					if payment_address.Publicaddress == "" && payment_address.Accountname == "" && payment_address.Accounttag == "" {
						log.Println("some account deposit address is still missing")
						go read_asset_deposit_address(payment_address.Assetid, account.Userid, account.Sessionid, account.Privatekey, account_deposit_address_receive_chan)
					}
				}
			}
		case v := <-single_direction_op_cmd_chan:
			switch v.op_code {
			case op_all_money_go_home:
				var allaccount []MixinAccountindb
				db.Find(&allaccount)
				for _, v := range allaccount {
					go all_money_gomyhome(v.Userid, v.Sessionid, v.Privatekey, v.Pin, v.Pintoken, admin_uuid_record.Uuid)
				}
				type Message struct {
					Action string
					Result string
				}
				m := Message{"Allmoneygohome", fmt.Sprintf("total %d account will send all balance to admin", len(allaccount))}
				b, _ := json.Marshal(m)
				v.Res_c <- b
			case op_all_snaps:
				var allbot_snaps []Snapshotindb
				db.Find(&allbot_snaps)
				type Snap struct {
					Snapshotid string
					Createdat  time.Time
					Amount     string
					Assetid    string
				}
				type Message struct {
					Action string
					Result []Snap
				}
				var result []Snap
				for _, v := range allbot_snaps {
					this := Snap{v.SnapshotId, v.SnapCreatedAt, v.Amount, v.AssetId}
					result = append(result, this)
				}

				m := Message{"Allsnap", result}
				b, _ := json.Marshal(m)
				v.Res_c <- b
			case assets_price:
				type Assetprice2client struct {
					Fullname string
					Symbol   string
					USDPrice float64
					BTCPrice float64
					Assetid  string
				}

				var asset_prices []Assetprice2client
				for _, v := range default_asset_id_group {
					var asset_info AssetInformationindb
					if db.Where(&AssetInformationindb{Assetid: v}).First(&asset_info).RecordNotFound() == false {
						var asset_price Assetpriceindb
						if db.Where(&Assetpriceindb{Assetid: v}).First(&asset_price).RecordNotFound() == false {
							usdprice, _ := strconv.ParseFloat(asset_price.Priceinusd, 64)
							btcprice, _ := strconv.ParseFloat(asset_price.Priceinbtc, 64)
							this := Assetprice2client{
								Fullname: asset_info.Name,
								Symbol:   asset_info.Symbol,
								Assetid:  asset_info.Assetid,
								USDPrice: usdprice,
								BTCPrice: btcprice,
							}
							asset_prices = append(asset_prices, this)
						}
					}
				}
				b, _ := json.Marshal(asset_prices)
				v.Res_c <- b
			}

		case v := <-charge_cmd_chan:
			if v.Method == "GET" {
				var charge_record ChargeRecordindb
				if db.Where(&ChargeRecordindb{Customerid: v.Customerid}).First(&charge_record).RecordNotFound() == false {
					var resp ChargeResponse
					resp.Currency = charge_record.Currency
					resp.Amount = charge_record.Amount
					resp.Customerid = charge_record.Customerid
					resp.Expired_after = charge_record.Expiredafter

					//looking for currency asset id from SYMBOL
					var currentcy_asset_info AssetInformationindb
					if db.Where(&AssetInformationindb{Symbol: charge_record.Currency}).First(&currentcy_asset_info).RecordNotFound() == false {
						var charge_mixinaccount_relation ChargeRelationWithMixinAccountindb
						if db.Where(&ChargeRelationWithMixinAccountindb{Chargerecordid: charge_record.ID}).First(&charge_mixinaccount_relation).RecordNotFound() == false {
							//find relation record with mixin account
							var mixin_account MixinAccountindb

							db.Find(&mixin_account, charge_mixinaccount_relation.Mixinaccountid)
							if mixin_account.ID != 0 {
								var payment_address DepositAddressindb
								//find the currency payment address, this is a deposit address for crypto asset, different with lightning charge invoice
								if db.Where(&DepositAddressindb{Accountrecord_id: mixin_account.ID, Assetid: currentcy_asset_info.Assetid}).First(&payment_address).RecordNotFound() == false {
									resp.Paymentmethod = PaymentMethod{
										Name:           currentcy_asset_info.Symbol,
										PaymentAddress: payment_address.Publicaddress,
										PaymentAccount: payment_address.Accountname,
										PaymentMemo:    payment_address.Accounttag,
									}
								}

								//fill market price
								var asset_price_record Assetpriceindb
								if db.Where(&Assetpriceindb{Assetid: currentcy_asset_info.Assetid}).First(&asset_price_record).RecordNotFound() == false {
									resp.Paymentmethod.Priceinusd = asset_price_record.Priceinusd
									resp.Paymentmethod.Priceinbtc = asset_price_record.Priceinbtc
								}

								resp.Paidstatus = 0
								resp.Receivedamount = 0
								var all_assets_payment_snapshots_indb []Snapshotindb
								db.Where(&Snapshotindb{UserId: mixin_account.Userid, AssetId: currentcy_asset_info.Assetid}).Find(&all_assets_payment_snapshots_indb)
								for _, v := range all_assets_payment_snapshots_indb {
									f, err := strconv.ParseFloat(v.Amount, 64)
									if err != nil {
										log.Fatal(err)
										continue
									} else {
										if f > 0 {
											resp.Receivedamount += f
										}
									}
								}

								if resp.Receivedamount >= resp.Amount {
									resp.Paidstatus = 2
								} else if resp.Receivedamount > 0 {
									resp.Paidstatus = 1
								}
							}
						}
					}
					resp.Webhookurl = v.Webhookurl
					v.Res_c <- resp
				}
			} else {
				var resp ChargeResponse
				resp.Currency = v.Currency
				resp.Amount = v.Amount
				resp.Customerid = v.Customerid
				resp.Expired_after = v.Expired_after
				resp.Paidstatus = 0
				resp.Receivedamount = 0
				resp.Webhookurl = v.Webhookurl
				v.Res_c <- resp
			}
		case v := <-req_cmd_chan:
			if v.Method == "GET" {
				payment_id := v.Reqid
				var req ClientReq
				var res PaymentRes
				response_c := v.Res_c
				db.Where(&ClientReq{Reqid: payment_id}).Find(&req)
				if req.ID != 0 {
					res.Reqid = v.Reqid
					var mixin_account MixinAccountindb
					db.Find(&mixin_account, req.MixinAccountid)
					if mixin_account.ID != 0 {
						var payment_addresses []DepositAddressindb
						db.Where(&DepositAddressindb{Accountrecord_id: mixin_account.ID}).Find(&payment_addresses)
						var all_method []PaymentMethod
						for _, v := range payment_addresses {
							var pv PaymentMethod
							var asset_info AssetInformationindb
							db.Where(&AssetInformationindb{Assetid: v.Assetid}).First(&asset_info)
							log.Println(asset_info)
							pv.Name = asset_info.Symbol
							pv.PaymentAddress = v.Publicaddress
							pv.PaymentAccount = v.Accountname
							pv.PaymentMemo = v.Accounttag

							var asset_price_record Assetpriceindb
							if db.Where(&Assetpriceindb{Assetid: v.Assetid}).First(&asset_price_record).RecordNotFound() == false {
								pv.Priceinusd = asset_price_record.Priceinusd
								pv.Priceinbtc = asset_price_record.Priceinbtc
							}

							all_method = append(all_method, pv)
						}

						res.Payment_methods = all_method

						var all_payment_snapshots_indb []Snapshotindb
						var all_payment_snapshots []Payment_Record
						var received_in_usd float64
						var received_in_btc float64
						db.Where(&Snapshotindb{UserId: mixin_account.Userid}).Find(&all_payment_snapshots_indb)
						for _, v := range all_payment_snapshots_indb {
							f, err := strconv.ParseFloat(v.Amount, 64)
							if err != nil {
								log.Println(err)
								continue
							} else {
								if f > 0 {
									this_snap := Payment_Record{
										Amount:     v.Amount,
										AssetId:    v.AssetId,
										CreatedAt:  v.SnapCreatedAt,
										SnapshotId: v.SnapshotId,
									}
									all_payment_snapshots = append(all_payment_snapshots, this_snap)
									var asset_price Assetpriceindb
									if db.Where(&Assetpriceindb{Assetid: v.AssetId}).First(&asset_price).RecordNotFound() == false {
										float_price_usd, _ := strconv.ParseFloat(asset_price.Priceinusd, 64)
										received_in_usd += f * float_price_usd
										float_price_btc, _ := strconv.ParseFloat(asset_price.Priceinbtc, 64)
										received_in_btc += f * float_price_btc
									}
								}
							}
						}
						res.Payment_records = all_payment_snapshots
						res.ReceivedinUSD = received_in_usd
						res.ReceivedinBTC = received_in_btc
						response_c <- res
					} else {
						response_c <- res
					}
				} else {
					response_c <- res
				}
			} else {
				unique_id := v.Reqid
				response_c := v.Res_c
				var res PaymentRes
				var free_mixinaccount MixinAccountindb
				if db.Where("client_reqid = ?", "0").First(&free_mixinaccount).RecordNotFound() == false {
					res.Reqid = v.Reqid
					new_req := ClientReq{
						Reqid:          unique_id,
						Callbackurl:    v.Callback,
						MixinAccountid: free_mixinaccount.ID,
					}
					db.Create(&new_req)
					free_mixinaccount.ClientReqid = new_req.ID
					db.Save(&free_mixinaccount)
					go search_userincome("", free_mixinaccount.Userid, free_mixinaccount.Sessionid, free_mixinaccount.Privatekey, my_snapshot_chan, global_progress_c, free_mixinaccount.Utccreated_at, time.Now(), time.Now().Add(time.Duration(v.Expired_after)*time.Minute))
					var payment_addresses []DepositAddressindb
					db.Where(&DepositAddressindb{Accountrecord_id: free_mixinaccount.ID}).Find(&payment_addresses)
					var all_method []PaymentMethod
					for _, v := range payment_addresses {
						var pv PaymentMethod
						var asset_info AssetInformationindb
						db.Where(&AssetInformationindb{Assetid: v.Assetid}).First(&asset_info)
						log.Println(asset_info)
						pv.Name = asset_info.Symbol
						pv.PaymentAddress = v.Publicaddress
						pv.PaymentAccount = v.Accountname
						pv.PaymentMemo = v.Accounttag

						var asset_price Assetpriceindb
						if db.Where(&Assetpriceindb{Assetid: v.Assetid}).First(&asset_price).RecordNotFound() == false {
							pv.Priceinusd = asset_price.Priceinusd
							pv.Priceinbtc = asset_price.Priceinbtc
						}
						all_method = append(all_method, pv)

					}
					res.Payment_methods = all_method
				} else {
					log.Println("no new user account")
				}
				response_c <- res
			}
		}
	}
}
