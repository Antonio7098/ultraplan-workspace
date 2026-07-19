# Sprint Requirements: Automated Sprint Review

> Project: `ultraplan-go`
> Sprint: `26-review-stage`
> Purpose: the authoritative, human-readable sprint contract. All other sprint artifacts must satisfy these requirements.

## Sprint Goal

Replace manual sprint review with a product-owned, evidence-grounded review stage that writes the current sprint-root `review.md` and is fully operable from both CLI and TUI through shared typed use cases.

## Required Outputs

| Output | Path | Description |
| ------ | ---- | ----------- |
| Requirements | `projects/ultraplan-go/sprints/26-review-stage/requirements.md` | This sprint contract defining scope, acceptance criteria, constraints, dependencies, and review expectations. |
| Sprint index | `projects/ultraplan-go/sprints/26-review-stage/sprint-index.md` | Selected contracts, evidence reports, reasoning templates, review protocols, source documents, and explicitly excluded context for the review-stage sprint. |
| Technical handbook | `projects/ultraplan-go/sprints/26-review-stage/technical-handbook.md` | Distillation of only the selected evidence needed to implement automated sprint review. |
| Area reasoning | `projects/ultraplan-go/sprints/26-review-stage/reasoning/architecture.md` | Architecture reasoning for Phase 3 review ownership, dependency direction, runtime boundaries, and CLI/TUI sharing. |
| Final reasoning | `projects/ultraplan-go/sprints/26-review-stage/reasoning.md` | Final implementation decisions, trade-offs, assumptions, risks, and expected evidence for the review stage. |
| Implementation plan | `projects/ultraplan-go/sprints/26-review-stage/plan.md` | Executable task plan that traces to `reasoning.md` and defines verification evidence. |
| Execute summary | `projects/ultraplan-go/sprints/26-review-stage/execute.md` | Sprint execution summary with task outcomes, diagnostics, and evidence pointers after implementation. |
| Flow state | `projects/ultraplan-go/sprints/26-review-stage/flow-state.json` | Versioned sprint-stage state extended through `review`, including status, verdict, artifact path, fingerprint, and diagnostics. |
| Execute run state | `projects/ultraplan-go/sprints/26-review-stage/.run-state.json` | Durable execute task state for the sprint implementation run when execute is used. |
| Review artifact | `projects/ultraplan-go/sprints/26-review-stage/review.md` | Current generated automated conformance review for this sprint, replacing any manual review when `review` runs. |
| Sprint review domain | `../ultraplan-go/internal/sprint/review.go` | Sprint-owned review scope, reviewer orchestration, deterministic checks, verdict synthesis, and atomic `review.md` writing. |
| Review validation | `../ultraplan-go/internal/sprint/review_validation.go` | Validator for generated `review.md`, reviewer coverage, contained citations, verdict values, and stale/missing evidence. |
| Sprint domain and flow updates | `../ultraplan-go/internal/sprint/domain.go` | Phase 3 review stage domain types, stage constants, statuses, verdicts, fingerprints, and flow-state fields. |
| Sprint service and flow updates | `../ultraplan-go/internal/sprint/service.go` | Shared sprint service use cases for review readiness, validation, prompt preview, execution, cancellation, and status integration. |
| Sprint prompt updates | `../ultraplan-go/internal/sprint/prompts.go` | Deterministic review prompt rendering for selected contracts and technical-handbook review requests. |
| Sprint persistence updates | `../ultraplan-go/internal/sprint/store_fs.go` | Atomic persistence and strict loading support required by review artifacts and flow-state updates. |
| CLI review wiring | `../ultraplan-go/internal/app/sprint_commands.go` | CLI support for `review`, `validate review`, `prompt review`, `flow --to review`, text output, JSON output, and exit-code mapping. |
| App review use cases | `../ultraplan-go/internal/app/usecases.go` | Typed app-level review operations shared by CLI and TUI, including progress, cancellation, status, and result DTOs. |
| TUI review support | `../ultraplan-go/internal/tui` | TUI readiness/status, dry-run, prompt preview, confirmation, live reviewer progress, cancellation, finding navigation, verdict display, and `review.md` preview. |
| Embedded review prompt default | `../ultraplan-go/internal/workspace` | Embedded default prompt asset for structured contract and handbook review, installable through `ultraplan defaults install`. |
| Embedded `review.md` template default | `../ultraplan-go/internal/workspace` | Embedded default output template for generated sprint review Markdown, installable through `ultraplan defaults install`. |
| Sprint review tests | `../ultraplan-go/internal/sprint/review_test.go` | Fake-runtime and deterministic unit tests for review scope, reviewer fan-out, structured output validation, deterministic checks, verdicts, persistence, stale inputs, and citation containment. |
| CLI/app review tests | `../ultraplan-go/internal/app/sprint_review_commands_test.go` | Command and shared-use-case tests for review, validate, prompt, flow, JSON output, exit codes, cancellation, and CLI/TUI parity DTOs. |
| TUI review tests | `../ultraplan-go/internal/tui` | Deterministic model/update tests for review readiness, confirmation, progress, cancellation, finding navigation, result display, and preview behavior. |

