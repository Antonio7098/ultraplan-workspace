# Sprint Plan: Read-only QA decomposition and synthesis

> Project: `ultraplan-go`
> Sprint: `36-read-only-qa`
> Source: `reasoning.md`
> **Inputs Used:** `projects/ultraplan-go/project-index.md`, `projects/ultraplan-go/roadmap.md`, `projects/ultraplan-go/docs/ARCHITECTURE.md`, `projects/ultraplan-go/docs/PRD.md`, `projects/ultraplan-go/docs/TRD.md`, `projects/ultraplan-go/sprints/36-read-only-qa/requirements.md`, `projects/ultraplan-go/sprints/36-read-only-qa/code-context.md`, `projects/ultraplan-go/sprints/36-read-only-qa/sprint-index.md`, `projects/ultraplan-go/sprints/36-read-only-qa/technical-handbook.md`, `projects/ultraplan-go/sprints/36-read-only-qa/reasoning/api-design.md`, `projects/ultraplan-go/sprints/36-read-only-qa/reasoning/architecture.md`, `projects/ultraplan-go/sprints/36-read-only-qa/reasoning/frontend.md`, `projects/ultraplan-go/sprints/36-read-only-qa/reasoning.md`

This plan executes `reasoning.md`. It does not reopen the selected phase model, state authorities, mapping rules, read-only policy, limits, public contracts, presentation model, or release gates.

The selected evidence reports were not reread individually. `technical-handbook.md`, the three area-reasoning documents, and final `reasoning.md` contain their decision-relevant findings. The exact `code-context.md` and prepared source references were sufficient for implementation seam planning, so no additional source file was treated as a new design authority.

## Reasoning Source

- **Sprint Reasoning:** `reasoning.md`
- **Sprint Index:** `sprint-index.md`
- **Technical Handbook:** `technical-handbook.md`
- **Area Reasoning:** `reasoning/api-design.md`, `reasoning/architecture.md`, `reasoning/frontend.md`

## Sprint Status

- **Status:** `implementation complete; real-runtime release gate blocked`
- **Owner:** `implementation agent`
- **Start Date:** `2026-08-24`
- **Completion Date:** `2026-08-24` (implementation; promotion remains blocked)

## Decisions To Execute

| Decision | Source Section | Execution Implication |
| --- | --- | --- |
| Separate verification phases and preserve review compatibility | `reasoning.md#decision-1-separate-verification-phases-and-preserve-review-compatibility` | Add `VerificationPhase` with `conformance-review`, `qa`, and reserved `repair`; keep planning order, `review`, `review.md`, `sprint.review`, smoke, verdicts, and existing JSON compatible; route `conformance-review` to the existing review implementation. |
| Keep four authorities separate and publish detailed state atomically | `reasoning.md#decision-2-keep-four-authorities-separate-and-publish-detailed-state-atomically` | Keep QA detail under `verification/`, flow state as a bounded pointer summary, run control as operational authority, and the checkout as source authority; publish private files in dependency order with writer fencing. |
| Build one deterministic map with exact primary ownership | `reasoning.md#decision-3-build-one-deterministic-map-with-exact-primary-ownership` | Normalize and sort every governed input, assign every changed path to exactly one primary shard, represent allowed overlap explicitly, and block unknown paths instead of omitting them. |
| Freeze configuration, budgets, scheduling, and retention | `reasoning.md#decision-4-freeze-configuration-budgets-scheduling-and-retention` | Validate immutable effective settings once, permit configuration only to lower product defaults, fingerprint values and sources, use bounded workers and queues, and block before attempt or storage limits are exceeded. |
| Enforce read-only investigation in code and verify target identity | `reasoning.md#decision-5-enforce-read-only-investigation-in-code-and-verify-target-identity` | Use Agentwrap read-only default-deny permissions, a product-owned explicit-argv check catalog, contained paths, bounded environment forwarding, and before/after identity manifests. |
| Orchestrate durable starts, cancellation, resume, and explicit recovery | `reasoning.md#decision-6-orchestrate-durable-starts-cancellation-resume-and-explicit-recovery` | Accept and claim runs before child work, distinguish attempt-wide and shard-local stops, cancel through run control, reuse only current semantic evidence, and keep recovery runtime-free. |
| Validate complete theories and synthesize without adjudication | `reasoning.md#decision-7-validate-complete-theories-and-synthesize-without-adjudication` | Validate all theory fields and outcomes, retain negative records, make final synthesis pure and deterministic, and cap parent-linked follow-up without producing issues, repairs, `qa.md`, or review verdicts. |
| Define one app, CLI, JSON, and HTTP contract | `reasoning.md#decision-8-define-one-app-cli-json-and-http-contract` | Expose typed app DTOs and closed operation kinds, keep query paths runtime-free, implement the fixed command and route families, and use stable bounded error mappings. |
| Keep TUI and browser presentation bounded and authority-free | `reasoning.md#decision-9-keep-tui-and-browser-presentation-bounded-and-authority-free` | Add focused sprint QA routes and views, render complete no-JavaScript snapshots, separate product and observer states, escape hostile content, and apply bounds before rendering. |
| Prove semantics offline and reserve dogfood for real boundaries | `reasoning.md#decision-10-prove-semantics-offline-and-reserve-dogfood-for-real-boundaries` | Use pure, temporary-workspace, fake-runtime, fake-process, adapter-parity, and race tests for normal evidence; use gated dogfood only for Agentwrap, process, filesystem, browser, Git, cancellation, and recovery boundaries. |

## Requirements / Contracts To Satisfy

