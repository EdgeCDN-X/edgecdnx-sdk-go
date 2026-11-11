package edgecdnx

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

type recordedRequest struct {
	Method string
	Path   string
	Auth   string
	Body   []byte
}

// newTestClient starts a test server that records the request and replies
// with status and respBody (JSON-encoded unless nil).
func newTestClient(t *testing.T, status int, respBody any, record *recordedRequest, opts ...Option) *Client {
	t.Helper()
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		record.Method = r.Method
		record.Path = r.URL.Path
		record.Auth = r.Header.Get("Authorization")
		if r.Body != nil {
			body, _ := io.ReadAll(r.Body)
			record.Body = body
		}
		if respBody == nil {
			w.WriteHeader(status)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(status)
		_ = json.NewEncoder(w).Encode(respBody)
	}))
	t.Cleanup(server.Close)

	client, err := NewClient(server.URL, opts...)
	if err != nil {
		t.Fatalf("NewClient: %v", err)
	}
	return client
}

func TestNewClientValidatesBaseURL(t *testing.T) {
	if _, err := NewClient("not-a-url"); err == nil {
		t.Fatal("expected error for base URL without scheme")
	}
	if _, err := NewClient("://bad"); err == nil {
		t.Fatal("expected error for unparsable base URL")
	}
	c, err := NewClient("https://api.example.com/")
	if err != nil {
		t.Fatalf("NewClient: %v", err)
	}
	if got := c.baseURL.String(); got != "https://api.example.com" {
		t.Fatalf("trailing slash not trimmed, got %q", got)
	}
}

func TestAPIKeyAuthHeader(t *testing.T) {
	var rec recordedRequest
	client := newTestClient(t, 200, []struct{}{}, &rec, WithAPIKey("ak_abc_secret"))
	if _, err := client.Projects.List(context.Background()); err != nil {
		t.Fatalf("List: %v", err)
	}
	if rec.Auth != "ApiKey ak_abc_secret" {
		t.Fatalf("unexpected Authorization header %q", rec.Auth)
	}
}

func TestBearerTokenAuthHeader(t *testing.T) {
	var rec recordedRequest
	client := newTestClient(t, 200, []struct{}{}, &rec, WithBearerToken("tok"))
	if _, err := client.Projects.List(context.Background()); err != nil {
		t.Fatalf("List: %v", err)
	}
	if rec.Auth != "Bearer tok" {
		t.Fatalf("unexpected Authorization header %q", rec.Auth)
	}
}

func TestTokenSourceAuthHeader(t *testing.T) {
	var rec recordedRequest
	client := newTestClient(t, 200, []struct{}{}, &rec, WithTokenSource(func(ctx context.Context) (string, error) {
		return "dynamic", nil
	}))
	if _, err := client.Projects.List(context.Background()); err != nil {
		t.Fatalf("List: %v", err)
	}
	if rec.Auth != "Bearer dynamic" {
		t.Fatalf("unexpected Authorization header %q", rec.Auth)
	}
}

func TestAPIErrorDecoding(t *testing.T) {
	var rec recordedRequest
	client := newTestClient(t, 404, map[string]string{"error": "project not found"}, &rec)
	_, err := client.Projects.Get(context.Background(), "missing")
	if err == nil {
		t.Fatal("expected error")
	}
	if !IsNotFound(err) {
		t.Fatalf("expected IsNotFound, got %v", err)
	}
	apiErr, ok := err.(*APIError)
	if !ok {
		t.Fatalf("expected *APIError, got %T", err)
	}
	if apiErr.Message != "project not found" {
		t.Fatalf("unexpected message %q", apiErr.Message)
	}
}

func TestAPIErrorWithReason(t *testing.T) {
	var rec recordedRequest
	client := newTestClient(t, 409, map[string]any{"error": "exists", "reason": "AlreadyExists"}, &rec)
	_, err := client.Zones.Create(context.Background(), "proj", CreateZoneRequest{Email: "a@b.com", Zone: "example.com"})
	if !IsConflict(err) {
		t.Fatalf("expected IsConflict, got %v", err)
	}
	apiErr := err.(*APIError)
	if apiErr.Reason != "AlreadyExists" {
		t.Fatalf("unexpected reason %q", apiErr.Reason)
	}
}

func TestUpdateAPIKeyRequestMarshalling(t *testing.T) {
	// Empty request marshals to an empty object (all fields unchanged).
	data, err := json.Marshal(UpdateAPIKeyRequest{})
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	if string(data) != "{}" {
		t.Fatalf("expected {}, got %s", data)
	}

	// ClearExpiry sends an explicit null.
	var req UpdateAPIKeyRequest
	req.ClearExpiry()
	data, err = json.Marshal(req)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	if string(data) != `{"expiresAt":null}` {
		t.Fatalf("expected explicit null expiresAt, got %s", data)
	}

	// SetExpiry sends the timestamp; other set fields are included.
	var full UpdateAPIKeyRequest
	full.DisplayName = Ptr("new name")
	full.Roles = &[]string{}
	full.SetExpiry(time.Date(2030, 1, 2, 3, 4, 5, 0, time.UTC))
	data, err = json.Marshal(full)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	var decoded map[string]any
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if decoded["displayName"] != "new name" {
		t.Fatalf("unexpected displayName %v", decoded["displayName"])
	}
	if decoded["expiresAt"] != "2030-01-02T03:04:05Z" {
		t.Fatalf("unexpected expiresAt %v", decoded["expiresAt"])
	}
	roles, ok := decoded["roles"].([]any)
	if !ok || len(roles) != 0 {
		t.Fatalf("expected empty roles array, got %v", decoded["roles"])
	}
}
