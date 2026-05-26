package report

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/fanboykun/smokery/apps/api/internal/model"
)

// DebugView returns a backend-developer-focused report
type DebugView struct {
	RunID    string        `json:"run_id"`
	Status   string        `json:"status"`
	Duration int64         `json:"duration_ms"`
	Failures []FailureInfo `json:"failures,omitempty"`
	Traces   []TraceInfo   `json:"traces,omitempty"`
}

type FailureInfo struct {
	Step       string `json:"step"`
	Assertion  string `json:"assertion"`
	Message    string `json:"message"`
	RequestURL string `json:"request_url"`
	Status     int    `json:"status"`
}

type TraceInfo struct {
	Step      string `json:"step"`
	TraceID   string `json:"trace_id,omitempty"`
	RequestID string `json:"request_id,omitempty"`
}

func GenerateDebugView(result *model.RunResult) *DebugView {
	view := &DebugView{RunID: result.RunID, Status: result.Status, Duration: result.Duration}
	for _, f := range result.Flows {
		for _, s := range f.Steps {
			collectStepInfo(s, view)
		}
	}
	for _, su := range result.Suites {
		for _, c := range su.Cases {
			collectStepInfo(c.Step, view)
		}
	}
	return view
}

func collectStepInfo(s model.StepResult, view *DebugView) {
	if s.Response.TraceID != "" || s.Response.RequestID != "" {
		view.Traces = append(view.Traces, TraceInfo{Step: s.Name, TraceID: s.Response.TraceID, RequestID: s.Response.RequestID})
	}
	for _, a := range s.Assertions {
		if !a.Passed {
			view.Failures = append(view.Failures, FailureInfo{
				Step: s.Name, Assertion: a.Type, Message: a.Message,
				RequestURL: s.Request.URL, Status: s.Response.Status,
			})
		}
	}
}

// CISummary returns a CI-friendly summary
type CISummary struct {
	Status   string   `json:"status"`
	Total    int      `json:"total"`
	Passed   int      `json:"passed"`
	Failed   int      `json:"failed"`
	Duration int64    `json:"duration_ms"`
	Failures []string `json:"failures,omitempty"`
}

func GenerateCISummary(result *model.RunResult) *CISummary {
	summary := &CISummary{Status: result.Status, Duration: result.Duration}
	for _, f := range result.Flows {
		for _, s := range f.Steps {
			summary.Total++
			if s.Status == "passed" {
				summary.Passed++
			} else {
				summary.Failed++
				summary.Failures = append(summary.Failures, s.Name)
			}
		}
	}
	for _, su := range result.Suites {
		for _, c := range su.Cases {
			summary.Total++
			if c.Step.Status == "passed" {
				summary.Passed++
			} else {
				summary.Failed++
				summary.Failures = append(summary.Failures, c.OperationID+":"+c.CaseType)
			}
		}
	}
	return summary
}

// GenerateJSONArtifact returns the full result as formatted JSON bytes
func GenerateJSONArtifact(result *model.RunResult) []byte {
	b, _ := json.MarshalIndent(result, "", "  ")
	return b
}

// GenerateMermaidDiagram generates a sequence diagram for flows
func GenerateMermaidDiagram(result *model.RunResult) string {
	var sb strings.Builder
	sb.WriteString("sequenceDiagram\n")
	sb.WriteString("    participant Client\n")
	sb.WriteString("    participant API\n")
	for _, f := range result.Flows {
		sb.WriteString(fmt.Sprintf("    Note over Client,API: Flow: %s\n", f.Name))
		for _, s := range f.Steps {
			arrow := "->>"
			if s.Status != "passed" {
				arrow = "-x"
			}
			sb.WriteString(fmt.Sprintf("    Client%sAPI: %s %s\n", arrow, s.Request.Method, s.Name))
			sb.WriteString(fmt.Sprintf("    API-->>Client: %d\n", s.Response.Status))
		}
	}
	return sb.String()
}
