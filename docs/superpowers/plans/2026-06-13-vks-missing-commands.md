# VKS Missing Commands Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Add 8 new `grn vks` commands (plus a `Patch` client method, a YAML-based kubeconfig package, docs, and changelog) so the CLI covers the kubeconfig, nodegroup version upgrade, auto-healing, metadata, events, and cluster-versions VKS API endpoints.

**Architecture:** Each command is a cobra `cobra.Command` in `go/cmd/vks/<name>.go`, registered in `go/cmd/vks/vks.go`, reusing the existing `createClient` / `outputResult` / `validator.ValidateID` helpers. Request-body and query-param construction is factored into **pure functions** so they can be unit-tested with `go test` without invoking `os.Exit`. The kubeconfig merge logic lives in a new `internal/kubeconfig` package and is fully unit-tested.

**Tech Stack:** Go 1.22, cobra, `gopkg.in/yaml.v3` (new dependency), `go test` + `net/http/httptest`.

---

## Conventions (from CLAUDE.md — follow exactly)

- All code/comments in **English**.
- VKS pagination is **0-based** (page 0 = first page).
- Use `--k8s-version` (NOT `--version`, which conflicts with the global version flag).
- IDs used in URLs must pass `validator.ValidateID()` first.
- **Do NOT auto commit/push.** Each task ends with a build/test verification step; the user commits manually at the end. Commit commands are intentionally omitted.
- Every change adds a changelog fragment via `./scripts/new-change` (Task 12).
- New commands require docs updates (Task 12).
- Build command: `cd go && CGO_ENABLED=0 go build -o /tmp/grn .`
- Test command: `cd go && go test ./...`

---

## File Structure

**New files:**
- `go/internal/kubeconfig/kubeconfig.go` — load/merge/write kubeconfig (pure, YAML).
- `go/internal/kubeconfig/kubeconfig_test.go` — merge unit tests.
- `go/cmd/vks/list_cluster_versions.go`
- `go/cmd/vks/upgrade_nodegroup_version.go`
- `go/cmd/vks/config_auto_healing.go`
- `go/cmd/vks/update_nodegroup_metadata.go`
- `go/cmd/vks/get_cluster_events.go`
- `go/cmd/vks/get_nodegroup_events.go`
- `go/cmd/vks/generate_kubeconfig.go`
- `go/cmd/vks/update_kubeconfig.go`
- `go/cmd/vks/builders_test.go` — unit tests for pure body/query builders.
- `docs/commands/vks/<each-command>.md` — 8 reference pages.

**Modified files:**
- `go/internal/client/client.go` — add `Patch`.
- `go/internal/client/client_test.go` — new, tests `Patch`.
- `go/cmd/vks/helpers.go` — add `buildEventsQuery` helper.
- `go/cmd/vks/vks.go` — register 8 commands.
- `go/go.mod`, `go/go.sum` — add `gopkg.in/yaml.v3`.
- `mkdocs.yml`, `docs/commands/vks/index.md`, `README.md`, `CLAUDE.md` — docs.

---

## Task 1: Add `Patch` method to the HTTP client

**Files:**
- Modify: `go/internal/client/client.go` (near the existing `Put` method, ~line 80)
- Test: `go/internal/client/client_test.go` (create)

- [ ] **Step 1: Write the failing test**

Create `go/internal/client/client_test.go`:

```go
package client

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/vngcloud/greennode-cli/internal/auth"
)

func TestPatchSendsPatchMethodAndBody(t *testing.T) {
	var gotMethod, gotBody string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotMethod = r.Method
		b, _ := io.ReadAll(r.Body)
		gotBody = string(b)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"ok":true}`))
	}))
	defer srv.Close()

	tm := auth.NewTokenManager("id", "secret")
	c := NewGreenodeClient(srv.URL, tm, 5*time.Second, false, false)

	_, err := c.Patch("/v1/thing", map[string]interface{}{"enableAutoHealing": true})
	if err != nil {
		t.Fatalf("Patch returned error: %v", err)
	}
	if gotMethod != http.MethodPatch {
		t.Errorf("method = %q, want PATCH", gotMethod)
	}
	if gotBody != `{"enableAutoHealing":true}` {
		t.Errorf("body = %q, want enableAutoHealing payload", gotBody)
	}
}
```

> Note: `NewGreenodeClient(baseURL, tokenManager, timeout, verifySSL, debug)`. If the token manager performs a network call to fetch a token against `srv.URL`, and the test fails on auth, adjust by pointing the token endpoint at the same test server or skip auth — inspect `internal/auth/token.go` first and mirror however other client behavior is exercised. If auth cannot be stubbed simply, replace this test with a direct test of the method string via a smaller seam, but do NOT delete the Patch coverage.

- [ ] **Step 2: Run test to verify it fails**

Run: `cd go && go test ./internal/client/ -run TestPatch -v`
Expected: FAIL — `c.Patch undefined (type *GreenodeClient has no field or method Patch)`.

- [ ] **Step 3: Add the `Patch` method**

In `go/internal/client/client.go`, immediately after the `Put` method:

```go
// Patch performs a PATCH request with a JSON body.
func (c *GreenodeClient) Patch(path string, body interface{}) (interface{}, error) {
	return c.request("PATCH", path, nil, body)
}
```

- [ ] **Step 4: Run test to verify it passes**

Run: `cd go && go test ./internal/client/ -run TestPatch -v`
Expected: PASS.

- [ ] **Step 5: Build**

Run: `cd go && CGO_ENABLED=0 go build -o /tmp/grn .`
Expected: builds with no errors.

---

## Task 2: Add `buildEventsQuery` pagination helper

The two events commands share query construction. VKS pagination is 0-based, so we do NOT default page to 1.

**Files:**
- Modify: `go/cmd/vks/helpers.go` (add function at end)
- Test: `go/cmd/vks/builders_test.go` (create)

- [ ] **Step 1: Write the failing test**

Create `go/cmd/vks/builders_test.go`:

```go
package vks

import (
	"reflect"
	"testing"
)

