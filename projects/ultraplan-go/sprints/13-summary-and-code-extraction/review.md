# Sprint Review: 13 Summary And Code Extraction

> Project: `ultraplan-go`
> Sprint: `13-summary-and-code-extraction`
> Plan: `projects/ultraplan-go/sprints/13-summary-and-code-extraction/plan.md`
> Reasoning: `projects/ultraplan-go/sprints/13-summary-and-code-extraction/reasoning.md`

## Status

- **Result:** complete
- **Completed:** 2026-06-12
- **Overall Verdict:** accepted; optional follow-ups remain

## Review Inputs

- `sprint-index.md` — 8 selected contracts, 11 evidence reports, 2 reasoning template slots (none created).
- `technical-handbook.md` — distilled evidence into patterns, anti-patterns, design pressures.
- `reasoning.md` — 7 final decisions (D1–D7) and 31 acceptance criteria (AC-01 through AC-31).
- `plan.md` — 10 tasks all marked complete, 4 evidence gates closed.
- `requirements.md` — sprint contract.
- Implementation diff under `/home/antonioborgerees/coding/ultraplan-go`:
  - **Modified:** `internal/app/app.go`, `internal/app/app_test.go`, `internal/app/study_commands.go`, `internal/codeextract/doc.go`, `internal/study/service.go`, `internal/study/summary.go`, `internal/study/summary_test.go`.
  - **New:** `internal/app/code_commands.go`, `internal/app/code_commands_test.go`, `internal/app/study_summary_commands_test.go`, `internal/codeextract/domain.go`, `internal/codeextract/parser.go`, `internal/codeextract/resolver.go`, `internal/codeextract/service.go`, `internal/codeextract/codeextract_test.go`.
- Verification evidence: `go test ./...` passes; `go build ./cmd/ultraplan` succeeds; `go vet ./...` is clean.

## Decision Conformance

| Decision From `reasoning.md` | Implemented? | Evidence | Notes |
| --- | --- | --- | --- |
| D1 — Study-owned deterministic summary matrix (`N/A` for inapplicable, empty for missing/invalid, total-desc sort) | yes | `internal/study/summary.go:29-92`; `summary_test.go:9-43` | CSV header `source,<dim...>,total`; rows sorted total-desc then source-asc; existing `findRating`, `SourceAppliesToDimension`, `SourceReportPath` reused. |
| D2 — Standalone `study <study> summary` with loud atomic writes | yes | `internal/app/study_commands.go:56-57,460-474`; `internal/study/summary.go:94-123`; `summary_test.go:70-89` | `atomicWriteFile` (CreateTemp→Write→Sync→Close→Rename) with `syncDir`; write failures surface as non-zero exit via `mapStudySummaryError` (study_commands.go:476-482). No runtime invocation. |
| D3 — `internal/codeextract` ownership (domain, parser, resolver, service); no global product packages | yes | `internal/codeextract/{domain,parser,resolver,service}.go`; grep confirms root `internal/` contains only `app, codeextract, platform, study, workspace` | `app/code_commands.go` is ~120 lines and only parses/invokes/renders/exits. |
| D4 — Source table parsing, citation forms (single/range/list/legacy dash), resolution order (source-relative, source-prefixed, basename fallback) | yes | `parser.go:20-150`; `resolver.go:62-97`; `codeextract_test.go:10-61` | Legacy `–` (en-dash) normalized to `-` (parser.go:60); unresolved for malformed specs, escapes, out-of-range, ambiguous basenames. |
| D5 — Deterministic text/JSON output; exit mapping (0/2/4/5/8) | yes | `service.go:128-167`; `code_commands.go:70-77`; `code_commands_test.go:53-72` | `StatusValidation`→`ExitValidation(5)`, `StatusPartial`→`ExitPartial(8)`; `errorCode()` produces stable category codes (app.go:52-71). |
| D6 — Local-only path handling; ignored dirs (`.git .hg .svn node_modules .ultraplan`); per-invocation basename cache | yes | `resolver.go:11-17,62-97,142-167`; `codeextract_test.go:14-16` | NUL byte check, absolute rejection, `..` escape rejection, symlink-aware `containedPath`. Cache scoped to resolver lifetime. |
| D7 — Verification gates (`go test ./...`, `go build ./cmd/ultraplan`); boundary review | yes | See *Verification Evidence* below | All gates green; no platform/product coupling. |

