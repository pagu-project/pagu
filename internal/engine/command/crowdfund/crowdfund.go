package crowdfund

import (
	"context"

	"github.com/pagu-project/pagu/internal/engine/command"
	"github.com/pagu-project/pagu/internal/entity"
	"github.com/pagu-project/pagu/internal/nowpayments"
)

const (
	CommandName            = "crowdfund"
	subCommandNameCreate   = "Create"
	subCommandNameInfo     = "info"
	subCommandNamePurchase = "purchase"
	subCommandNameClaim    = "claim"
	subCommandNameDisable  = "disable"
	subCommandNameReport   = "report"
)

type Crowdfund struct {
	ctx         context.Context
	nowPayments nowpayments.INowpayments
}

func NewCrowdfundCmd(ctx context.Context, nowPayments nowpayments.INowpayments) *Crowdfund {
	return &Crowdfund{
		ctx:         ctx,
		nowPayments: nowPayments,
	}
}

func (n *Crowdfund) GetCommand() *command.Command {
	subCmdCreate := &command.Command{
		Name:        subCommandNameCreate,
		Help:        "Create a new crowdfund campaign",
		Args:        []command.Args{},
		SubCommands: nil,
		AppIDs:      entity.AllAppIDs(),
		Handler:     n.handlerCreate,
		TargetFlag:  command.TargetMaskAll,
	}
	subCmdDisable := &command.Command{
		Name:        subCommandNameDisable,
		Help:        "Disable an active crowdfund campaign",
		Args:        []command.Args{},
		SubCommands: nil,
		AppIDs:      entity.AllAppIDs(),
		Handler:     n.handlerDisable,
		TargetFlag:  command.TargetMaskAll,
	}
	subCmdReport := &command.Command{
		Name:        subCommandNameReport,
		Help:        "Report of the crowdfund campaigns",
		Args:        []command.Args{},
		SubCommands: nil,
		AppIDs:      entity.AllAppIDs(),
		Handler:     n.handlerReport,
		TargetFlag:  command.TargetMaskAll,
	}

	subCmdInfo := &command.Command{
		Name:        subCommandNameInfo,
		Help:        "get information about the crowdfund campaigns",
		Args:        []command.Args{},
		SubCommands: nil,
		AppIDs:      entity.AllAppIDs(),
		Handler:     n.handlerInfo,
		TargetFlag:  command.TargetMaskAll,
	}
	subCmdPurchase := &command.Command{
		Name:        subCommandNamePurchase,
		Help:        "Purchase a package",
		Args:        []command.Args{},
		SubCommands: nil,
		AppIDs:      entity.AllAppIDs(),
		Handler:     n.handlerPurchase,
		TargetFlag:  command.TargetMaskAll,
	}
	subCmdClaim := &command.Command{
		Name:        subCommandNameClaim,
		Help:        "Claim your purchase",
		Args:        []command.Args{},
		SubCommands: nil,
		AppIDs:      entity.AllAppIDs(),
		Handler:     n.handlerClaim,
		TargetFlag:  command.TargetMaskAll,
	}

	cmdNetwork := &command.Command{
		Name:        CommandName,
		Help:        "Network related commands",
		Args:        nil,
		AppIDs:      entity.AllAppIDs(),
		SubCommands: make([]*command.Command, 0),
		Handler:     nil,
		TargetFlag:  command.TargetMaskAll,
	}

	cmdNetwork.AddSubCommand(subCmdCreate)
	cmdNetwork.AddSubCommand(subCmdDisable)
	cmdNetwork.AddSubCommand(subCmdReport)
	cmdNetwork.AddSubCommand(subCmdInfo)
	cmdNetwork.AddSubCommand(subCmdPurchase)
	cmdNetwork.AddSubCommand(subCmdClaim)

	return cmdNetwork
}
