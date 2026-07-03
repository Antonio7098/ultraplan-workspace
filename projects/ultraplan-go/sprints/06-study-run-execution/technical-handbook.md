# Sprint Technical Handbook: 06 Study Run Execution

> Project: `ultraplan-go`
> Sprint: `06-study-run-execution`
> Source: `sprint-index.md`
> **Inputs Used:** `projects/ultraplan-go/sprints/06-study-run-execution/sprint-index.md`, `templates/technical-handbook.md`, `projects/ultraplan-go/project-index.md`, `projects/ultraplan-go/sprints/06-study-run-execution/requirements.md`, `projects/ultraplan-go/docs/PRD.md`, `projects/ultraplan-go/docs/TRD.md`, `projects/ultraplan-go/docs/ARCHITECTURE.md`, `studies/go-cli-study/reports/final/01-project-structure.md`, `studies/go-cli-study/reports/final/05-error-handling.md`, `studies/go-cli-study/reports/final/06-io-abstraction.md`, `studies/go-cli-study/reports/final/10-logging-observability.md`, `studies/go-cli-study/reports/final/11-testing-strategy.md`, `studies/go-cli-study/reports/final/13-security.md`, `studies/go-cli-study/reports/final/14-performance.md`, `studies/go-cli-study/reports/final/15-philosophy.md`

This handbook distills the studies and reports selected by `sprint-index.md` for sprint reasoning. It does not decide architecture or implementation.

## Selected Studies And Reports

| Study / Report | Path | Relevant Finding | Confidence |
| --- | --- | --- | --- |
| `go-cli-study / 01-project-structure` | `studies/go-cli-study/reports/final/01-project-structure.md` | Mature Go CLIs keep entrypoints thin and product behavior in protected modules; `cmd/` imports business logic but business logic does not import `cmd/` (`main.go:16`, `cmd/root.go:112-127`, `internal/restic/repository.go:18`). | high |
| `go-cli-study / 05-error-handling` | `studies/go-cli-study/reports/final/05-error-handling.md` | Validation diagnostics should preserve cause chains with `%w`, expose programmatic categories when callers branch, and include actionable user guidance (`rclone/vfs/zip.go:21`, `gh-cli/pkg/ssh/ssh_keys.go:64`, `age/cmd/age/tui.go:47-54`). | high |
| `go-cli-study / 06-io-abstraction` | `studies/go-cli-study/reports/final/06-io-abstraction.md` | Testable CLIs inject filesystem and IO boundaries or keep IO behind small seams; direct `os.*` use in core paths makes deterministic tests harder (`pkg/iostreams/iostreams.go:551-568`, `internal/chezmoi/system.go:25`, `internal/fs/interface.go:10-31`). | high |
| `go-cli-study / 10-logging-observability` | `studies/go-cli-study/reports/final/10-logging-observability.md` | Useful diagnostics separate user output from operational detail and use stable structured fields where logs or status are machine-consumed (`internal/logging/logging.go:31-66`, `internal/slogs/keys.go:6-231`, `pkg/iostreams/iostreams.go:52-54`). | medium |
| `go-cli-study / 11-testing-strategy` | `studies/go-cli-study/reports/final/11-testing-strategy.md` | Elite projects rely on table-driven tests, fixtures, golden or snapshot checks, and centralized fakes where multiple tests share the same boundary (`chezmoi/internal/cmd/main_test.go:64-174`, `go-task/task_test.go:166-169`, `helm/internal/test/test.go:43`). | high |
| `go-cli-study / 13-security` | `studies/go-cli-study/reports/final/13-security.md` | Path, secret, and trust boundaries should be explicit; redaction types, credential scrubbing, and schema/path validation reduce accidental leakage and unsafe operations (`internal/options/secret_string.go:15-20`, `pkg/registry/transport.go:37-41`, `internal/config/json/validator.go:146`). | high |
| `go-cli-study / 14-performance` | `studies/go-cli-study/reports/final/14-performance.md` | Validation should avoid unnecessary scans and unbounded work; mature tools defer expensive work and stream or bound large operations (`gh-cli/pkg/cmdutil/factory.go:27-42`, `yq/pkg/yqlib/stream_evaluator.go:78-113`, `gdu/pkg/analyze/parallel.go:13`). | medium |
| `go-cli-study / 15-philosophy` | `studies/go-cli-study/reports/final/15-philosophy.md` | Strong Go CLIs accept complexity deliberately and keep behavior near domain ownership; factory/options, interface-driven abstraction, and explicit non-goals are recurring maintainability patterns (`pkg/cmd/factory/default.go:26-46`, `internal/chezmoi/system.go:25-45`, `VISION.md:97`). | medium |

