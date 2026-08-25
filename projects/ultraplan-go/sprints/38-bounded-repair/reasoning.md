# Sprint Reasoning

> **Inputs Used:** `projects/ultraplan-go/sprints/38-bounded-repair/requirements.md`, `projects/ultraplan-go/project-index.md`, `projects/ultraplan-go/roadmap.md`, `projects/ultraplan-go/docs/ARCHITECTURE.md`, `projects/ultraplan-go/docs/PRD.md`, `projects/ultraplan-go/docs/TRD.md`, `projects/ultraplan-go/sprints/38-bounded-repair/sprint-index.md`, `projects/ultraplan-go/sprints/38-bounded-repair/technical-handbook.md`, `projects/ultraplan-go/sprints/38-bounded-repair/reasoning/api-design.md`, `projects/ultraplan-go/sprints/38-bounded-repair/reasoning/architecture.md`, `projects/ultraplan-go/sprints/38-bounded-repair/reasoning/frontend.md`

## 1. Executive Summary

### Decisions in This Sprint

Sprint 38 will add repair as a sprint-owned `VerificationPhase` protocol. It will consume exactly one current adjudicated `QAIssue`, freeze all authority into an immutable packet, display that packet, durably accept one operation, persist a single-use confirmation bound to that acceptance, and only then permit mutable work. `internal/sprint` will own admission, scope, progress, reverification, cleanup interpretation, proof qualification, and semantic outcomes. Run control will continue to own operational acceptance, liveness, fencing, cancellation, and operational terminal state. Platform code will return process, filesystem, isolation, and identity facts without deciding whether a repair succeeded.

The runtime will edit a private isolated copy only. Product code will derive and strictly parse a bounded text patch, publish it as evidence, recheck target identity and ownership, and apply it through direct contained filesystem operations. It will never ask the runtime, Git, a shell, or a formatter to mutate the production checkout. Protected governance, evidence, state, test, configuration, and Git-control paths remain non-authorizable.

Reverification will be sequential: exact reproducer, affected primary shards, linked theory checks, boundary and follow-up shards, containing QA and smoke suites, then an independent Conformance Review delta. A narrower non-pass skips wider gates. Cleanup is part of the result, not a best-effort epilogue.

Manual mode will invoke this cycle once. Automatic mode will be an opt-in loop over the same cycle, available only after a current real manual proof. Product-owned lower-only budgets, deterministic progress facts, and stop rules will prevent another mutation when work stalls, repeats, widens, drifts, becomes uncertain, or needs a product decision. Operational lifecycle and semantic outcome will remain separate public facts.

CLI, JSON, TUI, and browser will use one app contract and the existing guarded durable-operation path. They will expose bounded canonical DTOs, separate packet review and confirmation, explicit manual-proof gating, canonical refresh after event loss, and the same cancellation and recovery behavior. No adapter will read private repair records or derive success.

### Confidence Assessment

- **High confidence:** package ownership, immutable packet and confirmation authority, isolated proposal, product-owned apply, strict persistence, fixed reverification order, one shared manual/automatic protocol, lifecycle/outcome separation, and adapter parity. These conclusions agree across all three area-reasoning artifacts, the project documents, and the selected contracts.
- **Medium confidence:** a crash-safe multi-file apply on an ordinary filesystem, the practical strictness of the supported patch subset, and the exact manual-proof fingerprint balance. The authority rules are final, but failure injection and real-runtime evidence must prove the implementation details.
- **Unresolved architecture decisions:** none. The final choice for packet preparation is a synchronous durable `repair-prepare` operation: it obtains writer-fenced ownership, publishes or reuses the packet without a goroutine or runtime, and terminates immediately. This closes the only implementation choice left open by API area reasoning.
- **Recommendation:** proceed to `plan.md`. Planning must decompose these decisions and evidence gates; it must not revisit ownership, confirmation order, apply authority, automatic-mode reuse, or result semantics.

## 2. Problem Framing

### Problem Statement

Sprint 37 can promote a current evidence-backed issue, but no current path may turn that issue into a bounded production change. Sprint 38 must add that mutation without allowing a model, adapter, stale worker, or damaged state file to expand scope, weaken the governing expectation, bypass confirmation, repeat a committed apply, or declare its own success.

### Why It Is Non-Trivial

Repair joins authorities that intentionally remain separate: QA decides whether an issue is current and eligible; a human authorizes one exact request; run control owns execution; the sprint service owns repair meaning; platform code owns mechanics; Conformance Review owns its verdict; and four interfaces observe the same run. A safe implementation must preserve every boundary while still providing one understandable workflow.

The mutable portion is small but crash-sensitive. Proposal generation, patch validation, production apply, actual-diff comparison, progressive checks, process cleanup, lock release, immutable result publication, and flow projection each create a failure boundary. A normal filesystem cannot atomically replace several production files as one transaction. Recovery therefore needs an apply journal and compensation evidence, and uncertainty must escalate rather than trigger another write.

Automatic mode adds no new authority, but it adds persistent budgets, absolute deadlines, progress normalization, repeated-patch detection, issue-set comparison, reopening counts, resume boundaries, and proof invalidation. Those facts must survive restart and cannot be inferred from model explanations or transient events.

### Constraints

- Repair stays in `internal/sprint` and remains separate from `PlanningStage`, QA adjudication, Conformance Review, and smoke authority.
- `internal/platform/process` supplies generic copy, identity, explicit executable/argument execution, cancellation, comparison, and cleanup mechanics only.
- `internal/runcontrol` owns durable operation acceptance, owner generation, liveness, cancellation, event order, and operational terminal arbitration only.
- The runtime writes only inside a private isolated copy. Product code alone may mutate the production checkout.
- Every mutable run requires one current packet and one single-use confirmation bound to packet, request, mode, limits, target, governed inputs, durable run identity, attempt, and fencing generation.
- Tests, requirements, acceptance criteria, governed planning inputs, QA or repair evidence, flow state, workspace configuration, and Git control are protected even when an issue claims they should change.
- Detailed records remain strict and private under `verification/`; `flow-state.json` contains only a bounded summary and digest-bound pointer.
- Every budget has a finite product default and immutable maximum. Workspace and `ULTRAPLAN_QA_*` configuration may lower values only.
- Browser, TUI, CLI, and JSON consume app DTOs. Events remain observational and browser disconnect remains independent from cancellation.
- No Git mutation, plugin system, alternate persistence authority, general issue tracker, remote worker, or content/retrieval system enters this sprint.

### Hidden Complexities

- Packet facts are distributed across QA issue, evidence, adjudication, map, shards, theories, review, smoke, check catalog, governed inputs, policy, and target records. Matching IDs without validating all parent digests would freeze inconsistent authority.
- Confirmation must include durable acceptance identity, so acceptance occurs before confirmation publication while dispatch remains blocked. A persistence failure can strand an accepted but unconfirmed operation and must terminate it without work.
- A valid proposal can become stale between isolated generation and apply. Full target identity and writer ownership must be checked before proposal, before apply, after apply, and after reverification.
- Pre-existing unrelated target changes belong to the frozen identity. Repair reports them but never absorbs, cleans, stages, or reverts them.
- Path containment alone does not protect tests, generated evidence, repository control, formatter side effects, links, or process descendants. Scope policy must combine path class, file identity, link checks, pre-image digests, actual-diff comparison, and process containment.
- A committed apply followed by interruption is not a resumable proposal boundary. Resume must continue from evidence and never apply that patch again.
- Cleanup timeout proves only that waiting stopped. Verification requires affirmative process-tree termination, workspace removal, compensation state, and lock disposition.
- A semantic `failed` repair can be an operationally completed run because the protocol produced its authoritative result. Conversely, operational completion cannot imply `verified`.
- A manual proof file can exist but still be stale because code, schema, maxima, policy, governed inputs, isolation capabilities, check catalog, runtime policy, or target identity changed.
- Automatic retention can prune old cycle bodies while observers hold cursors. Public queries must expose retention gaps rather than imply complete history.

