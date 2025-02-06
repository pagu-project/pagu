package cli

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/pagu-project/pagu/pkg/log"

	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/lipgloss"
	"github.com/pagu-project/pagu/internal/engine"
	"github.com/pagu-project/pagu/internal/entity"
	"github.com/spf13/cobra"
)

const PROMPT = "\n>> "

func HandleCliCommands(cmd *cobra.Command, botEngine *engine.BotEngine) {

	r, err := glamour.NewTermRenderer(
		glamour.WithColorProfile(lipgloss.ColorProfile()),
		glamour.TermRendererOption(glamour.WithAutoStyle()),
		glamour.WithPreservedNewLines(),
	)
	if err != nil {
		log.Error("error on rendering terminal", "error", err)
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

		mdRresponse, err := r.Render(fmt.Sprintf("%v\n%v", response.Title, response.Message))
		if err != nil {
			log.Error("error in rendering mark down", "error", err)
		}

		cmd.Printf(mdRresponse)
	}
}
