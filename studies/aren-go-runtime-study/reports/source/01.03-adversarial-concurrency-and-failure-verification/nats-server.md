# Source Analysis: nats-server

## 01.03 Adversarial Concurrency and Failure Verification

### Source Info

| Field | Value |
|-------|-------|
| Name | nats-server |
| Path | `studies/aren-go-runtime-study/sources/nats-server` |
| Language / Stack | Go 1.25+ (module `github.com/nats-io/nats-server/v2`), GitHub Actions CI, Antithesis SDK for runtime invariant assertions |
| Analyzed | 2026-08-26 |

## Summary

nats-server verifies concurrent and failure behavior through five complementary mechanisms rather than happy-path unit tests:

1. **Race-detector-gated CI with a parallel "no-race" suite.** Most CI jobs run `go test $RACE` (`scripts/runTestsOnTravis.sh:40`, `scripts/runTestsOnTravis.sh:48`, `scripts/runTestsOnTravis.sh:60`). Tests that are too slow or allocation-heavy under `-race` live in build-tag-separated files (`//go:build !race && !skip_no_race_tests ...` at `test/norace_test.go:14`, `server/norace_1_test.go:14`, `server/norace_2_test.go:14`) following a `TestNoRace*` naming convention (`test/norace_test.go:38-40`), executed by two dedicated CI jobs split via the `skip_no_race_1_tests` / `skip_no_race_2_tests` tags (`scripts/runTestsOnTravis.sh:19-33`).

2. **Restart/crash/recovery testing at every storage layer.** ~41 recovery/corruption tests in the filestore alone (e.g., `TestFileStoreNoPanicOnRecoverTTLWithCorruptBlocks` at `server/filestore_test.go:10690`, `TestFileStoreMsgBlockFirstAndLastSeqCorrupt` at `server/filestore_test.go:7430`), JetStream meta-group recovery (`TestJetStreamClusterMetaRecoveryLogic` at `server/jetstream_cluster_3_test.go:226`, meta snapshot rescue tests at `server/jetstream_cluster_3_test.go:13817-14001`), and a seeded 3-minute KV soak that randomly restarts servers mid-write (`TestLongKVPutWithServerRestarts`, `server/jetstream_cluster_long_test.go:37-190`).

3. **Fault-injection seams for transports and sockets.** An in-process raft transport hub supporting partitions, heals, and message hooks (`server/raft_transport_helpers_test.go:29-124`); a TCP proxy with controllable RTT/bandwidth (`netProxy`, `server/jetstream_helpers_test.go:1958-2090`); a `shortWriteConn` that truncates writes and returns `io.ErrShortWrite` to exercise partial-write handling (`server/server_test.go:1683-1700`); and fake socket/listener/injector plumbing used to fuzz the TLS handshake (`server/server_fuzz_test.go:32-247`).

4. **Invariant assertions embedded in production code** via the Antithesis SDK (`assert.Unreachable("WAL truncate lost commits", ...)` at `server/raft.go:4271`; commit-loss on truncate at `server/raft.go:4648`; filestore read/write/flush/sync error assertions at `server/filestore.go:5458`, `server/filestore.go:5467`, `server/filestore.go:8707-8732`), activated only under the `enable_antithesis_sdk` build tag and otherwise compiled to noops (`internal/antithesis/test_assert.go:15`, `internal/antithesis/noop.go:15-27`). Every `require_*` test helper doubles as an Antithesis assertion emitter (`server/test_test.go:60-201`), as do cluster wait helpers (`server/jetstream_helpers_test.go:1466-1759`) and `checkFor` timeouts (`server/server_test.go:57-64`).

5. **Tiered suite economics.** Presubmit is a matrix of ~16 focused jobs, each capped at 30 minutes with `-p=1 -failfast -count=1 -timeout=30m` (`scripts/runTestsOnTravis.sh:25-136`, `.github/workflows/tests.yaml:102-372`); expensive JetStream long tests are excluded from presubmit by the `include_js_long_tests` tag (`server/jetstream_cluster_long_test.go:15`) and run nightly with `-race -shuffle on -failfast -timeout=60m` (`.github/workflows/long-tests.yaml:33`).

The scale is large: roughly 3,550 test functions in the `server` package plus 363 in the integration `test` package (measured by `rg -c "^func Test" server/*_test.go`), including 156 `TestNoRace*` functions.

