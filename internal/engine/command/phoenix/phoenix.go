package phoenix

import (
	"context"
	"log"

	"github.com/pactus-project/pactus/crypto"
	"github.com/pagu-project/pagu/internal/engine/command"
	"github.com/pagu-project/pagu/internal/entity"
	"github.com/pagu-project/pagu/internal/repository"
	"github.com/pagu-project/pagu/pkg/amount"
	"github.com/pagu-project/pagu/pkg/client"
)

type PhoenixCmd struct {
	*phoenixSubCmds

	ctx           context.Context
	db            *repository.Database
	client        client.IClient
	faucetAmount  amount.Amount
	faucetAddress string
}

func NewPhoenixCmd(ctx context.Context, cfg *Config, db *repository.Database,
) *PhoenixCmd {
	client, err := client.NewClient(cfg.Client)
	if err != nil {
		log.Fatal("an error occurred on creating client for Phoenix Command", "error", err)
	}

	crypto.AddressHRP = "tpc"
	faucetAddress := cfg.PrivateKey.PublicKeyNative().AccountAddress().String()
	crypto.AddressHRP = "pc"

	// wlt, err := wallet.New(cfg.Wallet)
	// if err != nil {
	// 	log.Fatal("an error occurred on creating wallet for Phoenix Command", "error", err)
	// }

	return &PhoenixCmd{
		ctx:           ctx,
		client:        client,
		db:            db,
		faucetAddress: faucetAddress,
		faucetAmount:  cfg.FaucetAmount,
	}
}

func (p *PhoenixCmd) GetCommand() *command.Command {
	cmd := p.buildPhoenixCommand()
	cmd.PlatformIDs = entity.AllPlatformIDs()
	cmd.TargetFlag = command.TargetMaskAll

	p.subCmdFaucet.PlatformIDs = entity.AllPlatformIDs()
	p.subCmdFaucet.TargetFlag = command.TargetMaskAll

	p.subCmdWallet.PlatformIDs = entity.AllPlatformIDs()
	p.subCmdWallet.TargetFlag = command.TargetMaskAll

	return cmd
}
