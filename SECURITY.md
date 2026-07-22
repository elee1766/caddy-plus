# Security Policy

## Reporting a vulnerability

If you find a security issue in this repository's build pipeline or released
artifacts, please report it via
[GitHub private vulnerability reporting](../../security/advisories/new).
Do not open a public issue.

Note: vulnerabilities in Caddy itself or in the bundled plugins should be
reported upstream ([caddyserver/caddy](https://github.com/caddyserver/caddy/security)
or the respective plugin repository).

## Supply-chain hardening model

The release pipeline is designed so that no single compromise silently
produces a malicious artifact:

| Layer | Control |
|---|---|
| Action integrity | Every action pinned to a full-length commit SHA with a version comment, enforced by zizmor/actionlint CI on every workflow change; Dependabot proposes bumps with a 7-day cooldown. (GitHub's repo-level SHA-pinning enforcement is intentionally off: it is incompatible with the SLSA generator, which must be referenced by tag.) |
| Runtime network | `step-security/harden-runner` in block mode with per-job egress allowlists on every job |
| Credentials | No long-lived secrets. `GITHUB_TOKEN` only, with `permissions: {}` at workflow level and least-privilege per-job grants. Signing is cosign **keyless** via OIDC — there is no signing key to steal |
| Dependencies | Caddy and every plugin are pinned in `go.mod`; the committed `go.sum` hash-pins all transitive dependencies and is verified on every build (`-mod=readonly`, Go checksum database) |
| Build flags | `CGO_ENABLED=0 -trimpath -buildvcs=false -ldflags=-s -w -buildid=` (static, stripped, reproducible) |
| Provenance | SLSA L3 provenance from an isolated builder (`slsa-github-generator`) **and** GitHub artifact attestations for every binary and the container image |
| Transparency | SPDX SBOM (syft) for binaries and image, attached and attested |
| Tamper resistance | Immutable releases enabled; pipeline verifies its own signatures and provenance before a release leaves draft |
| Static analysis | `actionlint` + `zizmor` on every workflow change; OpenSSF Scorecard runs weekly |

## Verifying artifacts

See [README.md](README.md#verifying-a-release) for copy-paste verification
commands (`gh attestation verify`, `slsa-verifier`, `cosign verify-blob` /
`cosign verify`).

## Scope limitations

- Provenance and signatures prove **how and where** an artifact was built.
  They do not vouch for the upstream source code of Caddy or the plugins.
- Plugin versions are bumped manually (via PR); consumers should review the
  plugin list before adopting a release.
