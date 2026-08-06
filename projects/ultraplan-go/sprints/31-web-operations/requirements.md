# Sprint Requirements: Guarded Web Operations and SSE Progress

> Project: `ultraplan-go`
> Sprint: `31-web-operations`
> Purpose: the authoritative, human-readable sprint contract. All other sprint artifacts must satisfy these requirements.

## Sprint Goal

Expose existing guarded validation, planning, execute, review, smoke, verify, study run-loop, and cancellation operations through the local web surface with server-issued confirmations and bounded SSE progress, without adding a second workflow engine or durable web-owned state. Gracefully stopping the server must cancel every active server-owned operation through the same product/runtime cleanup paths used by explicit cancellation.

## Required Outputs

| Output | Path | Description |
| --- | --- | --- |
| Sprint requirements | `projects/ultraplan-go/sprints/31-web-operations/requirements.md` | This sprint contract derived from the project index, roadmap, PRD, TRD, architecture, Sprint 30 artifacts, and prior sprint reviews. |
| Architecture reasoning | `projects/ultraplan-go/sprints/31-web-operations/reasoning/architecture.md` | Decision record for operation capabilities, app/web boundaries, mutation locks, cancellation ownership, graceful server draining and cancel-all shutdown, and recovery through durable product state. |
| API design reasoning | `projects/ultraplan-go/sprints/31-web-operations/reasoning/api-design.md` | Decision record for confirmation, operation start/status/result/cancellation endpoints, SSE event schema, shutdown-visible operation states, error envelopes, and compatibility constraints. |
| Frontend reasoning | `projects/ultraplan-go/sprints/31-web-operations/reasoning/frontend.md` | Decision record for browser confirmations, progress, result, finding, cancellation, server-shutdown interruption, failure, and recovery views using server-rendered pages plus minimal JavaScript. |
| App operation use cases | `../ultraplan-go/internal/app/web_operations.go` | Typed prepare/start/status/cancel operation use cases over existing app operations, with normalized scope, fingerprints, runtime/model summaries, safe progress events, and canonical cancellation propagation. |
| App operation tests | `../ultraplan-go/internal/app/web_operations_test.go` | Covers confirmation binding, supported operation mapping, explicit and server-shutdown cancellation, typed errors, redaction, and CLI/TUI/web agreement with fake dependencies. |
| Sprint mutation lock support | `../ultraplan-go/internal/sprint/locks.go` | Product-owned per-sprint mutation exclusion for sprint flow, execute, review, smoke, and verify operations with actionable conflict diagnostics. |
| Sprint lock tests | `../ultraplan-go/internal/sprint/locks_test.go` | Covers lock acquire/release, conflict reporting, cancellation cleanup, server-shutdown cleanup ordering, stale-lock handling policy, and non-overlap with existing study locks. |
| Web operation hub | `../ultraplan-go/internal/web/operations.go` | Bounded ephemeral operation hub with operation IDs, cancellation functions, recent safe event buffers, subscribers, terminal results, short retention, server draining state, and cancel-all shutdown coordination. |
| Web operation routes | `../ultraplan-go/internal/web/routes.go` | Adds the confirmed-operation, status/result, SSE, and cancellation routes while preserving Sprint 30 read-only routes and API error behavior. |
| Web operation handlers | `../ultraplan-go/internal/web/handlers.go` | Maps HTTP requests to app operation use cases and renders JSON/HTML/SSE without importing product modules or CLI handlers. |
| Web confirmation and CSRF/session support | `../ultraplan-go/internal/web/security.go` | Extends local HTTP security with short-lived same-origin mutation protection and confirmation-token validation. |
| Web operation tests | `../ultraplan-go/internal/web/operations_test.go` | Covers operation lifecycle, retention, draining, rejection of new work during shutdown, cancellation of all active server-owned operations, bounded cleanup, result projection, conflict handling, redaction, and no durable web state. |
| Web route/API tests | `../ultraplan-go/internal/web/routes_test.go` | Covers POST/DELETE methods, prepare/start/status/cancel route shape, API envelopes, method rejection, server-draining responses, and compatibility with Sprint 30 read-only endpoints. |
| SSE tests | `../ultraplan-go/internal/web/sse_test.go` | Covers typed SSE events, monotonic event IDs, heartbeat comments, reconnect, buffer rollover, slow subscribers, disconnect behavior, server-shutdown cancellation events, and bounded stream closure. |
| Browser operation templates | `../ultraplan-go/internal/web/templates/*.html` | Updates pages/partials for operation scope confirmation, progress, results, findings, failures, cancellation, server shutdown, interruption, and recovery guidance. |
| Browser operation JavaScript | `../ultraplan-go/internal/web/static/*` | Minimal dependency-free progressive enhancement for prepare/start/cancel requests, EventSource subscription, reconnect, and durable-state refresh prompts. |
| Web documentation | `../ultraplan-go/docs/local-web.md` | Documents guarded browser operations, confirmation flow, SSE progress, explicit cancellation, browser disconnect behavior, graceful server shutdown cancellation, recovery, local-only trust limits, and durable state authority. |
| CLI/API documentation updates | `../ultraplan-go/docs/cli-reference.md` | Documents web operation behavior and server shutdown semantics as browser/API affordances over existing CLI/TUI use cases without changing CLI automation semantics. |
| Architecture documentation updates | `../ultraplan-go/docs/architecture.md` | Documents operation-hub ephemerality, server-owned run lifecycle, graceful shutdown draining and cancellation, SSE as progress-only transport, app/product ownership, mutation locks, and no web-owned workflow state. |

