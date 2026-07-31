# Sprint Reasoning: Local Web Foundation and Read-Only Dashboard

> Project: `ultraplan-go`
> Sprint: `30-web-foundations`
> Output: `projects/ultraplan-go/sprints/30-web-foundations/reasoning.md`
> **Inputs Used:** `projects/ultraplan-go/project-index.md`, `projects/ultraplan-go/sprints/30-web-foundations/requirements.md`, `projects/ultraplan-go/docs/ARCHITECTURE.md`, `projects/ultraplan-go/docs/PRD.md`, `projects/ultraplan-go/docs/TRD.md`, `projects/ultraplan-go/sprints/30-web-foundations/sprint-index.md`, `projects/ultraplan-go/sprints/30-web-foundations/technical-handbook.md`, `projects/ultraplan-go/sprints/30-web-foundations/reasoning/api-design.md`, `projects/ultraplan-go/sprints/30-web-foundations/reasoning/architecture.md`, `projects/ultraplan-go/sprints/30-web-foundations/reasoning/frontend.md`, `builtin:templates/sprint-reasoning.md`

This document decides. It synthesizes selected context, handbook evidence, area-specific reasoning, and contracts into final Sprint 30 decisions. It does not replace `sprint-index.md`, `technical-handbook.md`, or `reasoning/*.md`, and it does not execute implementation.

Requirement references use the order in `requirements.md`: `RO-*` is a Required Output row, `AC-*` is an Acceptance Criteria item, `NG-*` is a Non-Goal, and `C-*` is a Constraint. For example, `AC-3` is the third acceptance criterion and `C-6` is the sixth constraint.

## Sprint Purpose

- **Goal:** Add `ultraplan serve` as a loopback-only, read-only HTTP and browser interface over shared typed app queries, with a single-binary server-rendered dashboard, versioned JSON resources, bounded safe artifact previews, and deterministic lifecycle/security tests.
- **Non-Goals:** No browser-triggered validation or other mutations; no confirmations, operation handles, cancellation endpoint, SSE, or WebSockets; no runtime/smoke invocation; no arbitrary file editing; no database or browser-owned durable state; no hosted, LAN/public, remote, account, tenant, permission, or collaboration model; no frontend framework/build pipeline; no workflow-semantic changes; and no Git mutation (`NG-1` through `NG-10`, `C-5`, `C-11`).
- **Depends On:** Product Phase 3 completion or recorded exceptions, Sprints 24-25 shared local-interface composition, Phase 3 read-only review/smoke status surfaces, existing typed `internal/app` query/error behavior, and the Phase 4 package rules in the project Architecture, PRD, and TRD.

## Selected Context And Pre-Reasoning Artifacts

| Artifact | Path | How It Was Used |
| --- | --- | --- |
| Project Index | `projects/ultraplan-go/project-index.md` | Established the Phase 4 interface direction, available contract/evidence pool, source-of-truth docs, and prohibition on unselected context and Git mutation. |
| Sprint Requirements | `projects/ultraplan-go/sprints/30-web-foundations/requirements.md` | Supplied the authoritative Sprint 30 scope, outputs, acceptance criteria, constraints, dependencies, and verification commands. Its narrower read-only scope controls over broader PRD/TRD Phase 4 operation requirements. |
| Sprint Index | `projects/ultraplan-go/sprints/30-web-foundations/sprint-index.md` | Selected nine contracts, thirteen evidence reports, three area-reasoning artifacts, two required review protocols, and explicit exclusions for operations/SSE, remote serving, browser-owned state, frontend frameworks, and Git mutation. |
| Technical Handbook | `projects/ultraplan-go/sprints/30-web-foundations/technical-handbook.md` | Supplied high-confidence evidence for thin adapters, explicit injection, capability-limited interfaces, source-aware config, typed errors, cancellation, bounded work, redacted diagnostics, deterministic tests, and lazy initialization. |
| Architecture | `projects/ultraplan-go/docs/ARCHITECTURE.md` | Fixed module ownership and dependency direction: `cmd/ultraplan -> app+tui+web`, `web -> app`, and no direct web dependency on product/platform execution modules. |
| Product Requirements | `projects/ultraplan-go/docs/PRD.md` | Confirmed one product core across local interfaces, loopback/single-binary browser inspection, authoritative workspace state, and the Sprint 30 local-web-foundation rollout before guarded operations. |
| Technical Requirements | `projects/ultraplan-go/docs/TRD.md` | Supplied workspace/config precedence, CLI exit conventions, HTTP/browser security and testing requirements, single-binary build constraints, and the broader Phase 4 capabilities explicitly deferred by this sprint. |
| Prior Decision | None selected | `sprint-index.md` records no cataloged prior decision. Dependencies are requirements, not an implied architecture decision. |

## Area-Specific Reasoning Inputs

| Area | Reasoning Document | Key Conclusion | Evidence Basis | Impact On Final Decision |
| --- | --- | --- | --- | --- |
| API Design | `projects/ultraplan-go/sprints/30-web-foundations/reasoning/api-design.md` | Use synchronous `GET`/`HEAD` HTML and `/api/v1` resources over live typed app reads; use explicit transport DTOs/envelopes, opaque artifact references, fixed limits, safe typed-error projection, and no operation routes. | Handbook reports `01`-`08`, `10`-`14`; architecture and Sprint 30 requirements; explicit route, envelope, error, health, and limit trade-offs. | Fixes route/method semantics, API compatibility boundary, error codes, health meaning, collection/preview limits, and unknown `/api/` behavior. |
| Architecture | `projects/ultraplan-go/sprints/30-web-foundations/reasoning/architecture.md` | Refactor composition so `cmd/ultraplan` explicitly and lazily constructs TUI/web runners; keep a cohesive query-only app facade; make `internal/web` a transport adapter with no product state or broad filesystem/runtime capability. | Handbook evidence for thin entrypoints, constructor injection, capability boundaries, source-aware config, cancellation, bounded lifecycle, and real command-path tests. | Fixes package direction, lifecycle ownership, state ownership, artifact-resolution boundary, browser-launch behavior, and CLI/TUI non-regression constraints. |
| Frontend | `projects/ultraplan-go/sprints/30-web-foundations/reasoning/frontend.md` | Deliver complete server-rendered pages with progressive enhancement only, escaped Markdown/JSON source previews, semantic/accessibility behavior, responsive reflow, embedded assets, and no client state/cache/framework. | Handbook transport/security/performance/testing evidence plus PRD/TRD server-rendered single-binary constraints. | Fixes page structure and states, hostile-content treatment, no-JavaScript behavior, accessibility checks, responsive requirements, and frontend dependency budget. |

All three area conclusions are adopted. There is no area-level conflict to override. The broader PRD/TRD references to confirmation, operations, SSE, and cancellation are not adopted because `requirements.md`, `sprint-index.md`, and every area record explicitly defer them from Sprint 30.

