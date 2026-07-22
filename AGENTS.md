# AGENTS.md

## What this repo is

`caddy-plus` builds [Caddy](https://github.com/caddyserver/caddy) with a pinned
set of plugins via GitHub Actions. The buildable source is a small Go module:

- `main.go` — blank imports of caddy core + all plugins (the plugin list)
- `go.mod` / `go.sum` — pins every direct and transitive dependency
- `docker/Dockerfile` — scratch-based image (binary + CA certs, runs as nobody)
- `.github/workflows/detect.yml` — polls upstream caddy releases, bumps
  `go.mod`, commits, tags, dispatches release.yml
- `.github/workflows/release.yml` — build (go build) → finalize → SLSA
  provenance → sign/attest → docker → verify → publish
- `.github/workflows/lint.yml` — actionlint + zizmor
- `.github/workflows/scorecard.yml` — OpenSSF Scorecard

## Rules when editing

1. **Every action MUST be pinned to a full-length commit SHA** with a
   `# vX.Y.Z` version comment. Exception: `slsa-github-generator` reusable
   workflows, which must be referenced by tag (SLSA design). This is why
   GitHub's repo-level `sha_pinning_required` policy stays OFF — zizmor
   enforces pinning instead; do not re-enable the repo policy.
2. Never introduce `@latest` or unpinned versions anywhere — all Go
   dependencies are pinned via `go.mod`/`go.sum` (bumped by Dependabot PRs or
   by detect.yml for caddy core); builds use `-mod=readonly`.
3. Adding/removing a plugin = edit the blank imports in `main.go` and run
   `go mod tidy` (or let Dependabot propose version bumps).
4. Every job needs `step-security/harden-runner` with `egress-policy: block`
   and a minimal `allowed-endpoints` list if it touches the network.
5. Workflow-level `permissions: {}`; grant least privilege per job.
6. Use `persist-credentials: false` on checkout unless the job pushes
   (only `detect.yml` pushes, and it documents why).
7. Use env vars (`env:`) instead of `${{ ... }}` expressions inside `run:`
   blocks for any non-constant value (script-injection hygiene — zizmor
   enforces this).
8. No test suite exists and none should be added; the pipeline's "smoke check"
   step (`go version -m` against go.mod) is the sanity gate.
9. Versioning: tags mirror Caddy (`v2.11.4`); rebuilds use `-plus.N` suffix.
10. If you change anything listed above (layout, workflow names, rules),
    update this file.

## Validating changes locally

- YAML sanity: parse the workflows (e.g. `python3 -c 'import yaml,sys; yaml.safe_load(open(...))'`).
- Build sanity: `CGO_ENABLED=0 go build -trimpath -mod=readonly -o /tmp/caddy .`
  then `/tmp/caddy version` and `go version -m /tmp/caddy`.
- After workflow edits, CI runs actionlint + zizmor on the PR — both must pass.
- Full end-to-end validation requires dispatching a real release; coordinate
  with the maintainer before doing that.
