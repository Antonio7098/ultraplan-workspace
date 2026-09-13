package app

import (
	"bytes"
	"encoding/json"
	"html/template"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/Antonio7098/ultraplan-go/internal/sprint"
)

func TestQAInvestigator_4c191617c383(t *testing.T) {
	const (
		artifactPath    = ".ultraplan/performance/attempts/performance-v1-attempt-0123456789abcdef01234567/profile.json"
		artifactDigest  = "0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef"
		artifactContent = "synthetic-private-artifact-content"
	)
	record := sprint.PerformanceEvidenceRecord{
		AttemptID: "performance-v1-attempt-0123456789abcdef01234567",
		Name:      "profile.json",
		Path:      artifactPath,
		Digest:    artifactDigest,
		Size:      int64(len(artifactContent)),
	}
	evidence := PerformanceEvidenceResult{SchemaVersion: 1, Records: []sprint.PerformanceEvidenceRecord{record}, ReturnedCount: 1}

	encoded, err := json.Marshal(evidence)
	if err != nil {
		t.Fatal(err)
	}
	if bytes.Contains(encoded, []byte(artifactContent)) {
		t.Fatalf("control failed: JSON disclosed raw artifact contents")
	}

	_, testFile, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("locate test source")
	}
	templatePath := filepath.Join(filepath.Dir(testFile), "..", "web", "templates", "sprint.html")
	page, err := template.New("sprint").Parse(`{{define "layout/top"}}{{end}}{{define "layout/bottom"}}{{end}}`)
	if err != nil {
		t.Fatal(err)
	}
	page, err = page.ParseFiles(templatePath)
	if err != nil {
		t.Fatal(err)
	}
	data := struct {
		Page   string
		Sprint SprintSummary
		CSRF   string
	}{
		Page: "performance",
		Sprint: SprintSummary{
			Project: "synthetic-project",
			Slug:    "synthetic-sprint",
			Performance: PerformanceStatusResult{
				Project:              "synthetic-project",
				Sprint:               "synthetic-sprint",
				Enabled:              true,
				PerformanceAttemptID: record.AttemptID,
			},
			PerformanceEvidence: evidence,
		},
	}
	var rendered bytes.Buffer
	if err := page.ExecuteTemplate(&rendered, "page/sprint", data); err != nil {
		t.Fatal(err)
	}
	if strings.Contains(rendered.String(), artifactContent) {
		t.Fatalf("control failed: workbench disclosed raw artifact contents")
	}

	jsonExposesMetadata := bytes.Contains(encoded, []byte(`"path":"`+artifactPath+`"`)) &&
		bytes.Contains(encoded, []byte(`"digest":"`+artifactDigest+`"`))
	renderedExposesDigest := strings.Contains(rendered.String(), artifactDigest)
	if jsonExposesMetadata && renderedExposesDigest {
		t.Fatalf("ULTRAPLAN_QA_PREDICTED_FAILURE:TestQAInvestigator_4c191617c383 assertion: Assert the JSON response contains path and digest fields and rendered workbench output contains the digest.")
	}
}
