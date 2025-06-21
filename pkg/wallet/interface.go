package wallet

import "github.com/pagu-project/pagu/pkg/amount"

type IWallet interface {
	Balance() amount.Amount
	Address() string
	TransferTransaction(toAddress string, amt amount.Amount, memo string) (string, error)
	BondTransaction(pubKey, toAddress string, amt amount.Amount, memo string) (string, error)
	LinkToExplorer(txID string) string
}
