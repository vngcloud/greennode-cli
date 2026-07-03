package client

import (
	"encoding/json"
	"strings"
)

// redactDebugBody returns a copy of a JSON request/response body safe to print
// in --debug output: values of credential-bearing fields (e.g. the embedded
// kubeconfig, tokens, client certs/keys, secrets) are replaced with
// "[REDACTED]". Only used for debug logging — the real body is never modified.
// Non-JSON bodies are returned unchanged (VKS/vServer bodies are JSON).
func redactDebugBody(raw string) string {
	trimmed := strings.TrimSpace(raw)
	if trimmed == "" {
		return raw
	}
	var v interface{}
	if err := json.Unmarshal([]byte(trimmed), &v); err != nil {
		return raw
	}
	redactValue(v)
	out, err := json.Marshal(v)
	if err != nil {
		return raw
	}
	return string(out)
}

func redactValue(v interface{}) {
	switch t := v.(type) {
	case map[string]interface{}:
		for k, val := range t {
			if isSensitiveKey(k) {
				t[k] = "[REDACTED]"
				continue
			}
			redactValue(val)
		}
	case []interface{}:
		for _, item := range t {
			redactValue(item)
		}
	}
}

// isSensitiveKey reports whether a JSON key's value should be redacted in debug
// output. Matches known credential fields plus common secret-ish key patterns.
func isSensitiveKey(k string) bool {
	lk := strings.ToLower(k)
	switch lk {
	case "kubeconfig", "token", "client-certificate-data", "client-key-data",
		"client_secret", "clientsecret", "password":
		return true
	}
	return strings.Contains(lk, "secret") ||
		strings.Contains(lk, "password") ||
		strings.HasSuffix(lk, "token") ||
		strings.HasSuffix(lk, "key-data")
}
