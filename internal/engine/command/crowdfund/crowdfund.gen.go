// Code generated by command-generator. DO NOT EDIT.
package crowdfund

import (
	"github.com/pagu-project/pagu/internal/engine/command"
)

const (
	argNameCreateTitle     = "title"
	argNameCreateDesc      = "desc"
	argNameCreatePackages  = "packages"
	argNamePurchasePackage = "package"
	argNameClaimAddress    = "address"
)

type crowdfundSubCmds struct {
	subCmdCreate   *command.Command
	subCmdDisable  *command.Command
	subCmdReport   *command.Command
	subCmdInfo     *command.Command
	subCmdPurchase *command.Command
	subCmdClaim    *command.Command
}

func (c *CrowdfundCmd) buildSubCmds() *crowdfundSubCmds {
	subCmdCreate := &command.Command{
		Name:           "create",
		Help:           "Create a new crowdfunding campaign",
		Handler:        c.createHandler,
		ResultTemplate: "Crowdfund campaign '{{.campaign.Title}}' created successfully with {{ .campaign.Packages | len }} packages\n",
		Args: []*command.Args{
			{
				Name:     "title",
				Desc:     "The title of this crowdfunding campaign",
				InputBox: command.InputBoxText,
				Optional: false,
			},
			{
				Name:     "desc",
				Desc:     "A description of this crowdfunding campaign",
				InputBox: command.InputBoxMultilineText,
				Optional: false,
			},
			{
				Name:     "packages",
				Desc:     "The packages for this campaign in JSON format",
				InputBox: command.InputBoxMultilineText,
				Optional: false,
			},
		},
	}
	subCmdDisable := &command.Command{
		Name:           "disable",
		Help:           "Disable an existing crowdfunding campaign",
		Handler:        c.disableHandler,
		ResultTemplate: ``,
	}
	subCmdReport := &command.Command{
		Name:           "report",
		Help:           "View reports of a crowdfunding campaign",
		Handler:        c.reportHandler,
		ResultTemplate: ``,
	}
	subCmdInfo := &command.Command{
		Name:           "info",
		Help:           "Get detailed information about a crowdfunding campaign",
		Handler:        c.infoHandler,
		ResultTemplate: "**{{.campaign.Title}}**\n{{.campaign.Desc}}\n\nPackages:\n{{range .campaign.Packages}}\n- {{.Name}}: {{.USDAmount}} USDT to {{.PACAmount }}\n{{- end}}\n",
	}
	subCmdPurchase := &command.Command{
		Name:           "purchase",
		Help:           "Make a purchase in a crowdfunding campaign",
		Handler:        c.purchaseHandler,
		ResultTemplate: "Your purchase of {{ .purchase.USDAmount }} USDT to receive {{ .purchase.PACAmount }} successfully registered in our database.\nPlease visit {{ .paymentLink }} to make the payment.\n\nOnce the payment is done, you can claim your PAC coins using \"claim\" command.\n\nThanks\n",
		Args: []*command.Args{
			{
				Name:     "package",
				Desc:     "Select the crowdfunding package",
				InputBox: command.InputBoxChoice,
				Optional: false,
			},
		},
	}
	subCmdClaim := &command.Command{
		Name:           "claim",
		Help:           "Claim packages from a crowdfunding campaign",
		Handler:        c.claimHandler,
		ResultTemplate: "Thank you for supporting the Pactus blockchain!\n\nYou can track your transaction here: {{.txLink}}\nIf you have any questions or need assistance, feel free to reach out to our community.\n",
		Args: []*command.Args{
			{
				Name:     "address",
				Desc:     "Set your Pactus address",
				InputBox: command.InputBoxText,
				Optional: false,
			},
		},
	}

	return &crowdfundSubCmds{
		subCmdCreate:   subCmdCreate,
		subCmdDisable:  subCmdDisable,
		subCmdReport:   subCmdReport,
		subCmdInfo:     subCmdInfo,
		subCmdPurchase: subCmdPurchase,
		subCmdClaim:    subCmdClaim,
	}
}

func (c *CrowdfundCmd) buildCrowdfundCommand() *command.Command {
	crowdfundCmd := &command.Command{
		Emoji:       "🤝",
		Name:        "crowdfund",
		Help:        "Commands for managing crowdfunding campaigns",
		SubCommands: make([]*command.Command, 0),
	}

	c.crowdfundSubCmds = c.buildSubCmds()

	crowdfundCmd.AddSubCommand(c.subCmdCreate)
	crowdfundCmd.AddSubCommand(c.subCmdDisable)
	crowdfundCmd.AddSubCommand(c.subCmdReport)
	crowdfundCmd.AddSubCommand(c.subCmdInfo)
	crowdfundCmd.AddSubCommand(c.subCmdPurchase)
	crowdfundCmd.AddSubCommand(c.subCmdClaim)

	return crowdfundCmd
}
