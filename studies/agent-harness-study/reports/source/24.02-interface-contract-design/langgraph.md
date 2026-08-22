# Source Analysis: langgraph

## Interface Contract Design

### Source Info

| Field | Value |
|-------|-------|
| Name | langgraph |
| Path | `studies/agent-harness-study/sources/langgraph` |
| Language / Stack | Python (monorepo: core `libs/langgraph`, contract packages `libs/checkpoint`, conformance suite `libs/checkpoint-conformance`, implementations `libs/checkpoint-postgres`/`libs/checkpoint-sqlite`, tooling `libs/prebuilt`) |
| Analyzed | 2026-08-22 |

## Summary

LangGraph manages contracts through a layered architecture with an unusually strong enforcement story for its persistence layer. The base package (`langgraph-checkpoint`) owns the storage contracts — `BaseCheckpointSaver` (`libs/checkpoint/langgraph/checkpoint/base/__init__.py:176`), `SerializerProtocol` (`libs/checkpoint/langgraph/checkpoint/serde/base.py:15`), and `BaseStore` (`libs/checkpoint/langgraph/store/base/__init__.py:700`) — while the core framework (`libs/langgraph`) defines engine-side contracts: `BaseChannel` (`libs/langgraph/langgraph/channels/base.py:19`), `PregelProtocol` (`libs/langgraph/langgraph/pregel/protocol.py:25`), and the run-scoped `Runtime` bundle (`libs/langgraph/langgraph/runtime.py:125`). Dependency direction is consumer-owned: the checkpoint library re-declares a structural `ChannelProtocol` that "Mirrors langgraph.channels.base.BaseChannel" rather than importing it (`libs/checkpoint/langgraph/checkpoint/serde/types.py:39-40`), so the persistence layer has no dependency on the engine.

The standout mechanism is a dedicated conformance package, `libs/checkpoint-conformance`, that turns documented behavioral clauses of the checkpointer contract into executable capability suites. Nine capabilities (`put`, `put_writes`, `get_tuple`, `list`, `delete_thread`, `delete_for_runs`, `copy_thread`, `prune`, `delta_channel_history`) are auto-detected per implementation via method-override inspection (`libs/checkpoint-conformance/langgraph/checkpoint/conformance/capabilities.py:53-63,90-96`) and validated against spec tests that assert semantics, not just callability: write idempotency (`libs/checkpoint-conformance/langgraph/checkpoint/conformance/spec/test_put_writes.py:148-167`), namespace isolation (`test_put_writes.py:207-239`), special-channel handling for ERROR/INTERRUPT writes (`test_put_writes.py:170-204`), and pending-write clearing on subsequent checkpoints (`test_put_writes.py:242-263`).

Contracts evolve through explicit machinery: a version field on the checkpoint format (`Checkpoint.v`, `libs/checkpoint/langgraph/checkpoint/base/__init__.py:96`), open-ended metadata TypedDicts ("Marked as total=False to allow for future expansion", line 37-38), docstring-marked beta surfaces with failure-mode warnings (lines 66-69, 152-155), a versioned streaming API with v1/v2 overloads (`libs/langgraph/langgraph/pregel/protocol.py:107-149`), adapter shims for legacy serializers (`maybe_add_typed_methods`, `libs/checkpoint/langgraph/checkpoint/serde/base.py:40-48`), and runtime signature-sniffing for parameter drift on `put_writes(task_path=...)` (`libs/langgraph/langgraph/pregel/_loop.py:1497-1500`). The store and cache contracts enjoy runtime validation (namespace rules, TTL capability gates) but lack a conformance harness, making checkpointer substitutability markedly better engineered than store substitutability.

## Rating

**8 / 10** — Clear model with tests, explicit interfaces, and operational safeguards.

