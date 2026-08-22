# Source Analysis: temporal

## Dimension 22.01: Package and Module Boundaries

### Source Info

| Field | Value |
|-------|-------|
| Name | temporal (Temporal Server) |
| Path | `studies/agent-harness-study/sources/temporal` |
| Language / Stack | Go 1.26 single module `go.temporal.io/server` (`go.mod:1-3`), protobuf/gRPC, uber-fx DI |
| Analyzed | 2026-08-22 |

## Summary

Temporal is a single Go module whose boundaries live at the *package* level, not the module level. The top level is a documented layered layout — `api` (generated internal protos), `client` (inter-service RPC clients), `common` (shared infrastructure incl. persistence), `service` (frontend/history/matching/worker), `chasm` (component-state-machine library), `components` (nexus components), `temporal` (server assembly), `temporaltest`, `cmd`, `tools`, `schema`, `tests` (`AGENTS.md:26-47`). Composition is explicit: each service exposes an fx module (`service/frontend/fx.go:68`, `service/history/fx.go:58`, `service/matching/fx.go:31`, `service/worker/fx.go:47`) and the top-level `temporal` package wires them into nested fx apps via `TopLevelModule` (`temporal/fx.go:132-154`, per-service graphs at `temporal/fx.go:498-602`).

The dependency flow is mostly downward (`service → client/common/chasm/components → api`; verified by exhaustive import scan, non-test files), and Go's compiler guarantees no package-level import cycles. However, three real boundary inversions exist: (1) `common/persistence` imports `service/history/tasks` in 19 non-test files, so the persistence interface surface depends on a package under a service; (2) the `chasm` library imports `service/history/tasks`, `service/history/queues/*`, and `service/worker/scheduler` while `service/history` imports chasm back — a bidirectional coupling at module granularity that stays compiler-legal only because leaf packages differ; (3) frontend reaches directly into history/worker internals (`service/frontend/workflow_handler.go:84-88`). Cycle pressure is acknowledged in code comments with deliberate workarounds (`common/persistence/data_interfaces.go:24`, `service/history/api/worker_versioning_util.go:13`, `temporal/cluster_metadata_loader.go:15`).

Public API discipline is strong where it matters most: user-facing API types live in external modules `go.temporal.io/api` (`go.mod:67`) and `go.temporal.io/sdk` (`go.mod:69`), with a CI-enforced version policy (`cmd/tools/check-dependencies/main.go:1-7`, wired in `.github/workflows/check-release-dependencies.yml:28`); proto aliasing conventions are lint-enforced via `importas` (`.github/.golangci.yml:64-70`); `temporaltest` carries an explicit SemVer compatibility contract (`temporaltest/README.md:5-7`). There is no automated internal dependency-direction checker (depguard covers only one external package), and Go's `internal/` mechanism is used in only two places.

## Rating

**7 / 10** — Clear model with explicit interfaces and operational safeguards, but boundary inversions between shared layers and services prevent a higher score. The layered intent is real and documented; composition boundaries (fx modules, persistence `Factory`, plugin registration) are clean and testable. What holds it back from 8+: `common/persistence → service/history/tasks`, `chasm ↔ service/history` mutual reach-ins, frontend→history/worker direct imports, no lint rule enforcing internal package dependency direction, and near-zero use of Go's `internal/` visibility mechanism outside two directories.

## Evidence Collected

Every entry includes a file path with line numbers. Format: `path/to/file.go:NN`.

