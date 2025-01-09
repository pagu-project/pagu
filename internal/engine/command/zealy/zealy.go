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
		Name: "claim",
		Help: "Claim your Zealy reward",
		Args: []*command.Args{
			{
				Name:     "address",
				Desc:     "The Pactus address where the reward will be claimed",
				InputBox: command.InputBoxText,
				Optional: false,
			},
		},
		SubCommands: nil,
		AppIDs:      []entity.PlatformID{entity.PlatformIDDiscord},
		Handler:     z.claimHandler,
		TargetFlag:  command.TargetMaskMainnet,
	}

	subCmdStatus := &command.Command{
		Name:        "status",
		Help:        "Check the status of Zealy reward claims",
		Args:        nil,
		SubCommands: nil,
		AppIDs:      []entity.PlatformID{entity.PlatformIDDiscord},
		Handler:     z.statusHandler,
		TargetFlag:  command.TargetMaskModerator,
	}

	cmdZealy := &command.Command{
		Name:        "zealy",
		Help:        "Commands for managing Zealy campaign",
		Args:        nil,
		AppIDs:      []entity.PlatformID{entity.PlatformIDDiscord},
		SubCommands: make([]*command.Command, 0),
		Handler:     nil,
		TargetFlag:  command.TargetMaskMainnet | command.TargetMaskModerator,
	}

	cmdZealy.AddSubCommand(subCmdClaim)
	cmdZealy.AddSubCommand(subCmdStatus)

	return cmdZealy
}
