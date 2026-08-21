---
name: ultraplan-reconcile-review-smoke
description: Manually run the UltraPlan reconcile stage when given a project sprint path or project/sprint references. Use only when the user explicitly invokes $ultraplan-reconcile-review-smoke or directly asks to run this exact UltraPlan stage; do not invoke implicitly.
---

# UltraPlan Review And Smoke Reconciliation

Run this stage interactively while preserving UltraPlan's governed artifact chain.

## Operating contract

1. Treat a supplied sprint path as UltraPlan stage input, not as a Git target. For an input such as `projects/<project>/sprints/<sprint>/` or `.ultra/projects/<project>/sprints/<sprint>/`, find the workspace root, derive `<project>` and `<sprint>` from the path, and read the matching `project-index.md`. The sprint directory contains governed stage artifacts; when implementation access is required, resolve its repository from `Target Implementation Directory`, falling back to `Repository` only when the target field is absent. Resolve relative repository paths against the workspace root and verify the result before using it. Do not search nested source repositories for a similarly named skill, and do not ask what target to use merely because the supplied input is a directory.
2. If no sprint path was supplied, locate the workspace root and resolve the project and sprint from explicit references and the current location. Ask only when the project index is missing, a required implementation target cannot be resolved, or more than one project/sprint remains possible.
3. Run all UltraPlan commands from the resolved workspace root. Run `ultraplan project <project> status` and `ultraplan sprint <project> <sprint> status --json`. Treat files, the project index, and fresh CLI status as authoritative. Do not hand-edit flow-state JSON except in the explicitly authorized manual review or smoke reconciliation branches below, where every fingerprint, digest, verdict, timestamp, evidence identity, and completion identity must be updated coherently and immediately verified through validation and sprint status.
4. Check these prerequisites:

- an existing review.md or smoke.md result
- the sprint planning artifacts
- a resolvable target implementation directory

5. Inspect the prerequisite artifacts and report missing or internally inconsistent evidence. A stale review fingerprint is context for reconciliation, not a prerequisite failure: note the stored and current fingerprints in the proposed analysis, continue classifying the existing review findings against the current implementation, and do not ask to rerun review or stop reconciliation merely because they differ. Do not use `validate execute` as a reconciliation prerequisite after execution: that validator expects unchecked tasks for a new execute attempt, so a fully checked completed plan is expected to fail it; assess completed execution from execute evidence, run state, and the review dry-run instead. Ask before filling genuinely missing prerequisite artifacts or repairing inconsistent durable state.
6. If the target is already complete and valid, summarize that state and ask before regenerating or materially changing it.
7. If the user explicitly asks for a proposal, analysis, or discussion only, inspect all relevant evidence and return that without writing artifacts or advancing state. Otherwise, do the stage now; do not stop at a proposal.
8. The invoking agent owns the actual stage work. Except for the review stage, do not call an UltraPlan stage, flow, execute, verify, or smoke command to have the CLI or another runtime execute or complete the stage. CLI commands remain appropriate for discovery, effective-prompt resolution, dry-run previews, status inspection, validation, and post-write reconciliation. Review is the deliberate exception: invoke its governed CLI command because UltraPlan owns reviewer subagent fan-out, aggregation, and review state.
9. Use the embedded canonical prompt below together with the current CLI status and governed command preview.
10. Complete the stage-specific workflow, preserve unrelated user edits, and do not cross declared mutation boundaries.
11. Run the validation commands supported for the affected artifact, then use review/smoke dry runs and sprint status as the cross-stage reconciliation checks.
12. Run `ultraplan sprint <project> <sprint> status --json` after writes or governed review execution so flow-state, freshness, and artifact status are reconciled. Re-run project status if the project or sprint index changed.
13. Inspect downstream artifacts for references made stale by this change. Update directly coupled indexes/references when safe; otherwise report the exact dependent stage that must be revisited. Never delete or silently rewrite downstream decisions.
14. Finish with the artifact/result paths, validation outcome, state transition, and any remaining blocker or dependent stage.

## Stage workflow

Use this workflow after a review or smoke result needs human assessment and remediation. It is not a substitute for initially running the governed review or smoke process.

1. Establish the frozen and current evidence:
   - Read the complete review.md or smoke.md, flow-state.json, .run-state.json, requirements.md, sprint-index.md, reasoning.md, plan.md, and execute.md.
   - Resolve the implementation target from project-index.md and inspect its current status and diff. Preserve unrelated user changes.
   - Run a review dry-run when reviewing review evidence to obtain the current governed input fingerprint. Record whether the existing result is current or stale.
2. Classify every reported finding using implementation evidence:
   - genuine sprint defect: violates an acceptance criterion, explicit plan decision, mutation/security boundary, or claimed verification property;
   - genuine platform follow-up: real but pre-existing or broader than the sprint, so it must not block this sprint unless the selected contract explicitly made it a release gate;
   - superseded/already fixed;
   - unsupported or scope-expanding: speculative, contradicted by the plan, or explicitly deferred/non-goal.
   Cite concrete files and tests. Do not accept severity labels without reproducing the behavior or tracing the relevant code path.
