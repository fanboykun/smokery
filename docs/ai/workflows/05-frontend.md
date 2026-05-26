# Workflow 05 — Frontend Implementation

## Goal

Build the SvelteKit UI for managing projects, specs, operations, configurations, plan previews, runs, and reports.

## Primary UX flow

```text
Project list
  → Project detail
  → Spec import / analysis
  → Operation review
  → Environment/auth setup
  → Suite builder
  → Flow builder
  → Plan preview
  → Run live view
  → Report views
```

## Required route areas

```text
/projects
/projects/[id]/spec
/projects/[id]/operations
/projects/[id]/environments
/projects/[id]/flows/[fid]
/projects/[id]/suites/[sid]
/projects/[id]/plan
/projects/[id]/runs
/runs/[runId]
/runs/[runId]/report/[view]
/projects/[id]/settings
```

## Key components

### Spec and operation UI

- `OperationTable`
- `ClassificationOverrideForm`
- `SpecDiffViewer`
- `AnalysisSummaryCard`

### Flow builder

- `StepList`
- `StepEditor`
- `OperationPicker`
- `CaptureEditor`
- `ExpectationEditor`
- `FlowDiagramPreview`

### Suite builder

- `OperationSelectorEditor`
- `ListStrategyEditor`
- `ResponseShapeEditor`
- `EmptyResultPolicyEditor`
- `ConcurrencyEditor`

### Plan preview

- `CompileErrorList`
- `PlanSummaryCard`
- `GeneratedCaseTree`
- `DestructiveStepWarning`

### Run viewer

- `LiveEventFeed`
- `FlowResultTree`
- `SuiteResultTree`
- `RequestResponsePanel`
- `TraceContextCard`
- `FailureClassifier`

### Reports

- `AudienceReportTabs`
- `SequenceDiagram`
- `CaptureDependencyGraph`
- `LatencyChart`
- `RunTrendChart`
- `CommentThread`

## Frontend rules

- Use SvelteKit + TypeScript strict mode.
- Use TanStack Query for server state.
- Use Superforms + Zod for mutation forms.
- Use Monaco for raw JSON/YAML config editing.
- Use Mermaid for diagrams.
- Keep compiler errors visible and actionable.
- Do not allow run start if plan has blocking errors.

## Acceptance criteria

- User can import spec and view operations.
- User can override operation classification.
- User can configure environment and auth.
- User can build a basic flow.
- User can configure list sanity suite.
- User can preview a compiled plan.
- User can start a run and watch live progress.
- User can view a persisted report.
