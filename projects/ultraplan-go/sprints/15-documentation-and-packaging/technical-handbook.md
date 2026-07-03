# Sprint Technical Handbook: Documentation and Packaging

> Project: `ultraplan-go`
> Sprint: `15-documentation-and-packaging`
> Source: `sprint-index.md`
> **Inputs Used:** `projects/ultraplan-go/sprints/15-documentation-and-packaging/sprint-index.md`, `templates/technical-handbook.md`, `projects/ultraplan-go/project-index.md`, `projects/ultraplan-go/sprints/15-documentation-and-packaging/requirements.md`, `studies/go-cli-study/reports/final/01-project-structure.md`, `studies/go-cli-study/reports/final/02-command-architecture.md`, `studies/go-cli-study/reports/final/03-dependency-injection.md`, `studies/go-cli-study/reports/final/04-configuration-management.md`, `studies/go-cli-study/reports/final/05-error-handling.md`, `studies/go-cli-study/reports/final/06-io-abstraction.md`, `studies/go-cli-study/reports/final/07-state-context.md`, `studies/go-cli-study/reports/final/08-concurrency.md`, `studies/go-cli-study/reports/final/09-terminal-ux.md`, `studies/go-cli-study/reports/final/10-logging-observability.md`, `studies/go-cli-study/reports/final/11-testing-strategy.md`, `studies/go-cli-study/reports/final/12-extensibility.md`, `studies/go-cli-study/reports/final/13-security.md`, `studies/go-cli-study/reports/final/14-performance.md`, `studies/go-cli-study/reports/final/15-philosophy.md`

This handbook distills the studies and reports selected by `sprint-index.md` for sprint reasoning. It does not decide architecture or implementation.

## Selected Studies And Reports

| Study / Report | Path | Relevant Finding | Confidence |
| --- | --- | --- | --- |
| `01-project-structure` | `studies/go-cli-study/reports/final/01-project-structure.md` | Production CLIs keep entrypoints thin and route product behavior into protected/internal layers; examples include `cmd/restic/main.go:37-114`, `cmd/root.go:112-127`, and `cmd/gh/main.go:6`. | high |
| `02-command-architecture` | `studies/go-cli-study/reports/final/02-command-architecture.md` | Public command docs should mirror factory-built command surfaces and shared lifecycle behavior, not ad-hoc command internals; examples include `pkg/cmd/root.go:267-303`, `pkg/cmd/install.go:132-145`, and `cmd/cmd.go:240-340`. | high |
| `03-dependency-injection` | `studies/go-cli-study/reports/final/03-dependency-injection.md` | Release packaging and smoke procedures benefit from explicit composition roots and injectable dependencies rather than hidden globals; examples include `pkg/cmdutil/factory.go:16-43`, `internal/app/app.go:42-81`, and `cmd/restic/main.go:181-183`. | high |
| `04-configuration-management` | `studies/go-cli-study/reports/final/04-configuration-management.md` | Configuration docs must state precedence, validation point, and schema/version behavior explicitly; evidence includes flag restore at `internal/cmd/config.go:2253-2287`, `getConfig` at `internal/flags/flags.go:314-327`, and validation at `internal/config/config.go:609-641`. | high |
| `05-error-handling` | `studies/go-cli-study/reports/final/05-error-handling.md` | Recovery docs should distinguish user-facing diagnostics, operational logs, retry/fatal classes, and exit codes; evidence includes `cmd/age/tui.go:37-54`, `internal/errors/fatal.go:10`, `internal/ghcmd/cmd.go:44-49`, and `fs/fserrors/error.go:22-29`. | high |
| `06-io-abstraction` | `studies/go-cli-study/reports/final/06-io-abstraction.md` | CLI output, smoke evidence, and tests are more reliable when IO is injectable and captured via buffers; evidence includes `pkg/iostreams/iostreams.go:551-568`, `executor.go:553-564`, and `internal/ui/mock.go:10-53`. | high |
| `07-state-context` | `studies/go-cli-study/reports/final/07-state-context.md` | Long-running runs and release smoke should document cancellation, locks, and state recovery; evidence includes `pkg/cmd/install.go:333-347`, `task.go:89`, `internal/restic/lock.go:105`, and `internal/session/session.go:12-23`. | high |
| `08-concurrency` | `studies/go-cli-study/reports/final/08-concurrency.md` | Batch release behavior should preserve bounded workers and explicit cleanup rather than unbounded goroutines; evidence includes `task.go:87`, `pkg/analyze/parallel.go:13`, and `cmd/root.go:261-279`. | high |
| `09-terminal-ux` | `studies/go-cli-study/reports/final/09-terminal-ux.md` | User docs should distinguish interactive UX, non-TTY fallback, progress, and cancellation expectations; evidence includes `internal/cmd/prompt.go:124-137`, `internal/cmd/prompt.go:351-360`, and `signals.go:11-31`. | medium |
| `10-logging-observability` | `studies/go-cli-study/reports/final/10-logging-observability.md` | Smoke evidence and diagnostics should keep user output separate from logs and support structured/debug output without leaking secrets; evidence includes `internal/logging/logging.go:31-66`, `internal/slogs/keys.go:6-231`, and `fs/accounting/prometheus.go:78-108`. | high |
| `11-testing-strategy` | `studies/go-cli-study/reports/final/11-testing-strategy.md` | Release gates should favor offline, repeatable command/integration tests, fixtures, and golden-style evidence; examples include `acceptance/acceptance_test.go:26-29`, `task_test.go:166-169`, and `internal/test/test.go:43`. | high |
| `12-extensibility` | `studies/go-cli-study/reports/final/12-extensibility.md` | Runtime/tool extension should remain through deliberate extension boundaries, not undocumented shell or plugin behavior; evidence includes subprocess runtime at `internal/plugin/runtime_subprocess.go:65-79`, MCP tools at `internal/llm/agent/mcp-tools.go:106-129`, and extension loading at `pkg/cmd/root/root.go:192-203`. | medium |
| `13-security` | `studies/go-cli-study/reports/final/13-security.md` | Docs and evidence must redact secrets, document permission boundaries, and avoid unsafe shell/runtime claims; evidence includes `internal/options/secret_string.go:15-20`, `pkg/registry/transport.go:37-41`, `internal/permission/permission.go:44-108`, and `internal/llm/tools/bash.go:41-55`. | high |
| `14-performance` | `studies/go-cli-study/reports/final/14-performance.md` | Release docs should avoid unsupported performance promises while preserving bounded concurrency, streaming, lazy init, and profiling awareness; evidence includes `pkg/cmdutil/factory.go:27-42`, `age/internal/stream/stream.go:20`, and `gdu/pkg/analyze/parallel.go:13`. | medium |
| `15-philosophy` | `studies/go-cli-study/reports/final/15-philosophy.md` | Production-ready documentation should state non-goals and tradeoffs explicitly; evidence includes no-keyring philosophy at `age.go:18`, plugin-first philosophy at `AGENTS.md:88`, and extension rejection guidance at `VISION.md:97`. | medium |

