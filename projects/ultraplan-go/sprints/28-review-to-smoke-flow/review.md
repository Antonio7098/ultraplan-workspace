# Sprint 28 Review: Integrated Review-to-Smoke Verification Flow

> Project: `ultraplan-go`
> Sprint: `28-review-to-smoke-flow`
> Date: `2026-07-28`
> Reviewer: `sprint-review-protocol orchestrator (mini Max-M3 / 12 reviewer subagents)`
> Verdict: `fail`

## Review Context

This file is the canonical Sprint 28 conformance review produced by the `review-sprint-protocol` defined at `system/protocols/review-sprint-protocol.md`. It evaluates the integrated review-to-smoke verification flow implemented in `internal/sprint`, `internal/app`, `internal/tui`, and the supporting `internal/platform/runtime` and `internal/platform/process` packages. It is the replacement for the older manually produced review file and supersedes the failed agentwrap-driven review attempt recorded in `flow-state.json` (`review.verdict: blocked`). The previous attempt produced 7 of 12 valid structured reviewer results (cli-surface, configuration, documentation, errors, llm-runtime, security, testing) and failed schema/citation validation or OpenCode runtime errors for 5 (architecture, observability, persistence-and-migrations, workflows, handbook). This run executes the protocol's "Run one structured reviewer per selected contract and one for the technical handbook" rule via 12 bounded local-subagent reviewers, none of which depend on the agentwrap runtime.

## Input Fingerprint And Scope

### Target implementation

- Path: `/home/antonioborgerees/coding/ultraplan-go`
- Git revision: `f4d6d3848a89c7f3304794d5e6328c3e582e4dd3` (clean working tree; no uncommitted mutations)
- Last commit: `feat: add resumable review to smoke verification`
- Scope: full changed-path manifest from `git diff HEAD~1..HEAD` (28 files, ~3,700 lines delta), enumerated in `execute.md` lines 14-42
- Evidence of non-mutation: `git status` clean before/after review; `execute.md` lines 96-105 document prior mutation-boundary checks

### Governed sprint input digests

| File | SHA-256 |
| --- | --- |
| `projects/ultraplan-go/sprints/28-review-to-smoke-flow/requirements.md` | `8a242345ed0f07c83aad639fbcb0f5c1faee2c3e945b35bb3f515c330514f08c` |
| `projects/ultraplan-go/sprints/28-review-to-smoke-flow/sprint-index.md` | `a558bceb16fc038fb44e97f1d9adad2962f5e370a1c2f91563bec13ba1c4ce36` |
| `projects/ultraplan-go/sprints/28-review-to-smoke-flow/technical-handbook.md` | `844459a852b86430c4c862d74a60491d43072523499705dc4189ddb9c7fda7d9` |
| `projects/ultraplan-go/sprints/28-review-to-smoke-flow/reasoning.md` | `d79ab559db66600560b097c504a8255638bf0a56e90de6b8d966dd9ebe5e536f` |
| `projects/ultraplan-go/sprints/28-review-to-smoke-flow/reasoning/architecture.md` | `5ae036ec43177170ebbe2c6b120b426b1120b311b89cbc92587c94bfaad65d0e` |
| `projects/ultraplan-go/sprints/28-review-to-smoke-flow/plan.md` | `b7cf1f1d83931d9fbf805260d02d8d1007dd26bdf25a270d3e79d985b0e39718` |
| `projects/ultraplan-go/sprints/28-review-to-smoke-flow/execute.md` | `ee8eb86d87e4e4223800504dff4b78b76991c88e52f4da4e81ed2457e80a37b8` |
| `projects/ultraplan-go/sprints/28-review-to-smoke-flow/.run-state.json` `planFingerprint` | `b7cf1f1d83931d9fbf805260d02d8d1007dd26bdf25a270d3e79d985b0e39718` (matches `plan.md`) |

### Selected contracts and handbook

