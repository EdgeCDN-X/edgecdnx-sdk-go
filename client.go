// Package edgecdnx provides a Go client for the EdgeCDN-X API
// (github.com/EdgeCDN-X/edgecdnx-api).
//
// Create a client with NewClient and access resources through the
// service fields, for example:
//
//	client, err := edgecdnx.NewClient("https://api.edgecdnx.com", edgecdnx.WithAPIKey(token))
//	if err != nil { ... }
//	projects, err := client.Projects.List(ctx)
package edgecdnx

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

const defaultUserAgent = "edgecdnx-sdk-go"

// Credentials injects authentication information into outgoing requests.
type Credentials interface {
	// Apply sets authentication information (typically the Authorization
	// header) on the request before it is sent.
	Apply(req *http.Request) error
}

// APIKeyCredentials authenticates using an EdgeCDN-X API key
// (Authorization: ApiKey ak_<id>_<secret>). API keys are confined to a
// single project and cannot create projects.
type APIKeyCredentials struct {
	Token string
}

// Apply implements Credentials.
func (a APIKeyCredentials) Apply(req *http.Request) error {
	req.Header.Set("Authorization", "ApiKey "+a.Token)
	return nil
}

// BearerTokenCredentials authenticates using a static OIDC bearer token
// (Authorization: Bearer <token>).
type BearerTokenCredentials struct {
	Token string
}

// Apply implements Credentials.
func (b BearerTokenCredentials) Apply(req *http.Request) error {
	req.Header.Set("Authorization", "Bearer "+b.Token)
	return nil
}

// TokenSourceCredentials fetches a fresh OIDC bearer token for every request.
// Useful when tokens are short-lived and refreshed externally.
type TokenSourceCredentials struct {
	Source func(ctx context.Context) (string, error)
}

// Apply implements Credentials.
func (t TokenSourceCredentials) Apply(req *http.Request) error {
	if t.Source == nil {
		return fmt.Errorf("edgecdnx: token source is nil")
	}
	token, err := t.Source(req.Context())
	if err != nil {
		return fmt.Errorf("edgecdnx: fetch token: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+token)
	return nil
}

// Client is an EdgeCDN-X API client. Use NewClient to construct one.
type Client struct {
	baseURL    *url.URL
	httpClient *http.Client
	userAgent  string
	creds      Credentials

	// Projects manages EdgeCDN-X projects.
	Projects *ProjectsService
	// Zones manages DNS zones within a project.
	Zones *ZonesService
	// DNSRecords manages DNS records (DNS endpoints) within a zone.
	DNSRecords *DNSRecordsService
	// APIKeys manages API keys within a project.
	APIKeys *APIKeysService
	// Locations manages edge locations within a project.
	Locations *LocationsService
	// HealthCheckProfiles manages health check profiles within a project.
	HealthCheckProfiles *HealthCheckProfilesService
}

// Option customizes a Client.
type Option func(*Client)

// WithHTTPClient sets the underlying *http.Client.
func WithHTTPClient(hc *http.Client) Option {
	return func(c *Client) { c.httpClient = hc }
}

// WithUserAgent overrides the User-Agent header.
func WithUserAgent(ua string) Option {
	return func(c *Client) { c.userAgent = ua }
}

// WithCredentials sets a custom Credentials implementation.
func WithCredentials(creds Credentials) Option {
	return func(c *Client) { c.creds = creds }
}

// WithAPIKey authenticates with an EdgeCDN-X API key (ak_<id>_<secret>).
func WithAPIKey(token string) Option {
	return WithCredentials(APIKeyCredentials{Token: token})
}

// WithBearerToken authenticates with a static OIDC bearer token.
func WithBearerToken(token string) Option {
	return WithCredentials(BearerTokenCredentials{Token: token})
}

// WithTokenSource authenticates with a bearer token fetched per request.
func WithTokenSource(source func(ctx context.Context) (string, error)) Option {
	return WithCredentials(TokenSourceCredentials{Source: source})
}

// NewClient creates a Client for the EdgeCDN-X API at baseURL,
// for example "https://api.edgecdnx.com".
func NewClient(baseURL string, opts ...Option) (*Client, error) {
	u, err := url.Parse(strings.TrimRight(baseURL, "/"))
	if err != nil {
		return nil, fmt.Errorf("edgecdnx: invalid base URL %q: %w", baseURL, err)
	}
	if u.Scheme == "" || u.Host == "" {
		return nil, fmt.Errorf("edgecdnx: base URL %q must include scheme and host", baseURL)
	}

	c := &Client{
		baseURL:    u,
		httpClient: http.DefaultClient,
		userAgent:  defaultUserAgent,
	}
	for _, opt := range opts {
		opt(c)
	}

	c.Projects = &ProjectsService{client: c}
	c.Zones = &ZonesService{client: c}
	c.DNSRecords = &DNSRecordsService{client: c}
	c.APIKeys = &APIKeysService{client: c}
	c.Locations = &LocationsService{client: c}
	c.HealthCheckProfiles = &HealthCheckProfilesService{client: c}
	return c, nil
}

// do performs an API request. body (if non-nil) is JSON-encoded; a successful
// response is decoded into out (if non-nil). Non-2xx responses are returned
// as *APIError.
func (c *Client) do(ctx context.Context, method, path string, body, out any) error {
	var reqBody io.Reader
	if body != nil {
		buf := &bytes.Buffer{}
		if err := json.NewEncoder(buf).Encode(body); err != nil {
			return fmt.Errorf("edgecdnx: encode request body: %w", err)
		}
		reqBody = buf
	}

	req, err := http.NewRequestWithContext(ctx, method, c.baseURL.String()+path, reqBody)
	if err != nil {
		return fmt.Errorf("edgecdnx: create request: %w", err)
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", c.userAgent)
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	if c.creds != nil {
		if err := c.creds.Apply(req); err != nil {
			return err
		}
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("edgecdnx: %s %s: %w", method, path, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return newAPIError(resp)
	}

	if out == nil || resp.StatusCode == http.StatusNoContent {
		// Drain so the connection can be reused.
		_, _ = io.Copy(io.Discard, resp.Body)
		return nil
	}

	if err := json.NewDecoder(resp.Body).Decode(out); err != nil {
		return fmt.Errorf("edgecdnx: decode response of %s %s: %w", method, path, err)
	}
	return nil
}

// pathEscape escapes a path segment such as a project or resource name.
func pathEscape(segment string) string {
	return url.PathEscape(segment)
}

// Ptr returns a pointer to v. Convenient for optional request fields,
// for example edgecdnx.Ptr(300) for a DNS record TTL.
func Ptr[T any](v T) *T {
	return &v
}
