# Source Analysis: langfuse

## Dimension 24.01: Public API Surface

### Source Info

| Field | Value |
|-------|-------|
| Name | langfuse |
| Path | `studies/agent-harness-study/sources/langfuse` |
| Language / Stack | TypeScript monorepo (Next.js web app + worker; pnpm/turbo; Zod; Prisma/ClickHouse) |
| Analyzed | 2026-08-22 |

## Summary

Langfuse exposes its public API surface as HTTP endpoints implemented in a Next.js Pages Router application, with the wire contract maintained separately as Fern YAML definitions that are exported to OpenAPI specs and used to generate SDKs. The surface is not a single API but a family of distinct surfaces with different auth models:

1. **Project-scoped public REST API** under `/api/public/**` (73 route files), implemented in `web/src/pages/api/public/**` and contracted in `fern/apis/server/definition/**` (98 endpoint methods across the server API definitions). Auth is Basic (public key : secret key) or Bearer with scoped access levels (`packages/shared/src/server/auth/types.ts:58`).
2. **OpenTelemetry ingestion** at `/api/public/otel/v1/traces|metrics` accepting OTLP JSON and protobuf (`web/src/pages/api/public/otel/v1/traces/index.ts:33-114`).
3. **MCP server endpoint** at `/api/public/mcp` implementing Streamable HTTP transport for AI assistants (`web/src/pages/api/public/mcp/index.ts:65-183`).
4. **SCIM provisioning** at `/api/public/scim/**` (`web/src/pages/api/public/scim/Users/index.ts`).
5. **Self-hosted-only organizations/admin API**, contracted as a separate Fern API declaring "only available on self-hosted instances" with Bearer `ADMIN_API_KEY` auth (`fern/apis/organizations/definition/api.yml:3-40`).
6. **Operational endpoints** `/api/public/health` and `/api/public/ready` (`web/src/pages/api/public/health.ts:8-32`).

The distinguishing architectural choice is a contract-first pipeline: Fern YAML is the source of truth, hand-maintained for understandability (`CONTRIBUTING.md:504-513`), exported to OpenAPI specs served from `web/public/generated/{api,api-client,organizations-api}`, and generated SDK outputs are explicitly forbidden from hand edits (`AGENTS.md`, "Generated Files" section). A standardized route harness (`createAuthedProjectAPIRoute` + `withMiddlewares`) centralizes auth, rate limiting, validation, and error mapping, and repo-level agent instructions make Fern updates a hard requirement for any public contract change.

Stability tiers are explicit but unevenly applied: an `unstable/` URL namespace with its own machine-readable error contract and documented evolution caveat, `legacy/` Fern definitions for superseded v1 read endpoints, path-versioned `v2`/`v3` prefixes for some resources, and unprefixed paths implicitly meaning v1.

## Rating

**8 / 10** — Clear model with tests, explicit interfaces, and operational safeguards.

Rationale:
- The stable surface is identifiable without reading internals: every public route lives under one directory (`web/src/pages/api/public/**`), the contract lives under one directory (`fern/apis/server/definition/**`), and the two are tied together by mandated workflow (`AGENTS.md`: "Public API contract changes must update Fern sources in `fern/apis/**`"; `CONTRIBUTING.md:504-513`).
- The unstable tier is genuinely quarantined: URL namespace (`fern/apis/server/definition/unstable/evaluators.yml:9` base-path `/api/public/unstable`), dedicated error codes enum (`fern/apis/server/definition/unstable/errors.yml:2-26`), and a documented "may evolve" disclaimer (`fern/apis/server/definition/unstable/evaluators.yml:40-41`).
- Operational safeguards are wired into the route harness itself: per-resource rate limits (`web/src/features/public-api/server/createAuthedProjectAPIRoute.ts:47,367-378`), ingestion suspension enforcement (`web/src/pages/api/public/otel/v1/traces/index.ts:41-45`), audit logging on destructive operations (`web/src/pages/api/public/traces/index.ts:194-208`), and ClickHouse resource-exhaustion responses that carry migration advice (`web/src/features/public-api/server/withMiddlewares.ts:40-58,111-134`).
- It falls short of 9–10 because roughly 29% of public route files bypass the standard harness (21 of 73, see Evidence), the error body shape differs between stable (`{message, error}`) and unstable (`{message, code, details}`) tiers, URL versioning is partial (implicit v1, selective v2/v3 prefixes), and two Slack OAuth routes sit under `/api/public` while using browser session auth rather than API keys — accidental surface placement.

