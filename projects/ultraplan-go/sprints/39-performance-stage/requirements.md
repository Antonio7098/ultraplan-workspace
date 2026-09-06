# Sprint 39: Requirements-Driven Performance Stage

## Sprint Goal

Add an optional `performance` verification phase after execute and before Conformance Review. Projects must enable it explicitly. For an enabled project, the phase reads every performance target from the current sprint's validated `requirements.md`, establishes frozen benchmarks and a qualified baseline, profiles misses, tries bounded implementation changes in isolation, and publishes product-derived target verdicts without weakening correctness or the performance contract.

## Required Outputs

All implementation paths in this section are relative to the `ultraplan-go` implementation repository root. Runtime artifact paths are relative to the selected sprint directory.

### Project policy and requirements contract

- `internal/project/domain.go`: typed project performance policy with `disabled` and `enabled` modes. Missing policy means `disabled`.
- `internal/project/index.go`: strict parsing and validation of the optional project-index section below. The section controls activation only; it cannot contain target values.
- `internal/project/index_test.go`: missing, disabled, enabled, duplicate, unknown-field, invalid-mode, and attempted-target-override coverage.
- `internal/sprint/index.go`: integrate strict parsing and validation of the exact `Performance Targets` table below with the existing requirements validator.
- `internal/sprint/index_test.go`: table shape, fenced-example exclusion, cell validation, duplicate IDs, comparator and basis semantics, unit restrictions, gate behavior, sample bounds, project-policy disagreement, and stable normalization tests.
- `internal/workspace/init.go`, `internal/workspace/scaffold/templates/requirements.md`, and `internal/workspace/scaffold/prompts/create-requirements.md`: optional target-table guidance with no invented numeric thresholds.
- `internal/workspace/workspace_test.go`: exact generated section guidance and absence of fabricated targets.

The optional project-index contract is:

```markdown
## Performance Policy

- **Mode:** disabled
```

`Mode` is exactly `disabled` or `enabled`. The section may later gain lower-only operational limits through a separate governed requirement, but Sprint 39 must reject scenario, metric, comparator, value, unit, gate, sample, and basis fields outside sprint requirements.

An enabled sprint uses this exact requirements section and header:

```markdown
## Performance Targets

| ID | Scenario | Metric | Comparator | Value | Unit | Gate | Samples | Basis |
| --- | --- | --- | --- | ---: | --- | --- | ---: | --- |
| PERF-001 | `<bounded scenario>` | `<stable metric>` | `<=` | `<number>` | `ns/op` | `required` | `10` | `absolute` |
```

The v1 cell contract is:

- `ID`: unique `PERF-` identifier using uppercase ASCII letters, digits, and hyphens.
- `Scenario`: non-empty bounded text that identifies one runnable benchmark or measurement scenario.
- `Metric`: non-empty stable metric name. It is descriptive, not executable input.
- `Comparator`: exactly `<=`, `>=`, or `baseline`.
- `Value`: a finite decimal for `<=` and `>=`; exactly `-` for `baseline`.
- `Unit`: one supported versioned unit. Sprint 39 must support `ns/op`, `B/op`, `allocs/op`, `ms`, `MiB`, `ops/s`, `tasks/run`, and `%`.
- `Gate`: exactly `required` or `report`. `baseline` requires `report`.
- `Samples`: an integer within product-owned minimum and maximum bounds. `baseline` still requires enough samples for qualification.
- `Basis`: exactly `absolute`, `baseline`, or `current`. `<=` and `>=` use `absolute` or `baseline`; `baseline` uses `current`.
- `Basis=absolute`: compare the qualified candidate measurement directly with `Value` in `Unit`.
- `Basis=baseline`: `Unit` must be `%`; compare `100 * (candidate - baseline) / baseline` with `Value`. A latency reduction target can use `<= -10`; a throughput increase target can use `>= 15`. The benchmark manifest freezes the raw measurement unit, which must match between baseline and candidate. A zero or non-finite baseline is inconclusive.
- `Comparator=baseline`: record and report the qualified current value without imposing a candidate threshold.

