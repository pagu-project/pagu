package cli

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/pagu-project/pagu/internal/engine"
	"github.com/pagu-project/pagu/internal/entity"
	"github.com/pagu-project/pagu/pkg/markdown"
	"github.com/spf13/cobra"
)

const PROMPT = "\n>> "

func HandleCliCommands(cmd *cobra.Command, botEngine *engine.BotEngine) {
	markdown := markdown.NewCLIRenderer()
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

		res := fmt.Sprintf("%v\n%v", response.Title, response.Message)

		cmd.Print(markdown.Render(res))
	}
}
