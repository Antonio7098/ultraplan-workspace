# Source Analysis: pebble

## 01.03 Adversarial Concurrency and Failure Verification

### Source Info

| Field | Value |
|-------|-------|
| Name | pebble |
| Path | `studies/aren-go-runtime-study/sources/pebble` |
| Language / Stack | Go (LevelDB/RocksDB-style KV storage engine used by CockroachDB) |
| Analyzed | 2026-08-26 |

## Summary

Pebble verifies concurrency and failure invariants with a layered adversarial stack rather than happy-path unit tests. The core techniques are:

1. **A virtual filesystem as the single fault-injection seam.** All disk I/O flows through `vfs.FS`; `errorfs` wraps any FS to inject typed errors or latency (`vfs/errorfs/errorfs.go:279-296`, `vfs/errorfs/dsl.go:18-47`), and a crashable in-memory FS produces crash-consistent clones that simulate power loss at any instant, including partial retention of unsynced data (`vfs/mem_fs.go:129-160`).
2. **Metamorphic testing.** Random operation sequences are generated deterministically from a seed and replayed against many option configurations (standard + randomized); per-op histories are compared after normalizing multi-threaded output order (`metamorphic/meta.go:158-356`, `metamorphic/history.go:195-227`). Runs exercise restarts, error injection, WAL failover, multiple DB instances, and 1..GOMAXPROCS threads.
3. **Crash-recovery sweeps.** Tests iterate the "crash at the k-th write" index across an entire workload, clone the FS at the crash point, and require the store to reopen cleanly and match tracked KV state (`error_test.go:368-431`, `error_test.go:433-569`, `open_test.go:1914-2059`).
4. **Tiered execution and gating.** Presubmit runs unit tests under the `invariants` build tag plus a dedicated race-detector job; nightly jobs run per-package stress reruns, ASAN/MSAN/race instrumented builds, cross-version metamorphic tests, and QEMU s390x (`.github/workflows/ci.yaml`, `.github/workflows/stress.yaml`, `.github/workflows/instrumented.yaml`, `Makefile:34-84`).
5. **Reproducibility discipline.** Every randomized test derives behavior from seeds it logs; failures print seed, ops, options, history tail, and a ready-to-paste reproduction/reduction command (`metamorphic/meta.go:258-280`, `meta.go:331-351`). A delta-debugging reducer shrinks failing op sequences even when reproduction is flaky (`internal/metamorphic/reduce_test.go:232-261`).

Invariant checks (`invariants.Enabled`) are compiled into every test build and auto-enabled under `-race` (`internal/invariants/on.go:5`, `internal/invariants/invariants.go:16-17`), so the race build also gets assertion-dense code paths.

## Rating

**Rating: 9/10**

Rationale against the dimension's purpose ("repeatable Go techniques for proving invariants under races, panics, partial writes, restarts, resource exhaustion, scheduling variation"):

- Crash consistency is simulated repeatably and cheaply — no real disks, no process kills, no sleeps-as-sync-points — via `CrashClone` with configurable unsynced-data survival (`vfs/mem_fs.go:141-160`) and deterministic per-path randomness (`vfs/errorfs/latency.go:106-147`).
- Failure diagnosis is engineered end-to-end: seed printing, run-dir replay, child-process isolation of failing configs, automatic op reduction tolerant of nondeterminism, LSM diagram generation on reduction (`internal/metamorphic/reduce_test.go:207-222`).
- Expensive suites are separated from presubmit and gated by schedule, not hope: stress, ASAN/MSAN, cross-version, s390x all live in nightly workflows while race + invariants stay required on PRs (`.github/workflows/ci.yaml:82-93` vs `.github/workflows/stress.yaml:49-53`).
- Deductions from 10: no native Go fuzz targets exist anywhere in the tree (searched `func Fuzz` and case-insensitive "fuzz": zero hits); there is no clock seam, so timing-dependent failover tests rely on microsecond durations and hard caps (`open_test.go:1923-1931`, `latency.go:83-86`); some injectors rely on call-stack matching that their own docs call fragile (`vfs/errorfs/dsl.go:39-44`).

## Evidence Collected

Every entry includes file path with line numbers, relative to the source root `studies/aren-go-runtime-study/sources/pebble`.

