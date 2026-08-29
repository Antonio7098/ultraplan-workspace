# Sprint Reasoning

> **Inputs Used:** `projects/ultraplan-go/sprints/37-evidence-qa-smoke/requirements.md`, `projects/ultraplan-go/project-index.md`, `projects/ultraplan-go/roadmap.md`, `projects/ultraplan-go/docs/ARCHITECTURE.md`, `projects/ultraplan-go/docs/PRD.md`, `projects/ultraplan-go/docs/TRD.md`, `projects/ultraplan-go/sprints/37-evidence-qa-smoke/sprint-index.md`, `projects/ultraplan-go/sprints/37-evidence-qa-smoke/technical-handbook.md`, `projects/ultraplan-go/sprints/37-evidence-qa-smoke/reasoning/api-design.md`, `projects/ultraplan-go/sprints/37-evidence-qa-smoke/reasoning/architecture.md`, `projects/ultraplan-go/sprints/37-evidence-qa-smoke/reasoning/frontend.md`, `projects/ultraplan-go/all-evidence-qa-review.md`, `projects/ultraplan-go/all-evidence-smoke-report.md`, `projects/ultraplan-go/all-evidence-context-review.md`, `projects/ultraplan-go/all-evidence-reasoning-review.md`, `projects/ultraplan-go/all-evidence-requirements-review.md`, `projects/ultraplan-go/all-evidence-sprint-index-review.md`, `projects/ultraplan-go/all-evidence-technical-handbook-review.md`, `projects/ultraplan-go/sprints/35-durable-run-observability/reasoning.md`, `projects/ultraplan-go/sprints/35-durable-run-observability/plan.md`, `projects/ultraplan-go/sprints/35-durable-run-observability/execute.md`, `projects/ultraplan-go/sprints/35-durable-run-observability/review.md`, `projects/ultraplan-go/sprints/35-durable-run-observability/smoke.md`, `projects/ultraplan-go/sprints/35-durable-run-observability/flow-state.json`, `projects/ultraplan-go/sprints/36-read-only-qa/reasoning.md`, `projects/ultraplan-go/sprints/36-read-only-qa/plan.md`, `projects/ultraplan-go/sprints/36-read-only-qa/execute.md`, `projects/ultraplan-go/sprints/36-read-only-qa/review.md`, `projects/ultraplan-go/sprints/36-read-only-qa/smoke.md`, `projects/ultraplan-go/sprints/36-read-only-qa/flow-state.json`, `system/contracts/core/architecture.md`, `system/contracts/core/errors.md`, `system/contracts/core/configuration.md`, `system/contracts/core/observability.md`, `system/contracts/core/security.md`, `system/contracts/core/testing.md`, `system/contracts/core/documentation.md`, `system/contracts/surfaces/cli.md`, `system/contracts/runtime/llm.md`, `system/contracts/runtime/llm-evaluation-cost-safety.md`, `system/contracts/runtime/workflows.md`, `system/contracts/runtime/performance.md`, `system/contracts/runtime/persistence-and-migrations.md`, `system/protocols/architecture-review-protocol.md`, `system/protocols/review-sprint-protocol.md`, `system/protocols/deep-smoke-sprint-protocol.md`, `source-repos/opencode`, `source-repos/plannotator`, `source-repos/written`, `source-repos/judge0`, `internal/sprint/qa.go`, `internal/sprint/qa_state.go`, `internal/sprint/qa_test.go`, `internal/sprint/smoke.go`, `internal/sprint/smoke_author.go`, `internal/sprint/verify.go`, `internal/app/sprint_usecases.go`, `internal/app/operations.go`, `internal/app/operation_runner.go`, `internal/app/durable_operations.go`

## 1. Executive Summary

### Decisions in This Sprint

Sprint 37 will extend the Sprint 36 read-only QA path into an evidence-backed verifier. `internal/sprint` will own the rubric, evidence policy, verdict rules, adjudication, QA state, and canonical smoke-suite selection. `internal/platform/process` will own generic copy-workspace and child-process isolation. The app layer will expose bounded QA DTOs and operations. CLI, TUI, and browser adapters will only validate, invoke app use cases, and render returned facts.

The verifier will derive its scope from immutable governed planning artifacts plus the active implementation fingerprint. It will require typed checks for fact, negative, deterministic behavioral, LLM semantic, and adversarial prompt or content-containment behavior. Missing evidence, malformed evidence, provider failure, permission failure, timeout, cancellation, and incomplete adjudication are explicit non-pass outcomes.

Each worker will run against a fresh copy workspace. The original implementation target and governed project inputs will remain read-only and will not be named in child prompts, arguments, environment, or process metadata. Sprint 37 will schedule those copy-backed workers sequentially. This costs time but makes isolation attribution, cancellation, and deterministic evidence review much easier to trust.

An initially failed shard will receive three fresh evaluator calls against its immutable evidence bundle. Adjudication may change that shard to pass only when all three evaluator calls complete and at least two return locally validated pass verdicts. Otherwise the shard remains failed or becomes blocked when evaluation itself did not complete. A whole QA run passes only when coverage is complete and every required shard passes.

The smoke-suite path will be an additive `Suite: "smoke"` selection on the existing QA operation family. It will reuse the canonical `review`, `review.md`, `smoke`, `smoke.md`, `PlanningStage`, run-control, and smoke orchestration paths. It will not create a second smoke executor, apply repairs, mutate Git, or move external raw harness evidence into product-owned storage.

### Confidence Assessment

- **High confidence:** ownership and dependency direction; target and governed-input immutability; fail-closed result semantics; additive adapter design; canonical smoke reuse; Sprint 35 run-control authority; current-pointer query semantics. These follow directly from Sprint 35 and Sprint 36 behavior, the project architecture, the selected core contracts, and all three area-reasoning documents.
- **Medium confidence:** practical runtime cost of three-call semantic analysis and three-call failed-shard adjudication; the exact balance between copied workspace size and direct behavioral-test needs; malicious-provider behavior across every supported runtime adapter. The architecture is fixed, but measured limits and adapter-specific evidence must come from implementation tests and deep smoke.
- **Unresolved high-severity warnings:** none at the architecture level. High-severity implementation risks remain around symlink/path escapes, original-target leakage, stale fencing, evaluator incompleteness, and compatibility during private state versioning. Section 2 and Section 9 assign a mitigation and required proof to each risk.
- **Recommendation:** proceed to `plan.md`. The plan must implement and verify these decisions rather than revisit ownership, isolation, quorum, persistence, or smoke architecture.

## 2. Problem Framing

### Problem Statement

Sprint 36 can map changed paths into read-only QA shards and report bounded evidence, but it does not prove that the implemented sprint satisfies governed requirements under hostile content, failed providers, malicious write attempts, incomplete evidence, or independent evaluation. Sprint 37 must turn those maps into reviewable acceptance evidence without giving an analyzer or smoke author authority to alter the implementation it is judging.

### Why It Is Non-Trivial

The verifier must coordinate deterministic tests, source observations, LLM-backed semantic checks, process isolation, durable operation ownership, resumable product state, cancellation, and canonical smoke execution. These concerns cross package boundaries but cannot be collapsed into the CLI or a generic platform package.

