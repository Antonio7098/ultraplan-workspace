# Architecture Reasoning: Deep Smoke Harness Integration

> **Inputs Used:** `projects/ultraplan-go/sprints/27-deep-smoke/sprint-index.md`, `projects/ultraplan-go/sprints/27-deep-smoke/technical-handbook.md`, `projects/ultraplan-go/sprints/27-deep-smoke/requirements.md`, `projects/ultraplan-go/docs/ARCHITECTURE.md`, `projects/ultraplan-go/docs/PRD.md`, `projects/ultraplan-go/docs/TRD.md`, `system/reasoning/architecture_reasoning_template.md`

This area covers package ownership, dependency direction, state and process boundaries, and the shared application shape for review-gated external smoke. It does not design Sprint 28 `verify` orchestration, focused-rerun recovery, general issue management, automatic remediation, Git mutation, or a generic plugin system.

## Area Decisions

### Ownership and dependency direction

The existing module-driven architecture still fits, with a small boundary addition rather than a restructuring:

```text
cmd/ultraplan -> internal/app -> internal/sprint -> internal/project
                                      |          -> internal/workspace
                                      |          -> internal/platform/process
internal/tui --------------------------+

internal/platform/process -> no product modules
```

- `internal/sprint` owns the complete smoke workflow and its product meaning: current-review and freshness gating, catalog/manifest interpretation, harness discovery decoding, effective option resolution, scope selection, prerequisite classification, evidence validation, issue-reference import, deterministic verdict synthesis, `smoke.md` rendering/validation, and smoke flow/status state.
- `internal/platform/process` owns one narrow volatile boundary: run an executable with an argument vector, contained working directory, explicitly supplied environment, timeout, context cancellation, output limits, and process-tree cleanup; return exit, timing, bounded stream, truncation, and cleanup facts. Its request and result must not contain project, sprint, review, harness, suite, issue, evidence, or verdict types.
- `internal/project` continues to own parsing and validating catalog entries. It exposes the cataloged harness root and manifest reference as project data; it does not discover suites or execute the harness.
- `internal/workspace` supplies canonicalization and containment mechanics. Smoke-specific decisions about which returned paths are permitted remain in `internal/sprint`.
- `internal/app` exposes typed smoke readiness, preview, run, progress, cancellation, validation, status, and flow operations. CLI handlers only parse/render, while `internal/tui` only manages interaction and presentation. Neither interface may invoke the other or infer state from rendered text.
- `platform/runtime` and `agentwrap` remain unchanged. Review continues through agentwrap; smoke is supervised through `platform/process`, so the sprint does not create a second LLM runtime abstraction.

Focused files inside the existing packages are preferred: smoke domain/orchestration and validation files in `internal/sprint`, a runner in `internal/platform/process`, app use-case wiring in `internal/app`, and interface-specific rendering in existing CLI/TUI files. Subpackages, a generic workflow engine, and a plugin registry are not warranted.

### Product workflow

The smoke use case is one explicit pipeline with validation before every trust or mutation boundary:

1. Resolve workspace, project, sprint, catalog, effective smoke options, and allowed mutation roots.
2. Load and validate the current `review.md` plus its governed-input fingerprint. Permit normal execution only for a current passing/non-blocking review; require a confirmed and recorded diagnostic override otherwise.
3. Resolve the harness root and manifest from the project catalog. Canonicalize existing roots with symlink-aware resolution, reject escaping manifest/executable/cwd/evidence roots, and validate the supported manifest protocol and required capabilities before launch.
4. Execute the manifest's machine-readable discovery command through `platform/process`; strictly decode identity, protocol, levels, suites, tests, mappings, prerequisites, expected duration/cost, and evidence schema.
5. Select the narrowest sufficient scope deterministically: prefer a valid sprint-specific suite, then the smallest directly mapped suite set that covers the sprint, and use explicit tests only as diagnostic investigation unless discovery explicitly proves they satisfy required containing-suite coverage. Validate user level/suite/test overrides against discovery rather than passing arbitrary names through.
6. Classify unavailable required prerequisites as `blocked` and irrelevant scope as `not_applicable`; these terminal classifications do not launch a run. Otherwise show the effective scope, safe command representation, expected cost/duration, and mutation roots before guarded execution.
7. Run one harness process through the generic runner. UltraPlan owns no test-level worker pool; suite/test parallelism belongs to the harness. Propagate the operation context through process execution, progress projection, cancellation, bounded cleanup, and result collection.
8. Strictly decode the run response and validate run ID, selected coverage, result counts, duration, runtime/model metadata when supplied, external evidence identity and containment, and issue references. Exit code is evidence, not the verdict.
9. Synthesize exactly one of `pass`, `pass_with_open_issues`, `fail`, `blocked`, or `not_applicable` from validated facts. Missing evidence, required coverage, credentials, or cleanup certainty cannot produce a pass.
10. Render a candidate `smoke.md`, validate it in memory or from a same-directory temporary file, then atomically replace the current artifact. Update versioned flow state atomically only after the artifact commit succeeds, recording execution status separately from verdict. If the later flow-state write fails, status reports the reconciliation error and can reconstruct smoke state from the valid artifact; it must not delete or rewrite that artifact as compensation.

