> **Inputs Used:** `projects/ultraplan-go/sprints/37-evidence-qa-smoke/requirements.md`, `projects/ultraplan-go/sprints/37-evidence-qa-smoke/sprint-index.md`, `projects/ultraplan-go/sprints/37-evidence-qa-smoke/technical-handbook.md`, `projects/ultraplan-go/docs/ARCHITECTURE.md`, `projects/ultraplan-go/docs/PRD.md`, `projects/ultraplan-go/docs/TRD.md`, `studies/go-cli-study/reports/final/01-project-structure.md`, `studies/go-cli-study/reports/final/02-command-architecture.md`, `studies/go-cli-study/reports/final/03-dependency-injection.md`, `studies/go-cli-study/reports/final/04-configuration-management.md`, `studies/go-cli-study/reports/final/05-error-handling.md`, `studies/go-cli-study/reports/final/06-io-abstraction.md`, `studies/go-cli-study/reports/final/07-state-context.md`, `studies/go-cli-study/reports/final/08-concurrency.md`, `studies/go-cli-study/reports/final/10-logging-observability.md`, `studies/go-cli-study/reports/final/11-testing-strategy.md`, `studies/go-cli-study/reports/final/12-extensibility.md`, `studies/go-cli-study/reports/final/13-security.md`, `studies/go-cli-study/reports/final/14-performance.md`, `system/protocols/architecture-review-protocol.md`, `system/protocols/review-sprint-protocol.md`, `system/protocols/deep-smoke-sprint-protocol.md`

# Architecture reasoning

## Area Decisions

### Scope and invariants

Sprint 37 extends the existing read-only QA workflow into an evidence-producing workflow. It does not create a second execution system. The design must preserve these invariants:

- the original target is never an admitted write location for an analyzer or evaluator attempt;
- one durable operation owns a QA run, and only its current fenced owner may publish state;
- detailed evidence has one sprint-owned authority under `verification/`;
- `qa.md` is a bounded assessment generated from accepted evidence, not a raw transcript store;
- `flow-state.json` is a bounded status and pointer projection, not an evidence database;
- `qa --suite smoke` invokes the existing smoke authoring, discovery, execution, validation, and evidence path;
- CLI, TUI, and web code do not own QA policy or alternate state;
- generated patches are retained evidence and are never applied to the target.

The architecture favors narrow additions to the current package graph. A generic sandbox, workflow engine, evaluator plugin system, or second smoke runner would add abstractions without another proven caller.

### Package ownership and dependency direction

The package boundary is:

```text
cmd / internal/tui / internal/web
                  |
                  v
             internal/app
              /        \
             v          v
    internal/sprint  internal/runcontrol
             |
             v
 internal/platform/process
```

The arrows show allowed feature-level dependencies. Existing lower-level storage, model, configuration, and runtime dependencies remain behind their current package boundaries.

- `internal/platform/process` owns reusable mechanics: private workspace creation, explicit argv and cwd launch, environment filtering, process-group tracking, bounded output capture, timeout and cancellation propagation, descendant termination, workspace change capture, and cleanup. It returns facts. It does not import sprint packages, know requirement IDs, decide whether evidence is valid, promote issues, publish QA state, or select smoke scope.
- `internal/sprint` owns QA policy and evidence semantics: immutable evidence plans, attempt identity, writable-attempt admission, evidence validation, adjudication, issue deduplication and promotion, assessment construction, `qa.md`, detailed verification records, and translation of smoke results into QA evidence links. Isolation policy stays here even though process mechanics live below it.
- `internal/runcontrol` remains the only operational authority. It owns durable operation identity, ownership, fencing, cancellation, replay, recovery, terminal arbitration, and stale-owner rejection. Sprint code cannot create a competing lock or terminal-state mechanism.
- `internal/app` owns typed use-case sequencing. It resolves and validates inputs, starts or resumes the managed operation, delegates QA rules to `internal/sprint`, invokes the existing smoke use case when selected, and publishes through the current state stores while holding the runcontrol fence. It does not reimplement evidence rules.
- CLI and TUI adapters call the same `internal/app` QA use case and render its typed events and result. They may collect confirmation and presentation preferences, but they cannot select different adjudication or isolation behavior.
- `internal/web` may query projected status through typed application services. Any future QA mutation endpoint must also call an `internal/app` use case. It must not import `internal/sprint`, call a model runtime or process runner, parse CLI output, or write `qa.md`, `verification/`, or flow state itself.

