# Sprint Index: Integrated Review-to-Smoke Verification Flow

> Project: `ultraplan-go`
> Sprint: `28-review-to-smoke-flow`
> Purpose: selected context for this sprint. Must be a subset of `projects/ultraplan-go/project-index.md`.
> **Inputs Used:** `projects/ultraplan-go/project-index.md`, `projects/ultraplan-go/sprints/28-review-to-smoke-flow/requirements.md`, `projects/ultraplan-go/docs/ARCHITECTURE.md`, `projects/ultraplan-go/docs/PRD.md`, `projects/ultraplan-go/docs/TRD.md`

This document selects what must be read, distilled, reasoned through, or checked for this sprint. It does not make implementation decisions. All selections must come from the project index; no items may be included that are not listed in the project index.

## Sprint Scope

- **Sprint Goal:** Make `execute -> review -> smoke` a coherent, resumable verification flow with freshness detection, focused reruns, deterministic overall assessment, recovery guidance, and complete CLI/JSON/TUI parity.
- **Planned Output:** An integrated review-to-smoke flow through `smoke`, including versioned verification state, shared typed app use cases, CLI and TUI operations, deterministic status and assessment, focused reruns, recovery behavior, and tests.
- **Depends On:** Sprint 23 execute state and summary behavior, Sprint 26 automated review, Sprint 27 external smoke-harness integration, and the cataloged `ultraplan-go-smoke` harness.
- **Non-Goals:** General-purpose issue tracking; automatic product, test, finding, failure, or harness-issue fixes; Git mutation; a third sprint-root assessment artifact; copying detailed harness evidence into the sprint; wholesale replacement of Sprint 26 or Sprint 27 behavior; hosted or browser services; multi-user collaboration; cross-project or cross-sprint scheduling; and a TUI-owned workflow engine.

## Source Project Index

- `projects/ultraplan-go/project-index.md` is the authoritative source. Any file or item referenced below must appear there.

## Selected Contracts

Each contract applies as a flat whole to this sprint. All paths must appear in the project index's "Active Contract Pool" table.

| Contract | Why Selected |
| --- | --- |
| Architecture | Governs ownership of verification semantics in `internal/sprint`, shared app use cases, thin CLI/TUI surfaces, and product-neutral runtime/process boundaries. |
| Errors | Applies to deterministic verdict failures, blocked states, cancellation, timeout, malformed evidence, actionable recovery guidance, and CLI exit behavior. |
| Configuration | Applies to discovered or configured harness commands, review runtime selection, timeout settings, safe environment handling, and command preflight. |
| Observability | Applies to verification status, stage/run metadata, diagnostics, evidence links, staleness, and next-action visibility across text, JSON, and TUI. |
| Security | Applies to explicit executable/argv invocation, bounded environment forwarding, path containment, redaction, read-only verification, and prohibited product or Git mutation. |
| Testing | Applies to fake review runtimes, fake smoke harnesses, deterministic state and parity tests, race testing, and gated real-runtime evidence. |
| Documentation | Applies to readable and valid `review.md` and `smoke.md` summaries, recovery guidance, evidence links, command help, and stable artifact conventions. |
| CLI Surface | Applies to `verify`, `flow --to smoke`, status, focused rerun controls, diagnostic override handling, exit codes, and text/JSON parity. |
| LLM Runtime | Applies to the bounded agentwrap-backed review stage while preserving the generic runtime boundary and normal-test fake seam. |
| Workflows | Applies to ordered `execute -> review -> smoke` orchestration, gating, cancellation, resumability, reruns, stale propagation, and recovery. |
| Persistence And Migrations | Applies to versioned strict state formats, atomic writes, freshness metadata, and preservation of the last complete review and smoke artifacts. |

## Selected Evidence Reports

Copied from the project index's "Available Evidence Reports" table. These tell the technical handbook which reports to read – the project index is the authoritative source.

| Report | Path | Covers |
| --- | --- | --- |
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
| `15-philosophy` | `studies/go-cli-study/reports/final/15-philosophy.md` | Cross-cutting design philosophy, tradeoffs, maintainability |

## Selected Reasoning Templates

All paths must appear in the project index's "Available Reasoning Templates" table.

| Template | Output Path | Why Selected |
| --- | --- | --- |
| Architecture | `projects/ultraplan-go/sprints/28-review-to-smoke-flow/reasoning/architecture.md` | Reason through module ownership, dependency direction, shared CLI/TUI use cases, generic runtime/process boundaries, and durable verification state. |

## Prior Decisions To Carry Forward

No prior decisions are cataloged in the project index's "Prior Decisions" table. Sprint 23, Sprint 26, and Sprint 27 outputs remain dependencies through the sprint requirements, but no decision path is selectable here.

## Required Review Protocols

All paths must appear in the project index's "Review Protocols" table.

| Protocol | Path | Required Evidence |
| --- | --- | --- |
| Architecture Review | `system/protocols/architecture-review-protocol.md` | Evidence that verification semantics remain in `internal/sprint`, CLI and TUI share typed app use cases, and platform runtime/process packages remain product-neutral. |
| Sprint Review | `system/protocols/review-sprint-protocol.md` | Current `review.md` evidence covering selected contracts, handbook guidance, sprint decisions, plan execution, tests, findings, and deterministic verdict. |
| Deep Smoke Sprint | `system/protocols/deep-smoke-sprint-protocol.md` | Current `smoke.md` linked to matching cataloged-harness run and issue evidence, or a truthful blocked/not-applicable result. |

## Excluded Context

| Context | Reason Excluded | Revisit If |
| --- | --- | --- |
| Open-ended smoke investigation and harness maintenance | The sprint integrates product-owned smoke orchestration, focused reruns, evidence validation, and recovery; unrestricted investigation plus harness test/issue maintenance remains a separate action. | Requirements expand from bounded verification and focused reruns to maintaining or diagnosing the external harness itself. |
| General-purpose issue tracking | Relevant harness issues affect smoke verdicts and are linked, but assignment, scheduling, synchronization, and project-management behavior are non-goals. | A later project phase explicitly introduces issue-management requirements and catalog context. |
| Git mutation | Review, smoke, verify, status, freshness checks, and recovery must not add, commit, push, branch, merge, reset, checkout, or otherwise mutate Git state. | A future sprint explicitly scopes governed Git operations under newly selected contracts and requirements. |
| Automatic implementation fixes | Verification reports failures and next actions but must not automatically edit product source, product tests, review findings, smoke failures, or harness issues. | A future sprint explicitly defines a governed remediation workflow with mutation boundaries. |
| Third assessment artifact | Overall assessment must be derived deterministically from current review, smoke, flow state, and harness issues without another sprint-root artifact. | Product requirements intentionally add and govern another canonical assessment artifact. |
| Raw harness evidence in the sprint | Detailed run JSON, unrestricted stdout/stderr, per-test artifacts, and issue files remain under the cataloged external harness root. | The external harness contract or project artifact policy is explicitly revised. |
| Hosted, browser, and multi-user surfaces | This sprint is limited to CLI, JSON, and local TUI parity for the Phase 3 verification flow. | A later phase brings a local browser or hosted collaboration surface into scope. |
| Cross-project or cross-sprint verification scheduling | The sprint operates on one selected project and sprint and does not introduce a general scheduler. | A later sprint explicitly requires coordinated verification across projects or sprints. |

## Next Artifacts

- `technical-handbook.md` reads from the evidence reports listed above.
- `reasoning/*.md` captures area-specific reasoning.
- `reasoning.md` makes final sprint decisions.
- `plan.md` executes `reasoning.md`.
- `review.md` runs the selected review protocols against implementation.
