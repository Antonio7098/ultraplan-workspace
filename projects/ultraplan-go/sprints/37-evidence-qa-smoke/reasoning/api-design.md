# API Design: Evidence-producing QA and smoke integration

> **Inputs Used:** `projects/ultraplan-go/project-index.md`, `projects/ultraplan-go/roadmap.md`, `projects/ultraplan-go/docs/ARCHITECTURE.md`, `projects/ultraplan-go/docs/PRD.md`, `projects/ultraplan-go/docs/TRD.md`, `projects/ultraplan-go/sprints/37-evidence-qa-smoke/requirements.md`, `projects/ultraplan-go/sprints/37-evidence-qa-smoke/code-context.md`, `projects/ultraplan-go/sprints/37-evidence-qa-smoke/sprint-index.md`, `projects/ultraplan-go/sprints/37-evidence-qa-smoke/technical-handbook.md`, `system/reasoning/api-design-reasoning-template.md`, `studies/go-cli-study/reports/final/01-project-structure.md`, `studies/go-cli-study/reports/final/02-command-architecture.md`, `studies/go-cli-study/reports/final/04-configuration-management.md`, `studies/go-cli-study/reports/final/05-error-handling.md`, `studies/go-cli-study/reports/final/07-state-context.md`, `studies/go-cli-study/reports/final/09-terminal-ux.md`, `studies/go-cli-study/reports/final/10-logging-observability.md`, `studies/go-cli-study/reports/final/11-testing-strategy.md`, `studies/go-cli-study/reports/final/12-extensibility.md`, `studies/go-cli-study/reports/final/13-security.md`, `studies/go-cli-study/reports/final/14-performance.md`, `../ultraplan-go/internal/app/sprint_usecases.go`, `../ultraplan-go/internal/app/operations.go`, `../ultraplan-go/internal/app/operation_runner.go`, `../ultraplan-go/internal/app/sprint_commands.go`, `../ultraplan-go/internal/app/run_usecases.go`, `../ultraplan-go/internal/web/operation_handlers.go`, `../ultraplan-go/internal/web/qa_handlers.go`, `../ultraplan-go/internal/sprint/domain.go`, `../ultraplan-go/internal/sprint/qa.go`, `../ultraplan-go/internal/sprint/qa_state.go`, `../ultraplan-go/internal/sprint/qa_types.go`, `../ultraplan-go/internal/sprint/smoke.go`, `../ultraplan-go/internal/sprint/smoke_types.go`, `../ultraplan-go/internal/runcontrol/interfaces.go`, `../ultraplan-go/internal/runcontrol/model.go`, `../ultraplan-go/internal/runcontrol/lifecycle.go`

This area covers the local application, CLI/JSON, loopback HTTP, and durable-run contracts for Sprint 37. The API audience is internal adapters and stable local automation. It is not a public, partner, remote-worker, or multi-tenant API.

## Area Decisions

### 1. Keep one product API and two authorities

`internal/app` remains the only product-facing boundary used by CLI, TUI, and browser adapters. It returns bounded product facts and does not expose `internal/sprint` records, filesystem stores, process requests, runtime requests, or external harness response types.

The two existing authorities stay separate:

| Concern | Authority | App projection |
| --- | --- | --- |
| QA map, shard attempts, evidence, adjudication, issues, assessment, and `qa.md` | `internal/sprint` verification state | `QAResult` and focused QA query results |
| Acceptance, run identity, ownership, events, liveness, cancellation, and terminal arbitration | `internal/runcontrol` | Existing `RunUseCases` snapshots and events |
| Raw smoke runs, stdout/stderr, test artifacts, and harness issues | Manifest-selected external harness | Validated bounded links and identities only |
| Compatible smoke verdict, `smoke.md`, and flow projection | Existing `Service.RunSmoke` path | Existing smoke result plus a QA smoke-suite summary |

Progress events describe execution. They never count as adjudicated evidence and never set assessment. A browser reconnect first loads the durable run snapshot and current QA status, then resumes events from the durable cursor. It does not reconstruct product state from SSE.

### 2. Add one suite selector without creating a second command family

