package app

import (
	"context"
	"errors"
	"os"
	"strings"
	"testing"

	"github.com/Antonio7098/ultraplan-go/internal/sprint"
)

func TestQAInvestigator_765e01a64feb(t *testing.T) {
	// Control: classification retains distinct causes and categories before the
	// CLI boundary wraps them in an application exit class.
	cases := []struct {
		name     string
		err      error
		category string
	}{
		{name: "timeout", err: context.DeadlineExceeded, category: "timeout"},
		{name: "process", err: errors.New("process exited with status 1"), category: "process"},
		{name: "persistence", err: os.ErrNotExist, category: "persistence"},
		{name: "configuration", err: &sprint.PerformanceError{Code: "performance.config", Category: "configuration", Operation: "performance.run", Guidance: "repair configuration", Cause: errors.New("invalid performance configuration")}, category: "configuration"},
		{name: "cancellation", err: context.Canceled, category: "cancellation"},
	}
	for _, tc := range cases {
		classifiedErr := sprint.ClassifyPerformanceError("performance.run", tc.err)
		performanceErr, ok := sprint.AsPerformanceError(classifiedErr)
		if !ok || performanceErr.Category != tc.category {
			t.Fatalf("control %s category = %q, want %q", tc.name, performanceErr.Category, tc.category)
		}
	}

	// Verify the production boundary itself still has the reported collapse and
	// that cancellation remains its distinct passing control.
	source, err := os.ReadFile("sprint_commands.go")
	if err != nil {
		t.Fatal(err)
	}
	text := string(source)
	cancelBranch := `if errors.Is(runErr, context.Canceled) {
				return classified(ExitCancel, "sprint.performance: %w", runErr)
			}`
	validationBranch := `return classified(ExitValidation, "sprint.performance: %w", runErr)`
	if !strings.Contains(text, cancelBranch) {
		t.Fatalf("passing control missing cancellation ExitCancel branch")
	}
	if !strings.Contains(text, validationBranch) {
		t.Fatalf("performance failures no longer share the ExitValidation branch")
	}

	// The distinct timeout, process, persistence, and configuration categories
	// all reach the unconditional non-cancellation branch above.
	t.Fatalf("ULTRAPLAN_QA_PREDICTED_FAILURE:TestQAInvestigator_765e01a64feb")
}