## Sprint Technical Handbook Summary

- **Relevant Patterns:** Thin transports over typed app behavior; a visible process composition root; narrow constructor injection; read-only as absence of dangerous capabilities; source-aware config merge followed by validation; typed errors with boundary-specific rendering; root/request context propagation; sequential aggregation before measured concurrency; stable redacted diagnostics separate from user output; deterministic fakes plus real command/server-path tests; lazy serve-only initialization; bounded/incremental reads.
- **Important Trade-Offs:** A cohesive query facade adds mapping but limits capability; central composition improves lifecycle visibility but must stay thin; live reads maximize freshness but repeat work; sequential aggregation simplifies consistency but may increase latency; a versioned JSON API creates compatibility obligations; immutable startup config requires restart; complete representation goldens reveal drift but can become brittle.
- **Warnings / Anti-Patterns:** No transport-owned workflow/validation logic, mutable global registries, context service locators, raw/string-parsed errors, mixed injected/direct I/O, UI-only read-only enforcement, trusted-local-input assumptions, ownerless goroutines, whole-workspace materialization, handler-only testing, accidental plugin/public API contracts, or unredacted diagnostics mixed into responses.
- **Evidence Confidence:** High for package direction, explicit dependency ownership, error/cancellation/config patterns, capability restriction, bounded work, observability, and testing because the selected reports cite repeated concrete examples. Exact HTTP routes, constants, visual details, and loopback policy are Sprint 30 decisions constrained by project requirements rather than copied from studied repositories.

## Contracts Applied

| Contract / Requirement ID | Constraint | Decision Impact | Expected Evidence |
| --- | --- | --- | --- |
| Architecture; `RO-4`-`RO-15`; `AC-7`, `AC-12`, `AC-13`; `C-1`-`C-4` | Explicit composition, inward dependencies, transport-only web ownership, authoritative product state. | `cmd` constructs runners; web receives a query-only app facade and owns no product state or workflow logic. | Import review, composition tests, package-cycle-free `go test ./...`, Architecture Review. |
| Errors; `AC-1`, `AC-9`, `AC-11`; `C-2`, `C-9` | Preserve typed causes and provide safe actionable boundary mappings. | CLI keeps existing exit classes; API uses fixed status/code envelopes; HTML uses safe error pages; no raw causes or string parsing. | Error-mapping tables, leak assertions, help/startup tests, route tests. |
| Configuration; `AC-1`-`AC-4`; `C-3`, `C-6` | Existing precedence/discovery rules and pre-listen validation. | `--workspace` remains shared; `--listen` and `--open-browser` are explicit; effective config is immutable after startup; non-loopback is rejected before runner/listener use. | Explicit-set precedence tests, discovery tests, runner-not-called test on invalid bind. |
| Observability; `AC-6`, `AC-11`; `C-2`, `C-4` | Truthful status/health and redacted lifecycle/request diagnostics. | Cheap workspace-aware readiness; stable safe diagnostic fields on a separate channel; request IDs are generated by the server. | Health `200`/`503` tests, structured log capture/redaction checks, shutdown diagnostics. |
| Security; `AC-3`, `AC-10`, `AC-11`; `C-5`-`C-9` | Loopback, same-origin, bounded input, containment, safe rendering, no dangerous capability. | IP-literal loopback binding, strict Host/Origin checks, no CORS, security headers, opaque refs, bounded Markdown/JSON source previews, no operation routes. | Host/Origin/path/symlink/body/target/hostile-content/security-header tests. |
| Testing; `RO-8`, `RO-16`-`RO-20`; `AC-14`; `C-10` | Deterministic fakes, `httptest`, lifecycle/route/security/template/artifact coverage, race/build checks. | Focused semantic tests plus narrow normalized envelope goldens and actual command/listener paths. | Required focused tests and full commands listed in Expected Evidence. |
| Documentation; `RO-21`-`RO-23`; `AC-1`, `AC-8` | User-facing help and docs accurately state local/read-only behavior and architecture. | Add local-web guide and update CLI/architecture docs; explain bounds, previews, shutdown, browser-opening warning, and deferred operations. | Documentation review and `ultraplan serve --help` snapshot/semantic assertions. |
| CLI Surface; `RO-4`, `RO-6`, `RO-8`; `AC-1`-`AC-4`, `AC-13` | Stable command shape, workspace behavior, output discipline, exit mapping, no regression. | `serve` delegates through an injected runner; launcher warnings use diagnostics; existing commands/TUI stay unchanged. | Command-path tests for help, dispatch, output, cancellation, exit classes, CLI/TUI regression. |
| Performance; `RO-9`, `RO-13`; `AC-4`, `AC-10`, `AC-14`; `C-10` | Bound time, memory, concurrency, and scans; avoid costs for unrelated commands. | Lazy web construction, fixed HTTP limits/timeouts/concurrency, bounded collections/previews, sequential reads, no cache/polling. | Boundary tests, race tests, truncation assertions, review confirming no recursive source scan or unbounded goroutine. |

## Repos Studied / Source Evidence Used

The repositories below were not adopted wholesale. Their cited source examples, as distilled by `technical-handbook.md` from the selected reports, materially informed the decisions.

