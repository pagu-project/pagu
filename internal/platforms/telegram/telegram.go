package telegram

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/pagu-project/pagu/config"
	"github.com/pagu-project/pagu/internal/engine"
	"github.com/pagu-project/pagu/internal/engine/command"
	"github.com/pagu-project/pagu/internal/entity"
	"github.com/pagu-project/pagu/pkg/log"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	tele "gopkg.in/telebot.v3"
)

type Bot struct {
	ctx         context.Context
	cancel      context.CancelFunc
	engine      *engine.BotEngine
	botInstance *tele.Bot
	cfg         *config.Config
	target      string
}

type BotContext struct {
	Commands []string
}

var (
	argsContext = make(map[int64]*BotContext)
	argsValue   = make(map[int64]map[string]string)
)

func NewTelegramBot(botEngine *engine.BotEngine, token string, cfg *config.Config) (*Bot, error) {
	pref := tele.Settings{
		Token:  token,
		Poller: &tele.LongPoller{Timeout: 10 * time.Second},
	}

	tgb, err := tele.NewBot(pref)
	if err != nil {
		log.Error("Failed to create Telegram bot:", err)

		return nil, err
	}

	ctx, cancel := context.WithCancel(context.Background())

	return &Bot{
		engine:      botEngine,
		botInstance: tgb,
		cfg:         cfg,
		ctx:         ctx,
		cancel:      cancel,
		target:      cfg.BotName,
	}, nil
}

func (bot *Bot) Start() error {
	bot.deleteAllCommands()
	if err := bot.registerCommands(); err != nil {
		return err
	}

	go bot.botInstance.Start()
	log.Info("Starting Telegram Bot...")

	return nil
}

func (bot *Bot) Stop() {
	log.Info("Shutting down Telegram Bot")
	bot.cancel()
	bot.botInstance.Stop()
}

func (bot *Bot) deleteAllCommands() {
	for _, cmd := range bot.engine.Commands() {
		err := bot.botInstance.DeleteCommands(tele.Command{Text: cmd.Name})
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
	for i, cmd := range cmds {
		if !cmd.HasBotID(entity.BotID_Telegram) {
			continue
		}

		log.Info("registering new command", "name", cmd.Name, "desc", cmd.Help, "index", i, "object", cmd)

		btn := menu.Data(cases.Title(language.English).String(cmd.Name), cmd.Name)
		commands = append(commands, tele.Command{Text: cmd.Name, Description: cmd.Help})
		rows = append(rows, menu.Row(btn))
		if cmd.HasSubCommand() {
			subMenu := &tele.ReplyMarkup{ResizeKeyboard: true}
			subRows := make([]tele.Row, 0)
			for _, subCmd := range cmd.SubCommands {
				log.Info("adding command sub-command", "command", cmd.Name,
					"sub-command", subCmd.Name, "desc", subCmd.Help)

				subBtn := subMenu.Data(cases.Title(language.English).String(subCmd.Name), subCmd.Name)

				bot.botInstance.Handle(&subBtn, func(c tele.Context) error {
					if len(subCmd.Args) > 0 {
						return bot.handleArgCommand(c, []string{cmd.Name, subCmd.Name}, subCmd.Args)
					}

					return bot.handleCommand(c, []string{cmd.Name, subCmd.Name})
				})
				subRows = append(subRows, subMenu.Row(subBtn))
			}

			subMenu.Inline(subRows...)
			bot.botInstance.Handle(&btn, func(c tele.Context) error {
				_ = bot.botInstance.Delete(c.Message())

				return c.Send(cmd.Name, subMenu, tele.ModeMarkdown)
			})

			bot.botInstance.Handle(fmt.Sprintf("/%s", cmd.Name), func(ctx tele.Context) error {
				_ = bot.botInstance.Delete(ctx.Message())

				return ctx.Send(cmd.Name, subMenu, tele.ModeMarkdown)
			})
		} else {
			bot.botInstance.Handle(&btn, func(ctx tele.Context) error {
				if len(cmd.Args) > 0 {
					return bot.handleArgCommand(ctx, []string{cmd.Name}, cmd.Args)
				}

				_ = bot.botInstance.Delete(ctx.Message())

				return bot.handleCommand(ctx, []string{cmd.Name})
			})

			bot.botInstance.Handle(fmt.Sprintf("/%s", cmd.Name), func(ctx tele.Context) error {
				_ = bot.botInstance.Delete(ctx.Message())

				err := bot.handleCommand(ctx, []string{cmd.Name})

				return err
			})
		}
	}

	// initiate menu button
	_ = bot.botInstance.SetCommands(commands)
	menu.Inline(rows...)
	bot.botInstance.Handle("/start", func(ctx tele.Context) error {
		_ = bot.botInstance.Delete(ctx.Message())

		return ctx.Send("Pagu Main Menu", menu, tele.ModeMarkdown)
	})

	bot.botInstance.Handle(tele.OnText, func(ctx tele.Context) error {
		if argsContext[ctx.Message().Sender.ID] == nil {
			return nil
		}

		if argsValue[ctx.Message().Sender.ID] == nil {
			argsValue[ctx.Message().Sender.ID] = make(map[string]string)
		}

		return bot.parsTextMessage(ctx)
	})

	return nil
}

func (bot *Bot) parsTextMessage(ctx tele.Context) error {
	senderID := ctx.Message().Sender.ID
	cmd := findCommand(bot.engine.Commands(), senderID)
	if cmd == nil {
		return ctx.Send("Invalid command")
	}

	currentArgsIndex := len(argsValue[senderID])
	argsValue[senderID][cmd.Args[currentArgsIndex].Name] = ctx.Message().Text

	if len(argsValue[senderID]) == len(cmd.Args) {
		return bot.handleCommand(ctx, argsContext[senderID].Commands)
	}

	_ = bot.botInstance.Delete(ctx.Message())

	return ctx.Send(fmt.Sprintf("Please Enter %s", cmd.Args[currentArgsIndex+1].Name))
}

func (bot *Bot) handleArgCommand(ctx tele.Context, commands []string, args []*command.Args) error {
	msgCtx := &BotContext{Commands: commands}
	argsContext[ctx.Sender().ID] = msgCtx
	argsValue[ctx.Sender().ID] = nil
	_ = bot.botInstance.Delete(ctx.Message())

	return ctx.Send(fmt.Sprintf("Please Enter %s", args[0].Name))
}

// handleCommand executes a command with its arguments for the user.
// It combines the commands and arguments into a single string, calls the engine's Run method,
// clears the user's context, and sends the result back to the user.
func (bot *Bot) handleCommand(ctx tele.Context, commands []string) error {
	callerID := strconv.Itoa(int(ctx.Sender().ID))

	// Retrieve the arguments for the sender
	senderID := ctx.Message().Sender.ID
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
	_ = bot.botInstance.Delete(ctx.Message())

	// Clear the stored command context and arguments for the sender
	argsContext[senderID] = nil
	argsValue[senderID] = nil

	return ctx.Send(res.Message, tele.NoPreview, tele.ModeMarkdown)
}

func findCommand(commands []*command.Command, senderID int64) *command.Command {
	lastEnteredCommandIndex := len(argsContext[senderID].Commands) - 1
	enteredCommand := argsContext[senderID].Commands[lastEnteredCommandIndex]

	for _, cmd := range commands {
		if cmd.Name == enteredCommand {
			return cmd
		}

		for _, sc := range cmd.SubCommands {
			if sc.Name == enteredCommand {
				return sc
			}
		}
	}

	return nil
}
