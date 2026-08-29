# Source Analysis: langfuse

## 21.01 Plugin and Extension Points

### Source Info

| Field | Value |
|-------|-------|
| Name | langfuse |
| Path | `studies/agent-harness-study/sources/langfuse` |
| Language / Stack | TypeScript / Next.js (Pages Router) + pnpm/turbo monorepo, Prisma (Postgres), ClickHouse, Redis/BullMQ, Node 24 |
| Analyzed | 2026-08-28 |

## Summary

Langfuse is an LLM observability/evaluation platform, not an agent harness. It has **no third-party plugin system**. The closest analogue is an internal, compile-time registry for MCP (Model Context Protocol) tools (`web/src/features/mcp/server/registry.ts:75`) that follows a `McpFeatureModule` interface (`web/src/features/mcp/server/registry.ts:52`). Eleven feature modules (prompts, datasets, scores, observations, metrics, etc.) are statically imported and registered in `web/src/features/mcp/server/bootstrap.ts:34-50` — adding a new tool type or feature requires modifying `bootstrap.ts` and rebuilding. There is no dynamic discovery, no filesystem/marketplace loader, no `npm`/`pip` plugin resolution, and no isolation/sandboxing between features. Related extension surfaces (tRPC router registry `web/src/server/api/root.ts:68`, public REST routes `web/src/pages/api/public/`, automations/webhooks `packages/shared/src/domain/automations.ts:43`, and integrations for Slack/PostHog/Mixpanel/Blob Storage) are all closed, hard-coded discriminated unions. External extensibility is instead via API/SDK contracts (REST + OTEL ingestion + MCP over Streamable HTTP) and outbound webhooks, not in-process plugins. Overall maturity for this dimension is absent-by-design.

## Rating

**2/10 — Absent / ad-hoc (by design)**

Rationale: No mechanism exists to add a new tool type, provider, memory, eval, prompt, policy, or UI extension without modifying core code. The internal MCP registry is well-structured and tested with explicit interfaces, but it is not a public plugin API: no dynamic loading, no versioning/stability guarantee, no isolation, and no documentation for third-party authors. Outbound integrations (webhooks, Slack, GitHub Dispatch) are closed enums, not pluggable. This matches the 1-3 band ("Absent, implicit, ad-hoc, or unsafe"); scored 2 because the internal patterns are clean and could be promoted to a plugin API but currently are not.

## Evidence Collected

Every entry MUST include a file path with line numbers. Format: `path/to/file.ts:NN`.

