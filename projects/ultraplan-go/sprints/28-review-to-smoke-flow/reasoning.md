# Sprint Reasoning: Integrated Review-to-Smoke Verification Flow

> Project: `ultraplan-go`
> Sprint: `28-review-to-smoke-flow`
> Output: `projects/ultraplan-go/sprints/28-review-to-smoke-flow/reasoning.md`
> **Inputs Used:** `projects/ultraplan-go/project-index.md`, `projects/ultraplan-go/sprints/28-review-to-smoke-flow/requirements.md`, `projects/ultraplan-go/docs/ARCHITECTURE.md`, `projects/ultraplan-go/docs/PRD.md`, `projects/ultraplan-go/docs/TRD.md`, `projects/ultraplan-go/sprints/28-review-to-smoke-flow/sprint-index.md`, `projects/ultraplan-go/sprints/28-review-to-smoke-flow/technical-handbook.md`, `projects/ultraplan-go/sprints/28-review-to-smoke-flow/reasoning/architecture.md`, `studies/go-cli-study/reports/final/01-project-structure.md`, `studies/go-cli-study/reports/final/02-command-architecture.md`, `studies/go-cli-study/reports/final/03-dependency-injection.md`, `studies/go-cli-study/reports/final/04-configuration-management.md`, `studies/go-cli-study/reports/final/05-error-handling.md`, `studies/go-cli-study/reports/final/06-io-abstraction.md`, `studies/go-cli-study/reports/final/07-state-context.md`, `studies/go-cli-study/reports/final/08-concurrency.md`, `studies/go-cli-study/reports/final/09-terminal-ux.md`, `studies/go-cli-study/reports/final/10-logging-observability.md`, `studies/go-cli-study/reports/final/11-testing-strategy.md`, `studies/go-cli-study/reports/final/13-security.md`, `studies/go-cli-study/reports/final/15-philosophy.md`, `system/contracts/core/architecture.md`, `system/contracts/core/errors.md`, `system/contracts/core/configuration.md`, `system/contracts/core/observability.md`, `system/contracts/core/security.md`, `system/contracts/core/testing.md`, `system/contracts/core/documentation.md`, `system/contracts/surfaces/cli.md`, `system/contracts/runtime/llm.md`, `system/contracts/runtime/workflows.md`, `system/contracts/runtime/persistence-and-migrations.md`, `system/protocols/architecture-review-protocol.md`, `system/protocols/review-sprint-protocol.md`, `system/protocols/deep-smoke-sprint-protocol.md`, injected `builtin:templates/sprint-reasoning.md`

This document decides. It synthesizes the selected context, handbook evidence, architecture reasoning, contracts, and required protocols into final Sprint 28 decisions. It does not execute implementation.

## Sprint Purpose

- **Goal:** Deliver one coherent, resumable, read-only `execute -> review -> smoke` verification flow with freshness detection, deterministic assessment, focused reruns, recovery guidance, and semantic parity across CLI text, stable JSON, and TUI.
- **Non-Goals:** General issue tracking; automatic fixes; product, test, governed-input, documentation, or Git mutation; a third assessment artifact; raw harness evidence in the sprint; replacement of Sprint 26 review or Sprint 27 smoke; browser/hosted/multi-user surfaces; cross-sprint scheduling; or a generic workflow/plugin framework.
- **Depends On:** Sprint 23 execute state and summary behavior, Sprint 26 review orchestration, Sprint 27 smoke-harness integration, the project-index catalog, and the cataloged `ultraplan-go-smoke` harness.

## Selected Context And Pre-Reasoning Artifacts

| Artifact | Path | How It Was Used |
| --- | --- | --- |
| Sprint Index | `projects/ultraplan-go/sprints/28-review-to-smoke-flow/sprint-index.md` | Fixed the 11 selected contracts, 13 evidence reports, Architecture reasoning input, three review protocols, scope, dependencies, and excluded context. No unselected study or reasoning context is introduced. |
| Technical Handbook | `projects/ultraplan-go/sprints/28-review-to-smoke-flow/technical-handbook.md` | Supplied high-confidence patterns for thin surfaces, explicit state, typed failures, cancellation, boundary injection, explicit argv, structured status, bounded concurrency, and fake-backed tests; supplied warnings against shell strings, alternate UI state, globals, and speculative frameworks. |
| Architecture Reasoning | `projects/ultraplan-go/sprints/28-review-to-smoke-flow/reasoning/architecture.md` | Established `internal/sprint` ownership, one typed app projection, persisted facts plus read-time derivation, fingerprint manifests, last-complete artifact preservation, sequential orchestration, focused-rerun containment, narrow external seams, and per-sprint mutation exclusion. These conclusions are adopted without override. |
| Project Architecture | `projects/ultraplan-go/docs/ARCHITECTURE.md` | Constrained dependency direction and assigned sprint policy, app composition, TUI rendering, generic runtime execution, and generic process execution to their owning modules. |
| Product And Technical Requirements | `projects/ultraplan-go/docs/PRD.md`, `projects/ultraplan-go/docs/TRD.md` | Defined Phase 3 behavior, current `review.md`/`smoke.md`, review-before-smoke, deterministic product verdicts, external evidence containment, atomic state, fake seams, and CLI/TUI parity. |
| Prior Decision | None selected | `project-index.md` catalogs no prior decision artifact. Sprint 23/26/27 outputs are dependencies, not selectable prior decisions. |

## Area-Specific Reasoning Inputs

| Area | Reasoning Document | Key Conclusion | Evidence Basis | Impact On Final Decision |
| --- | --- | --- | --- | --- |
| Architecture | `projects/ultraplan-go/sprints/28-review-to-smoke-flow/reasoning/architecture.md` | Keep the module-driven architecture; consolidate verification policy and durable transitions in `internal/sprint`; expose one typed semantic projection through `internal/app`; keep runtime/process generic and injected; add no workflow framework. | Handbook thin-surface, DI, state/context, explicit-argv, testing, and deliberate-complexity evidence; project Architecture/PRD/TRD; accepted costs of read-time hashing, explicit attempt state, and containing-suite reruns. | All final decisions preserve these ownership and state boundaries. This document additionally fixes assessment precedence, override limits, schema compatibility, freshness manifests, and evidence requirements. |

## Sprint Technical Handbook Summary