## Relevant Patterns

- **Study-owned report validation:** The project-structure evidence favors protected application internals and unidirectional flow: CLI delegates to business logic, while internals never import CLI (`cmd/root.go:112-127`, `internal/restic/repository.go:18`). For this sprint, the pressure is to keep report paths, rating parsing, validation checks, and diagnostics with the study module rather than introducing global `validation`, `reports`, or `parsing` packages.
- **Validation results as domain values:** Error-handling evidence shows mature tools use structured error values, sentinel checks, typed errors, and preserved unwrap chains rather than plain strings (`go-task/errors/errors.go:47-50`, `helm/pkg/storage/driver/driver.go:39-48`, `restic/internal/errors/fatal.go:10`). This applies to validation checks that must expose pass/fail status, severity, path, expected value, observed value, and guidance without requiring callers to parse text.
- **Actionable diagnostics with preserved causes:** `%w` wrapping is a common high-quality pattern across rclone, gh-cli, helm, and opencode (`rclone/vfs/zip.go:21`, `gh-cli/pkg/ssh/ssh_keys.go:64`, `helm/pkg/storage/driver/secrets.go:76`, `opencode/internal/llm/agent/agent.go:238`). Sprint reasoning should preserve filesystem and parsing causes while also adding report path, check name, and corrective guidance.
- **Source-kind-aware validation rules:** Security and PRD/TRD context create a trust-boundary distinction between repository directory sources and embedded Markdown document sources. Evidence for explicit trust boundaries appears in k9s schema validation (`internal/config/json/validator.go:146`) and yq security preferences (`pkg/yqlib/security_prefs.go:3-7`); here, source kind is the boundary that should prevent Markdown documents from inheriting code-citation requirements by default.
- **Small test seams over broad abstraction:** IO evidence favors focused seams where tests need them: `IOStreams.Test()` for in-memory output (`pkg/iostreams/iostreams.go:551-568`), `System` for filesystem behavior (`internal/chezmoi/system.go:25`), and `MockTerminal` for terminal output (`internal/ui/mock.go:10-53`). For local report validation, reasoning should prefer minimal filesystem/test helpers already consistent with the repo over introducing broad platform abstractions.
- **Fixture-first validator tests:** Testing evidence repeatedly uses reproducible fixtures and golden comparisons for generated artifacts (`dive/cmd/dive/cli/testdata/snapshots/cli_build_test.snap:1`, `go-task/task_test.go:166-169`, `helm/internal/test/test.go:43`). This sprint needs deterministic fixtures for valid reports, missing files, missing sections, rating variants, ambiguous ratings, directory citation failures, Markdown citation exemptions, and final report failures.
- **Bounded local artifact processing:** Performance evidence warns against eager, unbounded processing and favors lazy or streaming work (`gh-cli/pkg/cmdutil/factory.go:27-42`, `age/internal/stream/stream.go:20`, `yq/pkg/yqlib/stream_evaluator.go:78-113`). Report validation should inspect expected report files and required content directly, not scan whole source repositories or resolve citations beyond shape checks.
- **Safe diagnostics and redaction discipline:** Security and observability evidence show secrets should be scrubbed before logging or display (`internal/options/secret_string.go:15-20`, `pkg/registry/transport.go:37-41`, `gh-cli/status.go:332-338`). Validation diagnostics may include paths, check names, and safe excerpts, but should avoid raw environment, runtime payloads, or report content beyond what is necessary for correction.

