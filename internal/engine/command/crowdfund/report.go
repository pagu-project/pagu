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
	count := c.db.GetTotalPurchasedPackages()
	amount := c.db.GetTotalCrowdfundedAmount()

	return cmd.RenderResultTemplate("count", count, "amount", amount)
}
