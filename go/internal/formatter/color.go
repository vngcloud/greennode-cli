package formatter

import (
	"io"
	"os"
	"strings"
)

// ANSI SGR codes used for status coloring.
const (
	ansiReset  = "\033[0m"
	ansiRed    = "\033[31m"
	ansiGreen  = "\033[32m"
	ansiYellow = "\033[33m"
)

// statusColors maps an upper-cased status value to an ANSI color. Values are
// the lifecycle/health states VKS and vServer return; anything not listed is
// left uncolored.
var statusColors = map[string]string{
	// Healthy / done — green.
	"ACTIVE": ansiGreen, "RUNNING": ansiGreen, "READY": ansiGreen,
	"AVAILABLE": ansiGreen, "HEALTHY": ansiGreen, "SUCCEEDED": ansiGreen,
	"COMPLETED": ansiGreen,
	// In progress — yellow.
	"CREATING": ansiYellow, "DELETING": ansiYellow, "UPDATING": ansiYellow,
	"UPGRADING": ansiYellow, "PENDING": ansiYellow, "PROVISIONING": ansiYellow,
	"RECONCILING": ansiYellow,
	// Failure — red.
	"ERROR": ansiRed, "FAILED": ansiRed, "UNHEALTHY": ansiRed, "DEGRADED": ansiRed,
}

// ColorEnabled decides whether output should be colored, mirroring the AWS CLI
// --color flag: "on" always, "off" never, "auto" (or empty) only when w is a
// terminal. Honors the NO_COLOR convention in auto mode.
func ColorEnabled(mode string, w io.Writer) bool {
	switch mode {
	case "on":
		return true
	case "off":
		return false
	default: // "auto" or unset
		if os.Getenv("NO_COLOR") != "" {
			return false
		}
		return isTerminal(w)
	}
}

func isTerminal(w io.Writer) bool {
	f, ok := w.(*os.File)
	if !ok {
		return false
	}
	fi, err := f.Stat()
	if err != nil {
		return false
	}
	return fi.Mode()&os.ModeCharDevice != 0
}

// colorCell wraps an already-padded cell in an ANSI color when color is enabled
// and raw is a recognized status. Padding is done by the caller on the raw text
// so column widths stay correct — the ANSI codes are zero-width and added
// around the padded string.
func colorCell(padded, raw string, color bool) string {
	if !color {
		return padded
	}
	if code, ok := statusColors[strings.ToUpper(strings.TrimSpace(raw))]; ok {
		return code + padded + ansiReset
	}
	return padded
}
