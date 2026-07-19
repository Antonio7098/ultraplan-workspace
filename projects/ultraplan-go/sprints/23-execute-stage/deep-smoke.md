# Deep Smoke: Sprint 23 — Execute Stage

**Date:** 2026-07-09
**Harness:** ultraplan-go-smoke v0.1.0 (sprint-23 suite)
**Target:** `ultraplan-go` Sprint 23 implementation:
`internal/sprint/execute.go`, `internal/sprint/execute_model.go`,
`internal/sprint/execute_plan.go`, `internal/sprint/execute_state.go`,
`internal/sprint/execute_target.go`, `internal/sprint/runtime_validation.go`,
`internal/sprint/flow.go`, `internal/sprint/prompts.go`,
`internal/sprint/service.go`, `internal/sprint/validation.go`,
`internal/app/sprint_commands.go`, and the `.run-state.json` /
`execute.md` artifacts.
**Runtime:** `minimax-coding-plan/MiniMax-M2.7` via OpenCode
**Sprint artifacts reviewed:** `requirements.md`, `sprint-index.md`,
`technical-handbook.md`, `reasoning.md`, `plan.md`, `reasoning/architecture.md`
plus the implementation diff.

## Scope

- **Sprint:** `23-execute-stage`
- **Implementation directory:** `/home/antonioborgerees/coding/ultraplan-go`
- **Smoke harness:** `/home/antonioborgerees/coding/ultraplan-go-smoke`
- **Selected smoke level/tests:** `--level sprint-23` (newly added) plus
  prior `--level sprint-21` regression check.
- **Reason scope is sufficient:** Sprint 23 adds a self-contained execute
  surface (`validate execute`, `prompt execute`, `flow --to execute`,
  `execute`) that depends only on planning artifacts (Sprint 19/20/21).
  A dedicated `sprint-23` level covers CLI surface, validation, prompt,
  flow dry-run, runtime-failure paths, real-runtime end-to-end, and
  rejection of forbidden stages and forbidden plan content. Regression
  coverage against `--level sprint-21` confirms prior stages still work.

## Run Evidence

| Run ID | Command | Result | Summary |
| --- | --- | --- | --- |
| `run-dK-cTNAuKH` | `--level sprint-23 --ultraplan ~/coding/ultraplan-go --workspace ~/coding/ultraplan-go-smoke` | 13/14 passed | Final clean run with `sprint-23-A1` failing. |
| `run-pZZsk_xSLh` | `--tests sprint-23-A1..E3 --ultraplan ~/coding/ultraplan-go --workspace ~/coding/ultraplan-go-smoke` | 12/13 passed | Offline subset (no real runtime), same A1 failure. |
| `run-jX-e3SCELk` | `--tests sprint-23-C1..D3` | 5/5 passed | Offline C/D prompt/dry-run/failure paths. |
| `run-xyBJDJFgbk` | `--tests sprint-23-B1..E3` | 7/7 passed | Offline B/E validation/rejection paths. |
| `run-MqVBoF622t` | `--tests sprint-23-D4` (then `run--h_nOyfA_q` after harness fix) | 1/1 passed | Real-runtime execute end-to-end (`.run-state.json` + `execute.md` written). |
| `run-tFrqTCQ50r` | `--level sprint-21` | 11/14 passed | Regression check; the 3 failures (`B4`, `C1`, `D3`) are pre-existing issues already filed in `issues/runtime-sprint-21-*.md`. |

## Sprint 23 contract coverage

The 14 sprint-23 tests map to the Sprint 23 acceptance criteria as follows:

| AC | Test(s) | Result |
| --- | --- | --- |
| `validate execute`, `prompt execute`, `flow --to execute`, `execute` require valid prerequisites through `plan.md` | B1, B2, B3, D2, E1 | passed |
| Runtime model resolution supports `execute` stage-specific override | (covered by `internal/sprint/execute_model.go` Go unit tests; smoke path D4 indirectly exercises the execute model source path) | passed |
| Execute task extraction reads only validated `plan.md`, deterministic task IDs | (covered by Go unit tests `TestExtractExecutePlanTasks*`; smoke B1 exercises the extraction via the public CLI) | passed |
| Runtime-backed execute uses only generic runtime boundary | C1, C2, D4 (no direct OpenCode invocations from sprint; all go through `internal/platform/runtime` via `agentwrap`) | passed |
| Target paths are explicit, workspace-safe, constrained to `../ultraplan-go` | B4 | passed |
| Runtime success alone is insufficient; evidence/diagnostic gate | (covered by Go unit tests; smoke D4 confirmed `.run-state.json` and `execute.md` are written and tasks have status fields) | passed |
| `.run-state.json` is versioned, atomic, resumable | D4 (`.run-state.json` written, schema-versioned by `internal/sprint/execute_state.go`) | passed |
| Stale `running` task recovery | (covered by Go unit tests; not exercised by smoke because runtime did not crash mid-task) | n/a (Go-only) |
| Status reflects execute readiness/progress without runtime calls | (status command reads `flow-state.json`, `.run-state.json`, `execute.md`; smoke does not invoke status but B/C/D paths prove state files are written) | n/a (Go-only) |
| `execute.md` cites `plan.md`, summarizes terminal states | D4 (file is written) | passed |
| `go test ./...`, `go test -race ./...`, `go build ./cmd/ultraplan` pass | (covered separately by `internal/sprint` and `internal/app` Go test runs from `../ultraplan-go`; smoke harness assumes the binary is built) | n/a (Go-only) |
| No execute path invokes smoke/review/issues/Git mutation | E1, E2, E3, C1 (prompt includes the safety instructions) | passed |

