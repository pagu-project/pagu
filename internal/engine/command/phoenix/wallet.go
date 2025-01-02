package phoenix

import (
	"github.com/pagu-project/pagu/internal/engine/command"
	"github.com/pagu-project/pagu/internal/entity"
)

func (p *PhoenixCmd) walletHandler(
	_ *entity.User,
	cmd *command.Command,
	_ map[string]string,
) command.CommandResult {
	return cmd.SuccessfulResultF(
		"Pagu Phoenix Wallet Address: %s\nBalance: %d", p.wallet.Address(), p.wallet.Balance())
}
