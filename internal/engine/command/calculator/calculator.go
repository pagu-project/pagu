package calculator

import (
	"github.com/pagu-project/pagu/internal/engine/command"
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

	return cmd
}
