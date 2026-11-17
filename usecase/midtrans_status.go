package usecase

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/capstone-orang-jago/parkin-be/dto"
)

type MidtransStatusProvider interface {
	GetTransactionStatus(ctx context.Context, orderID string) (*dto.MidtransStatusResponseDTO, error)
}

type httpMidtransStatusProvider struct {
	serverKey string
	baseURL   string
	client    *http.Client
}

func NewMidtransStatusProvider(serverKey string) MidtransStatusProvider {
	return &httpMidtransStatusProvider{
		serverKey: serverKey,
		baseURL:   "https://api.sandbox.midtrans.com",
		client:    &http.Client{Timeout: 10 * time.Second},
	}
}

func (provider *httpMidtransStatusProvider) GetTransactionStatus(ctx context.Context, orderID string) (*dto.MidtransStatusResponseDTO, error) {
	endpoint := fmt.Sprintf("%s/v2/%s/status", provider.baseURL, url.PathEscape(orderID))

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, err
	}

	authHeader := base64.StdEncoding.EncodeToString([]byte(provider.serverKey + ":"))
	req.Header.Set("Authorization", "Basic "+authHeader)
	req.Header.Set("Accept", "application/json")

	resp, err := provider.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("midtrans status request failed with status %d", resp.StatusCode)
	}

	var status dto.MidtransStatusResponseDTO
	if err := json.NewDecoder(resp.Body).Decode(&status); err != nil {
		return nil, err
	}

	return &status, nil
}
