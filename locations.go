package edgecdnx

import (
	"context"
	"net/http"

	v1alpha1 "github.com/EdgeCDN-X/edgecdnx-controller/api/v1alpha1"
)

// LocationsService manages edge locations within a project.
type LocationsService struct {
	client *Client
}

// CreateLocationRequest is the payload for LocationsService.Create.
type CreateLocationRequest struct {
	// Name is the location name (DNS-1123 subdomain).
	Name string `json:"name"`
	// BaseDomain is the base domain used for resources in this location.
	BaseDomain string `json:"baseDomain,omitempty"`
	// Labels are additional labels applied to the location.
	Labels map[string]string `json:"labels,omitempty"`
	// NodeGroups define the node groups serving this location.
	NodeGroups []v1alpha1.NodeGroupSpec `json:"nodeGroups,omitempty"`
	// GeoLookup configures geo-routing for this location.
	GeoLookup v1alpha1.GeoLookupSpec `json:"geoLookup,omitempty"`
	// Weight is the relative routing weight.
	Weight int32 `json:"weight,omitempty"`
	// FallbackLocations lists locations to fail over to.
	FallbackLocations []string `json:"fallbackLocations,omitempty"`
}

// UpdateLocationRequest is the payload for LocationsService.Update.
// Nil fields remain unchanged.
type UpdateLocationRequest struct {
	BaseDomain        *string                   `json:"baseDomain,omitempty"`
	Labels            *map[string]string        `json:"labels,omitempty"`
	NodeGroups        *[]v1alpha1.NodeGroupSpec `json:"nodeGroups,omitempty"`
	GeoLookup         *v1alpha1.GeoLookupSpec   `json:"geoLookup,omitempty"`
	Weight            *int32                    `json:"weight,omitempty"`
	FallbackLocations *[]string                 `json:"fallbackLocations,omitempty"`
}

func locationsPath(projectID string) string {
	return "/project/" + pathEscape(projectID) + "/locations"
}

// List returns all locations in the project.
func (s *LocationsService) List(ctx context.Context, projectID string) ([]v1alpha1.Location, error) {
	var locations []v1alpha1.Location
	if err := s.client.do(ctx, http.MethodGet, locationsPath(projectID), nil, &locations); err != nil {
		return nil, err
	}
	return locations, nil
}

// Get returns a single location by name.
func (s *LocationsService) Get(ctx context.Context, projectID, locationID string) (*v1alpha1.Location, error) {
	location := &v1alpha1.Location{}
	if err := s.client.do(ctx, http.MethodGet, locationsPath(projectID)+"/"+pathEscape(locationID), nil, location); err != nil {
		return nil, err
	}
	return location, nil
}

// Create creates a new location in the project.
func (s *LocationsService) Create(ctx context.Context, projectID string, req CreateLocationRequest) (*v1alpha1.Location, error) {
	location := &v1alpha1.Location{}
	if err := s.client.do(ctx, http.MethodPost, locationsPath(projectID), req, location); err != nil {
		return nil, err
	}
	return location, nil
}

// Update patches an existing location.
func (s *LocationsService) Update(ctx context.Context, projectID, locationID string, req UpdateLocationRequest) (*v1alpha1.Location, error) {
	location := &v1alpha1.Location{}
	if err := s.client.do(ctx, http.MethodPatch, locationsPath(projectID)+"/"+pathEscape(locationID), req, location); err != nil {
		return nil, err
	}
	return location, nil
}

// Delete deletes a location by name.
func (s *LocationsService) Delete(ctx context.Context, projectID, locationID string) error {
	return s.client.do(ctx, http.MethodDelete, locationsPath(projectID)+"/"+pathEscape(locationID), nil, nil)
}
