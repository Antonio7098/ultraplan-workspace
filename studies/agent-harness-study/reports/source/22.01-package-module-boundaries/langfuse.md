# Source Analysis: langfuse

## Package and Module Boundaries

### Source Info

| Field | Value |
|-------|-------|
| Name | langfuse |
| Path | `studies/agent-harness-study/sources/langfuse` |
| Language / Stack | TypeScript (Node 24, Next.js 16, React 19, tRPC v11, BullMQ, Prisma, ClickHouse) |
| Analyzed | 2026-08-21 |

## Summary

Langfuse is organized as a pnpm + Turbo monorepo with three application packages (`web`, `worker`, plus the `ee` workspace package), four shared infrastructure packages under `packages/`, and a small set of meta packages (ESLint config, TS config, custom ESLint plugin). The repository explicitly declares a one-way dependency direction in `.agents/AGENTS.md:42-50` (`web -> @langfuse/shared, @langfuse/ee`; `worker -> @langfuse/shared`; `@langfuse/ee -> @langfuse/shared`; `@langfuse/shared -> no imports from web, worker, or ee`) and the source tree consistently honors it: `@langfuse/shared` never imports from `web`, `worker`, or `ee` (verified by grep), and the worker never imports from `web`. Internal vs public API separation is handled at three layers — the `package.json#exports` map for `@langfuse/shared` (`packages/shared/package.json:17-66`), dedicated barrel files (`packages/shared/src/index.ts` for client-safe types, `packages/shared/src/server/index.ts` for server-only utilities, `packages/shared/src/db.ts` for Prisma), and per-package AGENTS.md docs that prescribe which import path to use in which context. The EE split is unusual: the `ee/` workspace package itself contains only ~14 lines of code (`ee/src/env.ts:1-9`, `ee/src/ee-license-check/index.ts:1-5`, empty `ee/src/index.ts`) and its sole export `isEeAvailable` is not imported anywhere in the repo; the actual EE business logic lives in colocated `web/src/ee/` (`web/src/ee/features/*`) and `worker/src/ee/` (`worker/src/ee/*`) subtrees. The boundary discipline is mostly enforced by convention and AGENTS.md guidance rather than by automated tools — no eslint rule prevents shared code from importing web/worker/ee paths, and no `dependency-cruiser`/`madge`-style boundary graph is run in CI. Module-internal cohesion within `@langfuse/shared` is high and well-bucketed (separate folders for `repositories`, `queries`, `services`, `auth`, `evals`, `ingestion`, `redis`, `clickhouse`, `ee`, `otel`), but the `server/index.ts` barrel re-exports roughly 90+ modules, exposing essentially all of it as a single API surface.

## Rating

**8 / 10 — Clear model with documented exports, explicit runtime/UI separation, subpath barrels for client/server/db/env, and queue contracts owned by shared. Deductions: (a) the `ee/` workspace package is a near-empty stub whose only declared export is unused, making the documented `web -> @langfuse/ee` dependency arrow misleading; (b) the `server/index.ts` barrel is essentially a wildcard re-export, blurring internal vs public API distinctions inside `@langfuse/shared`; (c) boundary enforcement is conventional (AGENTS.md) rather than automated (no eslint/cruiser boundary guard).**

## Evidence Collected

