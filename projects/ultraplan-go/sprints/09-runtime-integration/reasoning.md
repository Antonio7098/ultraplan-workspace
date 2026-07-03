# Sprint Reasoning: 09 Runtime Integration

> Project: `ultraplan-go`
> Sprint: `09-runtime-integration`
> Output: `projects/ultraplan-go/sprints/09-runtime-integration/reasoning.md`
> **Inputs Used:** `projects/ultraplan-go/sprints/09-runtime-integration/requirements.md`, `projects/ultraplan-go/project-index.md`, `projects/ultraplan-go/docs/ARCHITECTURE.md`, `projects/ultraplan-go/docs/PRD.md`, `projects/ultraplan-go/docs/TRD.md`, `projects/ultraplan-go/sprints/09-runtime-integration/sprint-index.md`, `projects/ultraplan-go/sprints/09-runtime-integration/technical-handbook.md`, `projects/ultraplan-go/sprints/09-runtime-integration/reasoning/runtime-integration.md`, `templates/sprint-reasoning.md`, `/home/antonioborgerees/coding/agentwrap/docs/README.md`, `/home/antonioborgerees/coding/agentwrap/docs/architecture.md`, `/home/antonioborgerees/coding/agentwrap/docs/core-api.md`, `/home/antonioborgerees/coding/agentwrap/docs/events-and-metadata.md`, `/home/antonioborgerees/coding/agentwrap/docs/errors-config-health.md`, `/home/antonioborgerees/coding/agentwrap/docs/wrappers.md`, `/home/antonioborgerees/coding/agentwrap/docs/opencode-adapter.md`, `/home/antonioborgerees/coding/agentwrap/docs/integration-guide.md`, `/home/antonioborgerees/coding/agentwrap/runtime.go`, `/home/antonioborgerees/coding/agentwrap/events.go`, `/home/antonioborgerees/coding/agentwrap/metadata.go`, `/home/antonioborgerees/coding/agentwrap/config.go`, `/home/antonioborgerees/coding/agentwrap/health.go`, `/home/antonioborgerees/coding/agentwrap/policy.go`, `/home/antonioborgerees/coding/agentwrap/permissions.go`, `/home/antonioborgerees/coding/agentwrap/validation.go`, `/home/antonioborgerees/coding/agentwrap/observability.go`, `/home/antonioborgerees/coding/agentwrap/errors.go`, `/home/antonioborgerees/coding/agentwrap/opencode/options.go`

This document decides. It synthesizes selected context, handbook evidence, area-specific reasoning, and contracts into final sprint decisions.

It does not replace `sprint-index.md`, `technical-handbook.md`, or `reasoning/*.md`.

## Sprint Purpose

- **Goal:** Integrate UltraPlan's generic runtime boundary with `github.com/Antonio7098/agentwrap` and `agentwrap/opencode` so later study execution can use OpenCode through agentwrap wrappers, health checks, policies, validation, permissions, and events without UltraPlan reimplementing runtime supervision.
- **Non-Goals:** Public `study run`, `synthesize`, `run-all`, `run-loop`, batch scheduling, durable workflow orchestration, per-study locks, stale-running recovery, active task cancellation, summary generation, code-reference extraction, target workflows, sprint planning, sprint execution, stable new public runtime JSON schemas, and any UltraPlan-owned OpenCode subprocess supervisor, native event decoder, retry engine, permission translator, or validation wrapper that duplicates agentwrap behavior.
- **Depends On:** Prior sprints 1-8 for CLI/app wiring, workspace/config/logging/health foundations, study metadata context, workspace and prompt locations, Markdown source applicability, report validation hooks, prompt composition outputs, and existing run-state ID/task metadata fields. The sprint also depends on PRD/TRD/ARCHITECTURE runtime rules and selected evidence reports in `sprint-index.md`.

## Selected Context And Pre-Reasoning Artifacts

| Artifact | Path | How It Was Used |
| --- | --- | --- |
| Sprint requirements | `projects/ultraplan-go/sprints/09-runtime-integration/requirements.md` | Treated as the authoritative sprint contract for outputs, acceptance criteria, non-goals, constraints, dependencies, and review expectations. Requirement IDs in this document refer to the ordered acceptance criteria as `AC-01` through `AC-30`. |
| Project index | `projects/ultraplan-go/project-index.md` | Confirmed project scope, target implementation directory, active contract pool, selected source docs, available evidence reports, and the explicit instruction not to invent a competing runtime contract. |
| Product requirements | `projects/ultraplan-go/docs/PRD.md` | Grounded product principles: runtime success is not product success, product workflows stay product-owned, runtime adapters execute agent work, health checks are required, canonical events must tolerate unknown future event types, and tests use fake runtimes by default. |
| Technical requirements | `projects/ultraplan-go/docs/TRD.md` | Grounded concrete agentwrap dependency, OpenCode adapter use, wrapper composition, request/config mapping, health, canonical events, metadata, retry/fallback, validation, permissions, errors, cancellation, and tests. |
| Architecture | `projects/ultraplan-go/docs/ARCHITECTURE.md` | Applied the dependency rule `study -> platform/runtime` and `platform/runtime -> no product modules`; confirmed runtime is generic platform execution infrastructure. |
| Sprint index | `projects/ultraplan-go/sprints/09-runtime-integration/sprint-index.md` | Used selected contracts, selected evidence reports, selected area reasoning template, excluded context, and required review protocols to limit scope. |
| Technical handbook | `projects/ultraplan-go/sprints/09-runtime-integration/technical-handbook.md` | Used distilled source evidence for thin commands, manual composition, decorator wrappers, merged config validation, typed errors, IO/fake seams, context propagation, structured events, permission boundaries, and behavior-focused tests. |
| Agentwrap docs and source | `/home/antonioborgerees/coding/agentwrap/docs/*.md`; `/home/antonioborgerees/coding/agentwrap/*.go`; `/home/antonioborgerees/coding/agentwrap/opencode/options.go` | Confirmed the current SDK surface: `Runtime.StartRun`, `Run.Events`, `Run.Wait`, `Run.Cancel`, `Capabilities`, `HealthChecker.CheckHealth`, `EffectiveConfig`, `PolicyRunner`, `BasicPolicy`, `ValidatingRuntime`, `ValidationSpec`, `ObservingRuntime`, `MemoryRunStore`, `PermissionPolicy`, `PermissionUnsupportedBestEffort`, `SDKError`, canonical event kinds, nil usage semantics, raw payload persistence policy, and OpenCode construction options. |
| Sprint reasoning template | `templates/sprint-reasoning.md` | Used as the structure for this consolidated decision artifact. |

