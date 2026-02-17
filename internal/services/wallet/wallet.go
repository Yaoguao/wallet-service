package wallet

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"
	"wallet-service/internal/domain/models"
	"wallet-service/internal/services"
	pgxdriver "wallet-service/pkg/pgx-driver"
	"wallet-service/pkg/pgx-driver/transaction"

	"github.com/google/uuid"
)

type SaverWallet interface {
	CreateWallet(ctx context.Context, id uuid.UUID, balance int64) (*models.Wallet, error)
}

type GetterWallet interface {
	GetWallet(ctx context.Context, id uuid.UUID) (*models.Wallet, error)
}

type OperationSaver interface {
	CreateOperation(ctx context.Context, tx pgxdriver.QueryExecuter, operation *models.Operation) error
}

type BalanceUpdaterWallet interface {
	IncreaseBalance(
		ctx context.Context,
		tx pgxdriver.QueryExecuter,
		walletID uuid.UUID,
		amount int64,
	) (int64, time.Time, error)
	DecreaseBalance(
		ctx context.Context,
		tx pgxdriver.QueryExecuter,
		walletID uuid.UUID,
		amount int64,
	) (int64, time.Time, error)
}

type WalletService struct {
	txManager transaction.Manager

	log                  slog.Logger
	walletSaver          SaverWallet
	walletGetter         GetterWallet
	walletBalanceUpdater BalanceUpdaterWallet

	operationSaver OperationSaver
}

func New(
	txManager transaction.Manager,
	log slog.Logger,
	walletSaver SaverWallet,
	walletGetter GetterWallet,
	walletBalanceUpdater BalanceUpdaterWallet,
	operationSaver OperationSaver,
) *WalletService {

	return &WalletService{
		txManager:            txManager,
		log:                  log,
		walletSaver:          walletSaver,
		walletGetter:         walletGetter,
		walletBalanceUpdater: walletBalanceUpdater,
		operationSaver:       operationSaver,
	}
}

func (ws *WalletService) CreateWallet(ctx context.Context, amount int64) (*models.Wallet, error) {
	const op = "services.wallet.CreateWallet"
	if amount < 0 {
		ws.log.Error("amount negative value")
		return nil, services.ErrAmountNegativeValue
	}

	id := uuid.New()

	wallet, err := ws.walletSaver.CreateWallet(ctx, id, amount)
	if err != nil {
		if errors.Is(err, transaction.ErrConflictingData) {
			ws.log.Debug("wallet already exist")
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return wallet, nil
}

func (ws *WalletService) GetWallet(ctx context.Context, id uuid.UUID) (*models.Wallet, error) {
	const op = "services.wallet.GetWallet"
	if id == uuid.Nil {
		return nil, services.ErrInvalidWalletID
	}

	wallet, err := ws.walletGetter.GetWallet(ctx, id)
	if err != nil {
		ws.log.Error("failed to get wallet",
			"wallet_id", id,
			"error", err,
		)
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return wallet, nil
}

func (ws *WalletService) Deposit(ctx context.Context, walletID uuid.UUID, amount int64) (*models.Wallet, error) {
	const op = "services.wallet.Deposit"

	if amount <= 0 {
		return nil, services.ErrAmountNegativeValue
	}

	var result *models.Wallet
	err := ws.txManager.ExecuteInTransaction(ctx, "deposit", func(tx pgxdriver.QueryExecuter) error {
		newBalance, updatedAt, err :=
			ws.walletBalanceUpdater.IncreaseBalance(ctx, tx, walletID, amount)

		if err != nil {
			return err
		}

		if err := ws.createOperationTx(ctx, tx, walletID, models.Deposit, amount); err != nil {
			return err
		}

		result = &models.Wallet{
			ID:        walletID,
			Balance:   newBalance,
			UpdatedAt: updatedAt,
		}

		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return result, nil
}

func (ws *WalletService) Withdraw(ctx context.Context, walletID uuid.UUID, amount int64) (*models.Wallet, error) {
	const op = "services.wallet.Withdraw"

	if amount <= 0 {
		return nil, services.ErrAmountNegativeValue
	}

	var result *models.Wallet
	err := ws.txManager.ExecuteInTransaction(ctx, "withdraw", func(tx pgxdriver.QueryExecuter) error {
		newBalance, updatedAt, err :=
			ws.walletBalanceUpdater.DecreaseBalance(ctx, tx, walletID, amount)

		if err != nil {
			return err
		}

		if err := ws.createOperationTx(ctx, tx, walletID, models.Withdraw, amount); err != nil {
			return err
		}

		result = &models.Wallet{
			ID:        walletID,
			Balance:   newBalance,
			UpdatedAt: updatedAt,
		}

		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return result, nil
}

func (ws *WalletService) createOperationTx(
	ctx context.Context,
	tx pgxdriver.QueryExecuter,
	walletID uuid.UUID,
	opType models.OperationType,
	amount int64,
) error {

	operation := &models.Operation{
		ID:       uuid.New(),
		WalletID: walletID,
		Type:     opType,
		Amount:   amount,
	}

	return ws.operationSaver.CreateOperation(ctx, tx, operation)
}