### Assumptions

- Sprint 37 supplies a current repair-eligible issue, accepted evidence, exact reproducer, affected shard and theory links, containing smoke, and current QA/adjudication fingerprints. Missing facts block packet preparation.
- Sprint 35 run control can be refactored to separate durable accept/claim from dispatch without weakening existing operations.
- Existing isolation and process mechanisms can reject links, compare bounded trees, terminate process trees, and prove workspace removal. Unsupported capability blocks repair.
- Direct contained filesystem replacement can stage bytes and preserve pre-images without invoking Git. Multi-file crash atomicity is not assumed.
- The existing Conformance Review capability can return a bounded focused delta while retaining independent verdict authority and leaving `review.md` unchanged.
- Public app, CLI, TUI, and web compatibility permits additive repair resources and optional fields; strict private repair records receive independent schema versions.
- Product code can compute stable fingerprints for executable code/protocol identity, schema, policy, maxima, governed inputs, approved checks, runtime policy, isolation capabilities, and target identity.

### Risks

| Risk | Consequence | Mitigation and required evidence |
|---|---|---|
| Inconsistent packet join | An issue gains stale or unintended authority | Validate every parent identity and digest as one snapshot; reject missing reproducer, scope, check, review, smoke, or target facts; add drift and cross-attempt tests |
| Accepted operation lacks confirmation | A run appears accepted without mutation authority | Split acceptance from dispatch; terminalize confirmation publication failure; restart tests must prove no worker or runtime starts |
| Multi-file apply is interrupted | Production may contain a partial patch | Persist an apply journal and pre-images, compensate in process, compare actual target state, and escalate uncertain crash recovery without retrying apply |
| Link or side-effect escape | Protected files can change despite allowed patch paths | Reject symlinks, hard links, special files, unsupported renames, formatter/shell execution, and actual changed paths outside the packet; use real filesystem attack tests |
| Fence and mutation lease disagree | A stale owner can publish or another owner can race terminalization | Check both controls at every applicable boundary, add `terminalizing`, compare ownership before release, and race-test stale owners and recovery |
| Reverification command coverage is incomplete | Repair cannot prove the result | Require immutable approved descriptors during packet admission; block rather than let the runtime invent a command |
| Cleanup is mistaken for a bounded wait | A verified result can hide live processes or locks | Require affirmative cleanup facts for each resource; map uncertainty to escalation and disqualify manual proof |
| Automatic progress is model-defined | Repeated or widening work can continue indefinitely | Normalize issue sets, severity, exact failure, patch digest, target identity, and new facts in product code; persist limits before scheduling another cycle |
| Manual proof is too weak or too fragile | Unsafe automatic admission or unnecessary invalidation | Version named fingerprint components, report component-level mismatches, and verify invalidation with real and fixture-based proof tests |
| Adapter or event drift | Interfaces show different authority or infer success | Use one canonical app fixture and status reload; test event gaps, late events, reconnect, restart, hostile text, and all terminal combinations |

## 3. Inputs and Source Grounding

### Sprint Index Summary

The sprint index selects all contracts needed for a production mutation: architecture, CLI, configuration, documentation, errors, LLM runtime/evaluation, observability, performance, persistence, security, testing, and workflows. It selects architecture, API, and frontend reasoning because the repair authority chain crosses package ownership, shared operations, and operator confirmation. It also requires Architecture Review, Sprint Review, and Deep Smoke Sprint evidence.

The exclusions are equally binding. Repair will not add a second smoke or review engine, general issue management, Git mutation, test or expectation weakening, broad refactoring, unbounded automation, content identity, alternate product persistence, hosted operation, or remote repair.

### Technical Handbook Findings Applied

- Keep one product path behind thin adapters. Repair uses one sprint service, one app contract, and one durable runner rather than adapter-specific commands or handlers.
- Validate authority before expensive or mutable work. Packet admission, effective limits, durable acceptance, confirmation, freshness, ownership, and proof gating all precede proposal or apply.
- Use typed semantic and operational facts. Closed blockers, outcomes, recovery actions, and lifecycle values survive wrapping and are never derived from message text.
- Propagate cancellation to work, then use a separate finite cleanup context. Cleanup uncertainty remains uncertainty.
- Localize goroutine ownership and keep reverification sequential. The shared runner is the only repair launch site.
- Treat subprocess isolation as a trust boundary, not a success authority. Runtime output is untrusted evidence; product code derives scope, applies bytes, and decides outcomes.
- Separate canonical status from progress streams. Reconnect and terminal events trigger a bounded status reload.
- Merge and validate lower-only configuration before confirmation, and persist effective values and safe source labels.
- Combine table tests, real command-path tests, failure injection, selective public-output fixtures, and gated real-runtime proof.

### Repos Studied / Source Evidence Used

| Repository or report | Why it mattered | Decisions influenced |
|---|---|---|
| `01-project-structure` with Helm `pkg/cmd/install.go` / `pkg/action/install.go` and gdu `cmd/gdu/app/app.go` | Shows thin command and presentation layers over owned behavior | D-01, D-09 |
| `02-command-architecture` with gh-cli factories, Helm command construction, and rclone wrappers | Shows parsing and command setup separated from execution lifecycle | D-02, D-09 |
| `03-dependency-injection` with gh-cli `pkg/cmdutil/factory.go` and restic repository/backend interfaces | Grounds explicit runtime, process, store, clock, and fence seams while warning against a god container | D-01, D-04 |
| `04-configuration-management` with go-task precedence, restic explicit flags, and k9s validation | Grounds merged effective lower-only limits and source-aware validation before work | D-06 |
| `05-error-handling` with rclone behavioral errors, restic fatal errors, and gh-cli exit classes | Grounds closed blockers, stable error categories, and separate CLI exit mapping | D-07, D-09 |
| `06-io-abstraction` with gh-cli test streams and restic terminal/filesystem interfaces | Grounds failure injection for persistence, process, output, and adapter behavior | D-04, D-10 |
| `07-state-context` with Helm signal cancellation and restic bounded cleanup lifetime | Grounds root cancellation followed by a distinct finite cleanup context | D-08 |
| `08-concurrency` with restic `errgroup`, k9s worker bounds, and OpenCode timed shutdown | Grounds one launch site, explicit waits, bounded cleanup, and no fire-and-forget repair work | D-04, D-08 |
| `09-terminal-ux` with chezmoi non-TTY prompts, lazygit stoppable views, and restic terminal progress | Grounds explicit non-interactive confirmation and observational progress | D-09 |
| `10-logging-observability` with k9s structured keys, restic backend logging, and Helm diagnostics | Grounds stable run/packet/cycle correlation and separation of safe public output from diagnostics | D-07, D-09 |
| `11-testing-strategy` with chezmoi txtar, gh-cli acceptance tests, and restic mocks | Grounds layered unit, integration, contract, race, and real-path evidence | D-10 |
| `12-extensibility` with Helm subprocess runtime/version dispatch and rclone registry warnings | Supports strict schema dispatch and a fixed runtime boundary while rejecting a repair plugin registry | D-01, D-03 |
| `13-security` with OpenCode permission gates, chezmoi private temporary directories, and restic redaction | Grounds visible confirmation, private isolation, explicit argv, path controls, and redaction | D-02, D-03, D-09 |
| `14-performance` with gh-cli lazy setup, lazygit streaming, restic bounded file work, and yq streaming | Grounds runtime-free preparation, finite budgets, incremental records, paging, and bounded retention | D-06, D-09 |

