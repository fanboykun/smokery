// Package sqlite implements all repo ports using a single SQLite database.
// Uses modernc.org/sqlite (pure Go, no CGO).
package sqlite

import (
	"context"
	"database/sql"
	"encoding/json"
	"time"

	"github.com/google/uuid"
	_ "modernc.org/sqlite"

	"github.com/fanboykun/smokery/apps/core/internal/model"
	"github.com/fanboykun/smokery/apps/core/internal/port"
)

const schema = `
CREATE TABLE IF NOT EXISTS projects (
	id TEXT PRIMARY KEY,
	name TEXT NOT NULL,
	description TEXT NOT NULL DEFAULT '',
	created_at TEXT NOT NULL DEFAULT (datetime('now')),
	updated_at TEXT NOT NULL DEFAULT (datetime('now'))
);

CREATE TABLE IF NOT EXISTS specs (
	id TEXT PRIMARY KEY,
	project_id TEXT NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
	version TEXT NOT NULL,
	title TEXT NOT NULL,
	raw BLOB,
	analysis BLOB,
	created_at TEXT NOT NULL DEFAULT (datetime('now'))
);

CREATE TABLE IF NOT EXISTS operations (
	id TEXT PRIMARY KEY,
	spec_id TEXT NOT NULL REFERENCES specs(id) ON DELETE CASCADE,
	operation_id TEXT NOT NULL,
	method TEXT NOT NULL,
	path TEXT NOT NULL,
	summary TEXT NOT NULL DEFAULT '',
	tags TEXT NOT NULL DEFAULT '[]',
	classification TEXT NOT NULL DEFAULT 'read',
	is_destructive INTEGER NOT NULL DEFAULT 0,
	overrides BLOB,
	created_at TEXT NOT NULL DEFAULT (datetime('now'))
);

CREATE TABLE IF NOT EXISTS runs (
	id TEXT PRIMARY KEY,
	project_id TEXT NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
	plan_id TEXT,
	status TEXT NOT NULL DEFAULT 'pending',
	started_at TEXT,
	finished_at TEXT,
	created_at TEXT NOT NULL DEFAULT (datetime('now'))
);

CREATE TABLE IF NOT EXISTS run_results (
	id TEXT PRIMARY KEY,
	run_id TEXT NOT NULL UNIQUE REFERENCES runs(id) ON DELETE CASCADE,
	result BLOB NOT NULL,
	created_at TEXT NOT NULL DEFAULT (datetime('now'))
);

CREATE TABLE IF NOT EXISTS comments (
	id TEXT PRIMARY KEY,
	run_id TEXT NOT NULL REFERENCES runs(id) ON DELETE CASCADE,
	author TEXT NOT NULL,
	body TEXT NOT NULL,
	created_at TEXT NOT NULL DEFAULT (datetime('now'))
);

CREATE TABLE IF NOT EXISTS artifacts (
	id TEXT PRIMARY KEY,
	run_id TEXT NOT NULL REFERENCES runs(id) ON DELETE CASCADE,
	type TEXT NOT NULL,
	path TEXT NOT NULL,
	created_at TEXT NOT NULL DEFAULT (datetime('now'))
);
`

// DB wraps a sql.DB with SQLite-specific initialization.
type DB struct {
	*sql.DB
}

// Open opens (or creates) a SQLite database at the given path and applies the schema.
func Open(path string) (*DB, error) {
	db, err := sql.Open("sqlite", path+"?_journal_mode=WAL&_busy_timeout=5000&_foreign_keys=on")
	if err != nil {
		return nil, err
	}
	db.SetMaxOpenConns(1) // SQLite doesn't support concurrent writes
	if _, err := db.Exec(schema); err != nil {
		db.Close()
		return nil, err
	}
	return &DB{db}, nil
}

func (d *DB) Ping(ctx context.Context) error {
	return d.DB.PingContext(ctx)
}

// --- Helpers ---

func newID() string { return uuid.New().String() }

func nowStr() string { return time.Now().UTC().Format(time.RFC3339) }

func parseTime(s string) time.Time {
	t, _ := time.Parse(time.RFC3339, s)
	return t
}

func parseTimePtr(s *string) *time.Time {
	if s == nil || *s == "" {
		return nil
	}
	t := parseTime(*s)
	return &t
}

func parseUUID(s string) uuid.UUID {
	id, _ := uuid.Parse(s)
	return id
}

func parseUUIDPtr(s *string) *uuid.UUID {
	if s == nil || *s == "" {
		return nil
	}
	id := parseUUID(*s)
	return &id
}

