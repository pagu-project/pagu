package crowdfund

import (
	"encoding/json"

	"github.com/pagu-project/pagu/internal/engine/command"
	"github.com/pagu-project/pagu/internal/entity"
)

func (c *CrowdfundCmd) createHandler(
	caller *entity.User,
	cmd *command.Command,
	args map[string]string,
) command.CommandResult {
	title := args["title"]
	desc := args["desc"]
	packagesJSON := args["packages"]

	packages := []entity.Package{}
	err := json.Unmarshal([]byte(packagesJSON), &packages)
	if err != nil {
		return cmd.RenderErrorTemplate(err)
	}

	if title == "" {
		return cmd.RenderFailedTemplate("The title of the crowdfunding campaign cannot be empty")
	}

	if len(packages) < 2 {
		return cmd.RenderFailedTemplate("At least 3 packages are required for the crowdfunding campaign")
	}

	campaign := &entity.CrowdfundCampaign{
		CreatorID: caller.ID,
		Title:     title,
		Desc:      desc,
		Packages:  packages,
	}
	err = c.db.AddCrowdfundCampaign(campaign)
	if err != nil {
		return cmd.RenderErrorTemplate(err)
	}

	return cmd.RenderResultTemplate("campaign", campaign)
}
