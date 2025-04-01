package main

import (
	"os"
	"os/signal"
	"syscall"

	pagucmd "github.com/pagu-project/pagu/cmd"
	"github.com/pagu-project/pagu/config"
	"github.com/pagu-project/pagu/internal/engine"
	whatsapp "github.com/pagu-project/pagu/internal/platforms/whatsup"
	"github.com/pagu-project/pagu/pkg/log"
	"github.com/spf13/cobra"
)

func runCommand(parentCmd *cobra.Command) {
	run := &cobra.Command{
		Use:   "run",
		Short: "Runs a mainnet instance of Pagu",
	}

	parentCmd.AddCommand(run)

	run.Run = func(cmd *cobra.Command, _ []string) {
		// Load configuration.
		configs, err := config.Load(configPath)
		pagucmd.ExitOnError(cmd, err)

		// Starting eng.
		eng, err := engine.NewBotEngine(configs)
		pagucmd.ExitOnError(cmd, err)

		log.InitGlobalLogger(configs.Logger)

		eng.Start()

		bot, err := whatsapp.NewWhatsUpBot(eng, configs)
		pagucmd.ExitOnError(cmd, err)

		err = bot.Start()
		pagucmd.ExitOnError(cmd, err)

		// Set up signal handling.
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
		go func() {
			<-c
			// When a signal is received, stop the bot and perform any other necessary cleanup.
			bot.Stop()
			eng.Stop()
			os.Exit(0)
		}()

		// Block the main goroutine until a signal is received.
		select {}
	}
}
