// Phase 2 Feature Types - Backend Contract
// These types define what the frontend expects from the backend API
// Backend team should implement endpoints that return these exact shapes

// ============================================================================
// TIER 1: REPORTING
// ============================================================================

export interface ContractComplianceViolation {
  id: string;
  operation_id: string;
  violation_type: 'response_schema_mismatch' | 'missing_required_field' | 'unexpected_field' | 'type_mismatch';
  severity: 'error' | 'warning';
  message: string;
  expected_schema: Record<string, any>;
  actual_value: Record<string, any>;
  location: string; // e.g., "response.body.user.id"
}

export interface ContractReport {
  id: string;
  run_id: string;
  total_violations: number;
  errors: ContractComplianceViolation[];
  warnings: ContractComplianceViolation[];
  compliance_score: number; // 0-100
  passed_assertions: number;
  failed_assertions: number;
}

export interface AnalystReport {
  id: string;
  run_id: string;
  summary: string;
  root_causes: {
    cause: string;
    impact: number; // percentage of failures caused by this
    affected_operations: string[];
  }[];
  recommendations: {
    title: string;
    description: string;
    priority: 'high' | 'medium' | 'low';
  }[];
  timeline_insights: {
    timestamp: string;
    operation_id: string;
    status: 'passed' | 'failed';
    duration_ms: number;
    error?: string;
  }[];
}

export interface QAReport {
  id: string;
  run_id: string;
  status: 'passed' | 'failed';
  total_tests: number;
  passed_tests: number;
  failed_tests: number;
  pass_rate: number; // percentage
  flaky_tests: string[];
  coverage_summary: {
    total_operations: number;
    tested_operations: number;
    coverage_percentage: number;
  };
  blockers: {
    operation_id: string;
    issue: string;
    severity: 'critical' | 'high' | 'medium';
  }[];
}

export interface CorrelationLink {
  type: 'trace' | 'log' | 'metric';
  name: string;
  url: string;
  trace_id?: string;
  span_id?: string;
  timestamp?: string;
}

export interface CorrelationReport {
  id: string;
  run_id: string;
  root_trace_id: string;
  links: CorrelationLink[];
  metrics: {
    p50_latency_ms: number;
    p95_latency_ms: number;
    p99_latency_ms: number;
    error_rate: number;
  };
  external_services: {
    name: string;
    status: 'healthy' | 'degraded' | 'unhealthy';
    latency_ms: number;
  }[];
}

// ============================================================================
// TIER 2: FAILURE CLASSIFICATION
// ============================================================================

export type FailureClassification = 
  | 'network_timeout'
  | 'network_error'
  | 'auth_failure'
  | 'rate_limit'
  | 'server_error'
  | 'contract_violation'
  | 'test_flaky'
  | 'test_broken'
  | 'infrastructure'
  | 'unknown';

export interface RunFailureClassification {
  run_id: string;
  classification: FailureClassification;
  confidence: number; // 0-100
  classified_at: string;
  classified_by?: string;
  notes?: string;
}

export interface TeamMember {
  id: string;
  name: string;
  email: string;
  role: 'viewer' | 'editor' | 'admin';
  avatar_url?: string;
  joined_at: string;
}

export interface RunAssignment {
  run_id: string;
  assigned_to: string; // TeamMember.id
  assigned_at: string;
  assigned_by: string;
  status: 'open' | 'in_progress' | 'resolved' | 'won_t_fix';
}

export interface FailureAction {
  id: string;
  run_id: string;
  action_type: 'classified' | 'assigned' | 'commented' | 'status_changed';
  actor_id: string;
  actor_name: string;
  created_at: string;
  details: Record<string, any>;
}

// ============================================================================
// TIER 3: SPEC EVOLUTION
// ============================================================================

export interface SpecVersion {
  id: string;
  project_id: string;
  version: string;
  uploaded_at: string;
  uploaded_by: string;
  summary: string;
  operation_count: number;
  schema_changes: number;
  breaking_changes: number;
}

export interface SpecDiff {
  from_spec_id: string;
  to_spec_id: string;
  changes: {
    type: 'added' | 'removed' | 'modified';
    operation_id: string;
    path: string;
    method: string;
    breaking: boolean;
    details: string;
  }[];
}

export interface ImpactAnalysis {
  spec_version_id: string;
  affected_flows: {
    flow_id: string;
    flow_name: string;
    affected_steps: number;
    impact: 'breaking' | 'compatible' | 'unknown';
  }[];
  affected_suites: {
    suite_id: string;
    suite_name: string;
    affected_operations: number;
    impact: 'breaking' | 'compatible' | 'unknown';
  }[];
  affected_runs: number;
  risk_assessment: 'low' | 'medium' | 'high' | 'critical';
}

// ============================================================================
// TIER 4: ANALYTICS
// ============================================================================

export interface LatencyDataPoint {
  timestamp: string;
  p50: number;
  p95: number;
  p99: number;
  operation_id?: string;
}

export interface LatencyAnalytics {
  range: string; // "7d", "30d", "90d"
  data: LatencyDataPoint[];
  slowest_operations: {
    operation_id: string;
    avg_latency: number;
    p99_latency: number;
  }[];
  fastest_operations: {
    operation_id: string;
    avg_latency: number;
  }[];
}

