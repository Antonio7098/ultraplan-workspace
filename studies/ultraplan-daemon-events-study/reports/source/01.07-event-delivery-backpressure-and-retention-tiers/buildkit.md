# Source Analysis: buildkit

## Event Delivery, Backpressure, and Retention Tiers

### Source Info

| Field | Value |
|-------|-------|
| Name | buildkit |
| Path | `studies/ultraplan-daemon-events-study/sources/buildkit` |
| Language / Stack | Go (gRPC daemon `buildkitd`, bbolt + containerd content-store, `golang.org/x/time/rate`) |
| Analyzed | 2026-09-03 |

## Summary

BuildKit splits events into two tiers: an ephemeral, non-blocking in-memory progress tier (`util/progress` pipe → `MultiWriter` fan-out → `MultiReader` replay → gRPC `Status` stream) and a durable build-history tier (`solver/llbsolver/history.Queue` in bbolt + content-store blobs). Hot-path writes never block on subscribers (`util/progress/progress.go:241-247` only locks, overwrites `dirty[ID]`, broadcasts). High-volume logs are clipped by size/speed caps with a 256 KiB circbuf tail (`util/progress/logs/logs.go:22-23,106-164`), gRPC payloads are split at ~1 MiB (`client/status.go:95-111`), and the TUI throttles rendering to ~100–150 ms (`util/progress/progressui/display.go:81-92`). Durable lifecycle facts (STARTED/COMPLETE records, length-prefixed status blobs, error blobs, provenance/traces) are committed to bbolt/leases/content-store at build end (`solver/llbsolver/history.go:59-335`, `solver/llbsolver/history/buildhistory.go:418-551,679-768`). Retention is bounded only for history (default 50 entries / 48 h, GC every 120 s, `solver/llbsolver/history/buildhistory.go:80-86,131-198`); the in-memory progress buffers (`MultiWriter.items`, `MultiReader.sent`, `progressReader.dirty` for unique log IDs, history `pubsub` goroutine-per-send) grow without bound and have no lag/queue-depth metrics.

## Rating

**5/10 — Present but inconsistent, weakly documented, fragile under flood.**

Rationale: coalescing by ID, non-blocking producer path, log clipping, gRPC chunking, display rate-limiting, and bbolt+content-store durability with pinned-entry-aware GC are real and tested (`util/progress/progress_test.go`, `solver/llbsolver/history/pubsub_test.go`). Downgraded because (a) `MultiWriter.items` and `MultiReader.sent` are append-only with no pruning, (b) history `pubsub.Send` spawns an unbounded goroutine per subscriber per event and blocks that goroutine when the 32-slot channel fills, (c) there are no metrics for progress lag, queue depth, or dropped/coalesced counts, and (d) durable history writes happen post-hoc in `recordBuildHistory` rather than as a WAL on the hot path, so a crash mid-build loses in-flight progress.

## Evidence Collected

Every entry MUST include a file path with line numbers. Format: `path/to/file.ts:NN`.

