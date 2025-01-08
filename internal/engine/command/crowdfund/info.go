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
	if c.activeCampaign == nil {
		return cmd.RenderFailedTemplate("No active campaign")
	}

	return cmd.RenderResultTemplate("campaign", c.activeCampaign)
}