The reports and repositories are evidence, not governing contracts. Sprint requirements, project documents, selected contracts, and shipped Sprint 35-37 authority boundaries take precedence.

### Area Reasoning Summary

- `reasoning/architecture.md` selects a focused `internal/sprint` protocol, strict verification persistence, accept-before-confirm-before-dispatch sequencing, isolated proposal generation, product-derived and product-applied patches, two ownership controls, sequential reverification, one shared cycle, manual-proof fingerprinting, terminalization, and conservative recovery.
- `reasoning/api-design.md` selects a separate repair app contract, repair-specific bounded read resources, generic guarded durable operations for writes, explicit manual/automatic request shapes, digest-bound paging, stable errors and exits, idempotent transitions, and separate lifecycle/outcome fields.
- `reasoning/frontend.md` selects extension of the existing sprint QA presentation, a read-only packet review followed by guarded confirmation, always-visible manual-proof gating, canonical refresh after event loss, complete no-JavaScript operation, bounded paged detail, hostile-text defenses, and presentation-specific accessibility.

The final reasoning adopts all three conclusions. It resolves packet preparation as a synchronous durable `repair-prepare` operation because AC-4 requires writer-fenced state publication and the operation performs a private state mutation. It does not dispatch a worker or runtime and needs no user confirmation because it cannot mutate production.

### External Sources Consulted

No sources beyond the sprint-index selection and the concrete repository references distilled in `technical-handbook.md` and the area reasoning were consulted. The selected evidence resolves the design pressures without widening governed context.

## 4. Decision Register

| ID | Decision | Rationale | Trade-offs | Contracts / Requirements | Evidence Expected |
|---|---|---|---|---|---|
| D-01 | Keep repair as one `internal/sprint` verification protocol; keep run control, platform mechanics, app mapping, and adapters in their current authority domains | Preserves product ownership and one-way dependencies | The sprint package grows and needs focused phase functions | AC-1, AC-3, AC-5, AC-8; ARCH-CORE-001/002, ARCH-ENTRY-001, ARCH-SHARED-001, WF-BOUNDARY-001 | Import review, service tests, adapter tests, architecture review |
| D-02 | Freeze one immutable current issue packet, then bind a single-use confirmation after durable acceptance and before dispatch | Makes every mutation traceable to viewed, current authority | Adds a prepared packet and accepted-but-unconfirmed boundary | AC-1, AC-2, AC-3; SEC-AUTHZ-001, SEC-DEFAULT-001, WF-IDEMPOTENCY-001, PERSIST-INTEGRITY-001 | Packet determinism, stale joins, changed-field confirmation, replay, restart, no-dispatch tests |
| D-03 | Generate in a private isolated copy; product code derives, validates, publishes, and directly applies a strict bounded patch without Git or shell | Removes model write authority and proves actual scope | Rejects binary, implicit rename, broad formatter, and other legitimate but unsupported repairs | AC-3, AC-4; SEC-FILES-001, SEC-INJECT-001, LLM-TOOL-001, EVAL-SAFETY-001 | Path/link attacks, malformed patches, target races, apply journal, rollback and actual-diff tests |
| D-04 | Extend strict verification persistence with immutable detail, a mutable digest-bound current pointer, flow-state schema migration, mutation lease, run-control fence, and `terminalizing` barrier | Keeps evidence auditable and stale writers powerless | More records and one explicit recovery state | AC-4, AC-8, AC-9; PERSIST-SCHEMA-001, PERSIST-ATOMIC-001, PERSIST-MIG-001, PERSIST-RECOVERY-001, WF-COMP-001 | Permission, digest, partial write, stale fence, crash, rollback, lock-release, migration and recovery tests |
| D-05 | Run the frozen reverification ladder sequentially and stop wider work after the first required non-pass | Preserves widening semantics and independent review authority | Slower than parallel verification | AC-5; WF-STATE-001, TEST-INT-001, LLM-LIFECYCLE-001 | Gate-order, skipped-gate, frozen-command, target drift, new-finding, and delta-review tests |
| D-06 | Implement one cycle used once by manual mode and looped by automatic mode only behind current manual proof and lower-only persisted budgets | Prevents a second mutation engine and makes boundedness demonstrable | The cycle result and policy state are richer | AC-3, AC-6, AC-7; CFG-SOURCE-001, CFG-TYPE-001, PERF-BOUND-001, EVAL-COST-001, WF-VERSION-001 | Manual one-cycle proof, proof invalidation, all budget/stop conditions, restart and resume counter tests |
| D-07 | Keep durable operation lifecycle and semantic repair outcome separate; commit exactly one immutable result from product facts | Prevents runtime completion, cancellation, or event order from manufacturing a verdict | Clients display two correlated states | AC-8, AC-10; ERR-CODE-001, ERR-DATA-001, OBS-TASK-001, WF-STATE-001 | Six outcomes, lifecycle combinations, late completion, duplicate cancellation, result CAS, exit/HTTP mappings |
| D-08 | Treat cancellation, cleanup, resume, startup recovery, and server shutdown as governed protocol stages | Cleanup and ownership uncertainty remain truthful | Recovery is conservative and may require human action | AC-8, AC-9; LLM-LIFECYCLE-001, WF-RETRY-001, WF-COMP-001, OBS-DIAG-001 | Process-tree cancellation, cleanup timeout, committed-apply resume, abandoned owner, shutdown and restart tests |
| D-09 | Expose one bounded versioned app projection through CLI, JSON, TUI, and server-rendered browser resources; events remain non-authoritative | Delivers parity without exposing private records | DTO, paging, and renderer work increase | AC-2, AC-3, AC-10; CLI-JSON-001, CLI-EXIT-001, CLI-NONINT-001, OBS-CORR-001, PERSIST-READ-001 | Canonical fixture, no-JS flow, CSRF/auth, reconnect/gap, hostile text, bounded paging and parity tests |
| D-10 | Gate automatic exposure and sprint acceptance on current real manual proof, full offline/race/build gates, three reviews, and real-runtime interface dogfood | Fixtures cannot prove runtime isolation, cleanup, or end-to-end authority | Release takes longer and may spend provider budget | AC-6, AC-11; TEST-SMOKE-001/002, TEST-E2E-001, DOC-OPS-001, EVAL-HUMAN-001 | `go test`, race, build, manual four-interface dogfood, bounded automatic stop/resume dogfood, review artifacts |

## 5. Detailed Architecture / Design Reasoning

### D-01: Ownership and the protocol boundary

`internal/sprint` will contain one visible top-level repair protocol with named phases for admission, packet publication, confirmation validation, proposal, apply, reverification, cleanup, proof qualification, and result derivation. Pure functions will normalize packets, classify paths, compare scope, calculate progress, choose stop reasons, and derive outcomes. The implementation must not turn these phases into a generic workflow framework or scatter them across adapters.

