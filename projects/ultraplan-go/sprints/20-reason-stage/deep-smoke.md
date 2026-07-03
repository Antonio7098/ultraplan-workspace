# Deep Smoke: 20-reason-stage

> Project: `ultraplan-go`
> Sprint: `20-reason-stage`
> Run ID: `run-lfUvHxbGeC`
> Runtime: `opencode/big-pickle`
> Date: 2026-06-15

## Scope

- Sprint: `20-reason-stage`
- Implementation directory: `/home/antonioborgerees/coding/ultraplan-go`
- Smoke harness: `/home/antonioborgerees/coding/ultraplan-go-smoke/`
- Selected smoke level/tests: `--level sprint-20` (12 tests)
- Reason scope is sufficient: Sprint 20 delivers `validate area-reasoning`, `validate reasoning`, `prompt area-reasoning`, `prompt reasoning`, `flow --to reasoning --dry-run`, non-dry-run `flow --to reasoning`, selected-area prerequisite handling, flow-state updates through `reasoning`, and unsupported future-stage rejection. The new sprint-20 smoke suite covers CLI help, valid and invalid area/final validation, runtime-free prompt previews, dry-run mutation safety, dry-run selected-area prerequisite failure, missing-runtime failure state recording, real OpenCode `flow --to reasoning`, and unsupported `smoke` flow target rejection.

## Run Evidence

| Run ID | Command | Result | Summary |
|---|---|---|---|
| `run-EyqAxqP0V3` | `OPENCODE_MODEL=opencode/big-pickle npx tsx src/cli.ts smoke --level sprint-20 --ultraplan ~/coding/ultraplan-go --workspace ~/coding/ultraplan-go-smoke` | 2/12 passed | Harness fixture bug: `sprint-index.md` omitted required excluded contexts, so product validation failed before reasoning assertions. |
| `run-t5YLAFeqyA` | `OPENCODE_MODEL=opencode/big-pickle npx tsx src/cli.ts smoke --level sprint-20 --ultraplan ~/coding/ultraplan-go --workspace ~/coding/ultraplan-go-smoke` | 11/12 passed | Fixture fixed; D2 assertion expected a lower-level missing-area phrase while CLI rendered `selected area reasoning failed validation`. D4 real-runtime flow passed. |
| `run-lfUvHxbGeC` | `OPENCODE_MODEL=opencode/big-pickle npx tsx src/cli.ts smoke --level sprint-20 --ultraplan ~/coding/ultraplan-go --workspace ~/coding/ultraplan-go-smoke` | 12/12 passed (47.7s) | All sprint-20 smoke probes passed. Real OpenCode D4 exited 0: `reasoning complete`; stages marked requirements, sprint-index, technical-handbook, area-reasoning, and reasoning complete while `plan` remained missing. |

## Findings

| ID | Severity | Status | Evidence | Required Action |
|---|---|---|---|---|
| sprint-20-smoke-added | info | resolved | Added `sprint-20` smoke suite with 12 tests; `run-lfUvHxbGeC` passed 12/12. | None. |
| sprint-20-fixture-excluded-context | info | resolved | `run-EyqAxqP0V3` failed because the smoke fixture omitted required excluded contexts (`implementation`, `review`, `issue`, `git`); after adding rows, those failures disappeared in `run-t5YLAFeqyA`. | None. Harness fixture corrected. |
| sprint-20-d2-assertion-specificity | info | resolved | `run-t5YLAFeqyA` failed only D2 because the test expected `missing selected area reasoning`, while the CLI correctly surfaced `selected area reasoning failed validation`; `run-lfUvHxbGeC` passed after asserting the CLI-level diagnostic. | None. Harness assertion corrected. |

## Open Issues

None. The selected sprint-20 smoke scope passed with real `opencode/big-pickle` runtime evidence.

## Resolved Issues

- Added `src/tests/sprint-20-reason.ts` and registered `--level sprint-20` in the smoke harness.
- Rebuilt `/home/antonioborgerees/coding/ultraplan-go/cmd/ultraplan/ultraplan` from current source before running smoke evidence.
- Corrected the sprint-20 harness fixture to satisfy selected-context validation and corrected the D2 assertion to match CLI-level diagnostics.

## Verdict

pass
