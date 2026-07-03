package client

import (
	"strings"
	"testing"
)

func TestRedactDebugBodyMasksKubeconfig(t *testing.T) {
	in := `{"status":"ACTIVE","kubeConfig":"apiVersion: v1\nclient-key-data: SUPERSECRETKEY","renewalWarning":false}`
	out := redactDebugBody(in)
	if strings.Contains(out, "SUPERSECRETKEY") {
		t.Errorf("kubeConfig value not redacted: %s", out)
	}
	if !strings.Contains(out, "[REDACTED]") {
		t.Errorf("expected [REDACTED] marker: %s", out)
	}
	if !strings.Contains(out, "ACTIVE") {
		t.Errorf("non-sensitive field lost: %s", out)
	}
}

func TestRedactDebugBodyNestedAndTokens(t *testing.T) {
	in := `{"data":{"token":"abc123","name":"ok"},"clientSecret":"shh","items":[{"bearerToken":"xyz"}]}`
	out := redactDebugBody(in)
	for _, leaked := range []string{"abc123", "shh", "xyz"} {
		if strings.Contains(out, leaked) {
			t.Errorf("secret %q leaked: %s", leaked, out)
		}
	}
	if !strings.Contains(out, "ok") {
		t.Errorf("non-sensitive value lost: %s", out)
	}
}

func TestRedactDebugBodyKeepsNonSensitive(t *testing.T) {
	in := `{"id":"cls-1","numNodes":3,"status":"CREATING"}`
	out := redactDebugBody(in)
	for _, keep := range []string{"cls-1", "3", "CREATING"} {
		if !strings.Contains(out, keep) {
			t.Errorf("non-sensitive %q should be kept: %s", keep, out)
		}
	}
	if strings.Contains(out, "[REDACTED]") {
		t.Errorf("nothing should be redacted here: %s", out)
	}
}

func TestRedactDebugBodyNonJSONPassthrough(t *testing.T) {
	in := "not json"
	if got := redactDebugBody(in); got != in {
		t.Errorf("non-JSON should pass through unchanged, got %q", got)
	}
}
