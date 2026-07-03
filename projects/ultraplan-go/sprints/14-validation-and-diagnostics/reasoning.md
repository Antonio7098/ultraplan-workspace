# Sprint Reasoning: Validation Command, Diagnostics, and JSON Stability

> Project: `ultraplan-go`
> Sprint: `14-validation-and-diagnostics`
> Output: `projects/ultraplan-go/sprints/14-validation-and-diagnostics/reasoning.md`
> **Inputs Used:** `projects/ultraplan-go/sprints/14-validation-and-diagnostics/requirements.md`, `projects/ultraplan-go/project-index.md`, `projects/ultraplan-go/docs/PRD.md`, `projects/ultraplan-go/docs/TRD.md`, `projects/ultraplan-go/docs/ARCHITECTURE.md`, `projects/ultraplan-go/sprints/14-validation-and-diagnostics/sprint-index.md`, `projects/ultraplan-go/sprints/14-validation-and-diagnostics/technical-handbook.md`, `templates/sprint-reasoning.md`

This document decides. It synthesizes selected context, handbook evidence, area-specific reasoning, and contracts into final sprint decisions.

It does not replace `sprint-index.md`, `technical-handbook.md`, or any future `reasoning/*.md` artifact.

## Sprint Purpose

- **Goal:** Deliver inspectable, automatable, and safe validation/status/health surfaces: `ultraplan study <study> validate`, `validate --json`, `study <study> status --json`, and `health --json` with stable schemas, redacted diagnostics, and deterministic tests.
- **Non-Goals:** No target/sprint workflow commands, browser UI, TUI, hosted service, metrics exporter, alerting integration, local API server, runtime adapter changes, execution/repair/retry/fallback/source cloning from validation, schema migrations, default real OpenCode smoke tests, report template or prompt semantic changes except diagnostics exposure, unsafe raw payload persistence, or general diagnostics framework.
- **Depends On:** Prior study-side capabilities for Markdown applicability, report validation, run-state/status, runtime health, single analysis/synthesis, run-all, durable run-loop, summary generation, and code-reference diagnostics; PRD/TRD stable JSON and diagnostics requirements.

## Selected Context And Pre-Reasoning Artifacts

| Artifact | Path | How It Was Used |
| --- | --- | --- |
| Sprint Requirements | `projects/ultraplan-go/sprints/14-validation-and-diagnostics/requirements.md` | Authoritative sprint contract for outputs, acceptance criteria, constraints, dependencies, review expectations, and required verification commands. |
| Project Index | `projects/ultraplan-go/project-index.md` | Confirmed project scope, target repository, active contract pool, selected study report catalog, and lack of formal prior decision artifacts. |
| PRD | `projects/ultraplan-go/docs/PRD.md` | Grounded product expectations: runtime success is not product success, explicit status/errors, first-class validation, machine-readable status/validate/inspect output, safe diagnostics, and deferred target/sprint workflows. |
| TRD | `projects/ultraplan-go/docs/TRD.md` | Grounded technical requirements: module-driven ownership, exit codes, JSON output envelope expectations, source applicability, report validation, run-state schema/version handling, agentwrap boundary, redaction, testing, and no real OpenCode requirement for normal tests. |
| Architecture | `projects/ultraplan-go/docs/ARCHITECTURE.md` | Grounded package boundaries: `internal/study` owns study validation/status semantics, `internal/app` wires and renders commands, `internal/platform/runtime` stays product-agnostic, and global technical-layer packages are avoided. |
| Sprint Index | `projects/ultraplan-go/sprints/14-validation-and-diagnostics/sprint-index.md` | Selected applicable contracts, evidence reports, excluded context, prior carry-forward constraints, and review protocols. |
| Technical Handbook | `projects/ultraplan-go/sprints/14-validation-and-diagnostics/technical-handbook.md` | Provided study evidence and source references for thin command wiring, stable JSON output, typed diagnostics, redaction, fake-first tests, bounded inspection, and scope control. |

## Area-Specific Reasoning Inputs

The sprint index selected an Architecture area reasoning artifact at `projects/ultraplan-go/sprints/14-validation-and-diagnostics/reasoning/architecture.md`, but no `reasoning/*.md` files are present for this sprint. There are therefore no area-specific conclusions to summarize or override.

| Area | Reasoning Document | Key Conclusion | Evidence Basis | Impact On Final Decision |
| --- | --- | --- | --- | --- |
| Architecture | Not present | No completed area-specific reasoning artifact exists. Final sprint decisions are made directly from requirements, project docs, sprint index, and technical handbook evidence. | Requirements, PRD, TRD, ARCHITECTURE, sprint index, and technical handbook. | Decisions explicitly encode module ownership, runtime separation, JSON rendering placement, durable-state handling, and no forbidden global packages so `plan.md` does not need to reopen architecture. |

