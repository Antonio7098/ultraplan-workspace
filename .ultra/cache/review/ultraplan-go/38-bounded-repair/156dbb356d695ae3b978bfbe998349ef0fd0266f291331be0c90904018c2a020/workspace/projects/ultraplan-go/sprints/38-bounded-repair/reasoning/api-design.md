> **Inputs Used:** `projects/ultraplan-go/sprints/38-bounded-repair/requirements.md`, `projects/ultraplan-go/sprints/38-bounded-repair/sprint-index.md`, `projects/ultraplan-go/sprints/38-bounded-repair/technical-handbook.md`, `projects/ultraplan-go/docs/ARCHITECTURE.md`, `projects/ultraplan-go/docs/PRD.md`, `projects/ultraplan-go/docs/TRD.md`, `system/reasoning/api-design-reasoning-template.md`, `studies/go-cli-study/reports/final/01-project-structure.md`, `studies/go-cli-study/reports/final/02-command-architecture.md`, `studies/go-cli-study/reports/final/03-dependency-injection.md`, `studies/go-cli-study/reports/final/04-configuration-management.md`, `studies/go-cli-study/reports/final/05-error-handling.md`, `studies/go-cli-study/reports/final/06-io-abstraction.md`, `studies/go-cli-study/reports/final/07-state-context.md`, `studies/go-cli-study/reports/final/08-concurrency.md`, `studies/go-cli-study/reports/final/09-terminal-ux.md`, `studies/go-cli-study/reports/final/10-logging-observability.md`, `studies/go-cli-study/reports/final/11-testing-strategy.md`, `studies/go-cli-study/reports/final/12-extensibility.md`, `studies/go-cli-study/reports/final/13-security.md`, `studies/go-cli-study/reports/final/14-performance.md`

# API design reasoning

This document reasons through the selected API Design template at `system/reasoning/api-design-reasoning-template.md`. It covers the internal application contract and its CLI, JSON, TUI, and loopback HTTP projections. It does not decide patch parsing, production application, or repair outcome rules. Those remain in `internal/sprint`. The API must make the authority sequence impossible to skip: freeze a packet, show the exact proposed authority, durably accept one operation, record a single-use confirmation bound to that acceptance, and only then start runtime or production work.

## Area Decisions

### One app contract owns every adapter path

Add an adapter-independent `RepairUseCases` contract in `internal/app/sprint_usecases.go`. CLI commands, TUI actions, HTML handlers, and JSON handlers call this contract or the existing generic operation contract. None of them may read `verification/`, construct confirmation digests, derive outcomes, or call `internal/sprint` directly.

The contract should expose these operations:

| Operation | Request identity | Result | Side effects |
| --- | --- | --- | --- |
| `PrepareRepair` | project, sprint, current issue ID, requested mode | bounded packet preview | Freezes or reuses the current immutable packet. It does not confirm, start a runtime, or mutate production. |
| `ConfirmRepair` | prepared packet digest, canonical request digest, target identity, mode, effective limits, durable acceptance | accepted repair run | Writes the immutable confirmation after durable acceptance. This method is called by the shared operation layer, not directly by an adapter. |
| `StartRepair` | accepted run identity and writer fence | canonical repair result | Claims the accepted operation and runs the shared repair protocol. It cannot create or alter confirmation authority. |
| `RepairStatus` | project, sprint, optional run ID | current bounded status | Read-only canonical reload. |
| `RepairPacket` | project, sprint, repair run ID | bounded packet DTO | Read-only packet projection without private evidence bodies. |
| `RepairCycles` | run ID, opaque cursor, page limit | bounded cycle page | Read-only paged summaries. |
| `RepairCycle` | run ID, cycle number | one bounded cycle DTO | Read-only scope, gate, cleanup, and evidence-reference summary. |
| `RepairResult` | run ID | terminal result DTO or explicit non-terminal state | Read-only canonical result. |
| `ResumeRepair` | run ID and current canonical confirmation facts | accepted operation result | Resumes only from a proven boundary and preserves consumed limits and deadline. |
| `CancelRepair` | run ID | cancellation acknowledgement plus refreshed status | Sends canonical idempotent cancellation to the current owner. |
| `RecoverRepair` | project, sprint, optional run ID | reconciliation result plus refreshed status | Reconciles abandoned authority without inferring success. |

