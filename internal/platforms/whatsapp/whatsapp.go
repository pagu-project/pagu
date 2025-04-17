package whatsapp

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/labstack/gommon/log"
	"github.com/pagu-project/pagu/config"
	"github.com/pagu-project/pagu/internal/engine"
	"github.com/pagu-project/pagu/internal/engine/command"
	"github.com/pagu-project/pagu/internal/entity"
	"github.com/pagu-project/pagu/pkg/session"
)

type Bot struct {
	ctx    context.Context
	cancel context.CancelFunc
	server *http.ServeMux
	engine *engine.BotEngine
	cmds   []*command.Command
	cfg    *config.Config
	hook   Webhook

	target         string
	sessionManager *session.SessionManager
}

type Webhook struct {
	verifyToken    string
	graphAPIToken  string
	webHookAddredd string
}

type InteractiveMessage struct {
	MessagingProduct string `json:"messaging_product"`
	RecipientType    string `json:"recipient_type"`
	To               string `json:"to"`
	Type             string `json:"type"`
	Interactive      any    `json:"interactive"`
}

type WebhookRequest struct {
	Object string  `json:"object"`
	Entry  []Entry `json:"entry"`
}

type Entry struct {
	ID      string   `json:"id"`
	Changes []Change `json:"changes"`
}

type Change struct {
	Value Value `json:"value"`
}

type Value struct {
	MessagingProduct string    `json:"messaging_product"`
	Metadata         Metadata  `json:"metadata"`
	Contacts         []Contact `json:"contacts"`
	Messages         []Message `json:"messages"`
	Field            string    `json:"field"`
}

type Metadata struct {
	DisplayPhoneNumber string `json:"display_phone_number"`
	PhoneNumberID      string `json:"phone_number_id"`
}

type Contact struct {
	Profile Profile `json:"profile"`
	WaID    string  `json:"wa_id"`
}

type Profile struct {
	Name string `json:"name"`
}

type Message struct {
	From        string      `json:"from"`
	ID          string      `json:"id"`
	Timestamp   string      `json:"timestamp"`
	Text        Text        `json:"text"`
	Type        string      `json:"type"`
	Interactive Interactive `json:"interactive"`
}

type Text struct {
	Body string `json:"body"`
}

type Interactive struct {
	Type      string    `json:"type"`
	ListReply ListReply `json:"list_reply"`
}

type ListReply struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
}

func (bot *Bot) renderPage(cmdName, destination string) (InteractiveMessage, error) {
	var command *command.Command
	rowsSubCmd := []any{}

	for _, cmd := range bot.cmds {
		if cmd.Name == cmdName {
			command = cmd

			break
		}
	}

	if command == nil {
		return InteractiveMessage{}, errors.New("render page for subcommand is not possible")
	}

	if command.Name == "help" || command.Name == "about" {
		return InteractiveMessage{}, errors.New("render page for connamd without subcommand is not possible")
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
		RecipientType:    "indivIDual",
		To:               destination,
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
	}, nil
}

func renderResult(result, destination string) map[string]any {
	return map[string]any{
		"messaging_product": "whatsapp",
		"recipient_type":    "indivIDual",
		"to":                destination,
		"type":              "text",
		"text": map[string]any{
			"body": result,
		},
	}
}

func renderRootCommand(destinatoin string) map[string]any {
	return map[string]any{
		"command":           "help",
		"messaging_product": "whatsapp",
		"recipient_type":    "individual",
		"to":                destinatoin,
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
							map[string]any{
								"id":          "1",
								"title":       "crowdfund",
								"description": "🤝 Commands for managing crowdfunding campaigns",
							},
							map[string]any{
								"id":          "2",
								"title":       "calculator",
								"description": "🧮 Perform calculations such as reward and fee estimations",
							},
							map[string]any{
								"id":          "3",
								"title":       "network",
								"description": "🌐 Commands for network metrics and information",
							},
							map[string]any{
								"id":          "4",
								"title":       "voucher",
								"description": "🎁 Commands for managing vouchers",
							},
							map[string]any{
								"id":          "5",
								"title":       "market",
								"description": "📈 Commands for managing market",
							},
							map[string]any{
								"id":          "6",
								"title":       "phoenix",
								"description": "🐦 Commands for working with Phoenix Testnet",
							},
							map[string]any{
								"id":          "7",
								"title":       "about",
								"description": "📝 About Pagu",
							},
							map[string]any{
								"id":          "8",
								"title":       "help",
								"description": "❓ Help for pagu command",
							},
						},
					},
				},
			},
		},
	}
}

func (bot *Bot) listCommand(command string) []string {
	for _, cmd := range bot.cmds {
		if cmd.Name == command {
			return []string{cmd.Name}
		}
		if cmd.HasSubCommand() {
			for _, subCmd := range cmd.SubCommands {
				if subCmd.Name == command {
					return []string{cmd.Name, subCmd.Name}
				}
			}
		}
	}

	return nil
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

func (bot *Bot) webhookHandler(w http.ResponseWriter, r *http.Request) {
	var resBody WebhookRequest

	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Printf("Error reading body: %v", err)
		http.Error(w, "Unable to read body", http.StatusBadRequest)

		return
	}
	defer r.Body.Close()

	if err := json.Unmarshal(body, &resBody); err != nil {
		log.Printf("Error unmarshalling response body: %v", err)
		http.Error(w, "Unable to parse request body", http.StatusBadRequest)

		return
	}

	if len(resBody.Entry) > 0 {
		for _, entry := range resBody.Entry {
			for _, change := range entry.Changes {
				if len(change.Value.Messages) == 0 {
					continue
				}
				message := change.Value.Messages[0]
				log.Printf("message from webhook: %+v\n", message)

				phoneNumberID := change.Value.Metadata.PhoneNumberID
				session := session.NewSession()

				var msg string
				if message.Type == "interactive" {
					msg = message.Interactive.ListReply.Title
					session.Commands = bot.listCommand(msg)
				} else {
					msg = message.Text.Body
					if !strings.EqualFold(msg, "help") && !strings.EqualFold(msg, "about") {
						session = bot.sessionManager.GetSession(phoneNumberID)
						session.Args = append(session.Args, msg)
					}
				}

				bot.sessionManager.OpenSession(phoneNumberID, *session)
				bot.sendCommand(r.Context(), phoneNumberID, message.From)
			}
		}
	}

	w.WriteHeader(http.StatusOK)
}

