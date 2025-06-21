package network

import (
	"github.com/pagu-project/pagu/internal/engine/command"
	"github.com/pagu-project/pagu/internal/entity"
)

func (c *NetworkCmd) walletHandler(
	_ *entity.User,
	cmd *command.Command,
	_ map[string]string,
) command.CommandResult {
	return cmd.RenderResultTemplate(
		"address", c.wallet.Address(),
		"balance", c.wallet.Balance())
}
