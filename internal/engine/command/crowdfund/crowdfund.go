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
	cmd := c.crowdfundCommand()

	subCmdCreate.AppIDs = []entity.PlatformID{entity.PlatformIDCLI, entity.PlatformIDDiscord}
	subCmdDisable.AppIDs = []entity.PlatformID{entity.PlatformIDCLI, entity.PlatformIDDiscord}
	subCmdReport.AppIDs = entity.AllAppIDs()
	subCmdInfo.AppIDs = entity.AllAppIDs()
	subCmdPurchase.AppIDs = entity.AllAppIDs()
	subCmdClaim.AppIDs = entity.AllAppIDs()

	subCmdCreate.TargetFlag = command.TargetMaskModerator
	subCmdDisable.TargetFlag = command.TargetMaskModerator
	subCmdReport.TargetFlag = command.TargetMaskMainnet
	subCmdInfo.TargetFlag = command.TargetMaskMainnet
	subCmdPurchase.TargetFlag = command.TargetMaskMainnet
	subCmdClaim.TargetFlag = command.TargetMaskMainnet

	activeCampaign := c.activeCampaign()
	if activeCampaign != nil {
		purchaseChoices := []command.Choice{}
		for index, pkg := range activeCampaign.Packages {
			choice := command.Choice{
				Name:  pkg.Name,
				Value: index,
			}

			purchaseChoices = append(purchaseChoices, choice)
		}
		subCmdPurchase.Args[0].Choices = purchaseChoices
	}

	return cmd
}
