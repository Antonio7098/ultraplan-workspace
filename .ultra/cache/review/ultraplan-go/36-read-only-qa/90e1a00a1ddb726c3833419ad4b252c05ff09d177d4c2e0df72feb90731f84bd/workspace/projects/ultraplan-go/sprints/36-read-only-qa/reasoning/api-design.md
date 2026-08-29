> **Inputs Used:** `projects/ultraplan-go/sprints/36-read-only-qa/requirements.md`, `projects/ultraplan-go/sprints/36-read-only-qa/sprint-index.md`, `projects/ultraplan-go/sprints/36-read-only-qa/technical-handbook.md`, `projects/ultraplan-go/docs/ARCHITECTURE.md`, `projects/ultraplan-go/docs/PRD.md`, `projects/ultraplan-go/docs/TRD.md`, `studies/go-cli-study/reports/final/02-command-architecture.md`, `studies/go-cli-study/reports/final/03-dependency-injection.md`, `studies/go-cli-study/reports/final/04-configuration-management.md`, `studies/go-cli-study/reports/final/05-error-handling.md`, `studies/go-cli-study/reports/final/07-state-context.md`, `studies/go-cli-study/reports/final/08-concurrency.md`, `studies/go-cli-study/reports/final/09-terminal-ux.md`, `studies/go-cli-study/reports/final/10-logging-observability.md`, `studies/go-cli-study/reports/final/11-testing-strategy.md`, `studies/go-cli-study/reports/final/13-security.md`, `studies/go-cli-study/reports/final/14-performance.md`, `system/reasoning/api-design-reasoning-template.md`

# API design: read-only QA

This area covers the typed app API, CLI and JSON commands, versioned HTTP mapping, durable operation integration, cancellation, reconnect, compatibility, and public limits for Sprint 36. It does not decide the internal map algorithm or persistence layout beyond the authority that public projections must respect.

## Area Decisions

### One canonical app projection

`internal/app` will expose QA through typed use cases and DTOs that do not expose `internal/sprint` structs. CLI, TUI, web handlers, server-rendered HTML, and durable-run presentation will all consume these types. Product decisions stay in `internal/sprint`; transport mapping stays in each adapter.

The app API will provide these capabilities:

```go
type SprintQAUseCases interface {
	QAMap(context.Context, QAMapRequest) (QAMapResult, error)
	QAStatus(context.Context, QAStatusRequest) (QAStatusResult, error)
	QATheory(context.Context, QATheoryRequest) (QATheoryResult, error)
	QASynthesis(context.Context, QASynthesisRequest) (QASynthesisResult, error)
	RunQA(context.Context, QARunRequest, func(QAProgress)) (QARunResult, error)
	ResumeQA(context.Context, QAResumeRequest, func(QAProgress)) (QARunResult, error)
	CancelQA(context.Context, QACancelRequest) (QACancelResult, error)
	RecoverQA(context.Context, QARecoveryRequest) (QARecoveryResult, error)
}
```

Cancellation will not introduce a second QA-specific owner or registry. `CancelQA` validates that the run belongs to the requested sprint, delegates the state transition to `RunUseCases.CancelRun`, and returns the resulting QA projection. Run control remains the canonical cancellation API.

`QAMap` is runtime-free and non-persisting when called for `qa --dry-run`. It resolves current governed inputs, returns the normalized map, and invokes no investigator. `QAStatus`, `QATheory`, and `QASynthesis` are read-only projections. They must not open a runtime, acquire execution ownership, repair state, or create a missing state file. `RecoverQA` is the explicit runtime-free mutation that reconciles stale attempt ownership, interrupted state, and contained state pointers.

Every result carries the same canonical fields where applicable:

```text
schema_version
project
sprint
phase
status
fresh
attempt_id
run_id
input_fingerprint
implementation_fingerprint
map_fingerprint
coverage
limits
shards
theory_outcomes
synthesis
progress
blocker
cancellation
terminal_result
next_action
```

The QA phase status vocabulary is verdict-neutral:

```text
missing
mapped
queued
running
synthesizing
completed
blocked
cancelled
interrupted
stale
invalid
```

`completed` means the bounded read-only investigation and synthesis finished. It does not mean pass, does not promote a theory, and does not alter the Conformance Review verdict. Confirmed, refuted, invalid, inconclusive, blocked, cross-shard, and not-applicable theory outcomes remain data within that status.