func marshalJSON(v any) []byte {
	b, _ := json.Marshal(v)
	return b
}

func unmarshalStrings(b []byte) []string {
	if len(b) == 0 {
		return nil
	}
	var s []string
	_ = json.Unmarshal(b, &s)
	return s
}

// --- ProjectRepo ---

type ProjectRepo struct{ db *DB }

func NewProjectRepo(db *DB) *ProjectRepo { return &ProjectRepo{db: db} }

var _ port.ProjectRepo = (*ProjectRepo)(nil)

func (r *ProjectRepo) Create(ctx context.Context, name, description string) (*model.Project, error) {
	id := newID()
	now := nowStr()
	_, err := r.db.ExecContext(ctx, "INSERT INTO projects (id, name, description, created_at, updated_at) VALUES (?, ?, ?, ?, ?)", id, name, description, now, now)
	if err != nil {
		return nil, err
	}
	return &model.Project{ID: parseUUID(id), Name: name, Description: description, CreatedAt: parseTime(now), UpdatedAt: parseTime(now)}, nil
}

func (r *ProjectRepo) Get(ctx context.Context, id uuid.UUID) (*model.Project, error) {
	row := r.db.QueryRowContext(ctx, "SELECT id, name, description, created_at, updated_at FROM projects WHERE id = ?", id.String())
	var p model.Project
	var idStr, createdAt, updatedAt string
	if err := row.Scan(&idStr, &p.Name, &p.Description, &createdAt, &updatedAt); err != nil {
		return nil, err
	}
	p.ID = parseUUID(idStr)
	p.CreatedAt = parseTime(createdAt)
	p.UpdatedAt = parseTime(updatedAt)
	return &p, nil
}

func (r *ProjectRepo) List(ctx context.Context) ([]model.Project, error) {
	rows, err := r.db.QueryContext(ctx, "SELECT id, name, description, created_at, updated_at FROM projects ORDER BY created_at DESC")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []model.Project
	for rows.Next() {
		var p model.Project
		var idStr, createdAt, updatedAt string
		if err := rows.Scan(&idStr, &p.Name, &p.Description, &createdAt, &updatedAt); err != nil {
			return nil, err
		}
		p.ID = parseUUID(idStr)
		p.CreatedAt = parseTime(createdAt)
		p.UpdatedAt = parseTime(updatedAt)
		out = append(out, p)
	}
	return out, nil
}

func (r *ProjectRepo) Update(ctx context.Context, id uuid.UUID, name, description string) (*model.Project, error) {
	now := nowStr()
	_, err := r.db.ExecContext(ctx, "UPDATE projects SET name = ?, description = ?, updated_at = ? WHERE id = ?", name, description, now, id.String())
	if err != nil {
		return nil, err
	}
	return r.Get(ctx, id)
}

func (r *ProjectRepo) Delete(ctx context.Context, id uuid.UUID) error {
	_, err := r.db.ExecContext(ctx, "DELETE FROM projects WHERE id = ?", id.String())
	return err
}

// --- SpecRepo ---

type SpecRepo struct{ db *DB }

func NewSpecRepo(db *DB) *SpecRepo { return &SpecRepo{db: db} }

var _ port.SpecRepo = (*SpecRepo)(nil)

func (r *SpecRepo) Create(ctx context.Context, in model.Spec) (*model.Spec, error) {
	id := newID()
	now := nowStr()
	_, err := r.db.ExecContext(ctx, "INSERT INTO specs (id, project_id, version, title, raw, analysis, created_at) VALUES (?, ?, ?, ?, ?, ?, ?)",
		id, in.ProjectID.String(), in.Version, in.Title, in.Raw, in.Analysis, now)
	if err != nil {
		return nil, err
	}
	in.ID = parseUUID(id)
	in.CreatedAt = parseTime(now)
	return &in, nil
}

func (r *SpecRepo) Get(ctx context.Context, id uuid.UUID) (*model.Spec, error) {
	row := r.db.QueryRowContext(ctx, "SELECT id, project_id, version, title, raw, analysis, created_at FROM specs WHERE id = ?", id.String())
	var s model.Spec
	var idStr, projStr, createdAt string
	if err := row.Scan(&idStr, &projStr, &s.Version, &s.Title, &s.Raw, &s.Analysis, &createdAt); err != nil {
		return nil, err
	}
	s.ID = parseUUID(idStr)
	s.ProjectID = parseUUID(projStr)
	s.CreatedAt = parseTime(createdAt)
	return &s, nil
}

