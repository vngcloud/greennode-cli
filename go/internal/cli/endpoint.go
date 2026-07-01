package cli

import (
	"fmt"
	"net/url"
	"os"
	"strings"
)

// trustedEndpointSuffix is the domain grn's own services live under. Requests to
// hosts outside it are flagged because grn sends a reusable IAM bearer token
// with every request (see WarnIfUntrustedEndpoint).
const trustedEndpointSuffix = ".vngcloud.vn"

// IsTrustedEndpoint reports whether endpointURL targets a host within the
// trusted vngcloud.vn domain. An empty value means no --endpoint-url override
// was given (the built-in region endpoint is used), which is trusted.
func IsTrustedEndpoint(endpointURL string) bool {
	if endpointURL == "" {
		return true
	}
	u, err := url.Parse(endpointURL)
	if err != nil {
		return false
	}
	host := u.Hostname()
	if host == "" {
		return false
	}
	return host == "vngcloud.vn" || strings.HasSuffix(host, trustedEndpointSuffix)
}

// WarnIfUntrustedEndpoint prints a security warning to stderr when endpointURL
// points outside the trusted vngcloud.vn domain. grn authenticates against the
// real IAM and sends the resulting reusable bearer token to whatever host
// --endpoint-url names, so a malicious or mistyped host can capture and replay
// that token. This warns; it does not block.
func WarnIfUntrustedEndpoint(endpointURL string) {
	if IsTrustedEndpoint(endpointURL) {
		return
	}
	u, _ := url.Parse(endpointURL)
	fmt.Fprintf(os.Stderr,
		"Warning: --endpoint-url %q is outside the trusted %s domain. grn will send your IAM bearer token to this host, and a bearer token can be replayed. Only use endpoints you trust.\n",
		u.Hostname(), strings.TrimPrefix(trustedEndpointSuffix, "."))
}