| Source / Repo / Report | Concrete Reference | Relevant Finding | Why It Matters For This Sprint | Used In Decision(s) |
| --- | --- | --- | --- | --- |
| gdu via `01-project-structure` | `gdu/cmd/gdu/app/app.go:30-49` | Multiple presentation adapters share a boundary rather than duplicate behavior. | Supports web as another app-backed local interface. | 1, 3, 6 |
| Helm via `02-command-architecture` | `helm/pkg/cmd/install.go:132-145` | Commands configure and delegate to reusable behavior. | Supports thin `serve` parsing/delegation rather than HTTP lifecycle in command dispatch. | 1, 2 |
| opencode via `03-dependency-injection` | `opencode/internal/app/app.go:42-81` | Visible application assembly clarifies dependency ownership. | Supports explicit runner construction and lifecycle ownership. | 1, 2 |
| gh-cli via `03-dependency-injection`, `12-extensibility`, `14-performance` | `gh-cli/pkg/cmd/factory/default.go:26-46`; `gh-cli/pkg/cmdutil/factory.go:16-43` | Narrow factories preserve test seams and can delay expensive resources. | Supports injected opener/listener seams and lazy serve-only setup without a registry. | 1, 2, 7 |
| chezmoi and restic via `04-configuration-management` | `chezmoi/internal/cmd/config.go:2253-2287`; `restic/internal/global/global.go:139,147` | Explicitly set flags must override merged config without defaults doing so accidentally. | Grounds `serve` option precedence and pre-listen validation. | 2 |
| go-task, restic, and gh-cli via `05-error-handling` | `go-task/errors/errors.go:47-50`; `restic/cmd/restic/main.go:199-209`; `gh-cli/internal/ghcmd/cmd.go:281-301` | Lower layers retain typed identity while outer boundaries classify and render. | Grounds separate CLI, HTML, and API error projections. | 3, 5 |
| restic via `06-io-abstraction` | `restic/internal/backend/backend.go:19-90`; `restic/internal/backend/mem/mem_backend.go:52-255` | Narrow capabilities permit deterministic substitutes and constrain I/O. | Supports query-only app access, app-owned artifact resolution, and fakes. | 1, 4, 8 |
| Helm via `07-state-context` | `helm/pkg/cmd/install.go:333-347`; `helm/pkg/action/install.go:284` | Signal cancellation should reach blocking work. | Grounds process-root context, request cancellation, and graceful shutdown. | 2, 7 |
| go-task, gdu, k9s, and opencode via `08-concurrency`, `14-performance` | `go-task/task.go:87,197`; `gdu/pkg/analyze/parallel.go:13,36`; `k9s/internal/pool.go:21-37`; `opencode/cmd/root.go:261-279` | Concurrency needs bounds, ownership, cancellation, and bounded cleanup; sequential work is a valid baseline. | Grounds sequential query aggregation, request concurrency limits, and joined shutdown. | 2, 7 |
| Helm and k9s via `10-logging-observability` | `helm/internal/logging/logging.go:31-71`; `k9s/internal/slogs/keys.go:6-231` | Diagnostics use stable fields and a channel separate from user output. | Grounds request/lifecycle logging and CLI output discipline. | 2, 5 |
| chezmoi, go-task, and Helm via `11-testing-strategy` | `chezmoi/internal/cmd/main_test.go:64-174`; `go-task/task_test.go:166-169`; `helm/internal/test/test.go:43` | Focused fakes should be combined with command paths and selective normalized goldens. | Grounds command, listener, handler, template, and envelope verification. | 8 |
| k9s, opencode, and Helm via `13-security` | `k9s/internal/config/json/validator.go:146`; `opencode/internal/options/secret_string.go:15-20`; `helm/pkg/registry/transport.go:37-41` | Treat inputs as hostile, centralize validation, and redact sensitive details. | Grounds Host/Origin/input/path checks and safe diagnostics/responses. | 4, 5, 6 |
| Helm and yq via `14-performance` | `helm/pkg/action/lazyclient.go:35-53`; `yq/pkg/yqlib/stream_evaluator.go:78-113` | Initialize expensive facilities lazily and process data incrementally. | Grounds lazy web startup and byte-bounded artifact reads. | 2, 4, 7 |
| `technical-handbook.md` | `Relevant Patterns`, `Trade-Offs`, `Anti-Patterns And Warnings`, `Design Pressures` | Consolidates all selected report findings and cautions for Sprint 30. | It is the selected evidence synthesis and prevents copying repository patterns without project-specific fit analysis. | 1-8 |

## Trade-Off And Debt Analysis

### Accepted Trade-Offs

| Trade-Off | Benefit | Cost / Constraint Accepted | Why Acceptable Now | Revisit Trigger |
| --- | --- | --- | --- | --- |
| Explicit `cmd` composition and lazy web construction | Avoids cycles/globals and protects unrelated commands from web cost/failure. | Adds runner contracts and more visible wiring; composition root can grow. | A third surface requires explicit lifecycle ownership, but no generic container is needed. | `main.go` begins owning behavior rather than construction, or a fourth side-effectful surface appears. |
| One cohesive query facade | Enforces read-only capability with manageable fake/setup cost. | Several query methods share one interface and may change together. | All methods serve one bounded inspection surface and expose no commands. | Independent consumers need materially different subsets, or mutation work begins. |
| Live per-request reads with `no-store` | Keeps workspace/product state authoritative and restart semantics trivial. | Repeated refresh repeats filesystem work and aggregate responses are not cross-request snapshots. | Sprint 30 has explicit refresh, fixed bounds, and no measured cache need. | Measured latency/load exceeds targets or repeated scans become visible. |
| Sequential aggregation and fail-whole-query errors | Deterministic ordering and simple cancellation/error semantics. | One failed aggregate section can suppress otherwise useful content; latency is additive. | No evidence yet justifies partial-result schemas or fan-out complexity. | Measurements show independent slow panels or users require partial availability. |
| Compatibility-controlled internal `/api/v1` | Gives bundled JavaScript stable envelopes without claiming a public platform. | DTO evolution requires compatibility review and fixtures. | The bundled UI still needs an explicit transport boundary. | External consumers are intentionally supported or a breaking field change is needed. |
| Fixed limits without pagination | Prevents unbounded scans/responses and keeps the initial UI/API small. | Large workspaces show truncated collections and require narrower detail navigation. | Core inspection scenarios fit bounded summaries. | Real workspaces routinely exceed 200 useful entries or need arbitrary-depth browsing. |
| Escaped Markdown/JSON source preview | Makes hostile embedded HTML/scripts inert and preserves auditability. | Less polished than rich Markdown rendering. | Inspection and safety matter more than rich rendering in the foundation. | A sanitized rich-preview requirement is selected with parser/sanitizer evidence. |
| Immutable server configuration | Predictable concurrency and complete pre-listen validation. | Changes require process restart. | Hot reload adds watchers, precedence, and race semantics with no Sprint 30 need. | Long-lived operation mode requires governed runtime config changes. |
| IP-literal loopback binding and exact origin | Avoids DNS rebinding/hostname ambiguity. | `localhost` is not accepted and IPv4/IPv6 origins are not interchangeable. | Explicit local security is preferable to convenience aliases. | A governed hostname/certificate/remote-access model is introduced. |
| Optional launcher failure is warning-only | A usable server remains successful when browser opening fails. | User may need to open the printed URL manually. | Browser opening is convenience, not server readiness. | Product requirements make launch success mandatory. |

### Potential Technical Debt

| Debt / Shortcut | Why It Might Accrue | Current Mitigation | Owner / Follow-Up |
| --- | --- | --- | --- |
| Fixed `200` item and `256 KiB` preview limits | Defaults are safety-driven rather than production-measured. | Return explicit count/byte/truncation metadata and document the limits. | Reassess from measured Sprint 30 use before pagination/configurability. |
| Sequential dashboard reads | More panels or large state files can increase latency. | Bound each result, propagate cancellation, avoid recursive source scans, instrument duration. | Add only measured, bounded app-level fan-out with an explicit partial-error contract. |
| Cohesive facade growth | Future web operations could be appended for convenience. | Name and keep the Sprint 30 surface query-only; inject no runtime/process/mutation dependency. | Sprint 31 must define a separate command/operation capability. |
| Internal `v1` compatibility burden | Bundled client and tests may freeze accidental fields. | Explicit web DTOs, narrow envelope goldens, semantic tests, and documented compatibility scope. | New version or coordinated migration decision for breaking changes. |
| Source-only Markdown preview | User expectations may shift toward rendered Markdown. | Label previews clearly and test hostile content remains inert. | Later frontend/security reasoning for parser and sanitizer selection. |
| Composition-root growth | More runners can make startup wiring hard to read. | Keep constructors small and `cmd` behavior-free; no service locator. | Architecture review when construction can no longer be understood locally. |
| No browser-level automation in the minimum suite | `httptest` and semantic HTML tests may miss browser-specific layout/focus issues. | Manual responsive/keyboard review is required; browser tests may be added without a frontend toolchain if practical. | Sprint 32 hardening or any observed browser/accessibility regression. |