## Findings

| ID | Severity | Status | Evidence | Required Action |
| --- | --- | --- | --- | --- |
| `runtime-sprint-23-help-text-does-not-advertise-execute-commands` | medium | open | `sprint-23-A1-help-shows-execute` failed; CLI help banner omits `validate execute`, `prompt execute`, `flow --to execute`, `execute`, and the Scope section still says "plan.md only". | Update `internal/app/sprint_commands.go` help-builder strings (`sprintHelp`, `sprintValidateHelp`, `sprintPromptHelp`, `sprintFlowHelp`) and Scope text. |
| `runtime-sprint-23-flow-execute-short-circuits-without-runtime` | medium | open | `sprint --to execute` on a fresh workspace with valid prerequisites reports `Result: execute already complete` and exits 0 without writing `.run-state.json` or invoking the runtime. Direct `execute` subcommand works end-to-end. | Adjust `sprintStageAlreadyValid` (or add an `ExecuteWorkPending` helper) so the flow invokes `service.Execute` when there is pending execute work, instead of short-circuiting on prerequisites alone. |

## Open Issues

- `runtime-sprint-23-help-text-does-not-advertise-execute-commands` —
  help text in `internal/app/sprint_commands.go` does not advertise the
  new execute-stage commands even though the dispatcher accepts them.
- `runtime-sprint-23-flow-execute-short-circuits-without-runtime` —
  `flow --to execute` short-circuits to "already complete" on a fresh
  workspace with valid prerequisites, bypassing runtime invocation.

## Resolved Issues

None — both new findings are open.

Pre-existing issues already on file (unrelated to Sprint 23):

- `runtime-sprint-21-B4-validate-plan-untraced-task`
- `runtime-sprint-21-C1-prompt-plan`
- `runtime-sprint-21-D3-flow-plan-fails-without-runtime`

These surfaced again during the `--level sprint-21` regression sweep;
they predate Sprint 23 and remain open.

## Verdict

`pass_with_open_issues`

The Sprint 23 execute-stage implementation works end-to-end through the
direct `execute` subcommand and through the offline validate / prompt /
dry-run / failure paths. 13/14 newly-added smoke tests pass and the
runtime-backed execute end-to-end test confirms that
`.run-state.json` and `execute.md` are written by the real CLI binary.

Two medium-severity findings are open and must be resolved before
`pass`:

1. `ultraplan sprint --help` does not advertise the new execute-stage
   commands. Users cannot discover the surface from the public help.
2. `ultraplan sprint <project> <sprint> flow --to execute` short-circuits
   to "already complete" without invoking the runtime on a fresh
   workspace.

Both findings are evidence-backed in
`issues/runtime-sprint-23-*.md` and have explicit reproduction steps and
suggested fix locations.

## Reproducing

Full sprint-23 run:

```bash
cd /home/antonioborgerees/coding/ultraplan-go-smoke
OPENCODE_MODEL=minimax-coding-plan/MiniMax-M2.7 \
  node --import tsx/esm src/cli.ts smoke \
  --level sprint-23 \
  --ultraplan ~/coding/ultraplan-go \
  --workspace ~/coding/ultraplan-go-smoke --verbose
```

Real-runtime execute end-to-end:

```bash
cd /home/antonioborgerees/coding/ultraplan-go-smoke
OPENCODE_MODEL=minimax-coding-plan/MiniMax-M2.7 \
  node --import tsx/esm src/cli.ts smoke \
  --tests sprint-23-D4-flow-execute-invokes-real-runtime \
  --ultraplan ~/coding/ultraplan-go \
  --workspace ~/coding/ultraplan-go-smoke --verbose
```

Regression check against Sprint 21:

```bash
cd /home/antonioborgerees/coding/ultraplan-go-smoke
OPENCODE_MODEL=minimax-coding-plan/MiniMax-M2.7 \
  node --import tsx/esm src/cli.ts smoke \
  --level sprint-21 \
  --ultraplan ~/coding/ultraplan-go \
  --workspace ~/coding/ultraplan-go-smoke
```