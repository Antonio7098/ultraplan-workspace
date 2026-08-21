> **Inputs Used:** `projects/ultraplan-go/sprints/34-shared-context/technical-handbook.md`, `projects/ultraplan-go/sprints/34-shared-context/requirements.md`, `projects/ultraplan-go/sprints/34-shared-context/sprint-index.md`, `projects/ultraplan-go/docs/ARCHITECTURE.md`, `projects/ultraplan-go/docs/PRD.md`, `projects/ultraplan-go/docs/TRD.md`, `studies/go-cli-study/reports/final/01-project-structure.md`, `studies/go-cli-study/reports/final/02-command-architecture.md`, `studies/go-cli-study/reports/final/03-dependency-injection.md`, `studies/go-cli-study/reports/final/05-error-handling.md`, `studies/go-cli-study/reports/final/06-io-abstraction.md`, `studies/go-cli-study/reports/final/07-state-context.md`, `studies/go-cli-study/reports/final/10-logging-observability.md`, `studies/go-cli-study/reports/final/11-testing-strategy.md`, `studies/go-cli-study/reports/final/13-security.md`, `studies/go-cli-study/reports/final/14-performance.md`, `studies/go-cli-study/reports/final/15-philosophy.md`, `system/reasoning/architecture_reasoning_template.md`

# Architecture: Shared Sprint Context

This area covers ownership and composition of the Phase 5 shared prompt prefix, transient resolution of source references, integration across downstream agent-backed sprint operations, and the boundaries that keep product semantics out of generic runtime and interface packages. It does not redesign `code-context`, add retrieval or caching, change product-artifact authority, or pull Sprint 35 durable-run work forward.

## Area Decisions

### Final conclusion

The existing module architecture fits partially and needs a small refactor before broad integration. `internal/sprint` is already the correct owner, but stage-local prompt construction must converge on one concrete, package-local shared-context renderer before the prefix is added to every call site. The renderer will prepare one immutable prefix per top-level governed operation; stage code will append its own instructions after a single explicit boundary. `internal/platform/runtime` will continue to receive an ordinary generic prompt and will not know that a prefix, sprint, requirements artifact, or context pack exists.

### Ownership and dependency direction

- `internal/sprint/prompt_context.go` will own loading the two stored artifacts, parsing only the reference fields needed from `code-context.md`, resolving referenced source bytes, enforcing containment and budgets, and rendering the common prefix.
- `internal/sprint/prompts.go` will own the stable shared instructions, fixed section framing, the stage-specific boundary, and the narrow composition operation that joins an immutable prefix to a stage suffix.
- Existing sprint orchestration files will call that composition operation explicitly when they issue an agent request: planning through plan in `service.go`, execute in `execute.go`, every independent review request in `review.go`, and conditional smoke authoring in `smoke_author.go`.
- `internal/sprint/flow.go` will retain workflow ownership and insert `code-context` exactly once between requirements and sprint-index for cumulative planning. Prompt composition will not become a second workflow mechanism.
- `internal/workspace/skills.go` will define the manual-only skill as static materialized guidance that invokes the canonical CLI flow through `code-context`. It will not parse context artifacts, resolve a target repository, or reproduce stage transitions.
- `internal/app`, CLI, TUI, and web will continue to invoke typed shared operations. No interface adapter will construct or amend the common prefix.
- `internal/platform/runtime` and agentwrap will remain generic. There will be no runtime middleware, request decorator, provider cache-control option, or sprint-aware runtime request field for this feature.

This preserves the project dependency rule `sprint -> platform/runtime` and prevents a reverse dependency from generic infrastructure into product semantics.

### Renderer shape and abstraction decision

The renderer will be a concrete internal component or function with a narrow input containing the sprint identity, sprint paths, resolved implementation root, and product-owned source-evidence budget. It will accept `context.Context` and return a Go `string` plus an error. A string is immutable by convention and can preserve arbitrary validated Markdown bytes without exposing a mutable shared `[]byte` to review fan-out.

