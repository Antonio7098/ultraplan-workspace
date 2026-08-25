# Sprint Smoke

Smoke status: `completed`
Verdict: `pass`
Date: `2026-08-24T17:01:53Z`

## Smoke Context

Project: `ultraplan-go`
Sprint: `36-read-only-qa`
Artifact: `projects/ultraplan-go/sprints/36-read-only-qa/smoke.md`

## Review Gate

Review verdict: `pass_with_findings`
Review fingerprint: `c0fd8452a95ae0461631b81a75140d2ff66a9b37ab21849cccaa3fcc6730dd63`
Diagnostic override: `false`
Override rationale: none

## Harness And Protocol

Harness: `ultraplan-go-smoke`
Protocol: `1.0`

## Smoke Authoring

Author run ID: `opencode-1`
Author model: `openrouter/stealth/ox-alpha`
Changed harness paths:
- none; existing traceable suite retained after agent inspection

## Coverage Mapping

Complete: `true`
Required coverage: `AC-36-01-map-bytes-deterministic`, `AC-36-02-changed-path-primary-ownership`, `AC-36-03-frozen-positive-limits`, `AC-36-04-target-non-mutation`, `AC-36-05-workspace-state-isolation`, `AC-36-06-conformance-review-single-capability`, `AC-36-07-cross-surface-status-agreement`, `AC-36-08-no-javascript-snapshot`, `AC-36-09-gated-real-runtime-dogfood`
Rationale: agent-authored real-boundary suite builds the sprint-36 worktree binary and proves the boundaries deterministic Go tests cannot replace: repeated real 'qa --dry-run' runs against the governed planning workspace produce byte-stable maps with one primary owner per changed path and frozen positive limits, qa status/dry-run leave flow-state.json, the verification state tree, and target Git identity byte-identical before and after, the conformance-review alias routes to the single existing review capability with identical envelope/fingerprint/exit while qa keeps its own sprint.qa operation, CLI and live-server QA status agree on phase/next-action for a shared fixture workspace with a no-JavaScript HTML snapshot of the same canonical facts, and Task 14's gated real-runtime dogfood remains preserved as deferred blocker evidence rather than claimed coverage

- `AC-36-01-map-bytes-deterministic` — none (mapped tests: sprint-36-live-qa-dry-run-deterministic-map-and-target-identity)
- `AC-36-02-changed-path-primary-ownership` — none (mapped tests: sprint-36-live-qa-dry-run-deterministic-map-and-target-identity)
- `AC-36-03-frozen-positive-limits` — none (mapped tests: sprint-36-live-qa-dry-run-deterministic-map-and-target-identity)
- `AC-36-04-target-non-mutation` — none (mapped tests: sprint-36-live-qa-dry-run-deterministic-map-and-target-identity, sprint-36-live-qa-workspace-state-isolation)
- `AC-36-05-workspace-state-isolation` — none (mapped tests: sprint-36-live-qa-workspace-state-isolation)
- `AC-36-06-conformance-review-single-capability` — none (mapped tests: sprint-36-live-conformance-review-alias-single-capability)
- `AC-36-07-cross-surface-status-agreement` — none (mapped tests: sprint-36-live-cross-surface-agreement-and-no-javascript-snapshot)
- `AC-36-08-no-javascript-snapshot` — none (mapped tests: sprint-36-live-cross-surface-agreement-and-no-javascript-snapshot)
- `AC-36-09-gated-real-runtime-dogfood` — none (mapped tests: sprint-36-gated-real-runtime-dogfood-blocker-evidence)

Tests:
- `sprint-36-gated-real-runtime-dogfood-blocker-evidence` (suite `sprint-36`): `AC-36-09-gated-real-runtime-dogfood`
- `sprint-36-live-conformance-review-alias-single-capability` (suite `sprint-36`): `AC-36-06-conformance-review-single-capability`
- `sprint-36-live-cross-surface-agreement-and-no-javascript-snapshot` (suite `sprint-36`): `AC-36-07-cross-surface-status-agreement`, `AC-36-08-no-javascript-snapshot`
- `sprint-36-live-qa-dry-run-deterministic-map-and-target-identity` (suite `sprint-36`): `AC-36-01-map-bytes-deterministic`, `AC-36-02-changed-path-primary-ownership`, `AC-36-03-frozen-positive-limits`, `AC-36-04-target-non-mutation`
- `sprint-36-live-qa-workspace-state-isolation` (suite `sprint-36`): `AC-36-04-target-non-mutation`, `AC-36-05-workspace-state-isolation`

## Selected Scope And Rationale

Scope kind: `suite`
Scope: `sprint-36`
Rationale: explicit suite override
Duration class: `medium`
Cost class: `local`
Diagnostic only: `false`

## Preconditions And Environment

Prerequisites: none
Environment: bounded allowlist; values not persisted
Evidence roots: `/home/antonioborgerees/coding/ultraplan/ultraplan-go-smoke/runs, /home/antonioborgerees/coding/ultraplan/ultraplan-go-smoke/issues`
Effective timeout: `30m0s` (source `manifest`)

## Safe Invocation

Argv: `"cli.mjs" "[ARG]" "[ARG]" "--project" "[ARG]" "--sprint" "[ARG]" "--workspace" "[ARG]" "--target" "[ARG]" "--scope-kind" "[ARG]" "--scope" "[ARG]"`

## Run Evidence

Run ID: `run-0fT-aadOwt`
Total: `5`
Passed: `5`
Failed: `0`
Skipped: `0`
Errors: `0`
Duration: `8.014s`
Runtime: `local-go`
Model: `none`

Executed tests:
- `sprint-36-live-qa-dry-run-deterministic-map-and-target-identity`: `passed`
- `sprint-36-live-qa-workspace-state-isolation`: `passed`
- `sprint-36-live-conformance-review-alias-single-capability`: `passed`
- `sprint-36-live-cross-surface-agreement-and-no-javascript-snapshot`: `passed`
- `sprint-36-gated-real-runtime-dogfood-blocker-evidence`: `passed`

### External Evidence Identity And Links

- `run` `/home/antonioborgerees/coding/ultraplan/ultraplan-go-smoke/runs/run-0fT-aadOwt.json` sha256 `c10bdf399681fb42db7b18b2839df446fab89b336bfcb726e5797cb8399ae726` size `76297` modified `2026-08-24T17:01:53Z`
- `summary` `/home/antonioborgerees/coding/ultraplan/ultraplan-go-smoke/runs/run-0fT-aadOwt-summary.md` sha256 `16c206f153de12f4887ea5d7a4729043d4d987f7801c3a38058e66956dd005a3` size `1267` modified `2026-08-24T17:01:53Z`

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
