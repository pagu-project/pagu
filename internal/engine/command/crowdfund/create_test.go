package crowdfund

import (
	"testing"

	"github.com/pagu-project/pagu/internal/engine/command"
	"github.com/pagu-project/pagu/internal/entity"
	"github.com/stretchr/testify/assert"
)

func TestCreate(t *testing.T) {
	td := setup(t)

	caller := &entity.User{DBModel: entity.DBModel{ID: 1}}
	cmd := &command.Command{}

	t.Run("Invalid Packages", func(t *testing.T) {
		args := map[string]string{
			"title":    "crowdfund-title",
			"desc":     "crowdfund-desc",
			"packages": "INVALID-JSON",
		}
		result := td.crowdfundCmd.createHandler(caller, cmd, args)
		assert.False(t, result.Successful)
		assert.Equal(t, result.Message, "invalid character 'I' looking for beginning of value")
	})

	t.Run("Empty title", func(t *testing.T) {
		args := map[string]string{
			"title":    "",
			"desc":     "",
			"packages": "[]",
		}
		result := td.crowdfundCmd.createHandler(caller, cmd, args)
		assert.False(t, result.Successful)
		assert.Equal(t, result.Message, "The title of the crowdfunding campaign cannot be empty")
	})

	t.Run("Empty Packages", func(t *testing.T) {
		args := map[string]string{
			"title":    "crowdfund-title",
			"desc":     "crowdfund-desc",
			"packages": "[]",
		}
		result := td.crowdfundCmd.createHandler(caller, cmd, args)
		assert.False(t, result.Successful)
		assert.Equal(t, result.Message, "At least 3 packages are required for the crowdfunding campaign")
	})

	t.Run("Ok", func(t *testing.T) {
		args := map[string]string{
			"title": "crowdfund-title",
			"desc":  "crowdfund-desc",
			"packages": `
			[
			   {"name": "package-1", "usd_amount": 100, "pac_amount": 100},
			   {"name": "package-2", "usd_amount": 200, "pac_amount": 200},
			   {"name": "package-3", "usd_amount": 300, "pac_amount": 300}
			]`,
		}
		result := td.crowdfundCmd.createHandler(caller, cmd, args)
		assert.True(t, result.Successful)
		assert.Equal(t, result.Message, "Crowdfund campaign 'crowdfund-title' created successfully with 3 packages")
	})
}
