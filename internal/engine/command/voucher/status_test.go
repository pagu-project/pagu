package voucher

import (
	"testing"
	"time"

	"github.com/pagu-project/pagu/internal/engine/command"
	"github.com/pagu-project/pagu/internal/entity"
	"github.com/stretchr/testify/assert"
)

func TestStatusNormal(t *testing.T) {
	td := setup(t)

	cmd := &command.Command{}
	caller := &entity.User{DBModel: entity.DBModel{ID: 1}}

	voucherCode := "12345678"
	testVoucher := td.createTestVoucher(t, WithCode(voucherCode), WithAmount(100e9))

	t.Run("one code status normal", func(t *testing.T) {
		args := make(map[string]string)
		args["code"] = voucherCode
		result := td.voucherCmd.statusHandler(caller, cmd, args)
		assert.True(t, result.Successful)
		assert.Contains(t, result.Message, "Code: 12345678")
		assert.Contains(t, result.Message, testVoucher.Recipient)
	})

	t.Run("wrong code", func(t *testing.T) {
		args := make(map[string]string)
		args["code"] = "000"
		result := td.voucherCmd.statusHandler(caller, cmd, args)
		assert.False(t, result.Successful)
		assert.Equal(t, result.Message, "An error occurred: voucher code is not valid, no voucher found")
	})

	t.Run("list vouchers status normal", func(t *testing.T) {
		td.createTestVoucher(t, WithAmount(50e9), WithTxHash("claimed_tx_hash"))
		td.createTestVoucher(t, WithAmount(20e9), WithValidMonths(1),
			WithCreatedAt(time.Now().AddDate(0, -2, 0)))

		args := make(map[string]string)
		result := td.voucherCmd.statusHandler(caller, cmd, args)
		assert.True(t, result.Successful)
		assert.Equal(t, result.Message, "Total Vouchers: 3\nTotal Amount: 170 PAC\n\n\n"+
			"Claimed: 1\nTotal Claimed Amount: 50 PAC\nTotal Expired: 1"+
			"\n")
	})
}
