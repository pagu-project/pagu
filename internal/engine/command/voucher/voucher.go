package voucher

import (
	"github.com/pagu-project/pagu/internal/engine/command"
	"github.com/pagu-project/pagu/internal/repository"
	"github.com/pagu-project/pagu/pkg/client"
	"github.com/pagu-project/pagu/pkg/wallet"
)

type VoucherCmd struct {
	*voucherSubCmds

	db            *repository.Database
	wallet        wallet.IWallet
	clientManager client.IManager
}

func NewVoucherCmd(db *repository.Database, wlt wallet.IWallet, cli client.IManager) *VoucherCmd {
	return &VoucherCmd{
		db:            db,
		wallet:        wlt,
		clientManager: cli,
	}
}

func (v *VoucherCmd) GetCommand() *command.Command {
	middlewareHandler := command.NewMiddlewareHandler(v.db, v.wallet)

	cmd := v.buildVoucherCommand()

	v.subCmdClaim.Middlewares = []command.MiddlewareFunc{middlewareHandler.WalletBalance}
	v.subCmdCreate.Middlewares = []command.MiddlewareFunc{middlewareHandler.OnlyModerator}
	v.subCmdCreateBulk.Middlewares = []command.MiddlewareFunc{middlewareHandler.OnlyModerator}
	v.subCmdStatus.Middlewares = []command.MiddlewareFunc{middlewareHandler.OnlyModerator}

	return cmd
}
