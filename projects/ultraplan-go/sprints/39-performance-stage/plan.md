# Sprint Plan: Requirements-Driven Performance Stage

> Project: `ultraplan-go`
> Sprint: `39-performance-stage`
> Source: `reasoning.md`
> **Inputs Used:** `projects/ultraplan-go/project-index.md`, `projects/ultraplan-go/roadmap.md`, `projects/ultraplan-go/docs/ARCHITECTURE.md`, `projects/ultraplan-go/docs/PRD.md`, `projects/ultraplan-go/docs/TRD.md`, `projects/ultraplan-go/sprints/39-performance-stage/requirements.md`, `projects/ultraplan-go/sprints/39-performance-stage/sprint-index.md`, `projects/ultraplan-go/sprints/39-performance-stage/technical-handbook.md`, `projects/ultraplan-go/sprints/39-performance-stage/reasoning/architecture.md`, `projects/ultraplan-go/sprints/39-performance-stage/reasoning/api-design.md`, `projects/ultraplan-go/sprints/39-performance-stage/reasoning/frontend.md`, `projects/ultraplan-go/sprints/39-performance-stage/reasoning.md`, `PRODUCT.md`, `DESIGN.md`, `docs/plans/integrated-roadmap.md`, `docs/plans/post-execution-qa-and-repair-loop.md`, `docs/plans/server-shutdown-run-cancellation-contract.md`, `docs/plans/retrieval-ready-content-plan.md`, `/home/antonioborgerees/coding/ultraplan/Aren/docs/phased-roadmap.md`, `/home/antonioborgerees/coding/ultraplan/Aren/docs/performance-engineering.md`

This plan executes `reasoning.md`. It does not reopen the sprint's architecture, authority, statistical policy, or product limits. Paths in implementation tasks are relative to the UltraPlan Go implementation repository root unless a workspace path is shown explicitly.

## Reasoning Source

- **Sprint Reasoning:** `reasoning.md`
- **Sprint Index:** `sprint-index.md`
- **Technical Handbook:** `technical-handbook.md`
- **Area Reasoning:** `reasoning/architecture.md`, `reasoning/api-design.md`, `reasoning/frontend.md`

## Sprint Status

- **Status:** executed; dogfooding deferred by user
- **Owner:** implementation agent
- **Start Date:** 2026-09-04
- **Completion Date:** 2026-09-06

## Decisions To Execute

| Decision | Source Section | Execution Implication |
| --- | --- | --- |
| Activation and target authority are strict and separate. | `reasoning.md#decision-1-activation-and-target-authority-are-strict-and-separate` | Parse activation only in `internal/project`; parse and normalize every target only from current `requirements.md` bytes in `internal/sprint`; reject policy disagreement and every alternate target source. |
| Performance is an optional verification phase with runtime-free admission. | `reasoning.md#decision-2-performance-is-an-optional-verification-phase-with-runtime-free-admission` | Insert `performance` before Conformance Review only for enabled projects. Disabled paths remain byte-preserving. Prepare and dry-run perform reads and deterministic validation only. |
| Benchmarks use closed product-owned runners and freeze before warmup. | `reasoning.md#decision-3-benchmarks-use-closed-product-owned-runners-and-freeze-before-warmup` | Support only `go-benchmark-v1` and `ultraplan-json-v1`. Product code constructs command descriptors, proves coverage, promotes benchmark-only changes, and freezes the complete manifest before measurement. |
| Qualification uses one complete median sample set and exact product comparisons. | `reasoning.md#decision-4-qualification-uses-one-complete-median-sample-set-and-exact-product-comparisons` | Run two default warmups, retain every measured sample, aggregate by exact median, qualify with `cv-v1`, and derive every inclusive comparator and outcome in product code. |
| Optimization is one-miss-at-a-time, isolated, and correctness-first. | `reasoning.md#decision-5-optimization-is-one-miss-at-a-time-isolated-and-correctness-first` | Profile one required miss, request one proposal in a fresh disposable copy, run frozen correctness before measurement, reject protected or overlapping changes, and promote only a strict qualified improvement with no required regression. |
| Every operational control has a finite default and hard maximum. | `reasoning.md#decision-6-every-operational-control-has-a-finite-default-and-hard-maximum` | Encode the frozen V1 limits below. Workspace and environment sources may only lower defaults. Persist counters and the absolute deadline across resume. |
| Private immutable evidence, digest freshness, and fenced recovery own truth. | `reasoning.md#decision-7-private-immutable-evidence-digest-freshness-and-fenced-recovery-own-truth` | Share only verification writer-fence and atomic-write mechanics. Keep performance records separate, publish immutable evidence before summaries, verify digest-bound pointers, and reconcile partial publication without inferring success. |
| One typed app boundary serves CLI, TUI, browser, and durable operations. | `reasoning.md#decision-8-one-typed-app-boundary-serves-cli-tui-browser-and-durable-operations` | Add bounded performance use cases and DTOs, reuse guarded operations and run control, keep operational lifecycle separate from product outcome, and reject stale confirmations before acceptance. |
| The browser and TUI present one sprint-scoped canonical workbench. | `reasoning.md#decision-9-the-browser-and-tui-present-one-sprint-scoped-canonical-workbench` | Render app-owned snapshots and bounded evidence metadata. Preserve server rendering, no-JavaScript controls, accessibility, and observer-independent cancellation. |
| Documentation, review, and dogfood trace the whole authority chain. | `reasoning.md#decision-10-documentation-review-and-dogfood-trace-the-whole-authority-chain` | Update all named docs, add one canonical cross-interface fixture family, run focused and repository-wide verification, then perform architecture review, sprint review, and gated real-runtime dogfood. |

## Requirements / Contracts To Satisfy