## Acceptance Criteria

- [ ] `POST /api/v1/operations/prepare` returns a normalized operation scope, affected workspace/target paths, runtime/model or harness information when relevant, mutation class, current governed-input fingerprint, expiry, and a server-issued confirmation token bound to the normalized request.
- [ ] Mutating or runtime-backed operations cannot start without a valid, unexpired, current confirmation token that matches the normalized start request and current input fingerprint.
- [ ] Supported browser operations are limited to existing app use cases for validation, prompt preview, dry-run, sprint flow, execute, review, smoke, verify, study run-loop, operation status/result inspection, and explicit cancellation.
- [ ] HTTP commands use `POST` for prepare/start and `DELETE` for cancellation; SSE is one-way progress only and cannot start, mutate, complete, fail, or cancel work.
- [ ] The operation hub is bounded and ephemeral: it stores only safe recent events, subscribers, cancellation handles, and terminal results for short retention, and server restart recovery rereads durable product state.
- [ ] Every operation started by the server is owned by the server lifecycle, including its root context, cancellation function, nested tasks, runtime or harness calls, process tree, locks, progress publication, cleanup, and durable terminal-state reconciliation; detached background runs that intentionally outlive the server are not supported.
- [ ] SSE emits typed event names with monotonic event IDs, heartbeat comments, safe redacted payloads, and reconnect behavior that resumes from buffered events when available or instructs the browser to refresh durable status when not.
- [ ] A disconnected or slow browser subscriber cannot block the underlying app operation, cannot change the operation result, and cancels only its own subscription.
- [ ] Explicit cancellation reaches the shared app operation context, product/runtime cleanup path, or smoke-harness cleanup path as applicable, and leaves durable state recoverable from CLI, TUI, and browser status views.
- [ ] Graceful server shutdown enters a draining state, rejects new operation starts, prevents queued work from beginning, records `reason: server_shutdown`, invokes every active operation's canonical cancellation function exactly once, propagates cancellation through nested tasks and owned process trees, waits for bounded cleanup and reconciliation, persists a truthful cancelled, interrupted, or cleanup-uncertain outcome, and only then closes SSE streams and exits.
- [ ] A graceful shutdown deadline that expires must escalate cleanup for remaining owned process trees and record interruption or cleanup uncertainty; it must not report success, release mutation locks as though cleanup succeeded, or leave active child processes intentionally detached.
- [ ] After a crash or forced termination, the next server start reconciles stale running operations and locks from durable product state, marks unresolved work interrupted or recovery-required, and never promotes stale running state to success.
- [ ] Conflicting sprint mutations and existing study run-loop mutations fail with actionable conflict diagnostics rather than running concurrently; sprint mutation exclusion is owned by `internal/sprint`, not by HTTP middleware alone.
- [ ] Structured API errors cover conflict, stale confirmation, invalid request, validation failure, runtime failure, cancellation, server draining, unavailable prerequisite, timeout, and internal failure without leaking secrets, raw provider payloads, unrestricted stderr, or absolute local paths.
- [ ] Browser pages show operation scope before confirmation and show truthful progress, final result, findings, failures, explicit or server-shutdown cancellation status, stale/reconnect guidance, and recovery actions after completion or interruption.
- [ ] `internal/web` does not import `internal/study`, `internal/project`, `internal/sprint`, runtime adapters, process adapters, or CLI handlers; HTTP handlers call typed `internal/app` use cases only.
- [ ] Existing CLI and TUI behavior remain supported and unchanged except for shared app-use-case additions required by the web surface.
- [ ] Normal verification passes with deterministic fakes: focused app operation tests, focused sprint lock tests, focused web operation/SSE/security/route/template tests, graceful-shutdown cancel-all and restart-reconciliation tests, `go test ./...`, `go test -race ./...`, and `go build ./cmd/ultraplan` from `../ultraplan-go`.

