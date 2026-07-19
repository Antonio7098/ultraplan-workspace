# Architecture Reasoning: Automated Sprint Review

> **Inputs Used:** `projects/ultraplan-go/sprints/26-review-stage/sprint-index.md`, `projects/ultraplan-go/sprints/26-review-stage/technical-handbook.md`, `projects/ultraplan-go/sprints/26-review-stage/requirements.md`, `projects/ultraplan-go/docs/ARCHITECTURE.md`, `projects/ultraplan-go/docs/PRD.md`, `projects/ultraplan-go/docs/TRD.md`, `system/reasoning/architecture_reasoning_template.md`

This area covers Phase 3 review ownership, dependency direction, runtime and persistence boundaries, reviewer concurrency, and the shared CLI/TUI application surface. The feature receives a validated sprint and target scope, freezes the governed inputs, gathers structured contract and handbook assessments, performs deterministic product checks, computes a product-owned verdict, and persists the current `review.md` and review flow state. Smoke execution, automatic remediation, issue tracking, Git mutation, and a new runtime abstraction are outside this area.

## Area Decisions

### Ownership and dependency direction

- `internal/sprint` owns the complete review workflow: readiness and selected-context resolution, immutable review-manifest construction, prompt construction, reviewer orchestration, structured result validation, deterministic conformance checks, finding normalization, verdict synthesis, `review.md` validation, and review/flow persistence. These operations transform sprint-owned state and therefore belong together rather than in global reviewer, validation, or workflow packages.
- `internal/sprint` may consume the existing `project` catalog API, workspace path policy, configuration values, clock, generic runtime, and sprint store. It must not import `study`, direct OpenCode packages, TUI packages, or CLI handlers. `internal/platform/runtime` remains unaware of projects, sprints, contracts, reviewers, findings, and verdicts.
- Runtime supervision continues through the existing agentwrap-backed platform boundary. Review adds a review-specific request/result translation inside `internal/sprint`; it does not create a competing provider interface or duplicate agentwrap retry, cancellation, permission, structured-event, or OpenCode process behavior.
- Embedded prompt and output-template lookup remains a workspace mechanism, while prompt meaning, required reviewer schema, and generated review content remain sprint behavior. This preserves the product/platform split without moving review semantics into the asset registry.

### Primary workflow shape

The primary unit is one explicit sprint review operation on the existing sprint service, supported by focused review files rather than a new subsystem or package tree:

1. Resolve project, sprint, target implementation, effective review configuration, selected contracts, and selected review protocols through existing catalog and path APIs.
2. Validate prerequisites through execute and reject duplicate, unknown, unreadable, missing, or escaping selected paths before starting runtime work.
3. Build an immutable invocation manifest containing workspace-relative governed inputs, selected scope, implementation identity, changed-path scope, model source, permission policy, and deterministic fingerprint.
4. In dry-run or prompt-preview mode, return this manifest and rendered requests without creating the runtime or writing review state.
5. For execution, issue one independent structured request for each selected contract and one for the technical handbook, under a per-invocation concurrency bound. Selected protocols constrain review and deterministic checks but do not create additional reviewer units unless a future requirement explicitly says so.
6. Collect every possible reviewer outcome. An individual malformed, missing, or failed result is recorded as a failed coverage unit and does not cancel healthy siblings. User cancellation, deadline expiry, or a stage-wide infrastructure failure stops scheduling, cancels active runs, and awaits bounded cleanup.
7. Validate returned schemas and contained citations before findings enter the product result. Run deterministic decision, plan-task, verification-command, coverage, citation, deviation, and missing-evidence checks independently of model prose.
8. Sort normalized results and findings by stable product keys, compute the verdict in product code, render and validate complete Markdown in staging, then atomically replace `review.md` and atomically update review flow state.

The current module shape still fits. The feature needs focused `review.go` and `review_validation.go` behavior plus extensions to existing domain, service, store, app, command, and TUI files; it does not justify splitting `internal/sprint` or introducing a generic workflow engine.

### State and consistency boundaries

