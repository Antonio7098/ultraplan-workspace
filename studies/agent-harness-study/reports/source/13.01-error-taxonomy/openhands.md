# Source Analysis: openhands

## Dimension 13.01: Error Taxonomy

### Source Info

| Field | Value |
|-------|-------|
| Name | openhands |
| Path | `studies/agent-harness-study/sources/openhands` |
| Language / Stack | Python 3.12 (FastAPI app server), TypeScript/React frontend; agent runtime delegated to external `software-agent-sdk` |
| Analyzed | 2026-08-21 |

## Summary

OpenHands does not have a single unified error taxonomy. Instead, each subsystem defines its own family of exception classes or enums, and the categories map cleanly onto the dimension's source list: **provider** errors (`AuthenticationError`, `RateLimitError`, `ProviderTimeoutError`, `ResourceNotFoundError`, `UnknownException` for git providers), **model/auth** errors (`LLMAuthenticationError`), **user/settings** errors (`MissingSettingsError`, `SessionExpiredError`), **infrastructure** status (`SandboxStatus.ERROR`), **policy/rate-limiting** (`RateLimitException` + `Retry-After` headers), and a dedicated **Slack integration taxonomy** with traceable string codes (`SLACK_ERR_00x`) plus a centralized user-message registry.

Classification is real and used for routing: HTTP status codes are translated into typed exceptions at one choke point (`handle_http_status_error` / `handle_http_error` in the shared `HTTPClient` base), FastAPI registers per-class exception handlers, integration managers dispatch on exception type to craft user-facing remediation messages, and sandbox `SandboxStatus` enum values are routed to distinct HTTP codes (409/410/503). The strongest part of the taxonomy is the Slack module: an enum of unique error codes with structured log context and a fallback message. The weakest parts are fragmentation (two unrelated `AuthError` classes, three duplicated `StartingConvoException` definitions, sibling errors with no common marker base beyond `ValueError`), one category (`LLMAuthenticationError`) whose raise site lives outside this repository, no central registry or documentation of the taxonomy, and string-sniffing classification in the frontend.

Can you tell from the error type whether to retry, escalate, or stop? Partially. Timeout and rate-limit categories carry retry semantics (`Retry-After` headers, user messages that say "try again"), auth/settings categories escalate to the user ("re-login", "set a valid LLM API key"), and infrastructure `ERROR` state stops message acceptance (503). But the decision is made ad hoc at each call site — there is no policy engine keyed off error class.

## Rating

**7 / 10** — Clear, tested model with explicit interfaces and operational safeguards, held back by fragmentation and weak documentation.

- Earns 7: typed classification by source exists and is exercised by tests (`tests/unit/integrations/protocols/test_http_client.py:126-220` covers every branch of the status→error mapping); dispatch is type-based, not string-based, on the backend; catch-all categories (`UnknownException`, `UNEXPECTED_ERROR`) prevent unclassified failures from crashing flows; sandbox status gating is explicit.
- Not 8+: no shared base classes or registry across subsystems (a generic `except ValueError` over-catches because provider errors derive from `ValueError`); duplicate class names across modules invite wrong-import bugs; `LLMAuthenticationError` is caught in seven places but never raised inside this repo; no documentation page describes the taxonomy; frontend classifies billing errors via substring matching.

## Evidence Collected

Every entry cites workspace-relative paths into the selected source directory.