Read-only intent is not proof. A worker can receive an original path through an environment variable, follow a symlink out of a copy, issue a write-capable tool call that happens to fail, or produce a plausible pass with missing evidence. The system therefore needs native path and permission boundaries, runtime policy enforcement, before-and-after identities, attributed write-attempt events, and fail-closed evidence validation.

LLM agreement is also not correctness. Three calls can share the same bias or bad evidence. The quorum is useful only after direct evidence has been collected, outputs have passed local schema validation, sessions are fresh, and the report preserves disagreement instead of flattening it.

Finally, Sprint 37 must add public progress and evidence without breaking existing QA scripts, app JSON v1, CLI names, TUI/browser behavior, or Sprint 35 run-control semantics.

### Constraints

- `internal/sprint` owns product QA and smoke policy. Generic platform packages may not import product packages or encode Sprint 37 semantics.
- `internal/platform/process` may provide product-neutral copy, permission, process, timeout, cleanup, and event-attribution mechanics behind a narrow consumer-owned port.
- The original implementation target and governed project inputs are read-only. Copy workspaces, private QA state, canonical QA summaries, and existing external smoke evidence roots are the only intentional write locations.
- Child prompts, arguments, environment, working-directory metadata, and emitted process metadata must not contain the original target path.
- Every user-influenced path must be normalized, bounded to an allowed root, and checked against traversal, absolute-path, and symlink escapes before work begins.
- LLM input and output use explicit versioned schemas. Provider schema support is not an authority; local validation is mandatory.
- Three analyzer calls are required only for LLM semantic and adversarial checks. Fact, negative, and deterministic behavioral checks use their declared direct observations and deterministic analyzers. A check passes only when every required direct observation completes and a strict majority of its declared analyzer verdicts passes.
- Failed-shard adjudication always uses exactly three fresh evaluator calls. No evaluator session may be reused from analysis or from another shard.
- Worker copy workspaces run sequentially in Sprint 37. Provider calls inside one semantic check may run only within the declared per-check cap and must preserve individual call identity.
- Run control, not QA state, remains authoritative for acceptance, ownership, fencing generation, cancellation, terminalization, and reconciliation.
- Private QA state advances to schema version 2 in a new namespace. Existing `qa-v1` records remain readable for compatibility but are never mutated into mixed-version records.
- The public app JSON envelope remains schema version 1. Sprint 37 adds optional fields and focused query results only.
- Existing `review`, `review.md`, `smoke`, `smoke.md`, `PlanningStage`, and external harness ownership remain unchanged.
- Production repair, patch application, general issue tracking, Git mutation, alternate persistence, hosted operation, and a new retrieval/content-identity subsystem are out of scope.

### Hidden Complexities

- A changed implementation may alter behavior outside the changed-file list. Coverage must combine changed paths, requirement ownership, affected tests, and the execution delta rather than equate diff coverage with requirement coverage.
- A copied tree can still expose the original through symlinks, inherited environment, command-line arguments, `/proc`-visible metadata, or diagnostics. Redaction alone is insufficient because the child must never receive the value.
- Native denial and runtime-policy denial prove different things. Native permissions protect the filesystem. Runtime events prove whether an LLM-backed worker attempted a forbidden action. Both are required.
- A concurrent user edit can change target identity during QA without being caused by QA. Event attribution distinguishes an attempted QA write from an external change; either case invalidates freshness, but only the former is a permission violation.
- Cancellation can race with a child result or durable terminal proposal. Sprint 35 fencing and deterministic terminal precedence must decide the winner before QA state projects it.
- Current-pointer records can be valid while their target run record is missing, malformed, or digest-mismatched. Recovery must fail closed and report the specific integrity error.
- Semantic evaluators can be influenced by hostile repository text. Rubric, evidence, and repository excerpts must be separate typed fields, and repository text must be labeled as untrusted data.
- Browser progress streams can drop events. The durable status query, not the event stream, is the recovery authority.
- Smoke authoring and QA verification have different write boundaries. Smoke authoring may write only manifest-allowed harness paths. QA analyzers may write only disposable copy workspaces and private QA state through the parent service.

### Assumptions

- Sprint 35 run-control APIs and fencing remain available and authoritative.
- Sprint 36 already provides a deterministic map, changed-path coverage, bounded shard identities, target identity, private QA records, app/CLI/TUI/browser entry points, and recovery hooks that Sprint 37 can extend additively.
- The runtime can expose permission capability and attributed tool events. If an adapter cannot do so, checks requiring that capability block before execution.
- A fresh session means a new provider/runtime conversation identity with no retained turns or tools from analysis, another evaluator, or another shard.
- Existing configuration supplies runtime and model identity. Sprint 37 does not introduce a silent model fallback or environment-variable lookup in business logic.
- External smoke harness evidence remains immutable and addressable by the canonical smoke service.
- Copying the bounded selected target is acceptable for Sprint 37 workloads. Limits will reject targets that exceed the declared file-count, byte, or timeout budgets rather than degrade into an unbounded scan.

### Risks

| Risk | Consequence | Mitigation and required evidence |
|---|---|---|
| Original target path leaks to a child | A hostile worker can address the real tree | Construct requests from opaque target identity and copy-relative paths; scrub inherited environment; assert the path is absent from prompts, argv, env, cwd metadata, logs, and events |
| Symlink or traversal escape | A copied workspace reaches protected files | Canonicalize roots, reject escaping links and absolute paths before copy or process launch, and add Linux integration tests for traversal and symlink cases |
| Runtime claims restrictions but does not enforce them | Write attempts or hidden tools may run | Require permission capability during preflight; block unsupported adapters; pair runtime policy with native filesystem permissions and attributed events |
| Shared evaluator bias | Majority produces a confident wrong verdict | Require direct observations, fresh sessions, local schema checks, preserved per-call rationale, deterministic behavioral tests, and Conformance Review |
| Evaluator/provider failure is mistaken for a product failure or pass | Misclassified acceptance | Record provider, timeout, cancellation, malformed-output, and missing-call outcomes separately; treat them as blocked non-pass |
| Sequential workers exceed practical duration | Poor operator experience | Apply explicit target, check, provider-call, timeout, and cleanup budgets; show progress; measure representative repositories before changing concurrency |
| State v2 breaks Sprint 36 recovery | Existing evidence becomes unreadable | Use a separate version namespace, retain read-only v1 decoding, reject unsupported versions explicitly, and test v1 read plus v2 write/recovery |
| Current pointer or summary diverges from run state | Stale or misleading status | Commit run record before atomically replacing the pointer, store digests, reconcile through run control, and expose freshness reasons |
| Browser/TUI present stale progress | Operator acts on obsolete state | Refresh from the authoritative focused status query after completion, reconnect, sequence gap, or visibility return |
| Smoke suite forks canonical behavior | Conflicting acceptance evidence | Route `Suite: "smoke"` through the existing smoke service and artifact path; add a test that legacy and QA-suite invocations select the same canonical route |

## 3. Inputs and Source Grounding

### Sprint Index Summary

The sprint index selects the Sprint 35 durable run-control chain, the Sprint 36 read-only QA chain, the project architecture and product requirements, all current evidence review reports, and the implementation seams around QA, smoke, app use cases, operations, durable ownership, CLI, TUI, and web adapters. It excludes hosted-operation work, alternate persistence, unrelated study subsystems, and a new retrieval/content-identity subsystem.

