package phoenix

import (
	"testing"

	"github.com/pagu-project/pagu/internal/entity"
	"github.com/stretchr/testify/assert"
)

func TestWallet(t *testing.T) {
	td := setup(t)

	caller := &entity.User{DBModel: entity.DBModel{ID: 1}}

	result := td.phoenixCmd.walletHandler(caller, td.phoenixCmd.subCmdWallet, nil)

	// Verify that calling the `Balance` API works fine.
	// For security reasons, we cannot use any existing private key in this test.
	assert.False(t, result.Successful)
	assert.Contains(t, result.Message, "account not found")
}