| Area | Evidence | File:Line |
|------|----------|-----------|
| Workspace package manifest | pnpm workspace includes `web`, `worker`, `packages/**`, `ee` | `pnpm-workspace.yaml:1-6` |
| Root package declares itself `private` and uses Turbo for task orchestration | top-level name/version/private/scripts | `package.json:2-3,6,27-44` |
| Shared package version + `private: true` | `@langfuse/shared` workspace package | `packages/shared/package.json:2-5` |
| EE package version + `private: true` | `@langfuse/ee` workspace package | `ee/package.json:2-3,5` |
| Web workspace deps include `@langfuse/ee` and `@langfuse/shared` | `web/package.json` dependencies | `web/package.json:56-57` |
| Worker workspace deps include `@langfuse/shared` (no `ee`) | `worker/package.json` dependencies | `worker/package.json:35` |
| EE workspace deps include `@langfuse/shared` only | `ee/package.json` dependencies | `ee/package.json:29` |
| Documented one-way dependency direction | root AGENTS.md "Dependency direction" | `.agents/AGENTS.md:42-50` |
| Queue payload schemas + queue-name enum owned by shared | file purpose | `packages/shared/src/server/queues.ts:1-604` |
| `@langfuse/shared` declares a curated `exports` map with 9 subpaths | `packages/shared/package.json#exports` | `packages/shared/package.json:17-66` |
| Server-only barrel for queue helpers, repositories, Redis, ClickHouse, etc. | barrel | `packages/shared/src/server/index.ts:1-157` |
| Client-safe root barrel for cross-runtime types and zod schemas | barrel | `packages/shared/src/index.ts:1-111` |
| DB barrel is intentionally NOT re-exported from the root index (client-safety guard) | `db.ts` header comment + `index.ts` content | `packages/shared/src/db.ts:1-3`; `packages/shared/src/index.ts:1-111` |
| Narrower `exports` subpaths for `auth/apiKeys`, `encryption`, `query`, `monitors`, `chatml`, `ee/ingestionMasking` | `packages/shared/package.json#exports` | `packages/shared/package.json:34-65` |
| `web` AGENTS.md codifies which shared subpath to use in which context | "Shared Package Imports" section | `web/AGENTS.md:25-39` |
| `worker` AGENTS.md codifies the same import-path rules | "Shared Package Imports" section | `worker/AGENTS.md:23-36` |
| Shared AGENTS.md documents the full export map and barrel alignment rules | "Export Entry Points" section | `packages/shared/AGENTS.md:30-55` |
| EE AGENTS.md states EE logic lives in `web/src/ee/*` and `worker/src/ee/*`, not in the `ee/` package | "Integration Notes" section | `ee/AGENTS.md:30-33` |
| `ee/src/index.ts` is an empty barrel | file contents | `ee/src/index.ts:1` |
| `ee/src/ee-license-check/index.ts` exports `isEeAvailable` | declaration | `ee/src/ee-license-check/index.ts:1-5` |
| `isEeAvailable` has zero importers in the entire repo | grep result | `grep -r isEeAvailable .` returns only the defining file |
| EE business logic colocated inside `web/src/ee/features/*` | `admin-api`, `audit-log-viewer`, `billing`, `in-app-agent`, `multi-tenant-sso`, `sso-settings`, `ui-customization`, `verified-domains` | `web/src/ee/features/` directory listing |
| EE business logic colocated inside `worker/src/ee/*` | `cloudSpendAlerts`, `cloudUsageMetering`, `dataRetention`, `meteringDataPostgresExport`, `usageThresholds` | `worker/src/ee/` directory listing |
| Worker imports `prisma` exclusively via the dedicated `db` subpath | cross-package grep | many files under `worker/src/**` import `@langfuse/shared/src/db` |
| Web imports server-only utilities via the dedicated `server` subpath | cross-package grep | many files under `web/src/**` import `@langfuse/shared/src/server` |
| Web imports frontend-safe types via the root `@langfuse/shared` barrel | `web/src/features/public-api/types/observations.ts` | `web/src/features/public-api/types/observations.ts:1-22` |
| Shared EE ingestion-masking helper is exported via a dedicated subpath | `packages/shared/package.json#exports` and implementation | `packages/shared/package.json:46-49`; `packages/shared/src/server/ee/ingestionMasking/applyIngestionMasking.ts` |
| Worker uses the `@langfuse/shared/src/server/ee/ingestionMasking` subpath explicitly | `worker/src/queues/otelIngestionQueue.ts:27` import block | `worker/src/queues/otelIngestionQueue.ts:27` |
| Worker has zero imports of `web/` or `@langfuse/web` | grep result | empty |
| `@langfuse/shared` has zero reverse imports of `web/`, `worker/`, or `ee/` | grep result | empty |
| `BaseError` hierarchy is owned by shared and re-exported via root barrel | `errors/BaseError.ts` and `errors/index.ts` | `packages/shared/src/errors/BaseError.ts:1-27`; `packages/shared/src/errors/index.ts:1-13` |
| Public REST + tRPC + EE router registration centralizes routing on the web side | tRPC root + public routes | `web/src/server/api/root.ts:14-61`; `web/src/server/api/trpc.ts:1-50`; `web/src/features/public-api/server/withMiddlewares.ts:1-80` |
| Public API contract types are isolated under `web/src/features/public-api/types/` | 18 type files | `web/src/features/public-api/types/` directory listing |
| Fern is the source of truth for public API; generated clients are excluded from the repo | `fern/apis/` + repo-wide generated ignore | `.agents/AGENTS.md:77-86`; `fern/apis/{client,organizations,server}/` |
| Web Next config transpiles `@langfuse/shared` from `src` for dev (turbopack alias) | `next.config.mjs` | `web/next.config.mjs:61,79-81` |
| Vitest workspace is limited to `web` and `worker` (no shared test project) | `vitest.workspace.ts` | `vitest.workspace.ts:1` |
| Turbo tasks declare build/typecheck/test graph and `@repo/eslint-plugin#build` as an internal dep | turbo task definitions | `turbo.json:10-105` |
| Workspace meta packages: shared TS config, shared ESLint config, custom ESLint plugin | `packages/{config-eslint,config-typescript,eslint-plugin}` | `packages/config-typescript/package.json:1-9`; `packages/eslint-plugin/package.json:1-*`; `packages/eslint-plugin/src/rules/{no-in-source-vitest,no-style-props,no-tailwind-overflow-scroll}.ts` |
| No boundary-enforcing eslint rule in `packages/config-eslint/base.js` or `next.js` | rule list | `packages/config-eslint/base.js:46-103`; `packages/config-eslint/next.js:84-121` |
| Web's only `no-restricted-imports` rule is about icon libraries, not package boundaries | `web/eslint.config.mjs` | `web/eslint.config.mjs` `no-restricted-imports` block |
| Base eslint only forbids the bare global `redis` import (forces `redis` from shared) | `packages/config-eslint/base.js` `no-restricted-globals` | `packages/config-eslint/base.js:72-78` |

