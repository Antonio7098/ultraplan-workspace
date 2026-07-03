# Sprint Reasoning: Documentation and Packaging

> Project: `ultraplan-go`
> Sprint: `15-documentation-and-packaging`
> Output: `projects/ultraplan-go/sprints/15-documentation-and-packaging/reasoning.md`
> **Inputs Used:** `projects/ultraplan-go/sprints/15-documentation-and-packaging/requirements.md`, `projects/ultraplan-go/sprints/15-documentation-and-packaging/reasoning.md`, `projects/ultraplan-go/sprints/15-documentation-and-packaging/sprint-index.md`, `projects/ultraplan-go/sprints/15-documentation-and-packaging/technical-handbook.md`, `projects/ultraplan-go/docs/PRD.md`, `projects/ultraplan-go/docs/TRD.md`, `projects/ultraplan-go/docs/ARCHITECTURE.md`, `projects/ultraplan-go/project-index.md`, `templates/sprint-reasoning.md`

This document decides. It synthesizes selected context, handbook evidence, area-specific reasoning, and contracts into final sprint decisions.

It does not replace `sprint-index.md`, `technical-handbook.md`, or `reasoning/*.md`.

## Sprint Purpose

- **Goal:** Prepare the first production-ready study-side release by completing user/operator documentation, packaging Linux and macOS CLI binaries with SHA-256 checksums, and recording offline plus gated OpenCode smoke-release evidence.
- **Non-Goals:** Do not implement target scaffolding, sprint planning, sprint execution, target/sprint validators, hosted service behavior, browser UI, multi-user collaboration, new runtime adapters, mandatory real OpenCode tests, schema migrations beyond already-shipped behavior, signing, notarization, tags, GitHub releases, or artifact upload.
- **Depends On:** Sprints 1-14 study-side capabilities, especially workspace/config/health, study initialization/listing, Markdown source applicability, prompt preview, analysis and synthesis runs, run-all, durable run-loop, summary/code extraction, validation/diagnostics/JSON stability, and the external OpenCode environment only when available for gated smoke.

## Selected Context And Pre-Reasoning Artifacts

| Artifact | Path | How It Was Used |
| --- | --- | --- |
| Project Index | `projects/ultraplan-go/project-index.md` | Established project scope, target implementation directory, active contract pool, available evidence reports, and the explicit deferred target/sprint scope. |
| Product Requirements | `projects/ultraplan-go/docs/PRD.md` | Grounded release documentation in study-side user workflows, local-first CLI scope, deferred target/sprint commands, artifact requirements, runtime success versus product validation, and launch expectations. |
| Technical Requirements | `projects/ultraplan-go/docs/TRD.md` | Grounded configuration, agentwrap/OpenCode boundaries, run state, retry/fallback, validation, logging, security, cancellation, locks, atomic writes, and testing/release constraints. |
| Architecture | `projects/ultraplan-go/docs/ARCHITECTURE.md` | Preserved module ownership: `internal/study` owns study workflows, `internal/workspace` owns workspace behavior, `internal/codeextract` owns extraction, and `internal/platform/runtime` remains generic runtime infrastructure. |
| Sprint Index | `projects/ultraplan-go/sprints/15-documentation-and-packaging/sprint-index.md` | Selected all applicable contracts, all go-cli-study evidence reports, required review protocols, non-goal exclusions, and the expected next artifact chain for this sprint. |
| Technical Handbook | `projects/ultraplan-go/sprints/15-documentation-and-packaging/technical-handbook.md` | Provided evidence patterns and warnings for thin CLI entrypoints, command-reference accuracy, config precedence, recovery docs, offline-first release gates, security-conscious evidence, bounded concurrency, and local packaging. |
| Sprint Reasoning Template | `templates/sprint-reasoning.md` | Supplied the required structure for final sprint decisions, trade-off analysis, expected evidence, assumptions, risks, implementation constraints, and plan handoff. |

## Area-Specific Reasoning Inputs

| Area | Reasoning Document | Key Conclusion | Evidence Basis | Impact On Final Decision |
| --- | --- | --- | --- | --- |
| Architecture | `reasoning/architecture.md` | Not available; the selected `reasoning/` directory does not exist for this sprint. Architecture decisions are made directly in this sprint reasoning from `ARCHITECTURE.md`, `sprint-index.md`, and `technical-handbook.md`. | `ARCHITECTURE.md` module rules; `technical-handbook.md` thin entrypoint and internal ownership evidence; `sprint-index.md` selected Architecture contract. | Final decisions require documentation and packaging to describe the single CLI entrypoint and current study-side modules without adding new product architecture or runtime ownership claims. |

