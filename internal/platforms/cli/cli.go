package cli

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/lipgloss"
	"github.com/pagu-project/pagu/internal/engine"
	"github.com/pagu-project/pagu/internal/entity"
	"github.com/pagu-project/pagu/pkg/log"
	"github.com/spf13/cobra"
)

const PROMPT = "\n>> "

func HandleCliCommands(cmd *cobra.Command, botEngine *engine.BotEngine) {
	r, err := glamour.NewTermRenderer(
		glamour.WithColorProfile(lipgloss.ColorProfile()),
		glamour.WithAutoStyle(),
		glamour.WithPreservedNewLines(),
	)
	if err != nil {
		log.Warn("error on rendering terminal", "Warn", err)
	}

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

		if r != nil {
			richRresponse, err := r.Render(res)
			if err != nil {
				log.Warn("error in rendering mark down", "Warn", err)
			} else {
				res = richRresponse
			}
		}

		cmd.Print(res)
	}
}
