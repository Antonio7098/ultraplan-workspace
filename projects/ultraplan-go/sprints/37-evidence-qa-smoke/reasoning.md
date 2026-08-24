# Sprint Reasoning: Evidence-producing QA and smoke integration

> Project: `ultraplan-go`
> Sprint: `37-evidence-qa-smoke`
> Output: `projects/ultraplan-go/sprints/37-evidence-qa-smoke/reasoning.md`
> **Inputs Used:** `projects/ultraplan-go/project-index.md`, `projects/ultraplan-go/roadmap.md`, `projects/ultraplan-go/docs/ARCHITECTURE.md`, `projects/ultraplan-go/docs/PRD.md`, `projects/ultraplan-go/docs/TRD.md`, `projects/ultraplan-go/sprints/37-evidence-qa-smoke/requirements.md`, `projects/ultraplan-go/sprints/37-evidence-qa-smoke/sprint-index.md`, `projects/ultraplan-go/sprints/37-evidence-qa-smoke/technical-handbook.md`, `projects/ultraplan-go/sprints/37-evidence-qa-smoke/reasoning/api-design.md`, `projects/ultraplan-go/sprints/37-evidence-qa-smoke/reasoning/architecture.md`, `projects/ultraplan-go/sprints/37-evidence-qa-smoke/reasoning/frontend.md`, `studies/go-cli-study/reports/final/01-project-structure.md`, `studies/go-cli-study/reports/final/02-command-architecture.md`, `studies/go-cli-study/reports/final/03-dependency-injection.md`, `studies/go-cli-study/reports/final/04-configuration-management.md`, `studies/go-cli-study/reports/final/05-error-handling.md`, `studies/go-cli-study/reports/final/06-io-abstraction.md`, `studies/go-cli-study/reports/final/07-state-context.md`, `studies/go-cli-study/reports/final/08-concurrency.md`, `studies/go-cli-study/reports/final/09-terminal-ux.md`, `studies/go-cli-study/reports/final/10-logging-observability.md`, `studies/go-cli-study/reports/final/11-testing-strategy.md`, `studies/go-cli-study/reports/final/12-extensibility.md`, `studies/go-cli-study/reports/final/13-security.md`, `studies/go-cli-study/reports/final/14-performance.md`, `system/protocols/architecture-review-protocol.md`, `system/protocols/review-sprint-protocol.md`, `system/protocols/deep-smoke-sprint-protocol.md`

This document decides the Sprint 37 architecture. It combines the selected context, study evidence, and completed area reasoning. It does not implement the feature.

## Sprint Purpose

- **Goal:** Add isolated evidence-producing QA, deterministic product-owned adjudication, bounded issue promotion, canonical `qa.md`, and `qa --suite smoke` while preserving target immutability and every existing smoke guarantee.
- **Non-Goals:** Production repair, generated-patch application, permanent test promotion, QA as a `PlanningStage`, replacement of Conformance Review or smoke compatibility, a general issue tracker, a generic sandbox or worker system, content identity, retrieval, alternate authored-artifact persistence, hosted operation, remote workers, or Git mutation.
- **Depends On:** A current acceptable Sprint 36 Conformance Review and required smoke evidence; Sprint 36 deterministic map, changed-path coverage, read-only investigation, cancellation, resume, invalidation, synthesis, and non-mutation proof; Sprint 35 durable run control; grounded planning; the current smoke implementation and external harness manifest; Agentwrap; and the existing platform process boundary.

## Selected Context And Pre-Reasoning Artifacts

| Artifact | Path | How It Was Used |
| --- | --- | --- |
| Sprint requirements | `requirements.md` | Supplied the binding outputs, acceptance criteria, non-goals, ownership constraints, entry gate, safety rules, and release commands. |
| Sprint index | `sprint-index.md` | Limited the work to the selected contracts, 14 evidence reports, three reasoning areas, and three review protocols. It also excluded repair, Git mutation, alternate persistence, and a second smoke authority. |
| Technical handbook | `technical-handbook.md` | Supplied the comparative patterns for thin adapters, explicit dependencies, bounded concurrency, typed failures, path safety, durable observation, fault injection, and compatibility testing. |
| Project architecture | `../../docs/ARCHITECTURE.md` | Fixed package ownership, dependency direction, state authority, the verification/repair split, and the rule that browser and TUI code consume app use cases. |
| Product requirements | `../../docs/PRD.md` | Fixed the user-visible QA capability, local-interface parity, Phase 5 order, and post-Sprint-39 deferrals. |
| Technical requirements | `../../docs/TRD.md` | Fixed Agentwrap and process boundaries, run-control reuse, state and migration requirements, smoke protocol authority, local web security, and release checks. |
| Product roadmap | `../../roadmap.md` | Fixed Sprint 37 promotion gates and the sequence from read-only QA to evidence QA, repair, and dogfood. |
| Architecture area reasoning | `reasoning/architecture.md` | Decided full-copy isolation, Linux fail-closed containment, product-mediated generated files, scoped identities, adjudication, immutable attempts, publication recovery, smoke delegation, and finite limits. |
| API design area reasoning | `reasoning/api-design.md` | Decided one typed QA operation, separate bounded summary/detail reads, durable operation reuse, stable machine errors, immutable cursors, CLI shape, and direct smoke aliasing. |
| Frontend area reasoning | `reasoning/frontend.md` | Decided server-rendered browser content, a feature-owned TUI view, current-versus-canonical labeling, hostile-text handling, bounded pagination, accessibility, and observation-only live updates. |
| Architecture Review protocol | `system/protocols/architecture-review-protocol.md` | Defines the post-implementation package, coupling, state, side-effect, error, testing, and maintainability audit. |
| Sprint Review protocol | `system/protocols/review-sprint-protocol.md` | Requires a current product-owned review after implementation and before smoke. |
| Deep Smoke Sprint protocol | `system/protocols/deep-smoke-sprint-protocol.md` | Defines the gated real-system evidence needed for containment, smoke parity, cancellation, and external evidence integrity. |

## Area-Specific Reasoning Inputs

| Area | Reasoning Document | Key Conclusion | Evidence Basis | Impact On Final Decision |
| --- | --- | --- | --- | --- |
| API Design | `reasoning/api-design.md` | Use one typed `SprintQARequest`, bounded summary and immutable detail projections, the existing durable operation routes, and the existing smoke operation for `qa --suite smoke`. | Thin command delegation, typed errors, explicit configuration, durable context propagation, bounded payloads, and smoke adapter evidence. | Fixes request fields, response authority, status versus assessment, cancellation, pagination, JSON compatibility, and alias behavior. |
| Architecture | `reasoning/architecture.md` | Keep QA policy in `internal/sprint`, strengthen generic process isolation, use a full local copy with Linux namespace containment, serialize writable attempts, publish immutable verification generations, and refactor smoke only enough to share one execution path. | Project module rules, live boundary inspection recorded by the area document, Restic-style interfaces and cleanup, security findings on explicit argv and private temp paths, and Helm-style delegation. | Fixes ownership, containment, identity, publication, adjudication, assessment, limits, and smoke parity. |
| Frontend | `reasoning/frontend.md` | Add QA to the existing sprint page and a dedicated TUI view over bounded app DTOs. Keep the no-JavaScript path complete and treat SSE/TUI events as observation only. | Terminal fallback, output separation, app boundary, hostile-input safety, bounded rendering, and behavior-focused test evidence. | Fixes current-versus-canonical presentation, action flow, accessibility, pagination, hostile content, reconnect, and recovery. |

