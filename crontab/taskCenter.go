package crontab

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"sort"
	"time"
	"xuper-info/common"

	"github.com/gin-gonic/gin"
	"github.com/robfig/cron/v3"
)

type watchXuperOSResponse struct {
	Code int         `json:"code"`
	Data map[int]int `json:"data"`
}

// 定时任务中心
func Run() {
	// 测试消息
	WatchXuperOSTxs()
	//定时组件
	c := cron.New(cron.WithSeconds())
	// c.AddFunc("*0 0 9 ? * MON-FRI", fangKe)
	c.AddFunc("0 0 9 * * 4", WatchXuperOSTxs)
	go c.Start()
	defer c.Stop()
	select {}
}

// func fangKe() {
// 	common.Res.Msg = "辛苦帮忙预约访客~\n* 鲁忠\n* 尹晓怡\n* 孙延福\n* 王钧辉"
// 	common.Res.C = new(gin.Context)
// 	common.Res.C.Set("Touser", "jingqi03")
// 	common.SendMsg()
// }

func WatchXuperOSTxs() {
	var host = ""
	var param = map[string]string{
		"net_id": "5",
		"chain":  "xuper",
	}
	uri, err := url.Parse(host)
	if err != nil {
		fmt.Println(err)
	}
	query := uri.Query()
	for k, v := range param {
		query.Set(k, v)
	}
	uri.RawQuery = query.Encode()

	response, err := http.Get(uri.String())
	if err != nil {
		fmt.Println(err)
	}
	result, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Println(err)
	}
	var res watchXuperOSResponse
	json.Unmarshal(result, &res)
	keys := make([]int, 0)
	str := ""
	if res.Code == 0 {
		for k := range res.Data {
			keys = append(keys, k)
		}
		sort.Ints(keys)
		keys = keys[len(keys)-7:]
		value := 0
		sum := 0
		for i := 0; i < len(keys); i++ {
			value = res.Data[keys[i]]
			sum += value
			str += fmt.Sprintf("* %s 交易数: %v\n", time.Unix(int64(keys[i]), 10).Format("2006-01-02"), value)
		}
		str += fmt.Sprintf("* 7日内总交易: %v", sum)
		common.Res.Msg = str
		common.Res.C = new(gin.Context)
		common.Res.C.Set("Touser", "")
		common.SendMsg()
	}
}
