# Design: bổ sung command còn thiếu cho `grn vks`

Ngày: 2026-06-13

## Mục tiêu

Bổ sung các command thiết yếu mà GreenNode CLI (`grn vks`) chưa có nhưng API VKS
(`context/api-docs/vks.json`, OpenAPI 3.0.3) đã hỗ trợ. Phạm vi được chốt dựa
trên đối chiếu trực tiếp với API thật và mức cần thiết khi vận hành, lấy cảm hứng
từ bộ command EKS của aws-cli (đặc biệt `update-kubeconfig`).

## Phạm vi (8 command)

| # | Command | HTTP | Endpoint |
|---|---|---|---|
| 1 | `upgrade-nodegroup-version` | POST | `/v1/clusters/{c}/node-groups/{ng}/upgrade-version` |
| 2 | `config-auto-healing` | PATCH | `/v1/clusters/{c}/auto-healing-config` |
| 3 | `update-nodegroup-metadata` | PATCH | `/v1/clusters/{c}/node-groups/{ng}/metadata` |
| 4 | `get-cluster-events` | GET | `/v1/clusters/{c}/events` |
| 5 | `get-nodegroup-events` | GET | `/v1/clusters/{c}/node-groups/{ng}/events` |
| 6 | `list-cluster-versions` | GET | `/v1/cluster-versions` |
| 7 | `generate-kubeconfig` | POST | `/v1/clusters/{c}/kubeconfig` |
| 8 | `update-kubeconfig` | GET | `/v1/clusters/{c}/kubeconfig` (+ merge file) |

### Ngoài phạm vi (đợt sau)
`list-nodegroup-images`, `list-nodes`, `get-upgrade-insight`, workspace
(get/create/reset-service-account), `register-fleet`/`unregister-fleet`,
`stop-poc`, `acknowledge-warning` đứng riêng (cảnh báo renewal đã được tích hợp
vào `update-kubeconfig`).

## Thay đổi hạ tầng chung

### a) Thêm method `Patch` vào client
`internal/client` — hàm `request()` đã nhận `method` string, chỉ cần wrapper:
```go
func (c *GreenodeClient) Patch(path string, body interface{}) (interface{}, error) {
    return c.request("PATCH", path, nil, body)
}
```

### b) Dependency YAML
Thêm `gopkg.in/yaml.v3` vào `go.mod` (pure Go) phục vụ merge kubeconfig.

### c) Helper query phân trang
Trong `cmd/vks/helpers.go`, thêm helper gom `action/type/page/page-size` thành
`map[string]string` (chỉ thêm key khi flag được `Changed`), tái dùng cho 2 lệnh
events.

## Năm command theo pattern CRUD sẵn có

Mỗi command = 1 file mới trong `cmd/vks/`, đăng ký trong `vks.go`, theo khuôn
`get_nodegroup.go` / `create_nodegroup.go`. Dùng lại `createClient`,
`outputResult`, `validator.ValidateID`, `parseLabels`, `parseTaints`.

**Nguyên tắc PATCH:** chỉ đưa field vào body khi flag được set
(`cmd.Flags().Changed(name)`), tránh ghi đè nhầm giá trị hiện có.

### 1. `upgrade-nodegroup-version` — file `upgrade_nodegroup_version.go`
- Flags: `--cluster-id` (req), `--nodegroup-id` (req), `--kubernetes-version` (req).
- Body: `{"kubernetesVersion": <v>}`.
- Validate cluster-id, nodegroup-id qua `ValidateID`.

### 2. `config-auto-healing` — file `config_auto_healing.go`
- Flags: `--cluster-id` (req), `--enable-auto-healing` (bool, req),
  `--max-unhealthy` (string), `--unhealthy-range` (string),
  `--timeout-unhealthy` (int).
- Body (chỉ field đã `Changed`, trừ `enableAutoHealing` luôn gửi):
  `{enableAutoHealing, maxUnhealthy?, unhealthyRange?, timeoutUnhealthy?}`.
- Method: `Patch`.

### 3. `update-nodegroup-metadata` — file `update_nodegroup_metadata.go`
- Flags: `--cluster-id` (req), `--nodegroup-id` (req), `--labels`, `--tags`,
  `--taints`.
- Body (chỉ field đã `Changed`): `{labels?, tags?, taints?}`.
  `labels`/`tags` parse qua `parseLabels`; `taints` qua `parseTaints`.
