> **Inputs Used:** `projects/ultraplan-go/project-index.md`, `projects/ultraplan-go/sprints/39-performance-stage/requirements.md`, `projects/ultraplan-go/sprints/39-performance-stage/sprint-index.md`, `projects/ultraplan-go/sprints/39-performance-stage/technical-handbook.md`, `projects/ultraplan-go/docs/ARCHITECTURE.md`, `projects/ultraplan-go/docs/PRD.md`, `projects/ultraplan-go/docs/TRD.md`, `system/reasoning/architecture_reasoning_template.md`, `system/contracts/core/architecture.md`, `system/contracts/core/security.md`, `system/contracts/runtime/workflows.md`, `system/contracts/runtime/performance.md`, `system/contracts/runtime/persistence-and-migrations.md`, `studies/go-cli-study/reports/final/01-project-structure.md`, `studies/go-cli-study/reports/final/02-command-architecture.md`, `studies/go-cli-study/reports/final/03-dependency-injection.md`, `studies/go-cli-study/reports/final/04-configuration-management.md`, `studies/go-cli-study/reports/final/05-error-handling.md`, `studies/go-cli-study/reports/final/06-io-abstraction.md`, `studies/go-cli-study/reports/final/07-state-context.md`, `studies/go-cli-study/reports/final/08-concurrency.md`, `studies/go-cli-study/reports/final/09-terminal-ux.md`, `studies/go-cli-study/reports/final/10-logging-observability.md`, `studies/go-cli-study/reports/final/11-testing-strategy.md`, `studies/go-cli-study/reports/final/12-extensibility.md`, `studies/go-cli-study/reports/final/13-security.md`, `studies/go-cli-study/reports/final/14-performance.md`, `docs/plans/integrated-roadmap.md`, `docs/plans/post-execution-qa-and-repair-loop.md`, `docs/plans/server-shutdown-run-cancellation-contract.md`, `docs/plans/retrieval-ready-content-plan.md`, `/home/antonioborgerees/coding/ultraplan/Aren/docs/phased-roadmap.md`, `/home/antonioborgerees/coding/ultraplan/Aren/docs/performance-engineering.md`, `internal/project/domain.go`, `internal/project/index.go`, `internal/sprint/domain.go`, `internal/sprint/verification_phase.go`, `internal/sprint/flow.go`, `internal/sprint/state.go`, `internal/sprint/qa_types.go`, `internal/sprint/qa_state.go`, `internal/sprint/qa_repair.go`, `internal/platform/process/process.go`, `internal/platform/process/isolation.go`, `internal/app/operations.go`, `internal/app/operation_runner.go`

# Architecture: Requirements-Driven Performance Stage

This area decides ownership, dependency direction, orchestration, mutation, persistence, freshness, and recovery for the optional `performance` verification phase. It does not choose transport DTO details, UI layout, or the statistical constants that the implementation must define and test.

## Area Decisions

### 1. Extend the sprint module after one narrow ownership refactor

The current module architecture fits the feature. `internal/project` already owns project-index semantics, `internal/sprint` owns verification and product verdicts, `internal/platform/process` owns bounded process execution and disposable copies, `internal/runcontrol` owns durable operational runs, and `internal/app` is the shared interface boundary. Performance should not become a new top-level product package.

Use focused files in `internal/sprint` for targets, benchmark descriptors, measurement, optimization, and state. Keep parsing, comparison, admission, workflow transitions, and terminal outcomes concrete inside that package. Add narrow injected collaborators only for volatile work: runtime proposal generation, process execution, isolation, clocks, environment capture, and durable writer validation.

One small refactor should happen before performance persistence is added. Replace the QA-named writer token and service fence with a verification-owned token carrying the existing run ID, operational attempt ID, and fencing generation. QA, repair, and performance then consume the same ownership fact. Do not make their stores or record schemas generic. Their state changes for different product reasons even though they use the same fence and atomic-write mechanics.

The reusable persistence code should stay unexported in `internal/sprint` and cover mechanics only: contained no-follow path resolution, private directories and files, strict versioned JSON decoding, immutable create, same-directory atomic replacement, digest calculation, directory sync, and a fence recheck immediately before rename. `QAStore` and the new performance store retain separate path builders, validators, publication methods, and recovery rules.

### 2. Keep activation and target authority in different owners

