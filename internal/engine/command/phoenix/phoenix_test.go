package phoenix

import (
	"context"
	"testing"
	"time"

	"github.com/pagu-project/pagu/internal/repository"
	"github.com/pagu-project/pagu/internal/testsuite"
	"github.com/pagu-project/pagu/pkg/amount"
)

type testData struct {
	*testsuite.TestSuite

	phoenixCmd *PhoenixCmd
	database   *repository.Database
}

func setup(t *testing.T) *testData {
	t.Helper()

	ts := testsuite.NewTestSuite(t)

	testDB := ts.MakeTestDB()
	cfg := &Config{
		Client:         "testnet1.pactus.org:50052",
		PrivateKey:     "TSECRET1RZSMS2JGNFLRU26NHNQK3JYTD4KGKLGW4S7SG75CZ057SR7CE8HUSG5MS3Z",
		FaucetAmount:   amount.Amount(1),
		FaucetFee:      amount.Amount(0),
		FaucetCooldown: 1 * time.Hour,
	}

	phoenixCmd := NewPhoenixCmd(context.Background(), cfg,
		testDB)

	_ = phoenixCmd.GetCommand()

	return &testData{
		TestSuite:  ts,
		phoenixCmd: phoenixCmd,
		database:   testDB,
	}
}
