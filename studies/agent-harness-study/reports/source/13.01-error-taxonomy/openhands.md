# Source Analysis: openhands

## Dimension 13.01: Error Taxonomy

### Source Info

| Field | Value |
|-------|-------|
| Name | openhands |
| Path | `studies/agent-harness-study/sources/openhands` |
| Language / Stack | Python (FastAPI backend) + React/TypeScript frontend |
| Analyzed | 2026-08-21 |

## Summary

OpenHands is a large Python/React project with no single, unified error taxonomy. Instead, the source has at least **four parallel error class hierarchies**, scattered across both the open-source `openhands/` package and the `enterprise/` add-on. The strongest evidence of source-based classification lives in three places:

1. The git-provider HTTP client maps HTTP status codes to four named `ValueError` subclasses (`AuthenticationError`, `ResourceNotFoundError`, `RateLimitError`, `UnknownException`) plus a `ProviderTimeoutError` for `httpx.TimeoutException` (`openhands/app_server/integrations/protocols/http_client.py:82-109`).
2. The Slack integration has the only formal enum-based taxonomy in the codebase: `SlackErrorCode` (`SLACK_ERR_001`…`SLACK_ERR_999`) carried by a single `SlackError` class with a centralized user-message dispatch table (`enterprise/integrations/slack/slack_errors.py:18-134`).
3. The webhook router classifies terminal conversation errors into six string categories for analytics dashboards (`budget_exceeded`, `model_error`, `runtime_error`, `timeout`, `user_cancelled`, `unknown`), documented in `openhands/analytics/EVENTS.md:25-35` and implemented at `openhands/app_server/event_callback/webhook_router.py:70-90`.

Beyond these, the codebase defines ~70+ bespoke exception classes across `openhands/app_server/errors.py`, `openhands/app_server/types.py`, `openhands/app_server/settings/llm_profiles.py`, `enterprise/server/auth/auth_error.py`, `enterprise/server/routes/org_models.py`, and the integration subpackages. The retry-vs-don't-retry decision is encoded via `tenacity.retry_if_exception_type` on `KeycloakConnectionError` (`enterprise/server/auth/saas_user_auth.py:292-298`, `enterprise/server/auth/token_manager.py:170-173`) — this is the only place where the source category explicitly drives control flow.

Overall, the taxonomy is **present but inconsistent**: source categories exist in patches, are partially used for routing (HTTP→exception, retryable vs. non-retryable, SlackErrorCode→message), and are partially documented (`EVENTS.md` only). The thin `OpenHandsError` base in `errors.py:6-43` covers only `AuthError`, `PermissionsError`, `SandboxError` — most other domain errors either inherit from `ValueError`, `Exception`, or `HTTPException`, making a `try/except OpenHandsError` impossible across the codebase.

## Rating

**6 / 10 — Present but inconsistent, weakly documented, fragile.**

The rating rests on:

- **Source classification exists** in the HTTP-client error mapper (`http_client.py:82-109`), Slack `SlackErrorCode` enum (`slack_errors.py:18-30`), and analytics classifier (`webhook_router.py:70-90`).
- **Routing/handling uses the taxonomy** in three places: HTTP status → exception, `tenacity` retry on `KeycloakConnectionError`, `SlackErrorCode` → user message. The GitHub manager also maps multiple error types to distinct user-facing messages (`github_manager.py:408-427`).
- **Categorical documentation** exists for analytics in `EVENTS.md:25` (six categories), but no single document enumerates the `OpenHandsError` / `AuthenticationError` / `SlackErrorCode` / status-classification families.
- **Extensibility is mixed**: a new `SlackErrorCode` enum value trivially adds a new code; a new `OpenHandsError` subclass works but is not centrally registered; the `_classify_error_type` function in `webhook_router.py:70-90` is a hard-coded `if`-chain and adding a category means editing the function — it is not driven by an enum or table.
- **Fragility**: there are no fewer than 79 `class *Error|Exception` declarations across the repo with at least **four different base classes** (`Exception`, `ValueError`, `HTTPException`, `LookupError`), and dozens of bare `except Exception:` handlers that swallow category information (`enterprise/integrations/github/github_manager.py:431-435`, `openhands/app_server/integrations/provider.py:174-308`).

