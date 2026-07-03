# Sprint Review: 07 Prompt Composition Generation

> Project: `ultraplan-go`
> Sprint: `07-prompt-composition-generation`
> Date: 2026-05-31
> Status: accepted

## Summary

Implemented deterministic study prompt composition and preview rendering for directory analysis, Markdown document analysis, and synthesis prompts. All six selected contracts (Architecture, Errors, Security, Testing, Documentation, CLI Surface) and the technical handbook guidance are satisfied. No runtime, agentwrap, OpenCode, provider, network, or subprocess invocation paths were introduced.

Prompt behavior is owned by `internal/study`. CLI parsing and preview rendering are owned by `internal/app`. Platform packages do not import `study`. No global `internal/prompts`, `internal/templates`, or `internal/reports` packages were created.

## Review Inputs

- `sprint-index.md`
- `technical-handbook.md`
- `reasoning.md`
- `plan.md`
- Required protocols from `sprint-index.md`: Architecture Review, Sprint Review
- Implementation diff: `internal/study/domain.go`, `internal/study/prompts.go`, `internal/study/prompts_test.go`, `internal/app/study_commands.go`, `internal/app/study_prompt_commands_test.go`
- Verification: `go test ./...` passed, `go build ./cmd/ultraplan` passed

## Decision Conformance

| Decision From `reasoning.md` | Implemented? | Evidence | Notes |
| --- | --- | --- | --- |
| Study owns prompt composition | yes | `internal/study/domain.go` (types), `internal/study/prompts.go` (builders), `internal/app/study_commands.go` (thin CLI) | No global prompt/template/report package; platform packages do not import study |
| Deterministic manifest-backed prompt results | yes | `domain.go:61-84` PromptManifest with JSON tags, workspace-relative paths, deterministic ordering | Manifests include kind, study, dimension, source, template paths, input paths, output path |
| Directory analysis prompts preserve repository isolation and code citations | yes | `prompts.go:109-119` source isolation rules, `prompts.go:26-28` DisableCodeCitations check | File-path-line citation requirements unless dimension disables |
| Markdown document prompts are embedded-document only | yes | `prompts.go:127-157` embedded doc with frontmatter stripping, `stripFrontmatter` helper | Forbids external file/repository/code inspection; defaults away from code citations |
| Synthesis uses applicable report manifests only | yes | `prompts.go:59-64` deterministic sorting, applicability filtering, missing-report failure | Inapplicable Markdown pairs excluded from synthesis |
| Preview command is runtime-free rendering only | yes | `study_commands.go:150,164` calls builders directly; no runtime/agentwrap/OpenCode paths | Help text at `study_commands.go:275` confirms no-runtime |
| Tests focus on determinism, source-kind semantics, errors, and no-runtime boundaries | yes | `prompts_test.go`, `study_prompt_commands_test.go` | All test cases use local fixtures; no network/provider dependencies |

## Contract Conformance

### Architecture

| Requirement | Status | Evidence |
| --- | --- | --- |
| Module boundaries must remain explicit | satisfied | Prompt types/behaviors in `internal/study`, CLI in `internal/app` |
| Dependency direction must point inward | satisfied | `study` imports `workspace`; `app` imports `study`/`workspace`; platform packages do not import `study` |
| Domain layer must remain pure | satisfied | `domain.go:46-84` types are transport/persistence-free; builders use `os.ReadFile` and workspace helpers |
| Use cases depend on ports when a real seam exists | satisfied | `BuildAnalysisPrompt`/`BuildSynthesisPrompt` take `PromptRequest`, return `PromptResult`; no agentwrap/runtime/network |
| Entrypoint adapters must stay thin | satisfied | `runStudyPrompt` parses args, resolves refs, calls builders, renders output |
| Module public APIs must stay small and stable | satisfied | 5 exported types, 3 PromptKind constants, 3 exported functions |
| Composition must happen through registrars | satisfied | `NewService` constructor, no god-object |
| Shared/platform code must stay domain-neutral | satisfied | `workspace/paths.go` provides generic path utilities |

### Errors