All three area conclusions are adopted. None is overridden. Where they overlap, the architecture decision owns product behavior, API design owns adapter contracts, and frontend reasoning owns presentation.

## Sprint Technical Handbook Summary

- **Relevant Patterns:** Thin adapters over one product operation; product policy separated from generic mechanics; manual dependency construction; merged-policy validation before side effects; root-context cancellation with a separate cleanup deadline; bounded fan-out; typed failure facts; injectable filesystem/process seams; golden compatibility tests; delegation to an existing protocol authority; explicit argv and canonical path checks; bounded output and storage; and progress as observation rather than truth.
- **Important Trade-Offs:** Full copies support dirty and non-Git targets but cost disk and time. Sequential writable work reduces races but increases duration. Rich immutable state improves audit and recovery but raises migration and publication cost. Strict evidence caps protect resources but can make a result inconclusive. Direct smoke delegation preserves compatibility but leaves smoke-specific fields in a QA projection.
- **Warnings / Anti-Patterns:** Do not put adjudication in `internal/platform/process`; treat neither command failure nor model prose as an issue; do not claim cleanup when descendants or removal are uncertain; do not create detached background contexts for work; do not use shell strings; do not trust lexical path checks across symlinks; do not read verification files from adapters; do not create a smoke registry or second operation store; and do not optimize before hard aggregate limits exist.
- **Evidence Confidence:** High for package boundaries, explicit dependency construction, typed errors, context propagation, finite concurrency, process seams, security basics, bounded resource use, and testing strategy. Medium for protocol adapter design and terminal presentation because the reports offer analogous systems rather than this exact QA workflow. Sprint-specific identity, assessment, publication, and containment choices therefore rely on requirements and area reasoning as well as study evidence.

## Contracts Applied

Requirement IDs below label the acceptance criteria in their listed order. `AC-01` is the Sprint 36 entry gate and `AC-27` is the release command gate.

| Contract / Requirement ID | Constraint | Decision Impact | Expected Evidence |
| --- | --- | --- | --- |
| Architecture; AC-01, AC-10, AC-14, AC-18, AC-19 | Product semantics stay in `internal/sprint`; generic mechanics stay in platform code; adapters use `internal/app`; smoke keeps its authority. | Use focused files in existing packages, no QA module, registry, scheduler, or second smoke path. | Import-boundary tests, package review, and smoke call-path parity. |
| Security; AC-02, AC-03, AC-04, AC-05, AC-06, AC-07 | Writes require proven identity, containment, path safety, explicit argv, bounded environment, process cleanup, and target immutability. | Full private copy, product-mediated additions, Linux namespace backend, read-only execution mount, default deny, before/after target manifests, and fail-closed unsupported platforms. | Symlink/path/race tests, denied-write fixtures, environment/argv assertions, descendant cleanup tests, and real target identity audit. |
| Configuration; AC-05, AC-09 | Every limit is finite, validated after precedence resolution, and cannot be widened by an adapter or model. | Inject one validated `QASettings`; writable concurrency is fixed at one for Sprint 37. | Config precedence tests, invalid-limit tests, aggregate budget tests, and documented defaults. |
| Errors; AC-02, AC-06, AC-07, AC-12, AC-15, AC-16, AC-23 | Blocked admission, stale state, invalid evidence, cancellation, cleanup uncertainty, and publication failure need distinct machine outcomes and recovery guidance. | Add typed errors and stable app/JSON codes; preserve wrapped causes; keep execution status separate from assessment. | `errors.Is`/`errors.As` tests, exit/HTTP mapping fixtures, and recovery text checks. |
| Observability; AC-06, AC-07, AC-21, AC-22, AC-23 | Run, attempt, shard, command, evidence, cancellation, cleanup, and assessment facts must correlate without leaking raw evidence or secrets. | Reuse Sprint 35 durable run IDs, fencing, replay, and terminal arbitration; persist only bounded safe QA events. | Event correlation tests, replay-gap fixtures, redaction tests, and cross-surface run inspection. |
| Persistence And Migrations; AC-07, AC-14, AC-15, AC-16 | Detailed state is versioned, contained, digest-linked, immutable by attempt, atomically published, fenced, and recoverable. | Add verification schema version 2 over Sprint 36 version 1, immutable attempt directories, a state/report publication journal, retained prior generation, and a bounded flow projection. | Migration fixtures, unknown-version rejection, injected failures at every publication boundary, stale-writer races, digest/path checks, and prior-generation recovery. |
| LLM Runtime; LLM Evaluation / Cost / Safety; AC-04, AC-05, AC-10, AC-11 | Investigator and adjudicator model output is untrusted, bounded, structured, and unable to write or promote directly. | Read-only Agentwrap requests return generated-file and evidence-plan proposals; product code validates and materializes them; deterministic code makes final decisions. | Permission-policy tests, malformed/model-claim fixtures, bounded output tests, and direct-promotion denial. |
| Workflows; AC-08, AC-09, AC-13, AC-21, AC-22 | Attempts, retries, follow-up, cancellation, resume, publication, and terminal arbitration are finite and durable. | Serialize writable attempts, stop scheduling on cancellation, preserve completed evidence, clean under a separate deadline, and publish only complete generations. | Cancellation/resume tests, retry bounds, cleanup uncertainty, stale owner tests, and terminal race tests. |
| Performance; AC-06, AC-09, AC-14, AC-20 | Copy, output, evidence, storage, duration, rendering, and concurrency must remain bounded. | Enforce per-item and aggregate limits, paged details, capped fields, no raw smoke import, and no full event accumulation in TUI/browser. | Limit-table tests, truncation semantics, oversized target/evidence fixtures, and bounded rendering checks. |
| CLI Surface; AC-17, AC-18, AC-20 | `qa`, focused shard, resume/restart, status, cancel, and `qa --suite smoke` need stable text/JSON behavior without changing existing smoke clients. | Add thin commands over app use cases and preserve the existing smoke operation kind and results for both spellings. | Help and flag goldens, JSON schema fixtures, exit-code tests, alias parity, and app/CLI agreement. |
| Testing; AC-24, AC-25, AC-27 | Normal tests are offline and deterministic; race, fault, security, parity, and gated dogfood evidence are required. | Split tests by ownership and use temporary targets, fake runtimes/processes, fake harnesses, race tests, and one gated real run. | Required Go commands, package tests, review protocols, and external dogfood records. |
| Documentation; AC-17, AC-18, AC-20, AC-23, AC-26 | Operators need accurate command, architecture, workflow, browser, recovery, schema, compatibility, and release guidance. | Update the seven required documentation files and freeze documented JSON/state fields. | Documentation review, executable examples, schema fixtures, and release checklist audit. |

