# Sprint requirements: Read-only QA decomposition and synthesis

> Project: `ultraplan-go`
> Sprint: `36-read-only-qa`
> Purpose: the authoritative, human-readable sprint contract. All other sprint artifacts must satisfy these requirements.

## Sprint goal

Add a deterministic, durable, and cross-surface read-only QA phase that maps changed behavior into bounded verification surfaces, records resumable theory outcomes, and synthesizes cross-shard findings without generating checks, promoting issues, repairing code, or mutating the target repository.

## Required outputs

| Output | Path | Description |
| --- | --- | --- |
| Verification phase model | `../ultraplan-go/internal/sprint/verification_phase.go` | Defines `VerificationPhase` independently from `PlanningStage`, including compatibility mapping for the existing review capability. |
| QA domain model | `../ultraplan-go/internal/sprint/qa_types.go` | Defines schema-versioned maps, shards, theories, attempts, budgets, outcomes, synthesis, fingerprints, and stable verification-scoped identifiers. |
| QA state store | `../ultraplan-go/internal/sprint/qa_state.go` | Loads, validates, migrates or rejects, atomically writes, resumes, cancels, and reconciles detailed QA state outside `flow-state.json`. |
| Deterministic QA mapper | `../ultraplan-go/internal/sprint/qa_map.go` | Builds bounded behavioral verification surfaces from governed sprint context, execute evidence, changed paths, review findings, tests, boundaries, and risk tags. |
| Read-only investigation prompts and policy | `../ultraplan-go/internal/sprint/qa_prompt.go` | Prepares bounded investigator and synthesizer requests with enforceable read-only permissions, explicit outcome conditions, budgets, and target fingerprints. |
| QA orchestration service | `../ultraplan-go/internal/sprint/qa.go` | Runs dry-run mapping, bounded shard investigations, focused reruns, cancellation, resume, invalidation, recovery, progress, and next-action derivation. |
| Global synthesis | `../ultraplan-go/internal/sprint/qa_synthesis.go` | Deduplicates and challenges theories, records interactions and negative outcomes, and creates only bounded follow-up shards without issue promotion. |
| QA domain and state tests | `../ultraplan-go/internal/sprint/qa_state_test.go` | Covers schemas, stable IDs, unknown versions, atomic writes, cancellation, reconciliation, resume, and fingerprint invalidation. |
| QA mapping tests | `../ultraplan-go/internal/sprint/qa_map_test.go` | Covers reproducibility, changed-path coverage, primary ownership, boundary overlap, budgets, risk inputs, and map fingerprints. |
| QA orchestration and synthesis tests | `../ultraplan-go/internal/sprint/qa_test.go` | Covers read-only permissions, theory outcomes, bounded follow-up, progress, cancellation, resume, recovery, and non-promotion. |
| CLI and JSON commands | `../ultraplan-go/internal/app/sprint_commands.go` | Exposes `qa --dry-run`, `qa`, focused shard execution, status, cancellation/recovery controls, and machine-readable output through shared sprint use cases. |
| CLI and JSON tests | `../ultraplan-go/internal/app/sprint_commands_test.go` | Freezes command help, validation, text/JSON schemas, exit behavior, read-only enforcement, and app/CLI agreement. |
| Shared QA use cases | `../ultraplan-go/internal/app/sprint_usecases.go` | Provides typed map, status, theory, synthesis, run, cancel, resume, and recovery results for CLI, TUI, and web adapters. |
| Shared use-case tests | `../ultraplan-go/internal/app/sprint_usecases_test.go` | Proves adapter-independent projections, durable run correlation, bounded progress, and consistent next actions. |
| TUI QA view | `../ultraplan-go/internal/tui/qa_view.go` | Shows maps, shards, progress, theory outcomes, synthesis, blockers, cancellation, and recovery through app use cases. |
| TUI QA tests | `../ultraplan-go/internal/tui/qa_view_test.go` | Covers keyboard operation, state vocabulary, cancellation, dropped-delivery recovery, bounded rendering, and parity fixtures. |
| Browser QA handlers | `../ultraplan-go/internal/web/qa_handlers.go` | Maps versioned HTTP requests to shared QA use cases and durable operation progress without importing sprint internals. |
| Browser QA tests | `../ultraplan-go/internal/web/qa_handlers_test.go` | Covers routes, JSON compatibility, authorization, hostile content, progress, cancellation, reconnect, restart recovery, and parity. |
| Sprint QA presentation | `../ultraplan-go/internal/web/templates/sprint.html` | Presents map coverage, shard progress, theories, synthesis, blockers, and recovery with a useful no-JavaScript snapshot. |
| Browser QA enhancement | `../ultraplan-go/internal/web/static/js/operations.js` | Adds bounded progressive enhancement for QA progress and cancellation over existing operation/run APIs without client-side authority. |
| QA command documentation | `../ultraplan-go/docs/cli-reference.md` | Documents Conformance Review terminology, compatibility aliases, QA commands, JSON behavior, outcomes, and recovery controls. |
| QA architecture documentation | `../ultraplan-go/docs/architecture.md` | Documents verification-phase ownership, state authority, deterministic mapping, read-only policy, synthesis, and adapter boundaries. |
| QA user workflow | `../ultraplan-go/docs/user-guide.md` | Explains mapping, investigation, theory outcomes, focused reruns, cancellation, resume, and interpretation of blocked or inconclusive work. |
| QA browser documentation | `../ultraplan-go/docs/local-web.md` | Documents browser visibility, guarded starts, durable progress, cancellation, refresh, reconnect, and recovery. |
| QA recovery documentation | `../ultraplan-go/docs/recovery.md` | Documents stale attempts, changed fingerprints, interrupted runs, invalid state, cancellation, and safe restart behavior. |
| QA JSON and state schemas | `../ultraplan-go/docs/phase3-json-schemas.md` | Documents additive Conformance Review metadata plus detailed QA state and public JSON compatibility rules. |
| QA release checks | `../ultraplan-go/docs/release-checklist.md` | Adds deterministic mapping, mutation isolation, cross-surface parity, race, build, documentation, and gated dogfood checks. |
| QA state root | `projects/<project>/sprints/<sprint>/verification/state.json` | Stores the schema version, current attempt pointer, canonical QA summary, freshness, verdict-neutral status, and artifact digests. |
| Attempt map | `projects/<project>/sprints/<sprint>/verification/attempts/<attempt-id>/map.json` | Stores the immutable input fingerprint, changed-path coverage, shard definitions, budgets, and deterministic identifiers for one attempt. |
| Shard outcomes | `projects/<project>/sprints/<sprint>/verification/attempts/<attempt-id>/shards/<shard-id>.json` | Stores resumable investigator attempts, theories, static evidence, outcomes, context expansion, commands, and stop reasons. |
| Synthesis outcome | `projects/<project>/sprints/<sprint>/verification/attempts/<attempt-id>/synthesis.json` | Stores deduplication, contradictions, cross-shard interactions, bounded follow-up requests, blockers, and next actions without issues or repairs. |