That selection fixes the integration point. Sprint 37 extends the existing QA operation and state path. It does not introduce a parallel workflow engine, a new top-level acceptance command, or a second smoke implementation.

### Technical Handbook Findings Applied

- Treat read-only verification as defense in depth. Native permissions, runtime capability checks, path containment, target identity snapshots, and attributed events cover different failure modes.
- Use three-valued evidence states in the internal model: pass, fail, and blocked. Provider errors and missing evidence cannot become failed assertions or successful empty results.
- Keep direct observations primary. LLM analyzers interpret bounded evidence; they do not replace deterministic tests or local validation.
- Make prompt and output schemas explicit and versioned. Separate rubric instructions from untrusted repository content.
- Preserve each analyzer and evaluator result. Majority is a deterministic decision rule, not proof of independent truth.
- Use fixed budgets for context, calls, retries, timeouts, file counts, bytes, and cleanup. Reject over-budget work before provider spend.
- Preserve Sprint 35 ownership and terminal precedence. QA state projects run-control truth and must not invent a competing lifecycle.
- Reuse canonical smoke discovery, selection, execution, evidence validation, artifact writing, and overall assessment.

### Repos Studied / Source Evidence Used

| Repository or report | Why it mattered | Decisions influenced |
|---|---|---|
| `source-repos/opencode` | Demonstrates explicit permission policy, session isolation, tool-event visibility, and runtime lifecycle concerns | D-03 worker isolation, D-04 structured LLM calls, D-08 observability and cancellation |
| `source-repos/judge0` | Demonstrates process isolation, bounded execution, cleanup, and the difference between an application guard and a hostile sandbox | D-03 copy workspaces, sequential scheduling, resource limits, fail-closed cleanup |
| `source-repos/plannotator` | Demonstrates typed model-facing contracts and local validation around generated structured results | D-04 check schema, analyzer/evaluator result validation, additive app DTOs |
| `source-repos/written` | Demonstrates explicit judging criteria and reviewable generated evidence rather than accepting prose at face value | D-02 evidence map, D-04 rubric and adjudication, D-09 reporting |
| `projects/ultraplan-go/all-evidence-context-review.md` | Confirms that selected context and governing artifacts must remain explicit and traceable | D-01 ownership, D-02 source selection, D-10 review evidence |
| `projects/ultraplan-go/all-evidence-reasoning-review.md` | Identifies the need for final decisions, rejected alternatives, and contract traceability | All decisions and Sections 6 through 10 |
| `projects/ultraplan-go/all-evidence-requirements-review.md` | Grounds requirement completeness and non-goal boundaries | D-02 through D-10 |
| `projects/ultraplan-go/all-evidence-sprint-index-review.md` | Confirms the selected Sprint 35/Sprint 36 lineage and excluded subsystems | D-01, D-05, D-06 |
| `projects/ultraplan-go/all-evidence-technical-handbook-review.md` | Warns against overstating isolation, evaluator independence, and research completeness | D-03, D-04, D-08 |
| `projects/ultraplan-go/all-evidence-qa-review.md` | Identifies gaps between deterministic QA mapping and acceptance-grade evidence | D-02, D-04, D-09 |
| `projects/ultraplan-go/all-evidence-smoke-report.md` | Grounds canonical smoke ownership, external evidence authority, and real-adapter proof | D-06 and D-10 |
| Sprint 35 reasoning, plan, execute, review, smoke, and flow state | Establish durable acceptance, ownership, fencing, terminal precedence, and reconciliation as shipped behavior | D-05 and D-08 |
| Sprint 36 reasoning, plan, execute, review, smoke, and flow state | Establish deterministic QA maps, read-only target handling, private QA state, adapter surfaces, and compatibility baseline | D-01, D-02, D-05, D-07 |
| Current `internal/sprint`, `internal/app`, and smoke source excerpts selected by the sprint index | Confirm concrete extension seams and canonical service paths | D-01, D-05, D-06, D-07 |

The studied repositories are source evidence, not contracts. UltraPlan's selected requirements, contracts, and shipped Sprint 35/Sprint 36 behavior take precedence where a studied approach differs.

### Area Reasoning Summary

- `reasoning/architecture.md` places QA policy in `internal/sprint` and generic process mechanics in `internal/platform/process`. It selects fresh per-worker copies, sequential scheduling, native permissions plus runtime policy, run-control authority, additive state evolution, and canonical smoke reuse.
- `reasoning/api-design.md` selects additive app, CLI, and web contracts. It keeps app JSON v1, adds a closed `Suite: "smoke"` selector and focused bounded queries, retains current-pointer-only public projection, uses a new private state namespace, and rejects smoke-suite resume.
- `reasoning/frontend.md` selects one bounded app projection rendered separately by TUI and browser adapters. The browser keeps a server-rendered no-JavaScript baseline, both adapters refresh authoritative state after progress gaps, and hostile content remains inert text.

The final decisions retain those conclusions with two clarifications. The copy-workspace boundary is not described as a hostile multi-tenant sandbox, and majority verdicts are not described as proof of independent correctness. Both claims require narrower language and stronger review evidence.

### External Sources Consulted

No additional external sources were consulted. The selected handbook, reports, studied repositories, current source excerpts, and contracts were sufficient. Adding unrelated sources would weaken the sprint-index boundary without resolving a specific decision.

## 4. Decision Register

