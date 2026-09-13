package app

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestQAInvestigator_340a1300e236(t *testing.T) {
	// Passing control: the normal classifier creates the typed code/category
	// consumed by the performance-specific HTTP branch.
	control := mapPerformanceUseCaseError("performance.result", errors.New("performance result is unavailable"))
	typedControl, ok := AsPerformanceUseCaseError(control)
	if !ok || typedControl.Code != "performance.persistence_unavailable" || typedControl.Category != "persistence" {
		t.Fatalf("typed control = %#v, recognized=%t", control, ok)
	}
	handlerSource, err := os.ReadFile(filepath.Join("..", "web", "handlers.go"))
	if err != nil {
		t.Fatal(err)
	}
	handlerText := string(handlerSource)
	if !strings.Contains(handlerText, "app.AsPerformanceUseCaseError(err)") ||
		!strings.Contains(handlerText, "case \"persistence\", \"process\", \"timeout\", \"cleanup\":") ||
		!strings.Contains(handlerText, "status = http.StatusServiceUnavailable") {
		t.Fatal("passing control did not locate performance-specific HTTP mapping")
	}

	// Exercise the real result entry point in a valid read-only, nonterminal
	// performance state. No runtime or mutation is needed to request a result.
	root := initializedWorkspace(t)
	projectRoot := filepath.Join(root, "projects", "alpha")
	writeFixtureFileContent(t, projectRoot, "# Project Index\n\n## Performance Policy\n\n- **Mode:** disabled\n", "project-index.md")
	writeFixtureFileContent(t, filepath.Join(projectRoot, "sprints", "01-ready"), "# Requirements\n\nNo performance targets are enabled.\n", "requirements.md")
	u := dashboardUseCases{root: root, readOnly: true}
	status, err := u.PerformanceStatus(context.Background(), PerformanceRequest{Project: "alpha", Sprint: "01-ready"})
	if err != nil || status.Phase != "ready" || !status.Fresh || status.Result != nil {
		t.Fatalf("nonterminal fixture status = %#v, err=%v", status, err)
	}
	_, resultErr := u.PerformanceResult(context.Background(), PerformanceRequest{Project: "alpha", Sprint: "01-ready"})
	if resultErr == nil {
		t.Fatal("nonterminal result unexpectedly available")
	}
	_, resultTyped := AsPerformanceUseCaseError(resultErr)

	// Contract check for the aggregate entry point: direct status consults run
	// control, but the summary helper delegates to the un-enriched projection.
	useCaseSource, err := os.ReadFile("sprint_usecases.go")
	if err != nil {
		t.Fatal(err)
	}
	useCaseText := string(useCaseSource)
	directStart := strings.Index(useCaseText, "func (u dashboardUseCases) PerformanceStatus(")
	directEnd := strings.Index(useCaseText[directStart:], "\nfunc ")
	summaryStart := strings.Index(useCaseText, "func performanceSummaryForSprint(")
	summaryEnd := strings.Index(useCaseText[summaryStart:], "\nfunc ")
	if directStart < 0 || directEnd < 0 || summaryStart < 0 || summaryEnd < 0 {
		t.Fatal("lifecycle entry-point fixture could not be bounded")
	}
	directBody := useCaseText[directStart : directStart+directEnd]
	summaryBody := useCaseText[summaryStart : summaryStart+summaryEnd]
	if !strings.Contains(directBody, "u.runs.Run(ctx") || !strings.Contains(directBody, "OperationalLifecycle") || !strings.Contains(directBody, "CancellationState") {
		t.Fatal("passing control did not locate direct durable lifecycle enrichment")
	}
	aggregateOmitsRunEnrichment := strings.Contains(summaryBody, "performanceStatusProjection(service") && !strings.Contains(summaryBody, ".Run(") && !strings.Contains(summaryBody, "PerformanceStatus(")
	if aggregateOmitsRunEnrichment && !resultTyped {
		t.Fatalf("ULTRAPLAN_QA_PREDICTED_FAILURE:TestQAInvestigator_340a1300e236 assertion=Assert that the aggregate result omits lifecycle or cancellation fields and that the unavailable-result error is not recognized as a PerformanceUseCaseError and receives generic web handling.")
	}
	if !aggregateOmitsRunEnrichment {
		t.Fatal("aggregate lifecycle path is enriched")
	}
	t.Fatal("unavailable-result error is typed")
}
