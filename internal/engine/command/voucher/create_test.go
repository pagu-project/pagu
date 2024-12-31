package voucher

import (
	"testing"

	"github.com/h2non/gock"
	"github.com/pagu-project/pagu/internal/engine/command"
	"github.com/pagu-project/pagu/internal/entity"
	"github.com/stretchr/testify/assert"
)

func TestCreateOne(t *testing.T) {
	td := setup(t)

	cmd := &command.Command{}
	caller := &entity.User{DBModel: entity.DBModel{ID: 1}}

	t.Run("more than 1000 PAC", func(t *testing.T) {
		args := make(map[string]string)
		args["amount"] = "1001"
		args["valid-months"] = "1"
		result := td.voucherCmd.createOneHandler(caller, cmd, args)
		assert.False(t, result.Successful)
		assert.Contains(t, result.Message, "stake amount is more than 1000")
	})

	t.Run("wrong month", func(t *testing.T) {
		args := make(map[string]string)
		args["amount"] = "100"
		args["valid-months"] = "1.1"
		result := td.voucherCmd.createOneHandler(caller, cmd, args)
		assert.False(t, result.Successful)
	})

	t.Run("normal", func(t *testing.T) {
		args := make(map[string]string)
		args["amount"] = "100"
		args["valid-months"] = "1"
		result := td.voucherCmd.createOneHandler(caller, cmd, args)
		assert.True(t, result.Successful)
		assert.Contains(t, result.Message, "Voucher created successfully!")
	})

	t.Run("normal with optional arguments", func(t *testing.T) {
		args := make(map[string]string)
		args["amount"] = "100"
		args["valid-months"] = "12"
		args["recipient"] = "Kayhan"
		args["description"] = "Testnet node"
		result := td.voucherCmd.createOneHandler(caller, cmd, args)
		assert.True(t, result.Successful)
		assert.Contains(t, result.Message, "Voucher created successfully!")
	})
}

func TestCreateBulk(t *testing.T) {
	td := setup(t)

	cmd := &command.Command{}
	caller := &entity.User{DBModel: entity.DBModel{ID: 1}}

	t.Run("normal", func(t *testing.T) {
		defer gock.Off()
		gock.New("http://foo.com").
			Get("/bar").
			Reply(200).
			BodyString("Recipient,Email,Amount,Validated,Description\n" +
				"foo.bar,a@gmail.com,1,2,Some Descriptions\n" +
				"foo.bar,b@gmail.com,1,2,Some Descriptions")

		args := make(map[string]string)
		args["file"] = "http://foo.com/bar"
		args["notify"] = "TRUE"

		result := td.voucherCmd.createBulkHandler(caller, cmd, args)

		assert.True(t, result.Successful)
		assert.Contains(t, result.Message, "Vouchers created successfully!")
	})
}
