# Sprint Requirements: Shared Context Integration and Grounded-Planning Release

> Project: `ultraplan-go`
> Sprint: `34-shared-context`
> Purpose: the authoritative, human-readable sprint contract. All other sprint artifacts must satisfy these requirements.

## Sprint Goal

Reuse the stored, validated `requirements.md` and reference-only `code-context.md` unchanged through one stable shared prompt prefix, resolve its repository-relative source references into transient prompt evidence, and apply that prefix across every downstream agent-backed sprint operation.

## Required Outputs

| Output | Path | Description |
| --- | --- | --- |
| Shared sprint-context renderer | `../ultraplan-go/internal/sprint/prompt_context.go` | Owns deterministic shared prompt composition, preserves exact requirements/code-context bytes, and injects current source text resolved from the context pack's references without rewriting the artifact. |
| Planning-stage prompt integration | `../ultraplan-go/internal/sprint/service.go` | Applies the shared context to sprint-index, technical-handbook, area-reasoning, final-reasoning, and plan runtime requests. |
| Shared prompt helpers | `../ultraplan-go/internal/sprint/prompts.go` | Keeps stable shared instructions and the stage-specific boundary explicit and reusable. |
| Execute prompt integration | `../ultraplan-go/internal/sprint/execute.go` | Applies the same shared context before execute-specific agent instructions. |
| Conformance Review integration | `../ultraplan-go/internal/sprint/review.go` | Applies the same shared context to each agent-backed review request without changing review fan-out or verdict ownership. |
| Smoke prompt integration | `../ultraplan-go/internal/sprint/smoke_author.go` | Applies the same shared context wherever smoke suite authoring invokes an agent. |
| Cumulative-flow integration | `../ultraplan-go/internal/sprint/flow.go` | Ensures `flow --to plan` executes `code-context` exactly once after requirements and before sprint-index. |
| Manual code-context skill | `../ultraplan-go/internal/workspace/skills.go` | Registers a manual-only `ultraplan-code-context` skill that delegates stage execution to the canonical CLI operation. |
| Shared-prefix tests | `../ultraplan-go/internal/sprint/prompt_context_test.go` | Proves exact content reuse, prefix equality, composition order, live-source permission, and exclusion of dynamic execution data. |
| Planning and flow integration tests | `../ultraplan-go/internal/sprint/sprint_index_test.go` | Covers shared-context injection into cumulative planning and exact-once `code-context` execution. |
| Handbook integration tests | `../ultraplan-go/internal/sprint/handbook_test.go` | Covers exact shared-prefix reuse in technical-handbook requests. |
| Reasoning integration tests | `../ultraplan-go/internal/sprint/reasoning_test.go` | Covers exact shared-prefix reuse in area and final reasoning requests. |
| Plan integration tests | `../ultraplan-go/internal/sprint/plan_test.go` | Covers exact shared-prefix reuse and the stage-specific boundary in plan requests. |
| Execute integration tests | `../ultraplan-go/internal/sprint/execute_plan_test.go` | Covers exact shared-prefix reuse in execute requests without regressing execution behavior. |
| Review integration tests | `../ultraplan-go/internal/sprint/review_test.go` | Covers exact shared-prefix reuse by independent Conformance Review agent requests. |
| Smoke integration tests | `../ultraplan-go/internal/sprint/smoke_test.go` | Covers exact shared-prefix reuse by agent-backed smoke authoring and existing recovery behavior. |
| Stage-skill tests | `../ultraplan-go/internal/workspace/skills_test.go` | Covers code-context skill resolution, dry-run/materialization, manual-only metadata, canonical delegation, force behavior, and all-skill synchronization. |
| Product overview | `../ultraplan-go/README.md` | Documents the grounded-planning workflow, shared-context guarantee, manual skill, and supported release behavior. |
| CLI reference | `../ultraplan-go/docs/cli-reference.md` | Documents code-context commands, cumulative flow placement, status behavior, and skill materialization. |
| User guide | `../ultraplan-go/docs/user-guide.md` | Explains generating, inspecting, rerunning, and reusing a context pack through downstream stages. |
| Architecture guide | `../ultraplan-go/docs/architecture.md` | Documents `internal/sprint` ownership, renderer order, byte-stability boundary, and generic runtime separation. |
| Recovery guide | `../ultraplan-go/docs/recovery.md` | Documents cancellation, failed generation, atomic rerun, restart, and browser recovery for code-context workflows. |
| Planning smoke guide | `../ultraplan-go/docs/planning-smoke.md` | Documents fake-runtime coverage and the gated real-repository requirements-to-plan dogfood procedure and evidence. |
| Stage skills guide | `../ultraplan-go/docs/stage-skills.md` | Documents manual invocation and canonical CLI delegation for `ultraplan-code-context`. |
| Local web guide | `../ultraplan-go/docs/local-web.md` | Documents shared readiness, progress, findings, artifact, rerun, cancellation, and recovery behavior in the browser. |
| Generated-workspace guide | `../ultraplan-go/internal/workspace/scaffold/templates/README.md` | Documents code-context placement, downstream reuse, and manual skill availability in initialized workspaces. |