## Sprint Technical Handbook Summary

- **Relevant Patterns:** Thin command wiring with domain-owned behavior; stable machine output through explicit rendering paths; typed diagnostics and exit-code mapping; redaction before logs/text/JSON; fake-first command tests; bounded local artifact/state inspection.
- **Important Trade-Offs:** Centralized JSON envelope improves schema consistency but must not absorb product semantics; typed validation findings improve automation but add model surface; safe redaction can reduce debugging detail and must preserve safe cause summaries; runtime-free validation/status misses runtime liveness that belongs in `health --json`.
- **Warnings / Anti-Patterns:** Do not place business logic in `RunE`; do not mix ANSI/progress/human diagnostics into JSON; do not hide config/runtime globals behind validation/status; do not panic or emit string-only parse failures; do not leak secrets, raw payloads, prompt bodies, Markdown bodies, or full stderr; do not recursively scan source repositories.
- **Evidence Confidence:** High for command wiring, IO/test seams, exit codes, redaction, bounded inspection, and schema tests because the handbook cites repeated patterns across mature Go CLIs. Medium for configuration/observability/philosophy because those patterns are broader and must be applied within UltraPlan's existing constraints.

## Contracts Applied

| Contract / Requirement ID | Constraint | Decision Impact | Expected Evidence |
| --- | --- | --- | --- |
| Architecture; AC45-46; C63-66 | Study behavior stays in `internal/study`; app parses/renders/maps exit codes; runtime stays product-agnostic; no global validation/report/diagnostics/status packages. | Validation orchestration and result model are study-owned; JSON envelope helper is app-owned and payload-agnostic. | Import/package review, no forbidden packages, command handlers remain thin. |
| Errors; AC33-34; TRD 7.2 | Validation findings must be actionable and map to exit code `5`; causes must remain wrapped and classifiable. | Introduce typed validation checks/status/severity/guidance and command-level exit mapping. | Validation and command tests assert exit code `0`/`5`, safe observed detail, guidance, and non-validation errors remain distinct. |
| Configuration; AC37, AC41; TRD 6, 19, 22 | Health/config/runtime diagnostics must use merged config and redact sensitive values. | `health --json` returns redacted runtime/config/workspace summaries through existing health boundary. | Health JSON tests seed secrets and assert absence; config/runtime diagnostics are redacted. |
| Observability; AC37-38, AC43 | Status and health need truthful, structured diagnostics, run/task metadata, retry/fallback, validation, and unknown usage/cost handling. | `status --json` exposes stable summaries from durable/local state without invoking runtime; unknown usage/cost remains unknown/null/known flags. | Status JSON tests assert task counts, lock diagnostics, run metadata, validation summaries, retry/fallback metadata, and unknown usage/cost. |
| Security; AC41; C69-70 | All visible diagnostics must pass redaction; paths should prefer workspace-relative output. | Renderers only receive safe/redacted diagnostic strings or sanitized payloads. | Redaction tests across validation, status, and health text/JSON. |
| Testing; AC42, AC47-49; C73 | Tests must be deterministic, offline, fake-first, and schema-aware. | Add unit and command tests for validation/status/health JSON and build/race verification. | `go test ./...`, `go test -race ./...`, `go build ./cmd/ultraplan`; command tests parse JSON. |
| Documentation and CLI Surface; AC36-40, AC44 | Help output and text/JSON behavior must be stable and documented by tests. | Help text documents `validate`, `validate --json`, and `status --json`; text mode remains deterministic and backward-compatible. | Help/output tests and review of command help. |
| LLM Runtime; AC35; C66 | Validation/status cannot invoke agentwrap/OpenCode/network/subprocess; health may use existing runtime health boundary. | Separate runtime-free validation/status service paths from runtime health JSON path. | Tests use fakes or nil runtime seams to assert validate/status do not invoke runtime. |
| LLM Evaluation / Cost / Safety; AC38; C68 | Unknown usage/cost must not become zero and unsafe runtime payloads must not be exposed. | Status result uses explicit known/unknown/null representation for usage/cost summaries. | Status tests assert unknown values are not serialized as numeric zero. |
| Workflows; AC38; C74 | Status must expose durable run-loop metadata and preserve prior decisions. | Status JSON includes task state summaries, lock diagnostics, retry/fallback/cancellation/run metadata, and validation summaries. | Status JSON tests and sprint review checks. |
| Performance; C72 | Validation/status must stay bounded to known artifacts and avoid recursive source repo scans. | Validation inspects workspace/study structure, discovered sources/dimensions, expected paths, reports, summary, and run state only. | Tests and review confirm no source cloning/runtime/recursive repo scan behavior. |
| Persistence And Migrations; AC31, C71 | Durable state must be read strictly and malformed/unsupported schema versions diagnosed, not migrated. | Validation/status parse `run-state.json` defensively and emit safe findings for missing/malformed/unsupported versions. | Tests for missing, malformed, and unsupported `run-state.json`. |

