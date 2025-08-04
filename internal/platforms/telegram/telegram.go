package telegram

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/pagu-project/pagu/internal/engine"
	"github.com/pagu-project/pagu/internal/engine/command"
	"github.com/pagu-project/pagu/internal/entity"
	"github.com/pagu-project/pagu/pkg/log"
	"github.com/pagu-project/pagu/pkg/markdown"
	tele "gopkg.in/telebot.v4"
)

type Bot struct {
	ctx      context.Context
	engine   *engine.BotEngine
	teleBot  *tele.Bot
	cfg      *Config
	botID    entity.BotID
	markdown markdown.Renderer
}

type BotContext struct {
	Commands []string
}

var (
	argsContext = make(map[int64]*BotContext)
	argsValue   = make(map[int64]map[string]string)
)

func NewTelegramBot(ctx context.Context, cfg *Config, botID entity.BotID, engine *engine.BotEngine) (*Bot, error) {
	pref := tele.Settings{
		Token:     cfg.BotToken,
		ParseMode: tele.ModeMarkdownV2,
		Poller:    &tele.LongPoller{Timeout: 10 * time.Second},
	}

	teleBot, err := tele.NewBot(pref)
	if err != nil {
		log.Error("Failed to create Telegram bot", "error", err)

		return nil, err
	}

	markdown := markdown.NewTelegramRenderer()

	return &Bot{
		engine:   engine,
		teleBot:  teleBot,
		cfg:      cfg,
		ctx:      ctx,
		botID:    botID,
		markdown: markdown,
	}, nil
}

func (bot *Bot) Start() error {
	bot.deleteAllCommands()
	if err := bot.registerCommands(); err != nil {
		return err
	}

	go bot.teleBot.Start()
	log.Info("Starting Telegram Bot...")

	return nil
}

func (bot *Bot) Stop() {
	log.Info("Shutting down Telegram Bot")
	bot.teleBot.Stop()
}

func (bot *Bot) deleteAllCommands() {
	for _, cmd := range bot.engine.Commands() {
		err := bot.teleBot.DeleteCommands(tele.Command{Text: cmd.Name})
		if err != nil {
			log.Error("unable to delete command", "error", err, "cmd", cmd.Name)
		}
	}
}

func (bot *Bot) registerCommands() error {
	rows := make([]tele.Row, 0)
	menu := &tele.ReplyMarkup{ResizeKeyboard: true}
	commands := make([]tele.Command, 0)

	cmds := bot.engine.Commands()
	for _, cmd := range cmds {
		log.Debug("registering new command", "name", cmd.Name)

		btn := menu.Data(cmd.Name, cmd.Name)
		commands = append(commands, tele.Command{Text: cmd.Name, Description: cmd.Help})
		rows = append(rows, menu.Row(btn))

		if cmd.HasSubCommand() {
			subMenu := &tele.ReplyMarkup{ResizeKeyboard: true}
			subRows := make([]tele.Row, 0)
			for _, subCmd := range cmd.SubCommands {
				log.Debug("adding sub-command", "command", cmd.Name, "sub-command", subCmd.Name)

				subBtn := subMenu.Data(subCmd.Name, cmd.Name+subCmd.Name)
				subRows = append(subRows, subMenu.Row(subBtn))

				bot.teleBot.Handle(&subBtn, func(tgCtx tele.Context) error {
					if len(subCmd.Args) > 0 {
						return bot.handleArgCommand(tgCtx, []string{cmd.Name, subCmd.Name}, subCmd.Args)
					}

					return bot.handleCommand(tgCtx, []string{cmd.Name, subCmd.Name})
				})
			}

			subMenu.Inline(subRows...)
			bot.teleBot.Handle(&btn, func(tgCtx tele.Context) error {
				_ = bot.teleBot.Delete(tgCtx.Message())

				return bot.sendMarkdown(tgCtx, cmd.Name, subMenu)
			})

			bot.teleBot.Handle(fmt.Sprintf("/%s", cmd.Name), func(tgCtx tele.Context) error {
				_ = bot.teleBot.Delete(tgCtx.Message())

				return bot.sendMarkdown(tgCtx, cmd.Help, subMenu)
			})
		} else {
			bot.teleBot.Handle(fmt.Sprintf("/%s", cmd.Name), func(tgCtx tele.Context) error {
				_ = bot.teleBot.Delete(tgCtx.Message())

				err := bot.handleCommand(tgCtx, []string{cmd.Name})

				return err
			})
		}
	}

	// initiate menu button
	_ = bot.teleBot.SetCommands(commands)
	menu.Inline(rows...)
	bot.teleBot.Handle("/start", func(tgCtx tele.Context) error {
		_ = bot.teleBot.Delete(tgCtx.Message())

		return bot.sendMarkdown(tgCtx, "Pagu Main Menu", menu)
	})

	bot.teleBot.Handle(tele.OnText, func(tgCtx tele.Context) error {
		if argsContext[tgCtx.Message().Sender.ID] == nil {
			return nil
		}

		if argsValue[tgCtx.Message().Sender.ID] == nil {
			argsValue[tgCtx.Message().Sender.ID] = make(map[string]string)
		}

		return bot.parsTextMessage(tgCtx)
	})

	return nil
}

