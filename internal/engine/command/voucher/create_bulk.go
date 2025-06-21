package voucher

import (
	"time"

	"github.com/jszwec/csvutil"
	"github.com/pagu-project/pagu/internal/engine/command"
	"github.com/pagu-project/pagu/internal/entity"
	"github.com/pagu-project/pagu/pkg/log"
)

type bulkRecorder struct {
	Recipient   string `csv:"recipient"`
	Email       string `csv:"email"`
	Amount      string `csv:"amount"`
	ValidMonths string `csv:"valid-months"`
	Description string `csv:"desc"`
}

func (c *VoucherCmd) createBulkHandler(
	caller *entity.User,
	cmd *command.Command,
	args map[string]string,
) command.CommandResult {
	var bulkRecorders []bulkRecorder
	err := csvutil.Unmarshal([]byte(args[argNameCreateBulkCsv]), &bulkRecorders)
	if err != nil {
		return cmd.RenderErrorTemplate(err)
	}

	for index, rec := range bulkRecorders {
		voucher, err := c.createVoucher(
			caller,
			args[argNameCreateBulkType],
			rec.Recipient,
			rec.Email,
			rec.Amount,
			rec.ValidMonths,
			rec.Description,
		)
		if err != nil {
			return cmd.RenderErrorTemplate(err)
		}

		go func() {
			sleepTime := time.Duration(index*15) * time.Minute
			sleepTime += 5 * time.Second
			time.Sleep(sleepTime)

			err := c.sendEmail(args[argNameCreateBulkTemplate], voucher)
			if err != nil {
				log.Error("unable to send bulk email", "error", err)
			}
		}()
	}

	return cmd.RenderResultTemplate()
}