## Relevant Patterns

- **Thin release surface with internal ownership:** Keep release documentation and packaging tied to the single CLI entrypoint and existing internal package ownership. The structure report shows thin entrypoints delegating inward (`cmd/gh/main.go:6`, `cmd/restic/main.go:37-114`) and protected application packages under `internal/` (`internal/restic/repository.go:18`, `internal/chezmoi/chezmoi.go:1-2`). For this sprint, docs should describe user-visible behavior without implying public Go APIs or target/sprint workflows beyond the shipped CLI.
- **Explicit command reference from command factories and help surfaces:** The command architecture report favors factory-created commands with consistent wiring (`pkg/cmd/install.go:132`, `pkg/cmd/root.go:267-303`, `cmd/restic/cmd_backup.go:35`). The CLI reference should be checked against implemented command construction/help/output modes rather than reconstructed from memory.
- **Layered configuration with documented precedence:** Mature CLIs make precedence explicit and validate after merge. Evidence includes flag save/restore for CLI-over-config at `internal/cmd/config.go:2253-2287`, generic precedence lookup at `internal/flags/flags.go:314-327`, and centralized validation at `internal/config/config.go:609-641`. The configuration guide should name defaults, workspace config, environment variables, command flags, redaction, and schema rejection in one place.
- **Actionable recovery and classified errors:** Recovery material should teach users how to interpret validation failures, locks, retries, fallbacks, and exit behavior. The error report highlights hint-based user messages (`cmd/age/tui.go:47-54`), fatal classification (`internal/errors/fatal.go:10`, `cmd/restic/main.go:205`), and retry/fatal interfaces (`fs/fserrors/error.go:22-29`, `fs/fserrors/error.go:92-99`).
- **Injectable/capturable IO for smoke evidence:** Evidence collection should capture exact commands and outputs without mixing logs and data. `IOStreams.Test()` returns in-memory buffers at `pkg/iostreams/iostreams.go:551-568`, go-task injects stdout/stderr via `executor.go:553-577`, and restic uses mock terminals at `internal/ui/mock.go:10-53`.
- **Context-aware cancellation and durable state language:** Run-loop and recovery docs should make cancellation, stale running tasks, lock refresh, and resume interpretation explicit. The state/context report points to signal-to-context wiring (`pkg/cmd/install.go:333-347`), errgroup cancellation (`task.go:89`), lock modeling (`internal/restic/lock.go:105`), and session state (`internal/session/session.go:12-23`).
- **Offline-first release gates with gated real-runtime smoke:** The testing report supports repeatable command/integration tests using testscript and fixtures (`acceptance/acceptance_test.go:26-29`, `internal/cmd/main_test.go:64-174`) and golden comparisons for stable output (`internal/test/test.go:43`, `task_test.go:166-169`). For this sprint, normal `go test`, race tests, and build checks should stay offline; OpenCode smoke should be documented as opt-in/gated.
- **Security-conscious evidence and docs:** Secret redaction and trust boundaries should be visible in docs and smoke evidence. Evidence includes `SecretString` redaction (`internal/options/secret_string.go:15-20`), auth header scrubbing (`pkg/registry/transport.go:37-41`), safe/blocked command policy (`internal/llm/tools/bash.go:41-55`), and permission prompts (`internal/permission/permission.go:44-108`).
- **Bounded concurrency and resource claims:** For run-all/run-loop docs and smoke evidence, describe bounded execution and cancellation honestly. Evidence includes semaphore limiting (`pkg/analyze/parallel.go:13`, `k9s/internal/pool.go:26-48`), errgroup fan-out (`task.go:87`), and wait-with-timeout cleanup (`cmd/root.go:261-279`).

