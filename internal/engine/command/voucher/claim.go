package voucher

import (
	"errors"
	"fmt"
	"time"

	"github.com/pagu-project/pagu/internal/engine/command"
	"github.com/pagu-project/pagu/internal/entity"
	"github.com/pagu-project/pagu/pkg/log"
)

func (v *VoucherCmd) claimHandler(
	caller *entity.User,
	cmd *command.Command,
	args map[string]string,
) command.CommandResult {
	code := args["code"]
	if len(code) != 8 {
		return cmd.RenderFailedTemplate("Voucher code is not valid, length must be 8")
	}

	voucher, err := v.db.GetVoucherByCode(code)
	if err != nil {
		return cmd.RenderFailedTemplate("Voucher code is not valid, no voucher found")
	}

	if voucher.CreatedAt.AddDate(0, int(voucher.ValidMonths), 0).Before(time.Now()) {
		return cmd.RenderFailedTemplate("Voucher is expired")
	}

	if voucher.IsClaimed() {
		return cmd.RenderFailedTemplate("Voucher code claimed before")
	}

	address := args["address"]
	valInfo, _ := v.clientManager.GetValidatorInfo(address)
	if valInfo != nil {
		err = errors.New("This address is already a staked validator")
		log.Warn(fmt.Sprintf("Staked validator found. %s", address))

		return cmd.RenderErrorTemplate(err)
	}

	pubKey, err := v.clientManager.FindPublicKey(address, false)
	if err != nil {
		log.Warn(fmt.Sprintf("Peer not found. %s", address))

		return cmd.RenderErrorTemplate(err)
	}

	memo := fmt.Sprintf("Voucher %s claimed by Pagu", code)
	txHash, err := v.wallet.BondTransaction(pubKey, address, memo, voucher.Amount)
	if err != nil {
		return cmd.RenderErrorTemplate(err)
	}

	if txHash == "" {
		return cmd.RenderFailedTemplate("Can't send bond transaction")
	}

	if err = v.db.ClaimVoucher(voucher.ID, txHash, caller.ID); err != nil {
		return cmd.RenderErrorTemplate(err)
	}

	return cmd.RenderResultTemplate("txHash", txHash)
}
