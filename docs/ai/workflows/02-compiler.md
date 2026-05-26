# Workflow 02 — Compiler Pipeline

## Goal

Compile project configuration into a validated, **self-contained** executable SmokePlan.

The compiler is the bridge between OpenAPI knowledge and runtime execution. It is the **only** layer allowed to read the OpenAPI document. Everything the runner needs at runtime must be baked into the SmokePlan during compilation.

## Inputs

- Project
- Spec analysis (parsed by `internal/spec`)
- Operation registry
- Operation overrides
- Environment
- Auth profiles
- Flows
- Suites
- Fixtures
- Defaults
- Retention/redaction policy

## Output

- `*model.SmokePlan` — the AST: ordered, validated, self-contained
- Compile warnings
- Compile errors
- Plan preview summary

## Pipeline

```text
1. Load OpenAPI/spec analysis
2. Build normalized operation registry
3. Apply operation overrides
4. Resolve environment
5. Resolve auth profiles
6. Resolve fixtures
7. Compile flows
8. Compile suites
9. Generate suite cases
10. Validate required inputs
11. Validate destructive permissions
12. Validate response shape config
13. Bake runtime metadata (schemas, content types, param formats)
14. Validate retention/redaction policy
15. Produce SmokePlan
```

## Flow compilation

For each flow:

- Resolve operation references.
- Resolve step order.
- Validate path/query/header/body inputs.
- Validate capture references.
- Validate dependency references.
- Validate cleanup behavior.
- Validate destructive permissions.
- Embed any response schemas the runner will validate against.

## Suite compilation

For each suite:

- Apply operation selector.
- Produce matched operation groups.
- Generate cases from strategies.
- Validate case limits.
- Validate response shape hints.
- Validate pagination/search/filter config.
- Set concurrency.

## Compile error shape

```json
{
  "stage": "flow_steps",
  "path": "flows.user_crud.steps.get_user.request.path.id",
  "severity": "error",
  "message": "path param id has no source",
  "entityType": "flow_step",
  "entityId": "get_user"
}
```

## Self-containment guarantee

After the compiler produces a SmokePlan, the runner must be able to execute it without reading the OpenAPI document. If the runner needs:

- A response JSON schema → it must be embedded in the plan
- A path parameter format/regex → embedded
- A list-shape selector → embedded
- A content-type → embedded

If the compiler has not embedded something, the runner cannot use it. There is no fallback to "look it up from the spec at runtime."

This is what allows the runner to be a dependency-free library shared by the API server and the CLI.

## Acceptance criteria

- Invalid configs fail before execution.
- Destructive operations are blocked unless allowed.
- Plan preview can show all generated cases.
- Runner can execute the SmokePlan **without** reading UI config or the OpenAPI spec.
- Same input config produces the same plan.
