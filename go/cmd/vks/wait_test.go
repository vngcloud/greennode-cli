package vks

import (
	"errors"
	"testing"

	"github.com/vngcloud/greennode-cli/internal/client"
)

func TestEvaluateActive(t *testing.T) {
	cases := []struct {
		name             string
		result           interface{}
		err              error
		wantDone, wantFail bool
	}{
		{"active", map[string]interface{}{"status": "ACTIVE"}, nil, true, false},
		{"creating", map[string]interface{}{"status": "CREATING"}, nil, false, false},
		{"error", map[string]interface{}{"status": "ERROR"}, nil, false, true},
		{"failed", map[string]interface{}{"status": "FAILED"}, nil, false, true},
		{"transient err", nil, errors.New("boom"), false, false},
	}
	for _, tc := range cases {
		done, failed, _ := evaluateActive(tc.result, tc.err)
		if done != tc.wantDone || failed != tc.wantFail {
			t.Errorf("%s: done=%v failed=%v, want done=%v failed=%v", tc.name, done, failed, tc.wantDone, tc.wantFail)
		}
	}
}

func TestEvaluateDeleted(t *testing.T) {
	notFound := &client.APIError{StatusCode: 404, Body: ""}
	otherAPIErr := &client.APIError{StatusCode: 500, Body: ""}

	cases := []struct {
		name             string
		result           interface{}
		err              error
		wantDone, wantFail bool
	}{
		{"gone 404", nil, notFound, true, false},
		{"still deleting", map[string]interface{}{"status": "DELETING"}, nil, false, false},
		{"came back active", map[string]interface{}{"status": "ACTIVE"}, nil, false, true},
		{"non-404 err transient", nil, otherAPIErr, false, false},
		{"plain transient err", nil, errors.New("boom"), false, false},
	}
	for _, tc := range cases {
		done, failed, _ := evaluateDeleted(tc.result, tc.err)
		if done != tc.wantDone || failed != tc.wantFail {
			t.Errorf("%s: done=%v failed=%v, want done=%v failed=%v", tc.name, done, failed, tc.wantDone, tc.wantFail)
		}
	}
}
