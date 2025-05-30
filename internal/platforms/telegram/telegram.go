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
	"github.com/pagu-project/pagu/pkg/utils"
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

	cmds := bot.engine.Commands()
	for i, cmd := range cmds {
		if !cmd.HasBotID(entity.BotID_Telegram) {
			continue
		}

		switch bot.target {
		case config.BotNamePaguMainnet:
			if !utils.IsDefinedOnBotID(cmd.TargetBotIDs, entity.BotID_Telegram) {
				continue
			}

		case config.BotNamePaguModerator:
			if !utils.IsDefinedOnBotID(cmd.TargetBotIDs, entity.BotID_Moderator) {
				continue
			}

		default:
			log.Warn("invalid target", "target", bot.target)

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
				switch bot.target {
				case config.BotNamePaguMainnet:
					if !utils.IsDefinedOnBotID(subCmd.TargetBotIDs, entity.BotID_Telegram) {
						continue
					}

				case config.BotNamePaguModerator:
					if !utils.IsDefinedOnBotID(subCmd.TargetBotIDs, entity.BotID_Moderator) {
						continue
					}

				default:
					log.Warn("invalid target", "target", bot.target)

					continue
				}

				log.Info("adding command sub-command", "command", cmd.Name,
					"sub-command", subCmd.Name, "desc", subCmd.Help)

				subBtn := subMenu.Data(cases.Title(language.English).String(subCmd.Name), cmd.Name+subCmd.Name)

				bot.botInstance.Handle(&subBtn, func(tgCtx tele.Context) error {
					if len(subCmd.Args) > 0 {
						return bot.handleArgCommand(tgCtx, []string{cmd.Name, subCmd.Name}, subCmd.Args)
					}

					return bot.handleCommand(tgCtx, []string{cmd.Name, subCmd.Name})
				})
				subRows = append(subRows, subMenu.Row(subBtn))
			}

			subMenu.Inline(subRows...)
			bot.botInstance.Handle(&btn, func(tgCtx tele.Context) error {
				_ = bot.botInstance.Delete(tgCtx.Message())

				return bot.sendMarkdown(tgCtx, cmd.Name, subMenu)
			})

			bot.botInstance.Handle(fmt.Sprintf("/%s", cmd.Name), func(tgCtx tele.Context) error {
				_ = bot.botInstance.Delete(tgCtx.Message())

				return bot.sendMarkdown(tgCtx, cmd.Name, subMenu)
			})
		} else {
			bot.botInstance.Handle(fmt.Sprintf("/%s", cmd.Name), func(tgCtx tele.Context) error {
				_ = bot.botInstance.Delete(tgCtx.Message())

				err := bot.handleCommand(tgCtx, []string{cmd.Name})

				return err
			})
		}
	}

	// initiate menu button
	_ = bot.botInstance.SetCommands(commands)
	menu.Inline(rows...)
	bot.botInstance.Handle("/start", func(tgCtx tele.Context) error {
		_ = bot.botInstance.Delete(tgCtx.Message())

		return bot.sendMarkdown(tgCtx, "Pagu Main Menu", menu)
	})

	bot.botInstance.Handle(tele.OnText, func(tgCtx tele.Context) error {
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
	cmd := findCommand(bot.engine.Commands(), senderID)
	if cmd == nil {
		return bot.sendMarkdown(tgCtx, "Invalid command")
	}

	currentArgsIndex := len(argsValue[senderID])
	argsValue[senderID][cmd.Args[currentArgsIndex].Name] = tgCtx.Message().Text

	if len(argsValue[senderID]) == len(cmd.Args) {
		return bot.handleCommand(tgCtx, argsContext[senderID].Commands)
	}

	_ = bot.botInstance.Delete(tgCtx.Message())

	return bot.sendMarkdown(tgCtx, fmt.Sprintf("Please enter `%s`:", cmd.Args[currentArgsIndex+1].Name))
}

func (bot *Bot) handleArgCommand(tgCtx tele.Context, commands []string, args []*command.Args) error {
	msgCtx := &BotContext{Commands: commands}
	argsContext[tgCtx.Sender().ID] = msgCtx
	argsValue[tgCtx.Sender().ID] = nil
	_ = bot.botInstance.Delete(tgCtx.Message())

	firstArg := args[0]

	if len(firstArg.Choices) > 0 {
		choiceMsg := fmt.Sprintf("Please select a `%s`:\n\n", firstArg.Name)
		choiceMenu := &tele.ReplyMarkup{ResizeKeyboard: true}
		choiceRows := make([]tele.Row, 0, len(firstArg.Choices))
		for _, choice := range firstArg.Choices {
			choiceMsg += fmt.Sprintf("- %s", choice.Desc)
			choiceBtn := choiceMenu.Data(choice.Name, choice.Name, choice.Value)
			choiceRows = append(choiceRows, choiceMenu.Row(choiceBtn))
			bot.botInstance.Handle(&choiceBtn, func(tCtx tele.Context) error {
				commands = append(commands, fmt.Sprintf("--%s=%v", firstArg.Name, choice.Value))

				return bot.handleCommand(tgCtx, commands)
			})
		}

		choiceMenu.Inline(choiceRows...)

		return bot.sendMarkdown(tgCtx, choiceMsg, choiceMenu)
	}

	// Commands with no choices
	return bot.sendMarkdown(tgCtx, fmt.Sprintf("Please enter `%s`:", firstArg.Name))
}

// handleCommand executes a command with its arguments for the user.
// It combines the commands and arguments into a single string, calls the engine's Run method,
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
	_ = bot.botInstance.Delete(tgCtx.Message())

	// Clear the stored command context and arguments for the sender
	argsContext[senderID] = nil
	argsValue[senderID] = nil

	return bot.sendMarkdown(tgCtx, res.Message, tele.NoPreview)
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

func (bot *Bot) sendMarkdown(tgCtx tele.Context, what interface{}, opts ...interface{}) error {
	opts = append(opts, tele.ModeMarkdownV2)

	return tgCtx.Send(what, opts)
}