## Non-Goals

- Adding new product workflows, new sprint stages, or a browser-specific workflow engine is not included.
- Detached or daemonized local operations that intentionally survive graceful server shutdown are not included; a future durable worker architecture must define separate leases, heartbeats, ownership, and reconciliation before survival across control-plane restart is allowed.
- Browser editing of arbitrary workspace files, generated artifacts, project docs, plans, review findings, or smoke issue records is not included.
- Hosted SaaS, LAN/public binding, remote access, accounts, authentication teams, permissions, tenants, collaboration, remote workers, and remote workspace synchronization are not included.
- WebSockets, bidirectional terminal/session transport, live agent chat, or browser-driven remote shell/process control is not included.
- A database-backed state model, browser-owned durable state, durable web job queue, or operation history replacing product run state is not included.
- New review, smoke, execute, planning, or study run-loop semantics are not included; the browser may only expose current app use cases.
- General-purpose issue tracking, issue assignment, project-management features, automatic product fixes, and automatic Git mutation are not included.
- A frontend framework, Node.js runtime, Vite, separate frontend process, service worker, or asset build pipeline is not included.
- Browser-level API compatibility hardening, accessibility release audit, gated real-runtime browser evidence, and release packaging are deferred to Sprint 32 unless needed to prove Sprint 31 behavior with fakes.

## Constraints

- `internal/web` may depend on `internal/app` operation/query interfaces and plain result types only; it must not import product modules, runtime/process adapters, or CLI command handlers directly.
- HTTP handlers must not invoke `ultraplan` as a subprocess, parse CLI output, duplicate product validators, duplicate workflow state machines, or decide review/smoke/execute verdicts.
- The server operation hub is transport-lifecycle state only; workspace artifacts, `flow-state.json`, `.run-state.json`, study run state, `review.md`, `smoke.md`, and external harness evidence remain authoritative.
- Confirmation tokens must be short-lived, same-server-issued, bound to normalized operation scope and current fingerprints, and invalidated by stale inputs or mismatched requests.
- Product-owned locks must guard mutations; HTTP middleware may enforce request policy but must not be the only concurrency guard.
- SSE buffers, active operations, subscribers, request bodies, event payloads, and retention windows must be explicitly bounded.
- Browser disconnect cancels only the SSE subscription; explicit cancellation must use the operation cancellation endpoint and shared app context.
- Graceful server shutdown is different from browser disconnect: it must stop accepting new mutations, cancel every active server-owned operation, use the canonical cancellation and cleanup path, and wait for bounded terminal-state reconciliation before process exit.
- Mutation locks and cleanup ownership must remain held until the affected operation has reached a truthful terminal or cleanup-uncertain state; shutdown must not release them early merely to complete HTTP server closure.
- The Phase 4 server remains loopback-only, same-origin, CSRF-protected, body-limited, timeout-bounded, Host/Origin-validated, security-header-protected, and redaction-safe.
- The UI must remain Go `html/template` plus embedded CSS/minimal dependency-free JavaScript; no frontend build step or browser-owned product state is allowed.
- Tests must use `httptest`, fake app use cases, fake runtimes, fake smoke harnesses, and deterministic temporary workspaces by default; real runtime/harness/browser tests are gated or deferred.
- Automatic Git add, commit, push, branch, merge, reset, or other Git mutation remains prohibited.

## Dependencies