## Sprint Technical Handbook Summary

- **Relevant Patterns:** Thin CLI entrypoints delegating to internal modules; command references checked against command factories/help surfaces; explicit configuration precedence and validation; classified recovery errors; capturable IO for smoke evidence; context-aware cancellation and durable state language; offline-first release gates with gated real-runtime smoke; secret redaction; bounded concurrency; explicit non-goals.
- **Important Trade-Offs:** Comprehensive release docs must stay strictly study-side; real OpenCode smoke stays gated rather than mandatory; local checksums provide auditability without full signing/publication; rich diagnostics must be balanced against secret and prompt minimization; recovery docs must be actionable without inventing new behavior.
- **Warnings / Anti-Patterns:** Do not claim deferred target/sprint workflows; do not document direct OpenCode supervision by UltraPlan; do not mix logs, data output, and unsafe debug traces in smoke evidence; do not put real runtime smoke into default tests; do not expose provider tokens, full prompts, full report bodies, or raw unsafe payloads; do not silently ship local `agentwrap` replace directives.
- **Evidence Confidence:** High for command/config/testing/security/runtime-boundary decisions because the selected handbook cites concrete Go CLI reports and source references. Medium for terminal UX, performance, extensibility, and philosophy evidence where the findings are more advisory but still relevant to documentation wording and non-goal framing.

## Contracts Applied

| Contract / Requirement ID | Constraint | Decision Impact | Expected Evidence |
| --- | --- | --- | --- |
| Architecture | Preserve module ownership, dependency direction, thin entrypoint, and product/platform separation. | Docs describe the shipped CLI and study-side modules without introducing target/sprint product behavior or direct OpenCode supervision. | README/docs review against `ARCHITECTURE.md`; `go build ./cmd/ultraplan`. |
| Errors | Diagnostics must be actionable and error classes must be represented honestly. | Recovery runbook and CLI reference document validation failures, missing artifacts, lock diagnostics, cancellation, retry/fallback interpretation, and relevant exit behavior. | `docs/recovery.md`, `docs/cli-reference.md`, and smoke evidence review. |
| Configuration | Config precedence, validation, redaction, runtime mapping, and schema rejection must be documented. | Configuration guide becomes the authoritative release doc for supported fields and precedence. | `docs/configuration.md`; review against `ultraplan.yml`, config tests, and help output. |
| Observability | Status, smoke evidence, retry/fallback metadata, runtime metadata, and health truthfulness must be auditable. | Smoke evidence records exact commands/results and safe runtime disposition without unsafe payloads. | `dist/smoke-evidence.md`; status/validate JSON documentation; review for redaction. |
| Security | No secrets, unsafe payloads, full sensitive env dumps, or unsupported permission claims in docs/artifacts. | Docs and evidence use redacted examples and keep OpenCode/provider requirements explicit. | Security review of docs and `dist/`; no tokens or raw unsafe payloads. |
| Testing | Normal verification must be offline, fake-first, repeatable, and race-safe. | Release gates are `go test ./...`, `go test -race ./...`, and `go build ./cmd/ultraplan`; real OpenCode smoke is opt-in. | Passing command logs in `dist/smoke-evidence.md`. |
| Documentation | User-facing docs must be accurate, maintainable, and aligned with public behavior. | Produce user guide, config guide, recovery runbook, OpenCode smoke instructions, release checklist, CLI reference, and README update. | Required docs exist and pass manual review against current command behavior. |
| CLI Surface | Public command surface, output modes, help, exit codes, and stable JSON surfaces must be clear. | CLI reference documents only current study-side commands and stable Sprint 14 JSON surfaces. | `docs/cli-reference.md`; help-output comparison; no unstable JSON promised as stable. |
| LLM Runtime | Runtime integration remains through `agentwrap` and `agentwrap/opencode`. | OpenCode docs describe agentwrap/OpenCode prerequisites and smoke path, not UltraPlan-owned OpenCode process management. | `docs/opencode-smoke.md`; release checklist local replace audit. |
| LLM Evaluation / Cost / Safety | Runtime metadata and evidence must be safe and audit-ready. | Smoke evidence includes pass/fail/skip and redacted metadata, not full prompts or provider secrets. | `dist/smoke-evidence.md` redaction review. |
| Workflows | Run-all, run-loop, status, retries, cancellation, locks, and resumability must be described accurately. | User guide and recovery runbook cover long-running workflows and recovery semantics. | `docs/user-guide.md`, `docs/recovery.md`, `docs/cli-reference.md`. |
| Performance | Bounded concurrency and resource claims must not become unsupported promises. | Docs describe configured parallelism and batch behavior without unmeasured throughput guarantees. | Documentation review; no unsupported performance claims. |
| Persistence And Migrations | Atomic writes, schema rejection, run-state durability, and lock handling must be documented. | Recovery/config docs include schema-version rejection, loud write failures, stale locks, and `--force-unlock` risk. | `docs/configuration.md`, `docs/recovery.md`, release checklist review. |
| REQ-1 to REQ-13 | Required outputs listed in `requirements.md` must be created. | Plan must produce all docs, README update, four binaries, checksums, and smoke evidence. | Required paths exist. |
| AC-1 to AC-17 | Acceptance criteria in `requirements.md` define release readiness. | Plan must run or record all required tests/builds/reviews and keep OpenCode smoke gated. | `dist/smoke-evidence.md`, docs review, binary/checksum inspection. |

