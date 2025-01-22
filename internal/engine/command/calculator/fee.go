package calculator

import (
	"github.com/pagu-project/pagu/internal/engine/command"
	"github.com/pagu-project/pagu/internal/entity"
	"github.com/pagu-project/pagu/pkg/amount"
)

func (c *CalculatorCmd) feeHandler(
	_ *entity.User,
	cmd *command.Command,
	args map[string]string,
) command.CommandResult {
	amt, err := amount.FromString(args["amount"])
	if err != nil {
		return cmd.RenderFailedTemplate("Invalid amount param")
	}

	fee, err := c.clientMgr.GetFee(amt.ToNanoPAC())
	if err != nil {
		return cmd.RenderErrorTemplate(err)
	}

	feeAmount := amount.Amount(fee)

	return cmd.RenderResultTemplate("amount", amt, "fee", feeAmount)
}
