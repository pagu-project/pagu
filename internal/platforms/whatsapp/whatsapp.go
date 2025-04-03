package whatsapp

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"slices"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/labstack/gommon/log"
	"github.com/pagu-project/pagu/config"
	"github.com/pagu-project/pagu/internal/engine"
	"github.com/pagu-project/pagu/internal/engine/command"
	"github.com/pagu-project/pagu/internal/entity"
)

const (
	COMMAND    = "command"
	SUBCOMMAND = "subCommand"
)

type Bot struct {
	ctx         context.Context
	cancel      context.CancelFunc
	botInstance *fiber.App
	engine      *engine.BotEngine
	cmds        []*command.Command
	cfg         *config.Config

	target         string
	sessionManager *SessionManager
}

func (bot *Bot) renderPage(cmdName, to string) InteractiveMessage {
	var (
		rowsSubCmd []any
		command    *command.Command
	)

	for _, cmd := range bot.cmds {
		if cmd.Name == cmdName {
			command = cmd
			break
		}
		if cmd.HasSubCommand() {
			for _, subCmd := range cmd.SubCommands {
				if subCmd.Name == cmdName {
					command = cmd
					break
				}
			}
		}
	}

	for indx, subCmd := range command.SubCommands {
		rowsSubCmd = append(rowsSubCmd, map[string]any{
			"id":          fmt.Sprintf("%v", indx),
			"title":       subCmd.Name,
			"description": subCmd.Help,
		})
	}

	return InteractiveMessage{
		MessagingProduct: "whatsapp",
		RecipientType:    "individual",
		To:               to,
		Type:             "interactive",
		Interactive: map[string]any{
			"type": "list",
			"body": map[string]any{
				"text": command.Help,
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
}

func renderResult(result, to string) map[string]any {
	return map[string]any{
		"messaging_product": "whatsapp",
		"recipient_type":    "individual",
		"to":                to,
		"type":              "text",
		"text": map[string]any{
			"body": result,
		},
	}
}

func (bot *Bot) checkCommand(command string) string {
	for _, cmd := range bot.cmds {
		if cmd.Name == command {
			return COMMAND
		}
		if cmd.HasSubCommand() {
			for _, subCmd := range cmd.SubCommands {
				if subCmd.Name == command {
					return SUBCOMMAND
				}
			}
		}
	}
	return ""
}

func (bot *Bot) findCommand(subCommand string) string {
	for _, cmd := range bot.cmds {
		for _, subCmd := range cmd.SubCommands {
			if subCmd.Name == subCommand {
				return cmd.Name
			}
		}
	}
	return ""
}

func (bot *Bot) findArgs(subCommand string) []string {
	for _, cmd := range bot.cmds {
		for _, subCmd := range cmd.SubCommands {
			if subCmd.Name == subCommand {
				args := []string{}
				for _, arg := range subCmd.Args {
					args = append(args, arg.Name)
				}
				return args
			}
		}
	}
	return nil
}

var (
	WEBHOOK_VERIFY_TOKEN string
	GRAPH_API_TOKEN      string
	PORT                 int
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

func (bot *Bot) webhookHandler(c *fiber.Ctx) error {
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
					fmt.Printf("----message : %+v\n", message)
					phoneNumberID := change.Value.Metadata.PhoneNumberID
					if message.Type == "interactive" {
						msg := message.Interactive.ListReply.Title
						fmt.Println("----msg : ", msg)
						switch bot.checkCommand(msg) {
						case COMMAND:
							bot.sessionManager.OpenSession(phoneNumberID, Session{
								commands: []string{msg},
								args:     nil,
							})
						case SUBCOMMAND:
							mainCommand := bot.findCommand(msg)
							bot.sessionManager.OpenSession(phoneNumberID, Session{
								commands: []string{mainCommand, msg},
								args:     nil,
							})
						default:
							fmt.Println("error in add session")
						}
						bot.sendCommand(phoneNumberID, message.From)
					} else {
						if strings.ToLower(message.Text.Body) == "help" || strings.ToLower(message.Text.Body) == "start" {
							bot.sessionManager.OpenSession(phoneNumberID, Session{
								commands: []string{"help"},
								args:     nil,
							})
							bot.sendHelpCommand(phoneNumberID, message.From)
						} else {
							msg := message.Text.Body
							session := bot.sessionManager.GetSession(phoneNumberID)
							args := session.args
							args = append(args, msg)
							session.args = args
							bot.sessionManager.OpenSession(phoneNumberID, *session)
							bot.sendCommand(phoneNumberID, message.From)
						}
					}
				}
			}
		}
	}

	return c.SendStatus(fiber.StatusOK)
}

