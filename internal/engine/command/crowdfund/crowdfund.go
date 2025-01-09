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
	*crowdfundSubCmds

	ctx         context.Context
	db          *repository.Database
	wallet      wallet.IWallet
	nowPayments nowpayments.INowPayments
}

func NewCrowdfundCmd(ctx context.Context,
	db *repository.Database,
	wallet wallet.IWallet,
	nowPayments nowpayments.INowPayments,
) *CrowdfundCmd {
	return &CrowdfundCmd{
		ctx:         ctx,
		db:          db,
		wallet:      wallet,
		nowPayments: nowPayments,
	}
}

func (c *CrowdfundCmd) activeCampaign() *entity.CrowdfundCampaign {
	return c.db.GetActiveCrowdfundCampaign()
}

func (c *CrowdfundCmd) GetCommand() *command.Command {
	middlewareHandler := command.NewMiddlewareHandler(c.db, c.wallet)
	cmd := c.buildCrowdfundCommand()

	cmd.AppIDs = entity.AllAppIDs()
	cmd.TargetFlag = command.TargetMaskModerator | command.TargetMaskMainnet

	c.subCmdCreate.AppIDs = []entity.PlatformID{entity.PlatformIDCLI, entity.PlatformIDDiscord}
	c.subCmdCreate.TargetFlag = command.TargetMaskModerator
	c.subCmdCreate.Middlewares = []command.MiddlewareFunc{middlewareHandler.OnlyModerator}

	c.subCmdDisable.AppIDs = []entity.PlatformID{entity.PlatformIDCLI, entity.PlatformIDDiscord}
	c.subCmdDisable.TargetFlag = command.TargetMaskModerator
	c.subCmdDisable.Middlewares = []command.MiddlewareFunc{middlewareHandler.OnlyModerator}

	c.subCmdReport.AppIDs = entity.AllAppIDs()
	c.subCmdReport.TargetFlag = command.TargetMaskMainnet

	c.subCmdInfo.AppIDs = entity.AllAppIDs()
	c.subCmdInfo.TargetFlag = command.TargetMaskMainnet

	c.subCmdPurchase.AppIDs = entity.AllAppIDs()
	c.subCmdPurchase.TargetFlag = command.TargetMaskMainnet

	c.subCmdClaim.AppIDs = entity.AllAppIDs()
	c.subCmdClaim.TargetFlag = command.TargetMaskMainnet

	activeCampaign := c.activeCampaign()
	if activeCampaign != nil {
		purchaseChoices := []command.Choice{}
		for index, pkg := range activeCampaign.Packages {
			choice := command.Choice{
				Name:  pkg.Name,
				Value: index + 1,
			}

			purchaseChoices = append(purchaseChoices, choice)
		}
		c.subCmdCreate.Args[0].Choices = purchaseChoices
	}

	return cmd
}
