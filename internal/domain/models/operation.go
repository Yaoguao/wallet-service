package models

import (
	"time"

	"github.com/google/uuid"
)

type OperationType string

const (
	Deposit  OperationType = "DEPOSIT"
	Withdraw OperationType = "WITHDRAW"
)

type Operation struct {
	ID        uuid.UUID     `db:"id"`
	WalletID  uuid.UUID     `db:"wallet_id"`
	Type      OperationType `db:"type"`
	Amount    int64         `db:"amount"`
	CreatedAt time.Time     `db:"created_at"`
}
