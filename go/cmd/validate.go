package cmd

import (
	"fmt"
	"slices"
	"strings"

	"github.com/spf13/cobra"
)

var (
	validOutputFormats = []string{"json", "text", "table"}
	validColorModes    = []string{"on", "off", "auto"}
)

// validateGlobalFlags rejects invalid values for global flags before any work
// runs, mirroring `aws`/`gcloud` which fail fast on an unknown value instead of
// silently falling back. Only validates flags the user explicitly set.
func validateGlobalFlags(cmd *cobra.Command) error {
	checks := []struct {
		flag  string
		label string
		valid []string
	}{
		{"output", "output format", validOutputFormats},
		{"color", "color mode", validColorModes},
	}
	for _, c := range checks {
		if !cmd.Flags().Changed(c.flag) {
			continue
		}
		v, _ := cmd.Flags().GetString(c.flag)
		if err := validateChoice(c.label, v, c.valid); err != nil {
			return err
		}
	}
	return nil
}

// validateOutputFormat returns an error when value is a non-empty,
// unrecognized output format.
func validateOutputFormat(value string) error {
	return validateChoice("output format", value, validOutputFormats)
}

// validateChoice returns an error when value is non-empty and not one of valid.
// The message echoes the bad value, suggests the closest valid option when one
// is near, and always lists the valid options. label names the thing being
// validated (e.g. "output format").
func validateChoice(label, value string, valid []string) error {
	if value == "" {
		return nil
	}
	if slices.Contains(valid, value) {
		return nil
	}
	msg := fmt.Sprintf("invalid %s: '%s'", label, value)
	if s := suggestClosest(value, valid); s != "" {
		msg += fmt.Sprintf(", maybe you meant '%s'", s)
	}
	msg += fmt.Sprintf(" (valid values: %s)", strings.Join(valid, ", "))
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
