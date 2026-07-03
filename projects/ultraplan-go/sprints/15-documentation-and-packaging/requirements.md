# Sprint Requirements: Documentation and Packaging

> Project: `ultraplan-go`
> Sprint: `15-documentation-and-packaging`
> Purpose: the authoritative, human-readable sprint contract. All other sprint artifacts must satisfy these requirements.

## Sprint Goal

Prepare the first production-ready study-side release by completing user/operator documentation, packaging Linux and macOS CLI builds with checksums, and recording offline plus gated OpenCode smoke-release evidence.

## Required Outputs

| Output | Path | Description |
| ---------- | -------- | --------------- |
| User guide | `/home/antonioborgerees/coding/ultraplan-go/docs/user-guide.md` | End-to-end supported study-side workflow from workspace setup through study init, prompt preview, run, run-all/run-loop, validate, summary, status, and code extraction. |
| Configuration guide | `/home/antonioborgerees/coding/ultraplan-go/docs/configuration.md` | Supported `ultraplan.yml` fields, precedence, redaction behavior, runtime/model settings, retry/fallback settings, and documented schema-version rejection behavior. |
| Recovery runbook | `/home/antonioborgerees/coding/ultraplan-go/docs/recovery.md` | Operator guidance for validation failures, missing artifacts, stale running tasks, lock diagnostics, `--force-unlock`, cancellation recovery, retry/fallback interpretation, and loud write failures. |
| OpenCode smoke instructions | `/home/antonioborgerees/coding/ultraplan-go/docs/opencode-smoke.md` | Environment prerequisites and step-by-step gated real-runtime smoke procedure using agentwrap/OpenCode without making it part of the default test suite. |
| Release checklist | `/home/antonioborgerees/coding/ultraplan-go/docs/release-checklist.md` | Repeatable release gate covering tests, race tests, build, packaging, checksums, smoke evidence, local `replace` audit, and supported-platform notes. |
| CLI command reference | `/home/antonioborgerees/coding/ultraplan-go/docs/cli-reference.md` | Documented public command surface and output modes for the study-side release, including stable JSON surfaces from Sprint 14. |
| README release update | `/home/antonioborgerees/coding/ultraplan-go/README.md` | Top-level install, quickstart, supported scope, documentation index, and explicit deferred target/sprint workflow notice. |
| Linux amd64 binary | `/home/antonioborgerees/coding/ultraplan-go/dist/ultraplan-linux-amd64` | Release build of `./cmd/ultraplan` for Linux amd64. |
| Linux arm64 binary | `/home/antonioborgerees/coding/ultraplan-go/dist/ultraplan-linux-arm64` | Release build of `./cmd/ultraplan` for Linux arm64. |
| macOS amd64 binary | `/home/antonioborgerees/coding/ultraplan-go/dist/ultraplan-darwin-amd64` | Release build of `./cmd/ultraplan` for macOS amd64. |
| macOS arm64 binary | `/home/antonioborgerees/coding/ultraplan-go/dist/ultraplan-darwin-arm64` | Release build of `./cmd/ultraplan` for macOS arm64. |
| Checksums | `/home/antonioborgerees/coding/ultraplan-go/dist/checksums.txt` | SHA-256 checksums for every release binary in `dist/`. |
| Smoke evidence | `/home/antonioborgerees/coding/ultraplan-go/dist/smoke-evidence.md` | Recorded results for offline release gates and gated OpenCode smoke status, including skipped reason when the real-runtime environment is unavailable. |

## Acceptance Criteria