`internal/project` gains a typed `PerformancePolicy` with `disabled` and `enabled` values. `ProjectIndex` carries the parsed policy and its source digest. Missing policy is exactly `disabled`. The project parser rejects duplicate sections, unknown keys, invalid modes, and target-like fields. It never parses or stores target rows.

`internal/sprint` owns a strict `Performance Targets` parser and deterministic normalizer. The parser is a Markdown-aware single pass that tracks fenced blocks, HTML comments, and blockquote context before accepting a real level-two heading. It accepts only the immediately associated table, validates the exact ordered header and every v1 cell, then produces rows in stable ID order while preserving normalized decimal text for digest and display. A loose regular expression over the full document is not sufficient because it can accept examples or quoted content.

Policy-aware requirements validation joins the two facts at the sprint service boundary. The content parser accepts bytes and a typed policy value; it does not discover or reread the project index. This permits the required mismatch errors without allowing project policy to supply targets. The resulting target packet binds every normalized row exactly once to the requirements digest, policy digest, parser version, and packet digest.

Operational configuration is resolved separately. Built-in defaults, workspace settings, and environment settings may lower finite limits, but the resolved values and sources form an operational-policy fingerprint. If a requirements-owned sample count exceeds an effective operational limit, admission blocks and reports both values. It never rewrites the target packet.

### 3. Add performance to verification order, not planning order

Add `VerificationPhasePerformance` before `VerificationPhaseConformanceReview` in `VerificationPhases()`. Do not add a `PlanningStage`, compatibility stage, or another flow engine. The verification orchestrator remains the single owner of this order:

```text
successful execute
  -> performance when enabled
  -> Conformance Review
  -> QA and smoke compatibility
  -> repair when admitted
```

For a missing or disabled policy, the transition remains `execute -> Conformance Review` and performs no performance target parse beyond mismatch validation, benchmark discovery, runtime construction, command execution, state creation, or artifact write. Existing artifact bytes remain unchanged.

For an enabled policy, execute readiness and current target-worktree identity are prerequisites. Only `passed` and `passed_with_reports` permit Conformance Review. A required target that is missed, inconclusive, or blocked stops later verification. Report-only and baseline outcomes stay visible but do not block when every required target passes.

`FlowState` gains one optional bounded performance projection and a schema migration that leaves old states semantically disabled. The projection contains current phase, freshness, outcome, bounded target counts, cancellation, state and result pointers with digests, current performance attempt ID, blocker, and next action. Detailed targets, samples, profiles, patches, cycles, and command results remain outside `flow-state.json`.

### 4. Freeze one attempt before measurement

One performance attempt owns an immutable chain of identities. Admission resolves all required facts before durable acceptance:

1. Enabled project policy and current validated requirements.
2. Normalized target packet.
3. Successful execute evidence and exact target-worktree identity.
4. Product-owned correctness command set.
5. Process and environment policy.
6. Effective lower-only limits and remaining budgets.
7. Isolation capabilities.
8. Runtime and model selection for proposal work only.

After durable acceptance, benchmark discovery may find an existing mapping or ask the runtime to propose missing benchmark files in a disposable copy. Product code validates and, when needed, promotes that benchmark-only patch. It then freezes the benchmark manifest before the first warmup. The manifest maps every target to exactly one closed measurement descriptor and binds benchmark paths and digests, parser version, explicit argv, working directory, environment keys, timeout, output identity, raw measurement unit, and correctness coverage.

The frozen attempt envelope is the target packet, benchmark manifest, correctness policy, process policy, environment policy, initial worktree identity, effective limits, and workflow version. Resume loads these values from private state. It accepts no replacement from flags, configuration, current files, or runtime output. Drift blocks resume rather than silently starting a different experiment under the old attempt ID.

### 5. Use closed benchmark descriptors and product-owned measurement

Sprint 39 supports two descriptor families only: standard Go benchmark output and the versioned UltraPlan JSON measurement envelope. Product code maps descriptor types to parser functions and `process.Request` values. Requirements text and runtime output can name a scenario, but neither can provide executable argv, a parser, a working directory, environment values, or a timeout.

Warmups and measured samples run through `internal/platform/process` with explicit argv, a fixed contained directory, an allowlisted environment, finite timeout, bounded stdout and stderr, context cancellation, and process-tree cleanup. Retain bounded raw sample facts before aggregation. Product code validates output identity and units, qualifies variance and environment compatibility, computes aggregates, applies inclusive comparators, and derives each target verdict. Runtime prose is evidence for a hypothesis only.

