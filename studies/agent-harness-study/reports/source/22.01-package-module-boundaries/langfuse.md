# Source Analysis: langfuse

## Dimension 22.01: Package and Module Boundaries

### Source Info

| Field | Value |
|-------|-------|
| Name | langfuse |
| Path | `studies/agent-harness-study/sources/langfuse` |
| Language / Stack | TypeScript / Node 24; pnpm workspace + Turborepo; Next.js (web), BullMQ worker, Prisma + ClickHouse |
| Analyzed | 2026-08-21 |

All citations below are relative to the source root (`sources/langfuse/`).

## Summary

Langfuse is a pnpm + Turborepo monorepo with four runtime packages — `web` (Next.js UI + tRPC + public REST), `worker` (queue consumers), `@langfuse/shared` (domain models, DB clients, queue contracts, repositories, server services), and `@langfuse/ee` (enterprise) — plus three tooling packages under `@repo/*` (eslint config, tsconfig, custom eslint plugin). The declared dependency graph is a clean DAG (`AGENTS.md:45-49`): `shared` depends on nothing internal; `ee` and `worker` depend only on `shared`; `web` depends on `shared` + `ee`. Verified by inspecting every `package.json`: no package imports upward, so circular dependencies are structurally impossible at the package level.

The main boundary mechanism is a two-tier export surface in `@langfuse/shared`: a client-safe root barrel (`packages/shared/src/index.ts`, 111 lines) versus a server-only barrel (`packages/shared/src/server/index.ts`, 157 lines), exposed through an explicit `exports` map (`packages/shared/package.json:17-66`). Queue contracts between the web producer and worker consumer are centralized as zod schemas + a `QueueName` enum in `packages/shared/src/server/queues.ts:342`.

Weaknesses: boundary safety inside `web` is convention-based (server-only imports in page files rely on Next.js stripping `getServerSideProps` code from client bundles rather than lint/compile-time checks); the `@langfuse/ee` package is vestigial — declared as a web dependency but imported nowhere, with an empty `index.ts` and a stale `./sso` export pointing at source that does not exist; EE feature code actually lives in `web/src/ee/**` and `worker/src/ee/**` directories. There are no dedicated separation tests or cycle-detection tooling.

## Rating

**6 / 10** — Present but inconsistent. The model is clearly documented (`AGENTS.md:45-49`) and structurally sound (verified acyclic DAG, explicit `exports` maps, two-tier barrels), but enforcement is mostly lint/typecheck/convention: no boundary tests, no cycle-detection tooling, a dead `ee` package whose manifest contradicts reality, and server/client separation inside `web` that depends on framework behavior instead of static guarantees. The rubric's 7–8 band requires "tests [and] operational safeguards" for boundaries, which are absent.

## Evidence Collected

