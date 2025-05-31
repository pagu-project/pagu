package cli

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/pagu-project/pagu/internal/engine"
	"github.com/pagu-project/pagu/internal/entity"
	"github.com/pagu-project/pagu/pkg/markdown"
	"github.com/spf13/cobra"
)

const PROMPT = "\n>> "

type Bot struct {
	ctx    context.Context
	cmd    *cobra.Command
	engine *engine.BotEngine
}

func NewCLIBot(ctx context.Context, cmd *cobra.Command, engine *engine.BotEngine) (*Bot, error) {
	return &Bot{
		ctx:    ctx,
		cmd:    cmd,
		engine: engine,
	}, nil
}

func (bot *Bot) Start() error {
	markdown := markdown.NewCLIRenderer()
	reader := bufio.NewReader(os.Stdin)

	for {
		bot.cmd.Print(PROMPT)

		input, _ := reader.ReadString('\n')
		input = strings.TrimSuffix(input, "\n")

		if strings.EqualFold(input, "exit") {
			bot.cmd.Println("exiting from cli")

			return nil
		}

		response := bot.engine.ParseAndExecute(entity.PlatformIDCLI, "0", input)

		res := fmt.Sprintf("%v\n%v", response.Title, response.Message)

		bot.cmd.Print(markdown.Render(res))
	}
}

func (bot *Bot) Stop() {}
