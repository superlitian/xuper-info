package cmds

import (
	"encoding/json"
	"xuper-info/common"
	"xuper-info/crontab"

	"github.com/spf13/cobra"
	"github.com/xuperchain/xuper-sdk-go/v2/xuper"
	"github.com/xuperchain/xuperchain/cmd/client/cmd"
)

func NewTxCmd(cli *xuper.XClient) *cobra.Command {
	cm := &cobra.Command{
		Use:   "tx",
		Short: "Operate tx command, query",
	}

	cm.AddCommand(newTxQueryCmd(cli))
	cm.AddCommand(NewTxTotalCmd())
	return cm
}

// NewTxQueryCmd query tx
func newTxQueryCmd(cli *xuper.XClient) *cobra.Command {
	cm := &cobra.Command{
		Use:   "query txid",
		Short: "query transaction based on txid",
		Run: func(c *cobra.Command, args []string) {
			t, err := cli.QueryTxByID(args[0])
			if err != nil {
				common.Res.Msg = err.Error()
			}
			tx := cmd.FromPBTx(t)
			output, err := json.MarshalIndent(tx, "", "  ")
			if err != nil {
				common.Res.Msg = err.Error()
			} else {
				common.Res.Msg = string(output)
			}
			common.SendMsg()
		},
	}

	return cm
}

func NewTxTotalCmd() *cobra.Command {
	cm := &cobra.Command{
		Use:   "total",
		Short: "query transaction total for last 7d",
		Run: func(c *cobra.Command, args []string) {
			crontab.WatchXuperOSTxs()
		},
	}
	return cm
}
