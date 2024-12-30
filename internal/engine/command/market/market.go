package market

import (
	"github.com/pagu-project/pagu/internal/engine/command"
	"github.com/pagu-project/pagu/internal/entity"
	"github.com/pagu-project/pagu/pkg/cache"
	"github.com/pagu-project/pagu/pkg/client"
)

const (
	CommandName      = "market"
	PriceCommandName = "price"
	HelpCommandName  = "help"
)

type Market struct {
	clientMgr  client.IManager
	priceCache cache.Cache[string, entity.Price]
}

func NewMarket(clientMgr client.IManager, priceCache cache.Cache[string, entity.Price]) *Market {
	return &Market{
		clientMgr:  clientMgr,
		priceCache: priceCache,
	}
}

func (m *Market) GetCommand() *command.Command {
	subCmdPrice := &command.Command{
		Name:        PriceCommandName,
		Help:        "Shows the last price of PAC coin on the markets",
		Args:        []command.Args{},
		SubCommands: nil,
		AppIDs:      entity.AllAppIDs(),
		Handler:     m.getPrice,
		TargetFlag:  command.TargetMaskMainnet,
	}

	cmdMarket := &command.Command{
		Name:        CommandName,
		Help:        "Pactus market data and information",
		Args:        nil,
		AppIDs:      entity.AllAppIDs(),
		SubCommands: make([]*command.Command, 0),
		Handler:     nil,
		TargetFlag:  command.TargetMaskMainnet,
	}

	cmdMarket.AddSubCommand(subCmdPrice)

	return cmdMarket
}