## Acceptance criteria

- [ ] Human-facing text calls the existing analytical `review` capability "Conformance Review"; `review`, `review.md`, current verdict rules, and existing JSON clients remain compatible, and any `conformance-review` alias invokes the same capability rather than a second review implementation.
- [ ] `VerificationPhase` represents `conformance-review`, `qa`, and future `repair` independently from `PlanningStage`; QA is not inserted into the planning artifact sequence or modeled as another planning stage.
- [ ] `qa --dry-run` builds a map without invoking investigators. For byte-identical governed inputs and the same implementation fingerprint, repeated runs produce byte-identical normalized maps, IDs, ordering, budgets, and fingerprints apart from explicitly excluded observation timestamps.
- [ ] Mapping consumes current requirements, code-context, sprint reasoning and plan, execute evidence and changed paths, selected contracts and protocols, Conformance Review findings, adjacent tests and known checks, package/interface/state boundaries, and applicable risk tags.
- [ ] Every changed path appears in exactly one primary shard. Cross-package, interface, producer/consumer, state-transition, public-API, and cross-cutting behavior may use explicit boundary shards and bounded overlap, but no changed path is orphaned or assigned to multiple primary shards.
- [ ] Each map and shard records positive numeric hard limits for changed and contextual paths, context expansion, behavioral concerns, theories, iterations, commands, wall-clock duration, concurrent investigators, and follow-up shards. Invalid limits fail validation, and hitting a limit stops or blocks work without silently broadening scope.
- [ ] Map, shard, theory, and attempt IDs are deterministic and scoped to a documented QA schema version. Unknown major versions fail closed with recovery guidance; supported migration or compatibility behavior is fixture-tested and does not claim global content identity.
- [ ] Each theory records a falsifiable claim, basis, affected verification surface, expectation references, severity if confirmed, confirmation condition, refutation condition, inconclusive condition, proposed safe evidence strategy, current implementation fingerprint, and final outcome.
- [ ] Durable theory outcomes distinguish at least `confirmed`, `refuted`, `invalid`, `inconclusive`, `blocked`, `cross_shard`, and `not_applicable`; negative and blocked outcomes remain inspectable and are not discarded merely because they do not support an issue.
- [ ] Investigators can read only assigned and approved context, request bounded context expansion, and run only approved non-mutating existing checks. Permission tests deny product or verification-code writes, generated tests, fixtures, probes, shell indirection, Git mutation, issue promotion, and repair.
- [ ] Before and after each investigator attempt, identity checks prove that production source, production tests, governed planning inputs, verification code, smoke-harness content, and Git state were not mutated. QA may write only its declared detailed state and normal durable operational records.
- [ ] Cancellation stops new scheduling, cancels active runtime work through the canonical path, records active shard attempts truthfully, and preserves completed outcomes. Browser or SSE disconnection stops observation only and never implies QA cancellation or success.
- [ ] A restart with unchanged fingerprints resumes incomplete work without rerunning completed valid shards. A governed-input, implementation, review, map, or shard fingerprint change invalidates affected work explicitly and never reuses stale outcomes as current evidence.
- [ ] Detailed QA state writes are atomic and recoverable. `flow-state.json` retains only bounded summary, freshness, current status, next action, and pointers or digests; it does not become the shard/theory attempt database.
- [ ] Global synthesis deduplicates equivalent theories, retains contradictions and negative outcomes, identifies cross-shard interactions, and may create only explicitly budgeted follow-up shards with recorded parent evidence. It does not create issue records, repair eligibility, production patches, or a passing Conformance Review verdict.
- [ ] A failed or blocked current Conformance Review remains failed or blocked after diagnostic QA. QA status and theory outcomes cannot overwrite, upgrade, or substitute for the independent Conformance Review verdict.
- [ ] CLI text, CLI JSON, TUI, browser HTML/JSON, and durable run inspection agree on map fingerprint, changed-path coverage, shard status, progress, theory outcome, synthesis status, blocker, cancellation state, terminal result, and next action for shared fixtures.
- [ ] Every runtime-backed QA execution is durably accepted before child work starts and is discoverable through Sprint 35 run-control surfaces. Progress and diagnostics remain bounded and redacted, and QA does not add a separate browser-only operation registry, event authority, or unbounded subscriber path.
- [ ] Browser refresh, session rotation, reconnect, observer restart, and dropped live delivery recover QA state from authoritative app and workspace state. Server-rendered pages remain useful without JavaScript, and hostile theory/evidence text is escaped.
- [ ] Deterministic tests cover stable mapping, changed-path coverage, boundary overlap, all theory outcomes, budget exhaustion, permission denial, cancellation, resume, fingerprint invalidation, atomic-write failure, unknown schema, synthesis deduplication, bounded cross-shard follow-up, and cross-surface agreement.
- [ ] A gated real-repository dogfood run produces a deterministic map, read-only shard outcomes, synthesis, durable cross-surface progress, cancellation/recovery evidence, and a clean target identity check. Missing runtime or environment prerequisites produce a truthful blocked result and do not satisfy this exit criterion.
- [ ] `go test ./...`, `go test -race ./...`, `go vet ./...`, `go build ./cmd/ultraplan`, and `git diff --check` pass in `../ultraplan-go`.

