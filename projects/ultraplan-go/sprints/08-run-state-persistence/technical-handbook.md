# Sprint Technical Handbook: 08 Run State Persistence

> Project: `ultraplan-go`
> Sprint: `08-run-state-persistence`
> Source: `sprint-index.md`
> **Inputs Used:** `projects/ultraplan-go/project-index.md`, `projects/ultraplan-go/sprints/08-run-state-persistence/sprint-index.md`, `templates/technical-handbook.md`, `projects/ultraplan-go/sprints/08-run-state-persistence/requirements.md`, `projects/ultraplan-go/docs/ARCHITECTURE.md`, `projects/ultraplan-go/docs/PRD.md`, `projects/ultraplan-go/docs/TRD.md`, `studies/go-cli-study/reports/final/01-project-structure.md`, `studies/go-cli-study/reports/final/02-command-architecture.md`, `studies/go-cli-study/reports/final/03-dependency-injection.md`, `studies/go-cli-study/reports/final/04-configuration-management.md`, `studies/go-cli-study/reports/final/05-error-handling.md`, `studies/go-cli-study/reports/final/06-io-abstraction.md`, `studies/go-cli-study/reports/final/07-state-context.md`, `studies/go-cli-study/reports/final/09-terminal-ux.md`, `studies/go-cli-study/reports/final/10-logging-observability.md`, `studies/go-cli-study/reports/final/11-testing-strategy.md`, `studies/go-cli-study/reports/final/13-security.md`, `studies/go-cli-study/reports/final/14-performance.md`, `studies/go-cli-study/reports/final/15-philosophy.md`

This handbook distills the studies and reports selected by `sprint-index.md` for sprint reasoning. It does not decide architecture or implementation.

## Selected Studies And Reports