| Contract / Requirement ID | Required Behavior | Evidence Planned |
| --- | --- | --- |
| Architecture; `REQ-PHASE`; `REQ-COMPAT` | QA is a verification phase, not a planning stage. Human labels say Conformance Review while review and smoke compatibility remains exact. | Phase enumeration, planning-order, alias-equivalence, unchanged review JSON/artifact, import-boundary, and architecture-review evidence. |
| Persistence And Migrations; Errors; Security; `REQ-STATE` | Detailed QA state is v1, `qa-v1` scoped, contained, private, normalized, atomic, digest-checked, fenced, recoverable, and separate from flow-state detail. | Strict decoding, unknown version, trailing JSON, mode, containment, symlink, rename, fsync, pointer-order, digest, stale-writer, retention, quota, and reconciliation tests. |
| Testing; Performance; `REQ-MAP`; `REQ-BOUNDS` | Mapping is byte-stable, grounded in every required input, gives each changed path one primary owner, records explicit bounded overlap, and exposes all effective limits. | Golden maps, repeated normalized-byte comparisons, ownership and overlap tables, unknown-path fixtures, risk-input traces, invalid-limit tables, and fingerprint invalidation. |
| Security; LLM Runtime; LLM Evaluation / Cost / Safety; `REQ-READONLY` | Investigators can inspect only approved context and request only product-owned existing checks. Unsupported enforcement or identity drift blocks promotion. | Permission request inspection, adversarial write/shell/Git/path/environment cases, descriptor tests, output limits, before/after identity comparisons, and clean target evidence. |
| Workflows; Observability; `REQ-RUN` | Runtime QA is durably accepted and claimed before child work, uses the canonical cancellation path, resumes current evidence without worker adoption, and reports cleanup or persistence uncertainty. | Run-control inventory, no-child-on-acceptance-failure, fencing races, queued/running cancellation, timeout, replay, restart, dropped observer, and runtime-free recovery fixtures. |
| Testing; Workflows; `REQ-THEORY`; `REQ-SYNTHESIS` | Every theory is falsifiable and retains one closed outcome; synthesis is deterministic, contradiction-preserving, bounded, parent-linked, and verdict-neutral. | Full outcome tables, malformed theory rejection, deterministic synthesis bytes, deduplication, contradictions, interactions, follow-up exhaustion, stale-input rejection, and forbidden-field scans. |
| Architecture; CLI Surface; Errors; `REQ-PARITY` | App DTOs are the only shared product projection. CLI text/JSON, TUI, browser HTML/JSON, and durable run detail agree on canonical QA facts. | One semantic parity fixture, app independence tests, help and exit fixtures, strict transport tests, durable correlation, and stable error-code tests. |
| Security; Documentation; `REQ-WEB` | Browser requests remain loopback-only, allowlisted, strictly decoded, confirmation-bound, same-origin, CSRF-protected, escaped, bounded, and recoverable without JavaScript. | Route, authorization, hostile-content, no-JavaScript, focus, mobile, reduced-motion, reconnect, replay-gap, session-rotation, and server-restart tests. |
| Testing; Documentation; `REQ-RELEASE` | Offline, race, vet, build, diff, protocol, and gated dogfood checks pass; missing real prerequisites report blocked and do not satisfy the gate. | Required commands, Architecture Review, Sprint Review, Deep Smoke Sprint evidence, dogfood records, and clean before/after target identity. |

## Fixed authority and storage contract

| Concern | Authority | Implementation rule |
| --- | --- | --- |
| Map, shard, theory, semantic attempt, synthesis, freshness, blocker, and QA next action | `internal/sprint` files under `verification/` | Load and validate through the sprint service only. Adapters never read these files. |
| Bounded QA summary | `flow-state.json` | Store current status, freshness, counts, next action, cancellation summary, and contained pointer or digest only. Preserve it across planning refreshes. |
| Acceptance, owner lease, fencing, cancellation, event order, liveness, and terminal result | Sprint 35 run control | Reuse existing app and run-control paths. Do not add a QA lifecycle database or adapter registry. |
| Source, tests, governed inputs, verification code, harness content, and Git identity | Contained workspace and target checkout | Hash before and after each investigator attempt. Never reset, clean, or repair drift. |

The detailed state layout is fixed:

```text
verification/state.json
verification/attempts/<attempt-id>/map.json
verification/attempts/<attempt-id>/shards/<shard-id>.json
verification/attempts/<attempt-id>/synthesis.json
```

Publication order is map or changed shard/synthesis record, then `verification/state.json`, then the bounded `flow-state.json` projection. Maps do not change after publication. Directories use `0700`; detailed files use `0600`. Status follows current pointers and digests and never promotes orphan files by scanning.

## Frozen limits

Workspace configuration may lower defaults but cannot exceed maxima. Non-positive or over-maximum values fail validation. Effective values and their sources enter map and policy fingerprints.

| Resource | Default | Maximum |
| --- | ---: | ---: |
| Changed paths in a map | 512 | 512 |
| Primary shards | 32 | 32 |
| Boundary shards | 8 | 8 |
| Follow-up shards per synthesis | 4 | 4 |
| Total planned shards including follow-up | 44 | 44 |
| Pending scheduler entries | 44 | 44 |
| Changed paths per primary shard | 32 | 64 |
| Contextual paths per shard | 64 | 128 |
| Context expansions per shard | 2 | 4 |
| Paths per expansion | 16 | 32 |
| Behavioral concerns per shard | 12 | 24 |
| Theories per shard | 12 | 24 |
| Investigator iterations per shard attempt | 4 | 8 |
| Approved commands per shard attempt | 8 | 16 |
| Runtime retries after the initial call | 1 | 2 |
| Concurrent investigators | 3 | 8 |
| Command wall-clock duration | 5 minutes | 10 minutes |
| Shard wall-clock duration | 20 minutes | 30 minutes |
| Whole QA run wall-clock duration | 60 minutes | 90 minutes |
| Cleanup duration | 30 seconds | 30 seconds |
| Captured output per command | 256 KiB | 512 KiB |
| Stored investigator output per shard attempt | 1 MiB | 2 MiB |
| Investigator or challenger prompt | 512 KiB | 1 MiB |
| Recent product progress records | 100 | 200 |
| Retained QA attempts per sprint | 8 | 8 |
| Total detailed QA state per sprint | 128 MiB | 128 MiB |

Presentation and delivery bounds are 40 summary shards, 24 theories per focused shard page, 200 TUI durable events, 100 browser QA operation rows, Sprint 35's 16 KiB durable event limit and 250 ms progress coalescing, and live-region aggregate announcements no more often than every two seconds. Every partial collection states `Showing X of Y` and uses stable ID order and bounded cursor navigation.

## Stop policy

| Scope | Conditions | Required result |
| --- | --- | --- |
| Whole attempt | Missing, stale, or invalid governed inputs; unsupported or corrupt state; writer-token loss; unavailable permission enforcement; target drift; cancellation; whole-run timeout; authoritative persistence failure | Stop new admission, cancel active work when applicable, preserve prior valid evidence, record the blocker and next action, and never report success. |
| One shard | Malformed investigator output; approved check unavailable or nonzero; denied or exhausted context expansion; iteration, command, output, or wall-clock exhaustion; closed negative theory outcome | Record the local stop truthfully and allow independent shards to continue. Partial synthesis names incomplete or blocked siblings. |
| Cleanup | Process-group, state, event, or resource cleanup exceeds 30 seconds | Record run terminal `cleanup_uncertain`, QA phase `interrupted`, blocker, and recovery action. |
| Retention or storage | A ninth retained attempt or a write over 128 MiB would start | Block before runtime work. `qa recover` may prune only validated non-current attempts while preserving current and newest last-complete attempts. |