Pulling these together: the dimension prompt's headline question — "Can you tell from the error type whether to retry, escalate, or stop?" — is answerable *yes* inside narrow subsystems (HTTP retries are governed by `KeycloakConnectionError`; Slack messages are governed by `SlackErrorCode`; analytics classification is governed by `_classify_error_type`) but **no** for the system as a whole, because no central taxonomy unifies those answers.

## Evidence Collected

Every entry MUST include a file path with line numbers. Format: `path/to/file.ts:NN`.

| Area | Evidence | File:Line |
|------|----------|-----------|
| Central HTTP-exception base classes | `OpenHandsError(HTTPException)` with `AuthError`, `PermissionsError`, `SandboxError` subclasses pinned to specific status codes | `openhands/app_server/errors.py:6-43` |
| LLM / settings / session error types | `MissingSettingsError(ValueError)`, `LLMAuthenticationError(ValueError)`, `SessionExpiredError(ValueError)` re-exported in `__all__` | `openhands/app_server/types.py:19-43` |
| Profile storage error types | `ProfileNotFoundError(LookupError)`, `ProfileLimitExceededError(ValueError)`, `ProfileAlreadyExistsError(ValueError)` | `openhands/app_server/settings/llm_profiles.py:70-93` |
| App-conversation export error types | `ConversationExportAlreadyRunning`, `ConversationExportLockUnavailable`, `ConversationExportTooLarge` | `openhands/app_server/app_conversation/app_conversation_service.py:21-30` |
| Git provider / HTTP error mapping | `AuthenticationError`, `RateLimitError`, `ProviderTimeoutError`, `ResourceNotFoundError`, `UnknownException` all `ValueError`s; `handle_http_status_error` switches on `401`/`404`/`429` | `openhands/app_server/integrations/service_types.py:171-198`; `openhands/app_server/integrations/protocols/http_client.py:82-109` |
| HTTP-error mapping unit tests | Per-status-code test cases `test_handle_http_status_error_401/404/429/other` and `test_handle_http_error_with_different_error_types` | `tests/unit/integrations/protocols/test_http_client.py:126-219` |
| Async / concurrency error types | `AsyncException` aggregates multiple exceptions, `RedisLockUnavailable` raised from `try_acquire_redis_lock` | `openhands/app_server/utils/async_utils.py:62-97`; `openhands/app_server/utils/redis_lock.py:21-33` |
| Enterprise auth error hierarchy | `AuthError(Exception)` → `NoCredentialsError`, `EmailNotVerifiedError`, `BearerTokenError`, `CookieError`, `TosNotAcceptedError`, `ExpiredError`, `TokenRefreshError` | `enterprise/server/auth/auth_error.py:1-46` |
| Enterprise auth error tests | Inheritance and instantiation tests confirm `issubclass(NoCredentialsError, AuthError)` and exception chaining | `enterprise/tests/unit/test_auth_error.py:9-60` |
| Slack integration error taxonomy | `SlackErrorCode` enum (`SLACK_ERR_001`…`SLACK_ERR_999`) carried by a single `SlackError`; `_USER_MESSAGES` dispatch table; `get_user_message()` lookup | `enterprise/integrations/slack/slack_errors.py:18-134` |
| Slack integration error routing | Centralized `handle_slack_error` logs and sends user message; outer `except Exception` rewraps as `SlackErrorCode.UNEXPECTED_ERROR` | `enterprise/integrations/slack/slack_manager.py:441-466`, `613-652` |
| Organization domain error hierarchy | `OrgCreationError` → `OrgNameExistsError`, `LiteLLMIntegrationError`, `OrgDatabaseError`; `OrgDeletionError` → `OrgAuthorizationError`, `OrphanedUserError`; flat `OrgNotFoundError`, `OrgMemberNotFoundError`, `RoleNotFoundError`, `InvalidRoleError`, `InsufficientPermissionError`, `CannotModifySelfError`, `LastOwnerError`, `MemberUpdateError`, `GitOrgAlreadyClaimedError` | `enterprise/server/routes/org_models.py:29-148`, `628-636` |
| User-app-settings error hierarchy | `UserAppSettingsError` → `UserNotFoundError`, `UserAppSettingsUpdateError` | `enterprise/server/routes/user_app_settings_models.py:9-...` |
| Invitation error hierarchy | `InvitationError` → `InvitationAlreadyExistsError`, `UserAlreadyMemberError`, `InvitationExpiredError`, `InvitationInvalidError`, `InsufficientPermissionError`, `EmailMismatchError` | `enterprise/server/routes/org_invitation_models.py:11-...` |
| Jira payload error taxonomy | `JiraPayloadParseError(Exception)` plus dataclass result types `JiraPayloadSuccess` / `JiraPayloadSkipped` / `JiraPayloadError` for parse outcomes | `enterprise/integrations/jira/jira_payload.py:55-85` |
| Jira conversation error types | `StartingConvoException(Exception)` and `RepositoryNotFoundError(Exception)` separated as precondition vs. conversation-creation failures | `enterprise/integrations/jira/jira_types.py:53-71` |
| Jira DC error hierarchy | `JiraDcServiceAccountError(ValueError)` → `JiraDcServiceAccountConfigError`, `JiraDcServiceAccountResolutionError`; `JiraDcUserTokenError(Exception)` | `enterprise/integrations/jira_dc/jira_dc_service_account.py:20-...`; `enterprise/integrations/jira_dc/jira_dc_user_token.py:26-...` |
| GitLab webhook control-flow exception | `BreakLoopException(Exception)` raised when rate-limited or conditions unmet | `enterprise/integrations/gitlab/webhook_installation.py:35-38`, `74-78`, `98-103`, `124-125`, `175-177` |
| Rate-limit exception (HTTPException subclass) | `RateLimitException(HTTPException)` carrying `RateLimitResult`; bound to a 429 handler | `enterprise/server/rate_limit.py:109-137` |
| OAuth device-flow error responses | `DeviceTokenErrorResponse` Pydantic model with `error` / `error_description` / `interval`; `_oauth_error` helper maps `invalid_grant`, `slow_down`, `expired_token`, `access_denied`, `authorization_pending`, `server_error` | `enterprise/server/routes/oauth_device.py:47-80`, `132-239` |
| Analytics conversation-error classifier | `_classify_error_type` returns one of `budget_exceeded`, `timeout`, `user_cancelled`, `model_error`, `runtime_error`, `unknown` from raw error-message text | `openhands/app_server/event_callback/webhook_router.py:70-90`, `153-172` |
| Analytics taxonomy documentation | `EVENTS.md` lists the six categories and the `error_type` property on `conversation errored` events | `openhands/analytics/EVENTS.md:25-35` |
| Analytics track_conversation_errored | `track_conversation_errored` records `error_type` and `error_message`; also fires `track_credit_limit_reached` when `error_type == 'budget_exceeded'` | `openhands/analytics/analytics_service.py:255-283`; `openhands/app_server/event_callback/webhook_router.py:156-172` |
| Retry classification by source category | `@retry(retry=retry_if_exception_type(KeycloakConnectionError), stop=stop_after_attempt(3))` on `SaasUserAuth.refresh`; `KeycloakPostError` deliberately excluded | `enterprise/server/auth/saas_user_auth.py:292-298` |
| Token-manager retry classification | Four `tenacity` retry sites gated on `KeycloakConnectionError` with a shared `_before_sleep_callback` for logging | `enterprise/server/auth/token_manager.py:88-89`, `170-173`, `270-272`, `558-560`, `585-...` |
| Route-layer error-type to HTTP-status mapping | `orgs.py` route catches `CannotModifySelfError`/`RoleNotFoundError`/`InvalidRoleError`/`InsufficientPermissionError`/`LastOwnerError`/`MemberUpdateError` and maps each to a specific 400/403/404/500 response | `enterprise/server/routes/orgs.py:1257-1298` |
| GitHub manager multi-error user messaging | Distinct user messages for `MissingSettingsError`, `LLMAuthenticationError`, `(AuthenticationError, ExpiredError, SessionExpiredError)`, and a final `except Exception` | `enterprise/integrations/github/github_manager.py:408-435` |
| Auth middleware error-type dispatch | `SetAuthCookieMiddleware` catches `EmailNotVerifiedError`, `NoCredentialsError`, `AuthError` and renders different JSONResponse + cookie-clearing side effects | `enterprise/server/middleware.py:30-107` |
| HTTPException direct raises in router | `webhook_router.py` raises plain `HTTPException(401, ...)` and `HTTPException(404, ...)` for session-key validation rather than the `AuthError` taxonomy | `openhands/app_server/event_callback/webhook_router.py:244-275`, `299`, `540`, `550`, `558` |
| Provider-handler exception leak | `provider.py:512-513` catches `AuthenticationError` and re-raises as `Exception('Git provider authentication issue when getting remote URL')`, losing the source category | `openhands/app_server/integrations/provider.py:512-513` |
| FastAPI global exception handler | One app-wide `@app.exception_handler(AuthenticationError)` mapping to a 401 JSON response | `openhands/app_server/app.py:63-68` |
| Centralized app-server error type catalogue | `errors.py` contains exactly four classes; the rest of the catalogue is scattered across `types.py`, `service_types.py`, `llm_profiles.py`, integration subpackages, and the enterprise add-on | `openhands/app_server/errors.py:6-43` |
| Test coverage of `AuthError` family | `test_auth_error_inheritance`, `test_*_instantiation`, `test_auth_error_with_cause` confirm subclass and chaining contracts | `enterprise/tests/unit/test_auth_error.py:9-60` |
| Test coverage of HTTP error mapper | `test_handle_http_status_error_{401,404,429,other}` plus `test_handle_http_error_with_different_error_types` exhaustively exercise the `http_client.py` classifier | `tests/unit/integrations/protocols/test_http_client.py:126-219` |
| Test coverage of Slack error dispatch | `TestHandleSlackError`, `TestReceiveMessagePayloadProcessingError`, `TestSendErrorFromPayload`, `TestSendErrorComment` exercise the `SlackErrorCode` → message pathway | `enterprise/tests/unit/test_slack_integration.py` (referenced) |
| Test coverage of error analytics classification | `tests/unit/test_analytics_service.py` covers the `track_conversation_errored` payload | `tests/unit/test_analytics_service.py` |
| Test coverage of remote-sandbox service errors | `tests/unit/app_server/test_remote_sandbox_service.py` includes a `TestErrorHandling` class | `tests/unit/app_server/test_remote_sandbox_service.py` |
| Test coverage of `MaintenanceTaskStatus` ERROR | `tests/unit/test_maintenance_task_runner_standalone.py` asserts `error_type == 'ValueError'` is captured in task `info` | `enterprise/tests/unit/test_maintenance_task_runner_standalone.py:388`, `510`, `556` |
| Resend Keycloak sync error taxonomy | `ResendSyncError` → `ResendAPIError` hierarchy; tests pass `error_type='rate_limit_exceeded'` to verify behavior | `enterprise/sync/resend_keycloak.py`; `enterprise/tests/unit/sync/test_resend_keycloak.py:270-355` |

