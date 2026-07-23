# Sprint 28 Execute Summary

Execution status: `complete`
Date: `2026-07-22`
Target: `/home/antonioborgerees/coding/ultraplan-go`

## Implementation Summary

Sprint 28 composes the existing review and smoke stages behind one product-owned `Service.Verify` transition. It adds strict flow-state schema version 2 with an atomic version-1 migration, active/last attempts, last-complete evidence, durable per-coverage review resume checkpoints, separate input and artifact identities, deterministic freshness and assessment, focused review merge rules, non-promoting diagnostic smoke, typed conflicts/recovery, stable CLI JSON, and shared TUI operations.

The existing generic agentwrap runtime and explicit-argv process runner remain product-neutral boundaries. The runtime adapter now retains the latest terminal model output in a separately bounded 96 KiB in-memory field for machine consumers while continuing to cap mapped/display event strings at 4 KiB. Review and smoke implementations were extended rather than replaced.

## Changed Implementation Files

- `cmd/ultraplan/main.go`
- `internal/app/app.go`
- `internal/app/app_test.go`
- `internal/platform/runtime/agentwrap.go`
- `internal/platform/runtime/runtime.go`
- `internal/platform/runtime/runtime_test.go`
- `internal/sprint/domain.go`
- `internal/sprint/state.go`
- `internal/sprint/service.go`
- `internal/sprint/flow.go`
- `internal/sprint/review.go`
- `internal/sprint/smoke.go`
- `internal/sprint/smoke_protocol.go`
- `internal/sprint/smoke_types.go`
- `internal/sprint/verify.go`
- `internal/sprint/review_test.go`
- `internal/sprint/smoke_test.go`
- `internal/sprint/verify_test.go`
- `internal/app/sprint_commands.go`
- `internal/app/sprint_commands_test.go`
- `internal/app/sprint_verify_commands_test.go`
- `internal/app/sprint_usecases.go`
- `internal/app/operations.go`
- `internal/app/tui_commands.go`
- `internal/app/tui_commands_test.go`
- `internal/tui/model.go`
- `internal/tui/views.go`
- `internal/tui/verify_test.go`

## Task Evidence

1. Existing Sprint 26/27 seams were confirmed: `internal/sprint` already owned review/smoke policy, runtime and process seams were injectable, and flow dispatch was the minimal convergence point. The predecessor was confirmed as `flow-state.json` schema version 1. The cataloged harness is machine manifest protocol v1 with explicit executable/argv, discovery/run, evidence roots, scope mappings, prerequisites, and issue records.
2. `VerificationStatus`, `VerificationStage`, `VerificationAttempt`, `DiagnosticOverride`, `EvidenceReference`, and `OverallAssessment` expose one typed semantic projection. `deriveAssessment` implements the fixed precedence table.
3. State schema v2 persists active/last attempts and last-complete review/smoke facts. Version 1 migrates atomically and conservatively; zero, unknown, contradictory, malformed, and trailing JSON fail closed. Per-service sprint mutation exclusion returns `ErrVerificationConflict`.
4. Review manifests now include project inputs, all governed sprint artifacts, selected catalog content, area reasoning, execute evidence, explicit missing markers, changed scope, and deterministic Git/non-Git target identity. Smoke identity covers current review, harness manifest/protocol, scope, prerequisites, effective timeout source, target, external evidence, and issue files. Canonical Markdown has separate digests.
   Reviewer prompts carry the frozen logical path and digest index instead of embedding every input in one argv value. Each attempt reads a private temporary snapshot containing only captured manifest bytes, with external-directory access denied; prompts remain below a 96 KiB fail-closed budget and post-run fingerprint reconciliation still rejects drift.
5. Both `flow --to review|smoke` and `verify --to review|smoke` call `Service.Verify`, which requires complete execute evidence, obtains a current review, applies the gate, and returns the earliest next action.
6. Focused review reruns retain only complete same-fingerprint coverage and promote a single fully validated artifact. Failed/cancelled/malformed/drift attempts preserve the prior canonical artifact and facts.
7. Existing smoke discovery, scope containment, explicit argv, bounded output, issue validation, and evidence validation remain in place. Narrow scopes and review overrides are diagnostic and cannot replace canonical containing-suite evidence.
8. CLI status, flow, review, smoke, and verify expose the shared projection in text/JSON. Focus, scope, timeout, confirmation, and rationale are parsed explicitly; progress remains on stderr.
9. TUI actions call app use cases for status, flow, review, smoke, and verify; confirmation shows scope and diagnostic rationale. Views show stage state, freshness, evidence, issues, assessment, and recovery.
10. Tests cover assessment precedence, migration/rejection, target/artifact/external-evidence freshness, focused merge, conflicts, last-complete preservation, scope containment, issue verdicts, CLI parsing/help/parity, and TUI actions/rendering. Boundary searches found no sprint/review/smoke semantics in platform runtime/process and no CLI/process execution path in TUI.
11. The first automated review feedback was remediated before retry: complete flow scheduling moved into `Service.Flow`; TUI/runtime factories are explicit composition dependencies; review citations are frozen-manifest-only; terminal state write errors propagate; test scopes are always diagnostic; manifest timeouts fail closed; status models expected smoke unavailability as partial; JSON failures are singular and structured; and stable smoke argv retains no argument values.
12. The operator-selected MiniMax M3 reviewer returned structured JSON larger than the generic runtime adapter's 4 KiB mapped-event limit. The adapter now preserves a separately bounded 96 KiB terminal output for machine parsing without expanding event payloads sent to progress, logs, or TUI surfaces; regression tests prove both bounds.
13. Review attempts now checkpoint each coverage item's state, validated structured result, and latest OpenCode session ID in `flow-state.json`. A compatible retry reuses completed coverage and invokes only unfinished reviewers with `SessionActionContinue`; a changed fingerprint/model starts fresh. `review --restart`, `verify --restart-review`, flow, and TUI restart operations explicitly discard resumable state and use fresh sessions. The current agentwrap contract already maps continuation to OpenCode's `--session`, so no agentwrap source change was required.

