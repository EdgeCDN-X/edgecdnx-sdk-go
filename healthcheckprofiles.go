package edgecdnx

import (
	"context"
	"net/http"

	v1alpha1 "github.com/EdgeCDN-X/edgecdnx-controller/api/v1alpha1"
)

// HealthCheckProfilesService manages health check profiles within a project.
type HealthCheckProfilesService struct {
	client *Client
}

// CreateHealthCheckProfileRequest is the payload for
// HealthCheckProfilesService.Create.
type CreateHealthCheckProfileRequest struct {
	// Name is the profile name (DNS-1123 subdomain).
	Name string `json:"name"`
	// Probes defines the health check probes (at least one is required).
	Probes []v1alpha1.HealthCheckProbeSpec `json:"probes"`
}

// UpdateHealthCheckProfileRequest is the payload for
// HealthCheckProfilesService.Update. A nil Probes field remains unchanged.
type UpdateHealthCheckProfileRequest struct {
	Probes *[]v1alpha1.HealthCheckProbeSpec `json:"probes,omitempty"`
}

func healthCheckProfilesPath(projectID string) string {
	return "/project/" + pathEscape(projectID) + "/healthcheckprofiles"
}

// List returns all health check profiles in the project.
func (s *HealthCheckProfilesService) List(ctx context.Context, projectID string) ([]v1alpha1.HealthCheckProfile, error) {
	var profiles []v1alpha1.HealthCheckProfile
	if err := s.client.do(ctx, http.MethodGet, healthCheckProfilesPath(projectID), nil, &profiles); err != nil {
		return nil, err
	}
	return profiles, nil
}

// Get returns a single health check profile by name.
func (s *HealthCheckProfilesService) Get(ctx context.Context, projectID, profileID string) (*v1alpha1.HealthCheckProfile, error) {
	profile := &v1alpha1.HealthCheckProfile{}
	if err := s.client.do(ctx, http.MethodGet, healthCheckProfilesPath(projectID)+"/"+pathEscape(profileID), nil, profile); err != nil {
		return nil, err
	}
	return profile, nil
}

// Create creates a new health check profile in the project.
func (s *HealthCheckProfilesService) Create(ctx context.Context, projectID string, req CreateHealthCheckProfileRequest) (*v1alpha1.HealthCheckProfile, error) {
	profile := &v1alpha1.HealthCheckProfile{}
	if err := s.client.do(ctx, http.MethodPost, healthCheckProfilesPath(projectID), req, profile); err != nil {
		return nil, err
	}
	return profile, nil
}

// Update patches an existing health check profile.
func (s *HealthCheckProfilesService) Update(ctx context.Context, projectID, profileID string, req UpdateHealthCheckProfileRequest) (*v1alpha1.HealthCheckProfile, error) {
	profile := &v1alpha1.HealthCheckProfile{}
	if err := s.client.do(ctx, http.MethodPatch, healthCheckProfilesPath(projectID)+"/"+pathEscape(profileID), req, profile); err != nil {
		return nil, err
	}
	return profile, nil
}

// Delete deletes a health check profile by name.
func (s *HealthCheckProfilesService) Delete(ctx context.Context, projectID, profileID string) error {
	return s.client.do(ctx, http.MethodDelete, healthCheckProfilesPath(projectID)+"/"+pathEscape(profileID), nil, nil)
}
