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
	ID        uuid.UUID     `json:"id" db:"id"`
	WalletID  uuid.UUID     `json:"wallet_id" db:"wallet_id"`
	Type      OperationType `json:"type" db:"type"`
	Amount    int64         `json:"amount" db:"amount"`
	CreatedAt time.Time     `json:"created_at" db:"created_at"`
}
