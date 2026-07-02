package cli

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strings"
)

// DryRunNotice prints the standard footer shown at the end of a --dry-run
// preview. verb is the action that would run, e.g. DryRunNotice("delete")
// prints "Run without --dry-run to delete."
func DryRunNotice(verb string) {
	fmt.Printf("\nRun without --dry-run to %s.\n", verb)
}

// PrintDryRun prints a consistent --dry-run preview for a mutating request: a
// header, the target being changed, the request body (keys sorted for stable
// output), and the standard footer. verb is the action (e.g. "update",
// "upgrade", "configure").
func PrintDryRun(verb, target string, body map[string]interface{}) {
	fmt.Println("=== DRY RUN ===")
	if target != "" {
		fmt.Printf("Would %s %s:\n", verb, target)
	}
	keys := make([]string, 0, len(body))
	for k := range body {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		fmt.Printf("  %s: %v\n", k, body[k])
	}
	DryRunNotice(verb)
}

// Confirm asks the user to confirm a destructive action and reports whether to
// proceed. It returns true immediately when force is true. The prompt is shown
// as "<prompt> [y/N]: "; only "y" or "yes" (case-insensitive) proceeds.
func Confirm(force bool, prompt string) bool {
	if force {
		return true
	}
	fmt.Printf("\n%s [y/N]: ", prompt)
	reader := bufio.NewReader(os.Stdin)
	answer, _ := reader.ReadString('\n')
	answer = strings.TrimSpace(strings.ToLower(answer))
	return answer == "y" || answer == "yes"
}
