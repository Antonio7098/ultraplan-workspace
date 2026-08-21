# Sprint Smoke

Smoke status: `completed`
Verdict: `pass`
Date: `2026-08-19T18:31:11Z`

## Smoke Context

Project: `ultraplan-go`
Sprint: `32-hardening-and-release`
Artifact: `projects/ultraplan-go/sprints/32-hardening-and-release/smoke.md`

## Review Gate

Review verdict: `pass_with_findings`
Review fingerprint: `6e6902e21705467ad32d24f765f6784c6b61fea396eb73f7e453ab1931affe5a`
Diagnostic override: `false`
Override rationale: none

## Harness And Protocol

Harness: `ultraplan-go-smoke`
Protocol: `1.0`

## Smoke Authoring

Author run ID: `opencode-1`
Author model: `minimax-coding-plan/MiniMax-M3`
Changed harness paths:
- none; existing traceable suite retained after agent inspection

## Selected Scope And Rationale

Scope kind: `suite`
Scope: `sprint-32`
Rationale: agent-authored real-boundary suite covers the hardened local web release: real namespaced downward template hierarchy with semantic landmarks, real embedded layered CSS and JS asset matrix with revalidation cache policy, real state-bearing Cache-Control: no-store, real substantial validation and sprint workflows through the hardened browser, real browser-driven prepare/start/cancel across CSRF/session/host policy, real no-JavaScript complete sprint journey, real fail-closed configuration for LAN / hostname / invalid binds, real outside-tree packaging serving every documented page/api/asset without the source checkout, real internal/web dependency closure that excludes product/runtime/process, real stable CLI surface for serve + health, real redaction across JSON/SSE/rendered HTML, real accessibility (semantic landmarks, visible focus, narrow-viewport reflow) at desktop and 375x667, real focused 'go test' / '-race' / build gates, protocol discovery self-verification that the gated real-runtime and smoke-harness evidence path returns a non-empty Sprint 32 suite, real network-layer API envelope contract (unknown /api/ paths and unsupported methods return typed JSON 404/405 with Allow and no-store), real loopback CSRF and Origin rejection enforcing typed 403 JSON for every signed / cross-origin mutation against /api/v1/operations, real SIGTERM-driven graceful shutdown with exit code 0 and a server_stopped lifecycle marker, real SSE events endpoint contract rejecting unknown-operation replays as typed JSON 404 without ever beginning the event stream, repeated real-Chromium dashboard/project/sprint/study navigation stability, real CLI health --json / /api/v1/health / rendered dashboard agreement on workspace readiness, real malicious Markdown artifact round-trip rendered without executable scripts, javascript: links, or img[onerror] elements, real numeric-loopback-only Host header rejection at the network layer with typed JSON 403 host_rejected envelopes (LAN-bind, hostname-bind, and 'Host: evil.example.test' are all refused), real Chromium refresh-and-reopen during an active operation keeping the server-owned work running (disconnect isolates the subscription, not the operation), real SIGKILL-then-restart forced-termination reconciliation that persists truthful state without inferring success from process absence or artifact presence, real shared capability vocabulary exercised across multiple operation kinds through the single '/api/v1/operations/prepare' surface (no route-specific branch), real accessibility color-independent state cues and 'aria-live' regions on operation surfaces, real request-smuggling and body-limit matrix (oversized body, malformed JSON, duplicate Content-Length, wrong-method body) rejected at the security middleware with typed JSON errors, and real Chromium EventSource reconnect with a stale Last-Event-ID receiving typed SSE events with monotonic ids and no fabricated replay events, real typed-view-model and template-discipline scan that forbids filesystem paths, internal package names, runtime stack markers, and Go map[string] type markers from leaking through any rendered HTML or JSON surface (preserves AC-32-11 typed view-model discipline across the projection boundary), real compatibility-fixture envelope matrix that walks every documented /api/v1 success endpoint and verifies the frozen {data, meta.api_version="v1"} envelope with no-store cache policy (preserves AC-32-03 compatibility-fixture coverage at the network layer), and real concurrent prepare/start burst that exercises five simultaneous operations to terminal state with unique ids and no orphan subscribers, locks, or active operations after reconciliation (preserves AC-32-20 race/no-orphan behavior under realistic browser burst load)
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

