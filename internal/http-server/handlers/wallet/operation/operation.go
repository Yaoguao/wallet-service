package operation

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"wallet-service/internal/domain/models"
	"wallet-service/internal/http-server/handlers"
	"wallet-service/internal/services"
	"wallet-service/internal/storage"
	"wallet-service/pkg/helpers"

	"github.com/google/uuid"
)

var (
	emptyOperation       = ""
	emptyAmount    int64 = 0
)

type WalletService interface {
	Deposit(ctx context.Context, walletID uuid.UUID, amount int64) (*models.Wallet, error)
	Withdraw(ctx context.Context, walletID uuid.UUID, amount int64) (*models.Wallet, error)
}

type request struct {
	WalletID      uuid.UUID `json:"wallet_id"`
	OperationType string    `json:"operation_type"`
	Amount        int64     `json:"amount"`
}

type response struct {
	Status string         `json:"status"`
	Wallet *models.Wallet `json:"wallet,omitempty"`
	Error  string         `json:"error,omitempty"`
}

func New(log *slog.Logger, ws WalletService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req request
		err := helpers.ReadJSON(w, r, &req)
		if err != nil {
			log.Error("failed to decode request body")
			handlers.ErrorResponse(w, r, http.StatusInternalServerError, "ops! decode json")
			return
		}

		if err := validateRequest(req.WalletID, req.OperationType, req.Amount); err != nil {
			log.Error("failed to validate request")
			handlers.BadRequestResponse(w, r, err)
			return
		}

		var wallet *models.Wallet

		switch req.OperationType {
		case "DEPOSIT":
			wallet, err = ws.Deposit(r.Context(), req.WalletID, req.Amount)
		case "WITHDRAW":
			wallet, err = ws.Withdraw(r.Context(), req.WalletID, req.Amount)
		default:
			log.Error("failed to validate request")
			handlers.BadRequestResponse(w, r, errors.New("unknow operation"))
			return
		}

		if err != nil {
			log.Error(err.Error())
			if errors.Is(err, storage.ErrInsufficientFunds) {
				handlers.ErrorResponse(w, r, http.StatusBadRequest, "insufficient funds")
				return
			}
			if errors.Is(err, services.ErrAmountNegativeValue) {
				handlers.ErrorResponse(w, r, http.StatusBadRequest, "amount negative value")
				return
			}

			handlers.ErrorResponse(w, r, http.StatusInternalServerError, err.Error())
			return
		}

		err = helpers.WriteJSON(
			w,
			http.StatusOK,
			helpers.Envelope{"data": response{Wallet: wallet, Status: "operation was completed successfully"}},
			nil)

		if err != nil {
			log.Error(err.Error())
			return
		}

	}
}

func validateRequest(walletID uuid.UUID, opType string, amount int64) error {
	if walletID == uuid.Nil {
		return errors.New("error invalid argument: wallet_id")
	}

	if opType == emptyOperation {
		return errors.New("error invalid argument: operation_type")
	}

	if amount == emptyAmount {
		return errors.New("error the value is zero: amount")
	}
	return nil
}
