package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	mixin "github.com/MooooonStar/mixin-sdk-go/network"
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
)

type Snapshot struct {
	SnapshotId string `json:"snapshot_id"`
	Amount     string `json:"amount"`
	Asset      struct {
		AssetId string `json:"asset_id"`
	} `json:"asset"`
	Source     string    `json:"source"`
	CreatedAt  time.Time `json:"created_at"`
	TraceId    string    `json:"trace_id"`
	UserId     string    `json:"user_id"`
	OpponentId string    `json:"opponent_id"`
	Data       string    `json:"data"`
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
}

type MixinAccountindb struct {
	gorm.Model
	Userid      string `gorm:"primary_key"`
	Sessionid   string
	Pintoken    string
	Privatekey  string
	Pin         string
	ClientReqid uint
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
	Callbackurl    string
	MixinAccountid uint
}

type Searchtaskindb struct {
	gorm.Model
	Starttime       time.Time
	Endtime         time.Time
	Lasttime        time.Time
	Yesterday2today bool
	Assetid         string
	Ongoing         bool
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

type Searchtask struct {
	start_t         time.Time
	end_t           time.Time
	last_t          time.Time
	yesterday2today bool
	max_len         int
	asset_id        string
	ongoing         bool
}

type Searchprogress struct {
	search_task Searchtask
	Error       error
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

func read_my_snap(req_task Searchtask, user_config BotConfig, result_chan chan *Snapshot, progress_chan chan Searchprogress, quit_c chan int) {
	req_task.last_t = req_task.start_t
	for {
		snaps, err := mixin.NetworkSnapshots(req_task.asset_id, req_task.last_t, req_task.yesterday2today, req_task.max_len, user_config.user_id, user_config.session_id, user_config.private_key)

		if err != nil {
			progress_chan <- Searchprogress{
				Error: err,
			}
			continue
		}

		var resp MixinResponse
		err = json.Unmarshal(snaps, &resp)

		if err != nil {
			progress_chan <- Searchprogress{
				Error: err,
			}
			continue
		}
		if resp.Error != "" {
			log.Println("Server return error", resp.Error, " for req:", req_task.asset_id, " start ", req_task.start_t)
			return
		}
		for _, v := range resp.Data {
			if v.UserId != "" {
				result_chan <- v
			}
		}
		len_of_snap := len(resp.Data)
		if len_of_snap == 0 {
			time.Sleep(60 * time.Second)
			continue
		}
		last_element := resp.Data[len(resp.Data)-1]
		req_task.last_t = last_element.CreatedAt
		p := Searchprogress{
			search_task: req_task,
		}
		if last_element.CreatedAt.After(req_task.end_t) && req_task.end_t.IsZero() == false {
			p.search_task.ongoing = false
			progress_chan <- p
			return
		}
		p.search_task.ongoing = true
		progress_chan <- p
		if len_of_snap < req_task.max_len {
			time.Sleep(60 * time.Second)
		}
	}
}

func user_interact(cmd_c chan string, output_c chan string) {
	scanner := bufio.NewScanner(os.Stdin)
	var cmd string

	for {
		select {
		case v := <-output_c:
			log.Println(v)
		}
		scanner.Scan()
		cmd = scanner.Text()
		cmd_c <- cmd
	}
}

func create_mixin_account(account_name string, predefine_pin string, user_id string, session_id string, private_key string, result_chan chan MixinAccountindb) {
	user, err := mixin.CreateAppUser(account_name, predefine_pin, user_id, session_id, private_key)
	if err != nil {
		log.Println(err)
	} else {
		new_user := MixinAccountindb{
			Userid:      user.UserId,
			Sessionid:   user.SessionId,
			Pintoken:    user.PinToken,
			Privatekey:  user.PrivateKey,
			Pin:         predefine_pin,
			ClientReqid: 0,
		}
		result_chan <- new_user
	}
}

func main() {
	var start_time2 = time.Date(2018, 4, 25, 0, 0, 0, 0, time.UTC)
	var my_snapshot_chan = make(chan *Snapshot, 1000)
	var progress_chan = make(chan Searchprogress, 1000)
	var quit_chan = make(chan int, 2)
	var user_cmd_chan = make(chan string, 10)
	var user_output_chan = make(chan string, 100)
	var mixin_account_chan = make(chan MixinAccountindb, 100)
	var mixin_deposit_chan = make(chan DepositNetResponse, 100)
	var checkremain_account_c = make(chan uint, 10)
	var checkaccount_deposit_c = make(chan MixinAccountindb, 10)
	var req_create_payment_chan = make(chan string, 10)
	var req_read_deposit_chan = make(chan MixinAccountindb, 10)
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
	req_task := Searchtask{
		start_t:         start_time2,
		max_len:         500,
		yesterday2today: true,
		asset_id:        CNB_ASSET_ID,
	}
	var searchtasks_array_indb []Searchtaskindb
	db.Find(&searchtasks_array_indb)
	log.Println("Total ", len(searchtasks_array_indb), " search task")
	snap_cnb_quit_c := make(chan int, 1)

	if len(searchtasks_array_indb) > 0 {
		for _, v := range searchtasks_array_indb {

			if v.Ongoing == true {
				log.Println(v.Ongoing, v.Starttime, v.Endtime, v.Lasttime)
				unfinished_req_task := Searchtask{
					start_t:         v.Starttime,
					end_t:           v.Endtime,
					last_t:          v.Lasttime,
					yesterday2today: v.Yesterday2today,
					asset_id:        v.Assetid,
					ongoing:         v.Ongoing,
				}
				go read_my_snap(unfinished_req_task, user_config, my_snapshot_chan, progress_chan, snap_cnb_quit_c)
			}
		}
	} else {
		go read_my_snap(req_task, user_config, my_snapshot_chan, progress_chan, snap_cnb_quit_c)
	}
	promot := "allsnap: read all snap\n"
	promot += "status: ongoing search task\n"
	promot += "your selection:"
	user_output_chan <- promot
	go user_interact(user_cmd_chan, user_output_chan)

	checkremain_account_c <- 1
	for {
		select {
		case pv := <-progress_chan:
			if pv.Error != nil {
				log.Println(pv.Error)
				continue
			}
			searchtaskindb := Searchtaskindb{}
			query_task := Searchtaskindb{
				Starttime: pv.search_task.start_t,
				Endtime:   pv.search_task.end_t,
				Assetid:   pv.search_task.asset_id,
			}
			db.Where(&query_task).First(&searchtaskindb)
			if searchtaskindb.CreatedAt.IsZero() {
				var this_record = Searchtaskindb{
					Starttime:       pv.search_task.start_t,
					Endtime:         pv.search_task.end_t,
					Lasttime:        pv.search_task.last_t,
					Yesterday2today: pv.search_task.yesterday2today,
					Assetid:         pv.search_task.asset_id,
					Ongoing:         pv.search_task.ongoing,
				}
				db.Create(&this_record)
			} else {
				searchtaskindb.Lasttime = pv.search_task.last_t
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
				}
				db.Create(&thisrecord)
			}
		case <-quit_chan:
			log.Println("finished")
			return
		case mixin_account := <-req_read_deposit_chan:
			for _, v := range default_asset_id_group {
				depositRecord := DepositAddressindb{
					Accountrecord_id: mixin_account.ID,
					Assetid:          v,
				}
				db.Create(&depositRecord)
				go read_asset_deposit_address(v, mixin_account.Userid, mixin_account.Sessionid, mixin_account.Privatekey, mixin_deposit_chan)
			}
		case new_user := <-mixin_account_chan:
			db.Create(&new_user)

			req_read_deposit_chan <- new_user
		case asset_deposit_address_result := <-mixin_deposit_chan:
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

		case unique_id := <-req_create_payment_chan:
			var notlinked_mixinaccount MixinAccountindb
			var result string
			db.Where("client_reqid = ?", "0").First(&notlinked_mixinaccount)
			if notlinked_mixinaccount.ID != 0 {
				new_req := ClientReq{
					Callbackurl:    unique_id,
					MixinAccountid: notlinked_mixinaccount.ID,
				}
				db.Create(&new_req)
				notlinked_mixinaccount.ClientReqid = new_req.ID
				db.Save(&notlinked_mixinaccount)
				result += fmt.Sprintf("new req created with record id: %v, user id: %v, with client request %v\n", notlinked_mixinaccount.ID, notlinked_mixinaccount.Userid, new_req.ID)
			} else {
				checkremain_account_c <- 1
				//no avaible mixin account, create one in blockin mode
				const predefine_pin string = "123456"
				user, err := mixin.CreateAppUser("tom", predefine_pin, user_config.user_id, user_config.session_id, user_config.private_key)
				if err != nil {
					log.Println(err)
				} else {
					new_user := MixinAccountindb{
						Userid:      user.UserId,
						Sessionid:   user.SessionId,
						Pintoken:    user.PinToken,
						Privatekey:  user.PrivateKey,
						Pin:         predefine_pin,
						ClientReqid: 0,
					}
					db.Create(&new_user)
					new_req := ClientReq{
						Callbackurl:    unique_id,
						MixinAccountid: new_user.ID,
					}
					db.Create(&new_req)
					new_user.ClientReqid = new_req.ID
					db.Save(&new_user)
					result += fmt.Sprintf("new req created with record id: %v, user id: %v, with client request %v\n", new_user.ID, new_user.Userid, new_req.ID)
					req_read_deposit_chan <- new_user
				}
			}
			user_output_chan <- result
		case <-timer1.C:
			checkremain_account_c <- 1
		case tocheck_account := <-checkaccount_deposit_c:
			for _, v := range default_asset_id_group {
				go read_asset_deposit_address(v, tocheck_account.Userid, tocheck_account.Sessionid, tocheck_account.Privatekey, mixin_deposit_chan)
			}
		case <-checkremain_account_c:
			var available_mixin_account int
			db.Model(&MixinAccountindb{}).Where("client_reqid = ?", "0").Count(&available_mixin_account)
			if available_mixin_account < 10 {
				for i := 20; i > available_mixin_account; i-- {
					const predefine_pin string = "123456"
					go create_mixin_account("tom", predefine_pin, user_config.user_id, user_config.session_id, user_config.private_key, mixin_account_chan)
				}
			}

		case v := <-user_cmd_chan:
			result := "\n"
			switch v {
			case "allsnap":
				var allsnap []Snapshotindb
				db.Find(&allsnap)
				for _, v := range allsnap {
					result += fmt.Sprintf("at %v with id: %v amount:%v asset %v to %v by %v\n", v.SnapCreatedAt, v.SnapshotId, v.Amount, v.AssetId, v.UserId, v.Source)
				}
			case "status":
				var alltask []Searchtaskindb
				db.Find(&alltask)
				total_ongoing := 0
				total_finished := 0
				for _, v := range alltask {
					if v.Ongoing {
						total_ongoing += 1
						result += fmt.Sprintf("search %v at %v from:%v to %v\n", v.Assetid, v.Lasttime, v.Starttime, v.Endtime)
					} else {
						total_finished += 1
					}
				}
				result += fmt.Sprintf("total %v ongoing", total_ongoing)
				result += fmt.Sprintf("total %v finished", total_finished)
			case "quit":
				quit_chan <- 1
			default:
				splited_string := strings.Split(v, " ")
				switch splited_string[0] {
				case "searchuser":
					user := splited_string[1]
					var users_snap []Snapshotindb
					db.Where(&Snapshotindb{UserId: user}).Find(&users_snap)
					for _, v := range users_snap {
						result += fmt.Sprintf("at %v with id: %v amount:%v asset %v to %v by %v\n", v.SnapCreatedAt, v.SnapshotId, v.Amount, v.AssetId, v.UserId, v.Source)
					}
				case "createreq":
					if len(splited_string) > 1 {
						req_create_payment_chan <- splited_string[1]
					}

				case "listreqs":
					var allreqs []ClientReq
					db.Find(&allreqs)
					for _, v := range allreqs {
						result += fmt.Sprintf("req id: %v %v %v\n", v.ID, v.Callbackurl, v.MixinAccountid)
					}
				case "searchreq":
					payment_id := splited_string[1]
					var req ClientReq
					db.Where(&ClientReq{Callbackurl: payment_id}).Find(&req)
					if req.ID != 0 {
						var mixin_account MixinAccountindb
						db.Find(&mixin_account, req.MixinAccountid)
						if mixin_account.ID != 0 {
							result += fmt.Sprintf("Record found : %v user id %v\n", req.Callbackurl, mixin_account.Userid)
							var payment_addresses []DepositAddressindb
							db.Where(&DepositAddressindb{Accountrecord_id: mixin_account.ID}).Find(&payment_addresses)
							for _, v := range payment_addresses {
								if v.Publicaddress != "" {
									result += fmt.Sprintf("Asset : %v Payment address %v\n", v.Assetid, v.Publicaddress)
								} else {
									result += fmt.Sprintf("Asset : %v Payment name %v tag %v\n", v.Assetid, v.Accountname, v.Accounttag)
								}
							}
						} else {
							result += fmt.Sprintf("Record found, but no payment channel is missing")
						}
						user_output_chan <- result
					} else {
						result += fmt.Sprintf("No matched record")
					}
				case "createuser":
					const predefine_pin string = "123456"
					go create_mixin_account("tom", predefine_pin, user_config.user_id, user_config.session_id, user_config.private_key, mixin_account_chan)
				case "listusers":
					var allaccount []MixinAccountindb
					db.Find(&allaccount)
					for _, v := range allaccount {
						result += fmt.Sprintf("user id: %v %v %v\n", v.ID, v.Userid, v.ClientReqid)
						var payment_addresses []DepositAddressindb
						db.Where(&DepositAddressindb{Accountrecord_id: v.ID}).Find(&payment_addresses)
						for _, add := range payment_addresses {
							if add.Publicaddress != "" {
								result += fmt.Sprintf("Asset : %v Payment address %v\n", add.Assetid, add.Publicaddress)
							} else {
								result += fmt.Sprintf("Asset : %v Payment name %v tag %v\n", add.Assetid, add.Accountname, add.Accounttag)
							}
						}
					}
				}
			}
			result += "allsnap: read all snap\n"
			result += "status: ongoing search task\n"
			result += "your selection:"
			user_output_chan <- result
		}
	}
}
