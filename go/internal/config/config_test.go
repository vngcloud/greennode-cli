package config

import (
	"sort"
	"testing"
)

func TestRegionNames(t *testing.T) {
	got := RegionNames()
	sort.Strings(got)
	if len(got) < 2 {
		t.Errorf("RegionNames = %#v, want >=2 regions", got)
	}
	found := map[string]bool{}
	for _, r := range got {
		found[r] = true
	}
	if !found["HCM-3"] || !found["HAN"] {
		t.Errorf("RegionNames = %#v, want HCM-3 and HAN", got)
	}
}