## Trade-Offs

| Trade-Off | Benefit | Cost | When It Matters |
| --- | --- | --- | --- |
| Comprehensive release docs vs. strict study-side scope | Users get a complete guide, config reference, recovery runbook, CLI reference, and smoke checklist for the first release. | Over-documenting deferred target/sprint workflows would create false product commitments. | README, user guide, and CLI reference must explicitly exclude target scaffolding, sprint planning/execution, hosted service, browser UI, and multi-user collaboration. |
| Gated real OpenCode smoke vs. mandatory default tests | Keeps normal release gates offline, fake-first, repeatable, and usable without credentials; aligns with testscript/fixture patterns (`acceptance/acceptance_test.go:26-29`). | Real-runtime behavior may be untested on machines without OpenCode/provider credentials; smoke evidence must record an honest skip. | `docs/opencode-smoke.md` and `dist/smoke-evidence.md` must separate offline pass/fail from gated real-runtime status. |
| Explicit config precedence vs. implementation complexity | Users can predict which value wins, matching evidence from CLI-over-config patterns (`internal/cmd/config.go:2253-2287`) and generic precedence (`internal/flags/flags.go:314-327`). | Docs must remain synchronized with actual config loading and validation behavior. | Configuration guide, release checklist, and smoke evidence should check help/config behavior before publishing. |
| Rich recovery guidance vs. avoiding implementation decisions | Operators can recover from validation failures, stale locks, missing artifacts, cancellation, and write failures using evidence-backed patterns. | The handbook and docs must not invent new recovery behavior or schema migrations. | Recovery runbook should describe observed/current behavior and defer any changes to sprint reasoning/implementation. |
| Checksums and local packaging vs. full publication pipeline | SHA-256 artifacts give auditable local release outputs without external release risk. | No signing, notarization, GitHub release publication, or artifact upload means distribution trust remains limited. | Release checklist must state local artifact scope and not imply signed/notarized/platform-published releases. |
| Structured diagnostic detail vs. secret minimization | Logs and smoke evidence can support debugging and auditability, as seen in structured logging/reporting patterns (`internal/logging/logging.go:31-66`, `internal/slogs/keys.go:6-231`). | Too much detail can leak prompts, provider tokens, full environment values, or raw runtime payloads. | Smoke evidence and docs must include commands/results/metadata but redact sensitive values. |

## Anti-Patterns And Warnings

