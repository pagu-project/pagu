package phoenix

import (
	"github.com/pagu-project/pagu/internal/engine/command"
	"github.com/pagu-project/pagu/internal/entity"
)

//nolint:unused // remove me after I am used
func (pt *PhoenixCmd) walletHandler(cmd *command.Command,
	_ entity.PlatformID, _ string, _ ...string,
) command.CommandResult {
	return cmd.SuccessfulResultF(
		"Pagu Phoenix Address: %s\nBalance: %d", pt.wallet.Address(), pt.wallet.Balance())
}
