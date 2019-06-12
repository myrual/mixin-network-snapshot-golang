package main

import (
	"encoding/json"
	"log"
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
)

func searchSnapshot(asset_id string, start_t time.Time, yesterday2today bool, max_len int, config BotConfig) SnapNetResponse {
	snaps, err := mixin.NetworkSnapshots(asset_id, start_t, yesterday2today, max_len, config.user_id, config.session_id, config.private_key)

	if err != nil {
		return SnapNetResponse{
			Error: err,
		}
	}

	var resp MixinResponse
	err = json.Unmarshal(snaps, &resp)

	if err != nil {
		return SnapNetResponse{
			Error: err,
		}
	}
	return SnapNetResponse{
		MixinRespone: resp,
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
		v := searchSnapshot(req_task.asset_id, req_task.last_t, req_task.yesterday2today, req_task.max_len, user_config)
		if v.Error != nil {
			progress_chan <- Searchprogress{
				Error: v.Error,
			}
			continue
		}
		if v.MixinRespone.Error != "" {
			log.Println("Server return error", v.MixinRespone.Error, " for req:", req_task.asset_id, " start ", req_task.start_t)
			return
		}
		for _, v := range v.MixinRespone.Data {
			if v.UserId != "" {
				result_chan <- v
			}
		}
		len_of_snap := len(v.MixinRespone.Data)
		if len_of_snap == 0 {
			time.Sleep(60 * time.Second)
			continue
		}
		last_element := v.MixinRespone.Data[len(v.MixinRespone.Data)-1]
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

func main() {
	var start_time2 = time.Date(2018, 4, 25, 0, 0, 0, 0, time.UTC)
	var my_snapshot_chan = make(chan *Snapshot, 1000)
	var progress_chan = make(chan Searchprogress, 1000)
	var quit_chan = make(chan int, 2)

	db, err := gorm.Open("sqlite3", "test.db")
	if err != nil {
		panic("failed to connect database")
	}
	defer db.Close()

	db.AutoMigrate(&Snapshotindb{})
	db.AutoMigrate(&Searchtaskindb{})

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

	total_found_snap := 0
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
				db.Model(&searchtaskindb).Update(Searchtaskindb{Lasttime: pv.search_task.last_t, Ongoing: pv.search_task.ongoing})
			}
			log.Println(pv.search_task.ongoing, pv.search_task.last_t)
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
			total_found_snap += 1
			if total_found_snap%100 == 0 {
				log.Println(total_found_snap, v.SnapshotId, v.CreatedAt)
			}
		case <-quit_chan:
			log.Println("finished")
			return
		}
	}
}