This is a focused contract, not a general workflow engine. The structure and command studies support thin adapters and one delegated operation path, especially `studies/go-cli-study/reports/final/01-project-structure.md` and `studies/go-cli-study/reports/final/02-command-architecture.md`. The dependency-injection report supports explicit runtime, store, clock, process, and writer-fence dependencies at the app composition point rather than globals or context values: `studies/go-cli-study/reports/final/03-dependency-injection.md`.

### Preparation, operation preparation, confirmation, and start are separate steps

The name "prepare" exists at two levels and the API must not blur them.

1. `PrepareRepair` performs domain admission and freezes `issue-packet.json`. It returns a packet preview and digest.
2. Existing `PrepareOperation` prepares `OperationRepairStart` or `OperationRepairResume`. It returns the canonical request, governed-input fingerprint, warnings, mutation class, and confirmation facts.
3. Existing `RunOperation` validates the current preparation, durably accepts the operation, obtains its run identity and fencing generation, and invokes `ConfirmRepair` with that acceptance.
4. The shared durable runner calls `StartRepair` only after the confirmation record is committed.

No goroutine or runtime starts between steps 1 and 3. If durable acceptance or confirmation persistence fails, `StartRepair` is never called. This sequencing satisfies the requirement that confirmation bind durable acceptance without giving adapters an acceptance token they can invent. It also retains the existing operation manager as the only source of run identity and writer ownership.

The confirmation authority is a digest over the packet digest, full target identity, canonical operation request, selected mode, effective limits, governed-input fingerprint, durable run ID, operational attempt ID, and fencing generation. A prepared-operation token may carry this material to the browser, but it is short-lived transport state, not the authority record. The immutable `confirmation.json` is the durable authority record.

### Manual and automatic authorization are different request shapes

`RepairMode` is a closed enum with `manual` and `automatic`. Manual is the default only for packet preview. Starting any mutation requires an explicit confirmation fact.

Automatic confirmation requires both `mode: "automatic"` and `automatic_opt_in: true`. A request with automatic mode and no opt-in is invalid. A request with manual mode and automatic opt-in is also invalid. CLI uses `--automatic --yes`; TUI and browser use a separate automatic action or control followed by the guarded confirmation action. A mode value inferred from configuration, a previous manual confirmation, a page visit, an SSE reconnect, or an operation replay cannot satisfy this field.

Request flags may select only a lower effective limit where requirements expressly permit request selection. They can never raise configured limits or product maxima. Prefer no per-request budget flags in the first manual path beyond the required mode choice. Automatic `--max-cycles` may lower the effective cycle count. The canonical request stores both the requested value and the effective lower value so the confirmation page can show what will run.

### Operation kinds and mutation classes stay explicit

Add closed operation kinds for `repair-start`, `repair-resume`, and `repair-recover`. Packet preparation remains a sprint-owned state mutation but never starts runtime or production work. It may use a dedicated synchronous `repair-prepare` operation kind if the existing operation policy requires all private-state writes to pass through operation preparation. It must not be mislabeled read-only.

`repair-start` and `repair-resume` use the strongest production-mutation class. `repair-recover` uses a verification-state mutation class and cannot apply a patch. Status, packet, cycle, and result calls are queries, not operation kinds. Cancellation continues through the durable run-control cancellation API rather than adding a second repair-specific owner registry.

Automatic mode reuses `repair-start` and the same sprint service. It is a field in the canonical request, not a second operation family or runner. This rejects a separate automatic API that could drift from manual scope, cleanup, or outcome checks.

### Public DTOs contain facts, not private records

Use one versioned `RepairResult` projection across CLI JSON, TUI models, and HTTP JSON. HTML view models may reshape it for presentation but cannot add product conclusions. The core projection contains:

| Group | Public fields |
| --- | --- |
| Identity | schema version, project, sprint, QA attempt ID, repair run ID, operation run ID, packet ID and digest, issue ID, root-cause-group ID |
| Authority | mode, confirmation state, confirmer label, confirmation time, target fingerprint, governed-input fingerprint, policy fingerprint, schema fingerprint |
| Current state | repair phase, freshness, run lifecycle, current cycle, current gate, semantic outcome, blocker, next action |
| Scope | bounded allowed-path count and preview, forbidden-path count and preview, actual changed-path count and preview, changed files and bytes, scope result |
| Limits | effective totals, consumed totals, remaining totals, deadline status, effective-source summaries |
| Verification | ordered gate summaries, pass/fail/blocked/skipped status, finding counts, issue-set delta, highest-severity delta |
| Cleanup | attempted, process-tree complete, workspace removed, lock reconciled, complete, uncertainty reason |
| Evidence | digest-bound packet, cycle, cleanup, reverification, result, and manual-proof references |
| Recovery | resumable flag, latest proven boundary, cancellation facts, recovery or escalation action |

Never expose full prompts, proposal patch bytes, production contents, raw provider payloads, unrestricted command output, environment values, or private evidence bodies. Text fields use the existing sanitizer and redaction path before entering app DTOs. Path and evidence lists are capped and report `total`, `returned`, and an opaque continuation cursor.

`RepairStatus` is the canonical current view. Progress events carry only run correlation, phase, cycle, gate, bounded counters, and sanitized text. Consumers must reload `RepairStatus` after reconnect, a replay gap, terminal notification, cancellation acknowledgement, or process restart.

The IO and observability reports support a shared projection and separate user output from durable diagnostics: `studies/go-cli-study/reports/final/06-io-abstraction.md` and `studies/go-cli-study/reports/final/10-logging-observability.md`. The terminal report supports a complete non-interactive path and interruptible progress without making the progress channel authoritative: `studies/go-cli-study/reports/final/09-terminal-ux.md`.

### Query bounds are fixed by product policy

All collection queries use opaque cursors bound to project, sprint, run ID, packet digest, collection kind, and current retained lower boundary. A cursor for another run or changed packet returns `stale_request`; an expired retained cycle returns a typed retention error with the current earliest cycle and next action.

The first implementation should use conservative fixed defaults and maxima from app policy:

| Query | Default | Hard maximum |
| --- | --- | --- |
| Cycle page | 20 | 100 |
| Path preview | 20 | 100 |
| Gate summaries per cycle | complete fixed ladder | fixed by the packet, not caller-controlled |
| Evidence references | 50 | 200 |
| Public text field | existing safe projection limit | existing app maximum |

Mutation budgets and query page sizes are different concepts. Query flags cannot change repair budgets. The performance report supports incremental bounded reads and warns against unbounded retained collections: `studies/go-cli-study/reports/final/14-performance.md`.

### HTTP uses repair resources for reads and shared guarded operations for writes

Add versioned read resources under the existing sprint hierarchy:

```text
GET /api/v1/projects/{project}/sprints/{sprint}/repair
GET /api/v1/projects/{project}/sprints/{sprint}/repair/runs/{run_id}/packet
GET /api/v1/projects/{project}/sprints/{sprint}/repair/runs/{run_id}/cycles
GET /api/v1/projects/{project}/sprints/{sprint}/repair/runs/{run_id}/cycles/{cycle}
GET /api/v1/projects/{project}/sprints/{sprint}/repair/runs/{run_id}/result
```

Mutations continue through the shared guarded-operation endpoints. `POST /api/v1/operations/prepare` accepts a repair operation request, and `POST /api/v1/operations` accepts the server-issued confirmation. Existing durable run status, event, and cancellation routes remain canonical. Server-rendered no-JavaScript forms may use sprint-scoped HTML routes, but their handlers must map to the same app requests and guarded operation calls.

Successful start returns `202 Accepted` with the durable run resource and canonical repair status link. Query success returns `200`. A semantic terminal outcome such as `failed` or `escalated` is a successful resource read with `200`; it is not rewritten as an HTTP transport failure. Invalid input is `400`, authorization or CSRF failure is `403`, unknown retained identity is `404`, stale or conflicting authority is `409`, expired retained history is `410`, valid but inadmissible repair is `422`, unavailable required capability is `503`, and persistence failure is `500`. Every error uses the existing versioned public envelope with a stable category, safe message, correlation ID, and next action.

### CLI commands mirror resources and preserve script behavior

Use one `repair` command family:

```text
ultraplan sprint <project> <sprint> repair prepare --issue <id> [--automatic] [--json]
ultraplan sprint <project> <sprint> repair start --run <repair-run-id> --yes [--automatic] [--max-cycles <lower-value>] [--json]
ultraplan sprint <project> <sprint> repair status [--run <repair-run-id>] [--json]
ultraplan sprint <project> <sprint> repair packet --run <repair-run-id> [--json]
ultraplan sprint <project> <sprint> repair cycles --run <repair-run-id> [--cursor <cursor>] [--limit <n>] [--json]
ultraplan sprint <project> <sprint> repair result --run <repair-run-id> [--json]
ultraplan sprint <project> <sprint> repair resume --run <repair-run-id> --yes [--json]
ultraplan sprint <project> <sprint> repair cancel --run <operation-run-id> [--json]
ultraplan sprint <project> <sprint> repair recover [--run <repair-run-id>] [--json]
```

`prepare` always renders the packet facts needed for review and never accepts `--yes`. `start` requires `--yes`, and automatic start also requires `--automatic`. `resume` requires a fresh confirmation because it may continue mutable work, but it cannot change the original mode or effective limits. The app rejects unsupported flags before operation preparation.

CLI JSON uses a single document with `schema_version`, `operation`, `status`, `result`, and optional `error`. Text progress goes to stderr. JSON stdout contains no progress lines. `verified` and `verified_with_findings` return exit code 0 after a start or resume. `failed`, `blocked`, `escalated`, `stalled`, cancellation, interruption, and cleanup uncertainty return non-zero using existing usage, configuration, validation, runtime, cancellation, and partial-result classes. Read-only status and result queries return zero when the query succeeds, even when they report a non-success semantic outcome.

### Idempotency is explicit at every transition

Packet preparation with identical current inputs returns the same active prepared packet and digest. It does not publish duplicate current packets. After a confirmed or terminal run, a new repair attempt gets a new run identity even if the packet content digest is unchanged.

Durable acceptance deduplicates the confirmation digest. Repeating the same confirmed start before dispatch returns the same accepted run. Repeating it after confirmation consumption returns a replay response that points to the canonical run and does not authorize new work. Repeating start for an already running or terminal accepted run returns its current canonical status.

Cancellation is idempotent. Resume is idempotent at the latest proven boundary and cannot repeat a committed apply. Recovery is repeatable and may only move state through conservative reconciliation. These rules must survive process restart because their facts live in durable operation state and strict repair records, not adapter memory.

### Operational lifecycle and repair outcome remain separate

The API exposes two closed fields:

| Field | Values | Authority |
| --- | --- | --- |
| `run_lifecycle` | accepted, running, cancellation_requested, cancelled, interrupted, cleanup_uncertain, terminal | durable operation manager |
| `outcome` | verified, verified_with_findings, failed, blocked, escalated, stalled | sprint repair result |

A repair result is absent until the sprint service commits one. Cancellation does not manufacture `failed`, and run-control completion does not manufacture `verified`. If cleanup uncertainty exists, the durable lifecycle is `cleanup_uncertain` and the repair result, when safely publishable by the current owner or recovery, is `escalated`. Late completion cannot replace either terminal authority.

Typed public errors should preserve at least these categories: `invalid_request`, `admission_rejected`, `stale_request`, `confirmation_required`, `confirmation_mismatch`, `confirmation_replayed`, `conflict`, `limit_exhausted`, `runtime_unavailable`, `check_unavailable`, `persistence_failure`, `writer_fenced`, `target_drift`, `cleanup_uncertain`, `cancelled`, `interrupted`, `unknown_schema`, and `retention_gap`. Safe recovery guidance is a field, not prose inferred by adapters.

The error report supports typed facts and exit mapping rather than message parsing: `studies/go-cli-study/reports/final/05-error-handling.md`. The state and concurrency reports support root-context propagation, durable ownership, bounded cleanup contexts, localized goroutine starts, and explicit waits: `studies/go-cli-study/reports/final/07-state-context.md` and `studies/go-cli-study/reports/final/08-concurrency.md`.

### Compatibility is additive and versioned