| Area | Evidence | File:Line |
|------|----------|-----------|
| Extension interface — MCP feature module | `McpFeatureModule` interface with `name`, `description`, `tools: RegisteredTool[]`, optional `isEnabled(context)` | `web/src/features/mcp/server/registry.ts:52-67` |
| Extension interface — registered tool | `RegisteredTool` combines `ToolDefinition` + `ToolHandler<any>` + optional `allowInAppAgentKey` | `web/src/features/mcp/server/registry.ts:22-32` |
| Extension interface — tool definition helper | `DefineToolOptions<TInput>` and `ToolDefinition` with Zod→JSON Schema conversion, `readOnlyHint`/`destructiveHint` annotations, `defineTool()` factory | `web/src/features/mcp/core/define-tool.ts:23-61,112-179` |
| Extension interface — server context | `ServerContext` (projectId, orgId, apiKeyId, accessLevel, publicKey) captured in closures, stateless per-request | `web/src/features/mcp/types.ts:27-54` |
| Plugin loader — static bootstrap | Eleven features statically imported and `toolRegistry.register()` called; adding a feature requires editing this file and rebuilding | `web/src/features/mcp/server/bootstrap.ts:13-46` |
| Plugin loader — auto-bootstrap on import | `bootstrapMcpFeatures()` invoked at module load via side-effect import `import "@/src/features/mcp/server/bootstrap"` in the API route | `web/src/features/mcp/server/bootstrap.ts:58`, `web/src/pages/api/public/mcp/index.ts:45` |
| Plugin loader — absence of dynamic loading | No `fs` scan, no `require`/`import()` of external paths, no env-driven plugin list, no marketplace/manifest loader found; search for `plugin`/`extension` yields only file-extension and ESLint plugin hits | `web/src/features/mcp/server/registry.ts:75-112` (singleton `Map` only) |
| Plugin lifecycle — registration validation | `ToolRegistry.register()` throws on duplicate feature name or duplicate tool name, logs via `logger.info`; no unregister/disable/reload, no versioning | `web/src/features/mcp/server/registry.ts:85-112` |
| Plugin lifecycle — enablement gate | `isEnabled?: (context) => boolean‖Promise<boolean>` checked in `getToolDefinitions()` and `getEnabledTool()`; used by `observations`/`metrics` gated on `LANGFUSE_MIGRATION_V4_ALLOW_PREVIEW_OPT_IN` | `web/src/features/mcp/server/registry.ts:62-67,120-175`, `web/src/features/mcp/features/observations/index.ts:55-57`, `web/src/features/mcp/features/metrics/index.ts:25-26` |
| Plugin lifecycle — per-request server | Fresh `Server` instance per request, `toolRegistry.getToolDefinitions(context)` for `ListTools`, `getEnabledTool(name, context)` for `CallTool`, then `server.close()` + `transport.close()`; no persistent plugin state | `web/src/features/mcp/server/mcpServer.ts:50-108`, `web/src/features/mcp/server/transport.ts:66-94` |
| Plugin isolation — none between features | Single `ToolRegistry` singleton with flat `Map<string, RegisteredTool>`; tools share `ServerContext` and `prisma`/`redis`; no sandbox, no process/VM boundary, no dependency isolation | `web/src/features/mcp/server/registry.ts:75-77,235` |
| Plugin isolation — auth scoping (partial) | `allowInAppAgentKey` flag + `canUseTool()` + `accessLevel === "project"` check; project-scoped API key enforced in MCP route | `web/src/features/mcp/server/registry.ts:177-183`, `web/src/pages/api/public/mcp/index.ts:92-99` |
| Extension point — tRPC router registry | All ~50 tRPC routers statically imported and registered in `appRouter`; adding a capability requires editing `root.ts` | `web/src/server/api/root.ts:1-129` |
| Extension point — public REST API | ~30 route directories under `web/src/pages/api/public/` + `withMiddlewares` pattern; no plugin hook, each route is a file | `web/src/pages/api/public/:` (directory), `web/src/features/public-api/server/withMiddlewares.ts:82-199` |
| Extension point — automations/actions (webhooks) | `ActionTypeSchema = enum(["WEBHOOK","SLACK","GITHUB_DISPATCH"])` — closed discriminated union `ActionConfigSchema`; new action type requires schema + router + worker changes | `packages/shared/src/domain/automations.ts:43,125-129`, `web/src/features/automations/server/router.ts:253-497` |
| Extension point — outbound webhook queue | `WebhookQueue` / `EntityChangeQueue` / `EventPropagationQueue` in shared queue contracts; worker delivers webhooks with HMAC + retry, supports `prompt-version` and `monitor-alert` envelopes | `packages/shared/src/server/queues.ts:247-264,375-376`, `packages/shared/src/domain/webhooks.ts:10-42` |
| Extension point — integrations | PostHog, Mixpanel, BlobStorage (S3), Slack, OTEL — each a dedicated Prisma model + tRPC router + worker queue (`posthogIntegrationRouter`, `blobStorageIntegrationRouter`, `QueueName.PostHogIntegrationQueue`, etc.) — not pluggable | `packages/shared/prisma/schema.prisma:1188,1201,1245,1784`, `web/src/server/api/root.ts:99-101`, `packages/shared/src/server/queues.ts:360-365` |
| Extension point — queue contracts | `QueueName` / `QueueJobs` enums + Zod schemas in `queues.ts`; adding a queue requires editing shared package + producer + consumer | `packages/shared/src/server/queues.ts:342-604` |
| Extension documentation | MCP `README.md` documents client consumption (Claude Code/Cursor) and stateless architecture, but does not document how to author/register a new `McpFeatureModule` as a third party; inline docstring in `bootstrap.ts` gives 3-step internal instructions only | `web/src/features/mcp/README.md:1-199`, `web/src/features/mcp/server/bootstrap.ts:9-11`, `web/src/features/mcp/server/registry.ts:8-11` |
| Tests for extensibility | Worker webhook tests (`webhooks.test.ts`, `webhook-redirect.test.ts`, `promptVersionProcessor.test.ts`), automation tRPC tests (`automations-trpc.servertest.ts`), `automations.test.ts` for webhook schema; no tests for dynamic plugin loading/isolation (none exists) | `worker/src/__tests__/webhooks.test.ts:35`, `worker/src/__tests__/webhook-redirect.test.ts:23`, `packages/shared/src/domain/automations.test.ts:3-30`, `web/src/__tests__/server/automations-trpc.servertest.ts:2328` |