The table is the sole target authority. Project policy, workspace configuration, environment variables, CLI flags, runtime prompts, benchmark output, prior measurements, and generated files must not add, remove, replace, or relax target rows or values.

### Performance domain, orchestration, and persistence

- `internal/sprint/performance.go`: `VerificationPhasePerformance` service for admission, target freezing, benchmark preparation, baseline measurement, profiling, bounded optimization, final measurement, correctness verification, publication, and recovery.
- `internal/sprint/performance_types.go`: policy, target, target packet, benchmark manifest, measurement, qualification, profile, hypothesis, proposal, comparison, freshness, target outcome, run outcome, and limit types.
- `internal/sprint/performance_targets.go`: normalized packet creation, requirements digest binding, stable ordering, duplicate detection, and target comparison semantics.
- `internal/sprint/performance_benchmark.go`: benchmark discovery and authoring protocol, allowed benchmark paths, manifest validation, benchmark promotion, identity freezing, and target coverage checks.
- `internal/sprint/performance_measure.go`: bounded warmups, repeated samples, parser dispatch, environment capture, variance qualification, aggregation, and comparison facts.
- `internal/sprint/performance_optimize.go`: one-miss-at-a-time profiling, bounded hypothesis and proposal attempts, isolated correctness and measurement, product-owned promotion, and convergence rules.
- `internal/sprint/performance_state.go`: versioned private persistence, immutable attempt evidence, atomic current-state publication, digest-bound pointers, writer fencing, retention, and bounded flow-state projection.
- Focused `_test.go` files for every component above, including table tests, integration tests, failure injection, race tests, cancellation, recovery, target drift, benchmark drift, environment drift, noisy measurements, and mutation containment.
- `internal/sprint/domain.go`, `flow.go`, `verification_phase.go`, `state.go`, and `state_database.go`: additive performance phase membership, readiness, migration, freshness, and ordering after execute and before Conformance Review.

The service must publish this exact artifact family:

- `performance.md`: canonical bounded human-readable target table, baseline and candidate summaries, verdict, change summary, blockers, evidence pointers, and next action.
- `verification/performance-state.json`: current schema version, phase, freshness, active or terminal run correlation, target counts, bounded counters, outcome, blocker, next action, and digest-bound pointers.
- `verification/performance/attempts/<attempt-id>/target-packet.json`: immutable normalized targets, requirements digest, project-policy digest, parser version, and packet digest.
- `verification/performance/attempts/<attempt-id>/benchmark-manifest.json`: immutable target-to-benchmark mapping, benchmark paths, commands, parser, benchmark digest, and coverage result.
- `verification/performance/attempts/<attempt-id>/environment.json`: bounded machine, OS, architecture, runtime, toolchain, process policy, and noise-control identity.
- `verification/performance/attempts/<attempt-id>/baseline.json`: warmups, samples, qualification facts, aggregate measurements, and baseline verdicts.
- `verification/performance/attempts/<attempt-id>/cycles/<cycle-number>/profile.json`: bounded profile evidence and one target-linked hypothesis.
- `verification/performance/attempts/<attempt-id>/cycles/<cycle-number>/proposal.patch`: proposed benchmark or implementation patch retained before promotion.
- `verification/performance/attempts/<attempt-id>/cycles/<cycle-number>/scope.json`: before and after identities, actual paths, changed bytes, and containment result.
- `verification/performance/attempts/<attempt-id>/cycles/<cycle-number>/correctness.json`: product-owned correctness command results.
- `verification/performance/attempts/<attempt-id>/cycles/<cycle-number>/measurements.json`: candidate samples, qualification, comparison facts, and per-target outcome.
- `verification/performance/attempts/<attempt-id>/cycles/<cycle-number>/cleanup.json`: disposable-copy and process-tree cleanup result.
- `verification/performance/attempts/<attempt-id>/result.json`: immutable terminal outcome, per-target outcomes, consumed limits, final identities, unresolved misses, and recovery or escalation action.

### Process, runtime, and mutation boundaries

