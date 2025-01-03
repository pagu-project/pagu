package crowdfund

import (
	"context"
	"testing"

	"github.com/pagu-project/pagu/internal/entity"
	"github.com/pagu-project/pagu/internal/repository"
	"github.com/pagu-project/pagu/internal/testsuite"
	"github.com/pagu-project/pagu/pkg/nowpayments"
	"github.com/pagu-project/pagu/pkg/wallet"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

type testData struct {
	*testsuite.TestSuite

	crowdfundCmd *CrowdfundCmd
	database     *repository.Database
	nowpayments  *nowpayments.MockINowpayments
	wallet       *wallet.MockIWallet
}

func setup(t *testing.T) *testData {
	t.Helper()

	ts := testsuite.NewTestSuite(t)
	ctrl := gomock.NewController(t)

	testDB := ts.MakeTestDB()
	mockNowpayments := nowpayments.NewMockINowpayments(ctrl)
	mockWallet := wallet.NewMockIWallet(ctrl)

	crowdfundCmd := NewCrowdfundCmd(context.Background(),
		testDB, mockWallet, mockNowpayments)

	return &testData{
		TestSuite:    ts,
		crowdfundCmd: crowdfundCmd,
		database:     testDB,
		nowpayments:  mockNowpayments,
		wallet:       mockWallet,
	}
}

type CampaignOption func(*entity.CrowdfundCampaign)

func WithTitle(title string) CampaignOption {
	return func(c *entity.CrowdfundCampaign) {
		c.Title = title
	}
}

func WithPackages(packages []entity.Package) CampaignOption {
	return func(c *entity.CrowdfundCampaign) {
		c.Packages = packages
	}
}

func (td *testData) createTestCampaign(t *testing.T, opts ...CampaignOption) *entity.CrowdfundCampaign {
	t.Helper()

	campaign := &entity.CrowdfundCampaign{
		Title: td.RandString(16),
		Desc:  td.RandString(128),
		Packages: []entity.Package{
			entity.Package{
				Name:      td.RandString(16),
				USDAmount: td.RandInt(1000),
				PACAmount: td.RandInt(1000),
			},
			entity.Package{
				Name:      td.RandString(16),
				USDAmount: td.RandInt(1000),
				PACAmount: td.RandInt(1000),
			},
			entity.Package{
				Name:      td.RandString(16),
				USDAmount: td.RandInt(1000),
				PACAmount: td.RandInt(1000),
			},
		},
	}

	// Apply options
	for _, opt := range opts {
		opt(campaign)
	}

	err := td.database.AddCrowdfundCampaign(campaign)
	require.NoError(t, err)

	return campaign
}
