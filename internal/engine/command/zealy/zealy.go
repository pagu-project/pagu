package zealy

import (
	"github.com/pagu-project/pagu/internal/engine/command"
	"github.com/pagu-project/pagu/internal/entity"
	"github.com/pagu-project/pagu/internal/repository"
	"github.com/pagu-project/pagu/pkg/wallet"
)

type ZealyCmd struct {
	db     *repository.Database
	wallet wallet.IWallet
}

func NewZealyCmd(db *repository.Database, wlt wallet.IWallet) *ZealyCmd {
	return &ZealyCmd{
		db:     db,
		wallet: wlt,
	}
}

func (z *ZealyCmd) GetCommand() *command.Command {
	subCmdClaim := &command.Command{
		Name:   "claim",
		Active: false,
		Help:   "Claim your Zealy reward",
		Args: []*command.Args{
			{
				Name:     "address",
				Desc:     "The Pactus address where the reward will be claimed",
				InputBox: command.InputBoxText,
				Optional: false,
			},
		},
		SubCommands:  nil,
		TargetBotIDs: []entity.BotID{entity.BotID_Discord},
		Handler:      z.claimHandler,
	}

	subCmdStatus := &command.Command{
		Name:         "status",
		Help:         "Check the status of Zealy reward claims",
		Args:         nil,
		SubCommands:  nil,
		TargetBotIDs: []entity.BotID{entity.BotID_Discord, entity.BotID_Moderator},
		Handler:      z.statusHandler,
	}

	cmdZealy := &command.Command{
		Name:         "zealy",
		Help:         "Commands for managing Zealy campaign",
		Args:         nil,
		TargetBotIDs: []entity.BotID{entity.BotID_Discord, entity.BotID_Moderator},
		SubCommands:  make([]*command.Command, 0),
		Handler:      nil,
	}

	cmdZealy.AddSubCommand(subCmdClaim)
	cmdZealy.AddSubCommand(subCmdStatus)

	return cmdZealy
}
