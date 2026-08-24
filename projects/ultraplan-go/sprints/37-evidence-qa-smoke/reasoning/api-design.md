# API design reasoning: evidence-producing QA

> **Inputs Used:** `projects/ultraplan-go/sprints/37-evidence-qa-smoke/requirements.md`, `projects/ultraplan-go/sprints/37-evidence-qa-smoke/sprint-index.md`, `projects/ultraplan-go/sprints/37-evidence-qa-smoke/technical-handbook.md`, `projects/ultraplan-go/docs/ARCHITECTURE.md`, `projects/ultraplan-go/docs/PRD.md`, `projects/ultraplan-go/docs/TRD.md`, `system/reasoning/api-design-reasoning-template.md`, `studies/go-cli-study/reports/final/01-project-structure.md`, `studies/go-cli-study/reports/final/02-command-architecture.md`, `studies/go-cli-study/reports/final/03-dependency-injection.md`, `studies/go-cli-study/reports/final/04-configuration-management.md`, `studies/go-cli-study/reports/final/05-error-handling.md`, `studies/go-cli-study/reports/final/06-io-abstraction.md`, `studies/go-cli-study/reports/final/07-state-context.md`, `studies/go-cli-study/reports/final/08-concurrency.md`, `studies/go-cli-study/reports/final/09-terminal-ux.md`, `studies/go-cli-study/reports/final/10-logging-observability.md`, `studies/go-cli-study/reports/final/11-testing-strategy.md`, `studies/go-cli-study/reports/final/12-extensibility.md`, `studies/go-cli-study/reports/final/13-security.md`, `studies/go-cli-study/reports/final/14-performance.md`

> **Selected Area Template:** API Design, from `system/reasoning/api-design-reasoning-template.md`.

This area covers the typed application contract and its CLI, JSON, TUI, and local HTTP projections. It does not define isolation internals, evidence adjudication algorithms, or persistence implementation. Those remain owned by `internal/sprint`. The API must make those product decisions observable without letting any adapter become a second authority.

## Area Decisions

### One application operation, separate read projections

Add one typed `SprintQARequest` application operation for ordinary QA and focused shard execution. It is an asynchronous durable operation whenever it can start a runtime, investigator command, or writable workspace. The operation enters Sprint 35 run control before child work begins and returns a stable run correlation.

The request has this logical shape:

| Field | Type | Rule |
| --- | --- | --- |
| `project` | string | Required, resolved by the existing project rules. |
| `sprint` | string | Required, resolved beneath the selected project. |
| `shard_id` | optional string | Selects one current mapped shard. Mutually exclusive with `suite`. |
| `suite` | optional enum | Empty for empirical QA. The only Sprint 37 value is `smoke`. |
| `resume` | bool | Resumes compatible current verification state. Mutually exclusive with `restart`. |
| `restart` | bool | Starts a new attempt while retaining prior attempts and the last complete report. Mutually exclusive with `resume`. |
| `dry_run` | bool | Runs admission and scope calculation without writable child work. |
| `expected_fingerprint` | string | Required for guarded adapter starts and checked immediately before acceptance. |

Callers cannot supply command argv, environment, working directories, writable paths, timeouts, output caps, evidence limits, issue limits, or adjudication instructions. Product configuration and the frozen evidence plan own those values. This prevents a transport request from widening QA permissions.

Expose two app-owned read models rather than one unbounded result:

- `SprintQASummary` is the current bounded projection used by status lists, CLI text, TUI overview, browser HTML, and reconnect.
- `SprintQADetail` returns one immutable attempt with paged evidence, adjudication, and issues.

Both models are plain application DTOs. They contain no `internal/sprint` types, filesystem handles, process requests, run-control repository objects, CLI formatting, or HTML. `internal/web` depends only on these app contracts.

### Stable summary contract

`SprintQASummary` contains:

| Field | Meaning |
| --- | --- |
| `schema_version` | Version of the adapter-facing result, starting at `1`. |
| `project`, `sprint` | Resolved scope. |
| `current_attempt_id` | Current verification attempt, absent before the first attempt. |
| `canonical_attempt_id` | Attempt that produced the retained complete `qa.md`, which may differ after a failed rerun. |
| `run_id` | Current or most recent durable operation correlation, if one exists. |
| `input_fingerprint`, `implementation_fingerprint` | Current identity used for freshness decisions. |
| `fresh` | Product-derived freshness, never inferred by an adapter. |
| `attempt_status` | `missing`, `running`, `completed`, `failed`, `cancelled`, or `cleanup_uncertain`. |
| `assessment` | `incomplete`, `pass`, `pass_with_findings`, `fail`, or `blocked`. |
| `coverage` | Bounded map, shard, theory, and required-suite counts. |
| `evidence_counts` | Accepted, rejected, invalid, inconclusive, and truncated counts. |
| `issue_counts` | Counts by severity plus repair-eligible and regression-candidate counts. |
| `smoke` | Bounded suite status, containing-suite identity, external run identity, verdict, and canonical evidence flag. |
| `blockers` | Bounded typed blockers with code, safe message, and recovery action. |
| `qa_artifact`, `smoke_artifact` | Contained application artifact references, not arbitrary paths. |
| `next_action` | Product-derived action shared by every adapter. |

`attempt_status` describes execution facts. `assessment` describes the product's current QA conclusion. A cancelled attempt does not create a `cancelled` assessment, and a successfully completed command does not imply `pass`.

The deterministic assessment precedence is:

1. `incomplete` before a current admissible attempt has enough required evidence.
2. `blocked` when the Sprint 36 gate, current Conformance Review, identity, containment, cleanup, required evidence, or containing smoke suite cannot be proven.
3. `fail` when current valid evidence promotes a blocker or high issue, or the required containing smoke suite fails.
4. `pass_with_findings` when all mandatory evidence is current and sufficient but valid promoted lower-severity issues remain.
5. `pass` only when all mandatory current evidence passes, no promoted issue remains, and every required containing suite is canonical.

Cancellation, timeout, stale-writer rejection, and publication failure stay visible in attempt and durable-run status. They cannot erase or silently upgrade the last complete canonical assessment. QA also cannot alter the independent Conformance Review verdict.

### Detail, evidence, and artifact reads

Read APIs use immutable attempt IDs and opaque cursors. The app layer validates all pointers and digests before returning a record.

The local HTTP mapping is:

```text
GET /api/v1/projects/{project}/sprints/{sprint}/qa
GET /api/v1/projects/{project}/sprints/{sprint}/qa/attempts/{attempt_id}
GET /api/v1/projects/{project}/sprints/{sprint}/qa/attempts/{attempt_id}/evidence?cursor=...&limit=...
GET /api/v1/projects/{project}/sprints/{sprint}/qa/attempts/{attempt_id}/issues?cursor=...&limit=...
```

The current endpoint returns `SprintQASummary`. Attempt detail includes adjudication totals, root-cause groups, rejected-evidence reasons, exact evidence references, cleanup facts, and regression-candidate references. Evidence and issue endpoints return bounded immutable records plus `next_cursor` and `truncated`.

Generated patch content is not embedded in summary or list results. An evidence record may return a contained artifact reference for the existing bounded artifact-preview use case. The preview use case revalidates the reference beneath the selected attempt's `patches/` directory and applies its own byte cap. Callers never submit a filesystem path.

Request bodies keep the existing 64 KiB server limit and identifier limit. List endpoints default to 25 records and reject limits above 100. Each rendered untrusted text field is capped at 8 KiB, and every list response has a 512 KiB encoded ceiling. Truncation is explicit. Product evidence validity still uses the durable record, not the shortened presentation.

### Guarded starts and durable observation

Do not add a QA-specific HTTP mutation route or operation registry. Browser starts use the existing guarded operation contract:

```text
POST   /api/v1/operations/prepare
POST   /api/v1/operations
GET    /api/v1/runs/{run_id}
GET    /api/v1/runs/{run_id}/events
DELETE /api/v1/runs/{run_id}
```

The normalized operation kind is `sprint.qa`. The confirmation binds the project, sprint, shard or suite, resume/restart choice, dry-run state, current input fingerprint, mutation class, and expiry. Preparation returns the bounded scope, expected writes, model/runtime source, limits, prerequisites, and blocker state. Confirmation never authorizes a different request.

The start response returns `202 Accepted` with `run_id`, durable lifecycle, run URL, events URL, QA summary URL, and current replay cursor. Durable acceptance and owner claim happen before QA work. If acceptance fails, no child starts.

Reads remain available after browser session rotation because run and QA visibility is workspace-scoped. Start and cancellation require a current same-origin session, CSRF proof, request-bound confirmation where applicable, and the existing local authorization check. An SSE disconnect stops observation only.

