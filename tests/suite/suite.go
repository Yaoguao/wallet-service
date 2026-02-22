package suite

import (
	"context"
	"io"
	"net/http"
	"os"
	"testing"
	"time"
	"wallet-service/internal/config"

	"github.com/stretchr/testify/require"
)

type Suite struct {
	T   *testing.T
	Cfg *config.Config

	Client  *http.Client
	BaseURL string
}

func New(t *testing.T) (context.Context, *Suite) {
	t.Helper()
	t.Parallel()

	cfg := config.MustLoadPath(configPath())

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	t.Cleanup(cancel)

	client := &http.Client{
		Timeout: time.Second * 10,
	}

	suite := &Suite{
		T:       t,
		Cfg:     cfg,
		Client:  client,
		BaseURL: "http://" + cfg.HTTPServer.Address + "/api/v1",
	}

	return ctx, suite
}

func configPath() (string, string) {
	const key = "CONFIG_PATH"
	const keyDB = "DSN_POSTGRES"

	cp := os.Getenv(key)
	dsn := os.Getenv(keyDB)

	if cp != "" && dsn != "" {
		return cp, dsn
	}

	return "../config/config.yml", "postgres://walletdb:wallet@local-host:5432/walletdb?sslmode=disable"
}

func (s *Suite) Request(
	ctx context.Context,
	method string,
	path string,
	body io.Reader,
) (*http.Response, error) {

	req, err := http.NewRequestWithContext(ctx, method, s.BaseURL+path, body)
	require.NoError(s.T, err)

	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	if err != nil {
		return nil, err
	}

	return s.Client.Do(req)
}