## Acceptance Criteria

- [ ] `code-context.md` stores references rather than copied source: every selected entry has a repository-relative `Path`, exact `Lines`, optional `Symbol`, and concrete `Rationale`; fenced source content is rejected.
- [ ] One renderer owned by `internal/sprint` composes downstream prompts in this exact order: stable shared planning instructions, sprint identity, exact `requirements.md`, exact reference-only `code-context.md`, source text resolved from those references, other shared context, then stage-specific instructions and output contract.
- [ ] For the same sprint and unchanged governed inputs, every compatible downstream request contains a byte-for-byte identical common prefix through the stage-specific boundary, including byte-for-byte exact stored requirements and code-context content plus deterministically rendered referenced source evidence.
- [ ] Timestamps, run IDs, stage names, output paths, attempts, sessions, and other dynamic execution data do not enter the common prefix; required dynamic metadata appears only after the stage-specific boundary or in runtime metadata.
- [ ] Sprint-index, technical-handbook, area-reasoning, final-reasoning, plan, execute, every independent Conformance Review request, and agent-backed smoke authoring receive the shared context whenever they invoke an agent.
- [ ] Every covered downstream prompt explicitly identifies resolved snippets as transient prepared evidence rather than content stored in `code-context.md` or an exclusive source boundary, and permits the agent to inspect additional live repository files.
- [ ] `flow --to plan` preserves the canonical `requirements -> code-context -> sprint-index -> technical-handbook -> area-reasoning -> reasoning -> plan` order and runs `code-context` exactly once.
- [ ] Explicit code-context rerun, runtime failure, validation failure, and cancellation preserve Sprint 33's truthful state transitions and atomic replacement guarantee: only `code-context.md` may be replaced, and the last valid artifact survives an unsuccessful rerun.
- [ ] CLI, TUI, and browser continue to agree on code-context readiness, progress, state, artifact preview, findings, rerun, cancellation, and recovery through shared application operations and durable sprint state.
- [ ] `ultraplan skills materialise code-context` and selection by `ultraplan-code-context` produce a manual-only skill whose execution workflow invokes the canonical `ultraplan sprint <project> <sprint> flow --to code-context` operation rather than reimplementing source selection or artifact generation.
- [ ] Skill dry-run is non-mutating, normal materialization preserves custom files, `--force` restores the built-in output, all-skill materialization includes code-context, and generated skill content and metadata remain synchronized with the embedded definition.
- [ ] Deterministic fake-runtime fixtures prove exact prefix equality across all representative downstream operations, dynamic-data exclusion, exact-once flow execution, cancellation, failed rerun preservation, and browser recovery.
- [ ] A temporary representative workspace completes a real implementation-repository requirements-to-plan flow with `code-context` inserted exactly once using a gated runtime; unavailable credentials or runtime prerequisites are recorded as blocked with the exact missing prerequisite, never as passed.
- [ ] README, CLI, user, architecture, recovery, planning-smoke, stage-skill, local-web, and generated-workspace documentation agree with executable behavior and do not claim provider cache hits.
- [ ] `go test ./...`, `go test -race ./...`, `go vet ./...`, `go build ./cmd/ultraplan`, and `git diff --check` pass in `../ultraplan-go`.
- [ ] No repository index, RAG or embedding system, retrieval subsystem, UltraPlan cache or cache key, provider-specific cache-control dependency, parallel JSON context manifest, automatic staleness detector, amendment workflow, content-identity/provenance feature, or QA/repair stage is introduced.