## Repos Studied / Source Evidence Used

| Source / Repo / Report | Concrete Reference | Relevant Finding | Why It Matters For This Sprint | Used In Decision(s) |
| --- | --- | --- | --- | --- |
| Chezmoi, project structure | `01-project-structure.md`, `chezmoi/main.go:26-34`, `internal/chezmoi/chezmoi.go:1-2` | Thin entrypoints and one-way protected package dependencies. | Supports keeping CLI/TUI/web out of QA policy and state logic. | 1, 7, 8 |
| Helm, command architecture | `02-command-architecture.md`, `helm/pkg/cmd/install.go:132-145` | Command setup delegates to an action rather than implementing the operation. | Supports thin QA commands and a compatibility spelling that reaches one smoke implementation. | 6, 7 |
| Gh CLI, dependency injection and IO | `03-dependency-injection.md`, `gh-cli/pkg/cmdutil/factory.go:16-43`; `06-io-abstraction.md`, `gh-cli/pkg/iostreams/iostreams.go:551-568` | Explicit factories and test IO make dependencies and output behavior replaceable. | Supports injected process/isolation/runtime seams and adapter fixtures. | 1, 2, 7, 9 |
| Restic, IO and state/context | `06-io-abstraction.md`, `restic/internal/fs/interface.go:10-31`, `internal/backend/backend.go:19-90`; `07-state-context.md`, `internal/restic/lock.go:290-305` | Focused external-boundary interfaces and separate bounded cleanup contexts. | Supports fault-injected copy/process tests and cleanup after work cancellation. | 2, 5, 9 |
| Go Task, context and errors | `07-state-context.md`, `go-task/task.go:89`; `05-error-handling.md`, `go-task/errors/errors.go:47-50` | Structured cancellation and machine-readable task failures. | Supports distinct attempt status, assessment, cancellation, and stable adapter mappings. | 4, 7 |
| K9s and Dive, configuration | `04-configuration-management.md`, `k9s/internal/config/json/validator.go:1-187`, `dive/cmd/dive/cli/internal/options/analysis.go:48-53` | Validate merged configuration and cross-field constraints before execution. | Supports one validated `QASettings` and preflight before writable side effects. | 1, 2, 9 |
| Opencode, permissions and shutdown | `13-security.md`, `opencode/internal/permission/permission.go:44-108`; `08-concurrency.md`, `opencode/cmd/root.go:261-279` | Permission decisions need typed enforcement, and shutdown waits need a deadline. | Supports default deny, explicit capabilities, and truthful cleanup uncertainty. | 2, 9 |
| Lazygit, safe commands and fakes | `13-security.md`, `lazygit/cmd_obj_builder.go:38`; `11-testing-strategy.md`, `lazygit/pkg/commands/oscommands/fake_cmd_obj_runner.go:17-26` | Explicit argv blocks shell interpretation and fake runners can freeze cwd/env/calls. | Supports evidence-plan execution and deterministic process tests. | 2, 9 |
| Helm, extensibility | `12-extensibility.md`, `helm/internal/plugin/metadata_v1.go:24-48`, `internal/plugin/runtime_subprocess.go:65-79` | Versioned adapters can preserve an existing subprocess protocol. | Supports wrapping smoke without parsing or replacing its manifest path. | 6 |
| Chezmoi, Helm, and Rclone testing | `11-testing-strategy.md`, `chezmoi/internal/cmd/main_test.go:64-174`, `helm/internal/test/test.go:43`, `rclone/cmd/bisync/bisync_test.go:1435-1479` | Scenario tests and reviewed goldens protect command and artifact compatibility. | Supports CLI, JSON, HTML, TUI, `qa.md`, and smoke parity fixtures. | 6, 8, 9 |
| Age and Restic, bounded processing | `14-performance.md`, `age/internal/stream/stream.go:20,195-219`, `restic/internal/archiver/buffer.go:24-46` | Chunked work and rejection of oversized retained buffers prevent unbounded memory. | Supports copy/output/patch/storage caps and explicit truncation. | 2, 3, 9 |
| Terminal UX and observability reports | `09-terminal-ux.md`, `restic/internal/ui/termstatus/status.go:197-205`; `10-logging-observability.md`, `helm/internal/logging/logging.go:31-71` | Progress is interruptible presentation, while diagnostics stay separate from data output. | Supports observation-only TUI/SSE and clean JSON stdout. | 7, 8 |

The remaining selected report findings are also applied: project layout from `01`, command routing from `02`, dependency construction from `03`, configuration precedence from `04`, typed errors from `05`, IO seams from `06`, context discipline from `07`, bounded concurrency from `08`, terminal fallback from `09`, structured logging from `10`, layered testing from `11`, versioned delegation from `12`, trust boundaries from `13`, and finite resource use from `14`.

## Trade-Off And Debt Analysis

### Accepted Trade-Offs

| Trade-Off | Benefit | Cost / Constraint Accepted | Why Acceptable Now | Revisit Trigger |
| --- | --- | --- | --- | --- |
| Full local copy instead of Git worktree | Represents dirty, uncommitted, and non-Git targets exactly enough for verification. | Copying large targets costs time and disk. | Git state cannot represent all required target states and QA may not mutate Git. | Dogfood shows copy limits block valid targets or measured cost dominates run time. |
| Linux-only proven writable containment at first | Uses read-only mounts, PID/mount/network namespaces, dropped capabilities, and parent-death behavior that can be tested. | Writable QA is blocked on unsupported or unproven platforms. | A false claim of isolation is worse than reduced availability. | A macOS or other backend proves the same write, process, descendant, and cleanup guarantees. |
| Product-mediated generated additions | Prevents direct model writes and expectation weakening. | Investigators cannot use arbitrary editors or modify existing copied files. | Sprint 37 needs evidence, not flexible implementation work. | Sprint 38 repair defines a separately confirmed mutation contract. |
| Sequential writable attempts | Simplifies target drift, workspace ownership, process cleanup, budgets, and publication fencing. | A many-shard run can take up to the total two-hour budget. | Independent writable isolation has not yet earned parallel execution. | Fault and race dogfood proves independent workspaces and descendant cleanup with useful speedup. |
| Strict caps and explicit inconclusive results | Keeps model output, commands, evidence, storage, JSON, TUI, and HTML bounded. | A decisive line outside retained output cannot be used. | Honest inconclusive evidence is safer than inferred success. | Measurements show a specific cap is too low without breaching aggregate budgets. |
| Rich immutable verification state plus small projections | Preserves lineage, rejected evidence, issues, and recovery while keeping common reads cheap. | Schema validation, migration, and publication are substantial work. | Evidence promotion cannot be audited from a count-only state file. | Post-Sprint-39 content/schema work supplies a better proven representation. |
| Preserve the prior complete `qa.md` on failed rerun | Operators retain the last complete report. | Interfaces must explain two attempt identities. | Replacing good evidence with a partial report would destroy useful state. | Never for correctness; presentation may simplify after user testing. |
| Direct smoke delegation | Preserves current manifest, selection, invocation, evidence, verdict, and compatibility. | QA accepts smoke-specific result asymmetry. | A uniform wrapper is less important than one authority. | Only a separately governed compatibility migration after parity and deprecation evidence. |

