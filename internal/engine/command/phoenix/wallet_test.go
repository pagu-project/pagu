package phoenix

import (
	"testing"

	"github.com/pagu-project/pagu/internal/entity"
	"github.com/pagu-project/pagu/pkg/utils"
	"github.com/stretchr/testify/assert"
)

func TestWallet(t *testing.T) {
	td := setup(t)

	caller := &entity.User{DBModel: entity.DBModel{ID: 1}}

	result := td.phoenixCmd.walletHandler(caller, td.phoenixCmd.subCmdWallet, nil)

	assert.True(t, result.Successful)
	assert.Contains(t, result.Message, utils.TestnetAddressToString(td.phoenixCmd.faucetAddress))
}
