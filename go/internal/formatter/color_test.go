package formatter

import (
	"bytes"
	"strings"
	"testing"
)

func TestColorEnabledModes(t *testing.T) {
	var buf bytes.Buffer // not an *os.File -> isTerminal() false
	if !ColorEnabled("on", &buf) {
		t.Error("on should force color even for a non-terminal writer")
	}
	if ColorEnabled("off", &buf) {
		t.Error("off should never color")
	}
	if ColorEnabled("auto", &buf) {
		t.Error("auto should be off when the writer is not a terminal")
	}
	if ColorEnabled("", &buf) {
		t.Error("empty (auto) should be off for a non-terminal writer")
	}
}

func TestColorEnabledNoColorEnv(t *testing.T) {
	t.Setenv("NO_COLOR", "1")
	// on is explicit and ignores NO_COLOR; auto honors it.
	if !ColorEnabled("on", &bytes.Buffer{}) {
		t.Error("on should still color even with NO_COLOR set")
	}
	if ColorEnabled("auto", &bytes.Buffer{}) {
		t.Error("auto should be off when NO_COLOR is set")
	}
}

// With color on, a status value in table output is wrapped in ANSI codes, and
// column alignment is preserved (padding is applied to the raw text).
func TestFormatTableColorsStatus(t *testing.T) {
	data := map[string]interface{}{"id": "c1", "status": "ACTIVE"}
	var buf bytes.Buffer
	formatTable(data, &buf, true)
	out := buf.String()

	if !strings.Contains(out, ansiGreen+"ACTIVE") {
		t.Errorf("ACTIVE should be green-wrapped, got:\n%q", out)
	}
	if !strings.Contains(out, ansiReset) {
		t.Errorf("colored output should reset, got:\n%q", out)
	}
}

// With color off, output is plain (no ANSI codes) — unchanged from before.
func TestFormatTableNoColorByDefault(t *testing.T) {
	data := map[string]interface{}{"id": "c1", "status": "ERROR"}
	var buf bytes.Buffer
	formatTable(data, &buf, false)
	if strings.Contains(buf.String(), "\033[") {
		t.Errorf("color off should emit no ANSI codes, got:\n%q", buf.String())
	}
}

// Non-status values are never colored, even with color on.
func TestFormatTableLeavesNonStatusUncolored(t *testing.T) {
	data := map[string]interface{}{"id": "c1", "name": "demo"}
	var buf bytes.Buffer
	formatTable(data, &buf, true)
	if strings.Contains(buf.String(), "\033[") {
		t.Errorf("no recognized status -> no color, got:\n%q", buf.String())
	}
}