### Potential Technical Debt

| Debt / Shortcut | Why It Might Accrue | Current Mitigation | Owner / Follow-Up |
| --- | --- | --- | --- |
| Linux `bwrap` dependency | Some local environments will lack the launcher or namespace support. | Capability preflight returns a precise blocked reason before workspace writes. | Sprint 39 measures blocked rates; a later focused platform backend may be proposed. |
| Full-tree manifests | Large targets make repeated hashing expensive. | Deterministic walk, entry/byte caps, streaming hashes, and measured dogfood timings. | Sprint 39 may justify an equivalent incremental manifest, but not weaker identity. |
| Verification schema version 2 size | Evidence, patches, adjudication, and issues may create migration burden. | Immutable attempts, explicit version dispatch, total storage caps, bounded state root, and golden fixtures. | Post-Sprint-39 content-contract work must preserve compatibility. |
| Product adjudication rule growth | New evidence kinds may create a large validator and grouping table. | Keep local record validation beside types and relational validation in the adjudicator; no generic rule engine. | Sprint 39 reviews false positives and complexity before adding evidence kinds. |
| Duplicate current/canonical UI concepts | Every adapter must render two timelines correctly. | Shared app DTOs, shared parity fixtures, and explicit field names. | UI hardening follows measured operator confusion, not a hidden merge of the concepts. |
| Smoke adapter metadata | Recording the invoking spelling may tempt behavior branches. | Treat source spelling as diagnostics only and parity-test every operational field. | Remove source metadata if it causes branching without diagnostic value. |

### Future Considerations

| Consideration | Deferred Until | Reason Deferred | What Should Be Preserved Now |
| --- | --- | --- | --- |
| Production repair and generated-patch application | Sprint 38 | Repair needs frozen promoted issues and separate confirmation. | Exact evidence links, repair eligibility, regression candidates, and unapplied patch bytes. |
| Parallel writable attempts | After Sprint 37/39 isolation evidence | Concurrency multiplies target drift, cleanup, storage, and stale-writer races. | Distinct workspace and attempt IDs; no API promise that attempts must remain sequential forever. |
| Additional containment backends | Proven platform-specific design | Best-effort parity would violate fail-closed admission. | Capability reporting and backend identity in evidence records. |
| Permanent regression test promotion | Governed repair or maintenance sprint | QA generated checks are evidence, not accepted product code. | Patch identity, expectation links, and regression-candidate classification. |
| Global content identity and provenance | After Sprint 39 | The schema must be shaped by real evidence and repair artifacts. | Scoped SHA-256 identities with explicit version and migration behavior. |
| Retrieval, alternate authored persistence, graph, cloud, remote workers | Post-Sprint-39 gates | None is needed to produce local evidence safely. | Clear authority boundaries and exportable contained references. |
| Smoke compatibility retirement | Separate compatibility decision | Sprint 37 must preserve `smoke` and `smoke.md`. | One execution path and exhaustive parity fixtures. |

## Decisions

### Decision 1: Keep one product-owned QA workflow with a runtime admission gate

- **Decision:** `internal/sprint` will own writable QA admission, evidence-plan policy, identities, adjudication, issue promotion, assessment, `qa.md`, verification state, and smoke projection. Every writable run will revalidate the current Sprint 36 evidence gate before workspace creation. `internal/platform/process` will own only generic copy, containment, process, cancellation, and cleanup mechanics. `internal/runcontrol` remains the only operational lifecycle authority. CLI, TUI, and web will reach QA through `internal/app`.
- **Rationale:** The proof of read-only QA safety is an input to each writable attempt, not a historical assumption. Keeping policy, mechanics, operations, and presentation in their existing owners prevents a second workflow engine and makes failure authority explicit.
- **Study / Source Grounding:** `technical-handbook.md` sections "Thin adapters over one product operation" and "Product policy and generic mechanics stay on opposite sides of an interface"; `01-project-structure.md` at `chezmoi/main.go:26-34` and Restic's repository interface split; `03-dependency-injection.md` on explicit composition roots. The Sprint 36 gate itself comes from `requirements.md` AC-01 and project documents, not from a studied repository.
- **Trade-Offs Accepted:** QA coordination adds focused files to the already large `internal/sprint` package. That is preferable to a premature QA package hierarchy or workflow framework.
- **Technical Debt / Future Impact:** If Sprint 38 makes the sprint package hard to navigate, split only along proven dependency boundaries. Do not move policy into platform code to reduce file count.
- **Alternatives Rejected:** A top-level QA module was rejected because sprint artifacts and verification semantics are sprint-owned. A generic workflow/sandbox system was rejected because there is one local use case. Treating chronology as admission proof was rejected because stale or missing Sprint 36 evidence must block.
- **Contracts Satisfied:** Architecture, Workflows, Observability, Security; AC-01, AC-02, AC-10, AC-14, AC-21.
- **Evidence Required:** Entry-gate table tests for every missing/stale prerequisite; import-boundary tests; no-child-start assertions; architecture review of package direction; durable run correlation tests.

### Decision 2: Use full-copy, product-mediated, fail-closed writable isolation

