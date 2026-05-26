package postgres

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/fanboykun/smokery/apps/core/internal/adapter/postgres/db"
	"github.com/fanboykun/smokery/apps/core/internal/model"
	"github.com/fanboykun/smokery/apps/core/internal/port"
)

type OperationRepo struct{ q *db.Queries }

func NewOperationRepo(pool *pgxpool.Pool) *OperationRepo {
	return &OperationRepo{q: db.New(pool)}
}

var _ port.OperationRepo = (*OperationRepo)(nil)

func (r *OperationRepo) Create(ctx context.Context, in model.Operation) (*model.Operation, error) {
	o, err := r.q.CreateOperation(ctx, db.CreateOperationParams{
		SpecID:         toPgUUID(in.SpecID),
		OperationID:    in.OperationID,
		Method:         in.Method,
		Path:           in.Path,
		Summary:        in.Summary,
		Tags:           marshalJSON(in.Tags),
		Classification: in.Classification,
		IsDestructive:  in.IsDestructive,
	})
	if err != nil {
		return nil, err
	}
	return operationToModel(o), nil
}

func (r *OperationRepo) ListBySpec(ctx context.Context, specID uuid.UUID) ([]model.Operation, error) {
	ops, err := r.q.ListOperationsBySpec(ctx, toPgUUID(specID))
	if err != nil {
		return nil, err
	}
	out := make([]model.Operation, len(ops))
	for i, o := range ops {
		out[i] = *operationToModel(o)
	}
	return out, nil
}

func (r *OperationRepo) UpdateClassification(ctx context.Context, id uuid.UUID, classification string, isDestructive bool) (*model.Operation, error) {
	o, err := r.q.UpdateOperationClassification(ctx, db.UpdateOperationClassificationParams{
		ID:             toPgUUID(id),
		Classification: classification,
		IsDestructive:  isDestructive,
	})
	if err != nil {
		return nil, err
	}
	return operationToModel(o), nil
}

func (r *OperationRepo) UpdateOverrides(ctx context.Context, id uuid.UUID, overrides []byte) (*model.Operation, error) {
	o, err := r.q.UpdateOperationOverrides(ctx, db.UpdateOperationOverridesParams{
		ID:        toPgUUID(id),
		Overrides: overrides,
	})
	if err != nil {
		return nil, err
	}
	return operationToModel(o), nil
}

func operationToModel(o db.Operation) *model.Operation {
	return &model.Operation{
		ID:             fromPgUUID(o.ID),
		SpecID:         fromPgUUID(o.SpecID),
		OperationID:    o.OperationID,
		Method:         o.Method,
		Path:           o.Path,
		Summary:        o.Summary,
		Tags:           unmarshalStrings(o.Tags),
		Classification: o.Classification,
		IsDestructive:  o.IsDestructive,
		Overrides:      o.Overrides,
		CreatedAt:      fromPgTimestamptz(o.CreatedAt),
	}
}