Malformed discovery/run data, timeout, cancellation, process failure, cleanup failure, or evidence validation failure do not replace the last valid `smoke.md`. Their execution state and safe diagnostics are recorded in flow/status state where possible. A truthful preflight `blocked` or `not_applicable` result is a completed smoke classification and may replace `smoke.md` after full artifact validation, as required by the sprint contract.

### Boundary contracts

The only new production interface should be the process runner because it isolates an external executable and enables deterministic fakes. It should be consumer-shaped and minimal:

```text
ProcessRequest:
  executable, argv, cwd, environment, timeout,
  stdout/stderr limits, progress sink

ProcessResult:
  started/finished time, duration, exit code/signal,
  bounded stdout/stderr, truncation facts, cleanup facts
```

The process API receives `context.Context` separately and treats the environment as a complete bounded map or list assembled by sprint/config code, not as permission to inherit the host environment implicitly. It executes directly without a shell. Progress delivery is bounded and non-blocking with respect to process draining and cleanup; terminal completion facts are retained even if an optional consumer is slow or disconnected.

Harness manifest, discovery, run, evidence, selection, progress, and verdict structs remain concrete sprint-domain types. They are versioned where persisted or exchanged, decoded strictly for required fields and supported major versions, and may tolerate unknown additive fields only when they do not alter required capability or safety semantics. There is one current harness implementation, so no harness interface, registry, factory, or generalized extension API is added.

Effective smoke options are resolved once before preflight and retain source metadata for safe diagnostics. Precedence is: manifest defaults, workspace smoke configuration, bounded environment-derived settings, flow/TUI request values, then explicit command overrides. TUI selections and CLI flags both populate the same typed app request; neither reads environment or manifest defaults independently. Validation occurs after the merge.

### State, mutation, and observability

Durable product state is limited to the current sprint `smoke.md` and versioned smoke fields in `flow-state.json`. Raw run JSON, unrestricted process streams, per-test artifacts, and issue records remain under the cataloged external harness. Normal smoke may mutate only those sprint outputs and manifest-declared harness run/issue roots; product source, tests, governed planning inputs, and Git state are excluded.

Ephemeral state includes the effective request, validated manifest/discovery response, selected scope, progress events, bounded stream excerpts, candidate summary, and evidence-validation result. Derived readiness, staleness, selected coverage, result totals, and next action are recomputed from governed inputs and durable artifacts rather than persisted in an alternate TUI model.

Typed progress uses product-level phases such as preflight, discovery, selection, running, validating evidence, writing artifact, completed, cancelled, and failed, with optional suite/test identity and counts. Harness-native event names stay behind the sprint decoder. User text, stable JSON results, and operational diagnostics are separate renderings; all omit or redact unrestricted environment values, credentials, unsafe argv values, raw model/user content, and unbounded process output.

### Verification shape

Core selection, gate, evidence, verdict, and artifact logic must be pure or testable with concrete fixtures. Process behavior is tested both with an in-memory recording runner and a real fake-harness child executable for timeout, signal, malformed output, output limits, and cleanup behavior. Command-level and TUI model/update tests call the same app use cases and compare semantic fields rather than parsing each other's output. Golden coverage is reserved for stable `smoke.md`, JSON envelopes, help, and selected presentation snapshots; verdict and containment rules use semantic assertions.

