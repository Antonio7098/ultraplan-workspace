# Sprint Smoke

Smoke status: `completed`
Verdict: `pass`
Date: `2026-08-20T11:28:51Z`

## Smoke Context

Project: `ultraplan-go`
Sprint: `33-code-context-stage`
Artifact: `projects/ultraplan-go/sprints/33-code-context-stage/smoke.md`

## Review Gate

Review verdict: `pass`
Review fingerprint: `9796f51fbc4076db878ae1a17e7bba99eb4b36b594d2a882f27190178d11fc1d`
Diagnostic override: `false`
Override rationale: none

## Harness And Protocol

Harness: `ultraplan-go-smoke`
Protocol: `1.0`

## Smoke Authoring

Author run ID: `opencode-3`
Author model: `minimax-coding-plan/MiniMax-M3`
Changed harness paths:
- none; existing traceable suite retained after agent inspection

## Selected Scope And Rationale

Scope kind: `suite`
Scope: `sprint-33`
Rationale: dedicated real-target suite exercises code-context domain/runtime atomicity, CLI/app safety, and generic web integration without claiming the Sprint 34 live-provider dogfood gate
Duration class: `long`
Cost class: `metered-runtime`
Diagnostic only: `false`

## Preconditions And Environment

Prerequisites: none
Environment: bounded allowlist; values not persisted
Evidence roots: `/home/antonioborgerees/coding/ultraplan/ultraplan-go-smoke/runs, /home/antonioborgerees/coding/ultraplan/ultraplan-go-smoke/issues`
Effective timeout: `30m0s` (source `manifest`)

## Safe Invocation

Argv: `"cli.mjs" "[ARG]" "[ARG]" "--project" "[ARG]" "--sprint" "[ARG]" "--workspace" "[ARG]" "--target" "[ARG]" "--scope-kind" "[ARG]" "--scope" "[ARG]"`

## Run Evidence

Run ID: `run-rdjxCagl7K`
Total: `3`
Passed: `3`
Failed: `0`
Skipped: `0`
Errors: `0`
Duration: `1.603s`
Runtime: `local-go`
Model: `none`

Executed tests:
- `sprint-33-domain-runtime-and-atomicity`: `passed`
- `sprint-33-cli-app-and-safe-preview`: `passed`
- `sprint-33-generic-web-contract`: `passed`

### External Evidence Identity And Links

- `run` `/home/antonioborgerees/coding/ultraplan/ultraplan-go-smoke/runs/run-rdjxCagl7K.json` sha256 `2dcda9a7399e46ea2ba1b67561c49cab35a4bc2643b54435c994bda5c06aaaa2` size `3542` modified `2026-08-20T11:28:50Z`
- `summary` `/home/antonioborgerees/coding/ultraplan/ultraplan-go-smoke/runs/run-rdjxCagl7K-summary.md` sha256 `cd5dc2ea5b6249fb13433edde4b38b480582e7aaf36be3bf61d1ada2547e11a9` size `964` modified `2026-08-20T11:28:50Z`

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