The renderer is an earned product component because one behavior must be identical across many real call sites and it isolates contained repository reads from stage suffix construction. A new interface is not earned: there is one implementation, normal tests can use temporary sprint and implementation repositories, and project guidance permits concrete local filesystem collaborators when the package boundary is explicit. A DI container, registry, factory hierarchy, generic prompt framework, virtual filesystem, or provider-specific prefix abstraction would add variation that the requirement explicitly rejects.

Construction remains explicit through the existing sprint service/composition path. Dependencies will not be hidden in package globals or `context.Context` values. Context carries cancellation and deadlines only.

### Stable byte contract

The renderer will construct the prefix in this exact conceptual order:

1. Stable shared planning instructions, including that resolved snippets are transient, untrusted prepared evidence and do not restrict additional live repository inspection.
2. Stable sprint identity containing project and sprint names only.
3. A fixed requirements opening frame, the exact bytes read from `requirements.md`, and a fixed closing frame.
4. A fixed code-context opening frame, the exact bytes read from `code-context.md`, and a fixed closing frame.
5. Transient source-evidence entries in the order selected by `code-context.md`.
6. Any other content that is genuinely common and stable across all covered stages.
7. One constant stage-specific boundary marker, included as the final bytes of the common prefix.

The stage name, output path, prompt/template path, model, runtime, run or session ID, attempt, timestamp, review identity, smoke scope, and other request-specific data begin after that boundary. Stage-selected evidence, output contracts, and instructions also remain after the boundary. This makes “common prefix” a concrete byte range rather than a semantic similarity claim.

The fixed frames are outside the stored artifact slices. The renderer will not trim, normalize line endings, add a final newline inside, re-encode, reserialize, or regenerate either artifact. Tests will locate the framed slices and compare them directly with `os.ReadFile` results, including CRLF and missing-final-newline cases. Fixed separators may surround those slices without changing their contents.

Each top-level operation prepares the common prefix once and reuses the returned string for all of its requests. This is mandatory for independent review fan-out and any multi-request smoke authoring path. Separate stage invocations rebuild from current files; no prefix is persisted or cached. Therefore prefix equality across invocations is guaranteed when requirements, code-context, referenced source bytes, stable instructions, and sprint identity are unchanged. This does not create source staleness detection or a provider cache guarantee.

### Reference parsing and transient evidence

- The parser will recognize only the validated context-pack fields needed by this sprint: repository-relative `Path`, exact `Lines`, optional `Symbol`, and `Rationale`. It will not introduce a second manifest or rewrite `code-context.md`.
- Entries are rendered in document order. Exact duplicates and overlaps are not silently merged or reordered because doing so would change authored selection semantics; each occurrence counts against the budget.
- Line selection is byte-oriented over the current file with deterministic line-number semantics. It will preserve selected source bytes and their existing line endings rather than scanning with an implementation-dependent token limit or normalizing text.
- Every excerpt receives stable path/range framing and is explicitly labeled as untrusted source evidence, not executable instructions and not content persisted in the context pack.
- Missing files, malformed or descending ranges, zero/out-of-range lines, unsupported file kinds, containment failures, cancellation, and budget overflow fail prompt construction before runtime launch. The renderer will not omit or truncate evidence while presenting the result as complete.
- Resolution is sequential for Sprint 34. Deterministic order, straightforward cancellation, and a bounded number of targeted reads are more valuable than speculative parallelism. Concurrency can be reconsidered only after measurement shows source preparation is material.

### Containment and file safety

Repository content and generated references are untrusted inputs. For every reference, resolution will:

1. Reject empty, absolute, volume-qualified, and lexically escaping paths.
2. clean and join the relative path under the configured implementation root;
3. resolve and compare the canonical root and candidate using `filepath.Rel`, rejecting `..` or cross-volume results;
4. reject symlinks in the referenced path and reject non-regular files;
5. open the file, verify the opened handle is still a regular file, read only through the bounded path, and compare pre/post identity and canonical location before accepting bytes.

