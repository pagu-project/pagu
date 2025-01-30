package wallet

import "github.com/pagu-project/pagu/pkg/amount"

type Config struct {
	Address  string        `yaml:"address"`
	Path     string        `yaml:"path"`
	Password string        `yaml:"password"`
	Fee      amount.Amount `yaml:"fee"`
	Network  string        `yaml:"network"`
}