## Rating

**8 / 10.**

Rationale: This is close to a reference implementation of adversarial concurrency verification for a Go network daemon. It maps specific techniques to specific invariants: races → `-race` CI plus a deliberately complementary no-race suite; data loss under truncation → production-code `assert.Unreachable` guards (`server/raft.go:4271`, `server/raft.go:4648`); bit rot / torn state → direct on-disk corruption mutation tests (`server/filestore_test.go:10690-10731`); network faults → partition-capable transport hub and RTT/bandwidth proxies (`server/raft_transport_helpers_test.go:67-85`, `server/jetstream_helpers_test.go:1958-2090`); partial writes → injected short-write conns (`server/server_test.go:1687-1700`); leaks → goroutine-count baselines and subscription-count polling (`server/norace_2_test.go:3550-3580`, `test/norace_test.go:530-562`). Long-running chaos is seeded and reproducible (`const Seed = 123456` at `server/jetstream_cluster_long_test.go:39` with `rng := rand.New(rand.NewSource(Seed))` at `server/jetstream_cluster_long_test.go:62`).

It loses two points because: (a) fuzz targets exist but no in-repo workflow continuously runs them (searched all of `.github/workflows/` for "fuzz" — no matches), so they depend on individual developers running `go test -fuzz`; (b) there is no clock seam — all timing is real wall-clock with polling loops, so scheduling variation is tolerated statistically rather than controlled, and restarts are graceful `Shutdown()`+restart rather than kill -9 crash simulation (`server/jetstream_cluster_long_test.go:163-165`); and (c) the race sanitizer is disabled for PR and main/release push events per the env expression at `.github/workflows/tests.yaml:12`, so the highest-churn path (PRs) gets no race-detector run of the JS/store/raft suites inside this repository (a `skipIfBuildkite` helper at `server/test_test.go:394-398` implies part of the pipeline is externalized).

## Evidence Collected

Every entry cites a file path with line numbers, relative to `studies/aren-go-runtime-study/sources/nats-server/`.