Measured samples and candidate comparisons run serially by default under one attempt to avoid adding scheduler noise to the result. Pure parsing, target preparation, digesting, and bounded status reads may run concurrently where they do not touch shared mutable state. Callers cannot select measurement parallelism.

### 6. Separate benchmark promotion from implementation promotion

There are two mutation paths, and they must not share an allowlist.

Benchmark authoring occurs before baseline. A fresh disposable copy may change only validated benchmark, test, fixture, or narrowly approved measurement-support paths. Product code proves complete and unambiguous target coverage, runs benchmark correctness checks, verifies the actual diff, and applies the patch through the existing product-owned mutation boundary. Once promoted, those bytes become part of the frozen manifest and cannot change during the attempt.

Optimization begins only from a qualified miss that permits optimization. Each cycle selects one target-linked symptom, records one bounded profile and hypothesis, and gives the runtime one fresh disposable copy based on the current accepted implementation identity. The runtime cannot write the canonical worktree. Optimization proposals may change only approved production paths. They cannot change requirements, policy, benchmark or test files, fixtures, expected outputs, correctness commands, parser code, configuration, prior evidence, planning artifacts, Git control files, or protected roots.

Product code accepts a candidate only after it verifies patch syntax, link and path safety, actual changed paths and bytes, current source identity, complete cleanup, the full correctness gate, a qualified improvement for the selected target, and no regression of an already-met required target. The product rechecks canonical identity immediately before apply and records an apply journal. Overlapping user changes block promotion and remain untouched.

An accepted optimization advances the attempt's implementation identity. It invalidates measurements from earlier candidate identities and any existing Conformance Review, QA, smoke, or repair evidence. Execute task completion remains historical provenance, not proof of the final bytes. The attempt then measures the new canonical identity with the same frozen benchmark manifest. A later mutation outside this attempt makes the terminal performance result stale.

### 7. Publish private evidence first and summaries last

The performance store owns `verification/performance-state.json` and `verification/performance/attempts/<attempt-id>/...`. Attempt evidence is immutable except for explicitly named current pointers and recovery journals. Every current pointer includes the referenced record path and digest. Status reads the bounded current state and verifies its pointers; it does not reconstruct current status by scanning every attempt directory.

Each publication step checks the writer fence before writing and again immediately before rename. Terminal publication follows this order:

1. Write immutable cycle evidence and `result.json`.
2. Write canonical `performance.md` from the immutable result facts.
3. Atomically replace `performance-state.json` with pointers to the exact result and report digests.
4. Atomically update the bounded performance projection in `flow-state.json`.
5. Commit the run-control terminal observation.

Cross-file atomicity is not claimed. The order makes partial publication detectable. Recovery revalidates the fence, frozen inputs, source identity, immutable record digests, apply journal, and cleanup facts. It may idempotently complete a missing later publication from a valid immutable result. It cannot infer `passed` from a report, process exit, absent owner, or partially written state. If source or frozen identities changed, or cleanup cannot be proved, recovery publishes `blocked`, `interrupted`, or `cleanup_uncertain` as applicable.

Run control remains authoritative for durable acceptance, leases, cancellation routing, ordered operational events, and exactly one operational terminal observation. The sprint module remains authoritative for performance phase and target outcomes. Operational success means a valid product result was committed; it does not mean the result was `passed`.

### 8. Make freshness a dependency identity, not a timestamp

The current performance result binds these facts:

- requirements and project-policy digests;
- execute evidence and final implementation identity;
- target packet and benchmark manifest digests;
- parser, sample, aggregation, correctness, process, and environment policy identities;
- workflow version and terminal result digest.

Freshness comparison reports stable reason codes for each changed dependency. It does not use modification time as authority. A later repair or implementation mutation marks performance stale before verified or merge-ready state can be reported. Repair reverification runs functional checks first, then the applicable frozen performance targets against the repaired identity. Pre-repair samples cannot be reused as current evidence.

The performance attempt also invalidates downstream evidence at the exact successful mutation publication boundary, not when a runtime merely proposes a patch. Benchmark promotion before baseline and implementation promotion during a cycle both update the worktree identity and mark earlier downstream verification stale. A rejected proposal changes no canonical freshness.

