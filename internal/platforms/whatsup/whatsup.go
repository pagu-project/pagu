package whatsup

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/labstack/gommon/log"
	"github.com/pagu-project/pagu/config"
	"github.com/pagu-project/pagu/internal/engine"
	"github.com/pagu-project/pagu/internal/engine/command"
	"github.com/pagu-project/pagu/internal/entity"
	"github.com/tidwall/pretty"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	tele "gopkg.in/telebot.v3"
)

type Bot struct {
	ctx         context.Context
	cancel      context.CancelFunc
	botInstance *fiber.App
	engine      *engine.BotEngine
	cfg         *config.Config
	target      string
}

type BotContext struct {
	Commands []string
}

var (
	argsContext = make(map[int64]*BotContext)
	argsValue   = make(map[int64]map[string]string)

	WEBHOOK_VERIFY_TOKEN string
	GRAPH_API_TOKEN      string
	PORT                 int

	storage = make(map[string]InteractiveMessage)
)

type InteractiveMessage struct {
	MessagingProduct string `json:"messaging_product"`
	RecipientType    string `json:"recipient_type"`
	To               string `json:"to"`
	Type             string `json:"type"`
	Interactive      any    `json:"interactive"`
}

type WebhookRequest struct {
	Object string `json:"object"`
	Entry  []struct {
		ID      string `json:"id"`
		Changes []struct {
			Value struct {
				MessagingProduct string `json:"messaging_product"`
				Metadata         struct {
					DisplayPhoneNumber string `json:"display_phone_number"`
					PhoneNumberID      string `json:"phone_number_id"`
				} `json:"metadata"`
				Contacts []struct {
					Profile struct {
						Name string `json:"name"`
					} `json:"profile"`
					WaID string `json:"wa_id"`
				} `json:"contacts"`
				Messages []struct {
					From      string `json:"from"`
					ID        string `json:"id"`
					Timestamp string `json:"timestamp"`
					Text      struct {
						Body string `json:"body"`
					} `json:"text"`
					Type        string `json:"type"`
					Interactive struct {
						Type      string `json:"type"`
						ListReply struct {
							Id          string `json:"id"`
							Title       string `json:"title"`
							Description string `json:"description"`
						} `json:"list_reply"`
					} `json:"interactive"`
				} `json:"messages"`
				Field string `json:"field"`
			} `json:"value"`
		} `json:"changes"`
	} `json:"entry"`
}

func webhookHandler(c *fiber.Ctx) error {
	var resBody WebhookRequest

	// log.Println("Incoming webhook message: ", string(c.Body()))
	if err := json.Unmarshal(c.Body(), &resBody); err != nil {
		log.Printf("Error unmarshalling response body: %v", err)
		return c.Status(fiber.StatusBadRequest).SendString("Unable to parse request body")
	}

	// Log incoming message for debugging

	// Check if there are entries and changes in the webhook
	if len(resBody.Entry) > 0 {
		for _, entry := range resBody.Entry {
			for _, change := range entry.Changes {
				// Ensure there are messages in the change
				if len(change.Value.Messages) > 0 {
					message := change.Value.Messages[0]
					if message.Type == "text" {
						log.Printf("Received text message: %+v", message)

						// Extract phone number ID
						phoneNumberID := change.Value.Metadata.PhoneNumberID

						// Send List Message response
						sendHelpCommand(phoneNumberID, message.From)

					} else if message.Type == "interactive" {
						log.Printf("Received interactive message: %+v", message)

						// Extract phone number ID
						phoneNumberID := change.Value.Metadata.PhoneNumberID

						// Send List Message response
						// sendHelpCommand(phoneNumberID, message.From)
						switch message.Interactive.ListReply.Title {
						case "crowdfund":
							sendCommand("crowdfund", phoneNumberID, message.From)
						case "calculator":
							sendCommand("calculator", phoneNumberID, message.From)
						case "network":
							sendCommand("network", phoneNumberID, message.From)
						case "voucher":
							sendCommand("voucher", phoneNumberID, message.From)
						case "market":
							sendCommand("market", phoneNumberID, message.From)
						case "phoenix":
							sendCommand("phoenix", phoneNumberID, message.From)
						case "about":
							sendCommand("about", phoneNumberID, message.From)
						case "help":
							sendCommand("help", phoneNumberID, message.From)
						}
					}
				}
			}
		}
	}

	return c.SendStatus(fiber.StatusOK)
}

