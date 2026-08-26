# Source Analysis: conc

## 01.01 Lifecycle Transition Ownership and Terminal Arbitration

### Source Info

| Field | Value |
|-------|-------|
| Name | conc |
| Path | `studies/aren-go-runtime-study/sources/conc` |
| Language / Stack | Go 1.20 (`go.mod:3`); stdlib-only runtime (`sync`, `context`, `sync/atomic`), testify for tests (`go.mod:5`) |
| Analyzed | 2026-08-26 |

Citation convention: all `sources/conc/...` paths below are relative to the study root `studies/aren-go-runtime-study/`.

## Summary

conc (sourcegraph/conc) is a scoped-concurrency library built around one idea that maps directly onto Aren's lifecycle question: **every goroutine has exactly one owner, and the owner's `Wait()` call is the single commit point where execution state, outcomes, and failure notifications become visible together**. The stated design goal is explicit: "all concurrency should be scoped... goroutines should have an owner and that owner should always ensure that its owned goroutines exit properly" (`sources/conc/README.md:55-62`).

There is no explicit lifecycle state machine — no running/failed/canceled enums, no transition events. Instead, ownership and terminal arbitration are enforced structurally:

- The primitive `WaitGroup` spawns goroutines via `Go()` (`sources/conc/waitgroup.go:28-34`) and gates on their exit in `Wait()`, then re-throws any caught panic (`sources/conc/waitgroup.go:38-43`). A panic crossing a goroutine boundary is captured by a per-owner `panics.Catcher` whose storage is an atomic pointer with first-wins CAS arbitration (`sources/conc/panics/panics.go:15-31`).
- Work errors are a separate channel: `ErrorPool` collects returned errors into a mutex-guarded append-only slice published by `Wait()` (`sources/conc/pool/error_pool.go:27-48,94-100`). Panics and returned errors never merge; a panicking task skips error/result recording entirely and terminates the batch as a re-thrown panic.
- Cancellation is delegated to stdlib `context`: `ContextPool` passes a derived ctx to tasks (`sources/conc/pool/context_pool.go:24-50`) and can cancel-on-first-error-or-panic (`sources/conc/context_pool.go:26-47`). Cancellation *request* (ctx signal) and confirmed *termination* (worker exit observed by `sync.WaitGroup`) are distinct; only `Wait()` returning confirms termination.
- Result publication uses pre-reserved slots: each task's outcome index is allocated at submission time under a mutex (`sources/conc/pool/result_pool.go:32-37,96-120`), so concurrent completions write disjoint slots and collection happens only after quiescence, guarded by a defensive `TryLock` panic assertion ("collect should not be called until all goroutines have exited", `sources/conc/pool/result_pool.go:122-126`).
- `Stream` adds ordered incremental publication: one dedicated callbacker goroutine applies callbacks sequentially in submission order (`sources/conc/stream/stream.go:121-138`), with a recover-and-repanic shim that injects an empty callback so a panicking worker cannot starve the callbacker (`sources/conc/stream/stream.go:72-81`).

The most instructive detail for Aren is the "leaky abstraction" fix in `ContextPool.Go` (`sources/conc/pool/context_pool.go:40-47`): when cancelling on error, the causal error is written into the collector *before* `cancel()` fires, so cancellation-induced secondary errors cannot race ahead of the cause in first-error mode. This is a hand-rolled ordering guarantee binding outcome recording to terminal notification.

## Rating

**8 / 10**

Rationale against the dimension's purpose (one component with authority to commit state/timing/event/outcome coherently; work errors, runtime failures, panics, cancellation races kept distinct):

