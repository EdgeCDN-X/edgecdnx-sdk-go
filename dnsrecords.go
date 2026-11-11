package edgecdnx

import (
	"context"
	"net/http"

	v1alpha1 "github.com/EdgeCDN-X/edgecdnx-controller/api/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// DNSRecordsService manages DNS records (DNS endpoints) within a zone.
type DNSRecordsService struct {
	client *Client
}

// CreateDNSRecordRequest is the payload for DNSRecordsService.Create.
type CreateDNSRecordRequest struct {
	// DNSName is the record name, relative to or including the zone.
	DNSName string `json:"dnsName"`
	// RoutingPolicy is one of Simple, Weighted, Failover, Geolocation or
	// RoundRobin. Defaults to Simple when empty.
	RoutingPolicy string `json:"routingPolicy,omitempty"`
	// RecordTTL is the record TTL in seconds. Required; use edgecdnx.Ptr(0)
	// for an explicit zero TTL.
	RecordTTL *int `json:"recordTTL"`
	// RecordType is one of A, AAAA, CNAME, TXT, MX, SRV or NS.
	RecordType string `json:"recordType"`
	// Targets holds the record targets. Required for Simple and Failover
	// routing policies.
	Targets []string `json:"targets,omitempty"`
	// RouteSelector selects routes for Weighted, Geolocation and RoundRobin
	// routing policies (matchLabels only).
	RouteSelector *metav1.LabelSelector `json:"routeSelector,omitempty"`
}

// UpdateDNSRecordRequest is the payload for DNSRecordsService.Update.
// Zero-valued fields are left unchanged.
type UpdateDNSRecordRequest struct {
	// ResourceVersion enables optimistic concurrency when set.
	ResourceVersion string                `json:"resourceVersion,omitempty"`
	DNSName         string                `json:"dnsName,omitempty"`
	RoutingPolicy   string                `json:"routingPolicy,omitempty"`
	RecordTTL       *int                  `json:"recordTTL,omitempty"`
	RecordType      string                `json:"recordType,omitempty"`
	Targets         []string              `json:"targets,omitempty"`
	RouteSelector   *metav1.LabelSelector `json:"routeSelector,omitempty"`
}

func dnsRecordsPath(projectID, zoneID string) string {
	return zonesPath(projectID) + "/" + pathEscape(zoneID) + "/dns-endpoints"
}

// List returns all DNS records in the zone, including read-only records
// derived from service host aliases.
func (s *DNSRecordsService) List(ctx context.Context, projectID, zoneID string) ([]v1alpha1.DNSEndpoint, error) {
	var records []v1alpha1.DNSEndpoint
	if err := s.client.do(ctx, http.MethodGet, dnsRecordsPath(projectID, zoneID), nil, &records); err != nil {
		return nil, err
	}
	return records, nil
}

// Get returns a single DNS record by its identifier.
func (s *DNSRecordsService) Get(ctx context.Context, projectID, zoneID, recordID string) (*v1alpha1.DNSEndpoint, error) {
	record := &v1alpha1.DNSEndpoint{}
	if err := s.client.do(ctx, http.MethodGet, dnsRecordsPath(projectID, zoneID)+"/"+pathEscape(recordID), nil, record); err != nil {
		return nil, err
	}
	return record, nil
}

// Create creates a new DNS record in the zone.
func (s *DNSRecordsService) Create(ctx context.Context, projectID, zoneID string, req CreateDNSRecordRequest) (*v1alpha1.DNSEndpoint, error) {
	record := &v1alpha1.DNSEndpoint{}
	if err := s.client.do(ctx, http.MethodPost, dnsRecordsPath(projectID, zoneID), req, record); err != nil {
		return nil, err
	}
	return record, nil
}

// Update patches an existing DNS record.
func (s *DNSRecordsService) Update(ctx context.Context, projectID, zoneID, recordID string, req UpdateDNSRecordRequest) (*v1alpha1.DNSEndpoint, error) {
	record := &v1alpha1.DNSEndpoint{}
	if err := s.client.do(ctx, http.MethodPatch, dnsRecordsPath(projectID, zoneID)+"/"+pathEscape(recordID), req, record); err != nil {
		return nil, err
	}
	return record, nil
}

// Delete deletes a DNS record by its identifier.
func (s *DNSRecordsService) Delete(ctx context.Context, projectID, zoneID, recordID string) error {
	return s.client.do(ctx, http.MethodDelete, dnsRecordsPath(projectID, zoneID)+"/"+pathEscape(recordID), nil, nil)
}