The sprint service will depend on narrow volatile seams for runtime proposal generation, process/isolation mechanics, clock/deadlines, persistence failure injection, target identity, and writer fencing. Pure policy remains concrete. `internal/platform/process` will never receive issue severity, acceptance rules, or outcome authority. `internal/runcontrol` will never decide semantic repair success. This follows the handbook's thin-entrypoint and explicit-composition findings and ARCH-CORE-001/002, ARCH-LAYER-002, and ARCH-SHARED-001.

The accepted debt is growth in the existing sprint package. Focused files and phase functions are preferable to a new package tree because state and policy already live there. A plugin engine, generic repair framework, and adapter-specific service are rejected because each would split or widen the most sensitive authority chain.

### D-02: Packet, preparation, acceptance, confirmation, and dispatch

`PrepareRepair` will run as a synchronous durable `repair-prepare` operation. It will validate current QA, review, smoke, evidence, map, theory, check-catalog, policy, governed-input, implementation, and target digests as one snapshot. It will normalize and sort all finite fields, compute allowed and forbidden paths from product policy, and immutably publish or reuse one packet. It will then terminate its durable operation. It will not launch a goroutine, initialize a runtime, acquire production mutation authority, or accept `--yes`.

Starting repair uses a separate durable operation. Existing operation preparation first presents the canonical request and confirmation facts. `RunOperation` revalidates them, durably accepts and claims a run without dispatch, then asks the repair service to publish `confirmation.json`. The confirmation digest covers packet digest, full target identity, canonical request, mode, explicit automatic opt-in, effective limits and safe sources, governed-input and policy fingerprints, durable operation run ID, operational attempt ID, and fencing generation. Only after immutable confirmation publication may the shared runner dispatch work.

An accepted operation whose confirmation cannot be committed terminates as an operational persistence failure with no semantic success and no worker. Confirmation is consumed once. Duplicate submission may return the canonical correlated run but cannot mint another owner. Resume requires fresh confirmation over frozen authority and latest proven boundary; callers cannot resend editable mode, scope, deadline, or maxima.

The extra durable prepare operation is an accepted cost. It satisfies writer-fenced private state mutation and gives crash/recovery visibility while remaining runtime-free. A single prepare-and-start command, confirmation before durable acceptance, or dispatch before confirmation is rejected because each leaves authority implicit or incomplete. This implements AC-1 and AC-2 and follows the handbook's admission-before-work and typed-state findings.

### D-03: Proposal, scope enforcement, and production apply

The runtime receives a private `0700` isolated workspace and only packet-approved bounded context. It can edit only the copy. Product code compares that workspace against its immutable baseline and derives a canonical bounded unified text patch. It rejects malformed hunks, NUL or binary data, special files, symlinks, hard links, unsupported mode changes, ambiguous paths, and implicit rename semantics. An authorized rename must be representable as an independently allowed delete and add; otherwise the run escalates.

Before apply, product code persists `proposal.patch`, parses it again, verifies every old and new path, checks file and byte totals, validates expected pre-image digests, rechecks target identity, mutation lease, and run-control fence, and stages replacement bytes privately. It applies through direct same-directory filesystem replacement. It never invokes Git, shell interpolation, runtime tools, command substitution, or formatters.

An apply journal records intended operations, pre-images, and completed replacements. In-process failure attempts compensation. Rollback failure, process death during a multi-file apply, or a state that cannot prove expected pre-images escalates and never schedules another mutation. After apply, product code recomputes target identity and actual changed paths/bytes. Actual changes must match the applied operation set, remain a subset of `allowed_paths`, and avoid every protected class.

This deliberately supports a narrower repair class than a general patch tool. The debt is that valid binary, rename-heavy, mode-only, migration, generated-code, or broad formatting fixes require future adjudication. Direct runtime writes, Git apply, shell patching, whole-root replacement, and pretending multi-file rename is atomic are rejected under AC-3 and AC-4.

### D-04: Persistence, fencing, terminalization, and recovery boundaries

The strict verification store will add the required artifact tree and schemas. Packet, confirmation, proposal, scope, reverification, cleanup, and result files are immutable, private `0600`, bounded, strict-decoded, digest-verified, path-contained, and unknown-version rejecting. `verification/repair-state.json` is the only mutable repair pointer. Publication writes immutable detail first, verifies it, atomically replaces current state second, then updates the bounded flow summary.

`flow-state.json` receives one optional repair summary and an explicit schema-version migration. Existing concern summaries must survive unrelated saves. There is no dual-format write. Crashes may leave unreferenced immutable files, but recovery cannot treat them as current merely because they exist.

The sprint mutation lease serializes target and verification-state mutation. The run-control fence identifies the accepted writer. Both are required where mutable target or current-state authority is involved. Immutable terminal records may be published after verified lease release only by the still-current run-control fence while canonical state is `terminalizing`.

Terminalization will stop work, prove process/isolation cleanup, recheck target identity, and publish writer-fenced `terminalizing` state while holding the mutation lease. It then releases and verifies its own lease. If release is proven, it publishes immutable cleanup/result and terminal current state under the run-control fence. If release is uncertain, the current owner records cleanup uncertainty and `escalated` where safe. Acquisition of any new mutation must acquire the lease and then reject `terminalizing`, closing the release/publication race. Recovery of an abandoned terminalizing run may record interruption or escalation but may not infer verification or publish manual proof.

The extra current state and flow migration are accepted complexity. One mutable JSON file, a PID-only lock, checking only a fence or only a lease, or last-writer-wins result replacement are rejected under AC-4, AC-8, and AC-9.

### D-05: Progressive reverification and independent review

The service will run the packet's immutable descriptors in this exact order:

1. Exact reproducer and expected failing condition.
2. Affected primary shard checks.
3. Linked theory confirmation and refutation checks.
4. Boundary, neighboring, and parent-linked follow-up shard checks.
5. Containing approved QA suites, including mapped smoke.
6. Focused Conformance Review delta over packet expectations and actual changed paths.

Each descriptor fixes executable, arguments, working directory, environment allowlist, timeout, output bound, evidence rule, and target checks. A narrower failure, block, stale identity, cancellation, or cleanup uncertainty records every wider gate as skipped with reason and next action. No gates run concurrently.

The focused review delta is stored in `reverification.json`; it does not overwrite `review.md`, alter the existing verdict, or claim global conformance. Reverification records new failures, reopened issues, issue-set and severity deltas, scope changes, contradictions, and stale target facts even when the original reproducer passes.

Sequential execution costs latency but preserves first-failure meaning and avoids expensive or misleading wider work. Parallel gates and runtime-selected substitute commands are rejected under AC-5. The handbook's structured-concurrency evidence supports explicit cancellation and waiting, not parallelizing a semantically ordered ladder.

### D-06: Manual-first cycle, automatic bounds, progress, and proof

One cycle function owns proposal, validation, evidence publication, apply, actual-scope comparison, reverification, and cleanup. Manual mode invokes it at most once and cannot emit `stalled`. Automatic mode calls the same function again only after validating proof, confirmation, scope, target, policy, ownership, deadline, cleanup, and every budget.

Product defaults and immutable maxima will cover total cycles, mutation cycles, original-issue reopenings, unchanged issue-set cycles, changed files per cycle/run, changed bytes per cycle/run, patch bytes, wall time, runtime attempts, model turns, command count/time, output bytes, retained cycles, and cleanup time. Workspace and environment values can lower only. Request selection is initially limited to automatic `max-cycles`, which can only lower the merged value. Effective values and safe source labels are frozen before confirmation.