## Answers to Dimension Questions

### 1. Are modules cleanly separated?

Yes at the coarse package level; partially at the sub-module level. The four workspace packages (`web`, `worker`, `packages/shared`, `ee`) have crisp, non-overlapping roles — Next.js app, queue worker, shared domain/contracts, EE license helper. Within `@langfuse/shared`, modules are clearly bucketed into `repositories/`, `queries/`, `services/`, `auth/`, `evals/`, `ingestion/`, `redis/`, `clickhouse/`, `ee/`, `otel/`, `llm/`, `data-deletion/`, `dataset-run-items/`, `webhooks/`, `pricing-tiers/`, `s3/`, `instrumentation/`, `outbound-url/`, `services/email/` etc. — a quick directory walk reveals stable, single-purpose groupings. However, two separation problems remain: (a) the `server/index.ts` barrel re-exports ~90+ modules essentially as a wildcard surface (`packages/shared/src/server/index.ts:1-157`), so the distinction between "internal" and "external" within `@langfuse/shared/server` is not enforced; anything exported anywhere under `src/server` becomes part of the public-ish server surface. (b) The `ee/` workspace package is essentially a stub (see `ee/src/index.ts:1` being empty and `ee/src/ee-license-check/index.ts:1-5` exporting a single unused symbol), but `web/package.json:56` still declares it as a dependency, producing a documented-but-stale `web -> @langfuse/ee` arrow in `.agents/AGENTS.md:42-50`. The actual EE business logic is colocated inside `web/src/ee/` and `worker/src/ee/` subtrees — a third arrangement not captured by the doc's package graph.

