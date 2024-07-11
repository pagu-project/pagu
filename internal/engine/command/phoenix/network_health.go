package phoenix

import (
	"time"

	"github.com/pagu-project/Pagu/internal/engine/command"
	"github.com/pagu-project/Pagu/internal/entity"
	utils2 "github.com/pagu-project/Pagu/pkg/utils"
)

//nolint:unused // remove me after I am used
func (pt *Phoenix) networkHealthHandler(cmd *command.Command,
	_ entity.AppID, _ string, _ ...string,
) command.CommandResult {
	lastBlockTime, lastBlockHeight := pt.clientMgr.GetLastBlockTime()
	lastBlockTimeFormatted := time.Unix(int64(lastBlockTime), 0).Format("02/01/2006, 15:04:05")
	currentTime := time.Now()

	timeDiff := (currentTime.Unix() - int64(lastBlockTime))

	healthStatus := true
	if timeDiff > 15 {
		healthStatus = false
	}

	var status string
	if healthStatus {
		status = "Healthy✅"
	} else {
		status = "UnHealthy❌"
	}

	return cmd.SuccessfulResult("Network is %s\nCurrentTime: %v\n"+
		"LastBlockTime: %v\nTime Diff: %v\nLast Block Height: %v",
		status, currentTime.Format("02/01/2006, 15:04:05"),
		lastBlockTimeFormatted, timeDiff, utils2.FormatNumber(int64(lastBlockHeight)))
}