| Area | Evidence | File:Line |
|------|----------|-----------|
| API-layer error hierarchy | `OpenHandsError(HTTPException)` base with baked-in status codes; subclasses `AuthError` (401), `PermissionsError` (403), `SandboxError` (500) | `studies/agent-harness-study/sources/openhands/openhands/app_server/errors.py:6-43` |
| User/settings/model error types | `MissingSettingsError`, `LLMAuthenticationError`, `SessionExpiredError`, all deriving from `ValueError` | `studies/agent-harness-study/sources/openhands/openhands/app_server/types.py:19-34` |
| Provider error taxonomy | `AuthenticationError`, `UnknownException`, `RateLimitError`, `ProviderTimeoutError`, `ResourceNotFoundError` as siblings under `ValueError` | `studies/agent-harness-study/sources/openhands/openhands/app_server/integrations/service_types.py:171-198` |
| Classification code (HTTP → typed error) | `handle_http_status_error`: 401→`AuthenticationError`, 404→`ResourceNotFoundError`, 429→`RateLimitError`, else `UnknownException`; `handle_http_error`: `TimeoutException`→`ProviderTimeoutError`, else `UnknownException` | `studies/agent-harness-study/sources/openhands/openhands/app_server/integrations/protocols/http_client.py:82-110` |
| Classification tests | One test per mapping branch incl. timeout subtypes (`ReadTimeout`, `WriteTimeout` → `ProviderTimeoutError`) | `studies/agent-harness-study/sources/openhands/tests/unit/integrations/protocols/test_http_client.py:126-220` |
| Raise-site pattern | Handlers *return* typed exceptions; callers re-raise them after token-refresh retry on 401 | `studies/agent-harness-study/sources/openhands/openhands/app_server/integrations/github/service/base.py:74-95` |
| Enterprise auth taxonomy | `AuthError` base + `NoCredentialsError`, `EmailNotVerifiedError`, `BearerTokenError`, `CookieError`, `TosNotAcceptedError`, `ExpiredError`, `TokenRefreshError` | `studies/agent-harness-study/sources/openhands/enterprise/server/auth/auth_error.py:1-46` |
| Slack error-code enum | `SlackErrorCode` with unique traceable codes (`SLACK_ERR_001`…`SLACK_ERR_999`) spanning session, redis, provider-timeout/provider-auth, LLM-auth, settings, scopes | `studies/agent-harness-study/sources/openhands/enterprise/integrations/slack/slack_errors.py:18-30` |
| Slack centralized handling | `SlackError` carries code + `message_kwargs` + `log_context`; `_USER_MESSAGES` registry maps every code to a user message with `UNEXPECTED_ERROR` fallback | `studies/agent-harness-study/sources/openhands/enterprise/integrations/slack/slack_errors.py:33-134` |
| Slack dispatch | `receive_message` catches `SlackError` → `handle_slack_error`; any other exception is wrapped as `UNEXPECTED_ERROR`; handler picks log level by code and sends ephemeral message | `studies/agent-harness-study/sources/openhands/enterprise/integrations/slack/slack_manager.py:441-466`, `613-652` |
| Integration manager dispatch | GitHub job starter catches `MissingSettingsError` / `LLMAuthenticationError` / `(AuthenticationError, ExpiredError, SessionExpiredError)` and maps each to a different remediation message; bare `except Exception` last resort | `studies/agent-harness-study/sources/openhands/enterprise/integrations/github/github_manager.py:408-435` |
| Cross-taxonomy translation | Slack route converts provider-typed `ProviderTimeoutError` into Slack's `SlackErrorCode.PROVIDER_TIMEOUT` | `studies/agent-harness-study/sources/openhands/enterprise/server/routes/integration/slack.py:426-436` |
| Infrastructure status taxonomy | `SandboxStatus` enum `RUNNING/PAUSED/ERROR/MISSING`; remote service maps provider strings incl. `'error'→ERROR`; Docker maps `'dead'→ERROR` and defaults unknown states to `ERROR` | `studies/agent-harness-study/sources/openhands/openhands/app_server/sandbox/sandbox_models.py:9-14`, `remote_sandbox_service.py:56-60`, `docker_sandbox_service.py:118-127` |
| Status-based routing to HTTP | `MISSING→410 GONE`, `ERROR→503`, other non-RUNNING→`409 CONFLICT` with resume hint; documented in OpenAPI responses | `studies/agent-harness-study/sources/openhands/openhands/app_server/app_conversation/app_conversation_router.py:430-438`, `509-529` |
| Sandbox wait-loop uses ERROR state | Poll loop raises `SandboxError` immediately when sandbox enters ERROR vs continuing to poll while STARTING | `studies/agent-harness-study/sources/openhands/openhands/app_server/sandbox/sandbox_service.py:117-140` |
| Policy/rate-limit taxonomy | Enterprise `RateLimitException(HTTPException)` carries `RateLimitResult` with `retry_after`; handler emits 429 + `X-RateLimit-*`/`Retry-After` headers; OSS middleware returns 429 + `Retry-After: 1` | `studies/agent-harness-study/sources/openhands/enterprise/server/rate_limit.py:32-47`, `109-137`; `studies/agent-harness-study/sources/openhands/openhands/app_server/middleware.py:115-131` |
| Callback failure taxonomy | `EventCallbackStatus` = ACTIVE/DISABLED/COMPLETED/ERROR; failures recorded as `EventCallbackResultStatus.ERROR` rows with `detail=str(exc)` instead of propagating | `studies/agent-harness-study/sources/openhands/openhands/app_server/event_callback/event_callback_models.py:33-37`; `sql_event_callback_service.py:241-250` |
| Implicit retry via taxonomy absence | Title processor logs transient HTTP errors without failing and returns `None`, leaving callback ACTIVE so later events retry | `studies/agent-harness-study/sources/openhands/openhands/app_server/event_callback/set_title_callback_processor.py:67-73`, `132-138` |
| Timeout prioritization | `Provider.get_repositories` re-raises `ProviderTimeoutError` but swallows all other per-provider exceptions | `studies/agent-harness-study/sources/openhands/openhands/app_server/integrations/provider.py:300-308` |
| Rate-limit routing inside services | GitLab service catches `RateLimitError` and downgrades to `WebhookStatus.RATE_LIMITED` result rather than raising | `studies/agent-harness-study/sources/openhands/enterprise/integrations/gitlab/gitlab_service.py:351-355` |
| Session-expiry raise site | Token manager raises `SessionExpiredError` during refresh flow | `studies/agent-harness-study/sources/openhands/enterprise/server/auth/token_manager.py:580` |
| Frontend classification (string-based) | `isBudgetOrCreditError` classifies budget/credit failures by substring regex over the raw error message | `studies/agent-harness-study/sources/openhands/frontend/src/utils/error-handler.ts:40-55` |
| Tests of Slack error handling | `TestHandleSlackError` verifies ephemeral delivery, graceful degradation when view creation fails, and per-code messages | `studies/agent-harness-study/sources/openhands/enterprise/tests/unit/test_slack_integration.py:1171-1250` |