## Answers to Dimension Questions

1. **Are errors classified by source?**  
   *Partially.* A handful of subsystems classify by source, but the system as a whole does not. The clearest examples:
   - HTTP layer maps status codes to `AuthenticationError` (401), `ResourceNotFoundError` (404), `RateLimitError` (429), `ProviderTimeoutError` (timeouts), `UnknownException` (other) — `openhands/app_server/integrations/protocols/http_client.py:82-109`.
   - LLM authentication has a dedicated `LLMAuthenticationError` — `openhands/app_server/types.py:25`.
   - Session expiry is a distinct type: `SessionExpiredError` — `openhands/app_server/types.py:31`.
   - Analytics classifier buckets error messages into `budget_exceeded` / `model_error` / `runtime_error` / `timeout` / `user_cancelled` / `unknown` — `openhands/app_server/event_callback/webhook_router.py:70-90`.
   - Slack taxonomy: provider (`PROVIDER_TIMEOUT`, `PROVIDER_AUTH_FAILED`, `LLM_AUTH_FAILED`), storage (`REDIS_STORE_FAILED`, `REDIS_RETRIEVE_FAILED`), auth (`SESSION_EXPIRED`, `USER_NOT_AUTHENTICATED`, `MISSING_SLACK_SCOPES`, `MISSING_SETTINGS`), generic (`UNEXPECTED_ERROR`) — `enterprise/integrations/slack/slack_errors.py:18-30`.  
   Counter-evidence: ~70 other exception classes use no source label and live in scattered modules, and `provider.py:512-513` deliberately *erases* the `AuthenticationError` category by re-raising `Exception('Git provider authentication issue when getting remote URL')`.

