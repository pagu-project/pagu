package voucher

import (
	"errors"
	"fmt"
	"time"

	"github.com/pagu-project/pagu/internal/engine/command"
	"github.com/pagu-project/pagu/internal/entity"
	"github.com/pagu-project/pagu/pkg/amount"
)

func (v *VoucherCmd) statusHandler(_ *entity.User, cmd *command.Command, args map[string]string) command.CommandResult {
	if args["code"] != "" {
		code, ok := args["code"]
		if !ok {
			return cmd.ErrorResult(errors.New("invalid code param"))
		}

		return v.statusVoucher(cmd, code)
	}

	return v.statusAllVouchers(cmd)
}

func (v *VoucherCmd) statusVoucher(cmd *command.Command, code string) command.CommandResult {
	voucher, err := v.db.GetVoucherByCode(code)
	if err != nil {
		return cmd.ErrorResult(errors.New("voucher code is not valid, no voucher found"))
	}

	isClaimed := "NO"
	txLink := ""
	if voucher.IsClaimed() {
		isClaimed = "YES"
		txLink = fmt.Sprintf("https://pacviewer.com/transaction/%s", voucher.TxHash)
	}

	return cmd.SuccessfulResultF("Code: %s\nAmount: %s\n"+
		"Expire At: %s\nRecipient: %s\nDescription: %s\nClaimed: %v\nTx Link: %s"+
		"\n",
		voucher.Code,
		voucher.Amount,
		voucher.CreatedAt.AddDate(0, int(voucher.ValidMonths), 0).Format("02/01/2006, 15:04:05"),
		voucher.Recipient,
		voucher.Desc,
		isClaimed,
		txLink)
}

func (v *VoucherCmd) statusAllVouchers(cmd *command.Command) command.CommandResult {
	vouchers, err := v.db.ListVoucher()
	if err != nil {
		return cmd.ErrorResult(err)
	}

	total := 0
	totalAmount := amount.Amount(0)
	totalClaimedAmount := amount.Amount(0)
	totalClaimed := 0
	totalExpired := 0

	for _, vch := range vouchers {
		total++
		totalAmount += vch.Amount

		if vch.IsClaimed() {
			totalClaimed++
			totalClaimedAmount += vch.Amount
		}
		if time.Until(vch.CreatedAt.AddDate(0, int(vch.ValidMonths), 0)) <= 0 {
			totalExpired++
		}
	}

	return cmd.SuccessfulResultF("Total Vouchers: %d\nTotal Amount: %s\n\n\n"+
		"Claimed: %d\nTotal Claimed Amount: %s\nTotal Expired: %d"+
		"\n",
		total,
		totalAmount,
		totalClaimed,
		totalClaimedAmount,
		totalExpired)
}
