# Sprint Plan: Deep Smoke Harness Integration

> Project: `ultraplan-go`
> Sprint: `27-deep-smoke`
> Source: `reasoning.md`
> **Inputs Used:** `projects/ultraplan-go/project-index.md`, `projects/ultraplan-go/roadmap.md`, `projects/ultraplan-go/sprints/27-deep-smoke/requirements.md`, `projects/ultraplan-go/docs/ARCHITECTURE.md`, `projects/ultraplan-go/docs/PRD.md`, `projects/ultraplan-go/docs/TRD.md`, `projects/ultraplan-go/sprints/27-deep-smoke/sprint-index.md`, `projects/ultraplan-go/sprints/27-deep-smoke/technical-handbook.md`, `projects/ultraplan-go/sprints/27-deep-smoke/reasoning/architecture.md`, `projects/ultraplan-go/sprints/27-deep-smoke/reasoning.md`, `system/protocols/architecture-review-protocol.md`, `system/protocols/deep-smoke-sprint-protocol.md`, `system/protocols/review-sprint-protocol.md`, `builtin:templates/sprint-plan.md`

This plan executes `reasoning.md`. It does not introduce architecture, scope, or workflow stages beyond Decisions 1-7. Requirement references use the numbering established by `reasoning.md`: `AC-01` through `AC-13` and `C-01` through `C-10`.

## Reasoning Source

- **Sprint Reasoning:** `reasoning.md`
- **Sprint Index:** `sprint-index.md`
- **Technical Handbook:** `technical-handbook.md`
- **Area Reasoning:** `reasoning/architecture.md`
- **Roadmap Scope:** `projects/ultraplan-go/roadmap.md#sprint-27-deep-smoke-harness-integration`
- **Required Protocols:** `system/protocols/architecture-review-protocol.md`, `system/protocols/deep-smoke-sprint-protocol.md`, `system/protocols/review-sprint-protocol.md`

## Sprint Status

- **Status:** complete
- **Owner:** implementation agent
- **Start Date:** 2026-07-22
- **Completion Date:** 2026-07-22

## Implementation Context Inspected

The target repository is `../ultraplan-go/`. Planning inspected the existing integration points so tasks can extend delivered behavior rather than presume new structures:

- `internal/sprint/review.go`, `domain.go`, `service.go`, `flow.go`, `state.go`, `artifacts.go`, and `validation.go` provide the current governed-stage, fingerprint, flow-state, artifact, and validation patterns.
- `internal/project/domain.go`, `index.go`, and `validation.go` parse the current project catalog but do not yet recognize the `Smoke Harnesses` section or its manifest/evidence columns.
- `internal/app/sprint_commands.go`, `sprint_usecases.go`, `operations.go`, `json_output.go`, and `tui_commands.go` provide command, shared-operation, JSON-envelope, and TUI composition points.
- `internal/tui/app.go`, `model.go`, and `views.go` provide confirmation, bounded progress history, cancellation, result, and preview patterns.
- `internal/workspace/init.go` owns embedded default assets; `internal/platform/config/config.go` owns effective configuration and currently has no smoke settings.
- `internal/platform/process` does not exist and is the one new volatile-boundary package authorized by Decision 1.

## Decisions To Execute

| Decision | Source Section | Execution Implication |
| --- | --- | --- |
| Preserve Module Ownership And One Shared Use Case | `reasoning.md#decision-1-preserve-module-ownership-and-one-shared-use-case` | Keep smoke meaning and persistence in `internal/sprint`, generic process mechanics in `internal/platform/process`, dependency construction and typed operations in `internal/app`, and interaction/rendering in CLI/TUI adapters. |
| Use A Strict Cataloged Protocol-v1 Preflight | `reasoning.md#decision-2-use-a-strict-cataloged-protocol-v1-preflight` | Resolve only the cataloged harness and manifest, support protocol major v1, validate required capabilities and identities, and reject malformed or escaping paths before launch. |
| Resolve Options Once, Gate Review, And Select Sufficient Scope Deterministically | `reasoning.md#decision-3-resolve-options-once-gate-review-and-select-sufficient-scope-deterministically` | Merge options once with source metadata, enforce current-review rules, validate overrides through discovery, and select the smallest discovery-proven sufficient containing-suite scope. |
| Supervise One Direct Process With Bounded Lifecycle And Environment | `reasoning.md#decision-4-supervise-one-direct-process-with-bounded-lifecycle-and-environment` | Execute one direct executable/argv process with explicit cwd/env, context, timeout, output caps, non-blocking progress, process-group cleanup, and observable cleanup facts. |
| Validate Evidence, Synthesize Verdicts, Then Commit The Summary | `reasoning.md#decision-5-validate-evidence-synthesize-verdicts-then-commit-the-summary` | Cross-check structured run/evidence identity before computing one of five verdicts; validate and atomically replace `smoke.md` before updating flow state. |
| Expose One Typed State And Progress Model Across Every Surface | `reasoning.md#decision-6-expose-one-typed-state-and-progress-model-across-every-surface` | Make CLI text/JSON, status, validation, flow, and TUI adapt the same typed readiness, preview, progress, result, error, artifact, and next-action values. |
| Verify In Layers And Gate Real Harness Use | `reasoning.md#decision-7-verify-in-layers-and-gate-real-harness-use` | Use pure tests, protocol fixtures, recording runners, a child fake harness, command/TUI tests, selected goldens, race/build checks, and an explicit opt-in real-harness lane. |