## Repos Studied / Source Evidence Used

| Source / Repo / Report | Concrete Reference | Relevant Finding | Why It Matters For This Sprint | Used In Decision(s) |
| --- | --- | --- | --- | --- |
| `go-cli-study` / project-structure report | `studies/go-cli-study/reports/final/01-project-structure.md`; `cmd/restic/main.go:37-114`, `cmd/root.go:112-127`, `cmd/gh/main.go:6` | Production CLIs keep entrypoints thin and delegate behavior inward. | Release docs and packaging should center `./cmd/ultraplan` as the single entrypoint and avoid public Go API promises. | Decisions 1, 4 |
| `go-cli-study` / command-architecture report | `studies/go-cli-study/reports/final/02-command-architecture.md`; `pkg/cmd/root.go:267-303`, `pkg/cmd/install.go:132-145` | Command references should mirror constructed command surfaces and lifecycle behavior. | CLI docs must be checked against actual help/command registration rather than inferred from roadmap language. | Decision 1 |
| `go-cli-study` / configuration-management report | `studies/go-cli-study/reports/final/04-configuration-management.md`; `internal/cmd/config.go:2253-2287`, `internal/flags/flags.go:314-327`, `internal/config/config.go:609-641` | Mature CLIs document precedence and validation after config merge. | Configuration guide must state defaults, workspace config, env vars, flags, redaction, and schema rejection. | Decision 2 |
| `go-cli-study` / error-handling report | `studies/go-cli-study/reports/final/05-error-handling.md`; `cmd/age/tui.go:37-54`, `internal/errors/fatal.go:10`, `fs/fserrors/error.go:22-29` | Recovery docs should distinguish actionable user diagnostics from retry/fatal classes. | Recovery runbook must describe validation, missing artifacts, stale locks, cancellation, retries/fallbacks, and write failures. | Decision 3 |
| `go-cli-study` / IO abstraction report | `studies/go-cli-study/reports/final/06-io-abstraction.md`; `pkg/iostreams/iostreams.go:551-568`, `executor.go:553-577` | Captured IO and separated stdout/stderr make evidence reliable. | Smoke evidence should record commands and summarized outputs cleanly without mixing unsafe debug data. | Decision 5 |
| `go-cli-study` / state-context report | `studies/go-cli-study/reports/final/07-state-context.md`; `pkg/cmd/install.go:333-347`, `internal/restic/lock.go:105`, `internal/session/session.go:12-23` | Long-running commands need cancellation, lock, and session-state recovery language. | Recovery and user docs must explain run-loop interruption, stale running tasks, and `--force-unlock` risk. | Decision 3 |
| `go-cli-study` / concurrency report | `studies/go-cli-study/reports/final/08-concurrency.md`; `pkg/analyze/parallel.go:13`, `task.go:87`, `cmd/root.go:261-279` | Batch execution should use bounded workers and cleanup. | Docs should describe configured parallelism and cancellation without unsupported throughput claims. | Decisions 1, 3 |
| `go-cli-study` / observability report | `studies/go-cli-study/reports/final/10-logging-observability.md`; `internal/logging/logging.go:31-66`, `internal/slogs/keys.go:6-231` | Logs and structured output should be separated and safe. | Smoke evidence and CLI docs must distinguish text, JSON, logs, diagnostics, and stable surfaces. | Decisions 1, 5 |
| `go-cli-study` / testing-strategy report | `studies/go-cli-study/reports/final/11-testing-strategy.md`; `acceptance/acceptance_test.go:26-29`, `internal/cmd/main_test.go:64-174`, `task_test.go:166-169` | Release gates should favor offline fixtures and repeatable command tests. | Default release evidence requires offline tests/build; OpenCode smoke remains gated/skippable. | Decisions 5, 6 |
| `go-cli-study` / security report | `studies/go-cli-study/reports/final/13-security.md`; `internal/options/secret_string.go:15-20`, `pkg/registry/transport.go:37-41`, `internal/permission/permission.go:44-108` | Secrets, auth headers, and permission boundaries need explicit redaction and trust boundaries. | Docs and evidence must avoid tokens, raw unsafe payloads, and direct OpenCode supervision claims. | Decisions 2, 5 |
| `go-cli-study` / performance report | `studies/go-cli-study/reports/final/14-performance.md`; `pkg/cmdutil/factory.go:27-42`, `gdu/pkg/analyze/parallel.go:13` | Bounded concurrency and lazy setup are useful, but unmeasured performance claims are risky. | Documentation should mention supported controls, not publish unsupported benchmarks. | Decisions 1, 6 |
| `go-cli-study` / philosophy report | `studies/go-cli-study/reports/final/15-philosophy.md`; `age.go:18`, `VISION.md:97` | Strong projects state non-goals and reject undocumented complexity. | README and docs must explicitly state deferred target/sprint workflows and publication limitations. | Decisions 1, 6 |

