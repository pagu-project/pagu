package engine

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/pagu-project/pagu/internal/engine/command"
	"github.com/pagu-project/pagu/internal/engine/command/calculator"
	"github.com/pagu-project/pagu/internal/engine/command/crowdfund"
	"github.com/pagu-project/pagu/internal/engine/command/market"
	"github.com/pagu-project/pagu/internal/engine/command/network"
	"github.com/pagu-project/pagu/internal/engine/command/phoenix"
	"github.com/pagu-project/pagu/internal/engine/command/voucher"
	"github.com/pagu-project/pagu/internal/entity"
	"github.com/pagu-project/pagu/internal/job"
	"github.com/pagu-project/pagu/internal/repository"
	"github.com/pagu-project/pagu/pkg/cache"
	"github.com/pagu-project/pagu/pkg/client"
	"github.com/pagu-project/pagu/pkg/log"
	"github.com/pagu-project/pagu/pkg/mailer"
	"github.com/pagu-project/pagu/pkg/nowpayments"
	"github.com/pagu-project/pagu/pkg/wallet"
)

type BotEngine struct {
	ctx       context.Context
	clientMgr client.IManager
	db        *repository.Database
	rootCmd   *command.Command
}

func NewBotEngine(ctx context.Context, botID entity.BotID, cfg *Config) (*BotEngine, error) {
	db, err := repository.NewDB(cfg.Database.URL)
	if err != nil {
		return nil, err
	}
	log.Info("database loaded successfully")

	mgr := client.NewClientMgr(ctx)
	if cfg.LocalNode != "" {
		localClient, err := client.NewClient(cfg.LocalNode)
		if err != nil {
			return nil, err
		}

		mgr.AddClient(localClient)
	}

	for _, nn := range cfg.NetworkNodes {
		client, err := client.NewClient(nn)
		if err != nil {
			log.Warn("error on adding new network client", "error", err, "addr", nn)
		}
		mgr.AddClient(client)
	}

	wlt, err := wallet.New(&cfg.Wallet)
	if err != nil {
		return nil, WalletError{
			Reason: err.Error(),
		}
	}
	log.Info("wallet opened successfully", "address", wlt.Address())

	mailer := mailer.NewSMTPMailer(&cfg.Mailer)

	nowPayments, err := nowpayments.NewNowPayments(ctx, &cfg.NowPayments)
	if err != nil {
		return nil, err
	}

	return newBotEngine(ctx, botID, cfg, db, mgr, wlt, mailer, nowPayments), nil
}

func newBotEngine(ctx context.Context,
	botID entity.BotID,
	cfg *Config,
	db *repository.Database,
	mgr client.IManager,
	wlt wallet.IWallet,
	mailer mailer.IMailer,
	nowPayments nowpayments.INowPayments,
) *BotEngine {
	// TODO: create an object and interface for me
	// price caching job
	priceCache := cache.NewBasic[string, entity.Price](10 * time.Second)
	priceJob := job.NewPrice(ctx, priceCache)
	priceJobSched := job.NewScheduler()
	priceJobSched.Submit(priceJob)
	go priceJobSched.Run()

	crowdfundCmd := crowdfund.NewCrowdfundCmd(ctx, db, wlt, nowPayments)
	calculatorCmd := calculator.NewCalculatorCmd(mgr)
	networkCmd := network.NewNetworkCmd(ctx, mgr, wlt)
	phoenixCmd := phoenix.NewPhoenixCmd(ctx, &cfg.Phoenix, db)
	voucherCmd := voucher.NewVoucherCmd(ctx, &cfg.Voucher, db, wlt, mgr, mailer)
	marketCmd := market.NewMarketCmd(mgr, priceCache)

	rootCmd := &command.Command{
		Emoji:        "🤖",
		Name:         "pagu",
		Help:         "Welcome to Pagu! Please select a command to start.",
		TargetBotIDs: entity.AllBotIDs(),
		SubCommands:  make([]*command.Command, 0),
	}

	subCommands := []*command.Command{
		crowdfundCmd.BuildCommand(botID),
		voucherCmd.BuildCommand(botID),
		calculatorCmd.BuildCommand(botID),
		networkCmd.BuildCommand(botID),
		phoenixCmd.BuildCommand(botID),
		marketCmd.BuildCommand(botID),
	}

	for _, cmd := range subCommands {
		if cmd.Active {
			rootCmd.AddSubCommand(botID, cmd)
		}
	}

	rootCmd.AddAboutSubCommand(botID)
	rootCmd.AddHelpSubCommand(botID)

	return &BotEngine{
		ctx:       ctx,
		clientMgr: mgr,
		db:        db,
		rootCmd:   rootCmd,
	}
}

func (be *BotEngine) Commands() []*command.Command {
	return be.rootCmd.SubCommands
}

// ParseAndExecute parses the input string and executes it.
// It returns an error if parsing fails or execution is unsuccessful.
func (be *BotEngine) ParseAndExecute(
	platformID entity.PlatformID,
	callerID string,
	input string,
) command.CommandResult {
	log.Debug("run command", "callerID", callerID, "input", input)

	var cmds []string
	var args map[string]string

	cmds, args, err := parseInput(input)
	if err != nil {
		return command.CommandResult{
			Message:    err.Error(),
			Successful: false,
		}
	}

	return be.executeCommand(platformID, callerID, cmds, args)
}