2. **Is the taxonomy used for handling?**  
   *Yes, in narrow subsystems.* Evidence:
   - `provider.py:304-306` lets `ProviderTimeoutError` propagate while catching every other exception: `except ProviderTimeoutError: raise; except Exception as e: logger.warning(...)`.
   - `saas_user_auth.py:292-298` and `token_manager.py:170-173` use `tenacity.retry_if_exception_type(KeycloakConnectionError)` to retry only transient connection failures, explicitly excluding deterministic `KeycloakPostError`.
   - `slack_manager.py:441-466` uses a single top-level `except SlackError as e:` to route all source-classified errors to `handle_slack_error`, with a fallback `except Exception` rewrapping as `SlackErrorCode.UNEXPECTED_ERROR`.
   - `orgs.py:1257-1298` maps each domain-error type to a specific HTTP status code (`CannotModifySelfError`→403, `RoleNotFoundError`→500, `InvalidRoleError`→400, `LastOwnerError`→400, `InsufficientPermissionError`→403, `ValueError`→400).
   - `github_manager.py:408-435` selects a different Slack user message per `MissingSettingsError` / `LLMAuthenticationError` / `(AuthenticationError, ExpiredError, SessionExpiredError)`.
   - `webhook_router.py:167-172` branches on the analytics error category to fire `track_credit_limit_reached` only when `error_type == 'budget_exceeded'`.  
   The taxonomy is **not** used centrally: `errors.py` defines only four `OpenHandsError` subclasses and the FastAPI app registers exactly one global exception handler for `AuthenticationError` (`openhands/app_server/app.py:63-68`).