The real harness test remains opt-in. Its absence, missing credentials, or unavailable network is reported as a skipped/blocked gated test, never converted into passing smoke evidence. Required verification includes `go test ./...`, `go test -race ./...`, and `go build ./cmd/ultraplan` in addition to the fake-harness matrix.

## Trade-Offs

| Decision | Benefit | Accepted Cost | Rejected Alternative |
| --- | --- | --- | --- |
| Keep all smoke meaning in `internal/sprint` | Preserves module ownership and makes verdict/state rules traceable together. | The sprint package gains several focused files and collaborators. | A global verification/workflow package would fragment sprint state and prematurely generalize review/smoke. |
| Add a narrow `platform/process` runner | Isolates shell safety, cancellation, output bounds, cleanup, and test substitution. | Sprint orchestration must decode and validate all harness protocol data itself. | A harness-aware runner would reverse dependency direction and leak suite/evidence/verdict semantics into platform code. |
| Use one shared app use case for CLI and TUI | Guarantees behavioral parity and one persistence path. | Interface renderers must adapt typed events/results separately. | TUI-to-CLI subprocess calls or stdout parsing would create a second workflow and brittle coupling. |
| Strictly validate supported protocol and required capabilities | Prevents partial loads and false passes across independently released components. | A newer harness may be blocked until compatibility is added. | Warning-only parsing could silently ignore required safety or coverage fields. |
| Prefer containing-suite evidence | Produces explainable sprint-level confidence. | Runs may be broader and more expensive than a single mapped test. | Treating a focused test pass as sprint proof would violate required coverage semantics. |
| Supervise one harness process without internal test parallelism | Gives one lifecycle owner and avoids duplicate scheduling/cancellation policy. | UltraPlan cannot optimize individual test concurrency. | Product-owned per-test workers would duplicate harness responsibility and complicate evidence identity. |
| Bound and redact captured streams | Protects memory, JSON output, Markdown, and secrets. | Some diagnostics are available only through linked harness evidence. | Complete in-process capture risks unbounded memory and accidental disclosure. |
| Validate candidate output before atomic replacement | Preserves the last valid artifact through runtime and persistence failures. | Requires temporary-file handling and explicit reconciliation if flow-state update fails after artifact commit. | Writing incrementally or replacing before validation can corrupt the canonical summary. |
| Keep concrete sprint protocol types | Makes the flow direct and avoids speculative extensibility. | Supporting another unrelated harness shape will require deliberate adaptation later. | A plugin registry or generic harness interface is unearned with one cataloged contract. |

## Evidence

