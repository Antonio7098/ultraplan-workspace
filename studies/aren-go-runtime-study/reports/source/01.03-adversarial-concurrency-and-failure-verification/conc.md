# Source Analysis: conc

## 01.03 Adversarial Concurrency and Failure Verification

### Source Info

| Field | Value |
|-------|-------|
| Name | conc (sourcegraph/conc) |
| Path | `studies/aren-go-runtime-study/sources/conc` |
| Language / Stack | Go 1.20 (module `github.com/sourcegraph/conc`), testify for assertions, golangci-lint, GitHub Actions CI |
| Analyzed | 2026-08-26 |

## Summary

conc is an in-process structured-concurrency library (WaitGroup, goroutine pools with errors/context/results, ordered Stream, parallel iterator). Its adversarial verification strategy is deliberately narrow but consistently applied: **(1)** every presubmit run executes the entire suite under the Go race detector (`Makefile:20`, `.github/workflows/go.yml:37`); **(2)** each pool variant carries a near-identical matrix of failure tests — panic propagation, concurrency-limit falsification via atomic counters, config-after-init misuse panics, and reuse-after-Wait state resets pinned to a real bug report (`pool/error_pool_test.go:121-134`, issue #128); **(3)** cancellation is exercised at error, panic, parent-cancel, and timeout boundaries; and **(4)** expensive benchmarks are split into dedicated workflows that gate on performance-regression alerts against the main branch (`Makefile:24`, `.github/workflows/bench.yml:31,63`). The library has no injection seams for clocks or schedulers — tests use real time, real contexts, and unseeded randomness, and three tests openly document intrinsic races mitigated only by sleeps. There are no fuzz targets, no leak detectors (despite "make it harder to leak goroutines" being goal #1 of `studies/aren-go-runtime-study/sources/conc/README.md:45`), and no stress/crash-recovery suites. Verified locally in this study: `go test -race ./...` passes clean across all five packages.

## Rating

**7 / 10**

Rationale against the dimension's purpose ("repeatable techniques for proving invariants under races, panics, partial writes, restarts, resource exhaustion"):

- Strong and directly transferable: race detector as a hard presubmit gate; atomic-counter falsification of the bounded-parallelism invariant at limits {1, 10, 100}; a systematic panic-propagation matrix covering every pool variant including the panic-during-cancellation interaction (`pool/context_pool_test.go:207-226`); regression tests tied to a concrete historical bug (issue #128 reuse leakage of errors/results).
- Weak where Aren Phase 1 also needs proof: no fuzzing, no goroutine-leak detection despite leak-avoidance being a stated goal, admitted-flaky timing-biased first-error tests instead of deterministic sequencing, and no fault-injection seams. These gaps cap it below 8.

## Evidence Collected

Every entry is cited workspace-relative to `studies/aren-go-runtime-study/sources/conc/`.

| Area | Evidence | File:Line |
|------|----------|-----------|
| Race-detector presubmit gate | `go test -race -v ./... -coverprofile ./coverage.txt` runs on every push/PR to main | `studies/aren-go-runtime-study/sources/conc/Makefile:20`; `studies/aren-go-runtime-study/sources/conc/.github/workflows/go.yml:37` |
| Race suite verified passing | Local run in this study: all 5 packages `ok` under `-race` | `studies/aren-go-runtime-study/sources/conc/waitgroup_test.go` … `stream/stream_test.go` (observed behavior) |
| Panic propagation invariant (WaitGroup) | Child panic re-raised by Wait; siblings still complete; first panic wins over later non-panics | `studies/aren-go-runtime-study/sources/conc/waitgroup_test.go:64-115` |
| First-panic-wins mechanism | `Catcher.recovered atomic.Pointer[Recovered]` + `CompareAndSwap(nil, &rp)` makes concurrent panics race-safe | `studies/aren-go-runtime-study/sources/conc/panics/panics.go:16,26-31` |
| Catcher concurrency test | 100 goroutines call Try simultaneously; exactly the designated panic (i==50) is recovered | `studies/aren-go-runtime-study/sources/conc/panics/panics_test.go:126-145` |
| Bounded-parallelism falsifier | Atomic counter asserts concurrency never exceeds MaxGoroutines and returns to 0, at limits 1/10/100 with 10x tasks | `studies/aren-go-runtime-study/sources/conc/pool/pool_test.go:66-90`; same pattern `pool/error_pool_test.go:96-119`, `pool/context_pool_test.go:228-253`, `pool/result_pool_test.go:86-115`, `pool/result_error_pool_test.go:104-132`, `pool/result_context_pool_test.go:198-230` |
| Panics do not exhaust pool capacity | Pool limited to 2 workers survives 10 sequential panicking tasks | `studies/aren-go-runtime-study/sources/conc/pool/pool_test.go:106-115`; limiter release on worker exit `pool/pool.go:150-153,170-174` |
| Cancellation on error boundary | `cancelOnError`: error registered *before* cancel so causal error isn't lost to racing siblings (documented "leaky abstraction" note) | `studies/aren-go-runtime-study/sources/conc/pool/context_pool.go:38-48` |
| Cancellation on panic boundary | Deferred recover→cancel→re-panic wrapper ensures panics also trigger cancellation | `studies/aren-go-runtime-study/sources/conc/pool/context_pool.go:30-36` |
| Cancel-on-panic end-to-end | Asserts both that Wait panics AND both sibling tasks observed ctx.Done() | `studies/aren-go-runtime-study/sources/conc/pool/context_pool_test.go:207-226`; result variant `pool/result_context_pool_test.go:121-142` |
| Parent cancel / timeout boundaries | Tasks return ctx.Err() after external cancel and WithTimeout deadline | `studies/aren-go-runtime-study/sources/conc/pool/context_pool_test.go:92-131` |
| Fail-fast composition | WithFailFast = first-error-only + cancel-on-error; second error asserted absent | `studies/aren-go-runtime-study/sources/conc/pool/context_pool.go:87-92`; `pool/context_pool_test.go:192-205` |
| Ordered results under shuffled completion | Random sleep up to 100ms per task; submission order still preserved via index-reservation aggregator | `studies/aren-go-runtime-study/sources/conc/pool/result_pool_test.go:66-84`; mechanism `pool/result_pool.go:32-37,96-120` |
| Stream starvation guard on task panic | Recovered task panic injects empty callback so single callbacker never blocks forever | `studies/aren-go-runtime-study/sources/conc/stream/stream.go:72-86`; tested `stream/stream_test.go:121-135` |
| Single-callbacker ordering invariant | t.Fatal if >1 callback runs concurrently; ≤5 concurrent tasks enforced | `studies/aren-go-runtime-study/sources/conc/stream/stream_test.go:71-99`; callbacker loop `stream/stream.go:121-138` |
| Reuse/state-reset regression (issue #128) | Errors/results from a previous generation must not appear after Wait+Go+Wait | `studies/aren-go-runtime-study/sources/conc/pool/error_pool_test.go:121-134`; `pool/result_pool_test.go:117-128`; `pool/result_error_pool_test.go:134-148`; `pool/result_context_pool_test.go:232-246` |
| State-reset implementation | errs slice nil'd in Wait, aggregator zeroed, initOnce reset so pool can be reused | `studies/aren-go-runtime-study/sources/conc/pool/error_pool.go:35-39`; `pool/result_pool.go:41-46`; `pool/pool.go:72-82` |
| API-misuse guard (config after init) | All pools panic if reconfigured after first Go(); tested before and after Wait | `studies/aren-go-runtime-study/sources/conc/pool/pool.go:108-115`; `pool/pool_test.go:49-64`; echoed in all four other pool tests |
| Leak containment (by construction) | ContextPool always cancels via defer in Wait; Stream closes queue then joins callbacker | `studies/aren-go-runtime-study/sources/conc/pool/context_pool.go:54-58`; `stream/stream.go:91-103` |
| Admitted-flaky timing bias | Three first-error tests document intrinsic test races and mitigate with 10ms/100ms sleeps instead of synchronization | `studies/aren-go-runtime-study/sources/conc/pool/context_pool_test.go:176-185`; `pool/result_error_pool_test.go:72-80`; `pool/result_context_pool_test.go:182-190` |
| Benchmark separation | `-benchtime 5s -timeout 0 -cpu 1 -benchmem`, excluded from unit runs via `-run=XXX` | `studies/aren-go-runtime-study/sources/conc/Makefile:24`; benchmarks `pool/pool_test.go:149-166`, `stream/stream_test.go:138-159`, `iter/iter_test.go:171-182` |
| Perf-regression gating in CI | Dedicated Benchmark workflow per PR and push-to-main; fails on alert vs cached main-branch JSON keyed by CPU model | `studies/aren-go-runtime-study/sources/conc/.github/workflows/bench.yml:31,63`; `.github/workflows/main.yml:27,39,52` |
| Iterator dynamic work distribution | Shared atomic index; goroutines pull next index in a loop — verified by 10k-element and huge-concurrency tests | `studies/aren-go-runtime-study/sources/conc/iter/iter.go:71-84`; `iter/iter_test.go:109-120`; concurrency>default proof `iter/iter_test.go:45-72` |
| Test seam for internal default | `export_test.go` exposes `DefaultMaxGoroutines` to external test package | `studies/aren-go-runtime-study/sources/conc/iter/export_test.go:1-3` |
| Empty-input edge case | ForEach/Map on empty slices never invoke f and do not panic | `studies/aren-go-runtime-study/sources/conc/iter/iter_test.go:78-98`; `iter/map_test.go:28-37,87-98` |

## Answers to Dimension Questions

### Which concurrency failures would ordinary coverage miss?

Four classes that conc explicitly targets:

1. **Data races**: invisible without `-race`; hence it is hard-coded in both local (`Makefile:20`) and CI (`go.yml:37`) test invocations rather than left to developer habit. The lock-free first-panic-wins design (`panics/panics.go:29`) would be nearly impossible to validate otherwise.
2. **Swallowed child-goroutine panics**: an unprotected `go func(){...}()` panic kills the process; coverage shows nothing. Every component asserts propagation through its terminal barrier (`waitgroup_test.go:67-74`, `pool/error_pool_test.go:81-94`, `stream/stream_test.go:101-119`, `iter/iter_test.go:89-98`).
3. **Bounded-resource violations**: line/branch coverage cannot detect "pool spawned 11 goroutines when limit was 10". The atomic-counter limit tests make the violation an explicit assertion failure (`pool/pool_test.go:66-90`).
4. **Ordering violations under scheduler noise**: results-order correctness only fails when completions interleave; the random-sleep test (`result_pool_test.go:66-84`) actively manufactures that interleaving.

Also notable: cross-generation state leakage (errors/results bleeding between Wait cycles) — found in the wild as issue #128 and now pinned by regression tests in all four affected pool types.

### Can the suite force cancellation at every lifecycle boundary?

Mostly, within one pool generation. Covered boundaries: parent-context cancellation (`context_pool_test.go:95-105`), timeout (`context_pool_test.go:107-117`), first-task-error (`context_pool_test.go:134-147`), fail-fast (`context_pool_test.go:192-205`), and panic-as-cancellation-trigger (`context_pool_test.go:207-226`) — the last being the most adversarial, since it verifies siblings actually observed `ctx.Done()` while the panic still propagates. Not covered (no evidence found): cancellation while `Go()` blocks on channel backpressure (`pool/pool.go:55-66`), cancellation during a reused pool's second generation, or a task ignoring context and outliving Wait — the library relies on tasks honoring ctx, and nothing falsifies misbehavior there.

### How does it prove absence of duplicate terminal events, leaks, and partial commits?

- **Duplicate terminal events**: structurally prevented, not just tested — `CompareAndSwap(nil, &rp)` keeps only the first panic (`panics/panics.go:29`), and `WithFirstError` returns `errs[0]` (`pool/error_pool.go:43-44`). Tests assert the second error is absent (`result_context_pool_test.go:194-195`). Caveat: which error is "first" depends on scheduling unless cancel-on-error synchronizes; the deliberate add-before-cancel ordering comment (`context_pool.go:40-45`) addresses the causal-loss variant, and the tests admit residual timing sensitivity.
- **Leaks**: proven indirectly and incompletely. Concurrency counters must return to 0 after Wait (`pool_test.go:87`), ContextPool cancels unconditionally via defer (`context_pool.go:56`), Stream closes its queue and joins the callbacker even on panic propagation (`stream.go:96-99`). But there is **no goleak-style or `runtime.NumGoroutine()` check anywhere** — a leaked goroutine from a user forgetting Wait() would go undetected, despite goal #1 of the source's own README being "Make it harder to leak goroutines" (`studies/aren-go-runtime-study/sources/conc/README.md:45,49`).
- **Partial commits**: no storage semantics here, but the analogous invariant — "every submitted sibling completes even when one task panics" — is asserted: non-panicking siblings' effects are checked after the propagated panic (`waitgroup_test.go:100-115`), and Stream guarantees all callbacks still execute when one task/callback panics (`stream.go:30-32` doc contract, `stream_test.go:110-119`).

### Which techniques are proportionate for Aren Phase 1 rather than infrastructure-scale theatre?

Proportionate and worth copying verbatim:
- `-race` as a non-negotiable presubmit flag (one-line cost, catches whole bug classes).
- Atomic-counter limit falsification parameterized over {1, 10, N} (cheap, deterministic assertions on final state).
- A uniform panic-propagation matrix repeated per component, including panic×cancellation interaction.
- Regression tests named after real issues (#128 pattern) to keep fixes from rotting.
- Benchmarks excluded from fast presubmits (`-run=XXX`, `Makefile:24`) and gated separately with regression alerts keyed to hardware (`bench.yml`, `main.yml`).

Theatre for Phase 1 (correctly absent from conc): fuzz targets for pure orchestration code, metamorphic testing, chaos/crash-recovery harnesses — these pay off for parsers/state machines/persistence, not for a WaitGroup wrapper.

Missing-but-cheap additions Aren should adopt on top: goleak in TestMain, watchdog timeouts converting hangs into diagnosable failures, and deterministic sequencing replacing sleep-biased ordering tests.

## Architectural Decisions

1. **Terminal-barrier failure model**: every component funnels failures (panics, errors) through a single `Wait()` call, making the terminal barrier the natural place to assert failure semantics (`waitgroup.go:36-52`, `pool/error_pool.go:33-48`, `stream/stream.go:89-103`).
2. **Lock-free first-writer-wins for panics**: an `atomic.Pointer` CAS means concurrent panics need no mutex and no dedup logic downstream (`panics/panics.go:15-31`).
3. **Index-reservation result aggregation**: submission order is decoupled from completion order by reserving slots under a mutex and writing by index, giving deterministic output under nondeterministic scheduling (`pool/result_pool.go:32-37,96-120`).
4. **Panic-isolated worker replacement**: a panicking worker releases the limiter slot on defer, and the pool spawns fresh workers, so panics can't permanently shrink capacity (`pool/pool.go:50-53,150-153`).
5. **Zero-value usability via lazy `sync.Once` init**, plus defensive panics on post-init reconfiguration instead of silent misconfiguration (`pool/pool.go:100-115`).

## Notable Patterns

- **Test matrix cloning**: identical subtests (`panics on configuration after init`, `limit`, `reuse`, `propagates panics`) replicated across all six pool types — high boilerplate, but guarantees uniform failure semantics and makes omissions visible in diffs.
- **Adversarial-comment honesty**: tests that cannot be made fully deterministic say so inline and explain the mitigation (`context_pool_test.go:176-185`, `result_error_pool_test.go:72-80`).
- **Starvation-proofing via synthetic events**: stream injects an empty callback on panic so the ordered consumer never deadlocks waiting for a callback that will never arrive (`stream/stream.go:73-81`).
- **External test packages everywhere** (`package *_test`): all tests exercise only public APIs, so the suite doubles as an API-drift detector.
- **Hardware-keyed benchmark baselines**: CI caches benchmark JSON by commit SHA + OS + CPU model to avoid comparing across machines (`main.yml:49-52`).

## Tradeoffs

- **Real time vs determinism**: using actual `time.Sleep`/`WithTimeout` keeps tests simple and seam-free, but couples pass/fail to machine load; slow CI runners shrink the safety margins (10–100ms biases). An injected clock would trade ~zero complexity for reproducibility.
- **Uniformity vs precision**: cloned matrices across pool variants maximize consistency but duplicate maintenance; a fix to one variant's expectations must be mirrored five times.
- **Structural prevention vs behavioral proof**: CAS-based first-panic-wins and index-slot reservation make some bugs impossible, reducing test burden — but they concentrate correctness in subtle primitives whose own verification rests almost entirely on the race detector.
- **Coverage vs leak assurance**: asserting counters return to 0 proves workers exited, yet the suite never proves the *library* leaks nothing under early-abandonment (no Wait called) — the exact scenario implied by the README's stated goal #1 (`studies/aren-go-runtime-study/sources/conc/README.md:45`).

## Failure Modes / Edge Cases

- **First-error ambiguity without cancel-on-error**: with `WithFirstError` alone, the "first" error is whichever task registers first under the scheduler, not causally meaningful; tests paper over this with sleeps (documented at `result_error_pool_test.go:72-80`).
- **Blocking `Go()` under saturation**: when all workers are busy, `Go()` blocks on the tasks channel (`pool/pool.go:38,62-63`); combined with cancel-on-error this can delay cancellation delivery until a worker frees up — untested.
- **Callback-channel capacity assumption**: Stream sizes its queue as `MaxGoroutines()+1` (`stream.go:112`) — correct today because pool workers ≤ MaxGoroutines, but silently breaks if pool internals change; no test pins this coupling.
- **Panic value `nil`**: `tryRecover` treats `recover() == nil` as no panic (`panics/panics.go:27`); `panic(nil)` in modern Go becomes `*runtime.PanicNilError`, so this is safe on Go ≥1.21 but was lossy historically — version-sensitive edge worth noting since `go.mod` declares 1.20.
- **Reuse misuse window**: after Wait resets `initOnce` (`pool.go:77-82`), configuration methods still panic correctly (tested `pool_test.go:57-63`), but calling `Wait()` twice concurrently is undefined and untested.

## Future Considerations

For Aren Phase 1 adoption, the concrete engineering work suggested by this source:

1. Add `goleak.VerifyTestMain` to every package to convert leak-prevention-by-convention into an enforced invariant (fills conc's biggest gap relative to its own stated goal).
2. Replace sleep-biased ordering tests with deterministic synchronization (e.g., an injected error-registration hook or sequencer channels), eliminating the three documented flaky-race tests.
3. Keep conc's presubmit shape: `-race` mandatory, benchmarks in a separate hardware-keyed workflow with `fail-on-alert` regression gates.
4. Add a hang-watchdog (test binary timeout with goroutine dump) so scheduling deadlocks produce diagnosable artifacts rather than CI timeouts.
5. Adopt the panic-propagation × cancellation interaction test (`context_pool_test.go:207-226`) as a required template for any Aren component with both panics and contexts.

## Questions / Gaps

- **No evidence found** for fuzz testing: searched for `Fuzz*` functions and `-fuzz` flags across all `_test.go` files and workflows — none exist.
- **No evidence found** for goroutine-leak detection: searched for `goleak`, `NumGoroutine`, `runtime.NumGoroutine` — no matches; leak claims rest on structural design only.
- **No evidence found** for fault injection seams: no clock/transport/storage/scheduler abstractions exist anywhere in the source; the only test seam is `iter/export_test.go:1-3`.
- Git history in this snapshot is shallow (single commit), so historical race/deadlock fixes could not be enumerated beyond the in-test references to issue #128 (`pool/error_pool_test.go:122` et al.).
- Search boundary: analysis restricted to `studies/aren-go-runtime-study/sources/conc/` per isolation rules; the dimension's referenced Aren acceptance docs were not read (cross-source access banned).

---

Generated by `01.03-adversarial-concurrency-and-failure-verification` against `conc`.
