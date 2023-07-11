package common

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

type Result struct {
	Msg string
	C   *gin.Context
}

var Res Result

type Text struct {
	Content string `json:"content"`
}

type TextMsg struct {
	Touser  string `json:"touser"`
	Toparty string `json:"toparty"`
	Totag   string `json:"totag"`
	Msgtype string `json:"msgtype"`
	Agentid string `json:"agentid"`
	Text    Text   `json:"text"`
}

type SendToken struct {
	AccessToken string `json:"access_token"`
}

func getToken() string {
	var host = ""
	var param = map[string]string{
		"corpid":     "",
		"corpsecret": "",
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
	r := &SendToken{}
	json.Unmarshal(result, &r)
	return r.AccessToken
}

func SendMsg() {
	var host = ""
	var accessToken = getToken()

	uri, err := url.Parse(host)
	if err != nil {
		fmt.Println(err)
	}
	query := uri.Query()
	query.Set("access_token", accessToken)
	uri.RawQuery = query.Encode()
	textMsg := &TextMsg{
		Touser:  Res.C.GetString("Touser"),
		Toparty: "",
		Totag:   "",
		Msgtype: "",
		Agentid: "",
		Text: Text{
			Content: Res.Msg,
		},
	}

	sendBody, err := json.Marshal(textMsg)
	if err != nil {
		fmt.Println(err)
	}
	sendData := string(sendBody)
	client := &http.Client{}
	request, err := http.NewRequest("POST", uri.String(), strings.NewReader(sendData))
	if err != nil {
		fmt.Println(err)
	}
	request.Header.Set("Content-Type", "application/json")
	response, err := client.Do(request)
	defer response.Body.Close()
	result, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Println(err)
	}
	if len(Res.Msg) > 10 {
		Res.Msg = Res.Msg[:10] + "..."
	}
	fmt.Printf("res: %s, time: %s, to_user: %s, msg: %s\n", string(result), time.Now().Format("2006-01-02 15:04:05"), Res.C.GetString("Touser"), Res.Msg)
}

func In(target string, arr []string) bool {
	for _, element := range arr {
		if target == element {
			return true
		}
	}
	return false
}

func GetAIToken(ak, sk string) string {
	var host = "https://aip.baidubce.com/oauth/2.0/token"
	var param = map[string]string{
		"grant_type":    "client_credentials",
		"client_id":     ak,
		"client_secret": sk,
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
	res := make(map[string]string, 0)
	json.Unmarshal(result, &res)
	return res["access_token"]
}

func HttpGet(host string, header map[string]string, params map[string]string) ([]byte, error) {
	client := http.Client{}
	uri, err := url.Parse(host)
	if err != nil {
		fmt.Println(err)
	}
	query := uri.Query()
	for k, v := range params {
		query.Set(k, v)
	}
	uri.RawQuery = query.Encode()
	req, err := http.NewRequest(http.MethodGet, uri.String(), nil)
	if err != nil {
		return nil, err
	}
	// 添加请求头
	req.Header.Add("Content-type", "application/json;charset=utf-8")
	for k, v := range header {
		req.Header.Add(k, v)
	}
	// 发送请求
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return b, nil
}