## Area-Specific Reasoning Inputs

| Area | Reasoning Document | Key Conclusion | Evidence Basis | Impact On Final Decision |
| --- | --- | --- | --- | --- |
| Runtime integration architecture | `projects/ultraplan-go/sprints/09-runtime-integration/reasoning/runtime-integration.md` | Proceed with a thin `internal/platform/runtime` boundary over agentwrap/opencode, keep study semantics out, compose `ObservingRuntime -> ValidatingRuntime -> PolicyRunner -> opencode.Runtime`, and reject direct OpenCode supervision, direct product-module agentwrap usage, broad registries, unsafe raw payload persistence, and validation-driven fallback by default. | Requirements, PRD/TRD/ARCHITECTURE, technical handbook patterns, and selected evidence reports on thin boundaries, composition roots, decorators, typed errors, config validation, fake tests, structured events, and permission trust boundaries. | This conclusion is adopted as the main sprint direction and expanded into final decisions for boundary ownership, composition, config/health, policy/permissions, events/errors, and verification. |

## Sprint Technical Handbook Summary

- **Relevant Patterns:** Thin runtime boundary behind thin commands; manual composition root with explicit wrapper stack; decorator/wrapper composition for runtime concerns; merged config then validate once; typed error classification over string matching; injectable IO and fake runtime seams; context as lifecycle; structured events/logs; explicit permission trust boundaries; behavior tests over private implementation tests.
- **Important Trade-Offs:** Thin UltraPlan runtime model versus exposing raw agentwrap everywhere; fixed wrapper composition versus task-specific wrapper order; centralized runtime config validation versus distributed adapter checks; agentwrap-owned retry/fallback versus UltraPlan-owned scheduling policy; safe minimal diagnostics versus rich native event persistence; config-driven permission posture versus hardcoded defaults; fake-only default tests versus gated real OpenCode smoke tests.
- **Warnings / Anti-Patterns:** Do not reimplement the runtime supervisor; do not string-match runtime errors; do not use global runtime config or service singletons; do not let CLI health invoke study workflows; do not drop unknown event kinds or unsafe payload facts; do not convert unknown usage/cost to zero; do not use `context.Background()` in work paths; do not make unsupported required permission features best-effort silently.
- **Evidence Confidence:** High overall. Most handbook findings were high confidence and came from multiple mature Go CLIs plus the project PRD/TRD/ARCHITECTURE. Direct Agentwrap docs/source inspection confirms the concrete SDK names and semantics needed for implementation, so remaining uncertainty is about fitting those APIs into the existing UltraPlan app/config code rather than about the runtime architecture.

## Contracts Applied

| Contract / Requirement ID | Constraint | Decision Impact | Expected Evidence |
| --- | --- | --- | --- |
| Architecture, AC-02, AC-03 | Runtime behavior must live in `internal/platform/runtime`, stay generic, and avoid product imports. | Decide a thin platform runtime boundary with no `internal/study`, `internal/codeextract`, or product-module dependency. | Import review; runtime package tests; Architecture Review Protocol confirmation. |
| LLM Runtime, AC-01, AC-04, AC-05, AC-11, AC-12, AC-13 | Use agentwrap root APIs and `agentwrap/opencode`; compose wrappers in the required order; use agentwrap validation and policy. | Decide agentwrap is the external runtime SDK and OpenCode adapter owner; UltraPlan maps requests/results only. | `go.mod` dependency; wrapper composition tests; validation failure tests; retry/fallback tests; no direct OpenCode process launch. |
| Configuration, AC-06, AC-07, AC-08 | Runtime config must map known fields, fail fast on invalid/unsupported fields, and redact sensitive values. | Decide minimal runtime config extensions with merged validation before runtime launch or health checks. | Config unit tests for defaults, invalid values, required health names, precedence where present, and redaction. |
| Observability, AC-10, AC-17, AC-18, AC-19, AC-20, AC-21, AC-23 | Preserve correlation metadata, canonical events, safe diagnostics, raw payload omission facts, result metadata, and unknown usage/cost knownness. | Decide event/result mapping must preserve agentwrap event kinds and safe metadata without persisting unsafe native bytes by default. | Event mapping tests for all canonical kinds, native extension, raw omission facts, usage unknownness, warnings, and metadata. |
| Errors, AC-22, AC-23 | Preserve agentwrap `SDKError` classifications and cause chains without string matching. | Decide runtime error mapping uses `errors.As` or agentwrap helpers and surfaces safe classified diagnostics. | Error mapping tests for representative SDK categories; review check for no string classification. |
| Security, AC-08, AC-14, AC-19 | Redact secrets; express sandbox and permissions through agentwrap; fail before launch when required policy features cannot be represented. | Decide permission posture is explicit and representability is validated before OpenCode launch. | Permission mapping tests; unsupported required feature tests; redaction tests; diagnostics review. |
| CLI Surface, AC-15, AC-16, AC-25 | `ultraplan health` may include runtime health but must not run study workflows or prompt execution. | Decide health command delegates only to platform runtime health/capability checks and emits actionable diagnostics. | Health command tests for success/failure/redaction/exit status and no study workflow invocation. |
| Workflows, AC-09, AC-13, AC-24, AC-26 | Context and cancellation must propagate; retry/fallback belongs to agentwrap; default tests must use fakes and avoid real subprocesses. | Decide all runtime run/health paths accept caller context, drain events, call `Wait`, and use agentwrap-compatible fakes in normal tests. | Fake runtime tests for cancellation, timeout, retry, fallback, event drain, and no OpenCode/credentials/network requirement. |
| Testing, AC-24, AC-25, AC-26, AC-27, AC-28, AC-29, AC-30 | Required fake/config/health coverage, optional gated smoke tests, existing tests passing, `go test ./...`, and `go build ./cmd/ultraplan`. | Decide default verification is deterministic and runtime-free; optional real OpenCode smoke path is environment-gated and skipped by default. | `go test ./...`; `go build ./cmd/ultraplan`; smoke test skip evidence if added. |

