package market

import (
	"github.com/pagu-project/pagu/internal/engine/command"
	"github.com/pagu-project/pagu/internal/entity"
	"github.com/pagu-project/pagu/pkg/cache"
	"github.com/pagu-project/pagu/pkg/client"
)

type MarketCmd struct {
	*marketSubCmds

	clientMgr  client.IManager
	priceCache cache.Cache[string, entity.Price]
}

func NewMarketCmd(clientMgr client.IManager, priceCache cache.Cache[string, entity.Price]) *MarketCmd {
	return &MarketCmd{
		clientMgr:  clientMgr,
		priceCache: priceCache,
	}
}

func (c *MarketCmd) GetCommand() *command.Command {

	cmd := c.buildMarketCommand()
	cmd.AppIDs = entity.AllAppIDs()
	cmd.TargetFlag = command.TargetMaskMainnet

	c.subCmdPrice.AppIDs = entity.AllAppIDs()
	c.subCmdPrice.TargetFlag = command.TargetMaskMainnet

	return cmd
}
