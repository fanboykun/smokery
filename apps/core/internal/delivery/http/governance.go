package http

import (
	"context"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/google/uuid"

	"github.com/fanboykun/smokery/apps/core/internal/app"
	"github.com/fanboykun/smokery/apps/core/internal/model"
)

// --- Types ---

type MemberListOutput struct{ Body []model.ProjectMember }
type MemberOutput struct{ Body model.ProjectMember }
type AuditLogOutput struct{ Body []model.AuditLogEntry }
type WebhookListOutput struct{ Body []model.Webhook }
type WebhookOutput struct{ Body model.Webhook }
type NotificationListOutput struct{ Body []model.NotificationRule }
type NotificationOutput struct{ Body model.NotificationRule }

type AddMemberInput struct {
	ProjectIDParam
	Body struct {
		UserName  string `json:"user_name" minLength:"1"`
		UserEmail string `json:"user_email" minLength:"1"`
		Role      string `json:"role" minLength:"1"`
		AddedBy   string `json:"added_by" minLength:"1"`
	}
}

type CreateWebhookInput struct {
	ProjectIDParam
	Body struct {
		Name   string   `json:"name" minLength:"1"`
		URL    string   `json:"url" minLength:"1"`
		Events []string `json:"events"`
	}
}

type CreateNotificationInput struct {
	ProjectIDParam
	Body struct {
		Name     string                      `json:"name" minLength:"1"`
		Channel  string                      `json:"channel" minLength:"1"`
		Triggers []model.NotificationTrigger `json:"triggers"`
	}
}

// --- Registration ---

func RegisterGovernance(api huma.API, svc *app.GovernanceService) {
	// Members
	huma.Get(api, "/api/projects/{project-id}/members", func(ctx context.Context, in *ProjectIDParam) (*MemberListOutput, error) {
		projectID, err := uuid.Parse(in.ProjectID)
		if err != nil {
			return nil, ErrBadRequest("members", "invalid project id")
		}
		return &MemberListOutput{Body: svc.ListMembers(ctx, projectID)}, nil
	})

	huma.Register(api, huma.Operation{
		OperationID: "add-member",
		Method:      http.MethodPost,
		Path:        "/api/projects/{project-id}/members",
		Summary:     "Add a project member",
	}, func(ctx context.Context, in *AddMemberInput) (*MemberOutput, error) {
		projectID, err := uuid.Parse(in.ProjectID)
		if err != nil {
			return nil, ErrBadRequest("members", "invalid project id")
		}
		m := svc.AddMember(ctx, projectID, in.Body.UserName, in.Body.UserEmail, in.Body.Role, in.Body.AddedBy)
		return &MemberOutput{Body: *m}, nil
	})

	// Audit Log
	huma.Get(api, "/api/projects/{project-id}/audit-log", func(ctx context.Context, in *ProjectIDParam) (*AuditLogOutput, error) {
		projectID, err := uuid.Parse(in.ProjectID)
		if err != nil {
			return nil, ErrBadRequest("audit-log", "invalid project id")
		}
		return &AuditLogOutput{Body: svc.GetAuditLog(ctx, projectID)}, nil
	})

	// Webhooks
	huma.Get(api, "/api/projects/{project-id}/webhooks", func(ctx context.Context, in *ProjectIDParam) (*WebhookListOutput, error) {
		projectID, err := uuid.Parse(in.ProjectID)
		if err != nil {
			return nil, ErrBadRequest("webhooks", "invalid project id")
		}
		return &WebhookListOutput{Body: svc.ListWebhooks(ctx, projectID)}, nil
	})

	huma.Register(api, huma.Operation{
		OperationID: "create-webhook",
		Method:      http.MethodPost,
		Path:        "/api/projects/{project-id}/webhooks",
		Summary:     "Create a webhook",
	}, func(ctx context.Context, in *CreateWebhookInput) (*WebhookOutput, error) {
		projectID, err := uuid.Parse(in.ProjectID)
		if err != nil {
			return nil, ErrBadRequest("webhooks", "invalid project id")
		}
		w := svc.CreateWebhook(ctx, projectID, in.Body.Name, in.Body.URL, in.Body.Events)
		return &WebhookOutput{Body: *w}, nil
	})

	// Notifications
	huma.Get(api, "/api/projects/{project-id}/notifications", func(ctx context.Context, in *ProjectIDParam) (*NotificationListOutput, error) {
		projectID, err := uuid.Parse(in.ProjectID)
		if err != nil {
			return nil, ErrBadRequest("notifications", "invalid project id")
		}
		return &NotificationListOutput{Body: svc.ListNotifications(ctx, projectID)}, nil
	})

	huma.Register(api, huma.Operation{
		OperationID: "create-notification",
		Method:      http.MethodPost,
		Path:        "/api/projects/{project-id}/notifications",
		Summary:     "Create a notification rule",
	}, func(ctx context.Context, in *CreateNotificationInput) (*NotificationOutput, error) {
		projectID, err := uuid.Parse(in.ProjectID)
		if err != nil {
			return nil, ErrBadRequest("notifications", "invalid project id")
		}
		n := svc.CreateNotification(ctx, projectID, in.Body.Name, in.Body.Channel, in.Body.Triggers)
		return &NotificationOutput{Body: *n}, nil
	})
}
