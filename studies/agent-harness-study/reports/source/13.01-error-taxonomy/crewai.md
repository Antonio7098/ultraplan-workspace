# Source Analysis: crewai

## Dimension 13.01 — Error Taxonomy

### Source Info

| Field | Value |
|-------|-------|
| Name | crewai |
| Path | `studies/agent-harness-study/sources/crewai` |
| Language / Stack | Python 3.10+ (multi-package monorepo: `crewai`, `crewai-core`, `crewai-files`, `crewai-tools`, `cli`, `devtools`) |
| Analyzed | 2026-08-21 |

## Summary

CrewAI distributes error types across roughly seven local taxonomies rather than one central hierarchy. The largest, most rigorous taxonomy lives in the A2A protocol module (`lib/crewai/src/crewai/a2a/errors.py`): a `IntEnum` of 22 codes, a dataclass base, and 22 typed subclasses, plus routing helpers `is_retryable_error` and `is_client_error`. The agent loop in `crew_agent_executor.py` has three distinct dispatch paths (`OutputParserError` → retry, `LLMContextLengthExceededError`/context-length substring → summarize or `SystemExit`, all else → re-raise). LLM provider modules (openai, anthropic, gemini, bedrock, azure) translate their native exceptions into `LLMContextLengthExceededError`, which the loop catches. There is no project-wide base class: every error extends `Exception`, `ValueError`, `RuntimeError`, `ImportError`, or `TypeError` independently. Detection of the most performance-critical error (context overflow) uses substring matching against `CONTEXT_LIMIT_ERRORS`, not `isinstance`. The two A2A routing helpers are exported but never referenced anywhere in the source. The MCP module classifies failures by a free-form `error_type: str` field on its `FailedEvent` types rather than a typed enum.

## Rating

**5/10** — Clear per-domain taxonomies (notably A2A) with explicit codes, subclasses, and tests, plus a real three-branch dispatch in the agent loop. But no central base class, no `utilities/exceptions/__init__.py` exports (just a docstring), substring-based detection for the most important error class, dead routing helpers, and no documented taxonomy page.

## Evidence Collected

Every entry includes a file path with line numbers from the selected source directory.

