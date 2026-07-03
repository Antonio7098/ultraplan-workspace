# Deep Smoke: Sprint 17 — Artifact Domain and Flow State

## Scope

- Sprint: `17-artifact-domain-and-flow-state`
- Implementation directory: `/home/antonioborgerees/coding/ultraplan-go`
- Smoke harness: `/home/antonioborgerees/coding/ultraplan-go-smoke/`
- Selected smoke level/tests: **None — not applicable**
- Reason scope is sufficient: Sprint 17 is a runtime-free sprint. The sprint-index explicitly excludes deep-smoke and the smoke harness: "Runtime-backed smoke evidence is not needed for a runtime-free status and flow-state sprint." The sprint touches only offline status inspection, flow-state persistence, and sprint discovery — no OpenCode, model calls, subprocess execution, or runtime behavior. The acceptance criteria are satisfied by `go test ./...`, `go test -race ./...`, and `go build ./cmd/ultraplan`, all of which pass.

## Run Evidence

No smoke harness run was executed. Sprint 17 has no sprint-specific smoke suite (`--level sprint-17` does not exist) and none is needed because the sprint's delivered surface is entirely runtime-free.

### Offline Verification (sprint acceptance criteria)

| Check | Result | Summary |
| --- | --- | --- |
| `go test ./...` | Passed | All packages pass including `internal/sprint`, `internal/app`, `internal/project`, `internal/study`, `internal/workspace` |
| `go test -race ./...` | Passed | No data races detected across all packages |
| `go build ./cmd/ultraplan` | Passed | Binary builds cleanly |
| `go run ./cmd/ultraplan sprint --help` | Passed | Help shows only `sprint <project> <sprint> status` with documented scope and non-goals |

### Sprint Command End-to-End (test workspace)

| Check | Result | Summary |
| --- | --- | --- |
| `TestSprintHelpIsRegistered` | Passed | Sprint command registered in top-level help; subcommand help renders correctly |
| `TestSprintStatusRefreshesStateAndRendersDeterministically` | Passed | Status output includes project, sprint, flow-state path, deterministic stage ordering, artifact paths; no ANSI escapes; `flow-state.json` written |
| `TestSprintStatusErrorsAndInvalidFlowStateExitFive` | Passed | Ambiguous refs exit 5; missing projects exit 5; malformed flow-state exits 5 with "flow state malformed" diagnostic; invalid state preserved |
| `TestSprintMalformedArgumentsUseUsageExit` | Passed | Missing `<sprint>` arg produces usage error with exit 2 |

## Findings

| ID | Severity | Status | Evidence | Required Action |
| --- | --- | --- | --- | --- |
| (none) | — | — | — | — |

No findings. Sprint 17 delivers runtime-free status inspection only. All acceptance criteria pass through offline verification. No live runtime behavior is in scope.

## Open Issues

- (none)

## Resolved Issues

- (none)

## Verdict

**pass** — Sprint 17 is a runtime-free sprint that explicitly excludes smoke investigation per its sprint-index and requirements. The sprint's acceptance criteria (`go test ./...`, `go test -race ./...`, `go build ./cmd/ultraplan`) all pass. No smoke harness run is required or applicable. The deep-smoke protocol defers to offline verification for this sprint scope.
