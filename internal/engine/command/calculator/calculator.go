package calculator

import (
	"github.com/pagu-project/pagu/internal/engine/command"
	"github.com/pagu-project/pagu/internal/entity"
	"github.com/pagu-project/pagu/pkg/client"
)

type CalculatorCmd struct {
	*calculatorSubCmds

	clientMgr client.IManager
}

func NewCalculatorCmd(clientMgr client.IManager) *CalculatorCmd {
	return &CalculatorCmd{
		clientMgr: clientMgr,
	}
}

func (c *CalculatorCmd) GetCommand() *command.Command {
	cmd := c.buildCalculatorCommand()
	cmd.PlatformIDs = entity.AllPlatformIDs()
	cmd.TargetFlag = command.TargetMaskMainnet

	c.subCmdReward.PlatformIDs = entity.AllPlatformIDs()
	c.subCmdReward.TargetFlag = command.TargetMaskMainnet

	c.subCmdFee.PlatformIDs = entity.AllPlatformIDs()
	c.subCmdFee.TargetFlag = command.TargetMaskMainnet

	return cmd
}
