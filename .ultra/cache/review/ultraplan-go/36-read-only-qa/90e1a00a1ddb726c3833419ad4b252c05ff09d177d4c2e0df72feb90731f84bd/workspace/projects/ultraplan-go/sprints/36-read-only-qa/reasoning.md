# Sprint Reasoning: Read-only QA decomposition and synthesis

> Project: `ultraplan-go`
> Sprint: `36-read-only-qa`
> Output: `projects/ultraplan-go/sprints/36-read-only-qa/reasoning.md`
> **Inputs Used:** `projects/ultraplan-go/project-index.md`, `projects/ultraplan-go/sprints/36-read-only-qa/requirements.md`, `projects/ultraplan-go/roadmap.md`, `projects/ultraplan-go/docs/ARCHITECTURE.md`, `projects/ultraplan-go/docs/PRD.md`, `projects/ultraplan-go/docs/TRD.md`, `projects/ultraplan-go/sprints/36-read-only-qa/sprint-index.md`, `projects/ultraplan-go/sprints/36-read-only-qa/technical-handbook.md`, `projects/ultraplan-go/sprints/36-read-only-qa/reasoning/api-design.md`, `projects/ultraplan-go/sprints/36-read-only-qa/reasoning/architecture.md`, `projects/ultraplan-go/sprints/36-read-only-qa/reasoning/frontend.md`

This document decides the Sprint 36 architecture. It narrows broader PRD and TRD examples to the read-only scope fixed by `requirements.md`. It does not replace the sprint index, handbook, or area reasoning.

Requirement labels such as `REQ-PHASE` are local traceability names for the corresponding acceptance criteria in `requirements.md`. They do not add requirements.

## Sprint Purpose

- **Goal:** Add a deterministic, durable, read-only QA phase that maps changed behavior into bounded verification shards, retains resumable theory outcomes, and synthesizes cross-shard evidence across CLI, JSON, TUI, browser, and durable-run views.
- **Non-Goals:** Generated tests, fixtures, probes, smoke scenarios, `qa.md`, issue promotion, evidence adjudication, repair eligibility, production repair, smoke-as-QA integration, alternate product persistence, content identity, retrieval, hosted operation, remote workers, and Git mutation.
- **Depends On:** Sprint 35 durable run acceptance, fencing, cancellation, replay, recovery, and retention; current execute evidence; a current Conformance Review record; Agentwrap permission enforcement; and the existing app operation boundary.

## Selected Context And Pre-Reasoning Artifacts

| Artifact | Path | How It Was Used |
| --- | --- | --- |
| Sprint requirements | `requirements.md` | Fixed the read-only scope, required files, compatibility rules, state layout, outcome vocabulary, cross-surface behavior, and release gates. |
| Sprint index | `sprint-index.md` | Selected the 13 contracts, 14 study reports, three area analyses, and three post-implementation review protocols. |
| Technical handbook | `technical-handbook.md` | Supplied comparative evidence for thin adapters, explicit configuration, typed errors, durable state, cancellation, bounded concurrency, command safety, testing, and hard resource limits. |
| Project architecture | `../../docs/ARCHITECTURE.md` | Fixed module ownership, dependency direction, filesystem product authority, app adapter boundaries, and the Phase 5 sequence. |
| Product requirements | `../../docs/PRD.md` | Fixed user-visible QA behavior, compatibility, local-only browser operation, KPI expectations, and later QA and repair gates. |
| Technical requirements | `../../docs/TRD.md` | Fixed schema, persistence, runtime, security, cancellation, adapter, and verification requirements. |
| Product roadmap | `../../roadmap.md` | Kept Sprint 36 before writable QA, adjudication, and repair, and carried forward the Sprint 35 release dependency. |
| Prior decision | None cataloged | The project index has no formal prior-decision record. Sprint 35 behavior and the current project documents are binding dependencies instead. |

## Area-Specific Reasoning Inputs

| Area | Reasoning Document | Key Conclusion | Evidence Basis | Impact On Final Decision |
| --- | --- | --- | --- | --- |
| Architecture | `reasoning/architecture.md` | Keep one sprint-owned QA state machine, one existing run-control system, a pointer-only flow summary, explicit writer fencing, product-owned check execution, and bounded scheduling. | Project architecture plus handbook evidence on one-way dependencies, durable state, cancellation, concurrency, security, and filesystem testing. | Fixes package ownership, detailed state, atomic publication, runtime safety, cancellation, resume, and limits. |
| API Design | `reasoning/api-design.md` | Expose typed app projections and a closed QA command and HTTP contract. Keep product outcomes separate from transport errors and Conformance Review verdicts. | CLI, DI, configuration, error, state, observability, security, and performance reports plus current operation APIs. | Fixes app capabilities, command family, operation kinds, public DTOs, error categories, and adapter parity. |
| Frontend | `reasoning/frontend.md` | Keep QA under sprint navigation, render authoritative bounded snapshots without JavaScript, and treat live delivery as disposable observation state. | Current TUI and web implementation plus terminal UX, state, observability, testing, security, and performance reports. | Fixes routes, labels, keyboard controls, reconnect behavior, accessibility, hostile-content handling, and presentation bounds. |

The final decisions adopt these conclusions and resolve their open seams. The additions are an explicit shard-detail use case and route, no request-level model or budget overrides, runtime-free recovery with no child work, one machine outcome spelling, explicit terminal-result handling, a non-persisting browser dry-run response, and numeric scheduler and storage bounds.

## Sprint Technical Handbook Summary

- **Relevant Patterns:** One domain authority behind thin adapters; lazy runtime construction for read-only paths; explicit merged configuration; cause-preserving error classes; separation of durable state, sessions, cancellation, and locks; bounded workers and cleanup; one progress model; product-enforced command policy; deterministic fakes plus gated real boundaries; and limits on data as well as workers.
- **Important Trade-Offs:** Detailed disk state adds schema and atomicity work; narrow process and runtime seams support fault tests without a false general persistence abstraction; collecting independent shard outcomes costs runtime but preserves negative evidence; explicit argv supports fewer checks than a shell; and fixed bounds can block large changes but make authority and resource use reproducible.
- **Warnings / Anti-Patterns:** Do not duplicate workflow rules in adapters, use package globals for QA dependencies, continue after invalid configuration, flatten errors into strings, equate context presence with cancellation, launch one goroutine per shard, bind cancellation to observers, mix diagnostics with JSON, trust raw paths or shell text, or treat pooling and truncation as hard bounds.
- **Evidence Confidence:** High for the general Go patterns because the handbook cites several mature repositories and exact source locations. Medium for the selected numeric values, multi-file publication order, writer token, and QA vocabulary because those are UltraPlan-specific decisions that require direct tests and dogfood rather than imitation.

## Contracts Applied

