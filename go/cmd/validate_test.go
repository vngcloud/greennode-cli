package cmd

import (
	"strings"
	"testing"
)

func TestValidateOutputFormatValid(t *testing.T) {
	for _, v := range []string{"json", "text", "table"} {
		if err := validateOutputFormat(v); err != nil {
			t.Errorf("%q should be valid, got %v", v, err)
		}
	}
}

func TestValidateOutputFormatEmptyIsValid(t *testing.T) {
	// Empty means the flag was not provided; default resolution handles it.
	if err := validateOutputFormat(""); err != nil {
		t.Errorf("empty should be valid, got %v", err)
	}
}

func TestValidateOutputFormatTypoSuggests(t *testing.T) {
	err := validateOutputFormat("tabel")
	if err == nil {
		t.Fatal("typo 'tabel' should be rejected")
	}
	msg := err.Error()
	if !strings.Contains(msg, "tabel") {
		t.Errorf("error should echo the bad value, got %q", msg)
	}
	if !strings.Contains(msg, "table") {
		t.Errorf("error should suggest 'table' (did you mean), got %q", msg)
	}
}

func TestValidateOutputFormatUnknownNoSuggestion(t *testing.T) {
	err := validateOutputFormat("xml")
	if err == nil {
		t.Fatal("'xml' should be rejected")
	}
	// Far from any valid value: no bogus suggestion, but valid list is shown.
	if !strings.Contains(err.Error(), "json") {
		t.Errorf("error should list valid formats, got %q", err.Error())
	}
}

func TestValidateChoiceColor(t *testing.T) {
	for _, v := range []string{"on", "off", "auto", ""} {
		if err := validateChoice("color mode", v, validColorModes); err != nil {
			t.Errorf("%q should be valid, got %v", v, err)
		}
	}
	err := validateChoice("color mode", "atuo", validColorModes)
	if err == nil {
		t.Fatal("typo 'atuo' should be rejected")
	}
	if !strings.Contains(err.Error(), "auto") {
		t.Errorf("error should suggest 'auto', got %q", err.Error())
	}
}

func TestSuggestClosest(t *testing.T) {
	opts := []string{"json", "text", "table"}
	cases := map[string]string{
		"tabel": "table",
		"jsno":  "json",
		"txt":   "text",
		"xml":   "", // too far -> no suggestion
	}
	for in, want := range cases {
		if got := suggestClosest(in, opts); got != want {
			t.Errorf("suggestClosest(%q) = %q, want %q", in, got, want)
		}
	}
}
