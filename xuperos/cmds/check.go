package cmds

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"xuper-info/common"

	"github.com/spf13/cobra"
)

func NewCheck() *cobra.Command {
	cmd := &cobra.Command{
		Use: "check [text]",

		RunE: func(cmd *cobra.Command, args []string) error {
			b, err := Check(args[0])
			if err != nil {
				return err
			} else {
				common.Res.Msg = b
			}
			common.SendMsg()
			return nil
		},
	}
	return cmd
}

// check 文本检查
func Check(text string) (string, error) {
	var host = "https://aip.baidubce.com/rest/2.0/solution/v1/text_censor/v2/user_defined"
	var accessToken = common.GetAIToken("", "")
	uri, err := url.Parse(host)
	if err != nil {
		fmt.Println(err)
	}
	query := uri.Query()
	query.Set("access_token", accessToken)
	uri.RawQuery = query.Encode()
	sendBody := http.Request{}
	sendBody.ParseForm()
	sendBody.Form.Add("text", text)
	sendData := sendBody.Form.Encode()
	client := &http.Client{}
	request, err := http.NewRequest("POST", uri.String(), strings.NewReader(sendData))
	if err != nil {
		fmt.Println(err)
	}
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	response, err := client.Do(request)
	defer response.Body.Close()
	result, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return "", err
	}
	return string(result), nil
}
