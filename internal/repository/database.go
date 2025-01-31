package repository

import (
	"errors"
	"strings"

	"github.com/pagu-project/pagu/internal/entity"
	"gorm.io/driver/mysql"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Database struct {
	gormDB *gorm.DB
}

func NewDB(path string) (*Database, error) {
	parts := strings.SplitN(path, ":", 2)
	if len(parts) != 2 {
		return nil, errors.New("invalid database URL format; expected format is 'dbtype:connection_string'")
	}

	dbType, connStr := parts[0], parts[1]

	var db *gorm.DB
	var err error
	conf := &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	}

	switch dbType {
	case "mysql":
		db, err = gorm.Open(mysql.Open(connStr), conf)
	case "sqlite":
		db, err = gorm.Open(sqlite.Open(connStr), conf)
	default:
		return nil, errors.New("unsupported database type: " + dbType)
	}

	if err != nil {
		return nil, ConnectionError{
			Message: err.Error(),
		}
	}

	if !db.Migrator().HasTable(&entity.User{}) ||
		!db.Migrator().HasTable(&entity.PhoenixFaucet{}) ||
		!db.Migrator().HasTable(&entity.Voucher{}) ||
		!db.Migrator().HasTable(&entity.ZealyUser{}) ||
		!db.Migrator().HasTable(&entity.Notification{}) ||
		!db.Migrator().HasTable(&entity.CrowdfundCampaign{}) ||
		!db.Migrator().HasTable(&entity.CrowdfundPurchase{}) {
		if err := db.AutoMigrate(
			&entity.User{},
			&entity.PhoenixFaucet{},
			&entity.ZealyUser{},
			&entity.Voucher{},
			&entity.Notification{},
			&entity.CrowdfundCampaign{},
			&entity.CrowdfundPurchase{},
		); err != nil {
			return nil, MigrationError{
				Message: err.Error(),
			}
		}
	}

	return &Database{
		gormDB: db,
	}, nil
}
