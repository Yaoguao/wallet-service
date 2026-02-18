package services

import "errors"

var (
	ErrAmountNegativeValue = errors.New("amount negative value")

	ErrInvalidWalletID = errors.New("invalid argument")
)
