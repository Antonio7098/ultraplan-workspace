# Sprint Smoke

Smoke status: `completed`
Verdict: `fail`
Date: `2026-08-25T18:30:22Z`

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

Author run ID: `opencode-1`
Author model: `openrouter/stealth/ox-alpha`
Changed harness paths:
- none; existing traceable suite retained after agent inspection

## Coverage Mapping

Complete: `true`
Required coverage: `AC-37-01-writable-admission-fail-closed`, `AC-37-09-smoke-parity-single-authority`, `AC-37-14-scope-exclusions`, `AC-37-08-state-versioning-fencing`, `AC-37-03-target-immutability`, `AC-37-04-frozen-evidence-plans`, `AC-37-06-adjudication-authority`, `AC-37-07-canonical-assessment-report`, `AC-37-10-cross-surface-agreement`, `AC-37-12-cancellation-recovery-truthful`, `AC-37-13-browser-security-nojs`, `AC-37-15-gated-real-runtime-dogfood`
Rationale: agent-authored real-boundary suite builds the sprint-37 worktree binary and proves the boundaries deterministic Go tests cannot replace: the real qa CLI keeps one closed suite selector ('smoke' only) with stable rejection of resume-with-suite and unknown suites, gates non-dry external-harness execution behind --yes, and keeps smoke-suite dry-run write-free; writable QA admission fails closed on a real fixture with unvalidated governed inputs via a typed qa.stale_input error before any verification-state or target write; two dry-run maps of the real repository agree byte-stably with one primary owner per changed path, sequential writable concurrency (concurrent_investigators=1), frozen positive limits, and unchanged before/after target Git identity; smoke --dry-run and qa --suite smoke --dry-run resolve one manifest-driven containing-suite authority at the real CLI boundary (same harness identity, non-diagnostic sprint-37 suite scope, contained external evidence roots, identical prepared next action, review gate, and no durable writes from either dry-run); CLI status and a live local server agree on canonical phase and next action with an inert no-JavaScript HTML snapshot, the live focused evidence/adjudication/issues/assessment/smoke-suite routes are probed strictly over app facts, out-of-attempt evidence IDs fail closed without internal vocabulary, runtime-free qa recover truthfully reports missing state without inventing freshness, and qa cancel rejects an unknown durable run fail-closed; and Task 12's gated release/dogfood evidence remains preserved as deferred blocker evidence rather than claimed coverage. Isolation mechanics, bounded execution/cleanup truth, adjudication quorum rules, state v2 fencing, durable run authority, and cancellation races remain proven deterministically in the product's own fault-injected suites

- `AC-37-01-writable-admission-fail-closed` — none (mapped tests: sprint-37-live-qa-cli-suite-contract-and-yes-gate, sprint-37-live-writable-admission-fail-closed-fixture)
- `AC-37-09-smoke-parity-single-authority` — none (mapped tests: sprint-37-live-qa-cli-suite-contract-and-yes-gate, sprint-37-live-smoke-parity-single-authority)
- `AC-37-14-scope-exclusions` — none (mapped tests: sprint-37-live-qa-cli-suite-contract-and-yes-gate, sprint-37-live-smoke-parity-single-authority)
- `AC-37-08-state-versioning-fencing` — none (mapped tests: sprint-37-live-writable-admission-fail-closed-fixture)
- `AC-37-03-target-immutability` — none (mapped tests: sprint-37-live-dry-run-map-stability-target-immutability)
- `AC-37-04-frozen-evidence-plans` — none (mapped tests: sprint-37-live-dry-run-map-stability-target-immutability)
- `AC-37-06-adjudication-authority` — none (mapped tests: sprint-37-live-cross-surface-focused-routes-and-no-javascript-snapshot)
- `AC-37-07-canonical-assessment-report` — none (mapped tests: sprint-37-live-cross-surface-focused-routes-and-no-javascript-snapshot)
- `AC-37-10-cross-surface-agreement` — none (mapped tests: sprint-37-live-cross-surface-focused-routes-and-no-javascript-snapshot)
- `AC-37-12-cancellation-recovery-truthful` — none (mapped tests: sprint-37-live-cross-surface-focused-routes-and-no-javascript-snapshot)
- `AC-37-13-browser-security-nojs` — none (mapped tests: sprint-37-live-cross-surface-focused-routes-and-no-javascript-snapshot)
- `AC-37-15-gated-real-runtime-dogfood` — none (mapped tests: sprint-37-gated-real-runtime-dogfood-blocker-evidence)

