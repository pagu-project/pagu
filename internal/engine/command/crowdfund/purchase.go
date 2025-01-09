package crowdfund

import (
	"fmt"
	"strconv"

	"github.com/pagu-project/pagu/internal/engine/command"
	"github.com/pagu-project/pagu/internal/entity"
)

func (c *CrowdfundCmd) purchaseHandler(
	user *entity.User,
	cmd *command.Command,
	args map[string]string,
) command.CommandResult {
	activeCampaign := c.activeCampaign()
	if activeCampaign == nil {
		return cmd.RenderFailedTemplate("No active campaign")
	}

	pkgIndex, _ := strconv.Atoi(args[argNamePurchasePackage])
	pkg := activeCampaign.Packages[pkgIndex]

	purchase := &entity.CrowdfundPurchase{
		UserID:    user.ID,
		USDAmount: pkg.USDAmount,
		PACAmount: pkg.PACAmount,
	}
	err := c.db.AddCrowdfundPurchase(purchase)
	if err != nil {
		return cmd.RenderErrorTemplate(err)
	}
	orderID := fmt.Sprintf("crowdfund/%d", purchase.ID)
	invoiceID, err := c.nowPayments.CreateInvoice(pkg.USDAmount, orderID)
	if err != nil {
		return cmd.RenderErrorTemplate(err)
	}

	purchase.InvoiceID = invoiceID
	err = c.db.UpdateCrowdfundPurchase(purchase)
	if err != nil {
		return cmd.RenderErrorTemplate(err)
	}

	return cmd.RenderResultTemplate(
		"purchase", purchase,
		"paymentLink", c.nowPayments.PaymentLink(purchase.InvoiceID))
}
