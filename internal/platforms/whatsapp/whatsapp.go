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

	"github.com/pagu-project/pagu/internal/engine"
	"github.com/pagu-project/pagu/internal/engine/command"
	"github.com/pagu-project/pagu/internal/entity"
	"github.com/pagu-project/pagu/internal/platforms/whatsapp/session"
	"github.com/pagu-project/pagu/pkg/log"
	"github.com/pagu-project/pagu/pkg/markdown"
)

type Bot struct {
	ctx            context.Context
	botID          entity.BotID
	server         *http.ServeMux
	engine         *engine.BotEngine
	cmds           []*command.Command
	cfg            *Config
	markdown       markdown.Renderer
	sessionManager *session.SessionManager
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

func (bot *Bot) findCommandByName(name string) *command.Command {
	if name == bot.engine.RootCmd().Name {
		return bot.engine.RootCmd()
	}
	for _, cmd := range bot.cmds {
		if cmd.Name == name {
			return cmd
		}
		for _, subCmd := range cmd.SubCommands {
			if subCmd.Name == name {
				return subCmd
			}
		}
	}

	return nil
}

func (bot *Bot) renderTextResult(result, destination string) map[string]any {
	result = bot.markdown.Render(result)

	return map[string]any{
		"messaging_product": "whatsapp",
		"recipient_type":    "individual",
		"to":                destination,
		"type":              "text",
		"text": map[string]any{
			"body": result,
		},
	}
}

func (bot *Bot) renderTextError(err error, destination string) map[string]any {
	errorMessage := fmt.Sprintf("❌ An error occurred: %s", err)

	return bot.renderTextResult(errorMessage, destination)
}

func (bot *Bot) renderInteractivePage(cmd *command.Command, destination string) map[string]any {
	if !cmd.HasSubCommand() {
		return bot.renderTextError(fmt.Errorf("command %s has no subcommands", cmd.Name), destination)
	}

	rows := []any{}
	for _, subCmd := range cmd.SubCommands {
		rows = append(rows, map[string]any{
			"id":          fmt.Sprintf("cmd:%s", subCmd.Name),
			"title":       subCmd.NameWithEmoji(),
			"description": subCmd.Help,
		})
	}

	text := bot.markdown.Render(cmd.Help)

	return map[string]any{
		"messaging_product": "whatsapp",
		"recipient_type":    "individual",
		"to":                destination,
		"type":              "interactive",
		"interactive": map[string]any{
			"type": "list",
			"header": map[string]any{
				"type": "text",
				"text": cmd.NameWithEmoji(),
			},
			"body": map[string]any{
				"text": text,
			},
			"action": map[string]any{
				"button": "View Options",
				"sections": []any{
					map[string]any{
						"title": "Menu",
						"rows":  rows,
					},
				},
			},
		},
	}
}

func (bot *Bot) renderInteractiveChoices(argName string, choices []command.Choice, destination string) map[string]any {
	text := fmt.Sprintf("Select a `%s`:\n\n", argName)

	rows := []any{}
	for _, choice := range choices {
		text += fmt.Sprintf("- %s\n", choice.Desc)

		rows = append(rows, map[string]any{
			"id":          fmt.Sprintf("choice:%s", choice.Value),
			"title":       choice.Name,
			"description": bot.markdown.Render(choice.Desc),
		})
	}

	text = bot.markdown.Render(text)

	return map[string]any{
		"messaging_product": "whatsapp",
		"recipient_type":    "individual",
		"to":                destination,
		"type":              "interactive",
		"interactive": map[string]any{
			"type": "list",
			"body": map[string]any{
				"text": text,
			},
			"action": map[string]any{
				"button": "Choose a `package`",
				"sections": []any{
					map[string]any{
						"title": "Menu",
						"rows":  rows,
					},
				},
			},
		},
	}
}

//nolint:gocognit // high complexity
func (bot *Bot) webhookHandler(receviedData []byte) error {
	var resBody WebhookRequest

	if err := json.Unmarshal(receviedData, &resBody); err != nil {
		return fmt.Errorf("unable to unmarshal body: %w", err)
	}

	//nolint:nestif // high complexity
	if len(resBody.Entry) > 0 {
		for _, entry := range resBody.Entry {
			for _, change := range entry.Changes {
				if len(change.Value.Messages) == 0 {
					continue
				}
				message := change.Value.Messages[0]

				phoneNumberID := change.Value.Metadata.PhoneNumberID
				session := bot.sessionManager.GetSession(phoneNumberID)

				if session == nil {
					session = bot.sessionManager.OpenSession(phoneNumberID)
				}

				switch message.Type {
				case "interactive":
					log.Debug("Received interactive message", "message", message)

					// TODO: this part is not tested.
					msgID := message.Interactive.ListReply.ID
					if strings.HasPrefix(msgID, "cmd:") {
						cmdName := strings.TrimPrefix(msgID, "cmd:")
						lastCmd := session.GetLastCommand()
						if lastCmd == cmdName {
							log.Warn("Received repeated command", "cmdName", cmdName)
						} else {
							log.Debug("Add command", "cmdName", cmdName)
							session.AddCommand(cmdName)
						}
					} else if strings.HasPrefix(msgID, "choice:") {
						value := strings.TrimPrefix(msgID, "choice:")
						log.Debug("Add arg value", "value", value)
						session.AddArgValue(value)
					} else {
						log.Warn("Received unknown interactive ID format", "id", msgID)
					}
				case "text":
					log.Debug("Received text message", "message", message)
					msg := message.Text.Body

					if len(session.Args) > 0 {
						log.Debug("Add arg value", "value", msg)
						session.AddArgValue(msg)
					}
				default:
					log.Warn("Received unknown message type", "message", message)
				}

				err := bot.sendCommand(bot.ctx, phoneNumberID, message.From)
				if err != nil {
					return err
				}
			}
		}
	}

	return nil
}

func (bot *Bot) sendCommand(ctx context.Context, phoneNumberID, destination string) error {
	log.Debug("Sending command", "phoneNumberID", phoneNumberID, "destination", destination)

	var responseMessage map[string]any

	session := bot.sessionManager.GetSession(phoneNumberID)

	//nolint:nestif // high complexity
	if len(session.Commands) > 0 {
		lastCmd := session.Commands[len(session.Commands)-1]
		cmd := bot.findCommandByName(lastCmd)

		if cmd == nil {
			log.Warn("Command not found", "cmdName", lastCmd)
			bot.sessionManager.CloseSession(phoneNumberID)

			return nil
		}

		if cmd.HasSubCommand() {
			// Command has subcommands, show the menu
			responseMessage = bot.renderInteractivePage(cmd, destination)
		} else if len(cmd.Args) > 0 && len(session.Args) < len(cmd.Args) {
			arg := cmd.Args[len(session.Args)]
			argName := cmd.Args[len(session.Args)].Name

			if len(arg.Choices) > 0 {
				responseMessage = bot.renderInteractiveChoices(argName, arg.Choices, destination)
			} else {
				what := fmt.Sprintf("Enter `%s`:", argName)
				responseMessage = bot.renderTextResult(what, destination)
			}

			log.Debug("Add arg name", "name", argName)
			session.AddArgName(argName)
		} else {
			res := bot.executeCommand(session, destination)

			if res.Successful {
				responseMessage = bot.renderTextResult(res.Message, destination)
			} else {
				responseMessage = bot.renderTextError(errors.New(res.Message), destination)
			}

			// Close the session after executing the command
			bot.sessionManager.CloseSession(phoneNumberID)
		}
	} else {
		rootCmd := bot.engine.RootCmd()
		responseMessage = bot.renderInteractivePage(rootCmd, destination)

		session.AddCommand(rootCmd.Name)
	}

	url := fmt.Sprintf("https://graph.facebook.com/v23.0/%s/messages", phoneNumberID)

	jsonData, err := json.Marshal(responseMessage)
	if err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", "Bearer "+bot.cfg.AccessToken)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return err
	}
	defer func() {
		_ = res.Body.Close()
	}()

	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to send list message: %s", res.Status)
	}

	return nil
}

