# Execute Summary: Deep Smoke Harness Integration

Project: `ultraplan-go`
Sprint: `27-deep-smoke`
Status: `complete`
Completed: `2026-07-22T18:22:40Z`
Plan: `projects/ultraplan-go/sprints/27-deep-smoke/plan.md`
Run state: `projects/ultraplan-go/sprints/27-deep-smoke/.run-state.json`
Implementation target: `/home/antonioborgerees/coding/ultraplan-go`
External harness: `/home/antonioborgerees/coding/ultraplan-go-smoke`

## Task Outcomes

| Task | Status | Implementation Evidence |
| --- | --- | --- |
| 1. Catalog and smoke configuration | complete | `internal/project`, `internal/platform/config`; cataloged external manifest, source-aware precedence, ceilings, redaction |
| 2. Generic direct process boundary | complete | `internal/platform/process`; direct argv, exact cwd/env, bounded streams/progress, timeout/cancellation, process-group cleanup |
| 3. Protocol-v1 preflight and selection | complete | `internal/sprint/smoke_protocol.go`, external `ultraplan-smoke.json` and `src/protocol.ts` |
| 4. Gate, run, evidence, and verdict | complete | `internal/sprint/smoke.go`, `smoke_types.go`; current-review gate, strict evidence checks, five product-owned verdicts |
| 5. Artifact and flow state | complete | smoke stage registration, required-section renderer/validator, atomic `smoke.md`, fingerprinted smoke state and reconciliation diagnostics |
| 6. App, CLI, JSON, status, validation | complete | guarded `sprint smoke`, static status, `validate smoke`, stable typed errors/results, shared operational preview |
| 7. TUI parity | complete | smoke artifact/navigation, readiness, preview, confirmation, bounded progress, cancellation, result and diagnostic override actions |
| 8. Layered fake-harness and safety matrix | complete | recording runner, direct child-process tests, protocol/selection/evidence/preservation tests, app/TUI scenarios, race gate |
| 9. Documentation and gated real harness | complete | embedded template/defaults, five documentation guides, independent adapter discovery, explicit opt-in test and truthful review-gate skip |

## Verification Evidence

| Command | Result | Notes |
| --- | --- | --- |
| `go test ./internal/platform/config ./internal/platform/process ./internal/project ./internal/sprint ./internal/app ./internal/tui ./internal/workspace` | pass | Focused configuration, catalog, lifecycle, protocol, persistence, CLI/app, TUI, and embedded-default tests. |
| `go test ./...` | pass | All repository packages passed without launching the real harness. |
| `go test -race ./...` | pass | All packages passed under the race detector, including process/progress/cancellation paths. |
| `go build ./cmd/ultraplan` | pass | CLI built successfully. |
| `git diff --check` | pass | Target repository changes have no whitespace errors. |
| `ultraplan ... validate plan` | pass | Completed plan checklist and traceability validate after supported nested-task normalization. |
| `ultraplan ... status --json` | pass | Execute state loads as complete and smoke readiness truthfully reports the missing-current-review gate. |
| `node_modules/.bin/tsx src/protocol.ts discover --target /home/antonioborgerees/coding/ultraplan-go` | pass | Protocol `1.0`, harness `ultraplan-go-smoke`, available prerequisites, and complete Sprint 27 `ultra-deep` mapping returned as one JSON object. |
| `ULTRAPLAN_REAL_SMOKE=1 go test ./internal/sprint -run TestRealSmokeHarness -v` | pass with truthful skip | Real execution was blocked before discovery/run because Sprint 27 lacks a current review: `smoke review_gate: a current review is required`. |
| `go list -f ... ./internal/platform/process ./internal/sprint ./internal/app ./internal/tui` | pass | Platform process imports only the standard library; ownership and dependency direction match Decisions 1 and 6. |

## Architecture And Protocol Evidence