## Answers to Dimension Questions

1. **Are errors classified by source?**
   Yes, per subsystem. Provider-source errors are a five-class family keyed off HTTP semantics (`...integrations/service_types.py:171-198`, classified in `http_client.py:82-110`). Model/LLM auth has its own class (`types.py:25`), user settings its own (`types.py:19`), enterprise identity its own hierarchy (`auth_error.py:1-46`), infrastructure a status enum (`sandbox_models.py:9-14`), policy violations a dedicated exception with retry metadata (`rate_limit.py:109-120`), and Slack an enum-coded taxonomy covering session/user/provider/LLM/infra sources (`slack_errors.py:18-30`). There is no cross-cutting "source" attribute on a common error base — classification is by inheritance tree position within each module, not a shared dimension field.

2. **Is the taxonomy used for handling?**
   Yes. Type-based dispatch appears at four layers: FastAPI exception handlers registered per class (`openhands/app_server/app.py:63-68` for `AuthenticationError`→401; `enterprise/saas_server.py:192-203` for `NoCredentialsError`/`ExpiredError`→401; `rate_limit.py:25-29,123-137` for `RateLimitException`→429+headers); integration managers whose except-chains pick distinct user guidance per category (`github_manager.py:408-435`, replicated across slack/jira/gitlab/bitbucket managers); status-gated routing where `SandboxStatus.ERROR` blocks messaging with 503 while PAUSED yields 409-with-resume-hint (`app_conversation_router.py:508-529`); and selective propagation where timeouts outrank other provider errors (`provider.py:304-306`).

3. **Are error categories documented?**
   Weakly. Every class has a one-line docstring (e.g., `http_client.py:87`, `service_types.py:184`), the Slack module carries a good module-level docstring explaining the code/message design (`slack_errors.py:1-7`), and OpenAPI response blocks document HTTP outcomes per route (`app_conversation_router.py:430-438`, `478-483`). But no document enumerates the taxonomy, explains which categories are retryable, or records the cross-repo contract with the SDK. A repo-wide search of `*.md` for "error code", "error taxonomy", "SLACK_ERR", and "error handling" returned only generic testing guidance (`AGENTS.md:270,274`) and an authoring checklist (`skills/add_agent.md:33`).

