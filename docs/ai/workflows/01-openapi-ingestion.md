# Workflow 01 — OpenAPI Ingestion and Operation Registry

## Goal

Import an OpenAPI spec, analyze it, classify operations, and build an operation registry that the rest of the platform can use.

## Inputs

- Project ID
- OpenAPI file or URL
- Optional spec version label
- Existing operation overrides

## Outputs

- Stored spec version
- Raw spec artifact
- Spec analysis JSON
- Operation registry
- Classification warnings
- Suggested suites and draft flows

## Required implementation stages

### 1. Accept spec input

Support upload or pasted content initially.

Later support URL import and scheduled sync.

### 2. Parse spec

Use `pb33f/libopenapi`.

Support OpenAPI 3.0 and 3.1.

Return structured parse errors.

### 3. Normalize operations

For each operation, extract:

- method
- path
- operationId
- tags
- summary
- description
- parameters
- request body schema
- response schemas
- security requirements
- examples

### 4. Classify operations

Classify into:

```text
health
list
single
lookup
search
create
update
delete
action
auth
report
export
internal
unknown
```

Classification should include confidence.

### 5. Detect list hints

For list-like operations, detect:

- response list paths
- pagination params
- search params
- enum filter params
- possible identity fields

### 6. Detect destructive behavior

Default method assumptions:

```text
GET/HEAD/OPTIONS = safe
POST/PUT/PATCH/DELETE = unsafe by default
```

Override for known safe POST searches.

### 7. Apply user overrides

Overrides must win over inferred values.

### 8. Persist analysis

Store:

- raw spec
- analysis JSONB
- operation registry view
- classification warnings

## Acceptance criteria

- Invalid OpenAPI input produces useful errors.
- Valid OpenAPI input produces operation registry.
- List-like GET endpoints are detected.
- Parameterized GET endpoints are distinguished from list endpoints.
- Mutating endpoints are marked destructive by default.
- User overrides are preserved across re-import.
