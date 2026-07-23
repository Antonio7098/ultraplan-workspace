# Sprint Requirements: Integrated Review-to-Smoke Verification Flow

> Project: `ultraplan-go`
> Sprint: `28-review-to-smoke-flow`
> Purpose: the authoritative, human-readable sprint contract. All other sprint artifacts must satisfy these requirements.

## Sprint Goal

Make `execute -> review -> smoke` a coherent, resumable verification flow with freshness detection, focused reruns, deterministic overall assessment, recovery guidance, and complete CLI/JSON/TUI parity.

## Required Outputs

| Output | Path | Description |
| --- | --- | --- |
| Requirements | `projects/ultraplan-go/sprints/28-review-to-smoke-flow/requirements.md` | This sprint contract for the integrated review-to-smoke flow. |
| Sprint index | `projects/ultraplan-go/sprints/28-review-to-smoke-flow/sprint-index.md` | Selected contracts, evidence, reasoning templates, and protocols for Sprint 28. |
| Technical handbook | `projects/ultraplan-go/sprints/28-review-to-smoke-flow/technical-handbook.md` | Distilled implementation guidance from the selected evidence. |
| Reasoning | `projects/ultraplan-go/sprints/28-review-to-smoke-flow/reasoning.md` | Final decisions, tradeoffs, risks, and evidence expectations for Sprint 28. |
| Implementation plan | `projects/ultraplan-go/sprints/28-review-to-smoke-flow/plan.md` | Executable task plan that traces to `reasoning.md`. |
| Execute summary | `projects/ultraplan-go/sprints/28-review-to-smoke-flow/execute.md` | Summary of Sprint 28 implementation execution and task evidence. |
| Flow state | `projects/ultraplan-go/sprints/28-review-to-smoke-flow/flow-state.json` | Versioned stage state through `smoke`, including review/smoke freshness and verdict metadata. |
| Execute run state | `projects/ultraplan-go/sprints/28-review-to-smoke-flow/.run-state.json` | Versioned task state for Sprint 28 execute work. |
| Sprint review | `projects/ultraplan-go/sprints/28-review-to-smoke-flow/review.md` | Automated Sprint 28 conformance review once implementation is complete. |
| Sprint smoke summary | `projects/ultraplan-go/sprints/28-review-to-smoke-flow/smoke.md` | Sprint 28 smoke summary linking external harness run and issue evidence, or a truthful blocked/not-applicable result. |
| Sprint flow integration | `../ultraplan-go/internal/sprint/flow.go` | Enforce stage ordering through `smoke` and integrate review/smoke state into flow execution. |
| Sprint state model | `../ultraplan-go/internal/sprint/state.go` | Persist review and smoke status, verdicts, artifacts, fingerprints, run IDs, staleness, and next actions in flow state. |
| Sprint domain model | `../ultraplan-go/internal/sprint/domain.go` | Define any required verify, freshness, rerun, overall-assessment, and recovery result types. |
| Sprint service use cases | `../ultraplan-go/internal/sprint/service.go` | Provide product-owned verification orchestration without moving workflow semantics into app, CLI, TUI, runtime, or process packages. |
| Review integration | `../ultraplan-go/internal/sprint/review.go` | Support current/stale review detection and focused review reruns while preserving final verdict rules. |
| Smoke integration | `../ultraplan-go/internal/sprint/smoke.go` | Support current/stale smoke detection, focused smoke reruns, containing-suite evidence requirements, issue-aware verdicts, and recovery from missing harness evidence. |
| Sprint validation | `../ultraplan-go/internal/sprint/validation.go` | Validate integrated review/smoke artifacts, freshness, stage consistency, and deterministic overall assessment. |
| CLI/app use cases | `../ultraplan-go/internal/app/sprint_usecases.go` | Expose typed app use cases for status, flow, review, smoke, verify, rerun, and recovery data shared by CLI and TUI. |
| CLI commands | `../ultraplan-go/internal/app/sprint_commands.go` | Add `ultraplan sprint <project> <sprint> verify [--to review|smoke]`, update `flow --to smoke`, and keep text/JSON output aligned with the typed use cases. |
| TUI model and actions | `../ultraplan-go/internal/tui/model.go` | Expose the complete normal Phase 3 workflow from execute status through review and smoke with gate explanations, rerun controls, evidence links, overall assessment, and recovery guidance. |
| TUI rendering | `../ultraplan-go/internal/tui/views.go` | Render review/smoke state, staleness, verdicts, linked evidence, next action, and recovery information without parsing CLI output. |
| Sprint tests | `../ultraplan-go/internal/sprint/verify_test.go` | Cover integrated flow ordering, freshness, stale propagation, focused reruns, containing-suite evidence, issue-aware smoke verdicts, interruption recovery, and deterministic assessment. |
| CLI/app tests | `../ultraplan-go/internal/app/sprint_verify_commands_test.go` | Cover `verify`, `flow --to smoke`, JSON/text parity, diagnostic override handling, exit codes, and status next-action rendering. |
| TUI tests | `../ultraplan-go/internal/tui/verify_test.go` | Cover end-to-end verification actions, confirmations, cancellation/recovery states, focused rerun controls, evidence links, and agreement with app use-case results. |