| Area | Evidence | File:Line |
|------|----------|-----------|
| LLM context error class | `LLMContextLengthExceededError(Exception)` | `lib/crewai/src/crewai/utilities/exceptions/context_window_exceeding_exception.py:16` |
| Substring-based detection | `CONTEXT_LIMIT_ERRORS: Final[list[str]]` (9 phrases) used by `_is_context_limit_error` | `lib/crewai/src/crewai/utilities/exceptions/context_window_exceeding_exception.py:4-13,33-44` |
| Empty exceptions package | `utilities/exceptions/__init__.py` contains only a docstring (no `__all__`, no re-exports) | `lib/crewai/src/crewai/utilities/exceptions/__init__.py:1` |
| Database error base + templates | `DatabaseOperationError(Exception)` plus `DatabaseError` template-holder class | `lib/crewai/src/crewai/utilities/errors.py:10-52` |
| Agent repository error | `AgentRepositoryError(Exception)` raised at agent lookup | `lib/crewai/src/crewai/utilities/errors.py:55-56` |
| Agent output parse error | `OutputParserError(Exception)` raised in 3 parse-failure paths | `lib/crewai/src/crewai/agents/parser.py:45-59,117-128` |
| Converter/validation error | `ConverterError(Exception)` raised from `to_pydantic` fallbacks | `lib/crewai/src/crewai/utilities/converter.py:28-39,77-82` |
| Tool parsing error | `ToolUsageError(Exception)` raised in 6 tool-call paths | `lib/crewai/src/crewai/tools/tool_usage.py:68-73,834,849,861` |
| Tool quota error | `ToolUsageLimitExceededError(Exception)` raised by `CrewStructuredTool.invoke`/`ainvoke` | `lib/crewai/src/crewai/tools/structured_tool.py:111-112,315,350` |
| Embedding dim error | `EmbeddingDimensionMismatchError(ValueError)` — comment explains why it is intentionally not a `RuntimeError` (would be silently swallowed by background-save plumbing) | `lib/crewai/src/crewai/memory/storage/backend.py:11-41` |
| JSON project error hierarchy | `JSONProjectError(ValueError)` and `JSONProjectValidationError(JSONProjectError)` | `lib/crewai/src/crewai/project/json_loader.py:23-32` |
| Skill parse error | `SkillParseError(ValueError)` for SKILL.md parse failures | `lib/crewai/src/crewai/skills/parser.py:35-36` |
| Skills registry error | `SkillNotCachedError(Exception)` for non-interactive cache miss | `lib/crewai/src/crewai/experimental/skills/registry.py:20-28` |
| Experimental flag error | `ExperimentalFeatureDisabledError(RuntimeError)` | `lib/crewai/src/crewai/experimental/skills/_flag.py:11` |
| Flow ref error | `InvalidRefError(ValueError)` | `lib/crewai/src/crewai/flow/runtime/_refs.py:11` |
| Flow expression error | `ExpressionError(ValueError)` | `lib/crewai/src/crewai/flow/expressions.py:50` |
| Flow script execution error | `FlowScriptExecutionDisabledError(RuntimeError)` | `lib/crewai/src/crewai/flow/runtime/_actions.py:43` |
| RAG client mismatch | `ClientMethodMismatchError(TypeError)` for sync/async misuse | `lib/crewai/src/crewai/rag/core/exceptions.py:4-26` |
| Event-context stack errors | `StackDepthExceededError`, `EventPairingError`, `EmptyStackError` driven by `EventContextConfig` + `MismatchBehavior` enum | `lib/crewai/src/crewai/events/event_context.py:12-18,29-38,176-180,202-205,220-223` |
| Handler-graph cycle error | `CircularDependencyError(Exception)` | `lib/crewai/src/crewai/events/handler_graph.py:15` |
| Optional dependency error | `OptionalDependencyError(ImportError)` (deprecated wrapper around `require`) | `lib/crewai/src/crewai/utilities/import_utils.py:14-46` |
| **A2A enum (most complete)** | `A2AErrorCode(IntEnum)` with 22 codes organized into JSON-RPC 2.0 (-32700 to -32603), A2A-specific (-32001 to -32007), and CrewAI custom (-32009 to -32018) | `lib/crewai/src/crewai/a2a/errors.py:25-99` |
| A2A default messages | `ERROR_MESSAGES: dict[int, str]` keyed by code | `lib/crewai/src/crewai/a2a/errors.py:101-124` |
| A2A base dataclass | `A2AError(Exception)` carrying `code: int`, `message: str | None`, `data: Any`, with `to_dict`/`to_response` | `lib/crewai/src/crewai/a2a/errors.py:127-162` |
| A2A subclasses (typed) | 22 `@dataclass` subclasses each pinning `code` via `field(init=False)` with structured fields (`task_id`, `retry_after`, `requested_types`, etc.) | `lib/crewai/src/crewai/a2a/errors.py:166-441` |
| A2A polling timeout | `A2APollingTimeoutError(A2AClientTimeoutError)` | `lib/crewai/src/crewai/a2a/errors.py:21-22` |
| A2A error response builder | `create_error_response` produces JSON-RPC envelopes | `lib/crewai/src/crewai/a2a/errors.py:443-461` |
| **A2A routing helpers (unused)** | `is_retryable_error(code)` returns true for `INTERNAL_ERROR`, `RATE_LIMIT_EXCEEDED`, `TASK_TIMEOUT`; `is_client_error(code)` for parse/validation/resource-not-found | `lib/crewai/src/crewai/a2a/errors.py:464-503` |
| A2A transport/content negotiation | `TransportNegotiationError`, `ContentTypeNegotiationError` | `lib/crewai/src/crewai/a2a/utils/transport.py:48`; `lib/crewai/src/crewai/a2a/utils/content_type.py:58` |
| A2A UI validation | `A2UIValidationError(Exception)` for A2UI extensions | `lib/crewai/src/crewai/a2a/extensions/a2ui/validator.py:76` |
| A2A HTTP exception re-export | Local `HTTPException` redefinition inside `server_schemes.py` | `lib/crewai/src/crewai/a2a/auth/server_schemes.py:52` |
| **Agent loop dispatch (sync)** | Catches `OutputParserError` → `handle_output_parser_exception`; `Exception` → re-raise litellm, detect context via substring, summarize or `SystemExit`, else re-raise | `lib/crewai/src/crewai/agents/crew_agent_executor.py:434-458` |
| **Agent loop dispatch (async)** | Same three-branch pattern mirrored | `lib/crewai/src/crewai/agents/crew_agent_executor.py:1272-1296` |
| Lite-agent dispatch | Mirrors the same three-branch pattern | `lib/crewai/src/crewai/lite_agent.py:932-959` |
| Experimental agent executor | Uses return-string router keys `"parser_error"`, `"context_error"` from the exception type | `lib/crewai/src/crewai/experimental/agent_executor.py:1435-1449` |
| Core agent converter catch | `except ConverterError` in `Agent.core` | `lib/crewai/src/crewai/agent/core.py:1770` |
| LiteAgent converter catch | `except ConverterError as e` | `lib/crewai/src/crewai/lite_agent.py:671` |
| Training converter catch | `except ConverterError` | `lib/crewai/src/crewai/utilities/training_converter.py:32` |
| JSON loader catches | Four `except JSONProjectError` sites | `lib/crewai/src/crewai/project/json_loader.py:655,1201,1963,2021` |
| LLM substring-based reclassify | `LLMContextLengthExceededError._is_context_limit_error(error_msg)` used to wrap arbitrary provider exceptions | `lib/crewai/src/crewai/llm.py:1051-1056`; mirror paths at `1244-1249`, `1395-1400`, `1667-1672`, `1869-1873`, `2004-2005` |
| Provider-specific raises | OpenAI: `lib/crewai/src/crewai/llms/providers/openai/completion.py:940,1083,1780,2203`; Anthropic: `:920,1413,1470,1831`; Gemini: `:1171,1264`; Bedrock: `:453,581`; Azure: `:462` |
| Database `sqlite3.Error` wrapping | All 5 storage operations wrap `sqlite3.Error` into `DatabaseOperationError` | `lib/crewai/src/crewai/memory/storage/kickoff_task_outputs_storage.py:64,111,160,201,222` |
| MCP error events (string `error_type`) | `MCPConnectionFailedEvent.error_type`, `MCPToolExecutionFailedEvent.error_type`, `MCPConfigFetchFailedEvent.error_type` typed as `str | None` with docstring examples like `"timeout"`, `"authentication"`, `"network"` | `lib/crewai/src/crewai/events/types/mcp_events.py:46-98` |
| MCP emitters | Concrete error_type values: `"not_connected"`, `"connection_failed"`, `"tool_error"` | `lib/crewai/src/crewai/mcp/tool_resolver.py`; `lib/crewai/src/crewai/mcp/client.py` |
| Failed-event taxonomy on event bus | `LLMCallFailedEvent`, `AgentExecutionErrorEvent`, `LiteAgentExecutionErrorEvent`, `ToolUsageErrorEvent`, `CrewKickoffFailedEvent`, `TaskFailedEvent`, `MethodExecutionFailedEvent`, `FlowFailedEvent`, `Memory*FailedEvent`, `Knowledge*FailedEvent`, `A2A*FailedEvent`, `AgentReasoningFailedEvent`, `AgentEvaluationFailedEvent`, `A2AAuthenticationFailedEvent`, `LiteAgentExecutionErrorEvent` | `lib/crewai/src/crewai/events/types/*.py` |
| Test coverage — context length | `pytest.raises(LLMContextLengthExceededError)` in `test_llm.py` | `lib/crewai/tests/test_llm.py:383-398` |
| Test coverage — embedding dim | `EmbeddingDimensionMismatchError` exercised in 8 assertions | `lib/crewai/tests/memory/test_dimension_mismatch.py` |
| Test coverage — JSON project | `JSONProjectValidationError` used in tests | `lib/crewai/tests/project/test_json_loader.py` |
| Test coverage — skill parse | `SkillParseError` used in tests | `lib/crewai/tests/skills/test_parser.py` |
| Test coverage — agent repository | `AgentRepositoryError` used in tests | `lib/crewai/tests/agents/test_agent.py` |
| Test coverage — tool usage limit | `ToolUsageLimitExceededError` flow with `max_usage_count` | `lib/crewai/tests/tools/test_tool_usage_limit.py` |
| Test coverage — `is_retryable_error` / `is_client_error` | None found | (no tests reference these helpers) |
| Top-level utilities re-exports | `utilities/__init__.py` only re-exports `LLMContextLengthExceededError` and `ConverterError` | `lib/crewai/src/crewai/utilities/__init__.py:3-26` |