## Requirements / Contracts To Satisfy

| Contract / Requirement ID | Required Behavior | Evidence Planned |
| --- | --- | --- |
| `AC-01` | Current `pass`/`pass_with_findings` review permits smoke; current `fail`/`blocked` requires an explicit recorded diagnostic override; missing, malformed, or stale review blocks. | Gate table tests plus CLI/JSON/TUI confirmation and error assertions. |
| `AC-02` | Catalog and protocol-v1 manifest discovery reject missing, malformed, unsupported, duplicate, capability-deficient, or path-escaping inputs. | Project-index, manifest, symlink, containment, and structured-diagnostic fixtures. |
| `AC-03` | Discovery returns structured levels, suites, tests, mappings, prerequisites, duration/cost, evidence schema, and protocol identity without prose execution. | Discovery schema fixtures and recording-runner assertions. |
| `AC-04` | Default and explicit scope selection remain deterministic and sufficient; narrow diagnostic evidence cannot replace required containing-suite evidence. | Selection and coverage table tests for sprint suite, suite set, level, suite, test, ties, unknown IDs, blocked coverage, and irrelevance. |
| `AC-05` | Harness execution uses direct argv, contained cwd, bounded environment, timeout, cancellation, and owned descendant cleanup. | Recording runner and real child-process tests, including timeout, signal escalation, descendants, caps, and slow progress. |
| `AC-06` | Run ID, safe argv, counts, optional runtime/model/usage, duration, evidence identity, paths, and relevant issue IDs are validated and imported truthfully. | Valid and adversarial run/evidence fixtures, absent optional metadata cases, and redaction assertions. |
| `AC-07` | A fully validated `smoke.md` is atomically replaced only for evidence-backed completion or truthful `blocked`/`not_applicable`; raw evidence remains external. | Markdown golden, validator tests, injected write failures, and last-valid-artifact preservation tests. |
| `AC-08` | Product code deterministically synthesizes `fail`, `blocked`, `not_applicable`, `pass_with_open_issues`, or `pass`. | Complete verdict table, including malformed/process/evidence failures that must not become verdict passes. |
| `AC-09` | Normal smoke mutates only `smoke.md`, approved flow/status state, and manifest-declared external run/issue roots. | Before/after mutation snapshots and architecture/security review. |
| `AC-10` | Status, `validate smoke`, text, stable JSON, and TUI agree on readiness, scope, state, verdict, links, staleness, diagnostics, and next action. Command-level flow advancement to smoke remains deferred. | Shared-use-case scenarios, command tests, JSON goldens/schema assertions, and no-launch status test. |
| `AC-11` | TUI provides readiness, selection, prerequisites, cost/duration, confirmation, progress, cancellation, results, issues, and preview through app use cases. | TUI model/update/view tests, narrow-terminal tests, and semantic parity cases. |
| `AC-12` | Fake harness coverage includes success, failure, blocked, not applicable, timeout, cancellation, malformed JSON, missing evidence, issue references, path escape, and redaction; real harness is opt-in. | Named fixture matrix, child-process suite, and gated integration evidence or truthful blocked reason. |
| `AC-13` | Full tests, race tests, and CLI build pass; unavailable real prerequisites are reported as blocked, not pass. | Recorded `go test ./...`, `go test -race ./...`, `go build ./cmd/ultraplan`, and gated-run result. |
| `C-01`, `C-02`, `C-03`, `C-10` | Preserve sprint/platform/app/TUI ownership and one-way module dependencies. | Import inspection, architecture review, and shared-use-case tests. |
| `C-04` | Keep review on `agentwrap`; use the generic process boundary for smoke without changing `platform/runtime`. | Dependency review and focused tests proving no runtime/self-CLI path was added. |
| `C-05`, `C-06` | Never execute shell/prose/generated command strings; enforce review-before-smoke and record any allowed diagnostic override. | Injection fixtures, exact argv assertions, gate scenarios, and safe preview output. |
| `C-07`, `C-08` | Missing prerequisites/capabilities/coverage cannot pass; generated state is versioned, redacted, editable where applicable, and atomically written. | Blocked/validation/version/redaction/atomic-failure tests. |
| `C-09` | Normal tests remain fake-first; live OpenCode or harness execution is explicit and gated. | Default-suite process logs and opt-in test documentation. |
| Architecture, CLI Surface, Configuration, Documentation, Errors | Package ownership, stable commands/output, deterministic options, readable artifacts, and actionable typed failures conform to selected contracts. | Required protocol review plus package, command, golden, and documentation checks. |
| LLM Evaluation / Cost / Safety, LLM Runtime, Observability, Performance | Optional runtime/cost facts remain truthful, review runtime remains unchanged, progress is typed, and lifecycle/output are bounded. | Metadata, redaction, progress, cap, timeout, race, and dependency tests. |
| Persistence And Migrations, Security, Testing, Workflows | Versioned state, path safety, layered tests, gate/selection/verdict/flow rules, and mutation limits are enforced. | Schema and containment fixtures, failure injection, fake-harness matrix, and selected protocol reviews. |