### CLI and JSON contract

Sprint 36 adds one command family:

```text
ultraplan sprint <project> <sprint> qa --dry-run [--json]
ultraplan sprint <project> <sprint> qa [--shard <shard-id>] [--json]
ultraplan sprint <project> <sprint> qa resume [--shard <shard-id>] [--json]
ultraplan sprint <project> <sprint> qa status [--json]
ultraplan sprint <project> <sprint> qa cancel --run <run-id> [--json]
ultraplan sprint <project> <sprint> qa recover [--json]
```

`qa --shard` selects one shard from the current map. It cannot supply paths, theories, commands, prompts, budgets, environment, or runtime permissions. A completed current shard is not rerun by plain `qa --shard`; the result explains that the shard is current. A stale, interrupted, blocked, or incomplete shard can be run when the current fingerprints permit it. Full restart is deliberately absent from the public Sprint 36 contract because resume and fingerprint invalidation already define safe reuse, while destructive checkpoint removal adds little value.

Text and JSON invoke the same app method and apply the same error classification. JSON uses the existing command envelope and adds fields only under the QA result:

```json
{
  "schema_version": 1,
  "operation": "sprint.qa",
  "status": "completed",
  "result": {}
}
```

CLI exit behavior is fixed by domain meaning:

| Condition | Existing exit category |
| --- | --- |
| Valid dry-run, current status, or completed QA, regardless of theory outcomes | success |
| Invalid flags or incompatible modes | `ExitUsage` |
| Missing/stale prerequisite, invalid fingerprint/state/schema, permission denial, or blocked budget | `ExitValidation` |
| Durable acceptance, persistence, or runtime infrastructure failure | `ExitRuntime` |
| Explicit cancellation or interruption after preserving partial outcomes | `ExitPartial` |

A blocked QA result remains present in JSON even when the command returns a nonzero validation exit. Confirmed theories alone never cause a failure exit because Sprint 36 does not adjudicate or promote issues.

Human-facing output calls the existing analytical capability "Conformance Review." The existing `review` command, `review.md`, JSON operation name `sprint.review`, verdicts, and exits remain unchanged. `conformance-review` is a command alias normalized to the same review request and handler; it does not get a separate operation kind, state record, or JSON schema.

### Closed operation API

The shared operation vocabulary adds these closed kinds:

```text
qa-status
qa-dry-run
qa-start
qa-resume
qa-recover
```

`qa-start` and `qa-resume` are runtime-backed. They must pass through the Sprint 35 prepare, fingerprint revalidation, durable acceptance, owner claim, event, cancellation, and terminal paths before child work starts. `qa-status` and `qa-dry-run` are runtime-free and read-only. `qa-recover` is runtime-free but mutates only detailed QA state and its bounded summary through the sprint-owned reconciliation path.

The normalized operation request adds only `Shard`. Transport decoders must reject caller-supplied expected fingerprints, paths, commands, prompts, budgets, permission rules, environment, run ownership, attempt IDs, and theory outcomes. The server populates `ExpectedFingerprint` from preparation, as it does for existing operations.

Duplicate starts with the same consumed confirmation identity resolve to the same durable run. A resume creates a new durable run but reuses the same QA attempt only when governed, implementation, review, map, policy, and selected-shard fingerprints remain current. It never adopts a lost worker and never reruns completed valid shards.

### HTTP resources and operations

Versioned read routes expose the canonical snapshot:

```text
GET /api/v1/projects/{project}/sprints/{sprint}/qa
GET /api/v1/projects/{project}/sprints/{sprint}/qa/map
GET /api/v1/projects/{project}/sprints/{sprint}/qa/theories/{theory-id}
GET /api/v1/projects/{project}/sprints/{sprint}/qa/synthesis
```

These handlers call app use cases and return the common versioned success/error envelopes. Lists in the status resource are summary-bounded. Full map, theory, and synthesis detail use the focused routes rather than expanding the status response without limit.

Starts continue through the existing endpoints:

```text
POST /api/v1/operations/prepare
POST /api/v1/operations
GET /api/v1/runs/{run-id}
GET /api/v1/runs/{run-id}/events
DELETE /api/v1/runs/{run-id}
```

A QA preparation uses the existing operation request shape:

```json
{
  "operation": {
    "kind": "qa-start",
    "scope": {
      "project": "ultraplan-go",
      "sprint": "36-read-only-qa"
    },
    "options": {
      "shard": "qa-shard-v1-example"
    }
  }
}
```