## Non-goals

- Generated tests, fixtures, probes, smoke scenarios, verification patches, or any other evidence-producing mutation; these belong to Sprint 37 and require proven isolation.
- Canonical `qa.md`, evidence adjudication, issue promotion, repair eligibility, regression-candidate promotion, or smoke-as-QA integration; these belong to Sprint 37 or later.
- Production repair, manual repair, automatic repair, repair loops, issue packets, or any mutation authorized by a theory outcome.
- Replacing, renaming, or removing `review`, `review.md`, `smoke`, `smoke.md`, current verdict rules, or existing Phase 3 compatibility contracts.
- Turning `flow-state.json` into detailed QA storage or making operational run records authoritative for governed artifacts, source, tests, or verification outcomes.
- A general-purpose QA package, workflow engine, scheduler, issue tracker, sandbox framework, remote worker, daemon, broker, or plugin system.
- Global content identity, provenance, retrieval, RAG, embeddings, knowledge graphs, alternate authored-artifact persistence, cloud authority, hosted service, multi-user collaboration, or remote exposure.
- Fixing Sprint 35 review findings as part of QA execution. Applicable findings are mapping and investigation inputs; any product fix requires a later governed repair scope.

## Constraints

- `internal/sprint` owns verification phases, QA semantics, mapping, detailed state, investigation policy, synthesis, and compatibility. Generic runtime and process packages must not import sprint QA types or decide QA outcomes.
- CLI, TUI, and HTTP remain adapters over typed `internal/app` use cases. `internal/web` must not import `internal/sprint`, parse CLI output, run investigators directly, or persist alternate QA truth.
- Conformance Review remains analytically read-only and keeps `review` and `review.md` compatibility. QA remains a separate empirical phase and cannot alter the current review verdict.
- Investigator target access is read-only. The allowed command policy must use explicit argv and bounded environment forwarding, reject shell wrappers and path escapes, and prohibit writes to production, tests, verification code, governed inputs, harness files, and Git state.
- Detailed QA state must be schema-versioned, atomic, path-contained, bounded, and fail closed on unknown major versions or invalid fingerprints. Verification-scoped identifiers must have documented migration behavior.
- Mapping and synthesis must be deterministic for unchanged normalized inputs. Nondeterministic timestamps, run IDs, sessions, and provider metadata must not affect stable map or shard identity.
- Concurrency, context, command count, output size, wall-clock time, retries, follow-up shards, progress buffers, and rendered history must have validated finite limits selected during reasoning and frozen in tests and documentation before implementation planning completes.
- Existing Sprint 35 durable run identity, acceptance, cancellation, event ordering, redaction, retention, replay, reconciliation, and product-authority separation remain binding. Sprint 36 must not weaken them or create a second operational state model.
- Normal tests must be deterministic and offline with fake runtimes and temporary workspaces. Real-runtime dogfood must be gated, bounded, read-only against the target, and reported truthfully.
- Automatic Git mutation remains prohibited.

