package get

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"wallet-service/internal/domain/models"
	"wallet-service/internal/http-server/handlers"
	"wallet-service/internal/storage"
	"wallet-service/pkg/helpers"

	"github.com/google/uuid"
)

type WalletGetter interface {
	GetWallet(ctx context.Context, id uuid.UUID) (*models.Wallet, error)
}

type response struct {
	Status string         `json:"status"`
	Wallet *models.Wallet `json:"wallet,omitempty"`
	Error  string         `json:"error,omitempty"`
}

func New(log *slog.Logger, wg WalletGetter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := helpers.ReadUUIDParam(r, "WALLET_UUID")
		if err != nil {
			if err != nil {
				log.Error("failed to decode request param")

				handlers.ErrorResponse(w, r, http.StatusInternalServerError, "error read param id")
				return
			}
		}

		if id == uuid.Nil {
			handlers.BadRequestResponse(w, r, errors.New("error invalid argument: wallet_id"))
			return
		}

		wallet, err := wg.GetWallet(r.Context(), id)
		if err != nil {
			if errors.Is(err, storage.ErrWalletNotFound) {
				log.Error("wallet not found")
			}
			handlers.ErrorResponse(w, r, http.StatusInternalServerError, "internal server error")
			return
		}

		err = helpers.WriteJSON(
			w,
			http.StatusOK,
			helpers.Envelope{"data": response{Wallet: wallet, Status: "success"}},
			nil)

		if err != nil {
			return
		}
	}
}
