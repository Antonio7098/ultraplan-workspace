package sprint

import (
	"bufio"
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math"
	"math/big"
	"regexp"
	"sort"
	"strconv"
	"strings"

	pprocess "github.com/Antonio7098/ultraplan-go/internal/platform/process"
)

const performanceJSONMarker = "ULTRAPLAN_PERFORMANCE_V1 "

type performanceJSONEnvelope struct {
	SchemaVersion int             `json:"schema_version"`
	TargetID      string          `json:"target_id"`
	Scenario      string          `json:"scenario"`
	Metric        string          `json:"metric"`
	Unit          PerformanceUnit `json:"unit"`
	Value         string          `json:"value"`
}

func ParsePerformanceOutput(descriptor PerformanceDescriptor, target PerformanceTarget, output string, truncated bool) (string, error) {
	if truncated {
		return "", fmt.Errorf("benchmark output was truncated")
	}
	if descriptor.Kind == PerformanceGoBenchmark {
		return parseGoBenchmarkValue(descriptor, output)
	}
	if descriptor.Kind == PerformanceJSONRunner {
		return parsePerformanceJSONValue(descriptor, target, output)
	}
	return "", fmt.Errorf("unsupported descriptor kind %q", descriptor.Kind)
}

func parseGoBenchmarkValue(descriptor PerformanceDescriptor, output string) (string, error) {
	var values []string
	scanner := bufio.NewScanner(strings.NewReader(output))
	for scanner.Scan() {
		fields := strings.Fields(scanner.Text())
		if len(fields) < 3 || benchmarkBaseName(fields[0]) != descriptor.Symbol {
			continue
		}
		for i := 1; i < len(fields); i++ {
			if fields[i] == string(descriptor.RawUnit) && i > 0 {
				canonical, _, err := canonicalPerformanceDecimal(fields[i-1])
				if err != nil {
					return "", fmt.Errorf("parse Go benchmark value: %w", err)
				}
				values = append(values, canonical)
			}
		}
	}
	if err := scanner.Err(); err != nil {
		return "", err
	}
	if len(values) != 1 {
		return "", fmt.Errorf("expected one %s value for %s, found %d", descriptor.RawUnit, descriptor.Symbol, len(values))
	}
	return values[0], nil
}

func benchmarkBaseName(value string) string {
	return regexp.MustCompile(`-[0-9]+$`).ReplaceAllString(value, "")
}

func parsePerformanceJSONValue(descriptor PerformanceDescriptor, target PerformanceTarget, output string) (string, error) {
	var envelopes []performanceJSONEnvelope
	scanner := bufio.NewScanner(strings.NewReader(output))
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		at := strings.Index(line, performanceJSONMarker)
		if at < 0 {
			continue
		}
		decoder := json.NewDecoder(strings.NewReader(strings.TrimSpace(line[at+len(performanceJSONMarker):])))
		decoder.DisallowUnknownFields()
		var envelope performanceJSONEnvelope
		if err := decoder.Decode(&envelope); err != nil {
			return "", fmt.Errorf("decode performance JSON: %w", err)
		}
		if decoder.Decode(&struct{}{}) == nil {
			return "", fmt.Errorf("performance JSON contains trailing data")
		}
		envelopes = append(envelopes, envelope)
	}
	if len(envelopes) != 1 {
		return "", fmt.Errorf("expected one performance JSON envelope, found %d", len(envelopes))
	}
	envelope := envelopes[0]
	if envelope.SchemaVersion != PerformanceSchemaVersion || envelope.TargetID != target.ID || envelope.Scenario != target.Scenario || envelope.Metric != target.Metric || envelope.Unit != descriptor.RawUnit {
		return "", fmt.Errorf("performance JSON identity does not match the frozen target and descriptor")
	}
	canonical, _, err := canonicalPerformanceDecimal(envelope.Value)
	return canonical, err
}

