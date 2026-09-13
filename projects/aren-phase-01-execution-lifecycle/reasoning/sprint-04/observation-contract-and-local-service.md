# Area reasoning template: Observation contract and local service

## Purpose

Define the versioned read-only contract and process-scoped loopback service used by the Phase 1 browser frontend. Preserve runtime authority. The service projects lifecycle and performance evidence but owns none of it.

## Required inputs

- completed Sprint 1 through Sprint 3 artifacts and realised code;
- Phase 1 PRD, observability mandate, and performance mandate;
- selected event, replay, API, security, privacy, and frontend evidence.

## Questions to resolve

1. Which runtime-owned facts enter the transport DTO, and which remain internal?
2. How does the schema distinguish canonical facts, diagnostics, and performance evidence?
3. How are versions, optional fields, unknown event kinds, and incompatible data handled?
4. How do initial snapshot and live updates meet without loss, duplication, or reordering?
5. What bounds apply to payloads, clients, update rate, buffers, and reconnection?
6. How does the service bind locally, stop with its owner, and avoid becoming a daemon?
7. How are sensitive values omitted, redacted, or represented structurally?

## Required analysis

Trace each displayed fact back to its canonical owner. Define a read path for active and completed runs, sequence-based live delivery, late connection, client abandonment, and terminal completion. Reject any design where transport state resolves outcomes or browser backpressure reaches lifecycle commitment.

## Required output

- authority and data-flow diagram;
- versioned DTO definitions;
- snapshot and live-delivery protocol;
- bounds and overload behaviour;
- loopback and lifecycle rules;
- privacy and diagnostic policy;
- compatibility and error semantics;
- rejected daemon and global-event-stream alternatives.

