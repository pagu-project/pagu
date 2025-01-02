package calculator

import (
	"errors"

	"github.com/pagu-project/pagu/internal/engine/command"
	"github.com/pagu-project/pagu/internal/entity"
	"github.com/pagu-project/pagu/pkg/amount"
)

func (c *CalculatorCmd) calcFeeHandler(
	_ *entity.User,
	cmd *command.Command,
	args map[string]string,
) command.CommandResult {
	amt, err := amount.FromString(args["amount"])
	if err != nil {
		return cmd.ErrorResult(errors.New("invalid amount param"))
	}

	fee, err := c.clientMgr.GetFee(amt.ToNanoPAC())
	if err != nil {
		return cmd.ErrorResult(err)
	}

	feeAmount := amount.Amount(fee)

	return cmd.SuccessfulResultF("Sending %s will cost %s with current fee percentage."+
		"\n> Note: Consider unbond and sortition transaction fee is 0 PAC always.", amt.String(), feeAmount.String())
}
