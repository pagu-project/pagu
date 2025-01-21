package main

import (
	pagucmd "github.com/pagu-project/pagu/cmd"
	"github.com/pagu-project/pagu/config"
	"github.com/pagu-project/pagu/internal/engine"
	"github.com/pagu-project/pagu/internal/platforms/cli"
	"github.com/pagu-project/pagu/internal/version"
	"github.com/pagu-project/pagu/pkg/log"
	"github.com/spf13/cobra"
)

var configPath string

func run(cmd *cobra.Command, _ []string) {
	configs, err := config.Load(configPath)
	pagucmd.ExitOnError(cmd, err)

	log.InitGlobalLogger(configs.Logger)

	botEngine, err := engine.NewBotEngine(configs)
	pagucmd.ExitOnError(cmd, err)

	botEngine.Start()

	cli.HandleCliCommands(cmd, botEngine)
}

func main() {
	rootCmd := &cobra.Command{
		Use:     "pagu-cli",
		Version: version.StringVersion(),
		Run:     run,
	}

	rootCmd.PersistentFlags().StringVarP(&configPath, "config", "c", "./config.yml", "config path ./config.yml")
	err := rootCmd.Execute()
	pagucmd.ExitOnError(rootCmd, err)
}