| ID | Decision | Rationale | Trade-offs | Contracts / Requirements | Evidence Expected |
|---|---|---|---|---|---|
| D-01 | Keep QA policy in `internal/sprint`; add generic isolation mechanics to `internal/platform/process`; keep adapters thin | Preserves inward dependencies and extends shipped seams | More explicit ports and composition work | R3, R10, R11, R35, R38, R41-R44; ARCH-CORE-001, ARCH-CORE-002, ARCH-LAYER-002, ARCH-ENTRY-001, ARCH-SHARED-001 | Import review, public-seam tests, adapter tests |
| D-02 | Build an immutable evidence map from governed planning inputs, implementation fingerprint, changed paths, requirement ownership, and typed checks | Prevents plausible but incomplete QA | Map invalidation and richer records | R1-R9, R12-R20, R23-R24, R33-R34, R46; PERSIST-SCHEMA-001, PERSIST-INTEGRITY-001, TEST-UNIT-001 | Determinism, invalidation, completeness, contradiction, and stale-input tests |
| D-03 | Run each worker sequentially in a fresh copy workspace; never disclose the original target path to children; enforce native and runtime restrictions | Makes target immutability reviewable and attribution deterministic | Copy cost and lower throughput | R10-R11, R21, R29, R35-R36, R47-R48; SEC-FILES-001, SEC-INPUT-001, EVAL-SAFETY-001, PERF-BOUND-001, PERF-CONC-001 | Permission, path-leak, symlink, malicious-write, timeout, cancellation, cleanup, and target-identity tests |
| D-04 | Use typed direct checks, three-call semantic/adversarial analysis, and three fresh failed-shard evaluators with local validation and strict majority | Combines deterministic proof with bounded semantic review | Provider cost; majority can share bias | R17-R18, R22, R25-R32, R47-R50; LLM-IO-001, LLM-PROMPT-001, LLM-RETRY-001, EVAL-HUMAN-001, EVAL-COST-001 | Schema, quorum, fresh-session, disagreement, provider-failure, hostile-content, and all-check smoke evidence |
| D-05 | Keep run control authoritative; evolve private QA state to v2 in a new namespace; retain read-only v1 compatibility; expose only current bounded projections | Avoids dual lifecycle truth and mixed-version mutation | Two readers during compatibility window | R37-R40, R42, R51; PERSIST-ATOMIC-001, PERSIST-MIG-001, PERSIST-DERIVED-001, WF-STATE-001 | v1-read/v2-write tests, atomic pointer tests, fencing races, crash recovery, retention, freshness |
| D-06 | Add `Suite: "smoke"` to the existing QA operation family and route it through canonical smoke orchestration; reject suite resume | Gives one QA entry point without forking smoke truth | Selector adds a public compatibility branch | R19, R45-R46, R53; TEST-SMOKE-001, TEST-SMOKE-002, WF-BOUNDARY-001, CLI-SHAPE-001 | Legacy smoke parity, QA-suite route identity, no-resume validation, real-adapter smoke |
| D-07 | Extend app JSON v1 additively with focused queries and bounded summaries; render one projection in CLI, TUI, and browser; keep no-JS browser support | Preserves automation and avoids adapter-owned policy | DTO growth and renderer work | R7, R38-R44, R51; CLI-IO-001, CLI-JSON-001, TEST-CONTRACT-001, ARCH-ENTRY-001 | CLI snapshots, JSON compatibility, TUI tests, browser route/template tests, event-gap refresh |
| D-08 | Make every failure structured and non-pass; propagate cancellation; preserve correlation, prompt, provider, permission, and cleanup facts with redaction | Empty or generic outcomes cannot support acceptance | More explicit error/state variants | R8, R13, R21, R25, R28-R30, R36, R40, R47-R48, R51; ERR-CORE-001, ERR-CODE-001, OBS-CORE-001, OBS-CORR-001, CLI-LIFE-001 | Failure matrix, redaction tests, cancellation race tests, dropped-event recovery |
| D-09 | Define fixed scan, evidence, provider-call, timeout, retry, concurrency, history, and response-size budgets before work | Prevents unbounded cost and memory use | Oversized targets block instead of degrading silently | R9, R12, R18, R21, R40, R51-R52; PERF-BUDGET-001, PERF-BOUND-001, PERF-COST-001, CFG-TYPE-001, CFG-START-001 | Boundary tests, over-budget rejection, progress, representative duration/cost records |
| D-10 | Require contract review, sprint review, and a hostile real-adapter deep smoke before acceptance | Unit tests cannot prove runtime policy or target immutability | Acceptance takes longer and may spend provider budget | R5, R15-R20, R22-R32, R46-R53; TEST-E2E-001, TEST-FAIL-001, DOC-ARCH-001, DOC-OPS-001 | Architecture Review, Sprint Review, Deep Smoke Sprint, exact run references |

## 5. Detailed Architecture / Design Reasoning

### D-01: Ownership and dependency shape

`internal/sprint` will define consumer-owned interfaces for the process isolation operations it needs. The product service will compose the evidence map, select checks, create worker requests, validate results, adjudicate failed shards, project status, and route the smoke suite. It will depend on run-control and platform behavior through narrow seams.

`internal/platform/process` will implement copy-root creation, validated relative-path materialization, structured child execution, environment construction, native permission application, timeout, cancellation, bounded output capture, cleanup, and process events. It will know nothing about requirements, shards, rubrics, smoke verdicts, or QA state.

The app layer will expose use-case DTOs and operation kinds. CLI, TUI, and web packages will not read QA files, invoke providers, inspect workspaces, or calculate verdicts. Composition remains in established registrars and bootstrap paths.

Grounding comes from the project architecture, Sprint 35 and Sprint 36 extension seams, `technical-handbook.md`, and ARCH-CORE-001/002, ARCH-LAYER-002, ARCH-ENTRY-001, ARCH-SHARED-001, LLM-BOUNDARY-001, and WF-BOUNDARY-001. The trade-off is an additional consumer-owned process port rather than direct use of concrete helpers. That is deliberate because subprocess execution, permissions, and cleanup are consequential seams. The rejected alternative is product policy inside the process package, which would reverse dependency direction and make generic code own acceptance semantics.

### D-02: Evidence map and completeness

The evidence map is immutable for one run and has an explicit schema version. Its identity covers the governed-input fingerprint, implementation fingerprint, map policy version, check catalog version, target identity, and sorted check definitions. Each requirement records ownership, source references, one or more checks, expected evidence kinds, affected paths or tests, and coverage state.

The check catalog supports five closed kinds:

- `fact`: prove a required artifact, field, route, state, or ownership relation exists.
- `negative`: prove a forbidden behavior or dependency is absent using a bounded search or explicit rejection test.
- `behavioral`: execute a deterministic test or command and validate stable fields, exit status, and artifacts.
- `semantic`: ask bounded LLM analyzers to judge evidence against a versioned rubric after direct source evidence is collected.
- `adversarial`: execute hostile content or permission scenarios and verify containment plus result classification.

Every check names required observations, analyzer count, timeout, retry budget, accepted result schema, and evidence references. Fact, negative, and behavioral checks use deterministic analyzers and direct observations. Semantic and adversarial checks use exactly three fresh LLM analyzer sessions unless the adversarial check is wholly deterministic; their strict pass majority is two of three. A check with one declared deterministic analyzer has a strict majority only when that analyzer passes. Zero-analyzer checks are invalid.

Map creation fails before provider work when a requirement has no owner, a check has no observation, a changed path has no shard, a selected source is missing, fingerprints are unavailable, or the catalog contradicts the requirements. Existing Sprint 36 artifact IDs remain stable where their meanings are unchanged; Sprint 37 fields are additive or use a new v2 private record type.

Grounding comes from the QA and context evidence reports, Sprint 36, Plannotator, Written, and PERSIST-SCHEMA-001, PERSIST-INTEGRITY-001, TEST-UNIT-001, and EVAL-SCOPE-001. The accepted debt is that requirement-to-check mapping remains product-authored rather than inferred by a new retrieval system. The rejected alternative is diff-only coverage, which misses requirements affected by behavior outside changed files.

### D-03: Isolation, scheduling, and target immutability

The parent process resolves and fingerprints the original target, validates all selected paths, and creates a fresh bounded copy for each worker. It passes only copy-relative paths and an opaque target identity to the child. It does not pass the original target path through the prompt, argv, environment, child working-directory description, metadata, or progress events.

The parent strips inherited environment to an allowlist, rejects absolute and escaping paths, rejects or safely materializes symlinks according to a single documented policy, applies native permissions to protect non-copy roots, and requests restricted runtime permissions. A worker may mutate its disposable copy only when a direct behavioral test requires it. It may not write private QA state, canonical project artifacts, the original target, governed project inputs, or external smoke roots.

Before and after each worker, the parent compares the original target and governed-input identities. It also inspects attributed runtime events for write-capable calls against protected roots. An attributed attempt is a permission failure even when the native layer denied the write. An unattributed concurrent target change invalidates the run as stale and reports that distinction. Cleanup failure is visible and non-pass; cleanup does not erase the primary failure.

Sprint 37 runs copy-backed workers sequentially. Within one semantic check, up to three declared provider calls may be in flight only if the runtime adapter and evidence collector preserve per-call identity, cancellation, and complete failure aggregation. The implementation may choose sequential provider calls first; it may not exceed the declared cap. Any later increase in worker concurrency requires new race, isolation, performance, and cancellation evidence.

