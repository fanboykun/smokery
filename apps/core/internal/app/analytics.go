package app

import (
	"context"
	"encoding/json"
	"sort"
	"time"

	"github.com/google/uuid"

	"github.com/fanboykun/smokery/apps/core/internal/model"
	"github.com/fanboykun/smokery/apps/core/internal/port"
)

type AnalyticsService struct {
	runs port.RunRepo
}

func NewAnalyticsService(r port.RunRepo) *AnalyticsService {
	return &AnalyticsService{runs: r}
}

// --- Latency ---

type LatencyPoint struct {
	Timestamp string  `json:"timestamp"`
	P50       float64 `json:"p50"`
	P95       float64 `json:"p95"`
	P99       float64 `json:"p99"`
}

type LatencyOp struct {
	OperationID string  `json:"operation_id"`
	AvgLatency  float64 `json:"avg_latency"`
	P99Latency  float64 `json:"p99_latency,omitempty"`
}

type LatencyAnalytics struct {
	Range             string       `json:"range"`
	Data              []LatencyPoint `json:"data"`
	SlowestOperations []LatencyOp  `json:"slowest_operations"`
	FastestOperations []LatencyOp  `json:"fastest_operations"`
}

func (s *AnalyticsService) Latency(ctx context.Context, projectID uuid.UUID, rangeStr string) (*LatencyAnalytics, error) {
	runs, err := s.runs.ListByProject(ctx, projectID)
	if err != nil {
		return nil, err
	}

	days := rangeToDays(rangeStr)
	cutoff := time.Now().AddDate(0, 0, -days)

	opLatencies := map[string][]int64{}
	var points []LatencyPoint

	for _, run := range runs {
		if run.CreatedAt.Before(cutoff) {
			continue
		}
		result := s.loadRunResult(ctx, run.ID)
		if result == nil {
			continue
		}
		var durations []int64
		for _, f := range result.Flows {
			for _, step := range f.Steps {
				durations = append(durations, step.Duration)
				opLatencies[step.Name] = append(opLatencies[step.Name], step.Duration)
			}
		}
		for _, su := range result.Suites {
			for _, c := range su.Cases {
				durations = append(durations, c.Step.Duration)
				opLatencies[c.OperationID] = append(opLatencies[c.OperationID], c.Step.Duration)
			}
		}
		if len(durations) > 0 {
			sort.Slice(durations, func(i, j int) bool { return durations[i] < durations[j] })
			points = append(points, LatencyPoint{
				Timestamp: run.CreatedAt.Format(time.RFC3339),
				P50:       float64(percentile(durations, 50)),
				P95:       float64(percentile(durations, 95)),
				P99:       float64(percentile(durations, 99)),
			})
		}
	}

	// Compute per-op averages
	type opAvg struct {
		id  string
		avg float64
		p99 float64
	}
	var ops []opAvg
	for id, lats := range opLatencies {
		sort.Slice(lats, func(i, j int) bool { return lats[i] < lats[j] })
		sum := int64(0)
		for _, l := range lats {
			sum += l
		}
		ops = append(ops, opAvg{id: id, avg: float64(sum) / float64(len(lats)), p99: float64(percentile(lats, 99))})
	}
	sort.Slice(ops, func(i, j int) bool { return ops[i].avg > ops[j].avg })

	analytics := &LatencyAnalytics{Range: rangeStr, Data: points}
	for i, op := range ops {
		if i >= 5 {
			break
		}
		analytics.SlowestOperations = append(analytics.SlowestOperations, LatencyOp{OperationID: op.id, AvgLatency: op.avg, P99Latency: op.p99})
	}
	for i := len(ops) - 1; i >= 0 && len(analytics.FastestOperations) < 5; i-- {
		analytics.FastestOperations = append(analytics.FastestOperations, LatencyOp{OperationID: ops[i].id, AvgLatency: ops[i].avg})
	}
	return analytics, nil
}

// --- Flaky Operations ---

type FlakyOp struct {
	OperationID    string  `json:"operation_id"`
	Path           string  `json:"path"`
	Method         string  `json:"method"`
	SuccessRate    float64 `json:"success_rate"`
	Runs           int     `json:"runs"`
	Failures       int     `json:"failures"`
	FlakinessScore float64 `json:"flakiness_score"`
	Trend          string  `json:"trend"`
}

type FlakyAnalytics struct {
	Range         string   `json:"range"`
	Operations    []FlakyOp `json:"operations"`
	CriticalFlaky []FlakyOp `json:"critical_flaky"`
}

