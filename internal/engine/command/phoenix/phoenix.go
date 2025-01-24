package phoenix

import (
	"context"

	"github.com/pagu-project/pagu/internal/engine/command"
	"github.com/pagu-project/pagu/internal/entity"
	"github.com/pagu-project/pagu/internal/repository"
	"github.com/pagu-project/pagu/pkg/amount"
	"github.com/pagu-project/pagu/pkg/client"
	"github.com/pagu-project/pagu/pkg/wallet"
)

type PhoenixCmd struct {
	*phoenixSubCmds

	ctx          context.Context
	wallet       wallet.IWallet
	db           *repository.Database
	clientMgr    client.IManager
	faucetAmount amount.Amount
}

func NewPhoenixCmd(ctx context.Context, wlt wallet.IWallet, faucetAmount amount.Amount,
	clientMgr client.IManager, db *repository.Database,
) *PhoenixCmd {
	return &PhoenixCmd{
		ctx:          ctx,
		wallet:       wlt,
		clientMgr:    clientMgr,
		db:           db,
		faucetAmount: faucetAmount,
	}
}

func (p *PhoenixCmd) GetCommand() *command.Command {
	middlewareHandler := command.NewMiddlewareHandler(p.db, p.wallet)

	cmd := p.buildPhoenixCommand()
	cmd.AppIDs = entity.AllAppIDs()
	cmd.TargetFlag = command.TargetMaskTestnet

	p.subCmdFaucet.Middlewares = []command.MiddlewareFunc{middlewareHandler.WalletBalance}
	p.subCmdFaucet.AppIDs = entity.AllAppIDs()
	p.subCmdFaucet.TargetFlag = command.TargetMaskTestnet

	p.subCmdWallet.AppIDs = entity.AllAppIDs()
	p.subCmdWallet.TargetFlag = command.TargetMaskTestnet

	return cmd
}
