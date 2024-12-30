package wallet

import (
	"github.com/pactus-project/pactus/types/tx/payload"
	"github.com/pactus-project/pactus/wallet"
	"github.com/pagu-project/pagu/config"
	"github.com/pagu-project/pagu/pkg/amount"
	"github.com/pagu-project/pagu/pkg/log"
)

type Wallet struct {
	*wallet.Wallet

	address  string
	password string
}

func Open(cfg *config.Wallet) (*Wallet, error) {
	wlt, err := wallet.Open(cfg.Path, false)
	if err != nil {
		return nil, err
	}

	return &Wallet{
		Wallet:   wlt,
		address:  cfg.Address,
		password: cfg.Password,
	}, nil
}

func (w *Wallet) BondTransaction(pubKey, toAddress, memo string, amt amount.Amount) (string, error) {
	opts := []wallet.TxOption{
		wallet.OptionMemo(memo),
	}
	tx, err := w.Wallet.MakeBondTx(w.address, toAddress, pubKey, amt.ToPactusAmount(), opts...)
	if err != nil {
		log.Error("error creating bond transaction", "err", err, "to",
			toAddress, "amount", amt)

		return "", err
	}
	// sign transaction
	err = w.Wallet.SignTransaction(w.password, tx)
	if err != nil {
		log.Error("error signing bond transaction", "err", err,
			"to", toAddress, "amount", amt)

		return "", err
	}

	// broadcast transaction
	res, err := w.Wallet.BroadcastTransaction(tx)
	if err != nil {
		log.Error("error broadcasting bond transaction", "err", err,
			"to", toAddress, "amount", amt)

		return "", err
	}

	err = w.Wallet.Save()
	if err != nil {
		log.Error("error saving wallet transaction history", "err", err,
			"to", toAddress, "amount", amt)
	}

	return res, nil // return transaction hash
}

func (w *Wallet) TransferTransaction(toAddress, memo string, amt amount.Amount) (string, error) {
	// calculate fee using amount struct.
	fee, err := w.Wallet.CalculateFee(amt.ToPactusAmount(), payload.TypeTransfer)
	if err != nil {
		log.Error("error calculating fee", "err", err, "client")

		return "", err
	}

	opts := []wallet.TxOption{
		wallet.OptionFee(fee),
		wallet.OptionMemo(memo),
	}

	// Use amt.Amount for transaction amount.
	tx, err := w.Wallet.MakeTransferTx(w.address, toAddress, amt.ToPactusAmount(), opts...)
	if err != nil {
		log.Error("error creating transfer transaction", "err", err,
			"from", w.address, "to", toAddress, "amount", amt)

		return "", err
	}

	// sign transaction.
	err = w.Wallet.SignTransaction(w.password, tx)
	if err != nil {
		log.Error("error signing transfer transaction", "err", err,
			"to", toAddress, "amount", amt)

		return "", err
	}

	// broadcast transaction.
	res, err := w.Wallet.BroadcastTransaction(tx)
	if err != nil {
		log.Error("error broadcasting transfer transaction", "err", err,
			"to", toAddress, "amount", amt)

		return "", err
	}

	err = w.Wallet.Save()
	if err != nil {
		log.Error("error saving wallet transaction history", "err", err,
			"to", toAddress, "amount", amt)
	}

	return res, nil // return transaction hash.
}

func (w *Wallet) Address() string {
	return w.address
}

func (w *Wallet) Balance() int64 {
	balance, _ := w.Wallet.Balance(w.address)

	return int64(balance)
}
