# Source Analysis: nomad

## 01.03 Adversarial Concurrency and Failure Verification

### Source Info

| Field | Value |
|-------|-------|
| Name | nomad |
| Path | `studies/aren-go-runtime-study/sources/nomad` |
| Language / Stack | Go 1.26 (module `github.com/hashicorp/nomad`), bbolt, hashicorp/raft, yamux, go-plugin; testing via `testing/quick`, `pgregory.net/rapid` v1.3.0, native Go fuzzing, `go.uber.org/goleak` v1.3.0 |
| Analyzed | 2026-08-26 |

## Summary

Nomad proves invariants with five complementary techniques rather than one big chaos framework.

First, the scheduler decision core is tested with property-based testing. The reconciler suites generate random cluster and job states with the `rapid` library and assert safety properties ("something bad never happens") on every run (`scheduler/reconciler/reconcile_cluster_prop_test.go:88-187`, `scheduler/reconciler/reconcile_node_prop_test.go:66-90`). A dedicated GitHub workflow runs them at `-rapid.checks=100000` on every PR that touches `scheduler/**` or `nomad/structs/**` (`.github/workflows/test-scheduler-prop.yml:27-29`) and uploads failing counterexamples as artifacts (`.github/workflows/test-scheduler-prop.yml:40-43`).

Second, everything that talks to the outside world has a test seam built for failure injection: an always-failing state DB (`client/state/db_error.go:25-28`), a mock driver whose task config carries explicit failure knobs like `start_error`, `start_block_for`, and `kill_after` (`drivers/mock/driver.go:240-254`), a fake Consul whose mutex can be held to block deregistration mid-shutdown (`command/agent/consul/unit_test.go:408-433`), and a scripted process stand-in that re-executes the test binary itself (`helper/testtask/testtask.go:19-66`). Tests kill real subprocesses with SIGKILL and even SIGSTOP/SIGCONT races to simulate a helper dying mid-start (`client/allocrunner/taskrunner/logmon_hook_unix_test.go:62-76`, `client/allocrunner/taskrunner/logmon_hook_unix_test.go:133-159`).

Third, crash recovery is exercised by round-tripping real snapshots. FSM snapshot/restore tests cover every major table (`nomad/fsm_test.go:2308-2600+`), client restart tests shut down a live client, mutate server state while it is down, plant a corrupted allocation in the persisted DB, then bring the client back up (`client/client_test.go:813-831`). Multi-server Raft failover runs in-process with election timeouts tightened to 50ms so killing a leader takes seconds instead of minutes (`nomad/testing.go:96-99`, `nomad/leader_test.go:97-132`).

Fourth, eventual consistency is handled by a shared polling toolkit with CI-aware budgets (`testutil/wait.go:24-51`, `testutil/wait.go:107-113`) plus negative assertions that prove nothing bad happens for a duration (`testutil/wait.go:92-103`).

Fifth, expensive work is separated from presubmit: package groups shard the unit suite (`ci/test-core.json:2-51`), slow tests are env-gated (`ci/slow.go:12-18`), root-required tests skip locally (`ci/skip_non_root.go:15-23`), e2e suites are gated behind env vars and separate workflows (`e2e/consulcompat/consulcompat_test.go:15-20`, `.github/workflows/test-e2e.yml:62-64`), and flaky failures are retried up to three times via gotestsum (`GNUmakefile:336`, `GNUmakefile:341`).

The main weakness: the core unit suite runs without the race detector. `-race` only applies to the `api`/`jobspec2` submodules and to e2e/integration targets (`GNUmakefile:355`, `GNUmakefile:363`, `GNUmakefile:373-395`), so data races in `nomad/`, `client/`, or `command/` would pass the main CI workflow (`.github/workflows/test-core.yaml:154-164`).

## Rating

**8 / 10**

Rationale against the dimension's purpose (repeatable techniques for races, panics, partial writes, restarts, resource exhaustion, scheduling variation):

