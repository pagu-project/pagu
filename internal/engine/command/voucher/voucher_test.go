package voucher

import (
	"context"
	"testing"
	"time"

	"github.com/pagu-project/pagu/internal/entity"
	"github.com/pagu-project/pagu/internal/repository"
	"github.com/pagu-project/pagu/internal/testsuite"
	"github.com/pagu-project/pagu/pkg/amount"
	"github.com/pagu-project/pagu/pkg/client"
	"github.com/pagu-project/pagu/pkg/mailer"
	"github.com/pagu-project/pagu/pkg/wallet"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

type testData struct {
	*testsuite.TestSuite

	voucherCmd    *VoucherCmd
	testDB        *repository.Database
	mockClientMgr *client.MockIManager
	mockWallet    *wallet.MockIWallet
	mockMailer    *mailer.MockIMailer
}

func setup(t *testing.T) *testData {
	t.Helper()

	ts := testsuite.NewTestSuite(t)
	ctrl := gomock.NewController(t)

	tempFile := ts.CreateTempFile("{{.Recipient}} {{.Amount}} {{.Code}}")

	testConfig := Config{
		Templates: map[string]string{
			"sample": tempFile,
		},
	}
	testDB := ts.MakeTestDB()
	mockClientMgr := client.NewMockIManager(ctrl)
	mockWallet := wallet.NewMockIWallet(ctrl)
	mockMailer := mailer.NewMockIMailer(ctrl)

	voucherCmd := NewVoucherCmd(context.Background(), &testConfig,
		testDB, mockWallet, mockClientMgr, mockMailer)
	voucherCmd.BuildCommand(entity.BotID_CLI)

	return &testData{
		TestSuite:     ts,
		voucherCmd:    voucherCmd,
		testDB:        testDB,
		mockClientMgr: mockClientMgr,
		mockWallet:    mockWallet,
		mockMailer:    mockMailer,
	}
}

type VoucherOption func(*entity.Voucher)

func WithCode(code string) VoucherOption {
	return func(v *entity.Voucher) {
		v.Code = code
	}
}

func WithAmount(amt amount.Amount) VoucherOption {
	return func(v *entity.Voucher) {
		v.Amount = amt
	}
}

func WithTxHash(txHash string) VoucherOption {
	return func(v *entity.Voucher) {
		v.TxHash = txHash
	}
}

func WithValidMonths(validMonths uint8) VoucherOption {
	return func(v *entity.Voucher) {
		v.ValidMonths = validMonths
	}
}

func WithCreatedAt(createdAt time.Time) VoucherOption {
	return func(v *entity.Voucher) {
		v.CreatedAt = createdAt
	}
}

// Helper for setting email in VoucherOption style.
func WithEmail(email string) VoucherOption {
	return func(v *entity.Voucher) {
		v.Email = email
	}
}

func WithRecipient(recipient string) VoucherOption {
	return func(v *entity.Voucher) {
		v.Recipient = recipient
	}
}

func (td *testData) createTestVoucher(t *testing.T, opts ...VoucherOption) *entity.Voucher {
	t.Helper()

	voucher := &entity.Voucher{
		ValidMonths: 1,
		Email:       td.RandEmail(),
		Recipient:   td.RandString(8),
		Amount:      td.RandAmount(),
		Creator:     uint(td.RandInt(100)),
		Code:        td.RandString(8),
	}

	// Apply options
	for _, opt := range opts {
		opt(voucher)
	}

	err := td.testDB.AddVoucher(voucher)
	require.NoError(t, err)

	return voucher
}
