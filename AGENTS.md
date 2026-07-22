# AGENTS.md

## What this repo is

`caddy-plus` builds [Caddy](https://github.com/caddyserver/caddy) with a pinned
set of plugins via GitHub Actions. There is **no application source code** in
this repo — it is a build/packaging pipeline:

- `plugins.txt` — the plugin manifest (`module@version` per line, `#` comments)
- `docker/Dockerfile` — scratch-based image (binary + CA certs, runs as nobody)
- `.github/workflows/detect.yml` — polls upstream caddy releases, tags, dispatches
- `.github/workflows/release.yml` — build → finalize → SLSA provenance →
  sign/attest → docker → verify → publish
- `.github/workflows/lint.yml` — actionlint + zizmor
- `.github/workflows/scorecard.yml` — OpenSSF Scorecard

## Rules when editing

1. **Every action MUST be pinned to a full-length commit SHA** with a
   `# vX.Y.Z` version comment. Exception: `slsa-github-generator` reusable
   workflows, which must be referenced by tag (SLSA design).
2. Never introduce `@latest` or unpinned versions anywhere — plugins.txt
   entries must match `^module@v\d…` (enforced by the pipeline).
3. Every job needs `step-security/harden-runner` with `egress-policy: block`
   and a minimal `allowed-endpoints` list if it touches the network.
4. Workflow-level `permissions: {}`; grant least privilege per job.
5. Use `persist-credentials: false` on checkout unless the job pushes
   (only `detect.yml` pushes, and it documents why).
6. Use env vars (`env:`) instead of `${{ ... }}` expressions inside `run:`
   blocks for any non-constant value (script-injection hygiene — zizmor
   enforces this).
7. No test suite exists and none should be added; the pipeline's "smoke check"
   step (`go version -m` against pinned versions) is the sanity gate.
8. Versioning: tags mirror Caddy (`v2.11.4`); rebuilds use `-plus.N` suffix.
9. If you change anything listed above (layout, workflow names, rules),
   update this file.

## Validating changes locally

- YAML sanity: parse the workflows (e.g. `python3 -c 'import yaml,sys; yaml.safe_load(open(...))'`).
- After workflow edits, CI runs actionlint + zizmor on the PR — both must pass.
- Full end-to-end validation requires dispatching a real release; coordinate
  with the maintainer before doing that.