## Tasks

- [x] **Task 1: Extend Catalog And Effective Smoke Configuration**

> Executes: Decisions 1-3; `AC-02`, `AC-03`, `AC-05`, `C-01`, `C-05`, `C-08`, `C-10`

  - [x] Extend `internal/project` catalog parsing and validation to recognize the indexed `Smoke Harnesses` table, including harness root, manifest, evidence roots, external-path intent, duplicate identity, and actionable missing/escaping diagnostics.
  - [x] Preserve project ownership of catalog data only; do not move manifest decoding, suite discovery, or smoke execution into `internal/project`.
  - [x] Add smoke configuration fields and source metadata through the existing config/app composition path for workspace values, bounded environment-derived settings, timeout/capture controls, and request overrides.
  - [x] Implement the exact effective-option precedence from Decision 3 and validate only the merged result.
  - [x] Resolve `OQ-01` and `OQ-02` below in implementation documentation/tests before accepting configurable bounds or environment forwarding.
  - [x] Add table tests for every source precedence combination, invalid merged settings, unknown settings, absolute external catalog roots, duplicate entries, and redacted effective diagnostics.

- [x] **Task 2: Add The Generic Direct Process Boundary**

> Executes: Decisions 1 and 4; `AC-05`, `AC-12`, `AC-13`, `C-02`, `C-05`, `C-09`, `C-10`

  - [x] Create `internal/platform/process` with a consumer-shaped request/result accepting `context.Context`, executable, argv, cwd, complete environment, timeout, stdout/stderr caps, and an optional bounded progress sink.
  - [x] Return start/finish/duration, exit/signal, bounded stdout/stderr, truncation/drop counts, and explicit descendant-cleanup success/failure/uncertainty without importing product-domain types.
  - [x] Execute directly without a shell and start one owned process group on supported Linux/macOS targets.
  - [x] On cancellation or timeout, terminate gracefully, wait up to the decided 5-second grace period, escalate remaining owned descendants, and always wait/reap; cleanup failure or uncertainty must remain separately observable.
  - [x] Ensure progress backpressure cannot block stdout/stderr draining, cleanup, or delivery of the terminal result.
  - [x] Add recording and child-process tests for exact executable/argv/cwd/env, success, nonzero exit, timeout, cancellation, ignored graceful signal, descendants, large/malformed output, truncation, disconnected/slow progress, cleanup failure, and race/leak behavior.

- [x] **Task 3: Implement Protocol-v1 Preflight, Discovery, And Scope Selection**

> Executes: Decisions 2 and 3; `AC-01`-`AC-04`, `AC-08`, `C-01`, `C-05`-`C-08`, `C-10`

  - [x] Add concrete smoke manifest, discovery, prerequisite, scope, coverage, and safe-command types in focused `internal/sprint` files; support protocol major v1 only.
  - [x] Add or update `/home/antonioborgerees/coding/ultraplan-go-smoke/ultraplan-smoke.json` as the external sprint deliverable with executable/fixed argv prefix, discovery/run forms, cwd, evidence roots, capabilities, declared environment names, and protocol/schema identity.
  - [x] Canonicalize and symlink-resolve the catalog root, manifest, executable, cwd, and evidence roots before execution, then enforce containment at each trust boundary.
  - [x] Strictly decode required v1 fields, reject unsupported major versions, missing capabilities, duplicate IDs, malformed relationships, and ambiguous required semantics, while permitting harmless additive fields.
  - [x] Run machine-readable discovery through the generic process boundary using the 30-second, 4 MiB stdout, and 1 MiB stderr defaults decided in `reasoning.md`; never derive executable input from README or Markdown.
  - [x] Implement the deterministic hierarchy: complete sprint-specific suite, then smallest stable directly mapped complete suite set, otherwise `blocked` when complete coverage cannot be proven.
  - [x] Validate explicit level/suite/test IDs and containing-suite sufficiency; classify missing prerequisites/coverage as `blocked` and discovery-proven irrelevance as `not_applicable` without launching the run process.
  - [x] Produce typed preview data containing scope/rationale, gate/override, prerequisites, optional runtime/model, duration/cost, safe argv, effective limits/sources, allowed mutations, and evidence destination.
  - [x] Add manifest/discovery/path fixtures and complete scope-selection tables named in Decisions 2 and 3.

