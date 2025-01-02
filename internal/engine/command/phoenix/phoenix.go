package phoenix

import (
	"context"
	"fmt"

	"github.com/pagu-project/pagu/internal/engine/command"
	"github.com/pagu-project/pagu/internal/entity"
	"github.com/pagu-project/pagu/internal/repository"
	"github.com/pagu-project/pagu/pkg/amount"
	"github.com/pagu-project/pagu/pkg/client"
	"github.com/pagu-project/pagu/pkg/wallet"
)

const (
	CommandName         = "phoenix"
	FaucetCommandName   = "faucet"
	WalletCommandName   = "wallet"
	StatusCommandName   = "status"
	HealthCommandName   = "health"
	NodeInfoCommandName = "node-info"
	HelpCommandName     = "help"
)

type PhoenixCmd struct {
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

func (pt *PhoenixCmd) GetCommand() *command.Command {
	middlewareHandler := command.NewMiddlewareHandler(pt.db, pt.wallet)

	subCmdStatus := &command.Command{
		Name:        StatusCommandName,
		Help:        "Phoenix Testnet statistics",
		Args:        []command.Args{},
		SubCommands: nil,
		AppIDs:      entity.AllAppIDs(),
		Middlewares: nil,
		Handler:     pt.networkStatusHandler,
		TargetFlag:  command.TargetMaskTestnet,
	}

	subCmdFaucet := &command.Command{
		Name: FaucetCommandName,
		Help: fmt.Sprintf("Get %f tPAC Coins on Phoenix Testnet for Testing your code or project", pt.faucetAmount.ToPAC()),
		Args: []command.Args{
			{
				Name:     "address",
				Desc:     "your testnet address [example: tpc1z...]",
				InputBox: command.InputBoxText,
				Optional: false,
			},
		},
		SubCommands: nil,
		AppIDs:      entity.AllAppIDs(),
		Middlewares: []command.MiddlewareFunc{middlewareHandler.WalletBalance},
		Handler:     pt.faucetHandler,
		TargetFlag:  command.TargetMaskTestnet,
	}

	cmdPhoenix := &command.Command{
		Name:        CommandName,
		Help:        "Phoenix Testnet tools and utils for developers",
		Args:        nil,
		AppIDs:      entity.AllAppIDs(),
		SubCommands: make([]*command.Command, 0),
		Handler:     nil,
		TargetFlag:  command.TargetMaskTestnet,
	}

	cmdPhoenix.AddSubCommand(subCmdFaucet)
	cmdPhoenix.AddSubCommand(subCmdStatus)

	return cmdPhoenix
}
