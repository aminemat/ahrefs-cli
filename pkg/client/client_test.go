package client

import (
	"context"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestNewClient(t *testing.T) {
	tests := []struct {
		name   string
		config Config
		want   string
	}{
		{
			name:   "default base URL",
			config: Config{APIKey: "test-key"},
			want:   BaseURL,
		},
		{
			name:   "custom base URL",
			config: Config{APIKey: "test-key", BaseURL: "https://custom.api.com"},
			want:   "https://custom.api.com",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := NewClient(tt.config)
			if c.baseURL != tt.want {
				t.Errorf("NewClient() baseURL = %v, want %v", c.baseURL, tt.want)
			}
			if c.apiKey != tt.config.APIKey {
				t.Errorf("NewClient() apiKey = %v, want %v", c.apiKey, tt.config.APIKey)
			}
		})
	}
}

func TestClient_Get(t *testing.T) {
	tests := []struct {
		name           string
		apiKey         string
		endpoint       string
		params         url.Values
		serverStatus   int
		serverBody     string
		wantErr        bool
		wantStatusCode int
	}{
		{
			name:           "successful request",
			apiKey:         "test-key",
			endpoint:       "/test",
			params:         url.Values{"target": []string{"example.com"}},
			serverStatus:   http.StatusOK,
			serverBody:     `{"result":"success"}`,
			wantErr:        false,
			wantStatusCode: http.StatusOK,
		},
		{
			name:           "unauthorized",
			apiKey:         "invalid-key",
			endpoint:       "/test",
			params:         nil,
			serverStatus:   http.StatusUnauthorized,
			serverBody:     `{"error":{"message":"Invalid API key"}}`,
			wantErr:        true,
			wantStatusCode: http.StatusUnauthorized,
		},
		{
			name:           "rate limited",
			apiKey:         "test-key",
			endpoint:       "/test",
			params:         nil,
			serverStatus:   http.StatusTooManyRequests,
			serverBody:     `{"error":{"message":"Rate limit exceeded"}}`,
			wantErr:        true,
			wantStatusCode: http.StatusTooManyRequests,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create test server
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				// Verify authorization header
				auth := r.Header.Get("Authorization")
				expectedAuth := "Bearer " + tt.apiKey
				if auth != expectedAuth {
					t.Errorf("Authorization header = %v, want %v", auth, expectedAuth)
				}

				// Verify accept header
				accept := r.Header.Get("Accept")
				if accept != "application/json" {
					t.Errorf("Accept header = %v, want application/json", accept)
				}

				// Return mock response
				w.WriteHeader(tt.serverStatus)
				w.Write([]byte(tt.serverBody))
			}))
			defer server.Close()

			// Create client with test server URL
			c := NewClient(Config{
				APIKey:  tt.apiKey,
				BaseURL: server.URL,
			})

			// Make request
			resp, err := c.Get(context.Background(), tt.endpoint, tt.params)

			// Check error
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			// Check status code
			if resp != nil && resp.StatusCode != tt.wantStatusCode {
				t.Errorf("Client.Get() statusCode = %v, want %v", resp.StatusCode, tt.wantStatusCode)
			}

			// Check response body for successful requests
			if !tt.wantErr && resp != nil {
				if string(resp.Body) != tt.serverBody {
					t.Errorf("Client.Get() body = %v, want %v", string(resp.Body), tt.serverBody)
				}
			}
		})
	}
}

func TestClient_NoAPIKey(t *testing.T) {
	c := NewClient(Config{
		APIKey: "",
	})

	_, err := c.Get(context.Background(), "/test", nil)
	if err == nil {
		t.Error("Client.Get() with empty API key should return error")
	}

	expected := "API key is required"
	if err.Error() != expected {
		t.Errorf("Client.Get() error = %v, want %v", err.Error(), expected)
	}
}

func TestAPIError(t *testing.T) {
	tests := []struct {
		name           string
		statusCode     int
		wantCode       string
		wantSuggestion bool
	}{
		{
			name:           "unauthorized",
			statusCode:     http.StatusUnauthorized,
			wantCode:       "AUTH_ERROR",
			wantSuggestion: true,
		},
		{
			name:           "rate limit",
			statusCode:     http.StatusTooManyRequests,
			wantCode:       "RATE_LIMIT_ERROR",
			wantSuggestion: true,
		},
		{
			name:           "bad request",
			statusCode:     http.StatusBadRequest,
			wantCode:       "VALIDATION_ERROR",
			wantSuggestion: true,
		},
		{
			name:           "not found",
			statusCode:     http.StatusNotFound,
			wantCode:       "NOT_FOUND",
			wantSuggestion: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Client{}
			err := c.parseError(tt.statusCode, []byte("test error"))

			apiErr, ok := err.(*APIError)
			if !ok {
				t.Fatalf("parseError() should return *APIError, got %T", err)
			}

			if apiErr.Code != tt.wantCode {
				t.Errorf("APIError.Code = %v, want %v", apiErr.Code, tt.wantCode)
			}

			if tt.wantSuggestion && apiErr.Suggestion == "" {
				t.Error("APIError.Suggestion should not be empty")
			}
		})
	}
}

func TestClient_Retries(t *testing.T) {
	attempts := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		attempts++
		if attempts < 3 {
			// Fail first 2 attempts with 500
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		// Succeed on 3rd attempt
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"result":"success"}`))
	}))
	defer server.Close()

	c := NewClient(Config{
		APIKey:     "test-key",
		BaseURL:    server.URL,
		MaxRetries: 3,
	})

	resp, err := c.Get(context.Background(), "/test", nil)
	if err != nil {
		t.Errorf("Client.Get() with retries should succeed, got error: %v", err)
	}

	if attempts != 3 {
		t.Errorf("Expected 3 attempts, got %d", attempts)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}
}

func TestClient_NoRetryOn4xx(t *testing.T) {
	attempts := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		attempts++
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"error":{"message":"bad request"}}`))
	}))
	defer server.Close()

	c := NewClient(Config{
		APIKey:     "test-key",
		BaseURL:    server.URL,
		MaxRetries: 3,
	})

	_, err := c.Get(context.Background(), "/test", nil)
	if err == nil {
		t.Error("Client.Get() with 400 should return error")
	}

	if attempts != 1 {
		t.Errorf("Expected 1 attempt (no retries on 4xx), got %d", attempts)
	}
}