func NewWhatsAppBot(ctx context.Context, cfg *Config, botID entity.BotID, engine *engine.BotEngine) (*Bot, error) {
	server := http.NewServeMux()
	cmds := engine.Commands()

	sessionManager := session.NewSessionManager(ctx,
		time.Duration(cfg.Session.SessionTTL*int(time.Second)),
		time.Duration(cfg.Session.CheckInterval*int(time.Second)),
	)

	markdown := markdown.NewWhatsAppRenderer()

	bot := &Bot{
		cmds:           cmds,
		engine:         engine,
		cfg:            cfg,
		server:         server,
		ctx:            ctx,
		botID:          botID,
		sessionManager: sessionManager,
		markdown:       markdown,
	}
	go bot.sessionManager.RemoveExpiredSessions()

	server.HandleFunc(cfg.WebHookPath, func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			receviedData, err := io.ReadAll(r.Body)
			if err != nil {
				log.Warn("unable to read body", "error", err)
				w.WriteHeader(http.StatusBadRequest)

				return
			}
			defer func() {
				_ = r.Body.Close()
			}()

			err = bot.webhookHandler(receviedData)
			if err != nil {
				log.Error("Webhook handler error", "error", err)
				w.WriteHeader(http.StatusInternalServerError)

				return
			}

			w.WriteHeader(http.StatusOK)
		} else {
			log.Warn("Received unknown method", "method", r.Method)

			http.NotFound(w, r)
		}
	})

	return bot, nil
}

func (bot *Bot) Start() error {
	server := &http.Server{
		Addr:         bot.cfg.WebHookAddress,
		Handler:      bot.server,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	go func() {
		log.Info("Server is listening", "address", bot.cfg.WebHookAddress)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Warn("Error starting server", "error", err)
		}
	}()

	log.Info("Starting WhatsApp Bot...",
		"sessionTTL", bot.cfg.Session.SessionTTL,
		"checkInterval", bot.cfg.Session.CheckInterval)

	return nil
}

func (*Bot) Stop() {
	log.Info("Shutting down WhatsApp Bot")
}

// executeCommand executes a session with its commands and arguments for the user.
// It combines the commands and arguments into a single line, execute the command line
// and returns the result.
func (bot *Bot) executeCommand(session *session.Session, callerID string) command.CommandResult {
	commandLine := session.GetCommandLine()

	log.Debug("Executing command", "commandLine", commandLine)

	// Call the engine's Run method with the full command string
	return bot.engine.ParseAndExecute(entity.PlatformIDTelegram, callerID, commandLine)
}
