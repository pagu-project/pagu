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
	config      *Config
	db          *repository.Database
	wallet      wallet.IWallet
	nowPayments nowpayments.INowpayments
}

func NewCrowdfundCmd(ctx context.Context,
	config *Config,
	db *repository.Database,
	wallet wallet.IWallet,
	nowPayments nowpayments.INowpayments,
) *CrowdfundCmd {
	return &CrowdfundCmd{
		ctx:         ctx,
		config:      config,
		db:          db,
		wallet:      wallet,
		nowPayments: nowPayments,
	}
}

func (c *CrowdfundCmd) GetCommand() *command.Command {
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
		Handler:     c.createHandler,
		TargetFlag:  command.TargetMaskAll,
	}
	subCmdDisable := &command.Command{
		Name:        "disable",
		Help:        "Disable an existing crowdfunding campaign",
		Args:        []command.Args{},
		SubCommands: nil,
		AppIDs:      entity.AllAppIDs(),
		Handler:     c.disableHandler,
		TargetFlag:  command.TargetMaskAll,
	}
	subCmdReport := &command.Command{
		Name:        "report",
		Help:        "View reports of a crowdfunding campaign",
		Args:        []command.Args{},
		SubCommands: nil,
		AppIDs:      entity.AllAppIDs(),
		Handler:     c.reportHandler,
		TargetFlag:  command.TargetMaskAll,
	}
	subCmdInfo := &command.Command{
		Name:        "info",
		Help:        "Get detailed information about a crowdfunding campaign",
		Args:        []command.Args{},
		SubCommands: nil,
		AppIDs:      entity.AllAppIDs(),
		Handler:     c.infoHandler,
		TargetFlag:  command.TargetMaskAll,
	}
	subCmdPurchase := &command.Command{
		Name:        "purchase",
		Help:        "Make a purchase in a crowdfunding campaign",
		Args:        []command.Args{},
		SubCommands: nil,
		AppIDs:      entity.AllAppIDs(),
		Handler:     c.purchaseHandler,
		TargetFlag:  command.TargetMaskAll,
	}
	subCmdClaim := &command.Command{
		Name:        "claim",
		Help:        "Claim packages from a crowdfunding campaign",
		Args:        []command.Args{},
		SubCommands: nil,
		AppIDs:      entity.AllAppIDs(),
		Handler:     c.claimHandler,
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