| Area | Evidence | File:Line |
|------|----------|-----------|
| Workspace definition | pnpm workspace members: `web`, `worker`, `packages/**`, `ee`; supply-chain pinning via `minimumReleaseAge` | `pnpm-workspace.yaml:1-8` |
| Task graph | Turbo `build.dependsOn: ["db:generate", "^build"]` — topological build order across packages | `turbo.json:11-16` |
| Documented dependency direction | `web -> shared, ee`; `worker -> shared`; `ee -> shared`; `shared -> nothing` | `AGENTS.md:45-49` |
| web deps | `"@langfuse/ee": "workspace:*"`, `"@langfuse/shared": "workspace:*"` | `web/package.json:56-57` |
| worker deps | only `"@langfuse/shared": "workspace:*"` among internal packages | `worker/package.json:35` |
| ee deps | `"@langfuse/shared": "workspace:*"`; also couples ee to `next` and `next-auth` | `ee/package.json:29-33` |
| shared has no upward deps | runtime deps contain no `workspace:*` entries; only `@repo/*` devDeps | `packages/shared/package.json:96-163` |
| No reverse imports | grep of `packages/shared/src` for `@langfuse/ee`, `web/`, `worker/` imports returned zero hits (searched `*.ts`, `*.tsx`) | `packages/shared/src/**` (negative result) |
| Explicit public surface | `exports` map with narrow subpaths: `.`, `./src/db`, `./src/env`, `./src/server`, `./encryption`, `./query`, `./monitors`, etc. | `packages/shared/package.json:17-66` |
| Two-tier barrel split | client-safe root barrel vs server-only barrel, documented as such in package guide | `packages/shared/src/index.ts:1-111`; `packages/shared/src/server/index.ts:1-157`; `packages/shared/AGENTS.md` ("Export Entry Points") |
| Server/client import discipline | 218 files under `web/src` import the client-safe root barrel; zero `.tsx` under `web/src/components` import the server barrel | `web/src/**` (counts via grep) |
| Type-only cross-barrel imports | client components use `import { type X }` so server types don't reach bundles | `web/src/features/datasets/components/DatasetAggregateTableCell.tsx:11`; `web/src/components/table/use-cases/models.tsx:6` |
| Shared queue contract ownership | zod payload schemas + `QueueName` enum + per-queue options map owned by shared, consumed by both producer and consumer | `packages/shared/src/server/queues.ts:21-42,342,420+` |
| Consumer uses shared contracts | worker ingestion processor typed by `TQueueJobTypes[QueueName.IngestionQueue]` imported from shared server barrel | `worker/src/queues/ingestionQueue.ts:3-24,40` |
| Boundary-ish lint rule | global `redis` banned; message forces explicit import path `'@langfuse/shared/src/server'` | `packages/config-eslint/base.js:72-76` |
| Import restriction (UI layer) | `no-restricted-imports` limits icon libraries in web | `web/eslint.config.mjs:36-51` |
| Custom plugin scope | `@repo/eslint-plugin` ships style/test-placement rules only — no dependency-boundary rules | `packages/eslint-plugin/src/index.ts:1-13` |
| CI enforcement | pipeline runs repo-wide `pnpm run lint` + `typecheck`; separate eslint-plugin test job | `.github/workflows/pipeline.yml:140-143,206,237-238` |
| Independent deployables | web and worker each have own Dockerfile/entrypoint | `web/Dockerfile`; `worker/Dockerfile`; `worker/entrypoint.sh` |
| Vestigial ee package | empty index (`wc -c` = 0 bytes); stale `./sso` export targets `dist/src/sso/index.js` but `ee/src/sso/` does not exist | `ee/src/index.ts`; `ee/package.json:12-14`; `find ee/src -type f` → only `ee-license-check/index.ts`, `env.ts` |
| Unused ee dependency | `"@langfuse/ee"` declared by web but zero imports found in `web/src` or `worker/src`; `isEeAvailable` also unimported | `web/package.json:56`; negative greps over `web/src`, `worker/src` |
| Real EE boundary is directories | EE-gated features live inside apps, marked by folder READMEs | `web/src/ee/README.md`; `web/src/ee/features/**`; `worker/src/ee/{cloudSpendAlerts,cloudUsageMetering,dataRetention,...}` |
| Convention-based server safety | page module imports `prisma` (runtime) and `getTracesByIdsForAnyProject` at top level, used only inside `getServerSideProps` — safe only because Next strips gSSP code from client bundles | `web/src/pages/project/[projectId]/evals/configs/[configId].tsx:1,6,16`; `web/src/pages/trace/[traceId].tsx:6,10-21` |

## Answers to Dimension Questions