This split keeps policy near the sprint domain while preventing sprint-specific meaning from leaking into the reusable process package. It also preserves one application path for every adapter.

### Frozen run identity and evidence plan

Before the first attempt, the application layer resolves the effective project, sprint, target, configuration, selected suite, runtime/model identities, and managed-operation identity. Sprint QA then builds an immutable evidence plan. The plan records:

- a plan ID and parent operation ID;
- the governed input fingerprint and target identity/scope;
- requirement and claim IDs to be tested;
- analyzer and evaluator attempt IDs, roles, and independence groups;
- the runtime/model identity expected for each attempt;
- whether each attempt may write, its timeout, and its evidence limits;
- approved deterministic commands and smoke coverage references when applicable;
- claim weights and non-negotiable pass gates used by the final assessment.

No adapter or attempt may add a claim after execution starts. A retry gets a new attempt ID linked to the same plan item. A changed target, governed input, effective configuration, runtime identity, or required smoke mapping invalidates the affected result instead of silently changing the plan.

Configuration is resolved once and passed as a validated value. Workers do not reread configuration. Secrets and unrelated host environment values are not copied into an attempt unless the plan explicitly admits them for that attempt.

### Attempt isolation

Every attempt declared writable receives a new private workspace. Workspaces are never reused across analyzers, evaluators, retries, or smoke tests. The process package performs these steps:

1. Create a run-owned temporary directory with restrictive permissions.
2. Materialize the frozen target snapshot into a fresh workspace without relying on Git worktree support.
3. Launch the admitted executable with an explicit argument vector, the workspace as cwd, a filtered environment, bounded streams, a context deadline, and process-group ownership.
4. On completion or cancellation, terminate descendants, wait for exit, record bounded process facts, compute workspace changes against the materialized snapshot, and remove the workspace under a short independent cleanup context.
5. Return cleanup status and retained change evidence to sprint orchestration.

The original target path is not placed in attempt prompts, argv, cwd, or environment. Before execution, after each attempt, and before final publication, the application compares the live target identity with the frozen identity. Any mismatch fails target-integrity admission, prevents the affected evidence from becoming trusted, and blocks final publication. UltraPlan reports the mismatch and does not try to restore or overwrite the target.

This is deliberate containment, not a hostile-code security sandbox. A child process with ambient host filesystem access may still discover absolute paths. The architecture prevents ordinary cwd-relative writes from reaching the target and detects target drift before trust or publication. If a supported runtime cannot provide these guarantees, writable attempts are unavailable and fail closed. The sprint will not claim container-grade or kernel-grade isolation without a separate requirement and threat model.

Workspace changes are normalized into evidence records. A generated patch is stored with its attempt identity, base target identity, changed-path list, digest, truncation status, and validation result. No QA path exposes an apply operation.

### Sequential scheduling and cancellation

Attempts run sequentially in Sprint 37. This is a correctness choice, not a performance accident. Sequential execution makes workspace ownership, process cleanup, evidence ordering, budget accounting, and stale-owner behavior inspectable before concurrency is introduced.

One caller context controls normal work. Cleanup uses a bounded context derived independently after cancellation so descendant termination, evidence closure, and workspace removal still run. Runcontrol cancellation wins terminal arbitration. A worker that finishes after ownership loss may retain diagnostic records, but fencing prevents it from publishing `qa.md`, verification indexes, or flow state.

Parallel execution remains deferred until tests prove independent workspaces, bounded resource use, deterministic result collation, complete descendant cleanup, and unchanged adjudication outcomes under concurrency.

### Evidence records and admission

Each analyzer, evaluator, deterministic check, and delegated smoke run produces an immutable attempt envelope. The envelope includes the plan and operation IDs, attempt and parent IDs, role, runtime/model identity, target and input fingerprints, workspace identity when present, start and finish times, outcome, bounded output references, cleanup result, mutation summary, evidence digests, and any issue observations.