Cycle boundaries persist consumed counters and the original absolute deadline before another cycle can be scheduled. Checks run before every proposal, apply, approved command, publication, retry, and resume. Progress is one of: exact failure removed, admitted issue set reduced, highest admitted severity reduced, or a new bounded fact changes the next permitted action. Rewording, reordered evidence, repeated explanations, and runtime confidence do not count.

Automatic work stops before another mutation. Unchanged proof state, repeated patch digest, repeated target identity, stagnation, cycle exhaustion, or reopening exhaustion yields `stalled` when no stronger outcome applies. Scope or issue growth, severity growth, design/product decisions, contradictions, target or governed-input drift, unsupported test changes, unknown schemas, uncertain evidence, and cleanup uncertainty yield `escalated`.

`verification/manual-repair-proof.json` is a bounded digest-bound pointer created only after a real manual production apply completes the entire ladder, proves cleanup, remains current, and ends `verified` or `verified_with_findings`. Its protocol fingerprint names repair schema, executable code identity, product maxima, effective policy, isolation capabilities, approved checks, governed inputs, runtime policy, final target identity, and qualifying result digest. Automatic admission reports component mismatches. Fixture, dry-run, isolated-only, hand-authored, and administrative proof is inadmissible.

The richer persisted policy is accepted cost. A second automatic engine, config-enabled automatic mode, reset-on-resume limits, or model-defined progress is rejected under AC-6 and AC-7.

### D-07: Closed semantic outcomes and operational lifecycle

The semantic result is a closed enum derived only from persisted product facts:

- `verified`: the issue is gone, every required gate passes, cleanup is proven, target is current, no admitted in-scope issue remains, and the delta review passes.
- `verified_with_findings`: the issue is gone, every gate completes, cleanup is proven, and only current non-blocking findings remain without severity or scope growth.
- `failed`: deterministic admissible evidence proves the issue remains or a required check fails after permitted work, with no missing prerequisite preventing that conclusion.
- `blocked`: an unavailable prerequisite, runtime, executable, permission, lock, persistence operation, or evidence source prevents a semantic conclusion.
- `escalated`: scope, severity, target, governed input, design, expectation, mutation safety, schema, evidence, or cleanup requires human adjudication.
- `stalled`: automatic mode reaches a repetition, stagnation, cycle, or reopening bound without a stronger terminal result.

Run control retains its own lifecycle vocabulary for accepted, running, cancellation requested, completed/failed, cancelled, interrupted, and cleanup uncertain facts. A completed operation may contain semantic `failed`, `blocked`, `escalated`, or `stalled` because the protocol completed truthfully. Infrastructure loss before product result leaves the semantic result absent. Cleanup uncertainty maps to semantic `escalated` only when the current fence or conservative recovery can publish safely. Cancellation never maps to verified or failed.

Result publication is immutable create-or-compare. Existing result wins over late runtime completion, duplicate confirmation, repeated cancellation, stale owners, recovery, and SSE order. CLI start/resume exits zero only for the two verified outcomes; successful read queries return zero and expose the stored outcome. HTTP resource reads return `200` for semantic terminal results and reserve non-2xx for transport/request failures.

Two fields create some client complexity, but a combined status is rejected because it would violate AC-8 and hide which subsystem owns each fact.

### D-08: Cancellation, cleanup, resume, shutdown, and recovery

Canonical cancellation reaches runtime calls, approved commands, waits, process groups, and the cycle coordinator. Cleanup uses a separate context that is not already cancelled but is bounded by the frozen cleanup limit. It proves process-tree termination, isolated workspace removal, apply compensation state, and lock disposition. Missing proof expires as cleanup uncertainty.

Resume requires the same packet, mode, confirmation lineage, target and governed-input identities, policy, consumed budgets, absolute deadline, and latest immutable cycle boundary. A committed production apply is never repeated. Resume may continue uncompleted reverification or cleanup when evidence and target facts permit; uncertainty escalates instead of replaying mutation.

Graceful server shutdown enters draining, rejects new mutations, requests canonical cancellation for every server-owned repair, waits for bounded cleanup and terminal persistence, closes streams, and exits. Browser disconnect affects only the observer. Startup recovery validates owner liveness, packet and confirmation digests, target/governed identities, pointers, apply journal, process facts, isolated workspace, counters, deadline, and locks. It may preserve valid evidence, finish cleanup, record interruption, or require human action. It never adopts a dead worker or infers success from process absence, a patch file, or changed production bytes.

Conservative escalation is accepted over aggressive automatic recovery. Reapplying a patch, resetting a deadline, force-releasing an unowned lock, or publishing success for a dead owner is rejected under AC-9.

### D-09: Shared app, CLI, TUI, browser, and public facts

Add a focused `RepairUseCases` app contract for prepare, packet, start integration, status, paged cycles, cycle detail, result, resume, cancel, and recover. Generic guarded operations own durable acceptance, confirmation, dispatch, cancellation, and observation. Adapters cannot read `verification/`, calculate digests, derive budgets, inspect production files, or call `internal/sprint` directly.

Public DTOs expose bounded identity, authority, freshness, lifecycle, outcome, scope totals/previews, limits, gate summaries, cleanup, evidence references, retention boundaries, blockers, and next action. They exclude patch bodies, production contents, full prompts, provider payloads, unsafe environment values, unrestricted command output, and private evidence. Lists include total, returned, and opaque digest-bound cursors. Progress events carry correlation and transient phase/cycle/gate facts only.

CLI uses the `sprint <project> <sprint> repair` family with `prepare`, `start`, `status`, `packet`, `cycles`, `result`, `resume`, `cancel`, and `recover`. `prepare` never accepts `--yes`; `start` and mutable resume require explicit `--yes`; automatic start also requires `--automatic`. JSON stdout remains one versioned document while progress stays on stderr.

HTTP reads use versioned repair resources. Mutations use the existing guarded operation endpoints. Browser packet review and confirmation are separate CSRF-protected forms and work without JavaScript. TUI uses separate packet, cycle, and result routes with a second guarded mutation action. Automatic controls always render with current proof availability and mismatch reasons. Both interfaces reload canonical status after reconnect, replay gap, terminal event, cancellation, restart, or route revisit.

App-level bounds/redaction combine with HTML escaping and TUI control-character handling. Accessibility requires visible labels, descriptive confirmation actions, keyboard operation, bounded live regions, non-color status meaning, reduced-motion behavior, and narrow-screen/terminal resilience.

The accepted cost is additive DTO, paging, renderer, and compatibility work. A client-side state machine, adapter-owned outcome mapping, private-file query, SSE authority, or separate mutation endpoint family is rejected under AC-10.

### D-10: Acceptance, documentation, and rollout

Implementation will update architecture, CLI reference, user guide, schema, recovery, local-web, and release-checklist documents. Documentation must state packet and confirmation authority, exact commands, lower-only limits, all outcomes, exit behavior, cancellation/recovery, no-JavaScript operation, manual-proof invalidation, and escalation rules.

Acceptance runs focused tests throughout implementation, then `go test ./...`, `go test -race ./...`, and `go build ./cmd/ultraplan`. Architecture Review must trace every mutation from packet digest through result and proof. Sprint Review must map AC-1 through AC-11 and selected contract clauses to code and evidence. Deep Smoke Sprint must first prove one real manual production repair through all four interfaces and the full cleanup/result chain.

Only after that proof is current may automatic mode be exposed or dogfooded. The automatic run must exercise at least one bounded stop path and prove consumed counters and deadline survive restart or resume. Failure of proof-pointer publication leaves the manual result valid but automatic unavailable. Sprint 39 remains responsible for broad convergence and production hardening.

