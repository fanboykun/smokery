# Backend Requirements for FE Excellence

To build a truly excellent frontend experience for power users, the backend needs these enhancements. They should be implemented in parallel or before FE work completes.

## Priority 1: Critical for MVP (Block if missing)

### 1. Enhanced `GET /api/projects` with Stats
Add per-project aggregation to reduce N+1 queries:
```json
{
  "id": "proj-123",
  "name": "API E-Commerce",
  "spec_count": 1,
  "last_run": {
    "id": "run-456",
    "status": "passed",
    "created_at": "2025-05-30T10:30:00Z",
    "duration_ms": 2450
  },
  "stats": {
    "total_runs": 47,
    "passed_runs": 45,
    "failed_runs": 2,
    "env_count": 3,
    "flow_count": 5,
    "suite_count": 8,
    "health_percentage": 95.7
  }
}
```

**Why**: Dashboard needs project health at a glance. Users shouldn't wait for N+1 queries.

### 2. Improved `POST /api/projects/{id}/plan/preview` Diagnostics
Return structured errors instead of generic "compilation failed":
```json
{
  "plan": {...},
  "diagnostics": {
    "errors": [
      {
        "code": "MISSING_AUTH_PROFILE",
        "severity": "error",
        "message": "Flow 'user-flow' references undefined auth profile 'oauth2'",
        "entity_id": "flow-123",
        "entity_type": "flow",
        "location": "flows[0].auth_profile_id"
      }
    ],
    "warnings": [
      {
        "code": "UNUSED_OPERATION",
        "severity": "warning",
        "message": "Operation 'GET /admin/stats' is not used by any flow",
        "entity_id": "op-789",
        "entity_type": "operation"
      }
    ],
    "summary": {
      "total_errors": 1,
      "total_warnings": 1,
      "is_compilable": false
    }
  }
}
```

**Codes to support**:
- `MISSING_AUTH_PROFILE` - Flow references undefined auth
- `MISSING_OPERATION` - Suite references undefined operation
- `INVALID_SELECTOR` - Suite selector matches 0 operations
- `UNUSED_OPERATION` - Operation not used anywhere
- `INVALID_ASSERTION` - Assertion syntax/logic error
- `SCHEMA_MISMATCH` - Type/value doesn't match schema

**Why**: FE needs to show users WHAT is wrong, not just fail. Power users need to debug in seconds.

### 3. WebSocket URL in Run Responses
Include connection details in `GET /api/runs/{id}`:
```json
{
  "id": "run-123",
  "status": "running",
  "created_at": "2025-05-30T10:30:00Z",
  "websocket_url": "wss://smokery.example.com/api/runs/run-123/events",
  "fallback_poll_interval_ms": 1000,
  "expires_at": "2025-05-30T11:30:00Z"
}
```

**Why**: Hardcoded `localhost:8080` breaks in production. FE needs dynamic URL from BE.

### 4. Detailed `GET /api/runs/{id}/result`
Include step-by-step execution details:
```json
{
  "id": "run-123",
  "status": "failed",
  "flows": [
    {
      "id": "flow-auth",
      "name": "User Authentication",
      "status": "passed",
      "duration_ms": 1200,
      "steps": [
        {
          "id": "step-1",
          "operation_id": "op-login",
          "name": "POST /auth/login",
          "status": "passed",
          "duration_ms": 450,
          "request": {
            "method": "POST",
            "url": "https://api.example.com/auth/login",
            "headers": {...},
            "body": {...}
          },
          "response": {
            "status": 200,
            "headers": {...},
            "body": {...}
          },
          "assertions": [
            {
              "id": "assert-1",
              "expression": "response.status == 200",
              "result": true,
              "actual_value": 200
            }
          ],
          "captures": [
            {
              "name": "access_token",
              "expression": "response.body.token",
              "value": "eyJhbGc..."
            }
          ]
        }
      ]
    }
  ]
}
```

**Why**: Live run viewer needs step-by-step breakdown. Users need to see exactly where/why tests fail.

---

## Priority 2: Power-User Features (Nice to have)

### 5. `GET /api/projects/{id}/overview`
New endpoint combining stats + recent activity:
```json
{
  "project_id": "proj-123",
  "name": "API E-Commerce",
  "recent_runs": [
    { "id": "run-456", "status": "passed", "created_at": "...", "duration_ms": 2450 }
  ],
  "recent_failures": [...],
  "created_at": "2025-05-01T00:00:00Z",
  "updated_at": "2025-05-30T10:30:00Z",
  "stats": {...}
}
```

### 6. Bulk Classification Endpoint
`POST /api/specs/{spec-id}/operations/bulk-classify` to classify N operations at once:
```json
{
  "classifications": [
    { "operation_id": "op-1", "classification": "smoke" },
    { "operation_id": "op-2", "classification": "critical" }
  ]
}
```

### 7. Classification Hints
`GET /api/specs/{spec-id}/operations/classification-hints` to suggest classifications:
```json
{
  "suggestions": [
    { "pattern": "GET /health", "suggested_classification": "smoke" },
    { "pattern": "POST /auth/login", "suggested_classification": "critical" }
  ]
}
```

### 8. `GET /api/auth-profile-types`
Schema for building auth forms dynamically:
```json
{
  "types": [
    {
      "name": "bearer_token",
      "fields": [
        { "name": "token", "type": "string", "required": true }
      ]
    },
    {
      "name": "basic_auth",
      "fields": [
        { "name": "username", "type": "string", "required": true },
        { "name": "password", "type": "string", "required": true }
      ]
    }
  ]
}
```

### 9. Flow/Suite Validation Endpoints
Real-time feedback on config:
- `POST /api/projects/{id}/flows/{id}/validate` - Check flow config
- `POST /api/projects/{id}/suites/{id}/preview-selector` - How many ops match the selector?
- `POST /api/projects/{id}/suites/{id}/preview-strategy` - How many test cases will be generated?

**Why**: FE needs real-time signals like "This selector matches 0 operations" before compilation.

### 10. Enhanced `GET /api/projects/{id}/runs` Query Params
Support filtering/sorting:
```
GET /api/projects/{id}/runs?status=failed&from=2025-05-01&to=2025-05-30&sort=created_at:desc
```

**Why**: Power users need to find failures quickly without scrolling thousands of runs.

---

## Priority 3: Polish (Later)

### 11. Server-Sent Events Fallback
`GET /api/runs/{id}/events` (SSE) as fallback when WebSocket unavailable.

### 12. Comments with Context
Add `context` field to `GET /api/runs/{id}/comments` to show run state at comment time.

### 13. Performance Headers
ETag, Cache-Control for immutable endpoints like `auth-profile-types`.

---

## Implementation Checklist

- [ ] Priority 1.1: Enhanced project stats endpoint
- [ ] Priority 1.2: Structured plan preview diagnostics
- [ ] Priority 1.3: Dynamic WebSocket URLs
- [ ] Priority 1.4: Step-by-step run result details
- [ ] Priority 2.x: Power-user features (parallel with FE work)
- [ ] Priority 3.x: Polish features (after MVP)
