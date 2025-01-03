package crowdfund

import (
	"context"

	"github.com/pagu-project/pagu/internal/engine/command"
	"github.com/pagu-project/pagu/internal/entity"
	"github.com/pagu-project/pagu/internal/repository"
	"github.com/pagu-project/pagu/pkg/nowpayments"
	"github.com/pagu-project/pagu/pkg/wallet"
)

type CrowdfundCmd struct {
	ctx         context.Context
	db          *repository.Database
	wallet      wallet.IWallet
	nowPayments nowpayments.INowpayments
}

func NewCrowdfundCmd(ctx context.Context,
	db *repository.Database,
	wallet wallet.IWallet,
	nowPayments nowpayments.INowpayments) *CrowdfundCmd {
	return &CrowdfundCmd{
		ctx:         ctx,
		db:          db,
		wallet:      wallet,
		nowPayments: nowPayments,
	}
}

func (n *CrowdfundCmd) GetCommand() *command.Command {
	subCmdCreate := &command.Command{
		Name: "create",
		Help: "Create a new crowdfunding campaign",
		Args: []command.Args{
			{
				Name:     "title",
				Desc:     "The title of this crowdfunding campaign",
				InputBox: command.InputBoxText,
				Optional: false,
			},
			{
				Name:     "desc",
				Desc:     "A description of this crowdfunding campaign",
				InputBox: command.InputBoxMultilineText,
				Optional: false,
			},
			{
				Name:     "packages",
				Desc:     "The packages for this campaign in JSON format",
				InputBox: command.InputBoxMultilineText,
				Optional: false,
			},
		},
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