- **Relevant Patterns:** Thin adapters over shared operations; explicit composition with narrow volatile-boundary seams; explicit lifecycle state; caller-context propagation; bounded and joined concurrency; typed failures with separate recovery rendering; injected IO/runtime/process effects; post-merge configuration validation; executable-plus-argv process invocation; diagnostics separate from result output; behavior tests with fakes and selected goldens.
- **Important Trade-Offs:** Shared semantic results add application types; persisted state can drift from editable artifacts; preserving last-complete evidence requires separate attempt state; focused reruns save time but cannot establish containing-scope confidence; rich failure categories add taxonomy; status hashing adds read cost; sequential gating favors correctness over speculative latency.
- **Warnings / Anti-Patterns:** No verdict/freshness policy in CLI or TUI; no globals or context service locator; no detached goroutines; no process-exit-only success; no shell strings from Markdown; no diagnostics on JSON stdout; no collapse of stale/blocked/cancelled/malformed/fail; no broad interfaces or generic workflow engine without demonstrated pressure.
- **Evidence Confidence:** High for architecture, lifecycle, cancellation, error, security, IO, configuration, observability, and testing decisions because multiple mature Go repositories provide concrete examples. Medium for philosophy-based complexity limits because those are design heuristics; they reinforce, but do not independently determine, scope.

## Contracts Applied

| Contract / Requirement ID | Constraint | Decision Impact | Expected Evidence |
| --- | --- | --- | --- |
| Architecture: `ARCH-CORE-001`, `ARCH-CORE-002`, `ARCH-LAYER-001`, `ARCH-LAYER-002`, `ARCH-ENTRY-001`, `ARCH-MODULE-001`, `ARCH-COMP-001`, `ARCH-SHARED-001` | Product ownership, inward dependencies, thin surfaces, narrow real seams, domain-neutral platform. | Verification policy stays in `internal/sprint`; app composes typed use cases; CLI/TUI render; runtime/process remain generic. | Architecture protocol review, import review, app/surface tests, fake substitution without private wiring. |
| Errors: `ERR-CORE-001`, `ERR-SHAPE-001`, `ERR-CODE-001`, `ERR-TRANS-001`, `ERR-RETRY-001`, `ERR-IO-001`, `ERR-TASK-001`, `ERR-DEP-001`, `ERR-DATA-001`, `ERR-REDACT-001`, `ERR-USER-001` | Preserve cause and machine category; distinguish operational outcomes; make recovery actionable and safe. | Typed blocked, failed, stale, malformed, cancelled, timed-out, conflict, and partial/incomplete results drive exits and next actions. | Error/exit matrix tests, cause-chain checks, redaction tests, status recovery assertions. |
| Configuration: `CFG-SOURCE-001`, `CFG-TYPE-001`, `CFG-COMPAT-001`, `CFG-SECRET-001`, `CFG-START-001`, `CFG-ENV-001`, `CFG-PUBLIC-001`, `CFG-OBS-001` | Explicit precedence, typed post-merge validation, safe environment, compatible durable config. | Review model, harness command, timeout, environment, rerun scope, and override are resolved once and frozen before execution. | Precedence and invalid-config tests; dry-run source display; allowlist/redaction checks. |
| Observability: `OBS-CORE-001`, `OBS-CORR-001`, `OBS-HEALTH-001`, `OBS-DIAG-001`, `OBS-DEBUG-001`, `OBS-TASK-001`, `OBS-ALERT-001`, `OBS-PII-001` | Truthful readiness and terminal state, correlated runs, structured diagnostics, no secret leakage. | Flow/review/smoke identities, fingerprints, attempts, evidence links, safe diagnostics, staleness, and next action are visible from one projection. | Structured event/log assertions; JSON cleanliness; correlation and redaction checks; status snapshots. |
| Security: `SEC-INPUT-001`, `SEC-INJECT-001`, `SEC-SECRETS-001`, `SEC-DEPS-001`, `SEC-FILES-001`, `SEC-NET-001`, `SEC-DESER-001`, `SEC-DEFAULT-001`; `SEC-AUTHN-001`/`SEC-AUTHZ-001` not triggered | Validate untrusted paths/state/evidence; fail closed; explicit argv; contained writes; redaction. Local flow introduces no identity or protected-object authorization boundary. | Cataloged manifest controls harness invocation; environment is allowlisted; paths are contained; verification cannot mutate product/Git. | Path/symlink escape, argv, environment, redaction, malformed-state, and mutation-boundary tests/review. |
| Testing: `TEST-SEAM-001`, `TEST-UNIT-001`, `TEST-INT-001`, `TEST-SMOKE-001`, `TEST-SMOKE-002`, `TEST-FAIL-001`, `TEST-DET-001`, `TEST-CONTRACT-001`, `TEST-E2E-001`, `TEST-MIGRATION-001` | Deterministic policy tests, public fake seams, compatibility/failure coverage, gated real evidence. | Test ownership follows sprint/app/TUI boundaries; normal tests need no OpenCode/network/real harness. | Targeted packages, full/race tests, build, fixtures, command journey tests, and gated harness evidence. |
| Documentation: `DOC-OWNER-001`, `DOC-ARCH-001`, `DOC-PUBLIC-001`, `DOC-OPS-001`, `DOC-EXAMPLE-001`, `DOC-GEN-001`, `DOC-AGENT-001` | One source of truth, valid help/examples, operator recovery, documented architecture/artifacts. | Help and summaries document gates, reruns, override, state, links, exits, and recovery; no third assessment file. | Help tests/review; validators for required review/smoke sections; documentation review. |
| CLI Surface: `CLI-SHAPE-001`, `CLI-HELP-001`, `CLI-IO-001`, `CLI-EXIT-001`, `CLI-JSON-001`, `CLI-CONFIG-001`, `CLI-LIFE-001`, `CLI-SAFE-001`, `CLI-NONINT-001` | Predictable commands, stable JSON/exits, stdout/stderr separation, noninteractive safety, cancellation. | `verify --to review|smoke`, `flow --to smoke`, status and reruns call shared use cases; JSON is final-result-only. | Command/help/exit/golden tests, non-TTY tests, cancellation/cleanup tests, JSON schema fixtures. |
| LLM Runtime: `LLM-BOUNDARY-001`, `LLM-IO-001`, `LLM-LIFECYCLE-001`, `LLM-RETRY-001`, `LLM-RUN-001`, `LLM-PROMPT-001`, `LLM-EXPOSE-001`, `LLM-OBS-001`, `LLM-EVAL-001`, `LLM-SAFETY-001`, `LLM-COST-001`; `LLM-TOOL-001` only if reviewer tools are enabled | Agentwrap stays generic; product validates structured review; lifecycle, prompt identity, cost, timeout, and retry are inspectable. | Existing review runtime seam is reused; an LLM cannot choose verdicts or write canonical output directly. | Fake-runtime structured result tests; malformed/missing reviewer tests; prompt identity and metadata checks; gated runtime evidence. |
| Workflows: `WF-SCOPE-001`, `WF-BOUNDARY-001`, `WF-STATE-001`, `WF-RETRY-001`, `WF-IDEMPOTENCY-001`, `WF-COMP-001`, `WF-VERSION-001` | Explicit ordered states, resume/rerun/cancellation semantics, idempotency, reconciliation, versioning. | One sprint-owned sequential orchestrator and mutation guard; no generic engine; focused promotion rules prevent false confidence. | Transition-table, interruption, conflict, duplicate-rerun, stale propagation, and journey tests. |
| Persistence: `PERSIST-SCHEMA-001`, `PERSIST-MIG-001`, `PERSIST-ATOMIC-001`, `PERSIST-INTEGRITY-001`, `PERSIST-READ-001`, `PERSIST-DERIVED-001`, `PERSIST-FIXTURE-001`, `PERSIST-RECOVERY-001` | Version strict state, define sources of truth, atomic replacement, bounded reads, safe migration/recovery. | Persist facts and last-complete records; derive freshness/assessment; support exactly one known predecessor migration and reject unknown versions. | Migration fixtures, malformed/unknown schema tests, atomic-failure tests, artifact-preservation and bounded-read tests. |
| `requirements.md` Acceptance Criteria 1-13 | Ordering, freshness, parity, reruns, issues, deterministic assessment, recovery, read-only behavior, fakes, and build gates. | Every final decision below maps directly to one or more criteria. | `internal/sprint`, app/command, and TUI tests plus full/race/build commands and required protocol reviews. |

