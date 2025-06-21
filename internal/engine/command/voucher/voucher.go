package voucher

import (
	"context"

	"github.com/pagu-project/pagu/internal/engine/command"
	"github.com/pagu-project/pagu/internal/entity"
	"github.com/pagu-project/pagu/internal/repository"
	"github.com/pagu-project/pagu/pkg/client"
	"github.com/pagu-project/pagu/pkg/mailer"
	"github.com/pagu-project/pagu/pkg/wallet"
)

type VoucherCmd struct {
	*voucherSubCmds

	ctx           context.Context
	db            *repository.Database
	wallet        wallet.IWallet
	clientManager client.IManager
	mailer        mailer.IMailer
	templates     map[string]string
}

func NewVoucherCmd(ctx context.Context, cfg *Config, db *repository.Database, wlt wallet.IWallet,
	clientManager client.IManager, mailer mailer.IMailer,
) *VoucherCmd {
	return &VoucherCmd{
		ctx:           ctx,
		db:            db,
		wallet:        wlt,
		clientManager: clientManager,
		mailer:        mailer,
		templates:     cfg.Templates,
	}
}

func (c *VoucherCmd) BuildCommand(botID entity.BotID) *command.Command {
	cmd := c.buildVoucherCommand(botID)

	templateChoices := []command.Choice{}
	for tmplName := range c.templates {
		choice := command.Choice{
			Name:  tmplName,
			Desc:  tmplName,
			Value: tmplName,
		}

		templateChoices = append(templateChoices, choice)
	}
	c.subCmdCreate.Args[0].Choices = templateChoices
	c.subCmdCreateBulk.Args[0].Choices = templateChoices

	typeChoices := []command.Choice{
		{Name: "Stake", Desc: "Voucher will be bonded to the validator address", Value: "0"},
		{Name: "Liquid", Desc: "Voucher will be transferred to the account address", Value: "1"},
	}
	c.subCmdCreate.Args[1].Choices = typeChoices
	c.subCmdCreateBulk.Args[1].Choices = typeChoices

	return cmd
}