| Study / Report | Path | Relevant Finding | Confidence |
| --- | --- | --- | --- |
| `01-project-structure` | `studies/go-cli-study/reports/final/01-project-structure.md` | Mature Go CLIs keep entrypoints thin and maintain unidirectional dependencies from CLI to business logic; `internal/` is the dominant application boundary (`main.go:26-34`, `cmd/gh/main.go:6`, `internal/chezmoi/chezmoi.go:1-2`, `cmd/root.go:112-127`). | high |
| `02-command-architecture` | `studies/go-cli-study/reports/final/02-command-architecture.md` | Strong command designs use factory-created commands that parse flags, construct dependencies, delegate, and format output rather than embedding business logic in `RunE` (`pkg/cmd/install.go:132-145`, `pkg/cmdutil/factory.go:16-43`, `internal/cmd/config.go:1833-1845`). | high |
| `03-dependency-injection` | `studies/go-cli-study/reports/final/03-dependency-injection.md` | Manual composition roots, constructor injection, and small interfaces improve testability; global config/service state is the recurring failure mode (`pkg/cmdutil/factory.go:16-43`, `executor.go:22-24`, `internal/chezmoi/system.go:25`, `fs/config.go:793`). | high |
| `04-configuration-management` | `studies/go-cli-study/reports/final/04-configuration-management.md` | Layered config is safest when precedence and validation happen in one explicit path; `Flag.Changed`/restore patterns prevent defaults from shadowing user choices (`internal/cmd/config.go:2253-2287`, `internal/flags/flags.go:314-327`, `internal/global/global.go:139,147`). | medium |
| `05-error-handling` | `studies/go-cli-study/reports/final/05-error-handling.md` | Robust CLIs wrap errors with `%w`, use sentinels or typed errors for decisions, and render user-facing diagnostics separately from operational logs (`fs/fserrors/error.go:22-29`, `errors/errors_task.go:13-32`, `cmd/age/tui.go:37-54`, `internal/ghcmd/cmd.go:281-301`). | high |
| `06-io-abstraction` | `studies/go-cli-study/reports/final/06-io-abstraction.md` | Injectable IO and filesystem seams enable command tests and persistence tests without real terminals or external services (`pkg/iostreams/iostreams.go:551-568`, `internal/chezmoi/system.go:25`, `internal/fs/interface.go:10-31`, `internal/ui/mock.go:10-53`). | high |
| `07-state-context` | `studies/go-cli-study/reports/final/07-state-context.md` | Long-running/stateful tools benefit from root context propagation and centralized application state, but context-as-service-locator and global fallbacks hide dependencies (`pkg/cmd/install.go:333-347`, `task.go:89`, `internal/app/app.go:25-40`, `fs/config.go:793-802`). | medium |
| `09-terminal-ux` | `studies/go-cli-study/reports/final/09-terminal-ux.md` | Status UX should be calm, concise, non-blocking, and TTY-aware; progress and status are valuable for long work, but CLI-first tools should keep output scriptable (`internal/cmd/prompt.go:20-256`, `pkg/iostreams/iostreams.go:514-516`, `internal/ui/terminal.go:10-34`). | medium |
| `10-logging-observability` | `studies/go-cli-study/reports/final/10-logging-observability.md` | Useful observability comes from structured fields, runtime debug control, and stdout/stderr separation; logs should not replace deterministic status output (`internal/slogs/keys.go:6-231`, `internal/logging/logging.go:31-66`, `pkg/iostreams/iostreams.go:52-54`). | medium |
| `11-testing-strategy` | `studies/go-cli-study/reports/final/11-testing-strategy.md` | Strong CLIs combine table-driven unit tests, command-level integration tests, golden/output regression tests, and reusable fakes (`internal/cmd/main_test.go:64-174`, `task_test.go:166-169`, `internal/backend/mock/backend.go:14-26`, `pkg/httpmock/stub.go:35-199`). | high |
| `13-security` | `studies/go-cli-study/reports/final/13-security.md` | Durable artifacts and diagnostics must avoid secrets, validate paths, and treat user/workspace inputs as trust boundaries (`internal/options/secret_string.go:15-20`, `pkg/registry/transport.go:37-41`, `gpgencryption.go:151-165`, `internal/config/json/validator.go:146`). | high |
| `14-performance` | `studies/go-cli-study/reports/final/14-performance.md` | Fast status commands should avoid eager initialization, recursive scans, and unbounded memory; lazy factories and streaming/bounded work are proven patterns (`pkg/cmdutil/factory.go:27-42`, `yq/pkg/yqlib/stream_evaluator.go:78-113`, `gdu/pkg/analyze/parallel.go:13`). | high |
| `15-philosophy` | `studies/go-cli-study/reports/final/15-philosophy.md` | Coherent tools deliberately reject non-goals and keep complexity near the capability that needs it; wrapper and interface patterns support dry-run/debug/test behavior without scattering policy (`internal/chezmoi/system.go:25-45`, `pkg/cmd/factory/default.go:26-46`, `internal/pubsub/broker.go:10-19`). | medium |

## Relevant Patterns