Tests:
- `sprint-37-gated-real-runtime-dogfood-blocker-evidence` (suite `sprint-37`): `AC-37-15-gated-real-runtime-dogfood`
- `sprint-37-live-cross-surface-focused-routes-and-no-javascript-snapshot` (suite `sprint-37`): `AC-37-06-adjudication-authority`, `AC-37-07-canonical-assessment-report`, `AC-37-10-cross-surface-agreement`, `AC-37-12-cancellation-recovery-truthful`, `AC-37-13-browser-security-nojs`
- `sprint-37-live-dry-run-map-stability-target-immutability` (suite `sprint-37`): `AC-37-03-target-immutability`, `AC-37-04-frozen-evidence-plans`
- `sprint-37-live-qa-cli-suite-contract-and-yes-gate` (suite `sprint-37`): `AC-37-01-writable-admission-fail-closed`, `AC-37-09-smoke-parity-single-authority`, `AC-37-14-scope-exclusions`
- `sprint-37-live-smoke-parity-single-authority` (suite `sprint-37`): `AC-37-09-smoke-parity-single-authority`, `AC-37-14-scope-exclusions`
- `sprint-37-live-writable-admission-fail-closed-fixture` (suite `sprint-37`): `AC-37-01-writable-admission-fail-closed`, `AC-37-08-state-versioning-fencing`

## Selected Scope And Rationale

Scope kind: `suite`
Scope: `sprint-37`
Rationale: agent-authored real-boundary suite builds the sprint-37 worktree binary and proves the boundaries deterministic Go tests cannot replace: the real qa CLI keeps one closed suite selector ('smoke' only) with stable rejection of resume-with-suite and unknown suites, gates non-dry external-harness execution behind --yes, and keeps smoke-suite dry-run write-free; writable QA admission fails closed on a real fixture with unvalidated governed inputs via a typed qa.stale_input error before any verification-state or target write; two dry-run maps of the real repository agree byte-stably with one primary owner per changed path, sequential writable concurrency (concurrent_investigators=1), frozen positive limits, and unchanged before/after target Git identity; smoke --dry-run and qa --suite smoke --dry-run resolve one manifest-driven containing-suite authority at the real CLI boundary (same harness identity, non-diagnostic sprint-37 suite scope, contained external evidence roots, identical prepared next action, review gate, and no durable writes from either dry-run); CLI status and a live local server agree on canonical phase and next action with an inert no-JavaScript HTML snapshot, the live focused evidence/adjudication/issues/assessment/smoke-suite routes are probed strictly over app facts, out-of-attempt evidence IDs fail closed without internal vocabulary, runtime-free qa recover truthfully reports missing state without inventing freshness, and qa cancel rejects an unknown durable run fail-closed; and Task 12's gated release/dogfood evidence remains preserved as deferred blocker evidence rather than claimed coverage. Isolation mechanics, bounded execution/cleanup truth, adjudication quorum rules, state v2 fencing, durable run authority, and cancellation races remain proven deterministically in the product's own fault-injected suites
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

Run ID: `run-eo-sf7iKtE`
Total: `6`
Passed: `4`
Failed: `2`
Skipped: `0`
Errors: `0`
Duration: `25.751s`
Runtime: `local-go`
Model: `none`

Executed tests:
- `sprint-37-live-qa-cli-suite-contract-and-yes-gate`: `passed`
- `sprint-37-live-writable-admission-fail-closed-fixture`: `passed`
- `sprint-37-live-dry-run-map-stability-target-immutability`: `passed`
- `sprint-37-live-smoke-parity-single-authority`: `failed`
- `sprint-37-live-cross-surface-focused-routes-and-no-javascript-snapshot`: `failed`
- `sprint-37-gated-real-runtime-dogfood-blocker-evidence`: `passed`

### External Evidence Identity And Links

