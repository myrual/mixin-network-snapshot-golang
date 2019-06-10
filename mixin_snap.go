package main

import (
	"encoding/json"
	"log"
	"time"

	mixin "github.com/MooooonStar/mixin-sdk-go/network"
)

const (
	userid      = "a932cac1-e05b-4095-b662-f5ab284050bf"
	sessionid   = "9788a620-9884-442c-bc4f-403099230971"
	private_key = `-----BEGIN RSA PRIVATE KEY-----
MIICXAIBAAKBgQDb/2/TvL1lB779WvKxnkCI87arFsIIEQn/QkiUIKg+R5m+//Ac
TrF8mk2j4qlTHrDDYAp4gSqUq6uJTglZpPAbWbSkQT5guK7BLa1pAIuWC6DSsRDR
5K2TxJcDfQQ4ajJElgEbJvfh7KJHrWtwIpmNHuf2SWAoUt2lkgYx+pxw9wIDAQAB
AoGAeWFbEsj8+kQm1WSbqPI2ixlkfMsmrQqxqFuq7ceC7DHXGzHdCdFFTglYswZ5
A/pX7sQLoucpVjPNhgk/UW2WZ4rYXtw/tf7c8hqL0SV+EAosAx3Hjbbmrr1lUfaH
lBeq9pPImwOm343s4cbFFJFokzxyaWVf9TCZwmM7hDqM75ECQQDvuT1TfZXP6dbl
BhQA0bXeitwmNocH0saBf/LXV+JLC1EyOJurSVvlchG5p53lmQ/wQ6Z6voz/MT3v
YVo0R+u/AkEA6u9VQMgMN3u5WC4n120XTiCGLgC/I4l6a5fsNVHApKulf8hdacGh
bbG7QtHHwbSYeY5fV4KrwVVfVnrK3FGoyQJAI1LfZ4MU5Tsm0D6SCgDc1LsPb44P
XabAW2q4JOUtUjOLtmPDBH1dzjR9yiaZzLA+OgAt8t5LNntSDgkBWrzSTwJAYwaO
gMfRnnFgJnMOCBfLgvrik/Fsn6YLG97liXP0J3TSRZJHDZS4XmxT6k5STKu6uUHx
nglOLCe4D9OiPkuNQQJBAJ9VuX6XiBZfy9pf/2Bi1eWwS0n6PhilBy7OGuoQMvie
KQ8cFV0xyyWQx3oEnv7vZOebud65vDC6ZI2F25cnQ3A=
-----END RSA PRIVATE KEY-----`
)

type Snapshot struct {
	SnapshotId string `json:"snapshot_id"`
	Amount     string `json:"amount"`
	Asset      struct {
		AssetId string `json:"asset_id"`
	} `json:"asset"`
	CreatedAt time.Time `json:"created_at"`

	TraceId    string `json:"trace_id"`
	UserId     string `json:"user_id"`
	OpponentId string `json:"opponent_id"`
	Data       string `json:"data"`
}
type BotConfig struct {
	user_id     string
	session_id  string
	private_key string
}

func searchSnapshot(task Searchtask, result_chan chan SnapNetResponse, config BotConfig) {
	snaps, err := mixin.NetworkSnapshots(task.asset_id, task.start_t, true, task.max_len, config.user_id, config.session_id, config.private_key)

	if err != nil {
		result_chan <- SnapNetResponse{
			Error:    err,
			MixinReq: task,
		}
		return
	}

	var resp MixinResponse
	err = json.Unmarshal(snaps, &resp)

	if err != nil {
		result_chan <- SnapNetResponse{
			Error:    err,
			MixinReq: task,
		}
		return
	}
	result_chan <- SnapNetResponse{
		MixinRespone: resp,
		MixinReq:     task,
	}
}

type Searchtask struct {
	start_t  time.Time
	end_t    time.Time
	max_len  int
	asset_id string
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
)

type SnapNetResponse struct {
	Error        error
	MixinReq     Searchtask
	MixinRespone MixinResponse
}

type MixinResponse struct {
	Data  []*Snapshot `json:"data"`
	Error string      `json:"error"`
}

func main() {
	var start_time2 = time.Date(2019, 6, 1, 0, 0, 0, 0, time.UTC)
	var network_result_chan = make(chan SnapNetResponse, 100)
	var task_chan = make(chan Searchtask, 100)
	var quit_chan = make(chan int, 2)

	var user_config = BotConfig{
		user_id:     userid,
		session_id:  sessionid,
		private_key: private_key,
	}
	task_chan <- Searchtask{
		start_t:  start_time2,
		end_t:    start_time2.AddDate(0, 0, 1),
		max_len:  500,
		asset_id: USDT_ASSET_ID,
	}
	total_task := len(task_chan)
	log.Println("go with ", total_task, " tasks")
	now := time.Now()
	for {
		select {
		case task := <-task_chan:
			log.Println(task.start_t, task.max_len, " for ", task.asset_id)
			go searchSnapshot(task, network_result_chan, user_config)

		case v := <-network_result_chan:
			total_task -= 1
			if v.Error != nil {
				log.Println("Net work error ", v.Error, " for req:", v.MixinReq.asset_id, " start ", v.MixinReq.start_t)
			} else {
				if v.MixinRespone.Error != "" {
					log.Println("Server return error", v.MixinRespone.Error, " for req:", v.MixinReq.asset_id, " start ", v.MixinReq.start_t)
				} else {
					len_of_snap := len(v.MixinRespone.Data)
					last_element := v.MixinRespone.Data[len(v.MixinRespone.Data)-1]
					log.Println("the last element is created at:", last_element.CreatedAt)
					if len_of_snap < 500 {
						log.Println("no enough record to search, pause")
						break
					} else {
						if last_element.CreatedAt.After(v.MixinReq.end_t) {
							log.Println("reach ", v.MixinReq.end_t)
							log.Println("total ", time.Now().Sub(now), " passed")
							return
						}
						task_chan <- Searchtask{
							start_t:  last_element.CreatedAt,
							end_t:    v.MixinReq.end_t,
							asset_id: v.MixinReq.asset_id,
							max_len:  v.MixinReq.max_len,
						}
						total_task += 1
						log.Println("search again ", last_element.CreatedAt)
					}
				}
			}

			if total_task == 0 {
				log.Println("finish all search")
				quit_chan <- 0
			}

		case <-quit_chan:
			return
		}
	}
}
