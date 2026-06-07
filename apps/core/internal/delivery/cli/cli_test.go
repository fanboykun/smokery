package cli_test

import (
	"bytes"
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/google/uuid"

	"github.com/fanboykun/smokery/apps/core/internal/adapter/memory"
	"github.com/fanboykun/smokery/apps/core/internal/app"
	"github.com/fanboykun/smokery/apps/core/internal/delivery/cli"
	"github.com/fanboykun/smokery/apps/core/internal/model"
	"github.com/fanboykun/smokery/apps/core/internal/runner"
)

func newTestServices() *cli.Services {
	store := memory.NewStore()
	return &cli.Services{
		Project:   app.NewProjectService(memory.NewProjectRepo(store)),
		Spec:      app.NewSpecService(memory.NewSpecRepo(store), memory.NewOperationRepo(store)),
		Operation: app.NewOperationService(memory.NewSpecRepo(store), memory.NewOperationRepo(store)),
		Run:       app.NewRunService(memory.NewRunRepo(store), noopJobs{}),
		Report:    app.NewReportService(memory.NewRunRepo(store)),
		Runner:    runner.New(runner.DefaultOptions()),
	}
}

func TestImportSpecCommand(t *testing.T) {
	specFile := writeTemp(t, "spec.json", minimalSpec)
	svcs := newTestServices()
	root := cli.NewRootCmd(svcs)

	out := &bytes.Buffer{}
	root.SetOut(out)
	root.SetArgs([]string{"import-spec", specFile})

	if err := root.Execute(); err != nil {
		t.Fatalf("import-spec failed: %v", err)
	}
	if !strings.Contains(out.String(), "Test API") {
		t.Errorf("expected title in output, got: %s", out.String())
	}
	if !strings.Contains(out.String(), "listUsers") {
		t.Errorf("expected operation in output, got: %s", out.String())
	}
}

func TestImportSpecCommandJSON(t *testing.T) {
	specFile := writeTemp(t, "spec.json", minimalSpec)
	svcs := newTestServices()
	root := cli.NewRootCmd(svcs)

	out := &bytes.Buffer{}
	root.SetOut(out)
	root.SetArgs([]string{"import-spec", "--json", specFile})

	if err := root.Execute(); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(out.String(), `"operation_id"`) {
		t.Errorf("expected JSON output, got: %s", out.String())
	}
}

func TestCompileCommand(t *testing.T) {
	specFile := writeTemp(t, "spec.json", minimalSpec)
	cfgFile := writeTemp(t, "config.yaml", minimalConfig)
	svcs := newTestServices()
	root := cli.NewRootCmd(svcs)

	out := &bytes.Buffer{}
	root.SetOut(out)
	root.SetArgs([]string{"compile", "--config", cfgFile, "--spec", specFile, "--format", "json"})

	if err := root.Execute(); err != nil {
		t.Fatalf("compile failed: %v", err)
	}
	if !strings.Contains(out.String(), "flow_plans") {
		t.Errorf("expected plan in output, got: %s", out.String())
	}
}

func TestCompileCommandMissingFlags(t *testing.T) {
	svcs := newTestServices()
	root := cli.NewRootCmd(svcs)
	root.SetOut(&bytes.Buffer{})
	root.SetErr(&bytes.Buffer{})
	root.SetArgs([]string{"compile"})

	if err := root.Execute(); err == nil {
		t.Error("expected error for missing flags")
	}
}

func TestRunCommand(t *testing.T) {
	planFile := writeTemp(t, "plan.json", minimalPlan)
	svcs := newTestServices()
	root := cli.NewRootCmd(svcs)

	out := &bytes.Buffer{}
	root.SetOut(out)
	root.SetErr(&bytes.Buffer{})
	root.SetArgs([]string{"run", planFile})

	// This will fail because the target server doesn't exist, but it should
	// execute without panicking and show a summary
	_ = root.Execute()
	output := out.String()
	if !strings.Contains(output, "Status:") && !strings.Contains(output, "status") {
		t.Errorf("expected status in output, got: %s", output)
	}
}

func TestReportCommand(t *testing.T) {
	resultFile := writeTemp(t, "result.json", sampleResult)
	svcs := newTestServices()
	root := cli.NewRootCmd(svcs)

	out := &bytes.Buffer{}
	root.SetOut(out)
	root.SetArgs([]string{"report", "--view", "ci", resultFile})

	if err := root.Execute(); err != nil {
		t.Fatalf("report failed: %v", err)
	}
	if !strings.Contains(out.String(), "passed") {
		t.Errorf("expected 'passed' in output, got: %s", out.String())
	}
}

func TestReportCommandMermaid(t *testing.T) {
	resultFile := writeTemp(t, "result.json", sampleResult)
	svcs := newTestServices()
	root := cli.NewRootCmd(svcs)

	out := &bytes.Buffer{}
	root.SetOut(out)
	root.SetArgs([]string{"report", "--view", "mermaid", resultFile})

	if err := root.Execute(); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(out.String(), "sequenceDiagram") {
		t.Errorf("expected mermaid diagram, got: %s", out.String())
	}
}

// --- helpers ---

func writeTemp(t *testing.T, name, content string) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), name)
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}
	return path
}

type noopJobs struct{}

func (noopJobs) EnqueueRun(context.Context, uuid.UUID, *model.SmokePlan) error { return nil }
func (noopJobs) CancelRun(context.Context, uuid.UUID) error                    { return nil }

const minimalSpec = `{
  "openapi": "3.1.0",
  "info": {"title": "Test API", "version": "1.0"},
  "paths": {
    "/users": {
      "get": {
        "operationId": "listUsers",
        "summary": "List users",
        "responses": {"200": {"description": "ok"}}
      }
    }
  }
}`

const minimalConfig = `
environments:
  - id: dev
    name: dev
    base_url: http://localhost:9999
flows:
  - id: f1
    name: test-flow
    environment: dev
    steps:
      - name: list
        operation_id: listUsers
`

const minimalPlan = `{
  "id": "test-plan",
  "environment": {"id": "dev", "name": "dev", "base_url": "http://127.0.0.1:1"},
  "flow_plans": [{
    "flow_id": "f1",
    "name": "test",
    "steps": [{"name": "s1", "method": "GET", "path": "/users", "assertions": [{"type": "status", "expected": 200}]}]
  }]
}`

const sampleResult = `{
  "run_id": "r1",
  "status": "passed",
  "duration_ms": 50,
  "started_at": "2024-01-01T00:00:00Z",
  "finished_at": "2024-01-01T00:00:00Z",
  "flows": [{
    "flow_id": "f1",
    "name": "test-flow",
    "status": "passed",
    "steps": [{
      "name": "step1",
      "status": "passed",
      "request": {"method": "GET", "url": "http://localhost/users"},
      "response": {"status": 200},
      "assertions": [{"type": "status", "passed": true}],
      "duration_ms": 10
    }]
  }]
}`
