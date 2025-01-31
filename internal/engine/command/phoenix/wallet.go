package phoenix

import (
	"github.com/pagu-project/pagu/internal/engine/command"
	"github.com/pagu-project/pagu/internal/entity"
	"github.com/pagu-project/pagu/pkg/amount"
	"github.com/pagu-project/pagu/pkg/utils"
)

func (c *PhoenixCmd) walletHandler(
	_ *entity.User,
	cmd *command.Command,
	_ map[string]string,
) command.CommandResult {
	faucetAddress := utils.TestnetAddressToString(c.faucetAddress)
	balInt, err := c.client.GetBalance(c.ctx, faucetAddress)
	if err != nil {
		return cmd.RenderErrorTemplate(err)
	}

	return cmd.RenderResultTemplate(
		"address", faucetAddress,
		"balance", amount.Amount(balInt))
}
