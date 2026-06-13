package client

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/vngcloud/greennode-cli/internal/auth"
)

func TestPatchSendsPatchMethodAndBody(t *testing.T) {
	var gotMethod, gotBody string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotMethod = r.Method
		b, _ := io.ReadAll(r.Body)
		gotBody = string(b)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"ok":true}`))
	}))
	defer srv.Close()

	tm := auth.NewTokenManager("id", "secret")
	// Pre-seed a static token so GetToken never calls the real IAM endpoint.
	tm.SetToken("test-token", time.Now().Add(1*time.Hour))

	c := NewGreenodeClient(srv.URL, tm, 5*time.Second, false, false)

	_, err := c.Patch("/v1/thing", map[string]interface{}{"enableAutoHealing": true})
	if err != nil {
		t.Fatalf("Patch returned error: %v", err)
	}
	if gotMethod != http.MethodPatch {
		t.Errorf("method = %q, want PATCH", gotMethod)
	}
	if gotBody != `{"enableAutoHealing":true}` {
		t.Errorf("body = %q, want enableAutoHealing payload", gotBody)
	}
}
