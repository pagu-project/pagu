// Code generated by command-generator. DO NOT EDIT.
package voucher

import (
	"github.com/pagu-project/pagu/internal/engine/command"
)

const (
	argNameClaimCode         = "code"
	argNameClaimAddress      = "address"
	argNameCreateAmount      = "amount"
	argNameCreateValidMonths = "valid-months"
	argNameCreateRecipient   = "recipient"
	argNameCreateDescription = "description"
	argNameCreateBulkFile    = "file"
	argNameCreateBulkNotify  = "notify"
	argNameStatusCode        = "code"
)

type voucherSubCmds struct {
	subCmdClaim      *command.Command
	subCmdCreate     *command.Command
	subCmdCreateBulk *command.Command
	subCmdStatus     *command.Command
}

func (c *VoucherCmd) buildSubCmds() *voucherSubCmds {
	subCmdClaim := &command.Command{
		Name:           "claim",
		Help:           "Claim voucher coins and bond them to a validator",
		Handler:        c.claimHandler,
		ResultTemplate: "Voucher claimed successfully!\n\nhttps://pacviewer.com/transaction/{{.txHash}}\n",
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
	}
	subCmdCreate := &command.Command{
		Name:           "create",
		Help:           "Generate a single voucher code",
		Handler:        c.createHandler,
		ResultTemplate: "Voucher created successfully! \nCode: {{.vch.Code}}\n",
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
	}
	subCmdCreateBulk := &command.Command{
		Name:           "create-bulk",
		Help:           "Generate multiple voucher codes by importing a file",
		Handler:        c.createBulkHandler,
		ResultTemplate: "Vouchers created successfully!\n",
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
	}
	subCmdStatus := &command.Command{
		Name:           "status",
		Help:           "View the status of vouchers or a specific voucher",
		Handler:        c.statusHandler,
		ResultTemplate: "Code: {{.voucher.Code}}\nAmount: {{.voucher.Amount}}\nExpire At: {{.voucher.CreatedAt.AddDate 0 (int .voucher.ValidMonths) 0 | formatDate \"02/01/2006, 15:04:05\"}}\nRecipient: {{.voucher.Recipient}}\nDescription: {{.voucher.Desc}}\nClaimed: {{.isClaimed}}\nTx Link: {{.txLink}}\n",
		Args: []*command.Args{
			{
				Name:     "code",
				Desc:     "The voucher code (8 characters)",
				InputBox: command.InputBoxText,
				Optional: true,
			},
		},
	}

	return &voucherSubCmds{
		subCmdClaim:      subCmdClaim,
		subCmdCreate:     subCmdCreate,
		subCmdCreateBulk: subCmdCreateBulk,
		subCmdStatus:     subCmdStatus,
	}
}

func (c *VoucherCmd) buildVoucherCommand() *command.Command {
	voucherCmd := &command.Command{
		Name:        "voucher",
		Help:        "Commands for managing vouchers",
		SubCommands: make([]*command.Command, 0),
	}

	c.voucherSubCmds = c.buildSubCmds()

	voucherCmd.AddSubCommand(c.subCmdClaim)
	voucherCmd.AddSubCommand(c.subCmdCreate)
	voucherCmd.AddSubCommand(c.subCmdCreateBulk)
	voucherCmd.AddSubCommand(c.subCmdStatus)

	return voucherCmd
}