## Dependencies

| Prior sprint / output | Required for | Notes |
| --- | --- | --- |
| Sprint 35 review, `projects/ultraplan-go/sprints/35-durable-run-observability/review.md` | QA map inputs and operational risks | Verdict is `pass_with_findings`. Current findings, including persistence-failure visibility, stale-writer race proof, cross-surface parity, and SSE bounds, remain theories or risk inputs rather than implicit repair scope. |
| Sprint 35 smoke, `projects/ultraplan-go/sprints/35-durable-run-observability/smoke.md` | Promotion into read-only QA | Verdict is `pass` with six real-boundary tests, complete mapped coverage, no smoke findings, and no open harness issues. Preserve the selected same-host durable run topology and recovery contract. |
| Sprint 35 run-control implementation in `../ultraplan-go/internal/runcontrol/` | Durable QA identity, progress, cancellation, and recovery | QA executions must use the existing workspace-wide acceptance and observation path rather than a new process-local operation store. |
| Sprint 34 shared context and Sprint 33 code-context | Grounded mapper and investigator context | Reuse validated requirements/code-context and current source references without introducing retrieval, caching, or a parallel context manifest. |
| Current execute and Conformance Review outputs | Changed paths, implementation evidence, findings, and expectations | Mapping requires current non-stale execute evidence and the current independent review result; stale or missing inputs block current QA rather than being guessed. |
| `../ultraplan-go/docs/plans/integrated-roadmap.md` | Phase 5 sequencing and Sprint 36 exit gate | Current implementation-repository plan was read for this requirements stage as required by the workspace roadmap. |
| `../ultraplan-go/docs/plans/post-execution-qa-and-repair-loop.md` | QA terminology, decomposition, state, command, and safety design | Current implementation-repository plan was read directly. Sprint 36 adopts only its terminology/compatibility and read-only Phase 1 scope. |
| Project index, roadmap, PRD, TRD, and Architecture | Scope, ownership, compatibility, and release gates | These are the authoritative project sources for deterministic mapping, read-only enforcement, state separation, and cross-surface behavior. |
| Agentwrap/OpenCode runtime and existing app operation seams | Bounded investigator execution | Agentwrap remains runtime supervision authority; UltraPlan supplies typed read-only policy, bounded prompts, validation, and product-owned QA state. |