| Contract / Requirement ID | Required Behavior | Evidence Planned |
| --- | --- | --- |
| AC-1 | Missing or disabled policy preserves the existing flow and bytes; enabled policy inserts performance before Conformance Review. | Project parser tests, disabled side-effect sentinels, byte snapshots, phase-order tests, old workspace fixtures. |
| AC-2 | The exact requirements table is the sole target authority and normalizes deterministically. | Hostile Markdown and cell matrix tests, golden target packets and digests, policy-disagreement and alternate-source tests. |
| AC-3 | Admission resolves every prerequisite deterministically; prepare and dry-run start no child work and write nothing. | Fake runtime, process, isolation, and store collaborators that fail on invocation; idempotency and conflict tests. |
| AC-4 | Every target maps to one closed descriptor and one frozen benchmark identity before warmup. | Descriptor conversion, parser identity, coverage, benchmark-only scope, manifest digest, promotion race, and drift tests. |
| AC-5 | Complete sample sets, exact identity checks, conservative qualification, and product comparisons cannot manufacture a pass. | Median, standard deviation, CV, equality, negative threshold, zero baseline, malformed output, noise, truncation, and drift fixtures. |
| AC-6 | Optimization is finite, target-linked, isolated, correctness-first, and promotes only a current contained improvement. | Fresh-copy, profile selection, protected-path, overlap, correctness ordering, faster-regression rejection, no-progress, and exhaustion tests. |
| AC-7 | Governed inputs, measuring instruments, correctness authority, prior evidence, configuration, Git, and unrelated paths remain immutable. | Rename, delete, symlink, hard-link, descendant, generated-file, Git-control, test/fixture, and evidence mutation adversaries. |
| AC-8 | Exact target and run outcomes gate later verification and agree across artifacts and interfaces. | Outcome precedence tables, flow-gate tests, app projection fixtures, canonical Markdown and JSON semantic assertions. |
| AC-9 | Requirements, policy, execute, benchmark, environment-policy, correctness-policy, and implementation drift make evidence stale. | Fingerprint reason tests, post-promotion and post-repair invalidation, stale status, rerun-action, and later-phase gate tests. |
| AC-10 | Durable acceptance, writer fencing, cancellation, resume, recovery, and interface parity preserve one authority and one terminal result. | Acceptance-before-child, owner-death, expired-lease, late-writer, cancellation race, partial publication, cleanup uncertainty, restart, and parity tests. |
| AC-11 | Unit, integration, adversarial, race, build, review, and real-runtime evidence prove the contract. | Focused package tests, `go test ./...`, `go test -race ./...`, build, three review protocols, and gated dogfood. |
| Architecture and Security | Product logic remains in `internal/project` and `internal/sprint`; platform and interface code cannot gain target, verdict, or mutation authority. | Dependency review, package-cycle checks, explicit argv tests, containment tests, architecture review. |
| CLI Surface, Errors, and Documentation | Commands, output, exits, errors, help, examples, schemas, and recovery guidance are stable and actionable. | CLI fixtures, stable error codes, stdout/stderr assertions, documentation checks, sprint review. |
| Configuration and Performance | Model selection is stage-specific; operational controls are finite, source-aware, lower-only, and cannot alter targets. | Config merge tables, source reporting, upward-override rejection, limit and deadline persistence tests. |
| LLM Runtime and LLM Evaluation / Cost / Safety | Agentwrap/OpenCode proposes bounded mappings, hypotheses, and patches but cannot decide commands, comparisons, limits, promotion, or success. | Runtime request and permission tests, malformed output, redaction, bounded attempts, gated real-runtime evidence. |
| Observability and Workflows | Correlation, phase progress, cancellation, recovery, and terminal arbitration are durable, bounded, and consistent. | Event payload bounds, replay-gap tests, correlation fixtures, cancellation and terminal race tests. |
| Persistence And Migrations | Private records are strict, versioned, immutable, contained, atomic, digest-bound, fenced, retained, and recoverable. | Migration fixtures, no-follow and permission tests, write-failure hooks, digest corruption, stale writers, retention tests. |
| Testing | Tests cover pure rules, volatile boundaries, full workflows, compatibility, failures, and real operation at the appropriate level. | Table tests, fakes, temporary worktrees, command and interface fixtures, race tests, review and dogfood. |

## Frozen V1 Implementation Policy

Implementation must copy these rules into typed product constants and tests. Changes require a governed reasoning amendment, not an implementation shortcut.

### Target table

```markdown
## Performance Targets

| ID | Scenario | Metric | Comparator | Value | Unit | Gate | Samples | Basis |
| --- | --- | --- | --- | ---: | --- | --- | ---: | --- |
| PERF-001 | `<bounded scenario>` | `<stable metric>` | `<=` | `<number>` | `ns/op` | `required` | `10` | `absolute` |
```

- Accept one to 100 rows. IDs match the exact `PERF-` prefix followed by uppercase ASCII letters, digits, or hyphens and are unique.
- Bound trimmed valid UTF-8 Scenario text to 1 through 256 bytes and Metric text to 1 through 128 bytes.
- Accept only `<=`, `>=`, or `baseline`; `required` or `report`; `absolute`, `baseline`, or `current`; and `ns/op`, `B/op`, `allocs/op`, `ms`, `MiB`, `ops/s`, `tasks/run`, or `%`.
- Numeric values use finite base-10 decimal syntax without exponents. Canonicalize signs, zero, leading zeros, and trailing fractional zeros while retaining exact rational threshold semantics.
- `baseline` requires `Value=-`, `Gate=report`, and `Basis=current`. Threshold comparators use `absolute` or `baseline`; baseline-relative rows use `%`.
- Sort normalized rows bytewise by target ID and bind the requirements digest, project-policy digest, parser version, normalized rows, and packet digest.

### Descriptor families

| Runner | Product-owned command shape | Accepted output |
| --- | --- | --- |
| `go-benchmark-v1` | `go test` argv for one contained package and one exact `Benchmark...` symbol | Standard Go benchmark output for `ns/op`, `B/op`, and `allocs/op`. |
| `ultraplan-json-v1` | `go test` argv for one contained package and one exact `Test...` measurement symbol | Exactly one marked V1 JSON envelope for `ms`, `MiB`, `ops/s`, `tasks/run`, and `%` source values. |

Runtime proposals may identify only runner kind, contained package locator, exact symbol, and target mapping. Product code supplies executable identity, argv, cwd, environment allowlist, timeout, parser, and output identity.

### Qualification and comparison

- Use two warmups by default, then exactly the requirements-owned sample count from 5 through 100.
- Run warmups and measured samples serially. Retain all bounded samples. Never delete outliers or request another set because the first set missed or failed qualification.
- Aggregate with `median-v1`. For even sets, use the exact arithmetic midpoint of the two central exact decimal values.
- Record arithmetic mean, sample standard deviation using `n-1`, and `CV = 100 * sample_standard_deviation / abs(mean)` under `cv-v1`.
- All-zero values have CV zero. A zero mean with any nonzero sample is inconclusive. The set qualifies only when every identity matches, cleanup is proven, and CV is at or below the effective threshold.