| Contract / Requirement ID | Constraint | Decision Impact | Expected Evidence |
| --- | --- | --- | --- |
| Architecture; `REQ-PHASE` | QA is a verification phase owned by `internal/sprint`, not a planning stage. Adapters remain thin. | Add `VerificationPhase`; keep `PlanningStages()` unchanged; expose QA through `internal/app`; retain the web import boundary. | Architecture import tests, phase enumeration tests, and review of package dependencies. |
| CLI Surface; Documentation; `REQ-COMPAT` | People see Conformance Review while `review`, `review.md`, existing JSON, verdicts, and clients remain compatible. | Add `conformance-review` only as an alias to the existing review handler and state. | Help and command fixtures proving identical handler, artifact, operation name, output, and exits. |
| Persistence And Migrations; Errors; `REQ-STATE` | Detailed QA state is versioned, contained, atomic, bounded, and separate from flow state. | Use schema v1 files under `verification/`; publish pointer-last; reject unknown majors and malformed state; make recovery explicit. | Strict-load, path, mode, rename, fsync, pointer, digest, stale-writer, and recovery tests. |
| Workflows; Observability; `REQ-RUN` | Runtime QA uses Sprint 35 durable acceptance, ownership, cancellation, events, and recovery. | Accept and claim before child work; pass an opaque writer token; use the existing run cancellation path; create a new run on resume. | Durable acceptance inventory, fencing races, cancellation, replay, restart, and persistence-degradation fixtures. |
| Security; LLM Runtime; LLM Evaluation / Cost / Safety; `REQ-READONLY` | Investigators cannot write or choose commands, paths, environment, prompts, permissions, or outputs. | Agentwrap gets read-only default-deny policy; product code executes only map-owned explicit-argv checks; before/after identity gates promotion. | Adversarial permissions, path and symlink escape, shell, Git, cache-write, unsupported-policy, redaction, and identity-drift tests. |
| Performance; Configuration; `REQ-BOUNDS` | Every map and shard has positive finite limits whose exhaustion is visible. | Freeze defaults and maxima; configuration may lower them only; include effective values and sources in fingerprints; add queue, retention, and storage caps. | Validation tables, boundary tests, queue saturation, oversized output, timeout, retention, quota, and deterministic fingerprint tests. |
| Testing; `REQ-MAP` | Unchanged normalized inputs produce stable map identity and every changed path has one primary owner. | Use sorted canonical inputs, deterministic IDs, explicit boundary overlap, and visible blocked ownership for unknown paths. | Byte comparisons, primary ownership assertions, overlap fixtures, risk-input fixtures, and map fingerprint tests. |
| Testing; Workflows; `REQ-THEORY` | Theories are falsifiable and all positive, negative, invalid, blocked, and cross-shard outcomes remain durable. | Validate the complete theory contract and preserve every closed outcome value. | Outcome table tests, malformed theory rejection, resume, invalidation, and synthesis retention tests. |
| Architecture; LLM Safety; `REQ-SYNTHESIS` | Synthesis is bounded and verdict-neutral. It cannot promote issues or repairs. | Use deterministic product-owned normalization, grouping, ordering, contradiction detection, and follow-up selection over validated records. | Pure synthesis fixtures, deduplication, contradiction, parent links, limit exhaustion, and forbidden-field scans. |
| Observability; CLI Surface; `REQ-PARITY` | CLI text, JSON, TUI, browser, and run inspection report the same canonical facts. | Define one app DTO vocabulary and one parity fixture. Keep phase, freshness, review verdict, cancellation, terminal result, and observation as separate axes. | Shared fixture checks across every adapter and durable run detail. |
| Security; Documentation; `REQ-WEB` | Browser requests remain local, guarded, strict, escaped, bounded, and recoverable without JavaScript. | Reuse prepare, confirmation, run, cancellation, and SSE APIs; add focused QA resources; keep server-rendered snapshots complete. | Host, Origin, CSRF, session, hostile-content, no-JS, reconnect, replay-gap, focus, mobile, and reduced-motion tests. |
| Testing; Documentation; `REQ-RELEASE` | Offline tests, race, vet, build, diff, and gated real-runtime dogfood all pass. | Treat missing dogfood prerequisites as blocked, not pass, and require clean target identity evidence. | Required commands, protocol reviews, dogfood artifacts, and release-checklist inspection. |

## Repos Studied / Source Evidence Used

| Source / Repo / Report | Concrete Reference | Relevant Finding | Why It Matters For This Sprint | Used In Decision(s) |
| --- | --- | --- | --- | --- |
| Chezmoi and yq, `01-project-structure` | `chezmoi/main.go:16`, `chezmoi/internal/chezmoi/chezmoi.go:1-2`, `yq/cmd/root.go:9` | Command code depends on protected domain code in one direction. | Supports sprint-owned semantics and adapter-only CLI, TUI, and web code. | Phase and ownership; shared app boundary. |
| Gh-cli, `02-command-architecture` and `03-dependency-injection` | `gh-cli/pkg/cmdutil/factory.go:16-43`, `gh-cli/pkg/cmd/factory/default.go:26-46` | Injected and lazy command dependencies keep status paths from opening unrelated services. | Supports one app capability with runtime-free map and status paths. | Shared app boundary; dry-run and recovery. |
| Chezmoi, restic, and opencode, `04-configuration-management` | `chezmoi/internal/cmd/config.go:2253-2287`, `restic/internal/global/global.go:139,147`, `opencode/internal/config/config.go:609-641` | Effective configuration needs explicit precedence and one validation pass. | Supports immutable effective QA settings and fail-closed limits. | Configuration and hard bounds. |
| Rclone and restic, `05-error-handling` | `rclone/fs/fserrors/error.go:22-192`, `restic/internal/errors/fatal.go:10-53` | Errors can retain causes while carrying retry, stop, or fatal behavior. | Supports stable QA blocker and transport categories without string matching. | State vocabulary; public errors; recovery. |
| Restic, `06-io-abstraction` | `restic/internal/fs/interface.go:10-31`, `restic/internal/backend/backend.go:19-90` | Narrow I/O boundaries permit deterministic faults, but broad abstractions have real cost. | Supports concrete filesystem QA storage with narrow atomic-write and identity seams. | Detailed state and testing. |
| Helm, go-task, opencode, chezmoi, and restic, `07-state-context` | `helm/pkg/cmd/install.go:333-347`, `go-task/task.go:438-469`, `opencode/internal/session/session.go:12-23`, `chezmoi/internal/chezmoi/boltpersistentstate.go:26-31`, `restic/internal/restic/lock.go:105,290-305` | Cancellation, active-work deduplication, sessions, durable state, locks, and cleanup solve different problems. | Prevents using process context, browser state, or sessions as restart authority. | Durable run reuse; cancellation; resume; cleanup. |
| K9s, gdu, opencode, and gh-cli, `08-concurrency` | `k9s/internal/pool.go:21,30,37`, `gdu/pkg/analyze/parallel.go:13,36`, `opencode/cmd/root.go:261-279`, `gh-cli/pkg/cmd/extension/manager.go:196-206` | Bounded admission and timed shutdown are safer than goroutine-per-item fan-out. | Supports finite workers, pending queue, whole-run deadline, and cancellation while queued. | Scheduler and hard bounds. |
| Restic and gh-cli, `09-terminal-ux` | `restic/internal/ui/termstatus/status.go:197-205`, `gh-cli/pkg/iostreams/iostreams.go:514-516` | One progress owner reduces rendering races, and observer output can stop while work continues. | Supports canonical progress plus observer-independent execution. | Cross-surface progress and disconnect behavior. |
| Helm, k9s, and restic, `10-logging-observability` | `helm/internal/logging/logging.go:31-71`, `k9s/internal/slogs/keys.go:6-231`, `restic/internal/backend/logger/log.go:22-77` | Structured fields and stream separation improve diagnosis but do not prove redaction. | Supports bounded app-owned events and explicit pre-persistence redaction. | Observability and public safety. |
| Chezmoi, lazygit, restic, and dive, `11-testing-strategy` | `chezmoi/internal/cmd/main_test.go:64-174`, `lazygit/pkg/commands/oscommands/fake_cmd_obj_runner.go:17-26`, `restic/internal/backend/mock/backend.go:14-26`, `dive/dive/image/docker/testing.go:12-34` | Script fixtures, recording fakes, fault seams, and selected real boundaries complement each other. | Supports deterministic offline tests plus separately gated dogfood. | Evidence and release gates. |
| Dive, gh-cli, go-task, and restic, `12-extensibility` | `dive/cmd/dive/cli/internal/command/adapter/analyzer.go:13-15`, `go-task/executor.go:20-24,91-122`, `restic/cmd/restic/main.go:77-106` | Narrow adapters and static registries provide substitution without plugins. | Supports a closed in-tree approved-check catalog instead of a QA framework. | Check policy and package scope. |
| Restic, opencode, and chezmoi, `13-security` | `restic/internal/backend/shell_split.go:45-76`, `opencode/internal/permission/permission.go:44-108`, `chezmoi/internal/cmd/gpgencryption.go:151-165` | Explicit arguments, permission gates, containment, private storage, and redaction make trust boundaries reviewable. | Supports no shell, default deny, explicit argv, `0600` state, and escaped output. | Read-only policy and browser safety. |
| Age, restic, rclone, and opencode, `14-performance` | `age/internal/stream/stream.go:20,195-219`, `restic/internal/archiver/file_saver.go:56-58`, `rclone/lib/pool/pool.go:17-24,52-53`, `opencode/internal/message/message.go:37-42` | Hard worker, queue, memory, stream, and disk-state bounds matter on large inputs. | Supports numeric limits for paths, prompts, output, history, storage, and duration. | Hard bounds and retention. |

