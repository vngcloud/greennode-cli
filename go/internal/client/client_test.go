package client

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
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

func TestFormatErrorSurfacesNestedErrorObject(t *testing.T) {
	// VKS returns errors as {"error": {"message": ...}} — a nested object, not a
	// string. The detail must still reach the user instead of being dropped.
	body := []byte(`{"error":{"message":"KubeConfig can only be requested when the cluster is ACTIVE."}}`)
	got := formatError(http.StatusBadRequest, body)
	if !strings.Contains(got, "cluster is ACTIVE") {
		t.Errorf("formatError = %q, want it to contain the nested error message", got)
	}
}

func TestFormatErrorUsesPlainStringMessage(t *testing.T) {
	body := []byte(`{"message":"boom"}`)
	got := formatError(http.StatusBadRequest, body)
	if !strings.Contains(got, "boom") {
		t.Errorf("formatError = %q, want it to contain %q", got, "boom")
	}
}