## Repos Studied / Source Evidence Used

| Source / Repo / Report | Concrete Reference | Relevant Finding | Why It Matters For This Sprint | Used In Decision(s) |
| --- | --- | --- | --- | --- |
| `01-project-structure` report | `studies/go-cli-study/reports/final/01-project-structure.md`; `cmd/root.go:255-259`, `cmd/restic/main.go:37-114`, `internal/restic/repository.go:18` | Mature CLIs keep entrypoints thin and dependencies one-way into internal logic. | Supports keeping CLI health thin and runtime generic in `internal/platform/runtime`. | Decisions 1, 4 |
| `02-command-architecture` report | `studies/go-cli-study/reports/final/02-command-architecture.md`; `pkg/cmdutil/factory.go:16-43`, `pkg/cmd/install.go:132-145` | Factory-created commands and thin `RunE` functions delegate to services. | Supports health command delegating to platform runtime instead of owning runtime behavior. | Decision 4 |
| `03-dependency-injection` report | `studies/go-cli-study/reports/final/03-dependency-injection.md`; `internal/app/app.go:42-81`, `executor.go:22-24`, `ui_concurrent.go:9-12` | Manual composition roots, constructor injection, functional options, and decorators are preferred; globals are risky. | Supports explicit agentwrap stack construction and fake injection without DI framework or global runtime state. | Decisions 1, 2, 7 |
| `04-configuration-management` report | `studies/go-cli-study/reports/final/04-configuration-management.md`; `internal/config/config.go:609-641`, `internal/flags/flags.go:314-327` | Merge config sources, then validate once with explicit field errors. | Supports runtime config validation before launch and redacted config diagnostics. | Decision 3 |
| `05-error-handling` report | `studies/go-cli-study/reports/final/05-error-handling.md`; `fs/fserrors/error.go:22-29`, `internal/errors/fatal.go:10`, `internal/ghcmd/cmd.go:44-49` | Runtime-facing errors should preserve typed/sentinel classification through wrapping and `errors.Is`/`errors.As`. | Supports preserving `agentwrap.SDKError` classifications and avoiding string matching. | Decision 6 |
| `06-io-abstraction` report | `studies/go-cli-study/reports/final/06-io-abstraction.md`; `pkg/iostreams/iostreams.go:551-568`, `internal/ui/terminal.go:10-36`, `executor.go:553-564` | IO and external-boundary seams make commands and tests deterministic. | Supports buffer-backed health tests and fake runtime seams. | Decisions 4, 7 |
| `07-state-context` report | `studies/go-cli-study/reports/final/07-state-context.md`; `internal/ghcmd/cmd.go:142`, `pkg/cmd/install.go:333-347`, `internal/session/session.go:12-23` | Long-running external operations need context propagation and explicit cancellation. | Supports passing caller context through health and runtime requests and avoiding `context.Background()` in work paths. | Decisions 1, 7 |
| `08-concurrency` report | `studies/go-cli-study/reports/final/08-concurrency.md`; `task.go:87`, `internal/repository/repository.go:567`, `cmd/root.go:261-279` | Structured goroutine lifecycle, cancellation, and channel draining avoid leaks. | Supports event draining and `Run.Wait` expectations in runtime tests. | Decisions 6, 7 |
| `10-logging-observability` report | `studies/go-cli-study/reports/final/10-logging-observability.md`; `internal/logging/logger.go:25-62`, `internal/pubsub/broker.go:10-19`, `pkg/iostreams/iostreams.go:52-54` | Structured logs and event/pubsub models support diagnostics without mixing user output and logs. | Supports canonical event mapping, safe diagnostics, and separating health output from logs. | Decisions 4, 6 |
| `11-testing-strategy` report | `studies/go-cli-study/reports/final/11-testing-strategy.md`; `internal/chezmoitest/chezmoitest.go:86-92`, `pkg/httpmock/stub.go:35-199`, `internal/backend/mock/backend.go:14-26` | Reliable CLIs use table-driven tests, fixtures, mocks/fakes, and gated integration tests. | Supports fake-runtime default coverage and optional environment-gated OpenCode smoke tests. | Decision 7 |
| `12-extensibility` report | `studies/go-cli-study/reports/final/12-extensibility.md`; `internal/llm/tools/tools.go:69-72`, `internal/llm/agent/mcp-tools.go:106-129` | Extension systems use explicit protocol boundaries and registries only when justified. | Supports rejecting a broad runtime plugin registry while relying on agentwrap's boundary. | Decisions 1, 2 |
| `13-security` report | `studies/go-cli-study/reports/final/13-security.md`; `internal/permission/permission.go:44-108`, `internal/llm/tools/bash.go:41-55`, `internal/options/secret_string.go:15-20` | Agent execution needs explicit permissions, command filtering, secret redaction, and fail-fast trust boundaries. | Supports permission policy mapping, sandbox config, redaction, and unsupported-feature rejection. | Decisions 3, 5, 6 |
| `15-philosophy` report | `studies/go-cli-study/reports/final/15-philosophy.md`; `pkg/cmd/factory/default.go:26-46`, `internal/pubsub/broker.go:10-19`, `internal/permission/permission.go:74-108` | Strong CLIs accept complexity deliberately and reject scope that does not serve the current tool. | Supports minimal runtime config and rejecting plugin systems, workflow engines, and direct OpenCode supervision. | Decisions 1, 2, 7 |