Run ID: `run-8kxoK3aQec`
Total: `33`
Passed: `33`
Failed: `0`
Skipped: `0`
Errors: `0`
Duration: `47.734s`
Runtime: `local-browser`
Model: `none`

Executed tests:
- `sprint-32-live-namespaced-hierarchy-renders`: `passed`
- `sprint-32-live-layered-css-asset-matrix`: `passed`
- `sprint-32-live-disposable-js-asset-matrix`: `passed`
- `sprint-32-live-cache-policy-asset-vs-state`: `passed`
- `sprint-32-live-substantial-project-validation-journey`: `passed`
- `sprint-32-live-browser-driven-sprint-prepare-start-cancel`: `passed`
- `sprint-32-live-no-javascript-complete-sprint-journey`: `passed`
- `sprint-32-live-fail-closed-config-invalid-bind`: `passed`
- `sprint-32-live-outside-tree-packaging-and-asset-matrix`: `passed`
- `sprint-32-live-import-boundary-excludes-product-layers`: `passed`
- `sprint-32-live-cli-surfaces-for-serve-and-workspace`: `passed`
- `sprint-32-live-redaction-across-rendered-surfaces`: `passed`
- `sprint-32-live-accessibility-narrow-viewport-rendering`: `passed`
- `sprint-32-live-keyboard-focus-visibility`: `passed`
- `sprint-32-live-focused-verification-gates`: `passed`
- `sprint-32-live-protocol-discovery-completes-sprint-32-mapping`: `passed`
- `sprint-32-live-api-unknown-route-json-envelope`: `passed`
- `sprint-32-live-csrf-host-policy-enforcement`: `passed`
- `sprint-32-live-graceful-shutdown-drains-on-signal`: `passed`
- `sprint-32-live-sse-event-endpoint-header-contract`: `passed`
- `sprint-32-live-page-navigation-roundtrip-stability`: `passed`
- `sprint-32-live-cli-html-cross-surface-readiness-agreement`: `passed`
- `sprint-32-live-artifact-markdown-renders-without-script-injection`: `passed`
- `sprint-32-live-host-header-rejection-at-network-layer`: `passed`
- `sprint-32-live-browser-refresh-during-operation-preserves-work`: `passed`
- `sprint-32-live-forced-termination-persists-truthful-cleanup`: `passed`
- `sprint-32-live-shared-capability-vocabulary-across-operation-kinds`: `passed`
- `sprint-32-live-accessibility-color-independent-state-and-live-regions`: `passed`
- `sprint-32-live-request-smuggling-and-body-limits-rejected`: `passed`
- `sprint-32-live-sse-replay-gap-monotonic-real`: `passed`
- `sprint-32-live-compatibility-fixture-envelope-matrix`: `passed`
- `sprint-32-live-typed-view-model-and-template-discipline`: `passed`
- `sprint-32-live-concurrent-operation-start-race`: `passed`

### External Evidence Identity And Links

- `run` `/home/antonioborgerees/coding/ultraplan/ultraplan-go-smoke/runs/run-8kxoK3aQec.json` sha256 `44764485879eb1465bc791fb932ca9bbcd7faa396c774864e88918ab98caf972` size `153376` modified `2026-08-19T18:31:11Z`
- `summary` `/home/antonioborgerees/coding/ultraplan/ultraplan-go-smoke/runs/run-8kxoK3aQec-summary.md` sha256 `5a2f9a5ab373da25742babb751e053d206e385c02ef1fa167da9c8091a719402` size `3674` modified `2026-08-19T18:31:11Z`

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
