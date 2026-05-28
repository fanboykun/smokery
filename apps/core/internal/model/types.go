package model

import (
	"encoding/json"
	"time"
)

// --- Project Config ---

type ProjectConfig struct {
	Environments []Environment `json:"environments"`
	AuthProfiles []AuthProfile `json:"auth_profiles"`
	Flows        []Flow        `json:"flows"`
	Suites       []Suite       `json:"suites"`
}

type Environment struct {
	ID      string            `json:"id"`
	Name    string            `json:"name"`
	BaseURL string            `json:"base_url"`
	Headers map[string]string `json:"headers,omitempty"`
}

type AuthProfile struct {
	ID     string         `json:"id"`
	Name   string         `json:"name"`
	Type   string         `json:"type"` // bearer, basic, api_key, custom
	Config map[string]any `json:"config"`
}

// --- Flow ---

type Flow struct {
	ID          string     `json:"id"`
	Name        string     `json:"name"`
	Description string     `json:"description,omitempty"`
	Environment string     `json:"environment"`    // environment ID ref
	Auth        string     `json:"auth,omitempty"` // auth profile ID ref
	Steps       []FlowStep `json:"steps"`
	Cleanup     []FlowStep `json:"cleanup,omitempty"`
}

type FlowStep struct {
	Name        string            `json:"name"`
	OperationID string            `json:"operation_id"`
	Params      map[string]any    `json:"params,omitempty"`
	Body        any               `json:"body,omitempty"`
	Headers     map[string]string `json:"headers,omitempty"`
	Captures    []Capture         `json:"captures,omitempty"`
	Assertions  []Assertion       `json:"assertions,omitempty"`
}

// --- Suite ---

type Suite struct {
	ID          string        `json:"id"`
	Name        string        `json:"name"`
	Description string        `json:"description,omitempty"`
	Environment string        `json:"environment"`
	Auth        string        `json:"auth,omitempty"`
	Selector    SuiteSelector `json:"selector"`
	Strategy    SuiteStrategy `json:"strategy"`
}

type SuiteSelector struct {
	Tags            []string `json:"tags,omitempty"`
	Classifications []string `json:"classifications,omitempty"` // list, read, etc.
	Paths           []string `json:"paths,omitempty"`
	Exclude         []string `json:"exclude,omitempty"`
}

type SuiteStrategy struct {
	DefaultList        bool   `json:"default_list"`
	Pagination         bool   `json:"pagination"`
	SearchFromResponse bool   `json:"search_from_response"`
	EnumFilters        bool   `json:"enum_filters"`
	EmptyResultPolicy  string `json:"empty_result_policy"` // allow, warn, fail
	MaxCasesPerOp      int    `json:"max_cases_per_op,omitempty"`
}

// --- Capture & Assertion ---

type Capture struct {
	Name   string `json:"name"`
	Source string `json:"source"` // body, header
	Path   string `json:"path"`   // JSONPath or header name
}

type Assertion struct {
	Type     string `json:"type"` // status, jsonpath, schema, list_shape, not_empty
	Expected any    `json:"expected,omitempty"`
	Path     string `json:"path,omitempty"`
}

// --- Smoke Plan (compiler output) ---

type SmokePlan struct {
	ID          string       `json:"id"`
	ProjectID   string       `json:"project_id"`
	Environment Environment  `json:"environment"`
	Auth        *AuthProfile `json:"auth,omitempty"`
	FlowPlans   []FlowPlan   `json:"flow_plans,omitempty"`
	SuitePlans  []SuitePlan  `json:"suite_plans,omitempty"`
	CompiledAt  time.Time    `json:"compiled_at"`
}

type FlowPlan struct {
	FlowID  string        `json:"flow_id"`
	Name    string        `json:"name"`
	Steps   []PlannedStep `json:"steps"`
	Cleanup []PlannedStep `json:"cleanup,omitempty"`
}

type SuitePlan struct {
	SuiteID string        `json:"suite_id"`
	Name    string        `json:"name"`
	Cases   []PlannedCase `json:"cases"`
}

type PlannedStep struct {
	Name           string            `json:"name"`
	Method         string            `json:"method"`
	Path           string            `json:"path"`
	Params         map[string]any    `json:"params,omitempty"`
	Body           any               `json:"body,omitempty"`
	Headers        map[string]string `json:"headers,omitempty"`
	Captures       []Capture         `json:"captures,omitempty"`
	Assertions     []Assertion       `json:"assertions,omitempty"`
	ResponseSchema json.RawMessage   `json:"response_schema,omitempty"`
}

type PlannedCase struct {
	OperationID       string        `json:"operation_id"`
	CaseType          string        `json:"case_type"` // default_list, pagination, search, enum_filter
	Step              PlannedStep   `json:"step"`
	Steps             []PlannedStep `json:"steps,omitempty"` // multi-step cases (e.g. search-from-response)
	EmptyResultPolicy string        `json:"empty_result_policy,omitempty"` // allow, warn, fail
}

// --- Run Result ---

type RunResult struct {
	RunID      string        `json:"run_id"`
	Status     string        `json:"status"` // passed, failed, error
	Flows      []FlowResult  `json:"flows,omitempty"`
	Suites     []SuiteResult `json:"suites,omitempty"`
	StartedAt  time.Time     `json:"started_at"`
	FinishedAt time.Time     `json:"finished_at"`
	Duration   int64         `json:"duration_ms"`
}

type FlowResult struct {
	FlowID  string       `json:"flow_id"`
	Name    string       `json:"name"`
	Status  string       `json:"status"`
	Steps   []StepResult `json:"steps"`
	Cleanup []StepResult `json:"cleanup,omitempty"`
}

type SuiteResult struct {
	SuiteID string       `json:"suite_id"`
	Name    string       `json:"name"`
	Status  string       `json:"status"`
	Cases   []CaseResult `json:"cases"`
}

type StepResult struct {
	Name       string            `json:"name"`
	Status     string            `json:"status"`
	Request    RequestMeta       `json:"request"`
	Response   ResponseMeta      `json:"response"`
	Captures   map[string]any    `json:"captures,omitempty"`
	Assertions []AssertionResult `json:"assertions"`
	Error      string            `json:"error,omitempty"`
	Duration   int64             `json:"duration_ms"`
}

type CaseResult struct {
	OperationID string     `json:"operation_id"`
	CaseType    string     `json:"case_type"`
	Step        StepResult `json:"step"`
}

type RequestMeta struct {
	Method  string            `json:"method"`
	URL     string            `json:"url"`
	Headers map[string]string `json:"headers,omitempty"`
	Body    any               `json:"body,omitempty"`
}

type ResponseMeta struct {
	Status    int               `json:"status"`
	Headers   map[string]string `json:"headers,omitempty"`
	Body      any               `json:"body,omitempty"`
	TraceID   string            `json:"trace_id,omitempty"`
	RequestID string            `json:"request_id,omitempty"`
}

type AssertionResult struct {
	Type     string `json:"type"`
	Expected any    `json:"expected,omitempty"`
	Actual   any    `json:"actual,omitempty"`
	Passed   bool   `json:"passed"`
	Message  string `json:"message,omitempty"`
}

// --- Compiler Error ---

type CompileError struct {
	Stage    string `json:"stage"`
	Path     string `json:"path"`
	Message  string `json:"message"`
	Severity string `json:"severity"` // error, warning
	Entity   string `json:"entity,omitempty"`
}