func parseCommandInput(cmdInput string) []string {
	cmds := make([]string, 0)

	tokens := strings.Split(cmdInput, " ")
	for _, token := range tokens {
		token = strings.TrimSpace(token)

		if token != "" {
			cmds = append(cmds, token)
		}
	}

	return cmds
}

func parseArgumentInput(argInput string) (map[string]string, error) {
	args := make(map[string]string)

	tokens := strings.Split(argInput, "--")
	for _, token := range tokens {
		token = strings.TrimSpace(token)

		if token != "" {
			parts := strings.SplitN(token, "=", 2)
			key := strings.TrimSpace(parts[0])

			if key == "" {
				return nil, fmt.Errorf("invalid argument format: %s", argInput)
			}

			if len(parts) == 1 {
				// Boolean argument
				args[key] = "true"
			} else {
				value := strings.TrimSpace(parts[1])
				value = strings.Trim(value, "\"'")
				value = strings.TrimSpace(value)

				args[key] = value
			}
		}
	}

	return args, nil
}

// parseInput parses the input string into commands and arguments.
// The input string should be in the following format:
// `command1 command2 --arg1=val1 --arg2=val2`
// It returns an error if parsing fails.
func parseInput(input string) ([]string, map[string]string, error) {
	// normalize input
	input = strings.ReplaceAll(input, "\t", " ")
	input = strings.TrimSpace(input)
	if input == "" {
		return nil, nil, errors.New("input string cannot be empty")
	}

	argIndex := strings.Index(input, "--")

	var cmdInput, argInput string
	if argIndex != -1 {
		cmdInput = input[:argIndex]
		argInput = input[argIndex:]
	} else {
		cmdInput = input
		argInput = ""
	}

	cmds := parseCommandInput(cmdInput)
	args, err := parseArgumentInput(argInput)
	if err != nil {
		return nil, nil, err
	}

	return cmds, args, nil
}

// executeCommand executes the parsed commands with their corresponding arguments.
// It returns an error if the execution fails.
func (be *BotEngine) executeCommand(
	platformID entity.PlatformID,
	callerID string,
	commands []string,
	args map[string]string,
) command.CommandResult {
	log.Debug("execute command", "callerID", callerID, "commands", commands, "args", args)

	cmd := be.getTargetCommand(commands)

	if cmd.Handler == nil {
		return cmd.RenderHelpTemplate()
	}

	caller, err := be.GetUser(platformID, callerID)
	if err != nil {
		log.Error("unable to GetUser", "error", err)

		return cmd.RenderErrorTemplate(fmt.Errorf("user is not defined in %s application", platformID))
	}

	return cmd.Handler(caller, cmd, args)
}

func (be *BotEngine) getTargetCommand(inCommands []string) *command.Command {
	targetCmd := be.rootCmd
	cmds := be.rootCmd.SubCommands

	for _, inCmd := range inCommands {
		found := false
		for _, cmd := range cmds {
			if cmd.Name != inCmd {
				continue
			}
			targetCmd = cmd
			if len(cmd.SubCommands) > 0 {
				cmds = cmd.SubCommands
				found = true

				break
			}
			found = true

			break
		}
		if !found {
			break
		}
	}

	return targetCmd
}

func (be *BotEngine) NetworkStatus() (*network.NetStatus, error) {
	netInfo, err := be.clientMgr.GetNetworkInfo()
	if err != nil {
		return nil, err
	}

	chainInfo, err := be.clientMgr.GetBlockchainInfo()
	if err != nil {
		return nil, err
	}

	supply := be.clientMgr.GetCirculatingSupply()

	return &network.NetStatus{
		ConnectedPeersCount: netInfo.ConnectedPeersCount,
		ValidatorsCount:     chainInfo.TotalValidators,
		TotalBytesSent:      int64(netInfo.MetricInfo.TotalSent.Bytes),
		TotalBytesReceived:  int64(netInfo.MetricInfo.TotalReceived.Bytes),
		CurrentBlockHeight:  chainInfo.LastBlockHeight,
		TotalNetworkPower:   chainInfo.TotalPower,
		TotalCommitteePower: chainInfo.CommitteePower,
		NetworkName:         netInfo.NetworkName,
		TotalAccounts:       chainInfo.TotalAccounts,
		CirculatingSupply:   supply,
	}, nil
}

func (be *BotEngine) GetUser(platformID entity.PlatformID, platformUserID string) (*entity.User, error) {
	existingUser, _ := be.db.GetUserByPlatformID(platformID, platformUserID)
	if existingUser != nil {
		return existingUser, nil
	}

	newUser := &entity.User{
		PlatformID:     platformID,
		PlatformUserID: platformUserID,
		Role:           entity.BasicUser,
	}
	if err := be.db.AddUser(newUser); err != nil {
		return nil, err
	}

	return newUser, nil
}

func (be *BotEngine) Stop() {
	be.clientMgr.Stop()
}

func (be *BotEngine) Start() {
	be.clientMgr.Start()
}

func (be *BotEngine) RootCmd() *command.Command {
	return be.rootCmd
}