## Answers to Dimension Questions

1. **Are errors classified by source?** Partially. There is no single central taxonomy, but each subsystem defines its own. The clearest source classification exists in: (a) the A2A protocol with 22 codes segmented into JSON-RPC standard / A2A-specific / CrewAI-custom, (b) the LLM layer where every provider translation funnels into `LLMContextLengthExceededError`, (c) the agent loop which has three buckets — parser (`OutputParserError`), context (substring-detected), unknown (everything else), (d) the event bus with paired `Started`/`Completed`/`Failed` events per subsystem. There is no project-wide base class or enum tying these together.

2. **Is the taxonomy used for handling?** Yes in the agent loop (`crew_agent_executor.py:434-458` and `:1272-1296`), and yes in `lite_agent.py:932-959` and `experimental/agent_executor.py:1435-1449`. The dispatcher branches on `OutputParserError` (recover), context-length (recover-or-abort), and `Exception` (re-raise after logging). The A2A module *defines* the most mature taxonomy but its routing helpers (`is_retryable_error`, `is_client_error`) are unused — no caller branches on them. Storage code uniformly maps `sqlite3.Error` to `DatabaseOperationError` (5 sites in `kickoff_task_outputs_storage.py`). The LLM layer translates every provider's native context exception into the unified `LLMContextLengthExceededError` before bubbling it up.

