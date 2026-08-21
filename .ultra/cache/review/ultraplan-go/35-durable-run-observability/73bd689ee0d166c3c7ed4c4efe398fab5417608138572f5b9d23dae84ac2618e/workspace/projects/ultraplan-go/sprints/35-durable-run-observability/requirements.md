# Sprint 35 Requirements: Durable Run Identity and Cross-Surface Observability

> Project: `ultraplan-go`
> Sprint: `35-durable-run-observability`
> Purpose: the authoritative, human-readable sprint contract. Later sprint artifacts must reason through the open decisions and satisfy the fixed observable outcomes.

## Sprint Goal

Make every accepted runtime-backed UltraPlan execution a durable, workspace-visible run that can be discovered, inspected, followed, cancelled when authorized, and reconciled consistently from CLI, TUI, and any local web-server instance attached to the same workspace.

The sprint repairs the operational control plane exposed by real use: CLI-started work must appear in the browser running-process indicator, a run opened through another server must retain useful history and continue with live events, and a browser session or server restart must not turn a legitimate run into an unexplained `404 Operation not retained` dead end.

This sprint does not change the authority of project, sprint, study, or Markdown artifacts. Durable operational run state is a separate product concern from the later gated decision about alternate persistence for authored product content.

## Required Outputs

| Deliverable | Required Outcome |
|---|---|
| Run identity and lifecycle model | A stable workspace-scoped run ID, attempt identity, operation/stage identity, timestamps, ownership/liveness facts, safe correlation fields, and one arbitrated terminal outcome exist independently of a browser session, HTTP server process, or event subscriber. |
| Durable run repository and event journal | Accepted runs and sanitized progress events are durably recorded before they are advertised as observable; events have monotonic per-run ordering and can be replayed after reconnect, session expiry, or server restart. The reasoning and plan stages select the concrete package and storage mechanism. |
| Workspace-wide active-run projection | CLI, TUI, dashboard, project, sprint, study, and local web surfaces derive running counts and run lists from the same workspace-wide lifecycle model rather than the current web process's in-memory operation map or only the page currently being viewed. |
| Stable run inspection | A durable run detail surface reports identity, target, lifecycle state, latest heartbeat, stage/task progress, safe diagnostics, event history, terminal result, and recovery guidance. Read visibility is not coupled to the browser session that started the work. |
| Replayable live event delivery | The web event transport resumes from a durable sequence cursor, makes truncation or retention gaps explicit, and continues from durable state when a different local server instance serves the same workspace. SSE may remain the transport; the transport is not the source of truth. |
| Ownership, lease, and reconciliation | Active executions publish enough owner/lease/heartbeat and process/runtime identity to distinguish live, stalled, interrupted, cleanup-uncertain, and terminal work conservatively, including PID reuse and abrupt owner death. Startup and periodic reconciliation never infer success from process absence or artifact presence. |
| Cross-surface cancellation | Authorized cancellation resolves the durable run, reaches the current owner through a selected coordination mechanism, is idempotent, and participates in the same single-terminal-outcome arbitration as completion, failure, timeout, interruption, and server shutdown. |
| Compatibility and migration | Current sprint/study locks, flow state, execute state, runtime checkpoints, operation URLs, and JSON clients have an explicit migration or compatibility story. Legacy or expired operation links resolve to a durable run, a retained tombstone, or precise recovery guidance rather than a generic dashboard error. |
| Operational telemetry | Structured logs, local metrics, diagnostics, and optional traces share run/attempt/stage/task/runtime correlation. Operator-visible health distinguishes storage failure, stale ownership, reconciliation backlog, replay gap, subscriber lag/drop, and cancellation uncertainty. |
| Verification and fault injection | Deterministic, race, integration, and browser tests prove CLI-to-web discovery, multiple local servers, refresh/reconnect replay, session expiry, server restart, abrupt owner death, stale leases, PID reuse, slow subscribers, bounded storage, duplicate cancellation, and terminal-state races. |
| Documentation | Architecture, PRD, TRD, CLI/API, local-web, operations, and recovery documentation explain the lifecycle model, authority boundaries, visibility guarantees, retention behavior, diagnostic workflow, compatibility policy, and intentionally unresolved design choices. |