Add `Suite string` to `app.QARequest` and pass it through `OperationRequest.Suite`, which already exists. The closed values are:

| Value | Meaning |
| --- | --- |
| Empty | Normal mapped QA. Sprint 37 may perform isolated evidence-producing investigation. Existing callers retain their current request shape. |
| `smoke` | Run the canonical smoke executor through QA. |

No other suite name is accepted. Suite names are case-sensitive stable identifiers. The API does not accept a smoke level, test, executable, argument vector, working directory, environment, timeout, evidence root, model, or parallelism for `qa --suite smoke`. Product and manifest configuration own those values.

`Suite: "smoke"` is valid for QA dry-run and start only. It is mutually exclusive with a focused shard. QA recovery, status, cancellation, and evidence queries infer the current suite or target run and reject a supplied suite. `qa resume --suite smoke` is rejected because the smoke harness does not have a compatible resumable unit. The safe action after interruption is a new smoke-suite run, while the prior durable run and external evidence remain inspectable.

Do not add operation kinds such as `qa-smoke-start`. Keep `qa-dry-run` and `qa-start`, with `Suite` included in normalization, confirmation, governed-input fingerprinting, durable target metadata, and event correlation. The existing `smoke-start` operation remains unchanged for compatibility.

For `Suite: "smoke"`, the sprint-owned adapter calls the existing `Service.RunSmoke` preparation, discovery, selection, invocation, validation, verdict, `smoke.md`, flow-state, and roadmap-reconciliation path. It may translate `SmokeProgress` and the validated result into QA summaries, but it must not duplicate any smoke rule. The QA invocation has one durable `qa-start` run and must not create a nested durable `smoke-start` run.

### 3. Extend the app boundary with bounded focused queries

Keep the existing `QAUseCases` mutations and `QAQueries` methods. Add the following queries:

```go
type QAQueries interface {
    QAMap(context.Context, QARequest) (QAResult, error)
    QAStatus(context.Context, QARequest) (QAResult, error)
    QAShard(context.Context, QARequest) (QAShardResult, error)
    QATheory(context.Context, QARequest) (QATheoryResult, error)
    QASynthesis(context.Context, QARequest) (QASynthesisResult, error)
    QAEvidence(context.Context, QARequest) (QAEvidenceResult, error)
    QAAdjudication(context.Context, QARequest) (QAAdjudicationResult, error)
    QAIssues(context.Context, QARequest) (QAIssuePage, error)
    QAIssue(context.Context, QARequest) (QAIssueResult, error)
    QAAssessment(context.Context, QARequest) (QAAssessmentResult, error)
    QASmokeSuite(context.Context, QARequest) (QASmokeSuiteResult, error)
}
```

Add `Evidence`, `Issue`, `After`, and `Limit` to `QARequest`. Identifiers select only records reachable from the current validated state pointer. They are not caller-provided paths. `After` is an opaque cursor from the previous response. The default page is 50 records and the hard maximum is 200. A non-positive limit selects the default; a value over 200 is rejected rather than silently widened.

Mutation methods reject `Evidence`, `Issue`, `After`, and `Limit`. Query methods reject fields they do not use. This keeps strict request validation and avoids one permissive request bag whose ignored fields differ by adapter.

`QAStatus` remains the compact cross-surface summary. It gains optional additive fields:

```text
assessment
canonical_report
evidence_counts
adjudication_summary
issue_counts
promoted_issues
regression_candidate_count
smoke_suite
current_failure
```

`promoted_issues` is a bounded summary list, not the issue database. `current_failure` is distinct from `canonical_report`: a failed or interrupted new attempt can be visible while the last complete report remains available. Detailed evidence, issue, patch preview, rejection, and root-cause facts come from focused queries.

### 4. Return summaries, not persistence records

All new result types include `schema_version`, `project`, `sprint`, `attempt_id`, and the current governed-input and implementation fingerprints where the fact depends on freshness. They use app-owned summary types.

`QAEvidenceResult` returns one evidence record selected by evidence ID. Its bounded contract includes:

- evidence ID, kind, shard ID, theory IDs, and expectation references
- frozen-plan digest and bounded plan summary
- exact executable identity and redacted argument summary, never a shell command string
- command status, exit code, timeout, cancellation, output byte counts, truncation, and output digests
- target, workspace, map, generated-patch, external-evidence, and policy identities
- containment and cleanup status as independent facts
- repeatability status and supporting attempt references
- adjudication disposition, rejection codes, and supporting issue IDs
- a generated-patch reference with digest and byte count, plus a sanitized preview capped by the existing app preview limit

The response does not return full stdout/stderr, inherited environment values, secrets, arbitrary filesystem paths, raw model payloads, or raw external harness JSON. A patch preview is loaded by evidence identity through the app use case; the caller cannot supply a path. `truncated` and original byte count remain explicit.

`QAAdjudicationResult` returns the current adjudication ID and digest, status, considered and rejected evidence counts, bounded rejection summaries, deterministic root-cause groups, requested follow-up, promotion reasons, assessment-input codes, and next action. Any model-proposed grouping is labeled as a proposal until product rules validate it.

`QAIssueSummary` contains the issue ID, root-cause group ID, title, severity, repair eligibility, regression-candidate classification, exact evidence IDs, affected expectation references, promotion reason, freshness, and next action. It has no assignee, schedule, mutable workflow state, remote synchronization field, or repair command.

`QAAssessmentResult` uses the existing closed assessment vocabulary: `incomplete`, `blocked`, `fail`, `not_applicable`, `pass_with_findings`, and `pass`. It also returns deterministic basis codes and current references to Conformance Review, QA evidence, adjudication, blockers, and required containing smoke evidence. Product code computes the value. Investigator output, model prose, a check exit code, a smoke harness issue, and progress events cannot set or upgrade it. QA cannot alter the independent Conformance Review verdict.

`QASmokeSuiteResult` is a bounded QA projection of the existing smoke result. It includes execution status, verdict, diagnostic-only status, selected containing scope, harness/protocol/run identity, current review fingerprint, counts, validated external evidence links, cleanup status, `smoke.md` reference, freshness, and next action. Raw artifacts remain external.

### 5. Preserve stable JSON while versioning durable state honestly

Keep the CLI envelope and app response `schema_version: 1`. Existing fields keep their names, meanings, zero-value behavior, and JSON types. New fields are additive and use `omitempty` where absence means that Sprint 37 evidence has not been produced. Existing clients that ignore unknown response fields continue to work.

Do not use that compatibility choice for private durable state. The new verification graph changes what a valid passing record means and old strict readers reject unknown fields. Therefore:

- publish new `verification/state.json` records as schema version 2
- read schema version 1 through an explicit compatibility decoder and project missing evidence, adjudication, issue, and assessment facts as unavailable or incomplete
- never infer a Sprint 37 pass from a version 1 record
- write version 2 only; do not silently rewrite historical attempt records during a read
- reject unknown major versions and invalid digests with typed recovery guidance
- retain existing `qa-v1` attempt, map, shard, theory, challenge, and synthesis IDs
- use a distinct deterministic namespace for new evidence, patch, adjudication, issue, and assessment IDs, without claiming global content identity

`flow-state.json` receives only additive bounded QA fields: assessment, counts, freshness, next action, and contained state/report pointers or digests. Detailed records do not enter flow state.

### 6. Keep CLI behavior explicit and scriptable

The accepted command forms are:

```text
ultraplan sprint <project> <sprint> qa [--shard <id>] [--json]
ultraplan sprint <project> <sprint> qa --dry-run [--json]
ultraplan sprint <project> <sprint> qa --suite smoke [--json]
ultraplan sprint <project> <sprint> qa --suite smoke --dry-run [--json]
ultraplan sprint <project> <sprint> qa status [--json]
ultraplan sprint <project> <sprint> qa resume [--shard <id>] [--json]
ultraplan sprint <project> <sprint> qa cancel --run <run-id> [--json]
ultraplan sprint <project> <sprint> qa recover [--json]
```

`--suite` takes exactly one value and is valid only on run or dry-run. Existing argument order, action words, `--shard`, `--run`, and `--json` stay valid. Normal QA starts isolated evidence-producing work only after its admission checks pass; help and text must stop describing every run as read-only. Dry-run remains runtime-free and write-free.

