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

func (c *VoucherCmd) createHandler(
	caller *entity.User,
	cmd *command.Command,
	args map[string]string,
) command.CommandResult {
	voucher, err := c.createVoucher(
		caller,
		args[argNameCreateType],
		args[argNameCreateRecipient],
		args[argNameCreateEmail],
		args[argNameCreateAmount],
		args[argNameCreateValidMonths],
		args[argNameCreateDescription],
	)
	if err != nil {
		return cmd.RenderErrorTemplate(err)
	}

	err = c.sendEmail(args[argNameCreateTemplate], voucher)
	if err != nil {
		return cmd.RenderErrorTemplate(err)
	}

	return cmd.RenderResultTemplate("voucher", voucher)
}

func (c *VoucherCmd) createVoucher(caller *entity.User,
	typStr string, recipient, email, amtStr, validMonthsStr, desc string,
) (*entity.Voucher, error) {
	existing := c.db.GetNonExpiredVoucherByEmail(email)
	if existing != nil {
		return nil, fmt.Errorf("email already has a non-expired voucher: %s", email)
	}

	typ, err := strconv.Atoi(typStr)
	if err != nil {
		return nil, fmt.Errorf("invalid type: %w", err)
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
		return nil, fmt.Errorf("invalid valid-months: %w", err)
	}

	_, err = mail.ParseAddress(email)
	if err != nil {
		return nil, fmt.Errorf("invalid email address: %s", email)
	}
	code := utils.RandomString(8, utils.CapitalAlphanumerical)

	voucher := &entity.Voucher{
		Creator:     caller.ID,
		Code:        code,
		Type:        uint8(typ),
		ValidMonths: uint8(expireMonths),
		Amount:      amt,
		Recipient:   recipient,
		Email:       email,
		Desc:        desc,
	}

	err = c.db.AddVoucher(voucher)
	if err != nil {
		return nil, err
	}
	log.Info("voucher created", "email", voucher.Email, "amount", voucher.Amount, "code", voucher.Code)

	return voucher, nil
}

func (c *VoucherCmd) sendEmail(tmplName string, voucher *entity.Voucher) error {
	tmplPath := c.templates[tmplName]

	data := map[string]string{
		"Code":        voucher.Code,
		"Amount":      voucher.Amount.String(),
		"ValidMonths": strconv.Itoa(int(voucher.ValidMonths)),
		"Recipient":   voucher.Recipient,
	}
	err := c.mailer.SendTemplateMailAsync(voucher.Email, tmplPath, data)
	if err != nil {
		return err
	}

	return nil
}