## Acceptance Criteria

- [ ] `ultraplan sprint <project> <sprint> flow --to smoke` and `ultraplan sprint <project> <sprint> verify --to smoke` always run or require a current review before smoke unless an explicit diagnostic override is selected.
- [ ] A governed input, selected catalog item, execute evidence, target implementation fingerprint, `review.md`, or `smoke.md` change marks prior review/smoke evidence stale; a stale review also makes smoke stale.
- [ ] Sprint status in text, JSON, and TUI shows review and smoke execution state, verdict, staleness, artifact path, smoke run ID or blocked/not-applicable reason, relevant open harness issues, and required next action.
- [ ] Focused review reruns and focused smoke level/suite/test reruns are available without allowing a passing narrow smoke test to replace required evidence from its containing suite.
- [ ] Relevant open harness issues prevent a clean smoke pass and are listed in both `smoke.md` and sprint status.
- [ ] The overall assessment is deterministic from current `review.md`, `smoke.md`, `flow-state.json`, and referenced harness issues; it cannot contradict review or smoke stage verdicts and does not create a third assessment artifact.
- [ ] Interrupted review or smoke can be rerun without corrupting `flow-state.json` or replacing the last complete `review.md` or `smoke.md` with partial output.
- [ ] Recovery paths cover stale inputs, malformed review/smoke summaries, missing harness evidence, externally edited artifacts, blocked environment, cancellation, and timeout.
- [ ] CLI text, stable JSON output, and TUI views agree on current state, verdict, staleness, evidence links, and next action.
- [ ] TUI verification actions call shared typed app use cases and do not invoke CLI handlers, shell out to `ultraplan`, parse terminal output, or persist alternate verification state.
- [ ] Product code does not modify product source, product tests, governed sprint inputs, or Git state during review, smoke, verify, status, freshness checks, or recovery inspection.
- [ ] Normal tests use fake review runtimes and a fake smoke harness; no normal test requires OpenCode, provider credentials, network access, or the real smoke harness.
- [ ] `go test ./...`, `go test -race ./...`, and `go build ./cmd/ultraplan` pass in `../ultraplan-go`.

## Non-Goals

- Implementing general-purpose issue tracking, issue assignment, issue scheduling, burndown, or remote issue synchronization.
- Automatically fixing product source, product tests, `review.md` findings, smoke failures, or harness issues.
- Performing Git add, commit, push, branch, merge, reset, checkout, or other Git state mutation.
- Adding a third sprint-root assessment artifact beyond current `review.md` and `smoke.md`.
- Copying raw smoke run JSON, unrestricted stdout/stderr, per-test artifacts, or harness issue files into the sprint directory.
- Replacing Sprint 26 review orchestration or Sprint 27 smoke harness integration wholesale; this sprint integrates, hardens, and composes those stages.
- Building hosted services, browser UI, multi-user collaboration, or cross-project/cross-sprint verification scheduling.
- Making the TUI the source of truth or a separate workflow engine.

## Constraints

