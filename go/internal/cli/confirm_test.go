package cli

import "testing"

func TestConfirmForceSkipsPrompt(t *testing.T) {
	if !Confirm(true, "delete this?") {
		t.Error("force=true should proceed without prompting")
	}
}
