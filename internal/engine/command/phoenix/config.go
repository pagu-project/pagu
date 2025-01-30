package phoenix

import (
	"github.com/pactus-project/pactus/crypto/ed25519"
	"github.com/pagu-project/pagu/pkg/amount"
)

type Config struct {
	Client       string              `yaml:"client"`
	PrivateKey   *ed25519.PrivateKey `yaml:"private_key"`
	FaucetAmount amount.Amount       `yaml:"faucet_amount"`
}
