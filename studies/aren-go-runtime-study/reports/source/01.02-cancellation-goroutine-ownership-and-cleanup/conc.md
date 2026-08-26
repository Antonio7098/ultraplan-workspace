# Source Analysis: conc

## 01.02 Cancellation, Goroutine Ownership, and Cleanup

### Source Info

| Field | Value |
|-------|-------|
| Name | conc (`github.com/sourcegraph/conc`) |
| Path | `studies/aren-go-runtime-study/sources/conc` |
| Language / Stack | Go 1.20; stdlib `sync`, `context`, `sync/atomic`; testify for tests only (`sources/conc/go.mod:3-5`) |
| Analyzed | 2026-08-26 |

## Summary

`conc` is a structured-concurrency toolkit built around a single primitive: a panic-catching `WaitGroup` (`sources/conc/waitgroup.go:22-25`) that owns every goroutine it spawns and republishes child panics at the join point. On top of it sit worker-pool variants (`Pool`, `ErrorPool`, `ContextPool`, and their `Result*` counterparts) and an ordered-callback `Stream`. The library's stated goals are "make it harder to leak goroutines" and "handle panics gracefully" (`sources/conc/README.md:45-46`), with the explicit opinion that "all concurrency should be scoped" and that the owner of a goroutine is always a `WaitGroup` that is `Wait()`-ed before going out of scope (`sources/conc/README.md:55-61`).

For this dimension the picture splits sharply:

1. **Goroutine ownership and cleanup ordering are strong.** There is exactly one raw `go func()` statement in the entire library (`sources/conc/waitgroup.go:30`); every other spawn goes through a tracked handle. Shutdown ordering is deliberate — e.g. `Stream.Wait` defers callbacker cleanup so a panicking task join cannot strand the callbacker (`sources/conc/stream/stream.go:94-99`).
2. **Cancellation is thin and purely cooperative.** Contexts exist only in `ContextPool`/`ResultContextPool` via one `context.WithCancel` per pool (`sources/conc/pool/pool.go:142`). The base `Pool`, `WaitGroup`, `Stream`, and `iter` have no cancellation surface at all: `Go()` blocks unboundedly when the pool is saturated, no submit path accepts a context, and there is no timeout escalation or forced teardown anywhere. Stopping work relies entirely on tasks selecting on `ctx.Done()`.
3. **Idempotency under concurrent shutdown callers is contractual, not mechanical.** `Pool.Wait` closes the shared task channel on every call (`sources/conc/pool/pool.go:75`), so two concurrent `Wait` callers double-close and panic. The single-owner discipline in the README is what keeps this safe in practice.

Cleanup failure cannot overwrite the primary failure in sequential use: first-panic-wins CAS storage (`sources/conc/panics/panics.go:29`), error-before-cancel publication ordering (`sources/conc/pool/context_pool.go:39-47`), and deferred joins after panic propagation all protect the original signal.

## Rating

**6 / 10**

Rationale against the dimension's concern areas:

- Goroutine ownership tree and join guarantees: excellent (9). Every spawn is tracked; `Wait` truthfully means "all spawned goroutines exited" (`sources/conc/waitgroup.go:36-43`, `sources/conc/pool/pool.go:70-82`).
- Cleanup ordering and truthful completion of failures: strong (8). Deferred cleanup survives panics (`sources/conc/stream/stream.go:94-99`); primary errors are published before cancellation fires (`sources/conc/pool/context_pool.go:39-47`).
- Cancellation propagation: weak (4). Only `ContextPool` propagates a derived context; nothing cancels blocked submitters, and `iter`/`MapErr` run to completion regardless of failures (`sources/conc/iter/map.go:53-62`).
- Idempotent/repeatable shutdown and leak testing: weak (3–4). Concurrent `Wait` double-closes a channel (`sources/conc/pool/pool.go:75`); no goleak dependency or leak test exists anywhere in the module (`sources/conc/go.mod:5`; grep for `goleak|SetFinalizer|AddCleanup` across `*.go`, `*.yml`, `*.mod` returned zero matches).

