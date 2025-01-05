package crowdfund

import (
	"testing"

	"github.com/pagu-project/pagu/internal/engine/command"
	"github.com/pagu-project/pagu/internal/entity"
	"github.com/stretchr/testify/assert"
)

func TestInf(t *testing.T) {
	td := setup(t)

	campaign := td.createTestCampaign(t)

	caller := &entity.User{DBModel: entity.DBModel{ID: 1}}
	cmd := &command.Command{}

	t.Run("No active Campaign", func(t *testing.T) {
		result := td.crowdfundCmd.infoHandler(caller, cmd, nil)
		assert.False(t, result.Successful)
		assert.Contains(t, result.Message, "Read error: record not found")
	})

	t.Run("No active Campaign", func(t *testing.T) {
		td.crowdfundCmd.config.ActiveCampaignID = campaign.ID

		result := td.crowdfundCmd.infoHandler(caller, cmd, nil)
		assert.True(t, result.Successful)
		assert.Contains(t, result.Message, campaign.Title)
		assert.Contains(t, result.Message, campaign.Desc)
	})
}