## Trade-Off And Debt Analysis

### Accepted Trade-Offs

| Trade-Off | Benefit | Cost / Constraint Accepted | Why Acceptable Now | Revisit Trigger |
| --- | --- | --- | --- | --- |
| Separate detailed files with pointer-last publication | Isolates shard writes, keeps maps immutable, and keeps summaries bounded. | Multi-file consistency needs digests, ordering, and reconciliation instead of one transaction. | Filesystem artifacts remain product authority and Sprint 35 SQLite has a different role. | A later persistence gate proves product SQLite has concrete value. |
| Opaque writer token between app and sprint | Blocks stale workers without importing run-control storage types into sprint. | Adds correlation and promotion checks to every write. | It preserves dependency direction and one operational authority. | Composition cannot provide a claimed fencing generation before child work. |
| Product-owned explicit-argv checks | Makes commands, cwd, environment, timeout, and output policy reviewable. | Some checks that require shell features remain unavailable and block rather than run. | Read-only assurance matters more than broad command compatibility. | Sprint 37 proves a stronger isolated writable environment. |
| Independent shard continuation after local failure | Preserves refuted, invalid, inconclusive, and unrelated completed evidence. | May spend runtime after one shard fails. | Whole-attempt safety failures still stop all admission. | Dogfood shows repeated systemic failures misclassified as local. |
| Fixed conservative limits | Makes mapping, authority, storage, and execution reproducible. | Legitimate large changes may block. | Silent adaptive expansion would change fingerprints and investigator authority. | Measured dogfood shows a repeated safe limit is too low. |
| New durable run for resume | Keeps ownership and terminal history honest while reusing semantic evidence. | One QA attempt may correlate to several run IDs. | Dead workers and sessions must not be adopted. | Run control later defines a proven cross-process adoption contract. |
| Server refresh after terminal delivery | Shows only persisted QA and run-control facts. | Feels slower than optimistic client completion. | Delivery is not authority, especially after cancellation or persistence races. | None while filesystem QA state remains authoritative. |
| Focused routes and one-column TUI | Bounds initial reads and reuses current navigation. | Adds route fixtures and more navigation. | Avoids a client application and a second TUI focus model. | Real users cannot inspect bounded evidence efficiently. |

### Potential Technical Debt

| Debt / Shortcut | Why It Might Accrue | Current Mitigation | Owner / Follow-Up |
| --- | --- | --- | --- |
| Pointer-only QA summary differs from embedded review and smoke state. | Flow-state code now supports two verification persistence shapes. | Validate and preserve the QA pointer explicitly during planning refreshes; document the authority split. | Sprint 36 implementation and architecture review. |
| Approved check descriptors can drift from real commands. | Renames or changed flags make saved descriptors stale. | Fingerprint the descriptor catalog and invalidate affected shards; never fall back to shell text. | Sprint 36 tests, then Sprint 37 evidence-producing QA. |
| Schema v1 has rejection but no migration history. | The first future major change will need migration machinery. | Freeze v1 fixtures and reserve pure `vN -> vN+1` migration functions. | First shipped QA schema change. |
| Fixed attempt and storage retention can require operator cleanup. | Rich outcomes and command summaries accumulate across maps. | Block before quota breach and make `qa recover` prune only validated non-current attempts. | Dogfood measurement and later persistence decision. |
| Full target identity hashing can be expensive. | Large trees and repeated checks increase I/O. | Stream sorted manifests, apply path bounds, and prefer correctness over an unproven cache. | Add caching only after identity and invalidation are proven. |
| Human labels can drift across TUI and browser. | Separate renderers may rename the same state differently. | Put canonical labels and next actions in app DTOs and freeze a parity fixture. | Sprint 36 adapter tests. |
| Deterministic synthesis limits free-form cross-shard interpretation. | Product-owned grouping can miss semantic relationships not declared by investigators. | Let bounded read-only challenger output become validated, fingerprinted synthesis input, then keep final grouping and follow-up selection pure. | Reassess after Sprint 36 dogfood without allowing adjudication. |

### Future Considerations

| Consideration | Deferred Until | Reason Deferred | What Should Be Preserved Now |
| --- | --- | --- | --- |
| Generated tests, fixtures, probes, and smoke-as-QA | Sprint 37 | Requires writable isolation, cleanup proof, and evidence adjudication. | Stable map, shard, theory, check-request, identity, and run-correlation contracts. |
| Canonical `qa.md`, issue packets, and promotion | Sprint 37 or later | Sprint 36 outcomes are diagnostic and verdict-neutral. | Durable negative outcomes, parent evidence, and explicit non-promotion. |
| Bounded repair | Sprint 38 | Requires adjudicated frozen issues and write isolation. | Independent `repair` phase value without implementation authority. |
| Product SQLite or alternate artifact persistence | Post-gate evidence | Multi-surface use has not shown a product need beyond filesystem artifacts. | App DTOs, schema versions, contained references, and no adapter filesystem access. |
| Global content identity, retrieval, graphs, cloud, or collaboration | After Sprint 39 gates | Outside current assurance scope. | Document that `qa-v1` IDs are verification-scoped only. |
| Request-level model or budget controls | Dogfood evidence | They would expand compatibility and fingerprint complexity before a demonstrated need. | Effective source reporting and immutable policy fingerprints. |

## Decisions

### Decision 1: Separate verification phases and preserve review compatibility

