# Sprint Smoke

Smoke status: `completed`
Verdict: `pass`
Date: `2026-08-27T13:12:35Z`

## Smoke Context

Project: `ultraplan-go`
Sprint: `38-bounded-repair`
Artifact: `projects/ultraplan-go/sprints/38-bounded-repair/smoke.md`

## Review Gate

Review verdict: `pass`
Review fingerprint: `cf3f60b35f9477b5b9f473814def0e9aa6cfaa1bcb2485ca5fb3b54a08aaa302`
Diagnostic override: `false`
Override rationale: none

## Harness And Protocol

Harness: `ultraplan-go-smoke`
Protocol: `1.0`

## Smoke Authoring

Author run ID: `manual-reconciliation-20260827T131235Z`
Author model: `manual-reconciliation`
Changed harness paths:
- `src/tests/sprint-38-bounded-repair.ts`

## Coverage Mapping

Complete: `true`
Required coverage: `AC-38-01`, `AC-38-02`, `AC-38-03`, `AC-38-04`, `AC-38-06`, `AC-38-07`, `AC-38-08`, `AC-38-09`, `AC-38-10`, `AC-38-11`
Rationale: The complete Sprint 38 mapping selects six non-empty real-boundary tests for CLI proof gates, deterministic status, cancellation safety, HTTP parity, effective-source reporting, and truthful deferred dogfood evidence.

## Selected Scope And Rationale

Scope kind: `suite`
Scope: `sprint-38`
Rationale: Run the complete mapped Sprint 38 external suite after repairing the five failures from `run-AGjlQw3uXv`.
Duration class: `long`
Cost class: `metered-runtime`
Diagnostic only: `false`

## Preconditions And Environment

Prerequisites: none
Environment: bounded allowlist; values not persisted
Evidence roots: `/home/antonioborgerees/coding/ultraplan/ultraplan-go-smoke/runs, /home/antonioborgerees/coding/ultraplan/ultraplan-go-smoke/issues`
Effective timeout: `30m0s` (source `manifest`)

## Safe Invocation

Argv: `"tsx" "src/cli.ts" "smoke" "--workspace" "[ARG]" "--ultraplan" "[ARG]" "--tests" "[ARG]" "--json"`

## Run Evidence

Run ID: `run-lLWhgRMowA`
Total: `6`
Passed: `6`
Failed: `0`
Skipped: `0`
Errors: `0`
Duration: `8.670s`
Runtime: `opencode`
Model: `minimax-coding-plan/MiniMax-M3`

Executed tests:
- `sprint-38-live-repair-cli-contract-and-manual-only-enforcement`: `passed`
- `sprint-38-live-repair-status-stale-and-target-immutability`: `passed`
- `sprint-38-live-repair-cancel-fails-closed-for-unknown-run`: `passed`
- `sprint-38-live-repair-cross-surface-status-agreement`: `passed`
- `sprint-38-live-repair-config-effective-sources-and-lower-only`: `passed`
- `sprint-38-gated-real-runtime-dogfood-blocker-evidence`: `passed`

### External Evidence Identity And Links

- `run` `/home/antonioborgerees/coding/ultraplan/ultraplan-go-smoke/runs/run-lLWhgRMowA.json` sha256 `6f8ffd3ea258049c5240ad6edc7c3a3335f1ad6a965391409cbd99106fea5a79` size `49945` modified `2026-08-27T13:12:35Z`
- `summary` `/home/antonioborgerees/coding/ultraplan/ultraplan-go-smoke/runs/run-lLWhgRMowA-summary.md` sha256 `6e9a94775a1d7daa6438c7a154da6b81dc5548c65cae5be8469c8c7f86838617` size `1368` modified `2026-08-27T13:12:35Z`

## Findings

- none

The failed `run-AGjlQw3uXv` result remains as superseded history. Its five findings were repaired or corrected as smoke-fixture defects before the passing containing-suite run.

## Open Issues

- none

## Resolved Issues

- `runtime-sprint-38-live-repair-cli-contract-and-manual-only-enforcement`
- `runtime-sprint-38-live-repair-status-stale-and-target-immutability`
- `runtime-sprint-38-live-repair-cross-surface-status-agreement`
- `runtime-sprint-38-live-repair-config-effective-sources-and-lower-only`
- `runtime-sprint-38-gated-real-runtime-dogfood-blocker-evidence`

## Mutation And Safety Check

The run used temporary fixture workspaces. Cleanup completed and no prohibited target mutation occurred. Raw output remains in the external run record.

## Verdict And Next Action

Verdict: `pass`
Next action: No smoke action required.