## Repos Studied / Source Evidence Used

| Source / Repo / Report | Concrete Reference | Relevant Finding | Why It Matters For This Sprint | Used In Decision(s) |
| --- | --- | --- | --- | --- |
| `01-project-structure` | `studies/go-cli-study/reports/final/01-project-structure.md`; chezmoi `main.go:26-34`, gh-cli `cmd/gh/main.go:6`, restic `cmd/restic/main.go:37-114` | Mature CLIs keep entrypoints thin and delegate to internal packages. | Supports keeping command handlers as wiring/rendering and validation semantics in `internal/study`. | Decisions 1, 2, 7 |
| `02-command-architecture` | helm `pkg/cmd/install.go:132-145`, gh-cli `pkg/cmdutil/factory.go:16-43`, chezmoi `internal/cmd/config.go:1833-1845` | Command factories and `RunE` bodies should delegate instead of owning business rules. | Prevents `study_commands.go` from duplicating validators or status logic. | Decisions 1, 2, 5, 7 |
| `03-dependency-injection` | gh-cli `pkg/cmdutil/factory.go:16-43`, go-task `executor.go:22-24`, restic `internal/restic/repository.go:18-66` | Manual composition and seams enable fake-first command tests. | Supports nil/fake runtime seams proving validation/status are runtime-free and health uses controlled fakes. | Decisions 2, 5, 7 |
| `04-configuration-management` | dive `cmd/dive/cli/internal/options/analysis.go:48-53`, opencode `internal/config/config.go:609-641`, restic `internal/global/global.go:139,147` | Effective config should be validated after merge and diagnostics should avoid precedence surprises. | Informs health JSON config/runtime summaries and redacted config diagnostics. | Decision 5 |
| `05-error-handling` | gh-cli `internal/ghcmd/cmd.go:44-49`, go-task `errors/errors.go:47-50`, restic `internal/errors/fatal.go:10` | Typed/sentinel errors and exit-code mapping are script-friendly. | Supports validation exit code `5`, structured failures, and preserving non-validation error categories. | Decisions 1, 2, 7 |
| `06-io-abstraction` | gh-cli `pkg/iostreams/iostreams.go:551-568`, go-task `executor.go:553-564`, restic `internal/ui/mock.go:10-53` | Output should be injectable and test-capturable. | Supports centralized JSON rendering and command tests that parse captured stdout. | Decisions 3, 7 |
| `07-state-context` | gh-cli `internal/ghcmd/cmd.go:142`, restic `cmd/restic/cleanup.go:24-38`, rclone `fs/config.go:793-802` | Runtime-bound work needs context; local inspection commands can remain state-focused. | Supports runtime-free validation/status and runtime-bound health separation. | Decisions 2, 4, 5 |
| `09-terminal-ux` | chezmoi `internal/cmd/prompt.go:20-256`, gh-cli `internal/prompter/`, rclone `lib/terminal/terminal.go:82-86` | Machine output must avoid decorations, progress, and interactive assumptions. | Supports ANSI-free JSON and deterministic concise text output. | Decisions 3, 7 |
| `10-logging-observability` | helm `internal/logging/logging.go:31-66`, k9s `internal/slogs/keys.go:6-231`, rclone `fs/accounting/prometheus.go:78-108` | Structured diagnostics should be separated from stdout data and redacted. | Supports stable status/health JSON diagnostics and safe result payloads. | Decisions 3, 4, 5, 6 |
| `11-testing-strategy` | chezmoi `internal/cmd/main_test.go:64-174`, gh-cli `acceptance/acceptance_test.go:26-29`, helm `internal/test/test.go:43` | Command-level tests, golden/schema assertions, and fakes catch output regressions. | Drives JSON schema field assertions, exit-code tests, redaction tests, and verification commands. | Decision 7 |
| `13-security` | restic `internal/options/secret_string.go:15-20`, helm `pkg/registry/transport.go:37-41`, gh-cli `status.go:332-338` | Redaction must happen before visible diagnostics/logs. | Drives redaction policy updates and tests across validation/status/health. | Decision 6 |
| `14-performance` | gh-cli `pkg/cmdutil/factory.go:27-42`, yq `pkg/yqlib/stream_evaluator.go:78-113`, gdu `pkg/analyze/parallel.go:13` | Avoid eager runtime setup and unbounded scans. | Supports bounded validation over known artifacts and lazy runtime-free status. | Decisions 2, 4, 7 |
| `15-philosophy` | lazygit `VISION.md:97`, gh-cli `docs/project-layout.md:68-69`, opencode `internal/pubsub/broker.go:10-19` | Scope control and deliberate complexity beat broad frameworks. | Supports rejecting a general diagnostics framework and keeping helpers narrow. | Decisions 1, 3, 6 |

