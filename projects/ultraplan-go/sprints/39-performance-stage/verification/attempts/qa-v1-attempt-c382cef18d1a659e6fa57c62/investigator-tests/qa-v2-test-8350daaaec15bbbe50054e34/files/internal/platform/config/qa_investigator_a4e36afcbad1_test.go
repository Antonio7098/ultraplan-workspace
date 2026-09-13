package config

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"testing"
)

func TestQAInvestigator_a4e36afcbad1(t *testing.T) {
	declaredFunctions := func(t *testing.T, path string) map[string]bool {
		t.Helper()
		file, err := parser.ParseFile(token.NewFileSet(), path, nil, 0)
		if err != nil {
			t.Fatalf("parse %s: %v", path, err)
		}
		functions := make(map[string]bool)
		for _, declaration := range file.Decls {
			if function, ok := declaration.(*ast.FuncDecl); ok {
				functions[function.Name.Name] = true
			}
		}
		return functions
	}

	t.Run("control_generic_config_infrastructure_is_allowed", func(t *testing.T) {
		functions := declaredFunctions(t, "config.go")
		if !functions["Load"] {
			t.Fatal("control failed: generic config Load function was not found")
		}
	})

	functions := declaredFunctions(t, "performance.go")
	forbiddenPolicy := []string{"DefaultPerformance", "validatePerformance"}
	var found []string
	for _, name := range forbiddenPolicy {
		if functions[name] {
			found = append(found, name)
		}
	}
	if len(found) != 0 {
		fmt.Println("assertion: Assert the current source defines product-specific performance defaults and validation under internal/platform.")
		fmt.Println("ULTRAPLAN_QA_PREDICTED_FAILURE:TestQAInvestigator_a4e36afcbad1")
		t.Fail()
	}
}
