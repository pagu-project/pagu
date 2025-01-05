package crowdfund

import (
	"github.com/pagu-project/pagu/internal/engine/command"
	"github.com/pagu-project/pagu/internal/entity"
)

// Add caller.Name here?
const infoResponseTemplate = `
**{{.campaign.Title}}**

{{.campaign.Desc}}

Packages:
{{range .campaign.Packages}}
- {{.Name}}{{end}}
`

func (c *CrowdfundCmd) infoHandler(
	_ *entity.User,
	cmd *command.Command,
	_ map[string]string,
) command.CommandResult {
	campaign, err := c.db.GetCrowdfundCampaign(c.config.ActiveCampaignID)
	if err != nil {
		return cmd.RenderErrorTemplate(err)
	}

	return cmd.RenderResultTemplate(infoResponseTemplate, "campaign", campaign)
}