The composite lands at 6: best-in-class for scoped ownership, materially incomplete for Phase-2-style pressure (long-lived work, partial startup failure, concurrent supervision).

## Evidence Collected

Every entry includes a file path with line numbers.

| Area | Evidence | File:Line |
|------|----------|-----------|
| Ownership root | `WaitGroup.Go` does `wg.Add(1)` before spawning and `defer wg.Done()` inside the child; panics routed to shared `Catcher` via `Try` | `sources/conc/waitgroup.go:28-34` |
| Panic republish at join | `Wait()` blocks then calls `pc.Repanic()`; `WaitAndRecover()` returns `*panics.Recovered` instead | `sources/conc/waitgroup.go:38-52` |
| First-failure-wins storage | `Catcher.tryRecover` uses `recovered.CompareAndSwap(nil, &rp)` so later panics never overwrite the first | `sources/conc/panics/panics.go:26-31` |
| Single raw `go` statement | Only raw goroutine spawn in the library; all other spawns route through handles | `sources/conc/waitgroup.go:30` |
| Pool worker lifecycle | Workers range over `p.tasks` until channel closed+drained; limiter token released via `defer` | `sources/conc/pool/pool.go:150-162` |
| Pool shutdown = close + drain + join | `Wait()` closes `p.tasks`, defers `initOnce` reset, then `handle.Wait()` | `sources/conc/pool/pool.go:72-82` |
| Pool reuse after Wait | `initOnce` reset enables re-initialization on next use | `sources/conc/pool/pool.go:77-79`; test `sources/conc/pool/pool_test.go:128-146` |
| Bounded submission | `Go()` blocks when limited; limiter is a plain `chan struct{}` used as semaphore | `sources/conc/pool/pool.go:37-66`, `sources/conc/pool/pool.go:164-173` |
| Config freeze after first use | `panicIfInitialized` panics if `p.tasks != nil` — configuration locked once pool starts | `sources/conc/pool/pool.go:108-115`; tests `sources/conc/pool/pool_test.go:49-64` |
| Context derivation and cancel guarantee | `WithContext` creates `context.WithCancel(ctx)` and stores both ctx and cancel; `ContextPool.Wait` defers `cancel()` "to avoid memory leakage" | `sources/conc/pool/pool.go:140-148`; `sources/conc/pool/context_pool.go:15-16`, `sources/conc/pool/context_pool.go:54-58` |
| Cancel-on-error ordering | Error added directly to ErrorPool *before* `cancel()` so sibling `context.Canceled` errors can't win `WithFirstError` (explicit "leaky abstraction warning" comment) | `sources/conc/pool/context_pool.go:39-47` |
| Panics also trigger cancellation | In `cancelOnError` mode, deferred recover → `p.cancel()` → re-panic | `sources/conc/pool/context_pool.go:26-36`; test `sources/conc/pool/context_pool_test.go:207-226` |
| Fail-fast composition | `WithFailFast` = `WithFirstError` + `WithCancelOnError` | `sources/conc/pool/context_pool.go:87-92`; test `sources/conc/pool/context_pool_test.go:192-205` |
| Parent-context cancellation flows through | Tasks receive derived ctx; canceled/deadline parent surfaces as returned error | `sources/conc/pool/context_pool_test.go:92-131` |
| Stream shutdown ordering | `Wait` defers `close(s.queue)` + `callbackerHandle.Wait()` so callbacker is joined even when `pool.Wait` panics; workers joined first | `sources/conc/stream/stream.go:89-103` |
| Callback starvation guard on panic | Task recovers, deposits empty callback into its channel, re-panics so callbacker never blocks | `sources/conc/stream/stream.go:73-81` |
| Bounded queue backpressure | Queue capacity `MaxGoroutines()+1`; per-task callback channels buffered size 1, recycled via `sync.Pool` | `sources/conc/stream/stream.go:110-117`, `sources/conc/stream/stream.go:140-154` |
| Callback panic isolation | Callbacker catches each callback panic with `Catcher.Try` and continues; repanics after loop ends — other callbacks still run | `sources/conc/stream/stream.go:121-138`; tests `sources/conc/stream/stream_test.go:101-135` |
| Result slot reservation | `nextIndex()` reserves slots under mutex before execution; results land positionally even out of order | `sources/conc/pool/result_pool.go:32-37`, `sources/conc/pool/result_pool.go:96-120` |
| Defensive concurrent-Wait detector | `collect` uses `mu.TryLock()` and panics "collect should not be called until all goroutines have exited" | `sources/conc/pool/result_pool.go:123-126` |
| iter bounded fan-out | At most `min(GOMAXPROCS, len(input))` goroutines pulling indices from an atomic counter; inline `wg.Wait()` | `sources/conc/iter/iter.go:59-84` |
| MapErr has no early exit | Errors accumulate under mutex; iteration always completes; combined via `errors.Join` | `sources/conc/iter/map.go:47-63` |
| Error collection reset on Wait | `errs` drained and nil-ed for reuse; `onlyFirstError` returns `errs[0]` | `sources/conc/pool/error_pool.go:35-48`; test `sources/conc/pool/error_pool_test.go:121-131` |
| Stated design goal (scoped concurrency) | README: owner of every goroutine is a `WaitGroup`; call `Wait` before it goes out of scope; goal list "harder to leak goroutines", "handle panics gracefully" | `sources/conc/README.md:45-61`, `sources/conc/README.md:105-110` |
| Test coverage of cancel semantics | Parent cancel/timeout propagation, cancel-on-error, cancel-on-panic, fail-fast suppression of secondary `context.Canceled` | `sources/conc/pool/context_pool_test.go:92-226` |
| Missing leak testing | No `goleak` in deps, no finalizers/cleanup hooks, no goroutine-count assertions in any test file | `sources/conc/go.mod:5`; zero matches for `goleak\|SetFinalizer\|AddCleanup` in `*.go/*.yml/*.mod` |