| Contract / handbook | Path | SHA-256 |
| --- | --- | --- |
| Architecture | `system/contracts/core/architecture.md` | `275e80f72b0127ef1ee92c1f76b02c0e6f2084adaaf97315bdb61ae42c8a7556` |
| Errors | `system/contracts/core/errors.md` | `d802fe296d19932c6e0115625a4df0e5caa987a8b96ccd32fcd5849ced9af9e2` |
| Configuration | `system/contracts/core/configuration.md` | `40451fce9c8d33ec214aec8c70d524d2b3c0524a2f23655fa81b49263c634d66` |
| Observability | `system/contracts/core/observability.md` | `9a48650b90f273955751ef7e833e80436be64e36d6803a6cadc2cf76d6d962fb` |
| Security | `system/contracts/core/security.md` | `5af338526c1f0898631c06ca76d651aaa59b391fa5b5d341c899086ff3d3d328` |
| Testing | `system/contracts/core/testing.md` | `25905e8c5826e9068e4c752cd9ea80997e85cb7a6c53ae1ba72f9a8b78eb42fd` |
| Documentation | `system/contracts/core/documentation.md` | `4624a2876d574313a816d018c68300eae180cda6c526ade0ee5811dbdeeb767d` |
| CLI Surface | `system/contracts/surfaces/cli.md` | `1d9023f483e128960de402e4402b6a1eb116cb542f4617910e42c96d9569f243` |
| LLM Runtime | `system/contracts/runtime/llm.md` | `1d81f5a844f9088ac4c1b43305300449bad6b01e83e131ff5964de4f844ccce2` |
| Workflows | `system/contracts/runtime/workflows.md` | `5d0200141beb53f8330a508ec35b90571b546ed56d2419e77b9e1b4d96048b01` |
| Persistence And Migrations | `system/contracts/runtime/persistence-and-migrations.md` | `249990b954b05e2274bd9c04febf92184d69f4380386626582adeaaee54974a3` |
| Sprint Review Protocol | `system/protocols/review-sprint-protocol.md` | `2e221ec3cfc3f7fe61b16020bfa1f1f5767a54acff387f846d89f1ebabc2ec2d` |

### Reviewer count, model source, concurrency, write boundary

- Reviewer count: 12 (11 contracts + 1 technical handbook)
- Model source: local subagent reviewers (no external runtime dependency)
- Concurrency: parallel dispatch in three batches of (5, 5, 2)
- Write boundary: atomic write of `projects/ultraplan-go/sprints/28-review-to-smoke-flow/review.md` only; no product, test, governed-input, or Git mutation

## Decision Conformance

The seven final decisions in `reasoning.md` (lines 109-207) are evaluated against the implementation:

| Decision | Conformance | Finding references |
| --- | --- | --- |
| 1. Keep verification policy product-owned | **pass** | `ARCH-INFO-001`..`ARCH-INFO-006` positive; `HB-001-001` positive; TUI/app remain thin adapters |
| 2. Version facts and derive current truth | **partial** | `PERSIST-MIG-001` (high), `PERSIST-INTEGRITY-001` (high), `PERSIST-ATOMIC-002` (medium), `PERSIST-ATOMIC-003` (low), `PERSIST-READ-002` (low), `PERSIST-INTEGRITY-002` (low) |
| 3. Fingerprint deterministic stage manifests | **pass** | All reviewers confirm manifest covers and stale propagation; `OBS-CORR-002` (high) is about request correlation, not manifest identity |
| 4. Use one ordered gate and deterministic assessment | **partial** | `WF-RETRY-001-Override-Blocked` (medium) flags `prepareSmokeStatic` accepting override for `ReviewVerdictBlocked`; `ERR-CORE-001` (medium) flags stale/incomplete collapse to `ExitValidation` only |
| 5. Promote focused reruns only with complete coverage | **partial** | `WF-IDEMPOTENCY-001` (medium) flags lack of idempotency key; merge tests pass per `TEST-003` |
| 6. Preserve evidence through atomic, cancellable attempts | **partial** | `PERSIST-ATOMIC-002` (medium) flags `LoadFlowState` hidden commit; `WF-COMP-001` (low) flags smoke.md/state write-order window; `ERR-TASK-001` (high) flags `attemptExpired` dead code |
| 7. Freeze safe execution and expose one semantic result | **partial** | `ERR-REDACT-001` (blocker) for SDK response redaction; `LLM-PROMPT-001-violation` (high) for prompt identity; `LLM-OBS-001-violation` (medium) and `LLM-RETRY-001-partial` (medium); `OBS-CORR-002` (high) for top-level correlation; `CLI-HELP-001` (low) and `CLI-EXIT-001` (low) for help completeness |

## Plan Execution

All 10 plan tasks (`plan.md` lines 100-180) are marked `complete` in `.run-state.json` (lines 16-25). Verification evidence recorded in `execute.md` (lines 62-94) shows:

| Command | Result |
| --- | --- |
| `go test ./internal/sprint` | pass |
| `go test ./internal/app` | pass |
| `go test ./internal/tui` | pass |
| `go test ./...` | pass |
| `go test -race ./...` | pass |
| `go build ./cmd/ultraplan` | pass |
| `go vet ./...` | pass |
| `git diff --check` | pass |
| focused/narrow reruns after review remediation and durable resume | pass |

The plan-execution gate is satisfied. The deferred items in `execute.md` (lines 108-114) are correctly recorded as gate-conditional real-runtime evidence, not plan gaps.

## Verification Evidence

| Evidence class | Source | Quality |
| --- | --- | --- |
| Domain transition tests | `internal/sprint/verify_test.go` | Table-driven across all 12 precedence rows; covers migration, focused merge, override non-promotion, mutation conflict |
| App/CLI parity tests | `internal/app/sprint_verify_commands_test.go`, `sprint_commands_test.go` | Negative parse paths, exit matrix, JSON/text discipline, help gating |
| TUI tests | `internal/tui/verify_test.go`, `tui_commands_test.go` | Use-case delegation, render summary; narrow per `HB-010-010` |
| Boundary/security tests | `internal/platform/runtime/runtime_test.go`, `internal/sprint/smoke_test.go` | Policy/mapping/event retention, fail-closed defaults |
| Mutation snapshots | `execute.md` lines 96-105 + `git status` clean | No product, test, governed-input, or Git mutation |
| Real smoke | `execute.md` line 110 records `blocked` is the truthful result for Sprint 28 once the gate is exercised; `smoke.md` exists as the canonical artifact |

Boundary searches (re-checked this run) confirm `internal/platform/runtime` and `internal/platform/process` contain no `sprint|review|smoke|harness|verdict` references, and `internal/tui` contains no `os/exec`/`cmd/ultraplan` invocations or terminal-output parsers.

## Contract Conformance

The 11 selected contracts are evaluated by one bounded reviewer each. The complete structured findings are appended in the Findings section. Coverage summary:

| Contract | Applicability | Result |
| --- | --- | --- |
| Architecture | direct | pass (6 positive confirmations; 2 low nits) |
| Errors | partial | fail (1 blocker; 3 high; 2 medium; 2 positive) |
| Configuration | direct | pass for previously-blocking items; 1 medium secret finding; 1 low env-mode + 1 low redaction completeness |
| Observability | direct | partial (1 high correlation; 3 medium; 1 low) |
| Security | direct | pass (8 positive; 1 low partial for override gate) |
| Testing | direct | pass (2 low partial gaps and 1 low TUI coverage) |
| Documentation | direct | partial (2 high stale-docs; 1 medium recovery; 1 low TUI help; 4 positive) |
| CLI Surface | direct | pass on most rules; 2 low gaps (examples and exit-code table) |
| LLM Runtime | direct | partial (1 high prompt identity; 2 medium; 6 positive) |
| Workflows | direct | partial (3 high; 1 medium; 1 low; 1 info) |
| Persistence And Migrations | direct | fail (2 high; 1 medium; 3 low) |

## Technical Handbook Conformance

The technical handbook reviewer (`HB-001`..`HB-012`) confirms positive conformance to: thin surfaces over shared operations, explicit inspectable lifecycle state, typed failures with separate recovery rendering, narrow injected runtime/process/IO seams, bounded joined reviewer fan-out, atomic state with one migration, configuration precedence and redaction, signal/cancellation propagation, and behavior-focused fake-backed tests. Three low observations remain: a thin TUI verification test surface (`HB-010-010`), a bounded non-deterministic join on the runtime events goroutine (`HB-011-011`), and the sprint-owned `targetIdentity` git invocation that bypasses the platform seam (`HB-012-012`).

## Applicability And Deferred Scope

| Item | Status | Source |
| --- | --- | --- |
| Sprint 26 review and Sprint 27 smoke replacement | `not_triggered` | `reasoning.md` Decision 1 rejects wholesale replacement |
| Browser / hosted / multi-user surfaces | `not_triggered` | `sprint-index.md` excludes Phase 4; `p hase 4` begins at Sprint 30 |
| Cross-sprint scheduling | `not_triggered` | out of scope per requirements |
| Real OpenCode runtime evidence | `deferred` to `blocked` per `smoke.md` | gated by `ULTRAPLAN_REAL_SMOKE=1` |
| General issue tracking | `not_triggered` | scope explicit; harness issues are evidence only |
| Automatic product / test / finding / harness mutation | `explicitly_deferred` | prohibited by constraint |
| Browser-side verification | `not_triggered` | Phase 4 |