- `run` `/home/antonioborgerees/coding/ultraplan/ultraplan-go-smoke/runs/run-eo-sf7iKtE.json` sha256 `4570e367bf6d073807fa1b95593fc7e3edaa0bc3927028f4bde166160b7d48ac` size `79541` modified `2026-08-25T18:30:21Z`
- `summary` `/home/antonioborgerees/coding/ultraplan/ultraplan-go-smoke/runs/run-eo-sf7iKtE-summary.md` sha256 `4ef3fa653efa596fd3d9e71ee12fa53cada3c57269eae5092a2bc12306ba9db8` size `2458` modified `2026-08-25T18:30:21Z`

## Findings

### `sprint-37-live-smoke-parity-single-authority` — sprint-37-live-smoke-parity-single-authority

- Severity: `high`
- Observed: a dry-run pair wrote durable smoke or QA state although both paths are documented write-free
- Working theory: The real-boundary smoke probe 'sprint-37-live-smoke-parity-single-authority' did not observe the contracted behavior because: a dry-run pair wrote durable smoke or QA state although both paths are documented write-free. The expected behavior is documented in the sprint-32 acceptance criteria, and the test would pass i…
- Supporting evidence: Sprint 37 fresh target build: ["go","build","-o","/tmp/ultraplan-sprint37-bin-WvYo5Q/ultraplan","./cmd/ultraplan"]

Sprint 37 fresh target build output: exit=0
stdout:

stderr:


ultraplan sprint ultraplan-go 37-evidence-qa-smoke smoke --dry-run --json: exit=0
stdout:
{"operation":"sprint.smoke","result":{"project":"u…
- Next investigation: Inspect the captured OpenCode runtime logs and prompt output under the test artifact directory; confirm the model/provider/timeout and the governed workspace fingerprint match the runtime expectation, then rerun the affected test with --tests <name> to isolate the regression.

### `sprint-37-live-cross-surface-focused-routes-and-no-javascript-snapshot` — sprint-37-live-cross-surface-focused-routes-and-no-javascript-snapshot

- Severity: `high`
- Observed: focused route /adjudication returned 503: {"error":{"code":"unavailable","message":"The service is unavailable."},"meta":{"api_version":"v1","request_id":"d3adc6dce77b0866d4cac4a0d7dca30a"}}
- Working theory: The real-boundary smoke probe 'sprint-37-live-cross-surface-focused-routes-and-no-javascript-snapshot' did not observe the contracted behavior because: focused route /adjudication returned 503: {"error":{"code":"unavailable","message":"The service is unavailable."},"meta":{"api_version":"v1","request_id":"d3adc6dce77b…
- Supporting evidence: Sprint 37 fresh target build: ["go","build","-o","/tmp/ultraplan-sprint37-bin-UZG8MY/ultraplan","./cmd/ultraplan"]

Sprint 37 fresh target build output: exit=0
stdout:

stderr:


ultraplan sprint ultraplan-go 37-evidence-qa-smoke qa status --json: exit=0
stdout:
{"operation":"sprint.qa","result":{"schema_version":1,"p…
- Next investigation: Inspect the captured OpenCode runtime logs and prompt output under the test artifact directory; confirm the model/provider/timeout and the governed workspace fingerprint match the runtime expectation, then rerun the affected test with --tests <name> to isolate the regression.


## Open Issues

- `runtime-sprint-37-live-smoke-parity-single-authority` (high, test `sprint-37-live-smoke-parity-single-authority`): `/home/antonioborgerees/coding/ultraplan/ultraplan-go-smoke/issues/runtime-sprint-37-live-smoke-parity-single-authority.md`
- `runtime-sprint-37-live-cross-surface-focused-routes-and-no-javascript-snapshot` (high, test `sprint-37-live-cross-surface-focused-routes-and-no-javascript-snapshot`): `/home/antonioborgerees/coding/ultraplan/ultraplan-go-smoke/issues/runtime-sprint-37-live-cross-surface-focused-routes-and-no-javascript-snapshot.md`

## Resolved Issues

- none

## Mutation And Safety Check

Only smoke.md, flow-state.json, manifest-declared harness authoring paths, and manifest-declared external evidence roots were approved for mutation. Product source and governed sprint inputs were identity-checked before and after authoring.

## Verdict And Next Action

Verdict: `fail`
Next action: Inspect linked evidence, fix the selected-smoke failures, and rerun the containing suite.
