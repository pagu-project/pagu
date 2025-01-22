package calculator

import (
	"errors"
	"fmt"
	"testing"

	"github.com/pagu-project/pagu/pkg/amount"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestFeeHandler(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	td := setup(t)
	cmd := td.calculatorCmd.subCmdFee

	t.Run("Invalid Amount Param", func(t *testing.T) {
		args := map[string]string{
			"amount": "invalid",
		}

		result := td.calculatorCmd.feeHandler(nil, cmd, args)
		assert.False(t, result.Successful)
		assert.Contains(t, result.Message, "Invalid amount param")
	})

	t.Run("Error from GetFee", func(t *testing.T) {
		args := map[string]string{
			"amount": "10",
		}

		amt, _ := amount.FromString("10")
		td.mockClientMgr.EXPECT().GetFee(amt.ToNanoPAC()).Return(int64(0), errors.New("some error"))

		result := td.calculatorCmd.feeHandler(nil, cmd, args)
		assert.False(t, result.Successful)
		assert.Contains(t, result.Message, "some error")
	})

	t.Run("Successful Fee Calculation", func(t *testing.T) {
		args := map[string]string{
			"amount": "10",
		}

		amt, _ := amount.FromString("10")
		td.mockClientMgr.EXPECT().GetFee(amt.ToNanoPAC()).Return(int64(100), nil)

		result := td.calculatorCmd.feeHandler(nil, cmd, args)
		assert.True(t, result.Successful)
		assert.Contains(t, result.Message, fmt.Sprintf("Sending %v will cost %v with current fee percentage", amt, amount.Amount(int64(100))))
	})
}
