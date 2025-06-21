package voucher

import (
	"github.com/pagu-project/pagu/internal/engine/command"
	"github.com/pagu-project/pagu/internal/entity"
	"github.com/pagu-project/pagu/pkg/amount"
)

func (c *VoucherCmd) reportHandler(_ *entity.User, cmd *command.Command, _ map[string]string) command.CommandResult {
	vouchers, err := c.db.ListVoucher()
	if err != nil {
		return cmd.RenderErrorTemplate(err)
	}

	total := 0
	totalClaimedAmount := amount.Amount(0)
	totalClaimed := 0
	totalExpired := 0

	for _, voucher := range vouchers {
		total++

		if voucher.IsClaimed() {
			totalClaimed++
			totalClaimedAmount += voucher.Amount
		}
		if voucher.IsExpired() {
			totalExpired++
		}
	}

	return cmd.RenderResultTemplate(
		"total", total,
		"totalClaimed", totalClaimed,
		"totalClaimedAmount", totalClaimedAmount,
		"totalExpired", totalExpired)
}