| Area | Evidence | File:Line |
|------|----------|-----------|
| Invariant build tier | `//go:build invariants \|\| race` enables `Sometimes`, `CloseChecker`, buffer mangling in race builds too | internal/invariants/on.go:5 |
| Invariant gating policy | `Enabled = buildtags.Race || buildtags.Invariants`; guidance to use `Sometimes()` so production paths stay covered | internal/invariants/invariants.go:16-23 |
| Finalizer leak assertions | `SetFinalizer` asserts cleanup of DBs, memtables, locks, cache entries; disabled under race due to detector bugs | internal/invariants/invariants.go:26-45, db.go:1861, db.go:2648, internal/base/directory_lock.go:88, internal/cache/cache.go:138 |
| Stale-buffer falsification | `MaybeMangle`/`BufMangler.MaybeMangleLater` trash returned buffers 10%/50% of the time to catch aliasing bugs | internal/invariants/on.go:86-118 |
| Invalidating iterator | Wraps iterators in invariant builds 10% of the time, overwriting last-returned KV with 0xFF | internal/invalidating/iter.go:20-24 |
| Race-only regression test | `TestNodeArenaEnd` exists solely because race builds check pointer alignment at arena boundary | internal/arenaskl/race_test.go:15-21 |
| Fault-injection FS | `ErrInjected`, 24 typed `OpKind`s, read/write classification, `Wrap`, `Any`, `Counter`, `Toggle` injectors | vfs/errorfs/errorfs.go:20-91, 183-296 |
| Injection DSL predicates | `PathMatch`, `CallStackIncludes` (self-described fragile), `Reads`, `Writes`, offset-matching | vfs/errorfs/dsl.go:33-74, 90-102 |
| Deterministic random injection | `Randomly(p, seed)` is per-file-path deterministic despite nondeterministic concurrency | vfs/errorfs/dsl.go:104-116 |
| Latency injection | `RandomLatency`: exponential distribution, capped at 20x mean, lifetime cap; keyed per-path PRNG | vfs/errorfs/latency.go:26-34, 83-104, 106-147 |
| Crash simulation primitive | `NewCrashableMem` + `CrashClone(CrashCloneCfg{UnsyncedDataPercent, RNG})` produce post-crash FS states | vfs/mem_fs.go:37-61, 129-160 |
| Crash-at-kth-write sweep | `TestDBWALRotationCrash` increments k until no k-th write exists, reopening from cloned state each iteration | error_test.go:368-431 |
| Compaction crash sweep | `TestDBCompactionCrash`: random concurrency, latency injection, `UnsyncedDataPercent: 10` torn crashes, reopen must succeed | error_test.go:433-569 |
| Concurrent randomized crash test | `TestWALFailoverRandomized` + `runRandomizedCrashTest`: concurrent batch writers, weighted crash ops, tri-state KV model | open_test.go:1914-2059 |
| Tri-state KV tracking | `kvUnset/kvMaybeSet/kvSet` handles the race between `CrashClone` and in-flight syncs; `validateState` checks recovered DB | open_test.go:1980-2039 |
| Manifest crash test | `TestCrashDuringManifestWrite_LargeKeys` crashes mid-manifest-write concurrently | version_set_test.go:515-537 |
| Checkpoint crash test | `TestCheckpointFlushWAL` verifies checkpointed WAL survives a zero-unsynced-data crash clone | checkpoint_test.go:375-410 |
| Ingest crash test | `TestIngestFileNumReuseCrash` uses a crash-wrapping FS around ingest | ingest_test.go:3043-3099 |
| FMV ratchet crash matrix | Format-major-version upgrades re-run under `CrashClone{UnsyncedDataPercent: 0}` | format_major_version_test.go:130-191 |
| Atomic marker crash tests | `atomicfs` marker tested across repeated crash clones | vfs/atomicfs/marker_test.go:170-196 |
| MemFS crash-concurrency test | `TestMemFSCrashCloneConcurrency` races readers/writers against concurrent `CrashClone` write-lock holders | vfs/mem_fs_test.go:146-206 |
| Metamorphic entrypoint | `TestMeta`, `TestMetaTwoInstance`, `TestMetaCockroachKVs`; goroutine-leak test deferred unless disabled | internal/metamorphic/meta_test.go:55-74 |
| Deterministic generation | Ops generated from PCG-seeded RNG; seed defaults to wall-clock but is printed on failure for exact replay | metamorphic/meta.go:150-176, 258-260 |
| Option-space exploration | `standardOptions` fixed configs (disable_wal, tiny caches/memtables, max_open_files=1...) plus equal count of `random-%03d` configs run in parallel | metamorphic/meta.go:286-301, 358-404, options.go:468-509 |
| Multi-threaded execution | `Execute` errgroup; ops hash to threads by receiver object; synchronization points computed per-op | metamorphic/meta.go:634-693, test.go:585-655 |
| Thread count randomization | `testOpts.Threads = rng.IntN(runtime.GOMAXPROCS(0)) + 1` in random configs | metamorphic/options.go:844 |
| History normalization | Every recorded line suffixed `#<op>`; `reorderHistory` restores op order so thread interleaving doesn't affect comparison | metamorphic/history.go:52-75, 195-227 |
| Restart-with-crash op | `restartDB` clones FS with `UnsyncedDataPercent: 0` before closing, then reopens — restarts see post-crash state | metamorphic/test.go:325-384 |
| Strict-FS option | `strictFS` opts into Sync writes + `vfs.NewCrashableMem()`; random configs flip it 50% | metamorphic/options.go:78-82, 833-841 |
| Error injection in metas | `InjectErrorsRate` run option; `RetryInjected` policy retries only `ErrInjected`; retryable iterators restore last-key position | metamorphic/meta.go:438-443, retryable.go:17-60 |
| Background-error tripwires | Event listeners exit the process on any unexpected background error, cloning in-memory state for post-mortem first | metamorphic/test.go:127-179 |
| Per-op watchdog | `OpTimeout` panics if any single op exceeds budget (4x/8x multipliers for slow op classes) | metamorphic/test.go:430-447 |
| Op mix | 48 op types incl. restart, checkpoint, ingest/excise, range keys, external ingest; second preset biases toward many versions of same prefix | metamorphic/config.go:21-71, 124-145 |
| Failure reduction | Delta-debugging reducer removes ops probabilistically, tolerates flaky repro, simplifies keys, emits LSM diagram | internal/metamorphic/reduce_test.go:232-261, 263-278 |
| Leak detection harness | `leaktest.AfterTest(t)` from crllib used at 432 test sites; disable flag `--noleaktest` for reducer runs | internal/metamorphic/meta_test.go:56-58, internal/metamorphic/metaflags/meta_flags.go:55-56 |
| Resource leak accounting | `DB.Close` errors on leaked iterators (with stacks), snapshots, and memtable reservations | db.go:1905-1968 |
| WAL fd-leak check | `wal/reader_test.go` asserts no files left open after `Copy` error | wal/reader_test.go:348-352 |
| Internal consistency checker | `DB.CheckLevels` validates sequence numbers, range-del stacking across the LSM | level_checker.go:29, 515 |
| Debug-check hook | `Options.DebugCheck` invoked on every new version install; typically `DebugCheckLevels` | options.go:606-609, 1480-1483 |
| Stress reruns | `stress -exec` wrapper reruns each test up to 1000 times; `stressmeta` runs TestMeta under stress | Makefile:62-71 |
| Presubmit vs nightly split | PR CI: invariants test, no-invariants, no-cgo, race, macOS, stress-new-tests; Nightly: per-package stress, race+asan+msan, crossversion, s390x | .github/workflows/ci.yaml:22-128, .github/workflows/tests.yaml, .github/workflows/stress.yaml, .github/workflows/instrumented.yaml |
| Stress budget tuning | Stateful packages (root, manifest, metamorphic, sstable, wal) get 30m/maxruns=1000; others 5m | scripts/stress.sh |
| New-test stress gate | PR job diffs new `func Test` names vs base branch and stresses them 1000 runs/10m before merge | scripts/stress-new-tests.sh, .github/workflows/ci.yaml:108-128 |
| Cross-version metamorphic | `ExtendPreviousRun` seeds ops/state from previous version's run; `TestMetaCrossVersion` chains release binaries; nightly + optional stress mode | metamorphic/meta.go:63-87, internal/metamorphic/crossversion/crossversion_test.go:38-70, Makefile:73-84 |
| Bit-flip diagnostics | Checksum failures report which bit flipped (`bitflip.CheckSliceForBitFlip`) — diagnosability built into production paths | record/record.go:391, sstable/block/block.go:192 |