- Reuse `internal/platform/process` for explicit argv execution, allowlisted environment, bounded output, timeouts, cancellation, process-tree cleanup, and executable identity. Do not add shell-string execution.
- Reuse the execute-target worktree and current disposable-copy isolation. An inability to prove source identity, isolation, or cleanup is `blocked` or `cleanup_uncertain`, never success.
- Add stage-specific performance model configuration through the existing runtime model mechanism. Runtime selection cannot affect targets, commands, parsers, sample policy, correctness gates, limits, or verdict rules.
- Add safe finite defaults and immutable product maxima for benchmark authoring attempts, warmups, samples, coefficient of variation, commands, command time, output bytes, profile bytes, proposal bytes, changed files, changed bytes, optimization cycles, wall time, retained attempts, and cleanup time. Workspace and environment settings may only lower operational limits.
- Product code owns benchmark command descriptors. The runtime may propose a benchmark mapping, but it cannot provide free-form commands that are executed without parsing and validation.
- Benchmark and implementation proposals are created in disposable copies. Product code validates the patch, protected paths, current target identity, actual diff, correctness checks, and frozen benchmark behavior before promoting it to the sprint worktree.

### Shared application and interfaces

- `internal/app/sprint_usecases.go`: adapter-independent prepare, dry-run, start, status, resume, cancel, recover, result, and bounded evidence-query DTOs.
- `internal/app/operations.go`: performance operation kinds, requests, governed inputs, mutation class, stable result fields, and confirmation facts when a production promotion is possible.
- `internal/app/operation_runner.go`: one durable performance runner shared by CLI, TUI, and browser with writer-fence ownership and canonical cancellation.
- `internal/app/sprint_commands.go`: `ultraplan sprint <project> <sprint> performance`, plus `--dry-run`, `status`, `resume`, `cancel`, and `recover`; text and JSON output use the same app result.
- `internal/app/sprint_commands_test.go`, `sprint_usecases_test.go`, and `durable_operations_test.go`: parsing, disabled behavior, admission, JSON schema, exit codes, cancellation, restart, stale writers, and exactly-one-terminal-result coverage.
- `internal/tui`: policy, target, phase, progress, target verdict, blocker, next action, cancellation, and recovery views using application DTOs only.
- `internal/web`: server-rendered policy, target, preparation, progress, result, cancellation, and recovery routes using application use cases only. JavaScript remains optional progressive enhancement.
- Additive canonical fixtures under `internal/testdata/` proving CLI, JSON, TUI, and browser parity without exposing full prompts, secrets, raw provider payloads, or unbounded measurements.

### Documentation and release evidence

- `docs/architecture.md`: performance phase ownership, requirements authority, isolation, mutation, persistence, and freshness boundaries.
- `docs/cli-reference.md`: exact performance commands, states, outcomes, exit behavior, target schema, and examples.
- `docs/user-guide.md`: enabling a project, writing targets in sprint requirements, interpreting baseline and target results, cancellation, recovery, and disabling the phase.
- `docs/phase3-json-schemas.md` or its current successor: private performance records and additive public projections.
- `docs/recovery.md`: interruption, stale ownership, target or benchmark drift, partial publication, cleanup uncertainty, and restart reconciliation.
- `docs/local-web.md`: browser operation and no-JavaScript behavior.
- `docs/release-checklist.md`: disabled-flow regression, parser fixtures, correctness regression, measurement qualification, isolation, race tests, interface parity, and real-runtime dogfood.

The implementation and review must re-read these current plans at planning time:

| Input | Responsibility in this sprint |
| --- | --- |
| `docs/plans/integrated-roadmap.md` | Preserve phase ordering, product-owned authority, and gated expansion. |
| `docs/plans/post-execution-qa-and-repair-loop.md` | Reuse verification-phase, isolation, evidence, mutation, freshness, and bounded-convergence rules. |
| `docs/plans/server-shutdown-run-cancellation-contract.md` | Preserve durable acceptance, cancellation, cleanup, terminal arbitration, and restart reconciliation. |
| `docs/plans/retrieval-ready-content-plan.md` | Keep performance evidence compatible with later content identity without implementing that future abstraction. |
| `../Aren/docs/phased-roadmap.md` | Preserve Aren's capability sequence and avoid using performance work to pull future runtime capabilities forward. |
| `../Aren/docs/performance-engineering.md` | Preserve baseline-first targets, workload-shaped benchmarks, variance qualification, profile-led optimization, correctness before speed, concurrency scaling, memory/runtime-task evidence, and the rule against invented numeric budgets. |

