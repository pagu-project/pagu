package main

import (
	"bufio"
	"os"
	"strings"

	"github.com/pactus-project/pactus/crypto"
	"github.com/pagu-project/pagu"
	pagucmd "github.com/pagu-project/pagu/cmd"
	"github.com/pagu-project/pagu/config"
	"github.com/pagu-project/pagu/internal/engine"
	"github.com/pagu-project/pagu/internal/entity"
	"github.com/pagu-project/pagu/pkg/log"
	"github.com/spf13/cobra"
)

var configPath string

const PROMPT = "\n>> "

func run(cmd *cobra.Command, _ []string) {
	configs, err := config.Load(configPath)
	pagucmd.ExitOnError(cmd, err)

	log.InitGlobalLogger(configs.Logger)

	if configs.Network == "Localnet" {
		crypto.AddressHRP = "tpc"
	}

	botEngine, err := engine.NewBotEngine(configs)
	pagucmd.ExitOnError(cmd, err)

	botEngine.RegisterAllCommands()
	botEngine.Start()

	reader := bufio.NewReader(os.Stdin)

	for {
		cmd.Print(PROMPT)

		input, _ := reader.ReadString('\n')
		input = strings.TrimSuffix(input, "\n")

		if strings.EqualFold(input, "exit") {
			cmd.Println("exiting from cli")

			return
		}

		response := botEngine.ParseAndExecute(entity.AppIDCLI, "0", input)

		cmd.Printf("%v\n%v", response.Title, response.Message)
	}
}

func main() {
	rootCmd := &cobra.Command{
		Use:     "pagu-cli",
		Version: pagu.StringVersion(),
		Run:     run,
	}

	rootCmd.PersistentFlags().StringVarP(&configPath, "config", "c", "./config.yml", "config path ./config.yml")
	err := rootCmd.Execute()
	pagucmd.ExitOnError(rootCmd, err)
}