## Verification Evidence

All commands ran from `/home/antonioborgerees/coding/ultraplan-go` on 2026-07-22.

| Command | Result |
| --- | --- |
| `go test ./internal/sprint` | pass |
| `go test ./internal/app` | pass |
| `go test ./internal/tui` | pass |
| `go test ./...` | pass |
| `go test -race ./...` | pass |
| `go build ./cmd/ultraplan` | pass; generated root binary removed after the check |
| `go vet ./...` | pass |
| `git diff --check` | pass |
| `go test ./internal/sprint -run 'TestReviewManifestExecutionAndArtifactPreservation\|TestReviewerPromptStaysBelowSubprocessArgumentBudget\|TestReviewFanOutUsesConfiguredBound'` | pass |
| `go test ./internal/sprint ./internal/app ./internal/tui` after review remediation | pass |
| `go test ./...` after review remediation | pass |
| `go test -race ./...` after review remediation | pass |
| `go vet ./...` after review remediation | pass |
| `go build -o /tmp/ultraplan-sprint28 ./cmd/ultraplan` after review remediation | pass; temporary binary removed |
| `go test ./internal/platform/runtime ./internal/sprint` after MiniMax M3 compatibility fix | pass |
| `go test ./...` after MiniMax M3 compatibility fix | pass |
| `go test -race ./...` after MiniMax M3 compatibility fix | pass |
| `go vet ./...` after MiniMax M3 compatibility fix | pass |
| `go test ./internal/sprint ./internal/app ./internal/tui` after durable review resume | pass |
| `go test ./...` after durable review resume | pass |
| `go test -race ./...` after durable review resume | pass |
| `go vet ./...` after durable review resume | pass |
| `go build -o /tmp/ultraplan-resume-review ./cmd/ultraplan` after durable review resume | pass |
| `git diff --check` after durable review resume | pass |
| `rg -n 'sprint|review|smoke|harness|verdict' internal/platform/runtime internal/platform/process` | no product-semantic matches |
| `rg -n 'runSprint|runForTest|cmd/ultraplan|os/exec' internal/tui` | no CLI/self-process matches |

Normal tests used fake review runtimes and fake process/harness seams; they required no provider credentials, network, OpenCode process, or real smoke harness.

## Safety And Mutation Evidence

- Review requests retain read-only sandbox/permission policy and caller context. Frozen input copies use a private owned temporary root, expose only captured manifest bytes, and are removed after joined reviewer work.
- Smoke retains direct executable plus argv invocation, contained cwd/evidence roots, allowlisted environment, bounded capture, timeout, cancellation, and cleanup reporting.
- Target identities are read-only. No implementation path invokes Git mutation.
- Failed attempts preserve the prior canonical `review.md`/`smoke.md`; artifact/state disagreement becomes stale or malformed rather than current.
- Interrupted review checkpoints contain only bounded structured results and opaque session IDs. Session IDs are persisted on first observation and whenever agentwrap reports a newer retry session; no raw provider stream is stored.
- Stable state and JSON exclude raw provider and harness streams. Diagnostics are bounded and newline/NUL checked.
- Git status was inspected before and after implementation; only the listed implementation files and owned Sprint 28 artifacts changed.

## Deviations And Deferrals

No reasoning or architecture deviation was required. Automated review findings were resolved in the planned sprint/app composition boundaries rather than waived.

Real agent-backed review and cataloged harness smoke are protocol gates rather than normal-test requirements. They are run after this execute record and reported in `review.md` and `smoke.md`. If the harness lacks a Sprint 28 complete mapping or prerequisites, `smoke.md` records a truthful `blocked` result; no scope is inferred from README prose.

The first governed review attempt exposed an argv transport limit: health and a compact direct OpenCode run succeeded, while the original all-content reviewer prompt exceeded the operating system's single-argument limit and every subprocess start failed. The compact frozen path/digest index and 96 KiB regression test resolve that implementation defect before the canonical review retry.

The first `minimax-coding-plan/MiniMax-M3` fan-out then exposed a separate generic-adapter limit: the model returned terminal review JSON larger than the intentionally bounded 4 KiB event payload, so aggregation saw truncated JSON. A direct M3 JSON probe confirmed provider/model correctness. The adapter now gives machine consumers a 96 KiB in-memory terminal result while leaving displayed/runtime event payloads capped and raw payload persistence disabled.

## Recovery Notes

- Unknown flow-state versions: restore a supported v1/v2 file or regenerate state explicitly.
- Stale or edited review: rerun review with current governed inputs.
- Stale or edited smoke/evidence/issues: rerun the required containing suite.
- Active mutation conflict: wait for or cancel the active attempt, then retry.
- Interrupted review: rerun the same review/verify/flow command to reuse validated coverage and retained OpenCode sessions. Use `review --restart` or `--restart-review` on verify/flow to discard the checkpoint and start all reviewers fresh.
- Diagnostic override: provide explicit confirmation and rationale; the result remains non-promoting.