## Answers to Dimension Questions

**Which concurrency failures would ordinary coverage miss?**

- *Stale-buffer reuse.* Iterator KV slices are valid only until the next operation; ordinary coverage passes because nothing overwrites the memory. Pebble falsifies this by mangling previously returned buffers randomly in invariant builds (`internal/invariants/on.go:106-118`) and wrapping ~10% of internal iterators with an invalidating iterator that overwrites the last returned KV (`internal/invalidating/iter.go:20-24`). A consumer holding a stale slice sees 0xCC bytes instead of silently correct data.
- *Races at allocation boundaries.* `TestNodeArenaEnd` exists purely because race builds perform alignment checks at the skiplist arena boundary (`internal/arenaskl/race_test.go:15-21`) — a class of bug invisible without `-race`.
- *Torn persistence windows.* Coverage never observes data loss between `write` and `sync`. `CrashClone(UnsyncedDataPercent)` materializes exactly those windows, including partial survival of unsynced blocks (`vfs/mem_fs.go:129-160`), and sweep loops walk every k-th-write crash point (`error_test.go:538-568`).
- *Scheduling-induced failover.* WAL-failover paths only trigger under slow primary I/O; `RandomLatency` injection with microsecond thresholds forces failover deterministically per path (`open_test.go:1920-1937`, `wal/failover_writer_test.go:40-126`).
- *Multi-threaded interleavings of public API calls.* The metamorphic runner executes one op sequence across 1..GOMAXPROCS threads with computed synchronization points, then compares normalized histories (`metamorphic/meta.go:634-693`, `test.go:585-655`), catching state-machine bugs in object handoff (batches applied from another goroutine, iterator close racing snapshot close) that single-threaded runs never reach.

