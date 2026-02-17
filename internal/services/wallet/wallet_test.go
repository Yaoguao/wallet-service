package wallet

import (
	"context"
	"testing"
	"time"
	"wallet-service/internal/services"
	"wallet-service/internal/services/wallet/mocks"
	pgxdriver "wallet-service/pkg/pgx-driver"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestWalletService_Deposit_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTxManager := mocks.NewMockManager(ctrl)
	mockBalanceUpdater := mocks.NewMockBalanceUpdaterWallet(ctrl)
	mockOperationSaver := mocks.NewMockOperationSaver(ctrl)

	ctx := context.Background()
	walletID := uuid.New()
	amount := int64(100)
	newBalance := int64(200)
	now := time.Now()

	mockTxManager.
		EXPECT().
		ExecuteInTransaction(ctx, "deposit", gomock.Any()).
		DoAndReturn(func(
			ctx context.Context,
			name string,
			fn func(tx pgxdriver.QueryExecuter) error,
		) error {
			return fn(nil)
		})

	mockBalanceUpdater.
		EXPECT().
		IncreaseBalance(ctx, gomock.Any(), walletID, amount).
		Return(newBalance, now, nil)

	mockOperationSaver.
		EXPECT().
		CreateOperation(ctx, gomock.Any(), gomock.Any()).
		Return(nil)

	service := &WalletService{
		txManager:            mockTxManager,
		walletBalanceUpdater: mockBalanceUpdater,
		operationSaver:       mockOperationSaver,
	}

	result, err := service.Deposit(ctx, walletID, amount)

	require.NoError(t, err)
	require.NotNil(t, result)
	require.Equal(t, newBalance, result.Balance)
}

func TestWalletService_Withdraw_InsufficientFunds(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTxManager := mocks.NewMockManager(ctrl)
	mockBalanceUpdater := mocks.NewMockBalanceUpdaterWallet(ctrl)
	mockOperationSaver := mocks.NewMockOperationSaver(ctrl)

	ctx := context.Background()
	walletID := uuid.New()
	amount := int64(100)

	mockTxManager.
		EXPECT().
		ExecuteInTransaction(ctx, "withdraw", gomock.Any()).
		DoAndReturn(func(
			ctx context.Context,
			name string,
			fn func(tx pgxdriver.QueryExecuter) error,
		) error {
			return fn(nil)
		})

	mockBalanceUpdater.
		EXPECT().
		DecreaseBalance(ctx, gomock.Any(), walletID, amount).
		Return(int64(0), time.Time{}, services.ErrInsufficientFunds)

	service := &WalletService{
		txManager:            mockTxManager,
		walletBalanceUpdater: mockBalanceUpdater,
		operationSaver:       mockOperationSaver,
	}

	_, err := service.Withdraw(ctx, walletID, amount)

	require.Error(t, err)
	require.Contains(t, err.Error(), "Withdraw")
}
