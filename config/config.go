package config

import (
	"errors"
	"fmt"
	"os"

	"github.com/pagu-project/pagu/internal/engine"
	"github.com/pagu-project/pagu/internal/entity"
	"github.com/pagu-project/pagu/internal/platforms/discord"
	"github.com/pagu-project/pagu/internal/platforms/grpc"
	"github.com/pagu-project/pagu/internal/platforms/telegram"
	"github.com/pagu-project/pagu/internal/platforms/whatsapp"
	"github.com/pagu-project/pagu/pkg/log"
	"github.com/pagu-project/pagu/pkg/utils"
	"gopkg.in/yaml.v3"
)

type Config struct {
	BotID    entity.BotID    `yaml:"bot_id"`
	Engine   engine.Config   `yaml:"engine"`
	GRPC     grpc.Config     `yaml:"grpc"`
	Discord  discord.Config  `yaml:"discord"`
	Telegram telegram.Config `yaml:"telegram"`
	WhatsApp whatsapp.Config `yaml:"whatsapp"`
	Logger   log.Config      `yaml:"logger"`
}

func Load(path string) (*Config, error) {
	payload, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	cfg := &Config{}
	if err := yaml.Unmarshal(payload, cfg); err != nil {
		return nil, err
	}

	// Check if the required configurations are set
	if err := cfg.BasicCheck(); err != nil {
		return nil, err
	}

	return cfg, nil
}

// BasicCheck validate presence of required config variables.
func (cfg *Config) BasicCheck() error {

	if cfg.Engine.Wallet.Address == "" {
		return errors.New("config: Wallet address dose not set")
	}

	// Check if the WalletPath exists.
	if !utils.PathExists(cfg.Engine.Wallet.Path) {
		return fmt.Errorf("config: Wallet does not exist: %s", cfg.Engine.Wallet.Path)
	}

	if len(cfg.Engine.NetworkNodes) == 0 {
		return errors.New("config: network nodes is empty")
	}

	return nil
}
