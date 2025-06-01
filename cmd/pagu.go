package main

import (
	"context"

	"github.com/pagu-project/pagu/config"
	"github.com/pagu-project/pagu/internal/engine"
	"github.com/pagu-project/pagu/internal/entity"
	"github.com/pagu-project/pagu/internal/platforms/cli"
	"github.com/pagu-project/pagu/internal/platforms/discord"
	"github.com/pagu-project/pagu/internal/platforms/telegram"
	"github.com/pagu-project/pagu/internal/platforms/whatsapp"
	"github.com/pagu-project/pagu/internal/version"
	"github.com/pagu-project/pagu/pkg/log"
	"github.com/spf13/cobra"
)

type BotInstance interface {
	Start() error
	Stop()
}

var configPath string

func main() {
	rootCmd := &cobra.Command{
		Use:     "pagu",
		Version: version.StringVersion(),
	}
	rootCmd.CompletionOptions.DisableDefaultCmd = true

	rootCmd.PersistentFlags().StringVarP(&configPath, "config", "c", "", "path to the config file")
	err := rootCmd.MarkPersistentFlagRequired("config")
	ExitOnError(rootCmd, err)

	runCommand(rootCmd)
	err = rootCmd.Execute()
	ExitOnError(rootCmd, err)
}

func runCommand(parentCmd *cobra.Command) {
	run := &cobra.Command{
		Use:   "run",
		Short: "Run an instance of Pagu",
	}

	parentCmd.AddCommand(run)

	run.Run = func(cmd *cobra.Command, _ []string) {
		ctx, cancel := context.WithCancel(context.Background())
		// load configuration.
		cfg, err := config.Load(configPath)
		ExitOnError(cmd, err)

		// Initialize global logger.
		log.InitGlobalLogger(&cfg.Logger)

		// starting eng.
		eng, err := engine.NewBotEngine(ctx, &cfg.Engine)
		ExitOnError(cmd, err)

		eng.Start()

		var bot BotInstance
		switch cfg.BotID {
		case entity.BotID_CLI:
			bot, err = cli.NewCLIBot(ctx, cmd, eng)
			ExitOnError(cmd, err)

		case entity.BotID_Discord:
			bot, err = discord.NewDiscordBot(ctx, &cfg.Discord, cfg.BotID, eng)
			ExitOnError(cmd, err)

		case entity.BotID_Moderator:
			bot, err = discord.NewDiscordBot(ctx, &cfg.Discord, cfg.BotID, eng)
			ExitOnError(cmd, err)

		case entity.BotID_Telegram:
			bot, err = telegram.NewTelegramBot(ctx, &cfg.Telegram, cfg.BotID, eng)
			ExitOnError(cmd, err)

		case entity.BotID_WhatsApp:
			bot, err = whatsapp.NewWhatsAppBot(ctx, &cfg.WhatsApp, cfg.BotID, eng)
			ExitOnError(cmd, err)

		case entity.BotID_Web:
			// TODO: implement me
		}

		err = bot.Start()
		ExitOnError(cmd, err)

		TrapSignal(func() {
			cancel()

			bot.Stop()
			eng.Stop()
		})

		// run forever
		select {}
	}
}