## Trade-Off And Debt Analysis

### Accepted Trade-Offs

| Trade-Off | Benefit | Cost / Constraint Accepted | Why Acceptable Now | Revisit Trigger |
| --- | --- | --- | --- | --- |
| Complete user/operator docs while keeping scope study-side only | Users get an actionable first release without waiting for future target/sprint features. | Docs must repeatedly exclude deferred workflows to avoid confusion. | The PRD/TRD and sprint requirements explicitly make first release study-side only. | PRD/TRD are revised to include target/sprint workflows. |
| Gated OpenCode smoke instead of mandatory real-runtime tests | Offline release gates remain repeatable without credentials, network, OpenCode, or provider setup. | Machines without runtime setup may ship with real smoke skipped. | Sprint constraints require fake-first normal verification and an explicit skip path. | CI/release environment provisions real OpenCode and provider credentials. |
| Local binaries and SHA-256 checksums instead of signed/notarized publication | Produces auditable local release artifacts quickly. | Distribution trust is limited; macOS artifacts are not notarized and no GitHub release is published. | Signing, notarization, tags, and uploads are explicit non-goals. | A separate release-publishing sprint/request is created. |
| Manual release commands/checklist rather than a large packaging subsystem | Keeps this sprint focused on docs and release evidence with minimal implementation risk. | Repetition risk if packaging is done often. | First release needs auditable artifacts, not a full distribution pipeline. | Releases become frequent enough to justify automation. |
| Rich smoke evidence with redaction | Provides reviewable proof while protecting secrets and prompts. | Evidence may omit some debug detail that would aid forensics. | Security constraints forbid secret/raw unsafe payload exposure. | A secure debug-retention policy is designed and enabled explicitly. |
| Document stable JSON only where Sprint 14 made it stable | Protects compatibility promises. | Some useful debug JSON remains undocumented or labeled unstable. | Acceptance criteria require not promising unstable JSON surfaces. | Future sprint promotes additional JSON surfaces to stable. |

### Potential Technical Debt

| Debt / Shortcut | Why It Might Accrue | Current Mitigation | Owner / Follow-Up |
| --- | --- | --- | --- |
| Packaging remains command/checklist driven | Manual cross-compilation and checksum commands can drift from documentation. | Release checklist records exact commands, target list, checksum generation, and verification. | Future release automation sprint if release cadence increases. |
| Documentation can drift from CLI help and stable JSON behavior | CLI behavior may change after docs are written. | Release checklist requires help/output comparison and stable JSON schema review before publishing artifacts. | Every release sprint. |
| OpenCode smoke may be skipped on machines without runtime credentials | Real integration regressions might not be caught by default gates. | `docs/opencode-smoke.md` and `dist/smoke-evidence.md` must record prerequisites and explicit skip reasons. | CI/runtime provisioning follow-up. |
| Local `agentwrap` replace directive could remain in development | A release built with a local replace has unclear provenance. | Release checklist requires auditing `go.mod` and documenting any intentional non-release development state. | Release manager before publishing. |
| Recovery docs may expose implementation-specific lock details | Operators need detail, but internals can change. | Phrase docs around public commands, artifacts, and safe recovery semantics rather than private code assumptions. | Revisit when lock/run-state schema changes. |

### Future Considerations

