package services

import "errors"

var (
	ErrAmountNegativeValue = errors.New("amount negative value")
	ErrWalletNotFound      = errors.New("wallet not found")

	ErrInvalidWalletID = errors.New("invalid argument")
)
