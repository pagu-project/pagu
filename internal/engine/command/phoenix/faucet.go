package phoenix

import (
	"github.com/pagu-project/pagu/internal/engine/command"
	"github.com/pagu-project/pagu/internal/entity"
)

func (p *PhoenixCmd) faucetHandler(
	caller *entity.User,
	cmd *command.Command,
	args map[string]string,
) command.CommandResult {
	// if len(args) == 0 {
	// 	return cmd.RenderFailedTemplate("Invalid wallet address")
	// }

	// toAddr := args["address"]
	// if len(toAddr) != 43 || toAddr[:3] != "tpc" {
	// 	return cmd.RenderFailedTemplate("Invalid wallet address")
	// }

	// if !p.db.CanGetFaucet(caller) {
	// 	return cmd.RenderFailedTemplate("Uh, you used your share of faucets today!")
	// }

	// txHash, err := p.wallet.TransferTransaction(toAddr, "Pagu Phoenix Faucet", p.faucetAmount)
	// if err != nil {
	// 	return cmd.RenderErrorTemplate(err)
	// }

	// if err = p.db.AddFaucet(&entity.PhoenixFaucet{
	// 	UserID:  caller.ID,
	// 	Address: toAddr,
	// 	Amount:  p.faucetAmount,
	// 	TxHash:  txHash,
	// }); err != nil {
	// 	return cmd.RenderErrorTemplate(err)
	// }

	return cmd.RenderResultTemplate(
		"amount", p.faucetAmount.ToPAC(),
		"txHash", "txHash")
}
