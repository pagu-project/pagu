package discord

import (
	"fmt"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/pactus-project/pactus/types/amount"
	"github.com/pactus-project/pactus/util"
	"github.com/pagu-project/Pagu/config"
	"github.com/pagu-project/Pagu/internal/engine"
	"github.com/pagu-project/Pagu/internal/engine/command"
	"github.com/pagu-project/Pagu/internal/entity"
	"github.com/pagu-project/Pagu/pkg/color"
	"github.com/pagu-project/Pagu/pkg/log"
	"github.com/pagu-project/Pagu/pkg/utils"
)

type Bot struct {
	cfg      *config.DiscordBot
	Session  *discordgo.Session
	engine   *engine.BotEngine
	botColor color.ColorCode
	target   string
}

func NewDiscordBot(botEngine *engine.BotEngine, cfg *config.DiscordBot, target string) (*Bot, error) {
	session, err := discordgo.New("Bot " + cfg.Token)
	if err != nil {
		return nil, err
	}

	return &Bot{
		Session: session,
		engine:  botEngine,
		cfg:     cfg,
		target:  target,
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

func (bot *Bot) Stop() error {
	log.Info("Stopping Discord Bot")

	return bot.Session.Close()
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

//nolint:gocognit // Complexity cannot be reduced
func (bot *Bot) registerCommands() error {
	bot.Session.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		bot.commandHandler(s, i)
	})

	beCmds := bot.engine.Commands()
	for i, beCmd := range beCmds {
		if !beCmd.HasAppID(entity.AppIDDiscord) {
			continue
		}

		switch bot.target {
		case config.BotNamePaguMainnet,
			config.BotNamePaguStaging:
			if !util.IsFlagSet(beCmd.TargetFlag, command.TargetMaskMainnet) {
				continue
			}

		case config.BotNamePaguTestnet:
			if !util.IsFlagSet(beCmd.TargetFlag, command.TargetMaskTestnet) {
				continue
			}

		case config.BotNamePaguModerator:
			if !util.IsFlagSet(beCmd.TargetFlag, command.TargetMaskModerator) {
				continue
			}
		}

		log.Info("registering new command", "name", beCmd.Name, "desc", beCmd.Help, "index", i, "object", beCmd)

		discordCmd := discordgo.ApplicationCommand{
			Type:        discordgo.ChatApplicationCommand,
			Name:        beCmd.Name,
			Description: beCmd.Help,
		}

		if beCmd.HasSubCommand() {
			for _, sCmd := range beCmd.SubCommands {
				switch bot.target {
				case config.BotNamePaguMainnet:
					if !util.IsFlagSet(sCmd.TargetFlag, command.TargetMaskMainnet) {
						continue
					}

				case config.BotNamePaguTestnet:
					if !util.IsFlagSet(sCmd.TargetFlag, command.TargetMaskTestnet) {
						continue
					}

				case config.BotNamePaguModerator:
					if !util.IsFlagSet(sCmd.TargetFlag, command.TargetMaskModerator) {
						continue
					}
				}

				log.Info("adding command sub-command", "command", beCmd.Name,
					"sub-command", sCmd.Name, "desc", sCmd.Help)

				subCmd := &discordgo.ApplicationCommandOption{
					Type:        discordgo.ApplicationCommandOptionSubCommand,
					Name:        sCmd.Name,
					Description: sCmd.Help,
				}

				for _, arg := range sCmd.Args {
					if arg.Desc == "" || arg.Name == "" {
						continue
					}

					log.Info("adding sub command argument", "command", beCmd.Name,
						"sub-command", sCmd.Name, "argument", arg.Name, "desc", arg.Desc)

					subCmd.Options = append(subCmd.Options, &discordgo.ApplicationCommandOption{
						Type:        setCommandArgType(arg.InputBox.Int()),
						Name:        arg.Name,
						Description: arg.Desc,
						Required:    !arg.Optional,
					})
				}

				discordCmd.Options = append(discordCmd.Options, subCmd)
			}
		} else {
			for _, arg := range beCmd.Args {
				if arg.Desc == "" || arg.Name == "" {
					continue
				}

				log.Info("adding command argument", "command", beCmd.Name,
					"argument", arg.Name, "desc", arg.Desc)

				discordCmd.Options = append(discordCmd.Options, &discordgo.ApplicationCommandOption{
					Type:        setCommandArgType(arg.InputBox.Int()),
					Name:        arg.Name,
					Description: arg.Desc,
					Required:    !arg.Optional,
				})
			}
		}

		cmd, err := bot.Session.ApplicationCommandCreate(bot.Session.State.User.ID, bot.cfg.GuildID, &discordCmd)
		if err != nil {
			log.Error("can not register discord command", "name", discordCmd.Name, "error", err)

			return err
		}
		log.Info("discord command registered", "name", cmd.Name)
	}

	return nil
}

func (bot *Bot) commandHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if i.GuildID != bot.cfg.GuildID {
		bot.respondErrMsg("Please send messages on server chat", s, i)

		return
	}

	var inputBuilder strings.Builder

	// Get the application command data
	discordCmd := i.ApplicationCommandData()
	inputBuilder.WriteString("/")
	inputBuilder.WriteString(discordCmd.Name)

	for _, opt := range discordCmd.Options {
		if opt.Type == discordgo.ApplicationCommandOptionSubCommand {
			inputBuilder.WriteString(" ")
			inputBuilder.WriteString(discordCmd.Name)
			for _, o := range opt.Options {
				inputBuilder.WriteString(fmt.Sprintf(" --%s=%s", o.Name, o.Value))
			}
		}
	}

	res := bot.engine.ParseAndExecute(entity.AppIDDiscord, i.Member.User.ID, inputBuilder.String())
	bot.respondResultMsg(res, s, i)
}

func (bot *Bot) respondErrMsg(errStr string, s *discordgo.Session, i *discordgo.InteractionCreate) {
	errorEmbed := &discordgo.MessageEmbed{
		Title:       "Error",
		Description: errStr,
		Color:       color.Yellow.ToInt(),
	}
	bot.respondEmbed(errorEmbed, s, i)
}

func (bot *Bot) respondResultMsg(res command.CommandResult, s *discordgo.Session, i *discordgo.InteractionCreate) {
	var resEmbed *discordgo.MessageEmbed
	if res.Successful {
		resEmbed = &discordgo.MessageEmbed{
			Title:       "Successful",
			Description: res.Message,
			Color:       color.Green.ToInt(),
		}
	} else {
		resEmbed = &discordgo.MessageEmbed{
			Title:       "Failed",
			Description: res.Message,
			Color:       color.Yellow.ToInt(),
		}
	}

	bot.respondEmbed(resEmbed, s, i)
}

func (*Bot) respondEmbed(embed *discordgo.MessageEmbed, s *discordgo.Session, i *discordgo.InteractionCreate) {
	response := &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{embed},
		},
	}

	err := s.InteractionRespond(i.Interaction, response)
	if err != nil {
		log.Error("InteractionRespond error:", "error", err)
	}
}

func (bot *Bot) UpdateStatusInfo() {
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

func setCommandArgType(inputBox int) discordgo.ApplicationCommandOptionType {
	switch inputBox {
	case 0:
		return discordgo.ApplicationCommandOptionString
	case 1:
		return discordgo.ApplicationCommandOptionInteger
	case 2:
		return discordgo.ApplicationCommandOptionAttachment
	case 3:
		return discordgo.ApplicationCommandOptionNumber
	case 4:
		return discordgo.ApplicationCommandOptionBoolean
	default:
		return discordgo.ApplicationCommandOptionString
	}
}
