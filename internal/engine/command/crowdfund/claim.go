package crowdfund

import (
	"github.com/pagu-project/pagu/internal/engine/command"
	"github.com/pagu-project/pagu/internal/entity"
	"github.com/pagu-project/pagu/pkg/log"
)

func (c *CrowdfundCmd) claimHandler(
	user *entity.User,
	cmd *command.Command,
	args map[string]string,
) command.CommandResult {
	purchases, err := c.db.GetCrowdfundPurchases(user.ID)
	if err != nil {
		return cmd.RenderErrorTemplate(err)
	}

	txID := ""
	for _, purchase := range purchases {
		if purchase.IsClaimed() {
			continue
		}

		isPaid, err := c.nowPayments.IsPaid(purchase.InvoiceID)
		if err != nil {
			log.Error("nowPayments failed", "error", err, "purchase", purchase)

			return cmd.RenderInternalFailure()
		}

		if isPaid {
			address := args[argNameClaimAddress]

			txID, err = c.wallet.TransferTransaction(address, "Crowdfund campaign", purchase.PACAmount)
			if err != nil {
				log.Error("wallet failed", "error", err, "address", address)

				return cmd.RenderErrorTemplate(err)
			}

			purchase.TxHash = txID
			purchase.Recipient = address
			err = c.db.UpdateCrowdfundPurchase(purchase)
			if err != nil {
				log.Error("nowPayments failed", "error", err, "purchase", purchase)

				return cmd.RenderInternalFailure()
			}

			// Ensure we log it always
			log.Warn("payment successful", "txID", txID, "address", address, "amount", purchase.PACAmount)

			break
		}
	}

	if txID == "" {
		return cmd.RenderFailedTemplate("No unpaid purchase for this user")
	}

	txLink := c.wallet.LinkToExplorer(txID)

	return cmd.RenderResultTemplate("txLink", txLink)
}