## Trade-Off And Debt Analysis

### Accepted Trade-Offs

| Trade-Off | Benefit | Cost / Constraint Accepted | Why Acceptable Now | Revisit Trigger |
| --- | --- | --- | --- | --- |
| Study-owned validation model instead of app-owned command structs | Preserves module ownership and lets validation/status share study semantics. | `internal/study/domain.go` gains additional public-ish result types. | Stable JSON and text need typed fields; study owns the artifacts being validated. | Multiple non-study modules need identical validation semantics. |
| App-owned JSON envelope helper instead of per-command envelope duplication | Consistent schema version, command, workspace, status, timestamp, and ANSI-free output. | Risk of over-centralizing diagnostics if helper grows payload logic. | Helper will remain payload-agnostic and only renders envelopes. | Helper starts importing `study` or `platform/runtime`, or embeds command-specific fields. |
| Typed findings instead of simple error strings | Enables stable JSON fields, exit code `5`, severity/status, guidance, and reviewable tests. | More model surface to maintain. | The sprint explicitly requires stable, parseable JSON and actionable diagnostics. | Findings become too generic or duplicate existing report validators. |
| Redacted safe summaries instead of raw runtime/config/report details | Prevents secret and payload leakage. | Debugging may require inspecting local files manually. | Requirements forbid leaking secrets, raw payload bytes, prompt bodies, Markdown bodies, and full stderr. | A future explicit debug-retention feature adds safe opt-in controls. |
| Runtime-free validation/status and runtime-bound health | Keeps validation/status deterministic, fast, offline, and scriptable. | Runtime liveness failures are not discovered by validation/status. | Health is the correct surface for runtime checks; validation/status inspect artifacts/state. | Product requires status to include live runtime/process inspection. |
| Strict JSON schema assertions | Protects automation callers and catches regressions. | Test maintenance increases when schemas intentionally evolve. | Versioned schemas are a core sprint output. | Schema evolution policy is introduced and requires compatibility tests. |

### Potential Technical Debt

| Debt / Shortcut | Why It Might Accrue | Current Mitigation | Owner / Follow-Up |
| --- | --- | --- | --- |
| Validation result type may grow as more artifact classes are added. | Stable JSON encourages adding fields for every future diagnostic. | Keep model minimal: check, status, severity, path, expected, observed safe detail, guidance, schema/version, and summary counts. | Future validation expansion sprint or decision log. |
| JSON envelope versioning may be too coarse if all commands share one version. | Health/status/validation payloads may evolve at different rates. | Use shared envelope schema version and command-specific result schema/version fields where needed. | Future schema compatibility decision. |
| Redaction of nested structures may become inconsistent. | New nested payloads can bypass string-level redaction if added casually. | Centralize safe diagnostic construction and apply existing redaction policy before render. | Security review after additional JSON surfaces. |
| Status validation summary may duplicate some validation-command checks. | Status needs a concise validation summary but must not become full validation execution if costly. | Status should include existing state/artifact summaries and lightweight validation summaries bounded to known artifacts. | Revisit if status latency grows or duplicates validation service excessively. |
| Unsupported durable schema handling may initially be diagnostic-only. | Users may need migrations later. | Requirements explicitly reject silent migration; clear guidance will recommend future migration/remediation. | Future migration sprint. |

### Future Considerations

| Consideration | Deferred Until | Reason Deferred | What Should Be Preserved Now |
| --- | --- | --- | --- |
| Schema migration commands | A migration sprint is planned and indexed. | This sprint only diagnoses unsupported durable versions. | Strict version checks and actionable guidance. |
| Repair from validation failures | A repair/execution sprint explicitly scopes it. | Validation must not invoke runtime, repair, prompts, or synthesis. | Findings should include enough safe guidance to support future repair prompts. |
| Live status dashboards, TUI, metrics exporter, API server | Product scope changes. | Current sprint is CLI-only. | Stable JSON result payloads that future surfaces can consume. |
| Additional runtime adapters | Runtime roadmap changes. | Sprint must not change OpenCode supervision or add adapters. | Health JSON should use generic platform runtime health summaries. |
| Broader diagnostics framework | Multiple independent modules need common diagnostics. | Current need is limited to validation, status, and health surfaces. | Narrow redaction/output helpers and typed domain findings. |

## Final Decisions

### Decision 1: Study-Owned Validation Result Model

