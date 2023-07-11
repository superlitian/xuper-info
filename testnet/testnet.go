package testnet

import (
	"xuper-info/common"
	"xuper-info/xuperos/cmds"

	"github.com/spf13/cobra"
	"github.com/xuperchain/xuper-sdk-go/v2/xuper"
)

func TestNetExecute() {

	client, err := xuper.New("14.215.179.74:37101", xuper.WithConfigFile("./conf/sdk.testnet.yaml"))
	if err != nil {
		common.Res.Msg = err.Error()
		common.SendMsg()
	}
	osCmd := &cobra.Command{
		Use: "testnet",
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