### 2. Do dependencies flow in one direction?

Yes, at the package level. Verified by negative grep: `@langfuse/shared` contains zero imports from `web`, `worker`, or `ee`; `worker` contains zero imports from `web`. The declared graph is `web -> shared, ee`; `worker -> shared`; `ee -> shared`. This is consistent in code (the imports I sampled — `worker/src/queues/ingestionQueue.ts:1-23`, `worker/src/app.ts:9-99`, `web/src/server/api/trpc.ts:1-94`, `web/src/initialize.ts:1-12`, `web/src/server/auth.ts:45-58` — only consume `@langfuse/shared`, `@langfuse/shared/src/server`, `@langfuse/shared/src/db`, or `@langfuse/shared/src/server/auth/apiKeys`; nothing inverts the arrow). Inside `@langfuse/shared` there are tighter intra-package direction hints — e.g. `db.ts` is explicitly excluded from `src/index.ts` to keep Prisma out of client bundles (`packages/shared/src/db.ts:1-3`), and `src/server/index.ts` is the server-only barrel. The "EE inside web" arrangement does introduce a second-order one-way convention: `web/src/` non-EE code may import from `web/src/ee/features/*` (e.g. `web/src/features/entitlements/server/getPlan.ts:1` imports from `@/src/ee/features/billing/utils/stripeCatalogue`), but the EE folder does not import from the non-EE `web/src/features/` subtree — only from `@langfuse/shared` and local `@/src/ee` paths.

### 3. Can modules be used independently?

- `@langfuse/shared` is the canonical example. It builds independently to `dist/` (`packages/shared/package.json:8-11,67-72`), exposes a curated set of subpath exports (`packages/shared/package.json:17-66`), and its `encryption` module (`packages/shared/src/encryption/index.ts`), `chatml` module (`packages/shared/src/utils/chatml/index.ts`), `query` (`packages/shared/src/features/query/index.ts`), and `monitors` (`packages/shared/src/features/monitors/index.ts`) all have their own narrow subpaths that consumers can pull without dragging the full server runtime. The `env` and `db` subpaths (`packages/shared/src/env.ts:1-80`, `packages/shared/src/db.ts:1-62`) are similarly usable in isolation. Note that `@langfuse/shared/src/server` re-exports a wide barrel — pulling one symbol from it can pull a lot of transitive code via `export *` chains.
- `@langfuse/ee` cannot be meaningfully used: it has only `env.ts` (9 lines) and `ee-license-check/index.ts` (5 lines, exporting `isEeAvailable` which no one imports). A consumer pulling `@langfuse/ee` gets ~14 lines of TypeScript.
- `web` and `worker` each ship as Dockerized applications (`web/Dockerfile`, `worker/Dockerfile`); they are independently runnable but `worker` requires shared Prisma schema and ClickHouse tables from `@langfuse/shared`, so worker is runnable-but-coupled.
- The meta packages (`packages/config-eslint`, `packages/config-typescript`, `packages/eslint-plugin`) are fully independent dev tooling, with no runtime impact.

### 4. Are public APIs distinguished from internal ones?

Three layers of distinction are in place:

1. **Subpath exports map.** `@langfuse/shared`'s `package.json#exports` (`packages/shared/package.json:17-66`) lists exactly 9 surfaces: `.` (client-safe types), `./src/db` (Prisma, server-only), `./src/env` (env, server-only), `./src/server` (server-only utilities), `./src/server/auth/apiKeys`, `./src/utils/chatml`, `./encryption`, `./src/server/ee/ingestionMasking`, `./query`, `./query/server`, `./monitors`, `./monitors/server`. Anything not in this list is unimportable from outside the package — this is enforced by Node's package-exports resolution at runtime.

