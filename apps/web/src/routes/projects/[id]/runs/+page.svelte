<script lang="ts">
  import { page } from '$app/stores';
  import { createQuery, createMutation, useQueryClient } from '@tanstack/svelte-query';
  import { api } from '$lib/api/client';
  import * as Card from '$lib/components/ui/card';
  import { Badge } from '$lib/components/ui/badge';
  import { Button } from '$lib/components/ui/button';
  import { Chart, Area, Axis, Svg } from 'layerchart';
  import Breadcrumb from '$lib/components/Breadcrumb.svelte';
  import { scaleTime, scaleLinear } from 'd3-scale';

  const queryClient = useQueryClient();
  const projectId = $page.params.id!;

  const runs = createQuery(() => ({
    queryKey: ['runs', projectId],
    queryFn: async () => {
      const { data, error } = await api.GET('/api/projects/{project-id}/runs', {
        params: { path: { 'project-id': projectId } },
      });
      if (error) throw error;
      return data ?? [];
    },
  }));

  const startRun = createMutation(() => ({
    mutationFn: async () => {
      const { data, error } = await api.POST('/api/projects/{project-id}/runs', {
        params: { path: { 'project-id': projectId } },
        body: {
          plan_id: crypto.randomUUID(),
          plan: {
            id: crypto.randomUUID(),
            project_id: projectId,
            environment: { id: 'default', name: 'dev', base_url: 'http://localhost:8080' },
            flow_plans: [],
            suite_plans: [],
            compiled_at: new Date().toISOString(),
          },
        },
      });
      if (error) throw error;
      return data!;
    },
    onSuccess: () => queryClient.invalidateQueries({ queryKey: ['runs', projectId] }),
  }));

  // Compute trend data: cumulative pass rate over time
  const trendData = $derived(
    (runs.data ?? [])
      .slice()
      .sort((a, b) => new Date(a.created_at).getTime() - new Date(b.created_at).getTime())
      .reduce<{ date: Date; passRate: number }[]>((acc, run, i) => {
        const passed = acc.filter((d) => d.passRate > 50).length;
        const total = i + 1;
        const isPass = run.status === 'completed' ? 1 : 0;
        const rate = Math.round(((passed + isPass) / total) * 100);
        acc.push({ date: new Date(run.created_at), passRate: rate });
        return acc;
      }, []),
  );

  function statusVariant(status: string) {
    if (status === 'completed') return 'default' as const;
    if (status === 'failed') return 'destructive' as const;
    return 'secondary' as const;
  }
</script>

<main class="mx-auto max-w-5xl px-6 py-8">
  <Breadcrumb crumbs={[{ label: projectId.slice(0, 8), href: `/projects/${projectId}` }, { label: 'Runs' }]} />
  <div class="mb-6 flex flex-wrap items-end justify-between gap-4">
    <div>
      <p class="text-xs font-bold uppercase tracking-widest text-primary">Project {projectId.slice(0, 8)}</p>
      <h1 class="text-3xl font-bold">Runs</h1>
      <p class="text-sm text-muted-foreground">Execution history and pass-rate trends.</p>
    </div>
    <div class="flex gap-2">
      <Button variant="outline" href="/projects/{projectId}">← Project</Button>
      <Button onclick={() => startRun.mutate()} disabled={startRun.isPending}>
        {startRun.isPending ? 'Starting…' : '▶ Start Run'}
      </Button>
    </div>
  </div>

  <!-- Trend chart -->
  {#if trendData.length >= 2}
    <Card.Root class="mb-6">
      <Card.Header><Card.Title class="text-base">Pass Rate Trend</Card.Title></Card.Header>
      <Card.Content>
        <div class="h-32">
          <Chart data={trendData} x="date" y="passRate" xScale={scaleTime()} yScale={scaleLinear()} yDomain={[0, 100]} padding={{ top: 8, bottom: 24, left: 32, right: 8 }}>
            <Svg>
              <Axis placement="left" format={(v) => `${v}%`} ticks={3} />
              <Axis placement="bottom" ticks={4} />
              <Area class="fill-primary/20 stroke-primary" line={{ class: 'stroke-primary stroke-2' }} />
            </Svg>
          </Chart>
        </div>
      </Card.Content>
    </Card.Root>
  {/if}

  <!-- Runs list -->
  {#if runs.isLoading}
    <p class="text-muted-foreground">Loading…</p>
  {:else if (runs.data ?? []).length === 0}
    <Card.Root>
      <Card.Content class="py-12 text-center text-sm text-muted-foreground">
        No runs yet. Start one above.
      </Card.Content>
    </Card.Root>
  {:else}
    <div class="space-y-2">
      {#each (runs.data ?? []).slice().reverse() as r (r.id)}
        <a href="/runs/{r.id}" class="flex items-center gap-4 rounded-lg border border-border p-3 transition-colors hover:bg-secondary">
          <Badge variant={statusVariant(r.status)}>{r.status}</Badge>
          <span class="flex-1 font-mono text-sm">{r.id.slice(0, 8)}</span>
          <span class="text-xs text-muted-foreground">{new Date(r.created_at).toLocaleString()}</span>
        </a>
      {/each}
    </div>
  {/if}
</main>
