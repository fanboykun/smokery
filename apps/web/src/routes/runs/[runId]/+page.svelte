<script lang="ts">
  import { page } from '$app/stores';
  import { createQuery } from '@tanstack/svelte-query';
  import { api } from '$lib/api/client';
  import { onMount } from 'svelte';
  import MermaidDiagram from '$lib/components/MermaidDiagram.svelte';
  import * as Card from '$lib/components/ui/card';
  import * as Tabs from '$lib/components/ui/tabs';
  import { Badge } from '$lib/components/ui/badge';
  import { Button } from '$lib/components/ui/button';

  interface RunEvent {
    type: string;
    run_id: string;
    data?: unknown;
  }

  let events = $state<RunEvent[]>([]);
  let statusFilter = $state('all');
  let stepsExpanded = $state(true);
  let activeTab = $state('timeline');
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

  const ciReport = createQuery(() => ({
    queryKey: ['run-ci', runId],
    queryFn: async () => {
      const { data, error } = await api.GET('/api/runs/{id}/report/ci', { params: { path: { id: runId } } });
      if (error) throw error;
      return data!;
    },
    enabled: run.data?.status === 'completed' || run.data?.status === 'failed',
  }));

  const runResult = createQuery(() => ({
    queryKey: ['run-result', runId],
    queryFn: async () => {
      const { data, error } = await api.GET('/api/runs/{id}/result', { params: { path: { id: runId } } });
      if (error) throw error;
      const raw = data!.result as string;
      try { return JSON.parse(raw); } catch {}
      try { return JSON.parse(atob(raw)); } catch {}
      return null;
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
    // WebSocket connection with fallback to polling
    let ws: WebSocket | null = null;
    let wsConnected = false;

    function connectWebSocket() {
      if (!(run.data as any)?.websocket_url) return; // Wait for websocket_url from API

      try {
        ws = new WebSocket((run.data as any).websocket_url);
        ws.onopen = () => {
          wsConnected = true;
          console.log('[v0] WebSocket connected');
        };
        ws.onmessage = (e) => {
          const event: RunEvent = JSON.parse(e.data);
          events = [...events, event];
        };
        ws.onerror = () => {
          wsConnected = false;
          console.log('[v0] WebSocket error, falling back to polling');
        };
        ws.onclose = () => {
          wsConnected = false;
        };
      } catch (err) {
        console.log('[v0] WebSocket failed:', err);
        wsConnected = false;
      }
    }

    // Try to connect when run data becomes available
    const checkConnection = () => {
      if (!wsConnected && (run.data as any)?.websocket_url) {
        connectWebSocket();
      }
    };

    const interval = setInterval(checkConnection, 1000);

    return () => {
      clearInterval(interval);
      ws?.close();
    };
  });

  function statusVariant(status: string) {
    if (status === 'completed') return 'default' as const;
    if (status === 'failed') return 'destructive' as const;
    return 'secondary' as const;
  }

  // Collect all steps from result for filtering
  interface FlatStep { method?: string; url?: string; status?: number; error?: string; stepStatus?: string; response?: unknown; request?: unknown; source: string }

  const allSteps = $derived.by(() => {
    if (!runResult.data) return [] as FlatStep[];
    const res = runResult.data;
    const steps: FlatStep[] = [];
    for (const suite of res.suites ?? []) {
      for (const c of suite.cases ?? []) {
        if (c.step) steps.push({ method: c.step.request?.method, url: c.step.request?.url, status: c.step.response?.status, error: c.step.error, stepStatus: c.step.status, response: c.step.response, request: c.step.request, source: suite.name });
      }
    }
    for (const flow of res.flows ?? []) {
      for (const step of flow.steps ?? []) {
        steps.push({ method: step.request?.method, url: step.request?.url, status: step.response?.status, error: step.error, stepStatus: step.status, response: step.response, request: step.request, source: flow.name });
      }
    }
    return steps;
  });

  // Group failures by status code
  const failuresByStatus = $derived.by(() => {
    const failures = allSteps.filter((s) => s.stepStatus !== 'passed');
    const groups: Record<string, FlatStep[]> = {};
    for (const f of failures) {
      const key = f.status ? String(f.status) : 'no-response';
      (groups[key] ??= []).push(f);
    }
    return groups;
  });

  const statusCodes = $derived(Object.keys(failuresByStatus).sort());

  const filteredSteps = $derived(
    statusFilter === 'all'
      ? allSteps
      : statusFilter === 'failed'
        ? allSteps.filter((s) => s.stepStatus !== 'passed')
        : allSteps.filter((s) => String(s.status) === statusFilter),
  );

  function shortUrl(url?: string) {
    if (!url) return '';
    return url.replace(/https?:\/\/[^/]+/, '');
  }
</script>

<main class="mx-auto max-w-5xl px-4 py-8 sm:px-6 overflow-x-hidden">
  <div class="mb-6 flex flex-wrap items-end justify-between gap-3">
    <div class="min-w-0">
      <p class="text-xs font-bold uppercase tracking-widest text-primary">Run Detail</p>
      <h1 class="truncate text-2xl font-bold font-mono">{runId.slice(0, 8)}</h1>
    </div>
    <div class="flex flex-wrap gap-2">
      {#if run.data}
        <Badge variant={statusVariant(run.data.status)}>{run.data.status}</Badge>
      {/if}
      <Button variant="outline" size="sm" href="/runs/{runId}/comments">Comments</Button>
    </div>
  </div>

  <!-- CI Summary -->
  {#if ciReport.data}
    <Card.Root class="mb-4">
      <Card.Header><Card.Title class="text-base">Summary</Card.Title></Card.Header>
      <Card.Content>
        <div class="flex flex-wrap gap-3 text-sm">
          <span>Status: <Badge variant={statusVariant(ciReport.data.status)}>{ciReport.data.status}</Badge></span>
          <span class="text-primary font-bold">{ciReport.data.passed} passed</span>
          <span class="text-destructive font-bold">{ciReport.data.failed} failed</span>
          <span class="text-muted-foreground">{ciReport.data.total} total • {ciReport.data.duration_ms}ms</span>
        </div>
        {#if ciReport.data.failures && ciReport.data.failures.length > 0}
          <details class="mt-3">
            <summary class="cursor-pointer text-xs font-bold text-destructive">Failed cases ({ciReport.data.failures.length})</summary>
            <div class="mt-2 max-h-48 overflow-y-auto space-y-1">
              {#each ciReport.data.failures as f}
                <p class="break-all text-xs text-muted-foreground font-mono">{f}</p>
              {/each}
            </div>
          </details>
        {/if}
      </Card.Content>
    </Card.Root>
  {/if}

  <!-- Results tabs -->
  <Tabs.Root value={activeTab} onValueChange={(v: string) => (activeTab = v)} class="mb-4">
    <Tabs.List class="grid w-full grid-cols-4">
      <Tabs.Trigger value="timeline">Timeline</Tabs.Trigger>
      <Tabs.Trigger value="debug">Debug</Tabs.Trigger>
      <Tabs.Trigger value="diagram">Diagram</Tabs.Trigger>
      <Tabs.Trigger value="events">Events</Tabs.Trigger>
    </Tabs.List>

    <!-- Timeline Tab -->
    <Tabs.Content value="timeline" class="mt-4">
      {#if allSteps.length > 0}
        <Card.Root class="mb-4">
          <Card.Header class="flex-row flex-wrap items-center justify-between gap-2">
            <Card.Title class="text-base">Step Results ({allSteps.length})</Card.Title>
        <div class="flex items-center gap-2">
          <select
            class="rounded-md border border-input bg-background px-2 py-1 text-xs"
            bind:value={statusFilter}
          >
            <option value="all">All</option>
            <option value="failed">Failed only</option>
            {#each statusCodes as code (code)}
              <option value={code}>{code} ({failuresByStatus[code].length})</option>
            {/each}
          </select>
          <button
            class="text-xs text-muted-foreground hover:text-foreground"
            onclick={() => stepsExpanded = !stepsExpanded}
          >
            {stepsExpanded ? '▼ Collapse' : '▶ Expand'}
          </button>
        </div>
      </Card.Header>
      {#if stepsExpanded}
        <Card.Content class="max-h-[32rem] overflow-y-auto space-y-1">
          {#each filteredSteps.slice(0, 50) as step}
            <div class="rounded-md px-2 py-1.5 text-xs {step.stepStatus === 'passed' ? 'bg-primary/10' : 'bg-destructive/10'}">
              <div class="flex flex-wrap items-center gap-1.5">
                <span class="shrink-0 font-mono font-bold">{step.method}</span>
                <span class="min-w-0 break-all">{shortUrl(step.url)}</span>
                {#if step.status}
                  <Badge variant={step.status < 400 ? 'default' : 'destructive'} class="ml-auto shrink-0 text-[0.6rem]">{step.status}</Badge>
                {/if}
              </div>
              {#if step.error}
                <p class="mt-0.5 break-all text-destructive">{step.error}</p>
              {/if}
              {#if step.stepStatus !== 'passed' && step.response}
                <details class="mt-1">
                  <summary class="cursor-pointer text-muted-foreground">Response body</summary>
                  <pre class="mt-1 max-h-32 overflow-auto whitespace-pre-wrap break-all rounded bg-background p-2 text-[0.65rem]">{JSON.stringify((step.response as any).body, null, 2)}</pre>
                </details>
              {/if}
            </div>
          {/each}
          {#if filteredSteps.length > 50}
            <p class="text-xs text-muted-foreground">… and {filteredSteps.length - 50} more</p>
          {/if}
        </Card.Content>
      {/if}
        </Card.Root>
      {/if}
    </Tabs.Content>

    <!-- Debug Tab -->
    <Tabs.Content value="debug" class="mt-4">
      {#if debugReport.data}
        <Card.Root class="mb-4">
          <Card.Header><Card.Title class="text-base">Debug Report</Card.Title></Card.Header>
      <Card.Content class="space-y-3">
        <div class="flex flex-wrap gap-3 text-sm">
          <span>Status: <Badge variant={statusVariant(debugReport.data.status)}>{debugReport.data.status}</Badge></span>
          <span class="text-muted-foreground">{debugReport.data.duration_ms}ms</span>
        </div>
        {#if debugReport.data.failures && debugReport.data.failures.length > 0}
          <div class="space-y-1">
            <p class="text-xs font-bold uppercase text-destructive">Failures</p>
            {#each debugReport.data.failures as f}
              <div class="rounded-md bg-destructive/10 p-2 text-xs">
                <p class="font-medium break-all">{f.step}: {f.assertion}</p>
                <p class="break-all text-muted-foreground">{f.message} — {f.request_url} ({f.status})</p>
              </div>
            {/each}
          </div>
        {/if}
        {#if debugReport.data.traces && debugReport.data.traces.length > 0}
          <div class="space-y-1">
            <p class="text-xs font-bold uppercase text-muted-foreground">Traces</p>
            {#each debugReport.data.traces as t}
              <div class="flex flex-wrap gap-2 text-xs text-muted-foreground">
                <span class="break-all">{t.step}</span>
                {#if t.trace_id}<Badge variant="outline" class="text-[0.6rem]">trace:{t.trace_id.slice(0, 8)}</Badge>{/if}
                {#if t.request_id}<Badge variant="outline" class="text-[0.6rem]">req:{t.request_id.slice(0, 8)}</Badge>{/if}
              </div>
            {/each}
          </div>
        {/if}
        </Card.Content>
        </Card.Root>
      {/if}
    </Tabs.Content>

    <!-- Diagram Tab -->
    <Tabs.Content value="diagram" class="mt-4">
      {#if mermaid.data}
        <Card.Root class="mb-4">
          <Card.Header><Card.Title class="text-base">Sequence Diagram</Card.Title></Card.Header>
      <Card.Content class="overflow-x-auto">
        <MermaidDiagram code={mermaid.data} />
      </Card.Content>
        </Card.Root>
      {/if}
    </Tabs.Content>

    <!-- Events Tab -->
    <Tabs.Content value="events" class="mt-4">
      {#if events.length > 0}
        <Card.Root>
          <Card.Header><Card.Title class="text-base">Live Events ({events.length})</Card.Title></Card.Header>
      <Card.Content class="max-h-64 space-y-1 overflow-y-auto">
        {#each events as ev, i (i)}
          <div class="flex flex-wrap gap-2 rounded-md bg-secondary/50 px-2 py-1 text-xs">
            <Badge variant="outline" class="text-[0.6rem]">{ev.type}</Badge>
            <span class="min-w-0 break-all text-muted-foreground">{JSON.stringify(ev.data ?? '')}</span>
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
    </Tabs.Content>
  </Tabs.Root>
</main>
