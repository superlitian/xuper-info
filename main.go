package main

import (
	"crypto/md5"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"xuper-info/ai"
	"xuper-info/common"
	"xuper-info/testnet"
	"xuper-info/util"
	"xuper-info/xasset"
	"xuper-info/xuperos"

	"github.com/gin-gonic/gin"
)

type Auth struct {
	Signature string `json:"signature" form:"signature"`
	Timestamp string `json:"timestamp" form:"timestamp"`
	Rn        string `json:"rn" form:"rn"`
	Echostr   string `json:"echostr" form:"echostr"`
}
type Message struct {
	FromUserId   string `json:"FromUserId"`
	FromUserName string `json:"FromUserName"`
	CreateTime   string `json:"CreateTime"`
	MsgType      string `json:"MsgType"`
	Content      string `json:"Content"`
}

type MessageBody struct {
	ToUserName string `json:"ToUserName"`
	AgentID    int    `json:"AgentId"`
	Encrypt    string `json:"Encrypt"`
}

const (
	//6Vrhv09zlZiczO9jCiJeCExals 测试
	//bAjhbmsNJoxpGT7p 正式
	accessToken    = "bAjhbmsNJoxpGT7p"
	encodingAESKey = "rtmfeR6kA8okl8Tn9OQ8wX"
)

var (
	// 命令集
	rootCmds = []string{"/xuperos", "/xasset", "/testnet"}
	// 关键词过滤
	words = []string{"help"}
)

func main() {
	r := gin.Default()
	// 如流消息入口
	r.POST("/receive", receive)
	r.Run()
}

var m Message

func receive(c *gin.Context) {
	s := c.PostForm("messageJson")
	if s != "" {
		c.JSON(http.StatusOK, gin.H{
			"msg": "success",
		})
		var msg MessageBody
		json.Unmarshal([]byte(s), &msg)
		b, _ := util.Base64URLDecode(msg.Encrypt)
		AESKey, _ := base64.StdEncoding.DecodeString(encodingAESKey + "==")
		str := util.AesDecrypt(b, AESKey)
		json.Unmarshal(str, &m)
		fmt.Printf("receive fromUser: %s,  msg: %s\n", m.FromUserId, m.Content)
		c.Set("Touser", m.FromUserId)
		common.Res.C = c
		if m.MsgType != "event" {
			flags := strings.Fields(m.Content)
			os.Args = os.Args[:0]
			//将字符串切片追加到os.Args中
			for _, value := range flags {
				os.Args = append(os.Args, value)
			}
			if !common.In(os.Args[0], rootCmds) && !common.In(os.Args[0], words) {
				go func() {
					ai.Unit(m.Content)
				}()
			} else {
				switch os.Args[0] {
				case "/xuperos":
					xuperos.OsExecute()
				case "/testnet":
					testnet.TestNetExecute()
				case "/xasset":
					xasset.XaExecute()
				}
			}
		}
	} else {
		var auth Auth
		err := c.ShouldBind(&auth)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
		}
		if authorize(auth) {
			c.String(http.StatusOK, auth.Echostr)
		}
	}
}

func authorize(auth Auth) bool {
	s := auth.Rn + auth.Timestamp + accessToken
	srcCode := md5.Sum([]byte(s))
	if fmt.Sprintf("%x", srcCode) == auth.Signature {
		return true
	} else {
		return false
	}
}