Rationale: The checkpointer contract is mature by every measure in the rubric: small consumer-owned abstractions, behavioral semantics encoded both in prose and in a runnable conformance suite, multiple independent production implementations (memory, sqlite sync/async, postgres sync/async/shallow) that all honor shared protocols such as `WRITES_IDX_MAP` (`libs/checkpoint/langgraph/checkpoint/base/__init__.py:795`, consumed at `libs/checkpoint-postgres/langgraph/checkpoint/postgres/base.py:536` and `libs/checkpoint-sqlite/langgraph/checkpoint/sqlite/__init__.py:476`), graceful degradation paths, and explicit evolution/versioning mechanisms. It falls short of 9–10 because: (1) the `BaseStore` and `BaseCache` contracts have no equivalent conformance harness — substitutability there rests on convention plus per-implementation tests; (2) declared schemas have drifted from usage in spots (e.g., `pending_sends` is written and read but absent from the `Checkpoint` TypedDict, `libs/checkpoint/langgraph/checkpoint/base/__init__.py:134,824` vs. definition at :92-123); (3) some compatibility is maintained via fragile runtime introspection (`_loop.py:1497-1500`).

## Evidence Collected

Every entry includes a file path with line numbers. All paths are relative to `sources/langgraph/`.

| Area | Evidence | File:Line |
|------|----------|-----------|
| Checkpointer contract | `BaseCheckpointSaver` defines sync/async method pairs; unimplemented methods raise `NotImplementedError` as explicit extension points with semantic docstrings | libs/checkpoint/langgraph/checkpoint/base/__init__.py:176-580 |
| Semantic guarantee: ID ordering | Checkpoint `id` documented as "both unique and monotonically increasing" | libs/checkpoint/langgraph/checkpoint/base/__init__.py:98-101 |
| Version monotonicity contract | `get_next_version` requires "monotonically increasing" versions; default int increment | libs/checkpoint/langgraph/checkpoint/base/__init__.py:692-711 |
| Cross-implementation write protocol | `WRITES_IDX_MAP = {ERROR: -1, SCHEDULED: -2, INTERRUPT: -3, RESUME: -4}` — "Each Checkpointer implementation should use this mapping in put_writes"; honored in memory, postgres (sync/aio/shallow), sqlite (sync/aio) | libs/checkpoint/langgraph/checkpoint/base/__init__.py:788-795; libs/checkpoint/langgraph/checkpoint/memory/__init__.py:500; libs/checkpoint-postgres/langgraph/checkpoint/postgres/base.py:536; libs/checkpoint-sqlite/langgraph/checkpoint/sqlite/__init__.py:464-476 |
| Serializer protocol | `@runtime_checkable SerializerProtocol` with `dumps_typed`/`loads_typed`; legacy adapters via `SerializerCompat` + `maybe_add_typed_methods` | libs/checkpoint/langgraph/checkpoint/serde/base.py:14-48 |
| Encryption decorator | `EncryptedSerializer` composes over any `SerializerProtocol` + `CipherProtocol`; tolerates legacy unencrypted blobs ("+" suffix convention) | libs/checkpoint/langgraph/checkpoint/serde/encrypted.py:8-36 |
| Store contract minimality | `BaseStore` ABC: only `batch`/`abatch` abstract; `get/search/put/delete/list_namespaces` are template-method wrappers over batch ops | libs/checkpoint/langgraph/store/base/__init__.py:724-746,767-769,917-927 |
| Runtime namespace validation | `_validate_namespace`: rejects empty/non-string/dotted labels and reserved root `"langgraph"` with typed error `InvalidNamespaceError` | libs/checkpoint/langgraph/store/base/__init__.py:541-542,1255-1275 |
| TTL capability gate | `supports_ttl: bool = False` class flag; `put(ttl=...)` raises `NotImplementedError` naming the class when unsupported | libs/checkpoint/langgraph/store/base/__init__.py:719,912-916 |
| TTL semantics as prose | TTL expiry "scheduled for deletion on a best-effort basis"; refresh-on-read/write timing specified | libs/checkpoint/langgraph/store/base/__init__.py:526-534 |
| Filter operator grammar | `SearchOp.filter` documents `$eq/$ne/$gt/$gte/$lt/$lte` operators with examples — a query-language contract | libs/checkpoint/langgraph/store/base/__init__.py:250-285 |
| Channel ABC | `BaseChannel(Generic[Value, Update, Checkpoint])` with lifecycle hooks `consume()`/`finish()`, default no-op implementations, and error contracts (`EmptyChannelError`, `InvalidUpdateError`) | libs/langgraph/langgraph/channels/base.py:19-121 |
| Dependency inversion | `ChannelProtocol` in checkpoint lib "Mirrors langgraph.channels.base.BaseChannel" without importing it | libs/checkpoint/langgraph/checkpoint/serde/types.py:39-55 |
| Cache contract | `BaseCache(ABC)`: six abstract methods (get/set/clear × sync/async), serializer injected via constructor | libs/checkpoint/langgraph/cache/base/__init__.py:15-48 |
| Engine-level protocol | `PregelProtocol(Runnable)` — state access, bulk update, streaming/invoke with `version: Literal["v1","v2"]` overloads showing deliberate API versioning | libs/langgraph/langgraph/pregel/protocol.py:25-269 |
| Context propagation contract | `Runtime` dataclass bundles context/store/stream_writer/heartbeat/execution_info/control; null-object defaults (`_no_op_stream_writer`, `DEFAULT_RUNTIME`) keep consumers safe outside graphs | libs/langgraph/langgraph/runtime.py:107-121,206-217,285-293 |
| Cooperative cancellation | `RunControl.request_drain()` + `GraphDrained` exception: drain at superstep boundary, checkpoint saved, resumable later | libs/langgraph/langgraph/runtime.py:79-104; libs/langgraph/langgraph/errors.py:54-64 |
| Cancellation semantics distinction | `NodeCancelledError` converts user-raised `asyncio.CancelledError` into a node failure so runs report `error` instead of silently succeeding (framework-initiated cancellation stays silent) | libs/langgraph/langgraph/errors.py:168-186 |
| Timeout/retry interaction | `NodeTimeoutError` deliberately does NOT inherit from built-in `TimeoutError` "so that the default RetryPolicy treats it as retryable"; carries `kind: idle\|run` discriminant | libs/langgraph/langgraph/errors.py:190-206 |
| Error code taxonomy | `ErrorCode` enum with per-code troubleshooting URLs embedded in messages | libs/langgraph/langgraph/errors.py:34-47 |
| Conformance suite runner | `validate()` detects capabilities per implementation, creates fresh saver per capability, returns `CapabilityReport` | libs/checkpoint-conformance/langgraph/checkpoint/conformance/validate.py:45-128 |
| Capability detection | `_is_overridden` compares impl method identity against `BaseCheckpointSaver` defaults; BASE vs EXTENDED capability split | libs/checkpoint-conformance/langgraph/checkpoint/conformance/capabilities.py:29-50,73-96 |
| Behavioral spec: idempotency | `test_put_writes_idempotent` asserts duplicate `(task_id, idx)` writes do not duplicate rows | libs/checkpoint-conformance/langgraph/checkpoint/conformance/spec/test_put_writes.py:148-167 |
| Behavioral spec: isolation | `test_put_writes_across_namespaces` asserts `checkpoint_ns` scoping of pending writes | libs/checkpoint-conformance/langgraph/checkpoint/conformance/spec/test_put_writes.py:207-239 |
| Behavioral spec: recency | `test_get_tuple_latest_when_no_checkpoint_id` pins "returns newest checkpoint" semantics | libs/checkpoint-conformance/langgraph/checkpoint/conformance/spec/test_get_tuple.py:30 |
| Registration/decoration API | `@checkpointer_test(name=..., skip_capabilities=..., lifespan=...)` factory registry with async-generator lifecycles for DB setup/teardown | libs/checkpoint-conformance/langgraph/checkpoint/conformance/initializer.py:59-100 |
| Conformance adopted downstream | AsyncSqliteSaver runs delta-channel conformance in its own test suite; InMemorySaver likewise | libs/checkpoint-sqlite/tests/test_conformance_delta.py:15-35; libs/checkpoint/tests/test_conformance_delta.py:16-20; libs/checkpoint-conformance/tests/test_validate_memory.py:8-19 |
| Graph-time validation | `validate_graph` rejects reserved names, unknown channels, unsubscribed input channels, dangling output/interrupt targets — fail-fast at compile time | libs/langgraph/langgraph/pregel/_validate.py:13-107 |
| Tool-call middleware contract | `ToolCallRequest` dataclass + immutable `.override()` (direct mutation deprecated with warning); `ToolCallWrapper` protocol documents multi-invoke `execute` for retry | libs/prebuilt/langgraph/prebuilt/tool_node.py:133-199,202-277 |
| Deferred tool validation | Validation deferred until `execute()` so interceptors can short-circuit unregistered tools (`tool: BaseTool \| None`) | libs/prebuilt/langgraph/prebuilt/tool_node.py:1030-1032,1091-1097 |
| Schema-time tool validation | Pydantic validation of tool args with injected-argument errors filtered out before surfacing `ToolInvocationError` | libs/prebuilt/langgraph/prebuilt/tool_node.py:1103-1113 |
| Interrupt bubble pass-through | `GraphBubbleUp` explicitly re-raised past tool error handling, with scenarios enumerated in comments | libs/prebuilt/langgraph/prebuilt/tool_node.py:1120-1130 |
| Deserialization security allowlist | msgpack deserialization restricted to allowlisted modules; `with_msgpack_allowlist` merges and clones; violations logged with remediation note | libs/checkpoint/langgraph/checkpoint/serde/jsonplus.py:117-119,128-153,174,240 |
| Graceful capability degradation | `with_allowlist` logs a warning and returns original serde if the serializer lacks allowlist support instead of failing | libs/checkpoint/langgraph/checkpoint/base/__init__.py:713-742 |
| Format versioning | `Checkpoint.v` ("Currently 1") and deprecated-utils section below `LATEST_VERSION = 2` fence | libs/checkpoint/langgraph/checkpoint/base/__init__.py:96,809-811 |
| Beta surface demarcation | DeltaChannel methods/types carry `!!! warning "Beta"` blocks describing instability and override risk | libs/checkpoint/langgraph/checkpoint/base/__init__.py:66-69,152-155,587-593 |
| Failure-mode warnings in contract | `prune`/`copy_thread`/`delete_for_runs` docstrings spell out how naive implementations silently corrupt DeltaChannel state and list safe strategies | libs/checkpoint/langgraph/checkpoint/base/__init__.py:340-347,361-371,387-414 |
| Default-walk template method | `get_delta_channel_history` provides correct default walking public APIs only ("the return contract is fixed here"); savers may optimize | libs/checkpoint/langgraph/checkpoint/base/__init__.py:607-649 |
| Feature detection for drift | Loop sniffs `signature(checkpointer.put_writes).parameters.get("task_path")` to support older savers | libs/langgraph/langgraph/pregel/_loop.py:1497-1500,469-483 |
| Streaming contract versioning | `StreamPart[StateT, OutputT]` v2 stream parts vs v1 dict payloads, selected by `version` literal | libs/langgraph/langgraph/pregel/protocol.py:107-149; libs/langgraph/langgraph/types.py (StreamMode/StreamPart referenced) |