- `internal/sprint` owns review and smoke stage semantics, flow ordering, freshness, rerun policy, validation, verdicts, overall assessment, and sprint artifact paths.
- `internal/platform/runtime` remains generic agentwrap-backed execution infrastructure and must not learn sprint, review, smoke, contract, handbook, harness, or verdict semantics.
- `internal/platform/process` remains a generic explicit executable/argv/cwd/env/timeout/cancellation boundary and must not parse smoke protocols or synthesize verdicts.
- Review must run before smoke by default; blocker/high applicable review findings stop default smoke execution unless an explicit diagnostic override is selected and recorded.
- Freshness fingerprints must cover governed planning inputs, selected contracts/protocols, execute evidence, target implementation identity, changed-path scope, and current review/smoke artifacts where relevant.
- Flow state and run state writes must be versioned, strict on load, and atomic; failed or interrupted runs must not corrupt the last complete Markdown artifacts.
- Smoke execution must use discovered or configured harness commands as executable plus argv only; it must not evaluate shell commands assembled from Markdown or README prose.
- Missing runtime, credentials, harness coverage, environment prerequisites, evidence paths, or malformed harness output must produce `blocked` or `fail`, never a false pass.
- Detailed smoke evidence remains under the external harness root cataloged in `project-index.md`; sprint `smoke.md` contains stable run IDs and links only.
- CLI and TUI must use shared typed app use cases; TUI must not call CLI command handlers, parse stdout, or maintain alternate verification persistence.
- Workspace-relative paths must be used in generated artifacts and diagnostics unless an absolute local path is required for explicit operator recovery.
- Embedded prompt/template defaults remain sufficient; workspace `prompts/` and `templates/` files are optional intentional overrides, not prerequisites.
- Automatic product source/test mutation and all Git mutation remain prohibited during Phase 3 verification operations.

## Dependencies

| Prior Sprint / Output | Required For | Notes |
| --- | --- | --- |
| Sprint 23 `execute.md` and `.run-state.json` behavior | Review prerequisites and execute evidence freshness | Execute state is the base stage before review and smoke. Prior review recorded follow-ups around execute status, redaction, and tests; Sprint 28 must not regress execute visibility or safety. |
| Sprint 26 automated review stage | Current review status, verdict, artifacts, validation, and TUI review operation | Sprint 26 review passed and established `review.md` as the product-owned current review artifact. |
| Sprint 27 deep smoke harness integration | Smoke harness discovery, execution, `smoke.md`, issue links, and TUI smoke operation | Sprint 27 artifacts exist, but no Sprint 27 `review.md` was present when these requirements were created; Sprint 28 should validate actual Sprint 27 behavior during planning/execution. |
| `projects/ultraplan-go/project-index.md` smoke harness catalog | Harness resolution and evidence containment | The cataloged harness is `ultraplan-go-smoke` at `/home/antonioborgerees/coding/ultraplan-go-smoke/`; detailed evidence remains there. |
| `projects/ultraplan-go/docs/ARCHITECTURE.md` | Module boundaries and dependency direction | Sprint behavior stays in `internal/sprint`; app/CLI/TUI are local surfaces over shared use cases. |
| `projects/ultraplan-go/docs/PRD.md` and `projects/ultraplan-go/docs/TRD.md` | Product and technical acceptance for Phase 3 | Phase 3 requires review-before-smoke, freshness, focused reruns, recovery, stable JSON, and TUI parity. |
| Fake review runtime and fake smoke harness test seams | Deterministic normal test suite | Real OpenCode and real harness tests remain gated. |

## Review Expectations

| What | How Verified |
| --- | --- |
| Review-before-smoke ordering | Unit and command tests for `flow --to smoke`, `verify --to smoke`, failed review gate, and explicit diagnostic override. |
| Freshness and stale propagation | Tests that mutate governed inputs, execute evidence, target fingerprint, `review.md`, and `smoke.md`, then compare CLI text, JSON, and TUI state. |
| Focused reruns and containing-suite evidence | Fake harness tests proving focused test reruns do not hide required suite failures or replace containing-suite evidence. |
| Issue-aware smoke verdicts | Fake harness tests with relevant open/resolved issues and validation that `smoke.md` plus status list the relevant issue IDs and next action. |
| Deterministic overall assessment | Golden or fixture tests showing assessment derivation from stage verdicts, staleness, flow state, and harness issues with no contradictory output. |
| Recovery behavior | Tests for cancellation, timeout, malformed artifacts, missing harness evidence, stale inputs, and interrupted rerun preserving prior complete artifacts. |
| CLI/JSON/TUI parity | Shared app-use-case tests plus command/TUI tests comparing state, verdict, staleness, artifacts, run IDs, issues, and next action. |
| Module boundaries | Code review and dependency checks confirming `internal/sprint` owns verification semantics and platform packages remain product-neutral. |
| Mutation boundaries | Tests or review evidence that review/smoke/verify/status do not modify product source, product tests, governed planning inputs, or Git state. |
| Build and test gate | `go test ./...`, `go test -race ./...`, and `go build ./cmd/ultraplan` from `../ultraplan-go`. |
