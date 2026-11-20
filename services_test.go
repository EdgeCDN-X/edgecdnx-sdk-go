package edgecdnx

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	v1alpha1 "github.com/EdgeCDN-X/edgecdnx-controller/api/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestProjectsService(t *testing.T) {
	ctx := context.Background()

	t.Run("List", func(t *testing.T) {
		var rec recordedRequest
		client := newTestClient(t, 200, []v1alpha1.Project{{ObjectMeta: metav1.ObjectMeta{Name: "proj"}}}, &rec)
		projects, err := client.Projects.List(ctx)
		if err != nil {
			t.Fatalf("List: %v", err)
		}
		if rec.Method != "GET" || rec.Path != "/projects" {
			t.Fatalf("unexpected request %s %s", rec.Method, rec.Path)
		}
		if len(projects) != 1 || projects[0].Name != "proj" {
			t.Fatalf("unexpected projects %+v", projects)
		}
	})

	t.Run("Get", func(t *testing.T) {
		var rec recordedRequest
		client := newTestClient(t, 200, v1alpha1.Project{ObjectMeta: metav1.ObjectMeta{Name: "proj"}}, &rec)
		project, err := client.Projects.Get(ctx, "proj")
		if err != nil {
			t.Fatalf("Get: %v", err)
		}
		if rec.Method != "GET" || rec.Path != "/projects/proj" {
			t.Fatalf("unexpected request %s %s", rec.Method, rec.Path)
		}
		if project.Name != "proj" {
			t.Fatalf("unexpected project %+v", project)
		}
	})

	t.Run("Create", func(t *testing.T) {
		var rec recordedRequest
		client := newTestClient(t, 201, v1alpha1.Project{ObjectMeta: metav1.ObjectMeta{Name: "my-project"}}, &rec)
		project, err := client.Projects.Create(ctx, CreateProjectRequest{Name: "My Project", Description: "demo"})
		if err != nil {
			t.Fatalf("Create: %v", err)
		}
		if rec.Method != "POST" || rec.Path != "/projects" {
			t.Fatalf("unexpected request %s %s", rec.Method, rec.Path)
		}
		var sent map[string]any
		if err := json.Unmarshal(rec.Body, &sent); err != nil {
			t.Fatalf("unmarshal sent body: %v", err)
		}
		if sent["name"] != "My Project" || sent["description"] != "demo" {
			t.Fatalf("unexpected body %s", rec.Body)
		}
		if project.Name != "my-project" {
			t.Fatalf("unexpected project %+v", project)
		}
	})
}

func TestZonesService(t *testing.T) {
	ctx := context.Background()

	t.Run("List", func(t *testing.T) {
		var rec recordedRequest
		client := newTestClient(t, 200, []v1alpha1.Zone{{ObjectMeta: metav1.ObjectMeta{Name: "example.com"}}}, &rec)
		zones, err := client.Zones.List(ctx, "proj")
		if err != nil {
			t.Fatalf("List: %v", err)
		}
		if rec.Method != "GET" || rec.Path != "/project/proj/zones" {
			t.Fatalf("unexpected request %s %s", rec.Method, rec.Path)
		}
		if len(zones) != 1 {
			t.Fatalf("unexpected zones %+v", zones)
		}
	})

	t.Run("Create", func(t *testing.T) {
		var rec recordedRequest
		client := newTestClient(t, 201, v1alpha1.Zone{ObjectMeta: metav1.ObjectMeta{Name: "example.com"}}, &rec)
		zone, err := client.Zones.Create(ctx, "proj", CreateZoneRequest{Email: "admin@example.com", Zone: "example.com"})
		if err != nil {
			t.Fatalf("Create: %v", err)
		}
		if rec.Method != "POST" || rec.Path != "/project/proj/zones" {
			t.Fatalf("unexpected request %s %s", rec.Method, rec.Path)
		}
		var sent map[string]any
		if err := json.Unmarshal(rec.Body, &sent); err != nil {
			t.Fatalf("unmarshal sent body: %v", err)
		}
		if sent["email"] != "admin@example.com" || sent["zone"] != "example.com" {
			t.Fatalf("unexpected body %s", rec.Body)
		}
		if zone.Name != "example.com" {
			t.Fatalf("unexpected zone %+v", zone)
		}
	})

	t.Run("Delete", func(t *testing.T) {
		var rec recordedRequest
		client := newTestClient(t, 204, nil, &rec)
		if err := client.Zones.Delete(ctx, "proj", "example.com"); err != nil {
			t.Fatalf("Delete: %v", err)
		}
		if rec.Method != "DELETE" || rec.Path != "/project/proj/zones/example.com" {
			t.Fatalf("unexpected request %s %s", rec.Method, rec.Path)
		}
	})
}

