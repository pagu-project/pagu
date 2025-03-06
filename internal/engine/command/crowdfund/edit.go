package crowdfund

import (
	"encoding/json"
	"strings"

	"github.com/pagu-project/pagu/internal/engine/command"
	"github.com/pagu-project/pagu/internal/entity"
	"github.com/pagu-project/pagu/pkg/utils"
)

func (c *CrowdfundCmd) editHandler(
	_ *entity.User,
	cmd *command.Command,
	args map[string]string,
) command.CommandResult {
	campaign := c.lastCampaign()
	if campaign == nil {
		return cmd.RenderFailedTemplateF("There is no campaign")
	}

	title := args[argNameEditTitle]
	desc := args[argNameEditDesc]
	packagesJSON := args[argNameEditPackages]
	disableStr := args[argNameEditDisable]

	if title != "" {
		campaign.Title = title
	}

	if desc != "" {
		campaign.Desc = strings.ReplaceAll(desc, `\n`, "\n")
	}

	if packagesJSON != "" {
		packages := []entity.Package{}
		err := json.Unmarshal([]byte(packagesJSON), &packages)
		if err != nil {
			return cmd.RenderErrorTemplate(err)
		}
		campaign.Packages = packages
	}

	campaign.Active = !utils.IsToggleEnabled(disableStr)

	err := c.db.UpdateCrowdfundCampaign(campaign)
	if err != nil {
		return cmd.RenderErrorTemplate(err)
	}

	return cmd.RenderResultTemplate("campaign", campaign)
}
