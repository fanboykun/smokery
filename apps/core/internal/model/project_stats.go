package model

import (
	"time"

	"github.com/google/uuid"
)

type LastRunInfo struct {
	ID        uuid.UUID `json:"id"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	Duration  int64     `json:"duration_ms"`
}

type ProjectHealthStats struct {
	TotalRuns        int     `json:"total_runs"`
	PassedRuns       int     `json:"passed_runs"`
	FailedRuns       int     `json:"failed_runs"`
	EnvCount         int     `json:"env_count"`
	FlowCount        int     `json:"flow_count"`
	SuiteCount       int     `json:"suite_count"`
	HealthPercentage float64 `json:"health_percentage"`
}

type ProjectWithStats struct {
	ID          uuid.UUID           `json:"id"`
	Name        string              `json:"name"`
	Description string              `json:"description"`
	CreatedAt   time.Time           `json:"created_at"`
	UpdatedAt   time.Time           `json:"updated_at"`
	SpecCount   int                 `json:"spec_count"`
	LastRun     *LastRunInfo        `json:"last_run,omitempty"`
	Stats       ProjectHealthStats  `json:"stats"`
}