- **Study-owned state near study behavior:** The structure report favors protected application internals and unidirectional flow (`internal/chezmoi/chezmoi.go:1-2`, `cmd/root.go:112-127`), while the project architecture says product modules own their state, validation, scheduling, and persistence. For this sprint, the pressure is to keep run-state schema, validation, status summaries, and persistence inside `internal/study`, not a global `internal/state` or `internal/reports` layer.
- **Thin status command, non-CLI state logic:** Command evidence favors factory-created commands that only parse flags, call business logic, and render results (`pkg/cmd/install.go:132-145`, `pkg/cmd/issue/list/list.go:47-118`, `cmd/restic/main.go:37-114`). The `ultraplan study <study> status` command should therefore be a rendering and diagnostic boundary, with load/validate/summarize logic outside command handlers.
- **Explicit durable state schema with version rejection:** Configuration and persistence-adjacent evidence emphasizes schema validation and explicit precedence/validation points (`internal/config/config.go:609-641`, `internal/config/json/validator.go:1-187`). Run-state loading should have one clear schema-version validation path and reject unsupported versions instead of silently coercing them.
- **Atomic file write seam with injectable filesystem/time where useful:** IO evidence shows that filesystem-heavy tools use `System` or `FS` interfaces and test constructors (`internal/chezmoi/system.go:25`, `internal/fs/interface.go:10-31`, `pkg/iostreams/iostreams.go:551-568`). The sprint can use local files directly where existing project style supports it, but atomic save/load needs testable seams for temp dirs, failures, and clock/timestamp determinism.
- **Resume validation before trusting completed state:** Product evidence says runtime success is not product success, and error/report evidence supports validation as a first-class gate. Relevant report validators and state/context examples point toward rechecking persisted claims before acting on them (`internal/checker/checker.go:70-91`, `pkg/storage/driver/driver.go:27-36`, `internal/restic/lock.go:105`). Completed tasks should be trusted only after existing report validation confirms required artifacts.
- **Deterministic summaries from persisted state, not runtime inspection:** Observability evidence separates diagnostics from user output and favors stable fields (`internal/slogs/keys.go:6-231`, `fs/accounting/prometheus.go:78-108`). For this sprint, status summaries should derive from `run-state.json` and known report metadata, not active processes, agentwrap stores, network, or recursive repository scans.
- **Actionable, wrapped diagnostics:** Error evidence favors `%w`, sentinels, typed errors, and hints (`fs/fserrors/error.go:22-29`, `pkg/cmdutil/errors.go:35-41`, `cmd/age/tui.go:47-54`). Missing state, malformed JSON, unsupported schema, invalid completed output, and failed atomic writes should retain cause chains and name the state/report path.
- **Secret-free persisted summaries:** Security evidence shows dedicated secret redaction types and header scrubbing (`internal/options/secret_string.go:15-20`, `pkg/registry/transport.go:37-41`, `status.go:332-338`). Config summary fields in durable state should be safe operational metadata, not provider credentials, full environment, tokens, or unsafe raw runtime payloads.
- **Table-driven plus command-output tests:** Testing evidence supports table-driven unit tests, fakes, golden/output checks, and behavior assertions (`internal/cmd/main_test.go:64-174`, `helm/internal/test/test.go:43`, `task_test.go:166-169`). State construction, atomic writes, malformed files, status counts, and command output need focused fixture tests with no runtime dependency.

## Trade-Offs

| Trade-Off | Benefit | Cost | When It Matters |
| --- | --- | --- | --- |
| Concrete local filesystem helpers vs. broad FS interface | Smaller implementation and easier alignment with current Go project style; avoids premature abstraction | Harder to simulate rename/write failures unless helper seams are chosen carefully | Atomic write and corrupted/missing-state tests need controlled failures; evidence from `internal/chezmoi/system.go:25` and `internal/fs/interface.go:10-31` supports abstraction only where side effects are core |
| Missing inapplicable Markdown pairs vs. explicit skipped tasks | Smaller state file and simpler counts if inapplicable pairs are absent | Users may have less visibility into why a pair is not present; status logic must avoid interpreting absence as missing work | Sprint requirements allow skipping/absence but demand deterministic summaries; applicability evidence in requirements and TRD must be reconciled with status UX |
| Rich task metadata now vs. minimal placeholders for future runtime | Future run-loop/agentwrap integration has a stable schema and fewer migrations | More fields to validate and test before runtime execution exists; risk of speculative fields | Requirements require agentwrap metadata placeholders and attempt/error fields, but runtime execution is out of scope |
| Strict schema rejection vs. permissive compatibility | Safer durable format, clear migration boundary, easier debugging | Older state files stop working until migration exists | The configuration/security evidence favors explicit schema validation (`internal/config/config.go:609-641`, `internal/config/json/validator.go:1-187`), and sprint requirements require unsupported-version rejection |
| Human concise status vs. machine-readable detail | Calm CLI output aligned with terminal UX evidence and sprint scope | Automation may later need stable JSON, which is explicitly out of scope now | `status` must be deterministic and useful without promising public JSON compatibility |
| Lazy status loading vs. eager workspace/study validation | Faster status for large workspaces and less accidental scanning | Some unrelated workspace/config errors may not surface during status | Performance evidence warns against eager init (`pkg/cmdutil/factory.go:27-42`) and recursive scans; status should inspect only what is needed |

## Anti-Patterns And Warnings