Add repair fields to the canonical app fixture without changing existing QA field meaning. New interfaces should be additive where existing test adapters implement narrower QA query contracts. Keep repair under a new `RepairUseCases` interface rather than expanding `QAUseCases` and breaking every implementation.

HTTP and CLI JSON remain schema version 1 only if the existing compatibility policy treats optional resource fields and new operation kinds as additive. Private repair records have their own strict schema versions and reject unknown versions. No migration is needed because repair artifacts do not exist before this sprint. Unknown future private schemas block repair; they are not treated as empty or old state.

The extensibility report supports explicit version dispatch but warns against registries and plugin systems that silently overwrite or widen trust: `studies/go-cli-study/reports/final/12-extensibility.md`. Repair needs closed enums and explicit operation registration, not a plugin API.

### Authorization is local but still object-scoped

CLI authority comes from the local process and explicit flags. TUI authority comes from the current guarded action. Browser mutations require the existing same-origin session, CSRF validation, request body limit, short-lived prepared-operation token, and current authorization at both confirmation and cancellation time.

Every request resolves project, sprint, repair run, packet, and durable operation ownership together. A run ID from another sprint, a packet from a historical QA attempt, or a cancellation request from an unauthorized browser session fails before side effects. Read visibility follows the existing workspace-local policy and never grants mutation authority.

The security report supports explicit permission checks, argument arrays, private temporary storage, schema validation, and redaction: `studies/go-cli-study/reports/final/13-security.md`. For this sprint, those patterns are mandatory because the API authorizes production mutation rather than an ordinary command.

## Trade-Offs

| Decision | Benefit | Cost | Rejected alternative |
| --- | --- | --- | --- |
| Separate packet preparation from generic operation preparation | Makes immutable packet authority visible before confirmation and keeps runtime acceptance reusable. | Users see a two-step prepare/start flow, and internal naming needs care. | A single `repair --issue` call is too easy to treat as implicit confirmation and cannot bind a viewed packet to durable acceptance cleanly. |
| Confirmation after durable acceptance but before dispatch | Binds the exact accepted run and writer generation while preserving acceptance-before-goroutine. | A failed confirmation leaves an accepted operation that must terminate truthfully. | Confirmation before durable acceptance cannot include acceptance identity. Starting before confirmation violates the mutation gate. |
| Dedicated `RepairUseCases` interface | Keeps repair DTOs cohesive and avoids breaking QA-only adapters. | Adds another app contract and test fake. | Expanding `QAUseCases` would couple read-only QA consumers to production mutation and create broad source breakage. |
| Shared generic operation endpoints for mutations | Reuses confirmation, CSRF, durable acceptance, cancellation, shutdown, and replay behavior. | Repair-specific HTML handlers still need mapping code. | Repair-only HTTP mutation endpoints would duplicate the most security-sensitive lifecycle code. |
| Distinct lifecycle and semantic outcome fields | Preserves the authority of run control and sprint repair. | Clients must display two related states. | One combined status would incorrectly map cancellation, cleanup uncertainty, or runtime completion to a repair verdict. |
| Opaque digest-bound pagination cursors | Prevents cross-run and stale collection reads while bounding responses. | Cursors are not human-editable and expired history needs a typed response. | Numeric offsets can silently move after retention and do not bind packet authority. |
| Explicit `automatic_opt_in` in addition to mode | Makes automatic consent unambiguous and prevents a manual token from being transformed. | Automatic clients send one redundant-looking field. | Inferring consent from config or `mode` alone is too weak for a higher-risk operation. |
| Fixed resource reads and generic operation writes | Gives repair stable resource URLs without creating a second run registry. | Clients use both repair and run resources. | Making SSE or browser operation objects authoritative would break reconnect and cross-process observation. |
| Additive public schema, strict private schema | Preserves existing clients while failing closed on state that controls mutation. | Public clients must ignore unknown optional fields, while private readers need explicit version code. | Lenient private decoding could accept authority fields with changed meaning. |

The selected reports do not prove that every API must use these exact names. They do support the underlying choices: thin delegated commands, explicit composition, post-merge configuration validation, typed errors, injectable IO, context propagation, bounded concurrency, canonical status reload, versioned metadata, and fail-closed permission checks. The names above fit the repository's existing app and operation vocabulary and avoid a new generic repair framework.

