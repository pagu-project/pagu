package engine

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/pagu-project/Pagu/config"
	"github.com/pagu-project/Pagu/internal/engine/command"
	"github.com/pagu-project/Pagu/internal/engine/command/calculator"
	"github.com/pagu-project/Pagu/internal/engine/command/market"
	"github.com/pagu-project/Pagu/internal/engine/command/network"
	phoenixtestnet "github.com/pagu-project/Pagu/internal/engine/command/phoenix"
	"github.com/pagu-project/Pagu/internal/engine/command/voucher"
	"github.com/pagu-project/Pagu/internal/engine/command/zealy"
	"github.com/pagu-project/Pagu/internal/entity"
	"github.com/pagu-project/Pagu/internal/job"
	"github.com/pagu-project/Pagu/internal/repository"
	"github.com/pagu-project/Pagu/pkg/amount"
	"github.com/pagu-project/Pagu/pkg/cache"
	"github.com/pagu-project/Pagu/pkg/client"
	"github.com/pagu-project/Pagu/pkg/log"
	"github.com/pagu-project/Pagu/pkg/notification"
	"github.com/pagu-project/Pagu/pkg/notification/zoho"
	"github.com/pagu-project/Pagu/pkg/wallet"
)

type BotEngine struct {
	ctx    context.Context
	cancel context.CancelFunc

	clientMgr client.IManager
	db        repository.IDatabase
	rootCmd   *command.Command

	// commands
	calculatorCmd *calculator.Calculator
	networkCmd    *network.Network
	phoenixCmd    *phoenixtestnet.Phoenix
	voucherCmd    *voucher.Voucher
	marketCmd     *market.Market
	zealyCmd      *zealy.Zealy
}

func NewBotEngine(cfg *config.Config) (*BotEngine, error) {
	ctx, cancel := context.WithCancel(context.Background())

	db, err := repository.NewDB(cfg.Database.URL)
	if err != nil {
		cancel()

		return nil, err
	}
	log.Info("database loaded successfully")

	mgr := client.NewClientMgr(ctx)

	if cfg.LocalNode != "" {
		localClient, err := client.NewClient(cfg.LocalNode)
		if err != nil {
			cancel()

			return nil, err
		}

		mgr.AddClient(localClient)
	}

	for _, nn := range cfg.NetworkNodes {
		c, err := client.NewClient(nn)
		if err != nil {
			log.Error("can't add new network node client", "err", err, "addr", nn)
		}
		mgr.AddClient(c)
	}

	var wlt wallet.IWallet
	if cfg.Wallet.Enable {
		// load or create wallet.
		wlt, err = wallet.Open(cfg.Wallet)
		if err != nil {
			cancel()

			return nil, WalletError{
				Reason: err.Error(),
			}
		}

		log.Info("wallet opened successfully", "address", wlt.Address())
	}

	if cfg.BotName == config.BotNamePaguModerator {
		zapToMailConfig := zoho.ZapToMailerConfig{
			Host:     cfg.Notification.Zoho.Mail.Host,
			Port:     cfg.Notification.Zoho.Mail.Port,
			Username: cfg.Notification.Zoho.Mail.Username,
			Password: cfg.Notification.Zoho.Mail.Password,
		}
		mailSender, err := notification.New(notification.NotificationTypeMail, zapToMailConfig)
		if err != nil {
			cancel()

			return nil, err
		}

		// notification job
		mailSenderJob := job.NewMailSender(db, mailSender, cfg.Notification.Zoho.Mail.Templates)
		mailSenderSched := job.NewScheduler()
		mailSenderSched.Submit(mailSenderJob)
		go mailSenderSched.Run()
	}

	return newBotEngine(ctx, cancel, db, mgr, wlt, cfg.Phoenix.FaucetAmount), nil
}

func (be *BotEngine) Commands() []*command.Command {
	return be.rootCmd.SubCommands
}

func (be *BotEngine) RegisterAllCommands() {
	be.rootCmd.AddSubCommand(be.calculatorCmd.GetCommand())
	be.rootCmd.AddSubCommand(be.networkCmd.GetCommand())
	be.rootCmd.AddSubCommand(be.voucherCmd.GetCommand())
	be.rootCmd.AddSubCommand(be.marketCmd.GetCommand())
	be.rootCmd.AddSubCommand(be.zealyCmd.GetCommand())
	be.rootCmd.AddSubCommand(be.phoenixCmd.GetCommand())

	be.rootCmd.AddHelpSubCommand()
}