Sprint QA admits an envelope only when all of these checks pass:

- its schema and required fields are valid;
- its plan, operation, target, input, runtime/model, and attempt identities match the frozen plan;
- its evidence paths remain inside declared verification or external smoke roots;
- writable work used the assigned fresh workspace;
- process exit, timeout, cancellation, cleanup, and truncation facts are explicit;
- the original target still matches its frozen identity;
- the current runcontrol owner still holds the publish fence.

Malformed, stale, mismatched, escaped, incomplete, or unfenced evidence remains diagnostic and is excluded from adjudication. Admission failures are typed and retain their cause. They are not converted into product findings.

### Deterministic adjudication

Adjudication is a pure sprint-domain operation over the frozen plan and admitted envelopes. It does not inspect live files, rerun tools, or depend on adapter behavior.

- An issue observation carries a stable claim ID, issue class, normalized location, severity, summary, supporting evidence references, and source attempt.
- Deduplication groups observations by claim, issue class, and normalized location. Similar wording alone does not merge unrelated claims.
- Every raw observation remains in its attempt record. The canonical issue links all supporting and dissenting observations, so deduplication never destroys provenance.
- Promotion requires support from at least two admitted, independent passes. Distinct retries from one failed workspace or one copied result do not count as independent.
- Claims that require evaluator judgment use a fixed odd number of at least three independent evaluator attempts. A strict majority of admitted outcomes is required. Missing, tied, or invalidated outcomes produce `inconclusive`, never an implied pass.
- Among evaluator-approved variants in one issue group, the canonical issue retains the highest valid severity. Evidence completeness breaks equal-severity ties, followed by stable attempt order. Dissent is retained.
- A finding that does not meet promotion rules remains visible as an unpromoted observation but cannot drive the final defect list or a passing claim.

The assessment computes its score from the frozen claim weights and accepted claim statuses. Mandatory integrity, identity, cleanup, and evidence-completeness gates override the numeric score. The assessment records the score, rationale, pass/fail decision, inconclusive claims, promoted issues, and links to every supporting detailed record. Model prose may explain the result, but deterministic product code chooses promotion status, score, and verdict.

### Persistence and publication

Detailed records live under a run-specific directory below the sprint's `verification/` tree. The run index points to plans, attempts, checks, issues, patches, adjudication, and assessment records by relative path and digest. Records use atomic replacement where a single file is mutable; append-like history uses distinct immutable files rather than repeated in-place edits.

Publication is a fenced transaction at the application boundary:

1. Validate that required plan work has a terminal admitted result.
2. Write and validate the detailed assessment and run index.
3. Render and atomically replace `qa.md` from that assessment.
4. Atomically update the bounded QA projection in `flow-state.json` through the existing state owner.

A complete assessment may publish either pass or fail. Infrastructure interruption, target drift, missing majority, incomplete cleanup, stale ownership, or invalid evidence does not replace the last complete `qa.md` with a partial artifact. Diagnostic attempt records may remain under `verification/` for recovery and audit.

The flow-state projection contains only bounded fields such as QA status, verdict, operation/run ID, target and governed-input fingerprints, assessment path and digest, summary counts, timestamps, and the next action. It does not contain prompts, raw output, patches, full issue bodies, or copied smoke evidence. Existing fields and status readers remain compatible; Sprint 37 adds optional fields rather than creating a second state format.

### One smoke execution path

`ultraplan sprint <project> <sprint> qa --suite smoke` delegates to the same typed smoke application service used by the canonical smoke command. Delegation includes the current review gate, manifest resolution, suite authoring rules, machine discovery, coverage validation, scope selection, confirmation data, environment policy, explicit argv execution, result validation, issue linkage, mutation checks, and verdict computation.

QA does not parse terminal output or call a CLI handler. It receives a structured smoke result containing run identity, target identity, selected coverage, counts, evidence paths, issue IDs, mutation summary, and verdict. QA validates the identity and stores links plus a bounded summary in its own attempt record. Raw streams and detailed smoke results remain in the external harness, which stays authoritative for smoke evidence.

