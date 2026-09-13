package config

import (
	"os"
	"strings"
	"testing"
)

func TestQAInvestigator_0c580cfa0ca8(t *testing.T) {
	readSource := func(t *testing.T, path string) string {
		t.Helper()
		content, err := os.ReadFile(path)
		if err != nil {
			t.Fatalf("read %s: %v", path, err)
		}
		return string(content)
	}

	t.Run("control_generic_config_infrastructure_is_allowed", func(t *testing.T) {
		if !strings.Contains(readSource(t, "config.go"), "func Load(") {
			t.Fatal("control failed: generic config Load function was not found")
		}
	})

	performanceSource := readSource(t, "performance.go")
	hasProductDefaults := strings.Contains(performanceSource, "func DefaultPerformance(")
	hasProductValidation := strings.Contains(performanceSource, "func validatePerformance(")
	if hasProductDefaults && hasProductValidation {
		println("assertion: Assert the current source defines product-specific performance defaults and validation under internal/platform.")
		println("ULTRAPLAN_QA_PREDICTED_FAILURE:TestQAInvestigator_0c580cfa0ca8")
		t.Fail()
	}
}