| Area | Evidence | File:Line |
|------|----------|-----------|
| Single Go module | Module path `go.temporal.io/server`, go 1.26.4; public API deps pinned to separate modules | `go.mod:1-3`, `go.mod:67-69` |
| Documented package structure | Project structure map listing api/chasm/client/cmd/common/components/config/docs/proto/schema/service | `AGENTS.md:26-47` |
| Pure leaf API layer | `api/` has zero imports of any internal package (exhaustive grep of all `.go` files) | `api/` (grep result; e.g. generated `api/persistence/v1/chasm.pb.go:1`) |
| Internal proto sources separated | `proto/internal` holds buf-managed internal protos, distinct from generated code in `api/` | `proto/internal/buf.yaml`, `Makefile:27` |
| Server facade hides assembly | `Server` interface exposes only Start/Stop; `NewServer` delegates to `NewServerFx(TopLevelModule)` | `temporal/server.go:16-20`, `temporal/server.go:44-46` |
| Top-level composition module | `TopLevelModule` provides server impl + per-service providers + shared modules (dynamicconfig, pprof, tracing, chasm, serialization) | `temporal/fx.go:132-154` |
| Per-service nested fx graphs | HistoryServiceProvider/MatchingServiceProvider/FrontendServiceProvider/WorkerServiceProvider each build an isolated `fx.New` app; skipped cleanly if service not requested | `temporal/fx.go:498-516`, `temporal/fx.go:518-534`, `temporal/fx.go:536-584`, `temporal/fx.go:586-602` |
| Explicit graph-propagation seam | `GetCommonServiceOptions` converts server-graph dependencies into service-graph options; comment acknowledges "This is a workaround" | `temporal/fx.go:386-393` |
| Service modules as public seams | Each service exports `var Module = fx.Options(...)` consumed only by `temporal` | `service/frontend/fx.go:68`, `service/history/fx.go:58`, `service/history/chasm_engine.go:71`, `service/history/queue_factory_base.go:77`, `service/matching/fx.go:31`, `service/worker/fx.go:47` |
| Inter-service client boundary | `Factory` interface vends typed clients for history/matching/frontend/admin; services never dial each other directly | `client/clientfactory.go:28-50` |
| Persistence interface boundary | Manager interfaces (ShardManager, ExecutionManager, TaskManager, MetadataManager, ...) hide datastore implementations | `common/persistence/data_interfaces.go:1106-1258` |
| Datastore factory abstraction | `Factory ... The actual datastore is implementation detail hidden behind this interface` | `common/persistence/client/factory.go:23-46` |
| Driver plugin registration | SQLite driver self-registers via `init()`; loaded through blank imports in the binary only | `common/persistence/sql/sqlplugin/sqlite/plugin.go:45`, `cmd/server/main.go:24-26` |
| Custom store extension points | `AbstractDataStoreFactory`, custom visibility store, custom archiver factories surfaced as server options | `temporal/fx.go:108-111`, `common/persistence/client/abstract_data_store_factory.go:16` |
| Upward edge: persistence→history tasks | 19 non-test files under `common/persistence` import `service/history/tasks`; interfaces use `tasks.Category`/`tasks.Task` directly | `common/persistence/data_interfaces.go:20`, `common/persistence/data_interfaces.go:201,368-432`; full list incl. `common/persistence/cassandra/util.go`, `common/persistence/sql/execution_tasks.go` |
| Upward edge: dynamicconfig→matching | Dynamic-config constants import matching counter package | `common/dynamicconfig/constants.go:13` |
| Upward edge: rpc interceptor→frontend | Telemetry interceptor imports `service/frontend/configs` | `common/rpc/interceptor/telemetry.go:20` |
| Bidirectional chasm↔history | Chasm core/lib import history task & queue packages while 91 non-test files under `service/history` import `chasm` | `chasm/tree.go:32`, `chasm/transition_history.go:5`, `chasm/lib/callback/component.go:13`, `chasm/lib/scheduler/scheduler.go:27`; reverse direction: grep of `service/history` |
| Frontend reaching into sibling services | Frontend handler imports `service/history/api`, `service/worker/{batcher,dummy,scheduler,workerdeployment}` and calls history helpers directly | `service/frontend/workflow_handler.go:84-88`, `service/frontend/workflow_handler.go:581,992`; also `service/frontend/fx.go:50-52` |
| Documented cycle avoidance (alias) | "Archetype is a type alias for chasm.Archetype to avoid circular dependency." | `common/persistence/data_interfaces.go:24` |
| Documented cycle avoidance (function type) | Function-type indirection "avoiding import cycles between history/api and worker/workerdeployment packages" | `service/history/api/worker_versioning_util.go:11-17` |
| Documented cycle avoidance (placement TODO) | Cluster metadata loader parked in `temporal` "to avoid a circular dependency" | `temporal/cluster_metadata_loader.go:15` |
| Lint-enforced alias conventions | `importas`: public pb `${1}pb`, internal spb `${1}spb`, no aliases for pb services; enforced repo-wide | `.github/.golangci.yml:64-70` |
| Lint-enforced scoped behavior rules | forbidigo bans `time.Now` in `chasm/lib` (use `ctx.Now(component)`), Unix* time funcs only restricted inside cassandra persistence, panic banned in app code | `.github/.golangci.yml:33-56,178-186` |
| External dependency policy tool | `check-dependencies` validates api/sdk versions per branch; runs in CI | `cmd/tools/check-dependencies/main.go:1-7`, `.github/workflows/check-release-dependencies.yml:28` |
| Public test harness w/ compat contract | `temporaltest` promises SemVer BC for its Go API (one named exception) | `temporaltest/README.md:5-7`, `temporaltest/server.go:11-30` |
| Hidden implementation detail | Lite server implementation tucked in `internal/` package | `temporaltest/internal/lite_server.go:1-3` |
| Test-only surface idiom | `export_test.go` pattern used to widen chasm API for tests only | `chasm/export_test.go` |
| Test-support isolation | Integration tests quarantined in `tests/` behind build tag policy ("integration" tag per `AGENTS.md:71`); test fixtures like `persistencetest/` live beside their layer | `tests/` dir listing; `common/persistence/persistencetest/queues.go` |
| Chasm component registry | Libraries register into a `Registry` (plugin-style extension within the library) | `chasm/registry.go:69-105`, `chasm/fx.go:5-11` |
| Thin binary entrypoint | `cmd/server/main.go` is CLI scaffolding + blank plugin imports; all wiring delegated to `temporal` package | `cmd/server/main.go:10-40` |

