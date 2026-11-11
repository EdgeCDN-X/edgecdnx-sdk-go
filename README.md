# edgecdnx-sdk-go

Go client library (SDK) for the [EdgeCDN-X API](https://github.com/EdgeCDN-X/edgecdnx-api).
It is intended to be reused by any component that needs to talk to the EdgeCDN-X
control plane (CLIs, operators, integrations, UIs, automation).

## Installation

```sh
go get github.com/EdgeCDN-X/edgecdnx-sdk-go
```

## Usage

```go
package main

import (
	"context"
	"fmt"
	"log"

	edgecdnx "github.com/EdgeCDN-X/edgecdnx-sdk-go"
)

func main() {
	ctx := context.Background()

	// Authenticate with a project API key (Authorization: ApiKey ak_<id>_<secret>) ...
	client, err := edgecdnx.NewClient("https://api.edgecdnx.com", edgecdnx.WithAPIKey("ak_..."))
	if err != nil {
		log.Fatal(err)
	}

	// ... or with an OIDC bearer token:
	// client, _ := edgecdnx.NewClient(baseURL, edgecdnx.WithBearerToken(token))
	// ... or with a per-request token source (for short-lived tokens):
	// client, _ := edgecdnx.NewClient(baseURL, edgecdnx.WithTokenSource(fetchToken))

	projects, err := client.Projects.List(ctx)
	if err != nil {
		log.Fatal(err)
	}
	for _, p := range projects {
		fmt.Println(p.Name, p.Spec.Description)
	}
}
```

Responses reuse the EdgeCDN-X custom resource types from
[`edgecdnx-controller/api/v1alpha1`](https://github.com/EdgeCDN-X/edgecdnx-controller),
so fields match exactly what the API returns.

## Resources

| Service                      | Operations                      | Endpoint                                                         |
| ---------------------------- | ------------------------------- | ---------------------------------------------------------------- |
| `client.Projects`            | List, Get, Create               | `/projects`                                                      |
| `client.Zones`               | List, Create, Delete            | `/project/{project}/zones`                                       |
| `client.DNSRecords`          | List, Get, Create, Update, Delete | `/project/{project}/zones/{zone}/dns-endpoints`                |
| `client.APIKeys`             | List, Get, Create, Update, Delete | `/project/{project}/apikeys`                                   |
| `client.Locations`           | List, Get, Create, Update, Delete | `/project/{project}/locations`                                 |
| `client.HealthCheckProfiles` | List, Get, Create, Update, Delete | `/project/{project}/healthcheckprofiles`                       |

Notes:

- **Projects** — creating a project enrolls it in the free plan automatically on
  the server side; billing (checkout, portal) is out of scope for this SDK. The
  API does not currently expose project update or delete endpoints.
- **Zones** — zone names are unique across the platform; the API exposes list,
  create and delete only.
- **API keys** — the secret token is returned only once by `Create`
  (`CreatedAPIKey.Token`). `UpdateAPIKeyRequest` distinguishes "leave expiry
  unchanged" (default), `SetExpiry(t)` and `ClearExpiry()` (never expires).
  API keys are confined to their project and cannot create projects.

### Examples

```go
// Create a project (free plan).
project, err := client.Projects.Create(ctx, edgecdnx.CreateProjectRequest{
	Name:        "My Project",
	Description: "demo project",
})

// Create a zone.
zone, err := client.Zones.Create(ctx, project.Name, edgecdnx.CreateZoneRequest{
	Email: "admin@example.com",
	Zone:  "example.com",
})

// Create a DNS record (DNS endpoint).
record, err := client.DNSRecords.Create(ctx, project.Name, zone.Name, edgecdnx.CreateDNSRecordRequest{
	DNSName:    "www.example.com",
	RecordType: "A",
	RecordTTL:  edgecdnx.Ptr(300),
	Targets:    []string{"192.0.2.1"},
})

// Create an API key (store key.Token securely — it is shown only once).
key, err := client.APIKeys.Create(ctx, project.Name, edgecdnx.CreateAPIKeyRequest{
	DisplayName: "ci",
	Roles:       []string{"admin"},
})

// Remove an API key expiry.
var update edgecdnx.UpdateAPIKeyRequest
update.ClearExpiry()
_, err = client.APIKeys.Update(ctx, project.Name, key.ID, update)
```

### Error handling

Non-2xx responses are returned as `*edgecdnx.APIError` with helpers for common
status codes:

```go
_, err := client.Projects.Get(ctx, "missing")
if edgecdnx.IsNotFound(err) {
	// handle 404
}

var apiErr *edgecdnx.APIError
if errors.As(err, &apiErr) {
	fmt.Println(apiErr.StatusCode, apiErr.Message, apiErr.Reason)
}
```

Also available: `IsBadRequest`, `IsUnauthorized`, `IsForbidden`, `IsConflict`.

## Development

```sh
go test ./...
go vet ./...
```