export interface FlakyOperation {
  operation_id: string;
  path: string;
  method: string;
  success_rate: number; // percentage
  runs: number;
  failures: number;
  flakiness_score: number; // 0-100
  trend: 'improving' | 'degrading' | 'stable';
}

export interface FlakyOperationsAnalytics {
  range: string;
  operations: FlakyOperation[];
  critical_flaky: FlakyOperation[]; // < 80% success rate
}

export interface HealthTrendDataPoint {
  timestamp: string;
  date: string;
  total_runs: number;
  passed_runs: number;
  failed_runs: number;
  pass_rate: number; // percentage
}

export interface HealthTrendsAnalytics {
  range: string;
  data: HealthTrendDataPoint[];
  current_health: number; // percentage
  trend: 'improving' | 'degrading' | 'stable';
  weekly_average: number;
  monthly_average: number;
}

// ============================================================================
// TIER 5: TEAM GOVERNANCE
// ============================================================================

export interface ProjectMember {
  id: string;
  project_id: string;
  user_id: string;
  user_name: string;
  user_email: string;
  role: 'viewer' | 'editor' | 'admin';
  added_at: string;
  added_by: string;
}

export interface AuditLogEntry {
  id: string;
  project_id: string;
  action: 'create' | 'update' | 'delete' | 'run' | 'classify' | 'assign' | 'comment';
  resource_type: 'flow' | 'suite' | 'environment' | 'run' | 'spec';
  resource_id: string;
  actor_id: string;
  actor_name: string;
  timestamp: string;
  changes?: {
    field: string;
    old_value?: any;
    new_value?: any;
  }[];
  metadata?: Record<string, any>;
}

export interface Webhook {
  id: string;
  project_id: string;
  name: string;
  url: string;
  events: ('run.completed' | 'run.failed' | 'spec.updated' | 'member.added')[];
  is_active: boolean;
  created_at: string;
  last_triggered_at?: string;
  last_error?: string;
}

export interface WebhookDelivery {
  id: string;
  webhook_id: string;
  event: string;
  payload: Record<string, any>;
  status: 'success' | 'failed' | 'pending';
  response_status?: number;
  error?: string;
  delivered_at?: string;
}

export type NotificationChannel = 'email' | 'slack' | 'pagerduty';

export interface NotificationRule {
  id: string;
  project_id: string;
  name: string;
  channel: NotificationChannel;
  channel_config: {
    email?: string;
    slack_webhook_url?: string;
    pagerduty_integration_key?: string;
  };
  triggers: {
    event: 'run.failed' | 'run.passed' | 'flaky_detected' | 'spec_updated';
    conditions?: {
      field: string;
      operator: 'eq' | 'gt' | 'lt' | 'contains';
      value: any;
    }[];
  }[];
  is_active: boolean;
  created_at: string;
}

// ============================================================================
// API ENDPOINT CONTRACT DOCUMENTATION
// ============================================================================

/*
TIER 1: REPORTING
GET /api/runs/{id}/report/contract -> ContractReport
GET /api/runs/{id}/report/analyst -> AnalystReport
GET /api/runs/{id}/report/qa -> QAReport
GET /api/runs/{id}/correlations -> CorrelationReport

TIER 2: FAILURE CLASSIFICATION
GET /api/team/members -> TeamMember[]
PUT /api/runs/{id}/failure-classification -> RunFailureClassification
PUT /api/runs/{id}/assigned-to -> RunAssignment
GET /api/runs/{id}/actions -> FailureAction[]

TIER 3: SPEC EVOLUTION
GET /api/projects/{id}/specs -> SpecVersion[]
GET /api/specs/{from_id}/diff/{to_id} -> SpecDiff
GET /api/projects/{id}/impact/spec/{spec_id} -> ImpactAnalysis

TIER 4: ANALYTICS
GET /api/projects/{id}/analytics/latency?range=7d -> LatencyAnalytics
GET /api/projects/{id}/analytics/flaky-operations?range=30d -> FlakyOperationsAnalytics
GET /api/projects/{id}/analytics/health-trends?range=90d -> HealthTrendsAnalytics

TIER 5: TEAM GOVERNANCE
GET /api/projects/{id}/members -> ProjectMember[]
POST /api/projects/{id}/members -> ProjectMember
PUT /api/projects/{id}/members/{user_id} -> ProjectMember
DELETE /api/projects/{id}/members/{user_id}

GET /api/projects/{id}/audit-log?limit=100&offset=0 -> AuditLogEntry[]

GET /api/projects/{id}/webhooks -> Webhook[]
POST /api/projects/{id}/webhooks -> Webhook
PUT /api/projects/{id}/webhooks/{id} -> Webhook
DELETE /api/projects/{id}/webhooks/{id}
GET /api/webhooks/{id}/deliveries -> WebhookDelivery[]

GET /api/projects/{id}/notifications -> NotificationRule[]
POST /api/projects/{id}/notifications -> NotificationRule
PUT /api/projects/{id}/notifications/{id} -> NotificationRule
DELETE /api/projects/{id}/notifications/{id}
*/
