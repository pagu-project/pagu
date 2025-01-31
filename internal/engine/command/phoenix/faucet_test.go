package phoenix

import (
	"testing"
	"time"

	"github.com/pagu-project/pagu/internal/entity"
	"github.com/pagu-project/pagu/pkg/utils"
	"github.com/stretchr/testify/assert"
)

func TestFaucet(t *testing.T) {
	td := setup(t)

	caller := &entity.User{DBModel: entity.DBModel{ID: 1}}
	randAddr := utils.TestnetAddressToString(td.RandAccAddress())

	t.Run("No Address", func(t *testing.T) {
		args := map[string]string{}
		result := td.phoenixCmd.faucetHandler(caller, td.phoenixCmd.subCmdFaucet, args)
		assert.False(t, result.Successful)
		assert.Contains(t, result.Message, "Please provide a valid address")
	})

	t.Run("Invalid Address", func(t *testing.T) {
		args := map[string]string{
			"address": "invalid-address",
		}
		result := td.phoenixCmd.faucetHandler(caller, td.phoenixCmd.subCmdFaucet, args)
		assert.False(t, result.Successful)
		assert.Contains(t, result.Message, "invalid separator index")
	})

	t.Run("Ok", func(t *testing.T) {
		args := map[string]string{
			"address": randAddr,
		}
		result := td.phoenixCmd.faucetHandler(caller, td.phoenixCmd.subCmdFaucet, args)
		if result.Successful {
			assert.Contains(t, result.Message, td.phoenixCmd.faucetAmount.String())

			time.Sleep(1 * time.Second)

			// Test cooldown period
			result2 := td.phoenixCmd.faucetHandler(caller, td.phoenixCmd.subCmdFaucet, args)
			assert.False(t, result2.Successful)
			assert.Contains(t, result2.Message, "Please try again in 59 minutes")
		} else {
			// In case the test wallet is empty
			assert.Contains(t, result.Message, "insufficient funds")
		}
	})
}