The implementation must fail closed if replacement or containment cannot be established. The no-symlink policy is deliberately stricter than lexical containment and avoids treating links as an implicit second source boundary. Hard links remain a residual local-filesystem limitation because Go cannot portably prove an inode's only directory ancestry; bounded read-only access and target-repository permissions limit impact, and the risk is recorded below.

The renderer will not execute source content or commands, recursively scan the repository, follow evidence pointers beyond selected references, or grant write permission to the implementation repository merely to build a prompt.

### Budget and performance policy

Source evidence will have a product-owned aggregate byte budget derived before rendering, with fixed framing and stage-suffix reserve accounted for explicitly. The budget is about rendered bytes, not reference count, so the sprint does not add a hard maximum excerpt count. Per-file and aggregate checked arithmetic will prevent a large or malicious range from causing unbounded allocation.

The renderer will buffer the final prefix once because exact byte comparison and reuse across request fan-out are primary contracts. Full buffering is acceptable only because all variable source evidence is bounded before assembly and files are read through bounded ranges. It will pre-size a buffer from checked known lengths where practical, but it will not add pooling, streaming assembly, indexing, or concurrency without measurements.

On overflow, construction returns an actionable typed error naming the offending reference and allowed/observed size. It does not truncate, prioritize, summarize, or persist a repaired context pack. The operator must rerun or explicitly edit the governed context selection through existing workflows.

### State, mutation, and compatibility

The renderer reads `requirements.md`, `code-context.md`, and selected implementation source files. It creates only ephemeral prompt text and has no durable state, cache key, source fingerprint artifact, database, or migration.

Existing atomic code-context replacement and last-valid-artifact behavior remain unchanged. A failed downstream prefix build cannot replace either governed input or advance the requested stage. Old sprints without code-context continue through the compatibility behavior already established by Sprint 33; Sprint 34 will not add a second fallback path inside individual stage builders.

The only workflow mutation is the required cumulative order in the canonical flow implementation. Direct stage invocation and cumulative flow must call the same product operation, so exact-once behavior is observable without duplicating a `code-context` runner.

### Error, cancellation, and observability design

Reference failures will use a small inspectable sprint-owned error classification for invalid path, containment, file kind, missing source, invalid range, changed-during-read, and budget overflow. Lower-level errors will be wrapped with `%w` and include safe repository-relative path/range context. Cancellation remains detectable with `errors.Is(err, context.Canceled)` or `context.DeadlineExceeded`; it will not be converted to a generic reference error.

The renderer checks context before each reference, during bounded read loops, and before returning the assembled prefix. It starts no detached goroutines and creates no background context. Runtime invocation receives the same operation context after successful construction.

Dynamic diagnostics remain outside functional prompt bytes. Existing logs/events may record operation, project, sprint, safe relative path, range, evidence byte count, duration, outcome, and error category, but never source contents, full prompts, secrets, or absolute paths by default. Run/session/attempt/timestamp fields stay in runtime metadata or structured diagnostics after the prefix boundary. No new logging framework or metrics subsystem is justified for this sprint.

### Verification obligations

- Unit tests for the renderer will cover exact stored-byte reuse, fixed order and boundary bytes, no-final-newline and CRLF inputs, deterministic reference order, duplicate/overlap behavior, line-range extraction, cancellation, changed-file detection, budget boundaries, and every containment/file-kind rejection.
- Cross-stage fake-runtime tests will capture representative sprint-index, handbook, area/final reasoning, plan, execute, each independent review request, and agent-backed smoke-author request. They will compare the entire prefix byte range, prove exactly one prefix per request, and prove dynamic/stage-specific values occur only after the boundary.
- Flow tests will assert `requirements -> code-context -> sprint-index -> technical-handbook -> area-reasoning -> reasoning -> plan` and exactly one code-context execution in `flow --to plan`.
- Regression tests will preserve failed/cancelled rerun behavior, last-valid code-context, interface state agreement, and conditional agent-free paths.
- Skill tests will inspect both generated files and prove manual-only metadata, dry-run behavior, customization preservation, `--force`, all-skill inclusion, and canonical CLI delegation.
- Gated dogfood will use a disposable workspace and current implementation repository. It may report pass only when runtime requests were sent, a valid context pack and final plan were produced, exact-once order was captured, and no prohibited location changed. Missing credentials, provider, or runtime capability is `blocked` with the exact prerequisite.