func TestBuildEventsQueryOnlyIncludesSetValues(t *testing.T) {
	got := buildEventsQuery("CREATE", "", 2, 50, map[string]bool{
		"action": true, "type": false, "page": true, "page-size": true,
	})
	want := map[string]string{
		"action":   "CREATE",
		"page":     "2",
		"pageSize": "50",
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("buildEventsQuery = %#v, want %#v", got, want)
	}
}

func TestBuildEventsQueryEmptyWhenNothingSet(t *testing.T) {
	got := buildEventsQuery("", "", 0, 0, map[string]bool{})
	if len(got) != 0 {
		t.Errorf("buildEventsQuery = %#v, want empty map", got)
	}
}
```

- [ ] **Step 2: Run test to verify it fails**

Run: `cd go && go test ./cmd/vks/ -run TestBuildEventsQuery -v`
Expected: FAIL — `undefined: buildEventsQuery`.

- [ ] **Step 3: Implement the helper**

Append to `go/cmd/vks/helpers.go` (it already imports `"fmt"`):

```go
// buildEventsQuery builds query params for events endpoints, including only
// flags the user explicitly set. `changed` maps flag name -> was it set.
// VKS pagination is 0-based, so page is passed through verbatim.
func buildEventsQuery(action, eventType string, page, pageSize int, changed map[string]bool) map[string]string {
	params := map[string]string{}
	if changed["action"] && action != "" {
		params["action"] = action
	}
	if changed["type"] && eventType != "" {
		params["type"] = eventType
	}
	if changed["page"] {
		params["page"] = fmt.Sprintf("%d", page)
	}
	if changed["page-size"] {
		params["pageSize"] = fmt.Sprintf("%d", pageSize)
	}
	return params
}
```

- [ ] **Step 4: Run test to verify it passes**

Run: `cd go && go test ./cmd/vks/ -run TestBuildEventsQuery -v`
Expected: PASS.

---

## Task 3: `list-cluster-versions` command

Simplest command — GET, no params. Establishes the file pattern.

**Files:**
- Create: `go/cmd/vks/list_cluster_versions.go`
- Modify: `go/cmd/vks/vks.go`

- [ ] **Step 1: Create the command file**

`go/cmd/vks/list_cluster_versions.go`:

```go
package vks

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var listClusterVersionsCmd = &cobra.Command{
	Use:   "list-cluster-versions",
	Short: "List available Kubernetes versions for VKS clusters",
	RunE:  runListClusterVersions,
}

func runListClusterVersions(cmd *cobra.Command, args []string) error {
	apiClient, err := createClient(cmd)
	if err != nil {
		return err
	}

	result, err := apiClient.Get("/v1/cluster-versions", nil)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	return outputResult(cmd, result)
}
```

- [ ] **Step 2: Register in vks.go**

In `go/cmd/vks/vks.go` `init()`, add a new section before the closing brace:

```go
	// Version & event commands
	VksCmd.AddCommand(listClusterVersionsCmd)
```

- [ ] **Step 3: Build**

Run: `cd go && CGO_ENABLED=0 go build -o /tmp/grn .`
Expected: builds clean.

- [ ] **Step 4: Verify the command is wired**

Run: `/tmp/grn vks list-cluster-versions --help`
Expected: help text shows the command and the global flags (`--output`, `--region`, etc.).

---

## Task 4: `upgrade-nodegroup-version` command

POST `/v1/clusters/{c}/node-groups/{ng}/upgrade-version`, body `{"kubernetesVersion": <v>}`. Uses `--k8s-version` per CLAUDE.md.

**Files:**
- Create: `go/cmd/vks/upgrade_nodegroup_version.go`
- Modify: `go/cmd/vks/vks.go`
- Test: `go/cmd/vks/builders_test.go`

- [ ] **Step 1: Write the failing test for the body builder**

Append to `go/cmd/vks/builders_test.go`:

```go
func TestBuildUpgradeNodegroupBody(t *testing.T) {
	got := buildUpgradeNodegroupBody("v1.29.0")
	if got["kubernetesVersion"] != "v1.29.0" {
		t.Errorf("body = %#v, want kubernetesVersion=v1.29.0", got)
	}
	if len(got) != 1 {
		t.Errorf("body has %d keys, want 1", len(got))
	}
}
```

- [ ] **Step 2: Run test to verify it fails**

Run: `cd go && go test ./cmd/vks/ -run TestBuildUpgradeNodegroupBody -v`
Expected: FAIL — `undefined: buildUpgradeNodegroupBody`.

- [ ] **Step 3: Create the command file**

`go/cmd/vks/upgrade_nodegroup_version.go`:

```go
package vks

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/vngcloud/greennode-cli/internal/validator"
)

var upgradeNodegroupVersionCmd = &cobra.Command{
	Use:   "upgrade-nodegroup-version",
	Short: "Upgrade the Kubernetes version of a node group",
	RunE:  runUpgradeNodegroupVersion,
}

func init() {
	f := upgradeNodegroupVersionCmd.Flags()
	f.String("cluster-id", "", "Cluster ID (required)")
	f.String("nodegroup-id", "", "Node group ID (required)")
	f.String("k8s-version", "", "Target Kubernetes version (required)")

	upgradeNodegroupVersionCmd.MarkFlagRequired("cluster-id")
	upgradeNodegroupVersionCmd.MarkFlagRequired("nodegroup-id")
	upgradeNodegroupVersionCmd.MarkFlagRequired("k8s-version")
}

func buildUpgradeNodegroupBody(k8sVersion string) map[string]interface{} {
	return map[string]interface{}{"kubernetesVersion": k8sVersion}
}

func runUpgradeNodegroupVersion(cmd *cobra.Command, args []string) error {
	clusterID, _ := cmd.Flags().GetString("cluster-id")
	nodegroupID, _ := cmd.Flags().GetString("nodegroup-id")
	k8sVersion, _ := cmd.Flags().GetString("k8s-version")

	if err := validator.ValidateID(clusterID, "cluster-id"); err != nil {
		return err
	}
	if err := validator.ValidateID(nodegroupID, "nodegroup-id"); err != nil {
		return err
	}

	body := buildUpgradeNodegroupBody(k8sVersion)

	apiClient, err := createClient(cmd)
	if err != nil {
		return err
	}

	result, err := apiClient.Post(
		fmt.Sprintf("/v1/clusters/%s/node-groups/%s/upgrade-version", clusterID, nodegroupID), body,
	)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	return outputResult(cmd, result)
}
```

- [ ] **Step 4: Register in vks.go**

In the "Version & event commands" section of `vks.go` `init()`:

```go
	VksCmd.AddCommand(upgradeNodegroupVersionCmd)