## Answers to Dimension Questions

### 1. Are interfaces small, coherent, and owned by the consumer side?

Mostly yes. The strongest examples: `BaseStore` requires implementations to provide exactly two methods (`batch`/`abatch`, `libs/checkpoint/langgraph/store/base/__init__.py:724-746`) with everything else derived; `BaseCache` has six abstract members mirroring get/set/clear (`libs/checkpoint/langgraph/cache/base/__init__.py:24-48`); `ManagedValue` is a single-method ABC (`libs/langgraph/langgraph/managed/base.py:18-21`). Consumer ownership shows up structurally: the checkpoint package declares its own `ChannelProtocol` mirroring the engine's `BaseChannel` (`libs/checkpoint/langgraph/serde/types.py:39-40`) so the storage layer never imports the engine, and `Runtime` documents that `config` injection remains the caller's choice rather than being baked in (`libs/langgraph/langgraph/runtime.py:131-135`). The one interface that strains the principle is `BaseCheckpointSaver` itself: 20+ public methods once sync/async pairs and the beta DeltaChannel surface are counted (`libs/checkpoint/langgraph/checkpoint/base/__init__.py:176-690`), though sensible defaults keep the mandatory core small (five BASE capabilities, `capabilities.py:29-38`).

### 2. Do contracts specify behavior, not just method signatures?

