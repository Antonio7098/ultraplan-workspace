---
name: ultraplan-review
description: Manually run the UltraPlan review stage when given a project sprint path or project/sprint references. Use only when the user explicitly invokes $ultraplan-review or directly asks to run this exact UltraPlan stage; do not invoke implicitly.
---

# UltraPlan Review

Run this stage interactively while preserving UltraPlan's governed artifact chain.

## Operating contract

1. Treat a supplied sprint path as UltraPlan stage input, not as a Git target. For an input such as `projects/<project>/sprints/<sprint>/` or `.ultra/projects/<project>/sprints/<sprint>/`, find the workspace root, derive `<project>` and `<sprint>` from the path, and read the matching `project-index.md`. The sprint directory contains governed stage artifacts; when implementation access is required, resolve its repository from `Target Implementation Directory`, falling back to `Repository` only when the target field is absent. Resolve relative repository paths against the workspace root and verify the result before using it. Do not search nested source repositories for a similarly named skill, and do not ask what target to use merely because the supplied input is a directory.
2. If no sprint path was supplied, locate the workspace root and resolve the project and sprint from explicit references and the current location. Ask only when the project index is missing, a required implementation target cannot be resolved, or more than one project/sprint remains possible.
3. Run all UltraPlan commands from the resolved workspace root. Run `ultraplan project <project> status` and `ultraplan sprint <project> <sprint> status --json`. Treat files, the project index, and fresh CLI status as authoritative; never hand-edit flow-state JSON.
4. Check these prerequisites:

- completed execute stage
- current planning artifacts and target implementation

5. Validate every prerequisite that has a sprint validation command. If anything is missing, invalid, stale, or internally inconsistent, show the exact gaps and ask whether to fill them. Do not fill prerequisite gaps until the user agrees. If they agree, run the corresponding earlier UltraPlan skills in canonical order, then return to this stage.
6. If the target is already complete and valid, summarize that state and ask before regenerating or materially changing it.
7. If the user explicitly asks for a proposal, analysis, or discussion only, inspect all relevant evidence and return that without writing artifacts or advancing state. Otherwise, do the stage now; do not stop at a proposal.
8. Use the current effective prompt and concrete paths:

    ultraplan sprint <project> <sprint> prompt review

   The resolved prompt can include workspace or project overrides and therefore takes precedence over the canonical prompt below.
9. Complete the stage-specific workflow, preserve unrelated user edits, and do not cross declared mutation boundaries.
10. Run `ultraplan sprint <project> <sprint> validate review` when supported. Fix validation findings within this stage rather than declaring success early.
11. Run `ultraplan sprint <project> <sprint> status --json` after writes or governed execution so flow-state, freshness, and artifact status are reconciled. Re-run project status if the project or sprint index changed.
12. Inspect downstream artifacts for references made stale by this change. Update directly coupled indexes/references when safe; otherwise report the exact dependent stage that must be revisited. Never delete or silently rewrite downstream decisions.
13. Finish with the artifact/result paths, validation outcome, state transition, and any remaining blocker or dependent stage.

## Stage workflow

For example, the input `projects/ultraplan-go/sprints/30-web-foundations/` resolves to project `ultraplan-go`, sprint `30-web-foundations`, and the target implementation directory declared in `projects/ultraplan-go/project-index.md`. It does not resolve to the workspace repository or to a nested source checkout.

Preview the review scope first:

    ultraplan sprint <project> <sprint> review --dry-run

Then run or resume the governed review:

    ultraplan sprint <project> <sprint> review

The CLI owns reviewer fan-out, frozen inputs, aggregation, verdict calculation, state reconciliation, and creation of the sprint-root `review.md`. Do not replace it with a single-agent ad hoc code review. After it finishes, read the generated `review.md` and the fresh review status, then summarize the verdict, findings by severity, evidence freshness, blockers, and smoke eligibility for the user.

Do not silently fix findings during the review stage; if fixes are requested, return to execute and rerun review so evidence and fingerprints remain authoritative.

## Canonical stage prompt

The current resolved CLI prompt wins over this embedded baseline.

# Automated Sprint Review

This is UltraPlan's embedded default for the product-owned `review` stage. A
workspace `prompts/review.md` is an optional override; it is never a required
planning input.

Review the frozen manifest assembled by UltraPlan. Work read-only against the
approved implementation target. Do not write files, mutate Git, install
dependencies, or expand sprint scope.

UltraPlan starts one independent reviewer for every selected contract and one
for the technical handbook. Each reviewer must return the structured result
requested by the runtime prompt, including an explicit applicability decision,
severity, action, and real path/line citations for every applicable finding.

Treat missing coverage, malformed results, unsafe citations, unsupported
permission enforcement, or changed inputs as review failure. UltraPlan owns
aggregation, deterministic verdict calculation, validation, atomic replacement
of sprint-root `review.md`, and review state in `flow-state.json`.