## Execution order and release rule

Tasks are ordered so phase identity, domain validation, state authority, and deterministic mapping exist before runtime work or adapters depend on them. Runtime-backed QA must not become releasable until durable acceptance, writer fencing, permission enforcement, target identity, cancellation, and state promotion are proven together. If implementation evidence contradicts a final reasoning decision, stop and revisit the governed reasoning stage before changing this plan or implementation.

## Tasks

- [x] **Task 1: Separate verification identity and freeze effective QA settings**
  > Executes: Decisions 1 and 4; Architecture; Configuration; CLI Surface; `REQ-PHASE`, `REQ-COMPAT`, `REQ-BOUNDS`
  - [x] Add `internal/sprint/verification_phase.go` with `VerificationPhaseConformanceReview`, `VerificationPhaseQA`, and reserved `VerificationPhaseRepair`, plus compatibility mapping to the existing review capability. Keep `PlanningStages()`, planning artifact order, `StageReview`, and `StageSmoke` unchanged.
  - [x] Refactor sprint service composition so verification runtime selection is keyed by `VerificationPhase` and QA receives a validated `QASettings` value without adding QA to `map[PlanningStage]StageRuntime`.
  - [x] Extend the existing typed config and effective-source reporting for the QA model and lower-only limits. Reject zero, negative, and over-maximum values instead of clamping.
  - [x] Keep dry-run and query service construction lazy and runtime-free. Prove they do not initialize Agentwrap, process execution, durable acceptance, or product state writes.
  - [x] Change human review labels to `Conformance Review` only where text is presented. Freeze existing review command, artifact, verdict, operation name, JSON fields, and smoke behavior in regression fixtures.
  - [x] **Stop condition:** phase tests prove QA is absent from planning order, review compatibility bytes remain unchanged, and every frozen limit and source validates before domain implementation proceeds.

- [x] **Task 2: Define the v1 QA domain, identifiers, and validation rules**
  > Executes: Decisions 2-4 and 7; Persistence And Migrations; Errors; Performance; `REQ-STATE`, `REQ-MAP`, `REQ-BOUNDS`, `REQ-THEORY`
  - [x] Add `internal/sprint/qa_types.go` with schema v1 state, maps, primary and boundary shards, coverage, budgets, effective sources, approved-check references, target identity, semantic attempts, investigator attempts, theories, evidence summaries, context requests, command summaries, synthesis, blockers, cancellation, and freshness types.
  - [x] Define closed QA phase values `missing`, `mapped`, `queued`, `running`, `synthesizing`, `completed`, `blocked`, `cancelled`, `interrupted`, `stale`, and `invalid`. Keep freshness, run lifecycle, cancellation, terminal result, and Conformance Review verdict as separate fields.
  - [x] Define exact theory outcomes `confirmed`, `refuted`, `invalid`, `inconclusive`, `blocked`, `cross_shard`, and `not_applicable`. Require claim, basis, verification surface, expectation references, severity if confirmed, confirmation/refutation/inconclusive conditions, safe evidence strategy, implementation fingerprint, attempt history, and final outcome.
  - [x] Implement deterministic `qa-v1` map, shard, theory, and semantic attempt IDs from canonical project, sprint, parent identity, and normalized content. Exclude timestamps, run IDs, sessions, provider metadata, worker order, and current time.
  - [x] Add stable typed QA errors and public categories without flattening causes. Reserve recovery guidance for unknown schema, invalid state, stale input, permission denial, budget exhaustion, conflict, persistence failure, and runtime unavailability.
  - [x] Add table tests in `internal/sprint/qa_state_test.go` for every enum, required field, invalid fingerprint, malformed ID, budget edge, and theory outcome.

- [x] **Task 3: Build contained atomic QA state and bounded flow summary**
  > Executes: Decisions 2, 4, and 6; Architecture; Persistence And Migrations; Security; Workflows; `REQ-STATE`, `REQ-RUN`
  - [x] Add `internal/sprint/qa_state.go` and contained helpers for the fixed verification paths. Reject absolute, lexical escape, symlink escape, unsafe IDs, and references outside the selected sprint root.
  - [x] Implement strict schema v1 decoding with unknown-field and trailing-JSON rejection. Do not invent v0 migration; return actionable unsupported-version recovery and reserve explicit pure migration seams for future shipped versions.
  - [x] Write normalized indented JSON plus newline using same-directory temporary files, flush, close, rename, and directory sync. Apply `0700` directories and `0600` files.
  - [x] Publish immutable map, changed shard or synthesis record, `verification/state.json`, and bounded `flow-state.json` projection in that order. Inject failure hooks at every step and preserve the prior valid file, pointer, digest, and last-complete evidence on failure.
  - [x] Add the bounded QA pointer summary to flow state while preserving strict loading, previous flow-state compatibility, planning refresh behavior, review, and smoke. Detailed maps, theories, commands, and history must never enter flow state.
  - [x] Validate writer tokens carrying run ID, operational attempt ID, and fencing generation before every runtime-backed promotion. Reject stale writers without importing run-control types into `internal/sprint`.
  - [x] Implement explicit cancellation markers, freshness and dependent-work invalidation, pointer/digest reconciliation, interrupted ownership recovery, and bounded retention pruning. Read-only status never repairs or prunes.
  - [x] Test unknown versions, malformed records, permissions, containment, symlinks, atomic failures, pointer order, digest mismatch, stale writers, cancellation, unchanged resume, partial invalidation, quota, retention, and flow-summary mismatch.
  - [x] **Stop condition:** an injected failure at any publication boundary leaves the last valid canonical state readable, and a stale writer cannot promote bytes under race testing.