- Property-based scheduler verification with CI-enforced scale and counterexample capture is exactly the kind of falsification the dimension asks for, and the generator design deliberately excludes irrelevant dimensions from shrinking (`scheduler/reconciler/reconcile_cluster_prop_test.go:614-617`). This is rare even in mature Go projects.
- Every infrastructure boundary named in the dimension (clocks, transports, storage, schedulers, process runners, randomness) has an inspectable seam with evidence.
- Deductions: no race detector on the core unit path; goroutine-leak detection exists in exactly one package; a single fuzz target; the counterexample upload path (`./scheduler/reconciler/testdata`) does not exist in the repo, so the persistence loop from CI failure to committed regression seed appears unfinished.

## Evidence Collected

| Area | Evidence | File:Line |
|------|----------|-----------|
| Property tests (cluster reconciler) | `TestAllocReconciler_PropTest` draws random reconcilers for batch/service jobs; `sharedSafetyProperties` asserts placements never exceed group count, canaries never exceed budget, no placements on stop/scale-to-zero | `scheduler/reconciler/reconcile_cluster_prop_test.go:23`, `scheduler/reconciler/reconcile_cluster_prop_test.go:131-163` |
| Property tests (reconnect dedup) | `TestAllocReconciler_ReconnectingProps` fails if a reconnecting alloc is both kept and stopped, or neither | `scheduler/reconciler/reconcile_cluster_prop_test.go:647-735` |
| Property test CI gating | Dedicated workflow runs `-rapid.checks=100000 -run PropTest`; uploads failures to artifacts | `.github/workflows/test-scheduler-prop.yml:27-29`, `.github/workflows/test-scheduler-prop.yml:40-43` |
| Generator hygiene | IDs generated outside rapid generators so shrinking focuses on meaningful dimensions | `scheduler/reconciler/reconcile_cluster_prop_test.go:614-617` |
| Node reconciler props | System/sysbatch jobs: place ≤ count × ready nodes; migrate/stop/update ≤ existing allocs | `scheduler/reconciler/reconcile_node_prop_test.go:66-90` |
| Fuzzing | `FuzzTaskScheduleCron` seeds epoch, y2k, y2038, and a last-day-of-month regression into the corpus | `nomad/structs/task_sched_test.go:15-41` |
| Metamorphic/differential tests | `quick.Check` asserts escape reader output equals a naive reference under arbitrary read chunking; second property equates chunked reads to single reads | `helper/escapingio/reader_test.go:248-256`, `helper/escapingio/reader_test.go:321-329`, `helper/escapingio/reader_test.go:299-302` |
| Storage fault injection | `ErrDB` implements `StateDB` where every method returns an error | `client/state/db_error.go:20-28`; used at `client/hostvolumemanager/host_volumes_test.go:28-78` |
| Corrupted-state injection | Test plants a corrupt allocation in client state DB before shutdown, expects removal during restore | `client/client_test.go:813-818` |
| Client crash-recovery matrix | `TestClient_SaveRestoreState`: stop one job before shutdown, mutate two more while down, restart, assert per-job convergence | `client/client_test.go:702-831` |
| Process runner fault knobs | Mock driver config: `StartErr`, `StartErrRecoverable`, `StartBlockFor`, `KillAfter`, `PluginExitAfter`; driver-level `ShutdownPeriodicAfter` to kill fingerprinting mid-run | `drivers/mock/driver.go:240-254`, `drivers/mock/driver.go:176-184` |
| Subprocess crash injection | Logmon killed with SIGKILL, reattach recovery asserted; second test SIGSTOPs logmon then kills it while `Prestart` is blocked in `Start` | `client/allocrunner/taskrunner/logmon_hook_unix_test.go:29-90`, `client/allocrunner/taskrunner/logmon_hook_unix_test.go:94-175` |
| Blocked-dependency shutdown | Fake Consul mutex held during agent shutdown; asserts shutdown waits 200ms–1s for deregistration and completes after unlock | `command/agent/consul/unit_test.go:377-443` |
| Restore blocking semantics (#1795) | Restored dead tasks stay pending until server contact, Update, or Kill arrives | `client/allocrunner/taskrunner/task_runner_linux_test.go:529-599`, `client/allocrunner/taskrunner/task_runner_linux_test.go:601-669` |
| FSM snapshot restore | Per-table restore tests (`Nodes`, `Jobs`, `Evals`, `Allocs`, `Indexes`, ...) compare restored store to original; pre-TTL upgrade path canonicalizes on apply | `nomad/fsm_test.go:2308`, `nomad/fsm_test.go:2449`, `nomad/fsm_test.go:2477` |
| State-store restore | Index table restoration tested separately (last-index per table) | `nomad/state/state_store_restore_test.go:201-217` |
| Leader failover | 3-server in-process Raft clusters; `TestLeader_LeftLeader` kills the leader and waits for a successor; heartbeat timers verified cleared on old leader and restored on new leader | `nomad/leader_test.go:97-132`, `nomad/heartbeat_test.go:226-300` |
| Fast elections harness | Test config tightens raft election/heartbeat/lease to 50ms and autopilot intervals to 50–100ms | `nomad/testing.go:96-99`, `nomad/testing.go:106-108` |
| Startup retry under port collisions | Server creation retries 10 times with jittered sleep and fresh ports; cleanup enforces a 1-minute shutdown deadline | `nomad/testing.go:154-196`, `nomad/testing.go:170-177` |
| Clock seam | `libtime.Clock` injected into checks client and user cache; tests override with `libtimetest.NewClockMock` | `client/serviceregistration/checks/client.go:46-52`, `client/serviceregistration/checks/client_test.go:78-80`, `helper/users/cache_test.go:66-69`, dep at `go.mod:149` |
| Process runner seam | `testtask` re-executes the test binary with `TEST_TASK=execute` and a script of sleep/echo/write/pgrp ops; no external binaries needed | `helper/testtask/testtask.go:19-66` |
| Real dependency binaries | `NewTestVaultFromPath` / `NewTestVaultDelayed` run real Vault with delayed startup for renewal-failure paths | `testutil/vault.go:48-53`, `testutil/vault.go:135-198` |
| Eventual-consistency toolkit | `Wait` polls 500×10ms ×4 in CI; `AssertUntil` asserts absence of an event over a window; `WaitForVotingMembers` blocks until autopilot promotes all peers | `testutil/wait.go:24-51`, `testutil/wait.go:92-103`, `testutil/wait.go:107-113`, `testutil/wait.go:243-269` |
| Goroutine leak detection | `goleak.VerifyNone(t)` guards rotator lifecycle tests | `go.mod:140`, `client/logmon/logging/rotator_test.go:16-249` |
| Concurrency monotonicity | Broadcaster test fails if any listener observes a decreasing `AllocModifyIndex`, including while listeners close concurrently | `client/structs/broadcaster_test.go:93-166` |
| Eval exactly-once-ish delivery | Ack/Nack/delivery-limit suite; per-job serialization invariant (`TotalUnacked` per job capped at 1) | `nomad/eval_broker_test.go:400-459`, `nomad/eval_broker_test.go:813-889`, `nomad/eval_broker_test.go:950-1083` |
| Race detector placement | Core unit target `test-nomad` has no `-race`; submodule/e2e/integration targets do | `GNUmakefile:334-346`, `GNUmakefile:348-357`, `GNUmakefile:359-399` |
| Suite sharding and gating | Package groups (`nomad`, `client`, `command`, `drivers`, `quick`) in CI matrix; `SkipSlow` gated by `NOMAD_SLOW_TEST=0` in CI; root-gate; port allocator | `ci/test-core.json:2-51`, `.github/workflows/test-core.yaml:124-164`, `.github/workflows/test-core.yaml:46`, `ci/slow.go:12-18`, `ci/skip_non_root.go:15-23`, `ci/ports.go:18-23` |
| Coverage enforcement | `make missing` fails for packages absent from `ci/test-core.json` | `GNUmakefile:483-486` |
| Flake mitigation | gotestsum `--rerun-fails=3` in CI, disabled for human-facing `make test` | `GNUmakefile:335-336`, `GNUmakefile:341`, `GNUmakefile:493-496` |
| Historical regressions as tests | Restore-blocking behavior pinned with issue reference #1795; leak fixes recorded in changelog (GH-11741, GH-20348) | `client/allocrunner/taskrunner/task_runner_linux_test.go:529-530`, `.changelog/11741.txt:2`, `.changelog/20348.txt:2` |

## Answers to Dimension Questions

### Which concurrency failures would ordinary coverage miss?

Line coverage would miss all four classes Nomad tests explicitly:

1. Lost or reordered updates. `client/structs/broadcaster_test.go:115-121` fails if any listener ever observes a decreasing `AllocModifyIndex` while listeners are being added and closed concurrently. A happy-path send/receive test passes even if the broadcaster drops intermediate values.
2. Duplicate or contradictory terminal decisions. The reconciler property suites fail when a reconnecting alloc ends up both kept and stopped (`scheduler/reconciler/reconcile_cluster_prop_test.go:706-715`) or when placements exceed group counts across thousands of generated states.
3. Double-delivery of work. The eval broker suite pins ack/nack semantics, nack timeout resets, pause/resume of nack timers, and the delivery limit after which re-delivery stops (`nomad/eval_broker_test.go:813-1083`).
4. Goroutine leaks. `goleak.VerifyNone` catches leaked rotator goroutines (`client/logmon/logging/rotator_test.go:24`). Note this technique is applied in only one package, so leaks elsewhere surface as changelog entries (GH-11741, GH-20348) rather than CI failures.

Data races specifically: ordinary coverage plus no `-race` means races in the core modules are invisible to the primary CI workflow until integration or e2e runs, which do use `-race` (`GNUmakefile:371-388`).

### Can the suite force cancellation at every lifecycle boundary?

For the client task lifecycle, mostly yes, and each boundary has a named test:

- Restore boundary: restored tasks block until servers are contacted, Update arrives, or Kill arrives (`client/allocrunner/taskrunner/task_runner_linux_test.go:529-599`); system tasks skip the wait (`task_runner_linux_test.go:601-669`).
- Prestart boundary: logmon crash before start (`logmon_hook_unix_test.go:29`) and crash during start (`logmon_hook_unix_test.go:94`).
- Restart policy boundary: tracker modes for delay, fail, no-restart-on-success, zero attempts, killed, externally triggered restart, recoverable vs unrecoverable start errors (`client/allocrunner/taskrunner/restarts/restarts_test.go:39-204`).
- Shutdown ordering boundary: service job with prestart sidecars and poststop hooks asserts kill order (`client/allocrunner/alloc_runner_test.go:2456-2459`) and shutdown delays (`client/allocrunner/alloc_runner_test.go:1050`).
- Agent shutdown against a stuck dependency: deregistration must finish before shutdown returns (`command/agent/consul/unit_test.go:413-443`).
- Script checks and tasklets mid-flight: `client/allocrunner/taskrunner/script_check_hook_test.go:163`, `client/allocrunner/taskrunner/tasklet_test.go:135`.

No evidence found for systematic cancellation coverage of server-side long-running loops (blocking queries, event broker subscribers) at each await point; I searched `nomad/*_test.go` for cancellation-at-boundary patterns and found endpoint-level tests but not a per-boundary cancellation matrix.

### How does it prove absence of duplicate terminal events, leaks, and partial commits?

- Duplicate terminal events: the reconciler properties bound every outcome category by input counts (place/canary/destructive/migrate each ≤ expected totals, `scheduler/reconciler/reconcile_cluster_prop_test.go:149-184`), and reconnect decisions are proven total and exclusive (`reconcile_cluster_prop_test.go:706-733`). The eval broker proves one unacked eval per job at a time (`nomad/eval_broker_test.go:442-459`). These are structural proofs over generated inputs, not spot checks.
- Leaks: goleak in the log rotator only (`client/logmon/logging/rotator_test.go:24`). Otherwise leak detection is reactive, visible from changelog entries GH-11741 and GH-20348 (`.changelog/11741.txt:2`, `.changelog/20348.txt:2`).
- Partial commits: server state writes go through bbolt transactions wrapped by `boltdd`, which de-duplicates writes and msgpack-encodes values (`helper/boltdd/boltdd.go:4-6`, `helper/boltdd/boltdd.go:43-50`), and correctness is checked by snapshot→restore round-trip equality per table (`nomad/fsm_test.go:2308+`). Client-side, `TestClient_SaveRestoreState` persists, mutates upstream, restores, and prunes a planted corrupt entry (`client/client_test.go:702-831`). For stream framing, the escaping-io metamorphic tests prove output equivalence under arbitrary read chunk boundaries (`helper/escapingio/reader_test.go:248-329`), which is the "partial read" analogue of a partial write.

### Which techniques are proportionate for Aren Phase 1 rather than infrastructure-scale theatre?

Proportionate, evidenced by Nomad getting high value per line of machinery:

1. Property tests over pure decision cores with explicit safety properties and a dedicated CI job (`scheduler/reconciler/reconcile_cluster_prop_test.go:23`, `.github/workflows/test-scheduler-prop.yml:27-29`). No cluster needed; it tests a function.
2. Always-fail fakes at storage seams (`client/state/db_error.go:20-28`). Trivial to write, catches error-path regressions that mocks-with-defaults hide.
3. Self-executing process stand-ins (`helper/testtask/testtask.go:54-66`) to avoid shipping fixture binaries.
4. Native fuzzing for parsers with regression seeds committed to the corpus (`nomad/structs/task_sched_test.go:23-33`).
5. In-process Raft with tightened timing constants (`nomad/testing.go:96-99`) so failover tests run in seconds.
6. A shared wait toolkit with a CI multiplier and `AssertUntil` for negative invariants (`testutil/wait.go:92-113`).

Not proportionate for Phase 1: nightly multi-region chaos fleets, custom fault-injection DSLs, or soak farms. Nomad itself gates its heavyweight scenarios behind env vars and separate scheduled workflows (`e2e/consulcompat/consulcompat_test.go:15-20`, `enos/`), and its own changelog shows the highest-cost bugs were caught by cheap techniques.

## Architectural Decisions

1. **In-process test servers instead of black-box binaries for most suites.** `nomad.TestServer` constructs a full server with real Raft, memdb, and mocked Consul inside the test process (`nomad/testing.go:58-138`). The exec-based `testutil.NewTestServer` that spawns a real agent binary exists (`testutil/server.go:166-266`) but skips itself when no binary is present (`testutil/server.go:167-170`). Decision consequence: failover and RPC tests are fast enough to run per-PR.
2. **Failure knobs live in production-shaped interfaces, not parallel test APIs.** The mock driver's failure knobs are codec-tagged fields of its normal `TaskConfig` (`drivers/mock/driver.go:230-266`), and `ErrDB` satisfies the same `StateDB` interface used in production (`client/state/db_error.go:20`). Tests exercise the exact wiring production uses.
3. **Property tests assert safety, not liveness.** Both reconciler prop suites document their assertions as "something bad never happens" (`scheduler/reconciler/reconcile_cluster_prop_test.go:88-89`, `scheduler/reconciler/reconcile_node_prop_test.go:66-67`), leaving liveness to example-based tests.
4. **Determinism engineered at the edges, polling tolerated in the middle.** Pure logic gets generators and mock clocks; anything crossing a goroutine or network boundary gets bounded polling with generous CI multipliers (`testutil/wait.go:107-118`) and retry-on-flake (`GNUmakefile:336`).
5. **Expensive verification is opt-in and path-triggered.** Scheduler property tests trigger only on `scheduler/**`/`nomad/structs/**` diffs (`.github/workflows/test-scheduler-prop.yml:3-13`); e2e and compat suites run in their own workflows with env gates (`.github/workflows/test-e2e.yml:62-100`).

## Notable Patterns

- **Counterexample-oriented generator design.** IDs and names come from a plain counter instead of rapid generators so that shrinking never wastes budget on renaming (`scheduler/reconciler/reconcile_cluster_prop_test.go:614-617`). Worth copying verbatim.
- **Biased generators for realistic distributions.** `weightedBool(30)` makes stops/draining uncommon but ever-present (`scheduler/reconciler/reconcile_cluster_prop_test.go:468-473`), and node status sampling repeats "ready" three times out of four to mirror real clusters (`scheduler/reconciler/reconcile_cluster_prop_test.go:423-428`).
- **Metamorphic reference implementations.** A naive string scanner defines correct behavior; the optimized reader must match it byte-for-byte under arbitrary chunking (`helper/escapingio/reader_test.go:290-319`).
- **Time-travel corpora in fuzz seeds.** Epoch, y2k, 2024 leap-year edge, and y2038 timestamps are seeded explicitly (`nomad/structs/task_sched_test.go:25-33`).
- **Negative assertions.** `AssertUntil` inverts the usual poll: the test fails if the event ever happens within the window (`testutil/wait.go:92-103`), used to prove things like "no unexpected deregistration."
- **Env-conditioned test economics.** `ci.Parallel` disables parallelism in CI for throughput (`ci/slow.go:22-30`), `TestMultiplier` stretches budgets 4x on CI runners (`testutil/wait.go:107-113`), and root-only tests skip gracefully off-CI (`ci/skip_non_root.go:15-23`).
- **Regression seeds as first-class corpus entries.** The cron fuzzer seeds the exact day-of-month case that once failed (`nomad/structs/task_sched_test.go:32-33`).

## Tradeoffs

- **Speed vs race coverage.** Running the ~5-group unit matrix without `-race` keeps the 30-minute job budget (`.github/workflows/test-core.yaml:126`, `GNUmakefile:343`) but leaves core-module races undetected in presubmit. The `-race` runs exist only for small modules and slower integration targets (`GNUmakefile:348-399`).
- **Timing-window assertions can flake.** `TestConsul_ShutdownBlocked` bounds shutdown duration between 200ms and 1s on wall-clock time (`command/agent/consul/unit_test.go:427-438`). On an overloaded runner this is a false failure; gotestsum retries mask it rather than fix it.
- **Polling-based determinism trades rigor for stability.** 10ms sleep polls (`testutil/wait.go:29`) can mask ordering violations that a synchronized harness would catch, and the broadcaster test's 10ms negative waits (`client/structs/broadcaster_test.go:79-83`) assume the goroutine has not merely been slow.
- **Retry-on-flake hides signal.** `--rerun-fails=3` (`GNUmakefile:341`) improves merge throughput at the cost of under-reporting genuine intermittent failures; Nomad acknowledges this by disabling retries for local runs (`GNUmakefile:493-496`).
- **Real-binary integration tests skip silently.** `NewTestServer` skips when no compiled agent is found (`testutil/server.go:167-170`), so a misconfigured environment quietly reduces coverage instead of failing.

## Failure Modes / Edge Cases

Covered by dedicated tests:

- Helper process dies between client restarts (reattach recovery): `client/allocrunner/taskrunner/logmon_hook_unix_test.go:29-90`.
- Helper process dies while being attached, using SIGSTOP/SIGCONT to widen the race deterministically: `client/allocrunner/taskrunner/logmon_hook_unix_test.go:133-175`.
- Dependency hangs during shutdown: `command/agent/consul/unit_test.go:377-443`.
- Persisted state contains garbage from a previous version or a crash: corrupt alloc pruning at `client/client_test.go:813-818`; pre-TTL snapshot canonicalization at `nomad/fsm_test.go:2344-2371`; state DB schema upgrades at `client/state/db_test.go:542`.
- Leader loss with in-flight timers: heartbeat timer transfer verified across failover (`nomad/heartbeat_test.go:267-300`).
- Driver stops working mid-deployment (fingerprinter shutdown knob): `drivers/mock/driver.go:176-184`.
- Start errors classified recoverable vs fatal feeding the restart tracker: `client/allocrunner/taskrunner/restarts/restarts_test.go:160-204`.
- Reads split at arbitrary byte boundaries (escape sequences spanning chunks): `helper/escapingio/reader_test.go:115-129`, `helper/escapingio/reader_test.go:248-329`.

Not covered (no evidence found): OOM/disk-full resource exhaustion simulation; partial disk writes beneath bbolt; fuzzing of RPC codec inputs beyond cron strings; server-side blocking-query cancellation matrices. I searched for `Fuzz`, `fault`, `diskfull`, and exhaustion-style fixtures across the module and found none in these areas.

## Future Considerations

For Aren, adopt in this order of value per effort:

1. Property-test pure decision cores with safety properties and counterexample capture, copying the idGenerator pattern (`scheduler/reconciler/reconcile_cluster_prop_test.go:614-617`) and biased-weight helpers (`reconcile_cluster_prop_test.go:468-473`).
2. Add an always-fail implementation of every storage interface and wire it into restore-path tests (`client/state/db_error.go:20-28` pattern).
3. Put the race detector on the primary unit-test lane, even if it means splitting the matrix further; Nomad's omission here is its clearest gap given that its concurrency-heavy code lives exactly in the un-raced packages.
4. Extend goleak beyond one package, starting with anything owning goroutines per allocation (broadcaster, watchers, brokers).
5. Seed native fuzz targets for every external-input parser, committing regression cases to the corpus (`nomad/structs/task_sched_test.go:23-33` pattern).
6. Tighten in-process consensus timings centrally so failover tests stay seconds-long (`nomad/testing.go:96-99` pattern).
7. Skip the theatre: no fault-injection frameworks, no nightly chaos fleet until the above exist and stay green.

## Questions / Gaps

- The property-test workflow uploads failures from `./scheduler/reconciler/testdata` (`.github/workflows/test-scheduler-prop.yml:43`), but that directory does not exist in the repo and no code writes to it. Either the replay loop was removed or never landed. Reproducing a CI property failure therefore relies on rapid's printed seed, not a committed artifact.
- Does any scheduled (non-PR) job run the core unit suite with `-race`? No evidence found in `.github/workflows/`; I searched all workflow files for `race` and found matches only in GNUmakefile targets invoked by e2e/integration jobs.
- The blocked-evals randomized generator (`nomad/blocked_evals_stats_test.go:308-380`) uses `*rand.Rand` without an obvious seed-persistence mechanism. Whether failures there are reproducible depends on the enclosing `quick.Config` seeding, which is not visible in the file.
- Root-required tests skip unless `CI=true` (`ci/skip_non_root.go:15-23`). What enforces they actually ran in CI (rather than being silently skipped due to a missing capability)? No evidence found of an assertion that the skip did not occur.
- Leak detection covers only `client/logmon/logging`. Given changelog evidence of leaks elsewhere (GH-11741, GH-20348), the gap between known historical leaks and current goleak coverage is unquantified.

---

Generated by `01.03-adversarial-concurrency-and-failure-verification` against `nomad`.