## Contract Conformance

### Architecture — pass

| Requirement | Status | Evidence |
| --- | --- | --- |
| ARCH-CORE-001 (module boundaries) | satisfied | `internal/study` owns summary (`summary.go`, `service.go::WriteSummary`); `internal/codeextract` owns extraction (`domain.go`, `parser.go`, `resolver.go`, `service.go`); `internal/app` is thin adapter. |
| ARCH-CORE-002 (dependency direction) | satisfied | `grep -r "ultraplan-go/internal/study\|ultraplan-go/internal/codeextract\|ultraplan-go/internal/app" internal/platform/` returns **zero matches**; `internal/platform/runtime` only imports `platform/config`. |
| ARCH-LAYER-001 (domain purity) | satisfied | `codeextract` domain types and parser are stdlib-only; no transport/SDK imports. |
| ARCH-LAYER-002 (use cases depend on ports) | satisfied | `study.Service` accepts `Runtime` port and `WithRuntime` option; no concrete provider SDK in use cases. |
| ARCH-ENTRY-001 (thin entrypoints) | satisfied | `runCode` (code_commands.go:19-78) ~60 lines, `runStudySummary` (study_commands.go:460-474) ~15 lines; no product logic in `internal/app`. |
| ARCH-MODULE-001 (small public API) | satisfied | `codeextract` exposes 8 domain types and 3 public functions (`Extract`, `RenderText`, `RenderJSON`); no exported repos/SQL/providers. |
| ARCH-COMP-001 (composition through registrars) | satisfied | `study.NewService(workspaceRoot, opts...)` with `Option` pattern; CLI uses the service, not raw internals. |
| ARCH-SHARED-001 (platform domain-neutral) | satisfied | `internal/platform/{config,filesystem,logging,runtime}` expose only generic services; no product business behavior. |

### Errors — pass

| Requirement | Status | Evidence |
| --- | --- | --- |
| ERR-CORE-001 (no failure disappears) | satisfied | `summary.go:88-90` wraps `atomicWriteFile` with `%w`; `code_commands.go:70-77` maps `StatusValidation`/`StatusPartial` to non-zero exits; `codeextract_test.go:44-61` exercises unresolved/out-of-range/escape. |
| ERR-CODE-001 (stable machine-usable codes) | satisfied | `app.go:52-71` defines `errorCode()` returning `validation.usage`, `validation.workspace`, `validation.reference`, `workflow.partial`, `config.invalid`, `provider.runtime`, `workflow.cancelled`, `internal.error`. |
| ERR-TRANS-001 (preserve cause/context) | satisfied | All classified errors use `fmt.Errorf("...: %w", err)`; `classedError.Unwrap()` returns `e.err` (app.go:45). |
| ERR-IO-001 (boundary preserves signal) | satisfied | `classedError{class, code, err}` keeps class+code at boundary; CLI output preserves user-facing error string. |
| ERR-DATA-001 (persistence loud/distinct) | satisfied | Summary write failure returns wrapped error (summary.go:88-90); `summary_test.go:70-89` verifies prior content is preserved. |
| ERR-USER-001 (user vs operator separation) | satisfied | Summary path → stdout; warnings → stderr (study_commands.go:465-472). Code text → stdout or `--output` file; errors → stderr via `fail()`. |

### Security — pass with minor follow-ups

