# Sprint Smoke

Smoke status: `completed`
Verdict: `pass`
Date: `2026-08-23T14:11:41Z`

## Smoke Context

Project: `ultraplan-go`
Sprint: `35-durable-run-observability`
Artifact: `projects/ultraplan-go/sprints/35-durable-run-observability/smoke.md`

## Review Gate

Review verdict: `pass_with_findings`
Review fingerprint: `06e6993b94d22e9f7a88db2b4f76463eb8a3213cc880862819cf22dbab4edbcc`
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
Required coverage: `AC-35-01-durable-acceptance-fail-closed`, `AC-35-03-workspace-inspection-any-surface`, `AC-35-08-liveness-reconciliation-truthful`, `AC-35-09-single-terminal-outcome`, `AC-35-04-retained-history-other-server`, `AC-35-06-monotonic-sequence-replay`, `AC-35-07-typed-replay-gap-cursor-ahead`, `AC-35-14-bounded-retention-gaps-explicit`, `AC-35-02-cross-surface-active-runs`, `AC-35-05-observer-stability-no-404`, `AC-35-11-compatibility-legacy-links`, `AC-35-10-idempotent-authorized-cancellation`, `AC-35-12-correlated-telemetry-diagnostics`, `AC-35-13-redaction-before-persistence`, `AC-35-15-gated-real-runtime-dogfood`
Rationale: agent-authored real-boundary suite builds the sprint-35 worktree binary and proves durable run identity across real boundaries deterministic Go tests cannot replace: one accepted operation is served with identical identity/lifecycle/terminal outcome by two independent local servers plus the real CLI run list/show commands against one shared workspace SQLite store; retained strictly-monotonic history replays from a second server via query cursors and Last-Event-ID while ahead/conflicting/malformed cursors return typed cursor_ahead/cursor_conflict/invalid_cursor errors backed by the durable last sequence; SIGKILL of the accepting server leaves the arbitrated outcome inspectable by a fresh replacement server, an anonymous browser detail page, and the legacy /api/v1/operations URL while pre-durable op_* ids receive typed 410 guidance; duplicate cancellation of a terminal run is idempotent without rewriting the single winning outcome and mutations still require fresh session+CSRF authority; run diagnostics --json and --support-export produce correlated redacted bundles; and the gated provider-backed multi-surface dogfood remains preserved as user-declined blocker evidence rather than claimed coverage

- `AC-35-01-durable-acceptance-fail-closed` — none (mapped tests: sprint-35-live-durable-acceptance-and-cross-surface-agreement)
- `AC-35-03-workspace-inspection-any-surface` — none (mapped tests: sprint-35-live-durable-acceptance-and-cross-surface-agreement)
- `AC-35-08-liveness-reconciliation-truthful` — none (mapped tests: sprint-35-live-durable-acceptance-and-cross-surface-agreement)
- `AC-35-09-single-terminal-outcome` — none (mapped tests: sprint-35-live-cancellation-idempotence-and-single-terminal-outcome, sprint-35-live-durable-acceptance-and-cross-surface-agreement)
- `AC-35-04-retained-history-other-server` — none (mapped tests: sprint-35-live-retained-replay-and-typed-cursors-across-servers)
- `AC-35-06-monotonic-sequence-replay` — none (mapped tests: sprint-35-live-retained-replay-and-typed-cursors-across-servers)
- `AC-35-07-typed-replay-gap-cursor-ahead` — none (mapped tests: sprint-35-live-retained-replay-and-typed-cursors-across-servers)
- `AC-35-14-bounded-retention-gaps-explicit` — none (mapped tests: sprint-35-live-retained-replay-and-typed-cursors-across-servers)
- `AC-35-02-cross-surface-active-runs` — none (mapped tests: sprint-35-live-abrupt-owner-death-and-session-independent-inspection)
- `AC-35-05-observer-stability-no-404` — none (mapped tests: sprint-35-live-abrupt-owner-death-and-session-independent-inspection)
- `AC-35-11-compatibility-legacy-links` — none (mapped tests: sprint-35-live-abrupt-owner-death-and-session-independent-inspection)
- `AC-35-10-idempotent-authorized-cancellation` — none (mapped tests: sprint-35-live-cancellation-idempotence-and-single-terminal-outcome)
- `AC-35-12-correlated-telemetry-diagnostics` — none (mapped tests: sprint-35-live-diagnostics-support-export-and-redaction)
- `AC-35-13-redaction-before-persistence` — none (mapped tests: sprint-35-live-diagnostics-support-export-and-redaction)
- `AC-35-15-gated-real-runtime-dogfood` — none (mapped tests: sprint-35-gated-real-runtime-dogfood-blocker-evidence)