### 9. Keep all expensive and retained work finite

Product-owned hard maxima apply to benchmark authoring calls, runtime attempts, warmups, samples, commands, command time, stdout and stderr bytes, profiles, patches, changed files, changed bytes, cycles, wall time, retained attempts, and cleanup. Effective settings can only reduce them. Counters and the absolute deadline are durable and never reset on resume.

Private records retain bounded numeric facts and digests rather than full provider payloads or unlimited command output. Retention prunes only validated old attempt directories, never the current attempt or evidence referenced by the current terminal result. Limit exhaustion produces `stalled`, `target_missed`, or `blocked` according to the proven facts. It cannot produce a pass.

No generic plugin system, parser registry exposed to users, distributed benchmark service, historical database, cache, fleet scheduler, or Aren integration is justified. Aren's roadmap keeps persistence, policies, daemon hosting, workflows, and broader execution types gated by Aren's own evidence. UltraPlan should reuse its current local capabilities without pulling those later Aren phases forward.

## Trade-Offs

| Decision | Benefit | Cost | Rejected alternative |
| --- | --- | --- | --- |
| Performance stays in `internal/sprint` | Keeps target semantics, flow, mutation, and verdicts beside sprint state | The sprint package grows and needs focused files | A top-level `internal/performance` package would split one product workflow across modules. |
| Generic writer ownership, separate product stores | Removes duplicated fencing knowledge while preserving distinct QA and performance schemas | Requires a small QA/repair rename before feature work | Reusing `QAStore` for performance would make paths and recovery semantics misleading; fully generic verification storage would add a broad framework. |
| Service-level join of project policy and sprint targets | Reports enabled/disabled mismatches without creating a second target source | Pure content validation needs an explicit policy argument | Reading project state inside the table parser would hide dependencies and make the parser hard to test. |
| Optional flow-state projection | Old states remain valid and status stays bounded | Readers and migrations must handle one more optional projection | Putting full performance state in `flow-state.json` would make it an attempt database. |
| Closed descriptor families | Product code can audit every command and parser | New measurement formats require a versioned code change | Runtime-authored commands or a plugin registry would grant ungoverned execution authority. |
| Serial measured samples | Produces more comparable environment evidence | Multi-target runs take longer | Parallel measurement is faster but can contaminate the values used for target verdicts. |
| Two separate promotion allowlists | Benchmark authority cannot be weakened by an optimizer | Promotion code must classify the proposal type explicitly | One broad patch policy would let optimization rewrite its measuring instrument. |
| Immutable evidence followed by pointer publication | Recovery can identify the last proven boundary | A crash can leave valid private evidence not yet visible in summaries | Treating several files as if they were one transaction would be false on the current filesystem model. |
| Digest-based freshness | Every stale result names the changed authority | Computing bounded source identities has a real cost | Timestamp freshness is cheaper but cannot prove content identity. |
| Profile one miss per cycle | Keeps hypotheses and effects attributable | It may miss an optimization that improves several targets at once | Broad redesign makes scope, convergence, and regression attribution hard to bound. |

## Evidence

