package crowdfund

import (
	"fmt"
	"strconv"

	"github.com/pagu-project/pagu/internal/engine/command"
	"github.com/pagu-project/pagu/internal/entity"
	"github.com/pagu-project/pagu/pkg/log"
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

	pkgNumber, _ := strconv.Atoi(args[argNamePurchasePackage])
	pkgIndex := pkgNumber - 1
	if pkgIndex == -1 || pkgIndex >= len(activeCampaign.Packages) {
		return cmd.RenderFailedTemplateF("Invalid package number: %d", pkgNumber)
	}
	pkg := activeCampaign.Packages[pkgIndex]

	purchase := &entity.CrowdfundPurchase{
		UserID:    user.ID,
		USDAmount: pkg.USDAmount,
		PACAmount: pkg.PACAmount,
	}
	err := c.db.AddCrowdfundPurchase(purchase)
	if err != nil {
		log.Error("database failed", "error", err, "purchase", purchase)

		return cmd.RenderInternalFailure()
	}

	orderID := fmt.Sprintf("crowdfund/%d", purchase.ID)
	invoiceID, err := c.nowPayments.CreateInvoice(pkg.USDAmount, orderID)
	if err != nil {
		log.Error("NowPayments failed", "error", err, "orderID", orderID)

		return cmd.RenderInternalFailure()
	}

	purchase.InvoiceID = invoiceID
	err = c.db.UpdateCrowdfundPurchase(purchase)
	if err != nil {
		log.Error("database failed", "error", err, "purchase", purchase)

		return cmd.RenderInternalFailure()
	}

	return cmd.RenderResultTemplate(
		"purchase", purchase,
		"paymentLink", c.nowPayments.PaymentLink(purchase.InvoiceID))
}
