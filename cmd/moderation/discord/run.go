package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/pagu-project/Pagu/internal/engine"

	"github.com/pagu-project/Pagu/internal/platforms/discord"
	"github.com/pagu-project/Pagu/pkg/log"

	pCmd "github.com/pagu-project/Pagu/cmd"
	"github.com/pagu-project/Pagu/config"
	"github.com/spf13/cobra"
)

func runCommand(parentCmd *cobra.Command) {
	run := &cobra.Command{
		Use:   "run",
		Short: "Runs a mainnet instance of Pagu",
	}

	parentCmd.AddCommand(run)

	run.Run = func(cmd *cobra.Command, _ []string) {
		// load configuration.
		configs, err := config.Load(configPath)
		pCmd.ExitOnError(cmd, err)

		// Initialize global logger.
		log.InitGlobalLogger(configs.Logger)

		// starting botEngine.
		botEngine, err := engine.NewBotEngine(configs)
		pCmd.ExitOnError(cmd, err)

		discordModerationBot, err := discord.NewDiscordBot(botEngine, configs.DiscordModerationBot, config.TargetMaskModerator)
		pCmd.ExitOnError(cmd, err)

		err = discordModerationBot.Start()
		pCmd.ExitOnError(cmd, err)

		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
		<-sigChan

		if err := discordModerationBot.Stop(); err != nil {
			pCmd.ExitOnError(cmd, err)
		}

		botEngine.Stop()
	}
}