## Evidence Collected

Every entry includes a file path with line numbers, workspace-relative to `sources/langfuse/`.

| Area | Evidence | File:Line |
|------|----------|-----------|
| Public REST API route tree | 73 route files under `src/pages/api/public/**`; 48 use the standard harness | `web/src/pages/api/public/traces/index.ts:35-220` |
| Contract-first workflow (mandated) | "Public API contract changes must update Fern sources in `fern/apis/**` and regenerated outputs. Never hand-edit `generated/**`." | `AGENTS.md` (root, Generated Files section); `CONTRIBUTING.md:504-513` |
| Fern server API definition set | 98 endpoint methods across 30+ service files; auth declared `basic` | `fern/apis/server/definition/api.yml:1-25` |
| Generated OpenAPI served from web app | Export targets `web/public/generated/api/openapi.yml`, `organizations-api`, `api-client` exist in repo | `CONTRIBUTING.md:511-513`, `web/public/generated/` |
| Route harness: typed route factory | `createAuthedProjectAPIRoute` config requires zod `responseSchema`, optional query/body schemas, rate limit resource, access levels, admin-key flag, events-only rejection | `web/src/features/public-api/server/createAuthedProjectAPIRoute.ts:37-92` |
| Auth model in harness | Basic or Bearer via `ApiAuthService.verifyAuthHeaderAndReturnScope`; access-level gating defaults to `["project"]` | `web/src/features/public-api/server/createAuthedProjectAPIRoute.ts:106-146,299-333` |
| Access levels incl. Bearer-public-key `scores` scope | `ApiAccessLevel = "organization" \| "project" \| "scores"`; scores POST accepts both project and scores levels | `packages/shared/src/server/auth/types.ts:58`; `web/src/pages/api/public/scores/index.ts:16-20` |
| Self-hosted admin-key auth on public routes | `isAdminApiKeyAuthAllowed` requires dual headers + timing-safe compare + cloud-region block | `web/src/features/public-api/server/createAuthedProjectAPIRoute.ts:49-62,164-251` |
| Separate ee Admin API guard | `AdminApiAuthService` blocks on cloud unless allowed; separate from public API-key auth | `web/src/ee/features/admin-api/server/adminApiAuth.ts:11-31` |
| Middleware error mapping (stable contract) | BaseError → `{message, error:name}`; ZodError → 400 `{message, error:issues}`; ClickHouseResourceError → 422 with migration advice | `web/src/features/public-api/server/withMiddlewares.ts:159-196,111-134` |
| Legacy-endpoint migration nudges | Error messages point to `/api/public/v2/observations` and metrics v2 replacements | `web/src/features/public-api/server/withMiddlewares.ts:45-58` |
| Unstable error contract | Machine-readable `code` enum + structured `details`; `unstablePublicEvalsErrorContract = "unstable-public-evals"` | `web/src/features/public-api/server/unstable-public-api-error-contract.ts:24-47,160-262` |
| Unstable error contract in Fern | `PublicApiErrorCode` enum docs: "SDKs, CLIs, and agents should branch on `code` rather than parsing `message`" | `fern/apis/server/definition/unstable/errors.yml:2-26` |
| Unstable route wrapper | `createUnstablePublicEvalsRoute`/`withUnstablePublicEvalsMiddlewares` force the unstable error contract | `web/src/features/public-api/server/unstable-public-evals-route.ts:29-46` |
| Versioning: v2/v3 prefixes | `base-path: /api/public/v2` (prompts, scores v2); `/v3/scores` endpoint under `/api/public` | `fern/apis/server/definition/prompts.yml:8`; `fern/apis/server/definition/scores-v3.yml:7,30`; `web/src/pages/api/public/v3/scores/index.ts:99-138` |
| Versioning: legacy definitions quarantined in Fern | `fern/apis/server/definition/legacy/{score-v1,observations-v1,metrics-v1}.yml` | `fern/apis/server/definition/legacy/score-v1.yml:1-14` |
| Breaking-change deployment gate | `rejectInEventsOnlyMode` returns 404 on legacy read/write routes when `LANGFUSE_MIGRATION_V4_WRITE_MODE=events_only` | `web/src/features/public-api/server/createAuthedProjectAPIRoute.ts:76-82,307-320` |
| v3 rejects unsupported filters with guidance | `userId`/`traceTags` filters rejected: "use v2 or omit this filter" | `web/src/pages/api/public/v3/scores/index.ts:10-23` |
| Cursor pagination introduced in v3 | `EncodedScoresCursorV3` cursor param replaces page-based meta | `web/src/pages/api/public/v3/scores/index.ts:7-8,129-135`; `fern/apis/server/definition/scores-v3.yml:36-39` |
| OTel ingestion endpoint | OTLP protobuf + JSON, gzip support, strict content-type check, 16MB warn threshold | `web/src/pages/api/public/otel/v1/traces/index.ts:64-134` |
| Ingestion version negotiation | `x-langfuse-ingestion-version` header; versions > 4 rejected with explicit message | `web/src/pages/api/public/otel/v1/traces/index.ts:136-160` |
| MCP endpoint (non-harness) | Streamable HTTP, BasicAuth project-scope only, own CORS/security layer, explicit comment why it bypasses `withMiddlewares` | `web/src/pages/api/public/mcp/index.ts:7-24,65-171` |
| SCIM surface | Dedicated RFC-style routes (Users, ResourceTypes, Schemas, ServiceProviderConfig) | `web/src/pages/api/public/scim/Users/index.ts`, `.../ServiceProviderConfig.ts` |
| Organizations API (self-hosted only) | Docs: "This admin API is only available on self-hosted instances… set `ADMIN_API_KEY`" | `fern/apis/organizations/definition/api.yml:3-40` |
| Health/readiness surface | `/api/public/health` with optional strictness query flags, 200/503 semantics, version payload | `web/src/pages/api/public/health.ts:15-31` |
| Harness bypass count | 21 of 73 public route files lack `withMiddlewares`/`createAuthedProjectAPIRoute` usage, e.g. prompts, ingestion, projects/organizations CRUD, SCIM, slack, v2 prompts | `web/src/pages/api/public/prompts.ts:18-77`; search over `web/src/pages/api/public/**` |
| Manual reimplementation example (prompts v1) | Inline `ApiAuthService` call, manual rate limit, manual error handling instead of harness | `web/src/pages/api/public/prompts.ts:31-77` |
| Misplaced session-auth route under /api/public | Slack install/oauth use `getServerAuthSession` cookie auth, not API keys | `web/src/pages/api/public/slack/install/index.ts:16-30` |
| Abstraction boundary discipline (evals) | Guidance: do not leak internal `EvalTemplate`/`JobConfiguration` names into public contract; split evaluators vs evaluation-rules | `web/AGENTS.md` ("Add/Change public API endpoint", eval naming rule) |
| Internal tRPC kept separate | UI backend at `/api/trpc` with own router registry, distinct from public REST | `web/AGENTS.md` High-Signal Entry Points; `web/src/pages/api/trpc/` |
| Rich embedded API docs + examples | Evaluator create endpoint documents naming/versioning behavior, recommended workflow, per-error recovery guidance, runnable request/response examples; 6 of ~30 endpoint-bearing Fern services include example blocks | `fern/apis/server/definition/unstable/evaluators.yml:12-80`; `grep -l "examples:" fern/apis/server/definition/**` |
| Shared pagination/error types | Copy-paste pagination template and common error set declared once | `fern/apis/server/definition/utils/pagination.yml:1-17`; `fern/apis/server/definition/commons.yml` errors list in `api.yml:15-21` |
| Test coverage of public surface | 170 server test files; per-family API tests incl. v1/v2/v3 score splits, otel, scim, mcp, unstable evals | `web/src/__tests__/server/traces-api.servertest.ts`, `scores-api-v1.servertest.ts`, `scores-api-v2.servertest.ts`, `scores-api-v3.servertest.ts`, `otel-api.servertest.ts`, `scim-api.servertest.ts`, `mcp-public-api-tools.servertest.ts`, `unstable-evals-api.servertest.ts` |
| Harness unit tests | Auth-error and error-contract behavior tested directly | `web/src/__tests__/server/unit/create-authed-project-api-route-auth-errors.servertest.ts`, `unit/unstable-public-api-error-contract.servertest.ts` |
| Dev-mode response validation safeguard | Handler responses validated against `responseSchema` when `NODE_ENV=development` | `web/src/features/public-api/server/createAuthedProjectAPIRoute.ts:448-454` |
| Rate-limit upgrade paths surfaced to clients | `rateLimitUpgradePath` attached to legacy trace listing; upgrade message injected into 429 bodies | `web/src/pages/api/public/traces/index.ts:74-77`; `web/src/features/public-api/server/rateLimitUpgradePaths.ts` |
| Audit trail on destructive ops | Trace deletion writes audit log entries per trace ID before processing | `web/src/pages/api/public/traces/index.ts:189-213` |

