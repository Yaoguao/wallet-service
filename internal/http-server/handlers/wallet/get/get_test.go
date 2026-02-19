package get

import (
	"context"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"
	"wallet-service/internal/domain/models"
	"wallet-service/internal/http-server/handlers/wallet/get/mocks"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"go.uber.org/mock/gomock"
)

func TestOperationHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockGetter := mocks.NewMockWalletGetter(ctrl)

	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	handler := New(logger, mockGetter)

	id := uuid.New()

	req := httptest.NewRequest(http.MethodGet, "/wallets/"+id.String(), nil)

	rr := httptest.NewRecorder()

	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("WALLET_UUID", id.String())

	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	mockGetter.
		EXPECT().
		GetWallet(gomock.Any(), id).
		Return(&models.Wallet{ID: id}, nil)

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rr.Code)
	}
}