| Area | Evidence | File:Line |
|------|----------|-----------|
| Race gating in CI | `$RACE` env computed per event/ref; most jobs run `go test $RACE` | `.github/workflows/tests.yaml:12`, `scripts/runTestsOnTravis.sh:40` |
| No-race suite convention | Comment: tests prefixed `TestNoRace`, run without `-race` | `test/norace_test.go:38-40` |
| No-race build tags | `!race && !skip_no_race_tests && !skip_no_race_1_tests` (and `_2` variants) | `test/norace_test.go:14`, `server/norace_1_test.go:14`, `server/norace_2_test.go:14` |
| No-race job split | Two jobs run `TestNoRace` halves via mutually exclusive skip tags | `scripts/runTestsOnTravis.sh:19-33` |
| No-race throughput invariant | 100k subs across 2 servers, 120k replies via 50 senders; asserts zero slow consumers | `test/norace_test.go:42-229` |
| Memory-ceiling invariants | Cluster memory < 80MB; reply-permission growth < 20MB over 10k requests | `test/norace_test.go:295-344`, `test/norace_test.go:231-293` |
| Goroutine leak checks | Baseline `runtime.NumGoroutine()` compared after shutdown | `server/norace_2_test.go:3550-3580`, `test/routes_test.go:53-57`, `test/gosrv_test.go:25-51`, `server/jetstream_leafnode_test.go:2141-2159`, `server/filestore_test.go:10554-10568` |
| Subscription/connection leak check | Polls `Varz`/`Routez` until connections==0 and route NumSubs==0 | `test/norace_test.go:530-562` |
| Raft fault-injection hub | `raftTransportHub.partition/heal/healPartitions/setAfterMsgHook` simulate net partitions & intercept messages | `server/raft_transport_helpers_test.go:63-124` |
| Network proxy fault injection | `netProxy` with configurable RTT and up/down rates; `updateRTT` at runtime | `server/jetstream_helpers_test.go:1958-2090` |
| Partial-write injection | `shortWriteConn.Write` caps writes at 10 bytes and returns `io.ErrShortWrite`; `TestClientWriteLoopStall` proves no stall | `server/server_test.go:1683-1700`, `server/server_test.go:1702-1744` |
| Delayed/capture conn seams | `delayedWriteConn`, `capturePingConn`, `chunkWriteConn` for gateway/MQTT timing faults | `server/gateway_test.go:5245`, `server/gateway_test.go:6510`, `server/mqtt_test.go:7010` |
| Fake socket/listener for TLS | `FakeSocket`, `FakeConn`, `FakeListener`, `ClientHelloInjector` mutate ClientHello bytes in flight | `server/server_fuzz_test.go:32-247` |
| Parser fuzz target | `FuzzParser` feeds chunked (half-split) protocol buffers across CLIENT/ROUTER/GATEWAY/LEAF client kinds | `server/parser_fuzz_test.go:57-97` |
| TLS fuzz target | `FuzzServerTLS` mutates TLS version, hello payload, corrupts cert bytes; expects handshake success or clean rejection | `server/server_fuzz_test.go:284-409` |
| Subject-collision fuzz | `FuzzSubjectsCollide` with curated true/false corpus | `server/subject_fuzz_test.go:21-71` |
| Config fuzz entry | `conf/fuzz.go` `Fuzz(data []byte) int` for config parsing | `conf/fuzz.go:18` |
| Production invariant: WAL truncation | `assert.Unreachable("WAL truncate lost commits")` when `index < n.commit` | `server/raft.go:4271-4280` |
| Production invariant: lost commits | `assert.Unreachable("Truncate to earlier entry would lose commits")` | `server/raft.go:4648` |
| Production invariant: raft term regression | `assert.AlwaysOrUnreachable(ar.term <= n.term, "Raft response term mismatch")` | `server/raft.go:3967` |
| Production invariants: filestore I/O | `assert.Unreachable` on read/write/flush/sync errors in msg blocks | `server/filestore.go:5458`, `server/filestore.go:5467`, `server/filestore.go:7610-7619`, `server/filestore.go:8707`, `server/filestore.go:8732` |
| Production invariants: JS cluster | Write-error and reset-state assertions on meta/stream/consumer paths | `server/jetstream_cluster.go:1761`, `server/jetstream_cluster.go:3809`, `server/jetstream_cluster.go:4766`, `server/jetstream_cluster.go:7423` |
| Antithesis build-tag switch | SDK active only with `enable_antithesis_sdk`; noop otherwise | `internal/antithesis/test_assert.go:14-15`, `internal/antithesis/noop.go:14-27` |
| Assertion-instrumented helpers | All `require_*` helpers fire `AssertUnreachable` before failing | `server/test_test.go:60-201` |
| Timeout-as-violation pattern | `checkFor` and cluster waiters emit `AssertUnreachable("Timeout in ...")` | `server/server_test.go:57-64`, `server/jetstream_helpers_test.go:1466-1469`, `server/jetstream_helpers_test.go:1545-1550` |
| Restart-based chaos soak | `TestLongKVPutWithServerRestarts`: seeded RNG, random server restarts every 5s, fails if no successful update for 5s | `server/jetstream_cluster_long_test.go:37-190` |
| Random stop/start raft soak | `TestLongNRGChainOfBlocks` randomly starts/stops nodes to exercise recovery and snapshots | `server/jetstream_cluster_long_test.go:192-196` |
| Exactly-once ledger under duplicates | `TestLongClusterCLFSOnDuplicates` (nightly) plus CLFS tests in presubmit clusters | `server/jetstream_cluster_long_test.go:897`, `server/jetstream_cluster_3_test.go:6936`, `server/norace_1_test.go:4005` |
| Filestore corruption/recovery suite | 9 dedicated `TestFileStore*Corrupt*` tests; 36 corruption-related hits; e.g., manually corrupt block sequences then re-run `recoverPerMessageState()` | `server/filestore_test.go:6871`, `server/filestore_test.go:6965`, `server/filestore_test.go:7430`, `server/filestore_test.go:10690-10733` |
| Storage-permutation matrix | `testFileStoreAllPermutations` runs each store test across {no cipher/AES/ChaCha} × {none/S2 compression} | `server/filestore_test.go:56-72` |
| Meta-layer crash recovery | `TestJetStreamClusterMetaRecoveryLogic`, meta snapshot recovery/rescue/single-survivor tests | `server/jetstream_cluster_3_test.go:226`, `server/jetstream_cluster_3_test.go:13780-14001` |
| Long-test gating | `include_js_long_tests` tag excludes long tests from default builds; nightly job runs them with `-race -shuffle on -timeout=60m` | `server/jetstream_cluster_long_test.go:15`, `.github/workflows/long-tests.yaml:20-33` |
| Presubmit tiering | ~16 jobs (store/js×6/raft/no-race×2/mqtt/msgtrace/jwt/non-js/non-server) each `timeout-minutes: 30` | `.github/workflows/tests.yaml:102-372` |
| Test determinism knobs | Seeded RNG in soak test and benchmarks (`rand.NewSource(int64(seed))`) | `server/jetstream_cluster_long_test.go:39,62`, `server/jetstream_benchmark_test.go:342,595` |
| Raft sped up for tests | `init()` overrides heartbeat/election timeouts (hbInterval=50ms) to shrink failure windows | `server/jetstream_helpers_test.go:49-59` |
| Poll-until helper | `checkFor(totalWait, sleepDur, f)` retries predicate; used pervasively instead of fixed sleeps | `server/server_test.go:44-64`, `server/test_test.go:30`, `test/test_test.go:30` |
| Lock-ordering discipline | Documented global lock order `jetStream -> jsAccount -> Server -> client -> Account` etc. to structurally prevent deadlocks | `locksordering.txt:1-24` |
| External CI hint | `skipIfBuildkite(t)` skips tests when running on Buildkite CI | `server/test_test.go:394-398` |