## Repos Studied / Source Evidence Used

| Source / Repo / Report | Concrete Reference | Relevant Finding | Why It Matters For This Sprint | Used In Decision(s) |
| --- | --- | --- | --- | --- |
| Helm via reports `01`, `02`, `07` | `helm/pkg/cmd/install.go:132-145,333-347` | Commands delegate to actions and propagate signal cancellation. | Supports thin CLI adapters and one context through long-running verification. | 1, 6, 7 |
| gh-cli via reports `02`, `03`, `06`, `10`, `11` | `gh-cli/pkg/cmdutil/factory.go:16-43`; `pkg/cmd/factory/default.go:26-46`; `pkg/iostreams/iostreams.go:52-54,551-568`; `acceptance/acceptance_test.go:26-29` | Explicit factories, injected streams, separated outputs, and behavior-level command tests. | Supports shared app use cases, test seams, JSON cleanliness, and command journey coverage. | 1, 7 |
| Restic via reports `03`, `07`, `08`, `13` | `restic/internal/backend/backend.go:19-90`; `internal/restic/lock.go:290-305`; `internal/repository/repository.go:567`; `internal/options/secret_string.go:15-20` | Narrow external boundaries, distinct cleanup policy, structured concurrency, and secret-safe values. | Supports generic process/runtime seams, bounded cleanup, review fan-out, and redaction. | 1, 6, 7 |
| Go-task via reports `05`, `07`, `08` | `go-task/errors/errors_task.go:13-32`; `go-task/task.go:87-89` | Typed recoverable errors and context-coupled fan-out preserve actionable failure and cancellation. | Supports typed next actions and bounded independent review collection. | 6, 7 |
| Chezmoi and k9s via reports `04`, `08` | `chezmoi/internal/cmd/config.go:2253-2287`; `k9s/internal/config/k9s.go:423-451`; `k9s/internal/pool.go:21,30,37` | Track explicit flag changes, validate merged config, and bound concurrent work. | Supports frozen effective configuration and bounded reviewer concurrency. | 7 |
| Lazygit via reports `06`, `09`, `11`, `13` | `lazygit/pkg/commands/oscommands/cmd_obj_runner.go:18-23`; `fake_cmd_obj_runner.go:17-26`; `cmd_obj_builder.go:38`; `pkg/tasks/tasks.go:31-435` | External commands and progress are testable behind explicit boundaries; argv avoids shell interpretation. | Direct support for fake harness tests, safe invocation, and transient progress separate from final artifacts. | 6, 7 |
| OpenCode via report `08` | `opencode/cmd/root.go:261-279` | Shutdown waits must be bounded. | Supports bounded process-tree cleanup and truthful cleanup failure. | 6 |
| Age and gh-cli via report `05` | `age/cmd/age/tui.go:37-54`; `gh-cli/internal/ghcmd/cmd.go:44-49` | Recovery hints and stable exit mapping can remain separate from underlying error identity. | Supports one typed result rendered appropriately in CLI, JSON, and TUI. | 7 |
| Helm and lazygit via report `11` | `helm/internal/test/test.go:43`; `lazygit/pkg/commands/oscommands/fake_cmd_obj_runner.go:17-26` | Goldens catch deliberate output drift while fakes isolate effects. | Supports semantic parity fixtures plus selected stable-rendering goldens. | 7 |
| Philosophy report across gh-cli/restic/lazygit | `studies/go-cli-study/reports/final/15-philosophy.md`; `lazygit/VISION.md:97` | Add complexity only for named pressure; preserve explicit non-goals. | Supports integration around existing stages and rejection of a plugin registry, database, or workflow engine. | 1 and all rejected alternatives |
| Sprint technical handbook | `projects/ultraplan-go/sprints/28-review-to-smoke-flow/technical-handbook.md`, Relevant Patterns, Trade-Offs, Anti-Patterns | Synthesizes all selected report evidence for this sprint. | Provides the direct evidence bridge from studied repositories to the final architecture and evidence plan. | 1-7 |

## Trade-Off And Debt Analysis

### Accepted Trade-Offs

