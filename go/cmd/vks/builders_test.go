package vks

import (
	"testing"
)

func TestBuildUpgradeNodegroupBody(t *testing.T) {
	got := buildUpgradeNodegroupBody("v1.29.0")
	if got["kubernetesVersion"] != "v1.29.0" {
		t.Errorf("body = %#v, want kubernetesVersion=v1.29.0", got)
	}
	if len(got) != 1 {
		t.Errorf("body has %d keys, want 1", len(got))
	}
}

func TestBuildAutoHealingBodyOnlyChangedOptionalFields(t *testing.T) {
	got := buildAutoHealingBody(true, "30%", "", 600, map[string]bool{
		"max-unhealthy": true, "unhealthy-range": false, "timeout-unhealthy": true,
	})
	if got["enableAutoHealing"] != true {
		t.Errorf("enableAutoHealing = %#v, want true", got["enableAutoHealing"])
	}
	if got["maxUnhealthy"] != "30%" {
		t.Errorf("maxUnhealthy = %#v, want 30%%", got["maxUnhealthy"])
	}
	if got["timeoutUnhealthy"] != 600 {
		t.Errorf("timeoutUnhealthy = %#v, want 600", got["timeoutUnhealthy"])
	}
	if _, ok := got["unhealthyRange"]; ok {
		t.Errorf("unhealthyRange should be absent when flag not set; got %#v", got)
	}
}

func TestBuildMetadataBodyIncludesOnlyChangedKeys(t *testing.T) {
	got := buildMetadataBody("env=prod", "", "dedicated=gpu:NoSchedule", map[string]bool{
		"labels": true, "tags": false, "taints": true,
	})
	labels, ok := got["labels"].(map[string]string)
	if !ok || labels["env"] != "prod" {
		t.Errorf("labels = %#v, want env=prod", got["labels"])
	}
	if _, ok := got["tags"]; ok {
		t.Errorf("tags should be absent when flag not set; got %#v", got)
	}
	taints, ok := got["taints"].([]Taint)
	if !ok || len(taints) != 1 || taints[0].Key != "dedicated" || taints[0].Effect != "NoSchedule" {
		t.Errorf("taints = %#v, want one dedicated=gpu:NoSchedule taint", got["taints"])
	}
}
