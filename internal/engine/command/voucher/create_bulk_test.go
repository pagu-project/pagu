package voucher

import (
	"testing"

	"github.com/pagu-project/pagu/internal/entity"
	"github.com/stretchr/testify/assert"
)

func TestCreateBulk(t *testing.T) {
	td := setup(t)

	caller := &entity.User{DBModel: entity.DBModel{ID: 1}}

	t.Run("normal", func(t *testing.T) {
		args := map[string]string{
			"template": "sample",
			"type":     "1",
			"csv": `recipient,email,amount,valid-months,desc
name1,test1@test.com,100,1,Some descriptions
name2,test2@test.com,100,1,Some descriptions
name3,test3@test.com,100,1,Some descriptions
`,
		}

		result := td.voucherCmd.createBulkHandler(caller, td.voucherCmd.subCmdCreateBulk, args)
		assert.True(t, result.Successful)
		assert.Contains(t, result.Message, "Vouchers are going to send to recipients")
	})
}