| Requirement | Status | Evidence |
| --- | --- | --- |
| No failure may disappear | satisfied | `ErrPromptInapplicable` sentinel, wrapped file read errors, classified CLI errors |
| Error codes must be stable and machine-usable | satisfied | `errorCode` maps classes to stable codes: `validation.usage`, `validation.reference`, `validation.workspace` |
| Error translation must preserve cause and context | satisfied | `%w` wrapping; `errors.As`/`errors.Is` in `mapStudyPromptError` |
| Known-fatal prerequisites must fail before work | satisfied | Missing templates/docs/reports fail before prompt rendering |
| Boundary adapters must preserve operational signal | satisfied | `classifiedError` with code/class/cause; stderr diagnostics |
| User-facing and operator-facing messages separated | satisfied | Prompt output to stdout; diagnostic errors to stderr |
| Canonical error shape | not_applicable | Deferred to observability/runtime execution sprint |
| Retryable errors | not_applicable | Deferred to runtime execution sprint |
| Owned task failure state | not_applicable | Deferred to workflow sprint |
| External dependency failures | not_applicable | Sprint has no external dependencies |
| Persistence failures | not_applicable | Sprint has no persistence |
| Secret redaction | not_applicable | Sprint has no secret material |

### Security

| Requirement | Status | Evidence |
| --- | --- | --- |
| Untrusted input must be validated at boundaries | satisfied | Study/dimension/source refs validated via service resolution before use |
| Injection-prone operations must use safe APIs | satisfied | `strings.Builder`, `fmt.Fprintf`, `os.ReadFile`; no `exec.Command` or eval |
| File and path access must be constrained | satisfied | `workspace.ResolveInside` validates paths stay within workspace root |
| No-runtime preview boundary | satisfied | No agentwrap, OpenCode, provider, network, or subprocess calls in preview path |
| Authentication/Authorization/Secrets/Deps/Deser/Net | not_applicable | Not triggered by this sprint's prompt-only scope |

**Minor observation**: `writePromptPreview` at `study_commands.go:250` uses `0o644` file permissions. Consider `0o600` for preview output files.

### Testing

| Requirement | Status | Evidence |
| --- | --- | --- |
| Collaborators must be replaceable through public seams | satisfied | `PromptRequest` struct with explicit fields; tests use fixtures |
| Business logic must have unit coverage | satisfied | Directory, Markdown, synthesis prompt unit tests |
| Failure paths must be tested explicitly | satisfied | Missing templates, documents, reports, inapplicable pairs, unknown refs all tested |
| Tests must remain deterministic | satisfied | `t.TempDir()` fixtures, `sort.SliceStable` ordering, no network/timing deps |
| Public API/CLI compatibility must have tests | satisfied | Command tests for stdout/output-file success, exit codes, help text |
| Wiring/persistence integration coverage | not_applicable | No persistence changes in sprint |
| Real-dependency smoke coverage | not_applicable | Deferred to runtime execution sprint per non-goals |
| End-to-end scenario tests | not_applicable | Not triggered by prompt-only scope |

### Documentation

| Requirement | Status | Evidence |
| --- | --- | --- |
| Important docs must have clear ownership and location | satisfied | Sprint artifacts at standard `projects/ultraplan-go/sprints/07-prompt-composition-generation/` |
| Public surfaces must be documented | satisfied | `study_commands.go:265-277` help text describes prompt-only preview correctly; no-runtime boundary stated |
| Architecture changes need decision context | not_applicable | No architecture changes |
| Operational/Example/Generated/Agent docs | not_applicable | Not triggered by sprint scope |

### CLI Surface

| Requirement | Status | Evidence |
| --- | --- | --- |
| Commands and flags must be predictable | satisfied | `ultraplan study <study> prompt analysis|synthesis <dimension> [source] [--output]` |
| Help/version output must exist | satisfied | `studyPromptHelp`, `studyHelp`, version command from prior sprints |
| Stdout/stderr must be separated | satisfied | Prompt output to stdout; errors to stderr |
| Exit codes must be stable | satisfied | `ExitValidation`, `ExitWorkspace` codes for different error classes |
| Machine-readable output must be stable | satisfied | JSON manifest with stable struct fields |
| Non-interactive mode must never hang | satisfied | No interactive prompts; all ops are deterministic file reads |
| Config precedence | not_applicable | Deferred per sprint-index excluded context |
| CLI lifecycle/cancellation | not_applicable | Single short-lived file read+render; no cancellation needed |

## Technical Handbook Conformance

| Pattern | Status | Evidence |
| --- | --- | --- |
| Study-owned prompt logic behind thin CLI wiring | followed | Builders in `prompts.go`, thin handlers in `study_commands.go:150,164` |
| Factory-style command construction | followed | Explicit `deps`, `root`, `service` parameters; no globals |
| Lazy, deterministic input loading | followed | On-demand `readWorkspaceFile`; no runtime/provider/network |
| Injectable IO for preview output and tests | followed | `deps.stdout` `io.Writer`; capturable in tests |
| Structured, actionable error paths | followed | `%w` wrapping, `ErrPromptInapplicable` sentinel, named paths |
| Golden/fixture tests for deterministic artifacts | followed | Tempdir fixtures, repeated-build determinism test |
| Explicit trust boundaries for no-runtime preview | followed | Code path calls only builders; help states no-runtime |

