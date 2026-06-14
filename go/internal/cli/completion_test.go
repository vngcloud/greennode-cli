package cli

import (
	"context"
	"errors"
	"reflect"
	"testing"
	"time"

	"github.com/spf13/cobra"
)

func TestExtractIDsPrefersKnownKeyAndDedups(t *testing.T) {
	// Multiple slice fields: the known wrapper "items" must win deterministically
	// over an unrelated slice field, regardless of map iteration order.
	r := map[string]interface{}{
		"errors": []interface{}{map[string]interface{}{"id": "e1"}},
		"items": []interface{}{
			map[string]interface{}{"id": "a"},
			map[string]interface{}{"id": "a"}, // duplicate → deduped
			map[string]interface{}{"id": "b"},
		},
	}
	for i := 0; i < 20; i++ { // repeat to expose map-order nondeterminism
		if got := ExtractIDs(r, "id"); !reflect.DeepEqual(got, []string{"a", "b"}) {
			t.Fatalf("iter %d: ExtractIDs = %#v, want [a b]", i, got)
		}
	}
}

func TestExtractIDs(t *testing.T) {
	r1 := map[string]interface{}{"items": []interface{}{
		map[string]interface{}{"id": "c1", "name": "a"},
		map[string]interface{}{"id": "c2", "name": "b"},
	}}
	if got := ExtractIDs(r1, "id"); !reflect.DeepEqual(got, []string{"c1", "c2"}) {
		t.Errorf("items/id = %#v", got)
	}
	r2 := map[string]interface{}{"volumeTypes": []interface{}{
		map[string]interface{}{"uuid": "u1"},
		map[string]interface{}{"id": "i2"},
	}}
	if got := ExtractIDs(r2, "uuid", "id"); !reflect.DeepEqual(got, []string{"u1", "i2"}) {
		t.Errorf("uuid|id = %#v", got)
	}
	r3 := []interface{}{map[string]interface{}{"uuid": "s1"}}
	if got := ExtractIDs(r3, "uuid"); !reflect.DeepEqual(got, []string{"s1"}) {
		t.Errorf("top-array = %#v", got)
	}
	if got := ExtractIDs(nil, "id"); got != nil {
		t.Errorf("nil = %#v, want nil", got)
	}
	if got := ExtractIDs(map[string]interface{}{"x": 1}, "id"); got != nil {
		t.Errorf("no slice = %#v, want nil", got)
	}
}

func TestFlagValuesFiltersByPrefix(t *testing.T) {
	f := FlagValues("json", "text", "table")
	got, dir := f(nil, nil, "t")
	if !reflect.DeepEqual(got, []string{"text", "table"}) {
		t.Errorf("got %#v", got)
	}
	if dir != cobra.ShellCompDirectiveNoFileComp {
		t.Errorf("dir = %v", dir)
	}
}

func TestFlagValuesFrom(t *testing.T) {
	f := FlagValuesFrom(func() []string { return []string{"HCM-3", "HAN"} })
	got, _ := f(nil, nil, "H")
	if !reflect.DeepEqual(got, []string{"HCM-3", "HAN"}) {
		t.Errorf("got %#v", got)
	}
}

func TestFlagFromAPISuccessAndError(t *testing.T) {
	ok := FlagFromAPI(func(ctx context.Context, cmd *cobra.Command) ([]string, error) {
		return []string{"alpha", "beta"}, nil
	})
	got, dir := ok(nil, nil, "al")
	if !reflect.DeepEqual(got, []string{"alpha"}) || dir != cobra.ShellCompDirectiveNoFileComp {
		t.Errorf("success = %#v %v", got, dir)
	}
	bad := FlagFromAPI(func(ctx context.Context, cmd *cobra.Command) ([]string, error) {
		return nil, errors.New("boom")
	})
	if got, _ := bad(nil, nil, ""); got != nil {
		t.Errorf("error path = %#v, want nil", got)
	}
}

func TestFlagFromAPITimesOut(t *testing.T) {
	old := completionTimeout
	completionTimeout = 20 * time.Millisecond
	defer func() { completionTimeout = old }()
	slow := FlagFromAPI(func(ctx context.Context, cmd *cobra.Command) ([]string, error) {
		<-ctx.Done()
		return nil, ctx.Err()
	})
	if got, _ := slow(nil, nil, ""); got != nil {
		t.Errorf("timeout path = %#v, want nil", got)
	}
}

func TestResourceCompleterRegistry(t *testing.T) {
	RegisterResourceCompleter("test:thing", FlagValues("x1", "x2"))
	got, _ := ResourceCompletion("test:thing")(nil, nil, "x")
	if len(got) != 2 {
		t.Errorf("registered lookup = %#v", got)
	}
	got2, dir := ResourceCompletion("test:missing")(nil, nil, "")
	if got2 != nil || dir != cobra.ShellCompDirectiveNoFileComp {
		t.Errorf("missing key = %#v %v", got2, dir)
	}
}
