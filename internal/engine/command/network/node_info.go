package network

import (
	"fmt"

	"github.com/pagu-project/pagu/internal/engine/command"
	"github.com/pagu-project/pagu/internal/entity"
	"github.com/pagu-project/pagu/pkg/amount"
	utils2 "github.com/pagu-project/pagu/pkg/utils"
)

func (n *NetworkCmd) nodeInfoHandler(_ *entity.User,
	cmd *command.Command, args map[string]string,
) command.CommandResult {
	valAddress := args[argNameNodeInfoValidator_address]

	peerInfo, err := n.clientMgr.GetPeerInfo(valAddress)
	if err != nil {
		return cmd.RenderErrorTemplate(err)
	}

	ip := utils2.ExtractIPFromMultiAddr(peerInfo.Address)
	geoData := utils2.GetGeoIP(n.ctx, ip)

	nodeInfo := &NodeInfo{
		PeerID:     peerInfo.PeerId,
		IPAddress:  peerInfo.Address,
		Agent:      peerInfo.Agent,
		Moniker:    peerInfo.Moniker,
		Country:    geoData.CountryName,
		City:       geoData.City,
		RegionName: geoData.RegionName,
		TimeZone:   geoData.TimeZone,
		ISP:        geoData.ISP,
	}

	// here we check if the node is also a validator.
	// if its a validator , then we populate the validator data.
	// if not validator then we set everything to 0/empty .
	val, err := n.clientMgr.GetValidatorInfo(valAddress)
	if err == nil && val != nil {
		nodeInfo.ValidatorNum = val.Validator.Number
		nodeInfo.AvailabilityScore = val.Validator.AvailabilityScore
		// Convert NanoPAC to PAC using the Amount type and then to int64.
		stakeAmount := amount.Amount(val.Validator.Stake).ToPAC()
		nodeInfo.StakeAmount = int64(stakeAmount) // Convert float64 to int64.
		nodeInfo.LastBondingHeight = val.Validator.LastBondingHeight
		nodeInfo.LastSortitionHeight = val.Validator.LastSortitionHeight
	} else {
		nodeInfo.ValidatorNum = 0
		nodeInfo.AvailabilityScore = 0
		nodeInfo.StakeAmount = 0
		nodeInfo.LastBondingHeight = 0
		nodeInfo.LastSortitionHeight = 0
	}

	var pip19Score string
	if nodeInfo.AvailabilityScore >= 0.9 {
		pip19Score = fmt.Sprintf("%v✅", nodeInfo.AvailabilityScore)
	} else {
		pip19Score = fmt.Sprintf("%v⚠️", nodeInfo.AvailabilityScore)
	}

	return cmd.RenderResultTemplate(
		"PeerID", nodeInfo.PeerID,
		"IPAddress", nodeInfo.IPAddress,
		"Agent", nodeInfo.Agent,
		"Moniker", nodeInfo.Moniker,
		"Country", nodeInfo.Country,
		"City", nodeInfo.City,
		"RegionName", nodeInfo.RegionName,
		"TimeZone", nodeInfo.TimeZone,
		"ISP", nodeInfo.ISP,
		"Number", utils2.FormatNumber(int64(nodeInfo.ValidatorNum)),
		"AvailabilityScore", pip19Score,
		"Stake", utils2.FormatNumber(nodeInfo.StakeAmount),
	)
}
