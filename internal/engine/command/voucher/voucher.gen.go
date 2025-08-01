// Code generated by command-generator. DO NOT EDIT.
package voucher

import (
	"github.com/pagu-project/pagu/internal/engine/command"
	"github.com/pagu-project/pagu/internal/entity"
)

const (
	argNameClaimCode          = "code"
	argNameClaimAddress       = "address"
	argNameCreateTemplate     = "template"
	argNameCreateType         = "type"
	argNameCreateRecipient    = "recipient"
	argNameCreateEmail        = "email"
	argNameCreateAmount       = "amount"
	argNameCreateValidMonths  = "valid-months"
	argNameCreateDescription  = "description"
	argNameCreateBulkTemplate = "template"
	argNameCreateBulkType     = "type"
	argNameCreateBulkCsv      = "csv"
	argNameStatusCode         = "code"
	argNameStatusEmail        = "email"
	argNameReportSince        = "since"
)

type voucherSubCmds struct {
	subCmdClaim      *command.Command
	subCmdCreate     *command.Command
	subCmdCreateBulk *command.Command
	subCmdStatus     *command.Command
	subCmdReport     *command.Command
}

func (c *VoucherCmd) buildSubCmds() *voucherSubCmds {
	subCmdClaim := &command.Command{
		Name:            "claim",
		Help:            "Claim your voucher",
		Handler:         c.claimHandler,
		ResultTemplate:  "Voucher claimed successfully!\n\nhttps://pacviewer.com/transaction/{{.txHash}}\n",
		TargetBotIDs:    entity.AllBotIDs(),
		TargetUserRoles: entity.AllUserRoles(),
		Args: []*command.Args{
			{
				Name:     "code",
				Desc:     "YOur voucher code",
				InputBox: command.InputBoxText,
				Optional: false,
			},
			{
				Name:     "address",
				Desc:     "Your Pactus address to receive the PAC coins",
				InputBox: command.InputBoxText,
				Optional: false,
			},
		},
	}
	subCmdCreate := &command.Command{
		Name:           "create",
		Help:           "Generate a single voucher code",
		Handler:        c.createHandler,
		ResultTemplate: "Voucher created successfully!\nRecipient:: {{.voucher.Recipient}}\nEmail: {{.voucher.Email}}\nAmount: {{.voucher.Amount}}\nCode: {{.voucher.Code}}\n",
		TargetBotIDs: []entity.BotID{
			entity.BotID_CLI,
			entity.BotID_Moderator,
		},
		TargetUserRoles: []entity.UserRole{
			entity.UserRole_Admin,
			entity.UserRole_Moderator,
		},
		Args: []*command.Args{
			{
				Name:     "template",
				Desc:     "The email template to use for the voucher",
				InputBox: command.InputBoxChoice,
				Optional: false,
			},
			{
				Name:     "type",
				Desc:     "Type of voucher (Stake or Liquid)",
				InputBox: command.InputBoxChoice,
				Optional: false,
			},
			{
				Name:     "recipient",
				Desc:     "The name of the recipient",
				InputBox: command.InputBoxText,
				Optional: false,
			},
			{
				Name:     "email",
				Desc:     "The email address to send the voucher to",
				InputBox: command.InputBoxText,
				Optional: false,
			},
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
				Name:     "description",
				Desc:     "A description of the voucher's purpose",
				InputBox: command.InputBoxText,
				Optional: false,
			},
		},
	}
	subCmdCreateBulk := &command.Command{
		Name:           "create-bulk",
		Help:           "Generate bulk voucher codes from JSON input",
		Handler:        c.createBulkHandler,
		ResultTemplate: "Vouchers are going to send to recipients.\n",
		TargetBotIDs: []entity.BotID{
			entity.BotID_CLI,
			entity.BotID_Moderator,
		},
		TargetUserRoles: []entity.UserRole{
			entity.UserRole_Admin,
			entity.UserRole_Moderator,
		},
		Args: []*command.Args{
			{
				Name:     "template",
				Desc:     "The email template to use for the voucher",
				InputBox: command.InputBoxChoice,
				Optional: false,
			},
			{
				Name:     "type",
				Desc:     "Type of voucher (Stake or Liquid)",
				InputBox: command.InputBoxChoice,
				Optional: false,
			},
			{
				Name:     "csv",
				Desc:     "The csv file containing the voucher information (`recipient,email,amount,valid-months,desc`)",
				InputBox: command.InputBoxFile,
				Optional: false,
			},
		},
	}
	subCmdStatus := &command.Command{
		Name:           "status",
		Help:           "View the status of a specific voucher",
		Handler:        c.statusHandler,
		ResultTemplate: "**Code**: {{.voucher.Code}}\n**Amount**: {{.voucher.Amount}}\n**Expire** At: {{.expireAt}}\n**Recipient**: {{.voucher.Recipient}}\n**Description**: {{.voucher.Desc}}\n**Claimed**: {{.isClaimed}}\n**Tx** Link: {{.txLink}}\n",
		TargetBotIDs: []entity.BotID{
			entity.BotID_CLI,
			entity.BotID_Moderator,
		},
		TargetUserRoles: []entity.UserRole{
			entity.UserRole_Admin,
			entity.UserRole_Moderator,
		},
		Args: []*command.Args{
			{
				Name:     "code",
				Desc:     "The voucher code (8 characters)",
				InputBox: command.InputBoxText,
				Optional: true,
			},
			{
				Name:     "email",
				Desc:     "The recipient email",
				InputBox: command.InputBoxText,
				Optional: true,
			},
		},
	}
	subCmdReport := &command.Command{
		Name:           "report",
		Help:           "The report of total vouchers",
		Handler:        c.reportHandler,
		ResultTemplate: "**Total Vouchers**: {{.total}}\n**Total Claimed**: {{.totalClaimed}}\n**Total Claimed Amount**: {{.totalClaimedAmount}}\n**Total Expired**: {{.totalExpired}}\n",
		TargetBotIDs: []entity.BotID{
			entity.BotID_CLI,
			entity.BotID_Moderator,
		},
		TargetUserRoles: []entity.UserRole{
			entity.UserRole_Admin,
			entity.UserRole_Moderator,
		},
		Args: []*command.Args{
			{
				Name:     "since",
				Desc:     "Since how many month ago",
				InputBox: command.InputBoxInteger,
				Optional: true,
			},
		},
	}

	return &voucherSubCmds{
		subCmdClaim:      subCmdClaim,
		subCmdCreate:     subCmdCreate,
		subCmdCreateBulk: subCmdCreateBulk,
		subCmdStatus:     subCmdStatus,
		subCmdReport:     subCmdReport,
	}
}

func (c *VoucherCmd) buildVoucherCommand(botID entity.BotID) *command.Command {
	voucherCmd := &command.Command{
		Name:            "voucher",
		Emoji:           "🎁",
		Active:          true,
		Help:            "Commands for managing vouchers",
		SubCommands:     make([]*command.Command, 0),
		TargetBotIDs:    entity.AllBotIDs(),
		TargetUserRoles: entity.AllUserRoles(),
	}

	c.voucherSubCmds = c.buildSubCmds()

	voucherCmd.AddSubCommand(botID, c.subCmdClaim)
	voucherCmd.AddSubCommand(botID, c.subCmdCreate)
	voucherCmd.AddSubCommand(botID, c.subCmdCreateBulk)
	voucherCmd.AddSubCommand(botID, c.subCmdStatus)
	voucherCmd.AddSubCommand(botID, c.subCmdReport)

	return voucherCmd
}
