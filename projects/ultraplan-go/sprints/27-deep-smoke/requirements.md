# Sprint Requirements: Deep Smoke Harness Integration

> Project: `ultraplan-go`
> Sprint: `27-deep-smoke`
> Purpose: the authoritative, human-readable sprint contract. All other sprint artifacts must satisfy these requirements.

## Sprint Goal

Implement review-gated deep smoke integration that discovers the cataloged external smoke harness, safely runs the narrowest sufficient smoke scope, keeps raw evidence in the harness, atomically writes the current sprint `smoke.md`, and exposes the same operation through CLI, JSON, status, validation, flow, and TUI use cases.

## Required Outputs

| Output | Path | Description |
| ------ | ---- | ----------- |
| Requirements | `projects/ultraplan-go/sprints/27-deep-smoke/requirements.md` | Sprint contract for the deep smoke harness integration sprint. |
| Sprint Index | `projects/ultraplan-go/sprints/27-deep-smoke/sprint-index.md` | Selected contracts, evidence reports, reasoning templates, review protocols, and source documents for the sprint. |
| Technical Handbook | `projects/ultraplan-go/sprints/27-deep-smoke/technical-handbook.md` | Distilled implementation guidance from selected evidence for external process execution, CLI/TUI parity, testing, security, observability, persistence, and workflow behavior. |
| Architecture Reasoning | `projects/ultraplan-go/sprints/27-deep-smoke/reasoning/architecture.md` | Area reasoning for package ownership, dependency direction, and the sprint/platform/process boundary for smoke. |
| Final Reasoning | `projects/ultraplan-go/sprints/27-deep-smoke/reasoning.md` | Final sprint decisions, trade-offs, assumptions, risks, and required verification evidence. |
| Implementation Plan | `projects/ultraplan-go/sprints/27-deep-smoke/plan.md` | Executable implementation task plan traced to `reasoning.md` decisions and acceptance evidence. |
| Execute Summary | `projects/ultraplan-go/sprints/27-deep-smoke/execute.md` | Summary of execute task outcomes and evidence for implementation work completed during the sprint. |
| Execute Run State | `projects/ultraplan-go/sprints/27-deep-smoke/.run-state.json` | Durable task state for controlled sprint execution, including terminal task states and diagnostics. |
| Flow State | `projects/ultraplan-go/sprints/27-deep-smoke/flow-state.json` | Versioned stage state through smoke, including execution status, verdict, artifact paths, and freshness metadata where supported. |
| Sprint Review | `projects/ultraplan-go/sprints/27-deep-smoke/review.md` | Automated conformance review for this sprint before smoke execution. |
| Smoke Summary | `projects/ultraplan-go/sprints/27-deep-smoke/smoke.md` | Current sprint smoke summary with review gate status, selected harness scope, run ID, linked evidence paths, issue references, verdict, and required next action. |
| Harness Manifest | `../ultraplan-go-smoke/ultraplan-smoke.json` | Versioned machine-readable smoke harness contract describing entrypoint, protocol version, discovery/run commands, evidence directories, and capabilities. |
| Harness Run Evidence | `../ultraplan-go-smoke/runs/` | External raw smoke run evidence produced by the harness and linked from `smoke.md`, not copied into the sprint root. |
| Harness Issue Evidence | `../ultraplan-go-smoke/issues/` | External open/resolved issue records referenced by `smoke.md` when relevant. |

## Acceptance Criteria