- `projects/ultraplan-go/sprints/39-performance-stage/requirements.md:11-67` separates project activation, requirements-owned targets, sprint orchestration, focused performance files, and private attempt evidence. `requirements.md:129-220` fixes admission, benchmark identity, qualification, mutation, outcomes, freshness, cancellation, recovery, and interface authority.
- `projects/ultraplan-go/sprints/39-performance-stage/technical-handbook.md:29-53` supports module-owned policy, thin interfaces, explicit dependencies, frozen validation, typed outcomes, bounded process work, propagated cancellation, bounded retention, and profile-led changes. `technical-handbook.md:112-152` identifies the unresolved architecture pressures addressed above.
- `projects/ultraplan-go/docs/ARCHITECTURE.md:227-299` assigns activation to `internal/project` and the complete performance protocol to `internal/sprint`. `ARCHITECTURE.md:303-327`, `ARCHITECTURE.md:370-388`, and `ARCHITECTURE.md:719-768` separate app adapters, operational run control, authored state, detailed verification state, and later persistence choices.
- `system/contracts/core/architecture.md:58-161` requires explicit module ownership, inward dependencies, and narrow ports only at real volatile seams. `architecture.md:246-268` forbids product policy in generic platform packages.
- `system/contracts/runtime/workflows.md:67-151` requires inspectable terminal states, explicit cancellation, protected retries, reconciliation of partial effects, and versioned long-running workflow behavior.
- `system/contracts/runtime/performance.md:66-99` requires measured optimization and finite expensive work. `performance.md:139-178` requires owned, cancellable concurrency and visible runtime cost drivers.
- `system/contracts/runtime/persistence-and-migrations.md:52-105` requires clear state ownership, explicit compatibility, atomic replacement, and reconciliation. `persistence-and-migrations.md:125-159` requires bounded status reads and declared freshness for derived summaries.
- `system/contracts/core/security.md:97-136` requires strict boundary parsing and structured subprocess APIs. `security.md:179-238` requires contained paths, private storage, and restricted dynamic execution.
- `studies/go-cli-study/reports/final/01-project-structure.md:32-40` and `studies/go-cli-study/reports/final/02-command-architecture.md:32-38` support thin adapters and one-way delegation to product logic. `studies/go-cli-study/reports/final/06-io-abstraction.md:32-38` supports injected side-effect boundaries rather than direct process or filesystem access from presenters.
- `studies/go-cli-study/reports/final/03-dependency-injection.md:32-58` supports manual composition, narrow interfaces, and avoiding package globals. `studies/go-cli-study/reports/final/04-configuration-management.md:32-73` supports explicit precedence, validation after source resolution, and immutable effective configuration.
- `studies/go-cli-study/reports/final/05-error-handling.md:32-38` supports typed product outcomes and separate user and operational rendering. `studies/go-cli-study/reports/final/09-terminal-ux.md:80-101` and `studies/go-cli-study/reports/final/10-logging-observability.md:32-38` support interruptible progress, non-TTY behavior, structured fields, and stdout/stderr separation without moving authority into the interface.
- `studies/go-cli-study/reports/final/07-state-context.md:32-48` and `studies/go-cli-study/reports/final/08-concurrency.md:32-46` support one cancellation lineage, localized concurrency, explicit joins, and bounded cleanup. Those findings argue against detached measurement workers and adapter-owned cancellation.
- `studies/go-cli-study/reports/final/11-testing-strategy.md:32-38` supports table tests, command-level workflows, fakes, and behavior-focused assertions. The architecture therefore exposes process, isolation, clock, environment, persistence, and runtime seams without interfaces for pure target arithmetic.
- `studies/go-cli-study/reports/final/12-extensibility.md:32-44` records the cost and security trade-offs of plugin systems. Its subprocess evidence supports containment, not a Sprint 39 extension ecosystem. `studies/go-cli-study/reports/final/13-security.md:32-40,113-125` supports explicit trust boundaries and argument arrays.
- `studies/go-cli-study/reports/final/14-performance.md:32-38,89-111` supports lazy setup, bounded streaming, bounded concurrency, and incremental work. It does not justify pools, caches, or parallel measurements without a measured bottleneck.
- `docs/plans/post-execution-qa-and-repair-loop.md:447-468,632-698,756-811` provides the existing separation of isolated writable work, private detailed state, bounded flow summaries, fingerprints, and sprint-owned verification behavior. Performance reuses these guarantees without merging its verdict with QA.
- `docs/plans/server-shutdown-run-cancellation-contract.md:42-145,181-207` requires canonical idempotent cancellation, bounded cleanup, conservative restart reconciliation, and one terminal outcome under races.
- `docs/plans/integrated-roadmap.md:91-169,460-486,692-715` preserves filesystem authority, shared app use cases, truthful observability, durable run identity, and the gated placement of performance work after repair foundations.
- `docs/plans/retrieval-ready-content-plan.md:94-136,603-677` requires authoritative Markdown and revision-aware evidence while forbidding a premature retrieval or universal artifact abstraction. Digest-bound performance records satisfy the current need.
- `/home/antonioborgerees/coding/ultraplan/Aren/docs/phased-roadmap.md:53-128,273-413` requires behavior before abstraction, explicit lifecycle ownership, cancellation race tests, and exactly one terminal outcome. `/home/antonioborgerees/coding/ultraplan/Aren/docs/performance-engineering.md:30-55,175-313,435-469,595-700` requires baseline-first evidence, workload-shaped benchmarks, bounded retention, profile-led optimization, and correctness after every change.
- `internal/sprint/verification_phase.go:5-52` already separates verification phases from planning stages and keeps only shipped compatibility mappings. `internal/sprint/domain.go:136-146` and `internal/sprint/qa_types.go:257-277,658-670` show the existing bounded flow projection pattern.
- `internal/sprint/qa_state.go:16-136,239-292` already proves private bounded storage, no-follow path checks, retention, and digest-bound references. `internal/sprint/qa_types.go:199-221` and `internal/sprint/qa_repair.go:1407-1416` show the QA-named writer ownership that should become verification-generic before performance reuses it.
- `internal/platform/process/process.go:21-151` provides explicit argv, fixed directories, bounded output, timeout, cancellation, and cleanup facts. `internal/platform/process/isolation.go:18-136,139-261` provides finite disposable-copy creation, source identity, contained execution, diffing, and cleanup without assuming Git.
- `internal/app/operations.go:80-166` and `internal/app/operation_runner.go:115-205` provide the existing guarded-operation and fenced-runner seams. Performance should extend these seams while leaving product verdicts in `internal/sprint`.