## Evidence

| Decision area | Report finding | Sprint-specific conclusion |
| --- | --- | --- |
| Package and adapter boundary | `studies/go-cli-study/reports/final/01-project-structure.md` finds consistent one-way imports and thin entrypoints, including Helm's command/action split and gdu's UI abstraction. | Put repair rules in `internal/sprint`, shared requests and projections in `internal/app`, and transport mapping in adapters. |
| Command contract | `studies/go-cli-study/reports/final/02-command-architecture.md` favors factory-built commands and shared wrappers over large handlers. | CLI parsing ends at typed requests; the shared operation runner owns durable repair dispatch. |
| Dependency construction | `studies/go-cli-study/reports/final/03-dependency-injection.md` favors explicit composition roots and constructor injection, while warning about context service locators and oversized containers. | Inject repair runtime, process, store, clock, and fences through existing service construction. Do not hide confirmation or ownership in globals. |
| Effective limits | `studies/go-cli-study/reports/final/04-configuration-management.md` finds explicit precedence and validation after merging all sources to be the reliable pattern. | Freeze effective lower-only limits before operation preparation and include their sources in confirmation and status. |
| Errors and exits | `studies/go-cli-study/reports/final/05-error-handling.md` supports wrapped typed errors, behavioral classification, safe user rendering, and exit mapping. | Preserve stable categories and recovery actions. Never parse error strings to derive outcomes or HTTP status. |
| Testable adapter IO | `studies/go-cli-study/reports/final/06-io-abstraction.md` finds injectable streams and boundary interfaces necessary for command-level failure tests. | Keep CLI text, JSON, progress, and errors on injected writers. Use fakes for operation and persistence failures. |
| Cancellation and cleanup | `studies/go-cli-study/reports/final/07-state-context.md` supports root context propagation and a separate bounded cleanup lifetime. | Start, resume, cancel, server shutdown, commands, process trees, and cleanup share canonical cancellation while cleanup records uncertainty if its bound expires. |
| Durable dispatch | `studies/go-cli-study/reports/final/08-concurrency.md` favors localized goroutine starts, explicit waits, bounded work, and no fire-and-forget tasks. | The shared runner is the only launch site. Acceptance and confirmation finish before launch, and shutdown waits within the cleanup bound. |
| Interactive and non-interactive parity | `studies/go-cli-study/reports/final/09-terminal-ux.md` supports non-TTY operation, visible interruptible progress, and clear cancellation. | `--yes` is explicit and non-interactive. TUI and browser confirmation remain separate guarded actions. Progress never replaces status reload. |
| Correlation and safe output | `studies/go-cli-study/reports/final/10-logging-observability.md` favors structured correlation and separation of user output from diagnostics. | Correlate packet, repair run, operation run, attempt, cycle, writer generation, and runtime IDs without emitting private payloads. |
| Contract testing | `studies/go-cli-study/reports/final/11-testing-strategy.md` supports table tests, command-path integration, fakes, and selective golden fixtures. | Freeze JSON envelopes and text help with fixtures, then use semantic tests for confirmation replay, stale requests, cursor binding, outcome mapping, and shutdown races. |
| Versioning and extension restraint | `studies/go-cli-study/reports/final/12-extensibility.md` supports explicit schema versions and subprocess isolation, while warning about registry collisions and plugin trust. | Use closed repair enums and strict version dispatch. Do not add a plugin or user-defined operation registry. |
| Mutation authorization | `studies/go-cli-study/reports/final/13-security.md` supports visible permission gates, explicit argv, private temporary directories, redaction, and fail-fast schema checks. | Treat confirmation mismatch, path scope, stale target, unknown schema, and unsafe output as hard admission failures before mutation. |
| Response and retention bounds | `studies/go-cli-study/reports/final/14-performance.md` favors lazy setup, streaming, incremental state, and bounded collections. | Keep queries paged, status summaries finite, packet preparation runtime-free, and help/status free of runtime initialization. |

The project documents fix the broader constraints. `projects/ultraplan-go/docs/ARCHITECTURE.md` assigns repair semantics to `internal/sprint`, app contracts to `internal/app`, and transport concerns to CLI, TUI, and web packages. `projects/ultraplan-go/docs/PRD.md` requires one product core across local interfaces and a run that outlives its observer. `projects/ultraplan-go/docs/TRD.md` requires guarded HTTP operations, durable acceptance before execution, conservative reconciliation, and strict separation between operational run facts and product outcomes.

