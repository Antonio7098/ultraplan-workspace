# Review Sprint Protocol

> Role: post-implementation sprint review orchestrator.
> Canonical output: `projects/<project>/sprints/<sprint>/review.md`.
> Contract and protocol paths are resolved dynamically from project and sprint indexes.

This protocol replaces the older manually coordinated review. It reviews the implemented sprint scope against selected contracts, technical-handbook guidance, sprint decisions, plan tasks, and verification evidence. It does not redesign the sprint, fix product code, mutate tests, change governed inputs, or mutate Git state.

## Purpose And Order

Run review after `execute` and before deep smoke:

```text
execute -> review -> smoke
```

Review runs first because it is cheaper and more deterministic than live smoke. Applicable blocker or high-severity findings stop the default flow before smoke.

## 1. Discover Context

Required identifiers:

```text
Project: <project>
Sprint: <sprint-slug>
Project index: projects/<project>/project-index.md
Sprint directory: projects/<project>/sprints/<sprint-slug>/
Implementation directory: resolved from project-index.md
```

Read `project-index.md`, `requirements.md`, `sprint-index.md`, `technical-handbook.md`, selected `reasoning/*.md`, `reasoning.md`, `plan.md`, `execute.md` when present, `.run-state.json` when present, and `flow-state.json`.

Resolve every selected contract and required review protocol through the project index. Do not maintain a hardcoded contract-name map. Fail preflight for missing, duplicate, unreadable, unknown, or escaping catalog paths.

## 2. Freeze Review Scope

Before reviewer execution, determine and display:

- target implementation root and revision/fingerprint;
- changed paths or explicit full-target scope;
- governed sprint input hashes;
- selected contracts and protocols;
- plan tasks and execute evidence;
- approved deterministic verification commands;
- reviewer count, model source, concurrency, timeout, and write boundary.

The fingerprint covers governed inputs, selected contract/protocol content, implementation scope/identity, and execute evidence. A later mismatch makes `review.md` stale. Reviewers are read-only; the only product write is the final atomic replacement of `review.md`.

## 3. Run Deterministic Checks

Check decision conformance, plan/execute evidence, approved verification commands, reviewer coverage, evidence containment, citation validity, deviations/deferred scope, and prohibited mutations. Never execute arbitrary shell text extracted from generated Markdown; approved commands run as explicit executable/argument arrays.

## 4. Run Independent Reviewers

UltraPlan owns bounded fan-out and does not assume a reviewer runtime can spawn subagents. Run one structured reviewer per selected contract and one for the technical handbook.

Each reviewer classifies applicability as `direct`, `partial`, `not_triggered`, or `explicitly_deferred`; reports conformance, severity, contained evidence, issue, and recommended action; and returns validated structured data without writing files. Missing/failed mandatory reviewers remain failures and are never silently dropped.

## 5. Synthesize `review.md`

Collect all available results even when one fails. Write sections for Review Context, Input Fingerprint And Scope, Decision Conformance, Plan Execution, Verification Evidence, Contract Conformance, Technical Handbook Conformance, Applicability And Deferred Scope, Findings, Deviations, and Final Assessment.

Product code computes the verdict:

- `pass`: all mandatory work completed; no applicable blocker/high finding;
- `pass_with_findings`: all mandatory work completed; only medium/low/info findings;
- `fail`: blocker/high finding, failed required verification, invalid evidence, or missing mandatory reviewer;
- `blocked`: required inputs, scope, runtime/model, or verification environment unavailable.

An LLM may summarize but cannot choose or override the verdict.

## 6. Validate And Replace

Validate required sections, selected-contract/handbook coverage, failed-reviewer disclosure, contained citations, severity/verdict consistency, fingerprint/scope, redaction, and absence of smoke claims. Write atomically. A failed/cancelled run must not corrupt the last complete artifact. A successful explicitly started review may replace an older manual or generated `review.md`.

## 7. TUI Parity

The TUI exposes the same typed review use case as the CLI: readiness, dry-run, prompt preview, confirmation, scope/contracts, progress, cancellation, findings, verdict, rerun, recovery, and `review.md` preview. It must not call CLI handlers or parse terminal output.

## Stop Conditions

Return `blocked` before execution for missing/invalid inputs, unresolved catalog entries, ambiguous implementation scope, or unavailable review runtime/model. Return `fail`, with a synthesized review where possible, for failed verification, invalid/missing mandatory reviewer evidence, or applicable blocker/high findings.