- **Do not introduce global technical-layer packages for product state:** The project structure report warns that global layers fracture product context, and the architecture explicitly keeps study state in `internal/study`; avoid `internal/state`, `internal/scheduler`, `internal/reports`, or platform imports of `study`.
- **Do not put state construction or validation in `RunE`:** Command architecture evidence flags long `RunE` handlers as logic leakage (`opencode cmd/root.go:49-183`, `yq cmd/evaluate_sequence_command.go:152`). Command handlers should render and route only.
- **Do not silently accept malformed, unknown-version, or partially written state:** Error and config evidence favors explicit validation and wrapped diagnostics (`pkg/storage/driver/driver.go:27-36`, `internal/config/config.go:609-641`). Silent fallback would hide data loss.
- **Do not trust completed tasks without artifact validation:** Product requirements and validation evidence reject runtime-success-only completion; persisted `completed` should be revalidated through existing report helpers before status/resume logic treats it as true.
- **Do not let config summaries persist secrets:** Security evidence shows credentials need type-level or transport-level scrubbing (`internal/options/secret_string.go:15-20`, `pkg/registry/transport.go:37-41`). Avoid env dumps, provider tokens, raw stderr, and unsafe runtime payloads.
- **Do not make status inspect runtime processes or invoke agentwrap/OpenCode:** This sprint is persistence/status only. Observability evidence supports deterministic status records, but runtime inspection and event sinks belong to later runtime work.
- **Do not recursively scan source repositories for status:** Performance evidence warns against unbounded traversal and eager init (`gdu/pkg/analyze/parallel.go:13`, `yq/pkg/yqlib/stream_evaluator.go:78-113`). Status should read known state and known study/report metadata.
- **Do not treat inapplicable Markdown pairs as missing or failed:** Scope evidence requires applicability filtering everywhere task pairs are created, validated, summarized, or synthesized. Inapplicable pairs must be absent or skipped, not failures.

## Examples Worth Inspecting

| Example | Path / Source | Why It Is Useful |
| --- | --- | --- |
| Go-enforced internal boundary | `internal/chezmoi/chezmoi.go:1-2`; `cmd/root.go:112-127` | Shows CLI-to-domain dependency direction with no reverse command imports. |
| Command factory and delegation | `pkg/cmd/install.go:132-145`; `pkg/cmd/issue/list/list.go:47-118`; `cmd/restic/main.go:37-114` | Useful for keeping `study status` command thin while state logic stays testable. |
| Factory DI with lazy dependencies | `pkg/cmdutil/factory.go:16-43`; `pkg/cmdutil/factory.go:27-42` | Useful for avoiding runtime/provider initialization during status commands. |
| IO test constructor | `pkg/iostreams/iostreams.go:551-568`; `ui_mock.go:27-33` | Good model for deterministic command-output tests. |
| Filesystem abstraction | `internal/chezmoi/system.go:25`; `internal/fs/interface.go:10-31` | Useful when deciding how much indirection atomic persistence tests require. |
| Typed/sentinel errors | `fs/fserrors/error.go:22-29`; `go-task/errors/errors_task.go:13-32`; `pkg/storage/driver/driver.go:27-36` | Useful for missing vs malformed vs unsupported state diagnostics. |
| Secret redaction | `internal/options/secret_string.go:15-20`; `pkg/registry/transport.go:37-41`; `status.go:332-338` | Useful for config summary and error metadata safety. |
| testscript/golden testing | `internal/cmd/main_test.go:64-174`; `task_test.go:166-169`; `helm/internal/test/test.go:43` | Useful for command behavior, output stability, and fixture-driven persistence tests. |
| Lazy initialization | `pkg/cmdutil/factory.go:27-42`; `helm/pkg/action/lazyclient.go:35-53` | Useful for status command performance and avoiding runtime side effects. |

## Design Pressures

- Run state is durable product state, so schema shape, validation, atomic write behavior, and diagnostics need more discipline than ephemeral command output.
- State behavior belongs to the study module, while CLI output belongs to the app command layer; this creates pressure to define a small, stable study API for status and persistence.
- Status must be deterministic and truthful without runtime execution, which pressures the design toward persisted status fields plus report validation, not live process or agentwrap inspection.
- Applicability filtering must remain consistent across task construction, synthesis dependency recording, completed-output revalidation, and status summaries.
- Atomic writes must preserve the previous valid file on failure while still surfacing enough error context for users to recover.
- Config summary is useful for diagnostics, but security evidence pressures it to be deliberately narrow and redacted.
- The task schema needs future runtime placeholders without letting later runtime concepts drive current implementation into speculative scheduler or agentwrap integration.
- Tests need to cover persistence failure modes and command output without requiring OpenCode, network, provider credentials, or large source repositories.