func (bot *Bot) parsTextMessage(tgCtx tele.Context) error {
	senderID := tgCtx.Message().Sender.ID
	path := argsContext[senderID].Commands
	cmd, err := bot.engine.FindCommandByPath(path)
	if err != nil {
		return bot.sendMarkdown(tgCtx, err.Error())
	}

	currentArgsIndex := len(argsValue[senderID])
	argsValue[senderID][cmd.Args[currentArgsIndex].Name] = tgCtx.Message().Text

	if len(argsValue[senderID]) == len(cmd.Args) {
		return bot.handleCommand(tgCtx, argsContext[senderID].Commands)
	}

	_ = bot.teleBot.Delete(tgCtx.Message())

	return bot.sendMarkdown(tgCtx, fmt.Sprintf("Enter `%s`:", cmd.Args[currentArgsIndex+1].Name))
}

func (bot *Bot) handleArgCommand(tgCtx tele.Context, commands []string, args []*command.Args) error {
	msgCtx := &BotContext{Commands: commands}
	argsContext[tgCtx.Sender().ID] = msgCtx
	argsValue[tgCtx.Sender().ID] = nil
	_ = bot.teleBot.Delete(tgCtx.Message())

	firstArg := args[0]

	if len(firstArg.Choices) > 0 {
		choiceMsg := fmt.Sprintf("Select a `%s`:\n\n", firstArg.Name)
		choiceMenu := &tele.ReplyMarkup{ResizeKeyboard: true}
		choiceRows := make([]tele.Row, 0, len(firstArg.Choices))
		for _, choice := range firstArg.Choices {
			choiceMsg += fmt.Sprintf("- %s\n", choice.Desc)
			choiceBtn := choiceMenu.Data(choice.Name, firstArg.Name, choice.Value)
			choiceRows = append(choiceRows, choiceMenu.Row(choiceBtn))

			bot.teleBot.Handle(&choiceBtn, func(tgCtx tele.Context) error {
				commands = append(commands, fmt.Sprintf("--%s=%v", firstArg.Name, choice.Value))

				return bot.handleCommand(tgCtx, commands)
			})
		}

		choiceMenu.Inline(choiceRows...)

		return bot.sendMarkdown(tgCtx, choiceMsg, choiceMenu)
	}

	// Commands with no choices
	return bot.sendMarkdown(tgCtx, fmt.Sprintf("Enter `%s`:", firstArg.Name))
}

// handleCommand executes a command with its arguments for the user.
// It combines the commands and arguments into a single line and execute the command line.
// clears the user's context, and sends the result back to the user.
func (bot *Bot) handleCommand(tgCtx tele.Context, commands []string) error {
	callerID := strconv.Itoa(int(tgCtx.Sender().ID))

	// Retrieve the arguments for the sender
	senderID := tgCtx.Message().Sender.ID
	args := argsValue[senderID]

	// Combine the commands into a single string
	fullCommand := strings.Join(commands, " ")

	// Append arguments as key-value pairs
	if len(args) > 0 {
		argPairs := []string{}
		for key, value := range args {
			argPairs = append(argPairs, fmt.Sprintf("--%s=%s", key, value))
		}
		fullCommand = fmt.Sprintf("%s %s", fullCommand, strings.Join(argPairs, " "))
	}

	// Call the engine's Run method with the full command string
	res := bot.engine.ParseAndExecute(entity.PlatformIDTelegram, callerID, fullCommand)
	_ = bot.teleBot.Delete(tgCtx.Message())

	// Clear the stored command context and arguments for the sender
	argsContext[senderID] = nil
	argsValue[senderID] = nil

	return bot.sendMarkdown(tgCtx, res.Message, tele.NoPreview)
}

func (bot *Bot) sendMarkdown(tgCtx tele.Context, what string, opts ...any) error {
	rendered := bot.markdown.Render(what)
	opts = append(opts, tele.ModeMarkdownV2)

	return tgCtx.Send(rendered, opts...)
}