The delayed automatic rollout is accepted. Simultaneous manual/automatic delivery, fixture-created proof, or release based only on unit tests is rejected under AC-6 and AC-11.

## 6. Contracts and Standards Traceability

| Contract / Standard | Applicable Decisions | How Compliance Will Be Proven |
|---|---|---|
| ARCH-CORE-001/002, ARCH-LAYER-002, ARCH-ENTRY-001, ARCH-SHARED-001 | D-01, D-03, D-09 | Import review, focused sprint service, generic platform facts, thin app/adapters, no reverse product imports |
| ERR-CORE-001, ERR-CODE-001, ERR-TRANS-001, ERR-STARTUP-001, ERR-DATA-001, ERR-REDACT-001 | D-02, D-04, D-07, D-08, D-09 | Stable blocker matrix, preserved causes/correlation, preflight rejection, distinct persistence/cleanup errors, safe public messages |
| CFG-SOURCE-001, CFG-TYPE-001, CFG-START-001, CFG-ENV-001, CFG-OBS-001 | D-06, D-09 | Precedence, defaults, maxima, lower-only workspace/environment/request tests and safe source reporting |
| OBS-CORE-001, OBS-CORR-001, OBS-DIAG-001, OBS-TASK-001, OBS-PII-001 | D-04, D-07, D-08, D-09 | Packet/run/attempt/fence/cycle correlation, canonical status, truthful terminals, bounded redacted events and diagnostics |
| SEC-AUTHZ-001, SEC-INPUT-001, SEC-INJECT-001, SEC-FILES-001, SEC-DESER-001, SEC-DEFAULT-001 | D-02, D-03, D-04, D-09 | Confirmation binding, strict schemas, explicit argv, path/link/file-class attacks, target races, CSRF/object scope, fail-closed unknowns |
| TEST-SEAM-001, TEST-UNIT-001, TEST-INT-001, TEST-SMOKE-001/002, TEST-FAIL-001, TEST-CONTRACT-001, TEST-E2E-001, TEST-MIGRATION-001 | D-01 through D-10 | Replaceable volatile seams, tables, real filesystem/process tests, race/failure injection, fixtures, migrations, four-interface dogfood |
| DOC-ARCH-001, DOC-PUBLIC-001, DOC-OPS-001, DOC-EXAMPLE-001, DOC-GEN-001 | D-09, D-10 | Updated authoritative docs, checked commands/JSON examples, recovery guidance, release checklist |
| CLI-SHAPE-001, CLI-IO-001, CLI-EXIT-001, CLI-JSON-001, CLI-LIFE-001, CLI-SAFE-001, CLI-NONINT-001 | D-02, D-07, D-09 | Help and JSON fixtures, stdout/stderr tests, explicit `--yes`, automatic opt-in, exits, cancellation, no hidden prompts |
| LLM-BOUNDARY-001, LLM-TOOL-001, LLM-IO-001, LLM-LIFECYCLE-001, LLM-RETRY-001, LLM-RUN-001, LLM-EXPOSE-001, LLM-SAFETY-001 | D-01, D-03, D-06, D-08, D-09 | Isolated fixed runtime, structured bounded result, no production root, cancellation/cleanup, bounded retries, app-only exposure |
| EVAL-SCOPE-001, EVAL-REG-001, EVAL-DATA-001, EVAL-SAFETY-001, EVAL-COST-001, EVAL-MODEL-001, EVAL-HUMAN-001 | D-03, D-05, D-06, D-10 | Adversarial scope regressions, safe fixtures, budget facts, runtime/policy fingerprinting, independent review and manual confirmation |
| WF-SCOPE-001, WF-BOUNDARY-001, WF-STATE-001, WF-RETRY-001, WF-IDEMPOTENCY-001, WF-COMP-001, WF-VERSION-001 | D-01, D-02, D-04, D-06, D-07, D-08 | One app-owned protocol, explicit phases/outcomes, replay dedupe, apply journal, resume boundaries, versioned automatic behavior |
| PERF-BUDGET-001, PERF-BOUND-001, PERF-CONC-001, PERF-COST-001 | D-05, D-06, D-09 | Sequential gates, finite cycle/runtime/command/output/retention/query limits, persisted consumption, duration/cost evidence |
| PERSIST-SCHEMA-001, PERSIST-MIG-001, PERSIST-ATOMIC-001, PERSIST-INTEGRITY-001, PERSIST-READ-001, PERSIST-DERIVED-001, PERSIST-RECOVERY-001 | D-02, D-04, D-06, D-08, D-09 | Strict private versions, flow migration, immutable-detail/current-pointer order, digests, bounded queries, apply recovery, proof staleness |
| Project architecture, PRD, TRD, and roadmap | D-01 through D-10 | Repair remains verification-scoped, manual-first, shared across local interfaces, before Sprint 39 dogfood and later content gates |
| Architecture Review Protocol | D-01 through D-05, D-07, D-08 | Trace ownership and every mutation boundary, including fencing, apply, terminalization, and independent review |
| Sprint Review Protocol | D-02 through D-10 | Trace AC-1 through AC-11, protected classes, all outcomes/stops, parity, docs, tests, race and build evidence |
| Deep Smoke Sprint Protocol | D-03, D-05, D-06, D-08, D-09, D-10 | Real manual four-interface run first; current proof then bounded automatic stop with restart/resume counters |

## 7. Trade-Offs, Debt, and Future Considerations

### Accepted Trade-Offs

- A two-operation prepare/start interaction is more deliberate than one command. It makes the frozen packet reviewable and keeps production confirmation exact.
- Durable acceptance before confirmation creates an accepted-but-unconfirmed recovery boundary. It is required because acceptance identity and fence generation are confirmation inputs.
- Strict text patch support rejects some valid repairs. A narrow auditable mutation set is preferable to Git, shell, binary, formatter, or implicit rename authority.
- Compensated per-file apply cannot provide crash atomicity across several files. The design records an apply journal and escalates uncertainty rather than claiming a transaction the filesystem does not provide.
- Sequential reverification is slower than parallel checks. It preserves the required widening order, first authoritative non-pass, and target identity clarity.
- Two ownership controls and a `terminalizing` state add lifecycle complexity. They prevent stale publication and close the lock-release/result-publication race.
- Separate operational lifecycle and semantic outcome require more client rendering. They prevent cancellation or runtime completion from becoming a product verdict.
- Automatic mode carries many persistent counters and fingerprints. The cost is necessary to prove lower-only bounds and resume behavior.
- Server-rendered, paged interfaces are less fluid than a client application. They retain a complete no-JavaScript path and avoid browser-owned authority.
- Real manual proof delays automatic delivery. It is the evidence gate required for a production mutation protocol.

### Known Technical Debt

- Multi-file apply remains compensation-based. A future transactional filesystem or differently governed mutation boundary could reduce uncertainty, but Sprint 38 will not invent one.
- Binary patches, mode-only changes, implicit renames, generated-code workflows, migrations, dependency upgrades, and broad formatting are unsupported repair classes.
- `internal/sprint` gains a substantial protocol. Focused files and pure decision helpers mitigate this, but future extraction should occur only after repeated concrete reuse.
- Flow-state gains another schema migration and bounded projection. It remains a summary rather than a verification store.
- Manual-proof code identity and isolation capability fingerprints require careful version ownership. Sprint 39 dogfood may reveal irrelevant invalidation inputs that should be narrowed through a governed change.
- Public repair history is intentionally bounded. Pruned evidence remains visible as a retention gap rather than a complete historical explorer.
- Automatic progress rules are conservative and may stall repairs that a human could continue. That is preferable to silent authority expansion.