The JSON envelope remains:

```json
{
  "schema_version": 1,
  "operation": "sprint.qa",
  "status": "ok",
  "result": {}
}
```

On error, the envelope adds the existing stable error object. Text progress stays on stderr. Stable text and JSON results stay on stdout. CLI exit behavior follows the existing classes:

| Exit | Meaning for Sprint 37 QA |
| --- | --- |
| `0` | Query or dry-run succeeded, cancellation was accepted or already requested, or execution completed with `pass`, `pass_with_findings`, or `not_applicable`. |
| `2` | Invalid flag combination, identifier, suite, or request shape. |
| `5` | Admission, freshness, evidence, adjudication, or assessment is blocked or failed by product validation. |
| `6` | Runtime, process, or required persistence is unavailable. |
| `7` | Reserved existing cancellation class where the outer command itself reports cancellation. |
| `8` | Interrupted, timed out, partially completed, or cleanup-uncertain work. |

A promoted issue is not by itself a command transport error. The deterministic assessment and its exit mapping decide acceptance. Status and focused read queries return zero when they successfully report a blocked or failed state.

### 7. Add read routes and reuse guarded operation routes

Preserve all existing GET routes:

```text
GET /api/v1/projects/{project}/sprints/{sprint}/qa
GET /api/v1/projects/{project}/sprints/{sprint}/qa/map
GET /api/v1/projects/{project}/sprints/{sprint}/qa/shards/{shard}
GET /api/v1/projects/{project}/sprints/{sprint}/qa/theories/{theory}
GET /api/v1/projects/{project}/sprints/{sprint}/qa/synthesis
```

Add these GET routes:

```text
GET /api/v1/projects/{project}/sprints/{sprint}/qa/evidence/{evidence}
GET /api/v1/projects/{project}/sprints/{sprint}/qa/adjudication
GET /api/v1/projects/{project}/sprints/{sprint}/qa/issues?after=<cursor>&limit=<n>
GET /api/v1/projects/{project}/sprints/{sprint}/qa/issues/{issue}
GET /api/v1/projects/{project}/sprints/{sprint}/qa/assessment
GET /api/v1/projects/{project}/sprints/{sprint}/qa/suites/smoke
```

Do not add a direct `POST .../qa` endpoint. Browser starts still use the existing prepare, confirmation, operation start, durable run detail/events, and cancellation routes. `mapOperationRequest` accepts `options.suite: "smoke"` only for `qa-start` and `qa-dry-run`; it accepts `options.shard` only for normal `qa-start` and `qa-resume`. It continues to reject caller-controlled runtime and smoke details.

The canonical request and confirmation fingerprint include suite and shard. A prepared smoke-suite confirmation displays the existing smoke dry-run selection, manifest-owned harness, external evidence roots, expected duration/cost class, mutation scope, and the fact that raw evidence remains external. A stale fingerprint requires preparation again.

GET handlers call only app queries. They use current server read policy but do not embed authorization rules in app DTOs. Mutations retain same-origin, CSRF, request-bound confirmation, current authorization, body limit, and strict one-value JSON decoding. Cancellation also verifies that the durable run target belongs to the requested project, sprint, and QA operation. Browser disconnection cancels observation only.

### 8. Make retry and recovery semantics observable

Read queries and recovery are naturally idempotent. Cancellation is idempotent: repeated authorized requests return the same durable run plus `requested: false` after the first accepted request, not an error.

Operation start uses the existing durable acceptance contract. Retrying delivery of one already accepted confirmation returns the existing run. Consuming a new confirmation creates a new operational run even when governed inputs match. The semantic QA attempt may remain the same, but its operational attempt ID and fencing generation differ.

Normal QA resume targets only the current semantic attempt, reuses completed valid evidence, and schedules no work whose current validated record is complete. It cannot adopt a prior worker. Stale writers fail with conflict. Smoke-suite retries create a new external smoke run through `RunSmoke`; they do not overwrite raw prior harness evidence.