```

- [ ] **Step 5: Run test + build**

Run: `cd go && go test ./cmd/vks/ -run TestBuildUpgradeNodegroupBody -v && CGO_ENABLED=0 go build -o /tmp/grn .`
Expected: test PASS, build clean.

- [ ] **Step 6: Verify wiring**

Run: `/tmp/grn vks upgrade-nodegroup-version --help`
Expected: shows `--cluster-id`, `--nodegroup-id`, `--k8s-version` as required.

---

## Task 5: `config-auto-healing` command

PATCH `/v1/clusters/{c}/auto-healing-config`. `enableAutoHealing` always sent; other fields only if the flag was changed.

**Files:**
- Create: `go/cmd/vks/config_auto_healing.go`
- Modify: `go/cmd/vks/vks.go`
- Test: `go/cmd/vks/builders_test.go`

- [ ] **Step 1: Write the failing test**

Append to `go/cmd/vks/builders_test.go`:

```go
func TestBuildAutoHealingBodyOnlyChangedOptionalFields(t *testing.T) {
	got := buildAutoHealingBody(true, "30%", "", 600, map[string]bool{
		"max-unhealthy": true, "unhealthy-range": false, "timeout-unhealthy": true,
	})
	want := map[string]interface{}{
		"enableAutoHealing": true,
		"maxUnhealthy":      "30%",
		"timeoutUnhealthy":  600,
	}
	if got["enableAutoHealing"] != want["enableAutoHealing"] ||
		got["maxUnhealthy"] != want["maxUnhealthy"] ||
		got["timeoutUnhealthy"] != want["timeoutUnhealthy"] {
		t.Errorf("body = %#v, want %#v", got, want)
	}
	if _, ok := got["unhealthyRange"]; ok {
		t.Errorf("unhealthyRange should be absent when flag not set; got %#v", got)
	}
}
```

- [ ] **Step 2: Run test to verify it fails**

Run: `cd go && go test ./cmd/vks/ -run TestBuildAutoHealingBody -v`
Expected: FAIL — `undefined: buildAutoHealingBody`.

- [ ] **Step 3: Create the command file**

`go/cmd/vks/config_auto_healing.go`:

```go
package vks

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/vngcloud/greennode-cli/internal/validator"
)

var configAutoHealingCmd = &cobra.Command{
	Use:   "config-auto-healing",
	Short: "Configure auto-healing for a VKS cluster",
	RunE:  runConfigAutoHealing,
}

func init() {
	f := configAutoHealingCmd.Flags()
	f.String("cluster-id", "", "Cluster ID (required)")
	f.Bool("enable-auto-healing", false, "Enable auto-healing (required)")
	f.String("max-unhealthy", "", "Max unhealthy nodes, e.g. \"30%\"")
	f.String("unhealthy-range", "", "Unhealthy range")
	f.Int("timeout-unhealthy", 0, "Unhealthy timeout in seconds")

	configAutoHealingCmd.MarkFlagRequired("cluster-id")
	configAutoHealingCmd.MarkFlagRequired("enable-auto-healing")
}

func buildAutoHealingBody(enable bool, maxUnhealthy, unhealthyRange string, timeoutUnhealthy int, changed map[string]bool) map[string]interface{} {
	body := map[string]interface{}{"enableAutoHealing": enable}
	if changed["max-unhealthy"] && maxUnhealthy != "" {
		body["maxUnhealthy"] = maxUnhealthy
	}
	if changed["unhealthy-range"] && unhealthyRange != "" {
		body["unhealthyRange"] = unhealthyRange
	}
	if changed["timeout-unhealthy"] {
		body["timeoutUnhealthy"] = timeoutUnhealthy
	}
	return body
}

func runConfigAutoHealing(cmd *cobra.Command, args []string) error {
	clusterID, _ := cmd.Flags().GetString("cluster-id")
	enable, _ := cmd.Flags().GetBool("enable-auto-healing")
	maxUnhealthy, _ := cmd.Flags().GetString("max-unhealthy")
	unhealthyRange, _ := cmd.Flags().GetString("unhealthy-range")
	timeoutUnhealthy, _ := cmd.Flags().GetInt("timeout-unhealthy")

	if err := validator.ValidateID(clusterID, "cluster-id"); err != nil {
		return err
	}

	changed := map[string]bool{
		"max-unhealthy":     cmd.Flags().Changed("max-unhealthy"),
		"unhealthy-range":   cmd.Flags().Changed("unhealthy-range"),
		"timeout-unhealthy": cmd.Flags().Changed("timeout-unhealthy"),
	}
	body := buildAutoHealingBody(enable, maxUnhealthy, unhealthyRange, timeoutUnhealthy, changed)

	apiClient, err := createClient(cmd)
	if err != nil {
		return err
	}

	result, err := apiClient.Patch(
		fmt.Sprintf("/v1/clusters/%s/auto-healing-config", clusterID), body,
	)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	return outputResult(cmd, result)
}
```

- [ ] **Step 4: Register in vks.go**

Add a "Config commands" line near the auto-upgrade section of `vks.go` `init()`:

```go
	VksCmd.AddCommand(configAutoHealingCmd)