Cancellation is an idempotent request against the existing run resource. It stops new scheduling, reaches the current fenced owner, and leaves cleanup and terminal arbitration to run control and the product operation. The HTTP response reports `requested`, `already_requested`, or `terminal`; it never predicts that cleanup succeeded.

A repeated browser start with the same unexpired confirmation token returns the already accepted run rather than creating a second attempt. A different confirmed start conflicts while the per-sprint mutation lease is held. CLI retries after an uncertain response inspect the workspace run list and QA summary before starting new work.

### CLI and JSON shape

The CLI remains a thin adapter over the same requests and results:

```text
ultraplan sprint <project> <sprint> qa [--shard <id> | --suite smoke] [--resume | --restart] [--dry-run] [--json]
ultraplan sprint <project> <sprint> qa status [--attempt <id>] [--json]
ultraplan sprint <project> <sprint> qa cancel --run <run-id> [--json]
```

`qa status` calls the read use case. `qa cancel` calls the shared durable cancellation use case. Neither reads verification files directly. Human output stays concise and sends diagnostics to stderr. JSON uses the existing command envelope and snake-case fields, with `schema_version: 1`, command identity, status, result, and typed safe error data. ANSI and progress output never enter JSON stdout.

Exit behavior follows existing classes:

| Outcome | Exit behavior |
| --- | --- |
| Current `pass` or `pass_with_findings` | Success. The assessment remains distinguishable in JSON and text. |
| `fail` | Validation/product failure. |
| Preflight or current assessment `blocked` | Validation or unavailable-prerequisite failure with recovery guidance. |
| User cancellation | Cancellation exit class. |
| Accepted asynchronous browser operation | HTTP `202`; terminal outcome comes from run inspection. |
| Invalid flags, unknown shard, or incompatible options | Usage error before acceptance. |
| State corruption, unsupported version, or path/digest failure | Filesystem/state failure, closed to execution. |

### Smoke compatibility uses the smoke authority

`qa --suite smoke` is not implemented by the ordinary writable-investigator path. The app adapter maps it to the existing typed smoke operation and `Service.RunSmoke` path. Both `smoke` and `qa --suite smoke` therefore use the same manifest discovery, authoring, selection, containing-suite logic, argv, environment, timeout, cancellation, cleanup, evidence validation, verdict, `smoke.md` publication, and external run identity.

The durable operation kind remains the existing smoke kind for both spellings. The request source may be recorded as bounded adapter metadata for diagnostics, but it cannot alter behavior. After the smoke operation validates its canonical result, sprint QA records only the validated suite link, identities, bounded adjudication facts, and canonical status. Raw harness output stays outside the UltraPlan workspace.

This choice preserves old clients and makes parity testable. A new QA-native smoke request, manifest parser, result normalizer, or cancellation path is rejected.

### Error contract

Keep the existing API error envelope and add stable QA error codes rather than exposing Go error strings or sprint internals. Required codes include:

```text
qa_invalid_request
qa_not_ready
qa_entry_gate_blocked
qa_stale_input
qa_unknown_shard
qa_attempt_conflict
qa_state_invalid
qa_state_unsupported
qa_evidence_invalid
qa_cleanup_uncertain
qa_confirmation_stale
qa_cancel_not_authorized
qa_cursor_invalid
```

Each error carries a safe message, retryability, current run or attempt reference when available, and one product-derived recovery action. Debug details, target absolute paths, command output, environment values, model prose, and raw external evidence are excluded. Adapters use `errors.Is` and `errors.As` against typed app errors, then map them to existing exit and HTTP classes.

Unknown JSON request fields are rejected. Unknown additive response fields are ignored by compatible clients. Unknown major persisted-state versions fail closed and appear as `qa_state_unsupported`; the API does not attempt an implicit migration.

## Trade-Offs

