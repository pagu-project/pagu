package calculator

import (
	"testing"

	pactus "github.com/pactus-project/pactus/www/grpc/gen/go"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestRewardHandler(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	td := setup(t)
	cmd := td.calculatorCmd.subCmdReward

	t.Run("Invalid Stake Param", func(t *testing.T) {
		args := map[string]string{
			"stake": "invalid",
			"days":  "10",
		}

		result := td.calculatorCmd.rewardHandler(nil, cmd, args)
		assert.False(t, result.Successful)
		assert.Contains(t, result.Message, "Invalid stake param")
	})

	t.Run("Stake Out of Range", func(t *testing.T) {
		args := map[string]string{
			"stake": "0.5",
			"days":  "10",
		}

		result := td.calculatorCmd.rewardHandler(nil, cmd, args)
		assert.False(t, result.Successful)
		assert.Contains(t, result.Message, "is invalid amount, minimum stake amount is 1 PAC and maximum is 1,000 PAC")
	})

	t.Run("Invalid Days Param", func(t *testing.T) {
		args := map[string]string{
			"stake": "10",
			"days":  "invalid",
		}

		result := td.calculatorCmd.rewardHandler(nil, cmd, args)
		assert.False(t, result.Successful)
		assert.Contains(t, result.Message, "invalid days param")
	})

	t.Run("Days Out of Range", func(t *testing.T) {
		args := map[string]string{
			"stake": "10",
			"days":  "366",
		}

		result := td.calculatorCmd.rewardHandler(nil, cmd, args)
		assert.False(t, result.Successful)
		assert.Contains(t, result.Message, "is invalid time, minimum time value is 1 and maximum is 365")
	})

	t.Run("Successful Reward Calculation", func(t *testing.T) {
		args := map[string]string{
			"stake": "10",
			"days":  "10",
		}

		td.mockClientMgr.EXPECT().GetBlockchainInfo().Return(
			&pactus.GetBlockchainInfoResponse{
				TotalPower: 1000,
			}, nil,
		).AnyTimes()

		result := td.calculatorCmd.rewardHandler(nil, cmd, args)
		assert.True(t, result.Successful)
		assert.Contains(t, result.Message, "stake")
		assert.Contains(t, result.Message, "reward")
	})
}
