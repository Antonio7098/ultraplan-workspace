# LLM Evaluation, Cost, And Safety Contract

## Purpose

This contract extends LLM runtime contracts with evaluation, cost, and safety requirements.

It governs:

- prompt/model/tool evaluation
- behavioural regression testing
- safety checks
- cost and latency tracking
- dataset and eval data handling
- model/provider/runtime change review

## Scope

This contract applies when a change:

- changes prompts, model selection, tool use, structured output, retrieval, memory, or agent behaviour
- adds AI features exposed to users or operators
- changes model/provider/runtime fallback, retries, temperature, context construction, or output validation
- creates or uses eval datasets

## Requirement Index

| ID | Title | Applies To | Severity If Violated |
| --- | --- | --- | --- |
| EVAL-SCOPE-001 | AI behaviour changes must define evaluation scope | prompt/model/tool changes | High |
| EVAL-REG-001 | Critical AI behaviours need regression examples | user-facing AI paths | High |
| EVAL-DATA-001 | Eval data must be classified and safe to use | eval datasets | High |
| EVAL-SAFETY-001 | Unsafe tool/model behaviour must be guarded | tool/agent paths | Blocker |
| EVAL-COST-001 | Cost and latency drivers must be observable | runtime/provider-backed paths | Medium |
| EVAL-MODEL-001 | Model/provider/runtime changes must be reviewed as behavioural changes | model/provider/runtime switches | High |
| EVAL-HUMAN-001 | Human review must exist for high-impact uncertain outputs | high-risk AI decisions | High |

## Core Principles

1. Prompt changes are code changes.
2. Model changes are behaviour changes.
3. AI regressions need examples, not vibes.
4. Eval data is user data unless proven otherwise.
5. Tool-using agents need explicit safety boundaries.
6. Cost and latency are runtime behaviours.
7. High-impact AI output needs appropriate human review or guardrails.

## Requirements

### EVAL-SCOPE-001: AI Behaviour Changes Must Define Evaluation Scope

**Rule**
Any meaningful AI behaviour change must define what behaviour should improve, remain stable, or be explicitly allowed to change.

**Required**
- state changed prompt/model/tool/retrieval/runtime behaviour
- identify expected improvement or reason for change
- identify behaviours that must not regress
- choose appropriate eval/smoke/manual review level

**Forbidden**
- editing prompts/models/tools with no expected behaviour statement
- treating green unit tests as proof of AI behaviour quality

**Evidence**
- reasoning note identifies eval scope and acceptance signals

### EVAL-REG-001: Critical AI Behaviours Need Regression Examples

**Rule**
Important AI behaviours must be protected with representative examples or scenario tests.

**Required**
- maintain small high-signal regression sets for critical paths
- test structured-output validity and semantic expectations where possible
- include failure/edge examples, not just ideal prompts

**Forbidden**
- relying only on ad hoc manual chat trials for critical AI behaviour
- asserting exact wording where semantic behaviour is what matters

**Evidence**
- eval/smoke results are attached to prompt/model/runtime changes where relevant

### EVAL-DATA-001: Eval Data Must Be Classified And Safe To Use

**Rule**
Evaluation datasets, traces, prompts, completions, and user examples must be classified and handled according to data/privacy rules.

**Required**
- classify eval data as synthetic, anonymized, internal, user-content, PII, or sensitive
- avoid production user content unless policy allows it
- redact/anonymize where possible
- define retention and access controls for eval artifacts

**Forbidden**
- copying raw production user content into evals without approval/policy
- committing sensitive eval data to public or broad repos

**Evidence**
- eval datasets identify source, classification, and retention

### EVAL-SAFETY-001: Unsafe Tool/Model Behaviour Must Be Guarded

**Rule**
Agents and AI features that can perform actions, access tools, generate instructions, or affect users must have safety boundaries appropriate to risk.

**Required**
- define allowed tools/actions and forbidden actions
- require approval for destructive, expensive, external, or sensitive actions where needed
- validate tool arguments before execution
- guard against prompt injection and untrusted instructions in retrieved/user content
- separate user instructions from system/developer/tool policy

**Forbidden**
- allowing model output to call arbitrary tools or shell/network/file operations without policy
- treating retrieved/user content as trusted instructions
- executing destructive actions without confirmation/authorization

**Evidence**
- tests or review notes cover unsafe tool/action attempts

### EVAL-COST-001: Cost And Latency Drivers Must Be Observable

**Rule**
Runtime/provider-backed AI paths must expose enough metadata to understand cost and latency.

**Required**
- capture runtime, provider, model, prompt version, token counts where available, retry count, duration, tool count, and fallback/runtime/provider selection where relevant
- define max context, max attempts, timeout, and concurrency limits
- identify unexpectedly expensive paths

**Forbidden**
- hidden retries/fallbacks that multiply cost silently
- unbounded context assembly or tool loops

**Evidence**
- logs/events/metrics include cost/latency metadata for important AI paths

### EVAL-MODEL-001: Model/Provider/Runtime Changes Must Be Reviewed As Behavioural Changes

**Rule**
Changing model, provider, runtime, fallback, temperature, context strategy, or structured-output mode must be treated as behaviourally significant.

**Required**
- document expected impact on quality, cost, latency, reliability, and output shape
- run eval/smoke coverage appropriate to path risk
- preserve rollback option where practical

**Forbidden**
- silently changing default model/provider/runtime in production
- assuming provider-equivalent output quality without verification

**Evidence**
- model/provider/runtime changes include eval/smoke and release notes where relevant

### EVAL-HUMAN-001: Human Review Must Exist For High-Impact Uncertain Outputs

**Rule**
AI outputs that materially affect users, money, legal/compliance status, safety, or external actions must have appropriate human review, confidence thresholds, or guardrails.

**Required**
- define which outputs can auto-execute and which require review
- expose uncertainty, citations/evidence, validation status, or rationale where useful
- prevent unsupported claims from becoming authoritative decisions

**Forbidden**
- fully automated high-impact decisions with no review, appeal, validation, or audit trail
- presenting unverified generated content as authoritative fact in high-risk contexts

**Evidence**
- high-impact flows show review/approval/validation/audit path

## Review Rejection Criteria

Reject an AI behaviour change if it:

- edits prompts/models/tools with no evaluation scope
- changes critical AI behaviour with no regression examples or smoke plan
- uses raw user content in evals without classification/policy
- lets model output execute unsafe tools/actions without guardrails
- hides model/provider/runtime cost or latency drivers
- silently changes model/provider/runtime defaults
- automates high-impact decisions without review/validation/audit path

---
