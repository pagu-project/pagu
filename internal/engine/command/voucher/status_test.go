package voucher

import (
	"testing"

	"github.com/pagu-project/pagu/internal/entity"
	"github.com/stretchr/testify/assert"
)

func TestStatusNormal(t *testing.T) {
	td := setup(t)

	caller := &entity.User{DBModel: entity.DBModel{ID: 1}}

	voucherCode := "12345678"
	testVoucher := td.createTestVoucher(t, WithCode(voucherCode), WithAmount(100e9))

	t.Run("one code status normal", func(t *testing.T) {
		args := map[string]string{
			"code": voucherCode,
		}
		result := td.voucherCmd.statusHandler(caller, td.voucherCmd.subCmdStatus, args)
		assert.True(t, result.Successful)
		assert.Contains(t, result.Message, "12345678")
		assert.Contains(t, result.Message, testVoucher.Recipient)
	})

	t.Run("wrong code", func(t *testing.T) {
		args := map[string]string{
			"code": "000",
		}
		result := td.voucherCmd.statusHandler(caller, td.voucherCmd.subCmdStatus, args)
		assert.False(t, result.Successful)
		assert.Contains(t, result.Message, "record not found")
	})
}