- [x] **Task 4: Orchestrate Review Gate, Run, Evidence Validation, And Verdict**

> Executes: Decisions 3-5; `AC-01`, `AC-04`-`AC-09`, `C-01`, `C-05`-`C-08`

  - [x] Reuse delivered review validation/fingerprinting to require a current `pass` or `pass_with_findings`; allow only current `fail`/`blocked` through an explicit confirmed diagnostic override and retain the original review verdict/fingerprint everywhere.
  - [x] Keep missing, malformed, or stale review non-overridable in Sprint 27 and return an actionable typed gate error before discovery/run work that is not needed.
  - [x] Construct run argv only from validated manifest constants and validated selected scope, use one harness process, and apply the manifest timeout or the existing 30-minute run default plus decided bounded overrides.
  - [x] Strictly decode the run response and cross-check protocol identity, run ID, requested/observed coverage, process exit, counts, duration, optional runtime/model/usage, summary/detail identity, optional hashes, fallback identity, and issue IDs.
  - [x] Re-canonicalize every returned evidence/issue path, verify it exists under declared roots and identifies the reported run/issue, and keep raw JSON, streams, test artifacts, and issue files external.
  - [x] Implement the exact five-verdict table in Decision 5; treat malformed data, timeout, cancellation, unexplained process failure, missing/mismatched evidence, and failed/uncertain cleanup as execution failures that cannot pass or replace the current summary.
  - [x] Add typed error categories and safe diagnostics for review gate, protocol, scope, prerequisite, containment, process, timeout, cancellation, cleanup, evidence, validation, and persistence/reconciliation failures.
  - [x] Add orchestration and evidence fixtures for all successful, non-passing, blocked, malformed, mismatch, missing metadata, issue, redaction, and cleanup scenarios required by Decisions 4-5.

- [x] **Task 5: Persist And Validate Smoke Artifact And Flow State**

> Executes: Decisions 5 and 6; `AC-06`-`AC-10`, `C-01`, `C-07`, `C-08`

  - [x] Register the `smoke` stage and `smoke.md` artifact in sprint domain/path/validation/status structures without introducing a parallel verification hierarchy.
  - [x] Render all Deep Smoke Sprint protocol sections and required fields: context, review gate/fingerprint/override, harness/protocol, scope/rationale, prerequisites/environment, safe argv, run/count/duration/runtime facts, external evidence identity/links, open/resolved issues, mutation check, verdict, and next action.
  - [x] Validate the complete candidate and reject placeholders, invalid verdicts, missing identity/count/next-action fields, copied raw streams/JSON, unsafe paths, stale evidence, and secret values before replacement.
  - [x] Atomically replace `smoke.md` through a same-directory temporary file with flush/close/rename only for an evidence-backed completed run or truthful preflight `blocked`/`not_applicable` classification.
  - [x] Preserve the last valid `smoke.md` on malformed run/process/evidence, timeout, cancellation, cleanup, validation, or pre-commit persistence failure.
  - [x] After artifact commit, atomically update versioned flow state with execution status separate from verdict, review/smoke fingerprint, artifact path, run/evidence ID, timestamps, override fact, and safe diagnostics.
  - [x] On flow-state write failure, keep the valid artifact, report reconciliation required, and make status/validation detect artifact/state mismatch without implementing Sprint 28 automated recovery.
  - [x] Add golden, validator, atomic-failure, split-state, stale/missing evidence, unsupported-version, and before/after mutation-snapshot tests.

- [x] **Task 6: Expose Shared App, CLI, JSON, Status, And Validation Operations**