| Target form | Numeric fact | Target outcome |
| --- | --- | --- |
| Threshold, `Basis=absolute`, `Gate=required` | Compare qualified final median directly with exact `Value`; equality passes. | `met` or `missed`; failed qualification is `inconclusive` or `blocked`. |
| Threshold, `Basis=baseline`, `Gate=required` | Compare exact `100 * (final - baseline) / baseline` with `Value`; equality passes and raw units match. | `met` or `missed`; zero or non-finite baseline and failed qualification are `inconclusive` or `blocked`. |
| Threshold, `Gate=report` | Calculate the same comparison without creating a blocking target verdict. | `report_only`; failed qualification remains `inconclusive` or `blocked` but does not block at target-gate level. |
| `Comparator=baseline`, `Basis=current`, `Gate=report` | Record the qualified final current median and retain the initial baseline separately if code changed. | `baseline_recorded`; failed qualification is `inconclusive` or `blocked` and non-blocking at target-gate level. |

| Precedence | Condition | Run outcome |
| ---: | --- | --- |
| 1 | Owned descendant or disposable-copy cleanup cannot be proven. | `cleanup_uncertain` |
| 2 | Durable cancellation wins and cleanup is proven. | `cancelled` |
| 3 | Current authority, persistence, correctness, required qualification, isolation, or recovery is blocked or inconsistent. | `blocked` |
| 4 | No-progress or proposal/runtime limits prevent the next proven convergence boundary. | `stalled` |
| 5 | Final qualified evidence contains a required miss after allowed work ends or optimization is not authorized. | `target_missed` |
| 6 | Every required target is current, qualified, and met; correctness and cleanup pass; at least one report or baseline row remains. | `passed_with_reports` |
| 7 | Every target is required and met with current identities, final correctness, and proven cleanup. | `passed` |

### Operational limits

| Control | Built-in default | Immutable hard maximum | Counting rule |
| --- | ---: | ---: | --- |
| Target rows | 100 | 100 | Every normalized requirements row; validation rejects excess. |
| Benchmark authoring runtime attempts | 2 | 4 | Every started authoring runtime request, including repaired structured output. |
| Optimization runtime attempts per cycle | 2 | 4 | Every started hypothesis or proposal runtime request in that cycle. |
| Warmups per target and measurement phase | 2 | 5 | Every started warmup command. |
| Requirements-owned measured samples | 100 operational ceiling | 100 schema maximum, minimum 5 | Every started measured command; the row value is never reduced. |
| Qualification CV threshold | 5% | 10% representable ceiling | Effective threshold may only be lowered from 5%; lower is stricter. |
| Total commands per attempt | 128 | 256 | Benchmark, coverage, warmup, sample, profile, correctness, and cleanup commands. |
| Time per command | 5 minutes | 15 minutes | Wall time from process start through descendant cleanup. |
| Captured stdout per command | 256 KiB | 1 MiB | Bytes observed before truncation. |
| Captured stderr per command | 256 KiB | 1 MiB | Bytes observed before truncation. |
| Retained profile per cycle | 2 MiB | 8 MiB | Profile bytes before acceptance. |
| Retained proposal patch per cycle | 512 KiB | 2 MiB | Serialized patch bytes before validation. |
| Changed production files per proposal | 12 | 32 | Unique actual paths after rename and delete normalization. |
| Changed production bytes per proposal | 512 KiB | 2 MiB | Added plus deleted bytes from the actual diff. |
| Optimization cycles per attempt | 4 | 8 | Every entered cycle, accepted or rejected. |
| Total attempt wall time | 60 minutes | 120 minutes | Absolute deadline from durable acceptance, including resumes. |
| Retained performance attempts per sprint | 3 | 8 | Terminal attempt directories; current and referenced attempts are protected. |
| Cleanup reserve | 30 seconds | 60 seconds | Separate bounded context after main cancellation or failure. |
| Evidence query page size | 50 | 200 | Metadata records only, with an attempt and filter-bound cursor. |

## Tasks

- [x] **Task 1: Verify dependency gates and generalize writer ownership narrowly**
  > Executes: Decisions 2 and 7; AC-1, AC-7, AC-10
  - [x] Before implementation, record evidence that Sprint 38 manual single-issue repair passes end to end, automatic repair limits remain fixed, and repair cannot weaken evidence, requirements, tests, or acceptance criteria. Stop Sprint 39 implementation if this gate is not proven.
  - [x] Re-read and record current digests for `docs/plans/integrated-roadmap.md`, `docs/plans/post-execution-qa-and-repair-loop.md`, `docs/plans/server-shutdown-run-cancellation-contract.md`, `docs/plans/retrieval-ready-content-plan.md`, `../Aren/docs/phased-roadmap.md`, and `../Aren/docs/performance-engineering.md` before code changes. Any material conflict returns to governed reasoning.
  - [x] Replace the QA-named writer token and service fence with a verification-owned token carrying operation run ID, operational attempt ID, and fencing generation. Update QA and repair without changing their schemas, paths, outcomes, or behavior.
  - [x] Extract only unexported mechanical helpers for contained no-follow resolution, private directory and file creation, strict JSON decode, immutable create, digesting, atomic replace, directory sync, and pre-rename fence checks. Do not create a generic verification store or workflow.
  - [x] Add compatibility and race tests proving QA and repair writer ownership, stale-writer rejection, atomic publication, and existing fixtures remain unchanged.

- [x] **Task 2: Implement project activation and strict requirements-owned targets**
  > Executes: Decision 1; AC-1, AC-2, AC-7
  - [x] Add typed `PerformancePolicy` modes and disabled-by-default behavior in `internal/project/domain.go`. Carry the policy source digest without adding target fields.
  - [x] Extend `internal/project/index.go` with a Markdown-aware exact `## Performance Policy` parser. Accept one `Mode` key only. Reject duplicate sections, duplicate keys, unknown keys, invalid modes, and scenario, metric, comparator, value, unit, gate, samples, basis, or other target-like fields.
  - [x] Add `internal/project/index_test.go` coverage for missing, disabled, enabled, duplicate, unknown-field, invalid-mode, and attempted-target-override cases, including fenced, quoted, commented, and inline pseudo-sections.
  - [x] Implement the context-aware `Performance Targets` scanner in `internal/sprint/index.go`. Track fenced blocks, HTML comments, blockquotes, heading level, duplicate headings, immediate table association, exact ordered headers, separators, row count, and trailing malformed rows.
  - [x] Implement the complete V1 cell grammar, text bounds, finite canonical decimal normalization, exact rational threshold representation, comparator/gate/basis/unit rules, stable ID ordering, duplicate detection, and deterministic normalized packet creation.
  - [x] Join typed project policy and target parsing at the sprint service validation boundary. Enabled policy requires one valid table; disabled policy plus a declaration reports the required mismatch without treating target rows as active.
  - [x] Add stable golden packet and digest fixtures plus `internal/sprint/index_test.go` tables for every accepted and rejected contract branch, including same-byte stability and changed-requirements invalidation.
  - [x] Update `internal/workspace/init.go`, `internal/workspace/scaffold/templates/requirements.md`, and `internal/workspace/scaffold/prompts/create-requirements.md` with optional exact table guidance that either leaves numeric cells for governed completion or permits a baseline report row. Add `internal/workspace/workspace_test.go` assertions that no numeric threshold is fabricated.