1. **Are modules cleanly separated?** Mostly yes at the package level: `shared` is a pure leaf (no internal deps, verified by `package.json` contents at `packages/shared/package.json:96-163` and negative import greps), and app-level concerns are split into `web` (request/response surface) vs `worker` (queues). Within `web`, features are modularized under `src/features/*` with per-feature README contracts (e.g., `web/src/features/search-bar/README.md`). However, separation is muddied by (a) the kitchen-sink nature of `@langfuse/shared` — domain types, Prisma/ClickHouse clients, LLM execution, email rendering built on React Email (`packages/shared/src/server/index.ts:8-14` exports email services; react is even a peerDependency at `packages/shared/package.json:165-168`) — and (b) the vestigial `ee` package vs directory-based EE code inside both apps.
2. **Do dependencies flow in one direction?** Yes. The graph is a strict DAG: `shared ← {ee, worker, web}`, `web → ee`. Confirmed from manifests (`web/package.json:56-57`, `worker/package.json:35`, `ee/package.json:29`) and the absence of any upward import in `packages/shared/src`. Cycles are impossible without manifest changes, though nothing automated detects them if someone added one.
3. **Can modules be used independently?** Partially. `web` and `worker` are independently buildable/deployable (separate Dockerfiles, `turbo.json:102-105` even special-cases `worker#dev`). `@langfuse/shared` builds standalone (`tsc`, `packages/shared/package.json:68`) and its `exports` map permits narrow subpath imports (`./encryption`, `./src/utils/chatml`, `./query` at `packages/shared/package.json:38-65`). But the package is heavy: installing it drags LangChain, AWS/Azure/GCP storage SDKs, BullMQ, dd-trace, etc. (`packages/shared/package.json:96-148`), so "use one util without the whole runtime" holds only at bundle/tree-shake level, not at install level. It is `private: true`, so this matters mainly internally.
4. **Are public APIs distinguished from internal ones?** Yes, explicitly for `shared`: `main`/`types` point at the root barrel (`packages/shared/package.json:12-13`), the `exports` map whitelists entrypoints, and the package guide mandates keeping `package.json#exports`, barrels, and docs aligned in the same PR (`packages/shared/AGENTS.md`, "Export surface change" playbook). A second, softer distinction separates frontend-safe (`src/index.ts`) from server-only (`src/server/index.ts`) surfaces. For the HTTP API, Fern definitions are the contract source (`fern/apis/**`, enforced by CI guidance in `web/src/features/README.md`). Inside `web`/`worker` themselves, there is no public/internal distinction beyond folder conventions.

## Architectural Decisions

- **Monorepo with a single leaf shared package** (`pnpm-workspace.yaml:3-6`, `AGENTS.md:45-49`): all cross-app contracts (DB schema ownership `packages/shared/prisma/schema.prisma`, ClickHouse migrations, queue payloads) concentrate in `@langfuse/shared` so producer (web) and consumer (worker) cannot drift — e.g., `bucketPrefix` on `IngestionEvent` exists precisely to prevent producer/consumer env drift (`packages/shared/src/server/queues.ts:26-32`).
- **Two-tier export barrels** separating frontend-safe from server-only symbols (`packages/shared/src/index.ts` vs `packages/shared/src/server/index.ts`), with narrow subpath exports for focused needs (`packages/shared/package.json:17-66`).
- **Convention + documentation over enforcement**: the dependency direction is codified in agent/engineering guides rather than in a machine-checked rule set; the custom eslint plugin covers style rules only (`packages/eslint-plugin/src/index.ts:5-12`).
- **EE as a licensing concept, not a real package boundary**: despite `@langfuse/ee` existing in the graph, EE functionality physically resides in `web/src/ee/**` / `worker/src/ee/**` gated by `isEeAvailable` (`ee/src/ee-license-check/index.ts:3-5`) — a decision that has left the package hollow.

## Notable Patterns

- **Type-only imports as a leak valve**: client-side `.tsx` files pull server-barrel/db types with `import { type ... }`, which the compiler erases (`web/src/features/datasets/components/DatasetAggregateTableCell.tsx:11`, `web/src/features/widgets/chart-library/Chart.tsx:17`) — disciplined, but nothing enforces the `type` keyword.
- **Contract-first queues**: every BullMQ queue has a named enum member, zod payload schema, and options row in one file (`packages/shared/src/server/queues.ts:342,420+`), and processors are statically typed against it (`worker/src/queues/ingestionQueue.ts:40`).
- **Error taxonomy shared across tRPC and REST**: `BaseError` subclasses live in `packages/shared/src/errors/` and translate to HTTP codes in middlewares (documented playbook in `web/AGENTS.md`, "Error handling").
- **Repo-level guardrails in CI**: repo-wide lint/typecheck plus a dedicated job testing the custom eslint plugin (`.github/workflows/pipeline.yml:140-143,237-238`).

