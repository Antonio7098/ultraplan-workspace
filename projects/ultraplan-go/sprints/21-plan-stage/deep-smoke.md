# Deep Smoke: 21-plan-stage

> Project: `ultraplan-go`
> Sprint: `21-plan-stage`
> Run ID: `run-M9jT1FR566`
> Runtime: `opencode/big-pickle`
> Date: 2026-06-15

## Scope

- Sprint: `21-plan-stage`
- Implementation directory: `/home/antonioborgerees/coding/ultraplan-go`
- Smoke harness: `/home/antonioborgerees/coding/ultraplan-go-smoke/`
- Selected smoke level/tests: `--level sprint-21` (14 tests)
- Reason scope is sufficient: Sprint 21 delivers `validate plan`, `prompt plan`, `flow --to plan --dry-run`, non-dry-run `flow --to plan`, reasoning prerequisite handling, missing-runtime failure state recording, real OpenCode `flow --to plan`, and unsupported future-stage rejection. The new sprint-21 smoke suite covers CLI help, valid and invalid plan validation (missing sections, missing reasoning, untraced tasks, forbidden content), runtime-free prompt previews with decision traceability, dry-run mutation safety, dry-run reasoning prerequisite failure, missing-runtime failure state recording, real OpenCode `flow --to plan`, and unsupported `smoke` and `implementation` flow target rejection.

## Run Evidence

| Run ID | Command | Result | Summary |
|---|---|---|---|
| `run-Lh24HJC7fQ` | `OPENCODE_MODEL=opencode/big-pickle npx tsx src/cli.ts smoke --level sprint-21 --ultraplan ~/coding/ultraplan-go --workspace ~/coding/ultraplan-go-smoke` | 13/14 passed | D2 failed: test asserted `missing final reasoning` but CLI rendered `plan prerequisites failed validation`. |
| `run-M9jT1FR566` | `OPENCODE_MODEL=opencode/big-pickle npx tsx src/cli.ts smoke --level sprint-21 --ultraplan ~/coding/ultraplan-go --workspace ~/coding/ultraplan-go-smoke` | 14/14 passed (35.8s) | All sprint-21 smoke probes passed. Real OpenCode D4 exited 0: `plan complete`; stages marked requirements, sprint-index, technical-handbook, area-reasoning, reasoning, and plan complete. |

## Findings

| ID | Severity | Status | Evidence | Required Action |
|---|---|---|---|---|
| sprint-21-smoke-added | info | resolved | Added `sprint-21` smoke suite with 14 tests; `run-M9jT1FR566` passed 14/14. | None. |
| sprint-21-d2-assertion-specificity | info | resolved | `run-Lh24HJC7fQ` failed only D2 because the test expected `missing final reasoning` while the CLI correctly surfaced `plan prerequisites failed validation`; `run-M9jT1FR566` passed after asserting the CLI-level diagnostic. | None. Harness assertion corrected. |

## Open Issues

None. The selected sprint-21 smoke scope passed with real `opencode/big-pickle` runtime evidence.

## Resolved Issues

- Added `src/tests/sprint-21-plan.ts` and registered `--level sprint-21` in the smoke harness.
- Rebuilt `/home/antonioborgerees/coding/ultraplan-go/cmd/ultraplan/ultraplan` from current source before running smoke evidence.
- Corrected the D2 assertion to match the CLI-level diagnostic string.

## Verdict

pass