## Answers to Dimension Questions

### Which concurrency failures would ordinary coverage miss?

Ordinary (non-adversarial) coverage would miss exactly the classes nats-server targets explicitly:

- **Slow-consumer backpressure at scale**: `TestNoRaceRouteSendSubs` drives 100k subscriptions and 120k replies through 50 sender connections and asserts zero slow consumers (`test/norace_test.go:118-217`) — impossible under `-race` timing, hence quarantined in the no-race suite.
- **Partial writes stalling the write loop**: `shortWriteConn` injects 10-byte truncated writes returning `io.ErrShortWrite` (`server/server_test.go:1687-1700`), and `TestClientWriteLoopStall` (`server/server_test.go:1702-1744`) falsifies the "write loop stalls after short write" invariant.
- **WAL truncation losing committed entries**: guarded directly in production by `assert.Unreachable("WAL truncate lost commits")` when a truncate index is below commit (`server/raft.go:4271-4280`) and the symmetric case at `server/raft.go:4648`.
- **Bit rot / corrupted on-disk metadata**: tests mutate bytes directly (remove a middle block, alias first/last sequence ranges, corrupt PSIM headers) then require recovery to succeed without panic (`server/filestore_test.go:10706-10731`, `server/filestore_test.go:6871`, `server/filestore_test.go:7430`).
- **Duplicate terminal events / double delivery**: CLFS (consumer ledger flow sequence) invariants are tested against duplicates both in nightly soak (`server/jetstream_cluster_long_test.go:897`) and presubmit cluster tests (`server/jetstream_cluster_3_test.go:6936`).
- **Protocol parser desync on chunked input**: `FuzzParser` splits arbitrary byte strings in half and parses sequentially across all four client kinds (`server/parser_fuzz_test.go:84-97`).
- **Goroutine and subscription leaks after churn**: goroutine-count baselines (`server/norace_2_test.go:3550-3580`) and Varz/Routez polling for zero residual subs/connections (`test/norace_test.go:530-562`).
- **Deadlocks**: prevented structurally by a documented lock-ordering contract (`locksordering.txt:1-24`), which tests then exercise by locking all raft nodes mid-test (`smGroup.lockAll/unlockAll`, `server/raft_helpers_test.go:101-124`).

### Can the suite force cancellation at every lifecycle boundary?

Largely yes for cooperative boundaries, partially for hard boundaries:

