package phoenix

import (
	"time"

	"github.com/pagu-project/pagu/pkg/amount"
)

type Config struct {
	Client         string        `yaml:"client"`
	PrivateKey     string        `yaml:"private_key"`
	FaucetAmount   amount.Amount `yaml:"faucet_amount"`
	FaucetFee      amount.Amount `yaml:"faucet_fee"`
	FaucetCooldown time.Duration `yaml:"faucet_cooldown"`
}
