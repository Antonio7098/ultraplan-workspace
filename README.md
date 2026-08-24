# UltraPlan Workspace

This workspace stores local UltraPlan studies, planning projects, runtime state, and generated artifacts.

Completed study tasks and sprint stages commit and push their owned files automatically. The workspace uses `git.stage_completion: commit-and-push`, pushes the current upstream branch, and falls back to `origin` for branches without an upstream. Unrelated working-tree and index changes stay out of stage commits.

## Health And Config

```sh
ultraplan health
ultraplan config show
ultraplan config show --json
```

## Studies

```sh
ultraplan study list
ultraplan study <study> list
ultraplan study <study> prompt analysis <dimension> <source>
ultraplan study <study> run <dimension> <source>
ultraplan study <study> synthesize <dimension>
ultraplan study <study> run-loop --parallel 1
ultraplan study <study> validate
ultraplan study <study> status
ultraplan study <study> summary
```

## Planning Projects

```sh
ultraplan project list
ultraplan project <project> status
ultraplan project <project> validate
ultraplan sprint <project> <sprint> status
ultraplan sprint <project> <sprint> validate requirements
ultraplan sprint <project> <sprint> validate sprint-index
ultraplan sprint <project> <sprint> prompt requirements
ultraplan sprint <project> <sprint> prompt sprint-index
ultraplan sprint <project> <sprint> flow --to requirements --dry-run
ultraplan sprint <project> <sprint> flow --to plan --dry-run
```

## Defaults

Prompts and templates are built into the CLI. Materialize editable copies only when you need local overrides.

```sh
ultraplan defaults install --dry-run
ultraplan defaults install
```
