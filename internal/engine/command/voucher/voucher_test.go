package voucher

import (
	"testing"
	"time"

	"github.com/pagu-project/pagu/internal/entity"
	"github.com/pagu-project/pagu/internal/repository"
	"github.com/pagu-project/pagu/internal/testsuite"
	"github.com/pagu-project/pagu/pkg/amount"
	"github.com/pagu-project/pagu/pkg/client"
	"github.com/pagu-project/pagu/pkg/wallet"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

type testData struct {
	*testsuite.TestSuite

	voucherCmd    *VoucherCmd
	database      *repository.Database
	clientManager *client.MockIManager
	wallet        *wallet.MockIWallet
}

func setup(t *testing.T) *testData {
	t.Helper()

	ts := testsuite.NewTestSuite(t)
	ctrl := gomock.NewController(t)

	testDB := ts.MakeTestDB()
	mockClientManager := client.NewMockIManager(ctrl)
	mockWallet := wallet.NewMockIWallet(ctrl)

	voucher := NewVoucherCmd(testDB, mockWallet, mockClientManager)

	return &testData{
		TestSuite:     ts,
		voucherCmd:    voucher,
		database:      testDB,
		clientManager: mockClientManager,
		wallet:        mockWallet,
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

func WithRecipient(recipient string) VoucherOption {
	return func(v *entity.Voucher) {
		v.Recipient = recipient
	}
}

func (td *testData) createTestVoucher(t *testing.T, opts ...VoucherOption) *entity.Voucher {
	t.Helper()

	voucher := &entity.Voucher{
		ValidMonths: 1,
		Amount:      td.RandAmount(),
		Creator:     uint(td.RandInt(100)),
		Code:        td.RandString(8),
	}

	// Apply options
	for _, opt := range opts {
		opt(voucher)
	}

	err := td.database.AddVoucher(voucher)
	require.NoError(t, err)

	return voucher
}