## Trade-Offs

| Trade-Off | Benefit | Cost | When It Matters |
| --- | --- | --- | --- |
| Keep validators in `internal/study` instead of a global validation package | Matches Go CLI boundary evidence and project architecture; keeps report semantics near study state (`internal/chezmoi/chezmoi.go:1-2`, `cmd/root.go:112-127`) | May duplicate future validation shapes if another module later needs similar checks | Matters now because report validation is study-owned and target/sprint workflows are non-goals |
| Structured validation diagnostics vs plain `error` strings | Enables later validate/status/summary commands to inspect check names, severity, path, expected/observed values, and guidance; mirrors typed error evidence (`go-task/errors/errors.go:47-50`, `helm/pkg/storage/driver/driver.go:39-48`) | More domain types and test assertions to maintain | Matters because acceptance criteria require deterministic diagnostics and later API reuse |
| Citation shape validation now vs full citation resolution later | Satisfies directory-source validation without scanning repositories; aligns performance guidance to avoid unnecessary scans (`yq/pkg/yqlib/stream_evaluator.go:78-113`, `gdu/pkg/analyze/parallel.go:13`) | Cannot prove cited files exist or snippets are correct until code extraction sprint | Matters because code reference extraction is explicitly deferred while directory report shape is required |
| Minimal filesystem seams vs broad FS interface | Keeps current sprint small while still testable with temp files/fixtures; follows interface guidance to add seams at volatile boundaries (`internal/chezmoi/system.go:25`, `internal/fs/interface.go:10-31`) | Some tests may use real temporary files instead of in-memory FS | Matters if validation is only reading expected local report files and not abstracting an entire workspace filesystem |
| Golden/fixture tests vs inline strings | Fixture diffs make report-shape regressions visible and mirror mature CLI testing (`helm/internal/test/test.go:43`, `rclone/cmd/bisync/bisync_test.go:1435-1479`) | Fixtures can become stale and require review discipline | Matters because validators parse Markdown structures where subtle formatting changes matter |
| Strict rating ambiguity handling vs permissive coercion | Prevents missing or ambiguous ratings from becoming false `0` scores; matches error-handling and product principle that output validation gates success | Some generated reports may fail validation and require manual or later repair | Matters because summary generation later depends on distinguishing missing, invalid, ambiguous, and real zero-like values |

## Anti-Patterns And Warnings

- **Do not parse errors or diagnostics by string:** Error-handling evidence warns against no hierarchy and stringly-typed sentinels (`k9s/internal/client/errors.go:9-14`, `yq`, `fzf`). Validation reasoning should define inspectable check results rather than requiring callers to scrape message text.
- **Do not panic on malformed report input:** Dive's panic-on-unmarshal pattern is explicitly called out as a parse-error anti-pattern (`dive/image/docker/manifest.go:18`). Missing sections, unreadable files, malformed ratings, and ambiguous ratings should become validation failures with causes, not panics.
- **Do not let direct IO or globals leak into validation core:** IO evidence flags hardcoded `os.Stdout`, `os.Stderr`, `os.Open`, and singleton state as testability risks (`cmd/ls/ls.go:42`, `internal/cmd/templatefuncs.go:296`, `fs/fshttp/http.go:38`). Validators should be deterministic and easy to test without terminal, runtime, network, or OpenCode.
- **Do not treat runtime or file existence success as product success:** The PRD principle says runtime success is not product success, and validation evidence supports treating artifact checks as first-class gates. A report file that exists but lacks rating, summary, required answers, or required citation shape must not be considered valid.
- **Do not scan whole source repositories during report validation:** Performance evidence warns against unbounded repository/tree work (`gdu/pkg/analyze/parallel.go:13`, `yq/pkg/yqlib/stream_evaluator.go:78-113`). This sprint should validate citation shape only and leave resolving cited files to the code extraction sprint.
- **Do not leak secrets or unsafe raw content in diagnostics:** Security evidence uses redaction wrappers and credential scrubbing (`internal/options/secret_string.go:15-20`, `pkg/registry/transport.go:37-41`). Diagnostics should include safe report paths and concise observed facts, not full environment values, API tokens, or large raw report dumps.
- **Do not turn inapplicable Markdown source/dimension pairs into failures:** Sprint 5 behavior and this sprint's requirements make inapplicable pairs caller-excluded/skipped. Validation should not reintroduce missing-report failures for pairs that should not have produced reports.