- **Decision:** Add or extend typed validation result structures in `internal/study/domain.go` with stable fields for schema/result version, overall status, summary counts, check name/ID, status, severity, artifact path, expected condition, observed safe detail, guidance, source/dimension/run-state context when applicable, and redacted diagnostics.
- **Rationale:** The sprint requires stable JSON, deterministic text, actionable failures, inapplicable artifacts, run-state diagnostics, redaction, and exit-code mapping. Those semantics belong to `internal/study` because they describe study artifacts and durable state, not CLI rendering.
- **Study / Source Grounding:** `technical-handbook.md` thin command wiring and typed diagnostics patterns; `01-project-structure`; `02-command-architecture`; `05-error-handling`; `ARCHITECTURE.md` module ownership; TRD validation result mapping and report validator requirements.
- **Trade-Offs Accepted:** More typed model surface is accepted to avoid string-only failures and unstable JSON. The model must stay minimal and not become a general diagnostics framework.
- **Technical Debt / Future Impact:** Future validation additions may pressure the type; preserve extensibility with stable optional fields and result schema/version rather than ad hoc command structs.
- **Alternatives Rejected:** App-local JSON structs are rejected because they duplicate study semantics and risk drift from validation behavior. Plain error strings are rejected because they cannot satisfy stable JSON, guidance, severity, path, inapplicable status, or schema tests.
- **Contracts Satisfied:** Architecture, Errors, Observability, Security, Testing, CLI Surface, Persistence And Migrations; requirements AC31-34, AC36, AC41-45, C63-74.
- **Evidence Required:** `internal/study/validation_command_test.go` covers valid, missing, invalid, malformed frontmatter, inapplicable pairs, run-state diagnostics, unsupported/malformed durable state, and secret-free diagnostics; command tests assert JSON fields and exit code `5` for validation failures.

### Decision 2: Runtime-Free Study Validation Command

- **Decision:** Implement `ultraplan study <study> validate` as thin CLI wiring in `internal/app/study_commands.go` that calls `internal/study` validation orchestration in `internal/study/validation_command.go`. The service validates study structure, source discovery, dimension files, expected applicable per-source reports, final reports, `summary.csv`, and optional `studies/<study>/.ultraplan/run-state.json`. It reports inapplicable Markdown source/dimension pairs as `inapplicable` or skipped, not failures, and never invokes runtime, network, source cloning, prompt generation, synthesis, repair, or subprocess paths.
- **Rationale:** PRD/TRD principles require runtime success to be insufficient without artifact validation, but this command is an inspection surface, not an execution or repair surface. Keeping it runtime-free makes it deterministic, offline, and safe for automation.
- **Study / Source Grounding:** `technical-handbook.md` fake-first inspection, bounded inspection, state/context separation, and performance warnings; `07-state-context`; `14-performance`; TRD 9A applicability filtering, TRD 15 report validators, TRD 13 run-state handling.
- **Trade-Offs Accepted:** Validation will not discover live runtime setup failures; users must run `ultraplan health --json` for that. The command only inspects known workspace/study artifacts and durable state.
- **Technical Debt / Future Impact:** Future repair flows should reuse safe findings but must be introduced as separate runtime-invoking commands or flags, not hidden in validation.
- **Alternatives Rejected:** Running agentwrap/OpenCode to repair or regenerate missing outputs is rejected by sprint non-goals and safety constraints. Recursively scanning source repositories for inferred outputs is rejected because it violates bounded inspection and performance constraints.
- **Contracts Satisfied:** Architecture, Errors, Observability, Security, Testing, Performance, Persistence And Migrations, CLI Surface; requirements AC31-35, AC40-45, C63-75.
- **Evidence Required:** Command tests prove runtime fakes are not invoked; validation tests cover inapplicable Markdown pairs, missing/invalid reports, malformed frontmatter, malformed/unsupported state, and deterministic text/JSON output.

### Decision 3: Stable JSON Envelope In `internal/app`

- **Decision:** Add `internal/app/json_output.go` as a payload-agnostic JSON rendering helper for `schema_version`, `command`, `workspace`, `status`, `generated_at`, and `result`. The helper emits only JSON to stdout for JSON mode, without ANSI, progress text, or mixed human diagnostics. Command-specific result schema/version fields remain in their owning payloads.
- **Rationale:** `validate --json`, `status --json`, and `health --json` need a shared stable envelope, while product-specific payloads must not move into a generic diagnostics package. `internal/app` is the right layer because it owns command rendering and output mode selection.
- **Study / Source Grounding:** `technical-handbook.md` stable machine output pattern; `06-io-abstraction`; `09-terminal-ux`; `10-logging-observability`; TRD 7.3 JSON output requirements; sprint index selected CLI Surface and Documentation contracts.
- **Trade-Offs Accepted:** Centralizing envelope rendering adds a shared helper, but constraining it to envelope-only behavior avoids framework creep.
- **Technical Debt / Future Impact:** If payload schema versions diverge, command payloads must carry their own `schema_version` or `result_schema_version` without changing envelope semantics unexpectedly.
- **Alternatives Rejected:** Per-command duplicated envelopes are rejected because schema drift would be likely across three sprint surfaces. A global `internal/diagnostics` or `internal/status` framework is rejected because the requirements explicitly forbid broad diagnostics/status packages and the evidence favors scope control.
- **Contracts Satisfied:** Architecture, Observability, Security, Testing, Documentation, CLI Surface; requirements AC36-40, AC42-44, AC46, C67-70.
- **Evidence Required:** Command-level tests parse JSON for validation/status/health success and failure cases, assert envelope fields, assert no ANSI escapes, and assert no mixed human text in JSON mode.