### Rejected Alternatives

- **Direct runtime writes to production:** rejected because model permission would replace product scope and apply authority.
- **Git, shell, or formatter-based apply:** rejected because side effects, control files, hooks, ignored files, and actual changed scope become harder to bound.
- **One-click prepare and start:** rejected because packet review, durable acceptance, and exact confirmation cannot be proven separately.
- **Confirmation before durable acceptance:** rejected because the digest would omit the accepted run, attempt, and fence that will execute.
- **A second automatic engine:** rejected because scope, cleanup, persistence, and outcome policy would drift from manual proof.
- **Automatic enablement from configuration or proof-file presence:** rejected because current proof and explicit per-run opt-in are required.
- **Parallel reverification gates:** rejected because wider work could run after a narrow non-pass and obscure first-failure meaning.
- **One combined operation/result status:** rejected because operational completion and semantic verification have different authorities.
- **One mutable repair JSON file:** rejected because immutable evidence, resume boundaries, digest integrity, and exactly-one result would be weaker.
- **PID lock as sole ownership:** rejected because it does not prove durable accepted generation or prevent stale publication.
- **Aggressive crash recovery that reapplies or completes work:** rejected because process absence and file presence cannot prove intent, apply, cleanup, or verification.
- **Adapter-specific repair endpoints or state machines:** rejected because confirmation, cancellation, and outcomes would diverge across interfaces.
- **Client-side browser application:** rejected because Sprint 38 needs guarded server-owned operations, not a second authoritative state model.
- **General plugin framework:** rejected because version, collision, timeout, and trust costs have no scoped Sprint 38 use case.

### Future Considerations

- Sprint 39 should measure false starts, scope rejections, repair convergence, proof invalidation, cleanup certainty, operator comprehension, and automatic stop quality across representative repositories.
- Broader patch classes require separate adjudication and a new security/recovery decision, not an extension hidden inside the parser.
- Better transactional mutation may be considered only if real crash evidence shows compensation is operationally inadequate.
- Manual-proof fingerprint components may be narrowed only with evidence that a component causes irrelevant invalidation without affecting protocol safety.
- Bounded concurrency may be considered for independent mechanics, but the reverification ladder remains ordered unless requirements change.
- Richer history or evidence browsing should use bounded app queries and retained-state policy, never direct private-file access.
- Content identity, retrieval, product SQLite authority, knowledge graphs, hosted operation, and remote repair remain behind the post-Sprint-39 gates.

## 8. Implementation Constraints

- Keep repair policy, records, scope, progress, reverification, proof, and outcomes in `internal/sprint` focused files.
- Keep generic copy, process, identity, comparison, cancellation, and cleanup mechanics in platform packages with no repair semantics.
- Use a synchronous durable `repair-prepare` operation for packet publication. It must not start a goroutine, runtime, or production mutation.
- Split mutable start acceptance from dispatch so confirmation publication completes before every worker and runtime launch.
- Build packets only from current digest-validated product records. Caller input may select an issue ID but cannot supply scope, commands, expectations, or evidence.
- Normalize finite collections and repository-relative paths deterministically. Reject empty, absolute, escaping, linked, protected, ambiguous, or unsupported paths.
- Persist packet, confirmation, proposal, scope, reverification, cleanup, and result records immutably with strict schemas, bounds, private permissions, no symlink following, and content digests.
- Publish immutable detail before current repair state and flow summary. Never make unreferenced crash leftovers current by file presence.
- Increment and migrate the strict flow-state schema once; preserve existing review, smoke, and QA summaries during unrelated saves.
- Require the active mutation lease and run-control fence at every applicable mutation and current-state boundary. A stale or cancelled writer cannot apply, publish, release another owner's lock, or refresh proof.
- Use `terminalizing` as a barrier around cleanup, lease release, and terminal publication. Recovery cannot infer verification from it.
- Give the runtime only an isolated writable copy and bounded packet-approved context. Never pass a writable production root.
- Derive the proposal from isolated filesystem differences. Parse and validate it independently before apply.
- Support only bounded text changes through direct contained filesystem operations. Do not use Git, shell commands, formatters, hooks, or runtime tools for production apply.
- Record pre-images and an apply journal. Do not claim multi-file atomicity. Escalate uncertain compensation or crash state.
- Compare full target identity before proposal, before apply, after apply, and after reverification. Compare actual changed paths and bytes after apply.
- Run frozen reverification commands sequentially in the required order. Do not let the runtime substitute commands or outcomes.
- Keep Conformance Review independent and leave `review.md` unchanged. Store only the focused bounded delta in repair evidence.
- Use one cycle implementation for manual and automatic modes. Manual permits one mutation cycle and never emits `stalled`.
- Require a current qualifying real manual proof and explicit automatic opt-in before automatic work.
- Define finite defaults and immutable maxima for every AC-7 budget. Workspace/environment/request values may only lower applicable values.
- Persist consumed counters and the absolute deadline before scheduling another cycle. Resume cannot reset either.
- Derive progress and stop conditions from normalized product facts, never model prose or events.
- Keep semantic outcomes closed and separate from existing durable lifecycle states. Commit one immutable result with create-or-compare semantics.
- Propagate cancellation to every owned worker and child. Use a distinct bounded cleanup context and require affirmative cleanup facts.
- Resume only from immutable proven boundaries and never repeat a committed apply.
- Use one `RepairUseCases` contract and shared guarded operations across CLI, JSON, TUI, and browser. Adapters cannot read private files or call the sprint service directly.
- Bound and redact every public field and collection. Escape HTML and strip terminal controls. Do not expose patches, production content, prompts, provider payloads, unsafe environment, or unrestricted output.
- Keep browser operation complete without JavaScript. Treat SSE and TUI progress as observational; canonical status refresh owns truth.
- Do not add Git mutation, test/requirement/evidence changes, broad refactors, migrations, plugins, alternate persistence, issue tracking, retrieval, hosted operation, or remote workers.

## 9. Expected Evidence and Validation Strategy

### Tests

