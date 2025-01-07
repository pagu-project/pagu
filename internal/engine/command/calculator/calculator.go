package calculator

import (
	"github.com/pagu-project/pagu/internal/engine/command"
	"github.com/pagu-project/pagu/internal/entity"
	"github.com/pagu-project/pagu/pkg/client"
)

type CalculatorCmd struct {
	clientMgr client.IManager
}

func NewCalculatorCmd(clientMgr client.IManager) *CalculatorCmd {
	return &CalculatorCmd{
		clientMgr: clientMgr,
	}
}

func (c *CalculatorCmd) GetCommand() *command.Command {
	subCmdCalcReward := &command.Command{
		Name: "reward",
		Help: "Calculate the PAC coins you can earn based on your validator stake",
		Args: []command.Args{
			{
				Name:     "stake",
				Desc:     "The amount of stake in your validator",
				InputBox: command.InputBoxInteger,
				Optional: false,
			},
			{
				Name:     "days",
				Desc:     "The number of days to calculate rewards for (range: 1-365)",
				InputBox: command.InputBoxInteger,
				Optional: false,
			},
		},
		SubCommands: nil,
		AppIDs:      entity.AllAppIDs(),
		Handler:     c.calcRewardHandler,
		TargetFlag:  command.TargetMaskMainnet,
	}

	subCmdCalcFee := &command.Command{
		Name:        "fee",
		Help:        "Return the estimated transaction fee on the network",
		SubCommands: nil,
		AppIDs:      entity.AllAppIDs(),
		Handler:     c.calcFeeHandler,
		TargetFlag:  command.TargetMaskMainnet,
	}

	cmdBlockchain := &command.Command{
		Name:        "calculate",
		Help:        "Perform calculations such as reward and fee estimations",
		Args:        nil,
		AppIDs:      entity.AllAppIDs(),
		SubCommands: make([]*command.Command, 0),
		Handler:     nil,
		TargetFlag:  command.TargetMaskMainnet,
	}

	cmdBlockchain.AddSubCommand(subCmdCalcReward)
	cmdBlockchain.AddSubCommand(subCmdCalcFee)

	return cmdBlockchain
}
