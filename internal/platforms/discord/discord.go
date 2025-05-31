package discord

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/pagu-project/pagu/internal/engine"
	"github.com/pagu-project/pagu/internal/engine/command"
	"github.com/pagu-project/pagu/internal/entity"
	"github.com/pagu-project/pagu/pkg/amount"
	"github.com/pagu-project/pagu/pkg/color"
	"github.com/pagu-project/pagu/pkg/log"
	"github.com/pagu-project/pagu/pkg/utils"
)

type Bot struct {
	ctx     context.Context
	botID   entity.BotID
	cfg     *Config
	Session *discordgo.Session
	engine  *engine.BotEngine
}

func NewDiscordBot(ctx context.Context, cfg *Config, botID entity.BotID, engine *engine.BotEngine) (*Bot, error) {
	session, err := discordgo.New("Bot " + cfg.Token)
	if err != nil {
		return nil, err
	}

	return &Bot{
		ctx:     ctx,
		Session: session,
		engine:  engine,
		cfg:     cfg,
		botID:   botID,
	}, nil
}

func (bot *Bot) Start() error {
	log.Info("starting Discord Bot...")

	err := bot.Session.Open()
	if err != nil {
		return err
	}

	bot.deleteAllCommands()

	return bot.registerCommands()
}

func (bot *Bot) Stop() {
	log.Info("Stopping Discord Bot")

	_ = bot.Session.Close()
}

func (bot *Bot) deleteAllCommands() {
	cmdsServer, _ := bot.Session.ApplicationCommands(bot.Session.State.User.ID, bot.cfg.GuildID)
	cmdsGlobal, _ := bot.Session.ApplicationCommands(bot.Session.State.User.ID, "")

	allCmds := []*discordgo.ApplicationCommand{}
	allCmds = append(allCmds, cmdsServer...)
	allCmds = append(allCmds, cmdsGlobal...)

	for _, cmd := range allCmds {
		err := bot.Session.ApplicationCommandDelete(cmd.ApplicationID, cmd.GuildID, cmd.ID)
		if err != nil {
			log.Error("unable to delete command", "error", err, "cmd", cmd.Name)
		} else {
			log.Info("discord command unregistered", "name", cmd.Name)
		}
	}
}

func (bot *Bot) registerCommands() error {
	bot.Session.AddHandler(func(session *discordgo.Session, interaction *discordgo.InteractionCreate) {
		bot.commandHandler(session, interaction)
	})

	cmds := bot.engine.Commands()
	for _, cmd := range cmds {
		if !cmd.CanBeHandledByBot(bot.botID) {
			continue
		}

		log.Info("registering new command", "name", cmd.Name)

		discordCmd := discordgo.ApplicationCommand{
			Type:        discordgo.ChatApplicationCommand,
			Name:        cmd.Name,
			Description: cmd.Help,
		}

		if cmd.HasSubCommand() {
			for _, subCmd := range cmd.SubCommands {
				if !subCmd.CanBeHandledByBot(bot.botID) {
					continue
				}

				log.Info("adding sub-command", "command", cmd.Name, "sub-command")

				discordSubCmd := &discordgo.ApplicationCommandOption{
					Type:        discordgo.ApplicationCommandOptionSubCommand,
					Name:        subCmd.Name,
					Description: subCmd.Help,
				}

				for _, arg := range subCmd.Args {
					log.Info("adding sub command argument", "command", cmd.Name, "sub-command", subCmd.Name, "argument", arg.Name)

					opt := &discordgo.ApplicationCommandOption{
						Type:        discordOptionType(arg.InputBox),
						Name:        arg.Name,
						Description: arg.Desc,
						Required:    !arg.Optional,
					}

					if len(arg.Choices) > 0 {
						opt.Choices = make([]*discordgo.ApplicationCommandOptionChoice, len(arg.Choices))

						for i, choice := range arg.Choices {
							opt.Choices[i] = &discordgo.ApplicationCommandOptionChoice{
								Name:  choice.Desc,
								Value: choice.Value,
							}
						}
					}

					discordSubCmd.Options = append(discordSubCmd.Options, opt)
				}

				discordCmd.Options = append(discordCmd.Options, discordSubCmd)
			}
		}

		cmd, err := bot.Session.ApplicationCommandCreate(bot.Session.State.User.ID, "", &discordCmd)
		if err != nil {
			log.Error("can not register discord command", "name", discordCmd.Name, "error", err)

			return err
		}
		log.Info("discord command registered", "name", cmd.Name)
	}

	return nil
}

func (bot *Bot) commandHandler(session *discordgo.Session, interaction *discordgo.InteractionCreate) {
	var inputBuilder strings.Builder
	args := make(map[string]string)

	// Get the application command data
	discordCmd := interaction.ApplicationCommandData()

	inputBuilder.WriteString(discordCmd.Name)

	for _, opt := range discordCmd.Options {
		if opt.Type == discordgo.ApplicationCommandOptionSubCommand {
			inputBuilder.WriteString(" ")
			inputBuilder.WriteString(opt.Name)

			for _, o := range opt.Options {
				args = parseArgs(&discordCmd, o, args)
			}
		}
	}

	for k, v := range args {
		inputBuilder.WriteString(fmt.Sprintf(" --%s=%s", k, v))
	}

	var callerID string
	if interaction.Member != nil {
		callerID = interaction.Member.User.ID
	} else if interaction.User != nil {
		callerID = interaction.User.ID
	} else {
		log.Warn("unable to obtain the callerID", "input", inputBuilder.String())

		return
	}

	res := bot.engine.ParseAndExecute(entity.PlatformIDDiscord, callerID, inputBuilder.String())
	bot.respondResultMsg(res, session, interaction)
}

