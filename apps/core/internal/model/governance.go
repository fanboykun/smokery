package model

import (
	"time"

	"github.com/google/uuid"
)

type ProjectMember struct {
	ID        uuid.UUID `json:"id"`
	ProjectID uuid.UUID `json:"project_id"`
	UserID    string    `json:"user_id"`
	UserName  string    `json:"user_name"`
	UserEmail string    `json:"user_email"`
	Role      string    `json:"role"`
	AddedAt   time.Time `json:"added_at"`
	AddedBy   string    `json:"added_by"`
}

type AuditLogEntry struct {
	ID           uuid.UUID         `json:"id"`
	ProjectID    uuid.UUID         `json:"project_id"`
	Action       string            `json:"action"`
	ResourceType string            `json:"resource_type"`
	ResourceID   string            `json:"resource_id"`
	ActorID      string            `json:"actor_id"`
	ActorName    string            `json:"actor_name"`
	Timestamp    time.Time         `json:"timestamp"`
	Changes      []AuditChange     `json:"changes,omitempty"`
}

type AuditChange struct {
	Field    string `json:"field"`
	OldValue string `json:"old_value"`
	NewValue string `json:"new_value"`
}

type Webhook struct {
	ID              uuid.UUID `json:"id"`
	ProjectID       uuid.UUID `json:"project_id"`
	Name            string    `json:"name"`
	URL             string    `json:"url"`
	Events          []string  `json:"events"`
	IsActive        bool      `json:"is_active"`
	CreatedAt       time.Time `json:"created_at"`
	LastTriggeredAt *time.Time `json:"last_triggered_at,omitempty"`
}

type NotificationRule struct {
	ID        uuid.UUID              `json:"id"`
	ProjectID uuid.UUID              `json:"project_id"`
	Name      string                 `json:"name"`
	Channel   string                 `json:"channel"`
	Config    map[string]interface{} `json:"channel_config"`
	Triggers  []NotificationTrigger  `json:"triggers"`
	IsActive  bool                   `json:"is_active"`
	CreatedAt time.Time              `json:"created_at"`
}

type NotificationTrigger struct {
	Event      string                   `json:"event"`
	Conditions []NotificationCondition  `json:"conditions,omitempty"`
}

type NotificationCondition struct {
	Field    string `json:"field"`
	Operator string `json:"operator"`
	Value    string `json:"value"`
}
