package voucher

import (
	"fmt"

	"github.com/pagu-project/pagu/internal/engine/command"
	"github.com/pagu-project/pagu/internal/entity"
)

func (v *VoucherCmd) statusHandler(_ *entity.User, cmd *command.Command, args map[string]string) command.CommandResult {
	if args[argNameStatusCode] != "" {
		voucher, err := v.db.GetVoucherByCode(args[argNameStatusCode])
		if err != nil {
			return cmd.RenderErrorTemplate(err)
		}

		return v.statusVoucher(cmd, voucher)
	}

	if args[argNameStatusEmail] != "" {
		voucher, err := v.db.GetVoucherByEmail(args[argNameStatusEmail])
		if err != nil {
			return cmd.RenderErrorTemplate(err)
		}

		return v.statusVoucher(cmd, voucher)
	}

	return cmd.RenderFailedTemplate("set email or code")
}

func (*VoucherCmd) statusVoucher(cmd *command.Command, voucher *entity.Voucher) command.CommandResult {
	isClaimed := "NO"
	txLink := ""
	if voucher.IsClaimed() {
		isClaimed = "YES"
		txLink = fmt.Sprintf("https://pacviewer.com/transaction/%s", voucher.TxHash)
	}

	voucherExpiryDate := voucher.CreatedAt.AddDate(0, int(voucher.ValidMonths), 0).Format("02/01/2006, 15:04:05")

	return cmd.RenderResultTemplate(
		"voucher", voucher,
		"expireAt", voucherExpiryDate,
		"isClaimed", isClaimed,
		"txLink", txLink)
}