**Can the suite force cancellation at every lifecycle boundary?**

Partially — by close/restart/crash rather than by context cancellation:

- Compaction cancellation is exercised through `ErrCancelledCompaction`, returned from multiple compaction phases and retried internally (`compaction.go:42-44`, `compaction.go:2394`, `2484`, `2662`, `2867`); the metamorphic suite whitelists it as non-fatal (`metamorphic/test.go:128`).
- Lifecycle boundaries are forced by generated ops: `OpDBRestart` closes and reopens mid-sequence on a crash-cloned FS (`metamorphic/config.go:30`, `test.go:325-384`), `OpDBClose`/iter/snapshot/batch close ops are part of the op space (`config.go:32-33,48-56`).
- Hard interruption at arbitrary points comes from crash simulation, not cancellation: any background goroutine's next FS write can hit a crash-cloned world (`error_test.go:450-458`).
- No evidence found of a seam that cancels `Open` mid-flight via context: searches for a clock/context-cancellation injection point in Open found none; `Open` has no exported ctx parameter, and no test forces cancellation between Open's internal phases except via injected FS errors (e.g., `TestOpenCrashWritingOptions`, `open_test.go:802`; `TestCrashOpenCrashAfterWALCreation`, `open_test.go:1230-1260`). This is a genuine gap relative to "every lifecycle boundary."

**How does it prove absence of duplicate terminal events, leaks, and partial commits?**

- *Partial commits:* the randomized crash framework tracks every key in a tri-state model (`kvUnset`/`kvMaybeSet`/`kvSet`, `open_test.go:1980-2015`) — `kvMaybeSet` acknowledges the race between a crash clone and in-flight syncs — and after recovery walks the full keyspace asserting engine contents agree with the model (`validateState`, `open_test.go:2018-2039`). Manifest durability gets a dedicated concurrent-crash test (`version_set_test.go:515-537`); checkpoints assert the flushed WAL is present and non-empty in the crash clone (`checkpoint_test.go:375-410`); ingest asserts file-number reuse after crash cannot resurrect stale files (`ingest_test.go:3062-3099`).
- *Leaks:* three layers — (1) goroutine leak detection via `leaktest.AfterTest` at 432 test sites (e.g., `open_test.go:1915`); (2) resource accounting at `Close` that errors on leaked iterators with stack traces, snapshots, and memtable reservations (`db.go:1905-1968`); (3) finalizer-based assertions on unclosed caches/locks in invariant builds (`internal/cache/cache.go:138`, `internal/base/directory_lock.go:88`).
- *Duplicate terminal events:* proven indirectly by metamorphic history equality — every op result line is suffixed with its op index (`history.go:63-75`), histories are reordered by op index and diffed line-by-line across configurations (`history.go:195-212`, `meta.go:312-355`); a duplicated or missing terminal event shifts the sequence and fails the compare. Double-close is additionally falsified by `CloseChecker` panics in invariant builds (`internal/invariants/on.go:35-43`). No direct idempotence harness for terminal-event delivery was found; the guarantee rests on the history-equality argument above.