func verificationHandler(c *fiber.Ctx) error {
	mode := c.Query("hub.mode")
	token := c.Query("hub.verify_token")
	challenge := c.Query("hub.challenge")

	if mode == "subscribe" && token == WEBHOOK_VERIFY_TOKEN {
		return c.Status(fiber.StatusOK).SendString(challenge)
	}

	return c.Status(fiber.StatusForbidden).SendString("Forbidden")
}

func sendHelpCommand(phoneNumberID, to string) {
	message := map[string]any{
		"command":           "help",
		"messaging_product": "whatsapp",
		"recipient_type":    "individual",
		"to":                to,
		"type":              "interactive",
		"interactive": map[string]any{
			"type": "list",
			"body": map[string]any{
				"text": "\n\n\npagu 🤖 \nUsage: \npagu [subcommand] \n\nAvailable Subcommands: ",
			},
			"action": map[string]any{
				"button": "View Options",
				"sections": []any{
					map[string]any{
						"title": "Menu",
						"rows": []any{
							map[string]any{"id": "1", "title": "crowdfund", "description": "🤝 Commands for managing crowdfunding campaigns"},
							map[string]any{"id": "2", "title": "calculator", "description": "🧮 Perform calculations such as reward and fee estimations"},
							map[string]any{"id": "3", "title": "network", "description": "🌐 Commands for network metrics and information"},
							map[string]any{"id": "4", "title": "voucher", "description": "🎁 Commands for managing vouchers"},
							map[string]any{"id": "5", "title": "market", "description": "📈 Commands for managing market"},
							map[string]any{"id": "6", "title": "phoenix", "description": "🐦 Commands for working with Phoenix Testnet"},
							map[string]any{"id": "7", "title": "about", "description": "📝 About Pagu"},
							map[string]any{"id": "8", "title": "help", "description": "❓ Help for pagu command"},
						},
					},
				},
			},
		},
	}

	jsonData, err := json.Marshal(message)
	result := pretty.Pretty(jsonData)
	fmt.Printf("### 1\n", string(result))
	if err != nil {
		log.Printf("Error marshalling list message: %s", err)
		return
	}

	url := fmt.Sprintf("https://graph.facebook.com/v18.0/%s/messages", phoneNumberID)

	// Send the request using net/http (not fiber.Client)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		log.Printf("Error creating request: %s", err)
		return
	}

	// Set headers
	req.Header.Set("Authorization", "Bearer "+GRAPH_API_TOKEN)
	req.Header.Set("Content-Type", "application/json")

	// Send the request using the default HTTP client
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Error sending list message: %s", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("Failed to send list message: %s", resp.Status)
	}
}

func sendCommand(command, phoneNumberID, to string) {

	cmd := storage[command]
	cmd.To = to
	storage[command] = cmd

	jsonData, err := json.Marshal(storage[command])
	result := pretty.Pretty(jsonData)
	fmt.Printf("### 2\n", string(result))

	if err != nil {
		log.Printf("Error marshalling list message: %s", err)
		return
	}

	url := fmt.Sprintf("https://graph.facebook.com/v18.0/%s/messages", phoneNumberID)

	// Send the request using net/http (not fiber.Client)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		log.Printf("Error creating request: %s", err)
		return
	}

	// Set headers
	req.Header.Set("Authorization", "Bearer "+GRAPH_API_TOKEN)
	req.Header.Set("Content-Type", "application/json")

	// Send the request using the default HTTP client
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Error sending list message: %s", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("Failed to send list message: %s", resp.Status)
	}
}

