package voucher

import (
	"github.com/pagu-project/pagu/internal/engine/command"
	"github.com/pagu-project/pagu/internal/repository"
	"github.com/pagu-project/pagu/pkg/client"
	"github.com/pagu-project/pagu/pkg/mailer"
	"github.com/pagu-project/pagu/pkg/wallet"
)

type VoucherCmd struct {
	*voucherSubCmds

	db            *repository.Database
	wallet        wallet.IWallet
	clientManager client.IManager
	mailer        mailer.IMailer
}

func NewVoucherCmd(db *repository.Database, wlt wallet.IWallet,
	clientManager client.IManager, mailer mailer.IMailer,
) *VoucherCmd {
	return &VoucherCmd{
		db:            db,
		wallet:        wlt,
		clientManager: clientManager,
		mailer:        mailer,
	}
}

func (v *VoucherCmd) GetCommand() *command.Command {
	cmd := v.buildVoucherCommand()

	return cmd
}