- [ ] `go test ./...` passes from `/home/antonioborgerees/coding/ultraplan-go` without requiring OpenCode, provider credentials, network access, or real subprocess smoke fixtures.
- [ ] `go test -race ./...` passes from `/home/antonioborgerees/coding/ultraplan-go`.
- [ ] `go build ./cmd/ultraplan` passes from `/home/antonioborgerees/coding/ultraplan-go`.
- [ ] The release binaries listed in Required Outputs exist under `dist/` and are built from `./cmd/ultraplan` with `GOOS`/`GOARCH` matching their filenames.
- [ ] `dist/checksums.txt` contains one SHA-256 entry for each required binary and no entries for unrelated files.
- [ ] `dist/smoke-evidence.md` records the exact offline verification commands, pass/fail results, build targets, checksum generation command, and OpenCode smoke result or explicit skip reason.
- [ ] OpenCode smoke instructions specify prerequisites, required environment/config, commands to run, expected artifacts, cleanup notes, and how to keep the smoke path gated outside normal tests.
- [ ] User documentation covers only current study-side capabilities: workspace/config/health, study initialization/listing, Markdown source applicability, prompt preview, single run, synthesis, run-all, run-loop/status, validation, summary generation, and code extraction.
- [ ] Configuration documentation names supported config fields and precedence order: built-in defaults, workspace config, environment variables, and command flags.
- [ ] Configuration documentation states how secrets are redacted and how unsupported config/schema versions are rejected or diagnosed.
- [ ] Recovery documentation includes `--force-unlock` semantics, stale lock risk, stale running task recovery interpretation, validation failure repair workflow, cancellation recovery, missing output handling, retry/fallback metadata interpretation, and atomic write failure behavior.
- [ ] CLI command reference documents stable public JSON surfaces introduced by Sprint 14 and does not describe unstable JSON as stable.
- [ ] README links every new document, includes a quickstart, and explicitly states that target scaffolding, sprint planning, sprint execution, hosted SaaS, browser UI, and multi-user collaboration are out of scope for this release.
- [ ] Documentation does not claim unsupported commands, target/sprint workflows, provider-agnostic real-runtime guarantees, automatic Git mutation, hosted service behavior, or browser UI support.
- [ ] Release checklist includes a dependency audit confirming no local `replace github.com/Antonio7098/agentwrap => ../agentwrap` remains in the releasable module unless intentionally documented as non-release development state.
- [ ] Release checklist includes a review step for prompt/version metadata, runtime metadata redaction, stable JSON schema documentation, and recovery docs before publishing artifacts.

## Non-Goals

- Implementing new product behavior beyond documentation, packaging scripts or commands if needed, and release verification artifacts.
- Adding target scaffolding, sprint planning, sprint execution, target/sprint validators, or target/sprint runtime task kinds.
- Adding a hosted service, browser UI, multi-user auth, organization permissions, or collaboration features.
- Adding new runtime adapters beyond the existing agentwrap/OpenCode path.
- Making real OpenCode/provider smoke tests mandatory for default `go test ./...` or `go test -race ./...`.
- Introducing schema migrations unless Sprint 14 already added a schema change that must be documented or tested for release.
- Publishing GitHub releases, pushing tags, signing binaries, notarizing macOS binaries, or uploading artifacts unless separately requested.

## Constraints

- Current release scope is study-side only; docs and packaging must not present deferred target/sprint workflows as available.
- Runtime integration must remain through `github.com/Antonio7098/agentwrap` and `github.com/Antonio7098/agentwrap/opencode`; UltraPlan must not document or implement direct OpenCode supervision as product-owned behavior.
- Normal verification must stay fake-first and offline; real OpenCode smoke is gated by explicit environment availability and must have a documented skip path.
- Package/module boundaries from `docs/ARCHITECTURE.md` remain binding: product behavior stays in `internal/study`, workspace behavior in `internal/workspace`, code extraction in `internal/codeextract`, and generic runtime behavior in `internal/platform/runtime`.
- Release artifacts must be deterministic enough to audit: explicit build commands, explicit target list, and SHA-256 checksums are required.
- Generated documentation and release artifacts must not contain secrets, provider tokens, full sensitive environment variables, raw unsafe runtime payloads, or embedded prompt/report bodies except minimal examples.
- Existing stable JSON surfaces from Sprint 14 must be documented as compatibility-sensitive; undocumented or debug-only shapes must not be promised as stable.
- `dist/` artifacts are release outputs only; packaging must not require users to adopt a server, database, or opaque artifact store.

