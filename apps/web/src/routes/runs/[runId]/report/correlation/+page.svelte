<script lang="ts">
  import { page } from '$app/stores';
  import { createQuery } from '@tanstack/svelte-query';
  import { mockGetCorrelationReport } from '$lib/api/mock-phase2';
  import * as Card from '$lib/components/ui/card';
  import { Badge } from '$lib/components/ui/badge';
  import { ExternalLink, Activity } from '@lucide/svelte';

  const runId = $page.params.runId!;

  const report = createQuery(() => ({
    queryKey: ['runs', runId, 'report', 'correlation'],
    queryFn: () => mockGetCorrelationReport(runId),
  }));
</script>

<main class="space-y-6 px-6 py-8">
  {#if report.isPending}
    <p class="text-muted-foreground">Loading correlation report…</p>
  {:else if report.isError}
    <div class="rounded-lg border border-destructive/50 bg-destructive/5 p-4">
      <p class="text-sm text-destructive">Failed to load report</p>
    </div>
  {:else if report.data}
    <!-- Trace Info -->
    <Card.Root class="border-blue-500/30 bg-blue-500/5">
      <Card.Header>
        <Card.Title class="flex items-center gap-2">
          <Activity class="size-5 text-blue-400" />
          Root Trace
        </Card.Title>
      </Card.Header>
      <Card.Content>
        <div class="flex items-center justify-between">
          <code class="font-mono text-sm text-blue-300">{report.data.root_trace_id}</code>
          <button class="text-xs text-blue-400 hover:text-blue-300 transition-colors">Copy</button>
        </div>
      </Card.Content>
    </Card.Root>

    <!-- Performance Metrics -->
    <div class="grid gap-4 md:grid-cols-4">
      <Card.Root>
        <Card.Content class="text-center py-6">
          <p class="text-3xl font-bold text-emerald-400">{Math.round(report.data.metrics.p50_latency_ms)}</p>
          <p class="text-xs text-muted-foreground mt-2">P50 Latency (ms)</p>
        </Card.Content>
      </Card.Root>
      <Card.Root>
        <Card.Content class="text-center py-6">
          <p class="text-3xl font-bold text-yellow-400">{Math.round(report.data.metrics.p95_latency_ms)}</p>
          <p class="text-xs text-muted-foreground mt-2">P95 Latency (ms)</p>
        </Card.Content>
      </Card.Root>
      <Card.Root>
        <Card.Content class="text-center py-6">
          <p class="text-3xl font-bold text-orange-400">{Math.round(report.data.metrics.p99_latency_ms)}</p>
          <p class="text-xs text-muted-foreground mt-2">P99 Latency (ms)</p>
        </Card.Content>
      </Card.Root>
      <Card.Root>
        <Card.Content class="text-center py-6">
          <p class="text-3xl font-bold text-red-400">{(report.data.metrics.error_rate * 100).toFixed(1)}%</p>
          <p class="text-xs text-muted-foreground mt-2">Error Rate</p>
        </Card.Content>
      </Card.Root>
    </div>

    <!-- Observability Links -->
    <Card.Root>
      <Card.Header><Card.Title>Observability Links</Card.Title></Card.Header>
      <Card.Content class="space-y-3">
        {#each report.data.links as link (link.url)}
          <a href={link.url} target="_blank" rel="noopener noreferrer" class="flex items-center justify-between rounded-lg border border-border p-4 hover:bg-muted transition-colors">
            <div>
              <p class="font-semibold">{link.name}</p>
              <p class="text-xs text-muted-foreground mt-1">
                {link.type.charAt(0).toUpperCase() + link.type.slice(1)}
                {#if link.trace_id}
                  • <code class="font-mono">{link.trace_id.slice(0, 8)}...</code>
                {/if}
              </p>
            </div>
            <ExternalLink class="size-4 text-muted-foreground flex-shrink-0" />
          </a>
        {/each}
      </Card.Content>
    </Card.Root>

    <!-- Service Health -->
    <Card.Root>
      <Card.Header><Card.Title>Service Health</Card.Title></Card.Header>
      <Card.Content class="space-y-3">
        {#each report.data.external_services as service (service.name)}
          <div class="flex items-center justify-between p-3 rounded-lg border border-border">
            <div class="flex items-center gap-3">
              <div
                class="size-3 rounded-full flex-shrink-0"
                class:bg-emerald-500={service.status === 'healthy'}
                class:bg-yellow-500={service.status === 'degraded'}
                class:bg-red-500={service.status === 'unhealthy'}
              ></div>
              <div>
                <p class="font-semibold">{service.name}</p>
                <p class="text-xs text-muted-foreground">{service.latency_ms}ms latency</p>
              </div>
            </div>
            <Badge
              class={service.status === 'healthy' ? 'bg-emerald-500/20 text-emerald-300' : service.status === 'degraded' ? 'bg-yellow-500/20 text-yellow-300' : 'bg-red-500/20 text-red-300'}
            >
              {service.status}
            </Badge>
          </div>
        {/each}
      </Card.Content>
    </Card.Root>
  {/if}
</main>