- [x] **Task 3: Add performance domain types, limits, configuration, and pure verdict rules**
  > Executes: Decisions 4 and 6; AC-2, AC-5, AC-8
  - [x] Create `internal/sprint/performance_types.go` with typed policy snapshots, targets, packets, descriptors, manifests, environments, measurements, qualifications, profiles, hypotheses, proposals, comparisons, freshness, target outcomes, run outcomes, limits, counters, phases, correlations, blockers, and artifact references.
  - [x] Encode all frozen V1 constants and combinations from this plan. Keep exact target threshold arithmetic separate from floating-point dispersion calculations and reject non-finite or overflowed observations.
  - [x] Implement `internal/sprint/performance_targets.go` for packet digest binding, stable ordering, exact absolute and baseline-relative comparison, equality behavior, target outcomes, run-outcome precedence, and target coverage accounting.
  - [x] Add stage-specific performance model and variant fields through the existing configuration mechanism. Runtime selection must not carry target, descriptor, parser, command, environment, sample, gate, limit, or verdict authority.
  - [x] Implement source-aware lower-only limit resolution. Reject zero, negative, overflowed, or upward settings. Block admission when an effective sample ceiling is below a requirements-owned row and report both values.
  - [x] Add table tests for defaults, hard maxima, lower-only merges, sources, exact boundaries, all outcomes, and limit exhaustion. Verify counters and absolute deadlines cannot reset on resume.

- [x] **Task 4: Add verification ordering, flow-state compatibility, and deterministic admission**
  > Executes: Decision 2; AC-1, AC-3, AC-8, AC-9
  - [x] Add `VerificationPhasePerformance` before Conformance Review in `internal/sprint/verification_phase.go` without adding a `PlanningStage` or compatibility stage.
  - [x] Extend `internal/sprint/domain.go`, `flow.go`, and `verification_phase.go` so disabled projects keep `execute -> Conformance Review`, while enabled projects require a fresh `passed` or `passed_with_reports` result before later verification.
  - [x] Add an optional bounded performance projection to `FlowState` and migrate old schemas to no projection. Update `state.go` and `state_database.go` validation, checkpoint preservation, database/file parity, contained pointer rules, and stage readiness.
  - [x] Create runtime-free admission in `internal/sprint/performance.go`. Resolve policy, target packet, execute success and evidence, sprint worktree identity, correctness descriptors, isolation and cleanup support, process and environment policy, runner/parser versions, limits, unrelated-change overlap, conflicting writers, and current freshness.
  - [x] Make prepare and dry-run return stable policy, targets, expected coverage, protected roots, correctness commands, effective limits and sources, missing prerequisites, mutation possibility, and next action without constructing a runtime, running a command, accepting a durable run, creating performance storage, or changing artifacts.
  - [x] Add disabled byte-preservation tests and side-effect sentinels. Add admission matrix, idempotency, stale execute, overlap, isolation, correctness-policy, limit, writer-conflict, and acceptance-before-child-work tests.

- [x] **Task 5: Build private performance persistence and recovery boundaries**
  > Executes: Decision 7; AC-7, AC-9, AC-10
  - [x] Create `internal/sprint/performance_state.go` with strict private schemas and contained path builders for the complete required artifact family under `verification/performance/attempts/<attempt-id>/` plus `verification/performance-state.json`.
  - [x] Implement immutable target packet, benchmark manifest, environment, baseline, cycle profile, proposal patch, scope, correctness, measurements, cleanup, and terminal result writes. Identical retries may succeed only when bytes match; conflicts fail closed.
  - [x] Enforce private `0700` directories, `0600` files, no-follow traversal, bounded retained bytes, digest-bound workspace-relative pointers, state schema validation, writer checks before write and rename, same-directory atomic current-state replacement, sync, and protected-current retention.
  - [x] Implement the publication order `immutable result -> performance.md -> performance-state.json -> flow-state.json projection -> run-control terminal observation`. Do not claim cross-file atomicity.
  - [x] Implement bounded status reads from current state and direct pointers. Old evidence discovery must use bounded metadata queries, never recursive status reconstruction.
  - [x] Implement conservative resume and recovery at proven boundaries. Verify frozen identities, counters, deadline, apply journal, pointers, and cleanup. Recovery may finish later publication from valid immutable evidence but cannot launch runtime work, reset limits, adopt workers, infer success, or replace a terminal result.
  - [x] Add strict schema, unknown-field, permissions, no-follow, immutable collision, digest corruption, partial write, rename, sync, quota, retention, unsupported active schema, owner-death, expired-lease, stale-writer, late-completion, cancellation, and exactly-one-result tests with race coverage.

- [x] **Task 6: Implement closed benchmark discovery, authoring, coverage, and freezing**
  > Executes: Decision 3; AC-3, AC-4, AC-7
  - [x] Create `internal/sprint/performance_benchmark.go` with the two closed V1 descriptor families and product-owned conversion to explicit `process.Request` values.
  - [x] Validate contained package locators, exact Go symbol forms, unique output identities, raw units, target mappings, working directories, executable identity, argument vectors, environment keys, timeouts, and parser versions.
  - [x] Discover existing valid descriptors without executing them during dry-run. For missing coverage, request a bounded structured runtime proposal in a fresh disposable copy.
  - [x] Permit benchmark authoring changes only in target-package `*_test.go`, package-local test-only helpers, and target-package `testdata/performance/` data that is not correctness authority. Use hunk and actual-diff validation, not suffix checks alone.
  - [x] Reject production changes, existing assertion or expected-output weakening, general fixtures, configuration, governance, evidence, repository-control files, Git state, symlinks, hard links, renames or deletes outside scope, and unapproved generated files.
  - [x] Run product-owned benchmark correctness and coverage checks. Require exactly one descriptor per target, no duplicate output identity, full coverage, proven cleanup, current source identity, and no unrelated overlap before benchmark promotion.
  - [x] Promote through the existing product mutation boundary, then freeze descriptors, benchmark file set and digests, target mapping, runner/parser versions, raw units, correctness coverage, environment policy, and manifest digest before warmup.
  - [x] Add descriptor, argv, environment, output identity, complete/partial/duplicate coverage, benchmark-scope adversary, promotion race, cleanup, freeze, and benchmark-drift tests.

