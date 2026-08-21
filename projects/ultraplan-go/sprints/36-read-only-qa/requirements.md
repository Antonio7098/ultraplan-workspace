# Sprint 36 Requirements: Read-Only QA Decomposition and Synthesis

> Project: `ultraplan-go`
> Sprint: `36-read-only-qa`
> Purpose: the authoritative, human-readable sprint contract. Later sprint artifacts must reason through the open decisions and satisfy the fixed observable outcomes.

## Sprint Goal

Establish a safe, deterministic, durable, and cross-surface QA capability that decomposes changed behavior into bounded verification surfaces, runs read-only investigations, records falsifiable theory outcomes, and synthesizes the result without generating checks, promoting issues, or mutating production or verification code.

The sprint must preserve the existing analytical `review` workflow and `review.md` artifact while presenting that capability to users as **Conformance Review**. QA is a separate verification phase, not another planning stage and not a renamed form of Conformance Review.

Sprint 36 may be specified and reasoned about now, but implementation and execution are gated on Sprint 35 satisfying its durable run identity and cross-surface observability release criteria. Sprint 36 must reuse that control plane rather than create a second run, progress, cancellation, or recovery authority.

## Required Outputs

| Deliverable | Required Outcome |
|---|---|
| Conformance Review compatibility | Existing `review` commands, integrations, and `review.md` consumers continue to work. User-facing surfaces distinguish Conformance Review from empirical QA, and any additive alias or metadata is compatibility-preserving. |
| Verification phase model | QA and later repair can be represented independently of `PlanningStage`, with explicit lifecycle, terminal outcomes, freshness, cancellation, restart, and recovery semantics. The reasoning and plan stages select the concrete representation and migration path. |
| Versioned QA state contract | Schema-versioned state represents QA maps, attempts, shards, theories, evidence references, cross-shard requests, synthesis, blockers, fingerprints, and next actions. Detailed verification state remains outside `flow-state.json`, which contains only canonical summaries, freshness, verdicts, and pointers. |
| Deterministic QA mapper | Governed inputs produce a stable, validated map of bounded behavioral verification surfaces. Every changed path has one bounded primary surface, relevant boundary behavior is represented, overlap is explicit and bounded, and unchanged inputs reproduce the same map and identifiers. |
| Bounded shard contract | Each shard identifies its behavioral scope, governed inputs, changed paths, relevant requirements and contracts, adjacent tests or known checks, risk tags, exclusions, dependencies, budgets, and completion or blocking conditions. Shards are not created mechanically one-per-file. |
| Read-only investigation | Investigators may inspect approved context, inspect additional repository evidence within policy, and execute explicitly safe non-mutating checks. They cannot write tests, fixtures, probes, production code, verification code, governed artifacts, Git state, or external systems. |
| Theory and evidence model | Investigations persist falsifiable theories, static evidence, commands and bounded results where permitted, refutations, uncertainty, invalid or blocked setup, recommended future checks, and stable fingerprints. A recorded outcome never silently becomes a repairable issue. |
| Global synthesis | A central synthesis step deduplicates related theories, combines evidence, exposes contradictions and cross-shard interactions, retains negative and invalid results, requests bounded follow-up, and produces inspectable next actions without issue promotion or repair. |
| Bounded cross-shard follow-up | Investigators can request focused follow-up across verification surfaces through explicit budgets and lifecycle rules. Cycles, fan-out, duplicated work, and unbounded investigation are prevented or surfaced as blocked. |
| Durable execution and recovery | Mapping, shard investigation, follow-up, and synthesis use Sprint 35 run identity, progress, cancellation, terminal arbitration, replay, and recovery. Completed valid work is resumable when fingerprints still match and invalidated explicitly when governed inputs change. |
| Cross-surface product experience | CLI, JSON, TUI, and browser expose the same QA readiness, map, shard progress, theory outcomes, synthesis status, blockers, cancellation, recovery, freshness, and next actions through shared application operations. |
| Compatibility with current verification | Existing Conformance Review, `smoke`, `smoke.md`, and `verify` behavior remains operational and authoritative for its current contract. Sprint 36 does not absorb smoke or claim evidence-producing QA. |
| Verification and fault injection | Deterministic, race, integration, recovery, and browser tests cover mapping stability, changed-path coverage, read-only enforcement, state validation, cancellation, restart, fingerprint invalidation, cross-shard bounds, synthesis, corrupt or partial state, and cross-surface agreement. |
| Documentation | Architecture, PRD, TRD, CLI/API, local-web, user, and recovery documentation explain the phase boundary, terminology, state ownership, read-only policy, deterministic mapping, theory outcomes, synthesis limits, compatibility, and deferred Sprint 37 capabilities. |

