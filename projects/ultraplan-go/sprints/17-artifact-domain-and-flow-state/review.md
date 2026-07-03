# Sprint Review: Sprint Artifact Domain and Flow State

> Project: `ultraplan-go`
> Sprint: `17-artifact-domain-and-flow-state`
> Date: `2026-06-14`
> Status: `passed`

## Implementation Summary

- Added `internal/sprint` as the sprint-owned package for planning-stage domain modeling, sprint discovery, artifact paths, strict `flow-state.json` load/write, filesystem artifact inspection, deterministic status derivation, and status refresh service.
- Added thin app wiring for `ultraplan sprint <project> <sprint> status` with help, deterministic plain text, stdout/stderr separation through existing app conventions, and exit-code mapping.
- Added sprint and command tests covering stage/status constants, unsupported values, direct-child sprint discovery, prefix resolution, artifact paths, strict flow-state validation, invalid-state preservation, atomic write preservation, status refresh, help, output ordering, diagnostics, and no ANSI output.

## Files Changed

- `/home/antonioborgerees/coding/ultraplan-go/internal/sprint/doc.go`
- `/home/antonioborgerees/coding/ultraplan-go/internal/sprint/domain.go`
- `/home/antonioborgerees/coding/ultraplan-go/internal/sprint/discovery.go`
- `/home/antonioborgerees/coding/ultraplan-go/internal/sprint/artifacts.go`
- `/home/antonioborgerees/coding/ultraplan-go/internal/sprint/state.go`
- `/home/antonioborgerees/coding/ultraplan-go/internal/sprint/store_fs.go`
- `/home/antonioborgerees/coding/ultraplan-go/internal/sprint/service.go`
- `/home/antonioborgerees/coding/ultraplan-go/internal/sprint/sprint_test.go`
- `/home/antonioborgerees/coding/ultraplan-go/internal/app/sprint_commands.go`
- `/home/antonioborgerees/coding/ultraplan-go/internal/app/sprint_commands_test.go`
- `/home/antonioborgerees/coding/ultraplan-go/internal/app/app.go`
- `/home/antonioborgerees/coding/ultraplan-go/internal/app/app_test.go`
- `projects/ultraplan-go/sprints/17-artifact-domain-and-flow-state/plan.md`
- `projects/ultraplan-go/sprints/17-artifact-domain-and-flow-state/review.md`
- `projects/ultraplan-go/sprints/17-artifact-domain-and-flow-state/.run-state.json`

## Verification

| Check | Result |
| --- | --- |
| `go test ./internal/sprint ./internal/app` | Passed |
| `go test ./...` | Passed |
| `go test -race ./...` | Passed |
| `go build ./cmd/ultraplan` | Passed |
| `./ultraplan sprint --help` | Passed; help shows only the sprint status command and non-goals. |

## Architecture Review

- `internal/sprint` owns sprint stage ordering, artifact paths, flow-state schema validation, status derivation, and state persistence.
- `internal/app` parses arguments, calls `sprint.Service`, renders text, and maps errors. Sprint business rules are not implemented in command handlers.
- `internal/sprint` does not import `internal/study`, runtime, OpenCode, agentwrap, shell, Git, or network packages.
- No global `internal/planning`, `internal/workflow`, `internal/stages`, `internal/validation`, `internal/reports`, or `internal/prompts` package was introduced.
- Persisted artifact paths are workspace-relative and validated as contained by the selected sprint root before acceptance.

## Non-Goals Confirmed

- No `sprint validate`, `sprint prompt`, `sprint flow`, JSON status output, implementation execution, smoke investigation, automated review, issue tracking, project-index mutation, or Git mutation command was added.
- Unsupported `implementation` and `review` stages are not accepted by strict flow-state validation.
- Status refresh inspects only expected local planning artifacts and writes only `flow-state.json`.

## Deviations And Deferrals

- No plan deviations were required.
- Tests simulate write preservation through a pre-rename failure hook and validate prior-state preservation. Direct OS-level flush/close failures are not exhaustively forced because Go file descriptor failure injection would require a broader filesystem abstraction than this sprint selected.
- The no-template detector remains intentionally narrow and conservative; full sprint-index semantic validation remains deferred to a later validation sprint.

## Review Result

Sprint 17 is complete. Required implementation scope, tests, verification commands, non-goals, and review evidence are satisfied.