func (s *AnalyticsService) FlakyOperations(ctx context.Context, projectID uuid.UUID, rangeStr string) (*FlakyAnalytics, error) {
	runs, err := s.runs.ListByProject(ctx, projectID)
	if err != nil {
		return nil, err
	}

	days := rangeToDays(rangeStr)
	cutoff := time.Now().AddDate(0, 0, -days)

	type opStats struct {
		total, failed int
		path, method  string
	}
	stats := map[string]*opStats{}

	for _, run := range runs {
		if run.CreatedAt.Before(cutoff) {
			continue
		}
		result := s.loadRunResult(ctx, run.ID)
		if result == nil {
			continue
		}
		for _, su := range result.Suites {
			for _, c := range su.Cases {
				st := stats[c.OperationID]
				if st == nil {
					st = &opStats{}
					stats[c.OperationID] = st
				}
				st.total++
				if c.Step.Status != "passed" {
					st.failed++
				}
			}
		}
	}

	analytics := &FlakyAnalytics{Range: rangeStr}
	for id, st := range stats {
		if st.failed == 0 {
			continue
		}
		successRate := float64(st.total-st.failed) / float64(st.total) * 100
		score := float64(st.failed) / float64(st.total) * 100
		op := FlakyOp{
			OperationID: id, Path: st.path, Method: st.method,
			SuccessRate: successRate, Runs: st.total, Failures: st.failed,
			FlakinessScore: score, Trend: "stable",
		}
		if score >= 70 {
			analytics.CriticalFlaky = append(analytics.CriticalFlaky, op)
		} else {
			analytics.Operations = append(analytics.Operations, op)
		}
	}
	return analytics, nil
}

// --- Health Trends ---

type HealthPoint struct {
	Timestamp  string  `json:"timestamp"`
	Date       string  `json:"date"`
	TotalRuns  int     `json:"total_runs"`
	PassedRuns int     `json:"passed_runs"`
	FailedRuns int     `json:"failed_runs"`
	PassRate   float64 `json:"pass_rate"`
}

type HealthTrends struct {
	Range          string        `json:"range"`
	Data           []HealthPoint `json:"data"`
	CurrentHealth  float64       `json:"current_health"`
	Trend          string        `json:"trend"`
	WeeklyAverage  float64       `json:"weekly_average"`
	MonthlyAverage float64       `json:"monthly_average"`
}

func (s *AnalyticsService) HealthTrends(ctx context.Context, projectID uuid.UUID, rangeStr string) (*HealthTrends, error) {
	runs, err := s.runs.ListByProject(ctx, projectID)
	if err != nil {
		return nil, err
	}

	days := rangeToDays(rangeStr)
	cutoff := time.Now().AddDate(0, 0, -days)

	// Group by date
	type dayStats struct{ total, passed int }
	byDate := map[string]*dayStats{}

	for _, run := range runs {
		if run.CreatedAt.Before(cutoff) {
			continue
		}
		date := run.CreatedAt.Format("2006-01-02")
		ds := byDate[date]
		if ds == nil {
			ds = &dayStats{}
			byDate[date] = ds
		}
		ds.total++
		if run.Status == "passed" || run.Status == "completed" {
			ds.passed++
		}
	}

	// Sort dates
	dates := make([]string, 0, len(byDate))
	for d := range byDate {
		dates = append(dates, d)
	}
	sort.Strings(dates)

	trends := &HealthTrends{Range: rangeStr}
	totalAll, passedAll := 0, 0
	for _, d := range dates {
		ds := byDate[d]
		rate := float64(0)
		if ds.total > 0 {
			rate = float64(ds.passed) / float64(ds.total) * 100
		}
		trends.Data = append(trends.Data, HealthPoint{
			Timestamp: d + "T00:00:00Z", Date: d,
			TotalRuns: ds.total, PassedRuns: ds.passed,
			FailedRuns: ds.total - ds.passed, PassRate: rate,
		})
		totalAll += ds.total
		passedAll += ds.passed
	}

	if totalAll > 0 {
		trends.CurrentHealth = float64(passedAll) / float64(totalAll) * 100
		trends.WeeklyAverage = trends.CurrentHealth
		trends.MonthlyAverage = trends.CurrentHealth
	}
	if len(trends.Data) >= 2 {
		last := trends.Data[len(trends.Data)-1].PassRate
		prev := trends.Data[len(trends.Data)-2].PassRate
		if last >= prev {
			trends.Trend = "improving"
		} else {
			trends.Trend = "degrading"
		}
	} else {
		trends.Trend = "stable"
	}
	return trends, nil
}

// --- Helpers ---

func (s *AnalyticsService) loadRunResult(ctx context.Context, runID uuid.UUID) *model.RunResult {
	stored, err := s.runs.GetResult(ctx, runID)
	if err != nil || stored == nil || len(stored.Result) == 0 {
		return nil
	}
	var rr model.RunResult
	if json.Unmarshal(stored.Result, &rr) != nil {
		return nil
	}
	return &rr
}

func rangeToDays(r string) int {
	switch r {
	case "7d":
		return 7
	case "30d":
		return 30
	case "90d":
		return 90
	default:
		return 30
	}
}

func percentile(sorted []int64, p int) int64 {
	if len(sorted) == 0 {
		return 0
	}
	idx := len(sorted) * p / 100
	if idx >= len(sorted) {
		idx = len(sorted) - 1
	}
	return sorted[idx]
}
