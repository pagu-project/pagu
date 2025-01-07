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
	ctx            context.Context
	db             *repository.Database
	wallet         wallet.IWallet
	nowPayments    nowpayments.INowpayments
	activeCampaign *entity.CrowdfundCampaign
}

func NewCrowdfundCmd(ctx context.Context,
	db *repository.Database,
	wallet wallet.IWallet,
	nowPayments nowpayments.INowpayments,
) *CrowdfundCmd {
	return &CrowdfundCmd{
		ctx:            ctx,
		activeCampaign: nil,
		db:             db,
		wallet:         wallet,
		nowPayments:    nowPayments,
	}
}

func (c *CrowdfundCmd) GetCommand() *command.Command {
	if c.activeCampaign == nil {
		return nil
	}

	purchaseChoices := []command.Choice{}
	for index, pkg := range c.activeCampaign.Packages {
		choice := command.Choice{
			Name:  pkg.Name,
			Value: index,
		}

		purchaseChoices = append(purchaseChoices, choice)
	}
	subCmdPurchase.Args[0].Choices = purchaseChoices

	return c.crowdfundCommand()
}