Grounding comes from OpenCode, Judge0, the handbook's isolation cautions, Sprint 36 target identity, SEC-FILES-001, EVAL-SAFETY-001, PERF-BOUND-001, and PERF-CONC-001. The accepted trade-off is copy and latency cost. The rejected alternatives are a shared writable workspace, which permits cross-worker contamination, and prompt-only read-only instructions, which provide no enforceable boundary.

### D-04: Verdicts and failed-shard adjudication

Each observation and analyzer result is immutable and carries check ID, shard ID, call ID, role, session ID, prompt ID and version, model/provider identity, input evidence digest, timestamps, duration, structured outcome, rationale, citations, retry count, and error when present. Local code validates the complete result. Provider-side structured output support does not replace local validation.

The check result rules are fixed:

- `pass`: every required direct observation completed successfully, every required analyzer call completed with a valid result, and a strict majority of declared analyzer verdicts is pass.
- `fail`: direct evidence proves the assertion false, or all required calls complete and the strict analyzer majority is fail.
- `blocked`: required evidence or a required call is missing, malformed, cancelled, timed out, permission-denied, provider-failed, unsupported, or otherwise incomplete.

The shard first aggregates its checks. If every check passes, the shard passes without adjudication. If a deterministic or analyzer result initially fails, the service freezes the complete shard evidence bundle and starts exactly three fresh evaluator sessions. Evaluators receive the rubric, requirement text, immutable observations, and analyzer outputs as typed data. They receive no analyzer session history and no authority to run tools or alter evidence.

Adjudication is admissible only when all three evaluator calls complete, validate locally, reference the same evidence digest, and at least two say pass. In that case the shard's final result is pass with an explicit `adjudicated` marker and preserved initial failure. If all calls complete without a pass majority, the shard is fail. If any evaluator call is incomplete or invalid, the shard is blocked. Deterministic security violations, target mutation, governed-input mutation, missing coverage, and stale fingerprints are not overridable by evaluator opinion.

The run passes only when every required shard has a final pass, all expected evidence is present, target and governed-input identities remain current, and run-control terminalization succeeds. The report preserves all disagreements and call-level errors.

Grounding comes from Written, Plannotator, the handbook's warning about correlated evaluators, LLM-IO-001, LLM-PROMPT-001, LLM-LIFECYCLE-001, EVAL-HUMAN-001, EVAL-REG-001, and EVAL-COST-001. The accepted trade-off is bounded provider cost for failed shards. The known debt is that fresh sessions do not guarantee statistical independence. The rejected alternative is a single judge, which is too fragile, and automatic trust in model confidence, which has no deterministic acceptance meaning.

### D-05: Run control and private state

Sprint 35 run control owns operation acceptance, aliases, attempts, leases, fencing generations, cancellation, events, terminal proposals, and deterministic terminal precedence. QA state records the corresponding run ID, operational attempt ID, and fencing generation. Every private write checks the active fence. Stale writers receive an explicit conflict and cannot replace the current pointer or canonical summary.

New mutable QA records use private schema version 2 and a new namespace. The logical layout contains immutable run records addressed by run ID, an atomically replaced current pointer with the run-record digest, and bounded retained history. Existing v1 records remain readable for status and recovery compatibility. New work never mutates a v1 record or writes a mixed-version object. An unsupported or corrupt version fails with a stable migration or integrity code.

The commit order is run record, digest validation, atomic current-pointer replacement, then derived canonical summary and flow projection. Recovery reads run-control truth first, validates pointer and record digests, reconciles interrupted ownership, applies retention only to non-current terminal records, and reports freshness reasons. The public query returns the current run and bounded summaries, not private persistence records or unbounded history.

Grounding comes from Sprint 35, Sprint 36, the selected `durable_operations.go` and QA state excerpts, PERSIST-ATOMIC-001, PERSIST-MIG-001, PERSIST-DERIVED-001, WF-STATE-001, and ERR-DATA-001. The accepted debt is a compatibility reader for v1 until an explicit later removal decision. The rejected alternative is adding a second QA lifecycle to flow state, which would conflict with run-control authority.

### D-06: Canonical smoke-suite integration

The app operation request gains a closed optional suite discriminator. Empty means normal mapped QA. `smoke` means execute Sprint 37's required evidence checks through the QA service and then use the existing smoke service for canonical discovery, selection, execution, evidence validation, issue linking, artifact writing, and overall assessment. Unknown suite values fail before side effects.

CLI exposes the selection on existing QA dry-run and start commands. The intended shape is `ultraplan qa start <project> <sprint> --suite smoke`, with the same confirmation, run-control acceptance, progress, cancellation, and status mechanics as other runtime-backed QA work. `qa resume --suite smoke` is rejected because the canonical containing smoke suite has no Sprint 37 checkpoint contract. Operators restart it; the current status explains that action.

The QA suite does not rename `review` or `smoke`, replace `review.md` or `smoke.md`, redefine `PlanningStage`, or store raw harness evidence under private QA state. A current acceptable Conformance Review remains the smoke gate unless the existing explicit diagnostic override path applies. Diagnostic smoke cannot overturn a failed or blocked review.

Grounding comes from the smoke report, selected `smoke.go`, `smoke_author.go`, and `verify.go` excerpts, TEST-SMOKE-001/002, WF-BOUNDARY-001, and CLI-SHAPE-001. The accepted trade-off is an additional selector on a stable command. The rejected alternative is a new `qa-smoke` executor, which would duplicate discovery, evidence validation, and acceptance truth.

### D-07: App, CLI, TUI, and browser contracts

The app remains the only adapter-independent projection. Existing app JSON stays at schema version 1. Optional additions include suite identity, evidence-completeness counts, analyzer and evaluator summaries, permission and isolation summaries, smoke-route status, and focused detail references. Existing fields and meanings do not change.

Focused queries return one bounded view at a time: current map/status, one shard, one theory or check, synthesis, or smoke summary. They validate that requested IDs belong to the current map. They do not return private records, prompts containing raw secrets, unbounded histories, or arbitrary filesystem content.

CLI sends final human or JSON results to stdout and progress, warnings, diagnostics, and errors to stderr. A completed run with a failed or blocked QA verdict returns non-zero. A successful read-only status query returns zero even when it reports a prior failed verdict, because the query itself succeeded and the machine-readable verdict remains explicit. Cancellation, timeout, validation, permission, provider, persistence, and internal errors retain distinct non-zero classifications.

TUI and browser render the same bounded app facts with separate renderer code. Neither adapter computes verdicts. The TUI shows map fingerprint, current phase, shard/check progress, analyzer and evaluator breakdown, evidence completeness, permission denials, blocker, and next action. It offers start, cancel, recover, and focused inspection; resume appears only when the app reports a valid resumable checkpoint.

The browser retains server-rendered forms and status pages without JavaScript. Optional progress streaming improves the view but is not authoritative. After reconnect, sequence gap, completion, or visibility return, the page fetches focused durable status. Hostile evidence text is escaped and rendered as inert text. No generated text is inserted as raw HTML or executable link content.

