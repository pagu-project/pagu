package phoenix

import (
	"github.com/pactus-project/pactus/types/tx"
	"github.com/pagu-project/pagu/internal/engine/command"
	"github.com/pagu-project/pagu/internal/entity"
	"github.com/pagu-project/pagu/pkg/utils"
)

func (c *PhoenixCmd) faucetHandler(
	caller *entity.User,
	cmd *command.Command,
	args map[string]string,
) command.CommandResult {
	if len(args) == 0 {
		return cmd.RenderFailedTemplate("Please provide a valid address to receive the faucet amount.")
	}

	lastFaucet, _ := c.db.GetLastFaucet(caller)
	if lastFaucet != nil {
		timeSinceLastFaucet := lastFaucet.ElapsedTime()
		if timeSinceLastFaucet <= c.faucetCooldown {
			timeToWait := c.faucetCooldown - timeSinceLastFaucet
			return cmd.RenderFailedTemplateF(
				"You have already used your faucet share. Please try again in %s.",
				utils.FormatDuration(timeToWait))
		}
	}

	blockchainInfo, err := c.client.GetBlockchainInfo(c.ctx)
	if err != nil {
		return cmd.RenderErrorTemplate(err)
	}

	toAddr := args["address"]
	receiverAddress, err := utils.TestnetAddressFromString(toAddr)
	if err != nil {
		return cmd.RenderErrorTemplate(err)
	}

	trx := tx.NewTransferTx(blockchainInfo.LastBlockHeight, c.faucetAddress, receiverAddress,
		c.faucetAmount.ToPactusAmount(),
		c.faucetFee.ToPactusAmount(),
		tx.WithMemo("Pagu Phoenix Faucet"))

	signBytes := trx.SignBytes()

	sig := c.privateKey.SignNative(signBytes)
	trx.SetSignature(sig)
	trx.SetPublicKey(c.privateKey.PublicKey())

	trxData, _ := trx.Bytes()
	txHash, err := c.client.BroadcastTransaction(c.ctx, trxData)
	if err != nil {
		return cmd.RenderErrorTemplate(err)
	}

	if err := c.db.AddFaucet(&entity.PhoenixFaucet{
		UserID:  caller.ID,
		Address: toAddr,
		Amount:  c.faucetAmount,
		TxHash:  txHash,
	}); err != nil {
		return cmd.RenderErrorTemplate(err)
	}

	return cmd.RenderResultTemplate(
		"amount", c.faucetAmount,
		"txHash", txHash)
}
