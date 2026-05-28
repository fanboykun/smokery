# smokery — OpenAPI Smoke Testing CLI

A standalone CLI tool that compiles and runs smoke tests against APIs using OpenAPI specs. Same compiler and runner as the Smokery web platform — no database required.

## Installation

```bash
# From source
make install-cli

# Or build manually
cd apps/core && go build -o smokery ./cmd/smokery
```

## Quick Start

```bash
# 1. Inspect an OpenAPI spec
smokery import-spec api.yaml

# 2. Compile a config + spec into an executable plan
smokery compile --config smoke.yaml --spec api.yaml -o plan.yaml

# 3. Run the plan
smokery run plan.yaml

# 4. Generate a report from saved results
smokery run plan.yaml --json > result.json
smokery report result.json --view debug
```

## Commands

### `smokery import-spec <spec-file>`

Parse an OpenAPI spec and list classified operations.

```bash
# Human-readable output
smokery import-spec petstore.yaml

# JSON output (for piping or saving)
smokery import-spec petstore.yaml --json

# Save to file
smokery import-spec petstore.yaml -o analysis.json
```

### `smokery compile`

Compile a project config + OpenAPI spec into an executable SmokePlan.

```bash
smokery compile --config smoke.yaml --spec api.yaml

# Output as JSON
smokery compile --config smoke.yaml --spec api.yaml --format json

# Save to file
smokery compile --config smoke.yaml --spec api.yaml -o plan.yaml
```

**Flags:**
- `--config` — Path to project config (YAML/JSON). Required.
- `--spec` — Path to OpenAPI spec file. Required.
- `--format` — Output format: `yaml` (default) or `json`.
- `-o, --out` — Write plan to file instead of stdout.

### `smokery run <plan-file>`

Execute a compiled SmokePlan against the target API.

```bash
# Run with CI-friendly summary
smokery run plan.yaml

# Full JSON result
smokery run plan.yaml --json

# Save artifacts (HTML report + JSON result)
smokery run plan.yaml -o ./reports/
```

**Flags:**
- `--json` — Output the full RunResult as JSON.
- `-o, --output` — Directory to save report artifacts (result.json + report.html).

**Exit codes:**
- `0` — All tests passed.
- `1` — One or more tests failed.

### `smokery report <result-file>`

Render a report view from a saved RunResult JSON file.

```bash
# CI summary (default)
smokery report result.json

# Backend debug view (failures + traces)
smokery report result.json --view debug

# Mermaid sequence diagram
smokery report result.json --view mermaid
```

**Views:**
- `ci` — Pass/fail counts, duration, failure list.
- `debug` — Failures with request URLs, status codes, trace IDs.
- `mermaid` — Sequence diagram of flow steps.

## Config Format

The project config is a YAML file defining environments, auth, flows, and suites:

```yaml
environments:
  - id: staging
    name: Staging
    base_url: https://api.staging.example.com
    headers:
      X-Source: smokery

auth_profiles:
  - id: bearer
    name: Bearer Token
    type: bearer
    config:
      token: "${API_TOKEN}"

flows:
  - id: user-crud
    name: User CRUD
    environment: staging
    auth: bearer
    steps:
      - name: create-user
        operation_id: createUser
        body: { name: "Test User", email: "test@example.com" }
        assertions:
          - type: status
            expected: 201
        captures:
          - name: user_id
            source: body
            path: id

      - name: get-user
        operation_id: getUser
        params: { userId: "{{user_id}}" }
        assertions:
          - type: status
            expected: 200

    cleanup:
      - name: delete-user
        operation_id: deleteUser
        params: { userId: "{{user_id}}" }

suites:
  - id: list-smoke
    name: List Endpoints
    environment: staging
    selector:
      classifications: [list]
    strategy:
      default_list: true
      pagination: true
      empty_result_policy: warn
```

## Examples

See [`/examples/configs/`](../../../examples/configs/) for full config examples and [`/examples/plans/`](../../../examples/plans/) for pre-compiled plans.

## CI Integration

```yaml
# GitHub Actions example
- name: Run smoke tests
  run: |
    smokery compile --config smoke.yaml --spec api.yaml -o plan.yaml
    smokery run plan.yaml -o ./reports/
  
- name: Upload report
  uses: actions/upload-artifact@v4
  if: always()
  with:
    name: smoke-report
    path: ./reports/
```