- **Decision:** Each writable shard attempt will receive a new private full copy of the current target. The copy algorithm will use descriptor-safe or equivalently race-resistant traversal, `Lstat`, deterministic relative paths, regular-file copies, no preserved hard links, validated internal symlinks only, private workspace permissions, and entry/byte caps. It will exclude Git administrative data while separately recording Git control-state identity. Product code will materialize only validated new test, fixture, or data-probe files. Before child execution, Linux containment will expose the copy read-only, expose only attempt-local cache/temp paths as writable, disable network, drop capabilities, own a PID namespace/process group, and prove descendant cleanup. Unsupported or uncertain containment returns `blocked` before any generated write.
- **Rationale:** Dirty and non-Git targets rule out commit-only identity. Prompt instructions and cwd checks cannot stop absolute-path writes or hostile test code. Product-mediated additions prevent investigators from weakening existing expectations, and an OS write boundary protects both target and copied production files during execution.
- **Study / Source Grounding:** `13-security.md` on explicit argv, private temporary directories, permission decisions, and the path-canonicalization warning; `06-io-abstraction.md` on testable filesystem/process seams; `07-state-context.md` on a separate cleanup context; `14-performance.md` on streaming and bounded work. The exact namespace and read-only mount design is sprint-specific because the studies do not prove writable agent containment.
- **Trade-Offs Accepted:** Initial writable support is limited to environments that prove the backend capabilities. Copies consume more disk and time than worktrees.
- **Technical Debt / Future Impact:** A platform backend interface may gain another implementation later, but Sprint 37 will not expose a best-effort mode. Full-copy performance must be measured before optimization.
- **Alternatives Rejected:** Git worktrees were rejected as the primary mechanism because they omit dirty and non-Git state. Direct agent writes were rejected because runtime permission metadata is weaker than product validation. Lexical path checks and process groups alone were rejected because symlinks, absolute paths, and detached descendants can escape them.
- **Contracts Satisfied:** Security, Architecture, Performance, Errors; AC-02 through AC-07, AC-09.
- **Evidence Required:** Copy and non-Git tests; dirty-target fixtures; identity and containment tests; symlink swap/path escape/device/socket/FIFO rejection; denied target and production-file writes; exact argv/cwd/env tests; cancellation and descendant cleanup; removal failure; target before/after manifest equality; gated real-repository containment evidence.

### Decision 3: Freeze bounded evidence plans and use verification-scoped identities

- **Decision:** Before generated bytes are written, product code will validate and freeze expectation/theory references, confirmation/refutation/inconclusive conditions, approved new paths, generated content digests, exact executable and argv, relative cwd, allowlisted environment, timeout, output caps, structured result requirements, repeatability class, and cleanup obligations. Commands will never come from Markdown or shell strings. SHA-256 identities will bind governed inputs, implementation manifest, Git control state when present, map, shard, theory, workspace, plan, command, patch, and validated external smoke evidence. IDs are deterministic within the versioned verification model and explicitly not global content IDs.
- **Rationale:** Evidence is useful only when it can be traced to the exact expectation, source state, workspace, command, and preserved patch. Freezing the plan prevents a generated check from changing its own success conditions after execution.
- **Study / Source Grounding:** `04-configuration-management.md` on validating complete merged options before work; `05-error-handling.md` on structured machine facts; `13-security.md` on argument arrays and redaction; `14-performance.md` on bounded buffers. Global identity design has no relevant study source because the selected reports do not cover content-addressed QA lineage; scoped identities follow the sprint requirements and architecture reasoning instead.
- **Trade-Offs Accepted:** Canonical serialization and digest validation add implementation and fixture work. IDs are intentionally local to verification state and may need explicit migration after Sprint 39.
- **Technical Debt / Future Impact:** The later content contract must map or preserve these references. It must not silently reinterpret them as workspace-global IDs.
- **Alternatives Rejected:** Commit SHA identity was rejected because targets can be dirty or non-Git. Path-only references were rejected because they do not bind content or attempt. A global provenance service was rejected as premature. Commands parsed from model prose were rejected as unsafe and non-reproducible.
- **Contracts Satisfied:** Security, Persistence And Migrations, LLM Evaluation / Cost / Safety, Performance; AC-03, AC-05, AC-06, AC-07, AC-08.
- **Evidence Required:** Canonical digest goldens; mismatch and stale identity tables; command plan validation; environment allowlist/redaction tests; patch digest tests; external evidence identity checks; target drift invalidation at every checkpoint.

### Decision 4: Make adjudication and assessment deterministic product decisions

- **Decision:** A global adjudicator in `internal/sprint` will validate expectation grounding, current fingerprints, setup, containment, cleanup, confirmation-condition fidelity, evidence sufficiency, repeatability, flakiness, severity, and external identity. Product code will recompute root-cause groups from normalized expectation, component, failure signature, shard set, and evidence links. Only this code can promote an issue, mark repair eligibility, classify a regression candidate, or compute assessment. One execution is sufficient only for a complete deterministic assertion without contradictory evidence. Variable evidence requires two matching executions in distinct workspaces. Contradiction produces flaky evidence. Truncation is allowed only when decisive structured evidence remains complete; otherwise the outcome is inconclusive.
- **Rationale:** A failed command proves an observation, not a root cause. Model output can help summarize or propose a grouping, but it cannot satisfy a product acceptance rule.
- **Study / Source Grounding:** `05-error-handling.md` distinguishes execution facts from behavioral classification; `11-testing-strategy.md` favors behavior assertions and fault fixtures; the handbook warning "Do not treat command failure or model prose as an issue" directly applies. The exact assessment precedence comes from sprint requirements and area reasoning, not a studied repository.
- **Trade-Offs Accepted:** Deterministic adjudication needs a sizeable fixture matrix and conservative `blocked` or `inconclusive` outcomes when proof is incomplete.
- **Technical Debt / Future Impact:** Rule growth must stay as explicit validators and normalization functions, not become a generic rules engine. Sprint 39 should measure false-positive, rejection, and inconclusive rates.
- **Alternatives Rejected:** Investigator self-promotion, command-exit promotion, harness-issue promotion, and model-owned assessment were rejected because all bypass evidence validation. Count-only assessment was rejected because stale, narrow, rejected, or cleanup-uncertain evidence could look complete.
- **Contracts Satisfied:** LLM Runtime, LLM Evaluation / Cost / Safety, Security, Testing; AC-10 through AC-14.
- **Evidence Required:** Promotion/rejection tables for stale, malformed, flaky, ungrounded, diagnostic, narrow, failed-setup, uncontained, truncated, and cleanup-uncertain evidence; repeatability fixtures; deterministic root-cause grouping; model-output distrust tests; assessment precedence goldens.

### Decision 5: Publish immutable attempts and recover the `state.json`/`qa.md` generation as one unit