Yes, unusually well. Behavior clauses appear at three levels:

- **Prose**: monotonic checkpoint IDs (`base/__init__.py:100-101`), monotonic channel versions (`:697`), best-effort TTL expiry (`store/base/__init__.py:532`), retry-policy interaction of `NodeTimeoutError` (`errors.py:193-195`), silent-corruption warnings for DeltaChannel-unaware pruning (`base/__init__.py:397-413`).
- **Executable specs**: idempotency, recency, namespace isolation, parent-chain integrity, and ERROR/INTERRUPT write indexing are asserted in `libs/checkpoint-conformance/.../spec/*.py`.
- **Cross-cutting protocols**: `WRITES_IDX_MAP` coordinates negative-index conventions across all seven saver implementations (`base/__init__.py:795`).

Where behavior cannot yet be pinned (beta surfaces), the contract says so explicitly rather than implying stability (`!!! warning "Beta"`, `base/__init__.py:66-69`).

### 3. Can providers, tools, stores, and runtimes be replaced safely?

**Checkpointers: yes, demonstrably.** Three independent backends (memory/sqlite/postgres) satisfy the same contract; the conformance suite verifies each (`libs/checkpoint-sqlite/tests/test_conformance_delta.py:26-31`); adapters bridge legacy serializers (`serde/base.py:29-48`); feature detection absorbs parameter drift (`_loop.py:1497-1500`); and unsupported optional features degrade gracefully with warnings rather than crashes (`base/__init__.py:737-742`). Two caveats: detection-by-signature would misclassify a saver accepting `task_path` via `**kwargs` (`_loop.py:1498-1499`), and Postgres does not yet run the conformance suites (no references found under `libs/checkpoint-postgres/tests/` — searched for `conformance` and `checkpointer_test`).