## Acceptance Criteria

- Sprint 35's release gate is satisfied before Sprint 36 implementation begins, and Sprint 36 uses the resulting durable run control plane rather than adding an alternate operation registry or progress authority.
- Existing `ultraplan sprint ... review` behavior and the canonical `review.md` path remain compatible. Supported user-facing surfaces label it Conformance Review and explain that it checks planned-versus-delivered conformance rather than performing empirical QA.
- QA lifecycle is modeled separately from `PlanningStage`; no planning-stage status is overloaded to represent mapping, investigation, synthesis, cancellation, or verification outcomes.
- For byte-identical governed inputs and configuration, repeated mapping produces the same ordered shards, primary-path assignments, relationships, stable IDs, and fingerprints.
- Every changed path from the accepted execute result belongs to exactly one bounded primary verification surface. Intentional secondary coverage and cross-cutting relationships are explicit, bounded, and non-authoritative.
- The mapper uses relevant requirements, code context, reasoning and plan, execute evidence and changed paths, selected contracts, Conformance Review findings, adjacent tests or known checks, risk tags, and package, interface, boundary, or state-transition evidence when available. Missing or invalid required inputs produce a typed blocked or failed result rather than an invented map.
- Mapping is based on coherent behavior and risk. A one-file-per-shard implementation does not satisfy the sprint unless repository evidence proves that each such file is independently the correct behavioral surface.
- Map and shard validation rejects uncovered changed paths, duplicate primary ownership, invalid references, unbounded overlap, cycles that violate the selected contract, unsupported schema versions, and budgets that cannot be enforced.
- Read-only investigators operate under an enforceable non-mutation policy. Any attempted or observed mutation of the implementation repository, verification code, governed sprint artifacts, Git state, or external systems fails closed and is reported durably.
- Allowed checks are demonstrably non-mutating, bounded by time and output, cancellable, and recorded with sufficient identity to reproduce or explain the observation. Merely instructing an agent not to write is not the enforcement mechanism.
- Investigators record each falsifiable theory with its origin, scope, expectation, evidence or refutation, outcome, confidence or uncertainty where applicable, fingerprint, and recommended next action. Supported outcomes distinguish at least confirmed, refuted, inconclusive, invalid, and blocked.
- Static inspection may support or refute a theory, but Sprint 36 does not represent unexecuted generated checks, speculative tests, or unavailable environments as empirical evidence.
- Cross-shard requests are explicit, deduplicated, budgeted, traceable to their originating theory, and guaranteed to terminate as completed, rejected, exhausted, cancelled, invalid, or blocked.
- Global synthesis retains refuted, inconclusive, invalid, and blocked results; it does not discard disagreement to manufacture a clean result. Contradictions and unresolved cross-shard interactions remain visible.
- Synthesis does not promote theories to canonical issues, assign repair eligibility, adjudicate generated evidence, or authorize mutation. Those capabilities remain gated to Sprint 37 or later.
- QA state uses deterministic, schema-versioned identifiers scoped to verification. These identifiers have explicit compatibility and invalidation behavior and do not become the workspace-wide content identity or provenance model.
- Detailed attempts, maps, shards, theories, evidence references, and synthesis do not turn `flow-state.json`, Sprint 35's run repository, browser memory, or event streams into a second source of workflow truth.
- Cancellation stops further mapping, investigation, follow-up, and synthesis work within a bounded interval, arbitrates one durable terminal outcome, and preserves valid completed records without presenting partial work as complete.
- Restart or observer reconnection reconstructs the same durable QA status. Matching completed shard work can resume or be reused only when its governed fingerprint remains valid; stale work is retained for diagnosis and excluded from the current synthesis.
- CLI human output, CLI JSON, TUI, dashboard, sprint detail, and QA detail views agree on readiness, run identity, phase, counts, current work, outcomes, blockers, freshness, cancellation, recovery, and next actions.
- A successful runtime call is insufficient for a successful QA phase: persisted output must satisfy schema, path containment, reference, budget, fingerprint, and lifecycle validation before it is accepted.
- Existing `review`, `smoke`, `smoke.md`, and `verify` compatibility and verdict semantics do not regress. No canonical `qa.md`, generated test or fixture, writable investigation workspace, evidence adjudicator, issue promotion, repair workflow, or automatic repair is introduced.
- `go test ./...`, `go test -race ./...`, `go vet ./...`, `go build ./cmd/ultraplan`, and `git diff --check` pass in `../ultraplan-go`; a gated real-runtime dogfood records produced evidence or the exact unmet prerequisite without claiming a pass from absent evidence.

