package main

import (
	"github.com/pagu-project/pagu"
	"github.com/pagu-project/pagu/cmd"
	"github.com/spf13/cobra"
)

var configPath string

func main() {
	rootCmd := &cobra.Command{
		Use:     "pagu-telegram",
		Version: pagu.StringVersion(),
	}

	rootCmd.PersistentFlags().StringVarP(&configPath, "config", "c", "./config.yml", "config path ./config.yml")
	runCommand(rootCmd)
	err := rootCmd.Execute()
	cmd.ExitOnError(rootCmd, err)
}
