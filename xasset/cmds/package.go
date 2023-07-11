package cmds

import (
	"encoding/json"
	"fmt"
	"xuper-info/common"

	"github.com/spf13/cobra"
)

type XaBalance struct {
	Account   XaAccount `json:"account"`
	Errmsg    string    `json:"errmsg"`
	Errno     int       `json:"errno"`
	RequestId string    `json:"request_id"`
}

type XaAccount struct {
	AccountId string `json:"account_id"`
	Realname  string `json:"realname"`
	Approved  bool   `json:"approved"`
	Business  string `json:"business"`
	AppId     string `json:"app_id"`
	Ak        string `json:"ak"`
	Sk        string `json:"sk"`
	Pause     bool   `json:"pause"`
	Status    string `json:"status"`
	Package   int    `json:"package"`
	Ctime     int    `json:"ctime"`
}

func NewPackageCmd() *cobra.Command {
	cm := &cobra.Command{
		Use:   "package",
		Short: "query package for accountId",
	}
	cm.AddCommand(newXaBalanceCmd())
	return cm
}

func newXaBalanceCmd() *cobra.Command {
	cm := &cobra.Command{
		Use:   "balance [accountId]",
		Short: "query transaction based on txid",
		Run: func(c *cobra.Command, args []string) {
			//args[0]
			host := ""
			header := map[string]string{
				"X-Bce-Account": args[0],
				"X-User-Id":     args[0],
			}
			b, err := common.HttpGet(host, header, nil)
			if err != nil {
				common.Res.Msg = err.Error()
			}
			var result XaBalance
			json.Unmarshal(b, &result)
			if result.Errno != 0 {
				common.Res.Msg = result.Errmsg
			} else if result.Account.Business == "" {
				common.Res.Msg = "该用户不存在，请检查参数"
			} else {
				common.Res.Msg = fmt.Sprintf("用户「%s」量包余额为「%v」次", result.Account.Business, result.Account.Package)
			}
			common.SendMsg()
		},
	}
	return cm
}