| Consideration | Deferred Until | Reason Deferred | What Should Be Preserved Now |
| --- | --- | --- | --- |
| Target scaffolding and sprint execution docs | Future product scope revision | Explicitly out of first study-side release. | README/docs must say these workflows are deferred. |
| Binary signing, macOS notarization, and GitHub Releases | Separate release-publishing request | Sprint non-goals exclude publication and signing. | Keep deterministic artifact names and checksums. |
| Packaging automation | Later release hardening | Manual commands are enough for first auditable release. | Release checklist should keep commands scriptable. |
| Additional runtime adapters | Future runtime expansion sprint | Sprint constraints require existing agentwrap/OpenCode path only. | Docs should describe runtime integration through agentwrap, not hard-code direct OpenCode ownership. |
| More stable JSON schemas | Future compatibility sprint | Only Sprint 14 stable surfaces are in scope. | CLI reference must separate stable from debug/unstable output. |
| Secure debug evidence archive | Future observability/security work | Current release must minimize sensitive evidence. | Smoke evidence should record omission/redaction facts. |

## Final Decisions

### Decision 1: Ship Study-Side Documentation Set Only

- **Decision:** Create `docs/user-guide.md`, `docs/cli-reference.md`, and the README release update around the current study-side workflow: workspace/config/health, study init/listing, Markdown source applicability, prompt preview, single run, synthesize, run-all, run-loop/status, validate, summary, and code extraction. The docs must explicitly state that target scaffolding, sprint planning, sprint execution, hosted SaaS, browser UI, and multi-user collaboration are out of scope for this release.
- **Rationale:** The PRD and sprint requirements define the first production release as study-side only. Documentation must make the usable workflow complete without accidentally committing to deferred target/sprint surfaces.
- **Study / Source Grounding:** `technical-handbook.md` relevant patterns for thin release surface and explicit command reference; project-structure evidence from `cmd/restic/main.go:37-114`, `cmd/root.go:112-127`, `cmd/gh/main.go:6`; command-architecture evidence from `pkg/cmd/root.go:267-303`; philosophy evidence from `age.go:18` and `VISION.md:97` supporting explicit non-goals.
- **Trade-Offs Accepted:** Comprehensive docs may repeat scope boundaries, but that is preferable to false product commitments. Performance and workflow claims will be descriptive, not benchmark promises.
- **Technical Debt / Future Impact:** Docs must be revisited if target/sprint workflows are later implemented or if command help changes. The current docs should preserve room for future additions without describing them as available.
- **Alternatives Rejected:** Documenting prototype target/sprint workflows was rejected because PRD/TRD and sprint requirements defer them. Writing only a README quickstart was rejected because acceptance criteria require a full user guide and CLI reference.
- **Contracts Satisfied:** Documentation, CLI Surface, Architecture, Workflows, Performance, REQ-1, REQ-6, REQ-7, AC-8, AC-12, AC-13, AC-14.
- **Evidence Required:** Required docs exist; README links every new document; CLI reference is reviewed against help/current command behavior; docs contain no unsupported target/sprint, SaaS, browser, multi-user, direct OpenCode, automatic Git mutation, or provider-agnostic real-runtime claims.

### Decision 2: Make Configuration and Runtime Boundary Documentation Authoritative

- **Decision:** Create `docs/configuration.md` with supported `ultraplan.yml` fields, precedence order of built-in defaults, workspace config, environment variables, and flags, plus redaction behavior, runtime/model/retry/fallback settings, agentwrap/OpenCode mapping, and schema-version rejection diagnostics.
- **Rationale:** Configuration is a release-critical operator contract. Users need to know which values win, how secrets are handled, and how runtime settings map through agentwrap instead of direct OpenCode ownership.
- **Study / Source Grounding:** `technical-handbook.md` configuration pattern; configuration-management report references `internal/cmd/config.go:2253-2287`, `internal/flags/flags.go:314-327`, and `internal/config/config.go:609-641`; security report references `internal/options/secret_string.go:15-20` and `pkg/registry/transport.go:37-41`; TRD sections 6 and 11 require config precedence and agentwrap mapping.
- **Trade-Offs Accepted:** The guide must stay synchronized with implementation and may avoid documenting internal/debug-only fields. That is acceptable because unsupported fields should not become public compatibility promises.
- **Technical Debt / Future Impact:** Future config schema changes require doc updates and release checklist review. Schema migrations are not introduced by this sprint.
- **Alternatives Rejected:** Treating `ultraplan.yml` as self-documenting was rejected because acceptance criteria require field, precedence, redaction, and schema-version behavior documentation. Documenting direct OpenCode environment/process management as product behavior was rejected because TRD requires agentwrap/opencode ownership.
- **Contracts Satisfied:** Configuration, Security, LLM Runtime, Persistence And Migrations, REQ-2, AC-9, AC-10, AC-15, AC-16.
- **Evidence Required:** `docs/configuration.md` exists; precedence order is explicit; secrets/redaction and unsupported schema behavior are documented; release checklist includes config/stable metadata review and local `replace` audit.