- [x] **Task 7: Implement bounded measurement, parsing, qualification, and comparison**
  > Executes: Decision 4; AC-4, AC-5, AC-8
  - [x] Create `internal/sprint/performance_measure.go` with parser dispatch for standard Go benchmark output and the exact marked UltraPlan JSON V1 envelope.
  - [x] Reject missing, duplicate, ambiguous, malformed, non-finite, overflowed, wrong-target, wrong-scenario, wrong-metric, wrong-unit, wrong-symbol, wrong-command, truncated, or insufficient output before aggregation.
  - [x] Capture bounded machine, OS, architecture, runtime, toolchain, process policy, noise-control, command, benchmark, target-packet, and implementation identities in `environment.json` and every sample set.
  - [x] Run two default warmups and exactly the declared samples through `internal/platform/process` with explicit argv, fixed contained cwd, allowlisted environment, timeout, bounded output, cancellation, descendant cleanup, and serial measurement.
  - [x] Retain raw bounded facts before calculating exact median, mean, sample standard deviation, CV, dynamic stability facts, qualification, absolute comparisons, baseline-relative percentages, and exact inclusive boundaries.
  - [x] Treat parser uncertainty, output truncation, excessive CV, dynamic instability, environment drift, benchmark drift, and zero baseline conservatively. Distinguish `inconclusive` from broken-authority or cleanup `blocked` cases and never produce `met` from either.
  - [x] Measure every target at baseline and every target again against the final canonical implementation, including a no-change attempt.
  - [x] Add pure arithmetic tests and fake-process integration tests for odd/even medians, exact decimals, equality, negative percentages, zero and mixed-zero means, CV boundary, all identities, wrong output, noise, drift, cancellation, timeout, truncation, cleanup, and serial execution.

- [x] **Task 8: Implement bounded isolated optimization and product-owned promotion**
  > Executes: Decision 5; AC-6, AC-7, AC-8, AC-9
  - [x] Create `internal/sprint/performance_optimize.go`. Select qualified required misses by stable target ID. Do not optimize baseline rows or report targets in V1.
  - [x] For each cycle, select the product-owned profile type from the unit, run one bounded profile, retain one target-linked hypothesis, and request one structured proposal in a fresh disposable copy based on the current accepted implementation identity.
  - [x] Allow optimization changes only to approved production source. Deny requirements, project policy, planning artifacts, benchmarks, tests, test helpers, fixtures, expected outputs, snapshots, parser/authority code, correctness policy, configuration, all verification evidence, repository controls, Git, links, and unapproved generated paths.
  - [x] Derive and validate the actual patch, changed files and bytes, before/after identities, scope, unrelated overlap, process descendants, and cleanup. Preserve every pre-existing unrelated change and block overlap.
  - [x] Run the full frozen correctness command set before candidate measurement. A failure rejects the candidate without performance promotion.
  - [x] Under the frozen manifest, require a qualified strict improvement in the selected comparator's passing direction. Equality is not improvement. Remeasure every already-met required target and reject any miss, inconclusive result, or block.
  - [x] Write and sync the apply journal, recheck canonical source identity immediately before apply, apply through the product mutation boundary, verify the resulting identity, and invalidate downstream Conformance Review, QA, smoke, and repair evidence only after successful mutation publication.
  - [x] Persist every rejected proposal and reason within bounds. End no-progress or exhausted work as `stalled`, `target_missed`, or `blocked` from proven facts, never as passed.
  - [x] Add unit, integration, failure-injection, cancellation, and race tests for profile selection, fresh copies, one-target attribution, path and diff adversaries, correctness-before-measure ordering, faster functional regression, strict improvement, required-target regression, overlap, apply-journal crash boundaries, cleanup, convergence, and exhaustion.

- [x] **Task 9: Complete service orchestration, freshness, publication, and cancellation**
  > Executes: Decisions 2, 5, 6, and 7; AC-3 through AC-10
  - [x] Complete `VerificationPhasePerformance` in `internal/sprint/performance.go` with ordered admission, durable acceptance handoff, target freeze, benchmark preparation, baseline, bounded optimization, final measurement, final correctness, publication, resume, cancellation, and recovery.
  - [x] Persist phase, counters, deadline, accepted mutations, current identities, last proven boundary, blocker, next action, and evidence pointers before moving to the next side effect.
  - [x] Bind freshness to requirements, policy, execute evidence, target packet, benchmark manifest, runner/parser, sample/aggregation, correctness, process, environment policy, initial and final implementation, workflow, and terminal-result digests.
  - [x] Emit stable freshness reason codes and the exact rerun action. Do not recapture dynamic environment observations during status or use modification times as authority.
  - [x] Centralize post-mutation invalidation so performance promotion invalidates downstream verification, and any later implementation or repair mutation makes performance stale before verified or merge-ready status.
  - [x] Integrate repair reverification so frozen functional checks run first and applicable frozen performance targets run against the repaired fingerprint. Never reuse pre-repair samples as current evidence.
  - [x] Route cancellation to proposal runtimes and process descendants, stop new work, use the separate cleanup reserve, fence late publications, and commit `cancelled` or `cleanup_uncertain` according to proof.
  - [x] Render the exact required `performance.md` sections from immutable result facts and verify that the report, private state, durable operation, and flow projection agree.
  - [x] Add end-to-end service tests for no-change pass, passed with reports, required miss, blocked authority, stalled convergence, cancellation, cleanup uncertainty, target/benchmark/environment drift, resume, recovery, later mutation, and all terminal races.

