# Aren Cancellation And Lifetimes Contract

> Project: `aren-phase-01-execution-lifecycle`  
> Status: governing Phase 1 contract; candidate for long-lived promotion after Phase 1 review

## Purpose

This contract defines truthful cooperative cancellation and ownership of Aren-created runtime support.

Cancellation is not a boolean and it is not equivalent to `context.Done()`. Aren must preserve the distinction between a stop request, accepted cancellation facts, propagation to work, work completion, terminal interpretation, and release of Aren-owned support resources.

## Requirement Index

| ID | Title | Severity If Violated |
| --- | --- | --- |
| AREN-CANCEL-001 | Cancellation is a protocol | Blocker |
| AREN-CANCEL-002 | One acceptance path | Blocker |
| AREN-CANCEL-003 | First accepted cause wins | High |
| AREN-CANCEL-004 | Acceptance precedes work-visible signal | Blocker |
| AREN-CANCEL-005 | Request is not terminality | Blocker |
| AREN-CANCEL-006 | Uncooperative work remains active | Blocker |
| AREN-CANCEL-007 | Aren-owned goroutines have owners | High |
| AREN-CANCEL-008 | User goroutines are not automatically adopted | High |
| AREN-CANCEL-009 | Terminal publication and support quiescence are distinct | High |
| AREN-CANCEL-010 | Teardown does not rewrite history | High |
| AREN-CANCEL-011 | Already-cancelled parents still follow the lifecycle | High |
| AREN-CANCEL-012 | Cancellation dispositions are serialized facts | High |

## Requirements

### AREN-CANCEL-001 — Cancellation Is A Protocol

Keep these boundaries conceptually distinct:

```text
request
    -> acceptance
    -> signal propagation
    -> possible work observation
    -> completed invocation
    -> terminal resolution/publication
    -> support release
```

Do not collapse them into one generic `cancelled` fact.

Aren may not always know whether work observed the signal. A matching returned error is contractual acknowledgement under the selected terminal policy, not proof of physical causation.

### AREN-CANCEL-002 — One Acceptance Path

Explicit controller cancellation, parent-context cancellation, and parent deadline expiry use one Aren-owned acceptance path.

A parent watcher, context helper, transport, or work function must not become a second acceptance authority.

Observation and waiting do not grant cancellation authority.

### AREN-CANCEL-003 — First Accepted Cause Wins

The first cancellation request accepted while the run is active fixes:

- the retained reason;
- the effective cancellation cause;
- the canonical cancellation-request occurrence.

Later requests must not:

- replace the cause or reason;
- append another cancellation-request event;
- restart cancellation support;
- allocate another event sequence solely for the repeated request.

### AREN-CANCEL-004 — Acceptance Precedes Work-Visible Signal

Canonical acceptance facts must exist before Aren makes the accepted cancellation signal visible to work.

The accepting operation must settle its required propagation before reporting `accepted`, and terminal publication must not overtake accepted propagation that is still outstanding.

User callbacks, blocking joins, or external work must not run while lifecycle mutation ownership is held merely to satisfy this ordering.

### AREN-CANCEL-005 — Request Is Not Terminality

Aren must not publish terminal cancellation merely because:

- cancellation was requested;
- cancellation was accepted;
- the work context is done;
- a deadline elapsed;
- an observer stopped waiting.

Terminal cancellation requires a completed invocation and the terminal-resolution rule defined in `03-execution-lifecycle.md`.

### AREN-CANCEL-006 — Uncooperative Work Remains Active

If work ignores cancellation and does not return, Aren does not force-kill the goroutine and must not falsely claim that the work stopped.

The run remains active with its accepted cancellation facts visible.

Phase 1 defines no maximum stop latency and no goroutine kill mechanism.

### AREN-CANCEL-007 — Aren-Owned Goroutines Have Owners

Every Aren-owned goroutine or asynchronous carrier has:

- a named owner;
- a creation point;
- a lifetime;
- a stop condition;
- a cancellation source where applicable;
- a join or release path;
- an error-path release story.

`context.Background()` is permitted only at an explicit lifetime boundary that truly owns work beyond a caller lifetime. Phase 1 has no such daemon-owned execution boundary.

Fire-and-forget Aren work is prohibited unless its bounded lifetime and irrelevance to correctness are explicitly documented and verified.

Prefer eliminating support goroutines when a simpler mechanism provides the same semantics.

### AREN-CANCEL-008 — User Goroutines Are Not Automatically Adopted

Aren owns the supervised invocation carrier and the support resources it creates.

Goroutines independently spawned by user work are not automatically adopted, joined, or made part of the run's quiescence guarantee.

Aren must not imply otherwise in its lifecycle or documentation.

### AREN-CANCEL-009 — Terminal Publication And Support Quiescence Are Distinct

Terminal publication proves:

- the supervised work invocation has completed;
- applicable cancellation propagation has settled;
- terminal lifecycle facts are complete and published.

It does not automatically prove that every Aren support carrier has already exited before every waiter returns.

Any post-terminal Aren-owned support must have a finite release path under ordinary scheduling and must not depend on:

- observer drainage;
- another future lifecycle event;
- new cooperation from already-finished user work;
- unrelated parent cancellation.

Actual release must be independently verifiable.

### AREN-CANCEL-010 — Teardown Does Not Rewrite History

Releasing contexts, unregistering parent integration, or stopping internal support after work completion is teardown.

Teardown must not:

- create another accepted cancellation request;
- append cancellation history;
- change the terminal outcome;
- replace the first accepted cause.

General user cleanup hooks and cleanup-result aggregation remain outside Phase 1.

### AREN-CANCEL-011 — Already-Cancelled Parents Still Follow The Lifecycle

Starting with an already-cancelled parent still follows the canonical creation and start lifecycle.

Aren must:

1. create the run and record creation;
2. commit start;
3. accept the parent cancellation through the canonical path;
4. settle propagation before invoking work;
5. invoke the supplied work once with a cancelled context;
6. resolve the work return using the ordinary terminal policy.

An already-cancelled parent does not create a special direct `created -> cancelled` edge.

### AREN-CANCEL-012 — Cancellation Dispositions Are Serialized Facts

A cancellation request reports one of:

- `accepted`;
- `already_requested`;
- `already_terminal`.

The disposition describes the request's decision under lifecycle authority.

A request that returns `accepted` does not promise that the run is still active by the time the caller receives the response; work may have completed immediately after acceptance.

## Relationship To Terminal Resolution

Terminal cancellation uses the effective cause retained by first acceptance and the exact rule in `03-execution-lifecycle.md`:

```text
cancellation accepted
AND returned error is non-nil
AND errors.Is(returnedError, acceptedEffectiveCause)
    -> cancelled
```

No generic `context.Canceled` fallback is implied for a distinct custom cause.

## Primary Governing Sources

- `projects/aren-phase-01-execution-lifecycle/docs/PRD.md`
- `projects/aren-phase-01-execution-lifecycle/docs/final-language-decision.md`
- `projects/aren-phase-01-execution-lifecycle/project-reasoning/reasoning.md`
- `projects/aren-phase-01-execution-lifecycle/project-reasoning/areas/03-cancellation-goroutine-ownership-and-cleanup.md`