> Executes: Decisions 1, 3, 5, and 6; `AC-01`, `AC-03`, `AC-04`, `AC-10`, `C-03`, `C-06`, `C-08`

  - [x] Add typed smoke readiness, preview, run, progress, cancellation, result, validation, and status operations to the existing `internal/app` composition path, injecting process/catalog/filesystem/clock collaborators without mutable globals.
  - [x] Add `ultraplan sprint <project> <sprint> smoke` flags for explicit level/suite/test selection, timeout, diagnostic review override, dry-run/preview, non-interactive confirmation semantics, and stable JSON output.
  - [x] Extend `status` and `validate smoke` to use the same typed sprint/app results and to expose readiness, effective scope/rationale, gate/override, prerequisites, execution status, verdict, staleness, artifact/evidence links, diagnostics, and next action; defer command-level flow advancement to smoke.
  - [x] Keep status and validation static: inspect catalog, manifest, artifact, fingerprint, and paths without launching discovery/run; label dynamic readiness unknown when discovery is required.
  - [x] Render human progress separately from operational diagnostics and stable JSON; include no ANSI, raw streams, native harness events, or secret values in JSON.
  - [x] Map typed smoke failures to stable machine codes and existing meaningful exit classes without collapsing cancellation, blocked, validation, and runtime failures into one textual error.
  - [x] Add command/help/JSON goldens and semantic scenarios proving review gate/override, selection, verdict, links, stale/reconciliation state, next action, no-launch status, and interface agreement.

- [x] **Task 7: Add TUI Smoke Operation And Preview**

> Executes: Decisions 1 and 6; `AC-01`, `AC-03`, `AC-04`, `AC-10`, `AC-11`, `C-03`, `C-06`, `C-08`, `C-09`

  - [x] Add smoke readiness/status and guarded actions to existing TUI navigation and operation composition using only shared app use cases.
  - [x] Present review gate/override, scope selection, rationale, prerequisites, optional runtime/model, duration/cost class, safe command, mutation/evidence roots, and confirmation before launch.
  - [x] Project product phases `preflight`, `discovery`, `selection`, `running`, `validating_evidence`, `writing_artifact`, `completed`, `cancelled`, and `failed` into existing bounded event/history behavior with optional suite/test/count/drop facts.
  - [x] Propagate TUI cancellation context through app/sprint/process cleanup and show the resulting execution state without synthesizing a separate verdict.
  - [x] Display final verdict, counts, issues, evidence links, diagnostics, required next action, and a contained `smoke.md` preview; persist no alternate smoke state.
  - [x] Add model/update/view tests for readiness, selection, confirmation, progress, dropped events, cancellation, blocked/not-applicable, all evidence-backed verdicts, issue summary, preview, errors, and narrow terminals.
  - [x] Add parity tests proving TUI and CLI receive the same semantic request/result values and that no TUI path invokes CLI handlers, shells out to `ultraplan`, or parses terminal text.

- [x] **Task 8: Complete The Layered Fake-Harness And Safety Matrix**

> Executes: Decision 7; `AC-02`, `AC-05`-`AC-13`, `C-05`, `C-08`-`C-10`

  - [x] Organize pure/table tests for review gate, precedence, scope/coverage, verdicts, fingerprints, validators, containment, error mapping, and redaction.
  - [x] Add versioned manifest/discovery/run/evidence fixtures for valid v1, additive fields, malformed JSON, unsupported version, capability/identity errors, count/scope/hash/path/issue mismatches, and absent optional metadata.
  - [x] Use an in-memory recording runner for orchestration and a real fake-harness child executable only for direct-process, timeout, cancellation, signal, descendants, output caps, malformed stream, slow-consumer, and cleanup behavior.
  - [x] Add semantic cross-surface scenarios for CLI text/JSON, status, validation, flow, and TUI covering success, selected-smoke failure, blocked prerequisite/coverage, not applicable, timeout, cancellation, state-write failure, stale/reconciliation, issues, and redaction.
  - [x] Reserve goldens for `smoke.md`, stable JSON envelopes, help, and selected TUI presentation; keep security, selection, evidence, and verdict assertions semantic.
  - [x] Prove normal tests do not contact the real harness, network, OpenCode, or ambient credentials and do not alter product source/tests, governed inputs, harness maintenance inputs, or Git state.
  - [x] Run focused package tests before the full suite and record failures with requirement/decision references rather than weakening assertions.

- [x] **Task 9: Document, Gate Real Harness Use, And Assemble Review Evidence**

