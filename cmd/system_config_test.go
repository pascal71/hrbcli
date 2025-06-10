package cmd

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestParseConfigValue(t *testing.T) {
	tests := []struct {
		in   string
		want interface{}
	}{
		{"true", true},
		{"false", false},
		{"42", 42},
		{"abc", "abc"},
	}

	for _, tt := range tests {
		got := parseConfigValue(tt.in)
		if got != tt.want {
			t.Errorf("parseConfigValue(%q) = %#v, want %#v", tt.in, got, tt.want)
		}
	}
}

func TestSystemConfigSet(t *testing.T) {
	var reqs []recordedReq
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		r.Body.Close()
		reqs = append(reqs, recordedReq{r.URL.Path, r.Method, body})
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	setupEnv(server.URL)
	defer os.Unsetenv("HARBOR_URL")

	cmd := newSystemConfigSetCmd()
	if err := cmd.RunE(cmd, []string{"read_only", "true"}); err != nil {
		t.Fatalf("run error: %v", err)
	}

	if len(reqs) != 1 {
		t.Fatalf("expected 1 request, got %d", len(reqs))
	}
	if reqs[0].Path != "/api/v2.0/configurations" || reqs[0].Method != http.MethodPut {
		t.Fatalf("unexpected request: %+v", reqs[0])
	}
	var m map[string]interface{}
	json.Unmarshal(reqs[0].Body, &m)
	if v, ok := m["read_only"].(bool); !ok || !v {
		t.Fatalf("value not set correctly: %v", m)
	}
}