| Requirement | Status | Evidence |
| --- | --- | --- |
| SEC-INPUT-001 (untrusted input validated) | satisfied | Citation parsing rejects malformed specs (parser.go:39-46); source table rows require non-empty name+path (parser.go:124-127). |
| SEC-INJECT-001 (safe APIs) | satisfied | `filepath.Clean`, `filepath.FromSlash`, `containedPath()` with `filepath.EvalSymlinks`; no shell-exec; no SQL. |
| SEC-SECRETS-001 (no secret leakage) | satisfied | `grep -n "opencode\|agentwrap\|http\.\|net\.\|os\.Getenv" internal/codeextract/*.go` returns no matches. Renderers only emit report/source/path/lines/snippet fields. |
| SEC-FILES-001 (path access constrained) | satisfied | NUL check (resolver.go:63), absolute rejection (66), `..` escape (70), symlink-aware containment (111-132), explicit output path 0o644 (code_commands.go:64), atomic write to explicit path (summary.go:88-123). |
| SEC-NET-001 (network bounded) | satisfied | No network calls; `internal/codeextract` only uses `os.ReadFile/Stat/MkdirAll/WriteFile`. |
| SEC-DESER-001 (deserialization restricted) | satisfied | Markdown table parsing uses bounded regex and line iteration; no `unsafe` or runtime object construction. |
| SEC-DEFAULT-001 (fail closed) | satisfied | Path escape, NUL byte, absolute path all return `*Diagnostic` not silent success. |

### Testing — pass with minor follow-ups

| Requirement | Status | Evidence |
| --- | --- | --- |
| TEST-SEAM-001 (collaborators replaceable) | satisfied | `codeextract.Extract(Request)`; `study.Service.WriteSummary(studyRef)`; both injectable via constructor options. |
| TEST-UNIT-001 (unit coverage) | satisfied | `summary_test.go` (3 cases); `codeextract_test.go` (2 cases); total 8 sprint-specific test functions. |
| TEST-INT-001 (wiring integration) | satisfied | `code_commands_test.go` and `study_summary_commands_test.go` exercise end-to-end CLI paths via `runForTest`. |
| TEST-FAIL-001 (failure paths tested) | satisfied | Summary write failure (chmod 0o500), code partial (`ExitPartial`), code validation (`ExitValidation`), path escape, out-of-range, malformed specs, ambiguous basename, ignored dirs. |
| TEST-DET-001 (deterministic) | satisfied | All tests use `t.TempDir()`; exact byte-equality assertions; no goroutine timing. |
| TEST-CONTRACT-001 (public contract compatibility) | satisfied | Exit codes and exact output bytes asserted; JSON `status` field asserted. |
| AC-30 (offline) | satisfied | No OpenCode, network, or provider mocks required; all fixtures local. |
| AC-31 (build) | satisfied | `go build ./cmd/ultraplan` succeeds. |

### Documentation — pass

| Requirement | Status | Evidence |
| --- | --- | --- |
| DOC-OWNER-001 (clear ownership) | satisfied | Sprint artifacts under `projects/ultraplan-go/sprints/13-summary-and-code-extraction/`. |
| DOC-ARCH-001 (decision context) | satisfied | `reasoning.md` provides 7 final decisions with trade-offs and rejected alternatives. |
| DOC-PUBLIC-001 (public surfaces documented) | satisfied | `codeextract/doc.go`, `study/doc.go`; `renderHelp()` exposes `code`; `studyHelp()` exposes `<study> summary`; `codeHelp()` shows usage/flags. |
| DOC-EXAMPLE-001 (examples maintained) | satisfied | Test fixtures and help text are kept in sync with command behavior. |
| DOC-GEN-001 (reproducible) | satisfied | `summary.csv` is reproducibly generated; sprint artifacts are hand-written but versioned in repo. |
| DOC-AGENT-001 (concise, current) | satisfied | Package docs state intent and ownership; reasoning and plan are the decision source. |

### CLI Surface — pass

| Requirement | Status | Evidence |
| --- | --- | --- |
| CLI-SHAPE-001 (predictable commands) | satisfied | `ultraplan code <report>...`; `ultraplan study <study> summary`. |
| CLI-HELP-001 (help/version output) | satisfied | Top-level `renderHelp()` (app.go:200-218); `codeHelp()` (code_commands.go:110-120); `studyHelp()` (study_commands.go:122-150). |
| CLI-IO-001 (stdout/stderr separated) | satisfied | Summary path → stdout; warnings → stderr; code text → stdout or `--output`; errors → stderr. |
| CLI-EXIT-001 (stable exit codes) | satisfied | Constants in app.go:16-26: 0/2/4/5/8; mapping in code_commands.go:70-77 and study_commands.go:476-482. |
| CLI-JSON-001 (machine-readable output) | satisfied | `code --json` produces stable structured fields (`status`, `reports`, `references`, `diagnostics`). |
| AC-01, AC-12, AC-20, AC-22, AC-24 | satisfied | Help exposure, multi-report order, JSON, output-file write, and exit mapping all in place. |

