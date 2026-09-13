package app

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/Antonio7098/ultraplan-go/internal/sprint"
)

func TestParseSprintPerformanceArgsRejectsAuthorityOverrides(t *testing.T) {
	for _, args := range [][]string{
		nil,
		{"--samples", "5"},
		{"--command", "go test ./..."},
		{"--parser", "custom"},
		{"--environment", "SECRET=value"},
		{"--limit", "commands=256"},
		{"status", "--yes"},
		{"resume"},
	} {
		if _, err := parseSprintPerformanceArgs(args); err == nil {
			t.Fatalf("accepted unsafe or unconfirmed arguments %q", args)
		}
	}
	for _, args := range [][]string{{"--dry-run"}, {"status", "--json"}, {"--yes"}, {"resume", "--yes"}, {"cancel"}, {"recover"}} {
		if _, err := parseSprintPerformanceArgs(args); err != nil {
			t.Fatalf("rejected %q: %v", args, err)
		}
	}
}

func TestPerformanceProjectionOmitsPrivateEvidence(t *testing.T) {
	prepared := sprint.PerformancePrepareResult{SchemaVersion: 1, Project: "alpha", Sprint: "39-performance", PolicyMode: "enabled", PolicyDigest: strings.Repeat("a", 64), Enabled: true, Ready: true, DryRun: true, MutationPossible: true, Limits: sprint.DefaultPerformanceLimits(), Packet: sprint.PerformanceTargetPacket{Targets: []sprint.PerformanceTarget{{ID: "PERF-A", Scenario: "journey", Metric: "latency", Comparator: sprint.PerformanceLessOrEqual, Value: "10", Unit: sprint.PerformanceMS, Gate: sprint.PerformanceRequired, Samples: 5, Basis: sprint.PerformanceAbsolute}}}, NextAction: "Start performance."}
	encoded, err := json.Marshal(performancePrepareProjection(prepared))
	if err != nil {
		t.Fatal(err)
	}
	text := string(encoded)
	for _, forbidden := range []string{`"raw_samples"`, `"command_output"`, `"profile"`, `"patch"`, `"prompt"`, `"provider_payload"`, `"environment_values"`, `"secret"`} {
		if strings.Contains(text, forbidden) {
			t.Fatalf("public DTO contains forbidden field %q: %s", forbidden, text)
		}
	}
	for _, required := range []string{`"policy_mode":"enabled"`, `"id":"PERF-A"`, `"samples":5`, `"limits"`} {
		if !strings.Contains(text, required) {
			t.Fatalf("public DTO is missing %q: %s", required, text)
		}
	}
}