The preparation response describes verification-state writes as the mutation class and target access as read-only. It reports effective limits, current map fingerprint, selected shard, runtime model source, prerequisites, durable refresh path, and the confirmation token. The browser cannot override the model or investigator policy in Sprint 36.

Successful runtime start is reported only after durable acceptance. HTTP disconnection or SSE loss stops observation only. Reconnect reads the current QA snapshot and resumes existing run events by durable sequence. A retention gap returns the existing typed gap/tombstone behavior plus links to the current QA status; it does not reconstruct product truth from an adapter buffer.

No QA-specific SSE registry, cancellation endpoint, browser operation history, or client-side workflow state is added. `operations.js` may enhance the server-rendered snapshot with the existing operation/run APIs, but the no-JavaScript page must show map freshness and coverage, shard progress, theory outcome counts, synthesis state, blocker, cancellation, terminal result, and next action.

### Request validation and security

The effective request order is:

```text
explicit typed request option > validated workspace QA configuration > product default
```

Only presence-aware options participate in precedence. Zero, empty string, and `false` are not treated as caller intent unless the transport recorded the field as present. The complete effective request is validated once in the app/product path so CLI and HTTP cannot disagree.

Public requests may choose only project, sprint, operation mode, and a map-owned shard ID. They cannot submit filesystem paths, shell text, executable names, arguments, environment names or values, prompts, policy, output locations, fingerprints, attempt IDs, or result content. Project and sprint resolution repeats lexical and symlink-aware containment. Unknown operation fields fail strict decoding.

The local browser keeps the existing loopback, Host, Origin, session, CSRF, body-limit, timeout, and short-lived confirmation controls. Read routes still return only bounded app DTOs. Mutating starts, resume, recovery, and cancellation require current authorization. Hostile theory and evidence text is treated as data and escaped in HTML.

Investigator command policy is not an HTTP or CLI extension point. Product code selects approved existing non-mutating checks as explicit executable/argv requests with a bounded environment. Shell wrappers, interpreter indirection, path escapes, Git commands, generated files, and writes to source, tests, governed inputs, verification code, or smoke harness content are denied before runtime dispatch. Unsupported permission enforcement produces a blocked QA result.

### Stable errors and product outcomes

Domain errors preserve causes and expose stable categories to the app boundary. Text, JSON, HTTP, TUI, and durable diagnostics render one public classification rather than matching error strings.

Public QA error codes are:

```text
qa.invalid_request
qa.not_ready
qa.stale_input
qa.invalid_fingerprint
qa.unsupported_schema
qa.invalid_state
qa.permission_denied
qa.budget_exhausted
qa.conflict
qa.not_found
qa.persistence_failed
qa.runtime_unavailable
```

The common JSON error object contains `code`, a bounded safe `message`, `retryable`, optional bounded `details`, and `next_action`. Wrapped Go messages, provider payloads, command output, absolute paths, secrets, and stack data stay out of public responses.

`blocked`, `inconclusive`, `refuted`, `invalid`, and `not_applicable` are product outcomes, not transport failures. HTTP returns the canonical resource with its product status when a bounded operation reached one of those outcomes. Malformed requests, stale confirmation, unsupported schema, authorization rejection, and persistence failure use the existing transport error envelope.

### Fixed public bounds

Sprint 36 does not accept per-request budget overrides. Workspace configuration may lower product defaults but may not exceed these maxima. Every map and status response reports the effective values and their source. A non-positive value fails validation. Exhaustion records a stop reason and blocks or completes the affected scope without widening it.

| Budget | Default | Maximum |
| --- | ---: | ---: |
| Changed paths in one map | 512 | 512 |
| Primary shards | 32 | 32 |
| Boundary shards | 8 | 8 |
| Changed paths per primary shard | 32 | 64 |
| Contextual paths per shard | 64 | 128 |
| Context expansions per shard | 2 | 4 |
| Paths per context expansion | 16 | 32 |
| Behavioral concerns per shard | 12 | 24 |
| Theories per shard | 12 | 24 |
| Investigator iterations per shard attempt | 4 | 8 |
| Approved commands per shard attempt | 8 | 16 |
| Runtime retries after the initial call | 1 | 2 |
| Concurrent investigators | 3 | 8 |
| Follow-up shards per synthesis | 4 | 4 |
| Command wall-clock duration | 5 minutes | 10 minutes |
| Shard wall-clock duration | 20 minutes | 30 minutes |
| Whole QA run wall-clock duration | 60 minutes | 90 minutes |
| Captured output per command | 256 KiB | 512 KiB |
| Stored investigator output per shard attempt | 1 MiB | 2 MiB |
| QA prompt bytes per investigator or synthesizer | 512 KiB | 1 MiB |
| Shards rendered in a summary | 40 | 40 |
| Theories rendered per focused shard page | 24 | 24 |
| Recent QA progress entries rendered | 100 | 200 |