| Area | Evidence | File:Line |
|------|----------|-----------|
| Coalescing buffer (dirty map keyed by ID, unread writes collapse) | `progressReader.dirty map[string]*Progress`; `Read` drains whole map, sorts by timestamp | `util/progress/progress.go:110-174` |
| Non-blocking hot-path write | `writeRawProgress` locks, `dirty[p.ID]=p`, `Broadcast`, `return nil` — never waits for reader | `util/progress/progress.go:241-247` |
| Writer lifecycle / EOF signal | `Close` removes writer from set and broadcasts; `Read` returns `io.EOF` only when writers==0 and dirty empty | `util/progress/progress.go:149-159,249-256` |
| Fan-out registry + replay on join | `MultiWriter.writers` set; `Add` replays sorted `items` to late joiner; `writeRawProgress` appends to `items` then synchronously fans out under one mutex | `util/progress/multiwriter.go:34-96` |
| Unbounded fan-out history | `ps.items = append(ps.items, p)` with no cap, prune, or compaction | `util/progress/multiwriter.go:86-96` |
| Per-subscriber replay + unbounded `sent` log | `MultiReader.sent` accumulates every batch (`sent=append(sent,p...)`); late `Reader()` replays from index `i` in 100-item chunks | `util/progress/multireader.go:29-137` |
| Slow-subscriber isolation (progress tier) | New `Reader()` gets its own `progressWriter`+closer; `handle` loop writes to each `w` synchronously under `mr.mu` — one slow `writeRawProgress` (which is itself non-blocking) does not stall others, but lock hold extends | `util/progress/multireader.go:29-137` |
| Job-level fan-out wiring | `Job.pr *progress.MultiReader`, `Job.pw progress.Writer`; `NewJob` creates `NewContext` + `NewMultiReader`; per-vertex `mpw: NewMultiWriter(WithMetadata("vertex",dgst))`; `connectProgressFromState` attaches each job `pw` once and writes `clientVertex` | `solver/jobs.go:373-389,583-594,665-715` |
| Status forwarding blocks on consumer channel | `Job.Status` loops `pr.Read` → builds `SolveStatus` → `select ctx.Done vs ch <- ss`; a full/blocked `ch` parks this goroutine while execution keeps writing (non-blocking) | `solver/progress.go:17-103` |
| gRPC Status bridge, buffered chan 8 + drain | `ch := make(chan *client.SolveStatus, 8)`; solver→sender via errgroup; sender `stream.SendMsg` applies gRPC flow control; `defer drain ch` on error | `control/control.go:594-625` |
| gRPC message-size bound (batch writer split) | `SolveStatus.Marshal` accumulates `logSize`; if `>1MiB` it emits current `StatusResponse`, clears vertexes/statuses, continues logs in next message | `client/status.go:66-129` |
| Client Status receive pushes to caller chan (unbounded expectation) | `stream.Recv` loop → `statusChan <- NewSolveStatus(resp)` with no select/drop; caller-supplied chan backpressures Recv loop | `client/solve.go:369-397` |
| Log drop/coalesce policy (size+speed clip + tail) | Defaults `2MiB` max size, `200KiB/s` max speed (env-overridable `BUILDKIT_STEP_LOG_MAX_SIZE/_SPEED`); over-limit bytes go to 256 KiB circbuf, clipped marker appended, `flushBuffer` writes tail at close | `util/progress/logs/logs.go:22-23,62-164` |
| Per-write log identity defeats ID coalescing | Each log chunk uses `identity.NewID()` as progress ID, so `dirty` map cannot collapse log flood — one entry per chunk until drained | `util/progress/logs/logs.go:137-144` |
| Display-side sampling/throttle | `UpdateFrom`: 150 ms ticker, `rate.Every(100ms)` limiter, `TTY_DISPLAY_RATE` override; `update()` buffers, `refresh()` paints | `util/progress/progressui/display.go:81-118` |
| History pubsub (bounded chan, goroutine-per-send) | `ch: make(chan T, 32)`; `Send` holds `mu` and `go c.send(v)` per subscriber; `send` blocks on `ch <- v` until consumed or `done` | `solver/llbsolver/history/pubsub.go:10-54` |
| Durable record write path | `Queue.Update` handles STARTED (in-memory `active` + pubsub) vs COMPLETE (bbolt `update` + pubsub); `update` marshals record, creates `ref_<id>` lease, pins logs/trace/error/result blobs, `b.Put` | `solver/llbsolver/history/buildhistory.go:531-551,418-487` |
| Durable status-blob batch writer | `ImportStatus` drains `ch chan *SolveStatus`, counts cached/completed/total/warnings, frames each `StatusResponse` as `u32len+proto` via `bufio.Writer` into content-store blob, commits to descriptor | `solver/llbsolver/history/buildhistory.go:679-768` |
| Durable error-blob writer | `ImportError` converts to gRPC status proto, writes to history content-store blob | `solver/llbsolver/history/buildhistory.go:637-677` |
| History replay (late-join durability) | `Queue.Status` reads bbolt record, opens `br.Logs` blob from history ns, decodes `u32len+proto` frames, re-emits `client.SolveStatus` on `st` chan; `Solver.Status` prefers history, falls back to live `Job.Status` | `solver/llbsolver/history/buildhistory.go:351-416`, `solver/llbsolver/solver.go:465-481` |
| History live listen + active snapshot | `Listen` subscribes to pubsub, snapshots `active`, replays bbolt COMPLETE events (filtered/limited), then streams live `sub.ch` until `EarlyExit`/ctx-done; ref-counts `refs`/`deleted` guard in-flight deletes | `solver/llbsolver/history/buildhistory.go:770-888` |
| Retention config + defaults | `HistoryConfig{MaxAge, MaxEntries}` TOML keys; default `48h / 50` when nil | `cmd/buildkitd/config/config.go:222-225`, `solver/llbsolver/history/buildhistory.go:80-86` |
| Retention/cleanup job | Goroutine `clearOrphans(); for { gc(); sleep 120s }`; `gc` skips pinned, requires `len >= MaxEntries`, sorts newest-first, deletes tail older than `MaxAge`; `delete` emits DELETED event, removes bbolt row + lease, `Delete` triggers `GarbageCollect` | `solver/llbsolver/history/buildhistory.go:131-137,152-198,200-266,553-565` |
| Graceful-stop close of pubsub | On `GracefulStop` with no active/finalizers, `ps.Close()`; finalizer completion also closes when last finalizer drains | `solver/llbsolver/history/buildhistory.go:139-147,489-529` |
| Job retention window for late Status | `Discard` keeps job entry 10 s (`time.Sleep(10s)` then delete) so a trailing `Status` can still read progress | `solver/jobs.go:859-898` |
| Vertex lifecycle terminal facts | `notifyStarted` writes Started vertex, returns closure writing Completed/Error + `pw.Close`; controller `Start`/`Status` emit Started/Completed `client.Vertex`/`progress.Status` pairs | `solver/jobs.go:1368-1388`, `util/progress/controller/controller.go:30-96` |
| GC backpressure guard (unrelated but bounding) | `throttledGC = throttle.After(1m)`, `throttledReleaseUnreferenced = throttle.After(5m)`; `Solve` schedules GC 1 s after return | `control/control.go:143-157,428-430,677-717` |
| Tests: coalescing expectation | `TestProgress` asserts `5 < len(items) <= 7` despite more writes — proves ID-collapse | `util/progress/progress_test.go:18-49` |
| Tests: pubsub fan-out/close/concurrency | Send/receive to 2 subs, close semantics, idempotent close, 10 subs × 100 msgs | `solver/llbsolver/history/pubsub_test.go:10-103` |