## Trade-Off And Debt Analysis

### Accepted Trade-Offs

| Trade-Off | Benefit | Cost / Constraint Accepted | Why Acceptable Now | Revisit Trigger |
| --- | --- | --- | --- | --- |
| Thin UltraPlan internal runtime model over raw agentwrap types | Keeps platform/runtime generic, prevents study leakage, and gives app/study callers a stable internal shape. | Requires careful field mapping and risks hiding useful agentwrap detail. | The mapping is bounded by sprint requirements and must remain field-oriented, not a competing SDK. | Product callers repeatedly need unmapped agentwrap fields or mapping code starts duplicating agentwrap behavior. |
| Required fixed wrapper order | Provides predictable observation, validation, policy, and OpenCode delegation and satisfies sprint AC-04. | Validation failures do not participate in policy retry/fallback by default because `ValidatingRuntime` wraps `PolicyRunner`. | This is the explicitly required production composition for this sprint. | A later reviewed task type needs validation-driven fallback and records a wrapper-order exception. |
| Central config validation before launch | Produces actionable preflight errors and prevents unsupported runtime/health/permission launches. | Config validation must know enough adapter capabilities to avoid accepting unrepresentable values. | Fail-fast behavior is required and safer than launching with weakened policy. | Agentwrap exposes dynamic capability negotiation that should replace static checks. |
| Agentwrap-owned retry/fallback | Avoids duplicate retry engines and keeps policy close to SDK classifications/events. | UltraPlan must translate policy metadata and cannot customize every retry detail locally. | Sprint goal explicitly delegates runtime supervision to agentwrap. | Future workflow-level budgets require coordination above agentwrap while preserving SDK-owned per-runtime attempts. |
| Safe diagnostic summaries instead of raw native payload persistence | Reduces secret leakage and artifact bloat while preserving audit facts. | Some deep debugging requires opt-in raw payload retention later. | Security and observability requirements prefer safe omission facts by default. | A debug mode is explicitly designed with retention, redaction, and user consent. |
| Fake-only default verification | Keeps `go test ./...` deterministic without OpenCode, credentials, network, or subprocesses. | Fakes can diverge from real adapter behavior. | Optional gated smoke tests can cover real adapter paths without destabilizing default tests. | Agentwrap/opencode API changes or integration defects appear that fakes cannot detect. |

### Potential Technical Debt

| Debt / Shortcut | Why It Might Accrue | Current Mitigation | Owner / Follow-Up |
| --- | --- | --- | --- |
| Runtime mapping layer growth | Future study execution may add metadata, validation, repair, and session fields until the platform facade resembles a full SDK. | Keep fields generic, map to agentwrap root types, and reject study-specific types in `internal/platform/runtime`. | Future execution sprint reviews runtime boundary before adding workflow-specific fields. |
| Static health/capability names in config validation | Agentwrap/opencode supported health/capability IDs may change. | Inspect real agentwrap APIs during implementation and centralize accepted-name mapping in runtime/config tests. | Runtime integration implementation and future agentwrap upgrade review. |
| Memory/fake run store for this sprint | It does not provide durable run inspection across process exits. | Sprint requirements explicitly exclude durable agentwrap run records; result metadata is shaped for later persistence. | Batch/run-loop or durable runtime-record sprint. |
| Health command output shape | Extending health output without a stable new JSON schema can create later compatibility work. | Reuse existing output conventions and avoid promising a new stable runtime JSON schema. | Public JSON schema or status/health API sprint. |
| Permission mapping subset | OpenCode adapter may not represent all desired path/tool/shell restrictions. | Unsupported required features fail before launch unless explicitly configured best-effort. | Security/runtime policy follow-up after real adapter capability inspection. |

### Future Considerations

| Consideration | Deferred Until | Reason Deferred | What Should Be Preserved Now |
| --- | --- | --- | --- |
| Public study execution wiring | Later study run/synthesize/run-all/run-loop sprint | This sprint builds runtime infrastructure only. | Generic request metadata fields for run ID, task ID, task kind, study, dimension, source, source kind, output path, runtime, provider, and model. |
| Durable agentwrap run records | Batch/run-loop persistence sprint | Requirements allow in-memory or fake stores for this sprint. | Safe result metadata, run IDs, session IDs, events, attempts, validation, policy, permissions, artifacts, usage, and warnings. |
| Validation-driven retry/fallback | Later reviewed exception | Required wrapper order does not retry validation failures by default. | Explicit wrapper constructor tests and documented default behavior. |
| Repair attempts | Later execution or repair sprint | Repair policy can pull in study report semantics and prompt construction. | Generic validation hooks and safe repair metadata passthrough if agentwrap returns it. |
| Additional runtimes | Future adapter sprint | Only OpenCode through agentwrap is selected now. | Runtime name validation and narrow constructor seams without a broad plugin registry. |
| Stable runtime health JSON | Public output compatibility sprint | Sprint explicitly excludes new stable runtime schemas unless already present. | Reuse existing output mode and avoid undocumented schema commitments. |

