package phoenix

import (
	"context"
	"testing"

	"github.com/pagu-project/pagu/internal/repository"
	"github.com/pagu-project/pagu/internal/testsuite"
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
	_, privateKey := ts.RandEd25519KeyPair()
	cfg := &Config{
		Client:       "testnet1.pactus.org:50052",
		PrivateKey:   privateKey,
		FaucetAmount: 5,
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
