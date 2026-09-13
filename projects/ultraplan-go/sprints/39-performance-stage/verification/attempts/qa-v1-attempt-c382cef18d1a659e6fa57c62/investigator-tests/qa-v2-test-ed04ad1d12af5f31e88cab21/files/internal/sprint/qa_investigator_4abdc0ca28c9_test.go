package sprint

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
)

func TestQAInvestigator_4abdc0ca28c9(t *testing.T) {
	packet := PerformanceTargetPacket{
		SchemaVersion: PerformanceSchemaVersion,
		ParserVersion: PerformanceParserVersion,
		PacketDigest:  "qa-investigator-packet",
		Targets: []PerformanceTarget{{
			ID:    "PERF-QA",
			Unit:  PerformanceNSPerOp,
			Basis: PerformanceAbsolute,
		}},
	}

	discover := func(t *testing.T, separator string) error {
		t.Helper()
		root := t.TempDir()
		source := "package sample\n\n// ultraplan:performance-v1 go-benchmark-v1 ns/op PERF-QA\n" + separator + "func BenchmarkQA() {}\n"
		if err := os.WriteFile(filepath.Join(root, "performance_test.go"), []byte(source), 0o600); err != nil {
			t.Fatal(err)
		}
		_, err := DiscoverPerformanceManifest(root, packet)
		return err
	}

	if err := discover(t, ""); err != nil {
		t.Fatalf("adjacent marker control was rejected: %v", err)
	}
	if err := discover(t, "var intervening = true\n"); err == nil {
		t.Fatal("non-comment source statement control was accepted")
	}

	blankErr := discover(t, "\n")
	commentErr := discover(t, "// unrelated comment\n")
	if blankErr == nil && commentErr == nil {
		fmt.Println("ULTRAPLAN_QA_PREDICTED_FAILURE:TestQAInvestigator_4abdc0ca28c9")
		t.Fatal("Assert DiscoverPerformanceManifest accepts both separated marker cases.")
	}
	if blankErr != nil {
		t.Fatalf("blank-line reproduction was not accepted: %v", blankErr)
	}
	t.Fatalf("comment-line reproduction was not accepted: %v", commentErr)
}
