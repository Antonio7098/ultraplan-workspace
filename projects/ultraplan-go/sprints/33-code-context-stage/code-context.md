# Sprint Code Context

## Sprint Scope

Sprint 33 adds the code-context planning stage itself. This compatibility artifact marks that newly introduced prerequisite as satisfied for the sprint that implemented it, allowing the already-completed downstream work and review evidence to remain usable.

## Inspected Repository Areas

- `internal/sprint/code_context.go` owns validation and execution of the new stage.
- `internal/sprint/domain.go` places code-context between requirements and sprint-index.
- `internal/sprint/state.go` provides compatibility behavior for persisted pre-code-context flows.

## Selected Source Excerpts

### Code-context stage ownership

- **Path:** `internal/sprint/code_context.go`
- **Lines:** `1-12`
- **Symbol:** `sprint` package
- **Rationale:** This file is the implementation introduced by Sprint 33 and demonstrates that code-context belongs to the sprint workflow module.

```go
package sprint

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)
```

## Relationships

The sprint service resolves the execution target, invokes the shared runtime with read-only permissions, validates the returned Markdown, and promotes the candidate artifact before later planning stages consume it.

## Constraints

The stage may inspect source but must not mutate the target repository. Its output is restricted to the sprint-owned `code-context.md` artifact and must remain compatible with previously persisted flows.

## Open Questions

None for this compatibility stub. Full live-provider dogfood remains subsequent-sprint work.
