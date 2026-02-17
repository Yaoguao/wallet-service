package postgres

import (
	"context"
	"log/slog"
	"time"
	"wallet-service/internal/domain/models"
	pgxdriver "wallet-service/pkg/pgx-driver"
	"wallet-service/pkg/pgx-driver/transaction"

	"github.com/Masterminds/squirrel"
	"github.com/google/uuid"
)

type WalletRepository struct {
	postgres *pgxdriver.Postgres
	log      slog.Logger
}

func NewWalletRepository(log slog.Logger, postgres *pgxdriver.Postgres) *WalletRepository {
	return &WalletRepository{
		postgres: postgres,
		log:      log,
	}
}

func (wr *WalletRepository) CreateWallet(ctx context.Context, id uuid.UUID, balance int64) (*models.Wallet, error) {
	const op = "storage.postgres.CreateWallet"

	query, args, err := wr.postgres.
		Insert("wallets").
		Columns("id", "balance").
		Values(id, balance).
		Suffix("RETURNING id, balance, created_at, updated_at").
		ToSql()
	if err != nil {
		return nil, transaction.HandleError(op, "insert", err)
	}

	wallet := &models.Wallet{}
	err = wr.postgres.Pool.QueryRow(ctx, query, args...).Scan(
		&wallet.ID, &wallet.Balance, &wallet.CreatedAt, &wallet.UpdatedAt,
	)
	if err != nil {
		return nil, transaction.HandleError(op, "insert", err)
	}

	return wallet, nil
}

func (wr *WalletRepository) GetWallet(ctx context.Context, id uuid.UUID) (*models.Wallet, error) {
	const op = "storage.postgres.GetWallet"

	query, args, err := wr.postgres.
		Select("wallets").
		Columns("id", "balance", "created_at", "updated_at").
		ToSql()
	if err != nil {
		return nil, transaction.HandleError(op, "insert", err)
	}

	wallet := &models.Wallet{}
	err = wr.postgres.Pool.QueryRow(ctx, query, args...).Scan(
		&wallet.ID, &wallet.Balance, &wallet.CreatedAt, &wallet.UpdatedAt,
	)
	if err != nil {
		return nil, transaction.HandleError(op, "insert", err)
	}

	return wallet, nil
}

func (wr *WalletRepository) IncreaseBalance(
	ctx context.Context,
	tx pgxdriver.QueryExecuter,
	walletID uuid.UUID,
	amount int64,
) (int64, time.Time, error) {

	const op = "storage.postgres.IncreaseBalance"

	query, args, err := wr.postgres.
		Update("wallets").
		Set("balance", squirrel.Expr("balance + ?", amount)).
		Where(squirrel.Expr("id = ?", walletID)).
		Suffix("RETURNING balance, updated_at").
		ToSql()

	if err != nil {
		return 0, time.Time{}, transaction.HandleError(op, "build_update", err)
	}

	var newBalance int64
	var updatedAt time.Time
	err = tx.QueryRow(ctx, query, args...).Scan(&newBalance, &updatedAt)
	if err != nil {
		return 0, time.Time{}, transaction.HandleError(op, "update", err)
	}

	return newBalance, updatedAt, nil
}

func (wr *WalletRepository) DecreaseBalance(
	ctx context.Context,
	tx pgxdriver.QueryExecuter,
	walletID uuid.UUID,
	amount int64,
) (int64, time.Time, error) {

	const op = "storage.postgres.DecreaseBalance"

	query, args, err := wr.postgres.
		Update("wallets").
		Set("balance", squirrel.And{
			squirrel.Expr("balance - ?", amount),
			squirrel.Expr("balance >= ?", amount),
		}).
		Where(squirrel.Expr("id = ?", walletID)).
		Suffix("RETURNING balance, updated_at").
		ToSql()

	if err != nil {
		return 0, time.Time{}, transaction.HandleError(op, "build_update", err)
	}

	var newBalance int64
	var updatedAt time.Time
	err = tx.QueryRow(ctx, query, args...).Scan(&newBalance, &updatedAt)

	if err != nil {
		wr.log.Debug(op, err.Error())
		return 0, time.Time{}, transaction.HandleError(op, "update", err)
	}

	return newBalance, updatedAt, nil
}