## Findings

Findings are grouped by their canonical contract. Each finding adopts the severity returned by the responsible reviewer. Citation paths use the prior agentwrap convention `target/internal/...` when the reviewer cited that form, and the absolute path when the reviewer preferred it. All paths resolve into the target implementation tree.

### Architecture (`ARCH`)

| ID | Sev | Title | Action |
| --- | --- | --- | --- |
| ARCH-INFO-001..006 | info | Policy ownership, typed projection, thin adapters, product-neutral platform, registrar composition, narrow ports | keep boundaries |
| ARCH-LOW-001 | low | Bounded `git` subprocess for target-identity fingerprint | document or move behind a narrow port |
| ARCH-LOW-002 | low | CLI error classification uses substring matching | switch to typed sentinels |

### Errors (`ERR`)

| ID | Sev | Title | Action |
| --- | --- | --- | --- |
| ERR-REDACT-001 | blocker | SDK response redaction depends on truncation only; secrets in `UserDetail`/`DebugDetail`/`ResponseBody` leak | apply `config.RedactValue` to user/debug/response fields before returning; add regression |
| ERR-CODE-001 | high | Operation error codes derived from substring matching | replace with typed sentinels |
| ERR-TASK-001 | high | `attemptExpired` is dead code; stale-running attempts are never reconciled | call from `VerificationStatus` or remove |
| ERR-RETRY-001 | high | `OperationError.Retryable` is an orphan flag without a retry policy | remove or back with bounded retry owner |
| ERR-TRANS-001 | medium | `mapSmokeError` stringifies the typed smoke error | wrap with `%w` so `errors.As` survives |
| ERR-CORE-001 | medium | `stale`/`incomplete` overall assessments lack distinct exit codes | map `AssessmentIncomplete` to a distinct `Exit*` |
| ERR-SHAPE-001 | info | Canonical fields present (remediated) | no action |
| ERR-USER-001 | info | User/operator surfaces separated via `displaySafe` and `config.RedactValue` | no action |

### Configuration (`CFG`)

| ID | Sev | Title | Action |
| --- | --- | --- | --- |
| CFG-START-001 | info | `internal/platform/config` package now present and wired (remediated) | mark blocker resolved |
| CFG-PUBLIC-001 | info | `Redact()` allowlist implemented and tested (remediated) | mark blocker resolved |
| CFG-SOURCE-001 | info | Precedence enforced and tracked per field | no action |
| CFG-TYPE-001 | info | Post-merge typed parsing and validation | no action |
| CFG-COMPAT-001 | info | Flow-state v1→v2 migration strict | no action |
| CFG-START-001-CFG | info | Invalid config fails startup | no action |
| CFG-OBS-001 | info | Effective config inspectable via `config show` | no action |
| CFG-SECRET-001 | medium | Default smoke env still includes `HOME`/`PATH`; allowlist redaction has gaps | tighten default; redact `ExtraArgs`; warn on sensitive env names |
| CFG-ENV-001 | low | No explicit environment/mode axis | document or add `environment` field |
| CFG-OBS-001-REDACT-COMPLETENESS | low | `Redact()` omits some sensitive strings and smoke env is unredacted | defensively redact `ExtraArgs` and warn on `HOME`/`PATH` |

### Observability (`OBS`)

| ID | Sev | Title | Action |
| --- | --- | --- | --- |
| OBS-CORR-002 | high | Top-level command/request correlation identifier missing | generate UUID at CLI entry; propagate through `Runtime`/event/JSON |
| OBS-HEALTH-002 | medium | Readiness is a single bool; no liveness/readiness/degraded separation | model readiness as enum/struct; add preflight command |
| OBS-CORE-003 | medium | Runtime `Result.Warnings` (event drops, permission, repair) not surfaced in operator status | persist key warnings into `ReviewStageState`/`SmokeStageState`; render |
| OBS-DEBUG-002 | medium | No `--verbose` flag exposed for sprint verify/flow/review/smoke | add `--verbose` with documented redaction |
| OBS-TASK-002 | low | Harness cleanup status is captured but not exposed as structured terminal state | add structured `Cleanup` outcome to attempt diagnostics |

### Security (`SEC`)

