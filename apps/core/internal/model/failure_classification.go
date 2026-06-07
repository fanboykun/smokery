package model

import (
	"time"

	"github.com/google/uuid"
)

type FailureClassification struct {
	ID             uuid.UUID `json:"id"`
	RunID          uuid.UUID `json:"run_id"`
	Classification string    `json:"classification"`
	Assignee       string    `json:"assignee,omitempty"`
	Note           string    `json:"note,omitempty"`
	Author         string    `json:"author"`
	CreatedAt      time.Time `json:"created_at"`
}