- Earned: single unambiguous commit gate (`Wait()` variants) in every type; panic vs error vs cancellation are three distinct channels with different semantics (rethrown value vs joined errors vs ctx signal); deterministic CAS-based arbitration for competing panics; tests explicitly cover simultaneous panic + cancellation (`sources/conc/pool/context_pool_test.go:207-226`), multi-panic races (`sources/conc/waitgroup_test.go:76-86`), and concurrent `Try` from 100 goroutines (`sources/conc/panics/panics_test.go:126-145`).
- Deducted: "first error" arbitration in `WithFirstError` is registration-order, not causally-owned — the tests themselves document intrinsic raciness and patch it with sleeps (`sources/conc/pool/result_error_pool_test.go:66-90`, `sources/conc/pool/context_pool_test.go:168-190`); cancellation has no internal representation in conc's own types (fully outsourced to `context`); there is zero intermediate observability of run progress outside `Stream`; secondary panics are silently discarded after the first.

## Evidence Collected

| Area | Evidence | File:Line |
|------|----------|-----------|
| Single-owner doctrine | README goal #1: scoped concurrency, `WaitGroup` owns spawned goroutines | `sources/conc/README.md:55-62` |
| Panic handling doctrine | Goal #2: `Wait()` panics if any child panicked, decorated with child stacktrace | `sources/conc/README.md:84-114` |
| Owner primitive | `WaitGroup{wg sync.WaitGroup; pc panics.Catcher}`; `Go` = Add + `defer Done` + `pc.Try(f)` | `sources/conc/waitgroup.go:22-34` |
| Terminal commit gate | `Wait()`: barrier on all children, then `Repanic()` publishes panic as last step | `sources/conc/waitgroup.go:38-43` |
| Recover-mode terminal readout | `WaitAndRecover()` returns `*panics.Recovered` instead of rethrowing | `sources/conc/waitgroup.go:47-52` |
| Panic arbitration (first wins) | `Catcher.recovered atomic.Pointer[Recovered]`; `tryRecover` uses `CompareAndSwap(nil, &rp)` | `sources/conc/panics/panics.go:15-17,26-31` |
| Panic payload fidelity | `NewRecovered` captures `runtime.Callers` frames + `debug.Stack` at recovery site | `sources/conc/panics/panics.go:52-61` |
| Panic→error bridge | `Recovered.AsError()` / `ErrRecovered.Unwrap()` expose panic cause as unwrappable error | `sources/conc/panics/panics.go:83-102` |
| Error channel | `ErrorPool.Go` wraps task: return value routed through `addErr`; mutex-guarded append | `sources/conc/pool/error_pool.go:27-31,94-100` |
| Error publication + reset | `Wait()` reads errs after pool drained, joins or returns first, resets slice for reuse | `sources/conc/pool/error_pool.go:35-48` |
| Worker lifecycle | `Pool.worker` drains `p.tasks` channel until closed; limiter token released via defer | `sources/conc/pool/pool.go:150-162,153` |
| Pool teardown sequence | `Wait()`: `close(p.tasks)` → workers drain → `handle.Wait()` → reset `initOnce` | `sources/conc/pool/pool.go:72-82` |
| Config immutability after start | `panicIfInitialized` panics if `tasks != nil` (post-`Go` and post-`Wait`) | `sources/conc/pool/pool.go:108-115` |
| Ordered results via slot reservation | `nextIndex()` reserves index at submission; `save()` writes disjoint slots under mutex | `sources/conc/pool/result_pool.go:32-37,96-120` |
| Premature-collection guard | `collect` does `mu.TryLock()` and panics if lockable → invariant: called only post-quiescence | `sources/conc/pool/result_pool.go:122-126` |
| Cancel-on-error + panic | defer: recover → `p.cancel()` → re-panic, so panics also trigger cancellation | `sources/conc/pool/context_pool.go:26-36` |
| Outcome-before-notification ordering | "Leaky abstraction warning": `addErr(err)` before `cancel()` so causal error registers first | `sources/conc/pool/context_pool.go:38-47` |
| Cancellation cleanup after drain | `ContextPool.Wait` defers `cancel()` until after `errorPool.Wait()` (all workers exited) | `sources/conc/pool/context_pool.go:52-58` |
| Fail-fast composition | `WithFailFast` = `WithFirstError` + `WithCancelOnError` | `sources/conc/pool/context_pool.go:84-92` |
| Stream ordered delivery | single `callbacker` goroutine consumes queue FIFO, executes callbacks sequentially with `Catcher` | `sources/conc/stream/stream.go:121-138` |
| Panic-starvation prevention | task worker recovers, sends empty callback down its channel, re-panics | `sources/conc/stream/stream.go:72-81` |
| Stream teardown under panic | `Wait` defers `close(s.queue)` + callbacker wait so teardown survives propagated panic | `sources/conc/stream/stream.go:91-103` |
| Test: multi-panic race | "one is caught by waitandrecover": two children panic, exactly one `Recovered` surfaces | `sources/conc/waitgroup_test.go:127-138` |
| Test: panic isolation | "nonpanics run successfully": siblings still complete when one child panics | `sources/conc/waitgroup_test.go:140-156` |
| Test: concurrent Try safety | 100 goroutines race `Try`, i==50 panics, value observed once | `sources/conc/panics/panics_test.go:126-145` |
| Test: panic + cancel simultaneity | "WithCancelOnError and panic": panicking task cancels 2 waiting tasks; both cancelled counted AND panic propagates | `sources/conc/pool/context_pool_test.go:207-226` |
| Test: same for result pool | panic triggers cancellation of siblings while panic still surfaces from Wait | `sources/conc/pool/result_context_pool_test.go:121-142` |
| Test: admitted first-error race | comment: intrinsic race patched with sleep; "a race in the test, not in the library" | `sources/conc/pool/result_error_pool_test.go:66-90` |
| Test: cancellation vs normal completion | "no WithCancelOnError": sibling completes normally despite other's error | `sources/conc/pool/context_pool_test.go:149-166` |
| Test: parent-cancel propagation | task returns `ctx.Err()` after external cancel; Wait reports `context.Canceled` | `sources/conc/pool/context_pool_test.go:92-105` |
| Test: panic doesn't exhaust pool | repeated panics with MaxGoroutines(2) still terminate cleanly | `sources/conc/pool/pool_test.go:106-115` |
| Test: reuse after terminal | issue #128 regression: second cycle must not see first cycle's errors/results | `sources/conc/pool/error_pool_test.go:121-134`, `sources/conc/pool/result_pool_test.go:117-128` |
| Test: config locked even after Wait | "after wait" subtests assert reconfiguration still panics post-terminal | `sources/conc/pool/pool_test.go:57-63` |