### Decision 3: Publish Recovery Runbook For Operational Failure Modes

- **Decision:** Create `docs/recovery.md` covering validation failures, missing artifacts, stale running tasks, lock diagnostics, conservative `--force-unlock` semantics, cancellation recovery, retry/fallback metadata interpretation, partial completion, and loud atomic write failure behavior.
- **Rationale:** The product principle is that runtime success is not product success and resumability must be explicit. Operators need safe, actionable recovery procedures before the first release.
- **Study / Source Grounding:** `technical-handbook.md` recovery, state-context, and classified-error patterns; error-handling evidence from `cmd/age/tui.go:37-54`, `internal/errors/fatal.go:10`, `fs/fserrors/error.go:22-29`; state/lock evidence from `pkg/cmd/install.go:333-347`, `internal/restic/lock.go:105`, `internal/session/session.go:12-23`; concurrency evidence from `task.go:87` and `cmd/root.go:261-279`; TRD sections 13, 14, 15, 20, and 21.
- **Trade-Offs Accepted:** The runbook will document safe recovery and risk, not provide new behavior. It must avoid encouraging users to force unlock active runs casually.
- **Technical Debt / Future Impact:** If lock/run-state schema changes, recovery docs need updates. Current docs should describe public artifacts and commands rather than private implementation internals.
- **Alternatives Rejected:** Omitting recovery docs was rejected because acceptance criteria require them and run-loop operations are release-critical. Providing manual state-file editing as the primary recovery path was rejected because it is unsafe and bypasses product validation.
- **Contracts Satisfied:** Errors, Observability, Workflows, Persistence And Migrations, Security, REQ-3, AC-11.
- **Evidence Required:** `docs/recovery.md` exists; includes `--force-unlock` risk, stale task interpretation, validation repair workflow, cancellation recovery, missing output handling, retry/fallback metadata, and atomic write failure behavior; review confirms no unsafe secret or raw payload guidance.

### Decision 4: Package Four Local Release Binaries From The Single CLI Entrypoint

- **Decision:** Build `./cmd/ultraplan` into `dist/ultraplan-linux-amd64`, `dist/ultraplan-linux-arm64`, `dist/ultraplan-darwin-amd64`, and `dist/ultraplan-darwin-arm64`, with target `GOOS`/`GOARCH` matching each filename. Do not publish, sign, notarize, tag, or upload artifacts in this sprint.
- **Rationale:** Requirements define local Linux/macOS artifacts as release outputs. Keeping packaging tied to the single CLI entrypoint preserves architecture and avoids introducing a distribution platform.
- **Study / Source Grounding:** `technical-handbook.md` thin CLI and local packaging trade-off; project-structure evidence from thin entrypoints; PRD distribution requirement for a single CLI binary; TRD launch and testing requirements.
- **Trade-Offs Accepted:** Local artifacts with checksums are auditable but not signed or notarized. That limitation is explicit and acceptable because publication/signing are non-goals.
- **Technical Debt / Future Impact:** Manual cross-compilation could be automated later. Artifact names and commands should remain stable to support future release automation.
- **Alternatives Rejected:** Building from package paths other than `./cmd/ultraplan` was rejected because Sprint 1 and requirements identify it as the single CLI entrypoint. Creating installer packages, Homebrew formulas, GitHub releases, signing, or notarization was rejected as out of scope.
- **Contracts Satisfied:** Architecture, Testing, Documentation, REQ-8, REQ-9, REQ-10, REQ-11, AC-3, AC-4.
- **Evidence Required:** `go build ./cmd/ultraplan` passes; all four binaries exist in `dist/`; smoke evidence records build commands and targets; release checklist states supported-platform notes and publication limitations.

### Decision 5: Generate Exact Checksums and Redacted Smoke Evidence

