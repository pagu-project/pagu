package crowdfund

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReport(t *testing.T) {
	td := setup(t)

	t.Run("ok", func(t *testing.T) {
		result := td.crowdfundCmd.reportHandler(nil, td.crowdfundCmd.subCmdReport, nil)
		assert.True(t, result.Successful)
	})
}