### Performance — pass

| Requirement | Status | Evidence |
| --- | --- | --- |
| PERF-BOUND-001 (expensive work bounded) | satisfied | Basename fallback scoped to `r.sources` roots; ignored dirs skipped via `filepath.SkipDir` (resolver.go:152-156). |
| PERF-CACHE-001 (cache correctness) | satisfied | `basenameCache` is per-resolver (per-invocation) `map[string][]string`; no persistence, no stale risk. |
| PERF-CONC-001 (concurrency safe) | not_triggered | No concurrency added; single-goroutine flows only. |
| PERF-COST-001 (cost drivers) | not_triggered | No runtime/provider paths exercised. |

### Persistence And Migrations — pass

| Requirement | Status | Evidence |
| --- | --- | --- |
| PERSIST-ATOMIC-001 (atomic write boundaries) | satisfied | `atomicWriteFile` in summary.go:94-123 uses CreateTemp+Write+Sync+Close+Rename; `syncDir` after rename (line 121). |
| PERSIST-DERIVED-001 (derived stores) | satisfied | `summary.csv` is a derived artifact from per-dimension reports; `summary.csv` regenerable deterministically. |
| PERSIST-SCHEMA-001 (schema ownership) | satisfied | `summary.csv` is owned by `internal/study` via `SummaryPath` (summary.go:25-27). |
| PERSIST-FIXTURE-001 (fixtures isolated) | satisfied | Tests use `t.TempDir()`; no production fixture leakage. |
| PERSIST-RECOVERY-001 (recovery) | satisfied | Existing `summary.csv` preserved on write failure (summary_test.go:70-89). |
| AC-08, AC-09, AC-22, AC-31 | satisfied | Atomic summary, loud failures, explicit `--output` write, build all green. |

## Plan Execution

| Plan Task | Status | Evidence |
| --- | --- | --- |
| T1: Inventory existing study/app surfaces | done | `service.go` adds `WriteSummary`; CLI reuses `study.NewService`. |
| T2: Study summary domain and service | done | `summary.go` + `service.go::WriteSummary`. |
| T3: Test summary semantics and writes | done | `summary_test.go` (3 cases). |
| T4: Wire standalone summary CLI | done | `study_commands.go:56-57,460-474`; `study_summary_commands_test.go`. |
| T5: Code extraction domain, parser, resolver | done | `codeextract/{domain,parser,resolver}.go`. |
| T6: Code extraction service and output | done | `codeextract/service.go` (Extract, RenderText, RenderJSON). |
| T7: Wire top-level code CLI | done | `app.go:128-129,209`; `app/code_commands.go`; `code_commands_test.go`. |
| T8: Test code extraction thoroughly | done | `codeextract_test.go` (2 cases); `code_commands_test.go` (2 cases). |
| T9: Architecture/boundary/non-goal review | done | No platform/product coupling; no global product packages; no runtime calls. |
| T10: Final verification and review handoff | done | `go test ./...` passes; `go build ./cmd/ultraplan` succeeds; this review. |

## Verification Evidence