3. **Are error categories documented?**  
   *Minimally.* The only dedicated document is `openhands/analytics/EVENTS.md:25-35`, which lists the six analytics categories. Individual error classes carry one-line docstrings (e.g. `errors.py:7-9` "General Error", `slack_errors.py:6-7` "Centralized error handling for Slack integration"). No `docs/` or `ARCHITECTURE.md` describes the overall error model, and the enterprise documentation (`enterprise/doc/architecture/*.md`, `enterprise/doc/design-doc/*.md`) does not mention error categorization at all.

4. **Can new error types be added without breaking existing handling?**  
   *Mixed.* 
   - **Safe**: Adding a new `SlackErrorCode` enum value and a `_USER_MESSAGES` entry is fully additive and covered by `tests/unit/test_slack_integration.py`. Adding a new `AuthError` subclass also works because of the `issubclass` tests in `tests/unit/test_auth_error.py:9-14`.
   - **Fragile**: Adding a new `OpenHandsError` subclass has no registration step — callers that want to handle it must already know about it or import the class directly. The `_classify_error_type` function in `webhook_router.py:70-90` is a hard-coded `if`-chain on substring matches; adding a new analytics category requires editing this function (no enum or table to extend).
   - **Inconsistent**: A new HTTP-status mapping requires editing both `http_client.py:82-99` and adding a test in `tests/unit/integrations/protocols/test_http_client.py:126-219`. The mapping is centralized, so the change is local, but the new `ValueError` class has to be defined in `service_types.py:171-198` and re-exported.

## Architectural Decisions

- **Multiple parallel hierarchies rather than one base class.** The codebase has at least five roots: `OpenHandsError(HTTPException)` (`errors.py:6`), `AuthError(Exception)` (`enterprise/server/auth/auth_error.py:1`), `ValueError`-based domain errors (`types.py:19-43`, `service_types.py:171-198`), `SlackError(Exception)` plus `SlackErrorCode` enum (`slack_errors.py:18-33`), and ad-hoc `Exception` subclasses (e.g. `ConversationExportAlreadyRunning`, `RedisLockUnavailable`, `BreakLoopException`). There is no single `BaseOpenHandsError` to catch.
- **HTTP status codes are the de facto taxonomy for provider errors.** `http_client.py:82-99` centralizes the mapping so that downstream callers see typed errors. This is the strongest source-based classification in the system.
- **Slack integration inverts the pattern with a `SlackErrorCode` enum.** Instead of multiple exception classes, Slack uses a single `SlackError(Exception)` carrying an enum member. The enum is closed (it ends at `SLACK_ERR_999` for `UNEXPECTED_ERROR`) and lookup of user messages falls through to a generic default if the code is missing — `slack_errors.py:123-134`.
- **Error categories for analytics are derived, not raised.** `_classify_error_type` does best-effort string matching on the message text rather than introspecting exception types — `webhook_router.py:70-90`. This means a new exception class with a novel message format silently falls into `runtime_error` unless someone updates the classifier.
- **Retry classification is explicit but local.** `tenacity.retry_if_exception_type(KeycloakConnectionError)` appears at five call sites in `saas_user_auth.py:292-298` and `token_manager.py:170-585`. There is no central "is this retryable?" predicate; each module hand-codes its own retry policy.
- **Auth uses two parallel hierarchies.** The open-source package defines `AuthError`/`PermissionsError` inheriting from `HTTPException` (`errors.py:18-39`), while the enterprise add-on defines a separate `AuthError(Exception)` with seven subclasses (`enterprise/server/auth/auth_error.py:1-46`). Middleware in `enterprise/server/middleware.py:30-107` is the only place where the two are unified, and even there it dispatches to different JSONResponse shapes per leaf type.
- **Routes, not a central handler, do most error→HTTP mapping.** `orgs.py:1257-1298` is the canonical example: a 40-line `try/except` block catching each domain error and converting to a specific status. There is no FastAPI `@app.exception_handler(OrgNotFoundError)` registered.
- **Swallowing exceptions is common but not category-preserving.** `provider.py:174-308` uses `except Exception as e: logger.warning(...)` extensively, which means a transient `ProviderTimeoutError` is logged at warning level and lost; only the re-raise on `except ProviderTimeoutError:` at line 304 preserves the category. `github_manager.py:431-435` uses a final `except Exception: ... 'unexpected error starting the job'`, which deliberately downgrades the error type to a user-visible "unexpected error" message.