### Future Considerations

| Consideration | Deferred Until | Reason Deferred | What Should Be Preserved Now |
| --- | --- | --- | --- |
| Confirmed operations, validation execution, cancellation, and SSE | Sprint 31 | Explicit Sprint 30 non-goal and materially larger security/lifecycle model. | Separate query facade, typed app results, request context propagation, stable error envelope, explicit composition. |
| Authentication, remote/LAN bind, TLS, users/tenants | A later governed remote-service phase | Loopback single-user trust model cannot safely evolve implicitly into remote access. | Reject all non-loopback binds and avoid claims that loopback equals tenant isolation. |
| Caching/snapshots/watchers | Measured read amplification | Adds invalidation, staleness, locking, and restart semantics. | `no-store`, authoritative-state language, bounded query seams, request duration fields. |
| Pagination/search/filtering | Repeated truncation in representative workspaces | Initial bounded navigation does not require arbitrary collection traversal. | Stable truncation/count metadata and canonical identifiers. |
| Rich Markdown rendering | Explicit frontend/security sprint | Requires a maintained parser/sanitizer and new hostile-content contract. | Keep source content plain and never introduce `template.HTML`/`innerHTML`. |
| Client framework/build pipeline | Demonstrated client interaction complexity | Current pages are read-only and fully server-renderable. | Versioned API DTOs, semantic HTML, progressive enhancement, no browser-owned product state. |
| Partial dashboard results/concurrent aggregation | Measured independent panel latency or availability need | Requires per-section errors, ordering, cancellation, and goroutine policy. | Keep aggregation app-owned and handlers thin. |
| Public local API guarantees | Explicit product/API commitment | Auth, release policy, pagination, and broad client compatibility are absent. | Document `/api/v1` as bundled-client/internal and use explicit DTOs. |

## Decisions

Sprint 30 adopts the eight final decisions below as the binding implementation direction. Together they fix explicit CLI/TUI/web composition, a query-only app boundary, loopback server lifecycle and limits, read-only HTML and `/api/v1` contracts, opaque bounded artifact previews, HTTP security and error projection, server-rendered accessible presentation, live authoritative state reads, and layered verification. `plan.md` must execute these decisions without reopening architecture or adding the deferred operation, SSE, remote-service, browser-state, frontend-framework, or Git-mutation scope.

## Final Decisions

### Decision 1: Explicit Surface Composition And Query-Only App Boundary

- **Decision:** `cmd/ultraplan` will remain the process composition root and explicitly construct `internal/app`, the existing TUI runner, and a lazy web runner. `internal/app/surfaces.go` will define narrow injected runner contracts; `internal/app/web_usecases.go` will expose one cohesive query-only facade with typed dashboard, project, sprint, study, validation-summary, health, and artifact-preview results. `internal/web` may import this app surface and plain results only. It will not import product modules, workspace internals, runtime/process adapters, or CLI handlers. No global registration, service locator, `init` callback, or context-carried dependency is permitted.
- **Rationale:** This is the smallest refactor that prevents `app -> web -> app` cycles, makes lifecycle/test ownership visible, structurally denies mutation capabilities, and preserves product-module ownership.
- **Study / Source Grounding:** `technical-handbook.md` patterns “Thin transport over typed application behavior,” “Explicit composition and narrow injection,” and “Read-only as a capability boundary”; gdu `app.go:30-49`, Helm `install.go:132-145`, opencode `app.go:42-81`, gh-cli `factory.go:16-43`, and restic backend boundaries via reports `01`, `02`, `03`, `06`, `12`, and `13`.
- **Trade-Offs Accepted:** Visible runner/facade wiring and some mapping ceremony are accepted to avoid hidden globals and broad product injection. The composition root must stay construction-only.
- **Technical Debt / Future Impact:** Facade creep is a risk. Sprint 31 must add a separately reviewed command/operation capability rather than append mutation methods to the read facade.
- **Alternatives Rejected:** Constructing web inside app is rejected because it creates a cycle; package-global runner registration is rejected because it hides mutable ownership; web-to-product imports are rejected because they duplicate application behavior; one interface per route is rejected as premature fragmentation; injecting all app/product services is rejected as broad capability coupling.
- **Contracts Satisfied:** Architecture, Security, Testing; `RO-4`, `RO-5`, `RO-6`, `RO-7`; `AC-7`, `AC-12`, `AC-13`; `C-1`-`C-5`.
- **Evidence Required:** Import inspection; composition tests proving independent injected runners and no globals; fake app-query tests; command-path non-regression tests; Architecture Review confirming ownership and dependency direction; package-cycle-free `go test ./...`.

### Decision 2: Serve Command, Immutable Configuration, And Owned Lifecycle

- **Decision:** Add `ultraplan serve` with `--listen <ip:port>` defaulting to `127.0.0.1:8080`, `--open-browser` defaulting to false, and the existing global `--workspace`. Accept only numeric IP literals for which Go reports loopback, with valid explicit ports; bracketed IPv6 such as `[::1]:8080` is valid and hostname aliases such as `localhost` are rejected. Existing config precedence remains built-in defaults, workspace config, environment, then explicitly set flags. App preflight resolves/validates workspace and immutable effective serve config before invoking the runner; bind policy is rechecked before listener acquisition. The server uses 5-second header-read, 15-second request-read, 30-second write, 60-second idle, and 10-second graceful-shutdown limits, with at most 32 in-flight requests. Root and request contexts propagate cancellation; all server goroutines are owned and completed/cancelled before return. Browser opening uses an injected launcher after successful listen; failure is a redacted warning and not fatal.
- **Rationale:** A long-lived command needs explicit ownership, deterministic preflight, bounded cleanup, and isolation from normal CLI/TUI startup. Numeric loopback-only input avoids DNS ambiguity.
- **Study / Source Grounding:** `technical-handbook.md` configuration, state/context, concurrency, observability, testing, and lazy-resource patterns; chezmoi/restic precedence references, Helm cancellation references, opencode/k9s bounded-shutdown references, and gh-cli/Helm lazy construction via reports `03`, `04`, `07`, `08`, `10`, `11`, and `14`.
- **Trade-Offs Accepted:** Restart is required for config changes; fixed bounds may need tuning; `localhost` convenience is sacrificed; launcher failure does not fail a healthy server.
- **Technical Debt / Future Impact:** Timeout/concurrency defaults need measurement. Any future config fields must preserve explicit-set tracking and validation-before-listen semantics.
- **Alternatives Rejected:** Non-loopback or hostname binding is rejected by scope/security; hot reload is rejected as unneeded watcher/race complexity; eager construction is rejected because unrelated commands must not pay web costs; detached/background lifecycle is rejected because cancellation and cleanup become ownerless; fatal launcher failure is rejected because opening is optional.
- **Contracts Satisfied:** Architecture, Configuration, Errors, Observability, Security, CLI Surface, Performance, Testing; `RO-4`, `RO-6`, `RO-8`, `RO-9`; `AC-1`-`AC-5`, `AC-13`, `AC-14`; `C-3`, `C-4`, `C-6`, `C-10`.
- **Evidence Required:** Help assertions; precedence/discovery tests; table tests for IPv4/IPv6 loopback and rejected hostnames/non-loopback/malformed addresses; proof runner/listener is not called after failed preflight; startup/listen failure and cancellation tests; launcher warning test; graceful-shutdown/handler completion tests; race test and CLI/TUI dispatch regression tests.