- The technical handbook's project-structure and command evidence shows that thin entrypoints and reusable operations preserve dependency direction, while large command bodies and TUI-only workflows reduce parity (`01-project-structure`, `02-command-architecture`; handbook lines 14-15, 31). This supports shared typed app use cases with CLI/TUI renderers at the edge.
- Constructor injection and fake-runner evidence supports creating the process dependency at the app composition root rather than using globals or context lookup (`03-dependency-injection`, `06-io-abstraction`; handbook lines 16, 19, 32). This is the basis for the sole new volatile-boundary interface.
- Helm's subprocess runtime and versioned metadata demonstrate the useful separation between host-owned protocol interpretation and generic subprocess execution, while K9s warning-only schema handling demonstrates the failure mode of permissive required-field validation (`12-extensibility`; handbook lines 25, 33, 47, 73-75).
- Explicit executable and argv construction is the selected security baseline; path validation and redaction remain independent trust-boundary controls (`13-security`; handbook lines 26, 57, 72). This supports direct execution, canonical containment checks, bounded environment forwarding, and safe display rather than shell interpolation.
- Signal-to-context propagation, explicit waits, and bounded cleanup are consistently favored over detached work (`07-state-context`, `08-concurrency`; handbook lines 20-21, 34-35, 59, 77-78). This supports one context-owned process lifecycle and no fire-and-forget progress or cleanup goroutines.
- Streaming and capped buffers reduce memory and responsiveness risk, while separate user, machine, and diagnostic channels prevent process output from corrupting stable JSON (`09-terminal-ux`, `10-logging-observability`, `14-performance`; handbook lines 22-23, 27, 35, 38, 49). This supports bounded stream capture and typed progress projection.
- Explicit precedence followed by validation of the merged configuration is evidenced by go-task, restic, and K9s (`04-configuration-management`; handbook lines 17, 37, 83). This supports resolving smoke defaults and overrides once in sprint/app code with source-aware diagnostics.
- Command integration tests, fake process runners, and explicit golden updates provide complementary behavioral and output coverage (`11-testing-strategy`; handbook lines 24, 39, 52, 75, 81-82). This supports the layered fake-runner, child-process fake harness, semantic, command, TUI, and golden test strategy.
- `projects/ultraplan-go/docs/ARCHITECTURE.md` assigns smoke semantics to `internal/sprint`, generic external execution to `internal/platform/process`, composition to `internal/app`, and rendering/interaction to `internal/tui` (lines 21-23, 31-38, 215-299, 343-364). The sprint requirements repeat these as mandatory constraints rather than optional patterns (requirements lines 56-67).
- `projects/ultraplan-go/docs/TRD.md` defines the smoke data flow, atomic artifact behavior, independent execution status/verdict, strict evidence links, review gate, and CLI/TUI parity (lines 1531-1611, 1704-1739, 1846-1855, 1927-1971). These requirements establish the validate-before-trust and shared-use-case pipeline.
- `projects/ultraplan-go/docs/PRD.md` states that runtime success is not product success, review precedes expensive proof, raw evidence stays external, and all Phase 3 surfaces share one product core (lines 46-57, 128-148, 664-693). This supports deriving verdicts from validated evidence rather than process exit or interface-specific state.

## Risks

- **Cross-platform descendant cleanup:** Linux and macOS process-tree behavior can differ, and cleanup certainty may be unavailable after timeout. The runner must expose cleanup facts; uncertain or failed cleanup is a classified non-passing execution result with diagnostics, not hidden success.
- **Symlink and time-of-check/time-of-use path changes:** Canonical containment checks can be invalidated if harness paths change during execution. Resolve and validate existing roots before launch, restrict cwd/evidence roots to manifest-declared locations, revalidate returned evidence paths after execution, and treat identity mismatch as failure.
- **Protocol evolution:** Strict major-version/capability checks can block a newer compatible harness. The accepted policy is fail closed for unsupported major versions or missing required capabilities, permit non-semantic additive fields, and include observed/supported versions in the corrective diagnostic.
- **Ambiguous sprint mappings:** Multiple equally narrow mapped suites can make selection non-obvious. Use deterministic ordering and select the smallest complete mapped suite set; if discovery cannot prove complete coverage, block and require an explicit suite/level rather than guessing or passing a narrow test.
- **Slow progress consumers:** TUI rendering or JSON event consumers could block stream draining and process cleanup. Progress delivery must be bounded/coalesced, optional intermediate events may be dropped with a count, and the terminal result must always be retained.
- **Artifact/state split commit:** Filesystem artifacts cannot form one transaction with external harness evidence and flow state. The chosen commit point is validated external evidence followed by atomic `smoke.md`, then atomic flow state. Status must detect and report/reconcile a valid artifact with stale state rather than rolling back evidence or the artifact.
- **Diagnostic override misuse:** A forced run after blocking review could be mistaken for approval. The request, artifact, JSON result, and flow state must preserve the original review verdict/fingerprint and explicit diagnostic-override fact; the override never changes review status.
- **Secret leakage through argv, environment, or evidence excerpts:** Names alone are not enough to determine sensitivity. Maintain explicit allowlists and redaction rules, render a separately sanitized argv/environment summary, omit raw streams from `smoke.md`, and test secret values embedded in URLs and argument values.
- **External evidence drift or deletion:** Linked run/issue files may change after summary creation. Record the strongest available identity metadata, preferring hashes and otherwise run ID plus stable path/schema/size/timestamp facts; validation must report stale or missing evidence instead of preserving a pass claim.
- **Reasoning handoff:** This area document must be explicitly referenced by the later sprint `reasoning.md`; otherwise these package and commit-boundary decisions can be lost during planning. That handoff is required, but editing final reasoning is outside this task's permitted output.
