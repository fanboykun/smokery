// Phase 2 type definitions for mock API data

export interface ContractReport {
  id: string;
  run_id: string;
  total_violations: number;
  errors: ContractViolation[];
  warnings: ContractViolation[];
  compliance_score: number;
  passed_assertions: number;
  failed_assertions: number;
}

export interface ContractViolation {
  id: string;
  operation_id: string;
  violation_type: string;
  severity: string;
  message: string;
  expected_schema?: unknown;
  actual_value?: unknown;
  location: string;
}

export interface AnalystReport {
  id: string;
  run_id: string;
  summary: string;
  root_causes: { cause: string; impact: number; affected_operations: string[] }[];
  recommendations: { title: string; description: string; priority: string }[];
  timeline_insights: { timestamp: string; operation_id: string; status: string; duration_ms: number; error?: string }[];
}

export interface QAReport {
  id: string;
  run_id: string;
  status: string;
  total_tests: number;
  passed_tests: number;
  failed_tests: number;
  pass_rate: number;
  flaky_tests: string[];
  coverage_summary: { total_operations: number; tested_operations: number; coverage_percentage: number };
  blockers: { operation_id: string; issue: string; severity: string }[];
}

export interface CorrelationReport {
  id: string;
  run_id: string;
  root_trace_id: string;
  links: { type: string; name: string; url: string; trace_id?: string }[];
  metrics: { p50_latency_ms: number; p95_latency_ms: number; p99_latency_ms: number; error_rate: number };
  external_services: { name: string; status: string; latency_ms: number }[];
}

export interface TeamMember {
  id: string;
  name: string;
  email: string;
  role: string;
  avatar_url: string;
  joined_at: string;
}

export interface RunFailureClassification {
  id: string;
  run_id: string;
  classification: string;
  classified_by: string;
  classified_at: string;
}

export interface RunAssignment {
  id: string;
  run_id: string;
  assigned_to: string;
  assigned_by: string;
  assigned_at: string;
}

export interface FailureAction {
  id: string;
  run_id: string;
  action_type: string;
  actor_id: string;
  actor_name: string;
  created_at: string;
  details?: Record<string, unknown>;
}

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
  changes: { type: string; operation_id: string; path: string; method: string; breaking: boolean; details: string }[];
}

export interface ImpactAnalysis {
  spec_version_id: string;
  affected_flows: { flow_id: string; flow_name: string; affected_steps: number; impact: string }[];
  affected_suites: { suite_id: string; suite_name: string; affected_operations: number; impact: string }[];
  affected_runs: number;
  risk_assessment: string;
}

export interface LatencyAnalytics {
  range: string;
  data: { timestamp: string; p50: number; p95: number; p99: number }[];
  slowest_operations: { operation_id: string; avg_latency: number; p99_latency: number }[];
  fastest_operations: { operation_id: string; avg_latency: number }[];
}

export interface FlakyOperationsAnalytics {
  range: string;
  operations: FlakyOperation[];
  critical_flaky: FlakyOperation[];
}

export interface FlakyOperation {
  operation_id: string;
  path: string;
  method: string;
  success_rate: number;
  runs: number;
  failures: number;
  flakiness_score: number;
  trend: string;
}

export interface HealthTrendsAnalytics {
  range: string;
  data: { timestamp: string; date: string; total_runs: number; passed_runs: number; failed_runs: number; pass_rate: number }[];
  current_health: number;
  trend: string;
  weekly_average: number;
  monthly_average: number;
}

export interface ProjectMember {
  id: string;
  project_id: string;
  user_id: string;
  user_name: string;
  user_email: string;
  avatar_url?: string;
  role: string;
  added_at: string;
  added_by: string;
}

export interface AuditLogEntry {
  id: string;
  project_id: string;
  action: string;
  resource_type: string;
  resource_id: string;
  actor_id: string;
  actor_name: string;
  timestamp: string;
  changes?: { field: string; old_value: string; new_value: string }[];
}

export interface Webhook {
  id: string;
  project_id: string;
  name: string;
  url: string;
  events: string[];
  is_active: boolean;
  created_at: string;
  last_triggered_at?: string;
}

export interface WebhookDelivery {
  id: string;
  webhook_id: string;
  event: string;
  status: string;
  response_code?: number;
  delivered_at: string;
}

export interface NotificationRule {
  id: string;
  project_id: string;
  name: string;
  channel: string;
  channel_config: Record<string, unknown>;
  triggers: { event: string; conditions?: { field: string; operator: string; value: string }[] }[];
  is_active: boolean;
  created_at: string;
}
