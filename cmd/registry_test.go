package cmd

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/pascal71/hrbcli/pkg/api"
	"github.com/spf13/viper"
)

type recordedReq struct {
	Path   string
	Method string
	Body   []byte
}

func setupEnv(url string) {
	os.Setenv("HARBOR_URL", url)
	viper.Reset()
	initConfig()
}

func TestNewRegistryCreateCmd(t *testing.T) {
	var reqs []recordedReq
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		r.Body.Close()
		reqs = append(reqs, recordedReq{r.URL.Path, r.Method, body})
		switch r.URL.Path {
		case "/api/v2.0/registries/ping":
			w.WriteHeader(http.StatusOK)
		case "/api/v2.0/registries":
			w.Header().Set("Location", "/api/v2.0/registries/1")
			w.WriteHeader(http.StatusCreated)
		case "/api/v2.0/registries/1":
			w.Header().Set("Content-Type", "application/json")
			w.Write([]byte(`{"id":1,"name":"myreg","url":"https://example.com","type":"docker-hub"}`))
		default:
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	defer server.Close()

	setupEnv(server.URL)
	defer os.Unsetenv("HARBOR_URL")

	cmd := newRegistryCreateCmd()
	cmd.Flags().Set("url", "https://example.com")
	cmd.Flags().Set("type", "docker-hub")
	if err := cmd.RunE(cmd, []string{"myreg"}); err != nil {
		t.Fatalf("run error: %v", err)
	}

	if len(reqs) != 3 {
		t.Fatalf("expected 3 requests, got %d", len(reqs))
	}
	var pingReq api.RegistryReq
	json.Unmarshal(reqs[0].Body, &pingReq)
	if pingReq.Name != "myreg" || pingReq.URL != "https://example.com" {
		t.Fatalf("unexpected ping request: %+v", pingReq)
	}
}

func TestNewRegistryUpdateCmd(t *testing.T) {
	var reqs []recordedReq
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		r.Body.Close()
		reqs = append(reqs, recordedReq{r.URL.Path, r.Method, body})
		switch r.URL.Path {
		case "/api/v2.0/registries/1":
			if r.Method == http.MethodGet {
				w.Header().Set("Content-Type", "application/json")
				w.Write([]byte(`{"id":1,"name":"myreg","url":"https://example.com","description":"old","type":"docker-hub","insecure":false}`))
			} else {
				w.WriteHeader(http.StatusOK)
			}
		default:
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	defer server.Close()

	setupEnv(server.URL)
	defer os.Unsetenv("HARBOR_URL")

	cmd := newRegistryUpdateCmd()
	cmd.Flags().Set("description", "newdesc")
	if err := cmd.RunE(cmd, []string{"1"}); err != nil {
		t.Fatalf("run error: %v", err)
	}

	if len(reqs) != 2 {
		t.Fatalf("expected 2 requests, got %d", len(reqs))
	}
	var updateReq api.RegistryReq
	json.Unmarshal(reqs[1].Body, &updateReq)
	if updateReq.Description != "newdesc" {
		t.Fatalf("description not updated: %+v", updateReq)
	}
}