## Answers to Dimension Questions

**Who is responsible for waiting for every goroutine started on behalf of a run?**
A single owning handle per concurrency scope. `conc.WaitGroup` is the root: its `Go` registers the child before spawn and `Wait` joins all of them (`sources/conc/waitgroup.go:28-43`). `Pool` embeds a `WaitGroup` as `handle` and joins all workers in `Wait` (`sources/conc/pool/pool.go:31`, `sources/conc/pool/pool.go:81`). `Stream` owns two handles: the inner pool's workers plus a dedicated `callbackerHandle` for its single callbacker goroutine, both joined in `Wait` (`sources/conc/stream/stream.go:40-44`, `sources/conc/stream/stream.go:91-103`). `iter.ForEachIdx` spawns and joins within the call itself (`sources/conc/iter/iter.go:80-84`). Responsibility is contractual: the API requires the user to call `Wait` ("you must call Wait() to clean up any spawned goroutines", `sources/conc/pool/pool.go:17-19`; README rule, `sources/conc/README.md:59-61`). Nothing enforces it mechanically — a forgotten `Wait` leaks workers silently.

**What does the API claim when cancellation is requested but work has not stopped?**
Nothing explicit, because cancellation is cooperative by construction. `Wait` claims only full completion: it "will block until all goroutines ... exit" (`sources/conc/waitgroup.go:36-37`), i.e. terminal publication happens strictly after the deepest operation returns — there is no "cancelled-but-running" state in the API. For `ContextPool`, cancellation only signals; tasks that ignore `ctx.Done()` keep `Wait` blocked indefinitely. One honest admission exists in docs: with `WithCancelOnError`, sibling tasks race cancellation and return `context.Canceled`, so "all errors after the first will likely be context.Canceled" — which is why `WithFirstError` exists (`sources/conc/pool/context_pool.go:60-77`). Tests pin the observable behavior: canceled siblings do publish `ctx.Err()` alongside the trigger error unless suppressed (`sources/conc/pool/context_pool_test.go:134-147` vs `192-205`).