2. **Barrel discipline within the package.** `src/index.ts` (`packages/shared/src/index.ts:1-111`) deliberately omits `./db` (`packages/shared/src/db.ts:1-3` and the comment "This is not imported in the index.ts file of this package, as we must not import this into FE code"), omits `./server` (kept separate for server use), and omits `./env` (separate subpath). The contract is documented in `packages/shared/AGENTS.md:30-55`.

3. **Per-package import-path conventions.** `web/AGENTS.md:25-39` and `worker/AGENTS.md:23-36` prescribe the exact subpath to use in each context (frontend-safe, server-only, direct Prisma, focused helpers). This is convention, not lint-enforced, but it is applied consistently in the code I sampled.

The `web` package further distinguishes public API contracts by isolating type definitions under `web/src/features/public-api/types/` (18 type files for endpoints like observations, scores, traces, sessions, etc.) and routing all public REST endpoints through a single `withMiddlewares` helper (`web/src/features/public-api/server/withMiddlewares.ts:1-80`) plus `createAuthedProjectAPIRoute` (`web/src/features/public-api/server/createAuthedProjectAPIRoute.ts`).

Where the boundary is weaker: `@langfuse/shared/src/server/index.ts:1-157` is a near-wildcard barrel — almost every file under `src/server/**` is re-exported. This means there is effectively no "internal-to-shared-server" boundary: a symbol used internally is by default reachable via `@langfuse/shared/src/server`. The barrel also re-exports some modules whose own internal helpers are reachable, e.g. `./services/DashboardService/types` is re-exported (`packages/shared/src/server/index.ts:132`) but the types module is also reachable from the client-safe root barrel via `packages/shared/src/index.ts:102-106`, so `DashboardService` types leak across both surfaces.

## Architectural Decisions

- **Single shared package for everything cross-runtime.** Instead of splitting `shared-types`, `shared-server`, `shared-domain` into separate packages, Langfuse keeps one `@langfuse/shared` package with two barrels (`src/index.ts`, `src/server/index.ts`) and a dedicated `db.ts` (`packages/shared/src/index.ts:1-111`, `packages/shared/src/server/index.ts:1-157`, `packages/shared/src/db.ts:1-62`). This trades fine-grained encapsulation for a single import-graph and a single set of build outputs.
- **Queue payload contracts owned by shared.** `packages/shared/src/server/queues.ts:1-604` defines both the `IngestionEvent`, `OtelIngestionEvent`, `BatchExportJobSchema`, `TraceQueueEventSchema`, `DatasetQueueEventSchema`, etc., and the `QueueName` enum. Both web and worker import them; this is the canonical boundary between the two runtimes.
- **EE code colocated, not extracted.** Despite having an `ee/` workspace package, EE business code lives in `web/src/ee/features/*` and `worker/src/ee/*` subtrees, behind runtime feature gates like `isEeAvailable`/`isEnterpriseLicenseAvailable` (`packages/shared/src/server/ee/licenseCheck/index.ts:1-26`). The `ee/` workspace package is a license-check stub (`ee/src/ee-license-check/index.ts:1-5`).
- **Public REST contract as the source of public API truth, with Fern as upstream.** `web/src/features/public-api/server/` and `web/src/features/public-api/types/` define the runtime contract; `fern/apis/{client,organizations,server}/definition/*.yml` defines the OpenAPI/AsyncAPI source; `generated/` is the build output that is gitignored/excluded per `.agents/AGENTS.md:77-86`.
- **Next.js Pages Router with `transpilePackages: ["@langfuse/shared", ...]`.** `web/next.config.mjs:61` plus `web/next.config.mjs:79-81` turbopack alias to `./packages/shared/src` mean the web container treats shared as a sibling source tree in dev. In production, both packages build independently (`packages/shared/package.json:67-72` runs `tsc`; `web/package.json:10` runs `next build` after `^build`).
- **Vitest workspace scoped to `web` and `worker`.** `vitest.workspace.ts:1` lists only `["web", "worker"]`. Tests in `packages/shared/src/server/automations.test.ts`, `packages/shared/src/server/queues.test.ts`, `packages/shared/src/features/monitors/types.test.ts`, `packages/shared/src/features/query/server/queryBuilder.intervalBucketing.test.ts` etc. run via `pnpm --filter @langfuse/shared run test` (`packages/shared/package.json:74`) and not from the root vitest workspace.
- **Errors as a shared subsystem.** The `BaseError` hierarchy in `packages/shared/src/errors/BaseError.ts:1-27` (with `LangfuseNotFoundError`, `InvalidRequestError`, `UnauthorizedError`, `ForbiddenError`, `MethodNotAllowedError`, `InternalServerError`, `LangfuseConflictError`, `ServiceUnavailableError`, `NotImplementedError`, `PayloadTooLargeError`) is exported from both the root barrel (`packages/shared/src/index.ts:89`) and consumed directly by `web/src/features/public-api/server/withMiddlewares.ts:6-10`.

