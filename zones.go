package edgecdnx

import (
	"context"
	"net/http"

	v1alpha1 "github.com/EdgeCDN-X/edgecdnx-controller/api/v1alpha1"
)

// ZonesService manages DNS zones within a project.
//
// Zone names are unique across the whole platform. The API currently does
// not expose zone get-by-id or update endpoints.
type ZonesService struct {
	client *Client
}

// CreateZoneRequest is the payload for ZonesService.Create.
type CreateZoneRequest struct {
	// Email is the zone administrator contact (SOA RNAME is derived from it).
	Email string `json:"email"`
	// Zone is the fully qualified zone name, for example "example.com".
	Zone string `json:"zone"`
}

func zonesPath(projectID string) string {
	return "/project/" + pathEscape(projectID) + "/zones"
}

// List returns all zones in the project.
func (s *ZonesService) List(ctx context.Context, projectID string) ([]v1alpha1.Zone, error) {
	var zones []v1alpha1.Zone
	if err := s.client.do(ctx, http.MethodGet, zonesPath(projectID), nil, &zones); err != nil {
		return nil, err
	}
	return zones, nil
}

// Create creates a new zone in the project.
func (s *ZonesService) Create(ctx context.Context, projectID string, req CreateZoneRequest) (*v1alpha1.Zone, error) {
	zone := &v1alpha1.Zone{}
	if err := s.client.do(ctx, http.MethodPost, zonesPath(projectID), req, zone); err != nil {
		return nil, err
	}
	return zone, nil
}

// Delete deletes a zone by name.
func (s *ZonesService) Delete(ctx context.Context, projectID, zoneID string) error {
	return s.client.do(ctx, http.MethodDelete, zonesPath(projectID)+"/"+pathEscape(zoneID), nil, nil)
}