- [x] **Task 4: Implement deterministic behavior mapping and exact path ownership**
  > Executes: Decisions 3 and 4; Testing; Performance; `REQ-MAP`, `REQ-BOUNDS`
  - [x] Add `internal/sprint/qa_map.go` to load and validate requirements, code-context, sprint and area reasoning, plan, execute state and changed paths, selected contracts and protocols, current Conformance Review findings, adjacent tests and known checks, package/interface/state boundaries, producer-consumer relationships, public APIs, and risk tags.
  - [x] Reuse current governed input and execute authorities rather than reparsing prose as alternate truth. Missing or stale execute or review evidence blocks mapping; a current failed or blocked review remains valid diagnostic input.
  - [x] Normalize paths, line endings, identifiers, references, check-catalog identity, implementation identity, review fingerprint, policy fingerprint, limits, and source labels. Sort every unordered collection before hashing, ID generation, and JSON persistence.
  - [x] Group by coherent behavior with exactly one primary shard for every changed path. Add boundary shards only for named cross-package, interface, producer-consumer, state-transition, public-API, or cross-cutting concerns and record their overlap explicitly.
  - [x] Turn unknown or unclassifiable changed paths into blocked primary shards. Enforce changed-path, shard, path-per-shard, context, concern, queue, and total-planned-shard limits without omission or adaptive widening.
  - [x] Make `qa --dry-run` return normalized map bytes and coverage without runtime construction, durable acceptance, or state writes.
  - [x] Add `internal/sprint/qa_map_test.go` goldens for byte reproducibility, IDs, ordering, exactly-one primary ownership, allowed overlap, unknown paths, risk inputs, required source traces, budgets, and implementation/review/policy fingerprint changes.

- [x] **Task 5: Enforce investigator permissions, approved checks, and target identity**
  > Executes: Decision 5; Security; LLM Runtime; LLM Evaluation / Cost / Safety; `REQ-READONLY`, `REQ-BOUNDS`
  - [x] Add `internal/sprint/qa_prompt.go` to create bounded investigator and challenger packets from the frozen map, assigned common context, approved shard paths, theory contract, implementation fingerprint, context and check IDs, budgets, and explicit stop conditions.
  - [x] Build Agentwrap requests with `read_only`, restricted permissions, default deny, required permission capability, and contained read, list, and search only. Block useful output promotion when enforcement is unsupported.
  - [x] Define a closed in-tree approved-check descriptor catalog with stable ID, executable, argv, cwd, environment-name allowlist, timeout, output limit, and descriptor fingerprint. Investigators may request only a map-owned ID.
  - [x] Reject caller executable, argv, environment, path, prompt, policy, output location, and result content. Reject shell wrappers, command substitution, redirection, interpreter indirection, Git commands, and checks known to write caches, coverage, generated files, fixtures, probes, source, tests, governed inputs, verification code, or harness content.
  - [x] Execute approved descriptors through `internal/platform/process` without adding QA semantics to the platform package. Revalidate implementation and descriptor identity immediately before launch and bound stdout, stderr, timeout, and safe command summaries.
  - [x] Compute sorted before/after identity manifests for production source, production tests, governed sprint and project inputs, verification implementation code, smoke-harness content, and Git HEAD/index/worktree when applicable. Record symlink identity and require contained resolution.
  - [x] Stop new admission and block promotion on identity drift. Record bounded relative categories and paths, preserve prior canonical outcomes, and never clean or reset the target.
  - [x] Add adversarial tests for every denied tool and path class, unsupported permission enforcement, environment redaction, symlink escape, check race, cache-writing commands, output exhaustion, and before/after drift.

- [x] **Task 6: Orchestrate bounded investigations, cancellation, resume, and recovery**
  > Executes: Decisions 4-6; Workflows; Observability; Errors; `REQ-RUN`, `REQ-STATE`, `REQ-READONLY`
  - [x] Add `internal/sprint/qa.go` with runtime-free map/status paths and runtime-backed full or focused shard execution. Focused requests accept only a current map-owned shard ID and do not rerun a completed current shard.
  - [x] Build an instance-scoped scheduler with at most the effective investigator count and 44 pending entries. Do not create one goroutine per shard. Propagate the accepted-run context and whole-run deadline to workers, runtime calls, and processes.
  - [x] Require a current writer token before installing queued/running state or starting child work. Separate attempt-wide safety stops from shard-local malformed output, check failure, expansion denial, budget exhaustion, and closed negative outcomes.
  - [x] Validate investigator responses before promotion. Apply iteration, theory, context expansion, path, command, retry, prompt, output, shard-duration, and whole-run limits and record the exact exhausted resource.
  - [x] Route cancellation through existing run control. Stop admission, cancel active runtime and process contexts, preserve completed promoted shards, and record queued and active shard truthfully. Observer disconnect must not affect work.
  - [x] Use a separate 30-second cleanup context for process groups, final detailed state, bounded summary, safe events, and local resources. Map cleanup timeout to run `cleanup_uncertain` and QA `interrupted`.
  - [x] Resume with a new durable run and writer token while reusing the semantic attempt and completed valid shards only when governed, implementation, review, map, catalog, policy, limit, and selected-shard fingerprints remain current. Never adopt sessions, handles, goroutines, or lost workers.
  - [x] Implement runtime-free `qa recover` under the sprint mutation lock. Reconcile stale ownership, interrupted markers, pointers, digests, flow summary, and retention; direct work needing runtime to `qa resume`.
  - [x] Emit bounded redacted semantic progress for map, queue, shard, synthesis, blocker, cancellation, and recovery through the existing operation event path.
  - [x] Test queue saturation, concurrency, cancellation queued and running, local versus global stops, timeout, malformed output, retry cap, stale token, restart, partial invalidation, observer loss, cleanup uncertainty, persistence degradation, and no-child recovery.

- [x] **Task 7: Synthesize validated theories without issue or verdict promotion**
  > Executes: Decision 7; Testing; Workflows; LLM Evaluation / Cost / Safety; `REQ-THEORY`, `REQ-SYNTHESIS`
  - [x] Add `internal/sprint/qa_synthesis.go` as a pure final synthesis function over validated shard records current for one map fingerprint.
  - [x] Normalize equivalence keys, group duplicates, order theories and evidence deterministically, retain refuted, invalid, blocked, inconclusive, cross-shard, and not-applicable records, and preserve contradictions and declared interactions.
  - [x] Permit a bounded read-only challenger to propose schema-validated challenge records. Fingerprint and persist those records as explicit inputs before the pure final synthesis step.
  - [x] Select at most four follow-up shards from remaining map budget. Require stable IDs, explicit parent theory and evidence references, bounded paths and concerns, and the same permission, identity, and state rules as initial shards.
  - [x] Make partial synthesis state incomplete siblings and blockers. `completed` means bounded work ended, never pass.
  - [x] Prohibit issue records, repair eligibility, patches, generated checks, `qa.md`, and Conformance Review verdict fields in types, output, and persistence. Preserve an existing failed or blocked review unchanged.
  - [x] Add deterministic-byte, deduplication, contradiction, negative-retention, interaction, challenger, parent-link, follow-up exhaustion, stale-input, and forbidden-field tests in `internal/sprint/qa_test.go`.

