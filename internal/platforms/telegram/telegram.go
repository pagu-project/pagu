package telegram

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/pactus-project/pactus/util"
	"github.com/pagu-project/Pagu/config"
	"github.com/pagu-project/Pagu/internal/engine"
	"github.com/pagu-project/Pagu/internal/engine/command"
	"github.com/pagu-project/Pagu/internal/entity"
	"github.com/pagu-project/Pagu/pkg/log"
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

//nolint:gocognit // Complexity cannot be reduced
func (bot *Bot) registerCommands() error {
	rows := make([]tele.Row, 0)
	menu := &tele.ReplyMarkup{ResizeKeyboard: true}
	commands := make([]tele.Command, 0)

	for i, beCmd := range bot.engine.Commands() {
		if !beCmd.HasAppID(entity.AppIDTelegram) {
			continue
		}

		switch bot.target {
		case config.BotNamePaguMainnet:
			if !util.IsFlagSet(beCmd.TargetFlag, command.TargetMaskMain) {
				continue
			}

		case config.BotNamePaguTestnet:
			if !util.IsFlagSet(beCmd.TargetFlag, command.TargetMaskTest) {
				continue
			}
		}

		log.Info("registering new command", "name", beCmd.Name, "desc", beCmd.Help, "index", i, "object", beCmd)

		btn := menu.Data(cases.Title(language.English).String(beCmd.Name), beCmd.Name)
		commands = append(commands, tele.Command{Text: beCmd.Name, Description: beCmd.Help})
		rows = append(rows, menu.Row(btn))
		if beCmd.HasSubCommand() {
			subMenu := &tele.ReplyMarkup{ResizeKeyboard: true}
			subRows := make([]tele.Row, 0)
			for _, sCmd := range beCmd.SubCommands {
				switch bot.target {
				case config.BotNamePaguMainnet:
					if !util.IsFlagSet(sCmd.TargetFlag, command.TargetMaskMain) {
						continue
					}

				case config.BotNamePaguTestnet:
					if !util.IsFlagSet(sCmd.TargetFlag, command.TargetMaskTest) {
						continue
					}

				case config.BotNamePaguModerator:
					if !util.IsFlagSet(sCmd.TargetFlag, command.TargetMaskModerator) {
						continue
					}
				}

				log.Info("adding command sub-command", "command", beCmd.Name,
					"sub-command", sCmd.Name, "desc", sCmd.Help)

				subBtn := subMenu.Data(cases.Title(language.English).String(sCmd.Name), sCmd.Name)

				bot.botInstance.Handle(&subBtn, func(c tele.Context) error {
					if len(sCmd.Args) > 0 {
						return bot.handleArgCommand(c, []string{beCmd.Name, sCmd.Name}, sCmd.Args)
					}

					return bot.handleCommand(c, []string{beCmd.Name, sCmd.Name})
				})
				subRows = append(subRows, subMenu.Row(subBtn))
			}

			subMenu.Inline(subRows...)
			bot.botInstance.Handle(&btn, func(c tele.Context) error {
				_ = bot.botInstance.Delete(c.Message())

				return c.Send(beCmd.Name, subMenu)
			})

			bot.botInstance.Handle(fmt.Sprintf("/%s", beCmd.Name), func(c tele.Context) error {
				_ = bot.botInstance.Delete(c.Message())

				return c.Send(beCmd.Name, subMenu)
			})
		} else {
			bot.botInstance.Handle(&btn, func(ctx tele.Context) error {
				if len(beCmd.Args) > 0 {
					return bot.handleArgCommand(ctx, []string{beCmd.Name}, beCmd.Args)
				}

				_ = bot.botInstance.Delete(ctx.Message())

				return bot.handleCommand(ctx, []string{beCmd.Name})
			})
		}
	}

	// initiate menu button
	_ = bot.botInstance.SetCommands(commands)
	menu.Inline(rows...)
	bot.botInstance.Handle("/start", func(c tele.Context) error {
		_ = bot.botInstance.Delete(c.Message())

		return c.Send("Pagu Main Menu", menu)
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

func (bot *Bot) handleArgCommand(ctx tele.Context, commands []string, args []command.Args) error {
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
			argPairs = append(argPairs, fmt.Sprintf("%s=%s", key, value))
		}
		fullCommand = fmt.Sprintf("%s %s", fullCommand, strings.Join(argPairs, " "))
	}

	// Call the engine's Run method with the full command string
	res := bot.engine.ParseAndExecute(entity.AppIDTelegram, callerID, fullCommand)
	_ = bot.botInstance.Delete(ctx.Message())

	// Clear the stored command context and arguments for the sender
	argsContext[senderID] = nil
	argsValue[senderID] = nil

	return ctx.Send(res.Message, tele.NoPreview)
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
