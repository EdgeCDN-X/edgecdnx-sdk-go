package edgecdnx

import (
	"context"
	"encoding/json"
	"net/http"
	"time"
)

// APIKeysService manages API keys within a project.
type APIKeysService struct {
	client *Client
}

// APIKey describes an API key as returned by the API. The secret token is
// only available in CreatedAPIKey at creation time.
type APIKey struct {
	// ID is the key identifier used in the token (ak_<id>_<secret>).
	ID string `json:"id"`
	// Name is the backing resource name.
	Name string `json:"name"`
	// DisplayName is the human-readable key name.
	DisplayName string `json:"displayName"`
	// Project is the project the key is confined to.
	Project string `json:"project"`
	// Roles lists the project roles assigned to the key.
	Roles []string `json:"roles"`
	// ExpiresAt is nil when the key never expires.
	ExpiresAt *time.Time `json:"expiresAt"`
	// CreatedAt is the creation timestamp.
	CreatedAt time.Time `json:"createdAt"`
}

// CreatedAPIKey is returned by APIKeysService.Create and includes the secret
// token, which is shown only once.
type CreatedAPIKey struct {
	APIKey
	// Token is the full API key token (ak_<id>_<secret>). Store it securely;
	// it cannot be retrieved again.
	Token string `json:"token"`
}

// CreateAPIKeyRequest is the payload for APIKeysService.Create.
type CreateAPIKeyRequest struct {
	// DisplayName is the human-readable key name (max 128 characters).
	DisplayName string `json:"displayName"`
	// ExpiresAt is optional; nil creates a key that never expires.
	ExpiresAt *time.Time `json:"expiresAt,omitempty"`
	// Roles lists project roles to assign to the key.
	Roles []string `json:"roles,omitempty"`
}

// UpdateAPIKeyRequest is the payload for APIKeysService.Update. Unset fields
// remain unchanged. Use SetExpiry/ClearExpiry to control the expiresAt field
// and note that setting Roles to an empty (non-nil) slice removes all roles.
type UpdateAPIKeyRequest struct {
	// DisplayName renames the key when non-nil.
	DisplayName *string
	// Roles replaces the role assignments when non-nil; an empty slice
	// removes all roles.
	Roles *[]string

	expiresAtSet bool
	expiresAt    *time.Time
}

// SetExpiry sets a new expiry time for the key.
func (r *UpdateAPIKeyRequest) SetExpiry(t time.Time) {
	r.expiresAtSet = true
	r.expiresAt = &t
}

// ClearExpiry removes the expiry, making the key never expire.
func (r *UpdateAPIKeyRequest) ClearExpiry() {
	r.expiresAtSet = true
	r.expiresAt = nil
}

// MarshalJSON implements json.Marshaler. It distinguishes omitted fields from
// explicit nulls, which the API treats differently for expiresAt.
func (r UpdateAPIKeyRequest) MarshalJSON() ([]byte, error) {
	payload := map[string]any{}
	if r.DisplayName != nil {
		payload["displayName"] = *r.DisplayName
	}
	if r.expiresAtSet {
		payload["expiresAt"] = r.expiresAt
	}
	if r.Roles != nil {
		payload["roles"] = *r.Roles
	}
	return json.Marshal(payload)
}

func apiKeysPath(projectID string) string {
	return "/project/" + pathEscape(projectID) + "/apikeys"
}

// List returns all API keys in the project.
func (s *APIKeysService) List(ctx context.Context, projectID string) ([]APIKey, error) {
	var keys []APIKey
	if err := s.client.do(ctx, http.MethodGet, apiKeysPath(projectID), nil, &keys); err != nil {
		return nil, err
	}
	return keys, nil
}

// Get returns a single API key by its identifier.
func (s *APIKeysService) Get(ctx context.Context, projectID, apiKeyID string) (*APIKey, error) {
	key := &APIKey{}
	if err := s.client.do(ctx, http.MethodGet, apiKeysPath(projectID)+"/"+pathEscape(apiKeyID), nil, key); err != nil {
		return nil, err
	}
	return key, nil
}

// Create creates a new API key. The returned token is shown only once.
func (s *APIKeysService) Create(ctx context.Context, projectID string, req CreateAPIKeyRequest) (*CreatedAPIKey, error) {
	key := &CreatedAPIKey{}
	if err := s.client.do(ctx, http.MethodPost, apiKeysPath(projectID), req, key); err != nil {
		return nil, err
	}
	return key, nil
}

// Update patches an existing API key.
func (s *APIKeysService) Update(ctx context.Context, projectID, apiKeyID string, req UpdateAPIKeyRequest) (*APIKey, error) {
	key := &APIKey{}
	if err := s.client.do(ctx, http.MethodPatch, apiKeysPath(projectID)+"/"+pathEscape(apiKeyID), req, key); err != nil {
		return nil, err
	}
	return key, nil
}

// Delete deletes an API key by its identifier, revoking it immediately.
func (s *APIKeysService) Delete(ctx context.Context, projectID, apiKeyID string) error {
	return s.client.do(ctx, http.MethodDelete, apiKeysPath(projectID)+"/"+pathEscape(apiKeyID), nil, nil)
}