## Notable Patterns

- **Centralized error→message dispatch table** (`_USER_MESSAGES` keyed by `SlackErrorCode`): `slack_errors.py:69-110`. Adding a new code requires only updating the table.
- **HTTP-status to typed-error switch** inside `HTTPClient.handle_http_status_error` and `handle_http_error`: `http_client.py:82-109`. The annotated return type `AuthenticationError | RateLimitError | ResourceNotFoundError | UnknownException` documents the classification at the type level.
- **Dataclass result-type union for parser outcomes** (`JiraPayloadSuccess | JiraPayloadSkipped | JiraPayloadError`): `jira_payload.py:78-85`. This is a non-exception variant of the taxonomy pattern that is preferable for parser-level decisions.
- **`Enum` + `IntEnum` for status fields**: `SandboxStatus`, `EventCallbackStatus`, `EventCallbackResultStatus`, `AppConversationStartTaskStatus`, `ConversationTrigger`, `StorageProvider`, `ProviderType`, `WebhookStatus`, `MaintenanceTaskStatus`, `JiraEventType`, `SlackErrorCode`. The codebase uses enums liberally but not for error *type* in the open-source package.
- **Per-error user-facing message strings** keyed by exception type: `github_manager.py:408-427` and `slack_errors.py:_USER_MESSAGES`. This pattern means the user-message is part of the error contract.
- **`tenacity`-gated retry on `KeycloakConnectionError`**: explicit choice to retry transient infra errors but not deterministic `KeycloakPostError` — `saas_user_auth.py:292-298`, `token_manager.py:170-173`.
- **Test-driven taxonomy**: per-status-code tests in `tests/unit/integrations/protocols/test_http_client.py:126-219`; per-subclass tests in `tests/unit/test_auth_error.py:9-60`; per-SlackErrorCode tests in `tests/unit/test_slack_integration.py`. The taxonomy is covered by tests but the test names reveal that the team does not have a single `test_error_taxonomy.py`.

## Tradeoffs

- **Centralization vs. locality**: the HTTP-status mapping is centralized (`http_client.py:82-99`) which is good for consistency, but most other error categories are defined next to the code that raises them, which makes it hard to enumerate the full taxonomy from one place.
- **Single exception class + enum vs. exception class hierarchy**: Slack chose the former (`SlackError` + `SlackErrorCode`); the rest of the codebase chose the latter. The Slack approach is more extensible and serializable but couples the message lookup table to a single module. The hierarchy approach is more "Pythonic" but makes `try/except` for "all source-X errors" require either catching the base class or listing every leaf.
- **Deriving error category from message text vs. exception type**: `_classify_error_type` (`webhook_router.py:70-90`) parses the message string. This is more robust when the original exception type is unavailable (e.g. across process boundaries) but is brittle to message changes and conflates semantically different errors that share keywords.
- **Bare `except Exception` for resilience vs. category preservation**: many call sites (e.g. `provider.py:174-308`, `slack_manager.py:458-466`, `github_manager.py:431-435`) catch everything and log warning. This prevents crashes but throws away the source category, making post-hoc analysis difficult.
- **HTTPException base vs. Exception base**: `OpenHandsError` extends `HTTPException` so it carries a status code, but this couples the error class to a transport-layer concern. The other hierarchies (`AuthError`, `SlackError`, `RedisLockUnavailable`) are transport-agnostic, requiring explicit HTTP mapping at the route layer.
- **No global `OpenHandsError` base**: by extending multiple bases (HTTPException, ValueError, Exception, LookupError), the system pays a small price in expressiveness (you can catch `ValueError` to mean "bad input" and still get provider errors) but loses the ability to add app-wide behavior to all errors (e.g. automatic Sentry breadcrumbs, automatic request-id tagging).