- **Decision:** Add `VerificationPhase` with `conformance-review`, `qa`, and reserved future `repair`. Keep `PlanningStage`, `PlanningStages()`, planning artifact order, `StageReview`, `StageSmoke`, `review`, `review.md`, `smoke`, `smoke.md`, review verdicts, and existing JSON unchanged. Human interfaces say `Conformance Review`. The `conformance-review` CLI alias normalizes to the existing review request and handler and retains `sprint.review` as its JSON operation.
- **Rationale:** QA has different identity, persistence, and lifecycle rules from authored planning stages. A label change must not fork the analytical review capability.
- **Study / Source Grounding:** `technical-handbook.md` cites one-way command-to-domain dependencies in chezmoi and yq and factory routing in gh-cli. The exact phase values and compatibility behavior come from `requirements.md`, `ARCHITECTURE.md`, and the architecture and API area decisions.
- **Trade-Offs Accepted:** A compatibility mapper is extra code, but it prevents a broad rename and duplicate state.
- **Technical Debt / Future Impact:** The reserved `repair` value declares vocabulary, not Sprint 36 behavior. Later repair work must add its own authority and gates.
- **Alternatives Rejected:** Adding QA to `PlanningStage`; renaming or removing `review`; creating a second Conformance Review implementation; wrapping smoke as QA in this sprint.
- **Contracts Satisfied:** Architecture, CLI Surface, Documentation, `REQ-PHASE`, `REQ-COMPAT`.
- **Evidence Required:** Phase type tests, unchanged planning sequence, alias-equivalence command tests, frozen review JSON and artifact fixtures, and architecture review.

### Decision 2: Keep four authorities separate and publish detailed state atomically

- **Decision:** `internal/sprint` owns QA semantics and the concrete detailed store. The authoritative layout is:

```text
verification/state.json
verification/attempts/<attempt-id>/map.json
verification/attempts/<attempt-id>/shards/<shard-id>.json
verification/attempts/<attempt-id>/synthesis.json
```

`verification/state.json` stores schema, current attempt, current fingerprints, verdict-neutral status and freshness, bounded counts, blocker, next action, cancellation summary, run correlation, and referenced digests. `flow-state.json` stores only a bounded QA projection and contained pointer or digest. Run control owns operational acceptance, lease, fencing, cancellation, event order, liveness, and terminal result. The checkout owns source identity. Neither state system absorbs another.

Use schema major `1`, ID marker `qa-v1`, directories `0700`, and files `0600`. Loaders reject missing or unknown versions, unknown closed-contract fields, trailing JSON, unsafe paths, invalid fingerprints, and digest mismatch. No synthetic v0 migration is added.

Each write uses a same-directory temporary file, normalized indented JSON plus newline, flush, close, rename, and directory sync. Referenced files become durable before `verification/state.json`; detailed state becomes durable before `flow-state.json`. Maps are immutable. Status never promotes orphan files by directory scan. Runtime-backed writes require a current opaque token containing run ID, operational attempt ID, and fencing generation.
- **Rationale:** Detailed resumable evidence does not fit bounded flow state or operational run events. Pointer-last publication and writer fencing preserve authority across partial writes and stale workers.
- **Study / Source Grounding:** Restic's filesystem and backend seams support fault testing, while the state-context report distinguishes persistent state, sessions, and locks. Atomic sequence, modes, schema, and writer token are UltraPlan-specific decisions grounded in the Persistence, Security, and Sprint 35 contracts.
- **Trade-Offs Accepted:** Multi-file publication has more failure modes than a transaction. It avoids rewriting unrelated shard evidence and prevents Sprint 35 SQLite from becoming product truth.
- **Technical Debt / Future Impact:** The first real schema change needs an explicit pure migration with fixtures. Pointer reconciliation remains a product concern.
- **Alternatives Rejected:** One large state file; detailed flow-state embedding; product truth in run events; product SQLite; a generic repository or virtual filesystem; read-time repair; direct `internal/sprint -> internal/runcontrol` imports.
- **Contracts Satisfied:** Architecture, Errors, Security, Persistence And Migrations, Workflows, `REQ-STATE`, `REQ-RUN`.
- **Evidence Required:** Strict decoding, version rejection, containment, symlink escape, permissions, normalized bytes, atomic failure at every publication step, prior-state preservation, digest reconciliation, flow-summary mismatch, and stale-writer race tests.

### Decision 3: Build one deterministic map with exact primary ownership

- **Decision:** `qa --dry-run` resolves and validates current requirements, code-context, sprint reasoning and plan, execute state and changed paths, selected contracts and protocols, current Conformance Review findings, adjacent tests and known checks, package and interface boundaries, state transitions, producer and consumer relationships, public APIs, and risk tags. It invokes no runtime, accepts no durable run, and writes no state.

Canonicalization normalizes paths, line endings, identifiers, ordered references, selected limits, implementation identity, check-catalog identity, review fingerprint, and policy fingerprint. It sorts every unordered collection before hashing or persistence. Timestamps, sessions, run IDs, provider metadata, worker order, and current time are excluded from stable identity.

Every changed path has exactly one primary shard. Boundary shards may overlap named paths only for an explicit cross-package, interface, producer-consumer, state-transition, public-API, or cross-cutting concern. Unknown paths become blocked primary shards. Map, shard, theory, and semantic attempt IDs derive from schema marker, project, sprint, and normalized parent inputs. They are stable only inside `qa-v1` and are not global content IDs.
- **Rationale:** Stable ownership is the prerequisite for reproducible investigation, focused reruns, and trustworthy coverage totals.
- **Study / Source Grounding:** The handbook's project-structure evidence supports behavior-owned boundaries, and its performance evidence supports bounded discovery. The exact mapping inputs and ownership rules come from Sprint 36 requirements and area architecture.
- **Trade-Offs Accepted:** Behavior-based deterministic grouping needs explicit classifiers and can block unfamiliar paths. Silent fallback or file-by-file duplication would be less trustworthy.
- **Technical Debt / Future Impact:** New language or repository layouts may require new deterministic classifiers and fixtures. They must not use runtime inference to alter stable mapping.
- **Alternatives Rejected:** Runtime-generated maps; one shard per file; multiple primary owners; orphan omission; timestamps or run metadata in IDs; global content identity claims.
- **Contracts Satisfied:** Architecture, Testing, Performance, `REQ-MAP`, `REQ-BOUNDS`.
- **Evidence Required:** Repeated byte comparisons, stable IDs and ordering, changed-path coverage, single primary ownership, explicit overlap, unknown-path blocking, risk-tag inputs, implementation invalidation, and map fingerprint goldens.

### Decision 4: Freeze configuration, budgets, scheduling, and retention

- **Decision:** Sprint 36 public requests may select only project, sprint, operation mode, and a current map-owned shard ID. They may not override model, variant, budgets, retries, commands, prompts, paths, environment, permissions, or policy. The validated effective workspace configuration selects the QA model and may lower product defaults. Existing environment-backed workspace configuration participates through the current configuration loader. Product maxima cannot be raised. Non-positive or over-maximum values fail instead of clamping. Effective values and sources are frozen into map and policy fingerprints.

Product budgets are:

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

Presentation and delivery bounds are 40 summary shards, 24 theories per focused shard page, 200 TUI durable events, 100 browser QA operation rows, existing 16 KiB durable events, existing 250 ms progress coalescing, and live-region aggregate announcements no more often than every two seconds. Stable cursor paging states `Showing X of Y`; display limits never imply evidence completeness.