## Final Decisions

### Decision 1: Generic Platform Runtime Boundary

- **Decision:** Implement runtime integration in `internal/platform/runtime` as a thin UltraPlan-internal boundary with generic request, result, health, validation, permission, event, metadata, and diagnostic-facing types. The package may depend on agentwrap root APIs and `agentwrap/opencode`, but it must not import `internal/study`, `internal/codeextract`, or other product modules.
- **Rationale:** Runtime execution is generic platform infrastructure. Product modules own study semantics, prompt/report rules, scheduling, synthesis gating, and state machines; platform runtime owns only execution-facing mapping and safe diagnostics.
- **Study / Source Grounding:** `technical-handbook.md` patterns "Thin Runtime Boundary Behind Thin Commands" and "Manual Composition Root With Explicit Wrapper Stack"; `01-project-structure` evidence on one-way dependencies; `03-dependency-injection` evidence on explicit composition; `ARCHITECTURE.md` lines 159-197 define runtime as generic and `study -> platform/runtime` with no product imports.
- **Trade-Offs Accepted:** Accept a narrow mapping layer rather than exposing raw agentwrap types everywhere. This adds translation work but prevents product/runtime coupling.
- **Technical Debt / Future Impact:** The mapping layer must stay thin. Future study execution may add metadata fields only as generic strings/values, not study structs.
- **Alternatives Rejected:** Product modules calling agentwrap directly is rejected because it spreads runtime mechanics into study code. A public UltraPlan runtime SDK is rejected because it would compete with agentwrap. A broad runtime plugin registry is rejected because only OpenCode via agentwrap is in scope.
- **Contracts Satisfied:** Architecture; LLM Runtime; Testing; AC-01, AC-02, AC-03, AC-09, AC-10, AC-28.
- **Evidence Required:** Runtime package import review; tests proving request metadata mapping; no product imports in `internal/platform/runtime`; Architecture Review Protocol confirmation.

### Decision 2: Agentwrap/OpenCode Composition And Ownership

- **Decision:** Add `github.com/Antonio7098/agentwrap` as the runtime SDK dependency and construct the default production stack as `agentwrap.ObservingRuntime -> agentwrap.ValidatingRuntime -> agentwrap.PolicyRunner -> opencode.Runtime`. OpenCode must be constructed only through `agentwrap/opencode` options such as executable, extra args, environment, and stderr limits. UltraPlan must not launch `opencode` directly or parse OpenCode stdout/stderr directly.
- **Rationale:** Agentwrap owns volatile runtime supervision: OpenCode process execution, native event projection, health enforcement, validation wrapper behavior, retry/fallback policy, permission translation, diagnostics, and SDK error classification.
- **Study / Source Grounding:** `TRD.md` sections 11.1-11.4 require agentwrap root API usage, the wrapper composition, and opencode options. Agentwrap docs/source confirm `agentwrap.Runtime`, `agentwrap.RunRequest`, `agentwrap.PolicyRunner`, `agentwrap.ValidatingRuntime`, `agentwrap.ObservingRuntime`, `agentwrap.NewMemoryRunStore`, and `opencode.NewRuntime` with `WithExecutable`, `WithExtraArgs`, `WithEnv`, and `WithStderrLimit`. `technical-handbook.md` patterns "Decorator/Wrapper For Cross-Cutting Runtime Concerns" and anti-pattern "Do Not Reimplement The Runtime Supervisor" support composing wrappers rather than duplicating them. `03-dependency-injection` and `15-philosophy` evidence supports visible composition and deliberate complexity.
- **Trade-Offs Accepted:** Validation failures are outside policy retry/fallback by default because of the required wrapper order. This is acceptable for this sprint and any exception must be reviewed later.
- **Technical Debt / Future Impact:** If future tasks require validation-triggered fallback, the architecture must record an explicit wrapper-order exception rather than silently changing the default.
- **Alternatives Rejected:** UltraPlan-owned OpenCode subprocess supervision is rejected because it duplicates agentwrap/opencode. Validation inside policy by default is rejected because it violates the selected production order. Adapter-specific direct stdout/stderr parsing is rejected because agentwrap diagnostics are the boundary.
- **Contracts Satisfied:** LLM Runtime; Errors; Observability; Security; AC-01, AC-04, AC-05, AC-11, AC-12, AC-13, AC-22.
- **Evidence Required:** `go.mod` dependency; stack composition tests; adapter construction tests; validation failure test where underlying runtime succeeds but logical result fails; retry/fallback tests proving policy delegation; review check for no direct OpenCode process handling.

### Decision 3: Minimal Runtime Config Mapping, Validation, And Redaction

