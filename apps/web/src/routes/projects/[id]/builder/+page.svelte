<script lang="ts">
  import { page } from "$app/stores";
  import { goto } from "$app/navigation";
  import { createQuery } from "@tanstack/svelte-query";
  import { SvelteFlow, Controls, Background, BackgroundVariant, type Edge, type Node } from "@xyflow/svelte";
  import "@xyflow/svelte/dist/style.css";
  import { api } from "$lib/api/client";
  import { createProjectConfigStore } from "$lib/stores/project-config";
  import { canvasGraphToProjectConfig } from "$lib/canvas/graph-to-config";
  import type { CanvasGraph, CanvasOperation, OperationNodeData, SuiteGeneratorNodeData } from "$lib/canvas/types";
  import OperationCanvasNode from "$lib/components/canvas/OperationCanvasNode.svelte";
  import SuiteGeneratorNode from "$lib/components/canvas/SuiteGeneratorNode.svelte";
  import * as Card from "$lib/components/ui/card";
  import { Badge } from "$lib/components/ui/badge";
  import { Button } from "$lib/components/ui/button";
  import { Input } from "$lib/components/ui/input";
  import { Loader2, AlertCircle, CheckCircle2, GitBranch, Play } from "@lucide/svelte";
  import type { components } from "$lib/api/v1";

  type CompilerOutput = components["schemas"]["PlanPreviewResponse"];
  type CanvasNode = Node<Record<string, unknown>>;
  type CanvasEdge = Edge<Record<string, unknown>>;

  const projectId = $page.params.id!;
  const config = createProjectConfigStore(projectId);
  const nodeTypes = { operationNode: OperationCanvasNode, suiteGeneratorNode: SuiteGeneratorNode };

  let search = $state("");
  let result = $state<CompilerOutput | null>(null);
  let previewError = $state<string | null>(null);
  let previewLoading = $state(false);
  let nodes = $state.raw<CanvasNode[]>(($config.canvas?.nodes as CanvasNode[] | undefined) ?? []);
  let edges = $state.raw<CanvasEdge[]>(($config.canvas?.edges as CanvasEdge[] | undefined) ?? []);
  let defaultEnvironment = $state($config.canvas?.default_environment ?? $config.environments[0]?.id ?? "");
  let defaultAuth = $state($config.canvas?.default_auth ?? "");

  const specs = createQuery(() => ({
    queryKey: ["specs", projectId],
    queryFn: async () => {
      const { data, error } = await api.GET("/api/projects/{project-id}/specs", { params: { path: { "project-id": projectId } } });
      if (error) throw error;
      return data ?? [];
    },
  }));

  const latestSpecId = $derived(specs.data?.at(-1)?.id);

  const operations = createQuery(() => ({
    queryKey: ["canvas-operations", latestSpecId],
    enabled: !!latestSpecId,
    queryFn: async () => {
      const { data, error } = await api.GET("/api/specs/{spec-id}/operations/canvas", { params: { path: { "spec-id": latestSpecId! } } });
      if (error) throw error;
      return (data ?? []) as CanvasOperation[];
    },
  }));

  const filteredOperations = $derived.by(() => {
    const needle = search.toLowerCase().trim();
    const ops = operations.data ?? [];
    if (!needle) return ops.slice(0, 30);
    return ops.filter((op) => [op.operation_id, op.path, op.method, op.summary, ...(op.tags ?? [])].some((value) => value?.toLowerCase().includes(needle))).slice(0, 30);
  });

  const graph = $derived<CanvasGraph>({
    version: 1,
    default_environment: defaultEnvironment,
    default_auth: defaultAuth || undefined,
    nodes: nodes.map((node) => ({ id: node.id, type: node.type as any, position: node.position, data: node.data })),
    edges: edges.map((edge) => ({ id: edge.id, type: (edge.type as any) || edgeType(edge.sourceHandle, edge.targetHandle), source: edge.source, target: edge.target, sourceHandle: edge.sourceHandle ?? undefined, targetHandle: edge.targetHandle ?? undefined, data: edge.data })),
  });

  const derivedConfig = $derived.by(() => canvasGraphToProjectConfig(graph, operations.data ?? [], { environments: $config.environments, auth_profiles: $config.auth_profiles }));
  const diagnostics = $derived(result?.diagnostics);
  const hasErrors = $derived((diagnostics?.errors?.length ?? 0) > 0);
  const totalCases = $derived(result?.plan?.suite_plans?.reduce((n, suite) => n + (suite.cases?.length ?? 0), 0) ?? 0);
  const totalSteps = $derived(result?.plan?.flow_plans?.reduce((n, flow) => n + (flow.steps?.length ?? 0), 0) ?? 0);

  function addOperationNode(op: CanvasOperation) {
    const id = `op-${op.operation_id}-${crypto.randomUUID().slice(0, 6)}`;
    nodes = [...nodes, {
      id,
      type: "operationNode",
      position: { x: 120 + nodes.length * 42, y: 120 + nodes.length * 18 },
      data: { operation_id: op.operation_id, label: op.operation_id, operation: op, destructive_acknowledged: !op.is_destructive } satisfies OperationNodeData & Record<string, unknown>,
    }];
  }

  function addSuiteNode() {
    const id = `suite-${crypto.randomUUID().slice(0, 6)}`;
    nodes = [...nodes, {
      id,
      type: "suiteGeneratorNode",
      position: { x: 180 + nodes.length * 36, y: 360 },
      data: {
        name: "List Endpoint Suite",
        selector: { tags: [], classifications: ["list"], paths: [], exclude: [] },
        strategy: { default_list: true, pagination: true, search_from_response: false, enum_filters: false, empty_result_policy: "allow", max_cases_per_op: 0 },
        matched_count: (operations.data ?? []).filter((op) => op.classification === "list").length,
        case_count: (operations.data ?? []).filter((op) => op.classification === "list").length * 2,
      } satisfies SuiteGeneratorNodeData & Record<string, unknown>,
    }];
  }

  function edgeType(sourceHandle?: string | null, targetHandle?: string | null) {
    if (targetHandle?.startsWith("suite:")) return "suiteSelection";
    if (sourceHandle === "flow-out" || targetHandle === "flow-in") return "sequence";
    return "dataLink";
  }

  function acknowledgeDestructive(nodeId: string) {
    nodes = nodes.map((node) => node.id === nodeId ? { ...node, data: { ...node.data, destructive_acknowledged: true } } : node);
  }

  function saveCanvas() {
    config.update((current) => ({ ...current, ...derivedConfig, canvas: graph }));
  }

  async function compilePreview() {
    saveCanvas();
    previewLoading = true;
    previewError = null;
    try {
      const { data, error } = await api.POST("/api/projects/{project-id}/plan/preview", {
        params: { path: { "project-id": projectId } },
        body: derivedConfig as any,
      });
      if (error) throw new Error((error as any).detail ?? "Compilation failed");
      result = data as CompilerOutput;
    } catch (error) {
      previewError = error instanceof Error ? error.message : "Compilation failed";
    } finally {
      previewLoading = false;
    }
  }

  async function startRun() {
    if (!result?.plan || hasErrors) return;
    const { data, error } = await api.POST("/api/projects/{project-id}/runs", {
      params: { path: { "project-id": projectId } },
      body: { plan_id: result.plan.id, plan: result.plan },
    });
    if (error) {
      previewError = (error as any).detail ?? "Failed to start run";
      return;
    }
    goto(`/runs/${data.id}`);
  }

  function methodClass(method: string) {
    const m = method.toUpperCase();
    if (m === "GET") return "bg-sky-500/20 text-sky-300";
    if (m === "POST") return "bg-emerald-500/20 text-emerald-300";
    if (m === "DELETE") return "bg-red-500/20 text-red-300";
    return "bg-amber-500/20 text-amber-300";
  }