## Examples Worth Inspecting

| Example | Path / Source | Why It Is Useful |
| --- | --- | --- |
| Thin CLI with domain-owned internals | `cmd/root.go:112-127`, `internal/restic/repository.go:18`, `internal/chezmoi/chezmoi.go:1-2` | Shows CLI-to-domain dependency direction and internal package protection for product behavior. |
| `%w` wrapping with typed/sentinel checks | `rclone/vfs/zip.go:21`, `gh-cli/pkg/ssh/ssh_keys.go:64`, `helm/pkg/storage/driver/driver.go:27-48` | Useful when reasoning about preserving filesystem/parsing causes while adding validation context. |
| Actionable user guidance | `age/cmd/age/tui.go:47-54`, `gh-cli/internal/ghcmd/cmd.go:281-301`, `go-task/errors/errors_task.go:13-32` | Shows how diagnostics can include hints without hiding the underlying error. |
| IO test constructor and buffers | `pkg/iostreams/iostreams.go:551-568`, `ui_mock.go:27-33`, `executor_test.go:146-151` | Useful for deterministic tests around diagnostics and output if CLI validation output later enters scope. |
| Filesystem abstraction and test FS | `internal/chezmoi/system.go:25`, `internal/chezmoitest/chezmoitest.go:86-92`, `internal/fs/interface.go:10-31` | Relevant if report validators need a small filesystem seam for fixtures. |
| Golden and fixture testing | `helm/internal/test/test.go:43`, `go-task/task_test.go:166-169`, `dive/cmd/dive/cli/cli_build_test.go:13-27` | Supports fixture-based validation tests for Markdown report shapes and rating formats. |
| Structured logging fields | `internal/slogs/keys.go:6-231`, `internal/logging/logging.go:31-66`, `pkg/yqlib/logger.go:5` | Useful if validation diagnostics are later surfaced through structured logs/status. |
| Security redaction | `internal/options/secret_string.go:15-20`, `pkg/registry/transport.go:37-41`, `gh-cli/status.go:332-338` | Relevant to deciding what validation diagnostics may safely expose. |
| Lazy and bounded processing | `gh-cli/pkg/cmdutil/factory.go:27-42`, `age/internal/stream/stream.go:20`, `yq/pkg/yqlib/stream_evaluator.go:78-113` | Supports keeping validation targeted to expected artifact files and avoiding unrelated scans. |

## Design Pressures

- Validation must be a product gate: report files are durable study artifacts, and runtime exit success must not be enough to mark work complete.
- The sprint folder name mentions execution, but sprint scope explicitly excludes runtime execution, prompt composition, agentwrap/OpenCode wiring, run-loop behavior, summary writing, and code extraction.
- Study ownership is a hard pressure: validation, rating parsing, report paths, and diagnostics belong in `internal/study` unless reasoning identifies a concrete reusable boundary.
- Directory and Markdown document sources create different validation obligations; directory sources need code citation shape by default, while Markdown document sources need rating, summary, and answers without default code citations.
- Rating parsing must preserve uncertainty: missing, invalid, ambiguous, and valid scores are different states, and ambiguous values must not be coerced into `0`.
- Diagnostics must be actionable and deterministic enough for tests while remaining safe for logs and future status/JSON surfaces.
- Validators should be future-usable by single-run, batch, synthesis, summary, and validate commands without depending on runtime execution or public stable CLI output in this sprint.
- Validation should process known report paths, not discover work by scanning all reports or repositories indiscriminately.
- Tests need to cover realistic Markdown report variants but should not require real runtimes, network access, provider credentials, or generated model output.