func TestDNSRecordsService(t *testing.T) {
	ctx := context.Background()

	t.Run("List", func(t *testing.T) {
		var rec recordedRequest
		client := newTestClient(t, 200, []v1alpha1.DNSEndpoint{{ObjectMeta: metav1.ObjectMeta{Name: "rec1"}}}, &rec)
		records, err := client.DNSRecords.List(ctx, "proj", "example.com")
		if err != nil {
			t.Fatalf("List: %v", err)
		}
		if rec.Method != "GET" || rec.Path != "/project/proj/zones/example.com/dns-endpoints" {
			t.Fatalf("unexpected request %s %s", rec.Method, rec.Path)
		}
		if len(records) != 1 {
			t.Fatalf("unexpected records %+v", records)
		}
	})

	t.Run("Get", func(t *testing.T) {
		var rec recordedRequest
		client := newTestClient(t, 200, v1alpha1.DNSEndpoint{ObjectMeta: metav1.ObjectMeta{Name: "rec1"}}, &rec)
		record, err := client.DNSRecords.Get(ctx, "proj", "example.com", "rec1")
		if err != nil {
			t.Fatalf("Get: %v", err)
		}
		if rec.Method != "GET" || rec.Path != "/project/proj/zones/example.com/dns-endpoints/rec1" {
			t.Fatalf("unexpected request %s %s", rec.Method, rec.Path)
		}
		if record.Name != "rec1" {
			t.Fatalf("unexpected record %+v", record)
		}
	})

	t.Run("Create", func(t *testing.T) {
		var rec recordedRequest
		client := newTestClient(t, 201, v1alpha1.DNSEndpoint{ObjectMeta: metav1.ObjectMeta{Name: "rec1"}}, &rec)
		_, err := client.DNSRecords.Create(ctx, "proj", "example.com", CreateDNSRecordRequest{
			DNSName:    "www.example.com",
			RecordTTL:  Ptr(300),
			RecordType: "A",
			Targets:    []string{"192.0.2.1"},
		})
		if err != nil {
			t.Fatalf("Create: %v", err)
		}
		if rec.Method != "POST" || rec.Path != "/project/proj/zones/example.com/dns-endpoints" {
			t.Fatalf("unexpected request %s %s", rec.Method, rec.Path)
		}
		var sent map[string]any
		if err := json.Unmarshal(rec.Body, &sent); err != nil {
			t.Fatalf("unmarshal sent body: %v", err)
		}
		if sent["dnsName"] != "www.example.com" || sent["recordTTL"] != float64(300) || sent["recordType"] != "A" {
			t.Fatalf("unexpected body %s", rec.Body)
		}
	})

	t.Run("CreateWithRouteSelector", func(t *testing.T) {
		var rec recordedRequest
		client := newTestClient(t, 201, v1alpha1.DNSEndpoint{}, &rec)
		_, err := client.DNSRecords.Create(ctx, "proj", "example.com", CreateDNSRecordRequest{
			DNSName:       "geo.example.com",
			RoutingPolicy: "Geolocation",
			RecordTTL:     Ptr(0),
			RecordType:    "A",
			RouteSelector: &metav1.LabelSelector{MatchLabels: map[string]string{"project": "proj"}},
		})
		if err != nil {
			t.Fatalf("Create: %v", err)
		}
		var sent map[string]any
		if err := json.Unmarshal(rec.Body, &sent); err != nil {
			t.Fatalf("unmarshal sent body: %v", err)
		}
		// Explicit zero TTL must be serialized, not omitted.
		if ttl, ok := sent["recordTTL"]; !ok || ttl != float64(0) {
			t.Fatalf("expected explicit recordTTL 0 in body %s", rec.Body)
		}
		selector, ok := sent["routeSelector"].(map[string]any)
		if !ok {
			t.Fatalf("expected routeSelector in body %s", rec.Body)
		}
		labels := selector["matchLabels"].(map[string]any)
		if labels["project"] != "proj" {
			t.Fatalf("unexpected routeSelector %v", selector)
		}
	})

	t.Run("Update", func(t *testing.T) {
		var rec recordedRequest
		client := newTestClient(t, 200, v1alpha1.DNSEndpoint{ObjectMeta: metav1.ObjectMeta{Name: "rec1"}}, &rec)
		_, err := client.DNSRecords.Update(ctx, "proj", "example.com", "rec1", UpdateDNSRecordRequest{
			Targets: []string{"192.0.2.2"},
		})
		if err != nil {
			t.Fatalf("Update: %v", err)
		}
		if rec.Method != "PATCH" || rec.Path != "/project/proj/zones/example.com/dns-endpoints/rec1" {
			t.Fatalf("unexpected request %s %s", rec.Method, rec.Path)
		}
		var sent map[string]any
		if err := json.Unmarshal(rec.Body, &sent); err != nil {
			t.Fatalf("unmarshal sent body: %v", err)
		}
		if _, ok := sent["recordTTL"]; ok {
			t.Fatalf("nil recordTTL must be omitted on update, body %s", rec.Body)
		}
	})

	t.Run("Delete", func(t *testing.T) {
		var rec recordedRequest
		client := newTestClient(t, 204, nil, &rec)
		if err := client.DNSRecords.Delete(ctx, "proj", "example.com", "rec1"); err != nil {
			t.Fatalf("Delete: %v", err)
		}
		if rec.Method != "DELETE" || rec.Path != "/project/proj/zones/example.com/dns-endpoints/rec1" {
			t.Fatalf("unexpected request %s %s", rec.Method, rec.Path)
		}
	})
}