Recovery is runtime-free. It validates pointers, digests, writer ownership, terminal run state, retained records, cleanup truth, and last-complete report references. It may mark work interrupted, cancelled, stale, invalid, or cleanup-uncertain. It never infers success, promotes an issue, reruns a command, adopts a worker, or deletes evidence.

### 9. Use one stable error shape and fail closed

Extend the current typed QA categories with specific stable codes for isolation, containment, cleanup uncertainty, evidence validation, adjudication, unsupported suite, and missing focused records. App and HTTP adapters project them through the existing safe error fields:

```text
code
category
operation
component
message
cause
guidance
retryable
```

The code is stable and machine-readable. `message` and `guidance` are bounded and safe. `cause` is redacted display text, not an error chain dump. No field contains raw command output, model output, environment values, secret-bearing paths, or unrestricted harness diagnostics.

HTTP mapping is deterministic:

| Status | Use |
| --- | --- |
| `400` | Malformed JSON, invalid identifier/cursor/limit, unknown suite, or forbidden option combination. |
| `404` | A focused current evidence or issue ID does not exist. |
| `409` | Stale confirmation, current-state mismatch, active-owner conflict, or stale writer. |
| `422` | A well-formed operation is blocked by governed prerequisites, evidence validation, or adjudication policy. |
| `503` | Required runtime, process, run-control, or persistence capability is unavailable. |
| `500` | Unexpected internal failure after safe classification. |

A successfully loaded QA state with assessment `blocked` or `fail` is still a `200` query response. Product outcome is data, not an HTTP transport failure.

### 10. Correlate operations without turning telemetry into evidence

Every writable QA and smoke-suite durable run records these bounded correlation fields where applicable:

```text
run_id
operational_attempt_id
fencing_generation
project
sprint
operation_kind
suite
semantic_attempt_id
map_id
shard_id
evidence_id
adjudication_id
```

Events use stable phase and category values with bounded safe summaries. They may report isolation preflight, workspace creation, plan freeze, command execution, cleanup, evidence publication, adjudication, assessment, and canonical publication. They do not include raw generated content or use issue IDs as metric labels.

Metrics use bounded labels such as operation kind, suite, phase, terminal result, error category, evidence disposition, and assessment. Counts cover accepted runs, duration, cancellation latency, cleanup uncertainty, rejected evidence, promoted issues, dropped event delivery, and reconciliation. Logs and events aid diagnosis; only validated verification records count as evidence.

### 11. Freeze the contract with cross-surface tests

Required API tests are:

- app contract tests for every result projection, bound, hostile string, omitted optional field, current-pointer restriction, and no raw persistence exposure
- request table tests for every valid and invalid suite, shard, resume, recovery, query, and cancellation combination
- CLI help and parser tests for old forms plus `--suite smoke`, including text/JSON separation and exit classes
- JSON golden tests proving old fields and types are unchanged and new fields are additive
- HTTP route, method, strict JSON, body-limit, cursor, object-scope, same-origin, CSRF, confirmation, stale-request, and cancellation tests
- durable acceptance tests proving one QA smoke-suite run, no nested operation, idempotent redelivery, fencing, replay, reconnect, and terminal arbitration
- parity tests that feed the same fixtures to `smoke` and `qa --suite smoke` and compare discovery, selection, argv, environment, timeout, cancellation, cleanup, evidence validation, verdict, `smoke.md`, flow projection, and external run identity
- schema tests for version 1 read compatibility, version 2 writes, unknown-major rejection, immutable historical records, invalid digests, partial publication, and prior valid report preservation
- assessment tests proving model output, investigator claims, command failures, diagnostic smoke, and narrow smoke cannot directly pass, promote, or replace required evidence
- one shared fixture inspected through app DTOs, CLI text/JSON, TUI, browser HTML/JSON, durable run detail, `qa.md`, `smoke.md`, and verification state

## Trade-Offs

