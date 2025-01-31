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

	assert.True(t, result.Successful)
	assert.Contains(t, result.Message, "tpc1r3666ga8venyjykp7ffyy96mccs82uh9ry3d2fk")
}