func TestAPIKeysService(t *testing.T) {
	ctx := context.Background()

	t.Run("List", func(t *testing.T) {
		var rec recordedRequest
		client := newTestClient(t, 200, []APIKey{{ID: "abc", DisplayName: "ci"}}, &rec)
		keys, err := client.APIKeys.List(ctx, "proj")
		if err != nil {
			t.Fatalf("List: %v", err)
		}
		if rec.Method != "GET" || rec.Path != "/project/proj/apikeys" {
			t.Fatalf("unexpected request %s %s", rec.Method, rec.Path)
		}
		if len(keys) != 1 || keys[0].ID != "abc" {
			t.Fatalf("unexpected keys %+v", keys)
		}
	})

	t.Run("Get", func(t *testing.T) {
		var rec recordedRequest
		client := newTestClient(t, 200, APIKey{ID: "abc"}, &rec)
		key, err := client.APIKeys.Get(ctx, "proj", "abc")
		if err != nil {
			t.Fatalf("Get: %v", err)
		}
		if rec.Method != "GET" || rec.Path != "/project/proj/apikeys/abc" {
			t.Fatalf("unexpected request %s %s", rec.Method, rec.Path)
		}
		if key.ID != "abc" {
			t.Fatalf("unexpected key %+v", key)
		}
	})

	t.Run("Create", func(t *testing.T) {
		var rec recordedRequest
		client := newTestClient(t, 201, CreatedAPIKey{APIKey: APIKey{ID: "abc"}, Token: "ak_abc_secret"}, &rec)
		expiry := time.Date(2030, 1, 1, 0, 0, 0, 0, time.UTC)
		key, err := client.APIKeys.Create(ctx, "proj", CreateAPIKeyRequest{
			DisplayName: "ci",
			ExpiresAt:   &expiry,
			Roles:       []string{"admin"},
		})
		if err != nil {
			t.Fatalf("Create: %v", err)
		}
		if rec.Method != "POST" || rec.Path != "/project/proj/apikeys" {
			t.Fatalf("unexpected request %s %s", rec.Method, rec.Path)
		}
		if key.Token != "ak_abc_secret" {
			t.Fatalf("unexpected token %q", key.Token)
		}
	})

	t.Run("Update", func(t *testing.T) {
		var rec recordedRequest
		client := newTestClient(t, 200, APIKey{ID: "abc", DisplayName: "renamed"}, &rec)
		var req UpdateAPIKeyRequest
		req.DisplayName = Ptr("renamed")
		key, err := client.APIKeys.Update(ctx, "proj", "abc", req)
		if err != nil {
			t.Fatalf("Update: %v", err)
		}
		if rec.Method != "PATCH" || rec.Path != "/project/proj/apikeys/abc" {
			t.Fatalf("unexpected request %s %s", rec.Method, rec.Path)
		}
		var sent map[string]any
		if err := json.Unmarshal(rec.Body, &sent); err != nil {
			t.Fatalf("unmarshal sent body: %v", err)
		}
		if _, ok := sent["expiresAt"]; ok {
			t.Fatalf("expiresAt must be omitted when unset, body %s", rec.Body)
		}
		if key.DisplayName != "renamed" {
			t.Fatalf("unexpected key %+v", key)
		}
	})

	t.Run("Delete", func(t *testing.T) {
		var rec recordedRequest
		client := newTestClient(t, 204, nil, &rec)
		if err := client.APIKeys.Delete(ctx, "proj", "abc"); err != nil {
			t.Fatalf("Delete: %v", err)
		}
		if rec.Method != "DELETE" || rec.Path != "/project/proj/apikeys/abc" {
			t.Fatalf("unexpected request %s %s", rec.Method, rec.Path)
		}
	})
}