```

- [ ] **Step 5: Run test + build**

Run: `cd go && go test ./cmd/vks/ -run TestBuildAutoHealingBody -v && CGO_ENABLED=0 go build -o /tmp/grn .`
Expected: test PASS, build clean.

- [ ] **Step 6: Verify wiring**

Run: `/tmp/grn vks config-auto-healing --help`
Expected: shows `--cluster-id`, `--enable-auto-healing` required; optional tuning flags.

---

## Task 6: `update-nodegroup-metadata` command

PATCH `/v1/clusters/{c}/node-groups/{ng}/metadata`. Body `{labels?, tags?, taints?}`, only the keys whose flags were set. Reuses `parseLabels` and `parseTaints`.

**Files:**
- Create: `go/cmd/vks/update_nodegroup_metadata.go`
- Modify: `go/cmd/vks/vks.go`
- Test: `go/cmd/vks/builders_test.go`

- [ ] **Step 1: Write the failing test**

Append to `go/cmd/vks/builders_test.go`:

```go
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
```

- [ ] **Step 2: Run test to verify it fails**

Run: `cd go && go test ./cmd/vks/ -run TestBuildMetadataBody -v`
Expected: FAIL — `undefined: buildMetadataBody`.

- [ ] **Step 3: Create the command file**

`go/cmd/vks/update_nodegroup_metadata.go`:

```go
package vks

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/vngcloud/greennode-cli/internal/validator"
)

var updateNodegroupMetadataCmd = &cobra.Command{
	Use:   "update-nodegroup-metadata",
	Short: "Update labels, tags, and taints of a node group",
	RunE:  runUpdateNodegroupMetadata,
}

func init() {
	f := updateNodegroupMetadataCmd.Flags()
	f.String("cluster-id", "", "Cluster ID (required)")
	f.String("nodegroup-id", "", "Node group ID (required)")
	f.String("labels", "", "Node labels as key=value pairs (comma-separated)")
	f.String("tags", "", "Tags as key=value pairs (comma-separated)")
	f.String("taints", "", "Node taints as key=value:effect (comma-separated)")

	updateNodegroupMetadataCmd.MarkFlagRequired("cluster-id")
	updateNodegroupMetadataCmd.MarkFlagRequired("nodegroup-id")
}

func buildMetadataBody(labelsStr, tagsStr, taintsStr string, changed map[string]bool) map[string]interface{} {
	body := map[string]interface{}{}
	if changed["labels"] {
		body["labels"] = parseLabels(labelsStr)
	}
	if changed["tags"] {
		body["tags"] = parseLabels(tagsStr)
	}
	if changed["taints"] {
		body["taints"] = parseTaints(taintsStr)
	}
	return body
}

func runUpdateNodegroupMetadata(cmd *cobra.Command, args []string) error {
	clusterID, _ := cmd.Flags().GetString("cluster-id")
	nodegroupID, _ := cmd.Flags().GetString("nodegroup-id")
	labelsStr, _ := cmd.Flags().GetString("labels")
	tagsStr, _ := cmd.Flags().GetString("tags")
	taintsStr, _ := cmd.Flags().GetString("taints")

	if err := validator.ValidateID(clusterID, "cluster-id"); err != nil {
		return err
	}
	if err := validator.ValidateID(nodegroupID, "nodegroup-id"); err != nil {
		return err
	}

	changed := map[string]bool{
		"labels": cmd.Flags().Changed("labels"),
		"tags":   cmd.Flags().Changed("tags"),
		"taints": cmd.Flags().Changed("taints"),
	}
	if !changed["labels"] && !changed["tags"] && !changed["taints"] {
		return fmt.Errorf("at least one of --labels, --tags, --taints must be provided")
	}
	body := buildMetadataBody(labelsStr, tagsStr, taintsStr, changed)

	apiClient, err := createClient(cmd)
	if err != nil {
		return err
	}

	result, err := apiClient.Patch(
		fmt.Sprintf("/v1/clusters/%s/node-groups/%s/metadata", clusterID, nodegroupID), body,
	)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	return outputResult(cmd, result)
}
```

- [ ] **Step 4: Register in vks.go**

In the nodegroup section of `vks.go` `init()`:

```go
	VksCmd.AddCommand(updateNodegroupMetadataCmd)
```

- [ ] **Step 5: Run test + build**

Run: `cd go && go test ./cmd/vks/ -run TestBuildMetadataBody -v && CGO_ENABLED=0 go build -o /tmp/grn .`
Expected: test PASS, build clean.

- [ ] **Step 6: Verify wiring**

Run: `/tmp/grn vks update-nodegroup-metadata --help`
Expected: shows the five flags.

---

## Task 7: `get-cluster-events` command

GET `/v1/clusters/{c}/events` with query params via `buildEventsQuery`.

**Files:**
- Create: `go/cmd/vks/get_cluster_events.go`
- Modify: `go/cmd/vks/vks.go`

- [ ] **Step 1: Create the command file**

`go/cmd/vks/get_cluster_events.go`:

```go
package vks

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/vngcloud/greennode-cli/internal/validator"
)

var getClusterEventsCmd = &cobra.Command{
	Use:   "get-cluster-events",
	Short: "Get the list of events for a VKS cluster",
	RunE:  runGetClusterEvents,
}

func init() {
	f := getClusterEventsCmd.Flags()
	f.String("cluster-id", "", "Cluster ID (required)")
	f.String("action", "", "Filter by action")
	f.String("type", "", "Filter by event type")
	f.Int("page", 0, "Page number (0-based)")
	f.Int("page-size", 50, "Page size")

	getClusterEventsCmd.MarkFlagRequired("cluster-id")
}

func runGetClusterEvents(cmd *cobra.Command, args []string) error {
	clusterID, _ := cmd.Flags().GetString("cluster-id")
	action, _ := cmd.Flags().GetString("action")
	eventType, _ := cmd.Flags().GetString("type")
	page, _ := cmd.Flags().GetInt("page")
	pageSize, _ := cmd.Flags().GetInt("page-size")

	if err := validator.ValidateID(clusterID, "cluster-id"); err != nil {
		return err
	}

	changed := map[string]bool{
		"action":    cmd.Flags().Changed("action"),
		"type":      cmd.Flags().Changed("type"),
		"page":      cmd.Flags().Changed("page"),
		"page-size": cmd.Flags().Changed("page-size"),
	}
	params := buildEventsQuery(action, eventType, page, pageSize, changed)

	apiClient, err := createClient(cmd)
	if err != nil {
		return err
	}

	result, err := apiClient.Get(
		fmt.Sprintf("/v1/clusters/%s/events", clusterID), params,
	)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	return outputResult(cmd, result)
}
```

- [ ] **Step 2: Register in vks.go**

In the "Version & event commands" section:

```go
	VksCmd.AddCommand(getClusterEventsCmd)
