package cli

import (
	"reflect"
	"testing"
)

func TestBuildEventsQueryOnlyIncludesSetValues(t *testing.T) {
	got := BuildEventsQuery("CREATE", "", 2, 50, map[string]bool{
		"action": true, "type": false, "page": true, "page-size": true,
	})
	want := map[string]string{
		"action":   "CREATE",
		"page":     "2",
		"pageSize": "50",
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("BuildEventsQuery = %#v, want %#v", got, want)
	}
}

func TestBuildEventsQueryEmptyWhenNothingSet(t *testing.T) {
	got := BuildEventsQuery("", "", 0, 0, map[string]bool{})
	if len(got) != 0 {
		t.Errorf("BuildEventsQuery = %#v, want empty map", got)
	}
}

func TestParseCommaSeparated(t *testing.T) {
	got := ParseCommaSeparated(" a, b ,,c ")
	want := []string{"a", "b", "c"}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("ParseCommaSeparated = %#v, want %#v", got, want)
	}
	if ParseCommaSeparated("") != nil {
		t.Errorf("ParseCommaSeparated(\"\") should be nil")
	}
}