## Acceptance Criteria

- [ ] `ultraplan sprint <project> <sprint> review` requires valid prerequisites through `execute` and fails preflight for missing, duplicate, unknown, unreadable, or path-escaping selected contracts or review protocols.
- [ ] Review resolves every selected contract and selected review protocol dynamically through `project-index.md`; it does not hardcode the active contract pool or protocol paths.
- [ ] Review computes and records a deterministic input fingerprint over governed sprint inputs, selected contracts/protocols, target implementation identity, changed-path scope, `plan.md`, `execute.md`, and `.run-state.json` where present.
- [ ] Review freezes its execution manifest while running so reviewer tasks cannot silently drift if inputs change mid-run.
- [ ] UltraPlan owns bounded reviewer fan-out and issues one independent structured agentwrap request per selected contract plus one technical-handbook request; it does not assume the runtime can spawn subagents.
- [ ] Reviewer runtime requests use read-only permissions, workspace/target-contained paths, configured review model resolution, context cancellation, bounded retries, and safe diagnostics.
- [ ] Malformed reviewer output, missing reviewer output, failed reviewer tasks, runtime failure, timeout, or cancellation are recorded and cannot be converted into a passing review.
- [ ] Deterministic product checks cover decision conformance, plan-task execution evidence, approved verification-command results, complete reviewer coverage, citation containment and valid line references, deviations, and missing evidence.
- [ ] Product code, not model prose, computes the final verdict using only `pass`, `pass_with_findings`, `fail`, or `blocked`.
- [ ] Applicable blocker or high-severity findings produce `fail`; only medium, low, or info findings may produce `pass_with_findings`.
- [ ] Runtime success is insufficient: the stage completes only after sprint-root `review.md` exists, has the required sections, contains no placeholders, passes review validation, and is reflected in `flow-state.json`.
- [ ] Re-running review atomically replaces the prior manual or generated `review.md` without modifying unrelated sprint artifacts, governed planning inputs, product source, product tests, or Git state.
- [ ] A failed or incomplete review run preserves the last complete `review.md` and exposes current failed state and diagnostics through status, JSON, and TUI.
- [ ] `flow-state.json` supports `review` after `execute` and records execution status, deterministic verdict, artifact path, timestamp, input fingerprint, staleness, and safe diagnostics.
- [ ] `ultraplan sprint <project> <sprint> status`, `validate review`, `prompt review`, `flow --to review`, and `review --json` expose review readiness, selected scope, progress/result state, verdict, staleness, and next action consistently.
- [ ] Prompt preview and dry-run show selected model source, reviewer count, bounded concurrency, target scope, selected contracts/protocols, expected writes, and permitted mutation scope without invoking runtime or writing `review.md`.
- [ ] TUI review functionality uses the same typed app use cases as CLI and supports readiness/status, prompt preview, dry-run, confirmation, live reviewer progress, cancellation, result display, finding navigation, verdict display, and `review.md` preview.
- [ ] TUI code does not invoke CLI handlers, parse terminal output, interpret native provider payloads, or persist alternate review state.
- [ ] Embedded review prompt and `review.md` template defaults are available without workspace `prompts/` or `templates/` files, and optional materialization remains explicit through `ultraplan defaults install`.
- [ ] Fake-runtime tests cover review success, pass-with-findings, blocker/high failure, blocked preflight, missing reviewer result, malformed structured output, citation path escape, invalid line reference, stale inputs, timeout, cancellation, redaction, and failed atomic write preservation.
- [ ] Required verification commands pass from `../ultraplan-go`: `go test ./...`, `go test -race ./...`, and `go build ./cmd/ultraplan`.