> Executes: Decision 7; `AC-07`, `AC-09`-`AC-13`, `C-08`-`C-10`

  - [x] Update embedded `smoke.md` defaults, CLI help, configuration reference, user workflow, TUI guidance, external manifest contract, evidence inspection, diagnostic override, cancellation, blocked/stale/reconciliation, and real-harness setup documentation; keep workspace overrides optional.
  - [x] Resolve `OQ-03` by defining one explicit opt-in environment/config gate and one documented real-harness test command consistent with repository test conventions; default tests must skip without claiming pass.
  - [x] Run the opt-in cataloged harness test when prerequisites exist and record protocol/run/evidence identity; otherwise record the exact blocked prerequisite without manufacturing smoke evidence.
  - [x] Run and record focused tests, `go test ./...`, `go test -race ./...`, and `go build ./cmd/ultraplan` in the target repository.
  - [x] Review the implementation against Architecture Review for ownership/dependency direction, simplicity, lifecycle, state, errors, observability, tests, and no speculative registry/workflow abstraction.
  - [x] Review against Deep Smoke Sprint for gate, catalog/protocol, sufficient scope, safe execution, evidence identity, mutation limits, artifact sections, verdict, and TUI parity.
  - [x] Review against Sprint Review for decision/task conformance, selected contracts, verification records, deviations, deferred scope, and evidence completeness.
  - [x] Record any deviation from `reasoning.md` before implementation continues; do not silently change protocol, gate, scope, process bounds, verdict, persistence order, or deferred Sprint 28 behavior.

## Evidence Checklist

- [x] Tests prove all `AC-01` through `AC-13` behaviors and `C-01` through `C-10` constraints.
- [x] Runtime/process diagnostics are bounded, classified, actionable, and redacted.
- [x] Manifest, discovery, run, evidence, and state schema/version fixtures are complete.
- [x] Path containment includes lexical, symlink, executable/cwd, and returned-evidence escape cases.
- [x] Review gate, diagnostic override, sufficient-scope selection, and all five verdicts have deterministic evidence.
- [x] `smoke.md`, stable JSON, help, and selected TUI goldens were explicitly reviewed.
- [x] Mutation snapshots prove only approved product state and external evidence roots change.
- [x] CLI, JSON, status, validation, flow, and TUI semantic parity is demonstrated.
- [x] Timeout, cancellation, descendant cleanup, output caps, slow progress, race, and leak evidence exists.
- [x] Documentation and embedded defaults are complete; no workspace prompt/template override is required.
- [x] The real-harness lane produced valid evidence or a truthful blocked record.
- [x] `go test ./...`, `go test -race ./...`, and `go build ./cmd/ultraplan` results are recorded.
- [x] Architecture Review, Deep Smoke Sprint, and Sprint Review protocol evidence is recorded.
- [x] Deviations from `reasoning.md` are recorded before implementation continues.

## Verification Commands

Commands run from `../ultraplan-go/` unless noted otherwise.

| Check | Command | Expected Result |
| --- | --- | --- |
| Catalog, protocol, selection, evidence, persistence | `go test ./internal/project ./internal/sprint` | Project catalog and sprint smoke unit/fixture/orchestration tests pass without a real harness. |
| Generic process boundary | `go test ./internal/platform/process` | Direct-process child fixtures pass for argv/env/cwd, bounds, timeout, cancellation, descendants, progress, and cleanup. |
| CLI/app/TUI parity | `go test ./internal/app ./internal/tui` | Command, JSON, operation, model/update/view, cancellation, and preview tests pass. |
| Full deterministic suite | `go test ./...` | All packages pass with fake runtimes/harnesses and no live prerequisite. |
| Race and lifecycle suite | `go test -race ./...` | All packages pass with no reported races; process/progress cleanup tests terminate. |
| CLI build | `go build ./cmd/ultraplan` | The `ultraplan` binary builds successfully. |
| Gated real harness | Command established and documented by Task 9 after resolving `OQ-03` | A protocol-v1 run links valid external evidence, or the test reports a truthful skip/blocked prerequisite and never pass. |
| Architecture review | Inspect imports under `internal/platform/process`, `internal/sprint`, `internal/app`, and `internal/tui` using the selected Architecture Review protocol | Platform has no product imports; smoke semantics and persistence remain sprint-owned; CLI/TUI share app use cases. |

## Assumptions And Open Questions