3. **Are error categories documented?** No project-level taxonomy doc. Categories are documented *per file* via docstrings (e.g., `events/event_context.py` describes the `MismatchBehavior` enum; `a2a/errors.py` has a module docstring classifying the JSON-RPC ranges; `memory/storage/backend.py` documents the rationale for `EmbeddingDimensionMismatchError` not being a `RuntimeError`). The Edge docs mention failure events only in the event-listener page (`docs/edge/en/concepts/event-listener.mdx`), listing `ToolUsageErrorEvent`, `LLMCallFailedEvent`, `AgentExecutionErrorEvent`, etc., as bus events without explaining the underlying exception hierarchy. The CLI template generator surfaces `respect_context_window` only as a commented-out JSON key.

4. **Can new error types be added without breaking existing handling?** Yes for the agent loop — `except Exception` catches anything new and `handle_unknown_error` logs and re-raises (`lib/crewai/src/crewai/utilities/agent_utils.py:632-657`). This is robust but at the cost of every new error being treated as fatal unless someone also adds a branch to the dispatch. For the A2A module, adding a new code requires: a new enum entry, a default in `ERROR_MESSAGES`, a new dataclass subclass — three sites to touch (`lib/crewai/src/crewai/a2a/errors.py:25-124,166-441`) plus the `is_retryable_error`/`is_client_error` allowlists if the new code needs routing semantics. For MCP, the `error_type` field is free-form (`str | None`) so adding a value is a no-op but no type system catches typos.

