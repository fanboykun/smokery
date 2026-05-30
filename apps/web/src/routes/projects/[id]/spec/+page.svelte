<script lang="ts">
  import { page } from '$app/stores';
  import { createQuery } from '@tanstack/svelte-query';
  import { api } from '$lib/api/client';
  import PageLayout from '$lib/components/PageLayout.svelte';
  import * as Card from '$lib/components/ui/card';
  import * as Tabs from '$lib/components/ui/tabs';
  import { Badge } from '$lib/components/ui/badge';
  import { Input } from '$lib/components/ui/input';
  import { ChevronDown, ChevronRight, Search } from '@lucide/svelte';

  const projectId = $page.params.id!;
  let searchQuery = $state('');
  let expandedOps = $state<Set<string>>(new Set());

  const specs = createQuery(() => ({
    queryKey: ['specs', projectId],
    queryFn: async () => {
      const { data, error } = await api.GET('/api/projects/{project-id}/specs', {
        params: { path: { 'project-id': projectId } },
      });
      if (error) throw error;
      return data ?? [];
    },
  }));

  const latestSpecId = $derived(specs.data?.at(-1)?.id);

  const operations = createQuery(() => ({
    queryKey: ['operations', latestSpecId],
    queryFn: async () => {
      const { data, error } = await api.GET('/api/specs/{spec-id}/operations', {
        params: { path: { 'spec-id': latestSpecId! } },
      });
      if (error) throw error;
      return data ?? [];
    },
    enabled: !!latestSpecId,
  }));

  const filteredOps = $derived.by(() => {
    const ops = operations.data ?? [];
    if (!searchQuery.trim()) return ops;
    const q = searchQuery.toLowerCase();
    return ops.filter(
      (op) =>
        op.operation_id.toLowerCase().includes(q) ||
        op.path.toLowerCase().includes(q) ||
        op.method.toLowerCase().includes(q),
    );
  });

  const opsByTag = $derived.by(() => {
    const grouped: Record<string, typeof filteredOps> = {};
    for (const op of filteredOps) {
      const tags = op.tags ?? ['Untagged'];
      for (const tag of tags) {
        if (!grouped[tag]) grouped[tag] = [];
        grouped[tag].push(op);
      }
    }
    return grouped;
  });

  function toggleOp(opId: string) {
    const next = new Set(expandedOps);
    if (next.has(opId)) {
      next.delete(opId);
    } else {
      next.add(opId);
    }
    expandedOps = next;
  }

  function methodColor(method: string): string {
    const m = method.toUpperCase();
    if (m === 'GET') return 'bg-blue-500/20 text-blue-300';
    if (m === 'POST') return 'bg-green-500/20 text-green-300';
    if (m === 'PUT') return 'bg-yellow-500/20 text-yellow-300';
    if (m === 'DELETE') return 'bg-red-500/20 text-red-300';
    if (m === 'PATCH') return 'bg-purple-500/20 text-purple-300';
    return 'bg-gray-500/20 text-gray-300';
  }
</script>

<PageLayout title="OpenAPI Spec" subtitle="Explore your imported specification" class="max-w-6xl mx-auto">
  {#if specs.isPending}
    <div class="text-center py-8 text-muted-foreground">Loading specs…</div>
  {:else if specs.isError}
    <Card.Root class="border-destructive/50 bg-destructive/5">
      <Card.Content class="py-4 text-sm text-destructive">Failed to load specs</Card.Content>
    </Card.Root>
  {:else if !latestSpecId}
    <Card.Root class="border-dashed">
      <Card.Content class="py-8 text-center text-sm text-muted-foreground">
        No specs imported yet. Upload an OpenAPI spec in your project settings.
      </Card.Content>
    </Card.Root>
  {:else}
    <!-- Search -->
    <Card.Root class="mb-6">
      <Card.Content class="pt-6">
        <div class="relative">
          <Search class="absolute left-3 top-3 size-4 text-muted-foreground" />
          <Input bind:value={searchQuery} placeholder="Search operations by ID, path, or method…" class="pl-10" />
        </div>
      </Card.Content>
    </Card.Root>

    <!-- Results -->
    {#if operations.isPending}
      <div class="text-center py-8 text-muted-foreground">Loading operations…</div>
    {:else if filteredOps.length === 0}
      <Card.Root class="border-dashed">
        <Card.Content class="py-8 text-center text-sm text-muted-foreground">
          No operations match "{searchQuery}"
        </Card.Content>
      </Card.Root>
    {:else}
      <!-- Grouped by tag -->
      <div class="space-y-4">
        {#each Object.entries(opsByTag) as [tag, ops] (tag)}
          <Card.Root>
            <Card.Header class="pb-3">
              <div class="flex items-center justify-between">
                <Card.Title class="text-base">{tag}</Card.Title>
                <Badge variant="secondary">{ops.length} operations</Badge>
              </div>
            </Card.Header>
            <Card.Content class="space-y-2">
              {#each ops as op (op.id)}
                <div class="rounded-lg border border-border bg-muted/30">
                  <button
                    onclick={() => toggleOp(op.id)}
                    class="w-full flex items-center gap-2 p-3 hover:bg-muted/50 transition-colors text-left"
                  >
                    <div class="flex-shrink-0">
                      {#if expandedOps.has(op.id)}
                        <ChevronDown class="size-4" />
                      {:else}
                        <ChevronRight class="size-4" />
                      {/if}
                    </div>
                    <Badge class={`${methodColor(op.method)} font-mono font-bold text-xs shrink-0`} variant="secondary">
                      {op.method.toUpperCase()}
                    </Badge>
                    <span class="font-mono text-sm flex-1 truncate">{op.operation_id}</span>
                    {#if op.is_destructive}
                      <Badge variant="destructive" class="text-xs shrink-0">Destructive</Badge>
                    {/if}
                  </button>

                  {#if expandedOps.has(op.id)}
                    <div class="border-t border-border px-3 py-3 space-y-3 bg-background/50 text-xs">
                      <!-- Path -->
                      <div>
                        <p class="font-semibold text-muted-foreground mb-1">Path</p>
                        <code class="block bg-background rounded px-2 py-1 font-mono text-[0.7rem] break-all overflow-auto max-h-20">
                          {op.path}
                        </code>
                      </div>

                      <!-- Request/Response Summary -->
                      {#if op.summary}
                        <div>
                          <p class="font-semibold text-muted-foreground mb-1">Summary</p>
                          <p class="text-muted-foreground">{op.summary}</p>
                        </div>
                      {/if}

                      <!-- Classification -->
                      {#if op.classification}
                        <div>
                          <p class="font-semibold text-muted-foreground mb-1">Classification</p>
                          <Badge variant="outline">{op.classification}</Badge>
                        </div>
                      {/if}

                      <!-- Info -->
                      <div class="flex gap-4 flex-wrap">
                        {#if op.tags && op.tags.length > 0}
                          <div>
                            <p class="font-semibold text-muted-foreground mb-1">Tags</p>
                            <div class="flex gap-1 flex-wrap">
                              {#each op.tags as t}
                                <Badge variant="outline" class="text-xs">{t}</Badge>
                              {/each}
                            </div>
                          </div>
                        {/if}
                      </div>
                    </div>
                  {/if}
                </div>
              {/each}
            </Card.Content>
          </Card.Root>
        {/each}
      </div>
    {/if}
  {/if}
</PageLayout>