| ID | Sev | Title | Action |
| --- | --- | --- | --- |
| SEC-INJECT-001 | info | Explicit argv + env allowlist | no action |
| SEC-INPUT-001 | info | CLI/manifest/citation boundary validation | no action |
| SEC-SECRETS-001 | info | Secret-bearing values bounded/redacted/excluded | no action |
| SEC-FILES-001 | info | Path containment | no action |
| SEC-AUTHZ-001 | low | Diagnostic override is a deliberate authorization gate | no action |
| SEC-NET-001 | info | External calls bounded timeouts, opaque session IDs | no action |
| SEC-DESER-001 | info | Strict JSON decoding | no action |
| SEC-DEFAULT-001 | info | Security-sensitive defaults fail closed | no action |
| SEC-DEPS-001 | info | No new dependencies | no action |
| SEC-AUTHN-001 | n/t | AuthN boundary not introduced | no action |

### Testing (`TEST`)

| ID | Sev | Title | Action |
| --- | --- | --- | --- |
| TEST-001..011 | low | Assessment precedence / migration / merge / concurrency / atomic write / smoke / manifest / real smoke gating / runtime boundary / CLI / TUI | no action (positive confirmations) |
| TEST-F-001 | low | Reviewer permission-enforcement branch untested | add fake runtime returning `UnsupportedCount > 0` |
| TEST-F-002 | low | Reviewer exhausted structured-output repair branch untested | make both repair calls fail; assert blocked |
| TEST-F-003 | info | Public seams honored | no action |
| TEST-F-004 | info | Flow-state coverage lives in `sprint_test.go` and `verify_test.go` | no action |
| TEST-F-005 | low | TUI verify coverage is narrow (confirmation, override rationale, focused rerun not asserted) | add coverage |

### Documentation (`DOC`)

| ID | Sev | Title | Action |
| --- | --- | --- | --- |
| DOC-PUBLIC-002-CLI-REFERENCE-STALE-VERIFY | high | `docs/cli-reference.md` still defers integrated `verify` | rewrite and add `verify` subsection |
| DOC-PUBLIC-003-USER-GUIDE-STALE-VERIFY-INTRO | high | `docs/user-guide.md` opening contradicts Sprint 28 verify delivery | rewrite intro and add verification section |
| DOC-OPS-002-RECOVERY-MISSING-VERIFY-SCENARIOS | medium | `docs/recovery.md` lacks verify recovery scenarios | add `## Verify Recovery` section |
| DOC-PUBLIC-001-TUI-HELPTEXT-UNDEFINED | low | TUI in-app `HelpText()` is a generic legend; no Sprint 28 actions | extend `HelpText()` or attach route-scoped help |
| DOC-EXAMPLE-001-CLI-HELP-NO-EXAMPLES | low | CLI help lacks `Examples:` sections | add example blocks to all sprint help |
| DOC-OWNER-001-CLI-REFERENCE-LISTS-WRONG-FLAGS | medium | `docs/cli-reference.md` lists outdated flag lists | update to match parser |
| DOC-AGENT-001-IMPLEMENTATION-REPO-LACKS-DOCS-INDEX | low | Implementation-repo README lacks canonical docs linkage | add `## Canonical docs` to README |
| DOC-OPS-001-VERIFY-HELP-MISSING-DIAGNOSTIC-LIMITS | info | Remediated: `sprintVerifyHelp` documents override limits | no action |
| DOC-PUBLIC-001-REVIEW-SMOKE-VALIDATORS-OK | info | Review/smoke validators enforce required sections | no action |
| DOC-PUBLIC-001-FLOW-STATE-SHAPE-OK | info | Flow-state schema is versioned | no action |
| DOC-ARCH-001-REASONING-OK | info | Architecture decisions recorded | no action |
| DOC-OPS-001-PUBLIC-FLAG-PAIRINGS | info | Flag pairings documented | no action |

### CLI Surface (`CLI`)

| ID | Sev | Title | Action |
| --- | --- | --- | --- |
| CLI-HELP-001-no-worked-examples | low | Sprint help lacks `Examples:` | add examples block |
| CLI-EXIT-001-exit-codes-undocumented | low | Exit codes not documented in help | document table in `renderHelp` and `sprintHelp` |
| CLI-SHAPE-001-positional-order | info | Positional order consistent and documented (remediated) | no action |
| CLI-IO-001-stdout-stderr-discipline | info | Stdout carries results; stderr carries progress | no action |
| CLI-LIFE-001-signal-and-cancellation | info | Signal/cancellation propagation explicit | no action |
| CLI-SAFE-001-and-CLI-NONINT-001 | info | Destructive and non-interactive controls explicit | no action |

