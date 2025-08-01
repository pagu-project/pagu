// Code generated by command-generator. DO NOT EDIT.
package calculator

import (
	"github.com/pagu-project/pagu/internal/engine/command"
	"github.com/pagu-project/pagu/internal/entity"
)

const (
	argNameRewardStake = "stake"
	argNameRewardDays  = "days"
	argNameFeeAmount   = "amount"
)

type calculatorSubCmds struct {
	subCmdReward *command.Command
	subCmdFee    *command.Command
}

func (c *CalculatorCmd) buildSubCmds() *calculatorSubCmds {
	subCmdReward := &command.Command{
		Name:            "reward",
		Help:            "Calculate the PAC coins you can earn based on your validator stake",
		Handler:         c.rewardHandler,
		ResultTemplate:  "Approximately you earn {{.reward}} PAC reward, with {{.stake}} stake 🔒 on your validator in {{.days}} days ⏰ with {{.totalPower}} total power ⚡ of committee.\n\n> Note📝: This number is just an estimation. It will vary depending on your stake amount and total network power.\n",
		TargetBotIDs:    entity.AllBotIDs(),
		TargetUserRoles: entity.AllUserRoles(),
		Args: []*command.Args{
			{
				Name:     "stake",
				Desc:     "The amount of stake in your validator",
				InputBox: command.InputBoxInteger,
				Optional: false,
			},
			{
				Name:     "days",
				Desc:     "The number of days to calculate rewards for (range : 1-365)",
				InputBox: command.InputBoxInteger,
				Optional: false,
			},
		},
	}
	subCmdFee := &command.Command{
		Name:            "fee",
		Help:            "Return the estimated transaction fee on the network",
		Handler:         c.feeHandler,
		ResultTemplate:  "Sending {{.amount}} will cost {{.fee}} with current fee percentage.\n",
		TargetBotIDs:    entity.AllBotIDs(),
		TargetUserRoles: entity.AllUserRoles(),
		Args: []*command.Args{
			{
				Name:     "amount",
				Desc:     "The amount of PAC coins to calculate fee for",
				InputBox: command.InputBoxInteger,
				Optional: false,
			},
		},
	}

	return &calculatorSubCmds{
		subCmdReward: subCmdReward,
		subCmdFee:    subCmdFee,
	}
}

func (c *CalculatorCmd) buildCalculatorCommand(botID entity.BotID) *command.Command {
	calculatorCmd := &command.Command{
		Name:            "calculator",
		Emoji:           "🧮",
		Active:          true,
		Help:            "Perform calculations such as reward and fee estimations",
		SubCommands:     make([]*command.Command, 0),
		TargetBotIDs:    entity.AllBotIDs(),
		TargetUserRoles: entity.AllUserRoles(),
	}

	c.calculatorSubCmds = c.buildSubCmds()

	calculatorCmd.AddSubCommand(botID, c.subCmdReward)
	calculatorCmd.AddSubCommand(botID, c.subCmdFee)

	return calculatorCmd
}
