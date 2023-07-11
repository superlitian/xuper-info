package cmds

import (
	"encoding/json"
	"xuper-info/common"

	"github.com/spf13/cobra"
	"github.com/xuperchain/xuper-sdk-go/v2/xuper"
)

func NewContractCmd(cli *xuper.XClient) *cobra.Command {
	var address string
	var account string
	cmd := &cobra.Command{
		Use: "contracts",
		RunE: func(cmd *cobra.Command, args []string) error {
			if address != "" {
				b, err := cli.QueryAddressContracts(address)
				if err != nil {
					return err
				} else {
					common.Res.Msg = b[address].String()
				}
			} else if account != "" {
				b, err := cli.QueryAccountContracts(account)
				if err != nil {
					return err
				} else {
					bytes, _ := json.Marshal(b)
					str := string(bytes)
					common.Res.Msg = str
				}
			} else {
				common.Res.Msg = "this query must use '--address' or '--account' option"
			}
			common.SendMsg()
			return nil
		},
	}

	cmd.Flags().StringVar(&address, "address", "", "address name to query contracts")
	cmd.Flags().StringVar(&account, "account", "", "account address to query contracts")
	return cmd
}