## Notable Patterns

- **Curated `exports` map with mixed client/server subpaths.** `@langfuse/shared` is the canonical example: 9 subpaths split by audience (`packages/shared/package.json:17-66`). Each subpath is documented in the package's AGENTS.md (`packages/shared/AGENTS.md:30-55`).
- **Tighter import paths than necessary to express "server-only".** Worker code imports `prisma` only via `@langfuse/shared/src/db` (e.g. `worker/src/queues/ingestionQueue.ts:23`) rather than `@langfuse/shared`, enforcing at the call-site that Prisma is never pulled into a frontend context.
- **Shared domain types separated from shared infrastructure.** `packages/shared/src/domain/` (traces, observations, scores, score-configs, dataset-items, dataset-run-items, webhooks, prompts, automations, table-view-presets, observation-field-groups) is distinct from `packages/shared/src/server/repositories/` (the ClickHouse/Prisma-backed implementations). Domain types live in the client-safe barrel; repositories live in the server-only barrel.
- **Module-per-area in shared server.** One folder per subsystem: `auth/`, `evals/`, `ingestion/`, `redis/`, `clickhouse/`, `data-deletion/`, `dataset-run-items/`, `datasets/`, `llm/`, `outbound-url/`, `pricing-tiers/`, `queries/`, `services/`, `tableMappings/`, `webhooks/`, `ee/`, `otel/`, `instrumentation/`, `s3/`, `media-deletion.ts`, `traceDeletionProcessor.ts`, `deletionGuard.ts`. Each has its own index.
- **Webpack web worker for browser-side JSON parsing** that pulls a shared utility (`web/src/workers/json-parser.worker.ts:8` imports `deepParseJsonIterative` from `@langfuse/shared`) — the shared utility is reusable across the main bundle and a Worker bundle.
- **Server-only and client-safe sub-paths for the same feature.** `@langfuse/shared/query` (`packages/shared/src/features/query/index.ts`) is the client-safe surface; `@langfuse/shared/query/server` (`packages/shared/src/features/query/server/index.ts`) is the server-only builder/executor; same split for `@langfuse/shared/monitors` vs `@langfuse/shared/monitors/server` (`packages/shared/package.json:50-65`).
- **Inline license gating via `isEnterpriseLicenseAvailable`.** `packages/shared/src/server/ee/licenseCheck/index.ts:11-26` is a pure function over `SharedEnv` that decides whether EE features are available. It is used only inside `packages/shared/src/server/ee/ingestionMasking/applyIngestionMasking.ts` (per grep) — confirming the gate is centralized rather than scattered.
- **Meta packages at `packages/config-*` and `packages/eslint-plugin`.** A `@repo/typescript-config/base.json` shared by all four packages, a `@repo/eslint-config` flat config, and a tiny `@repo/eslint-plugin` (`packages/eslint-plugin/src/rules/{no-in-source-vitest,no-style-props,no-tailwind-overflow-scroll}.ts`) for repo-specific rules.

