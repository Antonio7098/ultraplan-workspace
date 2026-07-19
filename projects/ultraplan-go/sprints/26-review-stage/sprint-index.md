# Sprint Index: Automated Sprint Review

> Project: `ultraplan-go`
> Sprint: `26-review-stage`
> Purpose: selected context for this sprint. Must be a subset of `projects/ultraplan-go/project-index.md`.
> **Inputs Used:** `projects/ultraplan-go/project-index.md`, `projects/ultraplan-go/sprints/26-review-stage/requirements.md`, `projects/ultraplan-go/docs/ARCHITECTURE.md`, `projects/ultraplan-go/docs/PRD.md`, `projects/ultraplan-go/docs/TRD.md`

This document selects what must be read, distilled, reasoned through, or checked for this sprint. It does not make implementation decisions. All selections must come from the project index - no items may be included that are not listed in the project index.

## Sprint Scope

- **Sprint Goal:** Replace manual sprint review with a product-owned, evidence-grounded review stage that writes the current sprint-root `review.md` and is fully operable from CLI and TUI through shared typed use cases.
- **Planned Output:** Automated review-stage domain, orchestration, validation, deterministic verdicts, atomic artifact and flow-state persistence, CLI/TUI operations, embedded defaults, and deterministic tests described by `requirements.md`.
- **Depends On:** Sprints 16-21 and 23-25 outputs identified in `requirements.md`, especially project catalog resolution, sprint artifact and flow-state handling, selected-context validation, handbook/reasoning/plan artifacts, execute evidence, and shared TUI operational controls. Carry forward Sprint 23 risks around truthful status/help, redaction, execute model drift, run-state classification, and execute progress visibility.
- **Non-Goals:** Smoke investigation or execution, integrated `verify`, general-purpose issue tracking, automatic fixes, Git mutation, hosted or browser operation, remote review execution, TUI-only workflows, and reimplementation of agentwrap/OpenCode supervision.

## Source Project Index

- `projects/ultraplan-go/project-index.md` - authoritative source. Any file or item referenced below must appear there.

## Selected Contracts

Each contract applies as a flat whole to this sprint. All paths must appear in the project index's "Active Contract Pool" table.

| Contract | Why Selected |
| -------- | ------------ |
| Architecture | Governs sprint ownership of review semantics, product/platform separation, dependency direction, and shared CLI/TUI use cases. |
| CLI Surface | Governs review, validation, prompt, flow, status, help, text/JSON output, and exit-code behavior. |
| Configuration | Governs review model resolution, command preflight, runtime mapping, effective-config diagnostics, and redaction. |
| Documentation | Governs embedded prompt/template defaults, help, generated review artifacts, and maintainable user-facing behavior. |
| Errors | Governs preflight, runtime, validation, cancellation, persistence, and partial-review failure classification and diagnostics. |
| LLM Evaluation / Cost / Safety | Governs safe, evidence-grounded model evaluation, usage metadata, reviewer validation, and the boundary between model findings and product verdicts. |
| LLM Runtime | Governs agentwrap/OpenCode integration, structured review requests, runtime capabilities, permissions, cancellation, and validated results. |
| Observability | Governs reviewer progress, status truthfulness, safe diagnostics, structured output, run metadata, and failed/incomplete review visibility. |
| Performance | Governs bounded reviewer fan-out, concurrency, cancellation, resource use, and handling of review inputs and artifacts. |
| Persistence And Migrations | Governs atomic `review.md` and flow-state replacement, durable schema changes, strict loading, fingerprints, and preservation of last complete state. |
| Security | Governs read-only runtime permissions, workspace and target path containment, citation safety, secret redaction, and mutation boundaries. |
| Testing | Governs deterministic unit, fixture, fake-runtime, command, race, TUI model/update, failure-injection, and gated integration coverage. |
| Workflows | Governs review orchestration, bounded retries, cancellation, frozen manifests, state transitions, freshness, recovery, and flow integration after execute. |

## Selected Evidence Reports

Copied from the project index's "Available Evidence Reports" table. These tell the technical handbook which reports to read - the project index is the authoritative source.

