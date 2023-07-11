package cmds

import "C"
import (
	"github.com/spf13/cobra"
	"github.com/xuperchain/xuper-sdk-go/v2/xuper"
	"xuper-info/common"
)

func NewBalanceCmd(cli *xuper.XClient) *cobra.Command {
	cmd := &cobra.Command{
		Use:"balance [account/address]",

		RunE: func(cmd *cobra.Command, args []string) error {
			b, err := cli.QueryBalance(args[0])
			if err != nil{
				return err
			}else {
				common.Res.Msg = b.String()
			}
			common.SendMsg()
			return nil
		},
	}
	return cmd
}