## Answers to Dimension Questions

1. **Which events must never be dropped?**
   STARTED/COMPLETE `BuildHistoryRecord`s (bbolt row + `ref_<id>` lease + content blobs for logs/trace/error/results/attestations) — `solver/llbsolver/history/buildhistory.go:418-551,679-768`. Terminal `client.Vertex` Completed/Error pairs (`solver/jobs.go:1368-1388`, `util/progress/controller/controller.go:51-75`) are durable only insofar as `ImportStatus` drains the `Job.Status` channel to the status blob in `solver/llbsolver/history.go:185-206` before `Update(COMPLETE)`. Live progress delivery itself is lossy by design (ID-collapse, log clip, display throttle); durability comes from the history commit, not the stream.

2. **Which events are ephemeral or short-lived?**
   Everything in `util/progress`: per-ID `Status` updates (collapsed in `dirty`), per-chunk `VertexLog`s (unique IDs, clipped), `VertexStatus` deltas, and gRPC `StatusResponse`s. Also `MultiWriter.items`/`MultiReader.sent` replays, `pubsub` 32-slot channels, and the 10 s post-`Discard` job window (`solver/jobs.go:890-896`). `control/control.go:598` Status `ch` (cap 8) and client `statusChan` (`client/solve.go:393-395`) are transient handoffs.

3. **Can a slow UI stall execution?**
   No for the build itself; yes for its own status tail. Producer `Write` never blocks (`util/progress/progress.go:241-247`); `MultiWriter.writeRawProgress` fans out synchronously but each downstream `WriteRawProgress` is also non-blocking (`util/progress/multiwriter.go:86-96`). A slow UI parks in `Job.Status` at `ch <- ss` (`solver/progress.go:97-101`) or in `stream.SendMsg` gRPC flow control (`control/control.go:616-620`), causing `dirty`/`sent`/`items` to accumulate. Execution (scheduler, cache, exporters) proceeds. For history `Listen`, a slow `srv.Send` blocks only that subscriber's `Listen` loop (`solver/llbsolver/history/buildhistory.go:873-887`); other subscribers have separate channels.

4. **How are lifecycle and terminal facts flushed under pressure?**
   Two mechanisms: (a) live: `notifyStarted` closure + `Controller` completion write + `pw.Close` signal EOF (`solver/jobs.go:1368-1388`, `util/progress/controller/controller.go:51-75`, `util/progress/progress.go:249-256`); `vertexStream.encore()` synthesizes Completed+Canceled for still-open vertices on EOF (`solver/progress.go:145-158`). (b) durable: `recordBuildHistory` calls `CloseProgress`, then concurrently drains `j.Status(ctx2, ch)` into `ImportStatus` (bounded only by the 300 s timeout at `solver/llbsolver/history.go:65-76,185-206`) and commits the bbolt COMPLETE record plus trace/error blobs. Log pressure is shed *before* flush via the 2 MiB/200 KiB/s clip + 256 KiB tail (`util/progress/logs/logs.go:106-164`); oversized gRPC batches are split at 1 MiB (`client/status.go:95-111`).