## Risks

| Risk | Consequence | Required control |
| --- | --- | --- |
| Policy-aware validation becomes an alternate target source | Project config can weaken or replace sprint commitments | Keep project policy typed as activation only; construct the packet solely from current requirements bytes and bind both digests. |
| Markdown scanning accepts examples or quoted targets | Ungoverned rows become executable commitments | Track fences, comments, heading level, quote context, exact headers, and immediate table association; test hostile Markdown fixtures. |
| Writer generalization expands into a verification framework | Sprint 39 refactors working QA and repair behavior unnecessarily | Generalize only token, fence, and mechanical record-write code; keep workflows, stores, schemas, and outcomes separate. |
| Benchmark authoring changes product behavior | The baseline measures a runtime-authored implementation change | Use a benchmark-only allowlist, inspect the actual diff, run coverage and correctness checks, then freeze all benchmark bytes before warmup. |
| Optimization changes its own measuring instrument | A candidate manufactures an improvement | Deny benchmark, test, fixture, parser, target, sample-policy, and correctness-policy paths during optimization. |
| Environment or benchmark drift is treated as a target miss or pass | Comparisons look precise but are not comparable | Bind environment and manifest identities to every sample set; classify drift and excessive variance as inconclusive or blocked. |
| Parallel work contaminates measurements | Scheduler and resource contention changes verdicts | Serialize warmups and measured samples under one stability policy; parallelize only non-measuring work with explicit bounds. |
| A late process or runtime writes after cancellation | Cancelled or cleanup-uncertain work becomes passed | Fence every publication, recheck before rename and promotion, and let one immutable product result precede run-control terminal commitment. |
| Crash occurs between promotion and evidence publication | Canonical source changes with no trustworthy explanation | Persist and fsync the apply journal before mutation, record before and after identities, and reconcile conservatively before any resume. |
| Cross-file terminal publication is mistaken for a transaction | Partial state or Markdown appears current | Publish immutable evidence first and summaries last; verify all pointer digests and complete only idempotent later steps during recovery. |
| Freshness invalidation misses repair or external source changes | Old measurements permit verified or merge-ready state | Compare content identities at status and gate boundaries and centralize post-mutation invalidation for performance, review, QA, smoke, and repair. |
| Status polling scans all retained samples | CLI, TUI, or browser reads become slow and memory-heavy | Read one bounded current state and digest-bound pointers; expose old evidence through bounded queries, not recursive reconstruction. |
| Retention deletes current proof | A terminal verdict can no longer be audited or recovered | Protect the current attempt and every record referenced by current state or result; prune only validated old attempt directories. |
| Disabled projects perform eager work | Existing flows change behavior or bytes without opt-in | Branch on typed policy before target parsing beyond mismatch validation, runtime setup, process calls, store creation, or flow projection writes. |
| Performance work pulls future Aren capabilities into UltraPlan | The sprint adds daemons, remote workers, or generic execution policy | Keep all work local and same-host through current runtime, process, isolation, run-control, and app boundaries. |

No architecture question remains open for final sprint reasoning. Concrete type names and qualification constants may be selected during final reasoning and planning, but they must preserve the ownership, ordering, identity, mutation, persistence, freshness, recovery, and boundedness decisions above.
