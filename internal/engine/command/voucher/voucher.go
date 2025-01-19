package voucher

import (
	"github.com/pagu-project/pagu/internal/engine/command"
	"github.com/pagu-project/pagu/internal/entity"
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
	cmd.AppIDs = entity.AllAppIDs()
	cmd.TargetFlag = command.TargetMaskMainnet | command.TargetMaskModerator

	v.subCmdClaim.AppIDs = []entity.PlatformID{entity.PlatformIDDiscord}
	v.subCmdClaim.TargetFlag = command.TargetMaskMainnet
	v.subCmdClaim.Middlewares = []command.MiddlewareFunc{middlewareHandler.WalletBalance}

	v.subCmdCreate.AppIDs = []entity.PlatformID{entity.PlatformIDDiscord}
	v.subCmdCreate.TargetFlag = command.TargetMaskModerator
	v.subCmdCreate.Middlewares = []command.MiddlewareFunc{middlewareHandler.OnlyModerator}

	v.subCmdCreateBulk.AppIDs = []entity.PlatformID{entity.PlatformIDDiscord}
	v.subCmdCreateBulk.TargetFlag = command.TargetMaskModerator
	v.subCmdCreateBulk.Middlewares = []command.MiddlewareFunc{middlewareHandler.OnlyModerator}

	v.subCmdStatus.AppIDs = []entity.PlatformID{entity.PlatformIDDiscord}
	v.subCmdStatus.TargetFlag = command.TargetMaskModerator
	v.subCmdStatus.Middlewares = []command.MiddlewareFunc{middlewareHandler.OnlyModerator}

	return cmd
}