- Server lifecycle: every test tears down via `s.Shutdown(); s.WaitForShutdown()` (e.g., `server/test_test.go:344-347`), and restart helpers rebuild servers from saved options/config files (`cluster.restartServer`, `server/jetstream_helpers_test.go:1425-1443`).
- Raft node lifecycle: the `stateMachine` interface exposes `stop()`/`restart()`, and `stateAdder.restart()` recreates the WAL and replays from zero-sum to verify idempotent application (`server/raft_helpers_test.go:347-390`); the driver loop exits on both `quitCh` and `QuitC` (`server/raft_helpers_test.go:246-271`).
- Filestore flusher lifecycle: tests close the flusher channel mid-flight and verify async truncate behavior (`server/filestore_test.go:10750-10764`).
- Connection-level cancellation: `expectDisconnect` asserts EOF within 200ms (`test/test.go:465-474`).
- Gap: cancellation is always graceful. There is no in-repo harness that sends SIGKILL mid-write to prove fsync-ordering crash consistency; torn-state effects are approximated by post-hoc byte corruption instead (e.g., `server/filestore_test.go:10690-10731`). Evidence of true kill -9 style tests: none found in this repository.

### How does it prove absence of duplicate terminal events, leaks, and partial commits?

- **Duplicates**: the CLFS tests assert consumer redelivery ledgers do not advance on duplicates (`server/jetstream_cluster_3_test.go:6936`) and survive leader churn overnight (`server/jetstream_cluster_long_test.go:897`); queue-subscriber weight updates are asserted monotonic under concurrent subscription storms (`test/norace_test.go:406-433`).
- **Leaks**: three layers — goroutine count deltas after shutdown (`server/norace_2_test.go:3574-3579`), observable connection counts via `/varz` (`test/norace_test.go:532-539`), and routed-subscription counts via `/routez` (`test/norace_test.go:553-562`). Notably these are hand-rolled; there is no `goleak` dependency (checked `go.mod` and imports — none found).
- **Partial commits**: filestore recovery tests reconstruct stream state from blocks after simulated damage and compare against expected sequences/state (`server/filestore_test.go:6001`, `server/filestore_test.go:10500`), while production-side `assert.Unreachable` calls flag any flush/sync/write error reaching invariant-breaking code (`server/filestore.go:8707-8732`). The Antithesis wrappers record violations without halting, accumulating evidence across randomized runs (`internal/antithesis/test_assert.go:31-38`).

### Which techniques are proportionate for Aren Phase 1 rather than infrastructure-scale theatre?

Proportionate and transferable to Aren:

1. Build-tag quarantine of race-hostile tests with a naming convention (`test/norace_test.go:14,38-40`) — cheap, immediate, and keeps `-race` meaningful.
2. A poll-until helper with timeout-as-assertion semantics (`server/server_test.go:44-64`) replacing fixed sleeps — improves determinism diagnosis at near-zero cost.
3. Direct on-disk corruption tests for any durable store (`server/filestore_test.go:10690-10733`) — high value per line of test code.
4. Partial-write conn wrappers (`server/server_test.go:1683-1700`) — trivially reusable for any transport layer.
5. Instrumenting test assertion helpers so failures emit structured, dedupable events (`server/test_test.go:60-201`) even if Aren never runs Antithesis.
6. One seeded restart-soak test per critical subsystem (`server/jetstream_cluster_long_test.go:37-190`) gated behind a nightly tag — bounded cost, catches recovery regressions.

Likely theatre for Phase 1: supercluster topologies (`server/jetstream_super_cluster_test.go`), six-way cipher/compression permutation matrices (`server/filestore_test.go:56-72`), and continuous TLS handshake fuzzing — valuable at NATS's scale but not required to falsify Aren's Phase 1 invariants.

## Architectural Decisions

1. **Complementary race/no-race suites instead of one suite.** Rather than weakening tests to survive `-race`, nats-server keeps two tiers: race-sensitive correctness tests run sanitized; throughput/backpressure tests run unsanitized under `TestNoRace*` conventions enforced by build tags (`scripts/runTestsOnTravis.sh:19-33`, `test/norace_test.go:14`). This accepts double maintenance cost in exchange for full coverage of both timing regimes.

2. **Assertions live in production code, not just tests.** Invariant violations (lost commits, unexpected I/O errors, state resets) call `assert.Unreachable`/`assert.AlwaysOrUnreachable` inline (`server/raft.go:4271`, `server/filestore.go:5458-5467`, `server/jetstream_cluster.go:1761`), compiled out unless the `enable_antithesis_sdk` tag is set (`internal/antithesis/noop.go:15-27`). The same binary therefore serves normal operation, instrumented test runs, and third-party fault-injection platforms without code divergence.