func NewWhatsUpBot(botEngine *engine.BotEngine, cfg *config.Config) (*Bot, error) {
	WEBHOOK_VERIFY_TOKEN = cfg.WhatsUp.WebHookToken
	GRAPH_API_TOKEN = cfg.WhatsUp.GraphToken
	PORT = cfg.WhatsUp.Port

	app := fiber.New()

	// Webhook handlers
	app.Post("/webhook", webhookHandler)
	app.Get("/webhook", verificationHandler)

	// Default route
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("<pre>Nothing to see here. Checkout README.md to start.</pre>")
	})

	ctx, cancel := context.WithCancel(context.Background())

	return &Bot{
		engine:      botEngine,
		cfg:         cfg,
		botInstance: app,
		ctx:         ctx,
		cancel:      cancel,
		target:      cfg.BotName,
	}, nil
}

func (bot *Bot) Start() error {
	bot.deleteAllCommands()
	if err := bot.registerCommands(""); err != nil {
		return err
	}
	go func() {
		log.Printf("Server is listening on port: %s", PORT)
		if err := bot.botInstance.Listen(fmt.Sprintf(":%v", PORT)); err != nil {
			log.Fatalf("Error starting server: %s", err)
		}
	}()
	log.Info("Starting WhatsUp Bot...")

	return nil
}

func (bot *Bot) Stop() {
	log.Info("Shutting down WhatsUp Bot")
	bot.cancel()
}

func (bot *Bot) deleteAllCommands() {

}

//nolint:gocognit // Complexity cannot be reduced
func (bot *Bot) registerCommands(to string) error {
	commands := make([]tele.Command, 0)

	cmds := bot.engine.Commands()
	for i, cmd := range cmds {
		rowsSubCmd := []any{}
		// if !cmd.HasBotID(entity.BotID_Telegram) {
		// 	continue
		// }

		// switch bot.target {
		// case config.BotNamePaguMainnet:
		// 	if !utils.IsDefinedOnBotID(cmd.TargetBotIDs, entity.BotID_Telegram) {
		// 		continue
		// 	}

		// case config.BotNamePaguModerator:
		// 	if !utils.IsDefinedOnBotID(cmd.TargetBotIDs, entity.BotID_Moderator) {
		// 		continue
		// 	}

		// default:
		// 	log.Warn("invalid target", "target", bot.target)

		// 	continue
		// }

		log.Info("registering new command", "name", cmd.Name, "desc", cmd.Help, "index", i, "object", cmd)

		commands = append(commands, tele.Command{Text: cmd.Name, Description: cmd.Help})
		if cmd.HasSubCommand() {
			for indx, subCmd := range cmd.SubCommands {
				rowsSubCmd = append(rowsSubCmd, map[string]any{
					"id":          fmt.Sprintf("%v", indx),
					"title":       subCmd.Name,
					"description": subCmd.Help,
				})
			}
		}
		storage[cmd.Name] = InteractiveMessage{
			MessagingProduct: "whatsapp",
			RecipientType:    "individual",
			To:               to,
			Type:             "interactive",
			Interactive: map[string]any{
				"type": "list",
				"body": map[string]any{
					"text": cmd.Help,
				},
				"action": map[string]any{
					"button": "View Options",
					"sections": []any{
						map[string]any{
							"title": "Menu",
							"rows":  rowsSubCmd,
						},
					},
				},
			},
		}
		if cmd.HasSubCommand() {
			// subMenu := &tele.ReplyMarkup{ResizeKeyboard: true}
			// subRows := make([]tele.Row, 0)
			// for _, subCmd := range cmd.SubCommands {
			// 	switch bot.target {
			// 	case config.BotNamePaguMainnet:
			// 		if !utils.IsDefinedOnBotID(subCmd.TargetBotIDs, entity.BotID_Telegram) {
			// 			continue
			// 		}

			// 	case config.BotNamePaguModerator:
			// 		if !utils.IsDefinedOnBotID(subCmd.TargetBotIDs, entity.BotID_Moderator) {
			// 			continue
			// 		}

			// 	default:
			// 		log.Warn("invalid target", "target", bot.target)

			// 		continue
			// 	}

			// 	log.Info("adding command sub-command", "command", cmd.Name,
			// 		"sub-command", subCmd.Name, "desc", subCmd.Help)

			// 	subBtn := subMenu.Data(cases.Title(language.English).String(subCmd.Name), cmd.Name+subCmd.Name)

			// 	// bot.botInstance.Handle(&subBtn, func(c tele.Context) error {
			// 	// 	if len(subCmd.Args) > 0 {
			// 	// 		return bot.handleArgCommand(c, []string{cmd.Name, subCmd.Name}, subCmd.Args)
			// 	// 	}

			// 	// 	return bot.handleCommand(c, []string{cmd.Name, subCmd.Name})
			// 	// })
			// 	subRows = append(subRows, subMenu.Row(subBtn))
			// }

			// // subMenu.Inline(subRows...)
			// // bot.botInstance.Handle(&btn, func(c tele.Context) error {
			// // 	_ = bot.botInstance.Delete(c.Message())

			// // 	return c.Send(cmd.Name, subMenu, tele.ModeMarkdown)
			// // })

			// // bot.botInstance.Handle(fmt.Sprintf("/%s", cmd.Name), func(ctx tele.Context) error {
			// // 	_ = bot.botInstance.Delete(ctx.Message())

			// // 	return ctx.Send(cmd.Name, subMenu, tele.ModeMarkdown)
			// // })
		} else {
			// bot.botInstance.Handle(fmt.Sprintf("/%s", cmd.Name), func(ctx tele.Context) error {
			// 	_ = bot.botInstance.Delete(ctx.Message())

			// 	err := bot.handleCommand(ctx, []string{cmd.Name})

			// 	return err
			// })
		}
	}

	// initiate menu button
	// _ = bot.botInstance.SetCommands(commands)
	// bot.botInstance.Handle("/start", func(ctx tele.Context) error {
	// 	_ = bot.botInstance.Delete(ctx.Message())

	// 	return ctx.Send("Pagu Main Menu", menu, tele.ModeMarkdown)
	// })

	// bot.botInstance.Handle(tele.OnText, func(ctx tele.Context) error {
	// 	if argsContext[ctx.Message().Sender.ID] == nil {
	// 		return nil
	// 	}

	// 	if argsValue[ctx.Message().Sender.ID] == nil {
	// 		argsValue[ctx.Message().Sender.ID] = make(map[string]string)
	// 	}

	// 	return bot.parsTextMessage(ctx)
	// })

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

	// _ = bot.botInstance.Delete(ctx.Message())

	return ctx.Send(fmt.Sprintf("Please Enter %s", cmd.Args[currentArgsIndex+1].Name))
}