## Non-Goals

- Repository indexing, retrieval, RAG, embeddings, knowledge graphs, or treating `code-context.md` as the only source downstream agents may inspect.
- An UltraPlan cache, cache key, provider-specific cache-control integration, or a guarantee or measurement of provider prompt-cache hits.
- A parallel JSON context manifest, automatic source-change detection, context staleness or amendment workflows, content identity/provenance, or empirical QA/repair.
- A new generic stage framework, workflow engine, plugin system, repository-index package, runtime contract, browser-specific product service, or route-specific code-context operation semantics.
- A second persisted source-excerpt artifact or parallel JSON manifest; source text resolved from `code-context.md` references exists only in runtime prompts.
- Hosted service, remote binding, accounts, multi-user collaboration, remote workers, Aren integration, cloud authority, or automatic Git mutation.
- Post-Sprint-34 gated work described in the roadmap, including retrieval, persistence/SQLite, knowledge-graph, cloud, or repair-loop implementation.

## Constraints

- `internal/sprint` must own the shared renderer and all product-stage composition; `internal/platform/runtime` must remain generic and must not import sprint semantics.
- The renderer must preserve stored requirements/code-context bytes exactly. It may parse reference fields from `code-context.md` solely to resolve contained repository-relative paths and line ranges, but must never regenerate the artifact or persist resolved source text.
- Reference resolution must remain contained within the configured implementation repository, reject missing/out-of-range references, use deterministic reference order, and bound injected source to the runtime prompt budget.
- Stable shared content must be deterministic for the same sprint and unchanged governed inputs. Stage identity, stage output contracts, run metadata, and other request-specific values must begin only at the explicit stage-specific boundary.
- CLI, TUI, and HTTP adapters must continue to use typed shared application use cases. `internal/web` must not compose prompts, inspect repositories, duplicate stage validation, invoke the CLI binary, or persist alternate workflow truth.
- Filesystem artifacts and sprint-owned flow/run state remain authoritative. Existing atomic persistence, path containment, redaction, bounded diagnostics, cancellation, operation-conflict, and restart-recovery guarantees are binding.
- The manual skill must delegate execution to the canonical CLI code-context operation and must not duplicate implementation-target resolution, prompt composition, repository access, validation, or state transitions.
- Downstream agents must retain live repository read access; the context pack may guide inspection but must not restrict it.
- Normal tests must be offline and deterministic with fake runtimes and temporary repositories. Real-runtime dogfood must be explicitly gated and its outcome reported truthfully.
- No production source, test, governed planning input, smoke-harness evidence, or Git state may be mutated by the dogfood run except artifacts and state inside its disposable temporary workspace.
- The roadmap owns Sprint 34 scope. Later content-identity, QA/repair, retrieval, persistence, graph, and cloud directions may not be pulled into this sprint.

## Dependencies

