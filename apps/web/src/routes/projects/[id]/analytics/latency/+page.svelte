<script lang="ts">
  import { page } from '$app/stores';
  import { createQuery } from '@tanstack/svelte-query';
  import { mockGetLatencyAnalytics } from '$lib/api/mock-phase2';
  import * as Card from '$lib/components/ui/card';
  import { Badge } from '$lib/components/ui/badge';
  import { Button } from '$lib/components/ui/button';
  import { Clock, TrendingUp, TrendingDown } from '@lucide/svelte';

  const projectId = $page.params.id!;
  let range = $state('7d');

  const analytics = createQuery(() => ({
    queryKey: ['analytics-latency', projectId, range],
    queryFn: () => mockGetLatencyAnalytics(projectId, range),
  }));
</script>

<main class="mx-auto max-w-5xl px-6 py-8 space-y-6">
  <div class="flex items-center justify-between">
    <div>
      <p class="text-xs font-bold uppercase tracking-widest text-primary">Analytics</p>
      <h1 class="text-2xl font-bold">Latency</h1>
    </div>
    <div class="flex gap-1">
      {#each ['7d', '30d', '90d'] as r}
        <Button variant={range === r ? 'default' : 'outline'} size="sm" onclick={() => (range = r)}>{r}</Button>
      {/each}
    </div>
  </div>

  {#if analytics.isPending}
    <p class="text-muted-foreground">Loading…</p>
  {:else if analytics.data}
    <!-- Latency Chart (table representation) -->
    <Card.Root>
      <Card.Header><Card.Title class="text-base">Response Time Trend</Card.Title></Card.Header>
      <Card.Content>
        <div class="grid grid-cols-3 gap-4 mb-4">
          <div class="text-center">
            <p class="text-xs text-muted-foreground">p50</p>
            <p class="text-2xl font-bold text-emerald-400">{Math.round(analytics.data.data.at(-1)?.p50 ?? 0)}ms</p>
          </div>
          <div class="text-center">
            <p class="text-xs text-muted-foreground">p95</p>
            <p class="text-2xl font-bold text-yellow-400">{Math.round(analytics.data.data.at(-1)?.p95 ?? 0)}ms</p>
          </div>
          <div class="text-center">
            <p class="text-xs text-muted-foreground">p99</p>
            <p class="text-2xl font-bold text-red-400">{Math.round(analytics.data.data.at(-1)?.p99 ?? 0)}ms</p>
          </div>
        </div>
        <!-- Sparkline bars -->
        <div class="flex items-end gap-1 h-24">
          {#each analytics.data.data as point}
            {@const maxP99 = Math.max(...analytics.data.data.map(d => d.p99))}
            <div class="flex-1 flex flex-col justify-end gap-px">
              <div class="bg-red-500/30 rounded-t" style="height: {(point.p99 / maxP99) * 100}%"></div>
            </div>
          {/each}
        </div>
        <div class="flex justify-between text-xs text-muted-foreground mt-1">
          <span>{new Date(analytics.data.data[0]?.timestamp ?? '').toLocaleDateString()}</span>
          <span>{new Date(analytics.data.data.at(-1)?.timestamp ?? '').toLocaleDateString()}</span>
        </div>
      </Card.Content>
    </Card.Root>

    <!-- Slowest Operations -->
    <Card.Root>
      <Card.Header><Card.Title class="flex items-center gap-2 text-base"><TrendingDown class="size-4 text-red-400" /> Slowest Operations</Card.Title></Card.Header>
      <Card.Content class="space-y-2">
        {#each analytics.data.slowest_operations as op (op.operation_id)}
          <div class="flex items-center justify-between rounded-lg border border-border p-3">
            <span class="font-mono text-sm">{op.operation_id}</span>
            <div class="flex gap-3 text-sm">
              <span class="text-muted-foreground">avg <span class="font-bold text-yellow-400">{op.avg_latency}ms</span></span>
              <span class="text-muted-foreground">p99 <span class="font-bold text-red-400">{op.p99_latency}ms</span></span>
            </div>
          </div>
        {/each}
      </Card.Content>
    </Card.Root>

    <!-- Fastest Operations -->
    <Card.Root>
      <Card.Header><Card.Title class="flex items-center gap-2 text-base"><TrendingUp class="size-4 text-emerald-400" /> Fastest Operations</Card.Title></Card.Header>
      <Card.Content class="space-y-2">
        {#each analytics.data.fastest_operations as op (op.operation_id)}
          <div class="flex items-center justify-between rounded-lg border border-border p-3">
            <span class="font-mono text-sm">{op.operation_id}</span>
            <span class="font-bold text-emerald-400">{op.avg_latency}ms</span>
          </div>
        {/each}
      </Card.Content>
    </Card.Root>
  {/if}
</main>
