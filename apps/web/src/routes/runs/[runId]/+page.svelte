<script lang="ts">
  import { page } from '$app/stores';
  import { createQuery } from '@tanstack/svelte-query';
  import { api } from '$lib/api/client';
  import { onMount } from 'svelte';

  interface RunEvent {
    type: string;
    run_id: string;
    data?: unknown;
  }

  let events = $state<RunEvent[]>([]);

  const run = createQuery(() => ({
    queryKey: ['run', $page.params.runId],
    queryFn: async () => {
      const { data, error } = await api.GET('/api/runs/{id}', { params: { path: { id: $page.params.runId! } } });
      if (error) throw error;
      return data!;
    },
    refetchInterval: 3000,
  }));

  const debugReport = createQuery(() => ({
    queryKey: ['run-debug', $page.params.runId],
    queryFn: async () => {
      const { data, error } = await api.GET('/api/runs/{id}/report/debug', { params: { path: { id: $page.params.runId! } } });
      if (error) throw error;
      return data!;
    },
    enabled: run.data?.status === 'completed' || run.data?.status === 'failed',
  }));

  const mermaid = createQuery(() => ({
    queryKey: ['run-mermaid', $page.params.runId],
    queryFn: async () => {
      const { data, error } = await api.GET('/api/runs/{id}/report/mermaid', { params: { path: { id: $page.params.runId! } } });
      if (error) throw error;
      return data!;
    },
    enabled: run.data?.status === 'completed' || run.data?.status === 'failed',
  }));

  onMount(() => {
    const runId = $page.params.runId!;
    const ws = new WebSocket(`ws://localhost:8080/ws/runs/${runId}`);
    ws.onmessage = (e) => {
      const event: RunEvent = JSON.parse(e.data);
      events = [...events, event];
    };
    return () => ws.close();
  });
</script>

<h1>Run {$page.params.runId}</h1>
{#if run.data}
  <p>Status: <strong>{run.data.status}</strong></p>
{/if}

{#if debugReport.data}
  <h2>Debug Report</h2>
  <pre>{JSON.stringify(debugReport.data, null, 2)}</pre>
{/if}

{#if mermaid.data}
  <h2>Sequence Diagram</h2>
  <pre>{mermaid.data}</pre>
{/if}

<h2>Live Events</h2>
<ul>
  {#each events as ev, i (i)}
    <li>{ev.type}: {JSON.stringify(ev.data ?? '')}</li>
  {/each}
</ul>