Durable events retain Sprint 35's 16 KiB encoded-event limit and coalescing behavior. QA does not raise subscriber, replay, retention, or run-store quotas. Lists that can exceed a focused-route bound use stable ID ordering and bounded cursor pagination rather than truncation presented as completeness.

### Testing contract

One semantic fixture will be projected through app results, CLI text, CLI JSON, TUI, HTTP JSON, server-rendered HTML, and durable run detail. Tests compare canonical fields and meanings, not incidental formatting.

Required contract cases include:

- Strict request decoding, conflicting modes, unknown fields, oversized bodies, invalid IDs, and stable error codes.
- `review` and `conformance-review` invoking the same handler and retaining `sprint.review` JSON compatibility.
- Runtime-free dry-run and status paths proving no runtime construction, durable run acceptance, or state write.
- Durable acceptance before child work, duplicate-start deduplication, explicit cancellation, dropped observer delivery, reconnect, replay gap, session rotation, and restart recovery.
- App, CLI, TUI, HTML, HTTP, and run-detail parity for status, fingerprints, coverage, shard progress, outcomes, synthesis, blockers, cancellation, terminal result, and next action.
- Adversarial path, symlink, shell, Git, environment, permission, hostile-content, redaction, and output-bound cases.
- Fuzz tests for operation JSON decoding, CLI argument combinations, verification-scoped IDs, cursor parsing, and path/command policy classification.

Offline tests use fake runtimes, controlled clocks, in-memory streams, temporary workspaces, and narrow failure hooks. Gated dogfood is the only real-runtime API test and reports unavailable prerequisites as blocked.

## Trade-Offs

### Separate query routes instead of one expanding response

Focused map, theory, and synthesis resources add several handlers, but they keep status bounded and permit no-JavaScript and reconnect views to fetch only what they need. A single nested response was rejected because theory and evidence growth would turn status polling into an unbounded read.

### Closed operation kinds instead of a generic QA action string

Named `qa-status`, `qa-dry-run`, `qa-start`, `qa-resume`, and `qa-recover` kinds require compatibility fixtures when the API grows. That cost is useful. It keeps preparation, mutation class, runtime requirements, and authorization reviewable. A free-form action or plugin registry was rejected because Sprint 36 has a fixed operation set.

### Existing run cancellation instead of a QA cancellation registry

Delegating cancellation to run control means callers may need one extra status read to see the resulting QA state. In return, there is one cancellation identity, one fenced owner, and one reconnect contract across every surface. A QA-local cancel map was rejected because it could disagree with durable run truth.

### Configuration can lower limits but callers cannot raise them

This is less flexible than per-request budgets. It prevents browser and CLI callers from broadening investigator authority after confirmation, keeps fingerprints stable, and makes test fixtures reproducible. A later sprint can add carefully bounded overrides if real dogfood demonstrates the need.

### Product outcomes remain distinct from command and HTTP failure

Clients must inspect QA status and theory outcomes rather than infer quality from exit zero or HTTP 200. This is intentional. Treating confirmed theories as command failure would quietly turn investigators into issue adjudicators, while treating blocked work as success would hide missing evidence. The API returns the result and uses existing exit/error categories only for control-flow failures and explicit blockers.

### Durable resume creates a new run

Resume does not preserve one operational run ID across process loss. It preserves the QA attempt and completed valid shard records, then creates a new accepted run with fresh ownership. This keeps run-control fencing honest and avoids pretending that a dead worker was adopted.

### Rejected alternatives