3. **Failure waits are expressed as retryable predicates with violation reporting.** `checkFor` converts "invariant not reached in N seconds" into both a test failure and an `AssertUnreachable` event with context (`server/server_test.go:57-64`); cluster-level waiters carry structured details (cluster/account/stream names) into the assertion payload (`server/jetstream_helpers_test.go:1494-1500`).

4. **Suite economics via tags + name prefixes + CI matrix.** Three orthogonal slicing mechanisms — build tags (`skip_js_tests`, `skip_store_tests`, ... at `server/filestore_test.go:14`, `server/jetstream_test.go:14`), test-name prefixes (`TestNoRace`, `TestJetStreamCluster`, `TestNRG`, `TestJWT` in `scripts/runTestsOnTravis.sh:48-130`), and per-job 30-minute timeouts (`.github/workflows/tests.yaml:106-372`) — let CI parallelize a ~3,900-test codebase into sharded ≤30-minute jobs, with `-p=1 -failfast` to keep output diagnosable.

5. **Time acceleration instead of time abstraction.** Rather than injecting a fake clock, tests globally shorten raft timers in `init()` (`server/jetstream_helpers_test.go:49-55`) and rely on real-time polling. This keeps production code free of clock seams at the cost of statistical (not guaranteed) determinism.

## Notable Patterns

- **Timeout-as-property**: every wait loop ends in an assertion event, so a flake produces machine-readable evidence ("Timeout in cluster.waitOnPeerCount" with cluster name, `server/jetstream_helpers_test.go:1466-1469`) usable for cross-run aggregation.
- **Seeded chaos**: the restart soak fixes its RNG seed (`Seed = 123456`, `server/jetstream_cluster_long_test.go:39,62`) making a 3-minute randomized failure replayable bit-for-bit; benchmarks similarly derive per-worker RNG streams from seeds (`server/jetstream_benchmark_test.go:595,722`).
- **Corpus-driven native fuzzing with domain-aware seeds**: fuzz targets seed with realistic protocol fragments and mutate along meaningful axes (client kind, TLS version, certificate byte offsets) rather than raw bytes (`server/parser_fuzz_test.go:58-82`, `server/server_fuzz_test.go:347-359`).
- **Transport substitution at two levels**: high-fidelity real TCP through `netProxy` (for latency/bandwidth realism, `server/jetstream_helpers_test.go:1970-2042`) versus fully in-memory `raftTransportHub` (for deterministic partition topology and message interception, `server/raft_transport_helpers_test.go:43-124`).
- **Permutation matrices for storage**: every filestore behavioral test automatically covers encryption × compression combinations (`server/filestore_test.go:56-72`).
- **Leak checks as explicit baselines**: `time.Sleep(1s)` to settle, snapshot `runtime.NumGoroutine()`, run workload, shut down, then `checkFor` equality (`server/norace_2_test.go:3550-3579`).

## Tradeoffs

- **Speed vs. sanitizer fidelity**: PR CI runs the heavy JS/store/raft suites *without* `-race` (the `RACE` expression yields `''` for pull_request events, `.github/workflows/tests.yaml:12`), trading race-detection latency on the highest-churn path for feedback speed; race runs occur on branch pushes and nightly-style jobs instead.
- **Realism vs. determinism**: real wall-clock polling (`checkFor`, `waitOnLeader`) exercises genuine scheduling but makes some tests load-sensitive; mitigated by generous timeouts (10–30s) and `-shuffle on` only in the nightly job (`.github/workflows/long-tests.yaml:33`), not presubmit.
- **Assertion accumulation vs. halt-on-first-failure**: Antithesis violations print but do not fail the run (`internal/antithesis/test_assert.go:31-35`), maximizing evidence per run at the cost of needing post-processing to notice violations.
- **Graceful-shutdown restarts vs. crash realism**: `Shutdown()+WaitForShutdown()+restartServer` (`server/jetstream_cluster_long_test.go:163-167`) is reproducible and portable but cannot falsify fsync/torn-page bugs that only a SIGKILL reveals; those are approximated by manual byte corruption (`server/filestore_test.go:10712-10723`).
- **Tag sprawl vs. single-command ergonomics**: the many skip tags (`skip_no_race_1_tests`, `skip_js_cluster_tests_2`, ...) make local reproduction of a CI shard straightforward (`go test ... -tags=skip_no_race_2_tests`) but the full mapping lives only in shell script conditionals (`scripts/runTestsOnTravis.sh:19-136`).