| Decision | Benefit | Cost and rejected alternative |
| --- | --- | --- |
| Add `Suite` to existing QA and operation requests | Small additive change; confirmation and durable correlation already understand operation options. | A dedicated `qa-smoke` operation would be easier to branch on, but it would create more aliases, run-target types, and parity paths without adding product meaning. |
| Keep empty suite as normal QA | Existing callers and JSON requests remain valid. | An explicit `suite: evidence` would be clearer in new code, but requiring it would break shipped callers. Renderers may display `default` or `evidence` while the wire value stays empty. |
| Reuse `Service.RunSmoke` | Existing manifest, containing-suite, safety, verdict, publication, and compatibility behavior remain authoritative. | The adapter must avoid nested durable acceptance and duplicate mutation locking. A second QA smoke executor was rejected because parity tests cannot make two rule implementations one authority. |
| Focused bounded query methods | Adapters can inspect evidence and issues without reading private files or loading the entire attempt. | More DTO and mapping code. Returning `QAState` directly was rejected because it leaks persistence layout, hostile content, and future schema changes. |
| Keep app/JSON schema version 1, move durable state to version 2 | Local automation gets additive compatibility while strict private state does not pretend a semantic change is old-schema compatible. | Two version policies require clear tests. Keeping durable state at version 1 was rejected because old strict readers cannot validate new records and could overstate completeness. |
| Current-pointer-only focused queries | Stale or orphaned records cannot be mistaken for current evidence. | Historical inspection needs durable run detail or a later explicitly scoped history API. Caller-selected attempt paths were rejected because they weaken containment and freshness. |
| Page issue lists at 50 with a 200 maximum | Predictable response size and stable browser/TUI behavior. | Clients may need another request. Unbounded status payloads were rejected because issue and evidence counts scale independently from screen size. |
| Separate current failure from last complete report | Operators see the failed attempt without losing the last validated `qa.md`. | Consumers must understand two facts. Replacing the canonical pointer at run start was rejected because a failed publication would erase useful accepted evidence. |
| Reject smoke-suite resume | The API does not claim resumability that the external harness cannot provide. | Operators rerun and receive a new external run ID. Pretending a fresh harness invocation resumes old child work was rejected as misleading. |
| Keep status reads successful for failed assessments | HTTP and CLI query success remain separate from product acceptance. | Scripts must inspect assessment. Turning every failed assessment into a transport error was rejected because it prevents reliable inspection and conflates state with request handling. |
| Use durable events only for operations | Replay, cancellation, and reconnect remain consistent across surfaces. | Product state still requires a second status read. Treating SSE or logs as evidence was rejected because delivery can drop and neither has adjudication authority. |

## Evidence

### Governed product evidence

- `projects/ultraplan-go/sprints/37-evidence-qa-smoke/requirements.md` requires one adapter-independent contract, bounded projections, strict browser requests, durable run correlation, cancellation/recovery, schema compatibility, and identical smoke authority through `smoke` and `qa --suite smoke`. It also states that raw smoke evidence remains external and that product code alone promotes issues and derives assessment.
- `projects/ultraplan-go/sprints/37-evidence-qa-smoke/sprint-index.md` selects API Design specifically to resolve additive app, CLI/JSON, browser request, durable-run, cancellation, recovery, and compatibility contracts. It excludes a second smoke implementation, general issue tracking, alternate persistence, remote workers, and Git mutation.
- `projects/ultraplan-go/docs/ARCHITECTURE.md` assigns QA, adjudication, assessment, canonical reports, and smoke compatibility to `internal/sprint`; shared interface contracts and composition to `internal/app`; HTTP transport and security to `internal/web`; and operational run authority to `internal/runcontrol`. This rules out web-owned QA state and direct verification-file APIs.
- `projects/ultraplan-go/docs/PRD.md` requires `qa [--dry-run|--shard <id>|--suite smoke]`, stable CLI/JSON automation, one product core across local surfaces, observer-independent runs, and preservation of `smoke` compatibility.
- `projects/ultraplan-go/docs/TRD.md` requires stable machine-readable errors, strict loopback HTTP security, guarded confirmation, durable event replay, typed app use cases, explicit cancellation, and no browser or operational shadow authority.

### Current implementation evidence