## Failure Modes / Edge Cases

- **Category erasure in `provider.py:512-513`**: catching `AuthenticationError` and re-raising as bare `Exception` means downstream observability loses the source category. The original `AuthenticationError` is also not chained with `from`, so its traceback is lost.
- **Hard-coded substring matching in `_classify_error_type`**: `webhook_router.py:78-90` — a typo in an error message ("bdgt exceeded" instead of "budget exceeded") silently changes the analytics bucket from `budget_exceeded` to `runtime_error`. The fallback case at line 90 swallows any new error type as `runtime_error`.
- **No enum for the analytics error category**: the six categories are passed as raw strings, and a typo at a call site (e.g. `error_type='budget_exceed'`) would not match the comparison in `webhook_router.py:167`.
- **Sliding taxonomy in HTTP-client mapper**: adding a new HTTP status (e.g. 503 Service Unavailable) requires editing `http_client.py:82-99` *and* adding a new exception class to `service_types.py:171-198`. There is no exhaustiveness check; an unhandled status returns `UnknownException` and downstream code only gets the error message.
- **Webhook router raises plain `HTTPException` instead of the centralized `AuthError`**: `webhook_router.py:244-275`, `299`, `540`, `550`, `558` mix `HTTPException(401, ...)` with the `AuthError` taxonomy, meaning the app-wide exception handler in `app.py:63-68` does *not* trigger for the webhook router's auth failures. The two patterns coexist.
- **Two `AuthError` classes**: `openhands/app_server/errors.py:18` (`AuthError(OpenHandsError)`) and `enterprise/server/auth/auth_error.py:1` (`AuthError(Exception)`) have the same name but unrelated hierarchies. A `try/except AuthError` that works in enterprise code will not catch the open-source one (or vice versa), and vice versa.
- **Catch-all `except Exception` in `set_response_cookie` and analytics**: `enterprise/server/middleware.py:30-107` and `analytics_service.py:90-97` deliberately swallow all exceptions. The trade-off is reliability, but the cost is that a bug in error handling logic can hide behind a generic warning log.
- **Tenacity retry is not always safe**: `saas_user_auth.py:292-298` retries `KeycloakConnectionError` 3 times, but a network black hole would block the request for ~3 seconds before failing. There is no circuit breaker layered on top.
- **`Unhandled webhook event type` is not an error**: `JiraPayloadParser.parse` returns `JiraPayloadSkipped(...)` for unhandled events (`jira_payload.py:128-129`), which is a *good* pattern — error taxonomy should not over-classify intentional skips. The `SlackErrorCode.UNEXPECTED_ERROR = 'SLACK_ERR_999'` follows the same spirit by reserving a high-numbered code for the truly unknown case.
- **Rate-limit `except` swallows Keycloak errors silently in places**: `token_manager.py:573-576` catches `KeycloakConnectionError` and logs without re-raising, which means a transient blip is invisible to the caller. Compare with `token_manager.py:170-173` where the same exception is wrapped in a retry.

## Future Considerations

