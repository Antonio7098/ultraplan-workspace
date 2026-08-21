# Sprint Smoke

Smoke status: `completed`
Verdict: `fail`
Date: `2026-08-17T19:02:40Z`

## Smoke Context

Project: `ultraplan-go`
Sprint: `31-web-operations`
Artifact: `projects/ultraplan-go/sprints/31-web-operations/smoke.md`

## Review Gate

Review verdict: `pass`
Review fingerprint: `100e6267272154b4d2649816ebd2d388b3df352626021bb2da4a42edc5b581fd`
Diagnostic override: `false`
Override rationale: none

## Harness And Protocol

Harness: `ultraplan-go-smoke`
Protocol: `1.0`

## Smoke Authoring

Author run ID: `opencode-1`
Author model: `minimax-coding-plan/MiniMax-M3`
Changed harness paths:
- `src/tests/sprint-31-web-operations.ts`

## Selected Scope And Rationale

Scope kind: `suite`
Scope: `sprint-31`
Rationale: agent-authored real-boundary suite exercises the built loopback server, Chromium prepare/start/cancel flow, real SSE EventSource, real SIGTERM draining and forced-termination reconciliation, real stale-fingerprint rejection, real cross-session isolation, real parallel capacity rejection, real product-owned sprint mutation lock conflict through the web surface, real allowlist and traversal-scope rejection, real confirmation-token replay rejection, real typed SSE event name and monotonic id projection, real redaction-before-retention, real error envelopes, real no-JavaScript fallback with CLI cancellation guidance, real import-boundary review, and real CLI-surface stability for the guarded web operation and SSE progress sprint
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

Run ID: `run-ru9q3mluk-`
Total: `24`
Passed: `18`
Failed: `6`
Skipped: `0`
Errors: `0`
Duration: `3m26.886s`
Runtime: `local-browser`
Model: `none`

Executed tests:
- `sprint-31-live-prepare-scope-and-binding`: `passed`
- `sprint-31-live-stale-fingerprint-rejection`: `passed`
- `sprint-31-live-operation-method-contract`: `passed`
- `sprint-31-live-sse-event-stream-monotonic`: `passed`
- `sprint-31-live-sse-replay-gap-recovery`: `passed`
- `sprint-31-live-subscriber-disconnect-isolation`: `failed`
- `sprint-31-live-graceful-shutdown-drains`: `passed`
- `sprint-31-live-deadline-escalation-cleanup-uncertain`: `passed`
- `sprint-31-live-restart-reconciliation`: `passed`
- `sprint-31-live-cross-session-isolation`: `passed`
- `sprint-31-live-capacity-rejection`: `failed`
- `sprint-31-live-html-confirmation-page`: `failed`
- `sprint-31-live-browser-prepare-start-cancel`: `failed`
- `sprint-31-live-no-javascript-operation`: `failed`
- `sprint-31-live-import-boundary`: `passed`
- `sprint-31-live-cli-surfaces-unchanged`: `passed`
- `sprint-31-live-redaction-before-retention`: `passed`
- `sprint-31-live-error-envelope`: `passed`
- `sprint-31-live-verification-gates`: `passed`
- `sprint-31-live-mutation-lock-conflict`: `passed`
- `sprint-31-live-allowlist-mutation-rejection`: `passed`
- `sprint-31-live-confirmation-replay-rejection`: `passed`
- `sprint-31-live-typed-sse-event-projection`: `passed`
- `sprint-31-live-no-javascript-cancellation-guidance`: `failed`

### External Evidence Identity And Links

- `run` `/home/antonioborgerees/coding/ultraplan/ultraplan-go-smoke/runs/run-ru9q3mluk-.json` sha256 `a92ab0d15c50878faeb3a4cad75dbc6390715777c5667a394e19b3540be2d321` size `90299` modified `2026-08-17T19:02:40Z`
- `summary` `/home/antonioborgerees/coding/ultraplan/ultraplan-go-smoke/runs/run-ru9q3mluk--summary.md` sha256 `fd0cabcd76eeddb36b36ecaea018bbfe17973123841ad4da7a2e00d24c50f10e` size `5816` modified `2026-08-17T19:02:40Z`

## Findings

The run exposed six failures. Four are most likely smoke-harness defects or stale
selectors; one is a credible product-capacity defect; and one may combine a
short-lived operation race with missing no-JavaScript recovery guidance. These
are working theories, not confirmed root causes.

### `sprint-31-live-subscriber-disconnect-isolation` — browser session was never established