- [x] **Task 8: Add typed QA use cases and durable operation integration**
  > Executes: Decisions 6 and 8; Architecture; Workflows; Observability; `REQ-RUN`, `REQ-PARITY`
  - [x] Extend `internal/app/sprint_usecases.go` with adapter-independent `QAMap`, `QAStatus`, `QAShard`, `QATheory`, `QASynthesis`, `RunQA`, `ResumeQA`, `CancelQA`, and `RecoverQA` requests and results. Do not expose sprint persistence structs.
  - [x] Carry schema version, project, sprint, phase, freshness, attempt and run correlation, input/implementation/map/policy fingerprints, coverage, limits, bounded shards, outcome totals, synthesis, progress, blocker, cancellation, terminal result, and next action where applicable.
  - [x] Add closed operation kinds `qa-status`, `qa-dry-run`, `qa-start`, `qa-resume`, and `qa-recover`; add only the map-owned `Shard` request option. Reject caller models, budgets, commands, prompts, permissions, environment, paths, fingerprints, ownership, attempt IDs, and theory content.
  - [x] Prepare canonical governed inputs, mutation class, target-read-only warning, effective limits, model source, map fingerprint, selected shard, prerequisites, and durable refresh path. Revalidate preparation fingerprints immediately before execution.
  - [x] Put `qa-start` and `qa-resume` through durable acceptance and owner claim before service construction or child work. Convert the accepted fence into the opaque sprint writer token. Keep status and dry-run read-only and recovery runtime-free but explicitly mutating.
  - [x] Make `CancelQA` validate sprint/run correlation, delegate to `RunUseCases.CancelRun`, and return current cancellation plus QA snapshot without claiming terminal cancellation early.
  - [x] Map bounded semantic progress to existing run-control events and retain Sprint 35 event, replay, redaction, cancellation, terminal, and persistence-degradation behavior.
  - [x] Add app use-case, preparation, stale-confirmation, no-child acceptance failure, duplicate start, resume correlation, cancellation race, bounded projection, next-action, and runtime-free query tests in `internal/app/sprint_usecases_test.go` and operation tests.

- [x] **Task 9: Expose the CLI, JSON, and review alias contracts**
  > Executes: Decisions 1 and 8; CLI Surface; Errors; Documentation; `REQ-COMPAT`, `REQ-PARITY`
  - [x] Extend `internal/app/sprint_commands.go` with `qa --dry-run`, `qa [--shard]`, `qa resume [--shard]`, `qa status`, `qa cancel --run`, and `qa recover`, all with optional `--json` and no public restart, suite, model, budget, command, or path controls.
  - [x] Add `conformance-review` as argument normalization to the existing `review` handler. Preserve `review`, `review.md`, `sprint.review`, JSON fields, verdict behavior, and exit behavior exactly.
  - [x] Use JSON schema version 1 and operation `sprint.qa`. Return bounded result data even when a blocked, cancelled, interrupted, runtime, or persistence condition produces a nonzero exit.
  - [x] Map valid dry-run, current status, and completed diagnostic QA to success regardless of theory outcomes; usage faults to `ExitUsage`; stale/invalid/denied/budget blockers to `ExitValidation`; runtime and persistence infrastructure to `ExitRuntime`; cancellation and interrupted partial work to `ExitPartial`.
  - [x] Render Conformance Review verdict and freshness separately from verdict-neutral QA phase. Use `Read-only QA completed`, never `QA passed`, and do not label confirmed theories as issues.
  - [x] Extend the runtime-backed CLI inventory so start and resume cannot bypass durable acceptance. Prove dry-run, status, cancel, and recovery use only their intended authorities.
  - [x] Add command help, flag conflict, ID, text, JSON, stable error, exit, runtime-free, durable acceptance, alias-equivalence, and app/CLI agreement fixtures in `internal/app/sprint_commands_test.go`.

- [x] **Task 10: Add bounded TUI QA navigation and recovery**
  > Executes: Decision 9; CLI Surface; Security; Performance; `REQ-PARITY`, `REQ-WEB`
  - [x] Add `RouteSprintQA`, `RouteSprintQAShard`, and `RouteSprintQATheory` to the existing route model and expose `QA` under sprint navigation. Relabel the current review action as `Conformance Review` without removing review, smoke, verify, or diagnostic actions.
  - [x] Add `internal/tui/qa_view.go` over app use cases only. Render phase, freshness, review verdict, fingerprints, coverage, blocker, next action, bounded shards, outcome totals, synthesis, cancellation, terminal result, and durable run links without reading verification files.
  - [x] Keep the one-column viewport. Use Enter for detail, Escape for back, arrows or `j`/`k` for selection, `r` for authoritative refresh, `c` for explicit active-run cancellation, and `q` for observation-only exit.
  - [x] Refresh QA and durable state after dropped local delivery, route entry, explicit refresh, and terminal messages. Never reconstruct product outcomes from the in-memory event slice.
  - [x] Apply the fixed 40-shard, 24-theory, 200-event, viewport, and stable paging bounds before rendering. Sanitize ANSI and non-printing controls before width calculations.
  - [x] Render all phases, outcomes, cancellation states, freshness, observer states, unknown counts, blockers, and errors as text without color-only meaning. Keep `80x24` and `40x12` usable and preserve non-TTY fallback to CLI text or JSON.
  - [x] Add `internal/tui/qa_view_test.go` coverage for routes, keys, help, selection, wide/narrow output, all vocabulary, bounds, hostile content, requested versus terminal cancellation, dropped delivery, recovery, and shared parity fixtures.