## Architectural Decisions

- **No central base class.** Each domain (`utilities/errors`, `utilities/exceptions/`, `a2a/errors.py`, `tools/`, `memory/storage/`, `flow/`, `rag/core/`, `events/event_context.py`, etc.) defines its own `Exception` subclass directly. The package `utilities/exceptions/__init__.py` is empty (one docstring line). The trade-off is locality: a domain's failure surface is self-explanatory, but cross-domain routing requires manual knowledge of every module.
- **A2A module is the gold standard** — `A2AErrorCode(IntEnum)` + dataclass base + 22 typed subclasses with structured fields + default-message dictionary + routing helpers + JSON-RPC serialization. This is the only place in the codebase that uses an `IntEnum` for error codes.
- **Provider abstraction via translation, not inheritance.** Each LLM provider (`openai`, `anthropic`, `gemini`, `bedrock`, `azure`) explicitly raises `LLMContextLengthExceededError` from native exceptions inside `except` blocks. The translation list lives in `CONTEXT_LIMIT_ERRORS` (`lib/crewai/src/crewai/utilities/exceptions/context_window_exceeding_exception.py:4-13`).
- **Substring detection over isinstance for context length.** The agent loop calls `is_context_length_exceeded(e)` which delegates to `LLMContextLengthExceededError._is_context_limit_error(str(e))` — 9 substring phrases. The rationale (implicit) is that litellm re-exposes many provider exceptions whose Python type is not stable, so string matching is more robust than class matching. The cost is fragility to provider message changes.
- **Event bus as parallel error surface.** Every async subsystem (LLM, agent, tool, MCP, memory, knowledge, A2A, flow, evaluation) has `Started`/`Completed`/`Failed` events. The agent loop still raises Python exceptions synchronously, but downstream listeners (CLI TUI, tracing, evaluation, telemetry) consume the failed-event stream. The two paths coexist without a unified envelope.
- **Embedding mismatch is intentionally `ValueError`, not `RuntimeError`.** Explicit comment in `memory/storage/backend.py:18-21` — a `RuntimeError` would be silently swallowed by background-save plumbing. This is a domain-specific override of a sensible default.

## Notable Patterns

- **Three-branch agent loop dispatch** duplicated across 4+ executors with identical structure (sync, async, lite, experimental). Centralization exists in helpers — `handle_output_parser_exception`, `handle_context_length`, `handle_unknown_error`, `is_context_length_exceeded` — all in `lib/crewai/src/crewai/utilities/agent_utils.py:632-749`.
- **Dataclass + `field(init=False)`** pattern in A2A to lock each subclass to a specific enum code: `code: int = field(default=A2AErrorCode.X, init=False)` (`lib/crewai/src/crewai/a2a/errors.py:166-441`).
- **Default messages via dict lookup** instead of conditional `__post_init__` chains, except where structured fields like `task_id` produce a richer message (`lib/crewai/src/crewai/a2a/errors.py:101-124,141-144,186-189,223-226`).
- **Configurable error policy** in event context: `MismatchBehavior(Enum)` with `WARN | RAISE | SILENT` and `EventContextConfig` (`lib/crewai/src/crewai/events/event_context.py:12-26`) lets deployments choose between strict and permissive.
- **Free-form `error_type: str` on MCP failure events** rather than a typed enum — the docstring gives examples (`"timeout"`, `"authentication"`, `"network"`, `"not_connected"`, `"api_error"`, `"connection_failed"`, `"tool_error"`, `"validation"`, `"server_error"`) with `etc.` indicating openness.
- **No base class for failure events** either — every `*FailedEvent` is a `BaseEvent` subclass with `error: str` plus optional `error_type: str | None`.

## Tradeoffs