## Dependencies

| Prior Sprint / Output | Required For | Notes |
| ------------------------ | ------------------- | --------- |
| Sprint 1 CLI shell and app composition | Buildable release binary | `./cmd/ultraplan` must remain the single CLI entrypoint. |
| Sprint 2 workspace/config/health | User, configuration, and health documentation | Docs must reflect workspace discovery, config precedence, redaction, and health behavior. |
| Sprint 3 study domain/listing | User guide and CLI reference | Study/source/dimension listing semantics must be documented accurately. |
| Sprint 4 study initialization | User guide and quickstart | YAML initialization, generated files, dry-run, force/no-clone behavior, and study README behavior must be described. |
| Sprint 5 Markdown sources/applicability | User guide and validation docs | Markdown source discovery, frontmatter stripping, and inapplicable-pair behavior must be documented. |
| Sprint 6 report validation/rating parsing | Validation and recovery docs | Validation failure guidance must match actual report and rating validator behavior. |
| Sprint 7 prompt composition | User guide and CLI reference | Prompt preview commands and source-isolation/document-isolation rules must be documented. |
| Sprint 8 run-state/status | Recovery and status docs | Run-state path, schema handling, runtime-free status, and resume validation behavior must be documented. |
| Sprint 9 agentwrap/OpenCode runtime integration | OpenCode smoke instructions and config docs | Smoke path, runtime health, permission policy, retry/fallback, metadata, and local `replace` audit depend on this sprint. |
| Sprint 10 single analysis/synthesis | User guide and smoke procedure | Single runtime task behavior and validation-as-success-gate must be documented and smoke-tested when available. |
| Sprint 11 run-all batch execution | User guide and release smoke | Bounded batch execution, partial completion, summary generation, and cancellation behavior must be documented. |
| Sprint 12 durable run-loop | Recovery runbook and release checklist | Lock handling, `--force-unlock`, stale running task recovery, cancellation/resume, and run-state persistence are required release docs. |
| Sprint 13 summary/code extraction | User guide and CLI reference | `summary.csv` and `ultraplan code` behavior must be documented. |
| Sprint 14 validation/diagnostics/JSON stability | CLI reference and release compatibility notes | Stable JSON schemas, validation command, diagnostics, and final release gates must be documented. |
| External OpenCode environment | Gated real-runtime smoke | Required only when available; absence must be recorded as an explicit skip in `dist/smoke-evidence.md`. |

## Review Expectations

| What | How Verified |
| ------------- | ----------------------- |
| Scope stays study-side only | Review README and docs for absent target/sprint command claims and compare against PRD/TRD deferred scope. |
| Documentation matches implemented CLI | Run relevant `ultraplan --help` and command help output or inspect tested help strings, then compare with `docs/cli-reference.md`. |
| User and operator docs are actionable | Manual review of quickstart, config examples, recovery procedures, and smoke instructions for exact commands, expected outputs, and failure handling. |
| Stable JSON documentation is accurate | Compare `docs/cli-reference.md` against Sprint 14 implementation/tests and ensure only stable surfaces are labeled stable. |
| Packaging artifacts exist | Inspect `dist/` for the four required binaries, `checksums.txt`, and `smoke-evidence.md`. |
| Checksums are complete | Recompute or verify SHA-256 hashes for every required binary and compare with `dist/checksums.txt`. |
| Offline release gate passes | Review `dist/smoke-evidence.md` and rerun `go test ./...`, `go test -race ./...`, and `go build ./cmd/ultraplan` when feasible. |
| OpenCode smoke is gated and honest | Review `docs/opencode-smoke.md` and `dist/smoke-evidence.md` for either real pass evidence or a clear environment-unavailable skip reason. |
| Dependency release readiness | Inspect `go.mod` for local `replace` directives affecting `github.com/Antonio7098/agentwrap` and verify release checklist disposition. |
| Security and redaction | Review docs and smoke evidence for absence of secrets, raw unsafe payloads, full prompts, full environment dumps, or provider tokens. |