```

- [ ] **Step 3: Build + verify**

Run: `cd go && CGO_ENABLED=0 go build -o /tmp/grn . && /tmp/grn vks get-cluster-events --help`
Expected: build clean; help shows `--cluster-id` required plus filter/paging flags.

---

## Task 8: `get-nodegroup-events` command

GET `/v1/clusters/{c}/node-groups/{ng}/events`. Same shape as Task 7 plus `--nodegroup-id`.

**Files:**
- Create: `go/cmd/vks/get_nodegroup_events.go`
- Modify: `go/cmd/vks/vks.go`

- [ ] **Step 1: Create the command file**

`go/cmd/vks/get_nodegroup_events.go`:

```go
package vks

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/vngcloud/greennode-cli/internal/validator"
)

var getNodegroupEventsCmd = &cobra.Command{
	Use:   "get-nodegroup-events",
	Short: "Get the list of events for a node group",
	RunE:  runGetNodegroupEvents,
}

func init() {
	f := getNodegroupEventsCmd.Flags()
	f.String("cluster-id", "", "Cluster ID (required)")
	f.String("nodegroup-id", "", "Node group ID (required)")
	f.String("action", "", "Filter by action")
	f.String("type", "", "Filter by event type")
	f.Int("page", 0, "Page number (0-based)")
	f.Int("page-size", 50, "Page size")

	getNodegroupEventsCmd.MarkFlagRequired("cluster-id")
	getNodegroupEventsCmd.MarkFlagRequired("nodegroup-id")
}

func runGetNodegroupEvents(cmd *cobra.Command, args []string) error {
	clusterID, _ := cmd.Flags().GetString("cluster-id")
	nodegroupID, _ := cmd.Flags().GetString("nodegroup-id")
	action, _ := cmd.Flags().GetString("action")
	eventType, _ := cmd.Flags().GetString("type")
	page, _ := cmd.Flags().GetInt("page")
	pageSize, _ := cmd.Flags().GetInt("page-size")

	if err := validator.ValidateID(clusterID, "cluster-id"); err != nil {
		return err
	}
	if err := validator.ValidateID(nodegroupID, "nodegroup-id"); err != nil {
		return err
	}

	changed := map[string]bool{
		"action":    cmd.Flags().Changed("action"),
		"type":      cmd.Flags().Changed("type"),
		"page":      cmd.Flags().Changed("page"),
		"page-size": cmd.Flags().Changed("page-size"),
	}
	params := buildEventsQuery(action, eventType, page, pageSize, changed)

	apiClient, err := createClient(cmd)
	if err != nil {
		return err
	}

	result, err := apiClient.Get(
		fmt.Sprintf("/v1/clusters/%s/node-groups/%s/events", clusterID, nodegroupID), params,
	)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	return outputResult(cmd, result)
}
```

- [ ] **Step 2: Register in vks.go**

In the "Version & event commands" section:

```go
	VksCmd.AddCommand(getNodegroupEventsCmd)
```

- [ ] **Step 3: Build + verify**

Run: `cd go && CGO_ENABLED=0 go build -o /tmp/grn . && /tmp/grn vks get-nodegroup-events --help`
Expected: build clean; help shows both IDs required.

---

## Task 9: `internal/kubeconfig` package (load / merge / write)

Pure YAML logic, fully unit-tested. This is the core of `update-kubeconfig`.

**Files:**
- Modify: `go/go.mod`, `go/go.sum` (add dependency)
- Create: `go/internal/kubeconfig/kubeconfig.go`
- Test: `go/internal/kubeconfig/kubeconfig_test.go`

- [ ] **Step 1: Add the YAML dependency**

Run: `cd go && go get gopkg.in/yaml.v3@v3.0.1`
Expected: `go.mod` now lists `gopkg.in/yaml.v3 v3.0.1`.

- [ ] **Step 2: Write the failing tests**

Create `go/internal/kubeconfig/kubeconfig_test.go`:

```go
package kubeconfig

import (
	"os"
	"path/filepath"
	"testing"
)

const incomingKubeconfig = `apiVersion: v1
kind: Config
clusters:
- name: vks-cluster
  cluster:
    server: https://1.2.3.4:6443
    certificate-authority-data: AAAA
contexts:
- name: vks-cluster-ctx
  context:
    cluster: vks-cluster
    user: vks-user
current-context: vks-cluster-ctx
users:
- name: vks-user
  user:
    token: secret-token
`

func TestMergeIntoEmptyFileCreatesContext(t *testing.T) {
	dir := t.TempDir()
	target := filepath.Join(dir, "config")

	res, err := Merge(target, incomingKubeconfig, "vks_c-123", true)
	if err != nil {
		t.Fatalf("Merge error: %v", err)
	}
	if res.ContextName != "vks_c-123" {
		t.Errorf("ContextName = %q, want vks_c-123", res.ContextName)
	}

	cfg, err := Load(target)
	if err != nil {
		t.Fatalf("Load error: %v", err)
	}
	if cfg.CurrentContext != "vks_c-123" {
		t.Errorf("current-context = %q, want vks_c-123", cfg.CurrentContext)
	}
	if findContext(cfg, "vks_c-123") == nil {
		t.Errorf("merged context not found in %#v", cfg.Contexts)
	}
}

func TestMergePreservesExistingContexts(t *testing.T) {
	dir := t.TempDir()
	target := filepath.Join(dir, "config")
	existing := `apiVersion: v1
kind: Config
clusters:
- name: other
  cluster:
    server: https://9.9.9.9
contexts:
- name: other-ctx
  context:
    cluster: other
    user: other-user
current-context: other-ctx
users:
- name: other-user
  user:
    token: other
`
	if err := os.WriteFile(target, []byte(existing), 0o600); err != nil {
		t.Fatal(err)
	}

	if _, err := Merge(target, incomingKubeconfig, "vks_c-123", false); err != nil {
		t.Fatalf("Merge error: %v", err)
	}

	cfg, _ := Load(target)
	if findContext(cfg, "other-ctx") == nil {
		t.Errorf("existing context was dropped")
	}
	if findContext(cfg, "vks_c-123") == nil {
		t.Errorf("new context missing")
	}
	if cfg.CurrentContext != "other-ctx" {
		t.Errorf("current-context = %q, want other-ctx (setCurrent=false)", cfg.CurrentContext)
	}
}