The Aren paths are resolved from the UltraPlan implementation repository. UltraPlan's optional performance phase is an orchestration mechanism; it does not change Aren's internal phase structure or enable performance work for Aren unless Aren's project policy selects it.

## Acceptance Criteria

### AC-1: activation is explicit per project and disabled by default

- A missing `Performance Policy` section and `Mode: disabled` both disable the phase.
- Disabled projects retain the existing post-execute transition and do not parse target rows, discover benchmarks, start a runtime, execute measurement commands, create performance state, or change existing artifact bytes. Requirements validation may detect a non-empty target section only to report the policy mismatch required by AC-2.
- `Mode: enabled` inserts performance after successful execute and before Conformance Review. Later verification cannot start until the current performance run has a non-blocking outcome.
- Project-index validation rejects duplicate policy sections, unknown keys, invalid modes, and any attempt to place target values in project policy.
- Existing workspaces, project indexes, sprint state, and CLI invocations remain compatible without migration work by the user.

### AC-2: sprint requirements are the sole target authority

- An enabled sprint fails requirements validation if `## Performance Targets` is absent, duplicated, empty, malformed, or contains no usable row.
- A disabled project with a target table fails validation with guidance to enable the project or remove the table. Targets are never silently ignored.
- Only an actual level-two heading and its immediately associated table count. Text inside fenced code, inline code, blockquotes, comments, or another section is not a target declaration.
- Every accepted row satisfies the exact v1 schema in Required Outputs. Unknown columns, missing columns, reordered columns, duplicate IDs, unsupported units, invalid comparator and basis combinations, non-finite values, and sample counts outside bounds fail validation.
- Target normalization is deterministic. Identical requirements bytes and policy produce the same ordered packet and digest.
- The target packet includes every row exactly once and binds the source requirements digest. Any requirements change makes prior preparation and measurement stale.
- An operational limit lower than a requirements-owned sample count blocks admission with both values reported. It cannot silently reduce the target's sample count.
- Configuration, flags, runtime output, benchmark files, stored baselines, and prior artifacts cannot override the packet. A model suggestion can only become a target through a separately governed edit to requirements before the attempt starts.
- Requirements generation may copy or derive targets from selected governed sources. When no numeric threshold is governed, it may emit a `baseline` report row or leave the section for human completion; it must not invent a numeric threshold.

### AC-3: admission and dry-run are deterministic and runtime-free

- Admission resolves the enabled project policy, valid current target packet, successful execute state, sprint worktree identity, execute evidence, correctness commands, isolation support, process policy, parser versions, and effective lower-only limits.
- Admission rejects missing or stale execute evidence, unresolved target worktree, an active conflicting writer, target overlap with uncommitted unrelated changes, missing correctness commands, unsupported measurement units, or unavailable isolation.
- `--dry-run` reports policy, target rows, expected benchmark coverage, protected roots, correctness commands, effective limits, missing prerequisites, and next action. It starts no runtime or command and mutates no artifact.
- Repeating prepare or dry-run against unchanged inputs is idempotent and produces the same facts.

### AC-4: benchmarks are isolated, covered, and frozen before baseline

