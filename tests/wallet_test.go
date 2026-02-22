package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"
	"wallet-service/internal/domain/models"
	"wallet-service/tests/suite"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDepositFlow(t *testing.T) {
	ctx, st := suite.New(t)

	req := struct {
		Amount int64 `json:"amount"`
	}{}

	payload, err := json.Marshal(req)
	require.NoError(t, err)

	resp, err := st.Request(ctx, http.MethodPost, "/wallets", bytes.NewReader(payload))
	require.NoError(t, err)

	var res struct {
		Data struct {
			Wallet models.Wallet `json:"wallet"`
		} `json:"data"`
	}

	err = json.NewDecoder(resp.Body).Decode(&res)
	require.NoError(t, err)

	walletCreated := res.Data.Wallet

	var deposit struct {
		WalletID      uuid.UUID `json:"wallet_id"`
		Amount        int64     `json:"amount"`
		OperationType string    `json:"operation_type"`
	}

	deposit.WalletID = walletCreated.ID
	deposit.Amount = 100
	deposit.OperationType = "DEPOSIT"

	payload, err = json.Marshal(deposit)
	require.NoError(t, err)

	resp, err = st.Request(
		ctx,
		http.MethodPost,
		"/wallets/operation",
		bytes.NewReader(payload),
	)

	require.NoError(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)

	resp, err = st.Request(ctx, http.MethodGet,
		"/wallets/"+walletCreated.ID.String(),
		nil,
	)

	require.NoError(t, err)

	err = json.NewDecoder(resp.Body).Decode(&res)
	require.NoError(t, err)

	walletUpdate := res.Data.Wallet

	assert.Equal(t, deposit.Amount, walletUpdate.Balance)
}