## Trade-Offs

| Decision | Benefit | Accepted Cost | Rejected Alternative |
| --- | --- | --- | --- |
| One sprint-owned renderer with explicit stage calls | One byte contract and visible product ownership across all downstream routes | Every agent-backed call site must be audited and tested | Duplicated stage-local assembly, which would drift; generic runtime middleware, which would reverse the dependency boundary |
| Concrete internal renderer without a new interface | Minimal API, traceable flow, and temporary-repository tests | Filesystem behavior is tested through the concrete boundary rather than a mock interface | DI framework, registry, service locator, or one-implementation interface |
| Immutable operation-scoped string | Review fan-out receives identical bytes and cannot mutate shared content | Prefix memory is retained for the operation | Rebuilding per reviewer, which permits source races; shared mutable byte slices; persisted cache |
| Fixed frames around exact artifact slices | Exact stored bytes and boundaries are independently testable | Prompt framing is a compatibility contract that requires deliberate changes | Trimming/normalizing Markdown or parsing and reserializing it |
| Current source resolution once per operation | Evidence reflects the live repository without a cache or staleness subsystem | Separate stages can differ if source changes between invocations | Persisted excerpts, source snapshots, hashes/provenance, or automatic amendment |
| Sequential contained reads | Deterministic ordering, simple cancellation, low implementation risk | Potentially higher latency for many references | Concurrent reads before measurement; repository indexing or retrieval |
| Bounded full buffering | Simple exact-byte assembly and cheap reuse | Memory is proportional to the configured prompt budget | Streaming final prompt, which complicates atomic failure and exact comparisons; unbounded `ReadFile` |
| Fail on any reference or budget error | The prompt never claims incomplete evidence is complete | One bad reference blocks the downstream stage | Silent omission, heuristic prioritization, or truncation without an explicit contract |
| Reject referenced symlinks | Strong, understandable repository containment | Valid repositories that select symlinked files must select a regular target path instead | Following links based only on lexical prefix checks |
| Dynamic facts outside the common prefix | Stable bytes remain compatible with useful diagnostics | Prefix logs cannot carry per-run correlation inline | Embedding timestamps, stage names, run IDs, or output paths in the shared block |
| Canonical CLI delegation for the manual skill | One implementation of target resolution, validation, state, and cancellation | The skill cannot customize stage mechanics | Reimplementing code-context generation in skill instructions |

The selected option is the smallest honest design: a focused `internal/sprint` refactor plus explicit integrations. A broader prompt pipeline could reduce apparent call-site repetition, but it would create hidden lifecycle behavior and a generic abstraction before a second product module has demonstrated the same semantics.

## Evidence

### Sprint and project constraints

- **Requirements finding:** `projects/ultraplan-go/sprints/34-shared-context/requirements.md` mandates one `internal/sprint` renderer, exact stored bytes, deterministic transient source evidence, identical common prefixes, all downstream agent-backed routes, contained bounded reads, exact-once flow, canonical manual-skill delegation, offline fakes, and no cache/retrieval/framework expansion. **Inference:** ownership, serialization, integration, and failure behavior must be one product contract rather than conventions repeated by stages.
- **Architecture finding:** `projects/ultraplan-go/docs/ARCHITECTURE.md` assigns sprint prompt rendering and flow to `internal/sprint`, generic execution to `internal/platform/runtime`, and interface composition to `internal/app`; it advises focused files before subpackages and concrete collaborators until volatility earns interfaces. **Inference:** `prompt_context.go` inside the existing package is preferable to a new package, service layer, runtime wrapper, or virtual filesystem.
- **Product finding:** `projects/ultraplan-go/docs/PRD.md` defines the code-context pack as durable prepared evidence while live source remains authoritative, requires one product core across local surfaces, and gates retrieval, provenance, persistence, and cloud work after Sprint 34. **Inference:** transient source resolution belongs in request preparation and must not create new authoritative state.
- **Technical finding:** `projects/ultraplan-go/docs/TRD.md` requires the generic runtime request to remain product-neutral, normal tests to use temporary workspaces and fake runtimes, product artifacts to remain authoritative, and Phase 5 to use one shared renderer with exact-prefix coverage. **Inference:** stage integrations should pass a complete prompt through the existing runtime request instead of changing agentwrap or platform contracts.

