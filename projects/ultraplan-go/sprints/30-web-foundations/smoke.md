# Sprint Smoke

Smoke status: `completed`
Verdict: `pass`
Date: `2026-08-16T13:55:04Z`

## Smoke Context

Project: `ultraplan-go`
Sprint: `30-web-foundations`
Artifact: `projects/ultraplan-go/sprints/30-web-foundations/smoke.md`

## Review Gate

Review verdict: `pass_with_findings`
Review fingerprint: `8eed35875b03a4c9bd2823e290f64a6487108ca6d8647648efbec1c19caceacb`
Diagnostic override: `false`
Override rationale: none

## Harness And Protocol

Harness: `ultraplan-go-smoke`
Protocol: `1.0`

## Smoke Authoring

Author run ID: `manual-codex-2026-08-16-sprint-30`
Author model: `openai/gpt-5.6-sol`
Changed harness paths:
- `package-lock.json`
- `package.json`
- `src/cli.ts`
- `src/protocol.ts`
- `src/scenarios/workspace-factory.ts`
- `src/tests/sprint-30-web-foundations.ts`
- `tsconfig.json`

## Selected Scope And Rationale

Scope kind: `suite`
Scope: `sprint-30`
Rationale: agent-authored real-boundary suite covers the built loopback server, real workspace HTTP/security behavior, Chromium rendering, refresh, no-JavaScript operation, and graceful shutdown beyond deterministic Go tests
Duration class: `medium`
Cost class: `local`
Diagnostic only: `false`

## Preconditions And Environment

Prerequisites: none
Environment: bounded allowlist; values not persisted
Evidence roots: `/home/antonioborgerees/coding/ultraplan/ultraplan-go-smoke/runs, /home/antonioborgerees/coding/ultraplan/ultraplan-go-smoke/issues`
Effective timeout: `30m0s` (source `manifest`)

## Safe Invocation

Argv: `"tsx" "[ARG]" "[ARG]" "--project" "[ARG]" "--sprint" "[ARG]" "--workspace" "[ARG]" "--target" "[ARG]" "--scope-kind" "[ARG]" "--scope" "[ARG]"`

## Run Evidence

Run ID: `run-O3OrwQhAkj`
Total: `3`
Passed: `3`
Failed: `0`
Skipped: `0`
Errors: `0`
Duration: `15.692s`
Runtime: `local-browser`
Model: `none`

Executed tests:
- `sprint-30-live-http-lifecycle`: `passed`
- `sprint-30-live-security-boundary`: `passed`
- `sprint-30-real-browser-rendering`: `passed`

### External Evidence Identity And Links

- `run` `/home/antonioborgerees/coding/ultraplan/ultraplan-go-smoke/runs/run-O3OrwQhAkj.json` sha256 `89f98c415a497bf0785862e630a758df31a3b4f0c6c110395614f583aada92ca` size `8286` modified `2026-08-16T13:34:01Z`
- `summary` `/home/antonioborgerees/coding/ultraplan/ultraplan-go-smoke/runs/run-O3OrwQhAkj-summary.md` sha256 `3e7baa3fdc13b41bc49df2fa915db90f94847ab104c4262706aa3758c88cfe8e` size `956` modified `2026-08-16T13:34:01Z`

## Findings

- none

## Open Issues

- none

## Resolved Issues

- none

## Mutation And Safety Check

Only smoke.md, flow-state.json, manifest-declared harness authoring paths, and manifest-declared external evidence roots were approved for mutation. Product source and governed sprint inputs were identity-checked before and after authoring.

## Verdict And Next Action

Verdict: `pass`
Next action: Deep smoke is complete; proceed to the next roadmap stage.
