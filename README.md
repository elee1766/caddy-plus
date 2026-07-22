# caddy-plus

Automated, supply-chain-hardened builds of [Caddy](https://caddyserver.com) with a
pinned set of plugins. Whenever `caddyserver/caddy` publishes a new **stable**
release, this repository builds it from source and publishes signed, attested
binaries and a container image.

## What you get

| Artifact | Where |
|---|---|
| `caddy-plus_<tag>_linux_amd64.tar.gz` | [Releases](../../releases) |
| `caddy-plus_<tag>_linux_arm64.tar.gz` | [Releases](../../releases) |
| `caddy-plus_<tag>_freebsd_amd64.tar.gz` | [Releases](../../releases) |
| `checksums.txt` (SHA-256) | [Releases](../../releases) |
| SLSA L3 provenance (`*.intoto.jsonl`) | attached to each release |
| cosign keyless bundles (`*.bundle`) | attached to each release |
| SPDX SBOM | attached to each release |
| Container image | `ghcr.io/elee1766/caddy-plus:<tag>` (multi-arch, also `:latest`) |

## Plugins

The build contents are defined by [`plugins.txt`](plugins.txt) — one
`module@version` per line, always pinned. The pipeline verifies at build time
that every plugin is present in the binary at the exact pinned version.

## Versioning

- Releases mirror the upstream Caddy version: Caddy `v2.11.4` → tag `v2.11.4`.
- Rebuilds of the same Caddy version (e.g. plugin bumps) use a suffix:
  `v2.11.4-plus.1`, `v2.11.4-plus.2`, …
- Only stable upstream releases are built; betas/RCs are skipped.
- Releases are immutable once published.

## Verifying a release

Every binary can be verified three independent ways:

```sh
TAG=v2.11.4
ASSET=caddy-plus_${TAG}_linux_amd64.tar.gz

# 1. GitHub artifact attestation
gh attestation verify "${ASSET}" -R elee1766/caddy-plus

# 2. SLSA L3 provenance
slsa-verifier verify-artifact \
  --provenance-path "caddy-plus_${TAG}.intoto.jsonl" \
  --source-uri github.com/elee1766/caddy-plus \
  --source-tag "${TAG}" \
  "${ASSET}"

# 3. cosign keyless signature
cosign verify-blob \
  --bundle "${ASSET}.bundle" \
  --certificate-identity-regexp '^https://github.com/elee1766/caddy-plus/\.github/workflows/release\.yml@refs/tags/' \
  --certificate-oidc-issuer https://token.actions.githubusercontent.com \
  "${ASSET}"
```

Container image:

```sh
cosign verify ghcr.io/elee1766/caddy-plus:${TAG} \
  --certificate-identity-regexp '^https://github.com/elee1766/caddy-plus/\.github/workflows/release\.yml@refs/tags/' \
  --certificate-oidc-issuer https://token.actions.githubusercontent.com

gh attestation verify oci://ghcr.io/elee1766/caddy-plus:${TAG} -R elee1766/caddy-plus
```

## Container usage

```sh
docker run -d -p 80:80 -p 443:443 -p 443:443/udp \
  -v $PWD/Caddyfile:/etc/caddy/Caddyfile:ro \
  -v caddy_data:/data -v caddy_config:/config \
  ghcr.io/elee1766/caddy-plus:latest
```

The image is `scratch`-based (just the static binary + CA certificates) and
runs as the unprivileged `nobody` user.

## How it works

- [`detect.yml`](.github/workflows/detect.yml) polls upstream releases every
  6 h, tags this repo, and dispatches the release pipeline.
- [`release.yml`](.github/workflows/release.yml) builds with
  [`xcaddy`](https://github.com/caddyserver/xcaddy) (pinned), smoke-checks the
  binary, and produces all signatures/attestations before publishing.
- All actions are pinned to full-length commit SHAs, every network-touching
  job runs behind [step-security/harden-runner](https://github.com/step-security/harden-runner)
  egress blocking, and there are no long-lived signing keys anywhere (cosign
  keyless via GitHub OIDC).

See [SECURITY.md](SECURITY.md) for the full hardening model.

## License

Apache-2.0 (same as Caddy). Built binaries contain Caddy and the plugins
listed in `plugins.txt`, which carry their own licenses.
