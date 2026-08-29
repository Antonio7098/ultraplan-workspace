# Sprint Smoke

Smoke status: `completed`
Verdict: `pass`
Date: `2026-08-25T18:38:00Z`

## Smoke Context

Project: `ultraplan-go`
Sprint: `37-evidence-qa-smoke`
Artifact: `projects/ultraplan-go/sprints/37-evidence-qa-smoke/smoke.md`

## Review Gate

Review verdict: `pass_with_findings`
Review fingerprint: `b6a384c03156477c9ce399aa4fa0f354351a5c1a0f58437752682f9a53bd55a6`
Diagnostic override: `false`
Override rationale: none

## Harness And Protocol

Harness: `ultraplan-go-smoke`
Protocol: `1.0`

## Smoke Authoring

Author run ID: `manual-reconciliation-20260825T183800Z`
Author model: `manual`
Changed harness paths:
- `src/tests/sprint-37-evidence-qa-smoke.ts`

The reconciliation corrected two unsupported assertions. The parity probe now compares before-and-after identities instead of treating an existing governed artifact as a dry-run write. The browser probe now uses the implemented `/smoke-suite` route and expects state-dependent focused queries to fail closed when its fixture deliberately has no QA or verification state.

## Coverage Mapping

Complete: `true`
Required coverage: `AC-37-01-writable-admission-fail-closed`, `AC-37-09-smoke-parity-single-authority`, `AC-37-14-scope-exclusions`, `AC-37-08-state-versioning-fencing`, `AC-37-03-target-immutability`, `AC-37-04-frozen-evidence-plans`, `AC-37-06-adjudication-authority`, `AC-37-07-canonical-assessment-report`, `AC-37-10-cross-surface-agreement`, `AC-37-12-cancellation-recovery-truthful`, `AC-37-13-browser-security-nojs`, `AC-37-15-gated-real-runtime-dogfood`
Rationale: the complete `sprint-37` containing suite builds the recorded Sprint 37 implementation worktree and tests the real CLI and local HTTP boundaries. Deterministic product tests retain ownership of isolation mechanics, adjudication rules, state publication, and cancellation races.

Tests:
- `sprint-37-live-qa-cli-suite-contract-and-yes-gate`: `passed`
- `sprint-37-live-writable-admission-fail-closed-fixture`: `passed`
- `sprint-37-live-dry-run-map-stability-target-immutability`: `passed`
- `sprint-37-live-smoke-parity-single-authority`: `passed`
- `sprint-37-live-cross-surface-focused-routes-and-no-javascript-snapshot`: `passed`
- `sprint-37-gated-real-runtime-dogfood-blocker-evidence`: `passed`

## Selected Scope And Rationale

Scope kind: `suite`
Scope: `sprint-37`
Rationale: this is the required non-diagnostic containing suite for Sprint 37.
Duration class: `long`
Cost class: `metered-runtime`
Diagnostic only: `false`

## Preconditions And Environment

Prerequisites: none
Environment: bounded allowlist; values not persisted
Evidence roots: `/home/antonioborgerees/coding/ultraplan/ultraplan-go-smoke/runs, /home/antonioborgerees/coding/ultraplan/ultraplan-go-smoke/issues`
Effective timeout: `30m0s` (source `manifest`)

## Safe Invocation

Argv: `npm run smoke -- --workspace [WORKSPACE] --ultraplan [TARGET] --suite sprint-37 --json`

## Run Evidence

Run ID: `run-WKMV6h36jg`
Total: `6`
Passed: `6`
Failed: `0`
Skipped: `0`
Errors: `0`
Duration: `14.511s`
Runtime: `local-go`
Model: `none`

### External Evidence Identity And Links

- `run` `/home/antonioborgerees/coding/ultraplan/ultraplan-go-smoke/runs/run-WKMV6h36jg.json` sha256 `e7379e435678592e915fd5e933ef11fdda45f7420c0fef4bd16360b359ade674`
- `summary` `/home/antonioborgerees/coding/ultraplan/ultraplan-go-smoke/runs/run-WKMV6h36jg-summary.md` sha256 `b4ed5028f60d9dbf4694ab6ebdbada78066b4d5f5606210013de3182089b0a15`

## Findings

No current smoke findings.

The earlier failed run `run-eo-sf7iKtE` is retained as superseded history. Its two findings were unsupported harness assertions, not implementation failures. The corrected containing-suite run passed both probes and all four unaffected tests.

## Open Issues

- none

## Resolved Issues

- `runtime-sprint-37-live-smoke-parity-single-authority`: superseded by the corrected before-and-after identity check and passing containing-suite evidence.
- `runtime-sprint-37-live-cross-surface-focused-routes-and-no-javascript-snapshot`: superseded by truthful missing-state expectations, the implemented `/smoke-suite` route, and passing containing-suite evidence.

The original issue files remain in the external evidence root as historical records.

## Mutation And Safety Check

The reconciliation changed only the manifest-declared harness authoring file, this summary, and smoke flow state. The complete suite confirmed unchanged implementation identity across its read-only probes and successful cleanup of temporary workspaces.

## Verdict And Next Action

Verdict: `pass`
Next action: Continue to the next governed sprint gate. Sprint 37's deferred real-runtime dogfood remains recorded as a prerequisite-bound deferred task, not fabricated evidence.
