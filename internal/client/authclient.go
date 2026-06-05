package client

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type AuthClient struct {
	baseURL    string
	httpClient *http.Client
}

type ValidateResponse struct {
	Status string `json:"status"`
	User   string `json:"user"`
	Error  string `json:"error,omitempty"`
}

func NewAuthClient(url string) *AuthClient {
	return &AuthClient{
		baseURL: url,
		httpClient: &http.Client{
			Timeout: 5 * time.Second,
		},
	}
}

func (c *AuthClient) ValidateToken(ctx context.Context, token string) error {
	req, err := http.NewRequestWithContext(ctx, "GET", c.baseURL+"/auth/validate", nil)
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("auth service unavailable: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		var errResp ValidateResponse
		json.NewDecoder(resp.Body).Decode(&errResp)
		if errResp.Error != "" {
			return fmt.Errorf("%s", errResp.Error)
		}
		return fmt.Errorf("token validation failed with status %d", resp.StatusCode)
	}

	return nil
}