3. If the user requested fixes, implement only genuine authorized defects in the target repository. Add focused regression tests, then run verification proportional to risk, including the sprint's required full test/race/vet/build/diff gates when closure is claimed.
4. Re-run the governed review normally when practical so UltraPlan owns aggregation, fingerprints, artifact publication, and state. If the generated verdict is still demonstrably wrong and the user explicitly authorizes manual reconciliation:
   - retain the automated result as superseded history;
   - add a concise manual findings and reconciliation section to review.md;
   - use only supported verdict values and the current dry-run fingerprint;
   - recompute the review artifact SHA-256;
   - update review status, verdict, fingerprint, artifact digest, last-complete identity, timestamps, and diagnostic provenance coherently and atomically;
   - never invent reviewer coverage, test evidence, or a passing command result;
   - run sprint status immediately and require it to report the review as completed, fresh, and digest-consistent.
5. Reconcile obsolete next-stage attempts only when their recorded blocker has actually disappeared. Preserve historical successful evidence. A prior failed smoke attempt caused solely by a now-resolved review gate may be cleared from current smoke state; do not erase a real harness, test, cleanup, or mutation failure.
6. Discover and inspect the smoke harness before claiming readiness:
   - Start with sprint smoke --dry-run --json from the workspace root. Read ready, verdict, review_verdict, review_fingerprint, scope, prerequisites, diagnostics, and next_action.
   - Resolve the harness through the project smoke configuration and manifest. The normal sibling checkout is derived from the manifest/target relationship and is often ../ultraplan-go-smoke; do not assume that path without checking discovery output and ultraplan-smoke.json.
   - Inspect the protocol-v1 manifest and discovery implementation, especially executable/cwd, commands, evidence roots, prerequisites, suites, tests, sprintMappings, requiredCoverage, complete, and notApplicable.
   - A complete mapping must reference non-empty tests whose declared coverage satisfies requiredCoverage. Do not create an empty complete mapping to bypass smoke.
   - Use notApplicable only when the sprint plan or requirements explicitly defer or exclude external/live-runtime smoke. It must be represented as a non-contradictory mapping (notApplicable true and complete false), with a truthful rationale naming the owning later sprint or gate.
   - If a harness change is needed, edit only declared authoring paths, run the harness build/tests, then repeat the UltraPlan smoke dry-run.
7. The next stage is ready only when the dry-run reports ready true, the current review verdict/fingerprint match durable state, and the selection is either runnable with satisfied prerequisites or truthfully not_applicable. A dry-run that is blocked, stale, diagnostic-only unexpectedly, or missing coverage is not ready even if review passed.
8. When current smoke execution is passing and the user requested publication or reconciliation, update the sprint-root smoke.md and smoke flow state as one coherent result:
   - Preserve the prior generated summary as superseded history when its findings or verdict are being replaced; never erase real failed-run, cleanup, mutation-safety, or open-issue evidence.
   - Populate every required smoke.md section from the actual selected scope, sanitized invocation, prerequisites, run counts, external evidence identities and hashes, findings/issues, cleanup and mutation checks, supported verdict, and one explicit next action. Keep raw streams and secrets in the external evidence roots.
   - Use pass only when all selected required tests passed, evidence identity is complete, cleanup is certain, no prohibited mutation occurred, and the current non-overridden review gate permits promotion. Use pass_with_open_issues only for truthful non-blocking open issues; do not upgrade a diagnostic override, failed, blocked, or incomplete run.
   - Recompute the smoke.md SHA-256 and the input fingerprint from the exact durable evidence identities. Reconcile status completed, verdict, artifact path and digest, smoke fingerprint, input fingerprint, current review fingerprint, run/author/evidence identities, timestamps, issues, override facts, diagnostics, active-attempt state, and last-complete identity coherently and atomically. Never invent missing run or evidence fields.
   - Run validate smoke and sprint status --json immediately. Require smoke to be completed, fresh, digest-consistent, evidence-fingerprint-consistent, and tied to the current review; require the overall sprint flow to expose the resulting terminal assessment. If supported reconciliation cannot establish all of those facts, leave the truthful prior state intact, report reconciliation as the blocker, and do not claim that smoke or the sprint is passing.
9. Finish by reporting: reconciled review.md, smoke.md, flow-state.json, and external evidence paths; current review fingerprint, smoke input fingerprint, and artifact digest; fixes and verification performed; smoke harness path and mapping decision; final smoke verdict and freshness; overall sprint-flow assessment; and the exact remaining next action.

## Canonical stage prompt

The current resolved CLI prompt wins over this embedded baseline.

# Review And Smoke Reconciliation

Assess a generated UltraPlan review or smoke result against its governed sprint scope and the actual implementation. Fix genuine in-scope defects, preserve unrelated work, reconcile durable evidence truthfully, and prove that the next verification stage is ready or report its exact blocker.