- Method: `Patch`.

### 4. `get-cluster-events` — file `get_cluster_events.go`
- Flags: `--cluster-id` (req), `--action`, `--type`, `--page`, `--page-size`.
- GET với query params; trả `{items, total, page, pageSize}` → `outputResult`.

### 5. `get-nodegroup-events` — file `get_nodegroup_events.go`
- Flags: `--cluster-id` (req), `--nodegroup-id` (req), `--action`, `--type`,
  `--page`, `--page-size`.
- Như trên, endpoint nodegroup.

### 6. `list-cluster-versions` — file `list_cluster_versions.go`
- Flags: chỉ flag chung (không cần cluster-id).
- GET `/v1/cluster-versions` → `outputResult`. Dùng để biết version hợp lệ trước
  khi chạy `upgrade-nodegroup-version` / nâng cluster.

## Kubeconfig (2 command)

API: `POST /kubeconfig {expirationDays}` **sinh** kubeconfig (trả **202 Accepted**,
async, status `NONE → CREATING → ACTIVE`); `GET /kubeconfig` trả YAML hoàn chỉnh
`{kubeConfig, expirationAt, expirationDays, status, renewalWarning}` với token
nhúng sẵn (khác EKS dùng exec/get-token).

### 7. `generate-kubeconfig` — file `generate_kubeconfig.go`
- Flags: `--cluster-id` (req), `--expiration-days` (int, default 30).
- POST `/kubeconfig {expirationDays}`.
- Vì trả 202 (async), in thông báo: kubeconfig đang được tạo, chờ status ACTIVE
  rồi chạy `update-kubeconfig`.
- (Tùy chọn, có thể đưa đợt sau) `--wait`: poll `GET /kubeconfig` đến khi
  `status == ACTIVE`.

### 8. `update-kubeconfig` — file `update_kubeconfig.go` (merge kiểu EKS)
Luồng:
1. GET `/kubeconfig`.
   - `status == NONE` → lỗi, gợi ý chạy `generate-kubeconfig` trước.
   - `status == CREATING` → thông báo đang tạo, thử lại sau (không phải lỗi cứng).
   - `status == ERROR` → báo lỗi.
   - `renewalWarning == true` → in cảnh báo, gợi ý `generate-kubeconfig` để gia hạn
     (vẫn tiếp tục merge).
2. Parse field `kubeConfig` (YAML đầy đủ).
3. Merge `clusters[]` / `contexts[]` / `users[]` vào target kubeconfig:
   - target = `--kubeconfig` | `$KUBECONFIG` (entry đầu) | `~/.kube/config`.
   - context name = `--alias` | `vks_<clusterId>`.
   - set `current-context` (trừ khi `--no-set-context`).
   - tạo file/dir nếu chưa có; giữ nguyên entry khác.
4. `--dry-run`: in những gì sẽ ghi, không ghi file.

- Flags: `--cluster-id` (req), `--kubeconfig`, `--alias`, `--no-set-context`,
  `--dry-run`.
- Logic merge tách ra package `internal/kubeconfig` (load → merge → write) để test
  độc lập, tương tự `awscli/customizations/eks/kubeconfig.py`.

## Đăng ký command

Trong `cmd/vks/vks.go`, `AddCommand` cho cả 8 command, gom nhóm theo comment:
cluster ops / nodegroup ops / kubeconfig / versions & events.

## Testing
- Unit test `internal/kubeconfig`: merge mới, ghi đè cùng context, giữ context
  khác, file rỗng/chưa tồn tại, set/không set current-context.
- Unit test build body cho `upgrade-version` / `config-auto-healing` /
  `update-nodegroup-metadata` (đặc biệt logic `Changed()` chỉ gửi field đã set).
- Unit test helper query phân trang.
- Theo khuôn test Go sẵn có trong repo.

## Rủi ro / lưu ý
- Tên `config-auto-healing` lệch quy ước hiện có (`set-auto-upgrade-config`,
  `delete-auto-upgrade-config`). Chấp nhận theo yêu cầu; cân nhắc thống nhất sau.
- `generate-kubeconfig` async (202): `update-kubeconfig` ngay sau đó có thể gặp
  `CREATING` — đã xử lý bằng thông báo thử lại.
- VKS kubeconfig là token nhúng (không exec) → khi token hết hạn phải
  `generate-kubeconfig` lại; không tự refresh như EKS.
