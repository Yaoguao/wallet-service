package postgres

import (
	"context"
	"log/slog"
	"wallet-service/internal/domain/models"
	pgxdriver "wallet-service/pkg/pgx-driver"
	"wallet-service/pkg/pgx-driver/transaction"
)

// TODO - methods: GetOperationsByWallet(ctx, walletID uuid.UUID, limit, offset int) ([]*Operation, error)(optional)

type OperationRepository struct {
	postgres *pgxdriver.Postgres
	log      *slog.Logger
}

func NewOperationRepository(log *slog.Logger, postgres *pgxdriver.Postgres) *OperationRepository {
	return &OperationRepository{
		postgres: postgres,
		log:      log,
	}
}

func (or *OperationRepository) CreateOperation(
	ctx context.Context,
	tx pgxdriver.QueryExecuter,
	operation *models.Operation,
) error {

	const op = "storage.postgres.CreateOperation"

	query, args, err := or.postgres.Insert("operation").
		Columns("wallet_id", "type", "amount").
		Values(operation.WalletID, operation.Type, operation.Amount).
		ToSql()

	if err != nil {
		return transaction.HandleError(op, "insert", err)
	}

	_, err = tx.Exec(ctx, query, args...)
	if err != nil {
		return transaction.HandleError(op, "insert", err)
	}

	return nil
}