### Decision 4: Status JSON From Durable Local State

- **Decision:** Extend `ultraplan study <study> status --json` to emit a stable result payload with run counts, task state summaries, lock diagnostics, run metadata summaries, validation summaries, retry/fallback/cancellation metadata where present, artifact summaries, and usage/cost represented with explicit unknown/null/known flags rather than zero when unavailable. Status remains runtime-free.
- **Rationale:** Status is an operator inspection surface. It must summarize what exists in durable state and artifacts without initiating runtime health checks or process inspection. Unknown usage/cost must remain truthful for automation.
- **Study / Source Grounding:** `technical-handbook.md` observability, state/context, bounded inspection, and LLM evaluation/cost/safety pressures; `10-logging-observability`; `07-state-context`; `14-performance`; TRD canonical events/run metadata/run-state/status requirements; PRD status user scenario.
- **Trade-Offs Accepted:** Status may report stale durable state if a separate live process changed outside the state file; lock diagnostics and timestamps must make that visible rather than invoking runtime.
- **Technical Debt / Future Impact:** If future live status is required, it should be a separate live-inspection feature or explicitly scoped extension, not implicit runtime invocation in current status.
- **Alternatives Rejected:** Treating unknown usage/cost as `0` is rejected because it misleads automation. Calling agentwrap/OpenCode during status is rejected because requirements carry forward runtime-free status.
- **Contracts Satisfied:** Observability, Workflows, LLM Evaluation / Cost / Safety, Performance, Persistence And Migrations, Security, Testing, CLI Surface; requirements AC38-45, C68-75.
- **Evidence Required:** `internal/app/study_status_commands_test.go` parses status JSON and asserts counts, task states, lock diagnostics, run metadata summaries, validation summaries, retry/fallback metadata, unknown usage/cost behavior, redaction, and no runtime invocation.

### Decision 5: Health JSON Through Existing Runtime Health Boundary

- **Decision:** Extend `ultraplan health --json` in `internal/app/health_commands.go` to return a stable redacted result with health-check IDs, statuses, runtime/config/workspace health summaries, safe diagnostics, and guidance. Health may use the existing platform runtime/agentwrap health boundary, but platform runtime must not import or know study semantics.
- **Rationale:** Health is the appropriate surface for runtime/config/workspace readiness. It must be machine-readable and safe, while validation/status remain local inspection surfaces.
- **Study / Source Grounding:** `technical-handbook.md` configuration management, dependency injection, runtime boundary, and observability patterns; `04-configuration-management`; `03-dependency-injection`; `10-logging-observability`; TRD health, config, agentwrap mapping, and runtime boundary requirements.
- **Trade-Offs Accepted:** Health JSON may expose less low-level runtime detail than native tools because diagnostics must be redacted and bounded.
- **Technical Debt / Future Impact:** Additional runtime adapters should fit the same platform health summary shape without adding study imports or command-specific branching in platform runtime.
- **Alternatives Rejected:** Duplicating OpenCode process checks in `internal/app` is rejected because TRD requires agentwrap/opencode to own runtime process behavior. Importing study into platform runtime is rejected by architecture constraints.
- **Contracts Satisfied:** Configuration, Observability, Security, LLM Runtime, Testing, CLI Surface; requirements AC37, AC39-45, C66, C69, C73.
- **Evidence Required:** `internal/app/health_commands_test.go` covers `health --json` success/failure, stable check IDs/statuses, runtime-health diagnostics, schema fields, redaction, and fake runtime behavior.

### Decision 6: Redaction Before Text Or JSON Rendering

