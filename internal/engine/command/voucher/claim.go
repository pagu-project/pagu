package voucher

import (
	"errors"
	"fmt"
	"time"

	"github.com/pagu-project/pagu/internal/engine/command"
	"github.com/pagu-project/pagu/internal/entity"
	"github.com/pagu-project/pagu/pkg/log"
)

func (c *VoucherCmd) claimHandler(
	caller *entity.User,
	cmd *command.Command,
	args map[string]string,
) command.CommandResult {
	code := args[argNameClaimCode]
	if len(code) != 8 {
		return cmd.RenderFailedTemplateF("Voucher code is not valid, length must be 8: `%v`", code)
	}

	voucher, err := c.db.GetVoucherByCode(code)
	if err != nil {
		return cmd.RenderFailedTemplateF("Voucher code is not valid, no voucher found: `%v`", code)
	}

	if voucher.CreatedAt.AddDate(0, int(voucher.ValidMonths), 0).Before(time.Now()) {
		return cmd.RenderFailedTemplate("Voucher is expired")
	}

	if voucher.IsClaimed() {
		return cmd.RenderFailedTemplate("Voucher code claimed before")
	}

	address := args[argNameClaimAddress]

	var txHash string
	switch voucher.Type {
	case entity.VoucherTypeStake:
		valInfo, _ := c.clientManager.GetValidatorInfo(address)
		if valInfo != nil {
			err = errors.New("this address is already a staked validator")
			log.Warn("Staked validator found", "address", address)

			return cmd.RenderErrorTemplate(err)
		}

		pubKey, err := c.clientManager.FindPublicKey(address, false)
		if err != nil {
			log.Warn("Peer not found", "address", address)

			return cmd.RenderErrorTemplate(err)
		}

		memo := fmt.Sprintf("Voucher %s claimed by Pagu", code)
		txHash, err = c.wallet.BondTransaction(pubKey, address, voucher.Amount, memo)
		if err != nil {
			return cmd.RenderErrorTemplate(err)
		}

		if txHash == "" {
			return cmd.RenderFailedTemplate("Can't send bond transaction")
		}

	case voucher.Type:
		memo := fmt.Sprintf("Voucher %s claimed by Pagu", code)
		txHash, err = c.wallet.TransferTransaction(address, voucher.Amount, memo)
		if err != nil {
			return cmd.RenderErrorTemplate(err)
		}

		if txHash == "" {
			return cmd.RenderFailedTemplate("Can't send transfer transaction")
		}

	default:
		return cmd.RenderFailedTemplateF("Invalid Voucher type: %d", voucher.Type)
	}

	if err = c.db.ClaimVoucher(voucher.ID, txHash, caller.ID); err != nil {
		return cmd.RenderErrorTemplate(err)
	}

	return cmd.RenderResultTemplate("txHash", txHash)
}