- **Decision:** Extend config only as needed for runtime name, executable, extra args, safe environment additions, provider/model, timeout, retry count, fallback provider/model, required health checks, required capabilities, sandbox, permission posture, stderr/diagnostic limits, and redacted diagnostics. Resolve config through existing precedence, then validate before runtime construction, launch, or health checks.
- **Rationale:** Runtime config must be expressive enough for agentwrap/opencode while staying minimal and safe. Invalid runtime names, unsupported required health/capability names, empty executable values, invalid timeouts/retries, and unsupported required permission policy features must fail before launch.
- **Study / Source Grounding:** `technical-handbook.md` pattern "Merged Config Then Validate Once"; `04-configuration-management` evidence around config precedence and validation; `13-security` evidence around secret redaction and trust boundaries; `TRD.md` sections 6.1-6.5 map runtime config to agentwrap/opencode options. Agentwrap source confirms `EffectiveConfig`, `ConfigLayer`, `CallerConfigLayer`, `ValidateEffectiveConfig`, current health IDs (`runtime_available`, `structured_output`, `workdir`, `config`, `provider`, `model`, `authentication`, `runtime_paths`), and capability IDs (`structured_events`, `raw_payloads`, `cancellation`, `artifacts`, `permissions`, `usage`, `validation_events`, and session capability variants).
- **Trade-Offs Accepted:** Central validation must know a small amount about adapter capabilities. This is safer than distributed checks that fail after launch.
- **Technical Debt / Future Impact:** Accepted health/capability names and permission mappings may need updates as agentwrap evolves. Keep mapping centralized and covered by tests.
- **Alternatives Rejected:** Hardcoding runtime defaults without config is rejected because requirements include provider/model/policy/health mapping. Accepting arbitrary health or permission strings best-effort is rejected because required unsupported features must fail fast. Storing secrets in config output or logs is rejected.
- **Contracts Satisfied:** Configuration; Security; LLM Runtime; AC-06, AC-07, AC-08, AC-14, AC-15.
- **Evidence Required:** Config tests for defaults, validation failures, health/capability name mapping, empty executable, timeout/retry bounds, fallback mapping, permission unsupported-feature rejection, environment redaction, and config output redaction.

### Decision 4: Runtime Health Through Platform Boundary Only

- **Decision:** Extend `ultraplan health` so runtime checks call the platform runtime health/capability boundary when requested or configured. Health must report actionable diagnostics for missing executable, missing required capability, unsupported health check, authentication/provider/model failures, runtime unavailable states, and redacted config/environment facts. It must not invoke study workflows, prompt execution, batch scheduling, report validation, or OpenCode directly.
- **Rationale:** Health is a preflight/query path, not a hidden execution path. Keeping it narrow preserves CLI clarity and prevents runtime integration from pulling public study execution into scope.
- **Study / Source Grounding:** `technical-handbook.md` patterns "Thin Runtime Boundary Behind Thin Commands" and "Injectable IO And Fake Runtime Seams"; `02-command-architecture` evidence on thin command delegation; `06-io-abstraction` evidence for buffer-backed command tests; `PRD.md` runtime health checks; `TRD.md` runtime health and CLI requirements.
- **Trade-Offs Accepted:** Health output is extended without defining a new stable public runtime JSON schema. It must reuse existing output conventions unless implementation explicitly documents and tests a schema.
- **Technical Debt / Future Impact:** Public JSON compatibility may need a later schema decision. The current sprint should avoid commitments beyond existing health output behavior.
- **Alternatives Rejected:** Running a prompt as a health check is rejected because it is execution, not health. Calling `opencode` directly from the CLI is rejected because `agentwrap/opencode` owns process behavior. Invoking study validation or prompt composition from health is rejected as out of scope.
- **Contracts Satisfied:** CLI Surface; Architecture; Configuration; Security; AC-15, AC-16, AC-25, AC-26, AC-28.
- **Evidence Required:** Health command tests for success, missing executable/unavailable runtime, unsupported health requirement, redacted diagnostics, stable runtime failure exit status, and no study workflow/prompt/report invocation.

### Decision 5: Policy, Fallback, Sandbox, And Permission Mapping

- **Decision:** Map retry counts, rate-limit behavior, fallback provider/model, sandbox mode, permission mode, and structured permission posture into agentwrap policy and permission structures. Retry, rate-limit, fallback, exhaustion, and permission metadata/events must come from agentwrap. Unsupported required permission features must fail before launch unless explicitly configured as best-effort.
- **Rationale:** Policy and permissions are runtime trust boundaries. UltraPlan should configure them and surface decisions, not implement a separate policy engine or silently weaken required posture.
- **Study / Source Grounding:** `technical-handbook.md` trade-off "Agentwrap-owned retry/fallback vs UltraPlan-owned scheduling policy" and pattern "Permission And Trust Boundaries As First-Class Runtime Inputs"; `13-security` evidence from opencode permission flow and bash allow/deny policy; `05-error-handling` evidence for classified retry/fatal semantics; `TRD.md` sections 14 and 22.1. Agentwrap source confirms `BasicPolicy`, `FallbackAlternative`, `PolicyDecision`, `RunMetadata.Policy.Decisions`, `RateLimitInfo`, `PermissionPolicy`, `PermissionTool`, `PermissionAction`, `PermissionUnsupportedFail`, and `PermissionUnsupportedBestEffort`; OpenCode docs confirm path-level rules are unsupported by the current subprocess adapter unless best-effort is explicitly selected.
- **Trade-Offs Accepted:** UltraPlan gives up local control over retry/fallback internals in exchange for a single SDK-owned mechanism and consistent metadata.
- **Technical Debt / Future Impact:** Future workflow-level retry budgets may coordinate above agentwrap, but per-runtime retry/fallback remains SDK-owned. Permission coverage may need follow-up after inspecting exact OpenCode adapter capabilities.
- **Alternatives Rejected:** Implementing retry/fallback in the OpenCode adapter is rejected. Silently downgrading required permission features is rejected. Hardcoded permissive runtime behavior is rejected because study-side source repositories are untrusted inputs.
- **Contracts Satisfied:** LLM Runtime; Workflows; Security; Observability; Errors; AC-06, AC-13, AC-14, AC-20, AC-23, AC-24.
- **Evidence Required:** Policy mapping tests for retry, rate limit, fallback, exhaustion, and metadata; permission mapping tests for supported and unsupported required features; event tests for `rate_limit`, `retry`, `fallback`, and `permission`; diagnostics review for redaction.

### Decision 6: Event, Metadata, Diagnostics, Usage, And Error Mapping

