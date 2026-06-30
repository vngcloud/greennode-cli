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

func TestParseStructFlagShorthand(t *testing.T) {
	got, err := ParseStructFlag("minSize=2,maxSize=10", "minSize", "maxSize")
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if got["minSize"] != 2 || got["maxSize"] != 10 {
		t.Errorf("got %#v, want minSize=2 maxSize=10 (ints)", got)
	}
}

func TestParseStructFlagShorthandStringsStay(t *testing.T) {
	got, err := ParseStructFlag("type=NEW,placementGroupName=pg-1")
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if got["type"] != "NEW" || got["placementGroupName"] != "pg-1" {
		t.Errorf("got %#v", got)
	}
}

func TestParseStructFlagJSON(t *testing.T) {
	got, err := ParseStructFlag(`{"minSize":2,"maxSize":10}`, "minSize", "maxSize")
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if got["minSize"].(float64) != 2 || got["maxSize"].(float64) != 10 {
		t.Errorf("got %#v", got)
	}
}

func TestParseStructFlagEmpty(t *testing.T) {
	got, err := ParseStructFlag("  ")
	if err != nil || got != nil {
		t.Errorf("empty should be (nil,nil), got %#v err %v", got, err)
	}
}

func TestParseStructFlagErrors(t *testing.T) {
	if _, err := ParseStructFlag("{bad json"); err == nil {
		t.Error("malformed JSON should error")
	}
	if _, err := ParseStructFlag("minSize"); err == nil {
		t.Error("shorthand pair without '=' should error")
	}
	if _, err := ParseStructFlag("minSize=abc", "minSize"); err == nil {
		t.Error("non-integer int field should error")
	}
}
