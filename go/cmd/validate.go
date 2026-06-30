package cmd

import (
	"fmt"
	"slices"
	"strings"

	"github.com/spf13/cobra"
)

var validOutputFormats = []string{"json", "text", "table"}

// validateGlobalFlags rejects invalid values for global flags before any work
// runs, mirroring `aws`/`gcloud` which fail fast on an unknown --output instead
// of silently falling back. Only validates flags the user explicitly set.
func validateGlobalFlags(cmd *cobra.Command) error {
	if cmd.Flags().Changed("output") {
		v, _ := cmd.Flags().GetString("output")
		if err := validateOutputFormat(v); err != nil {
			return err
		}
	}
	return nil
}

// validateOutputFormat returns an error when value is a non-empty,
// unrecognized output format. The message echoes the bad value, suggests the
// closest valid format when one is near, and always lists the valid formats.
func validateOutputFormat(value string) error {
	if value == "" {
		return nil
	}
	if slices.Contains(validOutputFormats, value) {
		return nil
	}
	msg := fmt.Sprintf("invalid output format: '%s'", value)
	if s := suggestClosest(value, validOutputFormats); s != "" {
		msg += fmt.Sprintf(", maybe you meant '%s'", s)
	}
	msg += fmt.Sprintf(" (valid formats: %s)", strings.Join(validOutputFormats, ", "))
	return fmt.Errorf("%s", msg)
}

// suggestClosest returns the option nearest to value by Levenshtein distance,
// or "" when the nearest is too far to be a likely typo (distance > 2).
func suggestClosest(value string, options []string) string {
	best := ""
	bestDist := 1 << 30
	for _, opt := range options {
		if d := levenshtein(value, opt); d < bestDist {
			bestDist, best = d, opt
		}
	}
	if bestDist > 2 {
		return ""
	}
	return best
}

func levenshtein(a, b string) int {
	ra, rb := []rune(a), []rune(b)
	prev := make([]int, len(rb)+1)
	curr := make([]int, len(rb)+1)
	for j := range prev {
		prev[j] = j
	}
	for i := 1; i <= len(ra); i++ {
		curr[0] = i
		for j := 1; j <= len(rb); j++ {
			cost := 1
			if ra[i-1] == rb[j-1] {
				cost = 0
			}
			curr[j] = min(prev[j]+1, curr[j-1]+1, prev[j-1]+cost)
		}
		prev, curr = curr, prev
	}
	return prev[len(rb)]
}