| Report | Path | Covers |
| ------ | ---- | ------ |
| `01-project-structure` | `studies/go-cli-study/reports/final/01-project-structure.md` | Project layout, cmd/internal/pkg, dependency direction, thin entrypoints |
| `02-command-architecture` | `studies/go-cli-study/reports/final/02-command-architecture.md` | Command routing, flags, help text, command organization, shell completion |
| `03-dependency-injection` | `studies/go-cli-study/reports/final/03-dependency-injection.md` | Dependency construction, seams, testability, constructor patterns |
| `04-configuration-management` | `studies/go-cli-study/reports/final/04-configuration-management.md` | Config loading, precedence, environment variables, paths |
| `05-error-handling` | `studies/go-cli-study/reports/final/05-error-handling.md` | Error wrapping, classification, user-facing diagnostics, sentinel errors |
| `06-io-abstraction` | `studies/go-cli-study/reports/final/06-io-abstraction.md` | Filesystem/stdin/stdout abstraction, test seams, interface design |
| `07-state-context` | `studies/go-cli-study/reports/final/07-state-context.md` | Context propagation, app state, cancellation |
| `08-concurrency` | `studies/go-cli-study/reports/final/08-concurrency.md` | Worker pools, cancellation, parallel execution, goroutine management |
| `09-terminal-ux` | `studies/go-cli-study/reports/final/09-terminal-ux.md` | Terminal output, progress indicators, human UX, color and formatting |
| `10-logging-observability` | `studies/go-cli-study/reports/final/10-logging-observability.md` | Logs, diagnostics, structured events, observability patterns |
| `11-testing-strategy` | `studies/go-cli-study/reports/final/11-testing-strategy.md` | Unit, integration, fixture, command-level tests, coverage strategy |
| `13-security` | `studies/go-cli-study/reports/final/13-security.md` | Path safety, secrets, command injection risks, sandboxing |
| `14-performance` | `studies/go-cli-study/reports/final/14-performance.md` | Startup latency, large repos, memory management, performance |
| `15-philosophy` | `studies/go-cli-study/reports/final/15-philosophy.md` | Cross-cutting design philosophy, tradeoffs, maintainability |

## Selected Reasoning Templates

All paths must appear in the project index's "Available Reasoning Templates" table.

| Template | Output Path | Why Selected |
| -------- | ----------- | ------------ |
| Architecture | `projects/ultraplan-go/sprints/26-review-stage/reasoning/architecture.md` | Reason through Phase 3 review ownership, dependency direction, runtime boundaries, persistence boundaries, and CLI/TUI sharing. |

## Prior Decisions To Carry Forward

No prior decision artifacts are cataloged in the project index. Sprint 23 carry-forward risks are recorded in Sprint Scope as dependency constraints sourced from this sprint's requirements rather than represented as uncataloged decisions.

## Required Review Protocols

All paths must appear in the project index's "Review Protocols" table.

| Protocol | Path | Required Evidence |
| -------- | ---- | ----------------- |
| Architecture Review | `system/protocols/architecture-review-protocol.md` | Evidence that review semantics remain in `internal/sprint`, platform runtime stays generic, app use cases are shared, TUI owns interaction only, and dependencies follow the documented direction. |
| Sprint Review | `system/protocols/review-sprint-protocol.md` | Evidence that all selected contracts and handbook guidance are covered, deterministic checks and verdict rules are applied, required verification passes, and the current valid `review.md` is produced. |

## Excluded Context

| Context | Reason Excluded | Revisit If |
| ------- | --------------- | ---------- |
| Implementation execution | The execute stage and product-code mutation by review are already delivered prerequisites; this sprint consumes execute evidence and implements review without extending execute task execution. | A review acceptance criterion requires a concrete execute-stage compatibility change needed to expose truthful prerequisite evidence. |
| Smoke investigation | `smoke.md`, harness discovery and invocation, smoke execution, and `flow --to smoke` are explicit non-goals for Sprint 26. | Sprint 27 begins deep-smoke implementation. |
| Integrated review-to-smoke verification | The `verify` convenience workflow is explicitly deferred from this sprint. | Sprint 28 integrates review and smoke verification. |
| Issue tracking | General-purpose issue records, assignment, scheduling, synchronization, and cross-sprint verification scheduling are deferred product scope. | A future indexed sprint explicitly introduces issue-management behavior. |
| Git mutation | Review must not add, commit, push, branch, merge, reset, checkout, or otherwise mutate Git state. | A future product requirement and contract explicitly authorize a bounded Git operation. |
| Automatic fixes | Review reports validated findings and deterministic verdicts but must not alter product source, product tests, governed planning inputs, or findings automatically. | A future indexed sprint explicitly introduces guarded remediation. |
| Deep Smoke Sprint protocol | This protocol applies to completed sprints requiring real-runtime smoke evidence, while this sprint implements review only. | Deep-smoke behavior enters an indexed sprint's scope. |
| Extensibility evidence | Plugin architecture and generalized extension points are not needed to implement the defined review stage through the existing agentwrap and app boundaries. | Requirements expand review to additional runtime or plugin mechanisms. |

## Next Artifacts

- `technical-handbook.md` reads from the evidence reports listed above.
- `reasoning/*.md` captures area-specific reasoning.
- `reasoning.md` makes final sprint decisions.
- `plan.md` executes `reasoning.md`.
- `review.md` runs the selected review protocols against implementation.
