package calculator

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/pagu-project/pagu/internal/engine/command"
	"github.com/pagu-project/pagu/internal/entity"
	"github.com/pagu-project/pagu/pkg/amount"
	"github.com/pagu-project/pagu/pkg/utils"
)

func (c *CalculatorCmd) rewardHandler(
	_ *entity.User,
	cmd *command.Command,
	args map[string]string,
) command.CommandResult {
	stake, err := amount.FromString(args["stake"])
	if err != nil {
		return cmd.RenderFailedTemplate("Invalid stake param")
	}

	minStake, _ := amount.NewAmount(1)
	maxStake, _ := amount.NewAmount(1000)
	if stake < minStake || stake > maxStake {
		return cmd.RenderErrorTemplate(
			fmt.Errorf("%v is invalid amount, minimum stake amount is 1 PAC and maximum is 1,000 PAC", stake))
	}

	numOfDays, err := strconv.Atoi(args["days"])
	if err != nil {
		return cmd.RenderErrorTemplate(errors.New("invalid days param"))
	}

	if numOfDays < 1 || numOfDays > 365 {
		return cmd.RenderErrorTemplate(
			fmt.Errorf("%v is invalid time, minimum time value is 1 and maximum is 365", numOfDays))
	}

	blocks := numOfDays * 8640
	info, err := c.clientMgr.GetBlockchainInfo()
	if err != nil {
		return cmd.RenderErrorTemplate(err)
	}

	reward := (stake.ToNanoPAC() * int64(blocks)) / info.TotalPower

	return cmd.RenderResultTemplate("stake", stake,
		"days", numOfDays,
		"totalPower", utils.FormatNumber(int64(amount.Amount(info.TotalPower).ToPAC())),
		"reward", utils.FormatNumber(reward))
}