The scheduler is run-scoped and uses at most the configured investigator count. It never creates one goroutine per shard. The pending queue cannot exceed the total planned shard maximum. Limit exhaustion records the resource and affected scope and never widens authority.

Before a ninth attempt or a state write over 128 MiB, the operation blocks without starting runtime work. Explicit `qa recover` may prune oldest validated non-current attempts while preserving the current attempt and newest last-complete attempt. Reads never prune.
- **Rationale:** The area analyses fixed most limits but left queue and storage growth unspecified. Numeric bounds are required for reproducibility, safety, rendering, and recovery.
- **Study / Source Grounding:** K9s and gdu support bounded admission, opencode supports timed shutdown, restic demonstrates bounded queues, and rclone couples pooling with a hard cap. Chezmoi, restic, and opencode support explicit configuration merge and validation. Exact values are Sprint 36 decisions.
- **Trade-Offs Accepted:** No per-request tuning and an eight-attempt history are less flexible. They prevent authority expansion and unbounded state before dogfood supplies measurements.
- **Technical Debt / Future Impact:** Later changes to defaults or maxima alter policy fingerprints and require fixture and documentation updates. Retention may need revisiting after real use.
- **Alternatives Rejected:** CLI-only model overrides; browser budget controls; silent clamping; adaptive widening; unbounded queues; pruning during status; relying only on Sprint 35's workspace quota.
- **Contracts Satisfied:** Configuration, Performance, Security, Workflows, `REQ-BOUNDS`, `REQ-STATE`.
- **Evidence Required:** Configuration precedence and source tests, invalid limits, queue saturation, cancellation while queued, all byte and duration limits, attempt and storage caps, safe pruning, paging totals, and fingerprint invalidation.

### Decision 5: Enforce read-only investigation in code and verify target identity

- **Decision:** Investigators receive a bounded prompt packet containing assigned common context, approved shard paths, theory contract, current implementation fingerprint, and map-owned check IDs. Agentwrap requests use `read_only`, restricted permissions, default deny, required permission capability, and only contained read, list, and search access. Unsupported enforcement blocks before useful output promotion.

Investigators have no shell tool. A check request contains only a stable map-owned ID. Sprint code resolves it to an immutable descriptor with executable, argv, cwd, environment-name allowlist, timeout, output cap, and descriptor fingerprint, then invokes the generic process boundary. Reject shell wrappers, command substitution, output redirection, interpreter indirection, Git, caller-supplied executable or argv, caller-supplied environment, and path or symlink escape. Checks known to write caches, coverage, generated output, or temporary files inside protected trees are not approved.

Before dispatch and after each investigator attempt, compute a sorted identity manifest for production source, production tests, governed sprint and project inputs, verification implementation code, smoke-harness content, and Git HEAD, index, and worktree state when applicable. Revalidate implementation identity immediately before process launch. Symlinks record link identity and must resolve inside approved roots. QA and run-control state are excluded. Non-Git targets record Git identity as `not_applicable`; gated implementation-repository dogfood requires Git identity.

Any drift stops admission, blocks the attempt, records the affected category and bounded relative paths, preserves prior canonical outcomes, and performs no cleanup or reset.
- **Rationale:** Prompt wording cannot prove read-only behavior. Permission enforcement, closed command descriptors, containment, and before/after identity provide layered proof.
- **Study / Source Grounding:** Restic's argument handling, opencode's permission gate, go-task's controlled execution, chezmoi's restricted temporary storage, and the handbook's security warnings support a code-owned boundary. The identity categories come from Sprint 36 requirements and architecture reasoning.
- **Trade-Offs Accepted:** Some useful checks remain blocked and identity hashing costs I/O. Safety and truthful blockage take priority.
- **Technical Debt / Future Impact:** The check catalog needs maintenance. Sprint 37 must not weaken this policy when it introduces isolated writable evidence production.
- **Alternatives Rejected:** Read-only prompts with shell access; caller command strings; broad environment forwarding; Git cleanup; writable tests or probes; assuming sandbox support without checking runtime capability.
- **Contracts Satisfied:** Security, LLM Runtime, LLM Evaluation / Cost / Safety, Errors, `REQ-READONLY`.
- **Evidence Required:** Runtime request inspection, unsupported-policy behavior, file and symlink escape corpus, shell and interpreter denial, Git denial, cache-writing checks, environment redaction, before/after manifest comparison, race-window revalidation, and clean dogfood identity.

### Decision 6: Orchestrate durable starts, cancellation, resume, and explicit recovery

- **Decision:** `qa-start` and `qa-resume` use the existing operation preparation and fingerprint revalidation path. App durably accepts and claims a run before constructing child work, then passes a cancellation-aware context and opaque writer token to sprint QA. Failure to persist acceptance or ownership prevents runtime start.

Whole-attempt stop conditions are unavailable or stale governed inputs, unsupported or corrupt state, writer-token loss, unavailable permission enforcement, target drift, cancellation, whole-run timeout, and authoritative persistence failure. Shard-local conditions include malformed investigator output, unavailable or failed approved checks, denied or exhausted context expansion, local budgets, and closed negative outcomes. Independent shards may continue after local failure. Partial synthesis must state blockers and cannot claim complete coverage.

Cancellation uses `RunUseCases.CancelRun`. It stops admission, cancels active investigator and process contexts, preserves completed promoted shards, and records active work truthfully. A separate 30-second cleanup context handles process groups, final state, safe events, and resources. Cleanup timeout sets run terminal result `cleanup_uncertain` and QA phase `interrupted` with a blocker and recovery action. It never reports success. Observer loss never cancels work.

Resume creates a new durable run and writer token. It reuses the semantic attempt and valid completed shards only when governed input, implementation, review, map, check catalog, policy, effective limits, and selected-shard fingerprints remain current. It never adopts sessions, process handles, or dead workers.

`qa recover` is a synchronous, runtime-free verification-state mutation. It uses the sprint mutation lock and normal browser prepare and confirmation authorization, but it creates no investigator runtime, child work, writer token, or durable run ID. It may reconcile stale ownership, interrupted markers, contained pointers and digests, bounded flow summary, and retention. If runtime work is needed, recovery stops and directs the user to `qa resume`.
- **Rationale:** Operational ownership and semantic evidence must correlate without pretending they are one database or lifecycle.
- **Study / Source Grounding:** The state-context report distinguishes context, sessions, persistent state, active deduplication, locks, and cleanup. K9s and opencode support bounded scheduling and shutdown. Sprint 35 supplies the binding acceptance and cancellation contract.
- **Trade-Offs Accepted:** Recovery has confirmation despite no runtime because it mutates verification state. Resume produces another run ID. Both choices keep authority explicit.
- **Technical Debt / Future Impact:** A run may become terminal while QA summary persistence fails. Status must show both facts and require recovery rather than infer completion.
- **Alternatives Rejected:** QA-local cancellation; process-local ownership; browser or TUI lifecycle cancellation; session adoption; read-time recovery; recovery that starts child work; universal fail-fast for local theory outcomes.
- **Contracts Satisfied:** Workflows, Observability, Errors, Persistence And Migrations, `REQ-RUN`, `REQ-STATE`.
- **Evidence Required:** Acceptance-before-child tests, stale writer races, local versus global stop matrix, cancellation queued and running, 30-second cleanup behavior, dropped observer delivery, restart resume, partial invalidation, runtime-free recovery, and terminal/state mismatch fixtures.

