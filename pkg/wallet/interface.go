package wallet

import "github.com/pagu-project/pagu/pkg/amount"

type IWallet interface {
	Balance() amount.Amount
	Address() string
	TransferTransaction(toAddress, memo string, amt amount.Amount) (string, error)
	BondTransaction(pubKey, toAddress, memo string, amt amount.Amount) (string, error)
	LinkToExplorer(txID string) string
}