- Every target maps to exactly one validated measurement descriptor. A descriptor names its scenario, parser, explicit argv, working directory, bounded environment, timeout, and expected output identity.
- The first parser set accepts standard Go benchmark output for `ns/op`, `B/op`, and `allocs/op`, plus a versioned UltraPlan JSON envelope for `ms`, `MiB`, `ops/s`, `tasks/run`, and `%` source values.
- Benchmark authoring occurs only in a disposable copy. A benchmark patch may touch only validated benchmark, test, fixture, or narrowly approved measurement-support paths and may not change product behavior.
- Product code runs benchmark correctness and coverage checks before promotion. Failed, ambiguous, duplicate, or partial target coverage blocks baseline.
- After promotion, product code freezes the benchmark manifest, file set, file digests, parser version, command descriptors, and target mapping. The same frozen identity measures baseline, each candidate, and the final worktree.
- Once baseline begins, neither the runtime nor a command can edit or replace the frozen targets, benchmark implementation, parser, sample policy, correctness commands, or environment policy for that attempt.

### AC-5: measurement qualification cannot manufacture a pass

- Each target uses bounded warmups followed by the declared number of measured samples. Raw samples and bounded output evidence are retained before aggregation.
- Product code rejects missing, duplicate, ambiguous, non-finite, overflowed, wrong-unit, wrong-scenario, wrong-benchmark, or insufficient sample data.
- Qualification records the aggregation method, dispersion, coefficient of variation, environment identity, command identity, benchmark digest, and target packet digest.
- Excessive variance, thermal or resource instability detected by policy, environment drift, parser uncertainty, or benchmark drift produces `inconclusive` or `blocked`. It never produces `met`.
- Absolute and baseline percentage comparisons follow the normalized target semantics exactly. Boundary equality honors the declared inclusive comparator.
- Runtime prose, profile interpretation, or a previous result cannot determine `met` or `missed`; product code derives every comparison from qualified numeric evidence.

### AC-6: optimization is bounded, isolated, and correctness-preserving

- Optimization starts only for a qualified missed `required` target or an explicitly selected missed `report` target. A baseline-only row does not authorize optimization.
- Each cycle addresses one target-linked symptom, records one bounded profile and hypothesis, and requests one proposal in a fresh disposable copy based on the current accepted worktree identity.
- The runtime cannot write to the canonical sprint worktree. Product code validates patch syntax, protected roots, allowed production paths, file count, changed bytes, generated artifacts, symlinks, repository-control files, and actual diff before considering promotion.
- A candidate runs the full product-owned correctness gate before its performance result can be considered. Any functional regression rejects the proposal regardless of performance improvement.
- A proposal is promotable only when it is current, contained, cleanup is proven, required correctness passes, the targeted metric improves under the frozen benchmark, and no already-met required target becomes missed or inconclusive.
- Promotion uses the product-owned mutation boundary and revalidates identity immediately before apply. Pre-existing unrelated changes are reported and preserved; overlapping changes block promotion.
- Cycles, runtime attempts, command counts, time, output, profile size, proposed bytes, changed files, changed bytes, and total wall time have finite maxima and persist across resume. Exhaustion ends `stalled` or `target_missed`; it never resets limits or converts a miss into a pass.

### AC-7: protected sources and evidence remain immutable

- Performance work cannot modify `requirements.md`, project policy, target packets, benchmark manifests after freeze, parsers, correctness commands, sample policy, prior measurement evidence, `flow-state.json` directly, planning artifacts, review or QA evidence, workspace configuration, or Git metadata.
- Benchmark authoring cannot change production code. Optimization cannot change benchmark or test code, expected outputs, fixtures used as correctness or measurement authority, or acceptance criteria.
- Protection covers direct writes, renames, deletes, hard links, symlinks, subprocess descendants, formatters, generators, and cleanup.
- The runtime cannot commit, branch, reset, checkout, stash, clean, push, edit hooks, or mutate the Git index.
- A violation is retained as bounded rejected evidence when safe, never promoted, and ends with a concrete blocker or failure reason.

### AC-8: outcomes and phase ordering are explicit