**Can cleanup block forever or overwrite the primary failure?**
Block forever: yes, by design. Worker exit requires the closed-and-drained task channel plus current-task completion (`sources/conc/pool/pool.go:159-161`); a task that never observes its context holds `Wait` open with no timeout escalation anywhere in the module. Overwrite: protected in sequential use through three mechanisms — (1) first-panic-wins CAS in `Catcher` (`sources/conc/panics/panics.go:29`); (2) publication-before-cancellation ordering, where the triggering error is recorded before `cancel()` so late `context.Canceled` errors can't displace it under `WithFirstError` (`sources/conc/pool/context_pool.go:39-47`); (3) cleanup-after-propagation ordering in `Stream.Wait`, where the defer runs even though `pool.Wait()` may panic upward (`sources/conc/stream/stream.go:94-99`). The residual overwrite hazard is structural: a concurrent second `Wait` caller panics on double-close (`sources/conc/pool/pool.go:75`), replacing normal completion with a runtime panic.

**How is shutdown made idempotent under concurrent callers?**
It largely isn't, beyond what the primitives give for free. Idempotent pieces: repeated `context.CancelFunc` calls (`context_pool.go:56` fires it once per `Wait`, tasks fire it again on error/panic — safe), repeated `WaitGroup.Wait` (plain `wg.Wait`, `sources/conc/waitgroup.go:39`). Not idempotent: `Pool.Wait` executes `close(p.tasks)` unconditionally (`sources/conc/pool/pool.go:75`), so two simultaneous callers double-close → panic; the `initOnce = sync.Once{}` reset (`sources/conc/pool/pool.go:79`) also races under concurrent callers. The design answer is ownership discipline instead of locking: exactly one supervisor goroutine drives the submit→Wait lifecycle (README's scoped stance, `sources/conc/README.md:55-61`). The nearest mechanical guard is `resultAggregator.collect`'s `TryLock` panic, which detects (but doesn't prevent) joining while workers still run (`sources/conc/pool/result_pool.go:123-126`).

## Architectural Decisions

1. **One canonical goroutine-owning type.** All pools, streams, and iterators compose `conc.WaitGroup` rather than reimplementing tracking (`sources/conc/pool/pool.go:31`, `sources/conc/stream/stream.go:40`, `sources/conc/iter/iter.go:80`). This makes the ownership tree trivially auditable: one raw `go` statement in the whole codebase (`sources/conc/waitgroup.go:30`).
2. **Panics are a first-class failure channel, unified with errors at the join point.** Children catch into a shared `Catcher` (`sources/conc/waitgroup.go:32`); `Wait` re-raises with stack metadata (`sources/conc/panics/panics.go:52-61`). This preserves "crash the process" semantics of unhandled panics while converting them into join-point failures — the README explicitly argues propagation over swallowing/logging (`sources/conc/README.md:84-110`).
3. **Cancellation is opt-in via a separate pool type, not a core capability.** Rather than threading contexts through everything, `WithContext` converts a configured pool into a `ContextPool` holding one derived context (`sources/conc/pool/pool.go:140-148`). Base pools stay allocation-light and context-free.
4. **Shutdown modeled as close-and-drain, not signal-and-abandon.** Closing `tasks` lets in-flight submissions drain and workers exit naturally (`sources/conc/pool/pool.go:75`, `sources/conc/pool/pool.go:159-161`), giving deterministic completion and enabling struct reuse via `initOnce` reset (`sources/conc/pool/pool.go:77-79`).
5. **Fail-fast implemented as application-level ordering, not scheduler magic.** Cancel-on-error records the causal error first, then cancels, then returns nil to suppress duplicate registration (`sources/conc/pool/context_pool.go:39-47`) — an explicit workaround for the race between cancellation side effects and error publication.
6. **Configuration immutability enforced by panic.** All `With*` methods call `panicIfInitialized` (`sources/conc/pool/pool.go:111-115`), eliminating a class of mid-flight reconfiguration races at the cost of eager validation.

## Notable Patterns

