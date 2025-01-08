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
	if c.activeCampaign == nil {
		return cmd.RenderFailedTemplate("No active campaign")
	}

	return cmd.RenderResultTemplate("campaign", c.activeCampaign)
}
