<script lang="ts">
  import { page } from '$app/stores';
  import { createQuery } from '@tanstack/svelte-query';
  import { api } from '$lib/api/client';
  import { onMount } from 'svelte';
  import MermaidDiagram from '$lib/components/MermaidDiagram.svelte';
  import * as Card from '$lib/components/ui/card';
  import { Badge } from '$lib/components/ui/badge';
  import { Button } from '$lib/components/ui/button';

  interface RunEvent {
    type: string;
    run_id: string;
    data?: unknown;
  }

  let events = $state<RunEvent[]>([]);
  const runId = $page.params.runId!;

  const run = createQuery(() => ({
    queryKey: ['run', runId],
    queryFn: async () => {
      const { data, error } = await api.GET('/api/runs/{id}', { params: { path: { id: runId } } });
      if (error) throw error;
      return data!;
    },
    refetchInterval: 3000,
  }));

  const debugReport = createQuery(() => ({
    queryKey: ['run-debug', runId],
    queryFn: async () => {
      const { data, error } = await api.GET('/api/runs/{id}/report/debug', { params: { path: { id: runId } } });
      if (error) throw error;
      return data!;
    },
    enabled: run.data?.status === 'completed' || run.data?.status === 'failed',
  }));

  const mermaid = createQuery(() => ({
    queryKey: ['run-mermaid', runId],
    queryFn: async () => {
      const { data, error } = await api.GET('/api/runs/{id}/report/mermaid', { params: { path: { id: runId } } });
      if (error) throw error;
      return data!;
    },
    enabled: run.data?.status === 'completed' || run.data?.status === 'failed',
  }));

  onMount(() => {
    const ws = new WebSocket(`ws://localhost:8080/ws/runs/${runId}`);
    ws.onmessage = (e) => {
      const event: RunEvent = JSON.parse(e.data);
      events = [...events, event];
    };
    return () => ws.close();
  });

  function statusVariant(status: string) {
    if (status === 'completed') return 'default' as const;
    if (status === 'failed') return 'destructive' as const;
    return 'secondary' as const;
  }
</script>

<main class="mx-auto max-w-5xl px-6 py-8">
  <div class="mb-6 flex flex-wrap items-end justify-between gap-4">
    <div>
      <p class="text-xs font-bold uppercase tracking-widest text-primary">Run Detail</p>
      <h1 class="text-2xl font-bold font-mono">{runId.slice(0, 8)}</h1>
    </div>
    <div class="flex gap-2">
      {#if run.data}
        <Badge variant={statusVariant(run.data.status)}>{run.data.status}</Badge>
      {/if}
      <Button variant="outline" href="/runs/{runId}/comments">Comments</Button>
    </div>
  </div>

  <!-- Debug Report -->
  {#if debugReport.data}
    <Card.Root class="mb-4">
      <Card.Header><Card.Title class="text-base">Debug Report</Card.Title></Card.Header>
      <Card.Content class="space-y-3">
        <div class="flex gap-4 text-sm">
          <span>Status: <Badge variant={statusVariant(debugReport.data.status)}>{debugReport.data.status}</Badge></span>
          <span class="text-muted-foreground">{debugReport.data.duration_ms}ms</span>
        </div>
        {#if debugReport.data.failures && debugReport.data.failures.length > 0}
          <div class="space-y-1">
            <p class="text-xs font-bold uppercase text-destructive">Failures</p>
            {#each debugReport.data.failures as f}
              <div class="rounded-md bg-destructive/10 p-2 text-xs">
                <p class="font-medium">{f.step}: {f.assertion}</p>
                <p class="text-muted-foreground">{f.message} — {f.request_url} ({f.status})</p>
              </div>
            {/each}
          </div>
        {/if}
        {#if debugReport.data.traces && debugReport.data.traces.length > 0}
          <div class="space-y-1">
            <p class="text-xs font-bold uppercase text-muted-foreground">Traces</p>
            {#each debugReport.data.traces as t}
              <div class="flex gap-2 text-xs text-muted-foreground">
                <span>{t.step}</span>
                {#if t.trace_id}<Badge variant="outline" class="text-[0.6rem]">trace:{t.trace_id.slice(0, 8)}</Badge>{/if}
                {#if t.request_id}<Badge variant="outline" class="text-[0.6rem]">req:{t.request_id.slice(0, 8)}</Badge>{/if}
              </div>
            {/each}
          </div>
        {/if}
      </Card.Content>
    </Card.Root>
  {/if}

  <!-- Mermaid Sequence Diagram -->
  {#if mermaid.data}
    <Card.Root class="mb-4">
      <Card.Header><Card.Title class="text-base">Sequence Diagram</Card.Title></Card.Header>
      <Card.Content>
        <MermaidDiagram code={mermaid.data} />
      </Card.Content>
    </Card.Root>
  {/if}

  <!-- Live Events -->
  {#if events.length > 0}
    <Card.Root>
      <Card.Header><Card.Title class="text-base">Live Events ({events.length})</Card.Title></Card.Header>
      <Card.Content class="max-h-64 space-y-1 overflow-y-auto">
        {#each events as ev, i (i)}
          <div class="flex gap-2 rounded-md bg-secondary/50 px-2 py-1 text-xs">
            <Badge variant="outline" class="text-[0.6rem]">{ev.type}</Badge>
            <span class="truncate text-muted-foreground">{JSON.stringify(ev.data ?? '')}</span>
          </div>
        {/each}
      </Card.Content>
    </Card.Root>
  {:else if run.data?.status === 'pending' || run.data?.status === 'running'}
    <Card.Root>
      <Card.Content class="py-8 text-center text-sm text-muted-foreground">
        Waiting for events…
      </Card.Content>
    </Card.Root>
  {/if}
</main>
