package cli

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/pagu-project/pagu/internal/engine"
	"github.com/pagu-project/pagu/internal/entity"
	"github.com/spf13/cobra"
)

const PROMPT = "\n>> "

func HandleCliCommands(cmd *cobra.Command, botEngine *engine.BotEngine) {
	reader := bufio.NewReader(os.Stdin)

	for {
		chatHistory := bytes.Buffer{}
		cmd.Print(PROMPT)

		input, _ := reader.ReadString('\n')
		input = strings.TrimSuffix(input, "\n")

		if strings.EqualFold(input, "exit") {
			cmd.Println("exiting from cli")

			return
		}

		response := botEngine.ParseAndExecute(entity.PlatformIDCLI, "0", input)

		chatHistory.WriteString(fmt.Sprintf("%v\n%v", response.Title, response.Message))

		// Pass response to Glow via stdin
		command := exec.Command("glow")
		if command.Err != nil {
			cmd.Printf("%v\n%v", response.Title, response.Message)

			continue
		}
		command.Stdin = strings.NewReader(chatHistory.String())
		command.Stdout = os.Stdout
		command.Stderr = os.Stderr
		_ = command.Run()
	}
}