| Check | Command | Result | Notes |
| --- | --- | --- | --- |
| Full Go test suite | `go test ./...` | pass | `internal/app`, `internal/codeextract`, `internal/platform/*`, `internal/study`, `internal/workspace` all pass; `cmd/ultraplan` and `internal/platform/filesystem` have no test files. |
| CLI build | `go build ./cmd/ultraplan` | pass | Binary `ultraplan` produced. |
| `go vet` | `go vet ./...` | pass | No warnings. |
| Architecture import review | `grep -r "ultraplan-go/internal/{study,codeextract,app}" internal/platform/` | pass | Zero matches — platform remains product-free. |
| Non-runtime review | `grep -n "opencode\|agentwrap\|http\.\|net\.\|os\.Getenv" internal/codeextract/*.go` | pass | Zero matches — no runtime/network/env access. |
| Public packages list | `ls internal/` | pass | Only `app, codeextract, platform, study, workspace` — no global `summary`, `reports`, `parser`, `resolver`, or `validation` package. |
| Top-level help | `renderHelp()` in `app.go:200-218` | pass | Lists `init-workspace`, `config`, `code`, `health`, `study`, `version`. |
| Study help | `studyHelp()` in `study_commands.go:122-150` | pass | Lists `<study> summary` at line 129, 142. |

## Deviations

| Deviation | Reason | Risk | Follow-Up |
| --- | --- | --- | --- |
| Stale protocol path `system/protocols/sprint-review-protocol.md` was referenced in `sprint-index.md` and `plan.md`; actual file is `review-sprint-protocol.md`. | Documentation was corrected after review. | Low — only documentation drift, not a behavior issue. | Done: `sprint-index.md` and `plan.md` now reference `review-sprint-protocol.md`. |
| No `reasoning/*.md` area-specific reasoning files were created for this sprint, even though `sprint-index.md` and `requirements.md` list them. | Reasoning.md (lines 36-42) explicitly records this absence and uses consolidated reasoning as the decision source. Plan treats it as closed. | Low — decision quality is preserved in `reasoning.md`. | If a future sprint reuses the same artifact list, either create the directory or update requirements/index to drop the references. |
| `codeextract.Extract` returns `Result{}, err` and aborts on a single unreadable report (service.go:14-18) rather than collecting a per-report diagnostic and continuing. | Reasoning.md:64 covers unresolved references, not filesystem failures. ExitWorkspace (4) is the prescribed handling for unreadable inputs (reasoning D5). | Low — reasoning categorizes unreadable as workspace failure, not partial; behavior matches existing app exit conventions. | If multi-report with mixed readability is desired, change `extractReport` to surface an `unreadable report` diagnostic and continue. Add a fixture test covering this case. |
| No dedicated `code --help` test in `code_commands_test.go`. | Top-level help test already covered the registered command. | Low — coverage existed indirectly. | Done: `TestCodeCommandHelp` added. |
| No explicit multi-report argument-order test in `code_commands_test.go`. | Argument order was deterministic by implementation review. | Low — flow was deterministic and reviewable. | Done: `TestCodeCommandProcessesReportsInArgumentOrder` added. |
| No `--output` write-failure test in `code_commands_test.go`. | Error path was simple but uncovered. | Low — one-line `os.WriteFile` mapped to `ExitWorkspace`. | Done: `TestCodeCommandOutputWriteFailure` added. |

## Applicability Table

| Contract Requirement or Pattern | Applicability | Scope Reason |
| --- | --- | --- |
| Configuration contract | explicitly_deferred | Sprint scope does not change config loading, precedence, health mapping, or secret config display. |
| LLM Runtime contract | explicitly_deferred | Sprint must not launch, bypass, extend, or recompose runtime execution. |
| LLM Evaluation / Cost / Safety | explicitly_deferred | No provider execution, cost metadata, or model policy in scope. |
| Workflows contract | explicitly_deferred | Run-loop orchestration is non-goal; compatibility with completed outputs only. |
| `07-state-context` report | explicitly_deferred | No direct run-state mutation; inspection of completed artifacts only. |
| `08-concurrency` report | explicitly_deferred | Deterministic local processing; no new worker pools or parallel scheduling. |
| `12-extensibility` report | explicitly_deferred | No plugin architecture; module boundaries covered by Architecture selections. |
| Target/sprint workflow | explicitly_deferred | Project and sprint non-goals unrelated to deterministic summary/extraction. |
| Hosted/browser/TUI/plugin/workflow | explicitly_deferred | Out of scope; not touched by this sprint. |
| Remote source resolution | explicitly_deferred | Code extraction local-only under parsed source roots. |
| Multi-report partial continuation on filesystem failure | partial | Continues for unresolved references; aborts on unreadable inputs as per ExitWorkspace convention. |
| `code --help` direct test | satisfied | Covered directly by `TestCodeCommandHelp`. |
| Persistent basename cache across invocations | partial | Cache is per-invocation only by design; persistence is explicitly out of scope. |
| Stable public JSON schema for `code` | partial | Sprint-local deterministic JSON only; release-wide stability deferred to Sprint 14 per reasoning. |