## Answers to Dimension Questions

### Q1: Can two goroutines publish conflicting terminal outcomes, and what prevents it?

Yes they can publish concurrently; conflicts are either deterministically arbitrated or structurally eliminated, per channel:

- **Panics**: arbitrated by `atomic.CompareAndSwap(nil, &rp)` in `tryRecover` (`sources/conc/panics/panics.go:29`) — strictly first-recovered wins, later panics are dropped. Deterministic winner, silent loser. Verified by `waitgroup_test.go:127-138` (two panics, one survivor) and `waitgroup_test.go:88-98` (non-panicking siblings never overwrite a stored panic).
- **Errors**: no conflict because the collector is an append-only log under a mutex (`sources/conc/pool/error_pool.go:94-100`) — every error survives; there is nothing to arbitrate unless `WithFirstError` is set, in which case the "winner" is simply `errs[0]`, i.e., first to register, which the test suite admits is scheduling-dependent (`sources/conc/pool/result_error_pool_test.go:70-78`).
- **Results**: conflict-free by construction — indices reserved at submission (`nextIndex`, `sources/conc/pool/result_pool.go:33,96-103`) give each task a private slot (`save`, `sources/conc/pool/result_pool.go:105-120`).

### Q2: Is a cancellation request represented separately from confirmed termination?