4. **Can new error types be added without breaking existing handling?**
   Mostly yes, with caveats. Safe paths exist: new `SlackErrorCode` members fall back to `UNEXPECTED_ERROR` messaging via `_USER_MESSAGES.get(code, ...)` (`slack_errors.py:123-125`); unknown provider statuses default to `SandboxStatus.ERROR` (`docker_sandbox_service.py:127`); unhandled backend exceptions are caught by broad handlers and converted to generic user messages (`github_manager.py:431-435`, `slack_manager.py:458-466`). Caveats: adding a subclass does nothing until someone adds an except-clause or handler — there is no exhaustive dispatch (e.g., a `match` on code) to update; the sibling classes share no marker base, so callers cannot catch "any provider error" without listing each class (`github_manager.py:422` must enumerate three classes); and name collisions (`AuthError` in both `errors.py:18` and `auth_error.py:1`) make wrong imports a realistic breakage mode.

## Architectural Decisions

1. **Classify once, at the transport boundary.** All git-provider HTTP traffic funnels through `HTTPClient.handle_http_status_error`/`handle_http_error` (`http_client.py:82-110`), so GitHub/GitLab/Bitbucket/Azure DevOps services inherit identical classification. This is the single most load-bearing piece of taxonomy in the OSS side, and it is directly unit-tested (`test_http_client.py:126-220`).

2. **Return-and-raise convention.** The classifier methods return exception instances which callers explicitly raise (`github/service/base.py:92-95`). This keeps the classifier pure/testable but relies on caller discipline; forgetting the `raise` silently drops errors.

3. **Errors as UX, not just control flow.** Integration managers treat each error category as a distinct user-communication requirement — missing settings → "re-login", LLM auth → "set a valid API key", expired session → session message (`github_manager.py:408-427`). The taxonomy's primary consumer is user-facing remediation text.

4. **Per-subsystem taxonomies instead of a global one.** Slack got a redesigned enum-code system with structured logging (`slack_errors.py:18-61`), while org management uses ~15 plain exception classes (`org_models.py:29-144`). There was evidently no mandate to unify; newer modules (Slack) show more maturity than older ones.

5. **Status enums for long-running infrastructure, exceptions for requests.** Sandbox health uses the `SandboxStatus` enum polled and persisted (`sandbox_service.py:100-140`), while request-scoped failures use exceptions. The two meet at routing boundaries (`app_conversation_router.py:509-529`).

6. **Fail-soft defaults for observability.** Event-callback failures become persisted ERROR results rather than raised exceptions (`sql_event_callback_service.py:241-250`), and rate-limit backend outages swallow the check rather than blocking traffic (`rate_limit.py:67-68`) — availability is chosen over strictness in both cases.

## Notable Patterns

- **Traceable error codes**: `SLACK_ERR_001`…`999` values are surfaced to end users ("ref: {code}") so support can correlate a user report with a log entry (`slack_errors.py:74-81`, `107-109`) — the only place in the repo with user-visible error correlation IDs.
- **Cross-taxonomy adapters**: typed provider errors are translated into Slack codes at the route layer (`slack.py:426-436`), showing deliberate boundary conversion rather than leaking one taxonomy into another.
- **Category-pair catching**: `(AuthenticationError, ExpiredError, SessionExpiredError)` grouped when they share handling (`github_manager.py:422`) — pragmatic use of the taxonomy without forcing artificial hierarchies.
- **Retry-after as data**: `RateLimitResult.retry_after` computed from window stats and emitted both as a header and available programmatically (`rate_limit.py:39-47`, `89-96`).
- **Graceful-degradation hooks**: skill loader treats org-level config fetch failures (including `AuthenticationError`) as "feature absent" and returns `None` (`skill_loader.py:180-194`).

## Tradeoffs

