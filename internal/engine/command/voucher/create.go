package voucher

import (
	"fmt"
	"net/mail"
	"strconv"

	"github.com/pagu-project/pagu/internal/engine/command"
	"github.com/pagu-project/pagu/internal/entity"
	"github.com/pagu-project/pagu/pkg/amount"
	"github.com/pagu-project/pagu/pkg/log"
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
	vch, err := v.createVoucher(
		caller,
		args[argNameCreateTemplate],
		args[argNameCreateRecipient],
		args[argNameCreateEmail],
		args[argNameCreateAmount],
		args[argNameCreateValidMonths],
		args[argNameCreateDescription],
	)
	if err != nil {
		return cmd.RenderErrorTemplate(err)
	}

	return cmd.RenderResultTemplate("voucher", vch)
}

func (v *VoucherCmd) createVoucher(caller *entity.User,
	tmplName, recipient, email, amtStr, validMonthsStr, desc string,
) (*entity.Voucher, error) {
	existing, err := v.db.GetNonExpiredVoucherByEmail(email)
	if err != nil {
		return nil, fmt.Errorf("failed to check existing vouchers: %w", err)
	}
	if existing != nil {
		return nil, fmt.Errorf("email already has a non-expired voucher")
	}

	tmplPath, exists := v.templates[tmplName]
	if !exists {
		return nil, fmt.Errorf("template not exists: %s", tmplName)
	}

	amt, err := amount.FromString(amtStr)
	if err != nil {
		return nil, fmt.Errorf("invalid amount: %w", err)
	}

	maxAmt, _ := amount.NewAmount(1000)
	if amt > maxAmt {
		return nil, fmt.Errorf("amount is more than 1000 PAC")
	}

	expireMonths, err := strconv.Atoi(validMonthsStr)
	if err != nil {
		return nil, fmt.Errorf("invalid valid-months param: %w", err)
	}

	_, err = mail.ParseAddress(email)
	if err != nil {
		return nil, fmt.Errorf("invalid email address: %s", email)
	}
	code := utils.RandomString(8, utils.CapitalAlphanumerical)

	vch := &entity.Voucher{
		Creator:     caller.ID,
		Code:        code,
		ValidMonths: uint8(expireMonths),
		Amount:      amt,
		Recipient:   recipient,
		Email:       email,
		Desc:        desc,
	}

	err = v.db.AddVoucher(vch)
	if err != nil {
		return nil, err
	}

	go func() {
		data := map[string]string{
			"Code":        vch.Code,
			"Amount":      vch.Amount.String(),
			"ValidMonths": strconv.Itoa(int(vch.ValidMonths)),
			"Recipient":   vch.Recipient,
		}
		err = v.mailer.SendTemplateMail(email, tmplPath, data)
		if err != nil {
			log.Warn("failed to send voucher email: %v", err)
		}
	}()

	return vch, nil
}
