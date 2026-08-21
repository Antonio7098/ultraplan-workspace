# Sprint Smoke

Smoke status: `completed`
Verdict: `blocked`
Date: `2026-08-21T08:14:24Z`

## Smoke Context

Project: `ultraplan-go`
Sprint: `34-shared-context`
Artifact: `projects/ultraplan-go/sprints/34-shared-context/smoke.md`

## Review Gate

Review verdict: `pass`
Review fingerprint: `45324b45bc251e8ba41ec2750c21ca64fd2e2bc024f1bb5c2e73d2ccf8091e0c`
Diagnostic override: `false`
Override rationale: none

## Harness And Protocol

Harness: `ultraplan-go-smoke`
Protocol: `1.0`

## Smoke Authoring

Author run ID: `opencode-2`
Author model: `minimax-coding-plan/MiniMax-M3`
Changed harness paths:
- `src/tests/sprint-34-shared-context.ts`

## Selected Scope And Rationale

Scope kind: `none`
Scope: `none`
Rationale: discovery cannot prove enumerated deep-smoke coverage for this sprint
Duration class: `long`
Cost class: `metered-runtime`
Diagnostic only: `false`

## Preconditions And Environment

Prerequisites: none
Environment: bounded allowlist; values not persisted
Evidence roots: `/home/antonioborgerees/coding/ultraplan/ultraplan-go-smoke/runs, /home/antonioborgerees/coding/ultraplan/ultraplan-go-smoke/issues`
Effective timeout: `30m0s` (source `manifest`)

## Safe Invocation

Argv: `none`

## Run Evidence

Run ID: `none`
Total: `0`
Passed: `0`
Failed: `0`
Skipped: `0`
Errors: `0`
Duration: `0s`
Runtime: `none`
Model: `none`

### External Evidence Identity And Links

- none (preflight classification)

## Findings

- none

## Open Issues

- none

## Resolved Issues

- none

## Mutation And Safety Check

Only smoke.md, flow-state.json, manifest-declared harness authoring paths, and manifest-declared external evidence roots were approved for mutation. Product source and governed sprint inputs were identity-checked before and after authoring.

## Verdict And Next Action

Verdict: `blocked`
Next action: Update the harness with required coverage IDs and non-empty sprint-specific tests, then rerun smoke.
