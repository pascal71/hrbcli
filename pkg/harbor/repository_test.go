package harbor

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/pascal71/hrbcli/pkg/api"
)

func TestRepositoryServiceListEscapesProjectName(t *testing.T) {
	var gotPath string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotPath = r.URL.EscapedPath()
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("[]"))
	}))
	defer server.Close()

	client := &api.Client{
		BaseURL:    server.URL,
		APIVersion: "v2.0",
		HTTPClient: server.Client(),
	}

	svc := NewRepositoryService(client)

	project := "my project/with/slash"
	if _, err := svc.List(project, nil); err != nil {
		t.Fatalf("List error: %v", err)
	}

	expected := fmt.Sprintf("/api/v2.0/projects/%s/repositories", url.PathEscape(project))
	if gotPath != expected {
		t.Fatalf("unexpected path: got %s want %s", gotPath, expected)
	}
}