## Answers to Dimension Questions

**1. What is the intended public API surface?**
Six surfaces, each with a declared contract: (a) the project-scoped REST API at `/api/public/**` contracted by `fern/apis/server/definition/**` with Basic/Bearer API-key auth (`fern/apis/server/definition/api.yml:6-25`); (b) OTLP ingestion at `/api/public/otel/v1/*` (`fern/apis/server/definition/opentelemetry.yml:4`); (c) the MCP endpoint (`web/src/pages/api/public/mcp/index.ts:1-25`); (d) SCIM under `/api/public/scim` (`fern/apis/server/definition/scim.yml:7`); (e) the self-hosted-only organizations admin API as a separate Fern API with Bearer admin key (`fern/apis/organizations/definition/api.yml:3-40`) plus the ee-guarded `/api/admin/**` implementation (`web/src/ee/features/admin-api/server/adminApiAuth.ts:11-31`); and (f) health/ready probes (`web/src/pages/api/public/health.ts:8-32`). The UI's tRPC backend (`web/src/pages/api/trpc/`) is deliberately outside this surface.

**2. Is the stable API easy to distinguish from internal implementation details?**
Mostly yes, with caveats. Directory boundaries are clean: public routes only in `web/src/pages/api/public/**`, contracts only in `fern/apis/**`, and repo policy forbids hand-editing generated output (`AGENTS.md`). The unstable tier is fenced by URL prefix, dedicated Fern files, and a distinct error contract (`web/src/features/public-api/server/unstable-public-evals-route.ts:29-46`). Naming discipline is enforced socially and structurally — the public evals contract hides internal `EvalTemplate`/`JobConfiguration` names behind evaluator/rule concepts (`web/AGENTS.md`). Caveats: the Slack OAuth routes under `/api/public` are session-cookie-authenticated browser flows (`web/src/pages/api/public/slack/install/index.ts:22-30`), so directory membership alone does not prove API-key surface; and ~21 route files bypass the standard harness, so a reader cannot assume uniform request handling from the route tree alone.

