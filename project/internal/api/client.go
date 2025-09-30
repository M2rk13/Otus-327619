package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type Client struct {
	httpClient *http.Client
	apiKey     string
	baseURL    string
}

func NewClient(apiKey string) *Client {
	return &Client{
		httpClient: &http.Client{Timeout: 10 * time.Second},
		apiKey:     apiKey,
		baseURL:    "https://api.exchangerate.host",
	}
}

type exchangeRateResponse struct {
	Success bool `json:"success"`
	Info    struct {
		Rate float64 `json:"rate"`
	} `json:"info"`
	Result float64 `json:"result"`
	Error  struct {
		Code string `json:"code"`
		Info string `json:"info"`
	} `json:"error"`
}

func (c *Client) GetConversionRate(from, to string, amount float64) (float64, float64, error) {
	url := fmt.Sprintf("%s/convert?from=%s&to=%s&amount=%f&access_key=%s", c.baseURL, from, to, amount, c.apiKey)

	req, err := http.NewRequest("GET", url, nil)

	if err != nil {
		return 0, 0, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := c.httpClient.Do(req)

	if err != nil {
		return 0, 0, fmt.Errorf("failed to execute request: %w", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return 0, 0, fmt.Errorf("external API returned status: %s", resp.Status)
	}

	var apiResp exchangeRateResponse

	body, err := io.ReadAll(resp.Body)
	resp.Body.Close()

	if err != nil {
		return 0, 0, fmt.Errorf("failed to read response: %w", err)
	}

	if err := json.Unmarshal(body, &apiResp); err != nil {
		return 0, 0, fmt.Errorf("failed to decode response: %s", string(body))
	}

	if !apiResp.Success {
		return 0, 0, fmt.Errorf("external API error: %s - %s", apiResp.Error.Code, apiResp.Error.Info)
	}

	return apiResp.Result, apiResp.Info.Rate, nil
}