**Which techniques are proportionate for Aren Phase 1 rather than infrastructure-scale theatre?**

Proportionate, high-leverage per unit of infrastructure:

1. Build-tag invariant tier with `race ⇒ invariants` (`internal/invariants/on.go:5`) — one file pair, raises assertion density everywhere.
2. An `errorfs`-style injectable FS with `ErrInjected`, predicate DSL, and a retry policy that retries only injected errors (`vfs/errorfs/errorfs.go:20-24`, `metamorphic/retryable.go:17-33`) — small, dependency-free, reusable across Aren transports/storage.
3. A crashable MemFS with `CrashClone` + k-th-write sweep loops for WAL/manifest/checkpoint recovery tests (`vfs/mem_fs.go:141-160`, `error_test.go:538-568`) — this replaces whole fault-injection labs.
4. A miniature metamorphic loop: seeded op generation → execution → serialized history → compare across 2–3 configs, plus seed logging on failure (`metamorphic/meta.go:158-207`, `258-280`).
5. Pervasive `leaktest.AfterTest` + close-time leak counters (`db.go:1905-1968`).
6. A dedicated race job in presubmit and stress reruns of *newly added* tests only (`scripts/stress-new-tests.sh`).

Defer as theatre for Phase 1: cross-release binary chaining (`Makefile:73-84`), ASAN/MSAN farms (cgo-specific), s390x QEMU (`Makefile:86-102`), and the full delta-debugging reducer — valuable later, none block proving Phase 1 invariants.

## Architectural Decisions

1. **The FS interface is the universal failure seam.** Errors, latency, and crashes are all modeled below `vfs.FS` (`vfs/vfs.go` interface; wrappers at `vfs/errorfs/errorfs.go:291-296` and `vfs/mem_fs.go:61`), so the production code contains no test hooks for I/O faults — injection composes externally. This keeps fault paths exercised through the same call sites production uses.
2. **Determinism is engineered per-key, not global.** Injectors derive randomness from `(seed, file path)` via `keyedPrng` (`vfs/errorfs/latency.go:106-147`, `dsl.go:104-116`), so injected behavior is reproducible per file regardless of goroutine interleaving — determinism despite concurrency rather than instead of it.
3. **Child-process isolation for metamorphic configs.** Each OPTIONS variant runs as a fresh `os.Args[0]` process (`metamorphic/meta.go:212-248`), containing panics/exit calls and enabling coverage instrumentation of the inner binary; the parent compares persisted histories.
4. **Assertions compiled out, never deleted.** `invariants` package swaps implementations by build tag (`internal/invariants/on.go` vs `off.go`); `Enabled` is deliberately auto-on under `-race` (`internal/invariants/invariants.go:16-17`) so the race build doubles as the assertion-dense build.
5. **Fail-fast with evidence capture.** Unexpected background errors clone in-memory state to disk before exiting (`metamorphic/test.go:122-135`); failures print seed/ops/options/history and the exact reduce command (`meta.go:258-280`, `331-351`).

## Notable Patterns

- **Tri-state expectation modeling under crash races:** `kvMaybeSet` encodes "we believe we synced, but the crash clone may disagree" (`open_test.go:2001-2009`) — a template for verifying at-least-once semantics without false positives from inherent races.
- **k-sweep crash testing:** increment the crash index until it exceeds the number of writes, guaranteeing every persistence window in the workload is visited (`error_test.go:410-430`, `538-568`).
- **Watchdog timers per op with class-based budgets:** panics identify the exact op string when exceeded (`metamorphic/test.go:430-447`) — turns hangs into diagnosable failures.
- **History normalization:** suffixing lines with `#<op>` and reordering post-hoc makes multithreaded output comparable without serializing execution (`history.go:63-75`, `195-212`).
- **Injector composition:** `Any(...)`, `Toggle`, `Counter` let tests layer random errors + latency and assert injected errors actually surfaced to users (`vfs/errorfs/errorfs.go:183-277`).
- **Reduction robust to flakiness:** decaying removal probability works even when the bug reproduces intermittently (`internal/metamorphic/reduce_test.go:232-251`).

