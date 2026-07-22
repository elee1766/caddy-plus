// caddy-plus: Caddy with a pinned set of plugins.
//
// The module set is defined here (blank imports) and pinned in go.mod/go.sum.
// This mirrors exactly what xcaddy generates, but committed for auditability.
package main

import (
	caddycmd "github.com/caddyserver/caddy/v2/cmd"

	// plug in Caddy modules here
	_ "github.com/caddyserver/caddy/v2/modules/standard"

	// caddy-plus plugins
	_ "github.com/aksdb/caddy-cgi/v2"
	_ "github.com/caddy-dns/cloudflare"
	_ "github.com/greenpau/caddy-security"
	_ "github.com/greenpau/caddy-trace"
	_ "github.com/hslatman/caddy-crowdsec-bouncer"
	_ "github.com/mholt/caddy-dynamicdns"
	_ "github.com/mholt/caddy-l4"
)

func main() {
	caddycmd.Main()
}
