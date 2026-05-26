package http

import (
	"context"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/google/uuid"

	"github.com/fanboykun/smokery/apps/core/internal/app"
)

// RegisterOperations registers operation list/update endpoints.
func RegisterOperations(api huma.API, svc *app.OperationService) {
	huma.Get(api, "/api/specs/{spec-id}/operations", func(ctx context.Context, in *SpecIDParam) (*OperationListOutput, error) {
		specID, err := uuid.Parse(in.SpecID)
		if err != nil {
			return nil, huma.Error400BadRequest("invalid spec id")
		}
		ops, err := svc.ListBySpec(ctx, specID)
		if err != nil {
			return nil, huma.Error500InternalServerError("failed to list operations", err)
		}
		return &OperationListOutput{Body: ops}, nil
	})

	huma.Put(api, "/api/operations/{id}/classification", func(ctx context.Context, in *UpdateClassificationInput) (*OperationOutput, error) {
		id, err := uuid.Parse(in.ID)
		if err != nil {
			return nil, huma.Error400BadRequest("invalid id")
		}
		op, err := svc.UpdateClassification(ctx, id, in.Body.Classification, in.Body.IsDestructive)
		if err != nil {
			return nil, huma.Error500InternalServerError("failed to update classification", err)
		}
		return &OperationOutput{Body: *op}, nil
	})

	huma.Register(api, huma.Operation{
		OperationID: "update-operation-overrides",
		Method:      http.MethodPut,
		Path:        "/api/operations/{id}/overrides",
		Summary:     "Update operation overrides",
	}, func(ctx context.Context, in *UpdateOverridesInput) (*OperationOutput, error) {
		id, err := uuid.Parse(in.ID)
		if err != nil {
			return nil, huma.Error400BadRequest("invalid id")
		}
		op, err := svc.UpdateOverrides(ctx, id, in.RawBody)
		if err != nil {
			return nil, huma.Error500InternalServerError("failed to update overrides", err)
		}
		return &OperationOutput{Body: *op}, nil
	})
}