If the canonical review gate, manifest, required coverage, environment, or harness evidence is unavailable, the smoke plan item is blocked. QA cannot bypass the gate, substitute normal verification commands, author a private suite, or call a provider for an offline product path. Narrow reruns follow the smoke protocol and do not establish closure without a passing containing required suite.

### Adapter compatibility

The current CLI flags, exit behavior, JSON field meanings, TUI controls, and web navigation remain stable. New evidence and pointer fields are additive. Human-readable output may gain progress and evidence references but cannot remove established fields or change their meaning.

CLI and TUI use the same request, progress event, cancellation, recovery, and result types from `internal/app`. JSON rendering is a projection of those typed results. The web adapter uses the same application query path for status and never reads detailed evidence to derive a competing verdict.

### Verification seams

Tests follow the ownership boundaries:

- pure unit tests cover identity validation, independence checks, deduplication, majority outcomes, strongest-finding retention, scoring, and gate overrides;
- process integration tests cover fresh workspace creation, target drift detection, patch retention, bounded output, timeout, cancellation, descendant cleanup, permissions, and cleanup failure;
- runcontrol tests cover stale writers, replay, cancellation races, and terminal arbitration;
- application tests cover publication order, atomic replacement, partial-run recovery, and adapter parity;
- smoke delegation contract tests prove that QA calls the canonical typed smoke service and cannot substitute discovery or verdict logic;
- CLI/TUI/web tests prove thin rendering and forbidden dependency directions.

Interfaces are introduced only at the existing filesystem, process, runtime, clock, state-store, and smoke-service boundaries where deterministic fault injection needs them. Concrete sprint policy remains concrete.

## Trade-Offs

| Decision | Benefit | Cost | Why it is accepted |
|---|---|---|---|
| Filesystem snapshot per writable attempt | Works for dirty or non-Git targets and gives every attempt a fresh base | Copying and diffing can be expensive | Correctness and target independence matter more than Sprint 37 throughput |
| Sequential attempts | Simple ownership, cleanup, budget, and evidence ordering | Longer QA runs | Concurrency is not justified until isolation and adjudication are proven |
| Fail-closed target fingerprint checks | Prevents mutated or stale targets from producing trusted evidence | Concurrent legitimate target edits invalidate a run | Requiring a rerun is safer than attributing or repairing an unknown mutation |
| Detection rather than a claimed hostile-code sandbox | Avoids an unearned cross-platform security framework | Does not stop a malicious child with ambient host access | Sprint 37 needs trustworthy evidence and accidental-mutation containment, not execution of untrusted binaries under a new threat model |
| Immutable detailed records plus bounded projections | Preserves auditability without bloating flow state or `qa.md` | More files and explicit indexing | Evidence provenance cannot fit safely in one mutable summary file |
| Deterministic adjudication with model-authored rationale only | Reproducible promotion, scoring, and verdicts | Requires explicit schemas and validation | A model cannot be the authority for whether its own evidence counts |
| Reuse of canonical smoke service | One behavior, evidence schema, and safety policy | QA must obey smoke readiness and review gates | Bypassing those gates would make `qa --suite smoke` incompatible and less trustworthy |
| Additive persistence and output changes | Existing automation and saved sprint state keep working | New optional fields require careful default handling | CLI output and persisted flow state are concrete external compatibility obligations |

Rejected alternatives:

- **Git worktrees as the isolation boundary.** They fail for non-Git targets, dirty state, unavailable Git, and write paths outside a worktree. A Git repository may optimize snapshot creation later, but correctness cannot depend on it.
- **A subprocess and temporary cwd as a security sandbox.** Process groups control lifecycle, not filesystem authority. The design states the weaker guarantee honestly and fails target-integrity admission on drift.
- **One shared workspace with reset between attempts.** Reset correctness is hard to prove after cancellation, descendant leakage, or cleanup failure. Fresh workspaces make independence observable.
- **Parallel attempts in the first release.** Parallelism would complicate operation fencing, cleanup, output bounds, budget enforcement, and deterministic ordering before those invariants have tests.
- **A generic sandbox, workflow engine, evaluator registry, or plugin framework.** Sprint 37 has one concrete QA workflow. These abstractions would hide policy and expand the security boundary without another demonstrated use.
- **Adapter-owned orchestration or state.** Separate CLI, TUI, or web paths would drift in cancellation, evidence admission, and verdict behavior.
- **A second smoke runner or copied smoke evidence.** This would split discovery, safety policy, issue history, and evidence authority.
- **Model-selected issue promotion or final verdict.** It would make retries and model changes alter product state without a deterministic rule.
- **Applying generated patches.** Sprint 37 is verification, not repair. Applying a patch would mutate the target and collapse the boundary planned for Sprint 38.
- **Automatically restoring a changed target.** The system cannot safely distinguish its own write from a concurrent user edit. It reports drift and stops.

