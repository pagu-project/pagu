package engine

import (
	"github.com/pagu-project/pagu/internal/engine/command/phoenix"
	"github.com/pagu-project/pagu/internal/engine/command/voucher"
	"github.com/pagu-project/pagu/internal/repository"
	"github.com/pagu-project/pagu/pkg/mailer"
	"github.com/pagu-project/pagu/pkg/nowpayments"
	"github.com/pagu-project/pagu/pkg/wallet"
)

type Config struct {
	NetworkNodes []string           `yaml:"network_nodes"`
	LocalNode    string             `yaml:"local_node"`
	Mailer       mailer.Config      `yaml:"mailer"`
	NowPayments  nowpayments.Config `yaml:"now_payments"`
	Database     repository.Config  `yaml:"database"`
	Wallet       wallet.Config      `yaml:"wallet"`
	Phoenix      phoenix.Config     `yaml:"phoenix"`
	Voucher      voucher.Config     `yaml:"voucher"`
}
