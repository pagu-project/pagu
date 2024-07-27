package engine

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/pagu-project/Pagu/config"
	"github.com/pagu-project/Pagu/internal/engine/command"
	"github.com/pagu-project/Pagu/internal/engine/command/calculator"
	"github.com/pagu-project/Pagu/internal/engine/command/market"
	"github.com/pagu-project/Pagu/internal/engine/command/network"
	phoenixtestnet "github.com/pagu-project/Pagu/internal/engine/command/phoenix"
	"github.com/pagu-project/Pagu/internal/engine/command/validator"
	"github.com/pagu-project/Pagu/internal/engine/command/voucher"
	"github.com/pagu-project/Pagu/internal/engine/command/zealy"
	"github.com/pagu-project/Pagu/internal/entity"
	"github.com/pagu-project/Pagu/internal/job"
	"github.com/pagu-project/Pagu/internal/repository"
	"github.com/pagu-project/Pagu/pkg/amount"
	"github.com/pagu-project/Pagu/pkg/cache"
	"github.com/pagu-project/Pagu/pkg/client"
	"github.com/pagu-project/Pagu/pkg/log"
	"github.com/pagu-project/Pagu/pkg/wallet"
)

type BotEngine struct {
	ctx    context.Context
	cancel context.CancelFunc

	clientMgr client.Manager
	db        repository.Database
	rootCmd   *command.Command

	// commands
	calculatorCmd *calculator.Calculator
	networkCmd    *network.Network
	phoenixCmd    *phoenixtestnet.Phoenix
	zealyCmd      *zealy.Zealy
	voucherCmd    *voucher.Voucher
	marketCmd     *market.Market
	validatorCmd  *validator.Validator
}

func NewBotEngine(cfg *config.Config) (*BotEngine, error) {
	ctx, cancel := context.WithCancel(context.Background())

	db, err := repository.NewDB(cfg.Database.URL)
	if err != nil {
		cancel()
		return nil, err
	}
	log.Info("database loaded successfully")

	cm := client.NewClientMgr(ctx)

	if cfg.LocalNode != "" {
		localClient, err := client.NewClient(cfg.LocalNode)
		if err != nil {
			cancel()
			return nil, err
		}

		cm.AddClient(localClient)
	}

	for _, nn := range cfg.NetworkNodes {
		c, err := client.NewClient(nn)
		if err != nil {
			log.Error("can't add new network node client", "err", err, "addr", nn)
		}
		cm.AddClient(c)
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

	return newBotEngine(ctx, cancel, db, cm, wlt, cfg.Phoenix.FaucetAmount), nil
}

func newBotEngine(ctx context.Context, cnl context.CancelFunc, db repository.Database, cm client.Manager,
	wlt wallet.IWallet, phoenixFaucetAmount amount.Amount,
) *BotEngine {
	rootCmd := &command.Command{
		Emoji:       "🤖",
		Name:        "pagu",
		Help:        "Root Command",
		AppIDs:      entity.AllAppIDs(),
		SubCommands: make([]*command.Command, 0),
	}

	// price caching job
	priceCache := cache.NewBasic[string, entity.Price](0 * time.Second)
	priceJob := job.NewPrice(priceCache)
	priceJobSched := job.NewScheduler()
	priceJobSched.Submit(priceJob)
	go priceJobSched.Run()

	netCmd := network.NewNetwork(ctx, cm)
	calcCmd := calculator.NewCalculator(cm)
	ptCmd := phoenixtestnet.NewPhoenix(ctx, wlt, phoenixFaucetAmount, cm, db)
	zealyCmd := zealy.NewZealy(db, wlt)
	voucherCmd := voucher.NewVoucher(db, wlt, cm)
	marketCmd := market.NewMarket(cm, priceCache)
	validatorCmd := validator.NewValidator(db)

	return &BotEngine{
		ctx:           ctx,
		cancel:        cnl,
		clientMgr:     cm,
		db:            db,
		rootCmd:       rootCmd,
		networkCmd:    netCmd,
		calculatorCmd: calcCmd,
		phoenixCmd:    ptCmd,
		zealyCmd:      zealyCmd,
		voucherCmd:    voucherCmd,
		marketCmd:     marketCmd,
		validatorCmd:  validatorCmd,
	}
}

func (be *BotEngine) Commands() []*command.Command {
	return be.rootCmd.SubCommands
}

func (be *BotEngine) RegisterAllCommands() {
	be.rootCmd.AddSubCommand(be.calculatorCmd.GetCommand())
	be.rootCmd.AddSubCommand(be.networkCmd.GetCommand())
	be.rootCmd.AddSubCommand(be.zealyCmd.GetCommand())
	be.rootCmd.AddSubCommand(be.voucherCmd.GetCommand())
	be.rootCmd.AddSubCommand(be.marketCmd.GetCommand())
	be.rootCmd.AddSubCommand(be.phoenixCmd.GetCommand())
	be.rootCmd.AddSubCommand(be.validatorCmd.GetCommand())

	be.rootCmd.AddHelpSubCommand()
}

func (be *BotEngine) Run(
	appID entity.AppID,
	callerID string,
	commands []string,
	args map[string]string,
) command.CommandResult {
	log.Debug("run command", "callerID", callerID, "commands", args, "inputs", args)

	cmd := be.getTargetCommand(commands)
	if !cmd.HasAppID(appID) {
		return cmd.FailedResult("unauthorized appID: %v", appID)
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

	cs, err := be.clientMgr.GetCirculatingSupply()
	if err != nil {
		cs = 0
	}

	return &network.NetStatus{
		ConnectedPeersCount: netInfo.ConnectedPeersCount,
		ValidatorsCount:     chainInfo.TotalValidators,
		TotalBytesSent:      netInfo.TotalSentBytes,
		TotalBytesReceived:  netInfo.TotalReceivedBytes,
		CurrentBlockHeight:  chainInfo.LastBlockHeight,
		TotalNetworkPower:   chainInfo.TotalPower,
		TotalCommitteePower: chainInfo.CommitteePower,
		NetworkName:         netInfo.NetworkName,
		TotalAccounts:       chainInfo.TotalAccounts,
		CirculatingSupply:   cs,
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