- Severity: `high`
- Observed: Chromium's initial `fetch` to `/api/v1/health` failed; the server logged `403 origin_rejected`, so the test never started an operation or exercised subscriber disconnect isolation.
- Working theory: the probe calls `fetch` from a newly created `about:blank` page before navigating it to the loopback origin. That supplies a non-allowlisted/null browser origin and leaves `context.cookies()` empty. This is probably a harness setup defect rather than evidence that disconnecting an SSE subscriber cancels work.
- Supporting evidence: the only application request after health was a `403 origin_rejected`; no prepare, start, SSE, disconnect, or terminal-operation events appear in the issue evidence. The probe reads cookies before its first navigation.
- Next investigation: navigate to a same-origin page before obtaining CSRF/session state, then rerun this single test and confirm the operation reaches a terminal state after the SSE page closes.

### `sprint-31-live-capacity-rejection` — ninth operation was accepted

- Severity: `high`
- Observed: all nine start requests returned `202`; the ninth should have returned `429 operation_capacity` while eight operations were active.
- Working theory: either the live operation-capacity guard is checked against the wrong population, or the `validation` operations finish quickly enough that fewer than eight remain active when the ninth request arrives. The former is a product defect; the latter means the probe does not create sustained concurrency.
- Supporting evidence: server logs contain nine distinct `operation_started` events and no `operation_capacity` rejection. The failure reproduced on retry. The probe starts eight validation operations concurrently, then submits the ninth only after all eight start responses complete.
- Next investigation: inspect terminal timestamps/states for the first eight operation IDs and the hub's active-count transition. If they remain active, fix capacity accounting; if they complete first, use a controlled long-running operation to test the boundary.

### `sprint-31-live-html-confirmation-page` — stale form selector prevented prepare

- Severity: `high`
- Observed: Chromium timed out waiting for `POST /api/v1/operations/prepare`; the page and static assets loaded successfully, but no prepare request was sent.
- Working theory: the probe searches for `form.flow-action` on the Run page, while the current Run template renders per-stage forms as `form.operation-form.stage-start`. The scripted click therefore does nothing. This is probably a stale harness selector, not a confirmation-panel failure.
- Supporting evidence: server logs stop after the Run page, CSS, and JavaScript requests. The probe's evaluated script silently skips clicking when `form.flow-action` is absent. The current template contains `flow-action` on the Overview page and `stage-start` on the Run page.
- Next investigation: select a visible stage panel and click its `form.stage-start` submit button using an asserted locator; then verify the request and confirmation fields.

### `sprint-31-live-browser-prepare-start-cancel` — same stale selector blocked the full browser flow

- Severity: `high`
- Observed: Chromium timed out waiting for the initial prepare request, so start, navigation, and DELETE cancellation were never exercised.
- Working theory: this shares the stale `form.flow-action` Run-page selector used by the confirmation test. The failure occurs before product operation behavior is reached.
- Supporting evidence: only the Run page and static assets appear in server logs; there is no `operation_prepared` event. The test uses the same optional scripted click and waits 30 seconds for a request that cannot occur when the selector is absent.
- Next investigation: correct the Run-page locator, fail immediately when it is absent, and rerun this single flow through prepare, confirm, operation navigation, and DELETE interception.

### `sprint-31-live-no-javascript-operation` — expected control is not part of the Run page

- Severity: `high`
- Observed: with JavaScript disabled, Chromium timed out waiting for `select#run-flow-target`; no form submission occurred.
- Working theory: the probe expects an older select-based Run UI. The current Run page exposes individual server-rendered `form.stage-start` forms and reserves a stage selector for the Overview page under the different ID `overview-flow-target`. This is a harness/template-contract mismatch.
- Supporting evidence: the server returned the Run page successfully. Current template source has no `run-flow-target`, while each nonterminal stage has a normal POST form suitable for no-JavaScript submission.
- Next investigation: submit a visible `stage-start` form with JavaScript disabled and continue through `/operations/prepare` and `/operations/start`; separately decide whether a select control is actually a normative UI requirement.

### `sprint-31-live-no-javascript-cancellation-guidance` — fast completion lands on an unguided not-retained page

- Severity: `high`
- Observed: after starting validation, the subsequent operation page rendered “Operation not retained” and did not contain the expected CLI cancellation guidance.
- Working theory: validation terminates and leaves ephemeral retention before the browser follows the redirect, exposing the not-retained branch. That branch tells the user to refresh the owning page but does not include the no-JavaScript CLI/interrupt guidance present in the normal shell `<noscript>` content. This may be a real recovery-content gap revealed by a short-lived-operation race.
- Supporting evidence: prepare and start succeeded, but the captured HTML title is `Operation not retained`. The handler's not-retained message only recommends refreshing the owning page; the generic shell guidance says to use the equivalent CLI command and normal interrupt handling.
- Next investigation: record the started operation ID and terminal/retention transitions, confirm whether the redirect loses the operation before first render, and verify the guidance requirement against both retained and not-retained operation pages.

