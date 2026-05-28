//go:build integration

package postgres_test

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	tcpostgres "github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"

	"github.com/fanboykun/smokery/apps/core/internal/adapter/postgres"
)

func setupTestDB(t *testing.T) *pgxpool.Pool {
	t.Helper()
	ctx := context.Background()

	container, err := tcpostgres.Run(ctx, "postgres:16",
		tcpostgres.WithDatabase("smokery_test"),
		tcpostgres.WithUsername("test"),
		tcpostgres.WithPassword("test"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).WithStartupTimeout(30*time.Second),
		),
	)
	require.NoError(t, err)
	t.Cleanup(func() { _ = container.Terminate(ctx) })

	connStr, err := container.ConnectionString(ctx, "sslmode=disable")
	require.NoError(t, err)

	pool, err := pgxpool.New(ctx, connStr)
	require.NoError(t, err)
	t.Cleanup(func() { pool.Close() })

	// Run migration
	migrationPath := filepath.Join("..", "..", "..", "db", "migrations", "000001_init.up.sql")
	migration, err := os.ReadFile(migrationPath)
	require.NoError(t, err)
	_, err = pool.Exec(ctx, string(migration))
	require.NoError(t, err)

	return pool
}

func TestProjectRepo_Integration(t *testing.T) {
	pool := setupTestDB(t)
	repo := postgres.NewProjectRepo(pool)
	ctx := context.Background()

	// Create
	p, err := repo.Create(ctx, "integration-test", "desc")
	require.NoError(t, err)
	assert.Equal(t, "integration-test", p.Name)
	assert.NotEmpty(t, p.ID)

	// Get
	got, err := repo.Get(ctx, p.ID)
	require.NoError(t, err)
	assert.Equal(t, p.ID, got.ID)

	// List
	list, err := repo.List(ctx)
	require.NoError(t, err)
	assert.Len(t, list, 1)

	// Update
	updated, err := repo.Update(ctx, p.ID, "new-name", "new-desc")
	require.NoError(t, err)
	assert.Equal(t, "new-name", updated.Name)

	// Delete
	err = repo.Delete(ctx, p.ID)
	require.NoError(t, err)
	_, err = repo.Get(ctx, p.ID)
	assert.Error(t, err)
}

func TestRunRepo_Integration(t *testing.T) {
	pool := setupTestDB(t)
	projectRepo := postgres.NewProjectRepo(pool)
	runRepo := postgres.NewRunRepo(pool)
	ctx := context.Background()

	// Setup project
	p, err := projectRepo.Create(ctx, "test-project", "")
	require.NoError(t, err)

	// Create run
	run, err := runRepo.Create(ctx, p.ID, nil)
	require.NoError(t, err)
	assert.Equal(t, "pending", run.Status)

	// Update status
	now := time.Now()
	updated, err := runRepo.UpdateStatus(ctx, run.ID, "running", &now, nil)
	require.NoError(t, err)
	assert.Equal(t, "running", updated.Status)

	// Save result
	result := []byte(`{"status":"passed","duration_ms":100}`)
	sr, err := runRepo.SaveResult(ctx, run.ID, result)
	require.NoError(t, err)
	assert.Equal(t, run.ID, sr.RunID)

	// Get result
	got, err := runRepo.GetResult(ctx, run.ID)
	require.NoError(t, err)
	assert.JSONEq(t, string(result), string(got.Result))

	// List by project
	runs, err := runRepo.ListByProject(ctx, p.ID)
	require.NoError(t, err)
	assert.Len(t, runs, 1)

	// DeleteOlderThan (future date should delete)
	count, err := runRepo.DeleteOlderThan(ctx, time.Now().Add(1*time.Hour))
	require.NoError(t, err)
	assert.Equal(t, 1, count)
}

func TestArtifactRepo_Integration(t *testing.T) {
	pool := setupTestDB(t)
	projectRepo := postgres.NewProjectRepo(pool)
	runRepo := postgres.NewRunRepo(pool)
	artifactRepo := postgres.NewArtifactRepo(pool)
	ctx := context.Background()

	p, _ := projectRepo.Create(ctx, "test", "")
	run, _ := runRepo.Create(ctx, p.ID, nil)

	// Create artifact
	a, err := artifactRepo.Create(ctx, run.ID, "json_report", "runs/abc/result.json")
	require.NoError(t, err)
	assert.Equal(t, "json_report", a.Type)

	// List
	arts, err := artifactRepo.ListByRun(ctx, run.ID)
	require.NoError(t, err)
	assert.Len(t, arts, 1)

	// DeleteByRun
	deleted, err := artifactRepo.DeleteByRun(ctx, run.ID)
	require.NoError(t, err)
	assert.Len(t, deleted, 1)

	arts, _ = artifactRepo.ListByRun(ctx, run.ID)
	assert.Len(t, arts, 0)
}