### LLM Runtime (`LLM`)

| ID | Sev | Title | Action |
| --- | --- | --- | --- |
| LLM-PROMPT-001-violation | high | Canonical `PromptReference` shape missing; prompt identity not propagated | introduce `PromptReference` (id/version/owner_kind/owner_id/purpose/checksum); propagate |
| LLM-OBS-001-violation | medium | Trace/correlation IDs missing on runtime `Request`/`Result`/`Event` | add `TraceID`/`ParentTraceID` and propagate |
| LLM-RETRY-001-partial | medium | Per-attempt retry observability missing `MaxAttempts`/`AttemptsRemaining`/`IssueKind`/`IssueCode`/`Retryable` | extend `AttemptSummary` with canonical fields |
| LLM-BOUNDARY-001-ok | info | Runtime boundary generic | no action |
| LLM-LIFECYCLE-001-ok | info | Lifecycle states explicit | no action |
| LLM-RUN-001-ok | info | Durable flow-state | no action |
| LLM-IO-001-ok | info | Structured review input/output | no action |
| LLM-EXPOSE-001-ok | info | Operational exposure via app/TUI use cases | no action |
| LLM-EVAL-001-ok | info | Structured review + bounded repair | no action |
| LLM-SAFETY-001-ok | info | Reviewer read-only permission policy | no action |
| LLM-COST-001-ok | info | Token/cost/latency metadata exposed | no action |

### Workflows (`WF`)

| ID | Sev | Title | Action |
| --- | --- | --- | --- |
| WF-VERSION-001 | medium | Workflow behavior not versioned (only state schema is) | add `WorkflowVersion` to `FlowState`/`VerificationAttempt`/`VerificationStatus` |
| WF-IDEMPOTENCY-001 | medium | No explicit idempotency key for retried smoke or reviewer steps | mint deterministic `(sprint, stage, coverage-or-scope, fingerprint)` key |
| WF-RETRY-001-Override-Blocked | medium | `prepareSmokeStatic` allows override for `ReviewVerdictBlocked` (gate mismatch with `Service.Verify`) | tighten to `ReviewFail` only, or unify the gate |
| WF-COMP-001 | low | Reconciliation only; smoke.md/state write-order window | document or reorder to record-replay |
| WF-RETRY-001-AttemptExpiry | info | `attemptExpired` defined but never called | call from `VerificationStatus` |

### Persistence And Migrations (`PERSIST`)

| ID | Sev | Title | Action |
| --- | --- | --- | --- |
| PERSIST-MIG-001 | high | Strict decode can reject a legacy v1 file before migration runs | v1-specific decode or parallel mirror struct |
| PERSIST-INTEGRITY-001 | high | Migration writes `"legacy-unverifiable"` as a digest, pinning migrated state stale | reject migrated record or store real hash; add hex-digest validation |
| PERSIST-ATOMIC-002 | medium | `LoadFlowState` performs a hidden commit | split into pure read + explicit `MigrateFlowState` |
| PERSIST-ATOMIC-003 | low | `syncDir` silently swallows fsync errors | propagate error |
| PERSIST-READ-002 | low | Markdown artifacts are hashed with no size cap | add bounded hash helper |
| PERSIST-INTEGRITY-002 | low | `execute_state.go` is loaded without `DisallowUnknownFields` | mirror the strict-decoder pattern |

### Technical Handbook (`HB`)

| ID | Sev | Title | Action |
| --- | --- | --- | --- |
| HB-001-001..HB-009-009 | info | Policy ownership, lifecycle state, typed failures, narrow boundaries, bounded concurrency, atomic state, config precedence, cancellation, fake tests | no action |
| HB-010-010 | low | TUI verification test surface is thin | add TUI tests for confirmations, override, evidence links |
| HB-011-011 | low | Runtime events-collection goroutine has bounded wait but no deterministic join | drain synchronously or join under `sync.WaitGroup` |
| HB-012-012 | low | `targetIdentity` uses `exec.Command git` directly within sprint | move behind a narrow platform/workspace seam |

## Deviations