- **Decision:** Update `internal/platform/config/redaction.go` or adjacent existing redaction policy so validation, status, and health diagnostics are redacted before text or JSON rendering. New diagnostics must avoid API keys, secret config values, sensitive env values, prompts, embedded Markdown bodies, unsafe raw payload bytes, and full native stderr; safe excerpts must be truncated and redacted.
- **Rationale:** The sprint adds new automation surfaces that can expose nested diagnostics. Redaction must be enforced as an output precondition, not left to callers or tests after rendering.
- **Study / Source Grounding:** `technical-handbook.md` redaction boundary pattern; `13-security`; restic `internal/options/secret_string.go:15-20`, helm `pkg/registry/transport.go:37-41`, gh-cli `status.go:332-338`; TRD security and diagnostics requirements.
- **Trade-Offs Accepted:** Some raw context is intentionally lost. The result must compensate with check names, paths, expected conditions, safe observed summaries, and guidance.
- **Technical Debt / Future Impact:** Nested structure redaction can become inconsistent; tests must seed realistic secret-like values in every new surface.
- **Alternatives Rejected:** Redacting only text output is rejected because JSON is equally user-visible and often logged by automation. Persisting or exposing raw payload bytes/full stderr is rejected by sprint constraints and security requirements.
- **Contracts Satisfied:** Security, Configuration, Observability, Testing; requirements AC33, AC37-43, C69-70, C73-75.
- **Evidence Required:** Validation, status, and health tests seed API-key-like values, env/config secrets, runtime stderr snippets, prompts or Markdown-like bodies, and assert they do not appear in stdout/stderr JSON or text.

### Decision 7: Fake-First Schema And Behavior Tests

- **Decision:** Add deterministic offline tests for study validation, validate command JSON/text/exit codes, status JSON, and health JSON. JSON tests parse output and assert stable envelope/result field names for success and failure. Runtime-dependent paths use fakes or nil seams. Required verification commands are `go test ./...`, `go test -race ./...`, and `go build ./cmd/ultraplan` from `/home/antonioborgerees/coding/ultraplan-go`.
- **Rationale:** The sprint's value is stable automation behavior and safe diagnostics; behavior/schema tests are the clearest reviewable evidence. Default tests must not require OpenCode, provider credentials, network, source cloning, or long sleeps.
- **Study / Source Grounding:** `technical-handbook.md` testing strategy, IO abstraction, dependency injection, terminal UX, and performance evidence; `11-testing-strategy`; `06-io-abstraction`; `03-dependency-injection`; `09-terminal-ux`.
- **Trade-Offs Accepted:** Schema tests add maintenance cost and can constrain field renames; this is acceptable because JSON stability is a sprint goal.
- **Technical Debt / Future Impact:** Golden-style assertions should focus on public schema/behavior, not incidental private ordering beyond deterministic output requirements.
- **Alternatives Rejected:** Manual CLI smoke testing only is rejected because it cannot protect schema stability. Real OpenCode/provider tests in the default suite are rejected by requirements and would make verification non-deterministic.
- **Contracts Satisfied:** Testing, CLI Surface, Security, Documentation, Architecture; requirements AC22-25, AC36-49, C73-75.
- **Evidence Required:** New and updated tests in `internal/study/validation_command_test.go`, `internal/app/study_validate_commands_test.go`, `internal/app/study_status_commands_test.go`, and `internal/app/health_commands_test.go`; successful `go test ./...`, `go test -race ./...`, and `go build ./cmd/ultraplan` recorded in `review.md`.

## Expected Evidence

| Evidence Type | Required Evidence | Source / Command / Review Check |
| --- | --- | --- |
| Tests | Study validation covers valid study, missing artifacts, invalid reports, malformed frontmatter, inapplicable Markdown pairs, run-state diagnostics, malformed/unsupported state, and secret-free diagnostics. | `/home/antonioborgerees/coding/ultraplan-go/internal/study/validation_command_test.go` |
| Tests | Validate command covers text/JSON output, stable envelope/result fields, exit code `0` and `5`, checkable failures, redaction, and no runtime invocation. | `/home/antonioborgerees/coding/ultraplan-go/internal/app/study_validate_commands_test.go` |
| Tests | Status JSON covers counts, task states, locks, run metadata, validation summaries, retry/fallback/cancellation metadata, unknown usage/cost, redaction, and no runtime invocation. | `/home/antonioborgerees/coding/ultraplan-go/internal/app/study_status_commands_test.go` |
| Tests | Health JSON covers success/failure, health-check IDs/statuses, runtime/config/workspace summaries, schema fields, guidance, and redaction. | `/home/antonioborgerees/coding/ultraplan-go/internal/app/health_commands_test.go` |
| Verification | Unit and command tests pass offline. | `go test ./...` from `/home/antonioborgerees/coding/ultraplan-go` |
| Verification | Race verification passes. | `go test -race ./...` from `/home/antonioborgerees/coding/ultraplan-go` |
| Verification | CLI builds. | `go build ./cmd/ultraplan` from `/home/antonioborgerees/coding/ultraplan-go` |
| Review | Architecture boundaries remain intact and no forbidden global packages are added. | Architecture Review protocol; import/package review. |
| Review | JSON output contains no ANSI, progress, or mixed human text and text output remains concise/deterministic. | Sprint Review protocol; command output tests. |
| Review | Diagnostics are redacted across validation, status, and health. | Sprint Review protocol; redaction tests. |
| Documentation | Help text documents `study <study> validate`, `study <study> validate --json`, and `study <study> status --json`. | `/home/antonioborgerees/coding/ultraplan-go/internal/app/study_commands.go` help tests/review. |
| Review Artifact | Implementation evidence, deviations, verification status, and review conclusions are recorded. | `projects/ultraplan-go/sprints/14-validation-and-diagnostics/review.md` |