## Open Questions For Reasoning

The sprint must answer these questions from repository evidence, the Sprint 35 implementation, failure-mode analysis, and focused experiments. The requirements intentionally do not preselect the answers.

1. Which package owns the verification application model, and how do current sprint services, CLI, TUI, and web adapters invoke it without creating a generic workflow engine or embedding product semantics in the runtime layer?
2. What concrete `VerificationPhase` representation preserves current `PlanningStage`, flow-state, command, JSON, and UI compatibility while leaving room for Sprint 37 evidence adjudication and Sprint 38 repair?
3. What is the smallest authoritative state and artifact layout for maps, attempts, shards, theories, evidence references, synthesis, and summaries? Which human-readable output, if any, is appropriate before canonical `qa.md` arrives in Sprint 37?
4. Exactly which inputs govern map identity and freshness, how are accepted execute changed paths obtained, and what truthful behavior applies when Git metadata or an expected planning artifact is unavailable?
5. Which behavioral shard kinds are initially supported, and what deterministic rules map files, packages, interfaces, state transitions, acceptance criteria, cross-cutting risks, and boundaries into primary and secondary coverage?
6. What stable-ID and fingerprint formats are needed for maps, shards, theories, observations, attempts, and synthesis, and how will schema evolution invalidate or migrate them without implying global content identity?
7. What default and configurable limits govern shard count, path overlap, context size, investigator turns, commands, elapsed time, output, cross-shard requests, follow-up depth, and parallelism?
8. Which repository inspections and commands qualify as safe non-mutating checks, and how is non-mutation enforced and audited across filesystem, Git, subprocess, network, environment, and runtime boundaries?
9. Does each investigator receive a snapshot, a protected view, the live implementation repository, or another constrained environment, and how is the selected model proven read-only under failure and hostile tool behavior?
10. What theory schema makes expectations falsifiable, distinguishes static support from empirical evidence, records refutation and uncertainty, and prevents an investigator from promoting its own conclusions?
11. How are duplicate theories, conflicting conclusions, root-cause hypotheses, and cross-shard interactions represented without introducing Sprint 37's adjudicator or issue model early?
12. What bounded algorithm schedules cross-shard follow-up, detects cycles and duplicates, accounts for budgets, and decides when synthesis proceeds with unresolved requests?
13. Which Sprint 35 run, attempt, stage, task, event, cancellation, and recovery identities map to QA maps, shard attempts, and synthesis, and which lifecycle component arbitrates terminal outcomes?
14. What summaries and pointers belong in sprint flow state, what remains in verification-owned state, and how are atomic commits ordered so no surface observes accepted progress before its authoritative record exists?
15. What command and API vocabulary exposes Conformance Review and QA while preserving `review`, existing scripts, routes, JSON clients, and generated skills? Which focused mapping, shard, restart, and JSON operations are required now?
16. How should TUI and browser views present map topology, active and completed shards, theory outcomes, contradictions, follow-up, blockers, freshness, and recovery without deriving independent state or hiding bounded output?
17. How are legacy workspaces, unsupported schema versions, partial writes, corrupt records, stale fingerprints, deleted paths, changed requirements, and interrupted synthesis detected, retained, migrated, or failed closed?
18. Which prompt, repository, command, path, environment, and output data must be redacted or bounded before durable recording and cross-surface display?
19. Which deterministic fixtures, fault-injection cases, and gated real-runtime repositories are sufficient to prove behavioral sharding, read-only containment, useful theory quality, resumability, and UI operability before Sprint 37?