### Decision 3: Read-Only HTML And Versioned JSON Route Contract

- **Decision:** Register the exact `GET` resources decided in `reasoning/api-design.md`: `/`, `/projects`, `/projects/{project}`, `/projects/{project}/sprints/{sprint}`, `/studies`, `/studies/{study}`, `/artifacts/{ref}`, `/api/v1/dashboard`, `/api/v1/projects`, `/api/v1/projects/{project}`, `/api/v1/projects/{project}/sprints/{sprint}`, `/api/v1/studies`, `/api/v1/studies/{study}`, `/api/v1/validations?scope=...&ref=...`, `/api/v1/artifacts/{ref}`, and `/api/v1/health`, plus embedded `/static/` assets. Successful `GET` resources support implicit `HEAD`; all other methods return `405` and `Allow: GET, HEAD`. Unknown `/api/` paths and unsupported versions always use JSON errors; unknown browser routes use the HTML error page. JSON success is `{data,meta}` with `api_version`, server-generated `request_id`, and `generated_at`; JSON errors are `{error:{code,message,details?},meta}`. Route-specific web DTOs, not app/product structs, define stable field names. `/api/v1` is compatibility-controlled for the bundled browser but is not promised as a public integration API.
- **Rationale:** Resource-oriented `GET` routes encode read-only/idempotent behavior, preserve separate HTML/JSON representations over one app result, and provide a deliberate compatibility seam without creating an extension platform.
- **Study / Source Grounding:** `technical-handbook.md` thin transport, boundary-error, extensibility, and testing patterns; gdu, Helm, restic, go-task, gh-cli, and normalized-golden examples via reports `01`, `02`, `03`, `05`, `11`, and `12`.
- **Trade-Offs Accepted:** Explicit DTO mapping and `v1` compatibility maintenance are accepted. Aggregate query failure remains whole-response rather than partial-panel success.
- **Technical Debt / Future Impact:** Breaking DTO changes require a new version or coordinated migration decision. Public client support, caller pagination, operations, and SSE need separate future decisions.
- **Alternatives Rejected:** RPC-style `/api/query` is rejected because it obscures resources/methods; direct serialization is rejected for coupling/disclosure risk; HTML fallback for unknown API paths is rejected as unstable machine behavior; `POST /validations`, operations, confirmations, SSE, and cancellation routes are rejected by Sprint 30 scope.
- **Contracts Satisfied:** Architecture, Errors, Observability, Security, Testing, Documentation; `RO-10`, `RO-11`, `RO-12`; `AC-6`, `AC-7`, `AC-9`, `AC-12`; `C-1`, `C-2`, `C-5`, `C-9`.
- **Evidence Required:** Table-driven route/method/HEAD/Allow tests; unknown `/api/` and unsupported-version JSON tests; semantic envelope/field/null-omission/content-type/cache tests; normalized goldens only for stable success/error envelope shapes; assertions that no operation/mutation/SSE routes are registered.

### Decision 4: Bounded Opaque Artifact Preview

- **Decision:** Listing/detail app queries issue opaque, URL-safe artifact references; HTTP never accepts absolute or workspace-relative file paths. The app use case resolves a reference against canonical workspace state, rejects stale/forged/escaping/symlink-escaping/non-allowlisted targets, permits only governed Markdown and JSON artifact classes, and reads no more than 256 KiB plus one byte needed to detect truncation. It returns safe workspace-relative display metadata and bounded bytes. `internal/web/artifacts.go` validates the returned media/size contract and renders or maps it without direct filesystem access. API output returns source text, media type, total/returned bytes, and truncation; JSON may include a parsed value only within the same bound. HTML renders Markdown and JSON as escaped source in `<pre><code>`; no workspace content becomes `template.HTML`, `innerHTML`, executable links, or scripts. All rejected references collapse to safe `404 not_found`.
- **Rationale:** Opaque app-issued references and app-owned containment keep workspace knowledge out of HTTP and deny arbitrary filesystem capability. Escaped bounded source is auditable and makes hostile content inert.
- **Study / Source Grounding:** `technical-handbook.md` narrow I/O/capability, hostile-input, bounded/incremental work, and security evidence; restic I/O interfaces, k9s validation, opencode redaction, Helm credential scrubbing, and yq incremental processing via reports `06`, `13`, and `14`.
- **Trade-Offs Accepted:** Rich Markdown rendering and arbitrary path convenience are sacrificed; previews truncate at a fixed safety threshold.
- **Technical Debt / Future Impact:** Users may request rendered Markdown or larger previews. Either requires a new security/performance decision, not silent relaxation. Opaque references may become stale and intentionally return 404.
- **Alternatives Rejected:** Caller-provided cleaned paths are rejected because lexical cleaning does not solve encoding/symlink/disclosure risks; direct web filesystem access is rejected by dependency/capability rules; whole-file reads are rejected as unbounded; rich rendering is rejected because parser/sanitizer behavior is not selected.
- **Contracts Satisfied:** Architecture, Security, Errors, Performance, Testing; `RO-13`, `RO-14`, `RO-18`; `AC-6`, `AC-10`-`AC-12`; `C-1`, `C-2`, `C-5`, `C-7`, `C-10`.
- **Evidence Required:** Opaque-reference round trips; stale/forged reference tests; traversal, encoded traversal, absolute path, separator/control character, symlink escape, extension allowlist, exact/over-limit, invalid JSON, hostile Markdown/HTML/script, escaped output, truncation metadata, `nosniff`, and no-direct-web-filesystem import review.

### Decision 5: Local HTTP Security, Error Projection, And Diagnostics