- **Decision:** Consume agentwrap canonical events and map them into UltraPlan runtime event/result/log fields without dropping event kind. Preserve safe metadata for attempts, policy decisions, validation, repair, sessions, permissions, cleanup, artifacts, usage, estimated cost, warnings, and classified errors. Unknown native events projected as `native_extension` must not fail a run by themselves. Unsafe raw native payload bytes must not be persisted or logged by default; record presence, source, encoding, safety, and omission reason where available. Unknown token, usage, or cost values remain unknown, not zero. Runtime errors must preserve `agentwrap.SDKError` classifications through `errors.As` or agentwrap helpers.
- **Rationale:** Later run inspection needs useful runtime facts, but security requires redaction and safe payload handling. Typed error and event metadata are more reliable than parsing human-readable runtime text.
- **Study / Source Grounding:** `technical-handbook.md` patterns "Typed Error Classification Over String Matching" and "Structured Events And Logs For Operator Diagnostics"; `05-error-handling` evidence on wrapping and typed classifications; `08-concurrency` evidence on event lifecycle and channel safety; `10-logging-observability` evidence on structured logs/pubsub; `TRD.md` sections 12, 14.2-14.3, and 19. Agentwrap source confirms `Event.Kind()`, `RawPayload.Safe`, `RunEventRecord` raw omission fields, `PersistencePolicy.PersistUnsafeRawPayloads`, pointer-based `Usage` token fields, `RunMetadata` summaries, and `SDKError` categories/fields with `agentwrap.ErrorAs`.
- **Trade-Offs Accepted:** Raw debugging power is reduced by default to avoid leaking secrets and unsafe native payloads. Event/result mapping has to preserve knownness and classifications carefully.
- **Technical Debt / Future Impact:** A future debug-retention mode may be added, but it must explicitly handle retention policy, consent, redaction, and artifact safety.
- **Alternatives Rejected:** Persisting raw native payload bytes by default is rejected. Dropping unknown event kinds is rejected. Converting missing usage/cost to zero is rejected. Classifying runtime errors by string matching is rejected.
- **Contracts Satisfied:** Observability; Errors; Security; LLM Evaluation / Cost / Safety; AC-17, AC-18, AC-19, AC-20, AC-21, AC-22, AC-23, AC-24.
- **Evidence Required:** Event mapping tests for lifecycle, session, message, progress, tool, artifact, permission, blocking, usage, warning, fatal_error, rate_limit, validation, retry, fallback, final_result, and native_extension; malformed event tests; raw payload omission tests; usage knownness tests; `SDKError` mapping tests; event drain tests.

### Decision 7: Deterministic Fake-First Verification

- **Decision:** Runtime tests must use agentwrap-compatible fakes by default and must not require OpenCode, provider credentials, network access, or real subprocess execution. Required coverage includes successful run, successful runtime exit with validation failure, runtime failure, timeout, cancellation, malformed event, rate limit, validation failure, retry, fallback, permission failure, observability metadata, config validation, health output, and redaction. Any real OpenCode smoke test must be explicitly environment-gated, skipped by default, and go through `agentwrap/opencode`.
- **Rationale:** This sprint adds infrastructure that must be reviewable and deterministic. Real runtime smoke tests are valuable but cannot be required for normal `go test ./...`.
- **Study / Source Grounding:** `technical-handbook.md` pattern "Behavioral Tests Over Implementation Detail Tests"; `11-testing-strategy` evidence on fakes/mocks/gated integration; `06-io-abstraction` evidence on buffer-backed IO; `TRD.md` testing sections 23.1-23.4; `PRD.md` tests and fixtures requirements.
- **Trade-Offs Accepted:** Fakes may not catch every real OpenCode behavior. Gated smoke coverage is optional and must not destabilize default verification.
- **Technical Debt / Future Impact:** Fake/runtime divergence should be revisited when agentwrap/opencode changes or when public study execution starts using the runtime heavily.
- **Alternatives Rejected:** Default tests that require OpenCode or credentials are rejected. Tests that assert private wrapper internals instead of observable behavior are rejected. Ungated real-runtime integration tests are rejected.
- **Contracts Satisfied:** Testing; Workflows; CLI Surface; AC-24, AC-25, AC-26, AC-27, AC-28, AC-29, AC-30.
- **Evidence Required:** `internal/platform/runtime/runtime_test.go`, `internal/platform/config/config_test.go`, and `internal/app/health_commands_test.go`; `go test ./...`; `go build ./cmd/ultraplan`; smoke skip evidence if any smoke test is added.

## Expected Evidence

| Evidence Type | Required Evidence | Source / Command / Review Check |
| --- | --- | --- |
| Dependency | Agentwrap dependency declared and runtime integration imports agentwrap root APIs plus `agentwrap/opencode`, not a competing public SDK. | `/home/antonioborgerees/coding/ultraplan-go/go.mod`; code review. |
| Architecture | Runtime behavior is in `internal/platform/runtime`; app health remains in `internal/app`; platform runtime has no product-module imports; no direct OpenCode process launch exists in UltraPlan. | Import review; Architecture Review Protocol. |
| Tests | Request mapping, wrapper composition, health mapping, policy/permission mapping, event handling, validation failure, retry/fallback, cancellation, safe diagnostics, redaction, and error classification covered with fakes. | `/home/antonioborgerees/coding/ultraplan-go/internal/platform/runtime/runtime_test.go`. |
| Tests | Runtime config defaults, validation, required health/capability names, timeout/retry/fallback validation, unsupported permission feature rejection, precedence where applicable, and redaction covered. | `/home/antonioborgerees/coding/ultraplan-go/internal/platform/config/config_test.go`. |
| Tests | Runtime health output, missing executable/unavailable runtime, unsupported health requirement, redacted diagnostics, stable exit status, and no study workflow invocation covered. | `/home/antonioborgerees/coding/ultraplan-go/internal/app/health_commands_test.go`. |
| Runtime | Canonical event kinds preserved, native extensions non-fatal, raw payload bytes omitted by default with omission facts, unknown usage/cost remains unknown, and safe metadata is exposed. | Runtime event/result tests and review of `events.go`. |
| Runtime | Runtime-facing errors preserve `agentwrap.SDKError` classifications through `errors.As` or agentwrap helpers and include safe diagnostics. | Error mapping tests and code review for no string matching. |
| Build | Default tests pass without OpenCode, credentials, network, or real subprocesses. | Run `go test ./...` in `/home/antonioborgerees/coding/ultraplan-go`. |
| Build | CLI builds successfully. | Run `go build ./cmd/ultraplan` in `/home/antonioborgerees/coding/ultraplan-go`. |
| Review | Scope did not expand into public study execution, batch scheduling, durable workflow orchestration, summary generation, code extraction, target workflows, or sprint execution. | Sprint Review Protocol and requirements review. |
| Documentation | Sprint implementation plan executes these decisions without reopening architecture; review records verification outcomes and carry-forward decisions. | `projects/ultraplan-go/sprints/09-runtime-integration/plan.md`; `projects/ultraplan-go/sprints/09-runtime-integration/review.md`. |

