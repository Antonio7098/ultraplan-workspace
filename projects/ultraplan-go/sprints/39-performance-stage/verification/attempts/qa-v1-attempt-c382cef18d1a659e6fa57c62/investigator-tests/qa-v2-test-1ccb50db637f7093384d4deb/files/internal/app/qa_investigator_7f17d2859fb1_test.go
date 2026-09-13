package app

import (
	"strings"
	"testing"

	"github.com/Antonio7098/ultraplan-go/internal/platform/config"
)

func TestQAInvestigator_7f17d2859fb1(t *testing.T) {
	loadWithCommands := func(value string) error {
		_, err := config.Load(config.LoadOptions{Env: func(key string) string {
			if key == "ULTRAPLAN_PERFORMANCE_COMMANDS" {
				return value
			}
			return ""
		}})
		return err
	}

	if err := loadWithCommands(""); err != nil {
		t.Fatalf("default performance.commands should load: %v", err)
	}
	if err := loadWithCommands("64"); err != nil {
		t.Fatalf("lower performance.commands value should load: %v", err)
	}
	if err := loadWithCommands("256"); err != nil {
		if !strings.Contains(err.Error(), "performance.commands") {
			t.Fatalf("documented maximum failed for an unrelated reason: %v", err)
		}
		t.Fatalf("ULTRAPLAN_QA_PREDICTED_FAILURE:TestQAInvestigator_7f17d2859fb1 assertion=Assert configuration loading rejects those documented maximum values. documented performance.commands hard maximum 256 was rejected: %v", err)
	}
}
