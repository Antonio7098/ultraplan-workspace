package app

import (
	"context"
	"errors"
	"path/filepath"
	"testing"
)

func TestQAInvestigator_6fe849fce7fe(t *testing.T) {
	// Control: the production classifier produces the typed shape consumed by
	// performance-specific transport mapping.
	control := mapPerformanceUseCaseError("performance.result", errors.New("performance result is unavailable"))
	typedControl, ok := AsPerformanceUseCaseError(control)
	if !ok || typedControl.Code != "performance.persistence_unavailable" || typedControl.Category != "persistence" {
		t.Fatalf("control performance mapping = %#v, typed=%t", control, ok)
	}

	root := initializedWorkspace(t)
	projectRoot := filepath.Join(root, "projects", "alpha")
	writeFixtureFileContent(t, projectRoot, "# Project Index\n\n## Performance Policy\n\n- **Mode:** disabled\n", "project-index.md")
	writeFixtureFileContent(t, filepath.Join(projectRoot, "sprints", "01-ready"), "# Requirements\n\nNo performance targets are enabled.\n", "requirements.md")

	_, err := (dashboardUseCases{root: root, readOnly: true}).PerformanceResult(context.Background(), PerformanceRequest{Project: "alpha", Sprint: "01-ready"})
	if err == nil {
		t.Fatal("expected a nonterminal performance result to be unavailable")
	}
	if _, ok := AsPerformanceUseCaseError(err); !ok {
		t.Fatalf("ULTRAPLAN_QA_PREDICTED_FAILURE:TestQAInvestigator_6fe849fce7fe assertion=Assert that the error is not recognized as a PerformanceUseCaseError and receives generic mapping.")
	}
}
