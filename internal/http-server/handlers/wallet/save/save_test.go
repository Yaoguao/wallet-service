package save

import (
	"errors"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"wallet-service/internal/domain/models"
	"wallet-service/internal/http-server/handlers/wallet/save/mocks"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestCreateWalletHandler(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockService := mocks.NewMockWalletSaver(ctrl)

		logger := slog.New(slog.NewTextHandler(io.Discard, nil))

		expectedWallet := &models.Wallet{
			ID:      uuid.New(),
			Balance: 100,
		}

		mockService.
			EXPECT().
			CreateWallet(gomock.Any(), int64(100)).
			Return(expectedWallet, nil)

		handler := New(logger, mockService)

		body := `{"amount":100}`
		req := httptest.NewRequest(http.MethodPost, "/wallets", strings.NewReader(body))
		w := httptest.NewRecorder()

		handler.ServeHTTP(w, req)

		require.Equal(t, http.StatusOK, w.Code)

		require.Contains(t, w.Body.String(), expectedWallet.ID.String())
	})

	t.Run("service error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockService := mocks.NewMockWalletSaver(ctrl)
		logger := slog.New(slog.NewTextHandler(io.Discard, nil))

		mockService.
			EXPECT().
			CreateWallet(gomock.Any(), int64(100)).
			Return(nil, errors.New("service error"))

		handler := New(logger, mockService)

		body := `{"amount":100}`
		req := httptest.NewRequest(http.MethodPost, "/wallets", strings.NewReader(body))
		w := httptest.NewRecorder()

		handler.ServeHTTP(w, req)

		require.Equal(t, http.StatusInternalServerError, w.Code)
	})

	t.Run("invalid json", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockService := mocks.NewMockWalletSaver(ctrl)
		logger := slog.New(slog.NewTextHandler(io.Discard, nil))

		handler := New(logger, mockService)

		body := `{"amount":}`
		req := httptest.NewRequest(http.MethodPost, "/wallets", strings.NewReader(body))
		w := httptest.NewRecorder()

		handler.ServeHTTP(w, req)

		require.Equal(t, http.StatusInternalServerError, w.Code)

	})
}