## Answers to Dimension Questions

**1. What can be extended via plugins?**

Nothing via a public plugin API. Internally, four surfaces are extensible only by modifying core code:
- **MCP tools/features:** Add a module implementing `McpFeatureModule` (`web/src/features/mcp/server/registry.ts:52`) and register it in `web/src/features/mcp/server/bootstrap.ts:34-50`. Example feature: `prompts` (`web/src/features/mcp/features/prompts/index.ts:42`) with 6 tools. The `defineTool()` helper (`web/src/features/mcp/core/define-tool.ts:112`) standardizes Zod→JSON Schema and validation. Currently 11 features registered.
- **tRPC procedures:** Add a router file under `web/src/features/*/server/` or `web/src/server/api/routers/` and register in `web/src/server/api/root.ts:68`.
- **Public REST endpoints:** Add files under `web/src/pages/api/public/*` using `withMiddlewares` (`web/src/features/public-api/server/withMiddlewares.ts:82`) and Fern API definitions (`fern/apis/`).
- **Automation actions/integrations:** Extend the closed union `ActionTypeSchema = ["WEBHOOK","SLACK","GITHUB_DISPATCH"]` (`packages/shared/src/domain/automations.ts:43`) — requires shared schema + web router + worker delivery changes. Similarly, adding a new integration (e.g., a new blob storage provider) requires Prisma model + router + queue + worker job.

External extensibility is via **API/SDK** (ingestion `web/src/pages/api/public/ingestion.ts:115`, OTEL `web/src/pages/api/public/otel/`, public CRUD `web/src/pages/api/public/*`), **MCP as a consumer-facing extension point** (`web/src/pages/api/public/mcp/index.ts:65` — AI assistants call tools), and **outbound webhooks** (`packages/shared/src/domain/webhooks.ts:10`, `packages/shared/src/server/queues.ts:247-271`) — not in-process plugins.

**2. Can plugins be loaded at runtime?**

No. Evidence:
- `ToolRegistry` is a singleton `Map` populated only via explicit `register()` calls in `bootstrap.ts:34-46` at import time (`web/src/features/mcp/server/bootstrap.ts:58`). The MCP route side-effects this via `import "@/src/features/mcp/server/bootstrap"` (`web/src/pages/api/public/mcp/index.ts:45`).
- No filesystem scan, no `import()` of user-supplied paths, no config/env-driven plugin list, no manifest/marketplace, no hot-reload/unregister. Search for `plugin` across the repo returns only ESLint/file-extension hits.
- Adding a tool type without core code modification is impossible — the registry rejects duplicate names but has no discovery of external modules. The only runtime gating is `isEnabled(context)` (`web/src/features/mcp/server/registry.ts:62`) for feature flags, not for loading new code.

**3. Are plugins isolated from each other?**

No isolation. All MCP features run in-process in the `web` Next.js server:
- Single `ToolRegistry` instance with flat `tools: Map<string, RegisteredTool>` (`web/src/features/mcp/server/registry.ts:75-77`); handlers share the same `ServerContext`, `prisma`, `redis`, and logger, with no sandbox/VM/worker boundary.
- Per-request `Server` is fresh (`web/src/features/mcp/server/mcpServer.ts:50`) but tools within that request are not isolated — a throwing or slow handler affects the request; no per-tool timeout/circuit-breaker shown.
- Security isolation is limited to auth scoping: `allowInAppAgentKey` filtering (`web/src/features/mcp/server/registry.ts:177-183`) and project-level API key enforcement (`web/src/pages/api/public/mcp/index.ts:92-98`), which scopes *access* but not *execution*.
- Webhook delivery does have retry and failure tracking (`web/src/features/automations/server/router.ts:48-69`, `worker/src/__tests__/webhooks.test.ts:35`) but that is for outbound HTTP isolation, not in-process plugin isolation.

**4. Are extension points documented and stable?**