## Acceptance Criteria

- A runtime-backed command receives a durable run ID before child execution begins, or fails closed without starting the child when required run persistence cannot be established.
- Two concurrent CLI-started runs in one workspace appear as two active runs in the browser top bar and workspace run list without requiring navigation to either owning sprint or study page.
- Work started through CLI, TUI, or one local web-server instance can be inspected from another supported local surface attached to the same workspace.
- Opening a current run through a different local server provides retained event history and subsequently committed live events, subject only to an explicit documented topology boundary selected during reasoning.
- Refreshing a page, rotating or expiring a browser session, restarting the observing server, or losing an SSE connection does not erase the run identity or turn a valid run into an unexplained 404.
- Every event visible to a subscriber has a stable run ID and monotonic per-run sequence. Reconnect resumes from a cursor; if the requested history is unavailable, the API emits a typed gap/retention response and a current durable snapshot.
- Running counts use lifecycle state from the shared run projection. No UI depends on a planning-stage status value that cannot represent running execution.
- Active-run truth includes a renewable lease or equivalent liveness contract. Reconciliation is conservative, observable, idempotent, and safe against stale PIDs and owner crashes.
- Exactly one terminal outcome wins for completion, failure, cancellation, timeout, interruption, cleanup uncertainty, and reconciliation races. Repeated observers and cancellation requests cannot rewrite that outcome inconsistently.
- Durable operational records do not become a shadow authority for `requirements.md`, `code-context.md`, flow-state outcomes, generated artifacts, Git/source state, or external smoke evidence.
- Raw provider payloads, secrets, full prompts, unsafe paths, and unbounded output are not persisted or streamed by default. Redaction occurs before durable recording and fan-out.
- Retention, compaction, disk limits, subscriber backpressure, and degraded-mode behavior are explicit and testable; event loss or persistence failure is never silent.
- Normal test suites pass under `go test ./...` and `go test -race ./...`; the CLI builds; a gated real-runtime dogfood demonstrates cross-surface observation and recovery.

## Open Questions For Reasoning

The sprint must answer these questions from repository evidence, failure-mode analysis, and focused experiments. The requirements intentionally do not preselect the answers.

1. What is the supported topology for this sprint: multiple processes on one host, multiple hosts sharing a workspace, or both? Which guarantees become conditional at the boundary?
2. Where should the execution-control responsibility live so it is not owned by `internal/web`, yet does not become a generic workflow engine or duplicate agentwrap?
3. What is the smallest durable representation that satisfies ordering, atomicity, concurrency, recovery, inspection, and bounded retention: filesystem snapshots plus append-only segments, SQLite, an agentwrap-backed store, or another local mechanism?
4. Is there a single coordinator/daemon, direct multi-process repository access with leases, or a hybrid? How is leader/owner identity established and recovered without a single hidden in-memory authority?
5. Does the process that accepts a run always own the worker? If an owner exits, should work be cancelled, adopted, or merely reconciled as interrupted, and which stages are safe to adopt?
6. What identifiers and hierarchy are needed between user operation, product run, stage execution, task attempt, agentwrap run, provider session, OS process, and external smoke-harness run?
7. What lease duration, heartbeat cadence, clock assumptions, fencing token, and PID/process-birth checks make liveness truthful without declaring slow healthy work dead?
8. What delivery guarantee is required for progress events, how are duplicate events handled, and what are the retention, compaction, snapshot, cursor-expiry, and replay-gap semantics?
9. What backpressure policy protects execution and disk while telling operators exactly which detail was sampled, compacted, or dropped?
10. Which run fields and historical events are readable across browser sessions, and what fresh authorization is required for cancel, retry, or other mutations?
11. How should legacy `/operations/{id}` URLs, API clients, current lock files, `.flow-state.json`, `.run-state.json`, and runtime session checkpoints map into the new model without fabricating history?
12. Which diagnostics are always local, which metrics endpoint or command is appropriate, and whether OpenTelemetry export should be built now, provided as an optional adapter, or deferred?
13. What happens when durable record persistence fails before start, mid-run, or at terminal commit? Which degraded modes are safe, and when must execution stop?
14. What data can be retained by default without exposing prompts, provider payloads, repository content, credentials, or sensitive paths, and how can users inspect redaction decisions?

