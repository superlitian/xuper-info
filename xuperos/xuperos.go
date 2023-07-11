package xuperos

import (
	"xuper-info/common"
	"xuper-info/xuperos/cmds"

	"github.com/spf13/cobra"
	"github.com/xuperchain/xuper-sdk-go/v2/xuper"
)

func OsExecute() {

	client, err := xuper.New("39.156.69.83:37100", xuper.WithConfigFile("./conf/sdk.yaml"))
	if err != nil {
		common.Res.Msg = err.Error()
		common.SendMsg()
	}
	osCmd := &cobra.Command{
		Use: "xuperos",
	}
	//查询余额
	osCmd.AddCommand(cmds.NewBalanceCmd(client))
	//查询交易
	osCmd.AddCommand(cmds.NewTxCmd(client))
	//背书服务
	osCmd.AddCommand(cmds.NewCheck())
	// 查询合约
	osCmd.AddCommand(cmds.NewContractCmd(client))
	osCmd.DisableSuggestions = true
	err = osCmd.Execute()
	if err != nil {
		common.Res.Msg = err.Error()
		common.SendMsg()
	}
}