Decisions that materially affect compatibility, authority, safety, determinism, schema evolution, or recoverability must be recorded in reasoning before implementation planning. The plan must include migration, rollback, and fault-injection work rather than treating them as follow-up polish.

## Non-Goals

- Generating or writing tests, fixtures, probes, smoke scenarios, patches, or any other verification code.
- Creating writable or isolated mutation-capable investigation workspaces; that boundary belongs to Sprint 37.
- Canonical `qa.md`, empirical evidence adjudication, expectation-grounding verdicts, flaky-check adjudication, issue promotion, repair eligibility, or regression-candidate classification.
- Absorbing or retiring the existing smoke protocol, changing `smoke.md` authority, or changing canonical-versus-narrow smoke evidence semantics.
- Manual repair, automatic repair, production mutation, issue packets, repair cycles, or re-verification after repair.
- Replacing Conformance Review with QA or removing the `review` command or `review.md` compatibility path.
- Workspace-wide content IDs, revision-aware provenance, amendment workflows, supersession, citation identity, or the gated content contract after Sprint 39.
- Repository indexing, retrieval, RAG, embeddings, a knowledge graph, alternate authored-content persistence, cloud authority, Aren integration, or hosted multi-user operation.
- Building a general-purpose scheduler, DAG engine, issue tracker, distributed queue, or provider-specific QA subsystem.
- Automatic Git mutation, network mutation, package installation, database migration, service deployment, or other externally stateful investigation actions.

## Constraints

- `internal/sprint` or a clearly justified adjacent product package owns QA semantics, mapping, theory policy, synthesis, artifacts, and state. `internal/platform/runtime` remains generic and must not import sprint or QA semantics.
- Sprint 35's durable run system owns operational run identity, lifecycle observation, replay, cancellation routing, and reconciliation. Verification-owned state owns QA domain outcomes; neither duplicates the other's authority.
- CLI, TUI, and HTTP are adapters over shared typed application operations. `internal/web` must not map shards, invoke investigators, synthesize theories, infer completion, or persist alternate QA truth.
- Governed input selection, ordering, normalization, IDs, fingerprints, mapping, and synthesis must be deterministic for unchanged inputs and schema/configuration.
- The mapper must cover every accepted changed path with exactly one primary behavioral surface, keep secondary relationships bounded, and reject invalid or unbounded maps before investigation starts.
- Non-mutation is a technical boundary with fail-closed enforcement and post-attempt verification, not solely a prompt instruction. Any uncertainty about containment is a blocked outcome.
- Investigators remain individually bounded, cancellable, and non-authoritative. They may report theories and observations but cannot decide canonical issue status or authorize follow-up beyond the central budgeted policy.
- Durable state writes use atomic semantics, path containment, schema validation, reference validation, bounded diagnostics, redaction, and explicit recovery. Runtime success never substitutes for artifact validity.
- Verification-scoped IDs and schemas must be explicitly versioned and migratable. They must not pre-commit the later workspace-wide content identity model.
- Detailed QA state remains outside `flow-state.json`; shared run events and browser projections are observations, not canonical theory or synthesis storage.
- Existing Conformance Review, smoke, and verify behavior remains available and authoritative for its current scope throughout this sprint.
- Normal tests are deterministic and offline. Real-runtime dogfood is gated, non-mutating, and truthful about unavailable prerequisites.
- The roadmap owns Sprint 36 scope. Evidence-producing QA, smoke integration, issue promotion, repair, content identity, retrieval, persistence, graph, cloud, and Aren work may not be pulled into this sprint.

