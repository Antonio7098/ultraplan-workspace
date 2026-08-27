# Sprint merge

- Source: `ultraplan/ultraplan-go/38-bounded-repair` at `32ffead1d5cbd9315c8580cfd774f6daad206ed9`
- Target: `main` from `3686b481fe6253a9c09c6cb8100bac5f8163d616`
- Merge commit: `9946038b5e2f3b26c139a603465b0ca1623e6788`

## Merge bounded QA repair and run observability

- Introduce a bounded post-execution QA repair loop in internal/sprint via new qa_repair.go, qa_repair_state.go, and qa_types.go, with test coverage in qa_repair_test.go, qa_state_test.go, and verify_test.go.
- Add run observability through internal/app/durable_operations.go, operation_runner.go, and operations.go so the repair flow can inspect and react to operation status.
- Extend the CLI surface in internal/app via sprint_commands.go, sprint_usecases.go, and sprint/smoke.go plus smoke_types.go to expose repair bounds and run observability.
- Wire the bounded repair and run views into internal/web (qa_handlers.go, run_handlers.go, operation_handlers.go, operations.go, routes.go, templates/run_qa.html, templates/sprint.html) and internal/tui (qa_view.go, model.go, app.go, views.go) with handler and view tests.
- Add QA bounds configuration in internal/platform/config/config.go and the new qa.go with config_test.go coverage so bounds are tunable per deployment.
- Refresh DESIGN.md, docs/architecture.md, docs/cli-reference.md, docs/local-web.md, docs/phase3-json-schemas.md, docs/recovery.md, docs/release-checklist.md, docs/user-guide.md, and docs/plans/integrated-roadmap.md plus the new post-execution-qa-and-repair-loop plan to document the flow.
- Tighten quality gates with new tests under internal/app (sprint_commands_test.go, run_control_inventory_test.go), internal/web (operations_test.go, qa_handlers_test.go, run_handlers_test.go), and internal/tui (qa_view_test.go).