Yes, but the representation of the request lives in stdlib `context`, not in conc's own types. Request = the shared derived context passed to tasks (`f(p.ctx)`, `sources/conc/pool/context_pool.go:38`) being canceled by the parent or by `cancelOnError` logic (`sources/conc/pool/context_pool.go:45`); confirmation = actual worker exit, observed by `sync.WaitGroup` inside `handle.Wait()` (`sources/conc/pool/pool.go:81`, `sources/conc/waitgroup.go:39`). `ContextPool.Wait` deliberately calls `cancel()` only after full drain, as memory-leak cleanup rather than signaling (`sources/conc/pool/context_pool.go:54-58`). There is no internal "cancel requested but not yet terminated" flag; a task that ends due to cancellation is indistinguishable at the collector level unless it chooses to return `ctx.Err()` (the pattern tests rely on, `sources/conc/pool/context_pool_test.go:95-105`).

### Q3: Can observers see a terminal state before its outcome or terminal event is available?

No — within conc's model, state and outcome are committed together at the `Wait()` boundary. Partial outcomes (errors, results) accumulate in private collectors with no public accessors during the run (`sources/conc/pool/error_pool.go:21-22,38`; `sources/conc/pool/result_pool.go:87-92`), and `resultAggregator.collect` actively enforces this with a `TryLock`-based panic if collection is attempted before quiescence (`sources/conc/pool/result_pool.go:123-126`). The single exception is `stream.Stream`: callbacks execute incrementally, in strict submission order, as tasks complete (`sources/conc/stream/stream.go:121-138`) — but those observers are user-supplied callbacks, not library-published events, and panics are still withheld until `Wait` (`sources/conc/stream/stream.go:89-103`).

### Q4: Which parts can Aren adopt without importing a framework-sized lifecycle model?

Adoptable small pieces (each independently liftable, none larger than ~30 LOC):

1. Owner + catcher pairing: `WaitGroup{wg, pc}` with `Go` = Add/Try/Done and `Wait` = barrier + Repanic (`sources/conc/waitgroup.go:22-43`).
2. First-wins panic arbitration via atomic pointer CAS (`sources/conc/panics/panics.go:15-31`).
3. `Recovered` payload with caller frames + stack captured at recovery (`sources/conc/panics/panics.go:52-74`) and the `AsError` bridge (`sources/conc/panics/panics.go:83-102`).
4. Pre-reserved outcome slots for ordered, conflict-free result aggregation (`sources/conc/pool/result_pool.go:96-120`).
5. Teardown discipline: close-work-channel → drain → barrier → publish, with deferred cleanup that survives panic propagation (`sources/conc/pool/pool.go:72-82`, `sources/conc/stream/stream.go:94-99`).
6. The outcome-registration-before-cancel ordering rule from `ContextPool.Go` (`sources/conc/pool/context_pool.go:40-47`).

Nothing here would over-generalise Aren's contract: there is no framework-sized lifecycle machinery to import. The largest surface is the six-type pool taxonomy (`Pool`/`ErrorPool`/`ContextPool` × optional results, `sources/conc/pool/*.go`), which is API convenience layered by embedding, not runtime complexity — Aren can take the mechanisms and skip the taxonomy. What conc lacks relative to Aren's needs: named states, timestamps/durations, event emission, persistence hooks, per-task identity beyond slot index, and hierarchical scoping.

## Architectural Decisions