- **Local taxonomies vs. global classification.** Each domain owns its exceptions, which is easy to extend but means callers must know multiple modules. The agent loop's `except Exception` fallback means an unhandled error is never lost but is also never classified.
- **Substring matching vs. type checking.** Catches unknown provider exception classes but breaks silently if a provider changes its error message text. `CONTEXT_LIMIT_ERRORS` would need a manual update.
- **A2A helper code is unused.** `is_retryable_error` and `is_client_error` exist but no code branches on them. This is either dead API surface, intended library surface, or work-in-progress; no evidence in source of intent.
- **Tool quota errors are uncaught.** `ToolUsageLimitExceededError` is raised in `structured_tool.py:315,350` but no `except` clause in the agent loop catches it. The agent will propagate the exception and abort the run.
- **`SystemExit` as abort signal.** `handle_context_length` raises `SystemExit` when `respect_context_window=False` (`utilities/agent_utils.py:747-749`). This makes the abort hard to catch in user code without `except SystemExit` and ties "should retry vs. stop" to a hard process-exit signal.

## Failure Modes / Edge Cases

- **Provider message drift.** A new OpenAI/Anthropic/Gemini error message that no longer contains any of the 9 `CONTEXT_LIMIT_ERRORS` substrings will fail to be classified as context-length-exceeded, fall through to `handle_unknown_error`, and abort.
- **MCP error_type typos.** `error_type="timout"` (typo) on an `MCPToolExecutionFailedEvent` will not match any consumer expecting `"timeout"`. There is no enum to constrain it.
- **Tool quota collisions.** Two tools with the same name from different sources could double-count `current_usage_count` against `max_usage_count` because the counter is per-`CrewStructuredTool` instance, not per-name. (`lib/crewai/src/crewai/tools/structured_tool.py:138-139`).
- **Embedding dim silently catches some raises.** Comment at `memory/storage/backend.py:18-21` notes that `RuntimeError` is silently swallowed; the class is therefore `ValueError`. If a future refactor moves it to `RuntimeError`, the actionable migration message disappears.
- **`A2AError` re-init.** Because `A2AError` is a dataclass but extends `Exception`, Python's dataclass+exception interaction allows construction without args except `code`. Callers must always pass `code`.
- **Empty exceptions package.** Importing from `crewai.utilities.exceptions` resolves to an empty module, so consumers must import from the leaf modules (`context_window_exceeding_exception`). This is documented but easy to miss.

## Future Considerations

- The 22-code A2A enum could be the seed for a unified `CrewAIErrorCode` enum plus a `BaseCrewAIError` class so that consumers can `except BaseCrewAIError` and inspect `error.code`.
- `is_retryable_error` and `is_client_error` should either be called somewhere (most natural place: an HTTP/a2a retry layer) or removed.
- `ToolUsageLimitExceededError` needs a catch site in the agent loop; the natural place is the same `except Exception` branch in `crew_agent_executor.py:444-458` with a dedicated handler.
- `CONTEXT_LIMIT_ERRORS` should either be kept current via provider-integration tests or replaced with `isinstance` against a stable union of types.
- `error_type` on MCP failure events would benefit from a `Literal[...]` (or small `Enum`) to make typos compile-time errors instead of runtime mismatches.
- A single `docs/errors.mdx` summarizing the dispatch matrix (which error → which handler → which outcome) would close the documentation gap.

## Questions / Gaps

- Is there an intent behind `is_retryable_error`/`is_client_error` being defined but unused? Searched `lib/` for callers: no references found in source or tests.
- Is the empty `utilities/exceptions/__init__.py` placeholder for a future unified export module, or is the convention "import from leaf modules"?
- Are MCP `error_type` strings expected to remain free-form or is there an upstream type definition in `crewai-tools`?
- Is `SystemExit` raised in `handle_context_length` an intentional public contract or an internal abort signal that should be migrated to a typed error?
- Where is `DatabaseOperationError` caught? Searched: it is raised in `kickoff_task_outputs_storage.py` but no `except DatabaseOperationError` exists in `lib/`. Callers must rely on `except Exception`.

---

Generated by `13.01-error-taxonomy` against `crewai`.