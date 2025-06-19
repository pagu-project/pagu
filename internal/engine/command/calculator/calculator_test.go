package calculator

import (
	"testing"

	"github.com/pagu-project/pagu/internal/entity"
	"github.com/pagu-project/pagu/internal/testsuite"
	"github.com/pagu-project/pagu/pkg/client"
	"go.uber.org/mock/gomock"
)

type testData struct {
	*testsuite.TestSuite

	calculatorCmd *CalculatorCmd
	mockClientMgr *client.MockIManager
}

func setup(t *testing.T) *testData {
	t.Helper()

	ts := testsuite.NewTestSuite(t)
	ctrl := gomock.NewController(t)

	mockClientMgr := client.NewMockIManager(ctrl)

	calculatorCmd := NewCalculatorCmd(mockClientMgr)
	calculatorCmd.BuildCommand(entity.BotID_CLI)

	return &testData{
		TestSuite:     ts,
		calculatorCmd: calculatorCmd,
		mockClientMgr: mockClientMgr,
	}
}