- **Fragmentation vs. autonomy.** Each subsystem owning its taxonomy avoids cross-module coupling but makes it impossible to answer repo-wide questions like "list all retryable errors" or to write one `except TaxonomyError` handler. Two unrelated `AuthError` types coexist (`errors.py:18`, `auth_error.py:1`).
- **ValueError bases cut both ways.** Deriving provider errors from `ValueError` (`service_types.py:171-198`) means existing `except ValueError` code catches them automatically, but also means taxonomy members cannot be distinguished from accidental `ValueError`s, and no common base exists to catch the family.
- **Catch-all safety vs. information loss.** Broad `except Exception` handlers convert unknowns into friendly messages (`github_manager.py:431-435`) but discard the distinction between bug and environment failure; only logs retain it.
- **Fail-soft callbacks vs. silent loss.** Recording callback errors as rows keeps processing alive (`sql_event_callback_service.py:241-250`), but nothing shown re-drives failed callbacks automatically — recovery depends on the next matching event arriving.
- **Out-of-repo contract vs. dead-looking code.** Catching `LLMAuthenticationError` in seven managers encodes a contract with the external software-agent-sdk (documented in `AGENTS.md`'s SDK references) even though no raise site exists here; readers auditing only this repo see an apparently dead branch.

## Failure Modes / Edge Cases

- **Unclassified provider statuses**: any non-401/404/429 HTTP error becomes `UnknownException` (`http_client.py:98-99`) — e.g., 502/503 gateway blips, which are typically retryable, get neither retry semantics nor a distinct category.
- **String-based frontend classification**: `isBudgetOrCreditError` substring-matches provider error text (`error-handler.ts:45-55`); wording changes upstream silently reclassify billing failures. The docstring itself admits provider-side "credits" mentions can be misread.
- **Duplicated definitions drift risk**: three independent `StartingConvoException` classes (`slack_types.py:136`, `jira_types.py:53`, `jira_dc_types.py:37`) must be kept semantically aligned by hand.
- **Classifier-return misuse**: if a future service forgets to `raise self.handle_http_status_error(e)` (pattern at `github/service/base.py:92-93`), errors vanish into returned values.
- **Callback ERROR is terminal-ish**: a callback left ACTIVE after a recorded ERROR result (`sql_event_callback_service.py:241-250`) will keep failing on every subsequent event unless the processor self-disables; no backoff is visible.
- **Default-to-ERROR masking**: Docker statuses not in the mapping table collapse to `SandboxStatus.ERROR` (`docker_sandbox_service.py:125-127`), which triggers the hard 503 stop path even for transient/unrecognized states.

## Future Considerations

- Introduce a shared marker base (or a `source:` field enum: model/provider/tool/policy/context/user/infrastructure/timeout) so families can be caught collectively and new categories can declare retryability once.
- Consolidate duplicate names (`AuthError`, `StartingConvoException`) behind canonical imports to eliminate wrong-import hazards.
- Document the taxonomy: one page enumerating categories, their retry/escalate/stop semantics, and the SDK-boundary contract for `LLMAuthenticationError`.
- Extend the `HTTPClient` classifier with a retryable bucket for 5xx/gateway errors instead of folding them into `UnknownException`.
- Replace frontend substring sniffing with structured error fields (the Slack `code` pattern is the in-repo precedent to copy).
- Add automatic re-drive/backoff for callbacks in `EventCallbackResultStatus.ERROR` state.

## Questions / Gaps

- Where exactly does the external software-agent-sdk raise `LLMAuthenticationError`? No evidence found inside this repository (searched all non-test `.py` files for a raise site; only definition `openhands/app_server/types.py:25`, imports, except-clauses, and test mocks). The contract is inferred from AGENTS.md references to the separate SDK repo, which is outside this study's isolation boundary.
- Is there a documented SLO or operational runbook keyed on error codes? None found; searched markdown for "error taxonomy"/"error code" with no hits beyond style guidance.
- Whether legacy V0/V1 server code (`openhands/server/`) maintains additional error classes was only partially surveyed; the V1 app server (`openhands/app_server/`) re-exports the studied types (`openhands/server/types.py:8-19`), suggesting convergence, but a full V0 audit was out of scope for this pass.

---

Generated by `13.01-error-taxonomy` against `openhands`.
