# Execute Summary

> Project: `ultraplan-go`
> Sprint: `34-shared-context`
> Status: `implementation complete; ready for governed review`
> Execution mode: `manual in-session execution under $ultraplan-execute`
> Executed: `2026-08-20`

Plan: `projects/ultraplan-go/sprints/34-shared-context/plan.md`
Run state: `projects/ultraplan-go/sprints/34-shared-context/.run-state.json`

## Target And Scope

- Implementation repository: `/home/antonioborgerees/coding/ultraplan/ultraplan-go`
- Recorded base revision: `131c842595e62d503fadd431d15af394c21b9e3c`
- Review scope: the current uncommitted Sprint 34 implementation diff and the two untracked shared-context files listed below.
- Git state: no Git mutation was performed.
- Preserved unrelated work: pre-existing edits in `internal/web/static/app.css`, `internal/web/static/app.js`, `internal/web/templates/sprint.html`, and `internal/web/templates_test.go` were not modified as part of Sprint 34.

## Task Counts

- pending: 0
- running: 0
- complete: 8
- deferred: 0
- failed: 0
- cancelled: 0

## Implementation Result

- Added one concrete `internal/sprint` renderer that preserves exact requirements/code-context bytes, renders contained source ranges sequentially in authored order, and ends the byte-stable prefix with one fixed stage boundary.
- Added fail-closed relative-path, symlink, regular-file, identity/replacement, line-range, UTF-8, cancellation, 256 KiB prefix, and 64 KiB suffix-reserve enforcement without persistence, recursion, indexing, retrieval, or caching.
- Wired the shared prefix explicitly into sprint-index, technical-handbook, area/final reasoning, plan, execute, every independent review request and repair, and agent-backed smoke authoring. Review fan-out and execute task runs reuse one immutable operation prefix.
- Kept `internal/platform/runtime`, agentwrap, flow/state authority, review verdict ownership, execute task/session semantics, smoke mutation scope, and adapter ownership unchanged.
- Audited the existing cumulative flow and added an exact canonical-order assertion; code-context remains scheduled once after requirements.
- Added the manual-only `ultraplan-code-context` stage skill with exact canonical CLI delegation and retained dry-run, customization-preservation, force-restoration, isolation, and synchronized metadata behavior.
- Updated every required release guide with the reference-only artifact model, exact shared prefix, transient/untrusted source evidence, live-inspection permission, limits, recovery, browser behavior, canonical skill, and truthful dogfood rules.

## Changed Implementation Paths

```text
README.md
docs/architecture.md
docs/cli-reference.md
docs/local-web.md
docs/planning-smoke.md
docs/recovery.md
docs/stage-skills.md
docs/user-guide.md
internal/sprint/execute.go
internal/sprint/prompt_context.go
internal/sprint/prompt_context_test.go
internal/sprint/prompts.go
internal/sprint/review.go
internal/sprint/review_test.go
internal/sprint/service.go
internal/sprint/smoke_author.go
internal/sprint/smoke_test.go
internal/sprint/sprint_index_test.go
internal/workspace/scaffold/templates/README.md
internal/workspace/skills.go
internal/workspace/skills_test.go
```

`internal/sprint/prompt_context.go` and `internal/sprint/prompt_context_test.go` are untracked implementation files and must be included by review rather than relying only on tracked `git diff` output.

## Verification Transcript

All commands were run manually from `/home/antonioborgerees/coding/ultraplan/ultraplan-go` against the final implementation tree.

| Command | Result |
| --- | --- |
| `go test ./internal/sprint` | pass |
| `go test ./internal/workspace ./internal/app ./internal/tui ./internal/web` | pass |
| `go test ./...` | pass |
| `go test -race ./...` | pass |
| `go vet ./...` | pass |
| `go build ./cmd/ultraplan` | pass |
| `git diff --check` | pass |

Focused evidence includes exact LF/CRLF/no-final-newline artifact slices, exact 256 KiB fit and one-byte overflow, authored duplicates/overlaps, deterministic source replacement detection, no recursion/mutation, cross-stage previews, real service runtime requests for every planning stage plus execute, independent reviewer fan-out prefix equality, smoke-author prefix placement, canonical flow order, and manual-skill materialization behavior.

## Review Preflight

- Manual Architecture Review preflight: **Approve with comments**. Ownership, dependency direction, explicit call-site integration, generic runtime separation, immutable operation reuse, containment, state neutrality, errors, test coverage, bounded performance, and exclusions are sound. Portable hard-link ancestry and provider-specific capacity remain documented residual risks.
- Manual Sprint Review preflight: **Pass with gated evidence blocked**. The implementation and deterministic evidence cover all five decisions and selected contracts; no blocker/high implementation finding was identified.
- These are execution preflight results only. No canonical `review.md` is claimed or created; governed reviewer fan-out and product-owned verdict synthesis remain the next stage.

## Gated Dogfood

Status: **blocked**.

Exact prerequisite: the active manual execute skill permits UltraPlan CLI use only for project status, sprint status, and execute-prompt resolution. The documented real-runtime dogfood requires `flow --to plan`, which is outside that granted execute-stage command permission. The flow was therefore not invoked. No fake-runtime result, constructed/unsent request, or artifact-only observation is represented as real-runtime success.

Rerun the disposable procedure in `docs/planning-smoke.md` under an explicitly authorized runtime/dogfood gate before claiming real-provider success or assessing provider-specific prompt capacity.

## Deviations, Residual Risks, And Next Stage

- No implementation deviation from `reasoning.md` is recorded.
- Portable hard-link ancestry cannot be proven by the chosen cross-platform containment checks; reads remain contained, regular-file-only, bounded, read-only, and explicitly untrusted.
- The 256 KiB prefix plus 64 KiB suffix product bound is tested exactly but not proven against the configured real provider because dogfood is blocked.
- No cache, retrieval/index, staleness, provenance, QA/repair, alternate persistence, hosted, issue, or automatic Git-mutation capability was added.
- Downstream dependency to revisit: `projects/ultraplan-go/sprints/35-durable-run-observability/requirements.md:96` currently says Sprint 34 supplied real-runtime dogfood evidence. That statement is stale because the governed outcome here is `blocked`; revise Sprint 35 requirements and any derived reasoning/plan before executing that sprint.
- Next stage: run the governed Sprint 34 review so it freezes the current implementation and untracked-file scope, consumes this evidence and exact blocker, and creates `review.md` atomically. Smoke remains downstream of an acceptable canonical review.
