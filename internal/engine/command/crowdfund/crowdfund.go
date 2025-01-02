package crowdfund

import (
	"context"

	"github.com/pagu-project/pagu/internal/engine/command"
	"github.com/pagu-project/pagu/internal/entity"
	"github.com/pagu-project/pagu/internal/nowpayments"
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
		Name:        "create",
		Help:        "Create a new crowdfunding campaign",
		Args:        []command.Args{},
		SubCommands: nil,
		AppIDs:      entity.AllAppIDs(),
		Handler:     n.createHandler,
		TargetFlag:  command.TargetMaskAll,
	}
	subCmdDisable := &command.Command{
		Name:        "disable",
		Help:        "Disable an existing crowdfunding campaign",
		Args:        []command.Args{},
		SubCommands: nil,
		AppIDs:      entity.AllAppIDs(),
		Handler:     n.disableHandler,
		TargetFlag:  command.TargetMaskAll,
	}
	subCmdReport := &command.Command{
		Name:        "report",
		Help:        "View reports of a crowdfunding campaign",
		Args:        []command.Args{},
		SubCommands: nil,
		AppIDs:      entity.AllAppIDs(),
		Handler:     n.reportHandler,
		TargetFlag:  command.TargetMaskAll,
	}
	subCmdInfo := &command.Command{
		Name:        "info",
		Help:        "Get detailed information about a crowdfunding campaign",
		Args:        []command.Args{},
		SubCommands: nil,
		AppIDs:      entity.AllAppIDs(),
		Handler:     n.infoHandler,
		TargetFlag:  command.TargetMaskAll,
	}
	subCmdPurchase := &command.Command{
		Name:        "purchase",
		Help:        "Make a purchase in a crowdfunding campaign",
		Args:        []command.Args{},
		SubCommands: nil,
		AppIDs:      entity.AllAppIDs(),
		Handler:     n.purchaseHandler,
		TargetFlag:  command.TargetMaskAll,
	}
	subCmdClaim := &command.Command{
		Name:        "claim",
		Help:        "Claim packages from a crowdfunding campaign",
		Args:        []command.Args{},
		SubCommands: nil,
		AppIDs:      entity.AllAppIDs(),
		Handler:     n.claimHandler,
		TargetFlag:  command.TargetMaskAll,
	}

	cmdCrowdfund := &command.Command{
		Name:        "crowdfund",
		Help:        "Commands for managing crowdfunding campaigns",
		Args:        nil,
		AppIDs:      entity.AllAppIDs(),
		SubCommands: make([]*command.Command, 0),
		Handler:     nil,
		TargetFlag:  command.TargetMaskAll,
	}

	cmdCrowdfund.AddSubCommand(subCmdCreate)
	cmdCrowdfund.AddSubCommand(subCmdDisable)
	cmdCrowdfund.AddSubCommand(subCmdReport)
	cmdCrowdfund.AddSubCommand(subCmdInfo)
	cmdCrowdfund.AddSubCommand(subCmdPurchase)
	cmdCrowdfund.AddSubCommand(subCmdClaim)

	return cmdCrowdfund
}