- [x] **Task 11: Add versioned QA HTTP resources and no-JavaScript browser views**
  > Executes: Decisions 8 and 9; Architecture; Security; Performance; Documentation; `REQ-PARITY`, `REQ-WEB`
  - [x] Add `internal/web/qa_handlers.go` with bounded app-only mappings for `GET /api/v1/projects/{project}/sprints/{sprint}/qa`, `/qa/map`, `/qa/shards/{shard-id}`, `/qa/theories/{theory-id}`, and `/qa/synthesis`.
  - [x] Extend strict operation decoding and compatibility fixtures for `qa-status`, `qa-dry-run`, `qa-start`, `qa-resume`, and `qa-recover`, including only the `shard` option. Keep starts and cancellation on existing operation and run endpoints.
  - [x] Preserve the web import boundary: production web code imports only `internal/app` and the standard library and does not inspect files, run investigators, or decide outcomes.
  - [x] Add QA overview, focused shard, and focused theory HTML routes under one sprint. Update `internal/web/templates/sprint.html` with complete server-rendered snapshots, Conformance Review separation, stable paging totals, guarded forms, durable run links, blockers, and recovery.
  - [x] Make browser dry-run use prepare and confirmation but return the non-persisted map result directly without a durable run or state write. Use Post/Redirect/Get only for start, resume, and recovery mutations.
  - [x] Update `internal/web/static/js/operations.js` only to serialize the same shard option, follow existing run events, coalesce bounded progress, and reload authoritative state after terminal delivery. Add no QA registry, client authority, or JavaScript-only action.
  - [x] Preserve loopback, Host, Origin, session, CSRF, body, timeout, confirmation, and strict decoding controls. Escape theory, evidence, error, path, and runtime text and keep raw provider, environment, command output, and absolute paths out of responses.
  - [x] Bound server view models before template execution, support 320 CSS pixels, keyboard and focus order, native controls, table captions and headers, visible focus, reduced motion, high contrast, and live announcements no more than every two seconds.
  - [x] Add `internal/web/qa_handlers_test.go` and contract fixtures for routes, methods, authorization, strict input, no-JavaScript states, hostile content, focus, mobile layout, bounds, cancellation, reconnect, replay gap, session rotation, server restart, terminal refresh, and cross-surface parity.

- [x] **Task 12: Build the offline fault, race, and parity matrix**
  > Executes: Decisions 2-10; Testing; Security; Observability; `REQ-STATE`, `REQ-MAP`, `REQ-READONLY`, `REQ-RUN`, `REQ-THEORY`, `REQ-SYNTHESIS`, `REQ-PARITY`
  - [x] Create one canonical semantic fixture and project it through sprint state, app DTOs, CLI text, CLI JSON, TUI, HTTP JSON, server HTML, and durable run detail. Compare fingerprints, coverage, shard progress, outcomes, synthesis, blocker, cancellation, terminal result, and next action.
  - [x] Use pure tests for canonicalization, IDs, maps, limits, theory validation, synthesis, and invalidation; temporary real workspaces for paths, symlinks, permissions, strict JSON, atomic writes, identity, and preservation; fake runtime and process seams for permissions, malformed output, queueing, timeouts, output caps, retries, cancellation, and cleanup.
  - [x] Add race tests for scheduler admission, cancellation versus promotion, target drift windows, writer fencing, state publication, observer replacement, terminal refresh, and dropped delivery recovery.
  - [x] Extend runtime permission adapter tests so complete QA path rules, required capabilities, unsupported enforcement, and default-deny translation reach Agentwrap unchanged.
  - [x] Extend Sprint 35 regressions for acceptance-before-child, explicit cancellation, event order, redaction, bounded progress, replay, restart recovery, persistence degradation, and product-authority separation.
  - [x] Add fuzz tests for verification-scoped IDs, strict operation JSON, CLI mode combinations, path containment, command classification, and cursor parsing where the existing test style supports them.
  - [x] Scan persisted and public output for forbidden issue, repair, verdict-promotion, provider payload, secret, environment, raw command-output, absolute-path, generated-check, and Git-mutation content.
  - [x] **Stop condition:** no timing-only, provider-required normal test, adapter-owned authority, stale-writer success, unbounded collection, or inferred success from run delivery may remain before documentation and release work closes.

- [x] **Task 13: Document the shipped read-only QA contract**
  > Executes: Decisions 1-10; Documentation; CLI Surface; Persistence And Migrations; `REQ-COMPAT`, `REQ-STATE`, `REQ-PARITY`, `REQ-WEB`, `REQ-RELEASE`
  - [x] Update `docs/cli-reference.md` with Conformance Review terminology, the compatibility alias, all QA commands and flags, JSON envelopes, exits, outcomes, focused shards, cancellation, resume, and recovery.
  - [x] Update `docs/architecture.md` with verification-phase ownership, the four authorities, state paths and publication order, writer tokens, deterministic mapping, approved-check policy, target identity, synthesis limits, and adapter boundaries.
  - [x] Update `docs/user-guide.md` with dry-run mapping, changed-path coverage, investigations, every theory outcome, focused execution, blockers, cancellation, resume, and interpretation of `Read-only QA completed`.
  - [x] Update `docs/local-web.md` with versioned resources, no-JavaScript views, guarded starts, durable progress, explicit cancellation, refresh, reconnect, replay gaps, session rotation, and recovery.
  - [x] Update `docs/recovery.md` with unknown schema, invalid pointer or digest, stale inputs, identity drift, stale writers, interrupted work, cleanup uncertainty, persistence degradation, cancellation, retention, safe pruning, resume, and blocked dogfood.
  - [x] Update `docs/phase3-json-schemas.md` with additive Conformance Review metadata, QA public DTOs and stable errors, exact state layout, schema v1 and `qa-v1` scope, migration policy, fingerprints, limits, and compatibility rules.
  - [x] Update `docs/release-checklist.md` with deterministic map bytes, exact primary ownership, read-only permission and identity evidence, state fault injection, parity, race, build, diff, protocol, and gated dogfood checks.
  - [x] Cross-check every documented command, value, path, limit, status, outcome, error, and route against executable fixtures. Mark Sprint 37 evidence production, smoke integration, adjudication, issues, and Sprint 38 repair as future work.