Decisions that materially affect public API compatibility, authority, topology, security, or recoverability must be recorded in reasoning before implementation planning. The plan must include migration, rollback, and fault-injection work rather than treating them as follow-up polish.

## Non-Goals

- Hosted SaaS, public network exposure, team accounts, or multi-user authorization.
- A remote-worker protocol unless the supported-topology reasoning proves it is necessary for the stated acceptance criteria.
- Replacing canonical Markdown, JSON flow state, Git/source workspaces, or external smoke evidence with a general database-backed product model.
- Selecting SQLite, Postgres, a daemon, WebSockets, OpenTelemetry, or any broker by architectural preference alone.
- Persisting unredacted provider-native streams, full prompts, arbitrary stdout/stderr, or unlimited event history.
- Building a general-purpose scheduler, workflow engine, issue tracker, or distributed queue.
- Automatic Git mutation, automatic repair, content identity, retrieval, or knowledge-graph expansion.
- Making browser disconnect synonymous with cancellation.

## Constraints

- The run must exist independently of the process that presents it, the browser session that requested it, the transport used to observe it, and any single in-memory subscriber hub.
- Existing product modules continue to own their workflow semantics, artifacts, locks, and terminal outcomes. The shared run capability records and projects execution; it must not decide stage-specific success.
- `internal/web` remains an interface adapter. It may own HTTP/SSE connection state and bounded delivery buffers, but not durable run identity, lifecycle authority, or cross-process discovery.
- Agentwrap remains the runtime-supervision boundary. UltraPlan correlates its product runs with agentwrap runs and consumes safe canonical events rather than duplicating provider/process supervision.
- Commands and cancellation remain explicit operations; event streaming remains an observation channel.
- Required operational state commits use atomic/durable semantics appropriate to the chosen store. A successful HTTP response or visible running badge must not precede the authoritative acceptance record.
- Default behavior remains local-first and secure. Any expanded listen, sharing, or cross-host assumption requires a separate explicit security review.
- Public JSON/API changes are versioned or compatibility-preserving, with typed errors and documented retention behavior.

## Dependencies

### Planning Inputs

- Sprint 31 guarded operations, SSE, cancellation, locking, reconnect, and shutdown behavior.
- Sprint 32 browser hardening, API compatibility, recovery, and interface-state agreement fixtures.
- Sprint 34 shared-context real-runtime dogfood, which supplied the cross-surface operational evidence for this sprint.
- Current implementation-repository plans, re-read during reasoning rather than copied from this requirement:
  - `../ultraplan-go/docs/plans/integrated-roadmap.md`
  - `../ultraplan-go/docs/plans/ultraplan-local-server-experiment-plan.md`
  - `../ultraplan-go/docs/plans/server-shutdown-run-cancellation-contract.md`
- Current agentwrap event, observing runtime, run-store, cancellation, session, cleanup, and process-identity capabilities.
- Existing workspace lock, sprint flow state, execute run state, study run state, runtime session checkpoint, and external smoke evidence formats.

## Review Expectations

Review must verify the chosen model against the failure matrix, not only the happy-path UI:

- trace a run from acceptance through owner lease, runtime attempt, durable events, surface projections, and terminal arbitration
- demonstrate which component is authoritative at every transition and how stale writers are fenced
- prove identical active counts and lifecycle state across CLI, JSON, TUI, dashboard, and detail pages
- exercise two local web servers plus CLI activity against one workspace
- kill and restart owners and observers at each commit boundary and inspect recovery
- test session expiry, old URLs, replay cursors, retention gaps, slow consumers, corrupt/torn records, disk-full behavior, and redaction
- show that existing artifact and flow-state behavior remains compatible and authoritative
- document any topology limitation as a precise supported contract, not an implementation footnote
- include operator-facing evidence: correlated logs, metrics/diagnostics, reconciliation output, and an exportable redacted support bundle