- Durable source state is the governed sprint artifacts, selected catalog entries, execute evidence, target implementation identity, and changed-path scope. The invocation fingerprint derives from their canonical identities and contents and is computed, not independently authored.
- Ephemeral state consists of the frozen manifest, queued/active reviewer work, progress events, validated reviewer results, deterministic check results, and staged Markdown. Workers receive immutable per-reviewer inputs and return owned results; they do not mutate a shared result slice or render the final artifact.
- Durable output state is only the canonical sprint-root `review.md` and the review fields in `flow-state.json`. No separate reviewer cache, alternate TUI state file, or review database is introduced. Recovery reruns review from a newly frozen manifest rather than trusting partial model output from a prior fingerprint.
- The store must stage and fully validate `review.md` before canonical replacement. Atomic replacement preserves the previous file if staging, validation, flush, close, or rename fails. Flow state is also replaced atomically and records execution status separately from verdict, because a completed review may truthfully have a `fail` verdict.
- Cross-file replacement is not treated as a filesystem transaction. The canonical artifact is replaced only after it is complete and valid; the subsequent flow-state update records its exact fingerprint and path. Status must reconcile artifact validity and fingerprint rather than claim current success from either file alone. Thus a flow-state write failure exposes a persistence failure without treating a stale state record as proof, while the artifact on disk is always a complete version rather than a partial file.
- Re-running review may replace only `review.md` and the review portion of flow state. Governed inputs, product source/tests, workspace configuration, other sprint artifacts, and Git state are read-only.

### Shared CLI and TUI boundary

- `internal/app` exposes typed operations for readiness/status, dry-run, prompt preview, start, progress subscription or callback, cancellation, and final result. Request and result DTOs carry scope, fingerprint, model source, reviewer counts, execution status, verdict, staleness, findings, diagnostics, artifact path, and next action without terminal formatting.
- CLI handlers only parse flags, invoke those operations, map classified outcomes to exit codes, and render text or JSON. JSON stdout contains only the final machine envelope; progress and diagnostics use an injected non-data channel.
- TUI commands call the same app operations and translate typed progress/results into model messages. The TUI may own confirmation, navigation, filtering, and presentation, but not manifest construction, reviewer scheduling, validation, verdicts, persistence, or recovery policy.
- Progress may be emitted as reviewer lifecycle facts while work runs. Unvalidated findings are not exposed as final findings as they arrive; finding navigation and final text/JSON use only validated, deterministically sorted aggregate results.

### Dependencies and abstractions

- Keep manual construction in the existing app composition root. Runtime construction remains lazy so help, status, validation, prompt preview, and dry-run do not pay runtime startup cost or fail on credentials they do not need.
- Reuse the existing generic runtime boundary and agentwrap request model. The review service supplies prompt, work directory, model, timeout, metadata, required capabilities, structured validation, and enforceable read-only permission policy.
- Use narrow sprint-owned seams only at volatile boundaries already needing deterministic tests: runtime execution, review/store persistence, clock, and progress delivery. Keep manifest construction, sorting, severity rules, verdict synthesis, and Markdown validation concrete pure functions inside `internal/sprint`.
- Context carries cancellation and correlation metadata only. Services, stores, configuration, and policy are constructor/request dependencies, not context values or process-wide globals.
- Do not add plugin registries, generic `Manager`/`Processor` types, a protocol-reviewer strategy hierarchy, or interfaces mirroring one concrete helper. The selected contracts are data-driven inputs to one stable review algorithm, not separate implementations.

### Failure and completion semantics

- Preflight failures prevent a run and preserve current review evidence. Missing runtime/credentials/environment after valid scope resolution is `blocked`; invalid governed input is a validation/preflight failure; runtime, timeout, cancellation, malformed output, incomplete coverage, citation failure, and persistence failure remain distinguishable typed outcomes.
- A reviewer task failure is data in a completed collection attempt and makes a passing verdict impossible. It is not silently reduced to a warning and does not erase successful sibling evidence.
- Execution status and conformance verdict remain separate throughout domain, app DTO, JSON, CLI, TUI, and flow state. Runtime success cannot imply review completion, and review completion can result in `pass`, `pass_with_findings`, `fail`, or `blocked` only.
- Review completes only after all required coverage units have terminal outcomes, deterministic checks have run, the verdict has been computed, staged Markdown has passed review validation, canonical replacement has succeeded, and flow state has been updated. Cancellation is never reported as a completed pass.

## Trade-Offs

### Chosen: collect independent reviewer outcomes with bounded fan-out

Collect-all execution spends more time and model budget after one reviewer fails, but it preserves actionable coverage and truthful per-contract diagnostics. Fail-fast `errgroup` behavior was rejected for ordinary reviewer failures because it can hide whether unaffected contracts conform. Stage cancellation still applies for user cancellation, deadline expiry, or infrastructure conditions that invalidate the whole invocation.

### Chosen: concurrent collection, deterministic post-processing