func parseArgs(
	rootCmd *discordgo.ApplicationCommandInteractionData,
	opt *discordgo.ApplicationCommandInteractionDataOption,
	result map[string]string,
) map[string]string {
	switch opt.Type {
	case discordgo.ApplicationCommandOptionString:
		result[opt.Name] = opt.StringValue()

	case discordgo.ApplicationCommandOptionInteger:
		result[opt.Name] = strconv.Itoa(int(opt.IntValue()))

	case discordgo.ApplicationCommandOptionNumber:
		v := strconv.FormatFloat(opt.FloatValue(), 'f', 10, 64)
		result[opt.Name] = v

	case discordgo.ApplicationCommandOptionBoolean:
		result[opt.Name] = strconv.FormatBool(opt.BoolValue())

	case discordgo.ApplicationCommandOptionAttachment:
		// TODO: handle multiple attachment
		for _, attachment := range rootCmd.Resolved.Attachments {
			result[opt.Name] = attachment.URL
		}

	case discordgo.ApplicationCommandOptionSubCommand,
		discordgo.ApplicationCommandOptionSubCommandGroup,
		discordgo.ApplicationCommandOptionUser,
		discordgo.ApplicationCommandOptionChannel,
		discordgo.ApplicationCommandOptionRole,
		discordgo.ApplicationCommandOptionMentionable:

		log.Warn("received unhandled option type", "type", opt.Type)
	}

	return result
}

func (bot *Bot) respondResultMsg(res command.CommandResult,
	session *discordgo.Session, interaction *discordgo.InteractionCreate,
) {
	var resEmbed *discordgo.MessageEmbed
	if res.Successful {
		resEmbed = &discordgo.MessageEmbed{
			Title:       res.Title,
			Description: res.Message,
			Color:       color.Green.ToInt(),
		}
	} else {
		resEmbed = &discordgo.MessageEmbed{
			Title:       res.Title,
			Description: res.Message,
			Color:       color.Yellow.ToInt(),
		}
	}

	bot.respondEmbed(resEmbed, session, interaction)
}

func (*Bot) respondEmbed(embed *discordgo.MessageEmbed,
	session *discordgo.Session, interaction *discordgo.InteractionCreate,
) {
	response := &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{embed},
		},
	}

	err := session.InteractionRespond(interaction.Interaction, response)
	if err != nil {
		log.Error("InteractionRespond error:", "error", err)
	}
}

func (bot *Bot) UpdateStatusInfo() {
	// TODO: fix me!
	log.Info("info status started")
	for {
		status, err := bot.engine.NetworkStatus()
		if err != nil {
			continue
		}

		err = bot.Session.UpdateStatusComplex(newStatus("validators count",
			utils.FormatNumber(int64(status.ValidatorsCount))))
		if err != nil {
			log.Error("can't set status", "err", err)

			continue
		}

		time.Sleep(time.Second * 5)

		err = bot.Session.UpdateStatusComplex(newStatus("total accounts",
			utils.FormatNumber(int64(status.TotalAccounts))))
		if err != nil {
			log.Error("can't set status", "err", err)

			continue
		}

		time.Sleep(time.Second * 5)

		err = bot.Session.UpdateStatusComplex(newStatus("height", utils.FormatNumber(int64(status.CurrentBlockHeight))))
		if err != nil {
			log.Error("can't set status", "err", err)

			continue
		}

		time.Sleep(time.Second * 5)

		circulatingSupplyAmount := amount.Amount(status.CirculatingSupply)
		formattedCirculatingSupply := circulatingSupplyAmount.Format(amount.UnitPAC) + " PAC"

		err = bot.Session.UpdateStatusComplex(newStatus("circ supply", formattedCirculatingSupply))
		if err != nil {
			log.Error("can't set status", "err", err)

			continue
		}

		time.Sleep(time.Second * 5)

		totalNetworkPowerAmount := amount.Amount(status.TotalNetworkPower)
		formattedTotalNetworkPower := totalNetworkPowerAmount.Format(amount.UnitPAC) + " PAC"

		err = bot.Session.UpdateStatusComplex(newStatus("total power", formattedTotalNetworkPower))
		if err != nil {
			log.Error("can't set status", "err", err)

			continue
		}

		time.Sleep(time.Second * 5)
	}
}

func discordOptionType(inputBox command.InputBox) discordgo.ApplicationCommandOptionType {
	switch inputBox {
	case command.InputBoxText,
		command.InputBoxMultilineText:
		return discordgo.ApplicationCommandOptionString
	case command.InputBoxInteger:
		return discordgo.ApplicationCommandOptionInteger
	case command.InputBoxFloat:
		return discordgo.ApplicationCommandOptionNumber
	case command.InputBoxFile:
		return discordgo.ApplicationCommandOptionAttachment
	case command.InputBoxToggle:
		return discordgo.ApplicationCommandOptionBoolean
	case command.InputBoxChoice:
		return discordgo.ApplicationCommandOptionString
	default:
		return discordgo.ApplicationCommandOptionString
	}
}
