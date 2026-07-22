# Sprint Index: Deep Smoke Harness Integration

> Project: `ultraplan-go`
> Sprint: `27-deep-smoke`
> Purpose: selected context for this sprint. Must be a subset of `projects/ultraplan-go/project-index.md`.
> **Inputs Used:** `projects/ultraplan-go/project-index.md`, `projects/ultraplan-go/sprints/27-deep-smoke/requirements.md`, `projects/ultraplan-go/roadmap.md`, `projects/ultraplan-go/docs/ARCHITECTURE.md`, `projects/ultraplan-go/docs/PRD.md`, `projects/ultraplan-go/docs/TRD.md`

This document selects what must be read, distilled, reasoned through, or checked for this sprint. It does not make implementation decisions. All selections must come from the project index - no items may be included that are not listed in the project index.

## Sprint Scope

- **Sprint Goal:** Implement review-gated deep smoke integration that discovers the cataloged external smoke harness, safely runs the narrowest sufficient scope, keeps raw evidence in the harness, atomically writes the current sprint `smoke.md`, and exposes the operation consistently through CLI, JSON, status, validation, flow, and TUI use cases.
- **Planned Output:** The versioned external harness manifest and evidence contract; product-owned harness discovery, review gating, scope selection, safe invocation, evidence validation, deterministic verdicts, atomic `smoke.md` and flow-state updates; CLI/TUI surfaces; and deterministic fake-harness plus gated real-harness verification described by `requirements.md`.
- **Depends On:** Sprint 23 execute evidence; Sprint 24 shared TUI foundation; Sprint 25 guarded progress, confirmation, and cancellation controls; Sprint 26 current `review.md` verdict and validation behavior, which passed with no carry-forward findings; and the project-index smoke harness catalog with its external manifest, `runs/`, and `issues/` paths.
- **Non-Goals:** Sprint 28 integrated `verify`, focused-rerun recovery, and full stale-result workflow; extending execute or review automation; general-purpose issue tracking; automatic fixes; Git mutation; copying raw harness evidence into the sprint; hosted, browser, multi-user, plugin-marketplace, or remote-smoke services; and requiring workspace prompt/template overrides.

## Source Project Index

- `projects/ultraplan-go/project-index.md` - authoritative source. Any file or item referenced below must appear there.

## Selected Contracts

Each contract applies as a flat whole to this sprint. All paths must appear in the project index's "Active Contract Pool" table.

| Contract | Why Selected |
| -------- | ------------ |
| Architecture | Governs `internal/sprint` ownership of smoke semantics, the generic `internal/platform/process` boundary, dependency direction, and shared CLI/TUI app use cases. |
| CLI Surface | Governs smoke flags and overrides, help, exit behavior, text/JSON output, status, validation, and flow integration. |
| Configuration | Governs preflight configuration, bounded environment forwarding, path and timeout settings, effective diagnostics, and redaction. |
| Documentation | Governs the embedded `smoke.md` template, generated summary readability, help, evidence links, and maintainable user-facing behavior. |
| Errors | Governs actionable classification of review-gate, manifest, discovery, process, timeout, cancellation, evidence, validation, and blocked failures. |
| LLM Evaluation / Cost / Safety | Governs safe reporting of runtime/model and cost-class metadata and truthful validation of real-system evidence without false passes. |
| LLM Runtime | Governs preservation of the existing agentwrap review boundary and prevents external smoke execution from becoming a competing LLM runtime contract. |
| Observability | Governs smoke readiness, progress, run metadata, result counts, safe diagnostics, evidence links, cancellation, and truthful status across interfaces. |
| Performance | Governs bounded process output and resource use, timeout behavior, responsive progress, and duration/cost visibility. |
| Persistence And Migrations | Governs atomic `smoke.md` and flow-state replacement, versioned state and manifest handling, strict loading, evidence identity, and preservation of the last valid artifact. |
| Security | Governs explicit argv execution, contained paths and working directories, bounded environment forwarding, secret redaction, evidence containment, and mutation boundaries. |
| Testing | Governs fake-harness, fixture, command, race, TUI model/update, failure-injection, and gated real-harness coverage. |
| Workflows | Governs review-before-smoke gating, scope selection, cancellation, state transitions, deterministic verdicts, blocked/not-applicable outcomes, and flow integration. |

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
| `12-extensibility` | `studies/go-cli-study/reports/final/12-extensibility.md` | Extension points, plugin architecture, package boundaries, API design |
| `13-security` | `studies/go-cli-study/reports/final/13-security.md` | Path safety, secrets, command injection risks, sandboxing |
| `14-performance` | `studies/go-cli-study/reports/final/14-performance.md` | Startup latency, large repos, memory management, performance |

