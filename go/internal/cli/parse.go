package cli

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
)

// ParseCommaSeparated splits a comma-separated string into a trimmed slice,
// dropping empty entries.
func ParseCommaSeparated(s string) []string {
	if s == "" {
		return nil
	}
	parts := strings.Split(s, ",")
	result := make([]string, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p != "" {
			result = append(result, p)
		}
	}
	return result
}

// ParseStructFlag parses a struct-valued CLI flag accepting either JSON
// ({"minSize":2,"maxSize":10}) or AWS-style shorthand (minSize=2,maxSize=10).
// Keys listed in intFields are coerced to int in the shorthand form; all other
// values stay strings. JSON is passed through as decoded. An empty/blank value
// returns (nil, nil). Returns an error on malformed JSON, a shorthand entry
// without '=', or a non-integer value for an int field.
func ParseStructFlag(value string, intFields ...string) (map[string]interface{}, error) {
	v := strings.TrimSpace(value)
	if v == "" {
		return nil, nil
	}
	if strings.HasPrefix(v, "{") {
		var m map[string]interface{}
		if err := json.Unmarshal([]byte(v), &m); err != nil {
			return nil, fmt.Errorf("invalid JSON: %w", err)
		}
		return m, nil
	}

	ints := map[string]bool{}
	for _, f := range intFields {
		ints[f] = true
	}
	out := map[string]interface{}{}
	for _, pair := range strings.Split(v, ",") {
		pair = strings.TrimSpace(pair)
		if pair == "" {
			continue
		}
		idx := strings.Index(pair, "=")
		if idx < 0 {
			return nil, fmt.Errorf("invalid shorthand %q: expected key=value", pair)
		}
		key := strings.TrimSpace(pair[:idx])
		val := strings.TrimSpace(pair[idx+1:])
		if ints[key] {
			n, err := strconv.Atoi(val)
			if err != nil {
				return nil, fmt.Errorf("%s must be an integer, got %q", key, val)
			}
			out[key] = n
		} else {
			out[key] = val
		}
	}
	return out, nil
}

// BuildEventsQuery builds query params for events endpoints, including only
// flags the user explicitly set. `changed` maps flag name -> was it set.
// VKS pagination is 0-based, so page is passed through verbatim.
func BuildEventsQuery(action, eventType string, page, pageSize int, changed map[string]bool) map[string]string {
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
