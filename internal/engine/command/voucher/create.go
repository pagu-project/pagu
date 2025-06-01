package voucher

import (
	"strconv"

	"github.com/pagu-project/pagu/internal/engine/command"
	"github.com/pagu-project/pagu/internal/entity"
	"github.com/pagu-project/pagu/pkg/amount"
	"github.com/pagu-project/pagu/pkg/utils"
)

type BulkRecorder struct {
	Recipient        string  `csv:"Recipient"`
	Email            string  `csv:"Email"`
	Amount           float64 `csv:"Amount"`
	ValidatedInMonth int     `csv:"Validated"`
	Description      string  `csv:"Description"`
}

func (v *VoucherCmd) createHandler(
	caller *entity.User,
	cmd *command.Command,
	args map[string]string,
) command.CommandResult {
	code := utils.RandomString(8, utils.CapitalAlphanumerical)
	for _, err := v.db.GetVoucherByCode(code); err == nil; {
		code = utils.RandomString(8, utils.CapitalAlphanumerical)
	}

	amt, err := amount.FromString(args["amount"])
	if err != nil {
		return cmd.RenderFailedTemplate("Invalid amount param")
	}

	maxStake, _ := amount.NewAmount(1000)
	if amt > maxStake {
		return cmd.RenderFailedTemplate("Stake amount is more than 1000")
	}

	expireMonths, err := strconv.Atoi(args["valid-months"])
	if err != nil {
		return cmd.RenderFailedTemplate("Invalid valid-months param")
	}

	vch := &entity.Voucher{
		Creator:     caller.ID,
		Code:        code,
		ValidMonths: uint8(expireMonths),
		Amount:      amt,
	}

	vch.Recipient = args["recipient"]
	vch.Desc = args["description"]

	err = v.db.AddVoucher(vch)
	if err != nil {
		return cmd.RenderErrorTemplate(err)
	}

	// TODO: fix me later
	//   1- Ensure recipient is a valid email
	//   2- define voucher template
	//   3- define data for templates
	//   4- Send email and test
	// v.mailer.SendTemplateMail(vch.Recipient, )

	return cmd.RenderResultTemplate("voucher", vch)
}