func TestMergeOverwritesSameContextName(t *testing.T) {
	dir := t.TempDir()
	target := filepath.Join(dir, "config")

	if _, err := Merge(target, incomingKubeconfig, "vks_c-123", true); err != nil {
		t.Fatal(err)
	}
	if _, err := Merge(target, incomingKubeconfig, "vks_c-123", true); err != nil {
		t.Fatal(err)
	}

	cfg, _ := Load(target)
	count := 0
	for _, c := range cfg.Contexts {
		if c.Name == "vks_c-123" {
			count++
		}
	}
	if count != 1 {
		t.Errorf("context vks_c-123 appears %d times, want 1", count)
	}
}
```

- [ ] **Step 3: Run tests to verify they fail**

Run: `cd go && go test ./internal/kubeconfig/ -v`
Expected: FAIL — `undefined: Merge`, `Load`, `findContext`, etc.

- [ ] **Step 4: Implement the package**

Create `go/internal/kubeconfig/kubeconfig.go`:

```go
// Package kubeconfig loads, merges, and writes Kubernetes kubeconfig files.
package kubeconfig

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// Config models the subset of a kubeconfig we manipulate. Unknown fields are
// preserved via yaml.Node passthrough on the entry "cluster"/"context"/"user".
type Config struct {
	APIVersion     string         `yaml:"apiVersion,omitempty"`
	Kind           string         `yaml:"kind,omitempty"`
	Clusters       []NamedEntry   `yaml:"clusters"`
	Contexts       []NamedEntry   `yaml:"contexts"`
	Users          []NamedEntry   `yaml:"users"`
	CurrentContext string         `yaml:"current-context,omitempty"`
	Preferences    map[string]any `yaml:"preferences,omitempty"`
}

// NamedEntry is a generic {name, <payload>} entry. The payload key differs by
// list (cluster/context/user); we keep it as a raw node to avoid lossy parsing.
type NamedEntry struct {
	Name    string         `yaml:"name"`
	Cluster *yaml.Node     `yaml:"cluster,omitempty"`
	Context *contextBody   `yaml:"context,omitempty"`
	User    *yaml.Node     `yaml:"user,omitempty"`
}

type contextBody struct {
	Cluster string `yaml:"cluster"`
	User    string `yaml:"user"`
}

// MergeResult reports what was applied.
type MergeResult struct {
	ContextName string
	Path        string
}

// Load reads a kubeconfig from disk. A missing file yields an empty Config.
func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return &Config{APIVersion: "v1", Kind: "Config"}, nil
	}
	if err != nil {
		return nil, err
	}
	var cfg Config
	if len(data) > 0 {
		if err := yaml.Unmarshal(data, &cfg); err != nil {
			return nil, fmt.Errorf("failed to parse kubeconfig %s: %w", path, err)
		}
	}
	if cfg.APIVersion == "" {
		cfg.APIVersion = "v1"
	}
	if cfg.Kind == "" {
		cfg.Kind = "Config"
	}
	return &cfg, nil
}

// Write serializes cfg to path with 0600 perms, creating parent dirs (0700).
func Write(path string, cfg *Config) error {
	if dir := filepath.Dir(path); dir != "" {
		if err := os.MkdirAll(dir, 0o700); err != nil {
			return err
		}
	}
	out, err := yaml.Marshal(cfg)
	if err != nil {
		return err
	}
	return os.WriteFile(path, out, 0o600)
}

// Merge parses incoming (a full kubeconfig YAML string) and grafts its first
// cluster/context/user into the file at target under the name contextName.
// If setCurrent is true, current-context is set to contextName.
func Merge(target, incoming, contextName string, setCurrent bool) (*MergeResult, error) {
	var src Config
	if err := yaml.Unmarshal([]byte(incoming), &src); err != nil {
		return nil, fmt.Errorf("failed to parse incoming kubeconfig: %w", err)
	}
	if len(src.Clusters) == 0 || len(src.Contexts) == 0 || len(src.Users) == 0 {
		return nil, fmt.Errorf("incoming kubeconfig is missing cluster/context/user entries")
	}

	dst, err := Load(target)
	if err != nil {
		return nil, err
	}

	clusterName := "cluster_" + contextName
	userName := "user_" + contextName

	cluster := src.Clusters[0]
	cluster.Name = clusterName
	user := src.Users[0]
	user.Name = userName
	ctx := NamedEntry{
		Name:    contextName,
		Context: &contextBody{Cluster: clusterName, User: userName},
	}

	dst.Clusters = upsert(dst.Clusters, cluster)
	dst.Users = upsert(dst.Users, user)
	dst.Contexts = upsert(dst.Contexts, ctx)
	if setCurrent {
		dst.CurrentContext = contextName
	}

	if err := Write(target, dst); err != nil {
		return nil, err
	}
	return &MergeResult{ContextName: contextName, Path: target}, nil
}

// upsert replaces an entry with the same name or appends it.
func upsert(list []NamedEntry, e NamedEntry) []NamedEntry {
	for i := range list {
		if list[i].Name == e.Name {
			list[i] = e
			return list
		}
	}
	return append(list, e)
}

