package cli

import (
	"bufio"
	"os"
	"strings"

	"github.com/pagu-project/pagu/internal/engine"
	"github.com/pagu-project/pagu/internal/entity"
	"github.com/spf13/cobra"
)

const PROMPT = "\n>> "

func HandleCliCommands(cmd *cobra.Command, botEngine *engine.BotEngine) {
	reader := bufio.NewReader(os.Stdin)

	for {
		cmd.Print(PROMPT)

		input, _ := reader.ReadString('\n')
		input = strings.TrimSuffix(input, "\n")

		if strings.EqualFold(input, "exit") {
			cmd.Println("exiting from cli")

			return
		}

		response := botEngine.ParseAndExecute(entity.PlatformIDCLI, "0", input)

		cmd.Printf("%v\n%v", response.Title, response.Message)
	}
}