## Non-Goals

- Implementing `smoke.md`, smoke harness discovery, smoke execution, `ultraplan sprint <project> <sprint> smoke`, or `flow --to smoke` behavior.
- Implementing `ultraplan sprint <project> <sprint> verify`; integrated review-to-smoke verification is Sprint 28 scope.
- Building a general-purpose issue tracker, `issues.md`, `issues.json`, assignment, scheduling, remote issue synchronization, or cross-sprint verification scheduling.
- Automatically fixing product code, product tests, sprint artifacts, review findings, or smoke failures.
- Performing Git add, commit, push, branch, merge, reset, checkout, or any automatic Git mutation.
- Adding hosted services, browser UI, multi-user collaboration, or remote review execution.
- Replacing the CLI or JSON surfaces with a TUI-only workflow.
- Reimplementing OpenCode supervision, provider execution, retry policy, validation wrapping, or structured event handling outside `github.com/Antonio7098/agentwrap` and `agentwrap/opencode`.
- Copying raw runtime payloads, unrestricted stdout/stderr, or external smoke evidence into sprint artifacts.
- Treating historical manual `review.md` files or legacy `deep-smoke.md` files as current generated Phase 3 evidence unless the new review command explicitly replaces `review.md`.

## Constraints

- `internal/sprint` owns review scope, selected-contract/protocol resolution, prompt construction, reviewer orchestration, deterministic checks, verdict synthesis, validation, `review.md`, and flow-state integration.
- `internal/platform/runtime` remains generic execution infrastructure and must not import or understand project, sprint, contract, handbook, verdict, or review semantics.
- `internal/app` exposes typed review use cases shared by CLI and TUI; CLI handlers stay thin and TUI must not shell out to `ultraplan` or parse CLI stdout.
- `internal/tui` owns terminal interaction only and must not own review state machines, validation rules, runtime execution, verdict synthesis, or artifact persistence.
- Review must use `github.com/Antonio7098/agentwrap` and `agentwrap/opencode` through the existing runtime boundary; do not invent a competing runtime contract.
- Review runtime requests must be read-only with explicit permission policy and path containment; review must not write product source, product tests, governed planning inputs, workspace config, or Git state.
- Generated artifacts must use workspace-relative paths, avoid secrets, and avoid absolute paths except where explicitly necessary for local diagnostics.
- All review and flow-state writes must be atomic and must preserve the last complete artifact/state on failed replacement.
- Missing required runtime, credentials, selected context, reviewer coverage, or environment must produce `blocked` or a failed preflight, never a false pass.
- Review verdicts are deterministic and cannot be chosen by an LLM; LLM output may supply findings and summaries only after structured validation.
- Normal tests must use fake runtimes and deterministic fixtures; real OpenCode review evidence remains gated and optional unless explicitly available.
- Product source changes must remain in the target implementation repository; this planning workspace sprint creation task may mutate only `projects/ultraplan-go/sprints/26-review-stage/requirements.md`.

## Dependencies