| Prior Sprint / Output | Required For | Notes |
| --- | --- | --- |
| Sprint 33, `projects/ultraplan-go/sprints/33-code-context-stage/review.md` | Operational code-context baseline | Final verdict is `pass` with no open Sprint 33 findings; preserve its accepted stage order, validation, compatibility, atomic rerun, cancellation, artifact-preview, and recovery behavior. |
| Sprint 33 implementation in `../ultraplan-go/internal/sprint/code_context.go` | Authoritative stored context pack | Supplies validated `code-context.md`, target resolution, output isolation, and rerun semantics; this sprint reuses rather than redesigns it. |
| Project index, roadmap, PRD, TRD, and Architecture | Scope, release gates, and ownership | Define Phase 5 composition order, downstream coverage, exact-byte requirement, package ownership, exclusions, and release checks. |
| Existing downstream stage implementations in `../ultraplan-go/internal/sprint/` | Shared renderer integration | Sprint-index through smoke must converge on one renderer without replacing stage-specific prompts or output contracts. |
| Existing stage-skill materializer in `../ultraplan-go/internal/workspace/skills.go` | Manual code-context skill | Reuse current manual-only metadata, dry-run, preservation, force, and generated-file conventions. |
| Sprint 32, `projects/ultraplan-go/sprints/32-hardening-and-release/review.md` | Interface parity and gated-evidence truthfulness | Preserve the shared app/web boundary; real-runtime or harness evidence must be produced or recorded with an exact blocker, never inferred as passed. |
| Configured implementation repository `../ultraplan-go/` and an available gated runtime | Real-repository dogfood | Use a disposable representative workspace; absence of runtime credentials or prerequisites is a truthful blocked result. |

## Review Expectations

| What | How Verified |
| --- | --- |
| Composition ownership and order | Architecture review of `internal/sprint/prompt_context.go` and call sites confirms one renderer, the mandated ordering, and no sprint semantics in generic runtime code. |
| Exact-byte and prefix stability | Golden/table-driven fixtures compare raw prompt bytes across all covered stages, compare embedded artifact slices directly with stored file bytes, and verify deterministic source injection from reference order. |
| Dynamic-data exclusion | Fixtures vary stage, output path, timestamp, run ID, attempt, and session metadata and prove the common prefix remains unchanged. |
| Downstream coverage and live inspection | Focused fake-runtime tests capture sprint-index, handbook, area/final reasoning, plan, execute, each review request, and smoke-author prompts; each contains the shared prefix and additional-source permission. |
| Flow sequencing | Cumulative-flow tests and a fake-runtime call log prove `code-context` appears once, immediately after requirements and before sprint-index, in `flow --to plan`. |
| Sprint 33 behavior preservation | Regression tests cover validation, cancellation, runtime failure, atomic rerun, last-valid-artifact preservation, compatibility, bounded preview, and restart recovery. |
| Interface agreement | Shared app, CLI/JSON, TUI, and browser contract tests compare readiness, progress, findings, artifact, rerun, cancellation, terminal result, and recovered state. |
| Manual skill | Skill resolution, dry-run, materialization, customization preservation, force, embedded rendering, metadata, and synchronization tests inspect generated `SKILL.md` and `agents/openai.yaml` and verify canonical CLI delegation. |
| Documentation completeness | Review each required documentation path against current CLI help, renderer behavior, recovery semantics, and tested dogfood commands; reject cache-hit claims or stale stage ordering. |
| Real-repository dogfood | Inspect command log, temporary workspace status, generated valid `code-context.md`, captured downstream prompts, exact-once runtime call sequence, and final valid plan; blocked evidence must name the exact prerequisite. |
| Release gates | Review logs for `go test ./...`, `go test -race ./...`, `go vet ./...`, `go build ./cmd/ultraplan`, and `git diff --check` from `../ultraplan-go`. |
| Scope exclusions | Diff and dependency review confirm no cache/retrieval/index/manifest/staleness/provenance/QA/repair/framework/cloud work, no route-specific web semantics, and no dogfood mutation outside its disposable workspace. |