| Prior Sprint / Output | Required For | Notes |
| --- | --- | --- |
| Sprint 30: `projects/ultraplan-go/sprints/30-web-foundations/requirements.md` and implementation outputs | Local web foundation | Sprint 31 builds on `ultraplan serve`, loopback security, `/api/v1` envelopes, embedded templates/static assets, artifact previews, and explicit `cmd -> app+tui+web` composition. |
| Sprint 30 query-only app/web boundary | Operation boundary design | Mutation-capable app operation use cases must be separate from the read-only facade rather than broadening web capability accidentally. |
| Sprints 24-25 TUI operational controls | Guarded local operation precedent | Browser operations must match the existing confirmation/progress/cancellation discipline exposed by TUI use cases. |
| Sprints 26-29 review, smoke, and verify use cases | Browser verification operations | Web review, smoke, verify, findings, and recovery views must call existing app use cases and preserve review-before-smoke, verdict, staleness, and harness-evidence semantics. |
| Sprint 23 execute use cases and `.run-state.json` | Browser execute operation/status | Web execute operations must preserve durable task state, deterministic task IDs, cancellation, diagnostics, and no Git mutation. |
| Study run-loop locking and cancellation behavior | Browser study run-loop operation | Web study operations must use existing study locks and durable state rather than adding HTTP-owned study scheduling. |
| `../ultraplan-go/docs/plans/server-shutdown-run-cancellation-contract.md` | Graceful shutdown ownership and cancellation | Defines the normative draining, cancel-all, process-tree cleanup, terminal-state, SSE closure, forced-termination reconciliation, and future durable-worker boundaries that Sprint 31 must implement. |
| `projects/ultraplan-go/docs/ARCHITECTURE.md` Phase 4 rules | Package layout and dependency direction | Defines `internal/web` ownership, operation hub ephemerality, and prohibition on web-owned workflow state. |
| `projects/ultraplan-go/docs/TRD.md` sections 7.5 and 18A | HTTP, SSE, operation, security, and testing scope | Sprint 31 implements guarded operations and SSE progress subset of Phase 4. |
| `projects/ultraplan-go/docs/PRD.md` browser operation scenario | Product behavior and non-goals | Confirms browser operations are local, guarded, and backed by durable workspace state. |

## Review Expectations

| What | How Verified |
| --- | --- |
| Confirmation preflight is binding and current | Run focused app/web confirmation tests and inspect fingerprint, normalized-scope, expiry, token-binding, and stale-input checks. |
| Operation endpoints are method-correct and scoped | Run route/API tests for prepare, start, status/result, SSE, cancellation, method rejection, server-draining rejection, and unknown `/api/` behavior. |
| SSE is progress-only and bounded | Run SSE tests for event IDs, event names, heartbeat, reconnect, buffer rollover, slow subscribers, disconnect, server-shutdown cancellation events, bounded closure, and absence of command semantics. |
| Explicit cancellation reaches product work | Run app/web cancellation tests with fake runtime and fake harness operations; inspect durable-state recovery assertions. |
| Graceful shutdown cancels all server-owned work | Start multiple fake long-running operations, initiate graceful shutdown, verify draining and rejection of new starts, verify each canonical cancel function is called once, verify nested runtime/harness/process cleanup and bounded wait, assert truthful durable terminal or interrupted state, confirm locks are released only after reconciliation, and confirm the server and SSE streams then close. |
| Crash and forced-stop recovery is truthful | Seed stale running state and locks, restart the server, and verify reconciliation reports interrupted or recovery-required rather than success. |
| Product mutation locks prevent conflicts | Run sprint lock and study-operation conflict tests; inspect diagnostics for actionable scope/lock information. |
| Web boundary remains transport-only | Inspect `internal/web` imports and handlers; confirm no direct product/runtime/process/CLI-handler dependencies. |
| Durable state remains authoritative | Review operation hub implementation and tests proving restart/refresh reads product state and no durable web job store is introduced. |
| Browser UI is truthful and progressively enhanced | Run template/static tests and inspect views for scope confirmation, progress, result, findings, failure, explicit cancellation, server-shutdown interruption, reconnect, and recovery states. |
| Security and redaction remain enforced | Run Host/Origin/CSRF/body-limit/security-header/redaction tests covering mutating routes and SSE responses. |
| Existing surfaces do not regress | Run focused app dispatch/TUI/CLI agreement tests and full `go test ./...`. |
| Sprint verification passes | Run `go test ./...`, `go test -race ./...`, and `go build ./cmd/ultraplan` from `../ultraplan-go`. |
