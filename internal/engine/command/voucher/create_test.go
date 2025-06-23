package voucher

import (
	"testing"
	"time"

	"github.com/pagu-project/pagu/internal/entity"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestCreate(t *testing.T) {
	td := setup(t)

	caller := &entity.User{DBModel: entity.DBModel{ID: 1}}

	t.Run("more than 1000 PAC", func(t *testing.T) {
		args := map[string]string{
			"type":         "1",
			"amount":       "1001",
			"valid-months": "1",
			"template":     "sample",
		}

		result := td.voucherCmd.createHandler(caller, td.voucherCmd.subCmdCreate, args)
		assert.False(t, result.Successful)
		assert.Contains(t, result.Message, "amount is more than 1000")
	})

	t.Run("invalid amount", func(t *testing.T) {
		args := map[string]string{
			"type":         "1",
			"email":        td.RandEmail(),
			"amount":       "invalid-amount",
			"valid-months": "1",
			"template":     "sample",
		}

		result := td.voucherCmd.createHandler(caller, td.voucherCmd.subCmdCreate, args)
		assert.False(t, result.Successful)
		assert.Contains(t, result.Message, "invalid amount")
	})

	t.Run("invalid email", func(t *testing.T) {
		args := map[string]string{
			"type":         "1",
			"email":        "invalid-email",
			"amount":       "100",
			"valid-months": "1",
			"template":     "sample",
		}

		result := td.voucherCmd.createHandler(caller, td.voucherCmd.subCmdCreate, args)
		assert.False(t, result.Successful)
		assert.Contains(t, result.Message, "invalid email address: invalid-email")
	})

	t.Run("wrong month", func(t *testing.T) {
		args := map[string]string{
			"type":         "1",
			"email":        td.RandEmail(),
			"amount":       "100",
			"valid-months": "1.1",
			"template":     "sample",
		}

		result := td.voucherCmd.createHandler(caller, td.voucherCmd.subCmdCreate, args)
		assert.False(t, result.Successful)
	})

	t.Run("normal", func(t *testing.T) {
		args := map[string]string{
			"type":         "1",
			"email":        td.RandEmail(),
			"recipient":    "Kayhan",
			"amount":       "100",
			"valid-months": "1",
			"template":     "sample",
			"description":  "Some descriptions",
		}

		td.mockMailer.EXPECT().SendTemplateMailAsync(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)

		result := td.voucherCmd.createHandler(caller, td.voucherCmd.subCmdCreate, args)
		assert.True(t, result.Successful)
		assert.Contains(t, result.Message, "Voucher created successfully!")
	})
}

func TestTestCreateWithExistingVoucher(t *testing.T) {
	td := setup(t)
	caller := &entity.User{DBModel: entity.DBModel{ID: 1}}

	t.Run("expired voucher", func(t *testing.T) {
		createdAt := time.Now().AddDate(0, -2, 0) // 2 months ago
		voucher := td.createTestVoucher(t,
			WithValidMonths(1), // 1 month validity
			WithCreatedAt(createdAt),
		)

		args := map[string]string{
			"type":         "1",
			"email":        voucher.Email,
			"recipient":    voucher.Recipient,
			"amount":       "100",
			"valid-months": "1",
			"template":     "sample",
			"description":  "Some descriptions",
		}

		td.mockMailer.EXPECT().SendTemplateMailAsync(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)

		result := td.voucherCmd.createHandler(caller, td.voucherCmd.subCmdCreate, args)
		assert.True(t, result.Successful)
		assert.Contains(t, result.Message, "Voucher created successfully!")
	})

	t.Run("non-expired voucher", func(t *testing.T) {
		createdAt := time.Now().AddDate(0, -1, 0) // 1 month ago
		voucher := td.createTestVoucher(t,
			WithValidMonths(2), // 2 months validity
			WithCreatedAt(createdAt),
			WithTxHash("tx-hash"),
		)
		args := map[string]string{
			"type":         "1",
			"email":        voucher.Email,
			"recipient":    voucher.Recipient,
			"amount":       "100",
			"valid-months": "1",
			"template":     "sample",
			"description":  "Some descriptions",
		}

		result := td.voucherCmd.createHandler(caller, td.voucherCmd.subCmdCreate, args)
		assert.False(t, result.Successful)
		assert.Contains(t, result.Message, "email already has a non-expired voucher")
	})

	t.Run("claimed voucher", func(t *testing.T) {
		createdAt := time.Now().AddDate(0, -1, 0) // 1 month ago
		voucher := td.createTestVoucher(t,
			WithValidMonths(2), // 2 months validity
			WithCreatedAt(createdAt),
		)
		args := map[string]string{
			"type":         "1",
			"email":        voucher.Email,
			"recipient":    voucher.Recipient,
			"amount":       "100",
			"valid-months": "1",
			"template":     "sample",
			"description":  "Some descriptions",
		}

		// Resend the email
		td.mockMailer.EXPECT().SendTemplateMailAsync(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)

		result := td.voucherCmd.createHandler(caller, td.voucherCmd.subCmdCreate, args)
		assert.False(t, result.Successful)
		assert.Contains(t, result.Message, "email already has a non-expired voucher")
	})
}

func TestCreateType(t *testing.T) {
	td := setup(t)

	caller := &entity.User{DBModel: entity.DBModel{ID: 1}}

	t.Run("stake voucher", func(t *testing.T) {
		email := td.RandEmail()
		args := map[string]string{
			"type":         "0",
			"email":        email,
			"recipient":    "Kayhan",
			"amount":       "100",
			"valid-months": "1",
			"template":     "sample",
			"description":  "Some descriptions",
		}

		td.mockMailer.EXPECT().SendTemplateMailAsync(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)

		result := td.voucherCmd.createHandler(caller, td.voucherCmd.subCmdCreate, args)
		assert.True(t, result.Successful)
		assert.Contains(t, result.Message, "Voucher created successfully!")

		voucher, _ := td.testDB.GetVoucherByEmail(email)
		assert.Equal(t, entity.VoucherTypeStake, voucher.Type)
	})

	t.Run("liquid voucher", func(t *testing.T) {
		email := td.RandEmail()
		args := map[string]string{
			"type":         "1",
			"email":        email,
			"recipient":    "Kayhan",
			"amount":       "100",
			"valid-months": "1",
			"template":     "sample",
			"description":  "Some descriptions",
		}

		td.mockMailer.EXPECT().SendTemplateMailAsync(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)

		result := td.voucherCmd.createHandler(caller, td.voucherCmd.subCmdCreate, args)
		assert.True(t, result.Successful)
		assert.Contains(t, result.Message, "Voucher created successfully!")

		voucher, _ := td.testDB.GetVoucherByEmail(email)
		assert.Equal(t, entity.VoucherTypeLiquid, voucher.Type)
	})
}
