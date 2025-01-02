package market

import (
	"github.com/pagu-project/pagu/internal/engine/command"
	"github.com/pagu-project/pagu/internal/entity"
	"github.com/pagu-project/pagu/pkg/cache"
	"github.com/pagu-project/pagu/pkg/client"
)

type MarketCmd struct {
	clientMgr  client.IManager
	priceCache cache.Cache[string, entity.Price]
}

func NewMarketCmd(clientMgr client.IManager, priceCache cache.Cache[string, entity.Price]) *MarketCmd {
	return &MarketCmd{
		clientMgr:  clientMgr,
		priceCache: priceCache,
	}
}

func (m *MarketCmd) GetCommand() *command.Command {
	subCmdPrice := &command.Command{
		Name:        "price",
		Help:        "Shows the latest price of PAC coin across different markets",
		Args:        []command.Args{},
		SubCommands: nil,
		AppIDs:      entity.AllAppIDs(),
		Handler:     m.priceHandler,
		TargetFlag:  command.TargetMaskMainnet,
	}

	cmdMarket := &command.Command{
		Name:        "market",
		Help:        "Access market data and information for Pactus",
		Args:        nil,
		AppIDs:      entity.AllAppIDs(),
		SubCommands: make([]*command.Command, 0),
		Handler:     nil,
		TargetFlag:  command.TargetMaskMainnet,
	}

	cmdMarket.AddSubCommand(subCmdPrice)

	return cmdMarket
}