- **Unify the four base classes behind a single `OpenHandsError` base** that is *not* an `HTTPException`. Wrap transport-level concerns in a separate "HTTP-mapped error" decorator or in a single FastAPI exception handler. This would enable a single `try/except OpenHandsError` for "all domain errors" and let each leaf carry its source category as an attribute.
- **Replace the substring-based `_classify_error_type` with an enum + table**: e.g. `class AnalyticsErrorCategory(str, Enum): BUDGET_EXCEEDED = 'budget_exceeded'` and a mapping from keyword regex to category, defined once in `openhands/analytics/`. The webhook router would then dispatch on enum membership, making typos a syntax error.
- **Promote the `SlackError`+`SlackErrorCode` pattern to the system level**: every other subsystem that catches multiple distinct errors (e.g. `github_manager.py:408-435`, `orgs.py:1257-1298`) could use a single `SourceError` carrying a `SourceErrorCode` enum and a per-code dispatch table for user messages, HTTP status, retry policy, and analytics category. This would make "retry vs. escalate vs. stop" a lookup instead of an `if`-chain.
- **Add an `is_retryable` or `retry_policy` attribute to every domain exception class** so `tenacity.retry_if_exception_type` can be replaced with `tenacity.retry_if(lambda e: e.retry_policy == RetryPolicy.TRANSIENT)`. This would put retry decisions on the error type itself rather than scattered across call sites.
- **Wire up centralized FastAPI exception handlers for the new taxonomy**: register `@app.exception_handler(MissingSettingsError)`, `@app.exception_handler(LLMAuthenticationError)`, etc., rather than relying on per-route `try/except` blocks (`orgs.py:1257-1298`). This is the path that `app.py:63-68` already takes for `AuthenticationError`.
- **Add a `docs/ERROR_TAXONOMY.md`** that enumerates the full set of error classes, their source category, their retry policy, and the HTTP response they map to. The six analytics categories are documented in `EVENTS.md:25-35`; the same treatment for the wider taxonomy would prevent drift.
- **Fix the `provider.py:512-513` category erasure** by re-raising as the same type (`raise AuthenticationError(...) from e`) instead of `raise Exception('Git provider authentication issue when getting remote URL')`. The current code loses the source category at exactly the place where classification matters most.
- **Consolidate the two `AuthError` classes**: either rename the enterprise one to `SaasAuthError` or have it import the open-source one. Currently a developer reading `except AuthError` cannot tell which one is meant without a context switch.
- **Move exception-handler registration into a `setup_error_handlers(app)` function** in `app.py` so the policy is discoverable in one place, mirroring `setup_rate_limit_handler` in `enterprise/server/rate_limit.py:25-29`.

## Questions / Gaps

- The `OpenHandsError` hierarchy in `errors.py:6-43` is conspicuously small (4 classes) compared to the actual catalogue (~70+). Is this intentional, with the rest of the catalogue living in `types.py` and the integration packages? Or is `errors.py` meant to grow? No design doc addresses this.
- The `SandboxError` class is the only `OpenHandsError` subclass used widely (10 call sites in `live_status_app_conversation_service.py` and the sandbox services), while `AuthError` and `PermissionsError` are used sparingly. The asymmetry suggests the `OpenHandsError` base is a stub rather than a designed hierarchy.
- `_classify_error_type` (`webhook_router.py:70-90`) is documented as "best-effort string matching per CONTEXT.md decision" in a comment, but no `CONTEXT.md` was found in the repo. The decision is therefore unverifiable from the source.
- There is no central `OpenHandsError` registered as a FastAPI exception handler. A caller that catches `OpenHandsError` cannot rely on app-wide behavior to map it to a response.
- The taxonomy treats "validation" and "user" categories as the same (e.g. `OrgNameExistsError` could be either "bad input" or "domain conflict"). The codebase does not draw a clean line between them, which is acceptable but worth surfacing.
- The enterprise `OrgCreationError`/`OrgDeletionError` bases and the enterprise `InvitationError` base have different leaf types under the same name (e.g. `InsufficientPermissionError` is defined twice: `org_models.py:121` and `org_invitation_models.py`). The duplicates are not unified.
- The `MissingSettingsError`, `LLMAuthenticationError`, `SessionExpiredError` types are defined in `openhands/app_server/types.py:19-43` but not in `errors.py`, so a `try/except OpenHandsError` will not catch them. The split between `errors.py` and `types.py` is undocumented.
- The `RateLimitException` in `enterprise/server/rate_limit.py:109-137` is a `HTTPException(429)` with a custom `result` attribute, but the open-source `RateLimitError` is a `ValueError`. The two are not related by inheritance even though they share a name.
- `webhook_router.py:299` raises `AuthError()` from the open-source `errors.py` after a same-owner mismatch, but other auth checks in the same file raise plain `HTTPException(401, ...)`. The two patterns coexist with no policy statement on when to use which.
- The taxonomy of `KeycloakError` (referenced in `token_manager.py:13`, `142`, `804`, `945`, `976`, `993`) is from an external library (`python-keycloak`); how it integrates with the internal `AuthError`/`MissingSettingsError` taxonomy is not documented.

---

Generated by `13.01-error-taxonomy` against `openhands`.