- **Decision:** Generate `dist/checksums.txt` with exactly one SHA-256 entry for each required binary and no unrelated files. Create `dist/smoke-evidence.md` recording exact offline verification commands, pass/fail results, build targets, checksum generation command, and gated OpenCode smoke result or explicit skip reason.
- **Rationale:** The release must be auditable. Checksums prove artifact identity, and smoke evidence proves which gates ran without requiring real runtime access in default tests.
- **Study / Source Grounding:** `technical-handbook.md` testing, IO, observability, and security patterns; testing evidence from `acceptance/acceptance_test.go:26-29`, `internal/cmd/main_test.go:64-174`, `task_test.go:166-169`; IO capture evidence from `pkg/iostreams/iostreams.go:551-568`; observability evidence from `internal/logging/logging.go:31-66`; security evidence from secret redaction references.
- **Trade-Offs Accepted:** Evidence will include enough detail for audit but omit secrets, full prompts, full report bodies, full sensitive environment dumps, and raw unsafe runtime payloads. Real OpenCode smoke may be skipped when prerequisites are unavailable.
- **Technical Debt / Future Impact:** Skipped real smoke can leave runtime integration less covered; future CI should provision gated runtime credentials if release policy requires it.
- **Alternatives Rejected:** Adding OpenCode smoke to default `go test ./...` was rejected because tests must remain offline and fake-first. Including full runtime transcripts or environment dumps was rejected for security. Recording checksums for all `dist/` files was rejected because acceptance criteria require only required binaries and no unrelated entries.
- **Contracts Satisfied:** Testing, Observability, Security, LLM Runtime, LLM Evaluation / Cost / Safety, REQ-12, REQ-13, AC-1, AC-2, AC-5, AC-6, AC-7.
- **Evidence Required:** `dist/checksums.txt` contains exactly four SHA-256 entries; `dist/smoke-evidence.md` records `go test ./...`, `go test -race ./...`, `go build ./cmd/ultraplan`, build targets, checksum command, and OpenCode pass/fail/skip with redaction notes.

### Decision 6: Use A Release Checklist As The Publishing Gate

- **Decision:** Create `docs/release-checklist.md` as the repeatable gate for tests, race tests, build, packaging, checksums, smoke evidence, local `replace` audit, supported-platform notes, prompt/version metadata review, runtime metadata redaction, stable JSON documentation, recovery docs, and security review.
- **Rationale:** This sprint produces local release artifacts but not a full publication pipeline. A checklist is the right control surface to keep release readiness explicit and repeatable without adding unnecessary automation.
- **Study / Source Grounding:** `technical-handbook.md` offline-first release gates, local packaging trade-off, local replace warning, stable JSON warning, and philosophy of explicit non-goals; testing-strategy and security reports; performance report caution against unsupported promises.
- **Trade-Offs Accepted:** Checklist-driven release can be manual, but it keeps first-release scope small and auditable. Automation can be introduced once repeated releases justify it.
- **Technical Debt / Future Impact:** Checklist steps may drift; future automation should derive from the checklist rather than replace the release criteria silently.
- **Alternatives Rejected:** Building a full release automation system was rejected as unnecessary for this sprint. Skipping local `replace` audit was rejected because agentwrap provenance is release-critical. Treating all JSON output as stable was rejected because only Sprint 14 stable surfaces should be compatibility-sensitive.
- **Contracts Satisfied:** Documentation, Testing, Security, CLI Surface, LLM Runtime, Performance, REQ-5, AC-12, AC-15, AC-16, AC-17.
- **Evidence Required:** `docs/release-checklist.md` exists and includes every required gate; checklist references docs and `dist/` artifacts; review confirms dependency audit, stable JSON review, redaction review, recovery docs review, and platform limitations are present.

## Expected Evidence

| Evidence Type | Required Evidence | Source / Command / Review Check |
| --- | --- | --- |
| Tests | Offline test suite passes without OpenCode, provider credentials, network, or real subprocess smoke fixtures. | `go test ./...` from `/home/antonioborgerees/coding/ultraplan-go`; record in `dist/smoke-evidence.md`. |
| Tests | Race test suite passes. | `go test -race ./...` from `/home/antonioborgerees/coding/ultraplan-go`; record in `dist/smoke-evidence.md`. |
| Build | Main CLI builds. | `go build ./cmd/ultraplan`; record in `dist/smoke-evidence.md`. |
| Packaging | Four binaries exist with correct Linux/macOS and amd64/arm64 target names. | Inspect `/home/antonioborgerees/coding/ultraplan-go/dist/`; record exact build commands and targets. |
| Checksums | `checksums.txt` contains exactly four SHA-256 entries for required binaries. | `sha256sum` or equivalent checksum generation/verification command; review no unrelated files are included. |
| Runtime Smoke | Real OpenCode smoke pass/fail is recorded, or skip reason states missing runtime/provider/network/config prerequisites. | `docs/opencode-smoke.md`; `dist/smoke-evidence.md`. |
| Documentation | Required docs exist and cover current study-side workflow only. | Review `docs/user-guide.md`, `docs/configuration.md`, `docs/recovery.md`, `docs/opencode-smoke.md`, `docs/release-checklist.md`, `docs/cli-reference.md`, and `README.md`. |
| CLI Surface | Stable JSON surfaces from Sprint 14 are documented and unstable/debug JSON is not promised as stable. | Review `docs/cli-reference.md` against implementation/tests/help output. |
| Security | Docs and evidence contain no secrets, full sensitive environment variables, raw unsafe runtime payloads, full prompts, or full report bodies. | Manual security/redaction review of docs and `dist/smoke-evidence.md`. |
| Dependency Readiness | Release checklist includes local `replace github.com/Antonio7098/agentwrap => ../agentwrap` audit and disposition. | Review `docs/release-checklist.md` and `go.mod` during implementation. |
| Review | Architecture and sprint review protocols can be run after implementation. | `system/protocols/architecture-review-protocol.md`; `system/protocols/sprint-review-protocol.md`. |