// findContext is a test/helper accessor.
func findContext(cfg *Config, name string) *NamedEntry {
	for i := range cfg.Contexts {
		if cfg.Contexts[i].Name == name {
			return &cfg.Contexts[i]
		}
	}
	return nil
}
```

- [ ] **Step 5: Run tests to verify they pass**

Run: `cd go && go test ./internal/kubeconfig/ -v`
Expected: all three tests PASS.

> If YAML round-tripping of `cluster`/`user` payloads loses data because the source uses different field layouts, the `*yaml.Node` passthrough preserves them as-is — verify by printing the written file in `TestMergeIntoEmptyFileCreatesContext` if a test fails, and keep the node-based approach.

- [ ] **Step 6: Tidy modules + build**

Run: `cd go && go mod tidy && CGO_ENABLED=0 go build -o /tmp/grn .`
Expected: `go.sum` updated, build clean.

---

## Task 10: `generate-kubeconfig` command

POST `/v1/clusters/{c}/kubeconfig` with `{expirationDays}`. Returns 202 (async); inform the user.

**Files:**
- Create: `go/cmd/vks/generate_kubeconfig.go`
- Modify: `go/cmd/vks/vks.go`

- [ ] **Step 1: Create the command file**

`go/cmd/vks/generate_kubeconfig.go`:

```go
package vks

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/vngcloud/greennode-cli/internal/validator"
)

var generateKubeconfigCmd = &cobra.Command{
	Use:   "generate-kubeconfig",
	Short: "Request generation of a kubeconfig for a VKS cluster",
	Long: "Requests the VKS API to generate (or renew) a kubeconfig for the cluster. " +
		"This is asynchronous; once the kubeconfig becomes ACTIVE, run 'grn vks update-kubeconfig'.",
	RunE: runGenerateKubeconfig,
}

func init() {
	f := generateKubeconfigCmd.Flags()
	f.String("cluster-id", "", "Cluster ID (required)")
	f.Int("expiration-days", 30, "Number of days until the kubeconfig expires")

	generateKubeconfigCmd.MarkFlagRequired("cluster-id")
}

func runGenerateKubeconfig(cmd *cobra.Command, args []string) error {
	clusterID, _ := cmd.Flags().GetString("cluster-id")
	expirationDays, _ := cmd.Flags().GetInt("expiration-days")

	if err := validator.ValidateID(clusterID, "cluster-id"); err != nil {
		return err
	}

	body := map[string]interface{}{"expirationDays": expirationDays}

	apiClient, err := createClient(cmd)
	if err != nil {
		return err
	}

	_, err = apiClient.Post(
		fmt.Sprintf("/v1/clusters/%s/kubeconfig", clusterID), body,
	)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Kubeconfig generation requested for cluster %s (expires in %d days).\n", clusterID, expirationDays)
	fmt.Println("Generation is asynchronous. Once it is ACTIVE, run:")
	fmt.Printf("  grn vks update-kubeconfig --cluster-id %s\n", clusterID)
	return nil
}
```

- [ ] **Step 2: Register in vks.go**

Add a "Kubeconfig commands" section in `vks.go` `init()`:

```go
	// Kubeconfig commands
	VksCmd.AddCommand(generateKubeconfigCmd)
```

- [ ] **Step 3: Build + verify**

Run: `cd go && CGO_ENABLED=0 go build -o /tmp/grn . && /tmp/grn vks generate-kubeconfig --help`
Expected: build clean; help shows `--cluster-id` required, `--expiration-days` default 30.

---

## Task 11: `update-kubeconfig` command

GET `/v1/clusters/{c}/kubeconfig`, then merge into the target file via the `internal/kubeconfig` package.

**Files:**
- Create: `go/cmd/vks/update_kubeconfig.go`
- Modify: `go/cmd/vks/vks.go`

- [ ] **Step 1: Create the command file**

`go/cmd/vks/update_kubeconfig.go`:

```go
package vks

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"github.com/vngcloud/greennode-cli/internal/kubeconfig"
	"github.com/vngcloud/greennode-cli/internal/validator"
)

var updateKubeconfigCmd = &cobra.Command{
	Use:   "update-kubeconfig",
	Short: "Fetch the cluster kubeconfig and merge it into your kubeconfig file",
	RunE:  runUpdateKubeconfig,
}

func init() {
	f := updateKubeconfigCmd.Flags()
	f.String("cluster-id", "", "Cluster ID (required)")
	f.String("kubeconfig", "", "Path to kubeconfig file (default: $KUBECONFIG or ~/.kube/config)")
	f.String("alias", "", "Context name to use (default: vks_<cluster-id>)")
	f.Bool("no-set-context", false, "Do not set the merged context as current-context")
	f.Bool("dry-run", false, "Print what would be written without modifying the file")

	updateKubeconfigCmd.MarkFlagRequired("cluster-id")
}

// resolveKubeconfigPath picks the target path: explicit flag, then first entry
// of $KUBECONFIG, then ~/.kube/config.
func resolveKubeconfigPath(flagPath string) (string, error) {
	if flagPath != "" {
		return flagPath, nil
	}
	if env := os.Getenv("KUBECONFIG"); env != "" {
		return strings.Split(env, string(os.PathListSeparator))[0], nil
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".kube", "config"), nil
}

func runUpdateKubeconfig(cmd *cobra.Command, args []string) error {
	clusterID, _ := cmd.Flags().GetString("cluster-id")
	kubeconfigPath, _ := cmd.Flags().GetString("kubeconfig")
	alias, _ := cmd.Flags().GetString("alias")
	noSetContext, _ := cmd.Flags().GetBool("no-set-context")
	dryRun, _ := cmd.Flags().GetBool("dry-run")

	if err := validator.ValidateID(clusterID, "cluster-id"); err != nil {
		return err
	}

	contextName := alias
	if contextName == "" {
		contextName = "vks_" + clusterID
	}
	targetPath, err := resolveKubeconfigPath(kubeconfigPath)
	if err != nil {
		return err
	}

	apiClient, err := createClient(cmd)
	if err != nil {
		return err
	}

	result, err := apiClient.Get(fmt.Sprintf("/v1/clusters/%s/kubeconfig", clusterID), nil)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	resMap, ok := result.(map[string]interface{})
	if !ok {
		return fmt.Errorf("unexpected kubeconfig response format")
	}

	status, _ := resMap["status"].(string)
	switch status {
	case "NONE", "":
		return fmt.Errorf("no kubeconfig exists for cluster %s. Run 'grn vks generate-kubeconfig --cluster-id %s' first", clusterID, clusterID)
	case "CREATING":
		return fmt.Errorf("kubeconfig for cluster %s is still being generated; try again shortly", clusterID)
	case "ERROR":
		return fmt.Errorf("kubeconfig for cluster %s is in ERROR state; re-run 'grn vks generate-kubeconfig'", clusterID)
	}

	if warn, _ := resMap["renewalWarning"].(bool); warn {
		fmt.Fprintf(os.Stderr, "Warning: this kubeconfig is nearing expiration. Run 'grn vks generate-kubeconfig --cluster-id %s' to renew.\n", clusterID)
	}

	rawYAML, _ := resMap["kubeConfig"].(string)
	if rawYAML == "" {
		return fmt.Errorf("kubeconfig response did not contain kubeconfig data")
	}

	if dryRun {
		fmt.Printf("=== DRY RUN ===\n")
		fmt.Printf("Would merge context %q into %s\n", contextName, targetPath)
		if !noSetContext {
			fmt.Printf("Would set current-context to %q\n", contextName)
		}
		return nil
	}

	res, err := kubeconfig.Merge(targetPath, rawYAML, contextName, !noSetContext)
	if err != nil {
		return err
	}

	fmt.Printf("Updated context %q in %s\n", res.ContextName, res.Path)
	return nil
}
```

- [ ] **Step 2: Register in vks.go**

In the "Kubeconfig commands" section:

```go
	VksCmd.AddCommand(updateKubeconfigCmd)