func (bot *Bot) verificationHandler(w http.ResponseWriter, r *http.Request) {
	mode := r.URL.Query().Get("hub.mode")
	token := r.URL.Query().Get("hub.verify_token")
	challenge := r.URL.Query().Get("hub.challenge")

	if mode == "subscribe" && token == bot.hook.verifyToken {
		w.WriteHeader(http.StatusOK)
		_, err := fmt.Fprint(w, challenge)
		if err != nil {
			log.Print(err)
		}

		return
	}

	http.Error(w, "Forbidden", http.StatusForbidden)
}

func (bot *Bot) sendCommand(ctx context.Context, phoneNumberID, destination string) {
	var (
		jsonData   []byte
		err        error
		commandRes []byte
		session    = bot.sessionManager.GetSession(phoneNumberID)
	)

	if len(session.Commands) > 0 {
		cmd, err := bot.renderPage(session.Commands[len(session.Commands)-1], destination)
		if err == nil {
			jsonData, err = json.Marshal(cmd)
			if err != nil {
				log.Printf("Error marshalling list message: %s", err)

				return
			}
		} else {
			args := bot.findArgs(session.Commands[len(session.Commands)-1])
			if len(args) > 0 {
				if len(session.Args) < len(args) {
					commandRes = []byte(fmt.Sprintf("Enter %s: ", args[len(session.Args)]))
				} else {
					for indx, arg := range session.Args {
						if indx < len(args) {
							session.Commands = append(session.Commands, fmt.Sprintf("--%s=%s", args[indx], arg))
						}
					}
					commandRes = bot.handleCommand(session.Commands)
				}
			} else {
				commandRes = bot.handleCommand(session.Commands)
			}

			cmd := renderResult(string(commandRes), destination)
			jsonData, err = json.Marshal(cmd)
			if err != nil {
				log.Printf("Error marshalling list message: %s", err)

				return
			}
		}
	} else {
		cmd := renderRootCommand(destination)
		jsonData, err = json.Marshal(cmd)
		if err != nil {
			log.Printf("Error marshalling list message: %s", err)

			return
		}
	}

	url := fmt.Sprintf("https://graph.facebook.com/v18.0/%s/messages", phoneNumberID)

	// Send the request using net/http (not fiber.Client)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(jsonData))
	if err != nil {
		log.Printf("Error creating request: %s", err)

		return
	}

	// Set headers
	req.Header.Set("Authorization", "Bearer "+bot.hook.graphAPIToken)
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
	server := http.NewServeMux()
	ctx, cancel := context.WithCancel(context.Background())

	cmds := botEngine.Commands()

	sessionManager := session.NewSessionManager(ctx)
	sessionManager.SetConfig(
		time.Duration(cfg.Session.SessionTTL*int(time.Second)),
		time.Duration(cfg.Session.CheckInterval*int(time.Second)),
	)

	webHook := Webhook{
		verifyToken:    cfg.WhatsApp.WebHookToken,
		graphAPIToken:  cfg.WhatsApp.GraphToken,
		webHookAddredd: cfg.WhatsApp.WebHookAddress,
	}

	bot := &Bot{
		cmds:           cmds,
		engine:         botEngine,
		cfg:            cfg,
		server:         server,
		ctx:            ctx,
		cancel:         cancel,
		target:         cfg.BotName,
		sessionManager: sessionManager,
		hook:           webHook,
	}
	go bot.sessionManager.RemoveExpiredSessions()

	// Webhook handlers
	server.HandleFunc("/webhook", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			bot.verificationHandler(w, r)
		} else if r.Method == http.MethodPost {
			bot.webhookHandler(w, r)
		} else {
			http.NotFound(w, r)
		}
	})

	server.HandleFunc("/", func(w http.ResponseWriter, _ *http.Request) {
		_, err := fmt.Fprint(w, "<pre>Nothing to see here. Checkout README.md to start.</pre>")
		if err != nil {
			log.Print(err)
		}
	})

	return bot, nil
}

func (bot *Bot) Start() error {
	server := &http.Server{
		Addr:         bot.hook.webHookAddredd,
		Handler:      bot.server,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	go func() {
		log.Printf("Server is listening on address: %s", bot.hook.webHookAddredd)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
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

// handleCommand executes a command with its arguments for the user.
// It combines the commands and arguments into a single string, calls the engine's Run method,
// clears the user's context, and sends the result back to the user.
func (bot *Bot) handleCommand(commands []string) []byte {
	fullCommand := strings.Join(commands, " ")

	// Call the engine's Run method with the full command string
	res := bot.engine.ParseAndExecute(entity.PlatformIDTelegram, "", fullCommand)

	return []byte(res.Message)
}
