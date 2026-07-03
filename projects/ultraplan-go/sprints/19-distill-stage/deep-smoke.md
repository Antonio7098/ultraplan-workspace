# Deep Smoke: 19-distill-stage

> Project: `ultraplan-go`
> Sprint: `19-distill-stage`
> Run ID: `run-iC9yH7B62t`
> Runtime: `opencode/big-pickle`
> Date: 2026-06-15

## Scope

- Sprint: `19-distill-stage`
- Implementation directory: `/home/antonioborgerees/coding/ultraplan-go`
- Smoke harness: `/home/antonioborgerees/coding/ultraplan-go-smoke/`
- Selected smoke level/tests: `--level sprint-19` (18 tests)
- Reason scope is sufficient: Sprint 19 delivers `validate technical-handbook`, `prompt technical-handbook`, and `flow --to technical-handbook` (including `--dry-run` and non-dry-run). The sprint-19 smoke suite covers CLI surface (help, unsupported stages), validate semantic checks (valid, missing, empty, placeholder, missing sections), prompt rendering (runtime-free, ANSI-free, required constraints), flow dry-run (no runtime, no writes, ANSI-free), flow non-dry-run failure path (runtime unavailable, flow-state failure recording, no later-stage completion), and real OpenCode runtime invocation through `flow --to technical-handbook`.

## Run Evidence

| Run ID | Command | Result | Summary |
|---|---|---|---|
| `run-r1Oqx1MsZC` | `OPENCODE_MODEL=minimax-coding-plan/MiniMax-M2.7 node --import tsx/esm src/cli.ts smoke --level sprint-19 --ultraplan ~/coding/ultraplan-go --workspace ~/coding/ultraplan-go-smoke` | 17/17 passed (272ms) | All CLI surface, validate, prompt, flow dry-run, flow non-dry-run failure, and unsupported-stage tests pass |
| `run-4u1DmdYEu6` | `OPENCODE_MODEL=opencode/big-pickle node --import tsx/esm src/cli.ts smoke --level sprint-19 --ultraplan ~/coding/ultraplan-go --workspace ~/coding/ultraplan-go-smoke` | 17/17 passed | Same 17-test suite passed with big-pickle selected, before real-runtime D4 coverage was added |
| manual | `~/coding/ultraplan-go/cmd/ultraplan/ultraplan sprint ultraplan-go 19 flow --to technical-handbook` from `/tmp/test-ws` after rebuilding `cmd/ultraplan/ultraplan` | exit 0 | Real runtime flow completed: `technical-handbook complete`; stages showed requirements, sprint-index, and technical-handbook complete with area-reasoning ready |
| `run-iC9yH7B62t` | `OPENCODE_MODEL=opencode/big-pickle npx tsx src/cli.ts smoke --level sprint-19 --ultraplan ~/coding/ultraplan-go --workspace ~/coding/ultraplan-go-smoke` | 18/18 passed (42.1s) | Adds D4 real-runtime coverage. D4 exited 0, wrote a valid handbook, and marked technical-handbook complete |

## Findings

| ID | Severity | Status | Evidence | Required Action |
|---|---|---|---|---|
| sprint-19-smoke-added | info | resolved | 18/18 tests pass in `run-iC9yH7B62t` | None — smoke harness now covers sprint 19 distill-stage surface and real runtime invocation |
| stale-cli-binary | medium | resolved | Before rebuilding, non-dry-run `flow --to technical-handbook` returned `sprint.flow: runtime is required for technical-handbook flow`; after `go build -o cmd/ultraplan/ultraplan ./cmd/ultraplan`, the same command completed with exit 0 | Rebuild `cmd/ultraplan/ultraplan` before runtime-facing smoke runs |

## Open Issues

None. All 18 smoke tests pass with `opencode/big-pickle`.

## Resolved Issues

- Sprint 19 smoke harness coverage was absent. Added `sprint-19-distill.ts` with 18 tests covering CLI surface, validate, prompt, flow dry-run, flow non-dry-run failure path, real runtime invocation, and unsupported-stage rejection. All tests pass in `run-iC9yH7B62t`.
- The apparent nil-runtime failure was caused by running a stale `cmd/ultraplan/ultraplan` binary. Rebuilding from current source fixed the real-runtime path.

## Verdict

pass