- **Decision:** Verification state will use schema version 2, with an explicit additive migration from Sprint 36 version 1. Attempt directories become immutable after terminal publication. Publication will stage patches and evidence, validate them from disk, write adjudication/issues/attempt metadata and an attempt-local report, sync and rename the complete attempt, then commit a root `verification/state.json` and sprint-root `qa.md` generation under the sprint mutation lease and stale-writer token. A small journal and retained prior pair will recover from a crash between the two root files. `flow-state.json` will update last and contain only bounded status, freshness, counts, assessment, next action, and contained pointers/digests. Unknown versions or fields, invalid digests, escaping pointers, over-limit records, partial generations, and stale writers fail closed.
- **Rationale:** One filesystem rename cannot replace two independent root files. Readers must never combine a new state pointer with an old report or lose the prior valid report because a rerun failed.
- **Study / Source Grounding:** `04-configuration-management.md` covers explicit migration and schema validation; `05-error-handling.md` supports typed recovery failures; `11-testing-strategy.md` supports scenario and golden fixtures. The multi-file publication journal is a project-specific response to AC-14 through AC-16 because the studies do not provide a directly equivalent artifact pair.
- **Trade-Offs Accepted:** The journal adds a small recovery state machine and directory-sync work. This cost is justified by the requirement to retain the last complete report while exposing a current failed attempt.
- **Technical Debt / Future Impact:** Schema 2 is verification-scoped. A later content contract must migrate explicitly. Run-control fencing facts are referenced but never duplicated or renewed by verification state.
- **Alternatives Rejected:** Independent atomic renames were rejected because they can expose mismatched generations. Putting all detail in `flow-state.json` was rejected because it mixes authorities and breaks bounds. Rewriting Sprint 36 history was rejected because old theory outcomes must remain intact.
- **Contracts Satisfied:** Persistence And Migrations, Workflows, Architecture, Errors; AC-14, AC-15, AC-16, AC-21, AC-23.
- **Evidence Required:** Version 1 to 2 migration fixtures; unknown-version/field failures; path and digest validation; stale-writer races; injected failure and process-death tests at each sync/rename/journal boundary; prior-generation restoration; canonical/current attempt projection tests.

### Decision 6: Implement smoke-as-QA by direct delegation to the existing smoke authority

- **Decision:** `qa --suite smoke` will map to the existing typed smoke operation and durable operation kind. It will call the same `Service.RunSmoke` execution path as `smoke`, with only a private refactor to avoid acquiring the sprint mutation lease twice. Both spellings will share authoring, manifest discovery, scope and containing-suite selection, argv, cwd, environment, timeout, cancellation, descendant cleanup, evidence validation, verdict, flow projection, external run ID, and `smoke.md` bytes. QA will store only validated links, identities, bounded adjudication facts, and canonical-versus-diagnostic status. Raw harness evidence stays external.
- **Rationale:** Smoke already owns protocol and evidence authority. Any QA-native parser or runner would create drift in the exact safety behavior this sprint must preserve.
- **Study / Source Grounding:** `12-extensibility.md` at Helm's `metadata_v1.go:24-48` and `runtime_subprocess.go:65-79` shows versioned delegation to an existing process authority. `02-command-architecture.md` supports thin aliases. The handbook explicitly warns against a second smoke registry.
- **Trade-Offs Accepted:** Smoke does not look exactly like ordinary shard evidence. The QA projection will retain that asymmetry rather than normalize away important smoke facts.
- **Technical Debt / Future Impact:** Compatibility remains. Retirement requires a later explicit decision after parity, client migration, and deprecation evidence.
- **Alternatives Rejected:** A QA-native smoke executor, duplicate manifest parser, generic executor registry, and copying raw harness evidence into verification state were rejected because each creates a second authority or violates storage ownership.
- **Contracts Satisfied:** Architecture, CLI Surface, Workflows, Security; AC-17 through AC-20.
- **Evidence Required:** Fixture and gated parity comparing suite/test selection, containing suite, argv, cwd, environment, timeout, cancellation, cleanup, external run ID, evidence links, verdict, flow state, durable run result, issue evidence behavior, and exact `smoke.md` digest.

### Decision 7: Expose one bounded app operation and keep durable run control authoritative

- **Decision:** Add a typed `SprintQARequest` with project, sprint, optional shard, optional `smoke` suite, mutually exclusive resume/restart, dry-run, and expected fingerprint. Callers cannot supply argv, environment, paths, timeouts, limits, or adjudication instructions. Runtime-backed starts use Sprint 35 durable acceptance, owner claim, fencing, events, cancellation, and terminal arbitration before child work. Read use cases return a constant-size `SprintQASummary` and immutable paged `SprintQADetail`. CLI commands are `qa`, `qa status`, and `qa cancel --run`; browser starts and cancellation reuse existing operation/run routes and confirmation security. Stable error codes map typed failures without exposing raw internals.
- **Rationale:** One operation keeps all adapters consistent. Separate summary and detail models prevent status and reconnect from loading unbounded evidence. Durable lifecycle and product assessment remain different facts.
- **Study / Source Grounding:** `02-command-architecture.md` on thin delegates; `03-dependency-injection.md` on explicit construction; `05-error-handling.md` on typed errors and exit mapping; `09-terminal-ux.md` and `10-logging-observability.md` on presentation versus authority and stdout/stderr separation.
- **Trade-Offs Accepted:** Clients must follow a run link and then refresh QA state. Immutable pagination adds DTO mapping and cursor validation.
- **Technical Debt / Future Impact:** API schema version, verification schema version, and run-control schema version remain separate. Callers must tolerate additive fields but not unsupported major versions.
- **Alternatives Rejected:** QA-specific run/SSE registries were rejected because run control already owns lifecycle. One unbounded response was rejected because evidence volume is hostile input. Caller-defined execution mechanics were rejected as a permission bypass.
- **Contracts Satisfied:** CLI Surface, Architecture, Observability, Security, Errors, Performance; AC-20 through AC-23.
- **Evidence Required:** App use-case tests for acceptance-before-child, fencing, correlation, bounds, pagination, authorization-independent reads, fresh cancellation authority, and next action; CLI help/JSON/exit goldens; HTTP strict decode, confirmation, CSRF, origin, idempotency, reconnect, restart, and replay-gap tests.

### Decision 8: Present QA through bounded server rendering and a feature-owned TUI view

- **Decision:** Add a QA section to the existing sprint page and `internal/tui/qa_view.go`, both over app DTOs. Every view will show current attempt, canonical attempt, freshness, durable run status, assessment, cleanup certainty, evidence quality, adjudication, bounded issues, smoke canonical status, blockers, and product-derived next action. The browser will work fully without JavaScript. Existing dependency-free JavaScript will add progress, reconnect, pagination, disclosure, and cancellation only. Hostile evidence renders as escaped plain text; terminal control sequences are stripped. Browser pages use semantic headings, native controls, a restrained live region, visible focus, reduced-motion behavior, and a 320-pixel layout. Detail pages default to 25 records, reject more than 100, cap one text field at 8 KiB, and disclose all truncation.
- **Rationale:** The difficult UI problem is truthful state, especially a failed current rerun beside an older canonical report. A client framework would not solve that and would add a second state model.
- **Study / Source Grounding:** `09-terminal-ux.md` on non-TTY fallback and interruptible progress; `10-logging-observability.md` on safe output separation; `11-testing-strategy.md` on behavior-oriented UI fixtures; `13-security.md` on untrusted input and path controls; `14-performance.md` on bounded rendering.
- **Trade-Offs Accepted:** The sprint page becomes denser and users page through large evidence sets. Native controls are less custom but easier to test and use without JavaScript.
- **Technical Debt / Future Impact:** Shared fixtures must prevent adapter drift. UI polish may follow operator testing, but adapters may never derive assessment or freshness.
- **Alternatives Rejected:** A QA-only browser application, global client store, rich rendering of investigator Markdown, raw external evidence display, and TUI execution of CLI commands were rejected because they duplicate state, weaken escaping, or bypass app use cases.
- **Contracts Satisfied:** Architecture, Security, Performance, Observability, Testing; AC-20, AC-21, AC-22, AC-23, AC-24.
- **Evidence Required:** No-JavaScript browser snapshots; TUI keyboard/focus/resize/pagination tests; hostile HTML/ANSI/invalid UTF-8 fixtures; cancellation and focus restoration; replay-gap refresh; mobile/reduced-motion/accessibility checks; shared app/CLI/TUI/web/Markdown/state parity fixtures.