## Assumptions And Risks

| Item | Type | Impact | Mitigation / Follow-Up |
| --- | --- | --- | --- |
| Current CLI help and Sprint 14 stable JSON surfaces are available in the target repo. | Assumption | Docs can be checked against actual behavior. | Implementation plan must inspect help/tests before finalizing CLI reference. |
| Real OpenCode/provider environment may be unavailable. | Risk | Runtime smoke may be skipped, reducing real integration confidence. | Document gated prerequisites and record explicit skip reason in `dist/smoke-evidence.md`. |
| Local `agentwrap` replace may remain during development. | Risk | Release provenance may be unclear if shipped unintentionally. | Release checklist must audit `go.mod` and document/remove local replace before releasable state. |
| Documentation may accidentally describe deferred prototype target/sprint behavior. | Risk | Users may rely on unsupported commands/workflows. | Review README/docs against PRD/TRD deferred scope and sprint-index excluded context. |
| Smoke evidence could leak secrets or unsafe runtime payloads. | Risk | Security incident or unusable release evidence. | Redaction review; record commands/results/summaries only; omit provider tokens, full env, prompts, raw payloads. |
| Manual packaging commands may be run from the wrong directory or with wrong target. | Risk | Artifacts may not match required binaries. | Smoke evidence must record exact working directory, build commands, GOOS/GOARCH targets, and checksum command. |
| Race tests may expose existing implementation issues outside documentation/packaging. | Risk | Release gate blocks until fixed or explicitly triaged. | Plan must run race tests and treat failures as release blockers unless user separately scopes remediation. |
| macOS binaries are cross-compiled but not notarized. | Assumption | Users may need to handle platform trust prompts manually. | Release checklist and README must state supported local artifact scope and no notarization. |

## Implementation Constraints

- Documentation must cover only current study-side capabilities and must not claim target scaffolding, sprint planning, sprint execution, hosted service, browser UI, multi-user collaboration, direct OpenCode supervision, automatic Git mutation, or provider-agnostic real-runtime guarantees.
- Packaging must build from `/home/antonioborgerees/coding/ultraplan-go` using `./cmd/ultraplan` as the single CLI entrypoint.
- Release artifacts must be written under `/home/antonioborgerees/coding/ultraplan-go/dist/` and must not require a server, database, artifact store, signing service, notarization service, GitHub release, or upload.
- `dist/checksums.txt` must include only the four required release binaries.
- Normal release verification must remain offline and fake-first; real OpenCode smoke must be explicitly gated and skippable with a recorded reason.
- Runtime docs must keep integration through `github.com/Antonio7098/agentwrap` and `github.com/Antonio7098/agentwrap/opencode`.
- Docs and evidence must redact secrets and avoid full prompts, full generated report bodies, full sensitive environment variables, provider tokens, and raw unsafe runtime payloads.
- Recovery docs must describe safe public recovery workflows and `--force-unlock` risk without recommending unsafe state corruption.
- Configuration docs must state precedence as built-in defaults, workspace config, environment variables, then command flags.
- Stable JSON documentation must be limited to compatibility-sensitive surfaces introduced by Sprint 14; debug or undocumented JSON must not be labeled stable.
- Release checklist must include local `replace github.com/Antonio7098/agentwrap => ../agentwrap` audit and disposition.
- Plan and implementation must not reopen architecture or add new runtime adapters, schema migrations, target/sprint packages, or product behavior beyond documentation, packaging commands/artifacts, and release verification.

## Plan Handoff

`plan.md` must execute these decisions. It must not invent architecture, scope, or decisions beyond this document.

The plan must carry forward:

- final decisions
- applicable contracts and requirement IDs
- expected evidence
- risks and assumptions
- required review protocols

## Phase Exit Criteria

- [x] Selected context was read and used.
- [x] Area-specific reasoning documents were completed where applicable or explicitly marked not applicable.
- [x] Area-specific reasoning conclusions are reflected or explicitly overridden in final decisions.
- [x] Contracts are explicitly mapped to decisions and expected evidence.
- [x] Final decisions are clear enough for `plan.md` to execute.
- [x] Expected evidence is specific and reviewable.