- [ ] `ultraplan sprint <project> <sprint> smoke` requires a current passing or non-blocking `review.md` by default and requires an explicit force/diagnostic override to run after a failed review.
- [ ] The smoke harness is discovered from the project catalog and its versioned manifest; missing, malformed, unsupported, or path-escaping manifests fail preflight with actionable diagnostics.
- [ ] Harness discovery returns structured levels, suites, tests, sprint mappings, prerequisites, expected duration/cost class, evidence directories, and protocol version without parsing README prose as executable commands.
- [ ] Smoke selection chooses the narrowest sufficient sprint-specific scope by default and supports explicit level, suite, and test overrides without letting a narrow investigation replace required containing-suite evidence.
- [ ] External harness execution uses explicit executable and argv, contained working directory, bounded environment forwarding, context cancellation, timeout, and owned descendant-process cleanup; no shell-interpolated command string is executed.
- [ ] Smoke imports and validates run ID, safe argv display, result counts, runtime/model metadata where provided, duration, evidence paths, hashes or identity metadata, and relevant open/resolved issue IDs.
- [ ] `projects/ultraplan-go/sprints/27-deep-smoke/smoke.md` is atomically replaced only after smoke completion or truthful blocked/not-applicable classification, passes smoke validation, and links raw harness evidence rather than copying it.
- [ ] Smoke verdict synthesis is deterministic: selected smoke failure produces `fail`, unavailable required prerequisites produce `blocked`, irrelevant smoke produces `not_applicable`, passing smoke with relevant open issues produces `pass_with_open_issues`, and clean successful smoke produces `pass`.
- [ ] Product source, product tests, governed planning artifacts, and Git state are not modified by normal review or smoke execution; only sprint `smoke.md`, flow/status state, and external harness run/issue evidence may change through the approved operations.
- [ ] `status`, `validate smoke`, `flow --to smoke`, text output, and stable JSON output expose smoke readiness, selected scope, execution state, verdict, artifact path, evidence links, stale/blocked diagnostics, and required next action consistently.
- [ ] TUI smoke readiness/status, scope selection, prerequisites, cost/duration class, confirmation, live suite/test progress, cancellation, result display, issue summary, and `smoke.md` preview use shared typed app use cases and do not call CLI handlers or parse terminal output.
- [ ] Fake-harness tests cover success, failure, blocked, not-applicable, timeout, cancellation, malformed discovery/run JSON, missing evidence, issue references, path escape, and redaction; gated real-harness tests are opt-in and skip truthfully when unavailable.
- [ ] `go test ./...`, `go test -race ./...`, and `go build ./cmd/ultraplan` pass in the target implementation repository, or any environment-specific inability to run a gated real harness is reported as blocked rather than pass.

## Non-Goals

- Do not implement the Sprint 28 integrated `verify` command, focused rerun recovery, or full stale-result workflow beyond the smoke freshness and review gate needed by this sprint.
- Do not add general-purpose issue tracking, issue assignment, remote issue synchronization, burndown/scheduling, or cross-sprint verification scheduling.
- Do not automatically fix product source, product tests, sprint artifacts, harness tests, or harness issue records in response to review or smoke findings.
- Do not add Git add, commit, push, branch, merge, reset, or any automatic Git mutation behavior.
- Do not copy raw smoke JSON, unrestricted stdout/stderr, per-test artifacts, or harness issue files into `projects/ultraplan-go/sprints/27-deep-smoke/`; keep detailed evidence in the external harness.
- Do not build a hosted service, browser UI, multi-user workflow, plugin marketplace, or remote smoke execution service.
- Do not make workspace `prompts/` or `templates/` files required for smoke; embedded defaults remain sufficient unless intentionally overridden.

## Constraints

