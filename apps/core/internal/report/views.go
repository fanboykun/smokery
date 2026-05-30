package report

import "github.com/fanboykun/smokery/apps/core/internal/model"

// --- Contract View ---

type ContractViolation struct {
	ID             string `json:"id"`
	OperationID    string `json:"operation_id"`
	ViolationType  string `json:"violation_type"`
	Severity       string `json:"severity"`
	Message        string `json:"message"`
	ExpectedSchema any    `json:"expected_schema,omitempty"`
	ActualValue    any    `json:"actual_value,omitempty"`
	Location       string `json:"location"`
}

type ContractView struct {
	RunID            string              `json:"run_id"`
	TotalViolations  int                 `json:"total_violations"`
	Errors           []ContractViolation `json:"errors"`
	Warnings         []ContractViolation `json:"warnings"`
	ComplianceScore  float64             `json:"compliance_score"`
	PassedAssertions int                 `json:"passed_assertions"`
	FailedAssertions int                 `json:"failed_assertions"`
}

func GenerateContractView(result *model.RunResult) *ContractView {
	view := &ContractView{RunID: result.RunID}
	var total, passed, failed int

	processStep := func(s model.StepResult, opID string) {
		for _, a := range s.Assertions {
			total++
			if a.Passed {
				passed++
			} else {
				failed++
				v := ContractViolation{
					ID:            opID + ":" + a.Type,
					OperationID:   opID,
					ViolationType: a.Type,
					Message:       a.Message,
					Location:      "response.body",
				}
				if a.Type == "schema" || a.Type == "status" {
					v.Severity = "error"
					v.ExpectedSchema = a.Expected
					v.ActualValue = a.Actual
					view.Errors = append(view.Errors, v)
				} else {
					v.Severity = "warning"
					view.Warnings = append(view.Warnings, v)
				}
			}
		}
	}

	for _, f := range result.Flows {
		for _, s := range f.Steps {
			processStep(s, s.Name)
		}
	}
	for _, su := range result.Suites {
		for _, c := range su.Cases {
			processStep(c.Step, c.OperationID)
		}
	}

	view.PassedAssertions = passed
	view.FailedAssertions = failed
	view.TotalViolations = len(view.Errors) + len(view.Warnings)
	if total > 0 {
		view.ComplianceScore = float64(passed) / float64(total) * 100
	}
	return view
}

// --- Analyst View ---

type RootCause struct {
	Cause              string   `json:"cause"`
	Impact             int      `json:"impact"`
	AffectedOperations []string `json:"affected_operations"`
}

type Recommendation struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Priority    string `json:"priority"`
}

type TimelineInsight struct {
	Timestamp   string `json:"timestamp"`
	OperationID string `json:"operation_id"`
	Status      string `json:"status"`
	DurationMs  int64  `json:"duration_ms"`
	Error       string `json:"error,omitempty"`
}

type AnalystView struct {
	RunID             string            `json:"run_id"`
	Summary           string            `json:"summary"`
	RootCauses        []RootCause       `json:"root_causes"`
	Recommendations   []Recommendation  `json:"recommendations"`
	TimelineInsights  []TimelineInsight  `json:"timeline_insights"`
}

func GenerateAnalystView(result *model.RunResult) *AnalystView {
	view := &AnalystView{RunID: result.RunID}

	// Collect timeline and group failures by cause
	causeMap := map[string][]string{}

	processStep := func(s model.StepResult, opID string) {
		view.TimelineInsights = append(view.TimelineInsights, TimelineInsight{
			Timestamp:   result.StartedAt.Add(0).Format("2006-01-02T15:04:05Z"),
			OperationID: opID,
			Status:      s.Status,
			DurationMs:  s.Duration,
			Error:       s.Error,
		})
		if s.Status != "passed" {
			cause := classifyCause(s)
			causeMap[cause] = append(causeMap[cause], opID)
		}
	}

	for _, f := range result.Flows {
		for _, s := range f.Steps {
			processStep(s, s.Name)
		}
	}
	for _, su := range result.Suites {
		for _, c := range su.Cases {
			processStep(c.Step, c.OperationID)
		}
	}

	// Build root causes
	totalFailures := 0
	for _, ops := range causeMap {
		totalFailures += len(ops)
	}
	for cause, ops := range causeMap {
		impact := 0
		if totalFailures > 0 {
			impact = len(ops) * 100 / totalFailures
		}
		view.RootCauses = append(view.RootCauses, RootCause{Cause: cause, Impact: impact, AffectedOperations: ops})
	}

	// Generate summary and recommendations
	view.Summary = generateSummary(result, causeMap)
	view.Recommendations = generateRecommendations(causeMap)
	return view
}

