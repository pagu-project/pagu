package repository

import (
	"github.com/pagu-project/pagu/internal/entity"
)

type IVoucher interface {
	AddVoucher(v *entity.Voucher) error
	GetVoucherByCode(code string) (entity.Voucher, error)
	ClaimVoucher(id uint, txHash string, claimer uint) error
	ListVoucher() ([]*entity.Voucher, error)
}

func (db *Database) AddVoucher(v *entity.Voucher) error {
	tx := db.Create(v)
	if tx.Error != nil {
		return WriteError{
			Message: tx.Error.Error(),
		}
	}

	return nil
}

func (db *Database) GetVoucherByCode(code string) (entity.Voucher, error) {
	var voucher entity.Voucher
	err := db.Model(&entity.Voucher{}).Where("code = ?", code).First(&voucher).Error
	if err != nil {
		return entity.Voucher{}, err
	}

	return voucher, nil
}

func (db *Database) ClaimVoucher(id uint, txHash string, claimer uint) error {
	tx := db.Model(&entity.Voucher{}).Where("id = ?", id).Update("tx_hash", txHash).Update("claimed_by", claimer)
	if tx.Error != nil {
		return WriteError{
			Message: tx.Error.Error(),
		}
	}

	return nil
}

func (db *Database) ListVoucher() ([]*entity.Voucher, error) {
	var vouchers []*entity.Voucher
	tx := db.Find(&vouchers)
	if tx.Error != nil {
		return nil, ReadError{
			Message: tx.Error.Error(),
		}
	}

	return vouchers, nil
}