| Decision | Benefit | Cost | Rejected alternative |
| --- | --- | --- | --- |
| Separate summary and detail DTOs | Status and reconnect stay cheap while exact evidence remains inspectable. | More mapping code and schema fixtures. | One complete response was rejected because evidence, patches, issues, and output can exceed safe UI and JSON bounds. |
| Existing generic operation routes for starts and cancellation | Preserves confirmation, durable acceptance, fencing, replay, and authorization. | QA clients must follow links to QA state after run completion. | QA-specific operation and SSE endpoints were rejected because they would create a second lifecycle authority. |
| Immutable attempt reads with cursors | Stable audit views and bounded responses. | Clients need pagination and must distinguish current from canonical attempts. | Offset pagination over the current mutable pointer was rejected because reruns would reorder results. |
| Product-selected execution limits | Prevents callers from widening permissions or cost. | Advanced callers cannot tune every command per request. | Caller-supplied argv, environment, timeouts, or writable paths were rejected as a containment bypass. |
| Idempotency by confirmation token and mutation conflict | Browser retries do not duplicate accepted work. | CLI ambiguity requires run inspection rather than blind retry. | A new global idempotency-key service was rejected as unnecessary workflow infrastructure. |
| Existing smoke operation for both command spellings | Gives direct parity for protocol, run identity, cancellation, evidence, and `smoke.md`. | Smoke results require a small QA projection after validation. | A uniform QA-native smoke executor was rejected because normalization would duplicate smoke authority. |
| Stable machine codes plus separate hints | CLI, TUI, browser, and automation agree while humans get useful recovery text. | Code-to-hint mappings need compatibility tests. | Free-form error strings were rejected because adapters would classify failures differently. |
| Fixed presentation bounds | Prevents hostile evidence from exhausting terminals, JSON clients, or browser rendering. | A decisive excerpt may be shortened in a view. | Returning raw output was rejected. Full bounded evidence remains available through governed artifact references. |

Writable attempts remain sequential in Sprint 37. Read-only queries and observers may run concurrently, but the public request contract does not expose a writable concurrency option. This gives the first implementation a simpler lease, cleanup, and publication story. Parallel writable attempts can be considered only after isolation and process cleanup pass race and fault tests.

## Evidence

The boundary decision follows the project-structure report's thin-entrypoint and one-way dependency findings. Chezmoi and yq keep command code outside protected product logic, while restic defines interfaces on the owning side. That supports `internal/web -> internal/app -> internal/sprint` and keeps QA policy out of transport code. See `studies/go-cli-study/reports/final/01-project-structure.md`, sections "Thin CLI Entry Point," "Unidirectional Dependency Flow," and "Domain Interface/Implementation Split."

The command architecture report finds that factory-built, thin command handlers scale better than large `RunE` functions. Its Helm example delegates command setup to an action at `helm/pkg/cmd/install.go:132-145`, and gh-cli centralizes dependencies in `pkg/cmdutil/factory.go:16-43`. This supports one typed request with separate renderers rather than four adapter-specific QA workflows. See `studies/go-cli-study/reports/final/02-command-architecture.md`, "Thin-Delegate Pattern" and "Factory Function Command Creation."

Explicit construction and focused interfaces are supported by `studies/go-cli-study/reports/final/03-dependency-injection.md`, especially "Centralized Composition Root," "Constructor Injection with Factory Function," and the warning against context service lookup. The API therefore carries request-scoped cancellation and identity through context, but dependencies stay constructor-injected.

The request rejects caller-controlled execution mechanics because merged policy needs one validation point before side effects. `studies/go-cli-study/reports/final/04-configuration-management.md` documents post-load validation and explicit precedence, including Dive's `PostLoad()` and K9s schema validation. QA configuration is resolved before request validation, and the caller may narrow scope but not broaden policy.

Stable error codes and separate human hints follow the typed-error and user-versus-operational findings in `studies/go-cli-study/reports/final/05-error-handling.md`. Go-task's structured task errors, Restic's fatal classification, and gh-cli's exit mapping show why adapters should inspect types rather than parse text.

The process, filesystem, and output seams needed to test these contracts are grounded in `studies/go-cli-study/reports/final/06-io-abstraction.md`. Restic's filesystem/backend interfaces and gh-cli's test IO constructor show how to fault external behavior and capture adapter output without bypassing the application API.

The durable asynchronous shape and cancellation rules follow `studies/go-cli-study/reports/final/07-state-context.md`, especially root-context propagation and Restic's separate bounded cleanup context. The result separates cancellation request, work termination, cleanup certainty, and canonical assessment instead of collapsing them into one flag.

Sequential writable execution and bounded observation follow `studies/go-cli-study/reports/final/08-concurrency.md`. Its strongest examples localize goroutine launch sites, bound fan-out, and put timeouts around shutdown waits. The same report warns that unwaited goroutines and unbounded `WaitGroup.Wait()` calls hide cleanup failure.