func TestLocationsService(t *testing.T) {
	ctx := context.Background()

	t.Run("List", func(t *testing.T) {
		var rec recordedRequest
		client := newTestClient(t, 200, []v1alpha1.Location{{ObjectMeta: metav1.ObjectMeta{Name: "eu-west"}}}, &rec)
		locations, err := client.Locations.List(ctx, "proj")
		if err != nil {
			t.Fatalf("List: %v", err)
		}
		if rec.Method != "GET" || rec.Path != "/project/proj/locations" {
			t.Fatalf("unexpected request %s %s", rec.Method, rec.Path)
		}
		if len(locations) != 1 {
			t.Fatalf("unexpected locations %+v", locations)
		}
	})

	t.Run("Get", func(t *testing.T) {
		var rec recordedRequest
		client := newTestClient(t, 200, v1alpha1.Location{ObjectMeta: metav1.ObjectMeta{Name: "eu-west"}}, &rec)
		location, err := client.Locations.Get(ctx, "proj", "eu-west")
		if err != nil {
			t.Fatalf("Get: %v", err)
		}
		if rec.Method != "GET" || rec.Path != "/project/proj/locations/eu-west" {
			t.Fatalf("unexpected request %s %s", rec.Method, rec.Path)
		}
		if location.Name != "eu-west" {
			t.Fatalf("unexpected location %+v", location)
		}
	})

	t.Run("Create", func(t *testing.T) {
		var rec recordedRequest
		client := newTestClient(t, 201, v1alpha1.Location{
			ObjectMeta: metav1.ObjectMeta{Name: "eu-west"},
			Spec: v1alpha1.LocationSpec{
				BaseDomain: "edge.example.com",
				NodeGroups: []v1alpha1.NodeGroupSpec{{
					Name:         "cache",
					NodeSelector: map[string]string{"pool": "edge"},
					Nodes:        []v1alpha1.NodeSpec{{Name: "node-1", Cname: "node-1.example.com"}},
				}},
			},
		}, &rec)
		location, err := client.Locations.Create(ctx, "proj", CreateLocationRequest{
			Name:       "eu-west",
			BaseDomain: "edge.example.com",
			Weight:     10,
			NodeGroups: []v1alpha1.NodeGroupSpec{{
				Name:         "cache",
				NodeSelector: map[string]string{"pool": "edge"},
				Nodes:        []v1alpha1.NodeSpec{{Name: "node-1", Cname: "node-1.example.com"}},
			}},
		})
		if err != nil {
			t.Fatalf("Create: %v", err)
		}
		if rec.Method != "POST" || rec.Path != "/project/proj/locations" {
			t.Fatalf("unexpected request %s %s", rec.Method, rec.Path)
		}
		var sent map[string]any
		if err := json.Unmarshal(rec.Body, &sent); err != nil {
			t.Fatalf("unmarshal sent body: %v", err)
		}
		if sent["baseDomain"] != "edge.example.com" {
			t.Fatalf("baseDomain missing from request: %s", rec.Body)
		}
		nodeGroups, ok := sent["nodeGroups"].([]any)
		if !ok || len(nodeGroups) != 1 {
			t.Fatalf("nodeGroups missing from request: %s", rec.Body)
		}
		nodeGroup, ok := nodeGroups[0].(map[string]any)
		if !ok {
			t.Fatalf("invalid nodeGroup in request: %s", rec.Body)
		}
		nodeSelector, ok := nodeGroup["nodeSelector"].(map[string]any)
		if !ok || nodeSelector["pool"] != "edge" {
			t.Fatalf("nodeSelector missing from request: %s", rec.Body)
		}
		nodes, ok := nodeGroup["nodes"].([]any)
		if !ok || len(nodes) != 1 {
			t.Fatalf("nodes missing from request: %s", rec.Body)
		}
		node, ok := nodes[0].(map[string]any)
		if !ok || node["cname"] != "node-1.example.com" {
			t.Fatalf("node cname missing from request: %s", rec.Body)
		}
		if location.Spec.BaseDomain != "edge.example.com" || len(location.Spec.NodeGroups) != 1 {
			t.Fatalf("new location DTO fields not decoded: %+v", location.Spec)
		}
		group := location.Spec.NodeGroups[0]
		if group.NodeSelector["pool"] != "edge" || len(group.Nodes) != 1 ||
			group.Nodes[0].Cname != "node-1.example.com" {
			t.Fatalf("new location DTO fields not decoded: %+v", location.Spec)
		}
		legacyPayload, err := json.Marshal(CreateLocationRequest{Name: "eu-west"})
		if err != nil {
			t.Fatalf("marshal legacy request: %v", err)
		}
		var legacyBody map[string]any
		if err := json.Unmarshal(legacyPayload, &legacyBody); err != nil {
			t.Fatalf("unmarshal legacy request: %v", err)
		}
		if _, ok := legacyBody["baseDomain"]; ok {
			t.Fatalf("unset baseDomain must be omitted, body %s", legacyPayload)
		}
	})

	t.Run("Update", func(t *testing.T) {
		var rec recordedRequest
		client := newTestClient(t, 200, v1alpha1.Location{ObjectMeta: metav1.ObjectMeta{Name: "eu-west"}}, &rec)
		_, err := client.Locations.Update(ctx, "proj", "eu-west", UpdateLocationRequest{
			BaseDomain: Ptr(""),
			Weight:     Ptr(int32(20)),
		})
		if err != nil {
			t.Fatalf("Update: %v", err)
		}
		if rec.Method != "PATCH" || rec.Path != "/project/proj/locations/eu-west" {
			t.Fatalf("unexpected request %s %s", rec.Method, rec.Path)
		}
		var sent map[string]any
		if err := json.Unmarshal(rec.Body, &sent); err != nil {
			t.Fatalf("unmarshal sent body: %v", err)
		}
		if sent["weight"] != float64(20) {
			t.Fatalf("unexpected body %s", rec.Body)
		}
		if baseDomain, ok := sent["baseDomain"]; !ok || baseDomain != "" {
			t.Fatalf("explicit empty baseDomain must be sent, body %s", rec.Body)
		}
		if _, ok := sent["labels"]; ok {
			t.Fatalf("nil labels must be omitted, body %s", rec.Body)
		}
	})

	t.Run("Delete", func(t *testing.T) {
		var rec recordedRequest
		client := newTestClient(t, 204, nil, &rec)
		if err := client.Locations.Delete(ctx, "proj", "eu-west"); err != nil {
			t.Fatalf("Delete: %v", err)
		}
		if rec.Method != "DELETE" || rec.Path != "/project/proj/locations/eu-west" {
			t.Fatalf("unexpected request %s %s", rec.Method, rec.Path)
		}
	})
}