- `internal/project` owns catalog parsing only. `internal/sprint` owns manifest meaning, gates, selection, evidence validation, verdicts, Markdown, and state.
- `internal/platform/process` is product-agnostic and invokes a direct executable without a shell. It exposes bounded lifecycle facts and supported Unix process-group cleanup.
- `internal/app` composes the service and presents typed operations. CLI and TUI adapt those results rather than duplicating smoke rules or invoking each other.
- Review freshness is checked before unnecessary discovery/run work. Failed/blocked review override remains diagnostic and cannot replace current smoke evidence.
- Manifest, executable, cwd, evidence roots, returned evidence, and issue files are canonicalized and contained. Structured identities, counts, scope, hashes, and issue status are validated before verdict synthesis.
- Candidate `smoke.md` is validated and atomically renamed before smoke flow state advances. Malformed/process/evidence failures preserve the last valid artifact.
- Normal smoke mutations are confined to `smoke.md`, governed flow state, and manifest-declared external evidence roots. Sprint 28 integrated verify/recovery remains excluded.

## Policy Resolutions And Deviations

- `OQ-01`: hard maxima are discovery 5 minutes, run 24 hours, cleanup 30 seconds, stdout 64 MiB, and stderr 16 MiB. Defaults remain 30 seconds, 30 minutes, 5 seconds, 4 MiB, and 1 MiB.
- `OQ-02`: the base environment allowlist is `PATH`, `HOME`, `TMPDIR`, `LANG`, and `LC_ALL`. Manifest names require configured allowlisting; values are never persisted.
- `OQ-03`: the opt-in gate is `ULTRAPLAN_REAL_SMOKE=1`, with `go test ./internal/sprint -run TestRealSmokeHarness -v` as the documented lane.
- `flow --to smoke` is intentionally non-mutating/dry-run only. Mutating integrated verification and automated reconciliation remain Sprint 28 scope.
- Post-completion `validate execute` is not a release gate: by design it rejects already-checked tasks as non-executable. The durable execute state itself loads successfully through `status --json`, while `validate plan` passes.
- The external legacy harness's project-wide `npm run build` reports pre-existing TypeScript errors in older test sources. The delivered adapter executes successfully through the existing `tsx` runtime and emits valid protocol-v1 discovery. No real smoke pass is claimed.

## Real Harness Blocker

The catalog, manifest, executable, target, adapter, and runtime prerequisite are available. A real run was not authorized because Sprint 27 has no current validated `review.md`. The product stopped at the required review gate, the opt-in test reported a skip, and no `smoke.md` or external run evidence was manufactured.

## Files Changed For Sprint 27

Target repository:

- `docs/cli-reference.md`
- `docs/configuration.md`
- `docs/planning-smoke.md`
- `docs/recovery.md`
- `docs/user-guide.md`
- `internal/app/config_commands.go`
- `internal/app/operations.go`
- `internal/app/sprint_commands.go`
- `internal/app/sprint_commands_test.go`
- `internal/app/sprint_usecases.go`
- `internal/app/tui_commands.go`
- `internal/app/usecases.go`
- `internal/platform/config/config.go`
- `internal/platform/config/config_test.go`
- `internal/platform/config/redaction.go`
- `internal/platform/process/process.go`
- `internal/platform/process/process_other.go`
- `internal/platform/process/process_test.go`
- `internal/platform/process/process_unix.go`
- `internal/project/domain.go`
- `internal/project/index.go`
- `internal/project/project_test.go`
- `internal/project/validation.go`
- `internal/sprint/artifacts.go`
- `internal/sprint/doc.go`
- `internal/sprint/domain.go`
- `internal/sprint/flow.go`
- `internal/sprint/service.go`
- `internal/sprint/smoke.go`
- `internal/sprint/smoke_protocol.go`
- `internal/sprint/smoke_test.go`
- `internal/sprint/smoke_types.go`
- `internal/sprint/state.go`
- `internal/tui/doc.go`
- `internal/tui/model.go`
- `internal/workspace/init.go`
- `internal/workspace/scaffold/templates/smoke.md`

External harness:

- `ultraplan-smoke.json`
- `src/protocol.ts`

Sprint evidence:

- `plan.md`
- `execute.md`
- `.run-state.json`

The pre-existing local `go.mod`/`go.sum` agentwrap replacement was preserved and is not attributed to Sprint 27.
