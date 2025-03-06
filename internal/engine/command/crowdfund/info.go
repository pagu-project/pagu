package crowdfund

import (
	"github.com/pagu-project/pagu/internal/engine/command"
	"github.com/pagu-project/pagu/internal/entity"
)

func (c *CrowdfundCmd) infoHandler(
	_ *entity.User,
	cmd *command.Command,
	_ map[string]string,
) command.CommandResult {
	activeCampaign := c.activeCampaign()
	if activeCampaign == nil {
		return cmd.RenderFailedTemplate("There is no active campaign")
	}

	return cmd.RenderResultTemplate("campaign", activeCampaign)
}
