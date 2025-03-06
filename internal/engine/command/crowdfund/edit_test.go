package crowdfund

import (
	"testing"

	"github.com/pagu-project/pagu/internal/entity"
	"github.com/stretchr/testify/assert"
)

func TestEdit(t *testing.T) {
	td := setup(t)

	caller := &entity.User{DBModel: entity.DBModel{ID: 1}}

	t.Run("No active Campaign", func(t *testing.T) {
		args := map[string]string{
			"package": "1",
		}
		result := td.crowdfundCmd.editHandler(caller, td.crowdfundCmd.subCmdEdit, args)
		assert.False(t, result.Successful)
		assert.Contains(t, result.Message, "There is no campaign")
	})

	_ = td.createTestCampaign(t)

	t.Run("Ok", func(t *testing.T) {
		args := map[string]string{
			"title": "crowdfund-title-edited",
			"desc":  "crowdfund-desc-edited",
			"packages": `
			[
			   {"name": "package-1-edited", "usd_amount": 100, "pac_amount": 100},
			   {"name": "package-2-edited", "usd_amount": 200, "pac_amount": 200},
			   {"name": "package-3-edited", "usd_amount": 300, "pac_amount": 300}
			]`,
			"disable": "false",
		}
		result := td.crowdfundCmd.editHandler(caller, td.crowdfundCmd.subCmdEdit, args)
		assert.True(t, result.Successful)

		resultInfo := td.crowdfundCmd.editHandler(caller, td.crowdfundCmd.subCmdInfo, nil)
		assert.Contains(t, resultInfo.Message, "edited")
	})

	t.Run("Disable Campaign", func(t *testing.T) {
		args := map[string]string{
			"disable": "true",
		}
		result := td.crowdfundCmd.editHandler(caller, td.crowdfundCmd.subCmdEdit, args)
		assert.True(t, result.Successful)

		resultInfo := td.crowdfundCmd.infoHandler(caller, td.crowdfundCmd.subCmdInfo, nil)
		assert.Contains(t, resultInfo.Message, "There is no active campaign")
	})

	t.Run("Enable Campaign", func(t *testing.T) {
		args := map[string]string{
			"disable": "false",
		}
		result := td.crowdfundCmd.editHandler(caller, td.crowdfundCmd.subCmdEdit, args)
		assert.True(t, result.Successful)

		resultInfo := td.crowdfundCmd.infoHandler(caller, td.crowdfundCmd.subCmdInfo, nil)
		assert.Contains(t, resultInfo.Message, "crowdfund-title-edited")
	})
}