### Direct report findings and sprint-specific conclusions

- **Report finding:** `studies/go-cli-study/reports/final/01-project-structure.md` finds a recurring thin-boundary/interior-logic split and one-way imports, including Helm command-to-action delegation and yq's `cmd -> yqlib` direction. **Inference:** sprint stages may orchestrate, but shared prompt policy stays in the sprint interior and never moves into runtime or web adapters.
- **Report finding:** `studies/go-cli-study/reports/final/02-command-architecture.md` favors explicit factories and thin delegates, warns about large lifecycle handlers, and finds no formal middleware chain in the cohort. **Inference:** visible calls to one renderer at each agent boundary are safer than invisible prompt middleware or annotations.
- **Report finding:** `studies/go-cli-study/reports/final/03-dependency-injection.md` finds manual composition roots, constructor/factory injection, minimal globals, and no DI frameworks; it warns that central config objects and context-carried services can become hidden coupling. **Inference:** use the existing sprint construction path, pass narrow values explicitly, and do not add a registry, container, package global, or context service locator.
- **Report finding:** `studies/go-cli-study/reports/final/05-error-handling.md` supports `%w` wrapping, sentinels or typed errors for actionable conditions, user/operational separation, and aggregation only where partial work is meaningful. **Inference:** classify reference failures without string parsing, preserve cancellation causes, and fail one prompt build atomically rather than aggregate unsafe reads after a containment failure.
- **Report finding:** `studies/go-cli-study/reports/final/06-io-abstraction.md` shows that complete I/O seams and injected streams improve tests, while one direct bypass defeats consistency; it also notes large filesystem interfaces impose real cost. **Inference:** all source reads must pass through the one renderer, but temporary real repositories are sufficient and a broad new filesystem interface is not earned.
- **Report finding:** `studies/go-cli-study/reports/final/07-state-context.md` ties operational maturity to root-context propagation and structured cancellation, and warns against `context.Background()` in long work and services hidden in context values. **Inference:** the renderer accepts the operation context, performs no detached work, and stores no service dependencies in context.
- **Report finding:** `studies/go-cli-study/reports/final/10-logging-observability.md` finds structured fields and strict functional-output/diagnostic separation to be the main observability quality divider, while output-abstraction bypasses contaminate behavior. **Inference:** dynamic correlation belongs in logs/runtime metadata, never in the canonical prompt prefix.
- **Report finding:** `studies/go-cli-study/reports/final/11-testing-strategy.md` supports golden output comparisons, deterministic fakes, realistic command pipelines, and behavior assertions; it warns that regex/substring checks miss output regressions and golden updates can bless mistakes. **Inference:** compare raw prefix bytes and exact artifact slices at both renderer and representative stage levels, with deliberate fixture updates.
- **Report finding:** `studies/go-cli-study/reports/final/13-security.md` associates strong tools with visible trust boundaries and centralized validation, and identifies raw path comparison as a canonicalization weakness. **Inference:** generated paths and source text are untrusted; containment must include canonical checks, symlink/file-kind policy, bounded reads, and explicit evidence framing.
- **Report finding:** `studies/go-cli-study/reports/final/14-performance.md` recommends lazy work, bounded streaming/concurrency, and profiling before optimization; it warns that buffering and queues remain dangerous when unbounded. **Inference:** prepare only for agent-backed operations, use a byte cap, keep sequential reads initially, and buffer once only after size checks.
- **Report finding:** `studies/go-cli-study/reports/final/15-philosophy.md` emphasizes explicit non-goals, simplicity-first decisions, and abstractions earned by real implementations. **Inference:** one focused renderer is justified by current duplication pressure, while caches, plugins, indexes, provider controls, persistence, and a generic prompt framework are not.

