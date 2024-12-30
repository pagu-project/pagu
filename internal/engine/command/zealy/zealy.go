package zealy

import (
	"github.com/pagu-project/pagu/internal/engine/command"
	"github.com/pagu-project/pagu/internal/entity"
	"github.com/pagu-project/pagu/internal/repository"
	"github.com/pagu-project/pagu/pkg/wallet"
)

const (
	CommandName       = "zealy"
	ClaimCommandName  = "claim"
	StatusCommandName = "status"
	HelpCommandName   = "help"
)

type Zealy struct {
	db     repository.IDatabase
	wallet wallet.IWallet
}

func NewZealy(db repository.IDatabase, wlt wallet.IWallet) *Zealy {
	return &Zealy{
		db:     db,
		wallet: wlt,
	}
}

func (z *Zealy) GetCommand() *command.Command {
	subCmdClaim := &command.Command{
		Name: ClaimCommandName,
		Help: "Claim your Zealy Reward",
		Args: []command.Args{
			{
				Name:     "address",
				Desc:     "Your Pactus address",
				InputBox: command.InputBoxText,
				Optional: false,
			},
		},
		SubCommands: nil,
		AppIDs:      []entity.AppID{entity.AppIDDiscord},
		Handler:     z.claimHandler,
		TargetFlag:  command.TargetMaskMainnet,
	}

	subCmdStatus := &command.Command{
		Name:        StatusCommandName,
		Help:        "Status of Zealy reward claims",
		Args:        nil,
		SubCommands: nil,
		AppIDs:      []entity.AppID{entity.AppIDDiscord},
		Handler:     z.statusHandler,
		TargetFlag:  command.TargetMaskModerator,
	}

	cmdZealy := &command.Command{
		Name:        CommandName,
		Help:        "Zealy Commands",
		Args:        nil,
		AppIDs:      []entity.AppID{entity.AppIDDiscord},
		SubCommands: make([]*command.Command, 0),
		Handler:     nil,
		TargetFlag:  command.TargetMaskMainnet | command.TargetMaskModerator,
	}

	cmdZealy.AddSubCommand(subCmdClaim)
	cmdZealy.AddSubCommand(subCmdStatus)

	return cmdZealy
}