## Answers to Dimension Questions

**1. Are modules cleanly separated?**
Largely yes at the top level, with three known leaks. The documented layout assigns clear responsibilities to each top-level directory (`AGENTS.md:26-47`). Services are isolated from one another through the `client.Factory` indirection (`client/clientfactory.go:28-36`) rather than direct calls, and each service is assembled as its own nested fx app (`temporal/fx.go:498-602`). The leaks: `common/persistence` depends on `service/history/tasks` across its interface definitions (`common/persistence/data_interfaces.go:20,201,368-432`; 19 importing files including `common/persistence/cassandra/util.go` and `common/persistence/sql/execution_tasks.go`); the chasm library both imports and is imported by history-side packages (`chasm/tree.go:32` vs 91 importing files in `service/history`); and frontend directly uses history/worker internals (`service/frontend/workflow_handler.go:84-88,581,992`). Separation between the *binary*, the *assembly*, and the *services* is clean; separation between the *shared kernel* and *services* is not strict.

**2. Do dependencies flow in one direction?**
Downward in the main path — verified by an exhaustive non-test import matrix: `api` imports nothing; `schema → {api, common}`; `client → {api, chasm, common}`; `components → {api, chasm, client, common}`; `service → {api, chasm, client, common, components}`; `temporal → {api, chasm, client, common, service}`; `cmd → {api, chasm, client, common, temporal, tools}`. But upward edges exist exactly at the shared-kernel boundary: `common → service` (persistence/dynamicconfig/rpc-interceptor files listed above), `chasm → service` (tasks, queues/errors, worker/scheduler). Two upward edges are confined to test-support code (`common/persistence/cassandra/test.go`, `common/persistence/sql/test_sql_persistence.go` importing `temporal`). Package-level cycles are impossible in Go and none exist; module-granularity cycles were actively dodged instead, as documented at `common/persistence/data_interfaces.go:24` and `service/history/api/worker_versioning_util.go:11-17`.

**3. Can modules be used independently?**
Partially. Fully standalone: `api` (zero internal imports) and nearly standalone `common/serviceerror` (imports only `api/*` and `common/log`). The `client` package is runtime-independent of services in non-generated code — the only service import is the test helper `client/history/historytest/clienttest.go`. Persistence consumers can swap datastores without touching services because everything goes through `Factory` (`common/persistence/client/factory.go:23-46`) and drivers self-register via `init()` (`common/persistence/sql/sqlplugin/sqlite/plugin.go:45`), with binaries opting in via blank imports (`cmd/server/main.go:24-26`). However, answering the dimension's headline question ("can you use the tool system without pulling in the entire runtime?"): you cannot consume `common/persistence` or `chasm` without dragging in `service/history/tasks` (and for chasm lib components, `service/history/queues/errors`, `service/history/queues/common`, `service/worker/scheduler` — see `chasm/lib/callback/component.go:13`, `chasm/lib/scheduler/scheduler.go:27`). The shared layers are not independently shippable. `temporaltest` *is* independently consumable and explicitly supported as such (`temporaltest/README.md:1-7`).

