package network

import (
	"github.com/pagu-project/pagu/internal/engine/command"
	"github.com/pagu-project/pagu/internal/entity"
	"github.com/pagu-project/pagu/pkg/amount"
	"github.com/pagu-project/pagu/pkg/utils"
)

func (c *NetworkCmd) statusHandler(
	_ *entity.User,
	cmd *command.Command,
	_ map[string]string,
) command.CommandResult {
	netInfo, err := c.clientMgr.GetNetworkInfo()
	if err != nil {
		return cmd.ErrorResult(err)
	}

	chainInfo, err := c.clientMgr.GetBlockchainInfo()
	if err != nil {
		return cmd.ErrorResult(err)
	}

	supply := c.clientMgr.GetCirculatingSupply()

	// Convert NanoPAC to PAC using the Amount type.
	totalNetworkPower := amount.Amount(chainInfo.TotalPower).ToPAC()
	totalCommitteePower := amount.Amount(chainInfo.CommitteePower).ToPAC()
	circulatingSupply := amount.Amount(supply).ToPAC()

	net := NetStatus{
		ValidatorsCount:     chainInfo.TotalValidators,
		CurrentBlockHeight:  chainInfo.LastBlockHeight,
		TotalNetworkPower:   int64(totalNetworkPower),
		TotalCommitteePower: int64(totalCommitteePower),
		NetworkName:         netInfo.NetworkName,
		TotalAccounts:       chainInfo.TotalAccounts,
		CirculatingSupply:   int64(circulatingSupply),
	}

	return cmd.RenderResultTemplate(
		"NetworkName", net.NetworkName,
		"ConnectedPeers", utils.FormatNumber(int64(net.ConnectedPeersCount)),
		"ValidatorsCount", utils.FormatNumber(int64(net.ValidatorsCount)),
		"AccountsCount", utils.FormatNumber(int64(net.TotalAccounts)),
		"CurrentBlockHeight", utils.FormatNumber(int64(net.CurrentBlockHeight)),
		"TotalPower", utils.FormatNumber(net.TotalNetworkPower),
		"TotalCommitteePower", utils.FormatNumber(net.TotalCommitteePower),
		"CirculatingSupply", utils.FormatNumber(net.CirculatingSupply),
	)
}