### Decision 9: Enforce fixed execution budgets and prove behavior with layered tests and dogfood

- **Decision:** Resolve one injected `QASettings` after normal precedence and validate it before admission. Sprint 37 defaults are: writable concurrency 1; 32 shards per run; 2 attempts per shard; 4 generated checks, 8 commands, and 32 evidence records per attempt; 8 adjudication follow-ups per run; 64 promoted issues per attempt; 256 KiB investigator output and generated file; 2 MiB generated patches per attempt; 1 MiB stdout and 512 KiB stderr per command; 128 argv elements and 16 KiB encoded argv; 32 environment names and 32 KiB encoded environment; 5 minutes per command; 20 minutes per attempt; 2 hours per run; 10 seconds cleanup grace; 1 command retry; 500,000 copy entries and 8 GiB copied bytes; 256 MiB retained per run; and 1 GiB retained verification state per sprint. Aggregate duration and storage caps override multiplied item limits. Existing smoke settings remain authoritative for smoke. Tests will follow package ownership, then run race, vet, build, diff, architecture review, sprint review, and gated deep smoke.
- **Rationale:** Per-item limits alone can multiply into an unsafe total. Sequential writes and aggregate caps make the first implementation auditable. Fault injection is mandatory because happy-path process and filesystem tests cannot prove cleanup or publication recovery.
- **Study / Source Grounding:** `04-configuration-management.md` on merged validation; `08-concurrency.md` on bounded launch sites and wait deadlines; `11-testing-strategy.md` on functional fakes, scenarios, and goldens; `14-performance.md` on finite buffers and concurrency. The numeric defaults are sprint-specific decisions from architecture reasoning, not values copied from the studies.
- **Trade-Offs Accepted:** Valid large investigations may stop incomplete and require a focused rerun. The test suite will be large because the safety case spans process, filesystem, state, runtime, adapters, and smoke.
- **Technical Debt / Future Impact:** Sprint 39 must measure copy time, blocked rate, truncation, retained bytes, evidence validity, and run duration before any default changes. Do not add pooling or parallelism without measurements.
- **Alternatives Rejected:** User-controlled widening, unlimited evidence retention, multiplying independent maxima without total caps, and default parallel writable work were rejected as unsafe. Real-runtime dependencies in normal tests were rejected as non-deterministic.
- **Contracts Satisfied:** Configuration, Performance, Testing, Workflows, Documentation; AC-09, AC-24 through AC-27.
- **Evidence Required:** Limit validation and aggregate budget tables; fault-injected creation/copy/command/cancellation/descendant/removal/publication tests; package race tests; adapter parity fixtures; gated real contained check with one adjudication rejection or promotion audit; unchanged target identity; workspace cleanup; smoke parity; `go test ./...`, `go test -race ./...`, `go vet ./...`, `go build ./cmd/ultraplan`, and `git diff --check` in `../ultraplan-go`.

## Expected Evidence

| Evidence Type | Required Evidence | Source / Command / Review Check |
| --- | --- | --- |
| Platform unit tests | Creation, copy, manifest, symlink, special-file, command, timeout, cancellation, descendant, and removal faults with no sprint imports. | `internal/platform/process/isolation_test.go`; `go test ./internal/platform/process` |
| Sprint unit tests | Entry gate, plan policy, target identity, evidence validation, repeatability, truncation, adjudication, grouping, assessment, migration, publication, recovery, and smoke projection. | `internal/sprint/qa*_test.go`; `go test ./internal/sprint` |
| Application tests | Durable acceptance, fencing, result bounds, summary/detail projections, run correlation, cancellation, recovery, and adapter-independent next action. | `internal/app/sprint_usecases_test.go`, `internal/app/sprint_commands_test.go` |
| Adapter tests | CLI help/JSON/exit behavior; TUI keyboard and hostile text; HTTP guards, strict decoding, reconnect, cancellation, restart; no-JavaScript HTML and bounded rendering. | `internal/app`, `internal/tui/qa_view_test.go`, `internal/web/qa_handlers_test.go` |
| Compatibility tests | Identical operational and artifact results for `smoke` and `qa --suite smoke`. | `internal/sprint/smoke_test.go` plus shared app/CLI fixtures |
| Persistence tests | Schema 1 to 2 migration, unknown major rejection, digest/path checks, immutable attempts, stale writers, dependency-ordered publication, injected write failures, and prior-state recovery. | `internal/sprint/qa_state_test.go` |
| Cross-surface fixtures | App, CLI text/JSON, TUI, browser HTML/JSON, durable run detail, `qa.md`, `smoke.md`, and verification state agree on identities, status, assessment, blockers, and next action. | Shared fixture suite named in plan tasks |
| Race and static checks | Parallel starts, stale owners, cancellation/completion, cleanup, terminal arbitration, vet, build, and whitespace checks. | `go test -race ./...`; `go vet ./...`; `go build ./cmd/ultraplan`; `git diff --check` |
| Runtime dogfood | One real generated check remains contained, target and Git identities remain unchanged, cleanup is certain, adjudication records a rejection or promotion, and smoke parity uses the same external evidence authority. | Gated Deep Smoke Sprint protocol and manifest-declared harness run |
| Review | Architecture boundaries, requirement conformance, security, issue bounds, state recovery, compatibility, and no prohibited scope. | Architecture Review and Sprint Review protocols |
| Documentation | Commands, architecture, user workflow, browser operation, recovery, JSON/state schemas, and release checks match executable behavior. | `docs/cli-reference.md`, `docs/architecture.md`, `docs/user-guide.md`, `docs/local-web.md`, `docs/recovery.md`, `docs/phase3-json-schemas.md`, `docs/release-checklist.md` |

## Assumptions And Risks