- Adding QA to `PlanningStage` or the existing stage runtime map. Verification identity and lifecycle would be conflated with authored planning order.
- Exposing sprint domain structs directly from app. Adapter compatibility would become coupled to persistence migrations and internal fields.
- Letting HTTP call sprint services or parse CLI JSON. That would create a second workflow path and violate the web import boundary.
- Persisting detailed theories in `flow-state.json`. Status reads would become unbounded and planning refreshes could overwrite verification detail.
- Adding `qa.md`, issue resources, repair endpoints, generated-check options, or smoke suites. Those capabilities are outside Sprint 36.
- Binding execution lifetime to an HTTP request or SSE subscriber. Disconnect would become accidental cancellation and restart recovery would be impossible.
- Publishing raw runtime logs as progress. Provider text is neither a stable schema nor a safe public payload.
- Accepting arbitrary command strings or shell-compatible checks. Compatibility with shell features is not worth weakening the read-only boundary.
- Returning only error strings. Adapters need stable machine categories while logs retain wrapped causes.

## Evidence

### Governed project evidence

`projects/ultraplan-go/sprints/36-read-only-qa/requirements.md` requires typed shared QA use cases, versioned browser handlers, CLI/JSON commands, durable run acceptance, explicit cancellation, reconnect recovery, stable compatibility, finite limits, read-only permissions, and cross-surface agreement. It also fixes the authority split among detailed QA state, bounded flow summary, run-control records, and the source repository.

`projects/ultraplan-go/docs/ARCHITECTURE.md` assigns QA semantics to `internal/sprint`, shared use cases to `internal/app`, and presentation only to CLI, TUI, and web. Its Sprint 35 and Phase 5 sections establish that run control owns operational identity while product state retains QA meaning. This directly supports separate app DTOs and reuse of existing operation/run APIs.

`projects/ultraplan-go/docs/PRD.md` and `projects/ultraplan-go/docs/TRD.md` establish the human goal, compatibility vocabulary, versioned local HTTP boundary, guarded starts, durable observation, cancellation, and no-browser-authority rule. Their broader future command examples include smoke-as-QA, but the current sprint requirements explicitly exclude that option. This document follows the narrower sprint contract.

### Report evidence and project inference

The reports are comparative evidence, not proof that their designs fit UltraPlan unchanged.

- `studies/go-cli-study/reports/final/02-command-architecture.md` finds that mature command handlers parse, acquire dependencies, invoke shared behavior, and render results. Gh-cli's factory at `pkg/cmdutil/factory.go:16-43` and Helm's delegation from `pkg/cmd/install.go:132-145` support thin QA command adapters. The decision to use app-owned QA DTOs is the UltraPlan-specific inference.
- `studies/go-cli-study/reports/final/03-dependency-injection.md` shows manual composition and lazy function dependencies in gh-cli at `pkg/cmd/factory/default.go:26-46`, while warning about global service lookup. This supports keeping dry-run and status independent of runtime construction. It does not define QA lifecycle or DTO names.
- `studies/go-cli-study/reports/final/04-configuration-management.md` records explicit flag-presence handling in restic at `internal/global/global.go:139,147`, precedence restoration in chezmoi at `internal/cmd/config.go:2253-2287`, and merged validation in opencode at `internal/config/config.go:609-641`. This supports presence-aware options and one post-merge validation pass. The exact QA limits above are sprint decisions, not report findings.
- `studies/go-cli-study/reports/final/05-error-handling.md` supports cause-preserving typed errors and top-level exit mapping through rclone `fs/fserrors/error.go:22-192`, restic `internal/errors/fatal.go:10-53`, and gh-cli `internal/ghcmd/cmd.go:44-49,281-301`. Stable QA JSON codes are an UltraPlan contract added because strings are insufficient across adapters.
- `studies/go-cli-study/reports/final/07-state-context.md` distinguishes root context cancellation, process-local sessions, persistent state, and cross-process locks. Helm `pkg/cmd/install.go:333-347` and restic `internal/restic/lock.go:105,290-305` support cancellation propagation plus separately bounded cleanup. The decision that resume creates a new run while reusing current QA state follows UltraPlan's durable run authority.
- `studies/go-cli-study/reports/final/08-concurrency.md` supports bounded admission and timed shutdown through k9s `internal/pool.go:21,30,37` and opencode `cmd/root.go:261-279`; it identifies gh-cli's goroutine-per-item work at `pkg/cmd/extension/manager.go:196-206` as a scaling risk. This supports fixed concurrency and queue limits rather than caller-selected fan-out.
- `studies/go-cli-study/reports/final/09-terminal-ux.md` shows single-owner progress in restic at `internal/ui/termstatus/status.go:197-205` and non-TTY suppression in gh-cli at `pkg/iostreams/iostreams.go:514-516`. The report does not provide durable reconnect semantics. UltraPlan therefore keeps renderer selection separate from operation lifecycle and emits canonical progress through app/run DTOs.
- `studies/go-cli-study/reports/final/10-logging-observability.md` supports structured fields, stream separation, and central wrappers through Helm `internal/logging/logging.go:31-71`, k9s `internal/slogs/keys.go:6-231`, and restic `internal/backend/logger/log.go:22-77`. Since the report does not prove redaction or parity, Sprint 36 requires allowlisted, redacted product events before persistence and delivery.
- `studies/go-cli-study/reports/final/11-testing-strategy.md` supports command-level scenario tests, normalized output, call-recording process fakes, and injected failures through chezmoi `internal/cmd/main_test.go:64-174`, lazygit `pkg/commands/oscommands/fake_cmd_obj_runner.go:17-26`, and restic `internal/backend/mock/backend.go:14-26`. Shared semantic parity fixtures and decoder fuzzing are Sprint 36 decisions that close gaps the report identifies.
- `studies/go-cli-study/reports/final/13-security.md` supports explicit argv, permission gates, canonical path checks, and redaction through restic `internal/backend/shell_split.go:45-76`, opencode `internal/permission/permission.go:44-108`, and chezmoi `internal/cmd/gpgencryption.go:151-165`. This directly rejects arbitrary shell input and caller-owned investigator policy.
- `studies/go-cli-study/reports/final/14-performance.md` supports hard resource limits through age's fixed-size stream at `internal/stream/stream.go:20,195-219`, restic's bounded queue at `internal/archiver/file_saver.go:56-58`, and rclone's capped pool at `lib/pool/pool.go:17-24,52-53`. The exact path, shard, output, duration, and rendering values are deliberately frozen here for Sprint 36.

