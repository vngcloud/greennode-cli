package formatter

import (
	"bytes"
	"strings"
	"testing"
)

// A detail (single-object) response with nested array fields must render as a
// key/value table — not be hijacked into an empty table by a nested array.
func TestFormatTableDetailObject(t *testing.T) {
	data := map[string]interface{}{
		"id":            "k8s-123",
		"name":          "demo",
		"status":        "ACTIVE",
		"listSubnetIds": []interface{}{}, // empty nested array — must NOT blank the output
	}
	var buf bytes.Buffer
	formatTable(data, &buf, false)
	out := buf.String()

	if strings.TrimSpace(out) == "" {
		t.Fatalf("detail object rendered empty table")
	}
	for _, want := range []string{"id", "k8s-123", "name", "demo", "status", "ACTIVE"} {
		if !strings.Contains(out, want) {
			t.Errorf("table output missing %q\n%s", want, out)
		}
	}
}

// A list response ({items:[...]}) must render as a multi-column table with a
// header row and one row per item.
func TestFormatTableListResponse(t *testing.T) {
	data := map[string]interface{}{
		"items": []interface{}{
			map[string]interface{}{"id": "c1", "name": "alpha"},
			map[string]interface{}{"id": "c2", "name": "beta"},
		},
		"total": float64(2),
	}
	var buf bytes.Buffer
	formatTable(data, &buf, false)
	out := buf.String()

	for _, want := range []string{"id", "name", "c1", "alpha", "c2", "beta"} {
		if !strings.Contains(out, want) {
			t.Errorf("list table missing %q\n%s", want, out)
		}
	}
	// header order must be deterministic (sorted): "id" column before "name"
	if strings.Index(out, "id") > strings.Index(out, "name") {
		t.Errorf("headers not sorted deterministically:\n%s", out)
	}
}

// Top-level array renders as a multi-column table too.
func TestFormatTableTopLevelArray(t *testing.T) {
	data := []interface{}{
		map[string]interface{}{"uuid": "s1"},
		map[string]interface{}{"uuid": "s2"},
	}
	var buf bytes.Buffer
	formatTable(data, &buf, false)
	out := buf.String()
	for _, want := range []string{"uuid", "s1", "s2"} {
		if !strings.Contains(out, want) {
			t.Errorf("top-array table missing %q\n%s", want, out)
		}
	}
}