- Target outcomes are `met`, `missed`, `baseline_recorded`, `report_only`, `inconclusive`, or `blocked`.
- Run outcomes are `passed`, `passed_with_reports`, `target_missed`, `blocked`, `cancelled`, `cleanup_uncertain`, or `stalled`.
- `passed` requires every current `required` target to be qualified and met, every required correctness command to pass, cleanup to be proven, and final target, benchmark, environment, and worktree identities to match.
- `passed_with_reports` has the same required-target and correctness guarantees while retaining non-blocking report or baseline findings.
- A required `missed`, `inconclusive`, or `blocked` target prevents Conformance Review, QA, verified, and merge-ready transitions. Report and baseline outcomes remain visible but do not block.
- `performance.md`, private state, durable operation result, CLI text and JSON, TUI, browser, and `flow-state.json` summary agree on the outcome and current fingerprints.
- A runtime cannot declare global success. Exactly one terminal result is committed from product-owned facts, and late completion or stale writers cannot replace it.

### AC-9: later mutation invalidates performance evidence

- Performance freshness binds requirements, project policy, execute evidence, target worktree, target packet, benchmark manifest, parser, correctness policy, environment policy, and final implementation digests.
- Any later implementation or repair mutation marks the performance result stale before it can claim verified or merge-ready state.
- Repair reverification includes applicable frozen performance targets after functional checks. A successful repair cannot reuse measurements from the pre-repair implementation.
- When a performance optimization changes the implementation, Conformance Review and empirical QA run against the resulting fingerprint rather than earlier execute evidence.
- Status and every interface show stale reason, changed fingerprint, blocking effect, and the exact rerun action.

### AC-10: cancellation, recovery, persistence, and interface parity preserve authority

- One durable run has one owner and writer fence. Durable acceptance and state publication occur before runtime or command execution.
- Cancellation is idempotent, reaches runtime and subprocess descendants, waits for bounded cleanup, and records `cancelled` or `cleanup_uncertain`; it never records a pass.
- Resume restores frozen inputs, current phase, consumed counters, deadlines, accepted mutations, and the last proven boundary. It rejects changed requirements, project policy, benchmark identity, correctness policy, target worktree, or incompatible schema.
- Startup reconciliation handles owner death, expired leases, partial state publication, late command completion, missing disposable copies, and cleanup uncertainty conservatively.
- State writes are private, contained, no-follow, atomic, versioned, digest-bound, and stale-writer fenced. Detailed evidence stays outside `flow-state.json`.
- CLI text, CLI JSON, TUI, browser HTML, browser JSON, durable operation records, `performance.md`, `performance-state.json`, and flow-state projection differ only in presentation and intentional boundedness.
- Graceful local-server shutdown uses the existing canonical cancellation and reconciliation contract. Browser refresh, disconnect, session expiry, SSE loss, or another observer does not cancel or complete the run.

### AC-11: tests and dogfood prove the contract

- Unit and fixture tests cover every target schema rule, normalization, digest, parser, aggregation, comparator, qualification, outcome, policy, freshness, and migration branch.
- Integration and adversarial tests cover benchmark scope, frozen identity, protected paths, process bounds, functional regressions, measurement noise, environment drift, proposal rejection, promotion races, cancellation, disk and permission failure, partial publication, cleanup uncertainty, and stale writers.
- Compatibility tests prove disabled projects follow the existing flow and that old project indexes and sprint states remain valid.
- Interface fixtures prove stable additive JSON and application projections for every phase and terminal outcome.
- Focused package tests pass during implementation, followed by `go test ./...`, `go test -race ./...`, and `go build ./cmd/ultraplan`.
- A gated real-runtime dogfood run uses at least one required absolute target, one required baseline-relative target, and one report-only or baseline target. It records both a rejected functional regression and a contained accepted improvement, or truthfully records why no safe improvement was accepted.

## Non-Goals

- Always-on benchmarking or implicit activation from the presence of benchmark files.
- Performance target values in project indexes, workspace configuration, environment variables, CLI flags, prior baselines, prompts, or runtime output.
- Inventing numeric targets when sprint requirements do not govern one.
- A general distributed benchmark service, historical performance database, fleet scheduler, cloud runner, or cross-project leaderboard.
- Unrestricted shell commands, runtime-selected parsers, arbitrary environment forwarding, or direct runtime writes to the canonical implementation worktree.
- Automatic benchmark rewriting after baseline, test weakening, correctness-policy changes, snapshot re-recording, acceptance-criteria changes, or target relaxation.
- Indefinite profiling, search, retries, or optimization.
- Statistical claims beyond the explicitly versioned qualification and comparison rules.
- Replacing Conformance Review, empirical QA, smoke, or bounded repair with performance results.
- Implementing the later workspace-wide content identity, retrieval, graph, cloud, or Aren integration direction.