The key conclusion is that the API should expose product state and durable operation identity separately. A context cancels live work; it cannot provide restart recovery. A renderer shows progress; it cannot own cancellation. A successful investigation run records outcomes; it cannot supply a Conformance Review verdict.

## Risks

- The fixed limits may prove too small for a real repository. Exhaustion must stay visible and blocked rather than prompting an automatic increase. Dogfood should measure which limit was hit before a later sprint changes it.
- `RecoverQA` is a runtime-free mutation and could be confused with read-only status. Help, confirmation text, app mutation class, and authorization tests must state that it changes verification state but never the target repository.
- A current failed or blocked Conformance Review is still a valid diagnostic input, while missing or stale review evidence blocks QA. If readiness code conflates verdict acceptability with freshness, it may incorrectly prevent diagnostic QA or allow stale input.
- Operation success and QA completion are separate from theory outcomes. Adapters may accidentally color a confirmed theory as a failed run or present `completed` as a pass. Shared fixtures must freeze the vocabulary.
- Focused shard execution can produce misleading partial summaries if the API omits map coverage and synthesis freshness. Every focused result must retain parent map identity and explain whether synthesis is current.
- Stable IDs are scoped to the QA schema version, not global content identity. Public docs and JSON must not imply that an ID survives an unsupported major migration.
- Strict request DTOs can break clients if fields are renamed casually. Sprint 36 should add fields additively within schema version 1 and fail closed only on incompatible major versions or unknown request fields.
- The existing operation envelope may not currently carry every QA option or result link. Adding `Shard` and QA kinds must update CLI, app, web, JavaScript, and compatibility fixtures in one change.
- Recovery after atomic-write failure can leave run control terminal while QA state still reports an active attempt. The status API must report both facts and a recovery action rather than guessing success or rewriting state during a read.
- Redaction can drift if runtime diagnostics bypass the app projection. Persist and publish only product-owned bounded event fields; keep raw provider data outside public progress.
- Browser cancellation races with terminal completion. The canonical run snapshot decides the terminal outcome, and the QA projection must preserve completed shard outcomes even when cancellation loses that race.
- Server-rendered summaries can become expensive on large attempts. The fixed summary and focused-route limits must apply before template rendering, not after building an unbounded view model.
- The exact before-and-after target identity mechanism still depends on sprint architecture reasoning. The API requires both fingerprints and a drift result, but must not invent a Git-only guarantee for non-Git targets.