// ParseAndExecute parses the input string and executes it.
// It returns an error if parsing fails or execution is unsuccessful.
func (be *BotEngine) ParseAndExecute(
	appID entity.AppID,
	callerID string,
	input string,
) command.CommandResult {
	var cmds []string
	var args map[string]string

	cmds, args, err := parseCommand(input)
	if err != nil {
		return command.CommandResult{
			Message:    err.Error(),
			Successful: false,
		}
	}

	return be.executeCommand(appID, callerID, cmds, args)
}

// parseCommand parses the input string into commands and arguments.
// The input string should be in the following format:
// `command1 command2 --arg1=val1 --arg2=val2`
// It returns an error if parsing fails.
func parseCommand(input string) ([]string, map[string]string, error) {
	if strings.TrimSpace(input) == "" {
		return nil, nil, errors.New("input string cannot be empty")
	}

	// Split input by spaces while preserving argument values
	parts := strings.Fields(input)

	// Prepare results
	cmds := make([]string, 0)
	args := make(map[string]string)

	// Iterate over parts to separate commands and arguments
	for _, part := range parts {
		if strings.HasPrefix(part, "--") {
			// Argument: split on '='
			argParts := strings.SplitN(part, "=", 2)
			key := strings.TrimPrefix(argParts[0], "--")
			if len(argParts) != 2 || strings.TrimSpace(key) == "" || strings.TrimSpace(argParts[1]) == "" {
				return nil, nil, fmt.Errorf("invalid argument format: %s", part)
			}
			args[key] = argParts[1]
		} else {
			cmds = append(cmds, part)
		}
	}

	return cmds, args, nil
}

// executeCommand executes the parsed commands with their corresponding arguments.
// It returns an error if the execution fails.
func (be *BotEngine) executeCommand(
	appID entity.AppID,
	callerID string,
	commands []string,
	args map[string]string,
) command.CommandResult {
	log.Debug("run command", "callerID", callerID, "commands", args, "inputs", args)

	cmd := be.getTargetCommand(commands)
	if !cmd.HasAppID(appID) {
		return cmd.FailedResultF("unauthorized appID: %v", appID)
	}

	if cmd.Handler == nil {
		return cmd.HelpResult()
	}

	caller, err := be.GetUser(appID, callerID)
	if err != nil {
		log.Error(err.Error())

		return cmd.ErrorResult(fmt.Errorf("user is not defined in %s application", appID.String()))
	}

	for _, middlewareFunc := range cmd.Middlewares {
		if err := middlewareFunc(caller, cmd, args); err != nil {
			log.Error(err.Error())

			return cmd.ErrorResult(errors.New("command is not available. please try again later"))
		}
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

func (be *BotEngine) GetUser(appID entity.AppID, callerID string) (*entity.User, error) {
	if u, _ := be.db.GetUserByApp(appID, callerID); u != nil {
		return u, nil
	}

	user := &entity.User{ApplicationID: appID, CallerID: callerID}
	if err := be.db.AddUser(user); err != nil {
		return nil, err
	}

	return user, nil
}

func (be *BotEngine) Stop() {
	log.Info("Stopping the Bot Engine")

	be.cancel()
	be.clientMgr.Stop()
}

func (be *BotEngine) Start() {
	log.Info("Starting the Bot Engine")

	be.clientMgr.Start()
}

func newBotEngine(ctx context.Context,
	cancel context.CancelFunc,
	db repository.IDatabase,
	mgr client.IManager,
	wlt wallet.IWallet,
	phoenixFaucetAmount amount.Amount,
) *BotEngine {
	rootCmd := &command.Command{
		Emoji:       "🤖",
		Name:        "pagu",
		Help:        "Root Command",
		AppIDs:      entity.AllAppIDs(),
		SubCommands: make([]*command.Command, 0),
	}

	// price caching job
	priceCache := cache.NewBasic[string, entity.Price](10 * time.Second)
	priceJob := job.NewPrice(priceCache)
	priceJobSched := job.NewScheduler()
	priceJobSched.Submit(priceJob)
	go priceJobSched.Run()

	netCmd := network.NewNetwork(ctx, mgr)
	calcCmd := calculator.NewCalculator(mgr)
	phoenixCmd := phoenixtestnet.NewPhoenix(ctx, wlt, phoenixFaucetAmount, mgr, db)
	voucherCmd := voucher.NewVoucher(db, wlt, mgr)
	marketCmd := market.NewMarket(mgr, priceCache)
	zealyCmd := zealy.NewZealy(db, wlt)

	return &BotEngine{
		ctx:           ctx,
		cancel:        cancel,
		clientMgr:     mgr,
		db:            db,
		rootCmd:       rootCmd,
		networkCmd:    netCmd,
		calculatorCmd: calcCmd,
		phoenixCmd:    phoenixCmd,
		voucherCmd:    voucherCmd,
		marketCmd:     marketCmd,
		zealyCmd:      zealyCmd,
	}
}