| ID | Assumption / Open Question | Impact | Required Resolution |
| --- | --- | --- | --- |
| `A-01` | Sprint 26 review artifacts expose a validated verdict and governed-input fingerprint reusable by smoke. | If absent, current-review gating cannot be computed safely. | Reuse delivered fields; if genuinely absent, record a bounded compatibility issue and do not weaken freshness. |
| `A-02` | The external harness can deliver a protocol-v1 manifest and structured discovery/run responses with containing-suite mappings. | Real smoke remains blocked without the contract or sufficient mappings. | Treat the manifest/mapping as required sprint output and validate it through fixtures before gated use. |
| `A-03` | Supported descendant cleanup targets for Sprint 27 are Linux and macOS. | Windows cannot receive a passing cleanup claim. | Keep cleanup facts explicit and classify unsupported/uncertain cleanup as non-passing. |
| `OQ-01` | `reasoning.md` requires positive bounded timeout/capture overrides but does not set hard maximum values. | Guessing ceilings could reject valid suites or permit unsafe resource use. | Resolved as bounded implementation policy: discovery 5 minutes, run 24 hours, cleanup 30 seconds, stdout 64 MiB, and stderr 16 MiB maximum; decided defaults remain 30 seconds, 30 minutes, 5 seconds, 4 MiB, and 1 MiB. Tests cover invalid and bounded overrides. |
| `OQ-02` | `reasoning.md` requires a minimal platform environment plus manifest-declared names but does not enumerate the base variable names. | Over-forwarding can leak secrets; under-forwarding can prevent direct execution. | Resolved to `PATH`, `HOME`, `TMPDIR`, `LANG`, and `LC_ALL`; a manifest name is forwarded only when also present in the configured allowlist, and values never enter persisted diagnostics. |
| `OQ-03` | The exact opt-in real-harness environment/config gate and test command are not named by current evidence. | Inventing a command in the plan would create an ungrounded public/test contract. | Resolved to `ULTRAPLAN_REAL_SMOKE=1 go test ./internal/sprint -run TestRealSmokeHarness -v`; optional workspace/project/sprint variables select the catalog context. The default suite does not launch the real harness. |

No core architecture question remains open. If any open question would alter Decisions 1-7, implementation stops and records the proposed deviation rather than guessing.

## Risks And Blockers

| Risk / Blocker | Source | Mitigation | Status |
| --- | --- | --- | --- |
| Harness lacks protocol-v1 structured contract or sufficient mappings. | `reasoning.md#assumptions-and-risks` | Delivered and independently exercised the protocol-v1 manifest/adapter and Sprint 27 mapping. | mitigated |
| Review artifact lacks required current fingerprint fields. | `reasoning.md#assumptions-and-risks` | Reused Sprint 26 validation and fingerprint fields without weakening freshness. | mitigated |
| Symlink/TOCTOU changes escape validated roots. | `reasoning.md#assumptions-and-risks` | Canonicalize before launch, constrain cwd/roots, and revalidate evidence identity after execution. | mitigated |
| Descendants survive timeout/cancellation or cleanup certainty is unknown. | `reasoning.md#potential-technical-debt` | Process-group child tests cover bounded cancellation, escalation, reaping, and descendant cleanup; uncertain cleanup cannot pass. | mitigated |
| Slow progress consumer blocks process draining. | `reasoning.md#assumptions-and-risks` | Bounded non-blocking dispatch records drops; the slow-consumer child test completes. | mitigated |
| Valid `smoke.md` commits but flow-state update fails. | `reasoning.md#potential-technical-debt` | The artifact remains authoritative and reconciliation is surfaced; automated recovery remains Sprint 28 scope. | accepted boundary |
| Diagnostic override is mistaken for review approval. | `reasoning.md#assumptions-and-risks` | Original verdict/fingerprint and override are retained; failed-review override runs are diagnostic and cannot commit current smoke evidence. | mitigated |
| Secrets leak through argv, URLs, environment, excerpts, JSON, TUI, or Markdown. | `reasoning.md#assumptions-and-risks` | Explicit environment allowlists, safe argv redaction, bounded diagnostics, and no raw-stream persistence are enforced. | mitigated |
| External evidence changes or disappears after summary creation. | `reasoning.md#assumptions-and-risks` | Strong SHA-256 identity and contained paths are recorded; full automated stale-result recovery remains Sprint 28 scope. | accepted boundary |
| Conservative defaults are unsuitable for a real suite. | `reasoning.md#potential-technical-debt` | Manifest/config/request precedence and bounded ceilings resolve `OQ-01`; tuning awaits gated runtime evidence. | mitigated |
| Scope expands into Sprint 28 `verify`, focused reruns, or complete recovery. | `requirements.md#non-goals`; roadmap Sprint 28 | Preserve typed seams and next actions, but reject added convenience orchestration or automated recovery in this sprint. | mitigated by plan boundary |

## Boundaries And Rejected Alternatives