## Open Issues

- `runtime-sprint-31-live-subscriber-disconnect-isolation` (high, test `sprint-31-live-subscriber-disconnect-isolation`): `/home/antonioborgerees/coding/ultraplan/ultraplan-go-smoke/issues/runtime-sprint-31-live-subscriber-disconnect-isolation.md`
- `runtime-sprint-31-live-capacity-rejection` (high, test `sprint-31-live-capacity-rejection`): `/home/antonioborgerees/coding/ultraplan/ultraplan-go-smoke/issues/runtime-sprint-31-live-capacity-rejection.md`
- `runtime-sprint-31-live-html-confirmation-page` (high, test `sprint-31-live-html-confirmation-page`): `/home/antonioborgerees/coding/ultraplan/ultraplan-go-smoke/issues/runtime-sprint-31-live-html-confirmation-page.md`
- `runtime-sprint-31-live-browser-prepare-start-cancel` (high, test `sprint-31-live-browser-prepare-start-cancel`): `/home/antonioborgerees/coding/ultraplan/ultraplan-go-smoke/issues/runtime-sprint-31-live-browser-prepare-start-cancel.md`
- `runtime-sprint-31-live-no-javascript-operation` (high, test `sprint-31-live-no-javascript-operation`): `/home/antonioborgerees/coding/ultraplan/ultraplan-go-smoke/issues/runtime-sprint-31-live-no-javascript-operation.md`
- `runtime-sprint-31-live-no-javascript-cancellation-guidance` (high, test `sprint-31-live-no-javascript-cancellation-guidance`): `/home/antonioborgerees/coding/ultraplan/ultraplan-go-smoke/issues/runtime-sprint-31-live-no-javascript-cancellation-guidance.md`

## Resolved Issues

- none

## Mutation And Safety Check

Only smoke.md, flow-state.json, manifest-declared harness authoring paths, and manifest-declared external evidence roots were approved for mutation. Product source and governed sprint inputs were identity-checked before and after authoring.

## Verdict And Next Action

Verdict: `fail`
Next action: Inspect linked evidence, fix the selected-smoke failures, and rerun the containing suite.

## Manual Remediation Assessment — 2026-08-18

This assessment was added manually after the recorded smoke run. At the user's
direction, the smoke suite was not rerun. The original run counts, findings,
open-issue records, and `fail` verdict above remain unchanged as historical
evidence; the statuses below describe the later working trees and focused test
evidence only.

### Addressed findings

- [addressed] `sprint-31-live-subscriber-disconnect-isolation` — The smoke probe
  now navigates to the loopback origin before reading its session cookie and
  CSRF response header, then opens the SSE stream from that same-origin page.
- [addressed] `sprint-31-live-capacity-rejection` — Inspection confirmed that
  the hub rejects a ninth concurrently active operation. The live probe now
  distinguishes a legitimate slot release by comparing terminal and creation
  timestamps, and a deterministic blocking Go test holds all eight slots while
  asserting `operation_capacity` for the ninth start.
- [addressed] `sprint-31-live-html-confirmation-page` — The smoke probe now uses
  an asserted visible `form.stage-start` locator matching the current Run page.
- [addressed] `sprint-31-live-browser-prepare-start-cancel` — The browser flow
  uses the current stage form and continues to assert prepare, confirmation,
  start, operation-ID navigation, and DELETE cancellation behavior.
- [addressed] `sprint-31-live-no-javascript-operation` — The smoke probe now
  submits the server-rendered stage form instead of the removed select control.
  The Run page exposes its stage links and ordinary POST forms without requiring
  JavaScript.
- [addressed] `sprint-31-live-no-javascript-cancellation-guidance` — Active
  operation pages now expose a CSRF-protected server-rendered cancellation form,
  and retained active or terminal pages expose an explicit status refresh link.
  The probe now accounts for fast operations reaching a terminal state before
  the operation page is rendered.

### Focused verification

The following command passed after the remediation:

```text
go test ./internal/web ./internal/app ./internal/study
```

### Residual note

The `Operation not retained` error branch directs the operator to refresh the
owning project, sprint, or study page for durable status but does not repeat the
generic CLI/interrupt guidance. This is a low-priority recovery-copy improvement
and is not evidence that cancellation or operation retention is malfunctioning.

Manual remediation status: `addressed_pending_smoke_rerun`
Next action: none requested; retain the original failed smoke evidence until a
future explicitly requested run supersedes it.
