package harbor

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/pascal71/hrbcli/pkg/api"
)

func TestLabelServiceAddToArtifactEscapesPath(t *testing.T) {
	var gotPath string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotPath = r.URL.EscapedPath()
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client := &api.Client{
		BaseURL:    server.URL,
		APIVersion: "v2.0",
		HTTPClient: server.Client(),
	}

	svc := NewLabelService(client)

	project := "my project"
	repo := "my/repo"
	ref := "sha256:abc"
	if err := svc.AddToArtifact(project, repo, ref, 1); err != nil {
		t.Fatalf("AddToArtifact error: %v", err)
	}

	expected := "/api/v2.0/projects/" + url.PathEscape(project) + "/repositories/" + url.PathEscape(repo) + "/artifacts/" + url.PathEscape(ref) + "/labels"
	if gotPath != expected {
		t.Fatalf("unexpected path: got %s want %s", gotPath, expected)
	}
}
