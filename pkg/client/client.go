package client

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

const (
	// BaseURL is the Ahrefs API v3 base URL
	BaseURL = "https://api.ahrefs.com/v3"

	// DefaultTimeout for HTTP requests
	DefaultTimeout = 60 * time.Second

	// DefaultMaxRetries for failed requests
	DefaultMaxRetries = 3
)

// Client is the Ahrefs API client
type Client struct {
	baseURL    string
	apiKey     string
	httpClient *http.Client
	maxRetries int
}

// Config holds client configuration
type Config struct {
	APIKey     string
	BaseURL    string
	Timeout    time.Duration
	MaxRetries int
}

// NewClient creates a new Ahrefs API client
func NewClient(cfg Config) *Client {
	if cfg.BaseURL == "" {
		cfg.BaseURL = BaseURL
	}
	if cfg.Timeout == 0 {
		cfg.Timeout = DefaultTimeout
	}
	if cfg.MaxRetries == 0 {
		cfg.MaxRetries = DefaultMaxRetries
	}

	return &Client{
		baseURL: cfg.BaseURL,
		apiKey:  cfg.APIKey,
		httpClient: &http.Client{
			Timeout: cfg.Timeout,
		},
		maxRetries: cfg.MaxRetries,
	}
}

// Request represents an API request
type Request struct {
	Method   string
	Endpoint string
	Params   url.Values
}

// Response represents an API response with metadata
type Response struct {
	StatusCode int
	Body       []byte
	Headers    http.Header
	Meta       ResponseMeta
}

// ResponseMeta contains metadata about the API response
type ResponseMeta struct {
	UnitsConsumed      int   `json:"units_consumed,omitempty"`
	RateLimitRemaining int   `json:"rate_limit_remaining,omitempty"`
	ResponseTimeMS     int64 `json:"response_time_ms"`
}

// Do executes an API request
func (c *Client) Do(ctx context.Context, req Request) (*Response, error) {
	if c.apiKey == "" {
		return nil, fmt.Errorf("API key is required")
	}

	// Build URL
	u, err := url.Parse(c.baseURL + req.Endpoint)
	if err != nil {
		return nil, fmt.Errorf("invalid endpoint: %w", err)
	}
	if req.Params != nil {
		u.RawQuery = req.Params.Encode()
	}

	var lastErr error
	for attempt := 0; attempt <= c.maxRetries; attempt++ {
		if attempt > 0 {
			// Exponential backoff
			backoff := time.Duration(attempt) * time.Second
			select {
			case <-time.After(backoff):
			case <-ctx.Done():
				return nil, ctx.Err()
			}
		}

		resp, err := c.doRequest(ctx, req.Method, u.String())
		if err == nil {
			return resp, nil
		}

		lastErr = err

		// Don't retry on client errors (4xx) except 429
		if resp != nil && resp.StatusCode >= 400 && resp.StatusCode < 500 && resp.StatusCode != 429 {
			break
		}
	}

	return nil, fmt.Errorf("request failed after %d retries: %w", c.maxRetries, lastErr)
}

// doRequest performs a single HTTP request
func (c *Client) doRequest(ctx context.Context, method, url string) (*Response, error) {
	startTime := time.Now()

	httpReq, err := http.NewRequestWithContext(ctx, method, url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers
	httpReq.Header.Set("Authorization", "Bearer "+c.apiKey)
	httpReq.Header.Set("Accept", "application/json")
	httpReq.Header.Set("User-Agent", "ahrefs-cli/0.1.0")

	httpResp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("HTTP request failed: %w", err)
	}
	defer httpResp.Body.Close()

	body, err := io.ReadAll(httpResp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	responseTime := time.Since(startTime)

	resp := &Response{
		StatusCode: httpResp.StatusCode,
		Body:       body,
		Headers:    httpResp.Header,
		Meta: ResponseMeta{
			ResponseTimeMS: responseTime.Milliseconds(),
		},
	}

	// Parse units consumed from headers if available
	if units := httpResp.Header.Get("X-API-Units-Consumed"); units != "" {
		var unitsInt int
		if _, err := fmt.Sscanf(units, "%d", &unitsInt); err == nil {
			resp.Meta.UnitsConsumed = unitsInt
		}
	}

	if httpResp.StatusCode >= 400 {
		return resp, c.parseError(httpResp.StatusCode, body)
	}

	return resp, nil
}

// APIError represents an error response from the API
type APIError struct {
	StatusCode int
	Code       string
	Message    string
	Suggestion string
	DocsURL    string
}

func (e *APIError) Error() string {
	return fmt.Sprintf("API error (%d): %s", e.StatusCode, e.Message)
}

// parseError attempts to parse an error response
func (c *Client) parseError(statusCode int, body []byte) error {
	apiErr := &APIError{
		StatusCode: statusCode,
	}

	// Try to parse JSON error response
	var errResp struct {
		Error struct {
			Code    string `json:"code"`
			Message string `json:"message"`
		} `json:"error"`
	}

	if err := json.Unmarshal(body, &errResp); err == nil && errResp.Error.Message != "" {
		apiErr.Code = errResp.Error.Code
		apiErr.Message = errResp.Error.Message
	} else {
		// Fallback to status text
		apiErr.Message = string(body)
		if len(apiErr.Message) == 0 {
			apiErr.Message = http.StatusText(statusCode)
		}
	}

	// Add suggestions based on status code
	switch statusCode {
	case http.StatusUnauthorized, http.StatusForbidden:
		apiErr.Code = "AUTH_ERROR"
		apiErr.Suggestion = "Check your API key. Run 'ahrefs config set-key <your-key>' to configure"
		apiErr.DocsURL = "https://docs.ahrefs.com/docs/api/reference/api-keys-creation-and-management"
	case http.StatusTooManyRequests:
		apiErr.Code = "RATE_LIMIT_ERROR"
		apiErr.Suggestion = "Rate limit exceeded. Wait before retrying or check your subscription limits"
		apiErr.DocsURL = "https://docs.ahrefs.com/docs/api/reference/limits-consumption"
	case http.StatusBadRequest:
		apiErr.Code = "VALIDATION_ERROR"
		apiErr.Suggestion = "Check request parameters. Use --describe flag to see valid options"
	case http.StatusNotFound:
		apiErr.Code = "NOT_FOUND"
		apiErr.Suggestion = "Endpoint or resource not found. Verify the target and endpoint"
	}

	return apiErr
}

// Get performs a GET request
func (c *Client) Get(ctx context.Context, endpoint string, params url.Values) (*Response, error) {
	return c.Do(ctx, Request{
		Method:   http.MethodGet,
		Endpoint: endpoint,
		Params:   params,
	})
}

// Post performs a POST request
func (c *Client) Post(ctx context.Context, endpoint string, params url.Values) (*Response, error) {
	return c.Do(ctx, Request{
		Method:   http.MethodPost,
		Endpoint: endpoint,
		Params:   params,
	})
}