**4. Are public APIs distinguished from internal ones?**
Yes, with a mix of mechanisms. (a) The true public API lives outside this repo in `go.temporal.io/api` and `go.temporal.io/sdk` (`go.mod:67,69`), and a CI gate enforces that release/cloud branches pin tagged semver versions of those modules (`cmd/tools/check-dependencies/main.go:1-8`, `.github/workflows/check-release-dependencies.yml:28`). (b) Naming/lint convention separates internal protobuf types (`${1}spb`) from public ones (`${1}pb`) and forbids unaliased imports (`.github/.golangci.yml:64-70`). (c) Go's `internal/` visibility is used sparingly — only `proto/internal` and `temporaltest/internal` exist — elsewhere "internal" is convention, not enforcement. (d) The `export_test.go` idiom keeps widened surfaces out of the shipped API (`chasm/export_test.go`). (e) `temporaltest` documents exactly which parts of its API are stable (`temporaltest/README.md:5-7`). What's missing: no depguard/import-boundary rules constrain imports *between* internal packages, so nothing mechanical stops a new file in `common` from importing `service/...` today.

## Architectural Decisions

- **Single module, package-level boundaries.** All code ships as one module (`go.mod:1`); the hard public-API split is pushed to external modules (`go.mod:67-69`). This simplifies refactoring internally but means boundary violations are caught only by review/lint, never by the linker.
- **Composition root with nested DI graphs.** The server builds one parent fx graph and one child graph per enabled service (`temporal/fx.go:498-602`), which enforces that services receive their dependencies through an explicit propagation object (`temporal/fx.go:386-393`) rather than reaching into siblings. The authors flag this propagation as a workaround, showing awareness of the tension.
- **Interfaces + init()-registration for pluggable infrastructure.** Persistence managers are interfaces (`common/persistence/data_interfaces.go:1106-1258`), datastores vendored via `Factory` (`common/persistence/client/factory.go:23-46`), SQL drivers registered in `init()` (`common/persistence/sql/sqlplugin/sqlite/plugin.go:45`) and activated by blank imports at the binary edge (`cmd/server/main.go:24-26`). Same shape for chasm libraries registering into a registry (`chasm/registry.go:69`).
- **Cycle avoidance by aliasing and function-typing rather than restructuring.** Three separate sites document deliberate workarounds (`common/persistence/data_interfaces.go:24`, `service/history/api/worker_versioning_util.go:11-17`, `temporal/cluster_metadata_loader.go:15`) — evidence that the team treats the compiler as the cycle guardrail of last resort and accepts local debt instead of re-layering.
- **Lint as boundary police for behavior, not structure.** Scoped forbidigo rules encode package-specific contracts (no `time.Now` in `chasm/lib`, Cassandra timestamp rules confined to the cassandra store, alias rules for pb/spb — `.github/.golangci.yml:33-70,178-186`), but no equivalent machine rule guards import *direction*.

## Notable Patterns

- **Facade + option pattern at every seam:** `Server` interface hiding `ServerFx` (`temporal/server.go:16-46`), `ServerOption`s feeding a provider struct that defaults every dependency (`temporal/fx.go:94-129,170-318`).
- **Service-name-keyed skip logic:** providers silently no-op when a service isn't requested, letting one binary host any subset of services (`temporal/fx.go:502-506`), consistent with the `Services`/`DefaultServices` lists (`temporal/server.go:23-41`).
- **Shared-kernel drift:** `service/history/tasks` functions as a de-facto shared kernel (persisted category IDs documented as DB-stable at `service/history/tasks/category.go:19-28`) that `common/persistence` and `chasm` both depend on — a gravity well pulling service internals into the shared layer.
- **Test-support subpackages co-located with their layer:** `chasmtest`, `persistencetest`, `client/history/historytest`, `common/testing/testhooks` — keeping test doubles adjacent to the boundary they mock.
- **Documentation tied to structure:** architecture docs mirror the package split (`docs/architecture/README.md:3-9`), and `AGENTS.md:26-47` codifies the directory contract for contributors and agents alike.