| Trade-Off | Benefit | Cost / Constraint Accepted | Why Acceptable Now | Revisit Trigger |
| --- | --- | --- | --- | --- |
| Persist facts, derive freshness and assessment | External edits cannot leave cached conclusions falsely authoritative; no third artifact. | Status reads must build manifests, hash content, and reconcile evidence. | Truthfulness is a core requirement and workspace scale is bounded enough to start conservatively. | Measured status latency violates product targets or large repositories make full fallback manifests impractical. |
| Preserve last complete artifact beside active/failed attempt | Cancellation, timeout, malformed output, and failed reruns cannot destroy trusted evidence. | State model and UI distinguish current attempt, last complete evidence, and freshness. | Recovery requirements explicitly demand this distinction. | A future append-only history model is selected and can preserve the same invariant more simply. |
| Sequential stage flow with bounded review fan-out | Review-before-smoke is deterministic and avoids waste; independent findings remain collectable. | No speculative smoke and potentially longer wall time. | Smoke is expensive and depends on review; only reviewers are independent. | Profiling proves a safe independent subflow with no gate or state dependency. |
| Focused reruns require complete containing-scope evidence for promotion | Speeds diagnosis without letting a narrow pass erase broader failure. | Some focused success still requires a broader rerun before verdict improves. | False clean passes are more costly than additional harness time. | Harness protocol gains a validated equivalence/coverage proof accepted by project governance. |
| One typed semantic projection with separate renderers | One truth drives CLI/JSON/TUI and later interfaces. | Projection can grow and each renderer still needs tests. | Parity is mandatory and duplicate policy is unacceptable. | Consumers demonstrably require incompatible use-case-specific result shapes. |
| One known forward migration, strict rejection otherwise | Existing Sprint 26/27 state can be upgraded without permanent dual truth. | Older/ambiguous state requires operator recovery instead of best-effort inference. | Inventing freshness would create false confidence. | A documented older schema is found in supported installations and can be migrated deterministically. |
| Per-sprint mutation exclusion | Prevents logically conflicting writes beyond what atomic rename alone protects. | Concurrent operators can receive a retryable conflict. | Verification mutations are infrequent and correctness dominates throughput. | Multi-operator scheduling becomes an explicit product requirement. |

### Potential Technical Debt

| Debt / Shortcut | Why It Might Accrue | Current Mitigation | Owner / Follow-Up |
| --- | --- | --- | --- |
| Conservative target fallback fingerprint | Non-Git targets require hashing all contained regular files and may include generated churn. | Stable ordering, VCS exclusion only, bounded reads, and performance tests; never trade correctness for an unsafe cache. | `internal/sprint`; optimize after measured status data. |
| Shared verification projection growth | Adding every renderer concern could create a god DTO. | Keep product semantics only; exclude raw runtime/harness payloads; split by use case when consumers diverge. | App/sprint architecture review in future interface sprints. |
| One-version migration coverage | Supported historical state may prove broader than currently cataloged. | Strict rejection with recovery guidance; fixtures for known predecessor and unknown versions. | Persistence follow-up only if real supported legacy state is discovered. |
| Platform-specific descendant cleanup | Process-tree termination differs on Linux/macOS and may be incomplete. | Generic result exposes cleanup outcome; bounded cleanup; no clean pass after cleanup failure where mutation safety is uncertain. | `internal/platform/process`; gated smoke on supported platforms. |
| Fingerprint and evidence terminology | Multiple identities (input fingerprint, artifact digest, run ID, evidence identity) can confuse operators. | Stable names in domain/JSON/help and parity fixtures. | Documentation and CLI review in Sprint 28. |

### Future Considerations

| Consideration | Deferred Until | Reason Deferred | What Should Be Preserved Now |
| --- | --- | --- | --- |
| Browser surface | Phase 4, Sprint 30+ | Explicitly excluded from Sprint 28. | Typed app use cases and product-owned durable state, with no CLI parsing dependency. |
| Verification history | A requirement for retained historical attempts/artifacts | Current contract requires current summaries and resumability, not a history product. | Stable run/attempt IDs and last-complete records without a new sprint-root hierarchy. |
| General issue management | A future explicitly governed product phase | Harness issues only influence evidence and verdicts here. | Stable external issue IDs/links and relevance facts, no assignment/scheduling fields. |
| Cross-sprint scheduling | Explicit multi-sprint orchestration scope | One selected sprint is the bounded context. | Per-sprint use cases, context propagation, and mutation exclusion. |
| Workflow/plugin framework | Multiple concrete workflows or adapters prove a common contract | Current pressure is one ordered flow and two established external boundaries. | Narrow consumer-owned runtime/process seams; no registry dependency. |
| Fingerprint cache | Measured hashing cost justifies it | Cache invalidation could compromise freshness. | Deterministic manifest format and explicit content identities. |

## Decisions

Sprint 28 adopts the seven final decisions below as the binding implementation direction. Together they fix product ownership, durable and derived state, freshness manifests, review-before-smoke gating, focused-rerun promotion, interruption-safe persistence, safe external execution, and CLI/JSON/TUI parity. `plan.md` must execute these decisions without reopening architecture.

## Final Decisions

### Decision 1: Keep Verification Policy Product-Owned