func classifyCause(s model.StepResult) string {
	if s.Response.Status == 0 {
		return "Network Timeout"
	}
	if s.Response.Status == 429 {
		return "Rate Limiting"
	}
	if s.Response.Status == 401 || s.Response.Status == 403 {
		return "Auth Failure"
	}
	if s.Response.Status >= 500 {
		return "Server Error"
	}
	for _, a := range s.Assertions {
		if !a.Passed && a.Type == "schema" {
			return "Schema Violation"
		}
	}
	return "Contract Violation"
}

func generateSummary(result *model.RunResult, causes map[string][]string) string {
	if len(causes) == 0 {
		return "All tests passed successfully."
	}
	s := "Run had failures: "
	for cause, ops := range causes {
		s += cause + " (" + itoa(len(ops)) + " ops), "
	}
	return s[:len(s)-2]
}

func generateRecommendations(causes map[string][]string) []Recommendation {
	var recs []Recommendation
	for cause := range causes {
		switch cause {
		case "Network Timeout":
			recs = append(recs, Recommendation{Title: "Increase Timeout", Description: "Consider increasing API gateway timeout or adding retry logic", Priority: "high"})
		case "Rate Limiting":
			recs = append(recs, Recommendation{Title: "Add Rate Limit Handling", Description: "Implement backoff or reduce request concurrency", Priority: "high"})
		case "Auth Failure":
			recs = append(recs, Recommendation{Title: "Check Auth Configuration", Description: "Verify auth tokens/credentials are valid and not expired", Priority: "high"})
		case "Server Error":
			recs = append(recs, Recommendation{Title: "Investigate Server Errors", Description: "Check server logs for 5xx root causes", Priority: "high"})
		default:
			recs = append(recs, Recommendation{Title: "Fix " + cause, Description: "Review failing assertions and update spec or implementation", Priority: "medium"})
		}
	}
	return recs
}

func itoa(n int) string {
	if n == 0 {
		return "0"
	}
	s := ""
	for n > 0 {
		s = string(rune('0'+n%10)) + s
		n /= 10
	}
	return s
}

// --- QA View ---

type QABlocker struct {
	OperationID string `json:"operation_id"`
	Issue       string `json:"issue"`
	Severity    string `json:"severity"`
}

type CoverageSummary struct {
	TotalOperations  int     `json:"total_operations"`
	TestedOperations int     `json:"tested_operations"`
	CoveragePercent  float64 `json:"coverage_percentage"`
}

type QAView struct {
	RunID           string          `json:"run_id"`
	Status          string          `json:"status"`
	TotalTests      int             `json:"total_tests"`
	PassedTests     int             `json:"passed_tests"`
	FailedTests     int             `json:"failed_tests"`
	PassRate        float64         `json:"pass_rate"`
	FlakyTests      []string        `json:"flaky_tests"`
	CoverageSummary CoverageSummary `json:"coverage_summary"`
	Blockers        []QABlocker     `json:"blockers"`
}

func GenerateQAView(result *model.RunResult) *QAView {
	view := &QAView{RunID: result.RunID, Status: result.Status}
	opsSet := map[string]bool{}

	processStep := func(s model.StepResult, opID string) {
		view.TotalTests++
		opsSet[opID] = true
		if s.Status == "passed" {
			view.PassedTests++
		} else {
			view.FailedTests++
			view.Blockers = append(view.Blockers, QABlocker{
				OperationID: opID,
				Issue:       s.Error,
				Severity:    "high",
			})
		}
	}

	for _, f := range result.Flows {
		for _, s := range f.Steps {
			processStep(s, s.Name)
		}
	}
	for _, su := range result.Suites {
		for _, c := range su.Cases {
			processStep(c.Step, c.OperationID)
		}
	}

	if view.TotalTests > 0 {
		view.PassRate = float64(view.PassedTests) / float64(view.TotalTests) * 100
	}
	view.CoverageSummary.TestedOperations = len(opsSet)
	view.CoverageSummary.TotalOperations = len(opsSet) // best we can do without spec context
	if view.CoverageSummary.TotalOperations > 0 {
		view.CoverageSummary.CoveragePercent = 100
	}
	return view
}
