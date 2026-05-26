package model

import (
	"time"

	"github.com/google/uuid"
)

// --- Project ---

type Project struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// --- Spec (imported OpenAPI spec) ---

type Spec struct {
	ID        uuid.UUID `json:"id"`
	ProjectID uuid.UUID `json:"project_id"`
	Version   string    `json:"version"`
	Title     string    `json:"title"`
	Raw       []byte    `json:"raw"`
	Analysis  []byte    `json:"analysis"`
	CreatedAt time.Time `json:"created_at"`
}

// --- Operation (registered from a spec) ---

type Operation struct {
	ID             uuid.UUID `json:"id"`
	SpecID         uuid.UUID `json:"spec_id"`
	OperationID    string    `json:"operation_id"`
	Method         string    `json:"method"`
	Path           string    `json:"path"`
	Summary        string    `json:"summary"`
	Tags           []string  `json:"tags"`
	Classification string    `json:"classification"`
	IsDestructive  bool      `json:"is_destructive"`
	Overrides      []byte    `json:"overrides"`
	CreatedAt      time.Time `json:"created_at"`
}

// --- Run ---

type Run struct {
	ID         uuid.UUID  `json:"id"`
	ProjectID  uuid.UUID  `json:"project_id"`
	PlanID     *uuid.UUID `json:"plan_id,omitempty"`
	Status     string     `json:"status"`
	StartedAt  *time.Time `json:"started_at,omitempty"`
	FinishedAt *time.Time `json:"finished_at,omitempty"`
	CreatedAt  time.Time  `json:"created_at"`
}

// --- StoredRunResult (db row carrying the run's RunResult JSONB) ---

type StoredRunResult struct {
	ID        uuid.UUID `json:"id"`
	RunID     uuid.UUID `json:"run_id"`
	Result    []byte    `json:"result"`
	CreatedAt time.Time `json:"created_at"`
}

// --- Comment ---

type Comment struct {
	ID        uuid.UUID `json:"id"`
	RunID     uuid.UUID `json:"run_id"`
	Author    string    `json:"author"`
	Body      string    `json:"body"`
	CreatedAt time.Time `json:"created_at"`
}

// --- Artifact ---

type Artifact struct {
	ID        uuid.UUID `json:"id"`
	RunID     uuid.UUID `json:"run_id"`
	Type      string    `json:"type"`
	Path      string    `json:"path"`
	CreatedAt time.Time `json:"created_at"`
}
