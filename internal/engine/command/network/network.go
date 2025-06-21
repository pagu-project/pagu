package network

import (
	"context"

	"github.com/pagu-project/pagu/internal/engine/command"
	"github.com/pagu-project/pagu/internal/entity"
	"github.com/pagu-project/pagu/pkg/client"
	"github.com/pagu-project/pagu/pkg/wallet"
)

type NetworkCmd struct {
	*networkSubCmds

	ctx       context.Context
	clientMgr client.IManager
	wallet    wallet.IWallet
}

func NewNetworkCmd(ctx context.Context, clientMgr client.IManager, wallet wallet.IWallet) *NetworkCmd {
	return &NetworkCmd{
		ctx:       ctx,
		clientMgr: clientMgr,
		wallet:    wallet,
	}
}

type NodeInfo struct {
	PeerID              string
	IPAddress           string
	Agent               string
	Moniker             string
	Country             string
	City                string
	RegionName          string
	TimeZone            string
	ISP                 string
	ValidatorNum        int32
	AvailabilityScore   float64
	StakeAmount         int64
	LastBondingHeight   uint32
	LastSortitionHeight uint32
}

type NetStatus struct {
	NetworkName         string
	ConnectedPeersCount uint32
	ValidatorsCount     int32
	TotalBytesSent      int64
	TotalBytesReceived  int64
	CurrentBlockHeight  uint32
	TotalNetworkPower   int64
	TotalCommitteePower int64
	TotalAccounts       int32
	CirculatingSupply   int64
}

func (n *NetworkCmd) BuildCommand(botID entity.BotID) *command.Command {
	cmd := n.buildNetworkCommand(botID)

	return cmd
}