- `../ultraplan-go/internal/app/sprint_usecases.go` currently exposes `QARequest` with project, sprint, shard, theory, and run ID plus bounded map/status/shard/theory/synthesis projections. It has no typed evidence, adjudication, issue, assessment, canonical report, or smoke-suite query. Adding fields and focused methods preserves its adapter-independent role.
- `../ultraplan-go/internal/app/operations.go` already has `OperationRequest.Suite`, canonical request fingerprinting, guarded preparation, typed operation errors, and QA operation kinds. Its QA validator currently rejects suite, so allowing only `smoke` is a narrow additive change.
- `../ultraplan-go/internal/app/sprint_commands.go` currently supports QA run, map, status, resume, cancel, recover, shard focus, JSON envelopes, durable acceptance, writer fencing, and stable exit classification. The parser lacks `--suite`, and the renderer still labels QA read-only.
- `../ultraplan-go/internal/web/operation_handlers.go` already strictly decodes one JSON value and maps `options.suite`, but its QA branch rejects suite. `../ultraplan-go/internal/web/qa_handlers.go` contains read-only app-backed QA routes and no direct `internal/sprint` dependency. These are the correct extension points.
- `../ultraplan-go/internal/app/run_usecases.go` already provides sanitized durable run list, snapshot, cursor events, idempotent cancellation, and health. A QA-specific event or cancellation registry would duplicate this contract.
- `../ultraplan-go/internal/sprint/qa_types.go` has schema version 1, deterministic `qa-v1` IDs, bounded settings, writer correlation, summary evidence, and typed errors. It lacks frozen evidence plans, patch/adjudication/issue/assessment identities, and its strict state semantics justify a version 2 durable root rather than pretending all additions are old-schema compatible.
- `../ultraplan-go/internal/sprint/smoke.go` makes `Service.RunSmoke` the path that records attempts, runs preparation/discovery/selection/process/evidence validation, writes `smoke.md`, updates flow state, and may reconcile roadmap delivery. This is why QA must adapt that method rather than call lower smoke helpers.
- `../ultraplan-go/internal/sprint/domain.go` already defines the assessment vocabulary `incomplete`, `blocked`, `fail`, `not_applicable`, `pass_with_findings`, and `pass`. Reusing it avoids a conflicting QA verdict vocabulary.
- `../ultraplan-go/internal/runcontrol/interfaces.go`, `../ultraplan-go/internal/runcontrol/model.go`, and `../ultraplan-go/internal/runcontrol/lifecycle.go` provide the existing durable acceptance, event order, ownership, fencing, cancellation, reconciliation, and terminal authority that Sprint 37 must retain.

### Report findings and sprint-specific inference

- Report finding: `studies/go-cli-study/reports/final/01-project-structure.md` and `02-command-architecture.md` favor thin adapters and shared execution behavior, and warn about separate CLI/TUI command systems. Sprint-specific inference: all QA and smoke-suite behavior belongs behind app use cases, not in renderers or a parallel browser workflow.
- Report finding: `studies/go-cli-study/reports/final/04-configuration-management.md` supports explicit precedence, validation after merge, and fixed effective invocation inputs. Sprint-specific inference: callers select only the closed suite or map-owned shard; runtime, process, and harness details remain product-owned and enter the confirmation fingerprint.
- Report finding: `studies/go-cli-study/reports/final/05-error-handling.md` supports typed classification, exit mapping, and retained partial failures. Sprint-specific inference: QA needs stable product error codes and separate current-failure versus last-complete-report fields.
- Report finding: `studies/go-cli-study/reports/final/07-state-context.md` shows that work cancellation and cleanup require separate bounded lifetimes. Sprint-specific inference: cancellation and cleanup status must be independent response facts, and cleanup uncertainty cannot map to success.
- Report finding: `studies/go-cli-study/reports/final/09-terminal-ux.md` and `10-logging-observability.md` separate stable automation output, live progress, and diagnostics. Sprint-specific inference: JSON/text remain on stdout, progress remains on stderr or durable events, and neither logs nor events become evidence.
- Report finding: `studies/go-cli-study/reports/final/11-testing-strategy.md` supports deterministic fixtures, goldens, fakes, fault injection, and command-level tests for different risks. Sprint-specific inference: one parity fixture must traverse app, CLI, browser, durable run, Markdown, and state contracts.
- Report finding: `studies/go-cli-study/reports/final/12-extensibility.md` favors narrow extension contracts and warns that subprocess isolation alone does not impose resource or filesystem limits. Sprint-specific inference: `Suite` is a closed product selector, not a plugin or caller-defined executor.
- Report finding: `studies/go-cli-study/reports/final/13-security.md` supports explicit argv, permission boundaries, private temporary data, and redaction. Sprint-specific inference: no QA API accepts executable details or arbitrary evidence paths, and focused records are selected by validated IDs.
- Report finding: `studies/go-cli-study/reports/final/14-performance.md` favors bounded workers, buffers, and streaming. Sprint-specific inference: status returns counts and bounded summaries, focused collections are paged, and patch/output bodies are never unbounded inline data.

