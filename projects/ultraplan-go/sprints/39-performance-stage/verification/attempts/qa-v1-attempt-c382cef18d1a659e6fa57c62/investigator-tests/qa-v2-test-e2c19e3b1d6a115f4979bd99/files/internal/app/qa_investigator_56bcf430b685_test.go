package app

import (
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"strings"
	"testing"
)

func TestQAInvestigator_56bcf430b685(t *testing.T) {
	data, err := os.ReadFile("sprint_usecases.go")
	if err != nil {
		t.Fatal(err)
	}
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, "sprint_usecases.go", data, 0)
	if err != nil {
		t.Fatal(err)
	}
	bodies := map[string]string{}
	for _, decl := range file.Decls {
		fn, ok := decl.(*ast.FuncDecl)
		if !ok || fn.Body == nil || (fn.Name.Name != "PerformanceStatus" && fn.Name.Name != "performanceSummaryForSprint") {
			continue
		}
		start := fset.Position(fn.Body.Pos()).Offset
		end := fset.Position(fn.Body.End()).Offset
		bodies[fn.Name.Name] = string(data[start:end])
	}

	// Controls: both claimed entry points are present, direct status consults
	// durable run control, and sprint summaries delegate to the raw projection.
	direct, directOK := bodies["PerformanceStatus"]
	summary, summaryOK := bodies["performanceSummaryForSprint"]
	if !directOK || !summaryOK {
		t.Fatalf("fixture entry points missing: direct=%t summary=%t", directOK, summaryOK)
	}
	if !strings.Contains(direct, "u.runs.Run(ctx") || !strings.Contains(direct, "OperationalLifecycle") || !strings.Contains(direct, "CancellationState") {
		t.Fatalf("control direct status does not enrich durable lifecycle: %s", direct)
	}
	if !strings.Contains(summary, "performanceStatusProjection(service") {
		t.Fatalf("control summary does not use the claimed projection entry point: %s", summary)
	}

	if !strings.Contains(summary, ".Run(") && !strings.Contains(summary, "PerformanceStatus(") {
		t.Fatalf("ULTRAPLAN_QA_PREDICTED_FAILURE:TestQAInvestigator_56bcf430b685 assertion=Confirm only if: A state with a non-empty operation run ID and an inspectable durable lifecycle produces an empty operational_lifecycle or cancellation_state through SprintSummaries or api_sprint_performance while direct PerformanceStatus supplies those fi")
	}
}
