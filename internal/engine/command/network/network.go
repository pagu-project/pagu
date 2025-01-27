package network

import (
	"context"

	"github.com/pagu-project/pagu/internal/engine/command"
	"github.com/pagu-project/pagu/internal/entity"
	"github.com/pagu-project/pagu/pkg/client"
)

type NetworkCmd struct {
	*networkSubCmds

	ctx       context.Context
	clientMgr client.IManager
}

func NewNetworkCmd(ctx context.Context, clientMgr client.IManager) *NetworkCmd {
	return &NetworkCmd{
		ctx:       ctx,
		clientMgr: clientMgr,
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

func (n *NetworkCmd) GetCommand() *command.Command {
	cmd := n.buildNetworkCommand()
	cmd.PlatformIDs = entity.AllPlatformIDs()
	cmd.TargetFlag = command.TargetMaskAll

	n.subCmdNodeInfo.PlatformIDs = entity.AllPlatformIDs()
	n.subCmdNodeInfo.TargetFlag = command.TargetMaskAll

	n.subCmdStatus.PlatformIDs = entity.AllPlatformIDs()
	n.subCmdStatus.TargetFlag = command.TargetMaskAll

	n.subCmdHealth.PlatformIDs = entity.AllPlatformIDs()
	n.subCmdHealth.TargetFlag = command.TargetMaskAll

	return cmd
}
