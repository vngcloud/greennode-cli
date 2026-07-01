package cli

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// DryRunNotice prints the standard footer shown at the end of a --dry-run
// preview. verb is the action that would run, e.g. DryRunNotice("delete")
// prints "Run without --dry-run to delete."
func DryRunNotice(verb string) {
	fmt.Printf("\nRun without --dry-run to %s.\n", verb)
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