**3. Does the API expose the right level of abstraction for agent harness users?**
Yes for the primary consumers. Agent integrations get three fit-for-purpose channels: managed SDKs generated from Fern (Python/TS configured in `fern/apis/server/generators.yml:1-32`), raw OTLP for instrumentations (`web/src/pages/api/public/otel/v1/traces/index.ts:90-114` supporting protobuf and JSON), and MCP for assistant-driven prompt management (`web/src/pages/api/public/mcp/index.ts:65-145`). Header-based negotiation (`x-langfuse-sdk-name/-version/-ingestion-version`, accepted up to version 4) decouples client release cycles from server capability (`web/src/pages/api/public/otel/v1/traces/index.ts:136-160`). The unstable evals API documents agent-oriented recovery guidance keyed to machine-readable error codes rather than prose parsing (`fern/apis/server/definition/unstable/evaluators.yml:34-41`; `fern/apis/server/definition/unstable/errors.yml:5-8`). One leak: the stable error contract still mixes `{message, error}` shapes where `error` is sometimes the exception class name and sometimes an issues array (`web/src/features/public-api/server/withMiddlewares.ts:159-196`), which is thinner than the unstable tier's typed envelope.

**4. Are examples sufficient to use the API correctly without reading internals?**
Within the Fern definitions, largely yes: endpoints carry docs blocks with behavioral contracts (e.g., evaluator version auto-increment semantics and preflight-failure recovery, `fern/apis/server/definition/unstable/evaluators.yml:12-38`) and named request/response examples (`:55-80`). Pagination and filter constraints are documented at parameter level with explicit 400 conditions (`fern/apis/server/definition/scores-v3.yml:36-94`). However: only six of the ~30 endpoint-bearing Fern service definitions carry `examples:` blocks (`health.yml`, `ingestion.yml`, `legacy/score-v1.yml`, `opentelemetry.yml`, `unstable/evaluation-rules.yml`, `unstable/evaluators.yml`) against 98 endpoint methods, so most stable endpoints lack in-spec examples; the human-facing docs, quickstarts, and SDK repositories live outside this monorepo (README only links out, `README.md:139`); and no runnable in-repo example scripts target the REST API — correctness verification rests on the server test suite (`web/src/__tests__/server/*.servertest.ts`) rather than consumer-facing samples.

## Architectural Decisions

