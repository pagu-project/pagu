package nowpayments

// INowPayments defines the interface for interacting with the nowpayments.io service.
// It provides methods for creating payments and checking their status.
type INowPayments interface {
	// CreateInvoice creates an invoice with the specified price in USD and an associated order ID.
	// It returns the generated invoice ID or an error if the operation fails.
	CreateInvoice(priceUSD int, orderID string) (string, error)

	// IsPaid checks the invoice status of a given invoice ID.
	// It verifies whether the invoice has been successfully paid and completed.
	IsPaid(invoiceID string) (bool, error)

	// PaymentLink creates the payment link for the given invoice ID.
	PaymentLink(invoiceID string) string
}