func (r *SpecRepo) ListByProject(ctx context.Context, projectID uuid.UUID) ([]model.Spec, error) {
	rows, err := r.db.QueryContext(ctx, "SELECT id, project_id, version, title, raw, analysis, created_at FROM specs WHERE project_id = ? ORDER BY created_at", projectID.String())
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []model.Spec
	for rows.Next() {
		var s model.Spec
		var idStr, projStr, createdAt string
		if err := rows.Scan(&idStr, &projStr, &s.Version, &s.Title, &s.Raw, &s.Analysis, &createdAt); err != nil {
			return nil, err
		}
		s.ID = parseUUID(idStr)
		s.ProjectID = parseUUID(projStr)
		s.CreatedAt = parseTime(createdAt)
		out = append(out, s)
	}
	return out, nil
}

// --- OperationRepo ---

type OperationRepo struct{ db *DB }

func NewOperationRepo(db *DB) *OperationRepo { return &OperationRepo{db: db} }

var _ port.OperationRepo = (*OperationRepo)(nil)

func (r *OperationRepo) Create(ctx context.Context, in model.Operation) (*model.Operation, error) {
	id := newID()
	now := nowStr()
	_, err := r.db.ExecContext(ctx, "INSERT INTO operations (id, spec_id, operation_id, method, path, summary, tags, classification, is_destructive, overrides, created_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)",
		id, in.SpecID.String(), in.OperationID, in.Method, in.Path, in.Summary, marshalJSON(in.Tags), in.Classification, boolToInt(in.IsDestructive), in.Overrides, now)
	if err != nil {
		return nil, err
	}
	in.ID = parseUUID(id)
	in.CreatedAt = parseTime(now)
	return &in, nil
}

func (r *OperationRepo) ListBySpec(ctx context.Context, specID uuid.UUID) ([]model.Operation, error) {
	rows, err := r.db.QueryContext(ctx, "SELECT id, spec_id, operation_id, method, path, summary, tags, classification, is_destructive, overrides, created_at FROM operations WHERE spec_id = ? ORDER BY created_at", specID.String())
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []model.Operation
	for rows.Next() {
		o, err := scanOperation(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, *o)
	}
	return out, nil
}

func (r *OperationRepo) UpdateClassification(ctx context.Context, id uuid.UUID, classification string, isDestructive bool) (*model.Operation, error) {
	_, err := r.db.ExecContext(ctx, "UPDATE operations SET classification = ?, is_destructive = ? WHERE id = ?", classification, boolToInt(isDestructive), id.String())
	if err != nil {
		return nil, err
	}
	return r.get(ctx, id)
}

func (r *OperationRepo) UpdateOverrides(ctx context.Context, id uuid.UUID, overrides []byte) (*model.Operation, error) {
	_, err := r.db.ExecContext(ctx, "UPDATE operations SET overrides = ? WHERE id = ?", overrides, id.String())
	if err != nil {
		return nil, err
	}
	return r.get(ctx, id)
}

func (r *OperationRepo) get(ctx context.Context, id uuid.UUID) (*model.Operation, error) {
	row := r.db.QueryRowContext(ctx, "SELECT id, spec_id, operation_id, method, path, summary, tags, classification, is_destructive, overrides, created_at FROM operations WHERE id = ?", id.String())
	return scanOperationRow(row)
}

type scanner interface {
	Scan(dest ...any) error
}

func scanOperation(rows *sql.Rows) (*model.Operation, error) {
	var o model.Operation
	var idStr, specStr, createdAt string
	var tags, overrides []byte
	var destructive int
	if err := rows.Scan(&idStr, &specStr, &o.OperationID, &o.Method, &o.Path, &o.Summary, &tags, &o.Classification, &destructive, &overrides, &createdAt); err != nil {
		return nil, err
	}
	o.ID = parseUUID(idStr)
	o.SpecID = parseUUID(specStr)
	o.Tags = unmarshalStrings(tags)
	o.IsDestructive = destructive != 0
	o.Overrides = overrides
	o.CreatedAt = parseTime(createdAt)
	return &o, nil
}

func scanOperationRow(row *sql.Row) (*model.Operation, error) {
	var o model.Operation
	var idStr, specStr, createdAt string
	var tags, overrides []byte
	var destructive int
	if err := row.Scan(&idStr, &specStr, &o.OperationID, &o.Method, &o.Path, &o.Summary, &tags, &o.Classification, &destructive, &overrides, &createdAt); err != nil {
		return nil, err
	}
	o.ID = parseUUID(idStr)
	o.SpecID = parseUUID(specStr)
	o.Tags = unmarshalStrings(tags)
	o.IsDestructive = destructive != 0
	o.Overrides = overrides
	o.CreatedAt = parseTime(createdAt)
	return &o, nil
}

