package crowdfund

import (
	"context"
	"fmt"

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
	return c.db.GetCrowdfundActiveCampaign()
}

func (c *CrowdfundCmd) GetCommand() *command.Command {
	middlewareHandler := command.NewMiddlewareHandler(c.db, c.wallet)
	cmd := c.buildCrowdfundCommand()

	cmd.PlatformIDs = entity.AllPlatformIDs()
	cmd.TargetFlag = command.TargetMaskModerator | command.TargetMaskMainnet

	c.subCmdCreate.PlatformIDs = []entity.PlatformID{entity.PlatformIDCLI, entity.PlatformIDDiscord}
	c.subCmdCreate.TargetFlag = command.TargetMaskModerator
	c.subCmdCreate.Middlewares = []command.MiddlewareFunc{middlewareHandler.OnlyModerator}

	c.subCmdDisable.PlatformIDs = []entity.PlatformID{entity.PlatformIDCLI, entity.PlatformIDDiscord}
	c.subCmdDisable.TargetFlag = command.TargetMaskModerator
	c.subCmdDisable.Middlewares = []command.MiddlewareFunc{middlewareHandler.OnlyModerator}

	c.subCmdReport.PlatformIDs = entity.AllPlatformIDs()
	c.subCmdReport.TargetFlag = command.TargetMaskMainnet

	c.subCmdInfo.PlatformIDs = entity.AllPlatformIDs()
	c.subCmdInfo.TargetFlag = command.TargetMaskMainnet

	c.subCmdPurchase.PlatformIDs = entity.AllPlatformIDs()
	c.subCmdPurchase.TargetFlag = command.TargetMaskMainnet

	c.subCmdClaim.PlatformIDs = entity.AllPlatformIDs()
	c.subCmdClaim.TargetFlag = command.TargetMaskMainnet

	activeCampaign := c.activeCampaign()
	if activeCampaign != nil {
		purchaseChoices := []command.Choice{}
		for index, pkg := range activeCampaign.Packages {
			choice := command.Choice{
				Name:  fmt.Sprintf("%s (%d USDT to %s)", pkg.Name, pkg.USDAmount, pkg.PACAmount),
				Value: fmt.Sprintf("%d", index+1),
			}

			purchaseChoices = append(purchaseChoices, choice)
		}
		c.subCmdPurchase.Args[0].Choices = purchaseChoices
	}

	return cmd
}
