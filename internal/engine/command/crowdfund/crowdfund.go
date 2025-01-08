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