func (bot *Bot) verificationHandler(c *fiber.Ctx) error {
	mode := c.Query("hub.mode")
	token := c.Query("hub.verify_token")
	challenge := c.Query("hub.challenge")

	if mode == "subscribe" && token == WEBHOOK_VERIFY_TOKEN {
		return c.Status(fiber.StatusOK).SendString(challenge)
	}

	return c.Status(fiber.StatusForbidden).SendString("Forbidden")
}

func (bot *Bot) sendHelpCommand(phoneNumberID, to string) {
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
	// result := pretty.Pretty(jsonData)
	// fmt.Printf("### 1\n", string(result)
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

func findArg(argsMap map[string][]string, arg string) string {
	for command, args := range argsMap {
		if slices.Contains(args, arg) {
			return command
		}
	}
	return ""
}

func (bot *Bot) sendCommand(phoneNumberID, to string) {
	var (
		jsonData   []byte
		err        error
		commandRes []byte
		session    = bot.sessionManager.GetSession(phoneNumberID)
	)

	if len(session.commands) == 1 {
		// commandRes = bot.handleCommand(session.commands)
		cmd := bot.renderPage(session.commands[0], to)
		jsonData, err = json.Marshal(cmd)
	} else if len(session.commands) == 2 {
		args := bot.findArgs(session.commands[1])
		if len(args) > 0 {
			if len(session.args) != len(args) {
				commandRes = []byte(fmt.Sprintf("Enter your %s : ", args[len(session.args)]))
			} else {
				for indx, arg := range session.args {
					session.commands = append(session.commands, fmt.Sprintf("--%s=%s", args[indx], arg))
				}
				commandRes = bot.handleCommand(session.commands)
			}
		} else {
			commandRes = bot.handleCommand(session.commands)
		}
		cmd := renderResult(string(commandRes), to)
		jsonData, err = json.Marshal(cmd)
	}

	// result := pretty.Pretty(jsonData)
	// fmt.Printf("### 2\n", string(result))

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

func NewWhatsAppBot(botEngine *engine.BotEngine, cfg *config.Config) (*Bot, error) {
	WEBHOOK_VERIFY_TOKEN = cfg.WhatsApp.WebHookToken
	GRAPH_API_TOKEN = cfg.WhatsApp.GraphToken
	PORT = cfg.WhatsApp.Port

	app := fiber.New()
	ctx, cancel := context.WithCancel(context.Background())

	cmds := botEngine.Commands()

	sessionManager := NewSessionManager()
	sessionManager.checkInterval = 600 * time.Second
	sessionManager.sessionTtl = 300 * time.Second

	bot := &Bot{
		cmds:           cmds,
		engine:         botEngine,
		cfg:            cfg,
		botInstance:    app,
		ctx:            ctx,
		cancel:         cancel,
		target:         cfg.BotName,
		sessionManager: sessionManager,
	}

	go bot.sessionManager.removeExpiredSession()

	// Webhook handlers
	app.Post("/webhook", bot.webhookHandler)
	app.Get("/webhook", bot.verificationHandler)

	// Default route
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("<pre>Nothing to see here. Checkout README.md to start.</pre>")
	})

	return bot, nil
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
	log.Info("Starting WhatsApp Bot...")

	return nil
}

func (bot *Bot) Stop() {
	log.Info("Shutting down WhatsApp Bot")
	bot.cancel()
}

func (bot *Bot) deleteAllCommands() {
}

//nolint:gocognit // Complexity cannot be reduced
func (bot *Bot) registerCommands(to string) error {
	return nil
}

func (bot *Bot) parsTextMessage() error {
	return nil
}

func (bot *Bot) handleArgCommand(commands []string, args map[string]string) []byte {
	// choiceMeg := fmt.Sprintf("Please Select a %s\nChoose the best option below based on your preference:\n", args[0])
	for key, val := range args {
		commands = append(commands, fmt.Sprintf("--%s=%s", key, val))
	}
	return bot.handleCommand(commands)
}

// handleCommand executes a command with its arguments for the user.
// It combines the commands and arguments into a single string, calls the engine's Run method,
// clears the user's context, and sends the result back to the user.
func (bot *Bot) handleCommand(commands []string) []byte {

	// Retrieve the arguments for the sender
	// fmt.Println("+++++commands : ", commands)
	// Combine the commands into a single string
	fullCommand := strings.Join(commands, " ")

	// Call the engine's Run method with the full command string
	res := bot.engine.ParseAndExecute(entity.PlatformIDTelegram, "", fullCommand)
	// _ = bot.botInstance.Delete(ctx.Message())

	// Clear the stored command context and arguments for the sender

	return []byte(res.Message)
}