Partially documented internally, not stable as a public extension API.
- **MCP consumer docs:** `web/src/features/mcp/README.md:1-199` documents how *clients* (Claude Code, Cursor) consume tools, plus stateless architecture and `ServerContext` shape, but not how to author a new feature as a third party. It explicitly warns "Tool availability and schemas may evolve … Clients are expected to tolerate schema changes" (`web/src/features/mcp/server/mcpServer.ts:27-28`), signaling instability.
- **Internal authoring docs:** Only docstrings: `bootstrap.ts:9-11` lists 3 steps ("Create feature module … Import … Call register"), and `registry.ts:8-11,40-51` shows an example `McpFeatureModule`. No versioning, deprecation policy, or compatibility guarantee.
- **API stability signals:** `MCP_SERVER_VERSION = "0.3.0-unstable"` (`web/src/features/mcp/server/mcpServer.ts:29`), Fern API definitions (`fern/apis/`) are versioned but for REST, not for MCP tools. `ActionConfigSchema` discriminated union is versioned via `AvailableWebhookApiSchema` (`packages/shared/src/domain/automations.ts:45-48`) but additions are breaking.
- **No public extension guide, no CHANGELOG for plugin API, no tests asserting stability of the registry contract** beyond duplicate-name checks (`web/src/features/mcp/server/registry.ts:85-99`).

> Can a third party add a new tool type without modifying core code? **No.** All tool/feature registration requires editing `web/src/features/mcp/server/bootstrap.ts` (and creating files under `web/src/features/mcp/features/`) and redeploying. There is no dynamic plugin loader, no external package hook, and no configuration-driven extension.

## Architectural Decisions

- **Stateless per-request MCP server** (`web/src/features/mcp/server/mcpServer.ts:50`, `web/src/features/mcp/server/transport.ts:66`): Fresh `Server` + `StreamableHTTPServerTransport` (stateless, `sessionIdGenerator: undefined`) per request, context captured in closures. Decision trades session state for simplicity and prevents cross-project leaks; eliminates need for plugin lifecycle beyond `register`.
- **Singleton flat registry** (`web/src/features/mcp/server/registry.ts:75-235`): Single `ToolRegistry` with `Map<name, feature>` + `Map<toolName, tool>`. Decision centralizes tool discovery and conflict detection (duplicate tool names throw) but precludes isolation and dynamic loading; also makes testing registry state order-dependent (no `clear()` method).
- **Closed discriminated unions for actions/integrations** (`packages/shared/src/domain/automations.ts:43,125-129`): `WEBHOOK|SLACK|GITHUB_DISPATCH` and `QueueName` enum (`packages/shared/src/server/queues.ts:342-380`) are exhaustive enums. Decision favors type safety and queue contract clarity over extensibility; every new outbound channel is a core change.
- **API-first extensibility over in-process plugins** (`web/src/pages/api/public/mcp/index.ts:65`, `web/src/features/public-api/server/withMiddlewares.ts:82`): Langfuse exposes REST/OTEL ingestion and MCP as the extension surface, delegating host-plugin needs to SDKs and webhooks. Decision avoids the security/operational cost of running third-party code in `web`/`worker`.
- **`defineTool()` with split Zod schemas** (`web/src/features/mcp/core/define-tool.ts:23-35,112-179`): `baseSchema` for JSON Schema generation + `inputSchema` for runtime validation, auto-wrapped with `wrapErrorHandling`. Decision standardizes tool authoring and error formatting (`web/src/features/mcp/core/error-formatting.ts:`, `web/src/features/mcp/core/errors.ts:19-49`) but adds authoring constraints (no union/intersection schemas allowed, `hasJsonSchemaUnion` check at `define-tool.ts:138`).

## Notable Patterns

- **Feature-module plugin pattern (internal):** Each domain exports `McpFeatureModule { name, description, tools, isEnabled? }` (e.g., `web/src/features/mcp/features/prompts/index.ts:42`, `web/src/features/mcp/features/datasets/index.ts:39`, `web/src/features/mcp/features/evals/index.ts:36`, `web/src/features/mcp/features/observations/index.ts:24`). Consistent folder layout: `features/<name>/index.ts` + `tools/*.ts` + `schema.ts`. Pattern is clean and could be promoted to a public plugin API with minimal changes (add dynamic import + manifest).
- **Registry-mediated `ListTools`/`CallTool`:** `mcpServer.ts:64-95` delegates `ListToolsRequestSchema` → `toolRegistry.getToolDefinitions(context)` (filtering by `isEnabled` + `allowInAppAgentKey`) and `CallToolRequestSchema` → `toolRegistry.getEnabledTool(name, context)` → `handler(args, context)`. Ensures discovery and invocation respect the same gates.
- **Outbound webhook as integration extension:** `WebhookOutboundEnvelopeSchema = discriminatedUnion("type", [promptVersionWebhookEnvelopeSchema, MonitorWebhookQueueEventSchema])` (`packages/shared/src/server/queues.ts:261-264`) + `WebhookInputSchema` (`queues.ts:266-271`) enqueued to `QueueName.WebhookQueue` and delivered by worker with HMAC signatures (`packages/shared/src/encryption/signature.ts:24`) and redirect handling (`worker/src/__tests__/webhook-redirect.test.ts:23`).
- **Entitlement/feature-flag gating:** `isEnabled` used for preview-gated features (`LANGFUSE_MIGRATION_V4_ALLOW_PREVIEW_OPT_IN` in `observations/index.ts:55` and `metrics/index.ts:25`), mirroring how `ee/` package could gate enterprise plugins, but currently only for internal flags.