- **Decision:** Integrate Sprint 26 review and Sprint 27 smoke with a small convergence refactor. `internal/sprint` exclusively owns stage ordering, gates, freshness, rerun promotion, validation, verdicts, overall assessment, next action, and artifact paths. `internal/app` composes typed status/flow/review/smoke/verify/rerun/recovery use cases. CLI and TUI only collect input, invoke those use cases, render, and map exits. Dependencies remain `cmd/ultraplan` and `internal/tui` -> `internal/app` -> `internal/sprint` -> allowed project/workspace/platform packages; platform runtime/process import no product semantics.
- **Rationale:** One product owner prevents CLI, JSON, and TUI from becoming contradictory workflow engines while preserving established stage implementations.
- **Study / Source Grounding:** `technical-handbook.md` Relevant Patterns cites Helm command delegation (`helm/pkg/cmd/install.go:132-145`), gh-cli factories (`gh-cli/pkg/cmdutil/factory.go:16-43`), and restic boundary interfaces (`restic/internal/backend/backend.go:19-90`). Architecture reasoning lines 11-15 and 32-37 reaches the same conclusion. `ARCHITECTURE.md` lines 226-293 and TRD 18.1 require this ownership.
- **Trade-Offs Accepted:** Typed requests/results and converged call sites add some code; this is accepted to achieve required parity. Keep interfaces only at runtime, process, clock, progress, and hard-to-fixture persistence seams.
- **Technical Debt / Future Impact:** The app projection could grow; raw evidence and rendering-only fields are excluded, and future surfaces must reuse use cases rather than enlarge `internal/sprint` with transport concerns.
- **Alternatives Rejected:** Surface-owned integration is rejected because it duplicates policy. A new verification package, generic workflow engine, event bus, or plugin registry is rejected because one sprint module already owns the state and no second implementation proves a stable abstraction. Wholesale review/smoke replacement is rejected as unnecessary regression risk.
- **Contracts Satisfied:** Architecture all selected IDs; `WF-BOUNDARY-001`; `CLI-SHAPE-001`; `LLM-BOUNDARY-001`; requirements Acceptance Criteria 1, 10, and 11.
- **Evidence Required:** Import/dependency review under Architecture Review Protocol; tests showing CLI and TUI call typed app use cases; platform package tests/types containing no sprint/review/smoke semantics; no CLI-handler or subprocess self-invocation from TUI.

### Decision 2: Version Facts And Derive Current Truth

- **Decision:** Extend the single sprint-owned `flow-state.json` to a strict new schema. Each review/smoke record separates execution status from completed verdict and stores active/last attempt facts, last-complete artifact path/digest, normalized input-manifest fingerprint, timestamps, safe diagnostics, and stage-specific run/evidence/issue/override metadata. Persisted facts are authoritative only for the attempt that produced them; current freshness, readiness, overall assessment, and next action are deterministic read-time projections. Support one explicit atomic migration from the immediately preceding supported schema; preserve known stage facts and complete references, mark unverifiable freshness stale, and reject older, unknown, or malformed versions with recovery guidance.
- **Rationale:** Editable Markdown and external harness evidence can change outside the process. Cached freshness or assessment would drift, while deriving all lifecycle facts from Markdown would lose interruption, override, and attempt information.
- **Study / Source Grounding:** `technical-handbook.md` Explicit Lifecycle State and Centralized State trade-off (lines 34-35, 58-60) support inspectable state but warn about drift. Architecture reasoning lines 18-23 defines persisted facts and read-time projections. TRD 18.5 and Persistence `PERSIST-SCHEMA-001`, `PERSIST-DERIVED-001` require this split.
- **Trade-Offs Accepted:** Reads perform bounded reconciliation and hashing; migration is intentionally conservative. Unknown legacy state blocks instead of being guessed fresh.
- **Technical Debt / Future Impact:** Status performance and terminology need measurement/documentation. The schema must not become a second issue store or duplicate raw harness evidence.
- **Alternatives Rejected:** A third `assessment.md`, verification database, or alternate TUI state is rejected by scope and source-of-truth rules. Markdown-only reconstruction is rejected because attempts and fingerprints are not recoverable. Persisting staleness/overall assessment as authoritative is rejected because external edits invalidate it.
- **Contracts Satisfied:** `WF-STATE-001`, `WF-VERSION-001`; all Persistence IDs, with `PERSIST-FIXTURE-001` applying to test-data separation; `OBS-HEALTH-001`; requirements Acceptance Criteria 2, 3, 6, 7, and 9.
- **Evidence Required:** Schema fixtures for old/current/unknown/malformed state; deterministic projection tests; external-edit tests; migration preservation tests; atomic-write failure tests; state/artifact disagreement recovery tests.

### Decision 3: Fingerprint Deterministic Stage Manifests

- **Decision:** Review input manifests contain normalized workspace-relative paths in lexical order, explicit missing markers, content digests for all governed planning inputs, selected catalog entries and selected contract/protocol contents, execute summary/run-state evidence, target implementation identity, and explicit changed-path scope. Target identity uses configured contained root plus trustworthy Git commit and changed-file content digests; without trustworthy Git metadata it uses a contained regular-file manifest excluding only VCS metadata. Timestamps are metadata, never freshness inputs. Review artifact content has a separate digest. Smoke input manifests include the current review fingerprint/verdict/artifact digest, harness catalog and versioned manifest identity, selected scope and containing-suite requirement, target identity, relevant prerequisite/config identities, and validated issue/evidence identities; `smoke.md` has its own digest. A changed governed input, selected item, execute evidence, target identity, `review.md`, or `smoke.md` makes the owning evidence stale; stale review always makes smoke stale.
- **Rationale:** This resolves exactly what changes invalidate evidence without self-referential fingerprints or timestamp-only freshness.
- **Study / Source Grounding:** The handbook identifies editable-artifact divergence and stale propagation as primary design pressures (lines 117-125); no studied repository defines UltraPlan's exact manifest. The normalization policy is therefore grounded primarily in requirements Acceptance Criterion 2, TRD 18.9.1/18.9.3, Sprint Review Protocol section 2, and architecture reasoning lines 20-23. Study evidence remains relevant to explicit state, not to the product-specific field list.
- **Trade-Offs Accepted:** Conservative fallback hashing may be slower and generated-file-sensitive; this is preferable to false freshness. Safe caching is deferred until measured.
- **Technical Debt / Future Impact:** Changed-scope provenance and symlink/large-file handling require careful fixtures. Future cache optimization must preserve the exact manifest invariant.
- **Alternatives Rejected:** Mtime freshness is rejected as nondeterministic and forgeable. Artifact path/existence alone is rejected as unable to detect edits. Full raw harness copying is rejected because the external harness remains evidence owner.
- **Contracts Satisfied:** `PERSIST-DERIVED-001`, `PERSIST-READ-001`, `SEC-FILES-001`, `SEC-DESER-001`, `TEST-DET-001`; requirements Acceptance Criteria 2, 3, and 6.
- **Evidence Required:** Golden manifest/fingerprint fixtures; mutations of every required input class; ordering, missing marker, path containment, symlink escape, Git/non-Git fallback, external artifact edit, stale-review-to-smoke propagation, and bounded-read tests.

### Decision 4: Use One Ordered Gate And Deterministic Assessment