1. **Contract-first with generated clients, enforced by process.** Fern YAML is the single source of truth; exports produce the OpenAPI specs powering the hosted reference (`CONTRIBUTING.md:504-513`), and repo policy makes updating `fern/apis/**` a blocking step for contract changes (`AGENTS.md`). This trades a second artifact to maintain for guaranteed doc/spec parity.
2. **One route harness as the stability chokepoint.** `createAuthedProjectAPIRoute` funnels authentication variants (Basic, Bearer with access levels, self-hosted admin key), rate limiting, zod validation, dev-time response checking, and migration-mode gating through one code path (`web/src/features/public-api/server/createAuthedProjectAPIRoute.ts:37-92,299-474`). New cross-cutting behavior (e.g., events-only-mode 404s) ships once and applies to all conforming routes (`:307-320`).
3. **Explicit instability instead of hidden instability.** Newer eval functionality is published under `/api/public/unstable` with a weaker guarantee stated in the spec itself and a richer, machine-readable error envelope (`fern/apis/server/definition/unstable/errors.yml:2-26`), letting the team iterate while keeping the stable tier conservative.
4. **Versioning by path prefix plus deployment-mode gates.** v2/v3 prefixes coexist with unprefixed legacy routes quarantined in a `legacy/` folder of the spec (`fern/apis/server/definition/legacy/`), and the v4 storage migration uses a runtime flag to 404 incompatible legacy endpoints rather than a new URL space (`web/src/features/public-api/server/createAuthedProjectAPIRoute.ts:76-82`).
5. **Separate trust domains for separate surfaces.** Project API keys, org-scoped keys, Bearer public-key `scores`-only keys (`packages/shared/src/server/auth/types.ts:58`), in-app agent keys (`allowInAppAgentKey`, `createAuthedProjectAPIRoute.ts:71-74`), and self-hosted admin keys all resolve through the same verify function into distinct scopes, with cloud/self-hosted gating for elevated auth (`:178-183`).

## Notable Patterns

- **Schema-in/schema-out handlers.** Every harness route declares zod `querySchema`/`bodySchema`/`responseSchema` (`createAuthedProjectAPIRoute.ts:43-45`); handler bodies receive parsed, typed values, and dev mode validates outgoing payloads against the response schema (`:448-454`).
- **Error taxonomy centralized in shared package.** `BaseError` subclasses with `httpCode` map uniformly to HTTP responses in both middleware stacks (`withMiddlewares.ts:159-168`), and the unstable stack additionally translates them into coded envelopes (`unstable-public-api-error-contract.ts:160-262`).
- **Migration pressure embedded in responses.** Legacy endpoints under load return 422 bodies that name their v2 replacements (`withMiddlewares.ts:45-58`), and rate-limited legacy listings advertise upgrade paths (`web/src/pages/api/public/traces/index.ts:74-77`) — deprecation communicated at runtime, not just in docs.
- **Per-resource rate limiting as API metadata.** Routes declare semantic resources (`legacy-ingestion`, `trace-delete`, `media-upload`, `public-api-legacy`) so limits can be tuned per operation family rather than per IP (`traces/index.ts:41,74,193`; `media/[mediaId].ts:38`).
- **Spec-as-documentation.** Fern `docs:` blocks carry behavioral contracts (naming/versioning rules, recovery guidance, recommended workflows) close to the schema they describe (`fern/apis/server/definition/unstable/evaluators.yml:12-41`).

## Tradeoffs

- **Dual maintenance burden:** every behavioral change touches route + zod types + Fern YAML + tests; the repo acknowledges this by mandating all four in its playbooks (`web/AGENTS.md` "Add/Change public API endpoint"). The payoff is reference/spec fidelity; the cost is higher change latency.
- **Harness adoption is opt-in:** consistency depends on convention. 52 of 73 route files comply; the rest (ingestion, prompts v1/v2, projects/organizations CRUD, SCIM, slack, health) hand-roll auth and errors. Some bypasses are justified (streaming MCP transport needs direct response control, `mcp/index.ts:18-20`; health checks need no auth), but prompts v1 duplicates the harness's auth/ratelimit/error logic line-for-line (`prompts.ts:31-77`), which will drift.
- **Two error dialects:** clients written against the stable tier parse `{message, error}` while unstable-tier clients parse `{message, code, details}` (`withMiddlewares.ts:159-196` vs `unstable-public-api-error-contract.ts:27-31`). The better contract exists but only in the quarantine zone, presumably pending a stable-tier migration.
- **Path versioning covers a subset:** observations/prompts/scores/datasets/metrics gained v2 (and scores a v3) while traces, models, sessions, media, etc. remain unprefixed; consumers cannot infer currency of a resource from its URL shape alone.
- **Examples investment is uneven:** exhaustive narrative docs in newer endpoints (`unstable/evaluators.yml`) versus bare or absent examples across most of the older surface.