| Prior Sprint / Output | Required For | Notes |
| --------------------- | ------------ | ----- |
| Sprint 16 Project Domain and Project Index | Dynamic selected-contract and review-protocol resolution | Project catalog parsing and path validation must be available and remain catalog-owned. |
| Sprint 17 Sprint Artifact Domain and Flow State | Review stage integration | Existing sprint discovery, artifact path rules, strict flow-state loading, and atomic flow-state writes are extended with `review`. |
| Sprint 18 Select Stage | Selected context validation | `sprint-index.md` must select only cataloged contracts, evidence, templates, and review protocols. |
| Sprint 19 Distill Stage | Technical-handbook review | `technical-handbook.md` provides selected evidence guidance and must be reviewed as its own coverage unit. |
| Sprint 20 Reason Stage | Decision conformance checks | `reasoning/*.md` and `reasoning.md` provide decisions, assumptions, risks, and expected evidence for deterministic review checks. |
| Sprint 21 Plan Stage | Plan-task conformance checks | `plan.md` provides task and evidence checklists that review must verify against execute results. |
| Sprint 23 Execute Stage | Review prerequisites and evidence | `execute.md` and `.run-state.json` provide task execution status, diagnostics, runtime metadata, and evidence. Prior Sprint 23 review identified carry-forward risks around status/help, redaction, execute model drift, run-state classification, and execute progress visibility that review should expose truthfully rather than hide. |
| Sprints 24 and 25 TUI Foundation and Operational Controls | TUI review operation | Review must be exposed through the existing local TUI patterns, confirmation policy, progress/cancellation plumbing, and shared app use cases. |
| Phase 3 PRD/TRD/Architecture updates | Review scope and constraints | Phase 3 documents define `review.md` as the current canonical review artifact, require review before smoke, and require full CLI/TUI parity. |
| `github.com/Antonio7098/agentwrap` and `agentwrap/opencode` | Runtime-backed reviewer execution | Review must use the existing generic runtime integration rather than direct OpenCode or provider execution. |

## Review Expectations

| What | How Verified |
| ---- | ------------ |
| Review command surface | Run `ultraplan sprint <project> <sprint> review`, `validate review`, `prompt review`, `flow --to review`, `status`, and JSON variants against fixtures and inspect command tests. |
| Dynamic catalog resolution | Tests and review inspect missing, duplicate, unknown, unreadable, and escaping selected contract/protocol paths. |
| Reviewer orchestration | Fake-runtime tests prove one structured request per selected contract plus one handbook request, bounded concurrency, cancellation, retry behavior, and no runtime subagent assumption. |
| Deterministic checks and verdicts | Unit tests cover decision, plan, verification, citation, coverage, severity, blocked, fail, pass, and pass-with-findings cases. |
| Artifact validation | `validate review` and unit tests verify required `review.md` sections, no placeholders, contained evidence paths, valid line references, reviewer coverage, and allowed verdict values. |
| Atomic artifact behavior | Failure-injection tests prove a failed review write does not corrupt the last complete `review.md` or flow state. |
| Flow and freshness | Status/flow tests prove `review` follows `execute`, records fingerprint/verdict/staleness, and rejects stale or changed governed inputs as current evidence. |
| CLI/TUI parity | Shared app-use-case tests and TUI model/update tests prove CLI and TUI expose the same scope, confirmation, progress, cancellation, result, verdict, staleness, and next-action data. |
| Architecture boundaries | Code review checks `internal/sprint` ownership, no `study` dependency from `sprint`, no product imports from `platform/*`, no TUI-to-CLI integration, and no direct OpenCode or shell execution from review product code. |
| Security and mutation boundaries | Tests and review inspect read-only reviewer permissions, path containment, secret redaction, safe diagnostics, and absence of source/test/Git/governed-input mutation. |
| Verification gate | Run `go test ./...`, `go test -race ./...`, and `go build ./cmd/ultraplan` from `../ultraplan-go`; record results in `execute.md` and final `review.md`. |
