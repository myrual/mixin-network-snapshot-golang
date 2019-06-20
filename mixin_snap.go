package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strconv"
	"time"

	mixin "github.com/MooooonStar/mixin-sdk-go/network"
	"github.com/gofrs/uuid"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

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

	scan_interval_in_seconds = 10
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
type payment_records struct {
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

type ClientReq struct {
	gorm.Model
	Reqid          string
	Callbackurl    string
	MixinAccountid uint
	Callbackfired  bool
}

type Searchtaskindb struct {
	gorm.Model
	Starttime         time.Time
	Endtime           time.Time
	Yesterday2today   bool
	Assetid           string
	Ongoing           bool
	Userid            string
	Sessionid         string
	Privatekey        string
	Includesubaccount bool
	Taskexpired_at    time.Time
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
	yesterday2today    bool
	max_len            int
	asset_id           string
	ongoing            bool
	userid             string
	sessionid          string
	privatekey         string
	includesubaccount  bool
	task_expired_after time.Time
}

type Searchprogress struct {
	search_task Searchtask
	Error       error
}

type PaymentReqhttp struct {
	Reqid    string `json:"reqid"`
	Callback string `json:"callback"`
}

type PaymentReq struct {
	Method   string
	Reqid    string
	Callback string
	Res_c    chan PaymentRes
}
type PaymentMethod struct {
	Name        string
	PublicKey   string
	AccountName string
	AccountTag  string
}
type PaymentRes struct {
	Reqid           string
	Payment_methods []PaymentMethod
	Payment_records []payment_records
	Balance         []Asset
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

	ADMIN_UUID     = "28ee416a-0eaa-4133-bc79-9676909b7b4e"
	PREDEFINE_PIN  = "198435"
	PREDEFINE_NAME = "tom"
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

		snaps, err = mixin.NetworkSnapshots(req_task.asset_id, req_task.start_t, true, req_task.max_len, req_task.userid, req_task.sessionid, req_task.privatekey)

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
			log.Println(r.URL.Query())
			keys, ok := r.URL.Query()["reqid"]
			log.Println(keys)
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
				Reqid:    p.Reqid,
				Callback: p.Callback,
				Res_c:    payment_res_c,
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
func paymentHandle(w http.ResponseWriter, req *http.Request) {

}

func user_interact(cmd_c chan PaymentReq, output_c chan string) {

	http.HandleFunc("/payment", makePaymentHandle(cmd_c))
	log.Fatal(http.ListenAndServe(":8080", nil))
	log.Println("after web")
}

func all_money_gomyhome(userid string, sessionid string, privatekey string, pin string, pintoken string) {
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
			log.Println(this_user.UserId, v.Assetid, v.Balance)
			if v.Balance == "0" {
				continue
			} else {
				trans_result, trans_err := this_user.Transfer(ADMIN_UUID, v.Balance, v.Assetid, "allmoneygomyhome", uuid.Must(uuid.NewV4()).String())
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

func search_userincome(asset_id string, userid string, sessionid string, privatekey string, in_result_chan chan *Snapshot, in_progress_c chan Searchprogress, created_at time.Time, end_at time.Time, search_expired_after time.Time) {
	req_task := Searchtask{
		start_t:            end_at,
		end_t:              created_at,
		max_len:            500,
		yesterday2today:    false,
		asset_id:           asset_id,
		userid:             userid,
		sessionid:          sessionid,
		privatekey:         privatekey,
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
		snaps, err = mixin.MyNetworkSnapshots(req_task.asset_id, req_task.start_t, req_task.max_len, req_task.userid, req_task.sessionid, req_task.privatekey)

		if err != nil {
			in_progress_c <- Searchprogress{
				Error: err,
			}
			continue
		}
		var f interface{}
		errjs := json.Unmarshal(snaps, &f)
		if errjs != nil {

		} else {
			log.Println(f)
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
			v.UserId = req_task.userid
			in_result_chan <- v
			log.Println(v)
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
func restore_searchsnap(user_config BotConfig, in_result_chan chan *Snapshot, in_progress_c chan Searchprogress, default_asset_id_group []string, searchtasks_array_indb []Searchtaskindb) {
	if len(searchtasks_array_indb) > 0 {
		for _, v := range searchtasks_array_indb {
			if v.Ongoing == true {
				log.Println(v.Ongoing, v.Starttime, v.Endtime, v.Userid, v.Assetid)
				unfinished_req_task := Searchtask{
					start_t:            v.Starttime,
					end_t:              v.Endtime,
					max_len:            500,
					yesterday2today:    v.Yesterday2today,
					asset_id:           v.Assetid,
					ongoing:            v.Ongoing,
					userid:             v.Userid,
					sessionid:          v.Sessionid,
					privatekey:         v.Privatekey,
					includesubaccount:  v.Includesubaccount,
					task_expired_after: v.Taskexpired_at,
				}
				if v.Includesubaccount == false {
					go search_userincome(v.Assetid, v.Userid, v.Sessionid, v.Privatekey, in_result_chan, in_progress_c, v.Endtime, time.Now(), time.Now().Add(time.Hour*4))
				} else {
					if v.Yesterday2today {
						go read_snap_to_future(unfinished_req_task, in_result_chan, in_progress_c)
					}
				}

			}
		}
	} else {
		botCreateAt := read_bot_created_time(user_config.user_id, user_config.session_id, user_config.private_key)
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
					userid:            user_config.user_id,
					sessionid:         user_config.session_id,
					privatekey:        user_config.private_key,
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
	var user_output_chan = make(chan string, 100)
	var new_account_received_chan = make(chan MixinAccountindb, 100)
	var payment_received_asset_chan = make(chan ClientReq, 100)
	var account_deposit_address_receive_chan = make(chan DepositNetResponse, 100)
	var should_create_more_account_c = make(chan uint, 10)

	timer1 := time.NewTimer(1 * time.Minute)

	default_asset_id_group := []string{XLM_ASSET_ID, EOS_ASSET_ID}
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

	var user_config = BotConfig{
		user_id:     userid,
		session_id:  sessionid,
		private_key: private_key,
	}
	var ongoing_searchtasks_indb []Searchtaskindb
	db.Find(&ongoing_searchtasks_indb)
	restore_searchsnap(user_config, my_snapshot_chan, global_progress_c, default_asset_id_group, ongoing_searchtasks_indb)
	go user_interact(req_cmd_chan, user_output_chan)

	should_create_more_account_c <- 1
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
				Userid:            pv.search_task.userid,
				Includesubaccount: pv.search_task.includesubaccount,
			}
			db.Where(&query_task).First(&searchtaskindb)
			if searchtaskindb.CreatedAt.IsZero() {
				var this_record = Searchtaskindb{
					Starttime:         pv.search_task.start_t,
					Endtime:           pv.search_task.end_t,
					Yesterday2today:   pv.search_task.yesterday2today,
					Assetid:           pv.search_task.asset_id,
					Ongoing:           pv.search_task.ongoing,
					Userid:            pv.search_task.userid,
					Sessionid:         pv.search_task.sessionid,
					Privatekey:        pv.search_task.privatekey,
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
			db.Where(&MixinAccountindb{Userid: v.UserId}).First(&matched_account)
			if matched_account.ID != 0 {
				if matched_account.ClientReqid != 0 {
					var matched_req ClientReq
					db.First(&matched_req, matched_account.ClientReqid)
					payment_received_asset_chan <- matched_req
				}
			}

		case v := <-payment_received_asset_chan:
			var account MixinAccountindb
			db.Find(&account, v.MixinAccountid)
			log.Println("The request id ", v.ID, " receive payment ", v.Callbackurl, " in mixin account record id", v.MixinAccountid, " uuid ", account.Userid)
			go all_money_gomyhome(account.Userid, account.Sessionid, account.Privatekey, account.Pin, account.Pintoken)
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
				}
			}

		case <-timer1.C:
			should_create_more_account_c <- 1

		case <-should_create_more_account_c:
			var free_mixinaccounts []MixinAccountindb
			db.Model(&MixinAccountindb{}).Where("client_reqid = ?", "0").Find(&free_mixinaccounts)
			available_mixin_account := len(free_mixinaccounts)
			if available_mixin_account < 10 {
				for i := 20; i > available_mixin_account; i-- {
					go create_mixin_account(PREDEFINE_NAME, PREDEFINE_PIN, user_config.user_id, user_config.session_id, user_config.private_key, new_account_received_chan)
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

		case v := <-req_cmd_chan:
			if v.Method == "GET" {
				log.Println("GET", v.Reqid)
				payment_id := v.Reqid
				var req ClientReq
				var res PaymentRes
				response_c := v.Res_c
				db.Where(&ClientReq{Reqid: payment_id}).Find(&req)
				if req.ID != 0 {
					log.Println("GET req record with ", v.Reqid)
					res.Reqid = v.Reqid
					var mixin_account MixinAccountindb
					db.Find(&mixin_account, req.MixinAccountid)
					if mixin_account.ID != 0 {
						var payment_addresses []DepositAddressindb
						db.Where(&DepositAddressindb{Accountrecord_id: mixin_account.ID}).Find(&payment_addresses)
						var all_method []PaymentMethod
						for _, v := range payment_addresses {
							var pv PaymentMethod
							switch v.Assetid {
							case EOS_ASSET_ID:
								pv.Name = "EOS"
							case XLM_ASSET_ID:
								pv.Name = "XLM"
							}
							pv.PublicKey = v.Publicaddress
							pv.AccountName = v.Accountname
							pv.AccountTag = v.Accounttag

							all_method = append(all_method, pv)
						}

						res.Payment_methods = all_method

						var all_payment_snapshots_indb []Snapshotindb
						var all_payment_snapshots []payment_records
						db.Where(&Snapshotindb{UserId: mixin_account.Userid}).Find(&all_payment_snapshots_indb)
						for _, v := range all_payment_snapshots_indb {
							this_snap := payment_records{
								Amount:     v.Amount,
								AssetId:    v.AssetId,
								CreatedAt:  v.SnapCreatedAt,
								SnapshotId: v.SnapshotId,
							}
							all_payment_snapshots = append(all_payment_snapshots, this_snap)
						}
						res.Payment_records = all_payment_snapshots
						response_c <- res
					} else {
						response_c <- res
					}
				} else {
					response_c <- res
				}
			} else {
				log.Println("POST", v.Reqid, v.Callback)
				unique_id := v.Reqid
				response_c := v.Res_c
				var res PaymentRes
				var free_mixinaccount MixinAccountindb
				db.Where("client_reqid = ?", "0").First(&free_mixinaccount)
				if free_mixinaccount.ID != 0 {
					res.Reqid = v.Reqid
					new_req := ClientReq{
						Reqid:          unique_id,
						Callbackurl:    v.Callback,
						MixinAccountid: free_mixinaccount.ID,
					}
					db.Create(&new_req)
					free_mixinaccount.ClientReqid = new_req.ID
					db.Save(&free_mixinaccount)
					log.Println("POST", new_req.Reqid, new_req.Callbackurl, free_mixinaccount.Userid)
					go search_userincome("", free_mixinaccount.Userid, free_mixinaccount.Sessionid, free_mixinaccount.Privatekey, my_snapshot_chan, global_progress_c, free_mixinaccount.Utccreated_at, time.Now(), time.Now().Add(time.Hour*4))
					var payment_addresses []DepositAddressindb
					db.Where(&DepositAddressindb{Accountrecord_id: free_mixinaccount.ID}).Find(&payment_addresses)
					var all_method []PaymentMethod
					for _, v := range payment_addresses {
						var pv PaymentMethod
						switch v.Assetid {
						case EOS_ASSET_ID:
							pv.Name = "EOS"
						case XLM_ASSET_ID:
							pv.Name = "XLM"
						}
						pv.PublicKey = v.Publicaddress
						pv.AccountName = v.Accountname
						pv.AccountTag = v.Accounttag

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