**Stores: partially.** Runtime gates exist — namespace validation (`store/base/__init__.py:1255-1275`), TTL capability flags raising typed errors (`:912-916`) — but there is no cross-implementation conformance harness. Postgres, SQLite, and InMemory store tests are authored independently per backend (e.g., `libs/checkpoint-postgres/tests/test_store.py`, `libs/checkpoint/tests/test_store.py`), so subtle divergences in filter-operator or TTL semantics between backends would not be caught mechanically.

**Tools: yes within the LangChain `BaseTool` ecosystem.** Schemas are validated at registration time (schema must be a Pydantic model, `libs/prebuilt/langgraph/prebuilt/tool_validator.py:143-166`) and args at invocation time with injected-parameter error filtering (`tool_node.py:1106-1113`); unregistered tools produce structured error messages rather than exceptions (`INVALID_TOOL_NAME_ERROR_TEMPLATE`, `tool_node.py:108-110`).

### 4. Are compatibility failures caught early by tests or validation?

Largely, for the checkpointer axis: graph wiring is validated eagerly at build time (`pregel/_validate.py:23-107` raises before any run), conformance suites execute in downstream CI, and schema-time Pydantic checks catch malformed tool inputs. However, several artifacts suggest failures historically surfaced late and were patched with runtime tolerance: the signature-sniffing shim (`_loop.py:1497`), the `EncryptedSerializer`'s tolerance of unencrypted legacy blobs (`serde/encrypted.py:29-30`), and the type-drift around `pending_sends` (see Failure Modes). The store/cache axes remain the weakest link — nothing mechanical prevents two stores from interpreting `$gt` filters or TTL refresh differently.