func (bot *Bot) handleArgCommand(ctx tele.Context, commands []string, args []*command.Args) error {
	msgCtx := &BotContext{Commands: commands}
	argsContext[ctx.Sender().ID] = msgCtx
	argsValue[ctx.Sender().ID] = nil
	// _ = bot.botInstance.Delete(ctx.Message())

	choiceMenu := &tele.ReplyMarkup{ResizeKeyboard: true}
	choiceRows := make([]tele.Row, 0)
	choiceMeg := fmt.Sprintf("Please Select a %s\nChoose the best option below based on your preference:\n", args[0].Name)
	for _, arg := range args {
		if len(arg.Choices) > 0 {
			for _, choice := range arg.Choices {
				choices := strings.Split(choice.Name, " ")
				choiceMeg += fmt.Sprintf("- %s : %s\n", choices[0], strings.Join(choices[1:], " "))
				choiceBtn := choiceMenu.Data(cases.Title(language.English).String(choices[0]), choices[0])
				choiceRows = append(choiceRows, choiceMenu.Row(choiceBtn))
				// bot.botInstance.Handle(&choiceBtn, func(c tele.Context) error {
				// 	choices = strings.Split(choices[0], "-")
				// 	commands = append(commands, fmt.Sprintf("--%s=%v", choices[0], choices[1]))

				// 	return bot.handleCommand(c, commands)
				// })
			}
		}
	}
	choiceMenu.Inline(choiceRows...)

	return ctx.Send(choiceMeg, choiceMenu)
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
	// _ = bot.botInstance.Delete(ctx.Message())

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
