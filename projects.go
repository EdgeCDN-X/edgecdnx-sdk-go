package edgecdnx

import (
	"context"
	"net/http"

	v1alpha1 "github.com/EdgeCDN-X/edgecdnx-controller/api/v1alpha1"
)

// ProjectsService manages EdgeCDN-X projects.
//
// Creating a project automatically enrolls it in the free plan on the
// server side; paid plans are handled separately through the billing
// endpoints, which are not part of this SDK. The API currently does not
// expose project update or delete endpoints.
type ProjectsService struct {
	client *Client
}

// CreateProjectRequest is the payload for ProjectsService.Create.
type CreateProjectRequest struct {
	// Name is the human-readable project name (3-63 characters). The API
	// slugifies it into the project identifier.
	Name string `json:"name"`
	// Description is an optional description (max 255 characters).
	Description string `json:"description,omitempty"`
}

// List returns the projects the authenticated principal can read.
func (s *ProjectsService) List(ctx context.Context) ([]v1alpha1.Project, error) {
	var projects []v1alpha1.Project
	if err := s.client.do(ctx, http.MethodGet, "/projects", nil, &projects); err != nil {
		return nil, err
	}
	return projects, nil
}

// Get returns a single project by its identifier (slug).
func (s *ProjectsService) Get(ctx context.Context, projectID string) (*v1alpha1.Project, error) {
	project := &v1alpha1.Project{}
	if err := s.client.do(ctx, http.MethodGet, "/projects/"+pathEscape(projectID), nil, project); err != nil {
		return nil, err
	}
	return project, nil
}

// Create creates a new project on the free plan. The caller becomes the
// project admin. API key principals cannot create projects.
func (s *ProjectsService) Create(ctx context.Context, req CreateProjectRequest) (*v1alpha1.Project, error) {
	project := &v1alpha1.Project{}
	if err := s.client.do(ctx, http.MethodPost, "/projects", req, project); err != nil {
		return nil, err
	}
	return project, nil
}