## Architectural Decisions

1. **Contract ownership inverted away from the engine.** Storage contracts live in a standalone `libs/checkpoint` package; the engine consumes them. Where types must be shared without coupling, structural mirrors replace imports (`ChannelProtocol`, `libs/checkpoint/langgraph/checkpoint/serde/types.py:39-55`).

2. **A dedicated conformance product, not just internal tests.** `libs/checkpoint-conformance` ships as a separate installable package with registration decorators, capability auto-detection, progress callbacks, and machine-readable reports (`initializer.py:59-100`, `capabilities.py:66-96`, `report.py:93-182`). Third-party checkpointer authors can validate against the same spec the maintainers use.

3. **Capability tiers over fat interfaces.** Five mandatory capabilities versus four optional ones (`capabilities.py:29-50`), detected via override inspection rather than declared flags or interface splitting — implementations advertise what they implement by implementing it.

4. **Template-method defaults as compatibility anchors.** Correct-but-slow default implementations (e.g., the ancestor walk in `get_delta_channel_history`, `base/__init__.py:620-649`) guarantee old savers keep working while fast paths are added; the comment "the return contract is fixed here" (`:610-611`) states intent directly.

5. **Versioned escape hatches over breaking changes.** Streaming got a `version="v2"` parameter with typed `StreamPart` outputs while v1 remains default (`protocol.py:119,134-135,148`); checkpoint format carries `v`; metadata dicts are `total=False` for forward expansion (`base/__init__.py:37-38,96`).

6. **Errors designed for their consumers.** Error classes encode recovery-relevant distinctions: drain-vs-cancel-vs-fail (`GraphDrained`, `NodeCancelledError`), retryable-vs-not via inheritance choices (`NodeTimeoutError`, `errors.py:190-199`), and stable string codes linking to troubleshooting docs (`ErrorCode`, `errors.py:34-47`).

## Notable Patterns

- **Null-object defaults**: `Runtime` fields default to no-op writers/heartbeats and a `DEFAULT_RUNTIME` singleton, so nodes work identically inside and outside graphs (`runtime.py:107-110,285-293`).
- **Decorator composition on Protocols**: `EncryptedSerializer(SerializerProtocol)` wraps any inner serializer and any cipher, with a "+"-suffixed type-tag scheme preserving backward reads (`serde/encrypted.py:17-36`).
- **Immutable request objects with deprecation rails**: `ToolCallRequest.__setattr__` warns and directs to `.override()` (`tool_node.py:151-168`) — evolving a mutable-era API toward immutability without breaking users.
- **Reserved namespaces enforced centrally**: reserved channel/node names (`pregel/_validate.py:23-33`) and the reserved `"langgraph"` store root (`store/base/__init__.py:1272-1275`) prevent collisions between user data and framework keys.
- **Security as a serializer concern**: deserialization allowlists with exact-symbol policy and explicit refusal of prefix wildcards (`jsonplus.py:240`).
- **Docstring-carried failure-mode analysis**: prune/copy/delete docstrings read like mini design docs, enumerating how naive implementations corrupt state and offering concrete safe algorithms (`base/__init__.py:387-414`).

## Tradeoffs

- **Override-detection vs explicit capability declaration**: `_is_overridden` (`capabilities.py:90-96`) keeps interfaces lean but couples capability reporting to implementation details — an implementation inheriting a complete default implementation of an EXTENDED capability would be reported as non-capable, and one overriding merely to add logging would be reported capable even if broken.
- **Fat-but-defaulted BaseCheckpointSaver**: convenience accrues to the base class (20+ methods), trading interface minimalism for discoverability and default correctness.
- **Signature sniffing vs version negotiation**: runtime `inspect.signature` probes (`_loop.py:1497-1500`) avoid version-pinning dependencies but are brittle against `**kwargs`-style implementations and invisible to static checking.
- **Prose-heavy semantics**: much behavioral specification lives in docstrings; where conformance tests cover it this is fine, but store/TTL/filter semantics currently rely on documentation alone.
- **Per-implementation duplication in store tests**: without a store conformance suite, postgres/sqlite/memory test suites drift independently (compare `libs/checkpoint-postgres/tests/test_store.py` vs `libs/checkpoint/tests/test_store.py`), costing maintenance and risking coverage gaps.

