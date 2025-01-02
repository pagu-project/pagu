package nowpayments

// INowpayments defines the interface for interacting with the nowpayments.io service.
// It provides methods for creating payments and checking their status.
type INowpayments interface {
	// CreatePayment initializes a payment with the specified price in USD and an associated order ID.
	// It returns the generated invoice ID for the payment or an error if the operation fails.
	CreatePayment(priceUSD int, orderID string) (string, error)

	// IsPaid checks the payment status of a given invoice ID.
	// It verifies whether the payment has been successfully completed.
	IsPaid(invoiceID string) (bool, error)
}