## Risks

### Material risks and mitigations

- **Byte drift at one call site:** a stage could append content before the boundary or bypass composition. Mitigation: one suffix API, captured-request tests for every route, and an assertion for exactly one boundary/prefix per request.
- **Stored-byte corruption by “helpful” formatting:** trimming or newline normalization can violate the central contract. Mitigation: append raw file bytes between fixed frames and compare slices directly in tests.
- **Repository replacement race:** portable path checks cannot provide filesystem-wide transactional reads. Mitigation: reject symlinks, validate before and after opening/reading, compare file identity and canonical location, stop on any mismatch, and keep the operation read-only and bounded. Hard-link ancestry remains a documented residual risk.
- **Prompt injection from source text:** selected code may contain instruction-like content. Mitigation: stable instructions declare snippets untrusted evidence, frame every excerpt, retain task permission policy, and never interpret source text as commands or prompt configuration.
- **Source changes between stages:** because evidence is transient and uncached, separate invocations can produce different prefixes after repository edits. This is intentional live-source behavior, not a staleness feature. Equality claims and dogfood evidence are valid only while referenced bytes are unchanged.
- **Budget mismatch with real provider limits:** a product byte bound may still map imperfectly to provider tokens. Mitigation: reserve suffix capacity, fail before launch on the product bound, exercise representative dogfood, and adjust the named bound without changing serialization or introducing provider cache semantics.
- **Central renderer growth:** adding stage-specific context to the common component could turn it into a god object. Mitigation: common deterministic material only; stage-specific selected evidence and contracts remain suffix responsibilities.
- **Partial integration:** review fan-out and conditional smoke paths are easy to miss. Mitigation: enumerate call sites from requirements, capture every independent request, and include tests where the conditional path is agent-free.
- **Cancellation latency during local reads:** standard filesystem calls are not context-aware. Mitigation: bounded range reads with context checks between chunks and references; no unbounded whole-repository read.
- **Overly strict symlink rejection:** a valid selected reference may point through a repository symlink. Mitigation: fail with the exact relative path and require selection of a regular contained target; do not weaken containment silently.
- **Golden fixture misuse:** an update command could normalize an unintended compatibility break. Mitigation: pair full golden comparisons with focused assertions for raw artifact slices, boundary count, and dynamic-data exclusion; review fixture diffs explicitly.
- **Skill behavior drift:** generated skill text could diverge from canonical CLI semantics. Mitigation: derive materialized content from one embedded definition and test metadata/content synchronization plus actual command delegation.
- **Dogfood overstatement:** a built but unsent prompt or blocked provider cannot prove real flow completion. Mitigation: retain command/runtime evidence and report unavailable credentials or capabilities as blocked, never passed.

### Assumptions and open verification questions

- The Sprint 33 artifact validator remains the authority for context-pack structure; Sprint 34 adds parsing needed for rendering but does not create a competing validator or compatibility grammar.
- The configured implementation root is resolved by existing project/sprint mechanisms before renderer invocation; the renderer independently enforces containment for each selected reference.
- The current runtime prompt model accepts a Go string containing validated UTF-8 Markdown. Arbitrary binary or invalid UTF-8 source is rejected as unsupported source evidence rather than normalized.
- Review must verify whether any representative selected context uses symlinks or files changing during preparation. The decided behavior is fail closed; observed incompatibility requires changing the selected reference, not silently relaxing policy.
- Dogfood must establish that the chosen aggregate byte bound leaves sufficient room for the largest real stage suffix. If it does not, the bound or selected evidence is adjusted explicitly while the fail-on-overflow decision remains unchanged.

The architecture decision is **Proceed after the small shared-renderer refactor**. Complexity introduced is low to medium: one focused product component, explicit integrations, and containment/error tests. Complexity removed is medium: duplicated prefix policy and the risk of divergent downstream prompts. The accepted trade-off is bounded buffering and visible call-site wiring in exchange for exact bytes, deterministic behavior, and preserved module boundaries.