- Packet table tests for one current issue, exact reproducer, finite normalized scope, every required fingerprint, deterministic bytes/digest, idempotent reuse, and every admission rejection in AC-1.
- Protected-class tests for requirements, acceptance, planning artifacts, review/smoke/QA/repair evidence, flow state, config, Git metadata/hooks/ignore policy, tests, snapshots, baselines, generated evidence, links, renames, deletes, descendants, and formatter side effects.
- Confirmation tests that change each digest input independently, duplicate and replay submissions, cross-run/sprint identities, manual/automatic mode, automatic opt-in, acceptance failure, confirmation write failure, and restart before dispatch.
- Isolation and patch tests for allowed copy writes, production-root denial, malformed hunks, traversal, absolute paths, NUL/binary data, symlinks, hard links where supported, special files, implicit renames, mode changes, file/byte/patch limits, and unsafe retained proposal evidence.
- Apply tests with target drift before proposal/apply/after apply, expected pre-image mismatch, stale fence, lost lease, failure before every file replacement, compensation failure, crash journal recovery, actual-scope mismatch, and pre-existing unrelated changes.
- Persistence tests for strict schemas, unknown fields/versions, `0700`/`0600`, containment, digest mismatch, immutable create-or-compare, partial writes, current-pointer rollback, flow-state migration/preservation, retention, cursor gaps, and stale writers.
- Reverification tests for exact fixed order, skipped wider gates, immutable descriptors, missing executable, timeout, truncation, cancellation, target drift, cleanup uncertainty, new/reopened issues, issue/severity delta, contradiction, containing smoke, and independent delta review.
- Outcome tables for all six semantic outcomes and all operational lifecycle combinations, including late completion, duplicate terminal proposals, cancellation, interruption, cleanup uncertainty, persistence loss, and recovery.
- Automatic tests for every default and maximum, lower-only precedence, every check point, total/mutation cycles, reopenings, unchanged issue set, repeated patch/target, stagnation, scope/severity growth, design decisions, unknown schema, uncertain evidence, deadlines, runtime/model/command/output limits, retention, and restart/resume preservation.
- Manual-proof tests for qualifying real facts, every fingerprint component, stale/mismatched/weaker proof, non-manual result, missing cleanup, dry-run/fixture/isolated-only/hand-authored attempts, proof publication failure, and admission recheck immediately before mutation.
- Durable-operation tests for prepare, accept-before-confirm, confirm-before-dispatch, one owner, writer generation, cancellation routing, graceful shutdown, restart reconciliation, exactly one result, and no fire-and-forget goroutines.
- App and adapter contract tests for canonical DTO parity, bounded cursor binding, stale and retention errors, text/JSON separation, CLI exits, HTTP methods/status, object authorization, CSRF/body limits, no-JavaScript forms, TUI routes, reconnect/replay gaps, and browser disconnect independence.
- Hostile-output tests for HTML, terminal control sequences, long paths, secret-bearing diagnostics, prompts, provider payloads, environment values, raw output, and pagination truncation facts.
- Race tests for confirmation replay, cancellation versus completion, lease/fence changes, apply versus target change, pointer publication, terminalization, SSE/event delivery, startup recovery, and server shutdown.

### Logs / Metrics

- Structured records correlate project, sprint, QA attempt, issue, packet digest, repair run, durable operation run, operational attempt, fence generation, cycle, gate, runtime, and process identity.
- Authority events record safe packet/target/policy fingerprints, mode, confirmation state, effective limit sources, ownership checks, and replay decisions without secrets or private content.
- Mutation events record isolated workspace identity, proposal digest and size, allowed/actual path counts, changed files/bytes, target identity transitions, apply-journal phase, scope result, and compensation state.
- Reverification events record gate order, descriptor identity, status, skip reason, duration, bounded output facts, issue/severity delta, and review-delta identity.
- Automatic events record consumed/remaining counters, absolute deadline, progress fact, repeated patch/target, stagnation, reopening, stop reason, and proof mismatch components.
- Cleanup events record cancellation source, process-tree result, workspace removal, compensation, lease release, uncertainty, and bounded cleanup duration.
- Public diagnostics expose stable category, safe message, correlation, blocker, and next action. They exclude prompts, production content, provider payloads, secrets, unsafe environment, and unrestricted output.

### Manual / Review Checks

- Run `go test ./...`, `go test -race ./...`, and `go build ./cmd/ultraplan` in the implementation repository.
- Run Architecture Review and trace one mutation from QA issue through packet, prepare operation, start acceptance, confirmation, dispatch, isolation, proposal, apply, actual diff, each gate, cleanup, terminalization, result, and proof.
- Run Sprint Review and map AC-1 through AC-11 plus every selected contract family to code, tests, docs, and retained evidence.
- Run Deep Smoke Sprint first with one real manual issue and production mutation. Exercise CLI text and JSON, TUI, browser no-JavaScript and enhanced behavior, canonical reconnect, cancellation/status where applicable, the full ladder, cleanup, result, and proof pointer.
- Confirm that the same packet, service, operation runner, scope guard, apply boundary, reverifier, cleanup path, and outcome derivation were used by all interfaces.
- After current manual proof exists, run one explicit automatic repair that reaches a bounded stop path and demonstrate that counters and deadline survive restart or resume.
- Inspect manual and automatic records for private permissions, bounded size, digest correctness, target freshness, actual paths, cleanup facts, terminal uniqueness, public redaction, and agreement with flow/app/interface projections.
- Verify documentation commands and JSON examples against executable fixtures and record exact real-runtime run IDs, protocol fingerprints, provider/model/runtime identity, durations/costs, review result, and evidence paths.

## 10. Handoff to Planning

`plan.md` should implement these decisions in dependency order: repair types and pure policy; strict persistence and flow migration; durable prepare plus accept/dispatch refactor; packet and confirmation authority; isolated proposal and strict patch derivation; product-owned apply journal and recovery; sequential reverification; cleanup and terminal outcomes; shared manual cycle; manual proof; lower-only automatic loop; app DTOs and guarded operations; CLI; TUI; browser; documentation; focused tests; full race/build gates; manual dogfood; then automatic admission and bounded-stop dogfood.

Every task must cite decision IDs, ACs, contract clauses, target files, and expected evidence. The plan must preserve manual-first sequencing: automatic implementation and exposure cannot be accepted until the shared manual protocol, all four interfaces, and real manual proof pass.

The plan must not reopen these decisions:

- Repair is one `internal/sprint` verification protocol; platform, run control, app, and adapters retain their distinct authorities.
- Packet preparation is a synchronous durable writer-fenced operation with no worker or runtime.
- Mutable start is accepted before confirmation and dispatched only after confirmation is immutable.
- Runtime writes only an isolated copy; product code derives and directly applies a strict bounded patch without Git or shell.
- Strict immutable detail, current pointer, flow migration, two ownership controls, and `terminalizing` are required.
- Reverification is sequential and Conformance Review remains independent.
- Manual and automatic modes share one cycle; automatic requires current real manual proof, explicit opt-in, persisted lower-only bounds, and product-derived progress.
- Operational lifecycle and semantic outcome remain separate, and exactly one immutable semantic result wins.
- Cancellation, cleanup, resume, recovery, and shutdown remain governed stages and never infer success.
- CLI, JSON, TUI, and browser share bounded app DTOs and guarded operations; progress events are not authoritative.
- Acceptance requires the complete offline/race/build gate, three reviews, real manual four-interface dogfood, then a proof-gated automatic bounded-stop dogfood.

## Decisions

The final decisions are D-01 through D-10 in Sections 4 and 5. They fix repair ownership, packet and confirmation authority, isolated proposal and product-owned apply, strict persistence and fencing, sequential reverification, manual-proof-gated automatic reuse, closed semantic outcomes, conservative recovery, shared bounded interfaces, and release evidence.

## Assumptions And Risks

The governing assumptions and risk register are in Section 2. Planning must preserve the listed mitigations, especially whole-snapshot packet validation, accept/confirm/dispatch ordering, apply journaling, link and protected-class defense, dual ownership checks, immutable command coverage, affirmative cleanup proof, deterministic progress, manual-proof fingerprinting, and canonical cross-interface refresh.

## Implementation Constraints

Section 8 is normative. No implementation path may give a model or adapter production-write authority, bypass durable confirmation or ownership, mutate governed/test/evidence/Git assets, parallelize the semantic gate order, reset automatic budgets, infer success during recovery, or expose private evidence through public interfaces.

## Expected Evidence

Section 9 defines the required unit, integration, failure-injection, race, migration, contract, interface, documentation, review, and real-runtime evidence. Sprint acceptance requires all offline tests plus one qualifying manual production repair through all four interfaces; automatic work remains unavailable until that proof is current and then must demonstrate a bounded stop whose consumed authority survives restart or resume.