Status and event delivery remain observational because terminal rendering can disappear, degrade for non-TTY output, or lose an observer. `studies/go-cli-study/reports/final/09-terminal-ux.md` documents interruptible progress and non-TTY fallback. `studies/go-cli-study/reports/final/10-logging-observability.md` supports structured correlation and strict stdout/stderr separation.

The compatibility test plan follows `studies/go-cli-study/reports/final/11-testing-strategy.md`. Testscript-style command scenarios, fault-capable fakes, and reviewed golden files fit help, JSON, `qa.md`, `smoke.md`, and HTML contracts. Field-level tests remain preferable for state transitions so snapshots do not freeze incidental implementation detail.

The smoke decision uses the adapter and versioned-metadata evidence in `studies/go-cli-study/reports/final/12-extensibility.md`, particularly Helm's subprocess runtime and metadata conversion. The relevant lesson is delegation to an existing authority. The registry and plugin material is a warning here, not permission to add a QA plugin system.

The guarded-request and path rules are backed by `studies/go-cli-study/reports/final/13-security.md`. Explicit argument arrays, private temporary directories, permission decisions, canonical path checks, and redaction are all mechanical boundaries. A request field or model instruction cannot substitute for them.

The page and payload limits follow the streaming and bounded-buffer findings in `studies/go-cli-study/reports/final/14-performance.md`. Fixed response caps and pagination prevent evidence volume from turning status inspection into unbounded memory use. The same report cautions against speculative pooling, so this API selects limits and streaming before low-level optimization.

The sprint requirements add constraints that the comparative reports do not cover: Sprint 36 admission evidence, target immutability, global-only adjudication, retained canonical reports after failed reruns, detailed verification state outside `flow-state.json`, smoke authority in the external harness, and Sprint 35 durable lifecycle ownership. `projects/ultraplan-go/docs/ARCHITECTURE.md`, `projects/ultraplan-go/docs/PRD.md`, and `projects/ultraplan-go/docs/TRD.md` confirm the app and web dependency direction, versioned local HTTP contract, loopback security policy, and Phase 5 sequencing.

## Risks

- A failed rerun creates two legitimate timelines: the current failed attempt and the last complete canonical report. Every DTO and view must label `current_attempt_id` and `canonical_attempt_id`; a single `attempt_id` field would be misleading.
- Confirmation-token replay must be atomic with durable acceptance. If the token is consumed before acceptance commits, a network retry could strand the request. If it is consumed after acceptance without deduplication, it could create duplicate work.
- Opaque cursors must bind the immutable attempt and record kind. Accepting a cursor under a different attempt risks cross-attempt evidence disclosure or inconsistent pagination.
- Fixed text caps can hide the line that explains a failure. Views must expose truncation and a governed artifact reference when retained evidence is valid. Truncated command output itself cannot become sufficient evidence unless the decisive condition remains within the retained bytes.
- `pass_with_findings` is safe only if severity, evidence sufficiency, containing suites, and blocker precedence come from product code. Adapters must never derive it from issue counts alone.
- Smoke parity can drift if the alias adds even a small request default before delegation. Tests must compare selected containing suite, explicit argv, cwd, environment, timeout, external run ID, cancellation, cleanup, evidence links, verdict, flow projection, and `smoke.md` for both spellings.
- Browser read visibility and mutation authorization are deliberately different. Handler tests must cover session rotation, observer restart, CSRF failure, confirmation expiry, replay gaps, and cancellation by a newly authorized session.
- Hostile evidence can contain ANSI controls, invalid UTF-8, Markdown, HTML, URLs, and very long lines. App DTOs should normalize invalid UTF-8 and retain plain text; each adapter must escape its own medium and suppress terminal control sequences.
- Additive JSON compatibility does not make persisted-state migration automatic. API schema version, verification-state version, and durable run-control version are separate concerns and must remain visibly separate in tests and diagnostics.
- Sequential writable execution may make large QA runs slow. That cost is accepted for Sprint 37 because parallel mutation would multiply workspace, descendant-process, storage, cleanup, and stale-writer risks before independent isolation has been proven.
- The exact product-level attempt, evidence, issue, patch, and response byte limits must be defined in one validated QA configuration and documented with an aggregate worst-case calculation. The API limits above only bound transport and rendering; they do not replace execution and persistence budgets.