## Tradeoffs

- **Type safety vs. extensibility:** Closed unions (`ActionTypeSchema`, `QueueName`, `ToolDefinition` name uniqueness via throw) give exhaustive switch coverage and queue contract safety, but block third-party additions without core PRs. An open registry (string-keyed with runtime validation) would enable plugins at the cost of weaker compile-time guarantees.
- **Simplicity vs. lifecycle richness:** Singleton registry with no `unregister`/`reload`/`version` keeps the MCP server stateless and simple, but means no hot-reload, no canary/rollout, no per-plugin health checks. A service like Sentry's MCP (referenced in `web/src/features/mcp/types.ts:5` "Following stateless design pattern from Sentry MCP server") makes the same tradeoff.
- **Security (API keys) vs. plugin sandboxing:** By forcing all external code to run *outside* Langfuse (via REST/MCP/webhooks), the system avoids arbitrary-code-execution risks entirely. The cost is no in-process customization (e.g., custom eval operators, prompt transforms) — users must run those externally and call back via API.
- **Zod→JSON Schema strictness:** `defineTool` rejects unions/intersections (`web/src/features/mcp/core/define-tool.ts:138-142`) to stay compatible with MCP's JSON Schema draft-7, simplifying client interop but limiting schema expressiveness for tool authors.
- **Monorepo coupling:** `web → @langfuse/shared → no web/worker/ee imports` (`AGENTS.md` dependency direction) keeps shared contracts stable, but means adding a plugin still touches multiple packages (shared schema + web registration + possibly worker consumer + Fern spec + generated clients).

## Failure Modes / Edge Cases

