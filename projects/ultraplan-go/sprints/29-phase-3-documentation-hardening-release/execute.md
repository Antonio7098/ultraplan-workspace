# Sprint Execution: Phase 3 Documentation, Hardening, and Release

> Project: `ultraplan-go`
> Sprint: `29-phase-3-documentation-hardening-release`
> Status: `implementation complete; changes uncommitted`
> Executed: `2026-07-28`

## Target

- Repository: `/home/antonioborgerees/coding/ultraplan-go`
- Starting revision: `47e4e71a3658b5e68104d0b5d2bf7975dc1aa221`
- Starting state: clean
- Ending state: Sprint 29 implementation changes present and uncommitted

## Implemented

- Added canonical runtime `PromptReference` identity with ID, version, owner, purpose, and SHA-256 checksum.
- Added request trace identity and propagated trace/prompt identity through runtime metadata.
- Added regression coverage for prompt and trace propagation.
- Made flow-state migration and verification status inspection read-only. They derive migrated/expired truth without writing during a status/load operation.
- Updated migration and expired-attempt tests to prove the read boundary.
- Removed `PATH` and `HOME` from the default smoke environment; retained only `TMPDIR`, `LANG`, and `LC_ALL`.
- Updated the configuration reference for the tightened environment policy.
- Replaced stale pre-Phase-3 release-checklist scope and added test, race, build, vet, diff, fake-runtime, fake-harness, and truthful-blocked gates.
- Added review, smoke, and verify CLI examples.
- Added stable Phase 3 JSON schema documentation for review, smoke, verify, and status.
- Added migration guidance from manual `review.md`/`deep-smoke.md` to generated verification artifacts.
- Linked the implementation README to Phase 3 operator docs and declared the adjacent planning workspace as the authoritative product-document source.
- Updated recovery guidance for read-only expired-attempt derivation.

The implementation also retains the immediately preceding hardening work at `47e4e71`, including SDK diagnostic redaction, typed smoke error cause preservation, expired-attempt recovery modeling, permission enforcement coverage, structured-output repair coverage, and corrected verify/recovery documentation.

## Changed Paths

```text
README.md
docs/cli-reference.md
docs/configuration.md
docs/phase3-json-schemas.md
docs/phase3-migration.md
docs/recovery.md
docs/release-checklist.md
internal/platform/config/config.go
internal/platform/runtime/runtime.go
internal/sprint/runtime_progress_test.go
internal/sprint/service.go
internal/sprint/smoke_types.go
internal/sprint/state.go
internal/sprint/verify.go
internal/sprint/verify_test.go
```

## Verification

The managed environment exposes the default Go build cache read-only, so the release gate used `GOCACHE=/tmp/ultraplan-go-s29-cache`. The binary was written to `/tmp/ultraplan-go-s29` to avoid an untracked repository artifact.

| Command | Result |
| --- | --- |
| `go test ./internal/platform/runtime ./internal/platform/config ./internal/sprint ./internal/app ./internal/tui` | pass |
| `go test ./...` | pass |
| `go test -race ./...` | pass |
| `go build -o /tmp/ultraplan-go-s29 ./cmd/ultraplan` | pass |
| `go vet ./...` | pass |
| `git diff --check` | pass |

All commands above were run in `/home/antonioborgerees/coding/ultraplan-go`.

## Scope Note

Per user direction, the removed dogfood task was not performed. No live OpenCode run or external smoke-harness campaign was launched. The normal fake-backed, race, build, vet, and formatting release gates passed.
