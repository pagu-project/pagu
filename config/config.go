package config

import (
	"errors"
	"fmt"
	"os"

	"github.com/pagu-project/pagu/internal/engine/command/phoenix"
	"github.com/pagu-project/pagu/pkg/log"
	"github.com/pagu-project/pagu/pkg/nowpayments"
	"github.com/pagu-project/pagu/pkg/utils"
	"github.com/pagu-project/pagu/pkg/wallet"
	"gopkg.in/yaml.v3"
)

type Config struct {
	BotName      string              `yaml:"bot_name"`
	NetworkNodes []string            `yaml:"network_nodes"`
	LocalNode    string              `yaml:"local_node"`
	Database     Database            `yaml:"database"`
	GRPC         *GRPC               `yaml:"grpc"` // ! TODO: config for modules should moved to the module.
	Wallet       *wallet.Config      `yaml:"wallet"`
	Logger       *log.Config         `yaml:"logger"`
	HTTP         *HTTP               `yaml:"http"`
	Phoenix      *phoenix.Config     `yaml:"phoenix"`
	Discord      *DiscordBot         `yaml:"discord"`
	Telegram     *Telegram           `yaml:"telegram"`
	WhatsApp     *WhatsApp           `yaml:"whatsapp"`
	Notification *Notification       `yaml:"notification"`
	NowPayments  *nowpayments.Config `yaml:"now_payments"`
}

type Database struct {
	URL string `yaml:"url"`
}

type DiscordBot struct {
	Token   string `yaml:"token"`
	GuildID string `yaml:"guild_id"`
}

type GRPC struct {
	Listen string `yaml:"listen"`
}

type HTTP struct {
	Listen string `yaml:"listen"`
}

type Telegram struct {
	BotToken string `yaml:"bot_token"`
}

type WhatsApp struct {
	WebHookToken string `yaml:"web_hook_token"`
	GraphToken   string `yaml:"graph_token"`
	Port         int    `yaml:"port"`
}

type Notification struct {
	Zoho *Zoho `yaml:"zoho"`
}

type Zoho struct {
	Mail ZapToMail `yaml:"mail"`
}

type ZapToMail struct {
	Host      string            `yaml:"host"`
	Port      int               `yaml:"port"`
	Username  string            `yaml:"username"`
	Password  string            `yaml:"password"`
	Templates map[string]string `yaml:"templates"`
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
	if cfg.Wallet.Address == "" {
		return errors.New("config: Wallet address dose not set")
	}

	// Check if the WalletPath exists.
	if !utils.PathExists(cfg.Wallet.Path) {
		return fmt.Errorf("config: Wallet does not exist: %s", cfg.Wallet.Path)
	}

	if len(cfg.NetworkNodes) == 0 {
		return errors.New("config: network nodes is empty")
	}

	return nil
}