- **Do not claim deferred product surfaces:** The sprint scope excludes target scaffolding, sprint planning/execution, hosted service, browser UI, and multi-user collaboration. Philosophy evidence shows strong projects state non-goals explicitly (`age.go:18`, `VISION.md:97`); docs should do the same.
- **Do not document direct OpenCode supervision as UltraPlan-owned behavior:** Runtime integration must stay through agentwrap/OpenCode. The extensibility/security reports show external tool/runtime boundaries should be explicit (`internal/llm/agent/mcp-tools.go:106-129`, `internal/permission/permission.go:44-108`).
- **Do not mix data output, logs, and debug traces in release evidence:** Logging evidence warns that debug output to stdout breaks scripting (`src/core.go:325`) and output abstractions can be bypassed (`cli.go:46`). Smoke evidence should clearly separate command, stdout/stderr summary, result, and redaction notes.
- **Do not put real OpenCode/provider smoke in default tests:** Testing evidence favors offline, repeatable command/integration tests; real external dependencies without gating are flagged as flaky (`dive` Docker fixture warning, network/external dependency caution). OpenCode smoke must remain opt-in and skippable with reason.
- **Do not expose secrets or unsafe runtime payloads:** Security evidence highlights secret redaction types and auth scrubbing (`internal/options/secret_string.go:15-20`, `pkg/registry/transport.go:37-41`). Avoid provider tokens, full sensitive environment dumps, full raw prompts, and unsafe report bodies in docs/evidence.
- **Do not promise performance characteristics not enforced by release gates:** Performance evidence supports bounded concurrency and streaming patterns, but also warns about eager init, unbounded accumulation, and goroutine leaks. Docs should describe supported behavior and known release checks, not unmeasured throughput.
- **Do not let local development `replace` directives silently ship:** The release checklist must include a `go.mod` audit for local `replace github.com/Antonio7098/agentwrap => ../agentwrap` because dependency provenance is release-critical.
- **Do not let packaging artifacts become opaque:** Checksums should cover exactly the four required binaries and no unrelated files; smoke evidence should record build targets and checksum command.

## Examples Worth Inspecting

| Example | Path / Source | Why It Is Useful |
| --- | --- | --- |
| Thin CLI entrypoint | `cmd/gh/main.go:6`, `cmd/restic/main.go:37-114`, `cmd/root.go:112-127` | Helps verify UltraPlan's README and CLI reference describe a single CLI entrypoint delegating to product packages. |
| Command construction and grouping | `pkg/cmd/root.go:267-303`, `pkg/cmd/install.go:132-145`, `cmd/restic/main.go:34-69` | Useful for comparing documented command surfaces and help organization against actual command registration. |
| Config precedence | `internal/cmd/config.go:2253-2287`, `internal/flags/flags.go:314-327`, `value_source.go:94-102` | Useful when writing precedence and override sections in `docs/configuration.md`. |
| Config validation | `internal/config/config.go:609-641`, `internal/config/json/validator.go:1-187`, `cmd/dive/cli/internal/options/analysis.go:48-53` | Useful for explaining when invalid config/schema values fail. |
| User-facing errors and hints | `cmd/age/tui.go:37-54`, `internal/ghcmd/cmd.go:281-301`, `go-task/errors/errors_task.go:13-32` | Useful for recovery runbook wording and actionable diagnostics. |
| Lock/session recovery | `internal/restic/lock.go:105`, `internal/restic/lock.go:290-305`, `internal/session/session.go:12-23` | Useful for stale lock, cancellation, and durable run-loop guidance. |
| Injectable IO and test capture | `pkg/iostreams/iostreams.go:551-568`, `executor.go:553-577`, `internal/ui/mock.go:10-53` | Useful for designing and documenting smoke evidence capture. |
| Offline integration tests | `acceptance/acceptance_test.go:26-29`, `internal/cmd/main_test.go:64-174`, `task_test.go:166-169` | Useful for release checklist structure and repeatable test evidence. |
| Secret redaction and auth scrubbing | `internal/options/secret_string.go:15-20`, `pkg/registry/transport.go:37-41`, `internal/ghcmd/cmd.go:281-301` | Useful for smoke evidence redaction review. |
| Bounded concurrency | `pkg/analyze/parallel.go:13`, `k9s/internal/pool.go:26-48`, `task.go:87` | Useful for run-all/run-loop documentation and avoiding unbounded promises. |

## Design Pressures