Grounding comes from the API and frontend area reasoning, CLI-IO-001, CLI-JSON-001, CLI-EXIT-001, CLI-LIFE-001, TEST-CONTRACT-001, and ARCH-ENTRY-001. The accepted cost is additive DTO and renderer work. The rejected alternative is adapter-specific filesystem reads, which would duplicate policy and expose private state.

### D-08 and D-09: Errors, observability, configuration, and bounds

All failures use stable machine codes and categories. At minimum the design distinguishes validation, path containment, permission capability, forbidden write attempt, stale target, stale fence, missing evidence, malformed evidence, assertion failure, provider dependency, retry exhaustion, timeout, cancellation, cleanup, persistence, migration, and internal failure. Boundary translation preserves causes plus run, attempt, shard, check, call, prompt, provider, and correlation identities after redaction.

Progress events are advisory and bounded. Durable status contains the authoritative phase, completed and total counts, current shard/check, retained evidence references, blocker, cancellation state, and next action. Dropped events trigger a focused refresh. Prompt identity and checksum, model/provider, retry count, duration, and cost counters are retained when relevant. Raw prompts, repository content, credentials, and unrestricted local paths are not logged.

The implementation will define typed limits for target files and bytes, per-file size, changed paths, shards, checks, direct-test output, LLM context, analyzer and evaluator calls, retries, call timeout, process timeout, cleanup grace, retained runs, response rows, and event buffering. Invalid or unsafe combinations fail before scans or provider calls. Sprint 37 must not add arbitrary environment reads or a silent model fallback. Safe effective settings may be shown through existing diagnostics.

Grounding comes from ERR-CORE-001, ERR-CODE-001, ERR-TRANS-001, ERR-STARTUP-001, OBS-CORE-001, OBS-CORR-001, OBS-TASK-001, CFG-TYPE-001, CFG-START-001, PERF-BUDGET-001, PERF-BOUND-001, and PERF-COST-001. No external study source changes these contract-driven decisions. The accepted cost is a larger explicit state machine and error matrix. The rejected alternative is generic `operation_failed` reporting, which cannot support retry, diagnosis, or acceptance.

### D-10: Acceptance and documentation

Acceptance requires deterministic unit and integration tests, race coverage where concurrent provider calls or cancellation are exercised, adapter compatibility tests, Architecture Review, Sprint Review, and Deep Smoke Sprint evidence. The deep smoke must use the real runtime adapter and must run every required Sprint 37 check against deterministic malicious fixtures, including a fixture that attempts protected writes and one that embeds prompt-like instructions in repository content.

Documentation updates will describe command shape, suite selection, read-only and disposable-write boundaries, configuration and limits, result states, exit behavior, cancellation and recovery, state compatibility, troubleshooting, and the limits of LLM majority. Generated help or JSON snapshots must be reproducible.

Grounding comes from TEST-UNIT-001, TEST-INT-001, TEST-FAIL-001, TEST-E2E-001, DOC-ARCH-001, DOC-PUBLIC-001, DOC-OPS-001, and the three required review protocols. The trade-off is slower acceptance. It is necessary because fake-only tests cannot prove adapter permission behavior or canonical smoke integration.

## 6. Contracts and Standards Traceability

| Contract / Standard | Applicable Decisions | How Compliance Will Be Proven |
|---|---|---|
| ARCH-CORE-001, ARCH-CORE-002, ARCH-LAYER-002, ARCH-ENTRY-001, ARCH-SHARED-001 | D-01, D-06, D-07 | Import graph review, narrow consumer-owned process port, thin adapter tests, no product imports from platform code |
| ERR-CORE-001, ERR-SHAPE-001, ERR-CODE-001, ERR-TRANS-001, ERR-STARTUP-001, ERR-DATA-001 | D-02, D-04, D-05, D-08 | Failure matrix verifies stable codes, cause preservation, fail-fast preflight, distinct missing/corrupt evidence, and no empty success |
| CFG-SOURCE-001, CFG-TYPE-001, CFG-COMPAT-001, CFG-START-001, CFG-PUBLIC-001 | D-06, D-08, D-09 | Typed validation and precedence tests; invalid limits and suites fail before side effects; diagnostics expose only safe fields |
| OBS-CORE-001, OBS-CORR-001, OBS-DIAG-001, OBS-TASK-001, OBS-PII-001 | D-03, D-04, D-05, D-08 | Run/attempt/shard/check/call correlation tests, durable status, dropped-event refresh, prompt identity, deterministic redaction |
| SEC-INPUT-001, SEC-FILES-001, SEC-INJECT-001, SEC-SECRETS-001, SEC-DESER-001 | D-02, D-03, D-06, D-08 | Traversal, symlink, absolute-path, argv, hostile-content, parser, and secret-leak tests |
| TEST-SEAM-001, TEST-UNIT-001, TEST-INT-001, TEST-SMOKE-001/002, TEST-FAIL-001, TEST-DET-001, TEST-CONTRACT-001, TEST-E2E-001 | D-01 through D-10 | Replaceable collaborators, deterministic unit/integration suites, canonical real-adapter smoke, negative cases, adapter snapshots, full-path scenario |
| DOC-OWNER-001, DOC-ARCH-001, DOC-PUBLIC-001, DOC-OPS-001, DOC-AGENT-001 | D-06, D-07, D-10 | Updated authoritative architecture, command, operator, and agent-facing documentation with no conflicting duplicate |
| CLI-SHAPE-001, CLI-HELP-001, CLI-IO-001, CLI-EXIT-001, CLI-JSON-001, CLI-LIFE-001, CLI-SAFE-001, CLI-NONINT-001 | D-06, D-07, D-08 | Help and JSON snapshots, stdout/stderr tests, stable exit classifications, confirmation/dry-run, cancellation, no hidden prompts |
| LLM-BOUNDARY-001, LLM-TOOL-001, LLM-IO-001, LLM-LIFECYCLE-001, LLM-RETRY-001, LLM-PROMPT-001, LLM-SAFETY-001 | D-01, D-03, D-04, D-08 | Restricted tool contract, local schema validation, fresh sessions, bounded retry, prompt identity, hostile-content containment |
| EVAL-SCOPE-001, EVAL-REG-001, EVAL-HUMAN-001, EVAL-MODEL-001, EVAL-DATA-001, EVAL-COST-001 | D-04, D-09, D-10 | Representative eval set, preserved disagreement and evidence, Conformance Review, provider/model/cost records, classified fixtures |
| WF-SCOPE-001, WF-BOUNDARY-001, WF-STATE-001, WF-RETRY-001, WF-IDEMPOTENCY-001, WF-COMP-001 | D-01, D-04, D-05, D-06 | Existing app-owned operation path, explicit phases, idempotent immutable records, bounded retries, cancellation, reconciliation; no new engine |
| PERF-BUDGET-001, PERF-BOUND-001, PERF-CONC-001, PERF-COST-001 | D-03, D-04, D-09 | Fixed caps, over-budget rejection, complete failure aggregation, cancellation, representative runtime and cost evidence |
| PERSIST-READ-001, PERSIST-SCHEMA-001, PERSIST-ATOMIC-001, PERSIST-MIG-001, PERSIST-INTEGRITY-001, PERSIST-DERIVED-001 | D-02, D-05 | Bounded reads, documented v2 ownership, immutable run then atomic pointer, v1 compatibility, digest and recovery tests |
| Project architecture, PRD, and TRD | D-01 through D-10 | Architecture Review and traceability from R1-R53 to implementation tasks and tests |
| Architecture Review Protocol | D-01, D-03, D-05, D-06 | Review ownership, inward dependencies, persistence, isolation, and canonical smoke path before acceptance |
| Sprint Review Protocol | D-02, D-04, D-07, D-08, D-09 | Review requirements, compatibility, failure behavior, evidence completeness, documentation, and regression coverage |
| Deep Smoke Sprint Protocol | D-03, D-04, D-06, D-10 | Real runtime adapter, every required check, hostile fixtures, denied writes, immutable target, external harness evidence identity |

