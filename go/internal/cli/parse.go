package cli

import (
	"fmt"
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