- The release is documentation-and-packaging heavy, so accuracy must come from implemented help, stable JSON surfaces, command behavior, and release gates rather than aspirational roadmap language.
- The project is study-side only for this release; all docs must resist leakage from deferred target/sprint workflow language.
- Packaging needs four local binaries and exact checksums, but publication, signing, notarization, tags, and uploads remain out of scope.
- OpenCode real-runtime smoke is valuable but environment-dependent; release evidence must be honest about pass, fail, or skip.
- Stable JSON surfaces from Sprint 14 are compatibility-sensitive; undocumented/debug-only shapes must not be promised as stable.
- Recovery docs must balance operator usefulness with not inventing new behavior around locks, schema migration, retry/fallback, cancellation, or missing artifacts.
- Security constraints require redaction and minimal evidence retention while still preserving enough information for auditability.
- Docs should preserve product/platform separation: product behavior in study docs, generic runtime notes around agentwrap/OpenCode, and no product-owned direct OpenCode supervision claims.
- Release check commands must remain feasible on a normal local development machine without network, OpenCode, provider credentials, or external smoke fixtures.

## Open Questions For Reasoning

- Which CLI JSON outputs are officially stable from Sprint 14, and which should be described only as text/debug/unstable behavior?
- What exact command examples should the user guide and CLI reference include so they match current help output and avoid unsupported target/sprint claims?
- How should `dist/smoke-evidence.md` redact environment and runtime metadata while still proving which gates ran and what result occurred?
- What is the exact skip reason format for gated OpenCode smoke when OpenCode, provider credentials, or network are unavailable?
- Should the release checklist require rerunning checksum verification after binary generation, or is one recorded generation command sufficient for this sprint?
- How should docs phrase `--force-unlock` risk so users understand stale-lock recovery without being encouraged to corrupt active runs?
- What local `replace` directive disposition is acceptable if agentwrap is still under local development at release time?
- Which implemented config fields and schema-version rejection paths should be documented, and are any fields internal/debug-only?
- How should recovery docs interpret retry/fallback metadata when a run partially completes but validation later fails?
- Are packaging commands manual release steps, scripts, or documented command snippets for this sprint?

## Evidence Pointers

- `studies/go-cli-study/reports/final/01-project-structure.md`: Inspect thin CLI and internal boundary evidence around `cmd/`, `internal/`, and unidirectional dependency flow.
- `studies/go-cli-study/reports/final/02-command-architecture.md`: Inspect factory functions, command grouping, lifecycle hooks, and command reference implications.
- `studies/go-cli-study/reports/final/04-configuration-management.md`: Inspect precedence, validation hooks, env var naming, XDG/config location, and schema/version cautions.
- `studies/go-cli-study/reports/final/05-error-handling.md`: Inspect user/operational error separation, fatal/retry classification, hints, and exit code mapping.
- `studies/go-cli-study/reports/final/06-io-abstraction.md`: Inspect injectable streams, mock terminals, filesystem abstractions, and smoke/test capture patterns.
- `studies/go-cli-study/reports/final/07-state-context.md`: Inspect context propagation, signal cancellation, explicit session/lock state, and stale/cancel recovery risks.
- `studies/go-cli-study/reports/final/08-concurrency.md`: Inspect errgroup, semaphore, stop-channel, timeout cleanup, and goroutine leak cautions.
- `studies/go-cli-study/reports/final/10-logging-observability.md`: Inspect stdout/stderr separation, structured logs, debug control, and evidence discipline.
- `studies/go-cli-study/reports/final/11-testing-strategy.md`: Inspect testscript, golden files, fixtures, integration boundaries, and external dependency warnings.
- `studies/go-cli-study/reports/final/13-security.md`: Inspect secret redaction, command execution trust boundaries, permission prompts, and credential scrubbing.
- `studies/go-cli-study/reports/final/14-performance.md`: Inspect lazy init, streaming, bounded concurrency, profiling, and performance anti-patterns.
- `studies/go-cli-study/reports/final/15-philosophy.md`: Inspect explicit non-goals, simplicity/extensibility tradeoffs, and warning signs around undocumented complexity.
- `projects/ultraplan-go/sprints/15-documentation-and-packaging/requirements.md`: Inspect required outputs, acceptance criteria, constraints, dependencies, and review expectations.

## Handoff To Reasoning

- Use this handbook as evidence input.
- Validate whether the observed patterns fit this project's constraints.
- Do not copy external patterns without sprint-specific reasoning.
- Keep final implementation decisions in sprint reasoning, not in this handbook.
- Treat cited study source paths and line numbers as the evidence base for documentation, packaging, smoke, and recovery tradeoffs.
