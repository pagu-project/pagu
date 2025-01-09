package voucher

import (
	"github.com/pagu-project/pagu/internal/engine/command"
	"github.com/pagu-project/pagu/internal/entity"
	"github.com/pagu-project/pagu/internal/repository"
	"github.com/pagu-project/pagu/pkg/client"
	"github.com/pagu-project/pagu/pkg/wallet"
)

type VoucherCmd struct {
	db            *repository.Database
	wallet        wallet.IWallet
	clientManager client.IManager
}

func NewVoucherCmd(db *repository.Database, wlt wallet.IWallet, cli client.IManager) *VoucherCmd {
	return &VoucherCmd{
		db:            db,
		wallet:        wlt,
		clientManager: cli,
	}
}

func (v *VoucherCmd) GetCommand() *command.Command {
	middlewareHandler := command.NewMiddlewareHandler(v.db, v.wallet)

	subCmdClaim := &command.Command{
		Name: "claim",
		Help: "Claim voucher coins and bond them to a validator",
		Args: []*command.Args{
			{
				Name:     "code",
				Desc:     "The voucher code",
				InputBox: command.InputBoxText,
				Optional: false,
			},
			{
				Name:     "address",
				Desc:     "Your Pactus validator address",
				InputBox: command.InputBoxText,
				Optional: false,
			},
		},
		SubCommands: nil,
		AppIDs:      []entity.PlatformID{entity.PlatformIDDiscord},
		Middlewares: []command.MiddlewareFunc{middlewareHandler.WalletBalance},
		Handler:     v.claimHandler,
		TargetFlag:  command.TargetMaskMainnet,
	}

	subCmdCreateOne := &command.Command{
		Name: "create-one",
		Help: "Generate a single voucher code",
		Args: []*command.Args{
			{
				Name:     "amount",
				Desc:     "The amount of PAC to bond",
				InputBox: command.InputBoxFloat,
				Optional: false,
			},
			{
				Name:     "valid-months",
				Desc:     "Number of months the voucher remains valid after issuance",
				InputBox: command.InputBoxInteger,
				Optional: false,
			},
			{
				Name:     "recipient",
				Desc:     "The recipient's name for the voucher",
				InputBox: command.InputBoxText,
				Optional: true,
			},
			{
				Name:     "description",
				Desc:     "A description of the voucher's purpose",
				InputBox: command.InputBoxText,
				Optional: true,
			},
		},
		SubCommands: nil,
		AppIDs:      []entity.PlatformID{entity.PlatformIDDiscord},
		Middlewares: []command.MiddlewareFunc{middlewareHandler.OnlyModerator},
		Handler:     v.createOneHandler,
		TargetFlag:  command.TargetMaskModerator,
	}

	subCmdCreateBulk := &command.Command{
		Name: "create-bulk",
		Help: "Generate multiple voucher codes by importing a file",
		Args: []*command.Args{
			{
				Name:     "file",
				Desc:     "File containing a list of voucher recipients",
				InputBox: command.InputBoxFile,
				Optional: false,
			},
			{
				Name:     "notify",
				Desc:     "Send notifications to recipients via email",
				InputBox: command.InputBoxToggle,
				Optional: false,
			},
		},
		SubCommands: nil,
		AppIDs:      []entity.PlatformID{entity.PlatformIDDiscord},
		Middlewares: []command.MiddlewareFunc{middlewareHandler.OnlyModerator},
		Handler:     v.createBulkHandler,
		TargetFlag:  command.TargetMaskModerator,
	}

	subCmdStatus := &command.Command{
		Name: "status",
		Help: "View the status of vouchers or a specific voucher",
		Args: []*command.Args{
			{
				Name:     "code",
				Desc:     "The voucher code (8 characters)",
				InputBox: command.InputBoxText,
				Optional: true,
			},
		},
		SubCommands: nil,
		AppIDs:      []entity.PlatformID{entity.PlatformIDDiscord},
		Middlewares: []command.MiddlewareFunc{middlewareHandler.OnlyModerator},
		Handler:     v.statusHandler,
		TargetFlag:  command.TargetMaskModerator,
	}

	cmdVoucher := &command.Command{
		Name:        "voucher",
		Help:        "Commands for managing vouchers",
		Args:        nil,
		AppIDs:      []entity.PlatformID{entity.PlatformIDDiscord},
		SubCommands: make([]*command.Command, 0),
		Handler:     nil,
		TargetFlag:  command.TargetMaskMainnet | command.TargetMaskModerator,
	}

	cmdVoucher.AddSubCommand(subCmdClaim)
	cmdVoucher.AddSubCommand(subCmdCreateOne)
	cmdVoucher.AddSubCommand(subCmdCreateBulk)
	cmdVoucher.AddSubCommand(subCmdStatus)

	return cmdVoucher
}
