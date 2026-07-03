# CLI Design Reasoning Template

## 1. Command Summary

### Command

```
<tool command subcommand>
```

User goal

<what the user/script is trying to do>

Command type

[ ] Read/query
[ ] Write/mutate
[ ] Generate
[ ] Validate/check
[ ] Init/scaffold
[ ] Migrate
[ ] Delete/destructive
[ ] Debug/diagnostic

---

## 2. Contract Mapping

| Contract              | Requirement IDs | Why Applies |
| --------------------- | --------------- | ----------- |
| cli.md                | <IDs>           | <reason>    |
| configuration.md      | <IDs>           | <reason>    |
| errors.md             | <IDs>           | <reason>    |
| testing.md            | <IDs>           | <reason>    |
| security.md           | <IDs>           | <reason>    |
| privacy-and-data.md   | <IDs>           | <reason>    |
| package-public-api.md | <IDs>           | <reason>    |

---

## 3. Interface Design

Usage

<command usage>

Positional args

- <arg>: <meaning>

Flags

- --<flag>: <meaning/default>

Why this shape?

<why args/flags/subcommands are structured this way>

---

## 4. Output Contract

Human output

stdout: <intended result>
stderr: <progress/warnings/errors>

JSON/machine output

[ ] Not needed
[ ] Needed — schema: <...>

Exit codes

0: success
<code>: <meaning>

---

## 5. Config and Environment

Config precedence: <flags/env/config/defaults>
Secrets touched? <yes/no>
Effective config display? <yes/no>

---

## 6. Safety

Destructive? <yes/no>
Dry run? <yes/no>
Confirmation required? <yes/no>
Non-interactive safe? <yes/no>

---

## 7. Testing

- help/version test: <...>
- stdout/stderr test: <...>
- JSON output test: <...>
- exit code test: <...>
- dry-run/confirmation test: <...>
- config precedence test: <...>

---

## 8. Decision

Decision: <proceed/revise/block>
Reason: <short rationale>
Main trade-off: <human UX vs automation stability/etc.>