The governed sprint requirements sharpen those constraints. AC-2 requires single-use confirmation over packet, target, request, mode, limits, and durable acceptance. AC-3 requires the same manual path through all adapters. AC-7 requires automatic mode to reuse that path and preserve limits across resume. AC-8 requires closed semantic outcomes. AC-9 requires canonical cancellation and recovery. AC-10 requires bounded, redacted, versioned parity across CLI, JSON, TUI, and browser.

## Risks

### Accepted operations that fail confirmation publication

Durable acceptance must precede confirmation because confirmation includes acceptance identity. A disk or writer-fence failure can therefore leave an accepted operation that never becomes confirmed. The operation manager must finish it as blocked or persistence failure without launching work. Recovery must recognize this exact boundary and must not manufacture a confirmation.

### Two identifiers may confuse users and clients

Repair run ID and durable operation run ID represent different records. Public DTOs must always return both with explicit names. CLI cancellation accepts the operation run ID; repair packet and cycle queries accept the repair run ID. Status should provide links or next commands with the correct identifier.

### Generic operation preparation may omit repair-specific facts

The existing confirmation model may not carry packet digest, target identity, mode, effective limits, and automatic opt-in. Extending it incompletely would produce a confirmation that looks valid but authorizes less than the requirements demand. Tests must compare every digest input and fail if any field changes.

### Additive interfaces can still drift

Separate `RepairUseCases`, `RepairQueries`, generic operations, and durable run queries could return overlapping but inconsistent facts. `RepairStatus` must remain the canonical product projection, and run queries must remain operational. The canonical fixture should assert their correlation and intentional differences.

### Resume confirmation could accidentally widen authority

Resume needs fresh user confirmation but cannot choose another mode, packet, target, policy, or budget. The request should carry only the run ID and confirmation proof. Product code reloads the frozen values. Allowing callers to resend editable limits or scope during resume would create an authority expansion path.

### HTTP status codes can hide semantic results

Mapping `failed`, `blocked`, `escalated`, or `stalled` to HTTP errors would encourage clients to discard the result body and infer meaning from transport. Queries must return the resource with `200`. Only request, authorization, conflict, availability, and persistence failures use non-2xx transport status.

### Retention can invalidate observer cursors during a run

Automatic cycles are retained within a fixed bound. A slow observer may request a pruned cycle. The API must return an explicit retention boundary and canonical status rather than an empty page or a partial sequence presented as complete.

### Confirmation replay races with cancellation and completion

Duplicate browser posts, CLI retries, late SSE events, cancellation, and operation completion can arrive concurrently. Durable deduplication and compare-and-set terminal publication must decide the winner. Adapter memory, request order, and event arrival order are not authority.

### Public text can leak retained runtime material

Issue claims, blocker details, command diagnostics, and provider errors are untrusted. Sanitization must happen before durable public events and again at public projection boundaries. Bounds alone do not provide redaction. Tests need hostile secrets, terminal control sequences, HTML, long paths, and provider payloads.

### Automatic mode may appear enabled when proof is merely present

The API should expose `automatic_available`, proof freshness, and bounded reasons. It must not expose a simple configuration boolean as availability. The start operation repeats proof validation immediately before confirmation and mutation.

### Open implementation question

The app contract needs one deliberate choice during implementation: whether packet preparation uses a dedicated durable operation or a synchronous writer-fenced app call. Either is acceptable only if preparation remains runtime-free, idempotent, serializes with conflicting mutations, publishes strict private state atomically, and cannot start work. The smaller choice is preferable. Existing mutation and writer-fence behavior should decide it; no second operation framework should be introduced.

### Decision

Proceed with a separate repair app contract, repair-specific read resources, and the existing guarded durable operation path for confirmation, execution, cancellation, and observation. Bind confirmation after durable acceptance and before dispatch. Keep manual and automatic modes in the same operation and sprint protocol, require an explicit automatic opt-in, and expose operational lifecycle separately from semantic repair outcome.
