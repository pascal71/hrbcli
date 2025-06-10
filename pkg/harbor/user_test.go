package harbor

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/pascal71/hrbcli/pkg/api"
)

func TestUserServiceCreateParsesLocation(t *testing.T) {
	var server *httptest.Server
	server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/api/v2.0/users":
			w.Header().Set("Location", server.URL+"/api/v2.0/users/7")
			w.WriteHeader(http.StatusCreated)
		case "/api/v2.0/users/7":
			w.Header().Set("Content-Type", "application/json")
			w.Write([]byte(`{"user_id":7,"username":"foo","email":"bar@example.com"}`))
		default:
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	defer server.Close()

	client := &api.Client{
		BaseURL:    server.URL,
		APIVersion: "v2.0",
		HTTPClient: server.Client(),
	}

	svc := NewUserService(client)
	user, err := svc.Create(&api.UserReq{Username: "foo", Email: "bar@example.com", Password: "x"})
	if err != nil {
		t.Fatalf("Create error: %v", err)
	}
	if user == nil || user.UserID != 7 {
		t.Fatalf("unexpected user: %+v", user)
	}
}
