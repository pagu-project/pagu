package crowdfund

import (
	"testing"

	"github.com/pagu-project/pagu/internal/entity"
	"github.com/stretchr/testify/assert"
)

func TestInfo(t *testing.T) {
	td := setup(t)

	caller := &entity.User{DBModel: entity.DBModel{ID: 1}}

	t.Run("No active Campaign", func(t *testing.T) {
		result := td.crowdfundCmd.infoHandler(caller, td.crowdfundCmd.subCmdInfo, nil)
		assert.False(t, result.Successful)
		assert.Contains(t, result.Message, "No active campaign")
	})

	t.Run("OK", func(t *testing.T) {
		testCampaign := td.createTestCampaign(t)

		result := td.crowdfundCmd.infoHandler(caller, td.crowdfundCmd.subCmdInfo, nil)
		assert.True(t, result.Successful)
		assert.Contains(t, result.Message, testCampaign.Title)
		assert.Contains(t, result.Message, testCampaign.Desc)
	})
}
