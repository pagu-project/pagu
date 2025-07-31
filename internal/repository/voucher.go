package repository

import (
	"errors"

	"github.com/pagu-project/pagu/internal/entity"
	"github.com/pagu-project/pagu/pkg/log"
	"gorm.io/gorm"
)

func (db *Database) AddVoucher(v *entity.Voucher) error {
	tx := db.gormDB.Create(v)
	if tx.Error != nil {
		return WriteError{
			Message: tx.Error.Error(),
		}
	}

	return nil
}

func (db *Database) GetVoucherByCode(code string) (*entity.Voucher, error) {
	var voucher entity.Voucher
	err := db.gormDB.Model(&entity.Voucher{}).Where("code = ?", code).First(&voucher).Error
	if err != nil {
		return nil, err
	}

	return &voucher, nil
}

func (db *Database) GetVoucherByEmail(email string) (*entity.Voucher, error) {
	// TODO: maybe more than 1?
	var voucher entity.Voucher
	err := db.gormDB.Model(&entity.Voucher{}).Where("email = ?", email).First(&voucher).Error
	if err != nil {
		return nil, err
	}

	return &voucher, nil
}

func (db *Database) ClaimVoucher(id uint, txHash string, claimer uint) error {
	tx := db.gormDB.Model(&entity.Voucher{}).Where("id = ?", id).Update("tx_hash", txHash).Update("claimed_by", claimer)
	if tx.Error != nil {
		return WriteError{
			Message: tx.Error.Error(),
		}
	}

	return nil
}

func (db *Database) ListVoucher() ([]*entity.Voucher, error) {
	var vouchers []*entity.Voucher
	tx := db.gormDB.Find(&vouchers)
	if tx.Error != nil {
		return nil, ReadError{
			Message: tx.Error.Error(),
		}
	}

	return vouchers, nil
}

// IsDuplicatedVoucher checks if the voucher is created with the same email, amount, recipient, and description.
func (db *Database) IsDuplicatedVoucher(voucher *entity.Voucher) bool {
	err := db.gormDB.Model(&entity.Voucher{}).
		Where("email = ? AND amount = ? AND recipient = ? AND desc = ?",
			voucher.Email, voucher.Amount, voucher.Recipient, voucher.Desc).
		Order("created_at DESC").
		First(&voucher).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false
		}
		log.Warn("failed to check duplicated voucher", "email", voucher.Email, "error", err)

		return true
	}

	return true
}
