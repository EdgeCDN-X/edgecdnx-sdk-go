package edgecdnx

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestDeleteDNSRecordPrecondition(t *testing.T) {
	for _, version := range []string{"", "123", "version+with&special=characters"} {
		t.Run(version, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Method != http.MethodDelete || r.URL.Path != "/project/proj/zones/example.com/dns-endpoints/rec" {
					t.Errorf("unexpected request: %s %s", r.Method, r.URL)
				}
				if got := r.URL.Query().Get("resourceVersion"); got != version {
					t.Errorf("resourceVersion = %q, want %q", got, version)
				}
				if version == "" && r.URL.RawQuery != "" {
					t.Errorf("unexpected query: %s", r.URL.RawQuery)
				}
				w.WriteHeader(http.StatusNoContent)
			}))
			defer server.Close()
			client, err := NewClient(server.URL)
			if err != nil {
				t.Fatal(err)
			}
			if err := client.DNSRecords.DeleteWithOptions(context.Background(), "proj", "example.com", "rec",
				DeleteDNSRecordOptions{ResourceVersion: version}); err != nil {
				t.Fatal(err)
			}
		})
	}
}

func TestRedirectIsAPIError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Location", "/redirect")
		w.WriteHeader(http.StatusFound)
	}))
	defer server.Close()
	client, err := NewClient(server.URL, WithHTTPClient(&http.Client{
		CheckRedirect: func(*http.Request, []*http.Request) error { return http.ErrUseLastResponse },
	}))
	if err != nil {
		t.Fatal(err)
	}
	if err := client.DNSRecords.Delete(context.Background(), "proj", "zone", "record"); !hasStatus(err, http.StatusFound) {
		t.Fatalf("expected APIError for redirect, got %v", err)
	}
}