func boolToInt(b bool) int {
	if b {
		return 1
	}
	return 0
}

// --- RunRepo ---

type RunRepo struct{ db *DB }

func NewRunRepo(db *DB) *RunRepo { return &RunRepo{db: db} }

var _ port.RunRepo = (*RunRepo)(nil)

func (r *RunRepo) Create(ctx context.Context, projectID uuid.UUID, planID *uuid.UUID) (*model.Run, error) {
	id := newID()
	now := nowStr()
	var planStr *string
	if planID != nil {
		s := planID.String()
		planStr = &s
	}
	_, err := r.db.ExecContext(ctx, "INSERT INTO runs (id, project_id, plan_id, status, created_at) VALUES (?, ?, ?, 'pending', ?)", id, projectID.String(), planStr, now)
	if err != nil {
		return nil, err
	}
	return &model.Run{ID: parseUUID(id), ProjectID: projectID, PlanID: planID, Status: "pending", CreatedAt: parseTime(now)}, nil
}

func (r *RunRepo) Get(ctx context.Context, id uuid.UUID) (*model.Run, error) {
	row := r.db.QueryRowContext(ctx, "SELECT id, project_id, plan_id, status, started_at, finished_at, created_at FROM runs WHERE id = ?", id.String())
	return scanRunRow(row)
}

func (r *RunRepo) ListByProject(ctx context.Context, projectID uuid.UUID) ([]model.Run, error) {
	rows, err := r.db.QueryContext(ctx, "SELECT id, project_id, plan_id, status, started_at, finished_at, created_at FROM runs WHERE project_id = ? ORDER BY created_at", projectID.String())
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []model.Run
	for rows.Next() {
		var run model.Run
		var idStr, projStr, createdAt string
		var planStr, startedAt, finishedAt *string
		if err := rows.Scan(&idStr, &projStr, &planStr, &run.Status, &startedAt, &finishedAt, &createdAt); err != nil {
			return nil, err
		}
		run.ID = parseUUID(idStr)
		run.ProjectID = parseUUID(projStr)
		run.PlanID = parseUUIDPtr(planStr)
		run.StartedAt = parseTimePtr(startedAt)
		run.FinishedAt = parseTimePtr(finishedAt)
		run.CreatedAt = parseTime(createdAt)
		out = append(out, run)
	}
	return out, nil
}

func (r *RunRepo) UpdateStatus(ctx context.Context, id uuid.UUID, status string, startedAt, finishedAt *time.Time) (*model.Run, error) {
	var startStr, finishStr *string
	if startedAt != nil {
		s := startedAt.UTC().Format(time.RFC3339)
		startStr = &s
	}
	if finishedAt != nil {
		s := finishedAt.UTC().Format(time.RFC3339)
		finishStr = &s
	}
	_, err := r.db.ExecContext(ctx, "UPDATE runs SET status = ?, started_at = ?, finished_at = ? WHERE id = ?", status, startStr, finishStr, id.String())
	if err != nil {
		return nil, err
	}
	return r.Get(ctx, id)
}

func (r *RunRepo) SaveResult(ctx context.Context, runID uuid.UUID, result []byte) (*model.StoredRunResult, error) {
	id := newID()
	now := nowStr()
	_, err := r.db.ExecContext(ctx, "INSERT INTO run_results (id, run_id, result, created_at) VALUES (?, ?, ?, ?)", id, runID.String(), result, now)
	if err != nil {
		return nil, err
	}
	return &model.StoredRunResult{ID: parseUUID(id), RunID: runID, Result: result, CreatedAt: parseTime(now)}, nil
}

func (r *RunRepo) GetResult(ctx context.Context, runID uuid.UUID) (*model.StoredRunResult, error) {
	row := r.db.QueryRowContext(ctx, "SELECT id, run_id, result, created_at FROM run_results WHERE run_id = ?", runID.String())
	var rr model.StoredRunResult
	var idStr, runStr, createdAt string
	if err := row.Scan(&idStr, &runStr, &rr.Result, &createdAt); err != nil {
		return nil, err
	}
	rr.ID = parseUUID(idStr)
	rr.RunID = parseUUID(runStr)
	rr.CreatedAt = parseTime(createdAt)
	return &rr, nil
}

