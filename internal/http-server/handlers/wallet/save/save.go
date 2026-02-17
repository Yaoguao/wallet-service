package save

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"wallet-service/internal/domain/models"
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
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			log.Error("failed to decode request body")

			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(response{Error: "invalid request body"})
			return
		}

		wallet, err := ws.CreateWallet(r.Context(), req.Amount)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(response{Error: "ops!"})
			return
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response{Wallet: wallet})
	}
}
