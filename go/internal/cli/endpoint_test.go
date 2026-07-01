package cli

import "testing"

func TestIsTrustedEndpoint(t *testing.T) {
	cases := map[string]bool{
		"":                                   true,  // no override -> built-in endpoint
		"https://vks.api.vngcloud.vn":        true,
		"https://hcm-3.api.vngcloud.vn/x":    true,
		"https://vngcloud.vn":                true,
		"http://attacker.com":                false,
		"https://evil.vngcloud.vn.attacker.com": false, // suffix must be a real domain boundary
		"http://localhost:8080":              false,
		"not-a-url ::::":                     false,
	}
	for in, want := range cases {
		if got := IsTrustedEndpoint(in); got != want {
			t.Errorf("IsTrustedEndpoint(%q) = %v, want %v", in, got, want)
		}
	}
}