## Failure Modes / Edge Cases

- **Oversized payloads:** JSON stringify overflow on huge responses is caught and converted to `PayloadTooLargeError` (`createAuthedProjectAPIRoute.ts:32-35,463-471`); request bodies are capped by Next config (4.5MB on ingestion/MCP, `ingestion.ts:29-35`, `mcp/index.ts:177-183`); OTel requests >16MB are logged with span counts for capacity monitoring but still processed (`otel/v1/traces/index.ts:120-134`).
- **Storage-migration breakage:** in `events_only` mode, legacy trace/observation/score read-write routes return 404 with a docs link instead of silently returning stale data (`createAuthedProjectAPIRoute.ts:76-82,307-320`; `traces/index.ts:44,78`; `scores/index.ts:59`).
- **Query-cost explosions on legacy reads:** ClickHouseResourceError is mapped to 422 with actionable advice and, on cloud, a pointer to the v2 replacement APIs (`withMiddlewares.ts:111-134`); deployments can also enforce `fromTimestamp` requirements or default date ranges/field projections via env flags (`traces/index.ts:80-115`).
- **Abuse containment:** ingestion suspension (quota exceeded) blocks writes at the harness boundary and inside OTel/MCP handlers (`otel/v1/traces/index.ts:41-45`; `mcp/index.ts:102-106`); unsupported future ingestion versions fail closed with the maximum supported value named (`otel/v1/traces/index.ts:147-160`).
- **Destructive operations:** multi-trace delete writes audit records per ID before invoking the deletion processor, preserving accountability if processing partially fails (`traces/index.ts:189-213`).
- **Known sharp edge:** the stable error contract leaks exception internals — `error` field carries the error class name for `BaseError`s (`withMiddlewares.ts:164-167`) but raw Zod issues arrays for validation failures (`:182-185`), and the harness's Zod instance differs from consumers', acknowledged by a duck-typed instanceof check in-repo (`:179-180`, `withMiddlewares.ts:202-204`).

## Future Considerations

- Migrate the remaining hand-rolled routes (especially `prompts.ts` and the projects/organizations CRUD group) onto `createAuthedProjectAPIRoute` so harness invariants (rate limiting, audit hooks, migration gates) apply uniformly; the v2 prompts delegation pattern (`web/src/pages/api/public/v2/prompts/index.ts:1` exporting a feature-local handler) shows the target shape.
- Promote the unstable `{code, details}` error envelope to the stable tier once breaking-change tolerance allows, unifying client error handling.
- Complete path versioning or adopt an explicit deprecation header scheme so resource currency is discoverable from the wire, not from Fern file layout.
- Raise example coverage in Fern definitions toward parity with the evaluators endpoint's standard (named examples + recovery guidance per error code).
- Relocate Slack OAuth routes out of `/api/public` (or rename the mount) so the directory boundary reliably means "API-key-authenticated public surface".

## Questions / Gaps

- **SDK source code is not in this repository.** The Fern generators write to external repos (`../../../generated/python`, `../../../langfuse-java`, `fern/apis/server/generators.yml:9-32`); the actual developer-facing client libraries could not be inspected here, so claims about SDK ergonomics are out of scope for this study.
- **Deprecation timeline mechanics are undocumented in-repo.** No evidence found of sunset headers, removal schedules, or per-version support windows; searched `fern/apis/**`, `CONTRIBUTING.md`, and route comments. Only runtime nudges (`withMiddlewares.ts:45-58`) signal deprecation.
- **Consumer-side contract testing** (e.g., generated-client round-trip tests against the running server) was not found in this repo; the 170 server test files exercise the server side (`web/src/__tests__/server/`), and generated outputs are excluded from hand edits, but nothing here verifies the shipped SDKs against live behavior.
- **`fern/apis/client` surface purpose is thin in-repo:** a single bearer-auth endpoint method (`grep` count: 1 method across `fern/apis/client/definition/`) whose role relative to the server API is not explained in-repo beyond the export command (`CONTRIBUTING.md:512`).

---

Generated by `dimension 24.01: Public API Surface` against `langfuse`.
