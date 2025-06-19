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

func (c *CalculatorCmd) BuildCommand(botID entity.BotID) *command.Command {
	cmd := c.buildCalculatorCommand(botID)

	return cmd
}