func QualifyPerformanceSamples(samples []PerformanceSample, cvLimit float64) (PerformanceQualification, error) {
	qualification := PerformanceQualification{Aggregation: "median-v1", Dispersion: "cv-v1", Samples: len(samples)}
	if len(samples) < PerformanceMinSamples || cvLimit < 0 || cvLimit > MaxPerformanceLimits().CVPercent {
		return qualification, fmt.Errorf("invalid sample count or CV limit")
	}
	values := make([]*big.Rat, len(samples))
	identity := samples[0]
	for i, sample := range samples {
		if !sample.CleanupComplete {
			qualification.Reason = "process cleanup is uncertain"
			return qualification, nil
		}
		if sample.TargetID != identity.TargetID || sample.Scenario != identity.Scenario || sample.Metric != identity.Metric || sample.Unit != identity.Unit || sample.CommandIdentity != identity.CommandIdentity || sample.BenchmarkDigest != identity.BenchmarkDigest || sample.PacketDigest != identity.PacketDigest || sample.EnvironmentID != identity.EnvironmentID {
			qualification.Reason = "sample identity drift"
			return qualification, nil
		}
		_, rat, err := canonicalPerformanceDecimal(sample.Value)
		if err != nil {
			return qualification, fmt.Errorf("sample %d: %w", i+1, err)
		}
		values[i] = rat
	}
	sort.Slice(values, func(i, j int) bool { return values[i].Cmp(values[j]) < 0 })
	median := new(big.Rat)
	mid := len(values) / 2
	if len(values)%2 == 1 {
		median.Set(values[mid])
	} else {
		median.Add(values[mid-1], values[mid]).Quo(median, big.NewRat(2, 1))
	}
	sum := new(big.Rat)
	for _, value := range values {
		sum.Add(sum, value)
	}
	meanRat := new(big.Rat).Quo(sum, big.NewRat(int64(len(values)), 1))
	mean, _ := new(big.Float).SetRat(meanRat).Float64()
	if math.IsInf(mean, 0) || math.IsNaN(mean) {
		return qualification, fmt.Errorf("sample mean is non-finite")
	}
	var squared float64
	allZero := true
	for _, value := range values {
		valueFloat, _ := new(big.Float).SetRat(value).Float64()
		if value.Sign() != 0 {
			allZero = false
		}
		delta := valueFloat - mean
		squared += delta * delta
	}
	stddev := math.Sqrt(squared / float64(len(values)-1))
	cv := 0.0
	if mean == 0 && !allZero {
		qualification.Reason = "zero mean with nonzero samples"
	} else if mean != 0 {
		cv = 100 * stddev / math.Abs(mean)
	}
	qualification.Median = exactPerformanceDecimal(median)
	qualification.Mean = exactPerformanceDecimal(meanRat)
	qualification.StdDev = stddev
	qualification.CVPercent = cv
	if qualification.Reason == "" && cv <= cvLimit {
		qualification.Qualified = true
	} else if qualification.Reason == "" {
		qualification.Reason = fmt.Sprintf("CV %.6g exceeds limit %.6g", cv, cvLimit)
	}
	return qualification, nil
}

func exactPerformanceDecimal(value *big.Rat) string {
	if value == nil {
		return ""
	}
	denominator := new(big.Int).Set(value.Denom())
	two, five := big.NewInt(2), big.NewInt(5)
	twos, fives := 0, 0
	for new(big.Int).Mod(denominator, two).Sign() == 0 {
		denominator.Quo(denominator, two)
		twos++
	}
	for new(big.Int).Mod(denominator, five).Sign() == 0 {
		denominator.Quo(denominator, five)
		fives++
	}
	if denominator.Cmp(big.NewInt(1)) != 0 {
		return value.RatString()
	}
	digits := twos
	if fives > digits {
		digits = fives
	}
	text := value.FloatString(digits)
	if strings.Contains(text, ".") {
		text = strings.TrimRight(strings.TrimRight(text, "0"), ".")
	}
	if text == "-0" || text == "" {
		return "0"
	}
	return text
}

func MeasurePerformanceTarget(ctx context.Context, runner pprocess.Runner, root string, descriptor PerformanceDescriptor, target PerformanceTarget, limits PerformanceLimits, env []string, packetDigest, benchmarkDigest, environmentID string) (PerformanceMeasurement, error) {
	measurement := PerformanceMeasurement{Target: target, Warmups: limits.Warmups}
	if runner == nil {
		return measurement, fmt.Errorf("performance process runner is required")
	}
	request, err := PerformanceProcessRequest(root, descriptor, limits, env)
	if err != nil {
		return measurement, err
	}
	for i := 0; i < limits.Warmups+target.Samples; i++ {
		result, runErr := runner.Run(ctx, request)
		if runErr != nil {
			return measurement, runErr
		}
		if !result.CleanupComplete {
			return measurement, fmt.Errorf("performance process cleanup is uncertain")
		}
		value, parseErr := ParsePerformanceOutput(descriptor, target, result.Stdout, result.StdoutTruncated)
		if parseErr != nil {
			return measurement, parseErr
		}
		if i < limits.Warmups {
			continue
		}
		sum := sha256.Sum256([]byte(result.Stdout))
		measurement.Samples = append(measurement.Samples, PerformanceSample{TargetID: target.ID, Scenario: target.Scenario, Metric: target.Metric, Unit: descriptor.RawUnit, Value: value, CommandIdentity: descriptor.OutputIdentity, BenchmarkDigest: benchmarkDigest, PacketDigest: packetDigest, EnvironmentID: environmentID, OutputDigest: hex.EncodeToString(sum[:]), OutputBytes: len(result.Stdout), CleanupComplete: result.CleanupComplete})
	}
	measurement.Qualification, err = QualifyPerformanceSamples(measurement.Samples, limits.CVPercent)
	return measurement, err
}

func performanceOutputLine(envelope performanceJSONEnvelope) string {
	data, _ := json.Marshal(envelope)
	return performanceJSONMarker + string(bytes.TrimSpace(data))
}

func parsePerformanceNumber(value string) (float64, error) {
	n, err := strconv.ParseFloat(value, 64)
	if err != nil || math.IsInf(n, 0) || math.IsNaN(n) {
		return 0, fmt.Errorf("non-finite performance number")
	}
	return n, nil
}