## Failure Modes / Edge Cases

- **Flaky baseline leak checks**: comparing absolute goroutine counts assumes a quiet process; the 1-second settle sleep (`server/norace_2_test.go:3551`) reduces but does not eliminate false positives on loaded machines.
- **Timing-dependent purge benchmarking**: `require_True(t, elapsed*50 > elapsed2)` compares wall-clock durations of purge operations (`server/norace_2_test.go:3543-3546`) — inherently environment-sensitive despite living in the no-race tier.
- **Skipped-but-present tests**: `TestNoRaceSlowProxy` is permanently skipped with `t.Skip()` (`test/norace_test.go:650-651`); dead adversarial coverage silently decays.
- **Fuzz corpus rot**: with no scheduled fuzzing workflow (no matches for "fuzz" under `.github/workflows/`), `FuzzParser`/`FuzzServerTLS` corpora grow only when developers opt in; regressions found by OSS-Fuzz-style services would not be caught in-repo.
- **Partial-tag test blindness**: running `go test ./server` locally compiles everything except tag-excluded files; a developer may believe a `TestNoRace*` test passes when it was simply not compiled (it requires an explicit `-run=TestNoRace` untagged build).
- **Assertion identifier collisions**: Antithesis message strings are the unique event IDs and include the test name but not the line (`internal/antithesis/test_assert.go:65-67`), so the same message fired from two lines of one test deduplicates incorrectly — acknowledged in-source.

## Future Considerations

- Wire existing fuzz targets into a scheduled workflow (even 2–5 minutes of `go test -fuzz=FuzzParser -fuzz=FuzzServerTLS` nightly) so corpus evolution does not depend on developer initiative; targets already exist at `server/parser_fuzz_test.go:57`, `server/server_fuzz_test.go:306`, `server/subject_fuzz_test.go:21`.
- Introduce a minimal clock seam (interface or package-level `now func() time.Time`) for TTL/expiry paths currently validated only by real sleeps (e.g., message-TTL recovery tests, `server/filestore_test.go:9427-9546`), enabling deterministic expiry-edge tests.
- Add a SIGKILL crash-consistency harness for the filestore (spawn subprocess, kill during flush window, reopen and verify) to close the gap between graceful-restart soaks and the byte-corruption approximation.
- Adopt a library-grade leak checker (goleak-style deferred check) alongside hand-rolled baselines to catch leaks in tests that forget explicit checks.
- Re-evaluate the PR-path race gap (`tests.yaml:12`): either run `-race` on a subset shard for PRs or document where sanitized PR coverage actually occurs (possibly the external Buildkite hinted at by `skipIfBuildkite`, `server/test_test.go:394-398`).

## Questions / Gaps

- **Where does race-sanitized coverage for PRs actually run?** The in-repo expression disables `$RACE` for pull_request events (`.github/workflows/tests.yaml:12`), and `skipIfBuildkite` implies an external pipeline (`server/test_test.go:394-398`), but the Buildkite configuration is outside this source tree. Searched: all files under `.github/workflows/`, `scripts/`. No in-repo answer.
- **Is any continuous fuzzing performed upstream?** No workflow references fuzz targets. Searched: `.github/workflows/*.yaml` and `*.yml` for "fuzz"/"-fuzz". If upstream relies on an external service, it is invisible here.
- **How are Antithesis assertion results consumed?** The SDK emits events but nothing in-repo aggregates or gates on them (`internal/antithesis/test_assert.go:31-38`); the consumption side (dashboards, thresholds) is not represented in this source.
- **Aren acceptance-matrix alignment**: the dimension cites `../../../../Aren/docs/phase-1-prd/03-validation-and-acceptance.md` as context, but source-isolation rules forbid reading outside `sources/nats-server/`. Mapping of these techniques onto Aren's specific Phase 1 acceptance criteria is therefore left to the synthesis stage with the dimension text as the only Aren-side input.

---

Generated by `01.03 Adversarial Concurrency and Failure Verification` against `nats-server`.
