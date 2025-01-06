package network

import (
	"github.com/pagu-project/pagu/internal/engine/command"
	"github.com/pagu-project/pagu/internal/entity"
	"github.com/pagu-project/pagu/pkg/amount"
	"github.com/pagu-project/pagu/pkg/utils"
)

func (n *NetworkCmd) statusHandler(
	_ *entity.User,
	cmd *command.Command,
	_ map[string]string,
) command.CommandResult {
	netInfo, err := n.clientMgr.GetNetworkInfo()
	if err != nil {
		return cmd.ErrorResult(err)
	}

	chainInfo, err := n.clientMgr.GetBlockchainInfo()
	if err != nil {
		return cmd.ErrorResult(err)
	}

	supply := n.clientMgr.GetCirculatingSupply()

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

	return cmd.SuccessfulResultF("Network Name: %s\nConnected Peers: %v\n"+
		"Validators Count: %v\nAccounts Count: %v\nCurrent Block Height: %v\nTotal Power: %v PAC\n"+
		"Total Committee Power: %v PAC\nCirculating Supply: %v PAC\n"+
		"\n> Note📝: This info is from one random network node. Non-calculator data may not be consistent.",
		net.NetworkName,
		utils.FormatNumber(int64(net.ConnectedPeersCount)),
		utils.FormatNumber(int64(net.ValidatorsCount)),
		utils.FormatNumber(int64(net.TotalAccounts)),
		utils.FormatNumber(int64(net.CurrentBlockHeight)),
		utils.FormatNumber(net.TotalNetworkPower),
		utils.FormatNumber(net.TotalCommitteePower),
		utils.FormatNumber(net.CirculatingSupply),
	)
}