### Decision 7: Validate complete theories and synthesize without adjudication

- **Decision:** Each theory records deterministic ID, falsifiable claim, basis, verification surface, expectation references, severity if confirmed, confirmation condition, refutation condition, inconclusive condition, safe evidence strategy, implementation fingerprint, attempt history, and one final machine outcome:

```text
confirmed
refuted
invalid
inconclusive
blocked
cross_shard
not_applicable
```

Hyphenated forms are prose only. Every outcome remains durable and inspectable.

Global synthesis consumes only validated records current for the map. Product code deterministically normalizes equivalence keys, groups duplicates, orders theories, retains contradictions and negative outcomes, records declared interactions, and selects at most four parent-linked follow-up shards from remaining budget. A bounded read-only challenger request may propose challenge records against frozen theory groups, but those records must pass the same schema, permission, identity, and fingerprint gates before becoming explicit synthesis inputs. The final synthesis function is pure and byte-stable for identical normalized inputs. Provider, session, run, worker, and observation metadata never affect synthesis identity.

Synthesis cannot create issues, repair eligibility, patches, generated checks, `qa.md`, or a Conformance Review verdict. A current failed or blocked Conformance Review remains unchanged. `completed` means the bounded investigation and synthesis ended, not pass.
- **Rationale:** QA needs to preserve what was disproved or could not be answered. Central deterministic grouping prevents investigators from promoting their own claims.
- **Study / Source Grounding:** The handbook's error and state evidence supports closed inspectable outcomes, while its performance evidence supports bounded records. The exact theory schema and non-promotion rules come from Sprint 36 requirements.
- **Trade-Offs Accepted:** Product-owned deterministic synthesis is less free-form than unconstrained model prose. Bounded challenger records retain semantic challenge without making runtime output the final authority.
- **Technical Debt / Future Impact:** Sprint 37 may adjudicate evidence, but it must consume these outcomes rather than rewrite historical theory state.
- **Alternatives Rejected:** Dropping refuted or invalid theories; investigator issue promotion; model prose as canonical synthesis; adaptive follow-up counts; review-verdict upgrade; repair actions.
- **Contracts Satisfied:** LLM Evaluation / Cost / Safety, Workflows, Testing, `REQ-THEORY`, `REQ-SYNTHESIS`.
- **Evidence Required:** Full theory validation, all outcome fixtures, deterministic synthesis bytes, deduplication, contradictions, interactions, parent evidence, follow-up exhaustion, stale-input rejection, and scans proving no issue, repair, or verdict fields are produced.

### Decision 8: Define one app, CLI, JSON, and HTTP contract

- **Decision:** `internal/app` exposes typed QA DTOs rather than sprint structs. The capability includes `QAMap`, `QAStatus`, `QAShard`, `QATheory`, `QASynthesis`, `RunQA`, `ResumeQA`, `CancelQA`, and `RecoverQA`. Query methods are read-only and do not create state or initialize runtime. `QAMap` is non-persisting in dry-run mode. `QAShard` closes the area-reasoning gap for focused TUI and browser views.

Canonical fields include schema version, project, sprint, phase, freshness, attempt and run correlation, input, implementation, map and policy fingerprints, coverage, limits, shards, theory outcomes, synthesis, progress, blocker, cancellation, terminal result, and next action where applicable.

QA phase values are:

```text
missing mapped queued running synthesizing completed blocked cancelled interrupted stale invalid
```

Freshness is separately `current` or `stale`. Phase `stale` means no current attempt may run or resume because authoritative fingerprints changed. Historical attempt terminal status remains inspectable. Run terminal results retain Sprint 35 values, including `succeeded`, `failed`, `cancelled`, `timed_out`, `interrupted`, `cleanup_uncertain`, and `persistence_degraded`. Execution failure is represented by the run terminal result plus a QA blocker; it is not a theory outcome. Cancellation responses mean the request was accepted and include current run cancellation plus QA snapshot. They do not claim terminal cancellation early.

CLI commands are:

```text
ultraplan sprint <project> <sprint> qa --dry-run [--json]
ultraplan sprint <project> <sprint> qa [--shard <shard-id>] [--json]
ultraplan sprint <project> <sprint> qa resume [--shard <shard-id>] [--json]
ultraplan sprint <project> <sprint> qa status [--json]
ultraplan sprint <project> <sprint> qa cancel --run <run-id> [--json]
ultraplan sprint <project> <sprint> qa recover [--json]
```

There is no public full restart or smoke-suite option in Sprint 36. JSON retains schema version 1 and operation `sprint.qa`. Valid dry-run, status, and completed diagnostic QA return success regardless of theory outcomes. Usage errors use `ExitUsage`; stale, invalid, denied, or budget-blocked work uses `ExitValidation`; runtime or persistence infrastructure uses `ExitRuntime`; cancellation or interrupted partial work uses `ExitPartial`. JSON still returns the bounded result on nonzero exits.

Closed operation kinds are `qa-status`, `qa-dry-run`, `qa-start`, `qa-resume`, and `qa-recover`. Runtime-backed start and resume use durable acceptance. Read-only dry-run and status do not. Recovery is runtime-free but mutating.

Versioned reads are:

```text
GET /api/v1/projects/{project}/sprints/{sprint}/qa
GET /api/v1/projects/{project}/sprints/{sprint}/qa/map
GET /api/v1/projects/{project}/sprints/{sprint}/qa/shards/{shard-id}
GET /api/v1/projects/{project}/sprints/{sprint}/qa/theories/{theory-id}
GET /api/v1/projects/{project}/sprints/{sprint}/qa/synthesis
```

Starts and cancellation reuse existing operation and run endpoints. Public errors use stable `qa.*` codes, bounded safe messages, retryability, bounded details, next action, and the existing outer request or run correlation identity. They exclude raw Go errors, absolute paths, provider payloads, command output, environment, secrets, and stacks. State locations are workspace-relative artifact references, app resource links, or opaque references, never absolute paths.
- **Rationale:** One typed product projection prevents command, TUI, and browser code from deriving different meanings or reading persistence directly.
- **Study / Source Grounding:** Gh-cli's factory and lazy dependencies support shared app wiring; rclone and restic support typed errors; restic progress supports one model with multiple renderers. Route and DTO details are UltraPlan decisions from API and frontend reasoning.
- **Trade-Offs Accepted:** Focused resources add handlers and fixtures. Closed operation kinds and additive JSON fields are easier to review and secure than generic actions.
- **Technical Debt / Future Impact:** Strict request DTOs require additive evolution. A future restart or suite option needs a new governed decision and compatibility fixtures.
- **Alternatives Rejected:** Exposing sprint structs; web imports of sprint; parsing CLI JSON; free-form operation actions; raw runtime logs; one unbounded nested status response; confirmed theories as command failures.
- **Contracts Satisfied:** Architecture, Errors, CLI Surface, Observability, Security, Documentation, `REQ-PARITY`, `REQ-WEB`.
- **Evidence Required:** App independence tests, command help and exits, JSON fixtures, strict decoder and fuzz tests, shard resource tests, operation preparation and stale confirmation, stable errors, redaction, and one cross-surface parity fixture.

