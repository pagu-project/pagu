package wallet

import (
	"fmt"

	"github.com/pactus-project/pactus/genesis"
	"github.com/pactus-project/pactus/wallet"
	"github.com/pagu-project/pagu/pkg/amount"
	"github.com/pagu-project/pagu/pkg/log"
)

type Wallet struct {
	*wallet.Wallet

	address  string
	password string
	fee      amount.Amount
}

func New(cfg *Config) (*Wallet, error) {
	wlt, err := wallet.Open(cfg.Path, false)
	if err != nil {
		return nil, err
	}

	return &Wallet{
		Wallet:   wlt,
		address:  cfg.Address,
		password: cfg.Password,
		fee:      cfg.Fee,
	}, nil
}

func (w *Wallet) BondTransaction(pubKey, toAddress, memo string, amt amount.Amount) (string, error) {
	opts := []wallet.TxOption{
		wallet.OptionFee(w.fee.ToPactusAmount()),
		wallet.OptionMemo(memo),
	}
	tx, err := w.Wallet.MakeBondTx(w.address, toAddress, pubKey, amt.ToPactusAmount(), opts...)
	if err != nil {
		log.Error("error creating bond transaction", "error", err, "to",
			toAddress, "amount", amt)

		return "", err
	}
	// sign transaction
	err = w.Wallet.SignTransaction(w.password, tx)
	if err != nil {
		log.Error("error signing bond transaction", "error", err,
			"to", toAddress, "amount", amt)

		return "", err
	}

	// broadcast transaction
	res, err := w.Wallet.BroadcastTransaction(tx)
	if err != nil {
		log.Error("error broadcasting bond transaction", "error", err,
			"to", toAddress, "amount", amt)

		return "", err
	}

	err = w.Wallet.Save()
	if err != nil {
		log.Error("error saving wallet transaction history", "error", err,
			"to", toAddress, "amount", amt)
	}

	return res, nil // return transaction hash
}

func (w *Wallet) TransferTransaction(toAddress, memo string, amt amount.Amount) (string, error) {
	opts := []wallet.TxOption{
		wallet.OptionFee(w.fee.ToPactusAmount()),
		wallet.OptionMemo(memo),
	}

	// Use amt.Amount for transaction amount.
	tx, err := w.Wallet.MakeTransferTx(w.address, toAddress, amt.ToPactusAmount(), opts...)
	if err != nil {
		log.Error("error creating transfer transaction", "error", err,
			"from", w.address, "to", toAddress, "amount", amt)

		return "", err
	}

	// sign transaction.
	err = w.Wallet.SignTransaction(w.password, tx)
	if err != nil {
		log.Error("error signing transfer transaction", "error", err,
			"to", toAddress, "amount", amt)

		return "", err
	}

	// broadcast transaction.
	res, err := w.Wallet.BroadcastTransaction(tx)
	if err != nil {
		log.Error("error broadcasting transfer transaction", "error", err,
			"to", toAddress, "amount", amt)

		return "", err
	}

	err = w.Wallet.Save()
	if err != nil {
		log.Error("error saving wallet transaction history", "error", err,
			"to", toAddress, "amount", amt)
	}

	return res, nil // return transaction hash.
}

func (w *Wallet) Address() string {
	return w.address
}

func (w *Wallet) Balance() amount.Amount {
	balance, _ := w.Wallet.Balance(w.address)

	return amount.Amount(balance.ToNanoPAC())
}

func (w *Wallet) LinkToExplorer(txID string) string {
	switch w.Network() {
	case genesis.Mainnet:
		return fmt.Sprintf("https://pacviewer.com/transaction/%s", txID)

	case genesis.Testnet:
		return fmt.Sprintf("https://phoenix.pacviewer.com/transaction/%s", txID)

	case genesis.Localnet:
		return txID

	default:
		return txID
	}
}
