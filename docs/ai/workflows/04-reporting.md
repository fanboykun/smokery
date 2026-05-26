# Workflow 04 — Reporting and Collaboration

## Goal

Transform persisted RunResults into audience-specific views, artifacts, diagrams, and collaboration records.

## Inputs

- RunResult
- Project report templates
- Retention/redaction policy
- Observability settings
- Comments and failure labels

## Outputs

- Backend Debug view
- Frontend Contract view
- Analyst Flow view
- QA Evidence view
- CI Summary
- JSON artifact
- HTML artifact
- Mermaid diagrams
- Trace/log links

## Required report views

### Backend Debug

Show:

- Failed requests
- Request/response details
- Assertions
- Schema errors
- Latency
- Request ID
- Trace ID
- Log/trace links

### Frontend Contract

Show:

- Flow sequence
- Endpoint dependencies
- Request/response shapes
- Captured values
- Example successful response shapes

### Analyst Flow

Show:

- Business-readable flow
- Step descriptions
- Expected outcomes
- Pass/fail state
- Minimal raw JSON

### QA Evidence

Show:

- Scenario result
- Assertions
- Evidence trail
- Repeatability metadata

### CI Summary

Show:

- Machine-readable status
- Counts
- Failed cases
- Artifact links

## Diagrams

Generate Mermaid for:

- Flow diagram
- Sequence diagram
- Capture dependency graph

## Collaboration

Support comments and failure labels on:

- Run
- Flow result
- Step result
- Suite result
- Case result

Failure states:

```text
new
acknowledged
investigating
expected
fixed
ignored
needs-config-update
```

Failure categories:

```text
backend bug
contract drift
test config issue
environment issue
data issue
auth issue
network issue
timeout issue
schema mismatch
assertion mismatch
unknown
```

## Acceptance criteria

- A run can produce backend debug and CI summary reports.
- Request ID and trace ID are visible when available.
- HTML and JSON artifacts can be stored in MinIO/S3.
- Comments can be attached to run/result entities.
- Mermaid sequence diagram can be generated from a flow run.