## Assumptions And Risks

| Item | Type | Impact | Mitigation / Follow-Up |
| --- | --- | --- | --- |
| Agentwrap docs/source match the sprint architecture: wrapper, health, validation, policy, permission, event, metadata, store/sink, opencode option, and `SDKError` concepts are present with concrete current names. | Confirmed assumption | Implementation can plan against the inspected APIs, but minor compile-time adjustments may still be needed if the sibling checkout changes before implementation. | Use the sibling Agentwrap checkout during implementation; keep tests anchored to public root package APIs and `agentwrap/opencode`. |
| Existing app/config/health foundations can accept minimal runtime additions. | Assumption | If current structures are narrower than expected, minor refactor may be needed. | Keep additions focused; do not broaden architecture unless required by compile/test feedback. |
| In-memory/fake run stores are sufficient for this sprint. | Assumption | No durable agentwrap run inspection across process exits. | Requirements exclude durable runtime records; preserve metadata for later persistence. |
| Runtime mapping layer can become too thick. | Risk | It may duplicate agentwrap or become a competing SDK. | Keep types narrow, generic, and mapped to agentwrap; review for duplicated policy/validation/event logic. |
| Health output can accidentally become execution. | Risk | It could invoke prompts/study workflows and violate scope. | Health tests must prove no prompt execution, batch scheduler, study workflow, or report validation path is invoked. |
| Permission mapping may reveal unsupported adapter features. | Risk | Users could expect stricter policy than OpenCode can enforce. | Fail before launch for unsupported required features; only allow best-effort when explicitly configured. |
| Event channels can block or goroutines can leak. | Risk | Runtime tests or future execution could hang. | Drain or consume events in tests and runtime wait paths; always call `Run.Wait` for started runs. |
| Raw diagnostics can leak secrets. | Risk | Unsafe payloads, prompts, env values, or stderr could enter logs/artifacts. | Redact sensitive config/env values; omit raw native bytes by default; store safe omission facts only. |
| Fake-runtime behavior can diverge from real OpenCode behavior. | Risk | Default tests might miss real adapter integration issues. | Optional gated smoke tests through `agentwrap/opencode`; update fakes when agentwrap changes. |

## Implementation Constraints

- Do not create runtime behavior outside `internal/platform/runtime` except config extensions in `internal/platform/config` and CLI health presentation in `internal/app`.
- Do not import `internal/study`, `internal/codeextract`, or other product modules from `internal/platform/runtime`.
- Do not define a competing public UltraPlan runtime SDK; agentwrap root package APIs are the external runtime boundary.
- Do not launch `opencode` directly, parse OpenCode stdout/stderr directly, or construct OpenCode behavior outside `agentwrap/opencode`.
- Preserve default wrapper order: `ObservingRuntime -> ValidatingRuntime -> PolicyRunner -> opencode.Runtime`.
- Propagate caller `context.Context` through health and runtime requests; avoid `context.Background()` in work paths.
- Resolve and validate runtime config before launch or health checks; reject unknown runtime names, unsupported required health/capability names, empty executable values, invalid timeouts/retries, and unsupported required permission features.
- Redact sensitive runtime config/environment values from config output, logs, health diagnostics, and runtime error diagnostics.
- Preserve agentwrap event kinds and `SDKError` classifications; do not string-match human-readable runtime errors.
- Do not persist or log unsafe raw native payload bytes by default; record safe omission facts when available.
- Keep unknown usage, token, and cost values unknown rather than converting them to zero.
- Delegate retry, rate-limit, fallback, validation wrapping, permission translation, health checks, OpenCode process handling, and native event projection to agentwrap.
- Use agentwrap-compatible fakes for normal tests; any real OpenCode smoke path must be explicitly environment-gated and skipped by default.
- Do not add public study execution, batch scheduling, durable workflow orchestration, summary generation, code extraction, target workflow, or sprint execution scope.

## Plan Handoff

`plan.md` must execute these decisions. It must not invent architecture, scope, or decisions beyond this document.

The plan must carry forward:

- final decisions
- applicable contracts and requirement IDs
- expected evidence
- risks and assumptions
- required review protocols

## Phase Exit Criteria

- [x] Selected context was read and used.
- [x] Area-specific reasoning documents were completed where applicable or explicitly marked not applicable.
- [x] Area-specific reasoning conclusions are reflected or explicitly overridden in final decisions.
- [x] Contracts are explicitly mapped to decisions and expected evidence.
- [x] Final decisions are clear enough for `plan.md` to execute.
- [x] Expected evidence is specific and reviewable.