```

- [ ] **Step 3: Build + verify**

Run: `cd go && CGO_ENABLED=0 go build -o /tmp/grn . && /tmp/grn vks update-kubeconfig --help`
Expected: build clean; help shows `--cluster-id`, `--kubeconfig`, `--alias`, `--no-set-context`, `--dry-run`.

---

## Task 12: Full test suite, docs, and changelog

**Files:**
- Create: `docs/commands/vks/{upgrade-nodegroup-version,config-auto-healing,update-nodegroup-metadata,get-cluster-events,get-nodegroup-events,list-cluster-versions,generate-kubeconfig,update-kubeconfig}.md`
- Modify: `docs/commands/vks/index.md`, `mkdocs.yml`, `README.md`, `CLAUDE.md`
- Modify: `.changes/next-release/` (via script)

- [ ] **Step 1: Run the full test suite + build**

Run: `cd go && go test ./... && CGO_ENABLED=0 go build -o /tmp/grn .`
Expected: all tests PASS, build clean.

- [ ] **Step 2: Verify all 8 commands appear**

Run: `/tmp/grn vks --help`
Expected: lists `config-auto-healing`, `generate-kubeconfig`, `get-cluster-events`, `get-nodegroup-events`, `list-cluster-versions`, `update-kubeconfig`, `update-nodegroup-metadata`, `upgrade-nodegroup-version` alongside existing commands.

- [ ] **Step 3: Write a docs page per command**

Create one Markdown file per command under `docs/commands/vks/`, following the existing page format (Description / Synopsis / Options / Examples) seen in `docs/commands/vks/delete-auto-upgrade-config.md`. Example for `upgrade-nodegroup-version.md`:

```markdown
# upgrade-nodegroup-version

## Description

Upgrade the Kubernetes version of a node group in a VKS cluster.

## Synopsis

```
grn vks upgrade-nodegroup-version
    --cluster-id <value>
    --nodegroup-id <value>
    --k8s-version <value>
```

## Options

`--cluster-id` (required)
: The ID of the cluster.

`--nodegroup-id` (required)
: The ID of the node group.

`--k8s-version` (required)
: The target Kubernetes version. Use `grn vks list-cluster-versions` to see valid versions.

## Examples

Upgrade a node group:

```bash
grn vks upgrade-nodegroup-version --cluster-id k8s-xxxxx --nodegroup-id ng-xxxxx --k8s-version v1.29.0
```
```

Write the remaining seven pages with their actual flags from Tasks 3–11 (no placeholders — list every flag with its real description).

- [ ] **Step 4: Update the docs index and nav**

In `docs/commands/vks/index.md`, add the 8 commands to the command table (match existing rows). In `mkdocs.yml`, add a nav entry for each new page under the VKS section (match existing entries).

- [ ] **Step 5: Update README.md and CLAUDE.md**

In `README.md`, add the new commands wherever VKS commands are listed. In `CLAUDE.md`, add the new files to the "Project structure" tree (the 8 command files + `internal/kubeconfig/`).

- [ ] **Step 6: Add changelog fragments**

Run (non-interactive):

```bash
cd /Users/lap16104/Documents/vks/greennode-cli
./scripts/new-change -t feature -c vks -d "Add kubeconfig commands (generate-kubeconfig, update-kubeconfig) for VKS clusters"
./scripts/new-change -t feature -c vks -d "Add upgrade-nodegroup-version, config-auto-healing, and update-nodegroup-metadata commands"
./scripts/new-change -t feature -c vks -d "Add list-cluster-versions, get-cluster-events, and get-nodegroup-events commands"
```

Expected: three JSON fragments created under `.changes/next-release/`.

- [ ] **Step 7: Final verification**

Run: `cd go && go vet ./... && go test ./... && CGO_ENABLED=0 go build -o /tmp/grn .`
Expected: clean vet, tests PASS, build clean.

> **Commit:** Per CLAUDE.md, do not auto-commit. Report completion and let the user commit/push (main is protected; work stays on the `feat/vks-missing-commands` branch).

---

## Self-Review Notes

- **Spec coverage:** All 8 commands from the design spec have tasks (Tasks 3–11), the `Patch` infra (Task 1), the YAML/kubeconfig package (Task 9), the events query helper (Task 2), and docs/changelog (Task 12). ✅
- **Naming:** `config-auto-healing` used consistently (command, file, builder). `--k8s-version` used for the version flag per CLAUDE.md (spec said `--kubernetes-version`; corrected to repo convention, body field remains `kubernetesVersion`).
- **Type consistency:** `buildEventsQuery`, `buildAutoHealingBody`, `buildMetadataBody`, `buildUpgradeNodegroupBody`, and `kubeconfig.Merge/Load/Write/Config/NamedEntry/MergeResult` signatures match between their definition tasks and their callers.
- **0-based pagination:** events `--page` defaults to 0 and is passed verbatim, per CLAUDE.md.
