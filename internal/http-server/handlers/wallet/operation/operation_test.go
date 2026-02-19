package operation

import (
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"wallet-service/internal/domain/models"
	"wallet-service/internal/http-server/handlers/wallet/operation/mocks"
	"wallet-service/internal/services"
	"wallet-service/internal/storage"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestOperationHandler(t *testing.T) {
	t.Run("success deposit", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockService := mocks.NewMockWalletService(ctrl)
		logger := slog.New(slog.NewTextHandler(io.Discard, nil))

		walletID := uuid.New()
		amount := int64(100)
		expectedWallet := &models.Wallet{
			ID:      walletID,
			Balance: 100,
		}

		mockService.
			EXPECT().
			Deposit(gomock.Any(), walletID, amount).
			Return(expectedWallet, nil)

		handler := New(logger, mockService)

		body := fmt.Sprintf(`{"wallet_id":"%s","operation_type":"DEPOSIT","amount":%d}`, walletID, amount)
		req := httptest.NewRequest(http.MethodPost, "/operations", strings.NewReader(body))
		w := httptest.NewRecorder()

		handler.ServeHTTP(w, req)

		require.Equal(t, http.StatusOK, w.Code)
		require.Contains(t, w.Body.String(), walletID.String())
		require.Contains(t, w.Body.String(), "operation was completed successfully")
	})

	t.Run("success withdraw", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockService := mocks.NewMockWalletService(ctrl)
		logger := slog.New(slog.NewTextHandler(io.Discard, nil))

		walletID := uuid.New()
		amount := int64(50)
		expectedWallet := &models.Wallet{
			ID:      walletID,
			Balance: 50,
		}

		mockService.
			EXPECT().
			Withdraw(gomock.Any(), walletID, amount).
			Return(expectedWallet, nil)

		handler := New(logger, mockService)

		body := fmt.Sprintf(`{"wallet_id":"%s","operation_type":"WITHDRAW","amount":%d}`, walletID, amount)
		req := httptest.NewRequest(http.MethodPost, "/operations", strings.NewReader(body))
		w := httptest.NewRecorder()

		handler.ServeHTTP(w, req)

		require.Equal(t, http.StatusOK, w.Code)
		require.Contains(t, w.Body.String(), walletID.String())
	})

	t.Run("invalid JSON", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockService := mocks.NewMockWalletService(ctrl)
		logger := slog.New(slog.NewTextHandler(io.Discard, nil))

		handler := New(logger, mockService)

		req := httptest.NewRequest(http.MethodPost, "/operations", strings.NewReader(`{"wallet_id":}`))
		w := httptest.NewRecorder()

		handler.ServeHTTP(w, req)

		require.Equal(t, http.StatusInternalServerError, w.Code)
	})

	t.Run("invalid operation type", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockService := mocks.NewMockWalletService(ctrl)
		logger := slog.New(slog.NewTextHandler(io.Discard, nil))

		handler := New(logger, mockService)

		walletID := uuid.New()
		reqBody := fmt.Sprintf(`{"wallet_id":"%s","operation_type":"UNKNOWN","amount":100}`, walletID)
		req := httptest.NewRequest(http.MethodPost, "/operations", strings.NewReader(reqBody))
		w := httptest.NewRecorder()

		handler.ServeHTTP(w, req)

		require.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("service insufficient funds", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockService := mocks.NewMockWalletService(ctrl)
		logger := slog.New(slog.NewTextHandler(io.Discard, nil))

		walletID := uuid.New()
		amount := int64(200)

		mockService.
			EXPECT().
			Withdraw(gomock.Any(), walletID, amount).
			Return(nil, storage.ErrInsufficientFunds)

		handler := New(logger, mockService)

		reqBody := fmt.Sprintf(`{"wallet_id":"%s","operation_type":"WITHDRAW","amount":%d}`, walletID, amount)
		req := httptest.NewRequest(http.MethodPost, "/operations", strings.NewReader(reqBody))
		w := httptest.NewRecorder()

		handler.ServeHTTP(w, req)

		require.Equal(t, http.StatusBadRequest, w.Code)
		require.Contains(t, w.Body.String(), "insufficient funds")
	})

	t.Run("service negative amount", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockService := mocks.NewMockWalletService(ctrl)
		logger := slog.New(slog.NewTextHandler(io.Discard, nil))

		walletID := uuid.New()
		amount := int64(-100)

		mockService.
			EXPECT().
			Deposit(gomock.Any(), walletID, amount).
			Return(nil, services.ErrAmountNegativeValue)

		handler := New(logger, mockService)

		reqBody := fmt.Sprintf(`{"wallet_id":"%s","operation_type":"DEPOSIT","amount":%d}`, walletID, amount)
		req := httptest.NewRequest(http.MethodPost, "/operations", strings.NewReader(reqBody))
		w := httptest.NewRecorder()

		handler.ServeHTTP(w, req)

		require.Equal(t, http.StatusBadRequest, w.Code)
		require.Contains(t, w.Body.String(), "amount negative value")
	})
}
