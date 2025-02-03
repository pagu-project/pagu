// Code generated by command-generator. DO NOT EDIT.
package phoenix

import (
	"github.com/pagu-project/pagu/internal/engine/command"
	"github.com/pagu-project/pagu/internal/entity"
)

const argNameFaucetAddress = "address"

type phoenixSubCmds struct {
	subCmdHealth *command.Command
	subCmdStatus *command.Command
	subCmdFaucet *command.Command
	subCmdWallet *command.Command
}

func (c *PhoenixCmd) buildSubCmds() *phoenixSubCmds {
	subCmdHealth := &command.Command{
		Name:           "health",
		Help:           "Check the network health status",
		Handler:        c.healthHandler,
		ResultTemplate: "Network is {{.Status}}\nCurrent Time: {{.CurrentTime}}\nLast Block Time: {{.LastBlockTime}}\nTime Difference: {{.TimeDiff}}\nLast Block Height: {{.LastBlockHeight}}\n",
		TargetBotIDs:   entity.AllBotIDs(),
	}
	subCmdStatus := &command.Command{
		Name:           "status",
		Help:           "View network statistics",
		Handler:        c.statusHandler,
		ResultTemplate: "Network Name: {{.NetworkName}}\nConnected Peers: {{.ConnectedPeers}}\nValidator Count: {{.ValidatorsCount}}\nAccount Count: {{.AccountsCount}}\nCurrent Block Height: {{.CurrentBlockHeight}}\nTotal Power: {{.TotalPower}} tPAC\nTotal Committee Power: {{.TotalCommitteePower}} tPAC\n\n> Note📝: This info is from a random network node. Some data may not be consistent.\n",
		TargetBotIDs:   entity.AllBotIDs(),
	}
	subCmdFaucet := &command.Command{
		Name:           "faucet",
		Help:           "Get tPAC test coins on Phoenix Testnet for testing your project",
		Handler:        c.faucetHandler,
		ResultTemplate: "You received {{.amount}} tPAC on Phoenix Testnet!\n\nhttps://phoenix.pacviewer.com/transaction/{{.txHash}}\n",
		TargetBotIDs:   entity.AllBotIDs(),
		Args: []*command.Args{
			{
				Name:     "address",
				Desc:     "Your testnet address [example: tpc1z...]",
				InputBox: command.InputBoxText,
				Optional: false,
			},
		},
	}
	subCmdWallet := &command.Command{
		Name:           "wallet",
		Help:           "Show the faucet wallet balance",
		Handler:        c.walletHandler,
		ResultTemplate: "Pagu Phoenix Wallet:\n\nAddress: {{.address}}\nBalance: {{.balance}}\n",
		TargetBotIDs:   entity.AllBotIDs(),
	}

	return &phoenixSubCmds{
		subCmdHealth: subCmdHealth,
		subCmdStatus: subCmdStatus,
		subCmdFaucet: subCmdFaucet,
		subCmdWallet: subCmdWallet,
	}
}

func (c *PhoenixCmd) buildPhoenixCommand() *command.Command {
	phoenixCmd := &command.Command{
		Emoji:        "🐦",
		Name:         "phoenix",
		Help:         "Commands for working with Phoenix Testnet",
		SubCommands:  make([]*command.Command, 0),
		TargetBotIDs: entity.AllBotIDs(),
	}

	c.phoenixSubCmds = c.buildSubCmds()

	phoenixCmd.AddSubCommand(c.subCmdHealth)
	phoenixCmd.AddSubCommand(c.subCmdStatus)
	phoenixCmd.AddSubCommand(c.subCmdFaucet)
	phoenixCmd.AddSubCommand(c.subCmdWallet)

	return phoenixCmd
}