- Do not add a harness-aware platform runner, plugin registry, generic workflow/verification package, or smoke subpackage tree.
- Do not parse README/Markdown/generated text into executable commands or invoke a shell-interpolated command string.
- Do not infer pass from process exit, incomplete coverage, missing prerequisites, malformed evidence, or uncertain cleanup.
- Do not let an explicit test or narrow diagnostic suite replace required containing-suite evidence without discovery-proven complete equivalence.
- Do not let the TUI invoke CLI handlers, shell out to `ultraplan`, parse terminal output, persist alternate state, or synthesize its own verdict.
- Do not inherit the full ambient environment, persist secret values, or copy raw JSON/streams/test artifacts/issues into the sprint root.
- Do not roll back a valid committed `smoke.md` because a later flow-state write fails.
- Do not implement Sprint 28 `verify`, focused-rerun recovery, full stale-result recovery, issue management, automatic remediation, Git operations, hosted/remote smoke, or plugin infrastructure.
- Do not extend command-level flow execution to the smoke stage in this plan; that deferred orchestration remains Sprint 28 scope.

## Review Inputs

Review should use:

- `sprint-index.md`
- `technical-handbook.md`
- `reasoning/architecture.md`
- `reasoning.md`
- this `plan.md`
- implementation diff
- verification and gated-harness evidence
- `system/protocols/architecture-review-protocol.md`
- `system/protocols/deep-smoke-sprint-protocol.md`
- `system/protocols/review-sprint-protocol.md`

## Execution Log

| Date / Step | Action | Evidence / Notes |
| --- | --- | --- |
| 2026-07-19 / planning | Created evidence-grounded implementation plan only. | Decisions 1-7, `AC-01`-`AC-13`, `C-01`-`C-10`, assumptions, risks, rejected alternatives, expected evidence, and three explicit implementation-policy open questions carried forward from the selected inputs. No implementation, smoke, review automation, issue tracking, Git operation, or run-state artifact was executed. |
| 2026-07-22 / implementation | Implemented catalog/config, generic process supervision, protocol-v1 orchestration, evidence/verdict persistence, CLI/JSON/status/validation, TUI operations, embedded defaults, documentation, and the external adapter/manifest. | Target changes are listed in `execute.md`; `/home/antonioborgerees/coding/ultraplan-go-smoke/src/protocol.ts` and `ultraplan-smoke.json` provide the external protocol boundary. |
| 2026-07-22 / policy resolution | Resolved `OQ-01`-`OQ-03` and retained the Sprint 28 boundary. | Bounds and environment are documented above. Command-level flow advancement remains non-mutating; integrated verification/recovery remains deferred. |
| 2026-07-22 / verification | Ran focused tests, full tests, race tests, CLI build, adapter discovery, import inspection, and the opt-in real-harness lane. | All deterministic Go gates passed. Adapter discovery returned protocol-v1 Sprint 27 coverage. The real run truthfully skipped because this sprint has no current `review.md`; no `smoke.md` or external run evidence was manufactured. |

## Execution Evidence

- Focused package tests passed for project catalog, config, direct process supervision, sprint smoke, app/CLI, TUI, and embedded workspace behavior.
- `go test ./...`, `go test -race ./...`, `go build ./cmd/ultraplan`, and `git diff --check` passed in the target repository.
- `node_modules/.bin/tsx src/protocol.ts discover --target /home/antonioborgerees/coding/ultraplan-go` returned a valid protocol-v1 identity, prerequisites, levels, suites, and complete `27-deep-smoke -> ultra-deep` mapping.
- `ULTRAPLAN_REAL_SMOKE=1 go test ./internal/sprint -run TestRealSmokeHarness -v` passed the test wrapper with a truthful skip: `smoke review_gate: a current review is required`.
- Import inspection confirms `internal/platform/process` imports only the standard library; sprint owns smoke semantics/persistence, app composes typed use cases, and TUI imports app rather than CLI handlers.
- The external harness's legacy TypeScript project-wide build still reports pre-existing test-suite type errors. The new adapter itself executes successfully through the repository's existing `tsx` runtime; this maintenance debt does not change the protocol gate or claim a real smoke pass.

## Completion Criteria

- [x] All nine tasks are complete or explicitly deferred with requirement and decision impact recorded.
- [x] Every `AC-01` through `AC-13` acceptance criterion has named passing evidence.
- [x] Every `C-01` through `C-10` constraint has named conformance evidence.
- [x] Verification commands were run, or environment-specific gated real-harness inability is documented as blocked rather than pass.
- [x] Evidence satisfies `reasoning.md#expected-evidence`, including unit, fixture, child-process, command, TUI, golden, race, build, mutation, and gated real-harness layers.
- [x] Required Architecture Review, Deep Smoke Sprint, and Sprint Review protocols can evaluate conformance without guessing intent.
- [x] No unresolved question silently changed protocol, gate, scope, environment, process bounds, verdict, persistence order, or module ownership.
- [x] No Sprint 28 or other explicit non-goal behavior entered the implementation.
- [x] `review.md` can evaluate decision, contract, task, verification, mutation-boundary, and deferred-scope conformance without rereading every study report.