5. **What bounds storage and in-memory growth?**
   Storage: `HistoryConfig.MaxEntries` (default 50) + `MaxAge` (default 48 h), enforced every 120 s, both thresholds required, pinned exempt (`solver/llbsolver/history/buildhistory.go:152-198`); orphan-blob sweep (`clearOrphans`, `200-239`); lease-scoped content-store + `GarbageCollect` on delete (`553-565`); worker GC throttled (`control/control.go:143-145,677-717`). In-memory: log clip caps per-stream bytes, gRPC split caps message bytes, display limiter caps paint rate — but `MultiWriter.items` (`util/progress/multiwriter.go:89`), `MultiReader.sent` (`util/progress/multireader.go:134`), and `dirty` entries for unique log IDs (`util/progress/logs/logs.go:141`) are unbounded for the life of the job; history `pubsub` channels are capped at 32 but each blocked send parks a goroutine (`solver/llbsolver/history/pubsub.go:22-54`). No evidence of max-bytes, drop-oldest, or spill-to-disk for those buffers.

## Architectural Decisions

- **Lossy live tier + durable history tier.** Real-time progress is best-effort/coalesced/clipped; truth is reconstructed from the committed history record + status blob. Evidence: `util/progress/progress.go:148-162,241-247` vs `solver/llbsolver/history/buildhistory.go:351-416,679-768`.
- **Producer-never-blocks invariant.** All progress `Write` paths are lock+map-set+broadcast, shifting backpressure to memory growth and late-join replay cost rather than execution stalls.
- **Per-vertex `MultiWriter` fan-in, per-job `MultiReader` fan-out.** Shared scheduler state (`solver/jobs.go:577-594,665-685`) lets many jobs observe one vertex DAG while each job gets isolated replay (`util/progress/multireader.go:29-108`).
- **Post-hoc durability instead of WAL.** Status blobs are materialized by re-reading the live stream at completion (`solver/llbsolver/history.go:185-206`), not by journaling each write. Simpler, but crash-before-commit loses the tail.
- **Conjunctive retention (count AND age).** GC deletes only records beyond `MaxEntries` *and* older than `MaxAge` (`solver/llbsolver/history/buildhistory.go:175-195`), avoiding surprise deletion of recent bursts.
- **Lease-pinned blobs + isolated `_history` namespace.** History content/leases live in `<ns>_history` (`solver/llbsolver/history/buildhistory.go:98-105`), so worker GC cannot reap referenced blobs prematurely.

## Notable Patterns

- **ID-keyed collapse:** `dirty map[string]*Progress` is the coalescing policy — same ID overwrites, distinct IDs accumulate. Status updates reuse IDs; logs mint new IDs per chunk.
- **Replay-on-subscribe:** both `MultiWriter.Add` (sorted `items` replay) and `MultiReader.Reader` (`sent` replay with catch-up state machine) give late attach the full prefix.
- **Clip-and-tail for logs:** hard byte + byte-per-second ceilings with circbuf tail preservation and explicit `[output clipped, log limit … reached]` marker.
- **Chunk-at-1 MiB gRPC batching:** `SolveStatus.Marshal` splits one logical status into N `StatusResponse`s purely on encoded log size.
- **Rate-limited presentation:** TUI consumes every message (`update`) but paints at most every ~100 ms (`displayLimiter`) plus 150 ms `refresh` ticker.
- **Ref-counted delete deferral:** `refs`/`deleted` maps in `Listen`/`delete` prevent removing a history record while a listener holds it.

## Tradeoffs

- **Liveness over completeness:** never blocking the solver on UI/gRPC keeps builds fast but makes memory the shock absorber — a flooded logger or stalled viewer grows `dirty`/`sent`/`items` without feedback to the producer.
- **Simplicity of unbounded replay vs cost:** full `sent`/`items` replay makes late `buildctl` attach trivial but costs O(build events) RAM per daemon lifetime (per job for `sent`, per vertex-state for `items` — neither is pruned on `Close`, since `MultiWriter.Close` is a no-op at `util/progress/multiwriter.go:98-100`).
- **Goroutine-per-send pubsub:** `Send` never blocks the publisher (`solver/llbsolver/history/pubsub.go:22-28`) at the price of one goroutine per event per subscriber; a 32-deep slow consumer leaks goroutines linearly with events.
- **Post-hoc history materialization:** avoids WAL write amplification on every progress write, but a daemon crash or 300 s `ImportStatus` timeout window can lose the very terminal facts the history tier promises.
- **Conjunctive GC:** safe against burst deletion, but a sustained high-build-rate daemon can hold `MaxEntries`-plus-recent records indefinitely (storage grows with build rate × blob size until age threshold passes).

