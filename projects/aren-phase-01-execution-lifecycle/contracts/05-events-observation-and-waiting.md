# Aren Events, Observation And Waiting Contract

> Project: `aren-phase-01-execution-lifecycle`  
> Status: governing Phase 1 contract; candidate for long-lived promotion after Phase 1 review

## Purpose

This contract defines canonical lifecycle history, replay, waiting, and independent observation.

The central rule is that Aren has one canonical per-run lifecycle record. Readers observe that record; they do not become producers, owners, acknowledgers, or controllers of execution.

## Requirement Index

| ID | Title | Severity If Violated |
| --- | --- | --- |
| AREN-OBS-001 | Canonical history belongs to the run | Blocker |
| AREN-OBS-002 | Event identity is `(run_id, sequence)` | High |
| AREN-OBS-003 | Lifecycle event vocabulary is closed | High |
| AREN-OBS-004 | Events describe committed truth | Blocker |
| AREN-OBS-005 | Recording is not delivery | High |
| AREN-OBS-006 | Observation cannot control execution | Blocker |
| AREN-OBS-007 | Readers have independent progress | High |
| AREN-OBS-008 | Cursor semantics are explicit | High |
| AREN-OBS-009 | Replay does not reconstruct execution | High |
| AREN-OBS-010 | Waiting is reusable observation | Blocker |
| AREN-OBS-011 | Observation completion is suffix-relative | High |
| AREN-OBS-012 | Lifecycle publication is defensively immutable | High |

## Requirements

### AREN-OBS-001 — Canonical History Belongs To The Run

Lifecycle events are recorded directly into one retained, in-memory, per-run canonical history by lifecycle authority.

A delivery mechanism, subscriber queue, CLI renderer, transport stream, log sink, or later telemetry system is not canonical history.

Phase 1 introduces no global event bus.

### AREN-OBS-002 — Event Identity Is `(run_id, sequence)`

Every canonical lifecycle event contains the run identity and a per-run sequence number.

Requirements:

- `run.created` has sequence `0`;
- subsequent events use contiguous increasing sequence numbers;
- sequence allocation and event commitment occur under the same lifecycle authority;
- ordering is defined by sequence, not timestamps.

Consumers may use `(run_id, sequence)` to identify repeated exposure of the same canonical event.

### AREN-OBS-003 — Lifecycle Event Vocabulary Is Closed

The Phase 1 lifecycle vocabulary is:

```text
run.created
run.started
run.cancellation_requested
run.succeeded
run.failed
run.cancelled
```

`run.cancellation_requested` records an occurrence while state remains `running`.

Creation and start occur once. Cancellation request occurs at most once. A completed run ends with exactly one terminal event.

Do not overload this canonical lifecycle stream with future token deltas, tool logs, arbitrary progress, telemetry spans, or transport housekeeping.

### AREN-OBS-004 — Events Describe Committed Truth

No event may announce a state, outcome, cancellation acceptance, or transition that the run did not actually commit.

The event is part of the lifecycle fact set, not a prediction or best-effort notification.

Observer delivery may lag canonical commitment. Canonical recording may not lag behind another public fact that claims the event's lifecycle meaning is already true.

### AREN-OBS-005 — Recording Is Not Delivery

Each canonical lifecycle occurrence is recorded once.

Replay or repeated reads may expose the same canonical event more than once.

Aren does not claim exactly-once observer delivery and does not require consumer acknowledgements in Phase 1.

Intentional replay is not duplicate canonical recording.

### AREN-OBS-006 — Observation Cannot Control Execution

Slow, absent, cancelled, or abandoned readers must not block or determine:

- work execution;
- cancellation acceptance or propagation;
- lifecycle transitions;
- terminal publication;
- waiters;
- other readers.

Do not place observer callbacks or blocking delivery inside lifecycle mutation ownership.

### AREN-OBS-007 — Readers Have Independent Progress

Each reader owns its own observation cursor or equivalent progress state.

Reading, cancelling, replaying, or abandoning one reader must not:

- advance another reader;
- mutate canonical history;
- mutate lifecycle state;
- cancel the run;
- consume the only copy of an event.

Concurrent mutation of one reader object by several goroutines is not automatically promised. If a concrete API offers that capability, it must be explicitly verified.

### AREN-OBS-008 — Cursor Semantics Are Explicit

A cursor names the **next sequence requested**.

Let `n` be the current committed history length and `c` the requested cursor.

```text
Default
    c = 0

0 <= c < n
    read inclusively from c

c = n and run is active
    wait for availability or observation-local cancellation

c = n and run is terminal
    normal exhaustion

c > n or invalid representation
    explicit invalid-cursor result
```

Do not silently clamp an invalid future cursor.

When event sequence `s` is returned, that traversal advances to `s + 1`.

Local observation cancellation without returned data does not advance the cursor and does not mutate the run.

### AREN-OBS-009 — Replay Does Not Reconstruct Execution

Replay reads retained canonical facts.

Replay does not:

- re-execute work;
- reconstruct authoritative lifecycle state from history;
- append events;
- create a new run;
- claim durable recovery.

Phase 1 replay exists only while the run and its in-memory history remain reachable.

### AREN-OBS-010 — Waiting Is Reusable Observation

Waiting observes independently completed terminal publication.

Requirements:

- wait blocks while the run is active;
- wait returns only after the complete terminal fact set is published;
- late wait returns immediately after completion;
- multiple waiters receive the same logical terminal outcome;
- waiting does not finalize the run;
- waiting does not exclusively consume the outcome;
- waiting does not depend on observer progress.

A context-cancellable wait may be introduced if it is actually needed, but local wait cancellation must remain distinct from run cancellation.

### AREN-OBS-011 — Observation Completion Is Suffix-Relative

A live reader completes normally when:

1. the run is terminal; and
2. the reader has exhausted the requested suffix.

If a reader starts at the current tail of an already-terminal run, its requested suffix is empty and may complete immediately. The terminal event must nevertheless already exist in canonical history and remain retrievable from an earlier cursor.

A reader must not report normal exhaustion merely because it temporarily reached the active tail.

### AREN-OBS-012 — Lifecycle Publication Is Defensively Immutable

Canonical lifecycle events and their Aren-owned payloads must not be mutable through supported public aliases.

Where payloads contain maps, slices, pointers, stack data, or other reference-bearing Aren-owned structures, the API must use defensive copies or another proved immutable-publication strategy.

This requirement does not extend deep immutability to arbitrary user result objects as defined by `03-execution-lifecycle.md`.

## Preferred Phase 1 Direction

Pull-oriented reading of retained history is preferred because it avoids making one producer queue per observer part of the correctness model.

This is a design preference, not a mandatory public API shape. Another mechanism may be used if it proves the same canonical-history, isolation, handoff, and abandonment guarantees without adding speculative machinery.

## Primary Governing Sources

- `projects/aren-phase-01-execution-lifecycle/docs/PRD.md`
- `projects/aren-phase-01-execution-lifecycle/project-reasoning/reasoning.md`
- `projects/aren-phase-01-execution-lifecycle/project-reasoning/areas/04-events-observation-waiting-and-replay.md`
