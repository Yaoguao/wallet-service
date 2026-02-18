package storage

import "errors"

var (
	ErrWalletExists   = errors.New("wallet already exists")
	ErrWalletNotFound = errors.New("wallet not found")

	ErrOperationExists   = errors.New("operation already exists")
	ErrOperationNotFound = errors.New("operation not found")

	ErrInsufficientFunds = errors.New("insufficient funds")
)