## Tradeoffs

- **Documentation-driven boundaries scale with discipline**: the model is excellent on paper (`AGENTS.md:45-49`, per-package `AGENTS.md` files) but costs nothing to violate silently until typecheck/lint catches a typo-level mistake — semantic violations (client importing server code that happens to typecheck) pass.
- **One big shared package vs many small ones**: concentrating everything in `@langfuse/shared` minimizes contract drift but makes it a grab-bag where changing an email template shares a version surface with queue schemas; the React peerDependency on a backend package (`packages/shared/package.json:165-168`) is a visible seam of this choice.
- **Next.js gSSP stripping as implicit isolation**: importing `prisma` directly into a page module (`web/src/pages/project/[projectId]/evals/configs/[configId].tsx:1`) keeps code out of the client bundle today, but couples correctness to framework bundling behavior rather than an explicit `server-only` guard (the `server-only` package appears exactly once, in `web/src/ee/features/in-app-agent/schema.ts`).

## Failure Modes / Edge Cases

- **Stale manifest rot**: `@langfuse/ee` declares a `./sso` subpath export targeting `dist/src/sso/index.js` while `ee/src/sso/` does not exist (`ee/package.json:12-14`); the package's `index.ts` is 0 bytes. Any consumer honoring `exports` strictly would fail on the documented entrypoint; today nobody notices because nobody imports the package.
- **Phantom dependency edge**: `web/package.json:56` keeps `@langfuse/ee` installed, so turbo still builds it (`^build` chain, `turbo.json:12`) — wasted build time and a misleading architecture diagram for newcomers.
- **Accidental client-bundle leaks**: if a developer moves a server call out of `getServerSideProps` in a page like `web/src/pages/trace/[traceId].tsx:6,21`, the prisma/clickhouse-backed import would attempt to bundle client-side; detection currently comes from build errors or review, not a rule.
- **Global-name collisions**: the ban on bare `redis` global (`packages/config-eslint/base.js:72-76`) shows a past failure mode where a hoisted Redis singleton leaked across modules; `publicHoistPattern` for `*prisma*` (`pnpm-workspace.yaml:82-83`) keeps similar pressure alive.

## Future Considerations

- Delete or re-materialize `@langfuse/ee`: either remove the unused dependency and stale exports, or move `web/src/ee/**` / `worker/src/ee/**` into it so the documented `web -> ee` edge reflects reality.
- Add mechanical boundary checks: `dependency-cruiser`/`madge --cyclic` in CI for the workspace DAG, plus `eslint-plugin-import` `no-restricted-paths` (or Next's `server-only` imports) to make the frontend-safe vs server-only barrel split machine-enforced rather than guide-enforced.
- Split `@langfuse/shared` along its natural seams (domain/types vs db vs queues vs llm-execution vs email) if independent consumption ever matters; today's `exports` map is a good scaffold for that refactor.
- Extend the existing `no-restricted-imports` pattern (already used for icons in `web/eslint.config.mjs:36-51`) to forbid `@langfuse/shared/src/db` outside `src/server/**`, `src/pages/api/**`, and test directories in web.

## Questions / Gaps

- **No evidence found for dedicated separation tests**: searched for cycle detectors (`madge`, `dpdm`), `dependency-cruiser` configs, and import-scanning tests across `package.json`, `turbo.json`, `.github/workflows/*`, and `scripts/` — none exist. The closest artifacts are the eslint-plugin RuleTester suite (`packages/eslint-plugin/src/rules/no-in-source-vitest.test.ts:1-25`) which tests style rules, not boundaries.
- **Whether `@langfuse/ee` was ever meaningfully imported** cannot be answered from a snapshot alone (would need git history); current state shows zero usage.
- Intra-package circularity (module-level cycles *within* `@langfuse/shared`, e.g., through its barrels) was not exhaustively verified — barrels re-export widely (`packages/shared/src/server/index.ts:1-157`) and no tooling reports cycles; only the package-level DAG is proven acyclic.

---

Generated by `22.01-package-and-module-boundaries` against `langfuse`.