## Tradeoffs

- **In-memory FS changes performance characteristics,** which alters exercised code paths/timings; Pebble knowingly accepts this and keeps exactly one standard disk-backed config, cloning in-memory state to disk only for post-mortems (`metamorphic/test.go:122-126`, `options.go:468-472`).
- **No clock seam.** Failover probes use real time at microsecond granularity (`open_test.go:1923-1931`); correctness depends on duration caps (20x mean, `latency.go:83-86`) and total-latency limits (`latency.go:89-100`) to avoid timeouts — simpler than clock abstraction, but tests remain mildly load-sensitive.
- **Call-stack predicates are fragile by design admission** (`vfs/errorfs/dsl.go:39-44`) — renames silently break targeted injections.
- **Metamorphic breadth vs runtime:** default runs generate 1,000–10,000 ops × ~2×N option configs in subprocesses (`meta.go:140-141`, `381-391`); this is why deep variants (stress, crossversion) are nightly-gated.
- **Finalizer-based leak detection is disabled under race builds** due to historical race-detector bugs (`internal/invariants/invariants.go:26-31`) — the two detectors never run together.

## Failure Modes / Edge Cases

- `_meta/<test>` failure artifacts are overwritten by subsequent invocations (`internal/metamorphic/meta_test.go:31-33`) — evidence loss unless preserved.
- Reduction runs with `--noleaktest` and 10s timeouts (`internal/metamorphic/reduce_test.go:170-176`): a reduced case can drop the original leak signal.
- Windows sleep granularity makes crash sweeps slow; compensated with a latency injector skip on Windows (`error_test.go:465-470`).
- `Prefetch` errors are not yet injectable (`vfs/errorfs/errorfs.go:520-523` TODO) — an untested failure surface.
- Crash-test validation tolerates `kvMaybeSet`, so a bug causing *extra* durable keys beyond maybe-set would be caught, but subtle WAL-prefix-loss bugs that keep keys inside the maybe-set envelope could pass silently; the model trades strictness for stability against inherent sync races (`open_test.go:2001-2014`).
- `history.Recordf` panics on multi-line formats and records-after-close (`history.go:52-61`) — harness bugs surface loudly rather than corrupting comparisons.

## Future Considerations

- Native Go fuzzing (`func Fuzz` targets) is absent; fuzzing block/WAL record parsing (which already has bit-flip diagnostic machinery at `record/record.go:391`, `sstable/block/block.go:192`) would be a natural extension and is cheap to adopt in Go 1.18+ toolchains.
- A controllable-clock seam for failover/election-style timers would remove the residual load sensitivity noted above; Pebble compensates with caps instead.
- Context-based cancellation of `Open` and compaction grant loops could reuse the existing injected-error retry plumbing; today only FS-error and crash-clone paths interrupt startup.

## Questions / Gaps

- **No native fuzz targets:** searched all `*.go` for `func Fuzz` and case-insensitive `fuzz` — zero results. Corruption tolerance is covered indirectly by datadriven corpora and checksum diagnostics, not by structured fuzzing.
- **No clock abstraction:** searched for `type Clock`/clock seams — none; timing control relies on injected latency and real `time.Sleep`.
- **Duplicate terminal-event idempotence** is proven only via metamorphic history equality and double-close panics; no direct test asserts "event X delivered exactly once" semantics.
- **Cancellation during `Open`:** no mechanism found to cancel Open mid-phase other than injected FS errors; whether that matters depends on Aren's embedding requirements.
- Search boundary: analysis confined to `studies/aren-go-runtime-study/sources/pebble`; referenced Aren acceptance docs were outside the permitted source directory and were not consulted.

---

Generated by `01.03-adversarial-concurrency-and-failure-verification` against `pebble`.