| Anti-Pattern | Status | Evidence |
| --- | --- | --- |
| Command handlers as prompt engines | avoided | Handlers parse args, call builders, render output only |
| Hidden globals for templates/config/IO | avoided | No package-level state |
| Bypass injectable output streams | avoided | `deps.stdout` interface; no direct `os.Stdout` |
| Break error chains or string matching | avoided | `%w` wrapping; `errors.Is` for sentinel checks |
| Treat Markdown docs like repositories | avoided | Document-only rules, frontmatter stripping, no-code-citation default |
| Speculative performance machinery | avoided | No concurrency, pools, or caching |
| Brittle implementation detail tests | avoided | Manifest fields and required content checks; no private field inspection |

## Plan Execution

| Plan Task | Status | Evidence |
| --- | --- | --- |
| Task 1: Inspect existing study/workspace APIs | done | Existing APIs reused; `workspace.ResolveInside`, frontmatter stripping, applicability helpers |
| Task 2: Define prompt domain types | done | `domain.go:46-84` adds `PromptRequest`, `PromptResult`, `PromptManifest`, `SourceReportInput`, `PromptKind` |
| Task 3: Implement template loading and path rendering | done | `prompts.go:20-53` loads templates; workspace-relative paths in manifests |
| Task 4: Implement directory analysis prompt builder | done | `prompts.go:102-119` with source isolation, citation rules, disabled-citation support |
| Task 5: Implement Markdown document prompt builder | done | `prompts.go:127-157` with frontmatter stripping, document-only rules, inapplicable handling |
| Task 6: Implement synthesis prompt builder | done | `prompts.go:51-95` with applicability filtering, deterministic sorting, missing-report failure |
| Task 7: Add prompt builder unit tests | done | `prompts_test.go` covers all builder types, manifests, missing inputs, determinism |
| Task 8: Add runtime-free preview CLI wiring | done | `study_commands.go:124-172` with analysis/synthesis subcommands and `--output` flag |
| Task 9: Add preview command tests | done | `study_prompt_commands_test.go` covers stdout, output file, all failure cases, no-runtime |
| Task 10: Verify scope, architecture, and build | done | `go test ./...` passed; `go build ./cmd/ultraplan` passed; architecture and scope review clean |

## Verification Evidence

| Check | Command / Method | Result | Notes |
| --- | --- | --- | --- |
| Unit tests | `go test ./...` from ultraplan-go | pass | All tests pass without network/provider/runtime dependencies |
| CLI build | `go build ./cmd/ultraplan` from ultraplan-go | pass | Binary builds successfully; removed after verification |
| Architecture review | Import/package inspection | pass | Study-owned prompt logic; app-owned CLI wiring; platform packages do not import study |
| Scope review | Diff and behavior inspection | pass | No runtime execution, run-loop, summary, code extraction, target/sprint workflow added |

## Deviations

| Deviation | Reason | Risk | Follow-Up |
| --- | --- | --- | --- |
| Preview command implementation extended `internal/app/study_commands.go` instead of creating `internal/app/study_prompt_commands.go` | Matches existing command organization, avoids unnecessary file splitting, behavior and package ownership remain aligned with the plan | Low - no behavioral or ownership difference | None |

## Residual Risks

| Risk | Mitigation | Status |
| --- | --- | --- |
| Citation-disable detection is prose-based until dimensions get structured metadata | Centralized in one helper (`DisableCodeCitations` bool check); tested | carried forward |
| Preview CLI shape may conflict with future runtime dry-run UX | Preview is intentionally minimal and backed by reusable study builders | carried forward |
| Missing synthesis reports block prompt generation strictly | Errors identify missing reports and explain requirements; partial mode deferred | accepted by design |
| Manifest paths could leak absolute local details | Workspace-relative paths used for manifests; reserved absolutes for diagnostics | closed |

## Review Protocol Results

| Protocol | Applied? | Findings |
| --- | --- | --- |
| Architecture Review | yes | Pass: prompt behavior study-owned, CLI wiring app-owned, no global prompt/template packages, platform packages do not import study, dependency direction correct |
| Sprint Review | yes | Pass: all required outputs exist, acceptance criteria met, no non-goal scope added, `go test ./...` and `go build ./cmd/ultraplan` pass |

## Final Assessment

- **Status:** accepted
- **Residual Risks:** Citation-disable detection is prose-based; partial synthesis preview remains deferred by design
- **Required Follow-Ups:** None