## Evidence

### Project and sprint contracts

- `projects/ultraplan-go/sprints/37-evidence-qa-smoke/requirements.md` requires fresh isolation for writable attempts, independent analyzer and evaluator records, majority adjudication, issue promotion only after independent support, bounded canonical assessment, target identity checks, retained-but-unapplied patches, runcontrol ownership, and smoke reuse.
- `projects/ultraplan-go/sprints/37-evidence-qa-smoke/technical-handbook.md` identifies the existing read-only QA path, strict state publication, durable writer fencing, process-group cleanup, and manifest-driven smoke as the extension points. Its main warning is that the missing layer is evidence semantics and isolation, not another command framework.
- `projects/ultraplan-go/docs/ARCHITECTURE.md` places typed use cases in `internal/app`, keeps CLI/TUI/web as adapters, assigns durable operations to runcontrol, and requires smoke execution to remain manifest-driven and external-evidence-backed.
- `projects/ultraplan-go/docs/PRD.md` fixes the scope at evidence-producing QA and smoke compatibility while explicitly deferring repair, broad refactors, and generic frameworks.
- `projects/ultraplan-go/docs/TRD.md` requires one state authority, fail-closed ownership and freshness, bounded state projections, cancellation and descendant cleanup, and no direct runtime/process execution from web handlers.

### Architecture and execution reports

- `studies/go-cli-study/reports/final/01-project-structure.md` and `studies/go-cli-study/reports/final/02-command-architecture.md` support thin adapters, library-first feature logic, shared typed execution paths, and CLI/TUI parity without handler reuse.
- `studies/go-cli-study/reports/final/03-dependency-injection.md` supports explicit composition and narrow capability seams. This is why process, clock, state, runtime, and smoke dependencies are injected at boundaries while adjudication stays concrete.
- `studies/go-cli-study/reports/final/04-configuration-management.md` supports resolving and validating one immutable effective configuration before workers start. It also supports redacted provenance rather than worker-side config reads.
- `studies/go-cli-study/reports/final/05-error-handling.md` supports typed failures, aggregation of partial results, actionable context, and preserving the last valid artifact when a later run fails.
- `studies/go-cli-study/reports/final/06-io-abstraction.md` supports narrow filesystem and process boundaries, explicit argv, in-memory fault injection for policy tests, and real integration tests for OS behavior.
- `studies/go-cli-study/reports/final/07-state-context.md` supports caller-context propagation, detached but bounded cleanup, atomic writes, operation-scoped ownership, and explicit partial state.
- `studies/go-cli-study/reports/final/08-concurrency.md` supports bounded scheduling, cancellation-aware workers, deterministic collation, and postponing parallel work until cleanup and race behavior are proven.
- `studies/go-cli-study/reports/final/10-logging-observability.md` supports operation, attempt, target, and runtime correlation IDs plus separation of diagnostic logs from durable evidence authority.
- `studies/go-cli-study/reports/final/11-testing-strategy.md` supports pure policy tests, fault-injected boundary tests, subprocess integration tests, race tests, and golden coverage for stable artifacts.
- `studies/go-cli-study/reports/final/12-extensibility.md` warns that subprocesses are not a sandbox and argues against speculative plugins and engines. It supports a narrow process capability used by one concrete sprint workflow.
- `studies/go-cli-study/reports/final/13-security.md` supports explicit argv, path containment, restrictive permissions, environment allowlisting, bounded output, target identity checks, process-group cleanup, and redaction.
- `studies/go-cli-study/reports/final/14-performance.md` supports bounded streams and work queues, incremental evidence handling, measured optimization, and concurrency limits instead of unbounded fan-out.