## Tradeoffs

- **Monorepo convenience vs. enforceability:** one module maximizes refactor speed but forfeits compile-time boundary enforcement; the project compensates with lint scoping and CI gates on the external API only.
- **Nested fx graphs vs. one flat graph:** per-service graphs give isolation and selective startup but forced the `GetCommonServiceOptions` bridge, which the authors themselves call "not ideal" (`temporal/fx.go:386-393`).
- **Chasm-as-library vs. chasm-in-history:** placing the chasm library at top level suggests independence, but its core (`chasm/tree.go:32`) and bundled libraries (`chasm/lib/callback/component.go:13`, `chasm/lib/scheduler/scheduler.go:27`) depend on history/worker internals — the library cannot be extracted without first re-homing task/queue error types.
- **init()-based plugin registration:** simple opt-in drivers, but registration order/omission errors appear only at startup rather than build time (mitigated by functional tests).

## Failure Modes / Edge Cases

- **Silent boundary erosion risk:** nothing mechanically prevents new upward imports from `common` or `chasm` into `service/*`; the existing 19-file inversion shows erosion already happened incrementally. A future contributor adding one more `service/history` import to `common/persistence` breaks no build and trips no linter.
- **Near-cycle fragility:** the documented workarounds (`data_interfaces.go:24`, `worker_versioning_util.go:11-17`) mean certain refactors (e.g., making `worker/workerdeployment` import `history/api` types directly) will flip a latent cycle into a compile error, surfacing as confusing cross-team breakage.
- **Frontend↔history coupling under change:** because frontend calls `history/api.ProcessInternalRawHistory` directly (`service/frontend/workflow_handler.go:992`), changes to history request/response shaping can break the frontend build even though the services are meant to communicate over gRPC.
- **Blank-import omission:** building a custom binary that forgets the sqlite/mysql/postgresql blank imports (`cmd/server/main.go:24-26`) yields a runtime factory error rather than a build failure.

## Future Considerations

- Enforce import direction mechanically: extend the existing depguard setup (`.github/.golangci.yml:57-63`) with rules denying `service/*` imports from `common/**` and `chasm/**` (with a curated allowlist for `service/history/tasks` until re-homed), turning today's convention into CI-checked invariant.
- Re-home the de-facto shared kernel: move `service/history/tasks` (or its persisted-ID core, `service/history/tasks/category.go:19-28`) and `service/history/queues/errors` into `common` or a neutral `tasks` package, dissolving the largest upward edge and the chasm↔history knot.
- Make wider use of Go `internal/` (e.g., `service/*/internal`) now that frontend→history direct calls (`workflow_handler.go:581,992`) could be routed through shared contracts instead.
- Resolve the parked items flagged in-code: relocate `temporal/cluster_metadata_loader.go` per its own TODO (`temporal/cluster_metadata_loader.go:15`) and revisit the fx propagation workaround (`temporal/fx.go:386-393`).

## Questions / Gaps

- No automated package-dependency-direction checker was found. Searched: `.github/.golangci.yml` (depguard limited to external uuid package), `Makefile` lint targets (`Makefile:390-412`), `cmd/tools/check-dependencies` (external-module policy only). If an internal boundary tool exists, it lives outside this source tree.
- Whether the chasm↔history entanglement is a deliberate end-state or migration debt is not stated anywhere in-repo; `docs/architecture/chasm.md` describes the model but no evidence was found addressing the library's dependency posture toward `service/history`.
- The `common → tests` import edge observed in the matrix traces only to test-fixture files (`common/persistence/cassandra/test.go`, `common/persistence/sql/test_sql_persistence.go`); whether these fixtures are compiled into release artifacts was not verified (no build-tag evidence found in the files inspected).
- Package-level acyclicity is asserted from Go compiler semantics plus the presence of cycle-avoidance comments; a full `go build ./...` was not executed in this sandboxed analysis environment.

---

Generated by dimension 22.01 (Package and Module Boundaries) against `temporal`.