## Dependencies

### Release Prerequisite

- Sprint 35, `projects/ultraplan-go/sprints/35-durable-run-observability/review.md`, must record a passing release verdict with its durable identity, cross-process observation, replay, cancellation, reconciliation, retention, and recovery guarantees satisfied before Sprint 36 implementation starts.
- The implemented Sprint 35 run repository, lifecycle service, event journal, application operations, and interface projections are the operational foundation for QA runs and shard attempts.

### Planning Inputs

- Project scope and sprint ordering:
  - `projects/ultraplan-go/roadmap.md`
  - `projects/ultraplan-go/docs/PRD.md`
  - `projects/ultraplan-go/docs/TRD.md`
  - `projects/ultraplan-go/docs/ARCHITECTURE.md`
- Current implementation-repository plans, re-read during reasoning rather than copied mechanically from this requirement:
  - `../ultraplan-go/docs/plans/integrated-roadmap.md`
  - `../ultraplan-go/docs/plans/post-execution-qa-and-repair-loop.md`
  - `../ultraplan-go/docs/plans/sprint-code-context-stage.md`
  - `../ultraplan-go/docs/plans/server-shutdown-run-cancellation-contract.md`
- Sprint 34's accepted shared context, current execute evidence and changed-path records, current Conformance Review findings and compatibility behavior, and existing smoke/verify contracts.
- Current CLI, JSON, TUI, local-web, application-service, runtime, cancellation, persistence, lock, recovery, redaction, and path-containment implementations in `../ultraplan-go`.

## Review Expectations

Review must verify the QA contract as a safe product capability, not merely the presence of new commands or successful agent calls:

- trace one QA run from governed input selection through deterministic mapping, shard attempts, theory records, cross-shard follow-up, synthesis, durable terminal arbitration, and every supported surface
- regenerate a map from unchanged inputs and compare raw ordered IDs, fingerprints, scopes, primary changed-path ownership, and relationships
- exercise multi-file behavior, interface boundaries, state transitions, cross-cutting risk, tests adjacent to changed code, and deliberately misleading file boundaries to prove behavioral rather than mechanical sharding
- attempt filesystem, Git, generated-test, external-system, package-manager, and governed-artifact mutations and prove that enforcement fails closed, records the violation, and leaves authoritative state unchanged
- inspect confirmed, refuted, inconclusive, invalid, and blocked theories plus contradictory and duplicate findings; prove that none silently becomes an issue or repair action
- exhaust cross-shard budgets, introduce cycles and duplicates, cancel during follow-up, and prove bounded termination and truthful synthesis
- interrupt and restart mapping, an investigator, follow-up, and synthesis at each persistence boundary; verify atomic recovery, matching-fingerprint reuse, stale-state invalidation, and one terminal outcome
- compare CLI human output, CLI JSON, TUI, dashboard, sprint detail, and QA detail for identical readiness, identity, counts, phase, progress, outcomes, blockers, freshness, recovery, and next actions
- corrupt and partially write state, change schema versions and governed inputs, remove referenced paths, expire retained events, and verify typed failure or recovery without fabricated success
- prove existing `review`, `review.md`, `smoke`, `smoke.md`, and `verify` compatibility and show that canonical `qa.md`, writable investigation, evidence adjudication, issue promotion, and repair remain absent
- inspect redaction, output bounds, command records, diagnostics, and the gated real-runtime evidence or exact blocker
- review the release logs for `go test ./...`, `go test -race ./...`, `go vet ./...`, `go build ./cmd/ultraplan`, and `git diff --check` in `../ultraplan-go`
