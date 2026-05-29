package sqlite

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/google/uuid"

	mdl "github.com/fanboykun/smokery/apps/core/internal/model"
)

func testSpec(title string, projectID uuid.UUID) mdl.Spec {
	return mdl.Spec{ProjectID: projectID, Version: "1.0", Title: title, Raw: []byte("{}"), Analysis: []byte("[]")}
}

func testOp(specID uuid.UUID) mdl.Operation {
	return mdl.Operation{SpecID: specID, OperationID: "listUsers", Method: "GET", Path: "/users", Summary: "List users", Tags: []string{"users"}, Classification: "list"}
}

func TestSQLiteAdapter(t *testing.T) {
	dir := t.TempDir()
	dbPath := filepath.Join(dir, "test.db")

	db, err := Open(dbPath)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	ctx := context.Background()

	// Test ProjectRepo
	projectRepo := NewProjectRepo(db)
	p, err := projectRepo.Create(ctx, "Test Project", "A test")
	if err != nil {
		t.Fatal(err)
	}
	if p.Name != "Test Project" {
		t.Fatalf("expected 'Test Project', got %q", p.Name)
	}

	got, err := projectRepo.Get(ctx, p.ID)
	if err != nil {
		t.Fatal(err)
	}
	if got.ID != p.ID {
		t.Fatal("ID mismatch")
	}

	list, err := projectRepo.List(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if len(list) != 1 {
		t.Fatalf("expected 1 project, got %d", len(list))
	}

	// Test SpecRepo
	specRepo := NewSpecRepo(db)
	spec, err := specRepo.Create(ctx, testSpec("spec", p.ID))
	if err != nil {
		t.Fatal(err)
	}
	specs, err := specRepo.ListByProject(ctx, p.ID)
	if err != nil {
		t.Fatal(err)
	}
	if len(specs) != 1 || specs[0].ID != spec.ID {
		t.Fatal("spec list mismatch")
	}

	// Test OperationRepo
	opRepo := NewOperationRepo(db)
	op, err := opRepo.Create(ctx, testOp(spec.ID))
	if err != nil {
		t.Fatal(err)
	}
	ops, err := opRepo.ListBySpec(ctx, spec.ID)
	if err != nil {
		t.Fatal(err)
	}
	if len(ops) != 1 || ops[0].ID != op.ID {
		t.Fatal("operation list mismatch")
	}

	// Test RunRepo
	runRepo := NewRunRepo(db)
	run, err := runRepo.Create(ctx, p.ID, nil)
	if err != nil {
		t.Fatal(err)
	}
	if run.Status != "pending" {
		t.Fatalf("expected 'pending', got %q", run.Status)
	}

	// Test CommentRepo
	commentRepo := NewCommentRepo(db)
	c, err := commentRepo.Create(ctx, run.ID, "user", "hello")
	if err != nil {
		t.Fatal(err)
	}
	comments, err := commentRepo.ListByRun(ctx, run.ID)
	if err != nil {
		t.Fatal(err)
	}
	if len(comments) != 1 || comments[0].ID != c.ID {
		t.Fatal("comment list mismatch")
	}

	// Test ArtifactRepo
	artRepo := NewArtifactRepo(db)
	a, err := artRepo.Create(ctx, run.ID, "json", "/path/to/file.json")
	if err != nil {
		t.Fatal(err)
	}
	arts, err := artRepo.ListByRun(ctx, run.ID)
	if err != nil {
		t.Fatal(err)
	}
	if len(arts) != 1 || arts[0].ID != a.ID {
		t.Fatal("artifact list mismatch")
	}

	// Cleanup
	if err := os.Remove(dbPath); err != nil {
		t.Log("cleanup:", err)
	}
}