Tests:
- `sprint-35-gated-real-runtime-dogfood-blocker-evidence` (suite `sprint-35`): `AC-35-15-gated-real-runtime-dogfood`
- `sprint-35-live-abrupt-owner-death-and-session-independent-inspection` (suite `sprint-35`): `AC-35-02-cross-surface-active-runs`, `AC-35-05-observer-stability-no-404`, `AC-35-11-compatibility-legacy-links`
- `sprint-35-live-cancellation-idempotence-and-single-terminal-outcome` (suite `sprint-35`): `AC-35-09-single-terminal-outcome`, `AC-35-10-idempotent-authorized-cancellation`
- `sprint-35-live-diagnostics-support-export-and-redaction` (suite `sprint-35`): `AC-35-12-correlated-telemetry-diagnostics`, `AC-35-13-redaction-before-persistence`
- `sprint-35-live-durable-acceptance-and-cross-surface-agreement` (suite `sprint-35`): `AC-35-01-durable-acceptance-fail-closed`, `AC-35-03-workspace-inspection-any-surface`, `AC-35-08-liveness-reconciliation-truthful`, `AC-35-09-single-terminal-outcome`
- `sprint-35-live-retained-replay-and-typed-cursors-across-servers` (suite `sprint-35`): `AC-35-04-retained-history-other-server`, `AC-35-06-monotonic-sequence-replay`, `AC-35-07-typed-replay-gap-cursor-ahead`, `AC-35-14-bounded-retention-gaps-explicit`

## Selected Scope And Rationale

Scope kind: `suite`
Scope: `sprint-35`
Rationale: agent-authored real-boundary suite builds the sprint-35 worktree binary and proves durable run identity across real boundaries deterministic Go tests cannot replace: one accepted operation is served with identical identity/lifecycle/terminal outcome by two independent local servers plus the real CLI run list/show commands against one shared workspace SQLite store; retained strictly-monotonic history replays from a second server via query cursors and Last-Event-ID while ahead/conflicting/malformed cursors return typed cursor_ahead/cursor_conflict/invalid_cursor errors backed by the durable last sequence; SIGKILL of the accepting server leaves the arbitrated outcome inspectable by a fresh replacement server, an anonymous browser detail page, and the legacy /api/v1/operations URL while pre-durable op_* ids receive typed 410 guidance; duplicate cancellation of a terminal run is idempotent without rewriting the single winning outcome and mutations still require fresh session+CSRF authority; run diagnostics --json and --support-export produce correlated redacted bundles; and the gated provider-backed multi-surface dogfood remains preserved as user-declined blocker evidence rather than claimed coverage
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

Run ID: `run--XvfeAIVbT`
Total: `6`
Passed: `6`
Failed: `0`
Skipped: `0`
Errors: `0`
Duration: `7.655s`
Runtime: `local-go`
Model: `none`

Executed tests:
- `sprint-35-live-durable-acceptance-and-cross-surface-agreement`: `passed`
- `sprint-35-live-retained-replay-and-typed-cursors-across-servers`: `passed`
- `sprint-35-live-abrupt-owner-death-and-session-independent-inspection`: `passed`
- `sprint-35-live-cancellation-idempotence-and-single-terminal-outcome`: `passed`
- `sprint-35-live-diagnostics-support-export-and-redaction`: `passed`
- `sprint-35-gated-real-runtime-dogfood-blocker-evidence`: `passed`

### External Evidence Identity And Links

- `run` `/home/antonioborgerees/coding/ultraplan/ultraplan-go-smoke/runs/run--XvfeAIVbT.json` sha256 `274948c757a601f3792223bc5d977b06e536230dd4140cb4a606936306150266` size `25727` modified `2026-08-23T14:11:41Z`
- `summary` `/home/antonioborgerees/coding/ultraplan/ultraplan-go-smoke/runs/run--XvfeAIVbT-summary.md` sha256 `cce4e11da5dddbd9d4b6c6c0e5ccc175523b1c81a16b90c77bd93fe85db3f661` size `1403` modified `2026-08-23T14:11:41Z`

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
