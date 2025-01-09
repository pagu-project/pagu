package network

import (
	"context"

	"github.com/pagu-project/pagu/internal/engine/command"
	"github.com/pagu-project/pagu/internal/entity"
	"github.com/pagu-project/pagu/pkg/client"
)

type NetworkCmd struct {
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
	subCmdNodeInfo := &command.Command{
		Name: "node-info",
		Help: "View information about a specific node",
		Args: []*command.Args{
			{
				Name:     "validator_address",
				Desc:     "The validator address",
				InputBox: command.InputBoxText,
				Optional: false,
			},
		},
		SubCommands: nil,
		AppIDs:      entity.AllAppIDs(),
		Handler:     n.nodeInfoHandler,
		TargetFlag:  command.TargetMaskAll,
	}

	subCmdHealth := &command.Command{
		Name:        "health",
		Help:        "Check the network health status",
		Args:        []*command.Args{},
		SubCommands: nil,
		AppIDs:      entity.AllAppIDs(),
		Handler:     n.healthHandler,
		TargetFlag:  command.TargetMaskAll,
	}

	subCmdStatus := &command.Command{
		Name:        "status",
		Help:        "View network statistics",
		Args:        []*command.Args{},
		SubCommands: nil,
		AppIDs:      entity.AllAppIDs(),
		Handler:     n.statusHandler,
		TargetFlag:  command.TargetMaskAll,
	}

	cmdNetwork := &command.Command{
		Name:        "network",
		Help:        "Commands for network metrics and information",
		Args:        nil,
		AppIDs:      entity.AllAppIDs(),
		SubCommands: make([]*command.Command, 0),
		Handler:     nil,
		TargetFlag:  command.TargetMaskAll,
	}

	cmdNetwork.AddSubCommand(subCmdNodeInfo)
	cmdNetwork.AddSubCommand(subCmdHealth)
	cmdNetwork.AddSubCommand(subCmdStatus)

	return cmdNetwork
}