- **Decision:** `verify --to smoke` and `flow --to smoke` use the same sequential transition function: require complete execute evidence, obtain a current review, then evaluate the smoke gate. Default smoke requires current review verdict `pass` or `pass_with_findings`. A diagnostic override is an explicit smoke-run request only: it requires confirmation and non-empty rationale, records actor-neutral request metadata and the blocked/stale/failing review identity, and never changes review verdict, freshness, smoke freshness propagation, or overall assessment. It grants no extra mutation or environment capability.
- **Rationale:** Review is cheaper and deterministic; one gate prevents command-specific behavior. Override supports diagnosis without laundering failed or stale review evidence.
- **Study / Source Grounding:** Handbook Sequential Flow vs Concurrency and configuration evidence (lines 44-46, 66) supports explicit ordered intent. Review Sprint Protocol lines 11-17 and Deep Smoke Protocol lines 19-31 define the gate and override. Architecture reasoning lines 27-30 adopts one transition path.
- **Trade-Offs Accepted:** No speculative smoke; diagnostic smoke may complete yet remain stale/non-promotable because its review is not current and acceptable.
- **Technical Debt / Future Impact:** Override metadata must remain small and safe; it is not an exception/waiver system.
- **Alternatives Rejected:** Continuing smoke silently after review failure is rejected as unsafe. Making `flow` and `verify` separate orchestrators is rejected as divergence. Allowing override to yield a clean overall result is rejected as contradictory.
- **Contracts Satisfied:** `WF-SCOPE-001`, `WF-STATE-001`, `CLI-SAFE-001`, `CFG-PUBLIC-001`; requirements Acceptance Criteria 1, 6, 9, and 12.
- **Evidence Required:** Shared transition-table tests for both commands; pass/pass-with-findings gate tests; stale/fail/blocker gate tests; explicit override/rationale/audit projection tests; proof override cannot improve freshness, verdict, or mutation policy.

The derived overall assessment uses this fixed precedence over current validated evidence:

| Condition, evaluated in order | Overall Assessment | Required Next Action |
| --- | --- | --- |
| Contradictory/malformed state or referenced canonical evidence cannot be validated | `blocked` | Repair/recover the named state or rerun the owning stage; never infer pass. |
| Review is current and `fail` | `fail` | Resolve review findings/evidence and rerun review; diagnostic smoke cannot change this. |
| Review is current and `blocked` | `blocked` | Satisfy the named review prerequisite and rerun review. |
| Required review or smoke is missing, running, cancelled, timed out, stale, or has only an unpromoted failed attempt | `incomplete` | Run/retry the earliest required non-current stage; retain last complete evidence only as historical recovery context. |
| Current acceptable review plus current smoke `fail` | `fail` | Resolve the classified product/harness/environment cause, run focused diagnosis as needed, then the containing required suite. |
| Current acceptable review plus current smoke `blocked` | `blocked` | Restore prerequisites/coverage/evidence and rerun smoke. |
| Current acceptable review plus current smoke `not_applicable` | `not_applicable` | No runtime smoke is applicable; retain the current review verdict and rationale. |
| Current acceptable review plus current smoke `pass_with_open_issues` | `pass_with_findings` | Resolve linked relevant issues and rerun narrow plus containing required scope for a clean pass. |
| Current review `pass_with_findings` plus current smoke `pass` | `pass_with_findings` | Address review findings if a clean assessment is required. |
| Current review `pass` plus current smoke `pass` | `pass` | No verification action required. |

Execution status is orthogonal to verdict: a successfully completed review or investigation may truthfully have a failing/blocked verdict. `stale` is freshness, not a verdict. The table is the only overall-assessment policy and produces no third artifact.

### Decision 5: Promote Focused Reruns Only With Complete Coverage

- **Decision:** A focused review rerun can promote a new canonical `review.md` only by producing one validated complete result: rerun results may replace selected reviewer/check results, while every retained result must match the same current review fingerprint and mandatory selected-contract/handbook coverage remains complete. Smoke scope precedence is sprint suite, directly relevant suite, explicit level/suite/test for diagnosis, then full harness only for broad closure. A focused passing smoke result is diagnostic until its required containing suite is current and passing. Relevant issue IDs are selected from validated harness metadata; any relevant open issue prevents `pass` and yields `pass_with_open_issues` or `fail` according to evidence. Resolution requires a fix outside verification, passing narrow rerun, passing containing suite, and resolved issue evidence.
- **Rationale:** Focused work reduces cost while canonical evidence remains complete and cannot hide an unaffected stale failure.
- **Study / Source Grounding:** Handbook Narrow Rerun vs Containing-Scope Confidence (lines 61-62), Deep Smoke Protocol sections 3 and 5, and architecture reasoning lines 29-30 directly support this policy. The exact issue-aware verdict rule comes from requirements Acceptance Criteria 4-5 and the selected protocol, not from the Go CLI study.
- **Trade-Offs Accepted:** A narrow success may not improve the canonical verdict immediately; issue relevance validation can block a clean result.
- **Technical Debt / Future Impact:** Review merge logic and harness issue relevance are subtle. UltraPlan stores only IDs, links, relevance/status facts, and fingerprints; it does not own issue lifecycle.
- **Alternatives Rejected:** Replacing canonical smoke with a passing test-level rerun is rejected as false confidence. Retaining review results from another fingerprint is rejected as stale coverage. Automatically changing harness issues or product code is rejected as prohibited mutation.
- **Contracts Satisfied:** `WF-IDEMPOTENCY-001`, `PERSIST-INTEGRITY-001`, `TEST-SMOKE-001`, `TEST-SMOKE-002`; requirements Acceptance Criteria 4, 5, 6, and 8.
- **Evidence Required:** Focused review merge tests; mandatory-reviewer coverage tests; stale retained-result rejection; scope precedence tests; narrow-pass/containing-suite-fail tests; open/resolved issue fixtures; `smoke.md` and status issue/link assertions.

### Decision 6: Preserve Evidence Through Atomic, Cancellable Attempts