## Architecture Review Protocol Notes

- **Behavior review (1):** Main workflows visible. `WriteSummary` flow in summary.go:29-92; `Extract` flow in service.go:12-47; `runCode` and `runStudySummary` are visibly thin.
- **Architecture fit (2):** Good fit. Module ownership matches reasoning: `internal/study` (summary), `internal/codeextract` (extraction), `internal/app` (CLI), `internal/platform` (generic).
- **Simplicity (3):** Same complexity. No speculative abstractions; no new framework. Atomic write helper added in summary.go is local and justified.
- **Cohesion (4):** Strong. `codeextract` is a focused bounded context with parser/resolver/service split. `study` summary is a self-contained file plus service method.
- **Coupling (5):** Acceptable. CLI depends on `codeextract` public API and `study.Service`; neither product module depends on `internal/app`. No globals.
- **DRY (6):** No concerning duplication. `summary.go` reuses `findRating`, `SourceAppliesToDimension`, `SourceReportPath`. `codeextract` owns its own parser.
- **State/side effects (7):** Clear and safe. Atomic summary write; explicit `--output` path; explicit per-invocation basename cache.
- **Paradigm (8):** Good fit for Go. Functions, methods on small types, no inheritance.
- **Error handling (9):** Strong. `classedError` with codes, `%w` wrapping, exit-class mapping, warnings separated to stderr.
- **Observability (10):** Acceptable. Warnings surface source/dimension/path/reason; no secret leakage.
- **Testing (11):** Strong. Exact-byte CSV, exit-code coverage, write-failure preservation, resolution edge cases.
- **Performance (12):** Fine. Bounded scans, ignored dirs, per-invocation cache.
- **Maintainability (13):** Easier. New module is self-contained; summary flow is hardened and reuses existing helpers.

## Review Protocol Results

| Protocol | Applied? | Findings |
| --- | --- | --- |
| Architecture Review | yes | All checks pass; boundary between product and platform preserved. |
| Sprint Review (review-sprint-protocol.md) | yes | All 31 acceptance criteria evidenced; 7 reasoning decisions implemented; verification gates green. |

## Final Assessment

- **Status:** accepted; optional follow-ups remain
- **Overall Result:** The sprint achieves its goal of deterministic study summary generation and code-reference extraction. Summary behavior lives in `internal/study`; extraction in `internal/codeextract`; CLI adapters are thin; platform remains product-free. All verification gates pass.

### Residual Risks

- Sprint-local JSON schema for `code --json` is not a release-wide contract (explicitly deferred to Sprint 14).
- Source table parser supports the required row shape with backticked paths; broader hand-written report variations remain future work.
- Snippets contain only requested lines under parsed source roots; if reports cite sensitive source lines, those snippets will be rendered to stdout or `--output`.
- `codeextract.Extract` aborts on any unreadable report rather than continuing with remaining reports (preserved as partial because reasoning treats unreadable as workspace failure, not partial).

### Required Follow-Ups

None — no blockers.

### Optional Improvements Later

- Decide whether unreadable reports should continue to fail fast with `ExitWorkspace` or become per-report diagnostics in a future sprint.
- Consider aligning perms for `--output` writes with the `atomicWriteFile` helper pattern if future write-atomicity requirements appear for extraction output.
- Add a "See also" pointer between `internal/codeextract/doc.go` and `internal/study/doc.go` to make the deliberate split visible.
- If summary columns ever need to grow, add a versioning line in the CSV header (analogous to `RunState.SchemaVersion`) to make PERSIST-MIG-001 trivially satisfied in future sprints.
