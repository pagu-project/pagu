package crowdfund

import (
	"encoding/json"

	"github.com/pagu-project/pagu/internal/engine/command"
	"github.com/pagu-project/pagu/internal/entity"
)

func (c *CrowdfundCmd) createHandler(
	_ *entity.User,
	cmd *command.Command,
	args map[string]string,
) command.CommandResult {
	title := args["title"]
	desc := args["desc"]
	packagesJSON := args["packages"]

	packages := []entity.Package{}
	err := json.Unmarshal([]byte(packagesJSON), &packages)
	if err != nil {
		return cmd.FailedResult(err.Error())
	}

	if title == "" {
		return cmd.FailedResult("The title of the crowdfunding campaign cannot be empty")
	}

	if len(packages) < 2 {
		return cmd.FailedResult("At least 3 packages are required for the crowdfunding campaign")
	}

	campaign := &entity.CrowdfundCampaign{
		Title:    title,
		Desc:     desc,
		Packages: packages,
	}
	err = c.db.AddCrowdfundCampaign(campaign)
	if err != nil {
		return cmd.FailedResult(err.Error())
	}

	return cmd.SuccessfulResultF(
		"Crowdfund campaign '%s' created successfully with %d packages",
		title, len(packages))
}
