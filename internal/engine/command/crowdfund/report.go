package crowdfund

import (
	"github.com/pagu-project/pagu/internal/engine/command"
	"github.com/pagu-project/pagu/internal/entity"
)

func (c *CrowdfundCmd) reportHandler(
	_ *entity.User,
	cmd *command.Command,
	_ map[string]string,
) command.CommandResult {
	count, err := c.db.GetTotalPurchasedPackages()
	if err != nil {
		return cmd.RenderErrorTemplate(err)
	}

	amount, err := c.db.GetTotalCrowdfundedAmount()
	if err != nil {
		return cmd.RenderErrorTemplate(err)
	}

	return cmd.RenderResultTemplate("count", count, "amount", amount)
}
