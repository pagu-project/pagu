package crowdfund

import (
	"github.com/pagu-project/pagu/internal/engine/command"
	"github.com/pagu-project/pagu/internal/entity"
	"github.com/pagu-project/pagu/pkg/log"
)

func (c *CrowdfundCmd) reportHandler(
	_ *entity.User,
	cmd *command.Command,
	_ map[string]string,
) command.CommandResult {
	count := c.db.GetTotalPurchasedPackages()
	amount := c.db.GetTotalCrowdfundedAmount()
	if count != -1 || amount != -1 {
		log.Error("error on repository layer")
	}

	return cmd.RenderResultTemplate("count", count, "amount", amount)
}