| Item | Type | Impact | Mitigation / Follow-Up |
| --- | --- | --- | --- |
| Sprint 36 detailed verification state uses schema version 1 and supplies all admission facts named by AC-01. | Assumption | Schema 2 migration and writable admission depend on it. | Validate exact version and required fields before writing schema 2; missing facts block rather than infer. |
| Linux `bwrap` or an equivalent selected launcher can prove mount, PID, network, capability, parent-death, and descendant guarantees on supported hosts. | Assumption | Without it, writable QA cannot run. | Capability preflight before generated writes; gated real-host test; actionable blocked status. |
| Full target manifests fit the 500,000-entry and 8 GiB copy limits for representative use. | Assumption | Larger targets block or require narrowed scope. | Measure dogfood and revise limits only with aggregate budget evidence. |
| Target or Git control state can change concurrently. | Risk | Evidence can become stale after a valid setup. | Recheck identity before copy, after copy, before first write, and after command; preserve but do not promote stale evidence. |
| Symlink replacement can race validation and copying. | Risk | A path-string implementation could read outside the target. | Use descriptor-relative traversal or equivalent no-follow/revalidation mechanics and adversarial race tests. |
| A descendant may escape process-group signaling. | Risk | Cleanup could be falsely reported as complete. | Use PID namespace/cgroup-equivalent containment, post-kill probing, bounded waits, and `cleanup_uncertain` on doubt. |
| Generated source is hostile even when its path is allowed. | Risk | Compiler/test execution can consume resources or attempt escape. | Read-only execution mount, no network, capability drop, exact commands, time/output/storage caps, and descendant cleanup. |
| Two repeated runs may share a host toolchain fault. | Risk | Matching failures may still reflect setup rather than product behavior. | Record setup/toolchain identity, validate setup, retain contradictory evidence, and allow bounded adjudication follow-up. |
| Truncation may remove the decisive diagnostic. | Risk | Evidence could be misclassified. | Require complete structured decisive output or mark inconclusive; show truncation in every projection. |
| A crash can occur between root state and report publication. | Risk | Readers could see mismatched current state and `qa.md`. | Journal the generation, retain the prior pair, sync directories, and fault-test every boundary. |
| Failed reruns can be mistaken for the retained canonical pass. | Risk | Operators could accept stale evidence. | Carry `current_attempt_id`, `canonical_attempt_id`, freshness, and retained-report labels through every DTO and artifact view. |
| Smoke alias defaults can drift despite shared code. | Risk | Compatibility can break without an obvious new parser. | Compare every operational field and exact `smoke.md` output under shared fixtures and gated harness execution. |
| Hostile text can contain invalid UTF-8, ANSI/OSC, HTML, Markdown, or long lines. | Risk | Terminal control injection, unsafe rendering, or resource exhaustion. | Normalize invalid UTF-8, strip terminal controls, use `html/template`, render plain text, cap fields, and test hostile fixtures. |
| Durable run success may coexist with QA `fail` or `blocked`. | Risk | Adapters may collapse orchestration completion into product pass. | Separate lifecycle, attempt status, and assessment types in DTOs and parity fixtures. |
| Fixed limits may be too strict or too loose. | Risk | Valid evidence may block, or aggregate cost may be high. | Aggregate caps always win; Sprint 39 measures and justifies any change. |

## Implementation Constraints

- Do not add QA or repair to `PlanningStage` or canonical planning artifact order.
- Do not start writable child work unless the current Sprint 36 gate and every isolation capability pass.
- Keep `internal/sprint` responsible for QA semantics and `internal/platform/process` free of sprint types and verdict logic.
- Keep `internal/web` dependent only on `internal/app`; preserve the import-boundary test.
- Use Agentwrap for investigator/adjudicator runtime work and the existing process seam for external commands.
- Use explicit executable/argv only. Never execute a shell string or command parsed from Markdown/model prose.
- Permit only product-validated new tests, fixtures, and data probes in the isolated copy. Do not modify existing production or test files.
- Exclude target and Git mutation. Detect drift, preserve user changes, and block promotion rather than reverting anything.
- Execute copied source read-only. Writable cache/temp paths must be attempt-local and declared before execution.
- Serialize writable attempts in Sprint 37. Read-only queries and observers may run concurrently.
- Stop new scheduling on cancellation, propagate cancellation to active runtime/process work, then perform bounded cleanup under a separate context.
- Treat cleanup uncertainty, stale identity, truncation of decisive evidence, narrow smoke, diagnostic smoke, malformed records, and unsupported state versions as unable to pass.
- Let only deterministic product code promote issues, classify repair eligibility/regression candidates, and compute assessment.
- Keep Conformance Review verdict independent and read-only to QA.
- Keep raw smoke JSON, output, artifacts, and harness issue files under manifest-declared external roots.
- Preserve `smoke`, `smoke.md`, existing clients, and the existing smoke operation kind.
- Keep detailed attempts under `verification/`; keep `flow-state.json` bounded.
- Publish dependency-ordered immutable attempts and recover the `state.json`/`qa.md` pair without losing the prior valid generation.
- Reuse Sprint 35 durable acceptance, fencing, liveness, replay, cancellation, retention, reconciliation, and terminal arbitration. Verification state must not duplicate those duties.
- Use bounded summary/detail app DTOs. Adapters must not read verification files, derive freshness, group issues, or compute assessment.
- Keep browser observation available across session rotation while requiring fresh authorization for starts and cancellation.
- Preserve a complete no-JavaScript browser path and escape hostile evidence for each output medium.
- Do not apply generated patches, repair production, promote permanent tests, mutate Git, add a generic issue tracker, or introduce post-Sprint-39 systems.
- Update all required documentation and pass the full test, race, vet, build, diff, review, and gated dogfood evidence set.

## Plan Handoff

`plan.md` must execute these decisions without reopening architecture. It must order the work so that generic isolation and smoke delegation seams exist before writable orchestration, and state/publication contracts exist before adapter projections claim current QA truth.

The plan must carry forward:

- all nine final decisions and their ownership boundaries;
- AC-01 through AC-27 mappings and the selected contracts;
- Linux fail-closed full-copy isolation and sequential writable attempts;
- product-mediated generated additions, frozen plans, scoped identities, and deterministic adjudication;
- schema version 2, immutable attempts, publication journal, prior-generation recovery, and bounded flow projection;
- direct smoke delegation and exhaustive parity evidence;
- typed app requests, bounded summary/detail reads, durable run reuse, stable errors, and guarded browser operation;
- current-versus-canonical presentation, hostile-text handling, no-JavaScript behavior, TUI operation, and accessibility;
- the fixed QASettings defaults and aggregate caps;
- offline fault, race, compatibility, cross-surface, review, and gated dogfood evidence;
- the Architecture Review, Sprint Review, and Deep Smoke Sprint protocols.

## Phase Exit Criteria

- [x] Selected context was read and used.
- [x] API Design, Architecture, and Frontend reasoning documents were completed and synthesized.
- [x] Area-specific conclusions are reflected without silent overrides.
- [x] Contracts and acceptance criteria are mapped to decisions and expected evidence.
- [x] Final decisions fix ownership, isolation, identity, adjudication, state, smoke, adapters, presentation, limits, and testing.
- [x] Expected evidence is specific and reviewable.
- [x] Rejected alternatives and accepted debt are explicit.
- [x] `plan.md` can execute the sprint without choosing new core architecture.