- [/] **Task 14: Run release gates, reviews, and gated real-repository dogfood** — Deferred: offline release checks passed; current Conformance Review and gated real-runtime dogfood remain post-execution prerequisites.
  > Executes: Decision 10; Testing; Security; Observability; Documentation; `REQ-RELEASE`
  - [x] Run focused sprint, app, TUI, web, runtime-policy, run-control regression, and race suites. Fix deterministic, authority, safety, compatibility, and leak failures before full release commands.
  - [x] Run `go test ./...`, `go test -race ./...`, `go vet ./...`, `go build ./cmd/ultraplan`, and `git diff --check` in `../ultraplan-go`; retain command output or CI links as implementation evidence.
  - [/] Run Architecture Review with the verification-phase split, dependency direction, state authorities, writer fencing, approved-check boundary, web import boundary, and absence of new frameworks or alternate persistence. Dry-run preflight is the execution handoff boundary; no verdict is claimed here.
  - [/] Run Sprint Review with deterministic mapping, strict state, theory and synthesis behavior, read-only enforcement, cancellation/resume/recovery, compatibility, parity, docs, and scope-exclusion evidence. Dry-run preflight is the execution handoff boundary; no verdict is claimed here.
  - [/] After current review prerequisites permit it, run the gated Deep Smoke Sprint protocol against `../ultraplan-go`. The current Conformance Review and real-runtime QA evidence do not yet exist, so this gate is blocked and remains incomplete.
  - [x] If runtime, Sprint 35 evidence, execute evidence, current Conformance Review evidence, browser, Git, or environment prerequisites are missing, record `blocked`. Do not count that as dogfood completion or replace it with fake success.
  - [x] Inspect the final diff and evidence for forbidden generated checks, fixtures, probes, smoke migration, `qa.md`, issues, repairs, product SQLite, plugins, daemons, remote workers, retrieval, content identity, cloud, or Git mutation.
  - [/] **Stop condition:** Sprint 37 cannot start until deterministic mapping, bounded read-only investigation, cancellation, resume, fingerprint invalidation, synthesis, cross-surface parity, and clean target identity have passing real-runtime evidence. Writable-isolation work remains deferred.

## Evidence Checklist

- [x] Phase and compatibility tests prove `VerificationPhase` is independent and `review` remains byte- and behavior-compatible under the Conformance Review label and alias.
- [x] Domain tests prove schema v1, `qa-v1` IDs, all phase and theory values, required theory fields, positive limits, and explicit fingerprint invalidation.
- [x] Mapping tests prove normalized byte stability, grounded input traces, exactly one primary owner per changed path, explicit bounded overlap, and visible unknown-path blockage.
- [x] State tests prove private contained paths, strict decoding, atomic publication order, prior-state preservation, pointer and digest reconciliation, retention, quota, cancellation, resume, and stale-writer rejection.
- [x] Runtime and process tests prove default deny, read-only enforcement, no shell or Git, closed explicit argv, bounded environment and output, context/check limits, cleanup, and no child work when acceptance or policy enforcement fails.
- [x] Identity tests prove no change to source, tests, governed inputs, verification code, harness content, or Git state and expose bounded drift without cleanup.
- [x] Theory and synthesis tests retain every outcome, deduplicate deterministically, preserve contradictions and interactions, cap parent-linked follow-up, and produce no issue, repair, patch, `qa.md`, or review verdict.
- [x] Durable-operation tests prove acceptance and claim before child work, canonical cancellation, event bounds, replay, persistence degradation, new-run resume, and observer-independent execution.
- [x] App, CLI, JSON, TUI, browser HTML/JSON, and durable run detail pass one canonical parity fixture.
- [x] TUI and browser tests cover no-JavaScript or non-TTY fallback, keyboard and focus behavior, narrow layouts, text-only status, hostile content, bounded rendering, reconnect, replay gaps, session rotation, and authoritative terminal refresh.
- [x] Documentation agrees with executable commands, paths, schemas, values, bounds, outcomes, errors, routes, and recovery actions.
- [/] Architecture Review, Sprint Review, and Deep Smoke Sprint evidence is not current. The real-runtime gate is truthfully blocked and remains incomplete.
- [x] Deviations from `reasoning.md` are recorded through the governed reasoning and planning stages before implementation continues.

## Verification Commands

Commands run from `../ultraplan-go` unless noted otherwise. Focused names are implementation deliverables and must exist before their task closes.

| Check | Command | Expected Result |
| --- | --- | --- |
| Sprint QA domain and state | `go test ./internal/sprint -run 'Test(VerificationPhase|QAState|QAID|QABudget|QATheory)'` | Phase compatibility, schemas, IDs, strict state, limits, theory validation, atomic writes, and invalidation pass offline. |
| Deterministic map | `go test ./internal/sprint -run 'TestQAMap'` | Repeated maps are byte-identical and every changed path has one primary owner with only explicit bounded overlap. |
| Investigation and synthesis | `go test ./internal/sprint -run 'TestQA(Investigation|Permission|Identity|Cancellation|Resume|Recovery|Synthesis)'` | Permission, identity, scheduling, stop, outcome, follow-up, cancellation, resume, recovery, and non-promotion contracts pass. |
| Runtime policy adapter | `go test ./internal/platform/runtime -run 'Test.*Permission'` | Default deny, path rules, required capability, and unsupported-policy behavior reach Agentwrap correctly. |
| App use cases and durable acceptance | `go test ./internal/app -run 'Test(QA|.*Durable.*QA|EveryRuntimeBackedCLIEntry)'` | Typed projections, preparation, no-child failure, run correlation, progress, cancellation, resume, and recovery pass. |
| CLI contracts | `go test ./internal/app -run 'TestSprint(QA|ReviewAlias)'` | Help, flags, text, JSON, exits, `sprint.qa`, and review alias compatibility match fixtures. |
| TUI QA views | `go test ./internal/tui -run 'TestQA'` | Routes, keys, vocabulary, bounds, hostile text, cancellation, refresh, and parity pass at wide and narrow sizes. |
| Web QA resources and presentation | `go test ./internal/web -run 'TestQA'` | Versioned routes, strict requests, authorization, no-JavaScript snapshots, bounds, safety, reconnect, recovery, and parity pass. |
| Web compatibility and import boundary | `go test ./internal/web -run 'Test(BrowserOperationKindContract|WebImportBoundary|APICompatibility|Security|SSE)'` | QA additions preserve operation compatibility, app-only imports, browser security, and existing run delivery semantics. |
| Focused race checks | `go test -race ./internal/sprint ./internal/app ./internal/tui ./internal/web` | Scheduler, state promotion, fencing, cancellation, terminal refresh, and observer paths are race-free. |
| Full deterministic suite | `go test ./...` | All normal tests pass without network, OpenCode, ambient credentials, or the external smoke harness. |
| Full race suite | `go test -race ./...` | All packages pass under the race detector with no stale-writer success or leaked workers. |
| Static analysis | `go vet ./...` | Vet reports no errors. |
| CLI build | `go build ./cmd/ultraplan` | The CLI builds successfully. |
| Diff hygiene | `git diff --check` | No whitespace errors are reported. |
| QA command help | `go run ./cmd/ultraplan sprint --help` | QA and Conformance Review command families and compatibility labels render without runtime construction. |
| Architecture and sprint review | `ultraplan sprint ultraplan-go 36-read-only-qa review` | Current Conformance Review evaluates the selected protocols and records an acceptable result or actionable findings without changing QA state. |
| Gated real-runtime dogfood | `ultraplan sprint ultraplan-go 36-read-only-qa smoke` | After review, dogfood proves map stability, read-only outcomes, synthesis, progress, cancellation/recovery, and clean identity, or reports a truthful blocker that leaves the gate incomplete. |