No reasoning-departure deviations were performed. The previous automated review (recorded in `flow-state.json`) had to be re-run from scratch because the prior agentwrap-driven attempt failed schema/citation validation or produced OpenCode fatal session errors for 5 of 12 reviewers. The protocol's "Run one structured reviewer per selected contract and one for the technical handbook" rule is satisfied by 12 local-subagent reviewers, none of which depend on the agentwrap runtime. The only operational deviation is the path-citation convention: some reviewers used `target/internal/...` (the prior convention) and others used absolute paths; both resolve to the same files and are accepted.

## Final Assessment

### Per-contract verdict

| Coverage | Applicable | Highest severity | Material findings |
| --- | --- | --- | --- |
| Architecture | yes | low | 2 low |
| Errors | yes | blocker | 1 blocker, 3 high, 2 medium |
| Configuration | yes | medium | 1 medium, 2 low (blockers remediated) |
| Observability | yes | high | 1 high, 3 medium, 1 low |
| Security | yes | low | 1 low |
| Testing | yes | low | 3 low |
| Documentation | yes | high | 2 high, 2 medium, 4 low (some remediated) |
| CLI Surface | yes | low | 2 low |
| LLM Runtime | yes | high | 1 high, 2 medium |
| Workflows | yes | medium | 3 medium, 1 low, 1 info |
| Persistence And Migrations | yes | high | 2 high, 1 medium, 3 low |
| Technical Handbook | yes | low | 3 low |

### Overall verdict

- `pass` requires all mandatory work completed and no applicable blocker/high finding.
- `pass_with_findings` requires all mandatory work completed and only medium/low/info findings.
- `fail` is triggered by **any applicable blocker or high finding**, failed required verification, invalid evidence, or missing mandatory reviewer.
- `blocked` requires unavailable inputs, scope, runtime/model, or verification environment.

This review surfaces **1 applicable blocker** (`ERR-REDACT-001` — SDK response redaction depends on truncation only) and **9 applicable high findings** (`ERR-CODE-001`, `ERR-TASK-001`, `ERR-RETRY-001`, `OBS-CORR-002`, `DOC-PUBLIC-002`, `DOC-PUBLIC-003`, `LLM-PROMPT-001-violation`, `PERSIST-MIG-001`, `PERSIST-INTEGRITY-001`). All 12 mandatory reviewers returned validated structured data. All required verification commands recorded in `execute.md` pass. The verdict is therefore:

**Verdict: fail**

### Required next action

1. `errors` contract: apply `config.RedactValue` to SDK `UserDetail`/`DebugDetail`/`ResponseBody`; replace substring-based code derivation with sentinels; wire `attemptExpired` into `VerificationStatus`; remove or back `OperationError.Retryable`.
2. `persistence-and-migrations` contract: route v1 decode through a non-strict path; either reject migrated v1 outright or store a real recomputed hash with typed `unverifiable` markers; split `LoadFlowState` into pure read + explicit migration; propagate `syncDir` errors.
3. `workflows` contract: gate `prepareSmokeStatic` to `ReviewFail` only (or unify with `Service.Verify`); add `WorkflowVersion` field; mint deterministic `IdempotencyKey` per attempt.
4. `llm-runtime` contract: introduce `PromptReference` and propagate prompt identity through metadata, events, and resume checkpoints.
5. `observability` contract: generate a top-level command/request correlation identifier and propagate through `Runtime`/`Event`/`Operation`.
6. `documentation` contract: rewrite `docs/cli-reference.md` and `docs/user-guide.md` to reflect the integrated `verify` command; add `docs/recovery.md` verify scenarios.
7. After remediation, re-run `go test ./...`, `go test -race ./...`, and `go build ./cmd/ultraplan`, then re-run this review. Promotion to `pass`/`pass_with_findings` requires the blocker and high finding inventory to be empty.

### Recovery notes

- No product, test, governed-input, project-doc, or Git mutation was performed during this review.
- `flow-state.json` reflects the prior agentwrap-driven attempt; the durable review-resume checkpoint was not used by this protocol run because the reviewers are local subagents, not agentwrap runtime instances. The blocked verdict in `flow-state.json` is the correct previous state and remains authoritative for the historical record; this `review.md` becomes the canonical artifact for the current scope.
- `smoke.md` exists and was not modified by this review.
- `execute.md` (target: `implementation complete`) and `.run-state.json` (target: `complete`) remain accurate for the implementation; the canonical verdict for Phase 3 acceptance is `review.md`.