1. **Commit-at-Wait**: all publication (state, outcome, panic) is deferred to a single blocking call per scope. `Wait()` is simultaneously the timing barrier, the outcome reader, and the notification dispatcher (`sources/conc/waitgroup.go:36-43`, `sources/conc/pool/error_pool.go:35-48`). Coherence is bought with total loss of intra-run visibility.
2. **Structural ownership over bookkeeping**: instead of tracking states, conc guarantees exit via `sync.WaitGroup` and forces the owner to collect (`sources/conc/README.md:59-62`). A run "exists" only between first `Go` and `Wait` return.
3. **Three distinct failure channels**: returned errors (values, collected), panics (control flow, rethrown at the boundary), cancellation (ctx signal). They interact (panic triggers cancel, `sources/conc/pool/context_pool.go:30-35`; cancel produces ctx errors recorded as ordinary task errors) but are never conflated in storage.
4. **First-wins arbitration for irreconcilable outcomes** (panics) vs **append-only accumulation for reconcilable ones** (errors, results) — chosen per channel based on whether multiple survivors make sense (`sources/conc/panics/panics.go:29`, `sources/conc/pool/error_pool.go:97`).
5. **Slot pre-allocation for deterministic ordering** of concurrent results (`sources/conc/pool/result_pool.go:33`), giving submission-order output regardless of completion order (verified under random sleeps, `sources/conc/pool/result_pool_test.go:66-84`).
6. **Recyclable scopes with frozen configuration**: pools reset their accumulators and `initOnce` after `Wait` for reuse (`sources/conc/pool/pool.go:77-79`, `sources/conc/pool/error_pool.go:39`, `sources/conc/pool/result_pool.go:44`) yet refuse reconfiguration forever (`sources/conc/pool/pool.go:111-115`, tested post-Wait at `sources/conc/pool/pool_test.go:57-63`) — terminal for identity, not for liveness.

## Notable Patterns

- **CAS-monotonic latch** (`panics.Catcher`): lock-free "first writer wins" using `CompareAndSwap(nil, ...)` (`sources/conc/panics/panics.go:29`) — the smallest possible terminal arbiter.
- **Recover–act–rethrow sandwich**: used twice — `ContextPool` cancels then re-panics (`sources/conc/pool/context_pool.go:30-35`), stream workers emit a placeholder callback then re-panic (`sources/conc/stream/stream.go:73-81`) — letting side effects ride along a panic path without swallowing it.
- **Deadlock-proof teardown**: stream `Wait` defers `close(queue)` before the blocking pool wait so a propagated panic cannot strand the callbacker (`sources/conc/stream/stream.go:94-99`); the empty-callback injection prevents a panicked producer from starving the consumer (`sources/conc/stream/stream.go:74-78`).
- **Invariant assertions at collection time**: `TryLock` panic in `collect` converts a would-be data race into an immediate programming-error crash (`sources/conc/pool/result_pool.go:124-126`).
- **Documented leaky abstraction**: the author explicitly flags breaking layering (`addErr` called directly instead of returning through the wrapper) to get error-registration/cancel ordering right, with rationale in-place (`sources/conc/pool/context_pool.go:40-44`) — a good precedent for annotating ordering-critical shortcuts.
- **Self-documented racy tests**: sleeps plus candid comments distinguishing test races from library races (`sources/conc/pool/result_error_pool_test.go:72-78`).

## Tradeoffs

- **Coherence vs observability**: commit-at-Wait makes state/outcome/notification atomic but offers no progress reporting, no partial completion status, no timeouts anywhere in the library; a hung task hangs `Wait` indefinitely (no watchdog exists in `sources/conc/pool/pool.go:150-162`).
- **Determinism vs information loss**: first-wins panic arbitration discards secondary panics including their stacks (`sources/conc/panics/panics.go:29`); `WithFirstError` similarly hides all but one error, and which one survives is scheduler-dependent (`sources/conc/pool/context_pool_test.go:176-186`).
- **Simplicity vs cancellation modeling**: outsourcing cancellation to `context` keeps conc tiny but means the library itself cannot answer "was this termination requested?" — callers infer from error values.
- **Reuse vs terminality**: recyclable pools (`sources/conc/pool/pool.go:77-79`) mean a completed scope isn't final, complicating any audit log of transitions; issue #128 shows sticky accumulators were a real bug (`sources/conc/pool/error_pool_test.go:121-134`).
- **Panic-as-control-flow across boundaries**: `Repanic` throws across `Wait` (`sources/conc/waitgroup.go:42`), so owners must handle two failure modalities (error values and panics) in one code path — expressive, but easy to misuse.

