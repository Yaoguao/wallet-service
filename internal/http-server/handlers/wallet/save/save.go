package save

import (
	"context"
	"log/slog"
	"net/http"
	"wallet-service/internal/domain/models"
	"wallet-service/internal/http-server/handlers"
	"wallet-service/pkg/helpers"
)

type WalletSaver interface {
	CreateWallet(ctx context.Context, amount int64) (*models.Wallet, error)
}

type request struct {
	Amount int64 `json:"amount"`
}

type response struct {
	Wallet *models.Wallet `json:"wallet,omitempty"`
	Error  string         `json:"error,omitempty"`
}

func New(log *slog.Logger, ws WalletSaver) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req request
		err := helpers.ReadJSON(w, r, &req)
		if err != nil {
			log.Error("failed to decode request body")

			handlers.ErrorResponse(w, r, http.StatusInternalServerError, "ops! decode json")
			return
		}

		wallet, err := ws.CreateWallet(r.Context(), req.Amount)
		if err != nil {
			handlers.ErrorResponse(w, r, http.StatusInternalServerError, "error create wallet")
			return
		}

		err = helpers.WriteJSON(
			w,
			http.StatusOK,
			helpers.Envelope{"data": response{Wallet: wallet}},
			nil)

		if err != nil {
			return
		}
	}
}