## 7. Trade-Offs, Debt, and Future Considerations

### Accepted Trade-Offs

- Sequential copy-backed workers are slower than shared parallel workers. Sprint 37 accepts that cost to make isolation, event attribution, cancellation, and cleanup deterministic.
- Three semantic analyzer calls and three failed-shard evaluator calls spend more time and provider budget than a single judge. The fixed odd quorum removes tie handling and preserves disagreement, but it does not guarantee independent truth.
- A separate private v2 namespace and read-only v1 decoder increase persistence code. They avoid mutating shipped records in place and make unsupported versions explicit.
- Current-pointer-only public projections limit historical exploration. They keep CLI, TUI, and browser responses bounded and prevent persistence details from becoming public contracts.
- `Suite: "smoke"` adds a branch to an existing operation family. It is preferable to another command and executor because the canonical smoke route remains singular.
- Native permission controls plus runtime policy and snapshots overlap. The overlap is intentional. Each mechanism catches a different class of failure.

### Known Technical Debt

- Fresh LLM sessions may still share provider, model, prompt, or training bias. Reports must not call them statistically independent.
- Copying a bounded target is an application isolation boundary, not a hostile multi-tenant sandbox. Stronger OS sandboxing remains a future security project if the threat model expands.
- The v1 compatibility reader remains until a later governed removal decision defines retention and migration completion.
- Worker concurrency is deliberately conservative. Performance data may justify bounded parallel workers later, but only with new race and isolation evidence.
- Requirement-to-check definitions remain product-authored. Sprint 37 does not add automated retrieval or semantic source identity.
- Browser history and raw evidence exploration remain intentionally narrow. Operators use focused queries and external canonical evidence roots.

### Rejected Alternatives

- **Shared writable worker tree:** rejected because one worker could contaminate another and target immutability would be hard to attribute.
- **Prompt-only read-only instructions:** rejected because they do not constrain filesystem or tool capabilities.
- **Original target mounted read-only into children:** rejected for Sprint 37 because it still discloses the original path and increases the chance of adapter or mount-policy mistakes.
- **Parallel copy-backed workers by default:** rejected until representative cost and race evidence justifies the extra scheduler complexity.
- **One LLM judge or model confidence threshold:** rejected because malformed or biased output would have too much authority and confidence is not a stable acceptance contract.
- **Evaluator override of security or freshness failures:** rejected because opinion cannot repair missing coverage, target mutation, stale inputs, or a broken permission boundary.
- **In-place private state migration:** rejected because interrupted writes could leave mixed semantics and break Sprint 36 recovery.
- **A new QA smoke command and executor:** rejected because it would fork canonical smoke selection, evidence validation, and artifact ownership.
- **Adapter-specific QA logic:** rejected because CLI, TUI, and browser verdicts would drift and private persistence would become a public dependency.
- **Treating provider failure as assertion failure:** rejected because it blames the implementation without evidence and can hide retryable infrastructure faults.

### Future Considerations

- Consider bounded parallel copy workers only after measured repository-size, provider-cost, cancellation, race, and isolation results support it.
- Consider stronger OS-level sandboxing if UltraPlan begins executing untrusted third-party binaries or hosted multi-tenant QA.
- Remove the v1 reader only through a versioned compatibility decision with retained-history evidence.
- Sprint 38 may consume failed evidence for bounded repair, but it must not weaken Sprint 37's immutable evidence, authority, or non-pass rules.
- A later reporting sprint may add paginated run history and richer evidence navigation through app-owned queries, not direct persistence access.
- Model, prompt, or evaluator-policy changes require a version bump, representative regression set, cost record, and review before becoming current.

## 8. Implementation Constraints

- Implement product decisions in `internal/sprint`; implement only generic copy/process mechanics in `internal/platform/process`.
- Expose consequential filesystem, process, runtime, clock, and persistence behavior through narrow consumer-owned seams. Do not patch private fields in tests.
- Keep CLI, TUI, and browser adapters free of filesystem reads, provider calls, persistence decoding, verdict calculation, and smoke selection logic.
- Preserve existing QA command names and meanings. Add only the closed optional smoke-suite selector and focused evidence queries.
- Preserve app JSON schema version 1 and all existing fields. New fields must be optional, bounded, and compatibility-tested.
- Write new mutable private QA records only as schema v2 in a new namespace. Keep v1 read-only. Never silently coerce unknown versions.
- Commit immutable run records before atomically replacing the current pointer. Validate digests during reads and recovery.
- Require a valid Sprint 35 writer fence for every private state mutation. Stale writers must not publish summaries or pointers.
- Build child requests without the original target path. Use copy-relative paths and an opaque target identity only.
- Validate root containment, traversal, absolute paths, and symlinks before copying or launching processes.
- Strip child environment to an allowlist. Use structured argv and never shell interpolation for user-influenced values.
- Require runtime permission capability for checks that depend on tool restrictions. Unsupported capability is blocked before provider work.
- Keep repository text, rubric instructions, and model outputs in separate typed fields. Treat repository content as untrusted data.
- Validate every model result locally against a versioned schema. Preserve each call identity, retry, error, and evidence digest.
- Use exactly three LLM analyzers for semantic checks and exactly three fresh evaluators for failed-shard adjudication. Do not make these counts silently configurable through environment variables.
- Do not allow adjudication to override target mutation, governed-input mutation, stale fingerprints, permission violations, missing coverage, or malformed evidence.
- Schedule copy-backed workers sequentially. Any provider-call concurrency must remain within a check's fixed cap and aggregate every result.
- Propagate cancellation through run control, worker processes, provider calls, cleanup, state projection, and adapter progress.
- Define typed practical caps for files, bytes, outputs, checks, calls, retries, timeouts, cleanup, history, events, and response rows. Reject over-budget work before side effects.
- Route smoke-suite work through existing smoke discovery, selection, execution, validation, artifact, and assessment code paths.
- Do not rename or replace `review`, `review.md`, `smoke`, `smoke.md`, or `PlanningStage`.
- Do not move external raw smoke evidence into product-owned QA persistence.
- Do not implement repair, patches, Git mutation, issue-tracker integration, hosted execution, alternate storage, or a retrieval subsystem.
- Escape browser evidence as inert text. Preserve server-rendered no-JavaScript status and action flows.
- Update authoritative command, architecture, operator, and compatibility documentation with the implementation.

## 9. Expected Evidence and Validation Strategy

### Tests