## Tradeoffs

- **Single barrel `server/index.ts` re-exports everything.** Pro: easy import path for `web` and `worker` code (one symbol, one import). Con: there is no internal boundary inside `@langfuse/shared/src/server` — any helper exported anywhere under `src/server/` becomes reachable. Pulling `@langfuse/shared/src/server` can pull a lot of code transitively.
- **EE is "in web and worker", not in the `ee/` package.** Pro: EE code can colocate with the OSS code it shadows (e.g. `web/src/ee/features/billing/server/stripeBillingService.ts:22-33` next to `web/src/features/organizations/`). Con: the `ee/` package exists as a workspace dependency that effectively no code uses, producing a stale-looking dependency arrow and a separate tsconfig/lint pipeline that runs but does very little.
- **No automated boundary checks.** The base eslint config forbids the bare global `redis` (`packages/config-eslint/base.js:72-78`) and the web config forbids `react-icons` other than `si`/`tb` (`web/eslint.config.mjs` `no-restricted-imports` block), but no rule enforces "do not import `web/` from `@langfuse/shared`" or "do not import `@langfuse/shared/src/server` from a client component." The rule set relies on convention plus per-package AGENTS.md guidance.
- **Generated clients are not in the repo.** Per `.agents/AGENTS.md:77-86`, `generated/*` is excluded; the Fern spec in `fern/apis/` is the source of truth. This means the public API contract is not directly inspectable in a fresh checkout — the runtime contract types in `web/src/features/public-api/types/` and the routes in `web/src/pages/api/public/` are the only in-repo expression of it.
- **Vitest workspace scoped to two packages.** Tests in `@langfuse/shared` are run via `pnpm --filter @langfuse/shared run test`, not from the root vitest workspace (`vitest.workspace.ts:1`). This is fine for module isolation but means there is no single command that runs the entire repo's tests, and no cross-package test fixture sharing.
- **Server packages include substantial client-unsafe dependencies.** `packages/shared/package.json:96-148` lists `@anthropic-ai/vertex-sdk`, `@aws-sdk/client-s3`, `@aws-sdk/client-cloudwatch`, `@azure/storage-blob`, `@clickhouse/client`, `@google-cloud/storage`, `@slack/oauth`, `@slack/web-api`, `bullmq`, `dd-trace`, `nodemailer`, `oci-common`, `oci-objectstorage`, `stripe`, etc. None of these leak into client bundles only because the subpath boundary (`@langfuse/shared` vs `@langfuse/shared/src/server`) is respected at import sites.

## Failure Modes / Edge Cases

- **If anyone imports `db.ts` via the root barrel, Prisma ships to the client.** The current barrel discipline (`packages/shared/src/index.ts:1-111` does not re-export `./db`) prevents this, but nothing automated enforces it; a future barrel change could regress it silently.
- **The `ee/` workspace package is dead weight today.** A change to `ee/src/index.ts` adding new exports is undetected because no consumer imports it. AGENTS.md's promise of "EE package consumed by web and worker" (`ee/AGENTS.md:7`) is true only for the workspace dependency edge, not for runtime code.
- **`server/index.ts` barrel growth.** As more server-only modules are added under `src/server/`, they all become part of one massive re-export. A symbol renamed or moved out of an internal helper file may still be reachable via the barrel, complicating deprecation.
- **No detection of cross-package test isolation regressions.** Vitest project boundaries are encoded in `vitest.workspace.ts` (only `web`, `worker`); `@langfuse/shared` tests are run only via `--filter`. If a shared test starts importing `web/` paths, there is no automated guard.
- **Next.js transpile + turbopack alias.** `web/next.config.mjs:79-81` aliases `@langfuse/shared` to `./packages/shared/src` for dev. If a developer makes a local edit to `packages/shared/src/**`, it propagates immediately to the web container without a rebuild — which is good for iteration, but means the dev experience can drift from production (which uses `dist/`).
- **Turbopack-only alias.** The turbopack alias in `web/next.config.mjs:79-81` covers dev; production webpack builds rely on `transpilePackages: ["@langfuse/shared", ...]` at `web/next.config.mjs:61` and on the built `dist/`. Two build paths means two resolution paths to keep aligned.