- **Decision:** Starting review/smoke atomically records a new attempt under a per-sprint mutation guard without removing the last-complete stage record or replacing canonical Markdown. All blocking work receives the caller context and is joined. Cancellation, timeout, malformed output, validation failure, process failure, or cleanup failure closes the attempt with a typed outcome and safe next action. Only a fully validated result is written to a temporary file, flushed/closed, atomically renamed to `review.md`/`smoke.md`, and then promoted with matching state; failures between artifact and state updates are reconciled by digest/fingerprint on the next read and never inferred current. Required descendant cleanup may use a separate bounded cleanup context and reports cleanup status separately from the primary outcome.
- **Rationale:** Atomic files alone do not preserve logical consistency or prevent concurrent valid operations from racing. Attempts and reconciliation make interruption resumable without destroying evidence.
- **Study / Source Grounding:** Restic cleanup and explicit state, Helm signal propagation, go-task structured context, and OpenCode bounded shutdown are cited in `technical-handbook.md` lines 34-38 and 102-105. Architecture reasoning lines 18-20, 30, and 41-49 defines the adopted shape.
- **Trade-Offs Accepted:** State transitions are more explicit and a concurrent operator may receive conflict. Cross-file atomicity is achieved through ordering plus reconciliation, not an unavailable filesystem transaction.
- **Technical Debt / Future Impact:** Cleanup portability and stale-lock recovery require gated evidence and conservative behavior.
- **Alternatives Rejected:** Overwriting canonical Markdown at attempt start is rejected because interruption destroys trusted evidence. Fire-and-forget progress/process goroutines are rejected because they leak and outlive state. Relying on atomic rename without mutation exclusion is rejected because logical races remain.
- **Contracts Satisfied:** `PERSIST-ATOMIC-001`, `PERSIST-RECOVERY-001`, `WF-COMP-001`, `CLI-LIFE-001`, `ERR-TRANS-001`; requirements Acceptance Criteria 7-9.
- **Evidence Required:** Fault-injected write/rename/state-update tests; cancellation/timeout/cleanup failure tests; previous-artifact byte preservation; concurrent mutation conflict and race tests; interrupted running-attempt recovery; no goroutine leak checks where practical.

### Decision 7: Freeze Safe Execution And Expose One Semantic Result

- **Decision:** Before execution, resolve and freeze typed effective configuration with highest-to-lowest precedence `explicit command flags -> environment -> workspace configuration -> built-in defaults`, consistent with the established config contract. Show non-secret sources in dry-run/confirmation. Review uses the existing generic agentwrap seam with read-only permissions, stable prompt identity, bounded reviewer concurrency, timeout, and validated structured results. Smoke resolves the cataloged versioned machine manifest and invokes a generic process runner with executable, argv, contained cwd, allowlisted environment, timeout/context, bounded output, and cleanup; Markdown/README prose is never executable. The shared result exposes stage execution state, verdict, freshness/reasons, artifacts, run/evidence IDs, relevant issues, override facts, assessment, and next action. Text/TUI may show bounded progress; stable JSON emits only the final result to stdout, while diagnostics/logs go to stderr or configured sinks.
- **Rationale:** Safe, inspectable boundaries make real integrations testable and keep all surfaces semantically aligned without leaking provider or harness internals.
- **Study / Source Grounding:** Explicit precedence (`chezmoi/internal/cmd/config.go:2253-2287`), injected runners (`lazygit/pkg/commands/oscommands/cmd_obj_runner.go:18-23`), explicit argv (`lazygit/cmd_obj_builder.go:38`), stream separation (`gh-cli/pkg/iostreams/iostreams.go:52-54`), and typed errors (`go-task/errors/errors_task.go:13-32`) are summarized in handbook lines 40-50. Review/Smoke protocols define product-specific validation.
- **Trade-Offs Accepted:** More preflight metadata and typed mappings; raw output is deliberately unavailable in stable sprint status and remains at its owning external source.
- **Technical Debt / Future Impact:** Stable JSON fields become compatibility obligations. Prompt/model/runtime changes require evaluation appropriate to risk.
- **Alternatives Rejected:** Shell execution, direct environment reads in product logic, unbounded stdout/stderr capture, runtime/process adapters containing verdict rules, JSON progress on ordinary stdout, and TUI parsing CLI output are rejected for security, testability, and parity reasons.
- **Contracts Satisfied:** Configuration, Observability, Security, CLI Surface, LLM Runtime, Errors, Documentation, and Testing selected IDs; requirements Acceptance Criteria 3, 8-13.
- **Evidence Required:** Preflight/config precedence tests; executable/argv/cwd/env capture with fake runner; path escape and redaction tests; malformed runtime/harness/evidence tests; semantic field parity across app/text/JSON/TUI; help/exit tests; fake-only normal suite and gated real evidence.

## Expected Evidence

| Evidence Type | Required Evidence | Source / Command / Review Check |
| --- | --- | --- |
| Domain and state tests | Ordering, transition table, assessment precedence, fingerprints, stale propagation, focused promotion, issue verdicts, migration, interruption, reconciliation, conflict, and read-only invariants. | `../ultraplan-go/internal/sprint/verify_test.go` and focused `go test` package runs. |
| App/CLI tests | `verify`, `flow --to smoke`, status, rerun/override flags, stable exits, stdout/stderr discipline, malformed/blocked recovery, and text/JSON semantic parity. | `../ultraplan-go/internal/app/sprint_verify_commands_test.go`; command fixtures/goldens where output is intentionally stable. |
| TUI tests | Shared use-case actions, gate explanations, confirmation, progress/cancellation, evidence links, reruns, recovery, narrow terminal rendering, and agreement with app results. | `../ultraplan-go/internal/tui/verify_test.go`; deterministic Bubble Tea messages, no interactive terminal. |
| Boundary/security tests | Fake review runtime and fake harness; explicit argv/cwd/env/timeout; path/symlink containment; bounded output; cleanup; redaction; no source/test/governed-input/Git mutation. | Fake seams and temporary workspaces; mutation snapshots; Security and Architecture review checks. |
| Full verification | All packages pass normally and under race detection; binary builds. | From `../ultraplan-go`: `go test ./...`, `go test -race ./...`, `go build ./cmd/ultraplan`. |
| Runtime evidence | Review/smoke attempts expose correlated safe metadata; stable JSON remains uncontaminated; last-complete artifacts survive failed attempts. | Status/JSON fixtures, structured log assertions, flow-state inspection, cancellation/timeout scenarios. |
| Real smoke | Matching cataloged harness run and issue evidence for runtime-facing claims, or a truthful protocol-valid `blocked`/`not_applicable` result. | `system/protocols/deep-smoke-sprint-protocol.md`; links from current `smoke.md` to external harness `runs/` and `issues/`. |
| Review | Contract, handbook, decision, plan, test, mutation, and architecture conformance with deterministic verdict. | `system/protocols/architecture-review-protocol.md` and `system/protocols/review-sprint-protocol.md`; current `review.md`. |
| Documentation | Help and operator guidance for commands, scope, config, override, statuses, exits, freshness, reruns, issue links, cancellation, timeout, and recovery; valid required sections in `review.md`/`smoke.md`. | Command help snapshots/review; stage validators; Documentation contract review. |