- `internal/sprint` owns smoke review-gating, harness discovery semantics, scope selection, evidence validation, verdict synthesis, `smoke.md`, validation, status, and flow integration.
- `internal/platform/process` may own only generic executable/argv/cwd/env/timeout/cancellation/process-tree cleanup behavior and must not import or understand project, sprint, review, smoke, verdict, or harness issue semantics.
- `internal/app` must expose CLI and TUI smoke behavior through shared typed use cases; TUI code must not invoke CLI command handlers, shell out to `ultraplan`, parse CLI stdout, or persist alternate smoke state.
- Runtime-backed review behavior remains through `agentwrap`; smoke harness execution is an external process boundary and must not introduce a competing LLM runtime contract.
- Harness execution must use explicit argv and must never evaluate shell commands assembled from Markdown, README prose, generated artifacts, or user-edited sprint content.
- Review runs before smoke by default; blocker/high review failures stop normal smoke unless an explicit diagnostic override is confirmed and recorded.
- Missing runtime credentials, required harness prerequisites, manifest support, or required coverage must produce `blocked`, never a false pass.
- Generated sprint artifacts must be editable Markdown or versioned JSON as appropriate, use workspace-relative paths where possible, avoid secrets, and be written atomically when product-owned.
- Normal tests must use fake runtimes and fake smoke harnesses; live OpenCode or live harness tests must be opt-in and skipped by default unless required environment is present.
- Implementation must preserve module dependency direction from `docs/ARCHITECTURE.md`: product modules may depend on platform packages, platform packages must not depend on product modules, and `study` must not depend on `project` or `sprint`.

## Dependencies

| Prior Sprint / Output | Required For | Notes |
| --------------------- | ------------ | ----- |
| Sprint 23 Execute Stage | Smoke prerequisites and execute evidence | Smoke validates and reports against sprint implementation evidence produced by the execute stage. |
| Sprint 24 TUI Foundation | TUI smoke surface | Smoke must appear in the existing local terminal UI without making the TUI a separate product workflow. |
| Sprint 25 Operational TUI Controls | Guarded actions, progress, and cancellation patterns | Smoke TUI behavior must reuse existing confirmation, progress, cancellation, and result-display patterns. |
| Sprint 26 Automated Sprint Review | Review-before-smoke gate | Sprint 26 review passed with no carry-forward findings; Sprint 27 depends on current `review.md` verdict and review validation behavior. |
| `projects/ultraplan-go/project-index.md` smoke harness catalog | Harness discovery | The harness root and manifest path must resolve from the project catalog, not from hardcoded developer-specific command prose. |
| `../ultraplan-go-smoke/ultraplan-smoke.json` | Harness protocol contract | The manifest must exist or be created/updated as this sprint's external harness deliverable before real smoke can pass. |
| `../ultraplan-go-smoke/runs/` and `../ultraplan-go-smoke/issues/` | Raw evidence and issue links | Detailed evidence stays in the harness and is linked from `smoke.md`. |

## Review Expectations

| What | How Verified |
| ---- | ------------ |
| Harness manifest discovery is catalog-driven and path-safe | Unit tests, fake harness fixtures, project-index validation, and code review of path containment and manifest resolution. |
| External process execution is safe | Code review confirms explicit executable/argv use, no shell interpolation, bounded env forwarding, timeout, cancellation, and process-tree cleanup; tests cover path escape and cancellation. |
| Review gate is enforced | CLI/JSON/TUI tests demonstrate passing/non-blocking review allows smoke, failed review blocks by default, and diagnostic override is explicit and recorded. |
| Scope selection is narrow but sufficient | Tests and review inspect sprint mapping, explicit override behavior, containing-suite rules, and selection rationale in `smoke.md`. |
| Smoke verdicts are truthful | Fake-harness tests cover `pass`, `pass_with_open_issues`, `fail`, `blocked`, and `not_applicable`; validation rejects false pass cases. |
| `smoke.md` is valid and evidence-linked | Smoke validator checks required sections, run ID, safe argv, result counts, verdict, evidence paths, issue references, no placeholders, and no copied raw evidence. |
| CLI, JSON, status, flow, and TUI agree | Command tests and TUI model/update tests compare readiness, scope, progress, cancellation, verdict, artifact path, evidence links, and next action through shared use-case results. |
| Mutation boundaries are preserved | Review confirms normal smoke does not change product source, product tests, governed sprint inputs, or Git state; only allowed sprint state/artifact and harness evidence paths are touched. |
| Verification commands pass | Review checks recorded output from `go test ./...`, `go test -race ./...`, `go build ./cmd/ultraplan`, fake-harness suites, and any gated real-harness run or truthful blocked skip. |