Parallel reviewers reduce wall-clock latency, but introduce ordering, quota, memory, and lifecycle complexity. The design accepts that complexity only at the external runtime boundary: scheduling is bounded, every started run is awaited, workers return isolated results, and final checks/rendering are sequential and stable. Fully sequential review was rejected because the units are independent and requirements explicitly call for bounded fan-out. Unbounded goroutines were rejected because they make resource use and cleanup unauditable.

### Chosen: extend `internal/sprint` before extracting abstractions

Keeping review near sprint state yields a larger package, but preserves product context and one-way dependencies. A new global review/workflow package was rejected because only sprint review owns these semantics and a generic engine would add lifecycle and extension obligations outside the sprint. A clean-architecture subpackage tree was also rejected until package size or dependency pressure demonstrates a concrete comprehension benefit.

### Chosen: typed app use cases over shared renderers

Typed DTOs require explicit mapping for CLI text, JSON, and TUI messages, but prevent terminal output from becoming an accidental integration protocol. TUI-to-CLI calls, stdout parsing, native provider payload interpretation, and alternate TUI persistence were rejected because they duplicate policy and undermine parity.

### Chosen: product-owned structured aggregation over model-authored review files

Having reviewers write sections directly could reduce product rendering code, but would permit nondeterministic ordering, unsafe paths, partial file mutation, and model-selected verdicts. Reviewers therefore return validated structured assessments; UltraPlan validates citations, applies deterministic checks and severity rules, and writes one artifact. Raw model prose is not a durable source of truth.

### Chosen: atomic per-file replacement with fingerprint reconciliation

Local filesystems do not provide a portable atomic transaction across `review.md` and `flow-state.json`. Introducing a database or transaction journal solely for two files was rejected as disproportionate. Instead, each file is independently staged and atomically replaced, only complete validated Markdown reaches the canonical path, flow state identifies the exact artifact fingerprint, and status treats mismatch or write failure as non-current evidence.

### Chosen: fresh rerun rather than partial reviewer-result persistence

Persisting each model result could reduce rerun cost, but would add a new durable schema, cache invalidation, secret/redaction obligations, and ambiguity when implementation inputs change. Sprint 26 persists terminal review facts and safe diagnostics only. A recovery run creates a new frozen manifest and reruns the selected coverage; bounded retries inside each request remain owned by agentwrap policy.

## Evidence

- Thin interface adapters over shared use cases are a high-confidence pattern in the handbook (`projects/ultraplan-go/sprints/26-review-stage/technical-handbook.md:31-32`). This supports placing behavior in `internal/sprint`, composition and typed operations in `internal/app`, and rendering alone in CLI/TUI adapters. The project architecture independently requires this exact dependency direction (`projects/ultraplan-go/docs/ARCHITECTURE.md:258-299`).
- Manual composition and lazy volatile dependencies keep ownership visible and simple commands fast (`projects/ultraplan-go/sprints/26-review-stage/technical-handbook.md:32-33`, `54-55`). This is the basis for explicit app construction and avoiding runtime initialization during help, status, validation, preview, and dry-run.
- Root cancellation, bounded fan-out, and awaited shutdown are supported by the state/context and concurrency evidence (`projects/ultraplan-go/sprints/26-review-stage/technical-handbook.md:20-21`, `36-37`, `61-62`). The sprint requirements make one request per selected contract plus one handbook request, bounded concurrency, cancellation, and complete failed/missing result accounting mandatory (`projects/ultraplan-go/sprints/26-review-stage/requirements.md:44-49`).
- The handbook identifies deterministic aggregation after concurrent work as necessary to avoid shared mutation and unstable output (`projects/ultraplan-go/sprints/26-review-stage/technical-handbook.md:37-38`, `94-100`). This supports isolated worker results, sequential validation/synthesis, stable sorting, and no concurrent artifact rendering.
- Wrapped causes plus semantic classification support boundary-specific exit, JSON, and user guidance (`projects/ultraplan-go/sprints/26-review-stage/technical-handbook.md:18`, `34-35`, `52-53`). The architecture therefore preserves preflight, blocked, runtime, timeout, cancellation, validation, reviewer coverage, persistence, and verdict as distinct facts rather than one error string.
- Injected streams and strict output separation make text, JSON, progress, and diagnostics independently testable (`projects/ultraplan-go/sprints/26-review-stage/technical-handbook.md:19`, `35`, `38-39`, `63-64`). This supports typed progress and separate machine-output rendering rather than direct worker writes.
- Mechanically enforced trust boundaries are stronger than prompt-only rules (`projects/ultraplan-go/sprints/26-review-stage/technical-handbook.md:25`, `40`, `67`). Requirements reinforce read-only agentwrap permissions, workspace/target containment, safe diagnostics, and prohibition of source, governed-input, and Git mutation (`projects/ultraplan-go/sprints/26-review-stage/requirements.md:47`, `53-54`, `79-89`).
- The handbook defines the frozen manifest and atomic artifact replacement as the review consistency boundary (`projects/ultraplan-go/sprints/26-review-stage/technical-handbook.md:96-100`). The TRD further requires fingerprinted freshness, separate execution status and verdict, complete reviewer coverage, contained line citations, deterministic verdicts, and preservation of the last complete artifact (`projects/ultraplan-go/docs/TRD.md:1706-1739`, `1834-1844`, `1929-1943`).
- Behavior-level tests with fake runtimes, command tests, synchronized I/O, and explicit fixtures are the selected evidence-backed strategy (`projects/ultraplan-go/sprints/26-review-stage/technical-handbook.md:39`, `68`, `86`). The architecture exposes runtime/store/clock/progress seams while keeping deterministic product rules pure, enabling the required success, severity, missing/malformed result, path escape, stale input, cancellation, redaction, atomic-write, race, and TUI tests (`projects/ultraplan-go/sprints/26-review-stage/requirements.md:61-62`).
- The product documents consistently define `review.md` as a product-owned stage after execute and before smoke, with one core shared by CLI and TUI (`projects/ultraplan-go/docs/PRD.md:46-57`, `128-135`, `837-850`; `projects/ultraplan-go/docs/TRD.md:115-138`, `1927-1971`). This confirms that review is an extension of the sprint module rather than a runtime, smoke, or interface concern.