- **Decision:** Enforce loopback-only listening, exact Host authority for the effective listener, same-origin pages/API, no permissive CORS, and exact Origin matching when Origin is present. Allow absent Origin for top-level navigation and local non-browser `GET`/`HEAD` after Host validation; reject `Origin: null`, malformed, cross-origin, and non-loopback origins. Reject bodies on Sprint 30 routes; return `413 request_too_large` for declared bodies over 64 KiB and `400 invalid_request` for other non-empty bodies. Limit request targets to 8 KiB and decoded identifiers/references to 128 bytes; reject control characters, separators, malformed/duplicate identifiers, duplicate or unknown query parameters. Apply CSP restricted to self-hosted embedded assets, `X-Content-Type-Options: nosniff`, `Referrer-Policy: no-referrer`, frame denial, and `Cache-Control: no-store`. Map typed app errors to only `invalid_request` 400, `request_rejected` 403, `not_found` 404, `method_not_allowed` 405, `request_too_large` 413, `internal_error` 500, and `unavailable` 503. Generate request IDs server-side. Structured diagnostics record safe normalized route/method/status/duration/bytes/code and lifecycle/shutdown fields on the diagnostic channel, never raw URLs/query values, artifact content, absolute paths, Host/Origin values, environment, secrets, or raw errors at normal level.
- **Rationale:** Loopback reduces exposure but does not make browser input or local processes trusted. Fixed validation, disclosure, and diagnostic boundaries are required even for read-only requests.
- **Study / Source Grounding:** `technical-handbook.md` error, observability, and security patterns and warnings; go-task/restic/gh-cli typed errors, Helm/k9s diagnostic separation, k9s validation, opencode secret wrappers, and Helm scrubbing via reports `05`, `10`, and `13`.
- **Trade-Offs Accepted:** Strict Host/Origin behavior may reject unusual proxies or aliases; rich debug information is kept out of normal responses/logs; no cookie/session/auth token is introduced for this read-only loopback sprint.
- **Technical Debt / Future Impact:** Sprint 31 mutations will require CSRF/session/confirmation decisions beyond this read-only Origin policy. This model is not an isolation guarantee against a compromised local account.
- **Alternatives Rejected:** Trusting loopback alone is rejected; permissive CORS is rejected; reflecting caller request IDs or errors is rejected; string-parsing errors is rejected; authentication/tenant systems are rejected as out of scope; exposing path existence through differentiated preview errors is rejected.
- **Contracts Satisfied:** Security, Errors, Observability, Performance, Testing; `RO-12`, `RO-14`, `RO-17`; `AC-3`, `AC-9`-`AC-11`; `C-5`-`C-10`.
- **Evidence Required:** Host/port/IPv6 and Origin matrix tests; absent-Origin and no-CORS tests; body/target/query/identifier boundary tests; all status/code mappings; leak/redaction assertions; security-header tests for HTML/JSON/static/error responses; request-ID generation and structured diagnostic capture; race tests.

### Decision 6: Server-Rendered Accessible Progressive Frontend

- **Decision:** Embed and parse `html/template` pages plus one small stylesheet and at most one small dependency-free external script. Parse templates once at serve startup and fail startup on parse errors; buffer page rendering before writing headers. Use a shared semantic shell with skip link, header, primary navigation, workspace context, `main`, footer/server status, unique title and `h1`, breadcrumbs, and earned partials for repeated statuses/findings/metadata/empty states. Initial HTML must fully support dashboard, project/sprint, study, validation, artifact, health, and error inspection without JavaScript. JavaScript may only perform explicit-user-triggered refresh of marked regions with stale/failure messaging and ordinary reload fallback; there is no polling, client router/store/cache, local storage, cookies, service worker, operation UI, or SSE. Use textual status labels, visible focus, native controls, semantic tables/definition lists, polite refresh status, responsive single-column reflow, local overflow for unavoidable tables/code, 200-percent zoom usability, and reduced-motion behavior.
- **Rationale:** Server rendering and progressive enhancement meet the single-binary/local-first requirement, minimize state and attack surface, and provide robust no-JavaScript and accessible inspection.
- **Study / Source Grounding:** `technical-handbook.md` transport separation, capability restriction, bounded work, diagnostics, and deterministic testing via reports `01`, `02`, `03`, `06`, `08`, `10`, `11`, `13`, and `14`. Exact accessibility and responsive rules come from `reasoning/frontend.md` and project requirements; no comparative repo source specifically determined typography or breakpoints because those are product presentation constraints, not study-derived architecture.
- **Trade-Offs Accepted:** Less app-like interactivity and source-only Markdown are accepted for simpler security, state, accessibility, and distribution. Buffered rendering delays streaming first bytes for bounded pages.
- **Technical Debt / Future Impact:** Browser automation may be expanded during hardening; framework/client-state adoption remains prohibited until selected by later reasoning. Template partials must not become a generic component framework.
- **Alternatives Rejected:** React/Vue/Vite or any Node build is rejected; JavaScript-owned navigation is rejected; timed polling is rejected; custom widgets are rejected; third-party assets/fonts/analytics are rejected; inline scripts and trusted workspace markup are rejected; card-only or unbounded-table layouts are rejected for operational density/mobile behavior.
- **Contracts Satisfied:** Architecture, Security, Performance, Testing, Documentation; `RO-15`, `RO-19`; `AC-6`, `AC-8`, `AC-10`-`AC-12`; `C-2`, `C-5`, `C-7`, `C-8`, `C-10`.
- **Evidence Required:** Embedded template parse test; representative page/empty/error/truncation fixtures; hostile string/URL/Markdown/JSON escaping assertions; semantic checks for title/headings/landmarks/skip link/navigation/table labels/status text/live regions; no-JavaScript navigation review; narrow/desktop/zoom/overflow and keyboard/focus review; checks for no inline/third-party assets or frontend build dependency.

### Decision 7: Live, Bounded, Sequential State Projection And Truthful Health