## Risks

| Risk | Consequence | Required control |
| --- | --- | --- |
| Calling `RunSmoke` beneath an already held sprint mutation lock deadlocks or duplicates side effects. | `qa --suite smoke` hangs or publishes twice. | The sprint-owned adapter must enter the same top-level smoke path with one lock/transaction boundary. Tests must prove one durable run, one harness invocation, one `smoke.md` commit, and one roadmap reconciliation. |
| Existing strict consumers reject additive JSON fields despite the schema staying at version 1. | Local integrations break even though field names were preserved. | Freeze representative old-client decoding fixtures and document that response decoders must ignore unknown fields. If a shipped strict consumer exists, gate each field behind its established compatibility mechanism. |
| State version 1 compatibility is mistaken for current evidence. | Old summary evidence could produce a false pass. | The v1 decoder must emit unavailable adjudication and an `incomplete` or `blocked` assessment. Only validated v2 records can satisfy Sprint 37 evidence requirements. |
| A focused ID can reach an orphaned or stale record. | UI presents historical evidence as current. | Resolve IDs only through the current validated state graph and include attempt/fingerprint facts in every focused result. Historical inspection remains visibly separate. |
| Bounded summaries omit the reason an issue was promoted or evidence was rejected. | Operators cannot audit assessment. | Always retain stable reason codes and exact evidence IDs in summaries; provide focused details and explicit truncation rather than silent omission. |
| Patch preview or hostile evidence leaks secrets or executable markup. | Terminal or browser output exposes sensitive content or enables injection. | Project through app sanitization, cap bytes, preserve truncation, escape HTML, never render embedded HTML, and never return caller-selected paths or raw environment/output data. |
| A browser retry starts duplicate work after losing the response. | Multiple isolated or smoke runs consume resources. | Durable acceptance precedes execution. Redelivery of the same accepted confirmation returns its run; a new confirmation visibly creates a new run. |
| Assessment and command exit code drift. | CLI, JSON, browser, and `qa.md` disagree on acceptance. | Centralize assessment and exit classification in app/product code and use one cross-surface fixture. Renderers must not infer pass from phase, command exit, issue count, or smoke verdict alone. |
| QA smoke projection drifts from the compatibility command. | `smoke` and `qa --suite smoke` report different authority. | Build the QA projection from the validated `RunSmoke` result and compare all meaningful fields in parity tests. Do not reimplement selection or verdict logic in app or web. |
| Smoke-suite interruption is called resumable even though the harness restarts. | Operators misunderstand evidence identity. | Reject smoke-suite resume, retain the interrupted durable run, and label the next action as a new run with a new external run ID. |
| New error categories expose internal details inconsistently. | Automation cannot classify failures and diagnostics may leak. | Map every product category once into stable app code, CLI exit, HTTP status, safe message, and recovery guidance. Test redaction with hostile causes. |
| Issue pagination order changes between requests. | Cursors skip or duplicate issues. | Order immutable current issue records by deterministic ID or publication sequence and bind cursors to the current attempt digest. Return conflict when the current attempt changes. |
| Suite is accepted by generic operation decoding for an unsupported QA action. | Confirmation and execution canonicalize different requests. | Share one validation table between preparation and execution, include suite in canonical serialization, and test every action/option pair through CLI and HTTP. |

The decision is to proceed with the additive contracts above. The main cost is more app DTO and projection code, but that cost keeps persistence, smoke, durable operation, and presentation authorities from collapsing into one unstable API.