## Risks

- **Cross-file commit interruption:** The process can stop after replacing valid `review.md` but before replacing flow state. Mitigation: embed and validate the exact input fingerprint in the artifact, make status reconcile both files, classify mismatch as non-current, and make rerun safe. No code path may infer current completion from file existence alone.
- **Reviewer cleanup leak:** Cancellation can return while an agentwrap run or event drain remains active. Mitigation: give every started run one owner, drain events, invoke cancellation, await `Wait`, bound cleanup, and exercise cancellation under `go test -race ./...`.
- **Collect-all cost and latency:** Continuing healthy reviewers after one failure increases runtime use. Mitigation: retain a configurable positive per-invocation concurrency bound, bounded agentwrap retries, cancellation, model/usage metadata, and clear progress; do not add speculative cross-run throttling until measurements show contention.
- **Permission-policy mismatch:** A runtime may not enforce every desired path rule. Mitigation: require the applicable agentwrap capabilities before launch, deny mutation/shell/network capabilities not needed by review, validate all input/output paths in product code, and block rather than downgrade required safety to prompt convention.
- **Input drift during review:** Files can change after fingerprinting. Mitigation: reviewers consume the frozen manifest identity and captured governed content; before persistence, re-check governed and target identity. A mismatch fails the invocation as stale and preserves the previous canonical review.
- **Diagnostic or citation disclosure:** Model output and runtime metadata can contain absolute paths, source content, credentials, or unsafe payloads. Mitigation: accept only the structured reviewer schema, normalize citations to contained workspace/target-relative paths, validate line ranges, redact safe diagnostics centrally, and omit raw native payloads from artifacts and normal JSON.
- **Package growth:** Adding orchestration, validation, rendering, and state extensions to `internal/sprint` can reduce readability. Mitigation: split by focused files and pure functions now; introduce a subpackage only if later concrete dependency or comprehension pressure outweighs the cost of crossing the module boundary.
- **CLI/TUI semantic drift:** Separate rendering paths can expose different status or next-action facts. Mitigation: make app DTOs the sole shared semantic result, test CLI/JSON/TUI against the same fixtures, and prohibit either interface from recomputing verdict, readiness, or staleness.
- **Assumption about implementation identity:** Git metadata may be absent. The architecture assumes the existing execute target reference plus explicit changed paths and a contained implementation fingerprint can provide a stable identity. If the implementation cannot produce a trustworthy identity, preflight must block; it must not silently weaken freshness checks.
- **Assumption about review resumption:** Sprint 26 does not require durable partial reviewer-result reuse. Recovery means a safe fresh invocation over a newly frozen manifest. If future measured cost requires partial reuse, it must be introduced as a versioned, fingerprint-keyed product state design rather than an unvalidated cache.
- **Open verification question:** Existing store and app APIs may already provide narrower atomic-write, progress, and cancellation seams than the named files imply. Implementation should reuse those concrete seams where they satisfy these decisions; this does not change ownership, durable state, or dependency direction and does not justify parallel abstractions.
