package voucher

import (
	"errors"
	"fmt"
	"net/mail"
	"strconv"

	"github.com/pagu-project/pagu/internal/engine/command"
	"github.com/pagu-project/pagu/internal/entity"
	"github.com/pagu-project/pagu/pkg/amount"
	"github.com/pagu-project/pagu/pkg/log"
	"github.com/pagu-project/pagu/pkg/utils"
)

func (v *VoucherCmd) createHandler(
	caller *entity.User,
	cmd *command.Command,
	args map[string]string,
) command.CommandResult {
	voucher, err := v.createVoucher(
		caller,
		args[argNameCreateRecipient],
		args[argNameCreateEmail],
		args[argNameCreateAmount],
		args[argNameCreateValidMonths],
		args[argNameCreateDescription],
	)
	if err != nil {
		return cmd.RenderErrorTemplate(err)
	}

	err = v.sendEmail(args[argNameCreateTemplate], voucher)
	if err != nil {
		return cmd.RenderErrorTemplate(err)
	}

	return cmd.RenderResultTemplate("voucher", voucher)
}

func (v *VoucherCmd) createVoucher(caller *entity.User,
	recipient, email, amtStr, validMonthsStr, desc string,
) (*entity.Voucher, error) {
	existing := v.db.GetNonExpiredVoucherByEmail(email)
	if existing != nil {
		return nil, fmt.Errorf("email already has a non-expired voucher: %s", email)
	}

	amt, err := amount.FromString(amtStr)
	if err != nil {
		return nil, fmt.Errorf("invalid amount: %w", err)
	}

	maxAmt, _ := amount.NewAmount(1000)
	if amt > maxAmt {
		return nil, errors.New("amount is more than 1000 PAC")
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

	voucher := &entity.Voucher{
		Creator:     caller.ID,
		Code:        code,
		ValidMonths: uint8(expireMonths),
		Amount:      amt,
		Recipient:   recipient,
		Email:       email,
		Desc:        desc,
	}

	err = v.db.AddVoucher(voucher)
	if err != nil {
		return nil, err
	}
	log.Info("voucher created", "email", voucher.Email, "amount", voucher.Amount, "code", voucher.Code)

	return voucher, nil
}

func (v *VoucherCmd) sendEmail(tmplName string, voucher *entity.Voucher) error {
	tmplPath := v.templates[tmplName]

	data := map[string]string{
		"Code":        voucher.Code,
		"Amount":      voucher.Amount.String(),
		"ValidMonths": strconv.Itoa(int(voucher.ValidMonths)),
		"Recipient":   voucher.Recipient,
	}
	err := v.mailer.SendTemplateMailAsync(voucher.Email, tmplPath, data)
	if err != nil {
		return err
	}

	return nil
}