- [x] **Task 10: Add shared application operations and CLI**
  > Executes: Decision 8; AC-3, AC-8, AC-10
  - [x] Extend `internal/app/sprint_usecases.go` with adapter-neutral prepare, dry-run, start, status, resume, cancel, recover, result, and bounded evidence-query requests and schema-versioned DTOs.
  - [x] Project all governed target fields, bounded aggregates, qualification, comparison, target outcome, run outcome, policy, phase, freshness reasons, separate operational lifecycle, identities, consumed/effective limits, one blocker, digest-bound references, and next action.
  - [x] Prove DTOs and events omit raw samples, command output, profiles, patches, prompts, provider payloads, environment values, secrets, and private structs.
  - [x] Add performance operation kinds and governed input fingerprints in `internal/app/operations.go`. Reprepare immediately before start or resume acceptance; reject stale facts with no run, command, runtime, artifact, or mutation side effect.
  - [x] Extend `internal/app/operation_runner.go` with one durable performance runner that acquires writer-fence ownership, uses canonical cancellation, and separates operational success from product outcome.
  - [x] Add start alias deduplication, resume binding to the existing attempt and boundary, writer-conflict handling, idempotent cancel and recover, one terminal proposal, and restart reconciliation tests in `internal/app/durable_operations_test.go`.
  - [x] Add `ultraplan sprint <project> <sprint> performance`, `--dry-run`, `status`, `resume`, `cancel`, and `recover` in `internal/app/sprint_commands.go`. Reject target, sample, command, parser, environment, higher-limit, and report-selection inputs.
  - [x] Implement calm text output, one ANSI-free JSON result on stdout, diagnostics on stderr, exact help, stable error codes, and outcome-specific exit mapping. Query exit status reflects query success rather than the stored product verdict.
  - [x] Add `sprint_usecases_test.go`, `sprint_commands_test.go`, operation, JSON schema, stdout/stderr, disabled, stale confirmation, cancellation, restart, stale writer, and exactly-one-terminal-result coverage.

- [x] **Task 11: Add the TUI performance workbench**
  > Executes: Decision 9; AC-1, AC-8, AC-9, AC-10
  - [x] Add a sprint-owned Performance destination using only app use cases and bounded DTOs. Do not read private JSON, parse `performance.md`, inspect samples, or derive outcomes from events.
  - [x] Preserve the existing execute-to-Conformance Review presentation when policy is disabled and no historical performance state exists. Show a compact gate on sprint overview only when enabled or historical state exists.
  - [x] Present policy, freshness, product phase, operational lifecycle, product outcome, required/report counts, blocker, later-verification gate, next action, ordered workflow path, target contracts, bounded measurement facts, limits, and evidence references.
  - [x] Keep lifecycle and product outcome separate and textual. Stale or historical met targets must not look current. Baseline rows must not expose optimization controls.
  - [x] Add guarded dry-run, start, resume, cancel, recover, result, evidence pagination, and refresh actions. Leaving the view or TUI must not cancel work.
  - [x] Support named key bindings, visible focus, one-column narrow rendering, textual status independent of color, bounded content, and cancellation acknowledgement.
  - [x] Add deterministic wide/narrow fixtures and tests for every canonical state, target selection, stale and historical state, guarded actions, cancellation, recovery, evidence gaps, and observer exit.

- [x] **Task 12: Add the server-rendered browser performance workbench**
  > Executes: Decisions 8 and 9; AC-3, AC-8, AC-9, AC-10
  - [x] Add `GET /projects/{project}/sprints/{sprint}/performance` and sprint-owned versioned JSON status, result, and bounded evidence metadata routes. Reuse existing operation prepare/start, durable run, event replay, and cancellation routes.
  - [x] Build explicit handler view models and namespaced feature-specific templates under the current page, layout, component, and presentation hierarchy. Reuse generic primitives only; do not create a generic verification workbench framework.
  - [x] Render the current authority panel, ordered path, expandable target records with all nine governed cells, qualification and bounded aggregates, blockers, stale reasons, next action, limits, digest-bound references, and paged evidence metadata.
  - [x] Implement start, resume, cancel, recover, filter, pagination, result, refresh, and allowlisted artifact-preview flows with ordinary links and forms. A disabled direct page is read-only and exposes no runtime or evidence controls.
  - [x] Preserve Host, Origin, CSRF, session, body, stream, path, redaction, and confirmation controls. Map malformed, missing, conflict, validation, capacity, and internal errors to the decided HTTP statuses and stable codes.
  - [x] Keep JavaScript optional. Durable events may trigger one debounced canonical snapshot fetch. Preserve open disclosures and focused subtrees, expose replay gaps, and never auto-retry mutations or derive verdicts from events.
  - [x] Use semantic headings, landmarks, lists, `details`, forms, labels, buttons, progress, textual outcomes, deliberate focus, one polite coalesced live region, visible errors, and reduced-motion behavior.
  - [x] Add `httptest`, template, route, security, no-JavaScript, hostile-content, bounded-rendering, accessibility, SSE reconnect, one-in-flight fetch, disclosure/focus, slow subscriber, session loss, refresh, and observer-independent cancellation tests.

- [x] **Task 13: Add canonical fixtures, documentation, and compatibility evidence**
  > Executes: Decision 10; AC-1 through AC-11
  - [x] Add one bounded canonical fixture family under `internal/testdata/` for disabled, ready, preparing, active, cancelling, passed, passed-with-reports, target-missed, blocked, cancelled, cleanup-uncertain, stalled, stale, interrupted/resumable, historical, evidence-gap, and recovery-conflict states.
  - [x] Use those fixtures to prove app DTO, CLI text, CLI JSON, TUI, browser HTML, browser JSON, durable operation, `performance.md`, private current state, result, and flow projection agree on authority while differing only in presentation and intentional bounds.
  - [x] Update `docs/architecture.md` with ownership, target authority, isolation, mutation, persistence, publication, freshness, and run-control separation.
  - [x] Update `docs/cli-reference.md` and `docs/user-guide.md` with exact commands, policy enable/disable steps, target schema, examples without invented thresholds, outcomes, gates, cancellation, recovery, stale evidence, and next actions.
  - [x] Update `docs/phase3-json-schemas.md` or its current successor with private schemas, public additive DTOs, optional flow projection, errors, cursors, bounds, and compatibility.
  - [x] Update `docs/recovery.md` with interruption, owner death, expired leases, target/benchmark/source drift, apply-journal states, partial publication, cancellation races, cleanup uncertainty, and restart reconciliation.
  - [x] Update `docs/local-web.md` with focused routes, guarded operations, no-JavaScript journeys, SSE behavior, observer independence, accessibility, and bounded evidence.
  - [x] Update `docs/release-checklist.md` with disabled regression, parser fixtures, correctness regression, qualification, isolation, mutation, persistence, race, interface parity, and real-runtime dogfood gates.
  - [x] Prove old project indexes, old flow states, disabled flows, and existing CLI invocations need no user migration and retain previous artifact bytes.