### Decision 9: Keep TUI and browser presentation bounded and authority-free

- **Decision:** QA remains under sprint navigation. Browser HTML routes are the QA overview, focused shard, and focused theory routes. TUI routes are `RouteSprintQA`, `RouteSprintQAShard`, and `RouteSprintQATheory`. The TUI keeps one column and current viewport behavior. `Enter` opens detail, `Escape` returns, arrows or `j` and `k` navigate, `r` refreshes, `c` requests cancellation for the visible active run, and `q` leaves without cancellation.

The overview renders QA phase, separate Conformance Review verdict and freshness, map and implementation fingerprints, coverage, blocker and next action, bounded shard progress, theory totals, synthesis and follow-up state, and attempt and durable-run links. `completed` renders as `Read-only QA completed`, never `QA passed`. Every theory outcome and state has text and does not rely on color. Unknown counts say `Unknown`.

Server-rendered HTML is complete without JavaScript. Browser dry-run uses prepare and confirmation to preserve strict request binding, then returns the non-persisted map result directly; it creates no durable run and refresh returns the unchanged canonical QA status. Post/Redirect/Get applies to state mutations, not the read-only dry-run result. Start, resume, and recovery use the existing guarded operation forms. JavaScript may submit the same forms, observe existing run APIs, coalesce progress, and reload authoritative state. It adds no actions or authority.

Reconnect fetches current run and QA snapshots, requests committed events after the last sequence, deduplicates by run ID and sequence, and resumes if retained. A replay gap remains an observation warning alongside a valid current QA snapshot. Terminal delivery triggers a state refresh before completion is shown. TUI delivery follows the same rule.

All theory, evidence, errors, paths, and runtime text are hostile input. HTML escapes it. TUI strips ANSI and non-printing control characters before measuring or rendering. Browser controls retain current loopback, Host, Origin, session, CSRF, body, timeout, and confirmation protections. QA pages work at 320 CSS pixels and TUI views at `80x24` and `40x12`. Full-screen TUI requires interactive input and output; CLI text and JSON are the non-TTY fallback.
- **Rationale:** Presentation state should answer what the user is viewing, not whether work is current, cancelled, or complete.
- **Study / Source Grounding:** Restic's single progress owner and gh-cli's non-TTY behavior support observer separation. The frontend area inspected current TUI routing, durable event recovery, server templates, operation forms, SSE reconnect, and bounded browser rows.
- **Trade-Offs Accepted:** Server refresh and focused routes add navigation and can feel slower. They avoid an authoritative client store, optimistic completion, and unbounded rendering.
- **Technical Debt / Future Impact:** TUI and browser labels must stay centralized. Current template organization may evolve, but Sprint 36 will not reorganize the frontend tree.
- **Alternatives Rejected:** Top-level QA tab; split-pane TUI; single-page app; browser state store; JavaScript-only controls; observer-triggered cancellation; color-only status; raw Markdown or terminal output; exhaustive initial rendering.
- **Contracts Satisfied:** Architecture, Observability, Security, Testing, Documentation, `REQ-PARITY`, `REQ-WEB`.
- **Evidence Required:** Keyboard and route tests, no-JS snapshots, mobile and narrow terminal fixtures, focus behavior, reduced motion, hostile content, ANSI stripping, non-TTY startup, reconnect, replay gap, session rotation, dropped delivery, bounded rows, and authoritative terminal refresh.

### Decision 10: Prove semantics offline and reserve dogfood for real boundaries

- **Decision:** Normal tests are deterministic and offline. Pure tests cover canonicalization, IDs, maps, budgets, theories, synthesis, and invalidation. Temporary workspaces cover containment, symlinks, strict JSON, modes, atomic writes, pointer order, identity, and preservation. Fake runtime and process seams cover permissions, malformed output, retries, queueing, cancellation, timeouts, output caps, process cleanup, and unsupported capability. One canonical fixture is projected through app, CLI text, CLI JSON, TUI, HTTP JSON, server HTML, and durable run detail. Race tests cover scheduling, cancellation, promotion, stale writers, and observer replacement.

Gated dogfood against `../ultraplan-go` must produce a deterministic map, read-only shard records, synthesis, durable progress, cancellation and recovery evidence, and a clean before/after target identity including Git. Missing runtime, Sprint 35 gate evidence, current execute evidence, current Conformance Review evidence, or environment prerequisites yields `blocked` and does not satisfy the exit gate.

Release requires:

```text
go test ./...
go test -race ./...
go vet ./...
go build ./cmd/ultraplan
git diff --check
```

Documentation updates cover CLI, architecture, user workflow, browser behavior, recovery, JSON and state schemas, and release checks. Architecture Review, Sprint Review, and Deep Smoke Sprint protocols remain mandatory after implementation.
- **Rationale:** Fakes make failure cases reproducible, while dogfood is needed for Agentwrap, process, filesystem, browser, and Git identity claims.
- **Study / Source Grounding:** Chezmoi script tests, lazygit command fakes, restic failure seams, and dive real fixtures support this split. Exact release commands and protocols come from Sprint 36 requirements and project governance.
- **Trade-Offs Accepted:** Gated dogfood is slower and environment-dependent. It cannot be replaced by fake success.
- **Technical Debt / Future Impact:** The planning checkout does not itself prove every Sprint 35 artifact gate. Implementation must check the dependency and report blockage rather than fabricate evidence.
- **Alternatives Rejected:** Network-dependent normal tests; real runtime in unit tests; dogfood pass on missing prerequisites; presentation-only snapshots; skipping race or identity evidence.
- **Contracts Satisfied:** Testing, Observability, Security, Documentation, `REQ-RELEASE`.
- **Evidence Required:** The complete test matrix, all release commands, protocol outputs, dogfood artifact references, run-control evidence, and clean target identity.

## Expected Evidence

| Evidence Type | Required Evidence | Source / Command / Review Check |
| --- | --- | --- |
| Domain tests | Verification phases, v1 schemas, deterministic IDs, theory validation, all outcomes, limits, invalidation, and deterministic synthesis. | `internal/sprint/qa_state_test.go`, `qa_map_test.go`, `qa_test.go` |
| Persistence tests | Strict loading, unknown version, modes, containment, symlink escape, atomic failure, pointer order, digests, cancellation, recovery, retention, quota, and stale fencing. | Temporary workspace and injected atomic-write hooks in `internal/sprint` tests |
| Runtime safety | Default deny, read-only sandbox, capability failure, no shell, explicit argv, bounded environment, command denial, target identity, timeout, output cap, and process cleanup. | Fake runtime/process tests plus gated real-runtime dogfood |
| App and CLI | Typed query and command results, durable acceptance inventory, exits, JSON envelope, review alias compatibility, focused shard, cancellation, resume, and recovery. | `internal/app/sprint_usecases_test.go`, `sprint_commands_test.go`, run-control inventory tests |
| TUI | Routes, keys, all state labels, bounded rendering, narrow terminals, hostile text, dropped delivery, cancellation, and non-TTY behavior. | `internal/tui/qa_view_test.go` and shared parity fixture |
| Browser | Versioned resources, guarded operations, no-JS snapshot, strict requests, hostile content, focus, mobile layout, reconnect, replay gaps, restart, session rotation, and parity. | `internal/web/qa_handlers_test.go`, template tests, operations contract tests |
| Durable operation | Acceptance before child work, fencing, event order, bounded progress, explicit cancellation, terminal races, replay, recovery, and persistence degradation. | Existing Sprint 35 run-control tests extended with QA fixtures |
| Review | Ownership, compatibility, state authority, non-promotion, adapter boundaries, and no unnecessary framework. | Architecture Review and Sprint Review protocols |
| Runtime dogfood | Stable map, current shard and synthesis records, cancellation/recovery trace, durable run links, and clean before/after identity. | Gated Deep Smoke Sprint evidence; unavailable prerequisites report `blocked` |
| Release | Tests, race, vet, build, and whitespace checks all pass. | `go test ./...`; `go test -race ./...`; `go vet ./...`; `go build ./cmd/ultraplan`; `git diff --check` |
| Documentation | Commands, terminology, state and JSON schemas, authority, workflow, browser recovery, failure recovery, and release gates agree with fixtures. | `docs/cli-reference.md`, `docs/architecture.md`, `docs/user-guide.md`, `docs/local-web.md`, `docs/recovery.md`, `docs/phase3-json-schemas.md`, `docs/release-checklist.md` |