## Open Questions For Reasoning

- Should inapplicable Markdown source/dimension pairs be absent from `Tasks`, represented as `skipped`, or recorded only in summary metadata?
- What exact fields belong in `ConfigSummary` so it is useful for diagnosing status while remaining secret-free and stable?
- What should stale active statuses (`running`, `validating`, `waiting`, `retrying`) reset to on resume, and which attempt/error/timestamp fields must be preserved?
- Should `LoadRunState` distinguish missing state through a sentinel error, a typed error, or a result enum, and how should the CLI map that to exit behavior?
- How much filesystem abstraction is necessary for reliable atomic-write failure tests without creating a broad platform filesystem API?
- Should schema version be a top-level integer only, or include a format identifier for clearer unsupported-file diagnostics?
- How should synthesis task dependencies identify applicable source report inputs: source name, report path, task ID, validation state, or a combination?
- What deterministic task ID format balances readability, stability, and future compatibility if dimension/source slugs change?
- Should status output include validation warnings for completed tasks that become invalid during resume validation, or should that remain in state diagnostics only?

## Evidence Pointers

- `studies/go-cli-study/reports/final/01-project-structure.md`: Inspect `Pattern 1`, `Pattern 2`, `Pattern 4`, and anti-patterns for module boundary and CLI/domain separation evidence.
- `studies/go-cli-study/reports/final/02-command-architecture.md`: Inspect `Thin-Delegate Pattern`, `Factory Function Command Creation`, and `RunE functions over 100 lines` for status command wiring discipline.
- `studies/go-cli-study/reports/final/03-dependency-injection.md`: Inspect `Constructor Injection with Factory Function`, `Interface Abstraction for Core Services`, and `Global State Usage` for testable state services.
- `studies/go-cli-study/reports/final/04-configuration-management.md`: Inspect `PostLoad Validation Hook`, `Flag.Changed tracking`, and config validation cautions for schema and config summary handling.
- `studies/go-cli-study/reports/final/05-error-handling.md`: Inspect `%w` wrapping, sentinel errors, typed errors, and user/operational separation for load/save/status diagnostics.
- `studies/go-cli-study/reports/final/06-io-abstraction.md`: Inspect `IOStreams with Test Constructor`, `System Interface for Filesystem`, and hardcoded `os.*` anti-patterns for persistence and command tests.
- `studies/go-cli-study/reports/final/07-state-context.md`: Inspect centralized app state, signal/context patterns, and service-locator cautions when deciding state service shape.
- `studies/go-cli-study/reports/final/09-terminal-ux.md`: Inspect non-TTY fallback, concise status/progress patterns, and stdout/stderr cautions for human status output.
- `studies/go-cli-study/reports/final/10-logging-observability.md`: Inspect structured keys, stdout/stderr separation, and debug control for deterministic status vs logs.
- `studies/go-cli-study/reports/final/11-testing-strategy.md`: Inspect testscript, golden files, centralized mocks, and behavior-vs-implementation warnings for sprint test plan.
- `studies/go-cli-study/reports/final/13-security.md`: Inspect secret redaction, path trust boundaries, private temp dirs, and schema validation for safe persisted state.
- `studies/go-cli-study/reports/final/14-performance.md`: Inspect lazy initialization, streaming/bounded work, and eager init warnings for fast status loading.
- `studies/go-cli-study/reports/final/15-philosophy.md`: Inspect explicit non-goals and deliberate complexity patterns to keep runtime execution, locks, and scheduler behavior out of this sprint.

## Handoff To Reasoning

- Use this handbook as evidence input.
- Validate whether the observed patterns fit this project's constraints.
- Do not copy external patterns without sprint-specific reasoning.
- Resolve the open questions before implementation choices are finalized.
- Keep runtime execution, worker pools, retry execution, locks, target workflows, and public JSON status out of this sprint unless later reasoning explicitly changes scope.
