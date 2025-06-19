package phoenix

import (
	"context"
	"time"

	"github.com/pactus-project/pactus/crypto"
	"github.com/pactus-project/pactus/crypto/ed25519"
	"github.com/pagu-project/pagu/internal/engine/command"
	"github.com/pagu-project/pagu/internal/entity"
	"github.com/pagu-project/pagu/internal/repository"
	"github.com/pagu-project/pagu/pkg/amount"
	"github.com/pagu-project/pagu/pkg/client"
	"github.com/pagu-project/pagu/pkg/log"
	"github.com/pagu-project/pagu/pkg/utils"
)

type PhoenixCmd struct {
	*phoenixSubCmds

	ctx            context.Context
	db             *repository.Database
	client         client.IClient
	privateKey     *ed25519.PrivateKey
	faucetAddress  crypto.Address
	faucetAmount   amount.Amount
	faucetFee      amount.Amount
	faucetCooldown time.Duration
}

func NewPhoenixCmd(ctx context.Context, cfg *Config, db *repository.Database,
) *PhoenixCmd {
	client, err := client.NewClient(cfg.Client)
	if err != nil {
		log.Fatal("phoenix: bad client", "error", err)
	}

	privateKey, err := utils.TestnetPrivateKeyFromString(cfg.PrivateKey)
	if err != nil {
		log.Fatal("phoenix: invalid private key", "error", err)
	}
	faucetAddress := privateKey.PublicKeyNative().AccountAddress()

	return &PhoenixCmd{
		ctx:            ctx,
		client:         client,
		db:             db,
		privateKey:     privateKey,
		faucetAddress:  faucetAddress,
		faucetAmount:   cfg.FaucetAmount,
		faucetFee:      cfg.FaucetAmount,
		faucetCooldown: cfg.FaucetCooldown,
	}
}

func (c *PhoenixCmd) BuildCommand(botID entity.BotID) *command.Command {
	cmd := c.buildPhoenixCommand(botID)

	return cmd
}