## Review expectations

| What | How verified |
| --- | --- |
| Phase and compatibility model | Architecture and API review prove `VerificationPhase` is separate from `PlanningStage`, `review`/`review.md` remain compatible, and no duplicate Conformance Review implementation exists. |
| Deterministic map | Golden and table-driven fixtures rerun unchanged inputs, compare normalized bytes and fingerprints, and assert exactly one primary shard for every changed path plus explicit bounded overlap. |
| Grounded coverage | Fixture inspection traces each shard to requirements, code-context, execute paths, review findings, selected contracts, adjacent tests, boundaries, and risk tags. |
| Read-only enforcement | Adversarial permission tests attempt file writes, generated tests, shell indirection, path/symlink escape, Git commands, harness mutation, and governed-input mutation; every attempt is denied and recorded without target drift. |
| Theory quality and outcomes | Fixtures require falsifiable conditions and cover confirmed, refuted, invalid, inconclusive, blocked, cross-shard, and not-applicable results while retaining negative evidence. |
| State durability and freshness | Fault injection covers interrupted atomic writes, unknown versions, cancellation, owner loss, restart, unchanged resume, changed-input invalidation, stale attempts, and pointer/digest reconciliation. |
| Bounded synthesis | Synthesis fixtures prove deduplication, contradiction retention, parent-linked follow-up, hard follow-up limits, and absence of issues, repair packets, or verdict promotion. |
| Cross-surface agreement | One fixture is inspected through app results, CLI text/JSON, TUI, browser HTML/JSON, and durable run detail; all canonical fields, blockers, and next actions match. |
| Sprint 35 contract preservation | Regression tests cover durable acceptance before execution, explicit cancellation, reconnect/replay, redaction, bounded progress, restart recovery, and separation from product artifact authority. |
| Browser safety and accessibility | Handler/template tests cover same-origin authorization, escaped hostile content, no-JavaScript snapshots, keyboard/focus behavior, reduced motion, reconnect, session rotation, and bounded rendering. |
| Documentation | Review CLI help, JSON schemas, architecture, user workflow, browser behavior, recovery, and release checklist against executable fixtures and exact artifact/state paths. |
| Dogfood and release gates | Inspect the real-repository map, shard and synthesis records, run-control evidence, before/after target identity, cancellation/recovery trace, and logs for all required test, race, vet, build, and diff checks. |
| Scope exclusions | Diff and dependency review confirm no generated checks, writable isolation, issue promotion, `qa.md`, smoke migration, repair, content identity, retrieval, alternate persistence, cloud, or Git mutation. |