- **Decision:** Every page/API request reads current authoritative workspace/product state through app queries. Do not add caches, watchers, snapshots, browser persistence, hidden preloads, recursive source-repository scans, background refresh, or handler-owned fan-out. Collection and finding responses return at most 200 entries and include `returned_count`, `total_count`, and `truncated`; no caller pagination is added. Requests propagate cancellation and operate under Decision 2 limits. Aggregate app queries are sequential and fail coherently; no partial panel error schema is introduced. `/api/v1/health` performs only cheap server readiness and lightweight configured-workspace query availability, returning 200 `ok` or 503 `unavailable`; it does not scan projects/studies, validate artifacts, check providers/runtimes, or invoke review/smoke.
- **Rationale:** Live bounded reads preserve source-of-truth semantics and simple recovery while avoiding cache invalidation and speculative concurrency. Health must report readiness, not overclaim product validity.
- **Study / Source Grounding:** `technical-handbook.md` state ownership, sequential baseline, bounded concurrency, lazy/incremental work, and whole-workspace warning; Helm context references, go-task/gdu/k9s bounds, and yq incremental processing via reports `07`, `08`, and `14`.
- **Trade-Offs Accepted:** Refresh repeats work, responses are point-in-time rather than cross-request snapshots, and one failed aggregate query fails the aggregate view. Truncation can hide older entries.
- **Technical Debt / Future Impact:** Add caching, pagination, concurrency, or partial results only from measurement and with explicit staleness/error contracts. Preserve count/truncation fields for that evolution.
- **Alternatives Rejected:** Web caches/snapshots are rejected due to invalidation and duplicate state; timed polling is rejected due to hidden load; unbounded responses are rejected; speculative fan-out is rejected due to goroutine/partial-result semantics; deep health is rejected as expensive and misleading.
- **Contracts Satisfied:** Architecture, Observability, Performance, Testing; `AC-4`, `AC-6`, `AC-14`; `C-2`, `C-4`, `C-10`.
- **Evidence Required:** Repeated-read non-mutation tests; collection/truncation tests at 0, exact, and over-limit sizes; cancellation tests; health 200/503 and no-deep-scan fake assertions; route duration diagnostics; review for no cache/watcher/polling/fan-out/recursive scan.

### Decision 8: Layered Verification, Documentation, And Existing-Surface Protection

- **Decision:** Verification will combine focused app fakes, web handler/template fakes, `httptest`, real listener/lifecycle tests, real serve command paths, temporary-workspace integration, semantic response assertions, narrow normalized envelope goldens, race tests, and a binary build. Normal tests require no browser process, OpenCode/runtime, smoke harness, Node.js, frontend server, database, or network service. Update `docs/local-web.md`, `docs/cli-reference.md`, and `docs/architecture.md` in the implementation repository to document command flags, loopback/read-only scope, state authority, limits, source-preview behavior, optional launcher warning, shutdown/troubleshooting, `/api/v1` bundled-client status, and `internal/web -> internal/app`. Existing CLI/TUI behavior and output remain unchanged except explicit runner composition.
- **Rationale:** Handler-only tests miss composition, listener, shutdown, and output regressions. Documentation and help are part of the surface contract, and the composition refactor is the primary regression risk.
- **Study / Source Grounding:** `technical-handbook.md` deterministic-boundary plus real-stack verification pattern; restic fakes, chezmoi command-path tests, go-task normalized fixtures, and Helm golden helpers via reports `06` and `11`.
- **Trade-Offs Accepted:** The suite is broader than handler unit tests and selective goldens require deliberate updates. Browser-level automation is not mandatory in this foundation if semantic/manual accessibility checks provide reviewable evidence.
- **Technical Debt / Future Impact:** Sprint 32 hardening may add browser automation and measured performance tests. Test seams must remain behavior-focused rather than proliferating interfaces.
- **Alternatives Rejected:** Handler-only testing is rejected; all-page goldens are rejected as brittle; real runtime/smoke dependencies are rejected as unrelated and nondeterministic; manual-only verification is rejected; changing existing CLI/TUI output for composition convenience is rejected.
- **Contracts Satisfied:** Testing, Documentation, CLI Surface, Architecture, Performance; `RO-8`, `RO-16`-`RO-23`; `AC-1`, `AC-13`, `AC-14`; `C-10`.
- **Evidence Required:** Focused app/web suites; `go test ./...`; `go test -race ./...`; `go build ./cmd/ultraplan`; help output checks; import/dependency review; desktop/narrow/keyboard/no-JavaScript review; completed Architecture Review and Sprint Review protocols; documentation path/content review.

## Expected Evidence

| Evidence Type | Required Evidence | Source / Command / Review Check |
| --- | --- | --- |
| Tests | Full deterministic suite passes. | From `../ultraplan-go`: `go test ./...` |
| Race | No races or lifecycle leaks detected by the race-enabled suite. | From `../ultraplan-go`: `go test -race ./...` |
| Build | Single Go binary builds with embedded templates/assets and no frontend build step. | From `../ultraplan-go`: `go build ./cmd/ultraplan` |
| App Command | Serve help, precedence, workspace resolution, invalid-bind preflight, runner invocation/failure, launcher warning, cancellation, exit mapping, and CLI/TUI non-regression. | Focused `go test ./internal/app -run 'Serve|Surface|Dispatch|TUI'` or the implementation's equivalent focused pattern. |
| Web Lifecycle | Loopback validation, actual listener startup/failure, context cancellation, timeout-bounded graceful shutdown, in-flight request completion/cancellation, health, and cleanup. | Focused `internal/web/server_test.go`; `httptest` and controlled listeners/fakes. |
| Routes/API | Every route/method/HEAD/Allow case, stable DTO/envelopes, unknown `/api/` JSON, content types, no-store, truncation, health, and absent operation routes. | Focused `internal/web/routes_test.go`; semantic assertions plus narrow normalized envelope goldens. |
| Security | Host/Origin, no-CORS, request/body/query/identifier limits, security headers, safe error projection, request-ID generation, redaction, and no disclosure. | Focused `internal/web/security_test.go`; log/response leak assertions. |
| Artifacts | Opaque-reference containment, traversal/symlink/extension/size rejection, bounded reads, Markdown/JSON escaping, invalid JSON, hostile content, and truncation. | Focused `internal/web/artifacts_test.go` plus app artifact-use-case tests. |
| Templates | All embedded templates parse and representative dashboard/detail/empty/error/artifact states render semantically and safely. | Focused `internal/web/templates_test.go`. |
| Runtime | Startup/listen/request/shutdown diagnostics include safe stable fields and omit raw sensitive values; health is truthful. | Captured structured diagnostic tests and manual `ultraplan serve` local check if review permits. No real runtime/provider is required. |
| Review | Dependency direction, transport-only web ownership, no operation/SSE scope leakage, unchanged CLI/TUI, bounded lifecycle, and single-binary frontend are confirmed. | `system/protocols/architecture-review-protocol.md` and `system/protocols/review-sprint-protocol.md`. |
| Accessibility/Responsive | No-JavaScript navigation, keyboard focus/order, semantic headings/landmarks/statuses, narrow reflow, local overflow, 200-percent zoom, and hostile-content inertness are confirmed. | Template semantic tests plus documented desktop/narrow/keyboard review checks. |
| Documentation | Serve help and three docs accurately describe flags, loopback/read-only limitations, API status, previews, shutdown, troubleshooting, and architecture boundary. | `../ultraplan-go/docs/local-web.md`, `../ultraplan-go/docs/cli-reference.md`, `../ultraplan-go/docs/architecture.md`, and `ultraplan serve --help`. |

## Assumptions And Risks

