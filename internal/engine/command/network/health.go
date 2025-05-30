package network

import (
	"time"

	"github.com/pagu-project/pagu/internal/engine/command"
	"github.com/pagu-project/pagu/internal/entity"
	"github.com/pagu-project/pagu/pkg/utils"
)

func (n *NetworkCmd) healthHandler(
	_ *entity.User,
	cmd *command.Command,
	_ map[string]string,
) command.CommandResult {
	lastBlockTime, lastBlockHeight, err := n.clientMgr.GetLastBlockTime()
	if err != nil {
		return cmd.RenderErrorTemplate(err)
	}

	lastBlockTimeFormatted := time.Unix(int64(lastBlockTime), 0).Format("02/01/2006, 15:04:05")
	currentTime := time.Now()
	timeDiff := currentTime.Unix() - int64(lastBlockTime)

	healthStatus := timeDiff <= 15

	var status string
	if healthStatus {
		status = "Healthy✅"
	} else {
		status = "UnHealthy❌"
	}

	return cmd.RenderResultTemplate(
		"Status", status,
		"CurrentTime", currentTime.Format("02/01/2006, 15:04:05"),
		"LastBlockTime", lastBlockTimeFormatted,
		"TimeDiff", timeDiff,
		"LastBlockHeight", utils.FormatNumber(int64(lastBlockHeight)),
	)
}