## Assumptions And Risks

| Item | Type | Impact | Mitigation / Follow-Up |
| --- | --- | --- | --- |
| Existing report validators can be reused or wrapped for command validation. | Assumption | If validators are too task-specific, implementation may duplicate rules. | Keep orchestration in `validation_command.go` and reuse existing helpers where possible; if gaps exist, add small study-owned helpers. |
| Existing status/run-state code exposes enough metadata for JSON summaries. | Assumption | If metadata is not accessible, status JSON may need additional study-domain summarizers. | Add focused study-owned summary functions without changing runtime or creating global status packages. |
| Existing redaction policy can be extended for new surfaces. | Assumption | If redaction is only string-based, nested result payloads may need explicit safe field construction. | Redact before render and test nested diagnostics with secret-like fixtures. |
| JSON schema vocabulary may be consumed by automation immediately. | Risk | Field renames after this sprint become breaking changes. | Version the envelope/result schema and assert field names in tests. |
| Validation/status could accidentally invoke runtime through shared factories. | Risk | Offline commands become slow, flaky, or unsafe. | Use nil/fake runtime seams in tests and review command construction for lazy dependencies. |
| Safe diagnostics may omit details users want for debugging. | Risk | Users may need extra manual inspection. | Include artifact paths, expected/observed safe summaries, guidance, and state paths; defer debug-retention to future explicit scope. |
| Unsupported run-state versions may frustrate users without migration. | Risk | Users get diagnostics but no automatic fix. | Requirements prohibit migrations; provide actionable guidance and follow-up migration scope later. |
| Area-specific architecture reasoning artifact is absent. | Risk | Plan may need architecture clarity. | This document records final architecture decisions directly so `plan.md` does not reopen them. |

## Implementation Constraints

- `internal/study` owns validation behavior, check semantics, study artifact/run-state inspection, applicability handling, and validation result types.
- `internal/app` owns CLI parsing, `--json` flag handling, help text, text/JSON rendering, and exit-code mapping only.
- `internal/platform/runtime` must not import `internal/study` or know studies, dimensions, sources, summaries, reports, validation-command semantics, or status payload semantics.
- Do not create global `internal/validation`, `internal/reports`, `internal/diagnostics`, `internal/status`, or workflow-engine packages.
- Validation and status must not invoke agentwrap, OpenCode, network calls, subprocesses, source cloning, synthesis, repair, retries, fallback execution, or prompt generation.
- Health JSON may use the existing platform runtime/agentwrap health boundary and fakeable seams, but must remain redacted and product-agnostic at the platform layer.
- JSON output must use a stable, versioned envelope with command, workspace, status, generated timestamp, and result payload; JSON mode must not emit ANSI, progress text, or human-only diagnostics.
- Unknown usage/cost must be represented as unknown/null/known flags and never coerced to numeric zero.
- All user-visible diagnostics must be redacted before text or JSON rendering.
- Paths in output should prefer workspace-relative paths; absolute paths are allowed only for necessary local diagnostics and must not expose secrets.
- Durable `run-state.json` must be read strictly; missing, malformed, or unsupported schema versions must produce clear diagnostics and no silent migration.
- Validation must be bounded to known workspace/study artifacts and must not recursively scan source repositories except through existing validators where already required.
- Tests must be deterministic, offline, fake-first, and must not require OpenCode, provider credentials, network access, source cloning, or long sleeps.
- Existing text behavior remains backward-compatible unless this sprint explicitly changes it and tests document the change.
- Implementation must not modify unrelated dirty files or revert user changes in the target repository.

## Plan Handoff

`plan.md` must execute these decisions. It must not invent architecture, scope, or decisions beyond this document.

The plan must carry forward:

- Final decisions for study-owned validation, runtime-free validation/status, app-owned JSON envelope, health JSON through runtime health boundary, redaction before rendering, and fake-first schema tests.
- Applicable contracts and requirement IDs from the Contracts Applied table.
- Expected evidence from the Expected Evidence table, including tests and verification commands.
- Risks and assumptions from this document.
- Required Architecture Review and Sprint Review protocols from `sprint-index.md`.

## Phase Exit Criteria

- [x] Selected context was read and used.
- [x] Area-specific reasoning documents were completed where applicable or explicitly marked not applicable.
- [x] Area-specific reasoning conclusions are reflected or explicitly overridden in final decisions.
- [x] Contracts are explicitly mapped to decisions and expected evidence.
- [x] Final decisions are clear enough for `plan.md` to execute.
- [x] Expected evidence is specific and reviewable.