### Review and smoke protocols

- `system/protocols/architecture-review-protocol.md` requires visible workflows, earned complexity, controlled state changes, clear side effects, narrow coupling, and tests that match package responsibilities. These checks favor the ownership split above over a generic engine.
- `system/protocols/review-sprint-protocol.md` establishes runcontrol-style frozen scope, deterministic verdict authority, contained evidence, mandatory reviewer completeness, atomic replacement, and stale-fingerprint rejection. QA follows the same architectural discipline without taking ownership of review.
- `system/protocols/deep-smoke-sprint-protocol.md` makes the external harness authoritative for detailed smoke evidence and fixes the review gate, authoring, machine discovery, scope, execution, mutation, issue, and rerun rules. QA therefore delegates instead of duplicating any part of that pipeline.

## Risks

| Risk | Required response | Residual risk |
|---|---|---|
| A child discovers and writes an absolute host path | Remove target references from admitted inputs, use a private cwd and filtered environment, compare target identity before trust and publication, and fail closed | This is not hostile-code containment; stronger OS isolation needs a separate threat model |
| Target fingerprinting is expensive or noisy for large dirty trees | Reuse the canonical target identity/scope rules, hash incrementally where current stores allow it, and measure before optimizing | Legitimate concurrent edits still invalidate a run |
| Cancellation leaves descendants or workspaces | Own process groups, use bounded independent cleanup, record cleanup status, and block admission when cleanup is incomplete | Platform-specific process behavior needs integration and smoke evidence |
| Attempt outputs or patches exhaust memory or disk | Bound captured streams, record truncation, stream digests and patch creation, cap retained evidence, and reject incomplete evidence where full content is required | Limits may need tuning from measured runs |
| Two observations are correlated rather than independent | Freeze independence groups in the plan and require distinct attempt and workspace identities | Model-provider correlation cannot be eliminated, only made explicit |
| Majority voting hides a severe minority finding | Retain dissent and raw observations; an inconclusive or integrity-gated claim cannot pass; choose the strongest evaluator-approved variant | A valid novel finding may remain unpromoted until another independent pass supports it |
| A stale worker publishes after cancellation or takeover | Require the current runcontrol fence for every canonical write and verify it again at final publication | Bugs below runcontrol remain high severity and require race coverage |
| `qa.md`, detailed records, and flow state diverge | Generate summaries from the validated assessment, store digests and relative pointers, and publish in fenced order | Multi-file publication is not one filesystem transaction, so reconciliation must detect an interrupted sequence |
| QA smoke behavior drifts from canonical smoke | Depend on the typed smoke service and add a contract test that rejects private discovery, execution, or verdict code | Changes to smoke result schemas require coordinated adapter updates |
| Existing CLI, JSON, or persisted state consumers break | Keep current meanings and add optional fields; test old state fixtures and established JSON output | Unknown external scripts may rely on undocumented prose formatting |
| The Sprint 36 admission gate is incomplete | Treat the current validated gate as mandatory input and fail closed if it is missing, stale, or inconsistent | Sprint 37 cannot make unvalidated prior evidence trustworthy after the fact |

Assumptions fixed for this sprint:

- writable attempts are sequential;
- an odd evaluator count of at least three is required wherever majority decides a claim;
- target drift blocks trust and publication regardless of who caused it;
- detailed smoke evidence remains external;
- repair and patch application remain Sprint 38 work;
- unsupported isolation or cleanup behavior fails closed instead of silently weakening guarantees.

Open implementation questions are limited to proof, not architecture choice:

- Each supported operating system must demonstrate descendant cleanup and private-workspace deletion. Until it does, writable QA is unsupported on that platform.
- The existing target resolver must demonstrate a stable identity for dirty and non-Git targets. If it cannot, QA blocks rather than inventing a Git-only fallback.
- The canonical smoke application service must expose a structured invocation seam that preserves its review gate. If it does not, Sprint 37 adds that narrow seam inside `internal/app`; it does not duplicate smoke orchestration.
- Evidence size limits must be measured against representative runs. Regardless of the chosen values, limits remain explicit, validated, and recorded in the frozen plan.