- [x] **Task 14: Run the release gate, required reviews, and gated dogfood**
  > Executes: Decision 10; AC-11
  - [x] Run focused package tests after each implementation slice, then run `go test ./...`, `go test -race ./...`, and `go build ./cmd/ultraplan` from the implementation repository.
  - [x] Run the Architecture Review protocol. Reject target or verdict logic outside owning modules, a generic verification store/workflow, runtime-selected execution, direct runtime mutation, future Aren capabilities, or any unearned abstraction.
  - [x] Run the Sprint Review protocol. Trace every requirement and target from the exact requirements row through packet, descriptor, manifest, raw samples, qualification, comparison, immutable result, `performance.md`, app DTO, and flow projection.
  - [x] Independently recalculate representative absolute and baseline-relative comparisons, equality boundaries, negative percentage thresholds, zero baseline, noisy sets, and non-finite rejection.
  - [ ] Run the gated Deep Smoke Sprint protocol with at least one required absolute target, one required baseline-relative target, and one report-only or baseline target. Deferred at user direction.
  - [ ] Dogfood a faster functional regression and prove rejection. Record one contained accepted improvement, or record the truthful blocker or safety reason that prevented promotion. Do not fabricate success. Deferred at user direction.
  - [ ] Exercise cancellation, recovery, stale ownership, target drift, benchmark drift, environment drift, cleanup uncertainty, later mutation invalidation, CLI/JSON/TUI/browser agreement, and exactly one terminal result in the gated run. Deterministic coverage passed; the gated run is deferred at user direction.
  - [ ] Record actual limit use, benchmark stability, evidence size, qualification behavior, optimization convergence, rejected changes, and unresolved risks for Sprint 40 hardening. Real-runtime measurements are deferred at user direction.

## Evidence Checklist

- [x] Missing and disabled policies preserve existing post-execute flow and artifact bytes.
- [x] Project policy tests reject every target-like field and invalid section form.
- [x] Requirements parser tests cover heading context, exact table shape, every cell rule, duplicates, row bounds, stable normalization, and policy disagreement.
- [x] Target packets bind exact requirements and policy digests and contain every row exactly once.
- [x] Prepare and dry-run are deterministic, runtime-free, command-free, and write-free.
- [x] Both V1 descriptor families have exact argv, environment, output identity, parser, coverage, and drift evidence.
- [x] Benchmark authoring changes only approved test measurement scope and freezes before warmup.
- [x] Every measured target has complete warmups, declared samples, raw bounded evidence, qualification, aggregation, and comparison facts.
- [x] Noisy, malformed, non-finite, truncated, drifted, or cleanup-uncertain evidence cannot pass.
- [x] A faster incorrect candidate is rejected and already-met required targets cannot regress.
- [x] Optimization cycles, commands, runtime calls, output, profiles, patches, files, bytes, wall time, retention, and cleanup are bounded and durable across resume.
- [x] Protected paths and evidence resist write, rename, delete, symlink, hard-link, descendant, generator, formatter, repository-control, and Git mutation attempts.
- [x] Private records are versioned, immutable, contained, private, atomic, digest-bound, fenced, retained, and recoverable.
- [x] Cancellation and terminal races commit exactly one product result and one correlated operational terminal observation.
- [x] Later implementation or repair changes make performance stale and block verified or merge-ready state until rerun.
- [x] App, CLI, JSON, TUI, browser, durable run, Markdown, private state, and flow projection use the same product facts.
- [x] Public DTOs, events, and views omit raw private evidence, prompts, secrets, unsafe environment values, and provider payloads.
- [x] Browser and TUI remain accessible, bounded, no-JavaScript capable where applicable, and independent from run ownership.
- [x] Documentation updates are complete and contain no invented performance target.
- [x] Focused tests, full tests, race tests, and build pass.
- [ ] Architecture Review and Sprint Review evidence is recorded; Deep Smoke evidence is deferred at the user's direction.
- [x] Any deviation from `reasoning.md` is recorded and governed before implementation continues.

## Verification Commands

| Check | Command | Expected Result |
| --- | --- | --- |
| Project policy | `go test ./internal/project` | Missing, disabled, enabled, duplicate, unknown, invalid, and override cases pass. |
| Workspace guidance | `go test ./internal/workspace` | Exact optional guidance is generated with no fabricated threshold. |
| Sprint target and performance unit tests | `go test ./internal/sprint` | Parser, packet, limits, descriptors, measurement, optimization, persistence, freshness, flow, cancellation, and recovery tests pass. |
| Application and CLI | `go test ./internal/app` | Use cases, operations, commands, JSON, exits, fencing, restart, and one-terminal-result tests pass. |
| Interface adapters | `go test ./internal/tui ./internal/web` | Canonical performance states, guarded actions, accessibility, security, no-JavaScript, SSE, and observer behavior pass. |
| Process and isolation regression | `go test ./internal/platform/process` | Explicit argv, bounds, cancellation, process-tree cleanup, identity, containment, and disposable-copy behavior remain valid. |
| Focused race checks | `go test -race ./internal/sprint ./internal/app ./internal/runcontrol` | No race, stale writer, duplicate terminal, cancellation, or recovery failure is reported. |
| Full tests | `go test ./...` | All repository tests pass. |
| Full race tests | `go test -race ./...` | All race-enabled tests pass. |
| Build | `go build ./cmd/ultraplan` | The CLI builds successfully. |
| Gated dogfood | Use the command and environment recorded by `system/protocols/deep-smoke-sprint-protocol.md`. | Absolute, baseline-relative, report/baseline, rejected regression, promotion or truthful non-promotion, cancellation, recovery, and interface evidence are recorded. |

## Risks And Blockers