## Risks And Blockers

| Risk / Blocker | Source | Mitigation | Status |
| --- | --- | --- | --- |
| Sprint 35 source behavior exists, but required prior artifact or dogfood evidence may be unavailable. | `reasoning.md#assumptions-and-risks` | Preflight current run-control capabilities and required evidence. Block current QA and dogfood rather than inventing proof. | open |
| A current failed or blocked Conformance Review is useful diagnostic input, while missing or stale review evidence is not. | `reasoning.md#assumptions-and-risks` | Test freshness separately from verdict acceptability and never alter the verdict. | open |
| App composition cannot supply a claimed fencing generation before child work. | `reasoning.md#assumptions-and-risks` | Block runtime QA until the opaque writer token is available and checked; prove no-child behavior and stale-writer rejection. | open |
| Multi-file publication can leave detailed state and flow summary out of sync. | `reasoning.md#assumptions-and-risks` | Publish pointer-last, store digests, preserve prior valid state, expose mismatch, and require explicit recovery. | accepted |
| An approved check writes despite a benign name or changes between mapping and launch. | `reasoning.md#assumptions-and-risks` | Use a closed fingerprinted descriptor catalog, immediate revalidation, sandboxing, adversarial fixtures, and before/after identity. | open |
| Target identity hashing is expensive or misses a dispatch race. | `reasoning.md#assumptions-and-risks` | Stream sorted manifests, recheck before launch, compare after each attempt, and avoid an unproven cache. | accepted |
| Cancellation races with final shard promotion. | `reasoning.md#assumptions-and-risks` | Define promotion ordering around validation, cancellation state, writer token, and atomic rename; race-test both winners. | open |
| Run terminal persistence succeeds while QA summary persistence fails. | `reasoning.md#assumptions-and-risks` | Present both authorities, mark QA invalid or stale, and offer recovery instead of inferring semantic success. | accepted |
| Fixed limits or 128 MiB storage block legitimate work. | `reasoning.md#assumptions-and-risks` | Report exact exhaustion, block before child work, measure dogfood, and change limits only through a fingerprinted governed revision. | accepted |
| Recovery pruning erases useful negative evidence. | `reasoning.md#assumptions-and-risks` | Prune only validated non-current attempts, preserve current and newest last-complete attempts, and record removals. | open |
| Human labels drift or `completed` is rendered as pass. | `reasoning.md#assumptions-and-risks` | Centralize canonical app fields and next actions and freeze one cross-surface parity fixture. | open |
| Challenger output varies across fresh runtime executions. | `reasoning.md#assumptions-and-risks` | Persist and fingerprint validated challenge records, then keep final synthesis pure and metadata-free. | accepted |
| Later Sprint 37 or 38 examples leak into this implementation. | `reasoning.md#assumptions-and-risks` | Inspect diff and docs for generated evidence, smoke integration, adjudication, issues, `qa.md`, repair, and writable isolation; defer all of them. | open |
| Browser or real-runtime prerequisites are unavailable. | `reasoning.md#expected-evidence` | Keep normal tests offline and record the gated real-boundary result as blocked without satisfying the exit criterion. | open |

## Review Inputs

Review should use:

- `sprint-index.md`
- `technical-handbook.md`
- `reasoning/api-design.md`
- `reasoning/architecture.md`
- `reasoning/frontend.md`
- `reasoning.md`
- this `plan.md`
- implementation diff
- deterministic map and state fixtures
- permission, process, identity, and fault-injection evidence
- cross-surface parity and browser accessibility evidence
- durable run cancellation, replay, fencing, and recovery evidence
- release command output and gated dogfood records
- `system/protocols/architecture-review-protocol.md`
- `system/protocols/review-sprint-protocol.md`
- `system/protocols/deep-smoke-sprint-protocol.md`

## Execution Log

| Date / Step | Action | Evidence / Notes |
| --- | --- | --- |
| 2026-08-23 / planning | Created the implementation plan from validated requirements, code context, sprint index, technical handbook, area reasoning, final reasoning, project docs, roadmap, and the resolved plan prompt. | Plan only. No implementation, review, smoke, issue, Git, or run-state work was executed. |
| 2026-08-24 / implementation | Implemented the read-only QA domain, deterministic map, strict private state, fenced orchestration, permission/check policy, identity checks, pure synthesis with bounded challenger inputs and follow-ups, durable operations, CLI, TUI, HTTP/HTML, stable public errors, and documentation in dedicated worktree `36-read-only-qa`. | Focused suites, `go test ./...`, `go test -race ./...`, `go vet ./...`, `go build ./cmd/ultraplan`, `node --check internal/web/static/app.js`, and `git diff --check` passed. No review or smoke verdict is claimed. |
| 2026-08-24 / gated dogfood | Evaluated the real-boundary preconditions after reconciling execute evidence. | Blocked: Sprint 36 has no current Conformance Review or real-runtime QA record, so Deep Smoke Sprint and clean real-runtime target-identity evidence remain incomplete. |

## Completion Criteria

- [x] All tasks are complete or explicitly deferred with a governed reason.
- [x] Verification commands were run, or real-boundary blockers are recorded without being counted as pass evidence.
- [x] Evidence satisfies every expectation from `reasoning.md`, including deterministic bytes, strict state, read-only identity, durable cancellation/resume/recovery, synthesis bounds, and adapter parity.
- [x] `review.md` can evaluate phase separation, compatibility, authority, safety, limits, non-promotion, and scope exclusions without guessing intent.
- [/] The gated real-repository run is blocked and does not yet prove a clean target identity. Sprint 36 remains incomplete for promotion into Sprint 37.