## Assumptions And Risks

| Item | Type | Impact | Mitigation / Follow-Up |
| --- | --- | --- | --- |
| Sprint 26 and 27 expose composable review/smoke seams inside `internal/sprint`. | Assumption | A larger refactor may otherwise be needed. | Plan begins with seam verification and converges existing behavior without changing decided ownership or replacing stages. |
| Immediately preceding flow-state schema is identifiable and deterministic to migrate. | Assumption | Ambiguous historical fields could be guessed incorrectly. | Migrate only known fields; mark unverifiable freshness stale; reject unknown versions with recovery. |
| Harness manifest provides structured scope, evidence, issue, and environment metadata required by protocol. | Assumption | Missing capabilities block safe smoke or issue-aware verdicts. | Return typed `blocked/coverage_missing` or malformed-contract recovery; never parse README prose as fallback. |
| Non-Git fingerprint fallback can be expensive or noisy. | Risk | Status latency and frequent stale results. | Stable bounded manifest; measure; optimize only with correctness-preserving cache identity. |
| Artifact/state update cannot be one filesystem transaction. | Risk | Crash may leave a new valid artifact with old state or vice versa. | Ordered atomic writes plus digest/fingerprint reconciliation; never infer current from path alone. |
| Focused review merge may retain obsolete findings. | Risk | Incomplete coverage could be promoted. | Require identical current fingerprint and complete mandatory reviewer/handbook coverage before promotion. |
| Harness issue relevance may be absent or ambiguous. | Risk | A clean pass could ignore an open issue. | Accept only validated structured relevance; ambiguity prevents clean pass and yields recovery guidance. |
| Concurrent operators can race stage mutations. | Risk | Lost logical updates despite atomic files. | Per-sprint mutation guard, conservative stale-lock handling, actionable conflict result, race tests. |
| Descendant cleanup varies by operating system. | Risk | Orphaned process or uncertain mutation boundary. | Bounded generic cleanup result, truthful non-pass, Linux/macOS gated evidence. |
| Shared projection may become oversized. | Risk | Tight coupling and difficult compatibility. | Include only semantic fields required by current consumers; no raw payloads; architecture review exported API. |
| Diagnostic override may be mistaken for acceptance. | Risk | Failing/stale review could appear approved. | Require rationale; visibly label override; preserve review/overall precedence; tests prohibit clean promotion. |
| Relevant open issues change externally after smoke. | Risk | Previously clean smoke can become stale or non-clean. | Include validated issue identities/statuses in smoke freshness and recompute assessment at read time. |

## Implementation Constraints

- Preserve the architecture reasoning conclusions; implementation may refine names but must not move policy out of `internal/sprint` or create a new workflow architecture.
- Compose existing Sprint 26 review and Sprint 27 smoke behavior; do not replace them wholesale.
- Use one `flow-state.json`; do not add `assessment.md`, a verification database, alternate TUI persistence, or raw harness copies.
- Keep execution status, verdict, and freshness as separate concepts. Use only protocol-defined review and smoke verdicts; use the overall-assessment table in this document.
- Implement the exact manifest/fingerprint ownership and stale propagation described in Decisions 2 and 3; timestamps never establish freshness.
- Permit smoke diagnostic override only under Decision 4; require explicit rationale and never let it improve review, freshness, or overall assessment.
- Promote focused reruns only under Decision 5; a narrow smoke pass never replaces current containing-suite evidence.
- Preserve last complete `review.md` and `smoke.md` until a replacement fully validates; use atomic same-directory replacement and reconciliation.
- Serialize verification mutations per sprint; status/freshness reads may use immutable snapshots and must not mutate state.
- Propagate caller context through every runtime/process/progress wait; join owned work; bound any separate cleanup context.
- Run the harness only from the cataloged versioned machine manifest with executable plus argv, contained cwd, allowlisted environment, timeout, bounded capture, and no shell interpretation.
- Keep `internal/platform/runtime` generic agentwrap infrastructure and `internal/platform/process` generic process infrastructure; neither may import product modules or synthesize verdicts.
- CLI and TUI must call shared typed app use cases. TUI must not invoke CLI handlers/binary, parse stdout, or interpret provider/harness payloads.
- Keep ordinary JSON output final-result-only and free of logs, progress, ANSI, and raw harness output; diagnostics go to stderr or configured sinks.
- Verification may write only owned flow state/current summaries and harness-declared external evidence. It must not mutate product source/tests/docs, governed inputs, or Git state.
- Normal tests must use fake runtime/harness/process seams and require no credentials, network, OpenCode, or real harness. Real dependencies remain gated.
- Generated paths and diagnostics are workspace-relative unless an absolute local path is required for explicit recovery.
- Embedded prompts/templates remain sufficient; workspace overrides are optional and, when present, validated intentional inputs.

## Plan Handoff

`plan.md` must execute these decisions. It must not invent architecture, scope, verdicts, assessment precedence, fingerprint fields, override behavior, rerun promotion rules, or persistence policy beyond this document.

The plan must carry forward:

- all seven final decisions;
- every applicable contract and requirement mapping;
- the deterministic assessment table and freshness rules;
- expected tests, runtime evidence, documentation, and three required review protocols;
- risks, assumptions, migration/reconciliation behavior, and mutation boundaries;
- targeted implementation in the required output files named by `requirements.md`.

## Phase Exit Criteria

- [x] Selected context was read and used.
- [x] Architecture area reasoning was completed and its conclusions are reflected without override.
- [x] Contracts are explicitly mapped to decisions and expected evidence.
- [x] Final decisions, rejected alternatives, trade-offs, debt, assessment precedence, and implementation constraints are specific enough for `plan.md`.
- [x] Expected evidence is specific and reviewable.
- [x] No core architecture decision remains open or unresolved.
