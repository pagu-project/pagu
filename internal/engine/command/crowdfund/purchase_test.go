package crowdfund

import (
	"testing"

	"github.com/pagu-project/pagu/internal/entity"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestPurchase(t *testing.T) {
	td := setup(t)

	caller := &entity.User{DBModel: entity.DBModel{ID: 1}}

	t.Run("No active Campaign", func(t *testing.T) {
		args := map[string]string{
			"package": "1",
		}
		result := td.crowdfundCmd.purchaseHandler(caller, td.crowdfundCmd.subCmdPurchase, args)
		assert.False(t, result.Successful)
		assert.Contains(t, result.Message, "No active campaign")
	})

	t.Run("Ok", func(t *testing.T) {
		_ = td.createTestCampaign(t)
		td.nowpayments.EXPECT().CreateInvoice(gomock.Any(), gomock.Any()).Return(
			"invoice-id", nil,
		)

		td.nowpayments.EXPECT().PaymentLink("invoice-id").Return(
			"payment-link",
		)

		args := map[string]string{
			"package": "1",
		}
		result := td.crowdfundCmd.purchaseHandler(caller, td.crowdfundCmd.subCmdPurchase, args)

		assert.True(t, result.Successful)
		assert.Contains(t, result.Message, "payment-link")
	})

	t.Run("Invalid Package Number", func(t *testing.T) {
		args := map[string]string{
			"package": "0",
		}
		result := td.crowdfundCmd.purchaseHandler(caller, td.crowdfundCmd.subCmdPurchase, args)

		assert.False(t, result.Successful)
		assert.Contains(t, result.Message, "Invalid package number: 0")
	})
}
