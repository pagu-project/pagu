package phoenix

import (
	"errors"

	"github.com/pagu-project/pagu/internal/engine/command"
	"github.com/pagu-project/pagu/internal/entity"
)

func (pt *Phoenix) faucetHandler(
	caller *entity.User,
	cmd *command.Command,
	args map[string]string,
) command.CommandResult {
	if len(args) == 0 {
		return cmd.ErrorResult(errors.New("invalid wallet address"))
	}

	toAddr := args["address"]
	if len(toAddr) != 43 || toAddr[:3] != "tpc" {
		return cmd.ErrorResult(errors.New("invalid wallet address"))
	}

	if !pt.db.CanGetFaucet(caller) {
		return cmd.FailedResult("Uh, you used your share of faucets today!")
	}

	txHash, err := pt.wallet.TransferTransaction(toAddr, "Phoenix Testnet Pagu PhoenixFaucet", pt.faucetAmount)
	if err != nil {
		return cmd.ErrorResult(err)
	}

	if err = pt.db.AddFaucet(&entity.PhoenixFaucet{
		UserID:  caller.ID,
		Address: toAddr,
		Amount:  pt.faucetAmount,
		TxHash:  txHash,
	}); err != nil {
		return cmd.ErrorResult(err)
	}

	return cmd.SuccessfulResultF("You got %f tPAC on Phoenix Testnet!\n\n"+
		"https://phoenix.pacviewer.com/transaction/%s", pt.faucetAmount.ToPAC(), txHash)
}