| Risk / Blocker | Source | Mitigation | Status |
| --- | --- | --- | --- |
| Sprint 38 manual repair exit evidence is absent or stale. | `reasoning.md#assumptions-and-risks` | Verify the gate before implementation. Block mutation work if repair, isolation, or promotion foundations are not proven. | closed: the retained proof is verified, production-applied, cleanup-complete, and records the complete ladder |
| Product-owned correctness descriptors cannot be resolved. | `reasoning.md#assumptions-and-risks` | Block admission with the missing policy identity and next action. Never invent fallback commands. | mitigated: deterministic admission fails closed with an actionable blocker |
| Strict Markdown scanning accepts pseudo-target content. | Decision 1 | Use a context-aware scanner and hostile fixtures for fences, comments, quotes, inline code, heading levels, and table association. | mitigated by exact scanners and hostile parser tests |
| Decimal normalization or floating-point use changes a threshold boundary. | Decisions 1 and 4 | Keep canonical decimal text and exact rational threshold arithmetic; use golden equality and sign fixtures. | mitigated by rational comparison and boundary tests |
| Benchmark authoring changes product behavior or correctness authority. | Decisions 3 and 5 | Separate benchmark and production allowlists, inspect hunks and actual diff, run correctness and coverage, then freeze bytes. | mitigated by bounded authoring scope, correctness and coverage checks, and immutable manifest freezing |
| Local noise produces an incorrect verdict. | Decision 4 | Serialize samples, keep complete sets, enforce CV and identity checks, retain dispersion, and return inconclusive instead of selecting favorable evidence. | mitigated by complete-set qualification and conservative outcomes |
| A tiny qualified gain fails to reproduce. | Decision 5 | Require strict direction, final full measurement, all-required regression checks, and no significance claim. Record the residual risk. | mitigated; residual environmental risk is recorded rather than overstated |
| Runtime or descendants write after cancellation. | Decisions 5 and 7 | Cancel the owned tree, reserve bounded cleanup time, recheck source identity and writer fence, and let cleanup uncertainty outrank pass. | mitigated by cancellation, fencing, cleanup reserve, and terminal precedence tests |
| A crash occurs around canonical mutation or summary publication. | Decision 7 | Sync an apply journal before mutation, publish immutable evidence first, verify pointers, and recover only proven later steps. | mitigated by immutable preimages, digest-bound journal reconciliation, resume, and conservative recovery |
| Retention removes current proof. | Decisions 6 and 7 | Protect the current attempt and every digest-bound reference; prune only validated terminal unreferenced attempts. | mitigated by protected-current retention and bounded metadata discovery |
| Operational success is presented as target success. | Decisions 8 and 9 | Keep lifecycle and outcome separate in DTOs and fixture every valid combination. | mitigated by separate DTO fields and canonical fixture coverage |
| Stale performance evidence appears current. | Decisions 7 and 9 | Give freshness priority in gating and presentation; show changed identities, blocking effect, historical reference, and exact rerun action. | mitigated by complete freshness binding and stale presentation tests |
| Similar QA, repair, and performance code becomes a broad framework. | Architecture reasoning | Share only writer-fence, atomic-write mechanics, and presentation primitives. Reject generic stores, workflows, or verdict models in architecture review. | closed by Architecture Review |
| Real dogfood cannot safely accept an improvement. | Decision 10 | Preserve rejected evidence and record the truthful reason. AC-11 permits truthful non-acceptance; it does not permit fabricated promotion. | deferred at user direction; no runtime claim is made |

## Review Inputs

Review should use:

- `sprint-index.md`
- `technical-handbook.md`
- `reasoning/architecture.md`
- `reasoning/api-design.md`
- `reasoning/frontend.md`
- `reasoning.md`
- this `plan.md`
- the implementation diff
- focused, full, race, build, and dogfood evidence
- `docs/plans/integrated-roadmap.md`
- `docs/plans/post-execution-qa-and-repair-loop.md`
- `docs/plans/server-shutdown-run-cancellation-contract.md`
- `docs/plans/retrieval-ready-content-plan.md`
- `/home/antonioborgerees/coding/ultraplan/Aren/docs/phased-roadmap.md`
- `/home/antonioborgerees/coding/ultraplan/Aren/docs/performance-engineering.md`
- `system/protocols/architecture-review-protocol.md`
- `system/protocols/review-sprint-protocol.md`
- `system/protocols/deep-smoke-sprint-protocol.md`

## Execution Log

| Date / Step | Action | Evidence / Notes |
| --- | --- | --- |
| 2026-09-04 / planning | Validated `requirements`, `sprint-index`, `technical-handbook`, and `reasoning`; confirmed the sprint plan stage was ready; generated this plan from the resolved plan prompt. | All four prerequisite validation commands returned `Validation: ok`; sprint status reported `plan` as `ready`. No implementation, review, smoke, Git, or multi-stage execution was started. |
| 2026-09-04 / admission | Verified the recorded Sprint 39 worktree and branch at baseline `37ac236`; confirmed Sprint 38 implementation ancestry and inspected its later manual repair proof. Re-read all six required source plans. | Sprint 38 proof is `verified`, production-applied, cleanup-complete, and complete-ladder. Source digests: integrated roadmap `0d29e9...`, post-execution QA `c19be1...`, shutdown `92ba1b...`, retrieval `34f8dc...`, Aren phased roadmap `519724...`, Aren performance engineering `6b1055...`. No material conflict was found. |
| 2026-09-04 / implementation | Added strict project activation, requirements-owned targets, exact verdict math, verification ordering, private fenced state, closed descriptor discovery, bounded serial measurement, disposable-copy execution, service publication, app operations, CLI, TUI and browser summaries, JSON status/result/evidence routes, configuration, and documentation. | Implementation is in the recorded Sprint 39 worktree. Raw measurements remain private. Disabled projects perform no performance writes. Benchmark commands run in a bounded disposable copy and must leave it unchanged before cleanup and publication. |
| 2026-09-06 / completion | Completed benchmark authoring, bounded optimization and promotion, apply-journal recovery, resume, retention, repair reverification, interface lifecycle and evidence controls, canonical fixtures, and documentation. Ran focused suites, `go test ./...`, `go test -race ./...`, `go build ./cmd/ultraplan`, and `git diff --check`; performed the Architecture Review, Sprint Review, and independent numeric recalculation. | All deterministic checks and both reviews passed. Deep Smoke and real-runtime dogfooding were not run at the user's explicit direction, and no runtime-performance claim is made. |

### Deferred execution scope

Only the gated Deep Smoke protocol and real-runtime dogfooding are deferred, at the user's explicit direction. All implementation, deterministic verification, architecture review, sprint review, and independent numeric review scope is complete. This execution makes no uncollected runtime-performance claim.

## Completion Criteria

- [x] The Sprint 38 dependency gate and all assumptions required for safe mutation are proven or the sprint remains blocked.
- [x] All tasks are complete or explicitly deferred through governed reasoning.
- [x] Every AC and selected contract has direct implementation and evidence coverage.
- [x] Disabled projects preserve the old flow and bytes with no eager performance work.
- [x] Enabled targets trace only to current requirements and use one frozen benchmark identity.
- [x] Qualification, optimization, correctness, limits, persistence, freshness, cancellation, and recovery satisfy the frozen V1 rules in this plan.
- [x] CLI, JSON, TUI, browser, durable operations, canonical artifacts, private state, and flow projection agree.
- [x] Verification commands were run or any deferral is documented with a blocker and next action.
- [x] Gated Deep Smoke and real-runtime dogfood evidence is deferred at the user's explicit direction.
- [x] Deviations from `reasoning.md` were recorded before implementation continued.
- [x] `review.md` evaluates target authority, frozen identity, numeric facts, mutation safety, terminal arbitration, and interface parity without guessing intent.