- Unit tests for deterministic map identity, sorted inputs, requirement ownership, check-catalog validation, changed-path coverage, affected tests, contradictions, stale fingerprints, and completeness.
- Unit tests for each check kind and for pass, fail, and blocked aggregation.
- Unit tests for one-of-one deterministic majority, two-of-three semantic majority, disagreement, malformed output, missing call, timeout, cancellation, provider failure, retry exhaustion, and inadmissible adjudication.
- Unit tests proving security, freshness, coverage, and mutation failures cannot be evaluator-overridden.
- Integration tests for fresh copy creation, allowed copy writes, protected-root denial, original-path omission, environment allowlisting, structured argv, symlink and traversal rejection, bounded output, timeout, cancellation, and cleanup.
- Integration tests comparing target and governed-input identities before and after successful, failed, cancelled, timed-out, and malicious runs.
- Run-control tests for acceptance aliases, owner claims, fencing generations, stale writers, cancellation races, crash recovery, terminal precedence, and summary projection.
- Persistence tests for v1 read compatibility, v2 writes, unsupported-version rejection, run-record digest validation, atomic current-pointer replacement, interrupted commits, retention, and reconciliation.
- App tests for additive JSON v1 fields, focused query ownership, bounded response sizes, suite validation, no-resume smoke behavior, and canonical smoke routing.
- CLI tests for legacy command compatibility, help, dry-run and confirmation, stdout/stderr separation, JSON shape, progress, stable non-zero classifications, cancellation, and status-query exit behavior.
- TUI tests for one bounded projection, start/cancel/recover controls, conditional resume, evidence drill-down, permission denials, blockers, and authoritative refresh.
- Browser tests for server-rendered no-JavaScript forms and status, escaped hostile content, optional progress streaming, reconnect or sequence-gap refresh, CSRF and route validation, and focused bounded queries.
- Canonical smoke parity tests proving legacy smoke and `Suite: "smoke"` reach the same discovery, selection, evidence validation, artifact, and assessment path.
- Race tests for cancellation, provider-call aggregation, event delivery, run-control fencing, and pointer publication where concurrency exists.

### Logs / Metrics

- Structured events carry correlation ID, run ID, operational attempt ID, fencing generation, shard ID, check ID, call ID, phase, stable code, component, and safe timestamp.
- LLM events add role, session ID, prompt ID/version/checksum, provider, model, retry count, duration, token or cost counters when available, and fallback identity when applicable.
- Process events add copy-workspace identity, safe relative command identity, timeout, exit status or signal, truncation, permission outcome, write-attempt attribution, and cleanup outcome. They must not contain the original target path.
- Durable status records expected/completed shards, checks, observations, analyzer calls, evaluator calls, evidence references, freshness reasons, target identity, permission summary, blocker, cancellation state, and next action.
- Metrics or review records capture repository file/byte counts, changed paths, checks, provider calls, retries, duration, retained runs, dropped events, and cleanup failures within bounded cardinality.
- Redaction tests inspect stdout, stderr, logs, events, errors, JSON, HTML, private summaries, and generated evidence for credentials, raw prompts, hostile content leakage, and original-target paths.

### Manual / Review Checks

- Run the Architecture Review Protocol against package ownership, dependency direction, consumer-owned ports, state v2, run-control authority, isolation, and canonical smoke routing.
- Run the Sprint Review Protocol against all R1-R53 mappings, non-goals, compatibility, negative paths, documentation, and the complete evidence manifest.
- Run the Deep Smoke Sprint Protocol with the real runtime adapter and every required Sprint 37 check.
- Include deterministic malicious fixtures that attempt writes to the original target and governed project inputs, use symlink/traversal paths, emit prompt-like repository content, return malformed model output, time out, cancel, and fail one evaluator call.
- Prove that write attempts are denied and attributed, original target and governed-input identities remain unchanged, each analyzer/evaluator has a fresh session identity, incomplete calls are non-pass, and all required evidence is present.
- Record the exact real-adapter smoke run, provider/model/prompt versions, selected containing suite, external evidence identities and digests, duration/cost, review fingerprint, and final assessment.
- Inspect CLI, TUI, and browser output manually for comprehensible progress, evidence completeness, disagreement, permission denial, blockers, cancellation, and next actions without exposing private records.
- Confirm docs and help explain that copy workspaces are disposable, the original target is read-only, majority is not proof, smoke-suite resume is unsupported, and failed or blocked QA does not apply repairs.

## 10. Handoff to Planning

`plan.md` should decompose these decisions in dependency order: domain schemas and policy; process isolation port and implementation; evidence-map and typed checks; analyzer/evaluator orchestration; run-control and v2 persistence; smoke-suite routing; app DTOs and focused queries; CLI, TUI, and browser rendering; observability, documentation, tests, review, and deep smoke. Each task must cite the applicable decision IDs, requirements, contract clauses, files, and expected evidence.

The plan must not reopen these decisions:

- QA policy belongs to `internal/sprint`; generic isolation mechanics belong to `internal/platform/process`.
- Workers use fresh copy workspaces and sequential scheduling. Children never receive the original target path.
- Native permissions, runtime policy, identity snapshots, and attributed events are all required.
- Fact, negative, behavioral, semantic, and adversarial checks use explicit schemas and direct evidence.
- Semantic/adversarial analysis uses three calls where LLM analysis is required; failed-shard adjudication uses exactly three fresh evaluators and a two-of-three pass majority after all calls complete.
- Incomplete, malformed, cancelled, timed-out, provider-failed, permission-denied, stale, or missing evidence is non-pass.
- Run control remains authoritative. New private QA state uses v2 in a new namespace with read-only v1 compatibility and atomic current-pointer publication.
- App JSON remains v1 and grows additively through bounded projections.
- `Suite: "smoke"` uses the existing QA operation family, rejects resume, and routes through canonical smoke behavior and artifacts.
- CLI, TUI, and browser adapters render app facts and do not own QA policy or private persistence.
- Production repair, Git mutation, issue tracking, hosted operation, alternate persistence, and retrieval/content identity remain out of scope.
- Acceptance requires Architecture Review, Sprint Review, and a hostile real-runtime Deep Smoke Sprint that executes every required check.

## Decisions

The final decisions are D-01 through D-10 in Sections 4 and 5. They fix ownership, immutable evidence mapping, copy-workspace isolation, analyzer and evaluator quorum, run-control authority, private state v2, canonical smoke routing, additive adapter contracts, bounded operations, and acceptance evidence. `plan.md` must implement those decisions without substituting a second lifecycle, smoke executor, or adapter-owned policy.

## Assumptions And Risks

The governing assumptions and risk register are in Section 2. Planning must preserve their mitigations, especially runtime permission preflight, original-path omission, path and symlink containment, target identity checks, stale fencing, local LLM-output validation, atomic pointer publication, authoritative refresh after event gaps, and fail-closed canonical smoke reuse. No architecture risk remains unanswered; implementation evidence must close each listed risk before acceptance.

## Implementation Constraints

Section 8 is normative. In particular, product policy stays in `internal/sprint`, generic isolation stays in `internal/platform/process`, workers run sequentially in fresh copies, children receive no original target path, private writes require the active Sprint 35 fence, state v2 does not mutate v1 records, app JSON stays v1, and `Suite: "smoke"` uses the canonical smoke path without resume.

## Expected Evidence

Section 9 defines the required unit, integration, race, compatibility, adapter, review, and deep-smoke evidence. Acceptance requires all three review protocols plus a real-runtime hostile-fixture run that executes every required check, proves protected writes were denied and attributed, proves protected identities stayed unchanged, records every analyzer and evaluator call, and leaves every incomplete or invalid result as non-pass.