## Future Considerations

- **Consolidate or retire `ee/`.** Either move EE logic into the `ee/` package (so the documented dependency arrow becomes real) or delete the package and update `.agents/AGENTS.md:42-50` and `web/package.json:56` to remove the misleading edge.
- **Add automated boundary enforcement.** A lightweight dependency-cruiser or eslint `no-restricted-imports` rule keyed on package paths ("`packages/shared/**` may not import from `web/**`, `worker/**`, or `ee/**`") would convert the documented contract into a checkable one.
- **Split `@langfuse/shared/src/server` into smaller barrels.** Group re-exports by audience (e.g. `./server/queue`, `./server/repository`, `./server/llm`) and let consumers pick. Reduces accidental surface area growth.
- **Promote internal-only modules in shared to a non-exported path.** Files that are only consumed by other shared files can move to e.g. `packages/shared/src/internal/**` and be excluded from the server barrel.
- **Wire `vitest.workspace.ts` to include shared.** Let the root `pnpm test` cover shared tests too, or document explicitly that shared tests must be invoked via `pnpm --filter @langfuse/shared run test`.
- **Versioned exports for the public REST contract.** The runtime types in `web/src/features/public-api/types/` and the Fern spec in `fern/apis/server/definition/` already exist; adding a CLI check that all types are present in Fern (or vice-versa) would harden the contract.

## Questions / Gaps

- **Is the `ee/` package intentionally a stub, or is there a planned migration of EE code from `web/src/ee/` and `worker/src/ee/` into it?** `ee/AGENTS.md:7` says it is "consumed by `web` and `worker`", but no runtime code consumes it. The doc could be stale or aspirational.
- **What is the contract between `packages/shared/src/features/query/types.ts` (client-safe) and `packages/shared/src/features/query/server/queryBuilder.ts` (server-only)?** They are exported under separate subpaths (`@langfuse/shared/query` vs `@langfuse/shared/query/server`) — is there an explicit contract test, or are they kept aligned by convention?
- **Is `@langfuse/shared/src/services/email/*` intentionally shared with `web` even though only `worker` sends email?** AGENTS.md says worker should prefer `@langfuse/shared/src/server` (`worker/AGENTS.md:23-36`), and indeed worker does (`worker/src/queues/notifications/handleNotificationJob.ts` etc. — not inspected in detail). No evidence of web importing the email services, but no rule prevents it.
- **No automated detection of `web/AGENTS.md:25-39` or `worker/AGENTS.md:23-36` rule violations.** "Use `@langfuse/shared/src/db` only in backend or test code; never route it into client bundles" — is there a lint rule or a test that verifies client bundles do not transitively pull in Prisma?
- **Is there a desired long-term boundary between `web/src/ee/features/*` and `web/src/features/*`?** The cross-import pattern (`web/src/features/entitlements/server/getPlan.ts:1` imports from `@/src/ee/features/billing/utils/stripeCatalogue`) suggests OSS code is allowed to read from EE. Does the team want this to remain one-way (EE never imports from OSS) or is two-way allowed?
- **Where does `@langfuse/shared/src/server/services/DashboardService/types` actually belong?** It is exported both via the server barrel (`packages/shared/src/server/index.ts:132`) and via the client-safe root barrel's selective re-export (`packages/shared/src/index.ts:102-106` for `ChartConfigSchema`, `DimensionSchema`, `MetricSchema`). The same module straddles both surfaces, which blurs the boundary.
- **Why is `vitest.workspace.ts:1` limited to `web` and `worker` rather than also including `packages/shared`?** No clear rationale found in AGENTS.md or `vitest.workspace.ts`.