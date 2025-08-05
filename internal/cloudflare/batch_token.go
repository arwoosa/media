package cloudflare

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/arwoosa/media/internal/cloudflare/dao"
)

var batchToken *dao.BatchToken

func getBatchToken(ctx context.Context) (*string, error) {
	// If token exists and is not expired, return it
	if batchToken != nil && time.Now().Before(batchToken.ExpiresAt) {
		return &batchToken.Token, nil
	}

	if err := checkConfig(); err != nil {
		return nil, err
	}

	url := fmt.Sprintf("https://api.cloudflare.com/client/v4/accounts/%s/images/v1/batch_token", accountID)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrCloudflareCallFailed, err)
	}

	req.Header.Set("Authorization", "Bearer "+apiToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrCloudflareCallFailed, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		// Read the response body
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("%w: received non-200 status code: %d, body: %s", ErrCloudflareCallFailed, resp.StatusCode, string(body))
	}

	var result struct {
		Result struct {
			Token     string    `json:"token"`
			ExpiresAt time.Time `json:"expiresAt"`
		} `json:"result"`
		Success bool `json:"success"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("%w: %w", ErrCloudflareCallFailed, err)
	}

	if !result.Success {
		return nil, fmt.Errorf("%w: cloudflare api returned success=false", ErrCloudflareCallFailed)
	}

	// Store the new token and its expiration
	batchToken = &dao.BatchToken{
		Token:     result.Result.Token,
		ExpiresAt: result.Result.ExpiresAt,
	}

	return &batchToken.Token, nil
}
