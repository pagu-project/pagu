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

func (c *CrowdfundCmd) lastCampaign() *entity.CrowdfundCampaign {
	return c.db.GetCrowdfundLastCampaign()
}

func (c *CrowdfundCmd) activeCampaign() *entity.CrowdfundCampaign {
	return c.db.GetCrowdfundActiveCampaign()
}

func (c *CrowdfundCmd) GetCommand() *command.Command {
	cmd := c.buildCrowdfundCommand()

	activeCampaign := c.activeCampaign()
	if activeCampaign != nil {
		purchaseChoices := []command.Choice{}
		for index, pkg := range activeCampaign.Packages {
			choice := command.Choice{
				Name:  pkg.Name,
				Desc:  fmt.Sprintf("**%s**: (%d USDT to %s)", pkg.Name, pkg.USDAmount, pkg.PACAmount.String()),
				Value: fmt.Sprintf("%d", index+1),
			}

			purchaseChoices = append(purchaseChoices, choice)
		}
		c.subCmdPurchase.Args[0].Choices = purchaseChoices
	}

	return cmd
}