## Constraints

- Performance remains a `VerificationPhase`, separate from `PlanningStage` and from Conformance Review, QA, smoke, and repair.
- The requirements document is authoritative for targets. `performance.md` and all JSON records are derived evidence.
- Project policy controls only whether the stage is active. Operational settings are safe defaults and lower-only limits, not target authority.
- All commands use explicit argv, a fixed working directory, allowlisted environment, bounded output, timeout, cancellation, and process-tree cleanup.
- Baseline and candidate measurements use the same frozen target packet, benchmark manifest, parser, sample policy, correctness commands, and compatible environment identity.
- Correctness is a hard prerequisite for accepting an optimization. Faster incorrect behavior is failure.
- Product code, not a model, owns target parsing, command descriptors, isolation, protected paths, patch application, measurement qualification, numeric comparisons, counters, freshness, and terminal outcomes.
- All numeric and duration controls have safe finite defaults and immutable maxima. Workspace and environment configuration may only lower operational limits.
- Full prompts, secrets, unsafe environment values, unbounded command output, provider payloads, and unrelated source content are not persisted or exposed publicly.
- No Sprint 39 implementation target is invented for UltraPlan itself. This sprint builds the mechanism; a target table is added to a sprint only when its requirements contain a governed performance commitment and the project policy is enabled.

## Dependencies

- Sprint 35 must continue to provide durable operation acceptance, ownership, writer fencing, cancellation, terminal arbitration, observation, and startup reconciliation.
- Sprint 36 must provide the separate `VerificationPhase` model and bounded flow-state projection.
- Sprint 37 must provide isolated writable copies, bounded evidence, adjudication concepts, and smoke/correctness integration that performance can reuse without collapsing the phases.
- Sprint 38 must provide the product-owned mutation boundary, protected-path enforcement, bounded proposal promotion, progressive reverification, and post-mutation freshness behavior.
- Manual repair must pass its Sprint 38 exit gate before Sprint 39 implementation starts. Automatic repair may remain disabled.
- The implementation repository must support current sprint-worktree identity, explicit command execution, disposable copies, process-tree cleanup, durable operations, private atomic state, and cross-surface application use cases. Missing capability blocks performance rather than weakening the contract.
- Current plan inputs are the six documents recorded under Required Outputs. Planning must cite their current bytes, not a prior sprint summary.

## Review Expectations

- Trace each target from the exact requirements row through normalized target packet, benchmark mapping, raw samples, qualification, comparison, canonical summary, and flow-state projection. There must be no alternate target source.
- Review disabled compatibility first, then enabled admission and measurement, then benchmark authoring, and only then optimization mutation.
- Verify that baseline and every candidate use identical frozen benchmark and parser identities and a compatible environment identity.
- Treat edits to requirements, targets, frozen benchmarks, tests, correctness policy, evidence, configuration, Git control files, or unapproved paths as security failures.
- Recalculate representative absolute and baseline-relative verdicts independently, including equality boundaries, negative percentage thresholds, noisy data, and non-finite values.
- Prove that a faster candidate with a correctness regression is rejected and that cycle exhaustion cannot produce success.
- Prove post-repair invalidation and the final ordering of performance, Conformance Review, and QA against the same implementation fingerprint.
- Compare all interfaces and persisted artifacts for each outcome. Differences must be presentation-only, with bounded summaries pointing to the same immutable evidence.
- Verify exactly-one-terminal-result behavior under cancellation, timeout, server shutdown, persistence failure, target drift, benchmark drift, cleanup uncertainty, late runtime completion, and restart.
- Reject any implementation that silently ignores a requirements target, derives verdicts from model prose, permits target overrides outside requirements, mutates the canonical worktree directly from the runtime, or treats inconclusive evidence as passing.