## Selected Reasoning Templates

All paths must appear in the project index's "Available Reasoning Templates" table.

| Template | Output Path | Why Selected |
| -------- | ----------- | ------------ |
| Architecture | `projects/ultraplan-go/sprints/27-deep-smoke/reasoning/architecture.md` | Reason through sprint/platform/process ownership, dependency direction, external harness and persistence boundaries, and shared CLI/TUI use cases. |

## Prior Decisions To Carry Forward

No prior decision artifacts are cataloged in the project index. Sprint 26's passing review with no carry-forward findings and the Sprint 23-26 dependencies are recorded in Sprint Scope as requirement-sourced constraints rather than represented as uncataloged decisions.

## Required Review Protocols

All paths must appear in the project index's "Review Protocols" table.

| Protocol | Path | Required Evidence |
| -------- | ---- | ----------------- |
| Architecture Review | `system/protocols/architecture-review-protocol.md` | Evidence that smoke semantics remain in `internal/sprint`, external execution remains generic in `internal/platform/process`, app use cases are shared, TUI owns interaction only, and dependencies follow the documented direction. |
| Deep Smoke Sprint | `system/protocols/deep-smoke-sprint-protocol.md` | Evidence of review gating, catalog and manifest discovery, narrowest-sufficient selection, safe harness execution, truthful verdicts, external raw evidence, and a valid linked `smoke.md`. |
| Sprint Review | `system/protocols/review-sprint-protocol.md` | Evidence that selected contracts and handbook guidance are covered, acceptance criteria and mutation boundaries are checked, required verification is recorded, and implementation is ready for smoke. |

## Excluded Context

| Context | Reason Excluded | Revisit If |
| ------- | --------------- | ---------- |
| Implementation execution | Sprint 27 consumes completed execute evidence but does not extend plan-task execution or use smoke to modify product source or tests. | A smoke acceptance criterion exposes a concrete compatibility defect that must be scheduled as separate implementation work. |
| Review automation | Sprint 26 already owns automated review; this sprint consumes its current validated verdict as a gate and does not redesign reviewer orchestration or `review.md` generation. | A required smoke gate cannot be expressed through the delivered review result and validation surface. |
| Integrated verification and focused-rerun recovery | The `verify` convenience command, complete freshness workflow, and focused rerun recovery are explicitly deferred to Sprint 28. | Sprint 28 begins integrated review-to-smoke verification. |
| Issue tracking | Smoke may import and display relevant external harness issue IDs, but assignment, scheduling, remote synchronization, and general issue-management behavior remain deferred. | A future indexed sprint explicitly introduces issue-management behavior. |
| Git mutation | Normal review and smoke must not add, commit, push, branch, merge, reset, checkout, or otherwise mutate Git state. | A future product requirement and contract explicitly authorize a bounded Git operation. |
| Automatic fixes | Smoke reports validated outcomes and next actions but must not alter product source, product tests, governed sprint artifacts, harness tests, or harness issue records in response to findings. | A future indexed sprint explicitly introduces guarded remediation or harness-maintenance actions. |
| Raw evidence duplication | Raw smoke JSON, unrestricted stdout/stderr, per-test artifacts, and issue files remain in the cataloged external harness rather than the sprint root. | The indexed harness contract or retention requirements explicitly change evidence ownership. |
| Hosted and remote surfaces | Hosted services, browser UI, multi-user workflows, plugin marketplace behavior, and remote smoke execution are outside this local CLI/TUI sprint. | A future indexed sprint explicitly brings one of these surfaces into scope. |
| Workspace prompt/template dependency | Embedded defaults are sufficient; workspace `prompts/` and `templates/` files remain optional intentional overrides. | A future requirement changes the embedded-default policy. |

## Next Artifacts

- `technical-handbook.md` reads from the evidence reports listed above.
- `reasoning/*.md` captures area-specific reasoning.
- `reasoning.md` makes final sprint decisions.
- `plan.md` executes `reasoning.md`.
- `review.md` runs the selected review protocols against implementation.