</script>

<div class="grid h-[calc(100vh-3.5rem)] grid-cols-[320px_minmax(0,1fr)_360px] bg-background">
  <aside class="border-r border-border bg-card/60 p-4">
    <div class="mb-4">
      <p class="text-xs font-bold uppercase tracking-widest text-primary">Canvas Builder</p>
      <h1 class="text-xl font-bold">Operation Graph</h1>
      <p class="mt-1 text-xs text-muted-foreground">Connect response fields to request inputs to compose smoke flows.</p>
    </div>

    <div class="space-y-3">
      <div class="grid grid-cols-2 gap-2">
        <select class="rounded-md border border-input bg-background px-2 py-2 text-xs" bind:value={defaultEnvironment}>
          <option value="">Environment</option>
          {#each $config.environments as env (env.id)}<option value={env.id}>{env.name}</option>{/each}
        </select>
        <select class="rounded-md border border-input bg-background px-2 py-2 text-xs" bind:value={defaultAuth}>
          <option value="">No auth</option>
          {#each $config.auth_profiles as auth (auth.id)}<option value={auth.id}>{auth.name}</option>{/each}
        </select>
      </div>

      <Button class="w-full" variant="outline" onclick={addSuiteNode}>Add Suite Generator</Button>
      <Input placeholder="Search operations..." bind:value={search} />

      <div class="max-h-[calc(100vh-18rem)] space-y-2 overflow-y-auto pr-1">
        {#if operations.isPending}<p class="text-xs text-muted-foreground">Loading operations...</p>{/if}
        {#each filteredOperations as op (op.operation_id)}
          <button class="w-full rounded-lg border border-border bg-background/70 p-2 text-left transition hover:border-primary/60 hover:bg-secondary" onclick={() => addOperationNode(op)}>
            <div class="flex items-center gap-2">
              <Badge class={`${methodClass(op.method)} font-mono text-[0.65rem]`}>{op.method}</Badge>
              <span class="min-w-0 truncate text-sm font-medium">{op.operation_id}</span>
            </div>
            <p class="mt-1 truncate font-mono text-[0.7rem] text-muted-foreground">{op.path}</p>
          </button>
        {/each}
      </div>
    </div>
  </aside>

  <main class="relative min-w-0">
    <SvelteFlow bind:nodes bind:edges {nodeTypes} fitView minZoom={0.25} maxZoom={1.4} class="canvas-flow">
      <Controls />
      <Background variant={BackgroundVariant.Dots} gap={18} size={1} />
    </SvelteFlow>
    <div class="pointer-events-none absolute left-4 top-4 rounded-lg border border-border bg-card/90 px-3 py-2 text-xs text-muted-foreground shadow-lg">
      <span class="font-semibold text-foreground">{nodes.length}</span> nodes · <span class="font-semibold text-foreground">{edges.length}</span> edges
    </div>
  </main>

  <aside class="space-y-4 overflow-y-auto border-l border-border bg-card/60 p-4">
    <Card.Root>
      <Card.Header>
        <Card.Title class="flex items-center gap-2 text-base"><GitBranch class="size-4" /> Derived Config</Card.Title>
        <Card.Description>Generated from the canvas before compiler preview.</Card.Description>
      </Card.Header>
      <Card.Content class="grid grid-cols-2 gap-2 text-center text-sm">
        <div class="rounded border border-border bg-background/60 p-2"><p class="text-xl font-bold">{derivedConfig.flows.length}</p><p class="text-xs text-muted-foreground">Flows</p></div>
        <div class="rounded border border-border bg-background/60 p-2"><p class="text-xl font-bold">{derivedConfig.suites.length}</p><p class="text-xs text-muted-foreground">Suites</p></div>
      </Card.Content>
    </Card.Root>

    {#each nodes.filter((node) => (node.data as any).operation?.is_destructive && !(node.data as any).destructive_acknowledged) as node (node.id)}
      <Card.Root class="border-red-500/50 bg-red-500/5">
        <Card.Content class="space-y-2 py-3 text-sm">
          <p class="font-medium text-red-300">Destructive operation needs acknowledgement</p>
          <p class="font-mono text-xs text-muted-foreground">{(node.data as any).operation.operation_id}</p>
          <Button size="sm" variant="destructive" onclick={() => acknowledgeDestructive(node.id)}>Acknowledge</Button>
        </Card.Content>
      </Card.Root>
    {/each}

    <div class="grid grid-cols-2 gap-2">
      <Button variant="outline" onclick={saveCanvas}>Save Canvas</Button>
      <Button onclick={compilePreview} disabled={previewLoading || !defaultEnvironment}>
        {#if previewLoading}<Loader2 class="size-4 animate-spin" />{/if}
        Compile
      </Button>
    </div>

    {#if previewError}
      <Card.Root class="border-destructive/50 bg-destructive/5"><Card.Content class="flex gap-2 py-3 text-sm text-destructive"><AlertCircle class="size-4" /> {previewError}</Card.Content></Card.Root>
    {/if}

    {#if diagnostics}
      <Card.Root class={hasErrors ? "border-destructive/50" : "border-emerald-500/50"}>
        <Card.Content class="space-y-3 py-4">
          <div class="flex items-center gap-2">
            {#if hasErrors}<AlertCircle class="size-5 text-destructive" />{:else}<CheckCircle2 class="size-5 text-emerald-400" />{/if}
            <div><p class="font-semibold">{hasErrors ? "Fix compiler errors" : "Plan is compilable"}</p><p class="text-xs text-muted-foreground">{totalSteps} flow steps · {totalCases} suite cases</p></div>
          </div>
          {#each diagnostics.errors ?? [] as err}<p class="rounded bg-destructive/10 p-2 text-xs text-destructive">{err.message}</p>{/each}
          {#each diagnostics.warnings ?? [] as warning}<p class="rounded bg-amber-500/10 p-2 text-xs text-amber-300">{warning.message}</p>{/each}
        </Card.Content>
      </Card.Root>
    {/if}

    <Button class="w-full gap-2" disabled={!result?.plan || hasErrors} onclick={startRun}><Play class="size-4" /> Run Compiled Plan</Button>
  </aside>
</div>

<style>
  :global(.canvas-flow) {
    background: radial-gradient(circle at 20% 20%, color-mix(in oklch, var(--primary) 10%, transparent), transparent 28%), var(--background);
  }
  :global(.svelte-flow__edge-path) {
    stroke-width: 2;
  }
</style>
