package engine

import (
	"fmt"
	"strconv"

	"github.com/kehiy/RoboPac/database"
	"github.com/kehiy/RoboPac/utils"
)

const (
	P2PCommandName            = "p2p-market"
	DepositAddressCommandName = "deposit-address"
	CreateOfferCommandName    = "create-offer"
)

func (be *BotEngine) RegisterP2PMarketCommands() {
	cmdDepositAddress := Command{
		Name:    DepositAddressCommandName,
		Desc:    "create a deposit address for P2P offer",
		Help:    "it will show your address if you already have an deposit address",
		Args:    []Args{},
		AppIDs:  []AppID{AppIdCLI, AppIdDiscord},
		Handler: be.depositAddressHandler,
	}

	cmdCreateOffer := Command{
		Name: CreateOfferCommandName,
		Desc: "create an offer for P2P market",
		Help: "",
		Args: []Args{
			{
				Name:     "total-amount",
				Desc:     "total amount of PAC",
				Optional: false,
			},
			{
				Name:     "total-price",
				Desc:     "total price which includes gas fee",
				Optional: false,
			},
			{
				Name:     "chain-type",
				Desc:     "e.g. BTCUSDT",
				Optional: false,
			},
			{
				Name:     "address",
				Desc:     "",
				Optional: false,
			},
		},
		AppIDs:  []AppID{AppIdCLI, AppIdDiscord},
		Handler: be.createOfferHandler,
	}

	cmdP2PMarket := Command{
		Name:        P2PCommandName,
		Desc:        "person to person market for pactus trading",
		Help:        "",
		Args:        nil,
		AppIDs:      []AppID{AppIdCLI, AppIdDiscord},
		SubCommands: []*Command{&cmdCreateOffer, &cmdDepositAddress},
		Handler:     nil,
	}

	be.Cmds = append(be.Cmds, cmdP2PMarket)
}

func (be *BotEngine) depositAddressHandler(_ AppID, callerID string, _ ...string) (*CommandResult, error) {
	u, err := be.db.GetUser(callerID)
	if err == nil {
		return MakeSuccessfulResult(
			"You already have a deposit address: %s", u.DepositAddress,
		), nil
	}

	addr, err := be.wallet.NewAddress(fmt.Sprintf("deposit address for %s", callerID))
	if err != nil {
		return MakeFailedResult(
			"can't make a new address: %v", err,
		), nil
	}

	err = be.db.AddUser(
		&database.DiscordUser{
			DiscordID:      callerID,
			DepositAddress: addr,
		},
	)
	if err != nil {
		return MakeFailedResult(
			"can't add discord user to database: %v", err,
		), nil
	}

	return MakeSuccessfulResult(
		"Deposit address crated for you successfully: %s", addr,
	), nil
}

func (be *BotEngine) createOfferHandler(source AppID, callerID string, args ...string) (*CommandResult, error) {
	u, err := be.db.GetUser(callerID)
	if err != nil {
		return nil, err
	}

	totalAmount, err := strconv.Atoi(args[0])
	if err != nil {
		return nil, err
	}

	totalPrice, err := strconv.Atoi(args[1])
	if err != nil {
		return nil, err
	}

	chainType := args[2]
	address := args[3]

	uBalance, err := be.clientMgr.GetBalance(u.DepositAddress)
	if err != nil {
		return nil, err
	}

	if float64(totalAmount) != utils.ChangeToCoin(uBalance) {
		return nil, fmt.Errorf("the deposit balance: %d is not equal to offered amount: %d",
			uBalance, totalAmount)
	}

	unitPrice := float64(totalPrice / totalAmount)

	offer := &database.Offer{
		TotalAmount: int64(totalAmount),
		TotalPrice:  int64(totalPrice),
		UnitPrice:   unitPrice,
		ChainType:   chainType,
		Address:     address,
		DiscordUser: *u,
	}

	if err = be.db.CreateOffer(offer); err != nil {
		return nil, err
	}

	return MakeSuccessfulResult(
		"Offer successfully created, your offer ID: %s", "TODO!!!!!!!",
	), nil
}