func TestHealthCheckProfilesService(t *testing.T) {
	ctx := context.Background()

	t.Run("List", func(t *testing.T) {
		var rec recordedRequest
		client := newTestClient(t, 200, []v1alpha1.HealthCheckProfile{{ObjectMeta: metav1.ObjectMeta{Name: "default"}}}, &rec)
		profiles, err := client.HealthCheckProfiles.List(ctx, "proj")
		if err != nil {
			t.Fatalf("List: %v", err)
		}
		if rec.Method != "GET" || rec.Path != "/project/proj/healthcheckprofiles" {
			t.Fatalf("unexpected request %s %s", rec.Method, rec.Path)
		}
		if len(profiles) != 1 {
			t.Fatalf("unexpected profiles %+v", profiles)
		}
	})

	t.Run("Get", func(t *testing.T) {
		var rec recordedRequest
		client := newTestClient(t, 200, v1alpha1.HealthCheckProfile{ObjectMeta: metav1.ObjectMeta{Name: "default"}}, &rec)
		profile, err := client.HealthCheckProfiles.Get(ctx, "proj", "default")
		if err != nil {
			t.Fatalf("Get: %v", err)
		}
		if rec.Method != "GET" || rec.Path != "/project/proj/healthcheckprofiles/default" {
			t.Fatalf("unexpected request %s %s", rec.Method, rec.Path)
		}
		if profile.Name != "default" {
			t.Fatalf("unexpected profile %+v", profile)
		}
	})

	t.Run("Create", func(t *testing.T) {
		var rec recordedRequest
		client := newTestClient(t, 201, v1alpha1.HealthCheckProfile{ObjectMeta: metav1.ObjectMeta{Name: "default"}}, &rec)
		_, err := client.HealthCheckProfiles.Create(ctx, "proj", CreateHealthCheckProfileRequest{
			Name:   "default",
			Probes: []v1alpha1.HealthCheckProbeSpec{{}},
		})
		if err != nil {
			t.Fatalf("Create: %v", err)
		}
		if rec.Method != "POST" || rec.Path != "/project/proj/healthcheckprofiles" {
			t.Fatalf("unexpected request %s %s", rec.Method, rec.Path)
		}
	})

	t.Run("Update", func(t *testing.T) {
		var rec recordedRequest
		client := newTestClient(t, 200, v1alpha1.HealthCheckProfile{ObjectMeta: metav1.ObjectMeta{Name: "default"}}, &rec)
		_, err := client.HealthCheckProfiles.Update(ctx, "proj", "default", UpdateHealthCheckProfileRequest{
			Probes: &[]v1alpha1.HealthCheckProbeSpec{{}},
		})
		if err != nil {
			t.Fatalf("Update: %v", err)
		}
		if rec.Method != "PATCH" || rec.Path != "/project/proj/healthcheckprofiles/default" {
			t.Fatalf("unexpected request %s %s", rec.Method, rec.Path)
		}
	})

	t.Run("Delete", func(t *testing.T) {
		var rec recordedRequest
		client := newTestClient(t, 204, nil, &rec)
		if err := client.HealthCheckProfiles.Delete(ctx, "proj", "default"); err != nil {
			t.Fatalf("Delete: %v", err)
		}
		if rec.Method != "DELETE" || rec.Path != "/project/proj/healthcheckprofiles/default" {
			t.Fatalf("unexpected request %s %s", rec.Method, rec.Path)
		}
	})
}