## Assumptions And Risks

| Item | Type | Impact | Mitigation / Follow-Up |
| --- | --- | --- | --- |
| Sprint 35 source behavior is present, but this planning checkout does not prove every prior artifact and dogfood gate. | Assumption | Runtime QA could start without its required operational foundation. | Preflight checks the current run-control capabilities and required gate evidence. Missing proof blocks dogfood and current QA. |
| Current failed or blocked Conformance Review is valid diagnostic input; missing or stale review is not. | Assumption | Readiness could block useful QA or accept stale findings. | Test freshness independently from verdict acceptability and never change the verdict. |
| Writer fencing can be supplied after durable claim and before child work. | Risk | Stale workers could overwrite current QA state. | If the token cannot be supplied and checked, block runtime QA. Add race tests. |
| Atomic rename covers one file, not an attempt transaction. | Risk | State and summary can disagree after failure. | Publish pointer-last, store digests, preserve prior valid records, and require explicit recovery. |
| Existing checks may write despite benign names. | Risk | QA could mutate protected content. | Closed catalog, adversarial corpus, sandboxing, and before/after identity. Unknown checks remain blocked. |
| Identity hashing may be expensive and still has a dispatch race. | Risk | Slow runs or undetected drift between checks. | Stream sorted manifests, recheck before dispatch, compare after, and block on drift. Do not add an unproven cache. |
| Cancellation can race with final shard promotion. | Risk | A valid completion could be lost or a cancelled shard marked complete. | Freeze promotion ordering around validation, atomic rename, writer token, and cancellation state; race-test both orders. |
| Run terminal persistence can succeed while QA summary persistence fails. | Risk | Adapters could claim semantic completion from run state. | Show both authorities, mark QA invalid or stale, and offer explicit recovery. |
| Fixed limits and 128 MiB storage may be too small. | Risk | Legitimate work blocks. | Report exact exhaustion, collect dogfood measurements, and revise only through a fingerprinted governed change. |
| Retention pruning could erase useful history. | Risk | Negative evidence becomes unavailable. | Prune only non-current attempts, preserve current and newest last-complete attempts, and record what recovery removed. |
| Machine and human status labels may drift. | Risk | Cross-surface parity fails or `completed` appears to mean pass. | Central app vocabulary and labels plus one parity fixture. |
| A challenger runtime can vary across fresh executions. | Risk | Synthesis inputs differ even when shard outcomes do not. | Persist and fingerprint validated challenge records; make final synthesis pure for identical normalized inputs; never include provider metadata in identity. |
| Project document metadata and long-range examples include later Sprint 37 and 38 behavior. | Risk | Plan could add smoke suites, restart, `qa.md`, or repair early. | Sprint 36 requirements and this reasoning control current scope. Documentation changes must distinguish current and future behavior. |

## Implementation Constraints

- Write QA semantics, mapping, detailed state, policy, orchestration, and synthesis only in `internal/sprint`.
- Do not add QA to `PlanningStage`, `PlanningStages()`, planning artifact order, or the planning-stage runtime map.
- Preserve `review`, `review.md`, `sprint.review`, verdict rules, smoke behavior, and existing JSON fields. `conformance-review` is one alias to the same implementation.
- Keep detailed maps, shard attempts, theories, command summaries, and synthesis out of `flow-state.json` and run-control events.
- Keep run-control and platform packages free of sprint QA types and outcome decisions.
- Keep `internal/web` imports limited to `internal/app` and the standard library.
- Add `QAShard` to the app and HTTP query contract so focused adapters never read files directly.
- Construct runtime dependencies lazily. Dry-run, status, shard, theory, and synthesis reads must not open runtime, accept runs, or write state.
- Treat recovery as a visible runtime-free verification-state mutation with no child work or durable run.
- Use schema v1, `qa-v1` scoped IDs, sorted canonical inputs, normalized JSON, pointer-last publication, private modes, contained paths, and explicit recovery.
- Require durable acceptance, claim, and writer token before runtime child work. Never adopt a lost worker or session.
- Deny investigator shell access, arbitrary argv, caller environment, Git, path escape, generated content, issue promotion, and repair.
- Apply before/after identity to source, tests, governed inputs, verification code, harness content, and Git when applicable.
- Validate every positive limit, source, maximum, and fingerprint before scheduling. Do not silently clamp or widen.
- Preserve every theory outcome. Use `cross_shard` and `not_applicable` as machine spellings.
- Keep QA phase, freshness, Conformance Review verdict, run lifecycle, cancellation, terminal result, and observer state as separate fields.
- Redact and bound values before QA persistence, durable events, logs, CLI output, TUI history, or browser delivery.
- Bound collections before constructing app, TUI, or template models. Paging must expose totals and stable order.
- Use the current template and static asset layout. Do not perform an unrelated frontend directory migration.
- Do not add `qa.md`, issues, repair endpoints, generated checks, smoke suites, product SQLite, plugins, daemons, remote workers, retrieval, content identity, or Git mutation.

## Plan Handoff

`plan.md` must execute these decisions without reopening architecture. It must sequence the small verification-phase and runtime-configuration refactor before QA domain work, then state and mapping, policy and orchestration, synthesis, app and CLI, TUI and browser, documentation, and release evidence.

The plan must carry forward:

- all ten final decisions and the local requirement trace labels
- the exact state layout, publication order, schema and identifier scope
- the complete product, scheduler, storage, delivery, and presentation bounds
- the approved-check and target-identity policy
- the authority split between QA state, flow summary, run control, and checkout
- the exact CLI, operation, app, HTTP, TUI, and browser contracts
- the risks and explicit block conditions
- Architecture Review, Sprint Review, and Deep Smoke Sprint protocols

## Phase Exit Criteria

- [x] Selected context was read and used.
- [x] Architecture, API Design, and Frontend reasoning documents were completed and synthesized.
- [x] Area conclusions are reflected and their open seams are resolved.
- [x] Contracts and named requirements map to decisions and evidence.
- [x] Final decisions are specific enough for `plan.md` to execute without architecture invention.
- [x] Expected evidence is concrete, bounded, and reviewable.
