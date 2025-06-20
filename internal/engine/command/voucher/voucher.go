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

func (v *VoucherCmd) BuildCommand(botID entity.BotID) *command.Command {
	cmd := v.buildVoucherCommand(botID)

	choices := []command.Choice{}
	for tmplName := range v.templates {
		choice := command.Choice{
			Name:  tmplName,
			Desc:  tmplName,
			Value: tmplName,
		}

		choices = append(choices, choice)
	}
	v.subCmdCreate.Args[0].Choices = choices

	return cmd
}
