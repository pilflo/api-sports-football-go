// Package api defines structs and functions for requesting API-Football enpoints.
package api

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"
)

// SubscriptionType is a custom type representing the subscription type to api-football.
type SubscriptionType string

const (
	// SubTypeAPISports is the api type if you subscribed to api-football through API Sports.
	SubTypeAPISports SubscriptionType = "APISports"
	// SubTypeRapidAPI is the api type if you subscribed to api-football through Rapid API.
	SubTypeRapidAPI SubscriptionType = "RapidAPI"
)

// ClientError for any client-related error.
type ClientError error

// ErrAPIKeyEmpty is returns when no API key could be found.
var ErrAPIKeyEmpty ClientError = fmt.Errorf("API Key must be non empty")

// Client represents the base client requester.
type Client struct {
	config     *config
	logger     *slog.Logger
	httpClient *http.Client
}

// tokenRoundTripper implements http.RoundTripper interface.
// It injects the token in headers.
type tokenRoundTripper struct {
	next   http.RoundTripper
	config *config
}

// authMiddleware returns an http.RoundTripper that can be used by an http.Client to inject Authorization header.
func authMiddleware(next http.RoundTripper, config *config) http.RoundTripper {
	return &tokenRoundTripper{next, config}
}

func (rt *tokenRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	apiKeyValue := os.Getenv(rt.config.apiKeyEnvVar)
	if apiKeyValue == "" {
		return nil, ErrAPIKeyEmpty
	}

	req.Header.Add(rt.config.apiKeyHTTPHeader, apiKeyValue)

	res, err := rt.next.RoundTrip(req)
	if err != nil {
		return nil, fmt.Errorf("failed to call next roundtripper: %w", err)
	}

	return res, nil
}

// NewClient returns a ready-to-use *Client for making requests to the API.
func NewClient(subType SubscriptionType) *Client {
	conf := newConfig(subType)

	httpClient := http.Client{
		Transport: authMiddleware(http.DefaultTransport, &conf),
	}

	return &Client{&conf, slog.Default(), &httpClient}
}

// WithCustomAPIURL allow to bring a custom slog.Logger to the library.
func (c *Client) WithCustomAPIURL(url string) *Client {
	c.config.basePath = url

	return c
}

// WithCustomLogger allow to bring a custom slog.Logger to the library.
func (c *Client) WithCustomLogger(logger *slog.Logger) *Client {
	c.logger = logger

	return c
}

func (c *Client) String() string {
	return fmt.Sprintf("Client [Type = %s, BasePath = %s, ApiKeyEnv = %s]", c.config.subType, c.config.basePath, c.config.apiKeyEnvVar)
}