- **Duplicate registration crashes startup:** `register()` throws `Feature 'x' is already registered` or `Tool 'y' conflicts` (`web/src/features/mcp/server/registry.ts:86-98`). No recovery — the `web` process fails to start; no isolation to keep other features alive. No test for duplicate registration handling beyond the throw.
- **Missing bootstrap import = empty registry:** If `import "@/src/features/mcp/server/bootstrap"` is removed from `web/src/pages/api/public/mcp/index.ts:45`, `ListTools` returns `[]` and all `CallTool` throw `Unknown tool` (`web/src/features/mcp/server/mcpServer.ts:89-91`). Silent failure, no warning log.
- **Registry state leakage across tests:** Singleton `toolRegistry` has no `clear()`/`reset()`; test ordering can cause duplicate-registration throws if features are re-registered. No dedicated registry unit tests found.
- **Feature flag drift:** `observations`/`metrics` gated on `LANGFUSE_MIGRATION_V4_ALLOW_PREVIEW_OPT_IN` (`web/src/features/mcp/features/observations/index.ts:55`). If the env var is mis-set, tools disappear without an explicit error to the client (just absent from `ListTools` and `Unknown tool` on direct call).
- **No per-tool timeout/isolation:** `registeredTool.handler(args, context)` (`web/src/features/mcp/server/mcpServer.ts:95`) runs in the request's event loop with no timeout or sandbox; a slow DB query or infinite loop blocks the MCP response. Transport cleanup in `finally { server.close() }` (`web/src/features/mcp/server/transport.ts:88-94`) mitigates leaks but not handler hangs.
- **Outbound webhook failure modes:** Worker webhook delivery retries with `lastFailingExecutionId` tracking (`packages/shared/src/domain/automations.ts:64,99`) and auto-disable after consecutive failures (`web/src/features/automations/server/router.ts:48-69, `getConsecutiveAutomationFailures`), but no dead-letter queue for webhooks shown; `Webhook does not return 2xx` errors surface as 500s in tests (`worker/src/__tests__/webhooks.test.ts:528,620`).
- **Schema evolution breakage:** `MCP_SERVER_VERSION = "0.3.0-unstable"` and self-describing warning (`web/src/features/mcp/server/mcpServer.ts:27-29`) acknowledge instability; clients that cache tool schemas will break when `defineTool` schemas change, as there is no version negotiation beyond the unstable string.

## Future Considerations

- **Promote internal registry to public plugin API:** Add dynamic discovery (e.g., scan `web/src/features/mcp/features/**/index.ts` or load from `LANGFUSE_MCP_PLUGINS` env list via `import()`), plus a `PluginManifest` with `name`, `version`, `peerLangfuseVersion`, and semver checks. Keep the `McpFeatureModule` interface (`web/src/features/mcp/server/registry.ts:52`) as the contract but stabilize it (remove `any` from `ToolHandler<any>`, add explicit error types).
- **Add lifecycle hooks:** Introduce `onRegister`, `onEnable`, `onDisable`, `onUnregister` to `McpFeatureModule`, plus `toolRegistry.unregister()` and hot-reload for dev. Add a `clear()` for tests and a `/api/public/mcp/health` that lists registered features/tool counts (currently only internal `getFeatureCount`/`getToolCount` at `registry.ts:217-226`).
- **Isolation and observability:** Run handlers with per-tool timeouts (e.g., `Promise.race` + `AbortSignal`), per-tool metrics (latency, error rate via OpenTelemetry already used in `withMiddlewares.ts:17`), and structured error-translation that preserves `UserInputError` vs `ApiServerError` (`web/src/features/mcp/core/errors.ts:19-49`) without leaking internals. Consider `worker_threads` or separate `worker` queue for expensive tools (already has `expensiveHint` in `define-tool.ts:47` but no enforcement).
- **Open the action/integration union:** Replace closed `ActionTypeSchema` enum (`packages/shared/src/domain/automations.ts:43`) with a string-keyed registry + Zod discriminated union built at runtime, allowing enterprise or community packages (`ee/`) to register new `Action` types without forking `shared`. Similarly, make `QueueName` extensible or introduce a `PluginQueue` namespace.
- **Documentation and stability:** Publish a "Writing an MCP feature" guide (beyond the 3-line bootstrap comment), pin `MCP_SERVER_VERSION` to semver, add a compatibility matrix, and add contract tests that snapshot `ToolDefinition` JSON Schemas to detect breaking changes. Reference Fern for REST versioning as a model.
- **UI extension point:** No UI plugin system exists. If needed, consider Next.js plugin slots (e.g., `ee/` already overlays `web`) and expose a `web/src/features/*/index.ts` component registry analogous to the MCP registry, with layer-based isolation via the existing `Layer` system (`web/AGENTS.md` layer guidance).

## Questions / Gaps

- **No third-party plugin authoring docs or examples:** Search finds no `CONTRIBUTING` guidance for plugins, no `examples/` for MCP/RPC plugins, and no `fern` or `README` describing a stable plugin contract. Search boundary: repo root, `web/src/features/mcp/`, `packages/shared/`, `ee/`, `AGENTS.md`.
- **No dynamic loader or marketplace:** No code, config key, or env var for discovering external plugins. Confirmed absence via grep for `plugin`/`extension` (only ESLint and file-extension hits) and inspection of `registry.ts`, `bootstrap.ts`, `mcpServer.ts`, `web/src/server/api/root.ts`, `ee/package.json`.
- **No plugin isolation tests:** No tests for registry isolation, handler sandboxing, or per-tool failure containment. Worker webhook tests (`worker/src/__tests__/webhooks.test.ts:35`) cover outbound delivery, not in-process plugin isolation.
- **No versioning/compatibility policy:** `MCP_SERVER_VERSION = "0.3.0-unstable"` (`web/src/features/mcp/server/mcpServer.ts:29`) and the self-describing warning imply intentional instability; no evidence of a stability or deprecation policy for the `McpFeatureModule` or `ToolDefinition` contracts.
- **Enterprise extensibility unclear:** `ee/` package exists (`ee/AGENTS.md`, `ee/package.json`) and `web → @langfuse/ee` dependency is allowed, but no evidence that `ee` can register MCP features or automation actions without modifying `web/src/features/mcp/server/bootstrap.ts` or `packages/shared/src/domain/automations.ts`. The mechanism, if any, is not documented.

---

Generated by `21.01-plugin-and-extension-points` against `langfuse`.
