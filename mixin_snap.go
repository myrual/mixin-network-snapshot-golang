package main

import (
	"encoding/json"
	"fmt"
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

func searchSnapshot(start_t time.Time, end_t time.Time, c chan *Snapshot, status_chan chan ScanProgress, config BotConfig) {
	var start_time = start_t
	log.Println(start_t, end_t)
	for {

		snaps, err := mixin.NetworkSnapshots("", start_time, true, 500, config.user_id, config.session_id, config.private_key)

		if err != nil {
			var errorStatus = ScanProgress{
				error_value:        err,
				error_stopped_time: start_time,
			}
			status_chan <- errorStatus
			return
		}

		var resp struct {
			Data  []*Snapshot `json:"data"`
			Error string      `json:"error"`
		}
		err = json.Unmarshal(snaps, &resp)

		if err != nil {
			var errorStatus = ScanProgress{
				error_value: err,
			}
			status_chan <- errorStatus
			return
		}

		if resp.Error != "" {
			fmt.Println("error in result")
			log.Fatal("read snapshot error", resp.Error)
			return
		}

		lastElement := resp.Data[len(resp.Data)-1]
		for _, v := range resp.Data {
			c <- v
		}
		var scanStatus = ScanProgress{
			lastest_scanned_time: lastElement.CreatedAt,
			status:               true,
		}
		status_chan <- scanStatus
		if lastElement.CreatedAt.After(end_t) {
			return
		}
		fmt.Println(lastElement.CreatedAt)
		start_time = lastElement.CreatedAt
	}
}

type Searchtask struct {
	start_t  time.Time
	end_t    time.Time
	asset_id string
}

type ScanProgress struct {
	lastest_scanned_time time.Time
	status               bool
	error_value          error
	error_stopped_time   time.Time
}

func taskReceiver(task_c chan Searchtask, result_c chan *Snapshot, quit_c chan int, status_c chan ScanProgress, config BotConfig) {
	for {
		select {
		case task := <-task_c:
			fmt.Println(task.start_t)
			go searchSnapshot(task.start_t, task.end_t, result_c, status_c, config)
		case <-quit_c:
			return
		}
	}
}
func snapReceiver(result_c chan *Snapshot, quit_c chan int) {
	for {
		select {
		case <-result_c:

		case <-quit_c:
			return
		}
	}
}

func create_task(start_time2 time.Time, end_time2 time.Time, c chan Searchtask) {
	var i int = 0
	const duration int = 720
	for {
		this_start := start_time2.Add(time.Minute * time.Duration(duration*i))
		this_end := this_start.Add(time.Minute * time.Duration(duration))
		if end_time2.After(this_end) {
			c <- Searchtask{start_t: this_start, end_t: this_end}
		} else {
			break
		}
		i += 1
	}
}

func main() {
	var start_time2 = time.Date(2018, 8, 11, 0, 0, 0, 0, time.UTC)
	var end_time2 = time.Date(2018, 8, 13, 0, 0, 0, 0, time.UTC)
	var snaps_chan = make(chan *Snapshot)
	var task_chan = make(chan Searchtask, 10)
	var quit_chan = make(chan int)
	var status_chan = make(chan ScanProgress, 10)

	var user_config = BotConfig{
		user_id:     userid,
		session_id:  sessionid,
		private_key: private_key,
	}

	go create_task(start_time2, end_time2, task_chan)
	go taskReceiver(task_chan, snaps_chan, quit_chan, status_chan, user_config)
	go snapReceiver(snaps_chan, quit_chan)
	fmt.Println("wait")
	for {
		select {
		case v := <-status_chan:
			if v.error_value != nil {
				log.Println(v.error_value)
			}
			if v.status == true {
				log.Println(v.lastest_scanned_time)
			}
		}
	}
}
