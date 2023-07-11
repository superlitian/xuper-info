package xasset

import (
	"github.com/spf13/cobra"
	"xuper-info/common"
	"xuper-info/xasset/cmds"
)

func XaExecute() {
	xaCmd := &cobra.Command{
		Use: "xasset",
	}
	//查询余额
	xaCmd.AddCommand(cmds.NewPackageCmd())

	xaCmd.DisableSuggestions = true
	err := xaCmd.Execute()
	if err !=nil{
		common.Res.Msg = err.Error()
		common.SendMsg()
	}
}