## Failure Modes / Edge Cases

- **Submit racing teardown**: `close(p.tasks)` in `Wait` (`sources/conc/pool/pool.go:75`) versus sends in `Go` (`sources/conc/pool/pool.go:45,62`) has no synchronization; concurrent `Go` during `Wait` would panic on send-to-closed-channel. Safe only under the single-scope-owner assumption of structured concurrency (`sources/conc/README.md:55-62`).
- **Panicking tasks vanish from collectors**: a task that panics never records an error or result; the entire batch terminates as a panic at `Wait`, and `ResultPool` returns only non-panicked results (`sources/conc/pool/result_pool.go:39-46`). Mixed error+panic batches surface only via panic.
- **Hung task = hung scope**: limiter tokens release only on worker exit (`sources/conc/pool/pool.go:153`); no timeout mechanism exists, so one blocked task blocks teardown and every queued successor.
- **Placeholder substitution in streams**: a panicking task contributes a synthesized empty callback (`sources/conc/stream/stream.go:78`), silently altering the callback sequence downstream consumers observe.
- **Secondary panic stacks lost** (see Tradeoffs): only the first recovered panic is retained.
- **Unsynchronized terminal readouts rely on happens-after**: `ErrorPool.Wait` reads `p.errs` without locking (`sources/conc/pool/error_pool.go:38`), correct solely because worker exit precedes it; any future API exposing mid-run reads would break this invariant — precisely what the `collect` `TryLock` guard protects in the result path (`sources/conc/pool/result_pool.go:124-126`).

## Future Considerations

For Aren Phase 1, the transferable findings are:

1. Make Aren's commit point as explicit as conc's `Wait` — one gate where state, outcome, timing, and notification flip together — rather than allowing providers/tools to publish piecemeal.
2. Keep conc's per-channel arbitration split: CAS-latch for irreconcilable terminal values, append-only log for reconcilable ones, private slots for ordered payloads.
3. Improve on conc where Aren's contract demands more: represent cancellation-request internally alongside confirmed-termination instead of delegating wholly to `context`, and consider retaining all panic records (not just the first) since Aren runs outlive a single process lifetime and diagnostics matter.
4. Adopt the ordering rule "record causal outcome before triggering cascading cancellation" (`sources/conc/pool/context_pool.go:40-47`) as a hard invariant in Aren's transition code.
5. If Aren needs observable progress, `stream.Stream`'s single sequencer goroutine (`sources/conc/stream/stream.go:121-138`) is the closest small analog to ordered event publication; anything richer (subscribers, replay) exceeds conc and should be designed fresh.

## Questions / Gaps

- The dimension cites Aren requirement docs (`../../../../Aren/docs/phase-1-prd/02-lifecycle-requirements.md`, `01-product-definition-and-scope.md`). Neither path resolves on disk from the study root (checked `/home/antonioborgerees/coding/Aren/docs/phase-1-prd/` — directory not found). Per source-isolation rules I did not search elsewhere; comparisons to Aren's contract rest on the dimension text alone. No evidence found in-source regarding exact Aren transition vocabulary.
- No evidence found in-source for duration/timestamp capture at transitions: conc records no timing metadata whatsoever (only perf notes in doc comments, e.g., `sources/conc/pool/pool.go:27-29`). Whether Aren needs transition timestamps must be decided from its PRD, not from this study.
- Behavior under `Wait` called concurrently from multiple goroutines, or `Go` after `Wait` without reuse intent, is undocumented and untested in-source; inferred unsafe from the unsynchronized `close`/send pair noted above.
- The `iter` package (`sources/conc/iter/iter.go:59-85`) adds a work-stealing-style atomic index over `WaitGroup` but introduces no new terminal-arbitration mechanics; covered here only to confirm no additional lifecycle machinery exists in the module.

---

Generated by `dimensions/01.01-lifecycle-transition-ownership-and-terminal-arbitration.md` against `conc`.