## Failure Modes / Edge Cases

- **Silent corruption via contract-violating savers**: the docstrings themselves document the worst case — a `"keep_latest"` prune that severs DeltaChannel ancestor chains causes delta channels to "silently reconstruct as empty (no error raised)" (`base/__init__.py:397-401`). The contract warns but cannot mechanically prevent third-party implementations from doing this unless they run the (optional) conformance suite.
- **Schema/type drift**: `copy_checkpoint` reads `checkpoint.get("pending_sends", [])` (`base/__init__.py:134`) and `empty_checkpoint` writes it (`:824`), yet the `Checkpoint` TypedDict (`:92-123`) declares no such key — the declared schema understates the wire format, weakening schema-time guarantees.
- **Capability misclassification**: signature-based detection of `task_path` support fails for implementations accepting it through `**kwargs` (`_loop.py:1497-1500`).
- **Legacy data tolerance as permanent tax**: `EncryptedSerializer` must forever handle unencrypted blobs (`serde/encrypted.py:29-30`); `maybe_add_typed_methods` wraps old serializers indefinitely (`serde/base.py:40-48`).
- **Uneven conformance adoption**: Postgres — arguably the most-used production saver — has no conformance-suite integration in its test directory (searched `libs/checkpoint-postgres/tests/` for `conformance`/`checkpointer_test`: no matches), so regressions there are caught only by bespoke tests.
- **Deprecated surface still load-bearing**: `ValidationNode` is deprecated (`tool_validator.py:43-46`) and `NodeInterrupt` deprecated in favor of `interrupt()` (`errors.py:110-114`), but both remain functional, requiring dual maintenance paths.

## Future Considerations

- Extend the conformance-package pattern to `BaseStore` (batch semantics, filter operators, TTL refresh) and `BaseCache`, closing the largest substitutability gap; the registration/lifecycle/reporting scaffolding in `initializer.py` and `report.py` is reusable as-is.
- Wire `PostgresSaver`/`AsyncPostgresSaver` into `validate()` alongside sqlite and memory to make conformance universal among first-party savers.
- Replace signature-sniffing with an explicit opt-in marker (e.g., a class attribute or capability method) for parameters like `task_path`, making feature detection deterministic and statically visible.
- Reconcile the `Checkpoint` TypedDict with actual persisted keys (`pending_sends`) or migrate readers off it, restoring schema honesty.
- Graduate or freeze the DeltaChannel beta surface: its warnings currently place substantial correctness burden on implementers of `prune`/`copy_thread`/`delete_for_runs` (`base/__init__.py:340-347,361-371,387-414`).

## Questions / Gaps

- **No evidence found** for a JS/TS mirror of these contracts in this checkout: `libs/sdk-js/` contains only a README (no source tree present), so cross-language contract parity could not be assessed.
- **No evidence found** for conformance-suite adoption beyond sqlite/memory: searches for `conformance` and `checkpointer_test` across `libs/checkpoint-postgres/tests/` returned no matches. If postgres coverage exists elsewhere (e.g., CI workflows invoking `validate()` externally), it was not observable within this source tree.
- Whether third-party (non-LangChain) checkpointer implementations exist and pass the conformance suite could not be verified from this repository alone; the suite's public packaging suggests intent, but external adoption evidence is out of scope for single-source analysis.
- The relationship between `Checkpoint.v = LATEST_VERSION = 2` (`base/__init__.py:811`) and the docstring claim "Currently `1`" (`:96`) appears inconsistent; migration logic keyed on `v` was not located in the searched files, so format-version dispatch behavior remains unclear.

---

Generated by `Dimension 24.02: Interface Contract Design` against `langgraph`.