| Item | Type | Impact | Mitigation / Follow-Up |
| --- | --- | --- | --- |
| Phase 3 is complete or exceptions are recorded before implementation starts. | Assumption | Missing app status/query behavior could expand Sprint 30 scope. | Plan preflight must confirm readiness; record blockers rather than duplicate Phase 3 behavior in web. |
| Existing app/product packages can expose required read projections without workflow-semantic changes. | Assumption | Missing typed queries could tempt CLI parsing or web-to-product imports. | Add typed app use cases only; preserve product ownership and test agreement with existing surfaces. |
| `127.0.0.1:8080` is an acceptable default for local use. | Assumption | Port collision can cause startup failure. | Return actionable startup error and document `--listen`; do not silently select a different port. |
| Loopback reduces but does not eliminate local disclosure. | Risk | Malicious local pages/processes or same-user processes may probe data. | Exact bind/Host/Origin, no CORS, opaque refs, safe errors, no secrets/paths, security headers; document single-user local trust boundary. |
| Strict IP-literal Host/Origin policy may surprise users expecting `localhost`. | Risk | Browser URL aliases can be rejected. | Open/print the canonical bound URL and document accepted forms; revisit only under explicit security reasoning. |
| Live reads amplify work on refresh. | Risk | Large workspaces may render slowly. | Fixed limits, no recursive source scans, cancellation, duration diagnostics, explicit refresh only; measure before cache/fan-out. |
| Opaque references can become stale between detail and preview requests. | Risk | Preview returns not found after workspace change. | Return safe 404 and require refresh; never redirect to a different artifact. |
| Symlink/encoding bugs can bypass lexical containment. | Risk | Arbitrary local file disclosure. | Canonical/symlink-aware app resolution plus encoded traversal and symlink test matrix; web never receives a broad reader. |
| DTO compatibility can drift with bundled JavaScript changes. | Risk | Browser/API mismatch or accidental disclosure. | Explicit transport DTOs, envelope goldens, semantic schema tests, compatibility review for `v1`. |
| Template escaping can be bypassed by trusted types or DOM APIs. | Risk | Workspace scripts/HTML execute. | Prohibit `template.HTML`/`innerHTML` for workspace data; static external scripts only; hostile-content tests and CSP. |
| Composition changes can regress CLI/TUI startup, help, output, or exit behavior. | Risk | Existing automation/local interface breaks. | Lazy construction and command-path non-regression tests; Architecture Review. |
| Shutdown can leak listener/handler/launcher goroutines. | Risk | Hung process/tests and races. | Explicit ownership, context propagation, 10-second shutdown bound, joined goroutines, race/lifecycle tests. |
| Fixed bounds may be too small or large for real workspaces. | Risk | Truncation or unnecessary resource use. | Make truncation visible, document values, instrument duration/bytes, revisit from evidence. |
| Broader PRD/TRD operation routes can leak into Sprint 30. | Risk | Read-only/security scope is violated. | Inject query capability only and test that POST/DELETE/operation/confirmation/SSE routes and UI controls are absent. |

## Implementation Constraints

- `cmd/ultraplan` owns process-level construction and signal context but no HTTP or product behavior.
- `internal/app` owns serve parsing/preflight, shared query composition, typed results, error identity, and CLI exit projection.
- `internal/web` depends only on `internal/app` use-case contracts/plain results plus standard/support packages; no product module, workspace implementation, runtime/process adapter, CLI handler, or arbitrary filesystem dependency.
- The web query facade contains no validate-now, execute, review, smoke, runtime, process, cancellation, persistence, Git, or other mutation capability.
- No package-global mutable runner, registry, service locator, `init` registration, or context-carried dependency.
- Construct/parse web-only resources lazily for `serve`; existing help/version/CLI/TUI paths must not initialize HTTP facilities.
- Resolve workspace/config and reject invalid/non-loopback listen values before runner/listener use; bind only numeric loopback IP literals.
- Use the exact route, envelope, error-code, method, request, collection, preview, timeout, shutdown, and concurrency contracts in Decisions 2-5 and 7.
- Unknown `/api/` paths always return structured JSON and never HTML.
- Every HTTP route is `GET`/`HEAD` inspection only; no POST, PUT, PATCH, DELETE, confirmation, operation, SSE, WebSocket, or runtime/smoke route exists.
- App-owned canonical resolution must precede artifact reads; references are opaque and app-issued; web never opens caller-selected paths.
- Workspace artifacts and product-owned run/flow/review/smoke state remain authoritative; no cache, snapshot, watcher, browser durable state, database, or alternate state machine.
- Propagate root/request contexts to all reads; own and bound every goroutine; return only after shutdown/cleanup completes or is truthfully reported.
- Parse embedded templates once at startup, buffer page rendering, and use contextual `html/template` escaping; never trust workspace markup.
- Embedded CSS and at most one small dependency-free external script only; no Node.js, Vite, frontend framework, hydration, client router/store, third-party asset, or build pipeline.
- Initial HTML must be complete and accessible without JavaScript; enhanced refresh is explicit-user-triggered and retains ordinary reload fallback.
- Keep diagnostics separate from HTML, JSON, and normal CLI output; apply the disclosure/redaction rules in Decision 5.
- Use explicit web DTOs/view models rather than serializing app/product structs or passing broad maps/models to templates.
- Tests default to app/runtime/process fakes, controlled listeners, temporary workspaces, and `httptest`; no real provider, smoke harness, browser process, Node runtime, or network service is required.
- Do not alter workflow, validation, review, smoke, execution, persistence, CLI, or TUI semantics beyond typed read queries and explicit runner composition.
- Do not perform automatic Git add, commit, push, branch, merge, reset, checkout, or any other Git mutation.

## Plan Handoff

`plan.md` must execute these decisions. It must not invent architecture, scope, route semantics, limits, state ownership, frontend technology, or security behavior beyond this document.

The plan must carry forward:

- all eight final decisions and their rejected alternatives
- applicable contracts and `RO-*`/`AC-*`/`NG-*`/`C-*` requirement references
- exact route, API envelope, error, limit, lifecycle, security, rendering, and accessibility constraints
- expected test, runtime, review, build, and documentation evidence
- assumptions, risks, mitigations, debt, and revisit triggers
- Architecture Review and Sprint Review protocols

## Phase Exit Criteria

- [x] Selected context was read and used.
- [x] Area-specific API Design, Architecture, and Frontend reasoning documents were completed and summarized.
- [x] All area-specific conclusions are reflected; broader Phase 4 operations/SSE are explicitly overridden by Sprint 30 scope.
- [x] Contracts and ordered requirement IDs are mapped to decisions and expected evidence.
- [x] Final decisions fix composition, lifecycle, routes, API/error contracts, limits, state ownership, artifact safety, frontend behavior, and verification without core architecture placeholders.
- [x] At least two alternatives are rejected with rationale; each decision records its own rejected alternatives.
- [x] Expected evidence is specific, reviewable, and tied to commands, paths, or protocols.
- [x] Trade-offs, potential debt, future considerations, assumptions, risks, and implementation constraints are recorded.
- [x] `plan.md` can execute these decisions without reopening architecture.