## Open Questions For Reasoning

- What exact section heading variants should count as a required heading, source information section, summary section, rating section, question/answer content, and final-report synthesis section?
- Should validation checks be represented as a flat list, grouped by artifact, or both, given later status/validate/summary reuse?
- Which diagnostics should be errors versus warnings, especially ambiguous ratings, optional sections, malformed but non-required content, and citation-shape absence when a dimension disables citations?
- How should a dimension explicitly disable directory code citation requirements without introducing a broad new dimension schema in this sprint?
- What rating parser result shape best preserves missing, invalid, ambiguous, and valid states without coupling summary generation into this sprint?
- Should report path helpers accept full domain objects, scalar names/numbers, or both, while keeping path generation deterministic and testable?
- How much report content, if any, should validation diagnostics include as an observed value before it becomes unsafe or noisy?
- Should final report validation require a parseable numeric rating summary for each source, or only the presence of rating summary content in this sprint?
- Where should Markdown source kind be attached in diagnostics so later callers can distinguish citation exemptions from missing evidence?

## Evidence Pointers

- `studies/go-cli-study/reports/final/01-project-structure.md`: Inspect `Pattern 1`, `Pattern 2`, `Pattern 4`, and the evidence index entries for `cmd/root.go:112-127`, `internal/restic/repository.go:18`, and `internal/chezmoi/chezmoi.go:1-2` when reasoning about package ownership and dependency direction.
- `studies/go-cli-study/reports/final/05-error-handling.md`: Inspect `Pattern 1`, `Pattern 3`, `Pattern 4`, `Pattern 7`, and anti-patterns around panic/string errors when designing validation diagnostics and error wrapping.
- `studies/go-cli-study/reports/final/06-io-abstraction.md`: Inspect `System Interface for Filesystem`, `IOStreams with Test Constructor`, and direct `os.*` warnings when deciding how much filesystem abstraction validation tests need.
- `studies/go-cli-study/reports/final/10-logging-observability.md`: Inspect stdout/stderr separation, structured key patterns, and safe diagnostics if validation output is wired into logging/status later.
- `studies/go-cli-study/reports/final/11-testing-strategy.md`: Inspect table-driven, fixture, golden-file, and centralized mock sections to shape validator and rating parser tests.
- `studies/go-cli-study/reports/final/13-security.md`: Inspect path validation, schema validation, secret redaction, and credential scrubbing sections to constrain diagnostic contents and path handling.
- `studies/go-cli-study/reports/final/14-performance.md`: Inspect lazy initialization, streaming, and bounded work sections before adding any report discovery, citation resolution, or repository traversal behavior.
- `studies/go-cli-study/reports/final/15-philosophy.md`: Inspect factory/options, interface abstraction, deliberate complexity, and explicit non-goals when choosing the smallest maintainable validator surface.
- `projects/ultraplan-go/sprints/06-study-run-execution/requirements.md`: Use for sprint scope, acceptance criteria, non-goals, and required outputs, not as external evidence.
- `projects/ultraplan-go/docs/PRD.md`, `projects/ultraplan-go/docs/TRD.md`, `projects/ultraplan-go/docs/ARCHITECTURE.md`: Use for project context, package ownership, and validation requirements, not as selected study evidence.

## Handoff To Reasoning

- Use this handbook as evidence input.
- Validate whether the observed patterns fit this project's constraints.
- Do not copy external patterns without sprint-specific reasoning.
- Keep final implementation decisions in sprint reasoning and plan artifacts, not in this handbook.