- **Recover–act–rethrow wrappers** to make deferred side effects part of panic propagation: `ContextPool.Go` cancels then re-panics (`sources/conc/pool/context_pool.go:30-35`); `Stream.Go` deposits an empty callback then re-panics to unblock the consumer (`sources/conc/stream/stream.go:73-81`).
- **Semaphore-as-channel limiter** with release-on-worker-exit via defer (`sources/conc/pool/pool.go:153`, `sources/conc/pool/pool.go:170-173`), keeping the concurrency bound correct even when tasks panic (worker dies, token freed, replacement spawned — tested in `sources/conc/pool/pool_test.go:106-115`).
- **Bounded buffering everywhere**: queue capacity tied to pool width (`MaxGoroutines()+1`, `sources/conc/stream/stream.go:112`); single-slot callback channels recycled through `sync.Pool` (`sources/conc/stream/stream.go:142-154`). Backpressure surfaces as blocking `Go`, never as unbounded memory.
- **Positional result slots** decoupling completion order from submission order (`nextIndex`/`save`, `sources/conc/pool/result_pool.go:96-120`) — cleanup of results never depends on scheduling order.
- **Single-consumer ordered drain**: one callbacker goroutine consumes a queue of per-task channels sequentially, isolating callback panics so producers never starve (`sources/conc/stream/stream.go:121-138`; test `sources/conc/stream/stream_test.go:121-135`).
- **Atomic-counter work stealing in iterators** avoids per-element channel/goroutine cost: N goroutines pull indices until exhausted (`sources/conc/iter/iter.go:71-78`).

## Tradeoffs

- **Contractual lifetime management vs. enforcement.** Scoped concurrency keeps hot paths fast (~300ns/task claimed, `sources/conc/pool/pool.go:27-29`) but shifts leak risk onto API discipline; a missed `Wait` leaks silently because nothing tracks outstanding handles.
- **Cooperative-only cancellation vs. responsiveness.** No timeouts, no forced teardown, no ctx-aware submit. Simple and composable, but long-lived or misbehaving tasks stall `Wait` indefinitely — the exact Phase-2 pressure scenario.
- **Panic-as-control-flow at boundaries.** Repanic-at-join is semantically clean but makes `Wait` a panic site, forcing callers into `defer recover` or `WaitAndRecover` patterns; combined with config-freeze panics, several misuse modes surface as panics rather than errors.
- **Reuse supported inconsistently.** `Pool` resets `initOnce` (`sources/conc/pool/pool.go:79`), `ErrorPool` resets `errs` (`sources/conc/pool/error_pool.go:38-39`), `ResultPool` resets `agg` (`sources/conc/pool/result_pool.go:44`) — but `Stream` never resets `initOnce`, so post-`Wait` reuse sends on a closed queue (`sources/conc/stream/stream.go:69` vs unclosed-over state at `sources/conc/stream/stream.go:97`). A `ContextPool` reused after `Wait` runs subsequent rounds on an already-canceled context (cancel fired at `sources/conc/pool/context_pool.go:56`).
- **Blocking `Go` couples submitters to pool throughput.** With a limiter, submission blocks without cancellation escape (`sources/conc/pool/pool.go:54-66`) — a producer cannot abandon a saturated pool.

## Failure Modes / Edge Cases

- **Concurrent `Wait` callers → double close panic.** Both pass `init()` (shared Once), both execute `close(p.tasks)` (`sources/conc/pool/pool.go:72-82`). Also makes the documented reuse reset itself racy. No test exercises concurrent `Wait`.
- **`Go` racing `Wait` → send on closed channel.** Contract says all submissions precede `Wait` (`sources/conc/pool/pool.go:17-19`), but nothing detects violation; a straggler `Go` after `Wait` starts panics at `p.tasks <- f` (`sources/conc/pool/pool.go:45`, `62`).
- **Forgotten `Wait` on `ContextPool` leaks the derived context.** `cancel` runs only in `Wait`'s defer (`sources/conc/pool/context_pool.go:54-58`); skipping `Wait` defeats Go's context-vet rule (cancel must be called) with no library-side guard or finalizer.
- **Unresponsive task stalls shutdown forever.** Workers exit only after the current task returns (`sources/conc/pool/pool.go:159-161`); the callbacker similarly blocks forever on `<-callbackCh` if a task neither returns nor panics (`sources/conc/stream/stream.go:126-128`). No watchdog, deadline, or escalation exists.
- **`panicIfInitialized` check is racy.** `p.tasks != nil` read without synchronization (`sources/conc/pool/pool.go:111-115`) can race `init()` from a concurrent first `Go`.
- **`MapErr` cannot short-circuit.** Every element is processed even after failures; cancellation-aware iteration requires dropping to `pool.ContextPool` manually (`sources/conc/iter/map.go:53-63`).
- **Secondary-error suppression depends on timing.** Without `WithFirstError`, a cancel-on-error run legitimately reports both the cause and sibling `context.Canceled`s (`sources/conc/pool/context_pool_test.go:134-147`); the `WithFirstError` test itself documents a sleep-based race workaround (`sources/conc/pool/context_pool_test.go:176-186`).
- **Partial-startup failure is untested.** No test spawns tasks where some block forever while others fail, then asserts shutdown; the closest is panic-propagation under limit (`sources/conc/pool/pool_test.go:106-115`).

