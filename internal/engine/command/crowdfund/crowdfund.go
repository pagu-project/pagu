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
		TargetFlag:  command.TargetMaskModerator,
	}
	subCmdDisable := &command.Command{
		Name:        "disable",
		Help:        "Disable an existing crowdfunding campaign",
		Args:        []command.Args{},
		SubCommands: nil,
		AppIDs:      entity.AllAppIDs(),
		Handler:     c.disableHandler,
		TargetFlag:  command.TargetMaskModerator,
	}
	subCmdReport := &command.Command{
		Name:        "report",
		Help:        "View reports of a crowdfunding campaign",
		Args:        []command.Args{},
		SubCommands: nil,
		AppIDs:      entity.AllAppIDs(),
		Handler:     c.reportHandler,
		TargetFlag:  command.TargetMaskModerator,
	}
	subCmdInfo := &command.Command{
		Name:        "info",
		Help:        "Get detailed information about a crowdfunding campaign",
		Args:        []command.Args{},
		SubCommands: nil,
		AppIDs:      entity.AllAppIDs(),
		Handler:     c.infoHandler,
		TargetFlag:  command.TargetMaskModerator,
	}
	subCmdPurchase := &command.Command{
		Name: "purchase",
		Help: "Make a purchase in a crowdfunding campaign",
		Args: []command.Args{
			{
				Name:     "package",
				Desc:     "Select the crowdfunding package",
				InputBox: command.InputBoxChoice,
				Optional: false,
				Choices: []command.Choice{
					{Name: "Package 1", Value: 1},
					{Name: "Package 2", Value: 2},
					{Name: "Package 3", Value: 3},
				},
			},
		},
		SubCommands: nil,
		AppIDs:      entity.AllAppIDs(),
		Handler:     c.purchaseHandler,
		TargetFlag:  command.TargetMaskModerator,
	}
	subCmdClaim := &command.Command{
		Name:        "claim",
		Help:        "Claim packages from a crowdfunding campaign",
		Args:        []command.Args{},
		SubCommands: nil,
		AppIDs:      entity.AllAppIDs(),
		Handler:     c.claimHandler,
		TargetFlag:  command.TargetMaskModerator,
	}

	cmdCrowdfund := &command.Command{
		Name:        "crowdfund",
		Help:        "Commands for managing crowdfunding campaigns",
		Args:        nil,
		AppIDs:      entity.AllAppIDs(),
		SubCommands: make([]*command.Command, 0),
		Handler:     nil,
		TargetFlag:  command.TargetMaskModerator,
	}

	cmdCrowdfund.AddSubCommand(subCmdCreate)
	cmdCrowdfund.AddSubCommand(subCmdDisable)
	cmdCrowdfund.AddSubCommand(subCmdReport)
	cmdCrowdfund.AddSubCommand(subCmdInfo)
	cmdCrowdfund.AddSubCommand(subCmdPurchase)
	cmdCrowdfund.AddSubCommand(subCmdClaim)

	return cmdCrowdfund
}
