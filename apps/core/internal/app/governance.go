package app

import (
	"context"
	"sync"
	"time"

	"github.com/google/uuid"

	"github.com/fanboykun/smokery/apps/core/internal/model"
)

// GovernanceService provides members, audit, webhooks, notifications.
// Uses in-memory storage as a placeholder until DB tables are added.
type GovernanceService struct {
	mu            sync.RWMutex
	members       []model.ProjectMember
	auditLog      []model.AuditLogEntry
	webhooks      []model.Webhook
	notifications []model.NotificationRule
}

func NewGovernanceService() *GovernanceService {
	return &GovernanceService{}
}

// --- Members ---

func (s *GovernanceService) ListMembers(_ context.Context, projectID uuid.UUID) []model.ProjectMember {
	s.mu.RLock()
	defer s.mu.RUnlock()
	var out []model.ProjectMember
	for _, m := range s.members {
		if m.ProjectID == projectID {
			out = append(out, m)
		}
	}
	return out
}

func (s *GovernanceService) AddMember(_ context.Context, projectID uuid.UUID, userName, userEmail, role, addedBy string) *model.ProjectMember {
	s.mu.Lock()
	defer s.mu.Unlock()
	m := model.ProjectMember{
		ID: uuid.New(), ProjectID: projectID,
		UserID: uuid.New().String(), UserName: userName, UserEmail: userEmail,
		Role: role, AddedAt: time.Now(), AddedBy: addedBy,
	}
	s.members = append(s.members, m)
	return &m
}

// --- Audit Log ---

func (s *GovernanceService) GetAuditLog(_ context.Context, projectID uuid.UUID) []model.AuditLogEntry {
	s.mu.RLock()
	defer s.mu.RUnlock()
	var out []model.AuditLogEntry
	for _, e := range s.auditLog {
		if e.ProjectID == projectID {
			out = append(out, e)
		}
	}
	return out
}

func (s *GovernanceService) RecordAudit(_ context.Context, entry model.AuditLogEntry) {
	s.mu.Lock()
	defer s.mu.Unlock()
	entry.ID = uuid.New()
	entry.Timestamp = time.Now()
	s.auditLog = append(s.auditLog, entry)
}

// --- Webhooks ---

func (s *GovernanceService) ListWebhooks(_ context.Context, projectID uuid.UUID) []model.Webhook {
	s.mu.RLock()
	defer s.mu.RUnlock()
	var out []model.Webhook
	for _, w := range s.webhooks {
		if w.ProjectID == projectID {
			out = append(out, w)
		}
	}
	return out
}

func (s *GovernanceService) CreateWebhook(_ context.Context, projectID uuid.UUID, name, url string, events []string) *model.Webhook {
	s.mu.Lock()
	defer s.mu.Unlock()
	w := model.Webhook{
		ID: uuid.New(), ProjectID: projectID,
		Name: name, URL: url, Events: events,
		IsActive: true, CreatedAt: time.Now(),
	}
	s.webhooks = append(s.webhooks, w)
	return &w
}

// --- Notifications ---

func (s *GovernanceService) ListNotifications(_ context.Context, projectID uuid.UUID) []model.NotificationRule {
	s.mu.RLock()
	defer s.mu.RUnlock()
	var out []model.NotificationRule
	for _, n := range s.notifications {
		if n.ProjectID == projectID {
			out = append(out, n)
		}
	}
	return out
}

func (s *GovernanceService) CreateNotification(_ context.Context, projectID uuid.UUID, name, channel string, triggers []model.NotificationTrigger) *model.NotificationRule {
	s.mu.Lock()
	defer s.mu.Unlock()
	n := model.NotificationRule{
		ID: uuid.New(), ProjectID: projectID,
		Name: name, Channel: channel, Triggers: triggers,
		IsActive: true, CreatedAt: time.Now(),
	}
	s.notifications = append(s.notifications, n)
	return &n
}