## Future Considerations

Concrete follow-ups this analysis supports, sized for Aren's Phase 2 executor work:

1. **Cancellation-aware submission.** Add a `TryGo`/`GoCtx` variant that selects on `ctx.Done()` alongside the limiter/tasks send (`sources/conc/pool/pool.go:54-66`) so producers can abandon saturated pools — currently impossible without leaking the submit attempt.
2. **Mechanically idempotent shutdown.** Wrap the `close(p.tasks)` + join sequence in a `sync.Once` (or document + lint the single-owner rule) to remove the concurrent-`Wait` double-close hazard at `sources/conc/pool/pool.go:75`.
3. **Leak detection in CI.** Adopt `go.uber.org/goleak` as a test dependency and add a partial-startup suite (task blocked past deadline + failing sibling + panic) asserting goroutine counts return to baseline; today the module has zero leak instrumentation (`sources/conc/go.mod:5`).
4. **Normalize reuse semantics.** Either reset `Stream.initOnce` symmetrically with `Pool` (`sources/conc/stream/stream.go:111-116` vs `sources/conc/pool/pool.go:79`) or document `Stream`/`ContextPool` as single-shot; current behavior differs silently across types.
5. **Deadline escalation helper.** Since the library refuses to force-stop, provide a wrapper that layers `context.WithTimeout` over a run and converts expiry into a published error, making "cancellation requested but work not stopped" an explicit, observable outcome rather than an indefinite hang.
6. **Guard the init race.** Replace the unsynchronized `p.tasks != nil` check (`sources/conc/pool/pool.go:112`) with an atomic flag or fold it into the same `initOnce`.

## Questions / Gaps

- **How does this map to Aren's lifecycle requirements?** The dimension file cites `../../../../Aren/docs/phase-1-prd/02-lifecycle-requirements.md`, but Source Isolation Rules forbid reading outside `sources/conc`. Mapping conc's contract (submit→Wait, cooperative cancel) onto Aren's specific requirements was therefore not performed here. Searched: only the injected prompt text.
- **No evidence found** for any timeout-escalation, forceful teardown, or process/closer management: grep for `goleak|SetFinalizer|AddCleanup|WithTimeout|AfterFunc` across `sources/conc/**` yields only the test-local `context.WithTimeout` at `sources/conc/pool/context_pool_test.go:109`. Processes and external closers are out of scope for this library entirely.
- **Concurrent-supervisor behavior is unverified empirically.** The double-close analysis at `sources/conc/pool/pool.go:75` is derived from code reading; no test demonstrates it. Running such a test was considered but modifying/executing study-source code falls outside this analysis task's workflow.
- **Performance claims** (~1µs startup, ~300ns/task, `sources/conc/pool/pool.go:27-29`; ~500ns/stream task, `sources/conc/stream/stream.go:34-37`) are documentation-only benchmarks (`BenchmarkPool`, `sources/conc/pool/pool_test.go:149-166`; `BenchmarkStream`, `sources/conc/stream/stream_test.go:138-159` exist but assert no numbers). Treated as stated intent, not verified behavior.

---

Generated by `01.02-cancellation-goroutine-ownership-and-cleanup` against `conc`.
