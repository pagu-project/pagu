package repository

import (
	"errors"
	"time"

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

func (db *Database) GetVoucherByCode(code string) (entity.Voucher, error) {
	var voucher entity.Voucher
	err := db.gormDB.Model(&entity.Voucher{}).Where("code = ?", code).First(&voucher).Error
	if err != nil {
		return entity.Voucher{}, err
	}

	return voucher, nil
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

// GetNonExpiredVoucherByEmail returns a non-expired voucher for the given email, or nil if none exists.
func (db *Database) GetNonExpiredVoucherByEmail(email string) *entity.Voucher {
	var voucher entity.Voucher
	err := db.gormDB.Model(&entity.Voucher{}).
		Where("email = ?", email).
		Order("created_at DESC").
		First(&voucher).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil
		}
		log.Warn("failed to fetch non-expired voucher by email '%s': %v", email, err)

		return nil
	}

	expireAt := voucher.CreatedAt.AddDate(0, int(voucher.ValidMonths), 0)
	if expireAt.After(time.Now()) {
		return &voucher
	}

	return nil
}
