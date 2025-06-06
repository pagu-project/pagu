package voucher

import (
	"testing"

	"github.com/pagu-project/pagu/internal/entity"
	"github.com/stretchr/testify/assert"
)

func TestCreateOne(t *testing.T) {
	td := setup(t)

	caller := &entity.User{DBModel: entity.DBModel{ID: 1}}

	t.Run("more than 1000 PAC", func(t *testing.T) {
		args := map[string]string{
			"amount":       "1001",
			"valid-months": "1",
		}

		result := td.voucherCmd.createHandler(caller, td.voucherCmd.subCmdCreate, args)
		assert.False(t, result.Successful)
		assert.Contains(t, result.Message, "Stake amount is more than 1000")
	})

	t.Run("wrong month", func(t *testing.T) {
		args := map[string]string{
			"amount":       "100",
			"valid-months": "1.1",
		}

		result := td.voucherCmd.createHandler(caller, td.voucherCmd.subCmdCreate, args)
		assert.False(t, result.Successful)
	})

	t.Run("normal", func(t *testing.T) {
		args := map[string]string{
			"amount":       "100",
			"valid-months": "1",
		}

		result := td.voucherCmd.createHandler(caller, td.voucherCmd.subCmdCreate, args)
		assert.True(t, result.Successful)
		assert.Contains(t, result.Message, "Voucher created successfully!")
	})

	t.Run("normal with optional arguments", func(t *testing.T) {
		args := map[string]string{
			"amount":       "100",
			"valid-months": "12",
			"recipient":    "Kayhan",
			"description":  "Some descriptions",
		}

		result := td.voucherCmd.createHandler(caller, td.voucherCmd.subCmdCreate, args)
		assert.True(t, result.Successful)
		assert.Contains(t, result.Message, "Voucher created successfully!")
	})
}