## Failure Modes / Edge Cases

- **Stalled viewer → memory growth, not build stall.** Attach a never-reading `Status` consumer: `ch (cap 8)` fills (`control/control.go:598`), `Job.Status` parks at `ch <- ss` (`solver/progress.go:97-101`), `dirty` (unique log IDs) + `MultiReader.sent` + `MultiWriter.items` grow; build completes and history commit still succeeds via the separate `ImportStatus` drain. First visible degradation is daemon RSS + goroutine count, then slower `Add`/replay for new subscribers (O(n) sorted replay under a single mutex at `util/progress/multiwriter.go:48-59`).
- **Log flood → silent truncation.** Sustained output above 200 KiB/s or 2 MiB total is clipped to the marker + final 256 KiB tail (`util/progress/logs/logs.go:106-164`); the durable status blob therefore contains a *truncated* log, with no sequence numbers to detect the gap downstream beyond the marker string.
- **History pubsub slow consumer → goroutine pile-up.** Each `Send` to a full 32-slot `channel` parks its goroutine in `send` (`solver/llbsolver/history/pubsub.go:49-54`); `Close` only closes `done`, it does not drain parked senders. Long-lived `ListenBuildHistory` with a stalled client leaks ~1 goroutine per build event.
- **Crash between COMPLETE and blob commit → orphaned/partial record.** `update` (bbolt row) and blob commits are not atomic (`solver/llbsolver/history/buildhistory.go:418-487`); `clearOrphans` deletes records whose lease resources vanished (`200-239`), which can discard a record that is otherwise valid.
- **No backpressure observability.** No counters for coalesced writes, clipped bytes, pubsub drops/parks, `sent`/`items` sizes, or Status chan occupancy were found; operators cannot alert on the exact signals this dimension requires.
- **Multiwriter cycle panic as backpressure-adjacent hazard.** Edge-merge cycles trigger `panic("multiwriter loop detected")` (`util/progress/multiwriter.go:39-46`) — a fail-fast choice that converts a wiring bug into a daemon crash.

## Future Considerations

- Bound `MultiWriter.items` and `MultiReader.sent` (e.g., ring with “replay-from-checkpoint” or spill status blobs to content-store incrementally instead of holding all in RAM); make `MultiWriter.Close` actually release `items`.
- Replace `pubsub.Send` goroutine-per-event with non-blocking send + drop-oldest + `dropped` counter, or per-subscriber ring with lag metric; add `queue_depth`, `dropped_total`, `replay_bytes` metrics.
- Add sequence numbers/byte offsets to `VertexLog` framing so truncation (log clip, 1 MiB gRPC split, status-blob replay) is detectable by consumers.
- Consider incremental/WAL durability for terminal vertex facts (fsync on Completed/Error) so crash-before-`ImportStatus` does not lose lifecycle truth; or document the at-most-once window explicitly.
- Expose history GC stats (scanned/deleted/pinned/skipped, oldest retained age) and progress-buffer gauges for dashboards.

## Questions / Gaps

- No evidence found for WAL/log growth metrics, progress lag histograms, or queue-depth gauges — searched `solver/`, `util/progress/`, `util/throttle/`, `solver/llbsolver/metrics.go`, `control/control.go` GC logging.
- No evidence found for a drop/coalesce policy on `VertexStatus` beyond ID-collapse in `dirty` — searched `util/progress/*.go`, `solver/progress.go`, `util/progress/progressui/`.
- No evidence found for compaction of `MultiWriter.items`/`MultiReader.sent` or retention limits on in-memory progress — `Close` paths (`util/progress/multiwriter.go:98-100`, `solver/jobs.go:854-898`) do not prune them.
- Open: what is the intended max `StatusResponse` count per `SolveStatus.Marshal` split under adversarial log sizes? The loop at `client/status.go:66-129` is unbounded in `out` length.
- Open: is the 300 s `ImportStatus` timeout in `solver/llbsolver/history.go:75` sufficient for very large status blobs, and what happens to the COMPLETE record if it fires mid-drain?

---
Generated by `Dimension 01.07: Event Delivery, Backpressure, and Retention Tiers` against `buildkit`.