func (r *RunRepo) DeleteOlderThan(ctx context.Context, before time.Time) (int, error) {
	res, err := r.db.ExecContext(ctx, "DELETE FROM runs WHERE created_at < ?", before.UTC().Format(time.RFC3339))
	if err != nil {
		return 0, err
	}
	n, _ := res.RowsAffected()
	return int(n), nil
}

func scanRunRow(row *sql.Row) (*model.Run, error) {
	var run model.Run
	var idStr, projStr, createdAt string
	var planStr, startedAt, finishedAt *string
	if err := row.Scan(&idStr, &projStr, &planStr, &run.Status, &startedAt, &finishedAt, &createdAt); err != nil {
		return nil, err
	}
	run.ID = parseUUID(idStr)
	run.ProjectID = parseUUID(projStr)
	run.PlanID = parseUUIDPtr(planStr)
	run.StartedAt = parseTimePtr(startedAt)
	run.FinishedAt = parseTimePtr(finishedAt)
	run.CreatedAt = parseTime(createdAt)
	return &run, nil
}

// --- CommentRepo ---

type CommentRepo struct{ db *DB }

func NewCommentRepo(db *DB) *CommentRepo { return &CommentRepo{db: db} }

var _ port.CommentRepo = (*CommentRepo)(nil)

func (r *CommentRepo) Create(ctx context.Context, runID uuid.UUID, author, body string) (*model.Comment, error) {
	id := newID()
	now := nowStr()
	_, err := r.db.ExecContext(ctx, "INSERT INTO comments (id, run_id, author, body, created_at) VALUES (?, ?, ?, ?, ?)", id, runID.String(), author, body, now)
	if err != nil {
		return nil, err
	}
	return &model.Comment{ID: parseUUID(id), RunID: runID, Author: author, Body: body, CreatedAt: parseTime(now)}, nil
}

func (r *CommentRepo) ListByRun(ctx context.Context, runID uuid.UUID) ([]model.Comment, error) {
	rows, err := r.db.QueryContext(ctx, "SELECT id, run_id, author, body, created_at FROM comments WHERE run_id = ? ORDER BY created_at", runID.String())
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []model.Comment
	for rows.Next() {
		var c model.Comment
		var idStr, runStr, createdAt string
		if err := rows.Scan(&idStr, &runStr, &c.Author, &c.Body, &createdAt); err != nil {
			return nil, err
		}
		c.ID = parseUUID(idStr)
		c.RunID = parseUUID(runStr)
		c.CreatedAt = parseTime(createdAt)
		out = append(out, c)
	}
	return out, nil
}

// --- ArtifactRepo ---

type ArtifactRepo struct{ db *DB }

func NewArtifactRepo(db *DB) *ArtifactRepo { return &ArtifactRepo{db: db} }

var _ port.ArtifactRepo = (*ArtifactRepo)(nil)

func (r *ArtifactRepo) Create(ctx context.Context, runID uuid.UUID, typ, path string) (*model.Artifact, error) {
	id := newID()
	now := nowStr()
	_, err := r.db.ExecContext(ctx, "INSERT INTO artifacts (id, run_id, type, path, created_at) VALUES (?, ?, ?, ?, ?)", id, runID.String(), typ, path, now)
	if err != nil {
		return nil, err
	}
	return &model.Artifact{ID: parseUUID(id), RunID: runID, Type: typ, Path: path, CreatedAt: parseTime(now)}, nil
}

func (r *ArtifactRepo) ListByRun(ctx context.Context, runID uuid.UUID) ([]model.Artifact, error) {
	rows, err := r.db.QueryContext(ctx, "SELECT id, run_id, type, path, created_at FROM artifacts WHERE run_id = ? ORDER BY created_at", runID.String())
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []model.Artifact
	for rows.Next() {
		var a model.Artifact
		var idStr, runStr, createdAt string
		if err := rows.Scan(&idStr, &runStr, &a.Type, &a.Path, &createdAt); err != nil {
			return nil, err
		}
		a.ID = parseUUID(idStr)
		a.RunID = parseUUID(runStr)
		a.CreatedAt = parseTime(createdAt)
		out = append(out, a)
	}
	return out, nil
}

func (r *ArtifactRepo) DeleteByRun(ctx context.Context, runID uuid.UUID) ([]model.Artifact, error) {
	arts, err := r.ListByRun(ctx, runID)
	if err != nil {
		return nil, err
	}
	_, err = r.db.ExecContext(ctx, "DELETE FROM artifacts WHERE run_id = ?", runID.String())
	if err != nil {
		return nil, err
	}
	return arts, nil
}
