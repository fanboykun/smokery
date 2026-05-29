<script lang="ts">
  import { page } from '$app/stores';
  import { createQuery, createMutation, useQueryClient } from '@tanstack/svelte-query';
  import { api } from '$lib/api/client';
  import * as Card from '$lib/components/ui/card';
  import { Badge } from '$lib/components/ui/badge';
  import { Button } from '$lib/components/ui/button';
  import { Input } from '$lib/components/ui/input';
  import { Textarea } from '$lib/components/ui/textarea';
  import { Separator } from '$lib/components/ui/separator';
  import Breadcrumb from '$lib/components/Breadcrumb.svelte';

  let filter = $state('');
  let selectedId = $state('');
  let overrideJson = $state('');
  let editClassification = $state('');
  let editDestructive = $state(false);

  const specs = createQuery(() => ({
    queryKey: ['specs', $page.params.id],
    queryFn: async () => {
      const { data, error } = await api.GET('/api/projects/{project-id}/specs', {
        params: { path: { 'project-id': $page.params.id! } },
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

  const filtered = $derived(
    (operations.data ?? []).filter((op) => {
      const needle = filter.toLowerCase().trim();
      if (!needle) return true;
      return [op.operation_id, op.path, op.method, op.classification, ...(op.tags ?? [])].some((v) =>
        v.toLowerCase().includes(needle),
      );
    }),
  );

  const grouped = $derived(
    filtered.reduce<Record<string, typeof filtered>>((acc, op) => {
      const tag = op.tags?.[0] ?? 'untagged';
      (acc[tag] ??= []).push(op);
      return acc;
    }, {}),
  );

  const selected = $derived((operations.data ?? []).find((op) => op.id === selectedId));

  $effect(() => {
    if (selected) {
      editClassification = selected.classification;
      editDestructive = selected.is_destructive;
      overrideJson = selected.overrides || '{}';
    }
  });

  const queryClient = useQueryClient();

  const updateClassification = createMutation(() => ({
    mutationFn: async () => {
      const { data, error } = await api.PUT('/api/operations/{id}/classification', {
        params: { path: { id: selectedId } },
        body: { classification: editClassification, is_destructive: editDestructive },
      });
      if (error) throw error;
      return data!;
    },
    onSuccess: () => queryClient.invalidateQueries({ queryKey: ['operations', latestSpecId] }),
  }));

  const updateOverrides = createMutation(() => ({
    mutationFn: async () => {
      const { data, error } = await api.PUT('/api/operations/{id}/overrides', {
        params: { path: { id: selectedId } },
        body: overrideJson,
        bodySerializer: (b) => b as unknown as string,
      });
      if (error) throw error;
      return data!;
    },
    onSuccess: () => queryClient.invalidateQueries({ queryKey: ['operations', latestSpecId] }),
  }));

  function methodColor(method: string) {
    switch (method.toUpperCase()) {
      case 'GET': return 'text-sky-400';
      case 'POST': return 'text-emerald-400';
      case 'PUT': return 'text-amber-400';
      case 'PATCH': return 'text-orange-400';
      case 'DELETE': return 'text-red-400';
      default: return 'text-muted-foreground';
    }
  }

  const tagNames = $derived(Object.keys(grouped));

  // Auto-scroll detail panel when selection changes
  let detailPanel: HTMLElement | undefined = $state();
  $effect(() => {
    if (selectedId && detailPanel) {
      detailPanel.scrollIntoView({ behavior: 'smooth', block: 'nearest' });
    }
  });
</script>

<main class="mx-auto max-w-7xl px-6 py-8">
  <Breadcrumb crumbs={[{ label: $page.params.id?.slice(0, 8) ?? '', href: `/projects/${$page.params.id}` }, { label: 'Operations' }]} />
  <div class="mb-6 flex flex-wrap items-end justify-between gap-4">
    <div>
      <p class="text-xs font-bold uppercase tracking-widest text-primary">Project {$page.params.id?.slice(0, 8)}</p>
      <h1 class="text-3xl font-bold">Operations</h1>
      <p class="text-sm text-muted-foreground">Classify operations, inspect hints, and configure overrides.</p>
    </div>
    <div class="flex gap-2">
      <Button variant="outline" href="/projects/{$page.params.id}">← Overview</Button>
      <Button href="/projects/{$page.params.id}/builder">Open Builder</Button>
    </div>
  </div>

  {#if specs.isLoading}
    <p class="text-muted-foreground">Loading specs…</p>
  {:else if !latestSpecId}
    <Card.Root>
      <Card.Content class="py-8 text-center">
        <p class="text-muted-foreground">No spec imported yet. Import an OpenAPI spec from the project overview to see operations.</p>
        <Button class="mt-4" href="/projects/{$page.params.id}">Go to Overview</Button>
      </Card.Content>
    </Card.Root>
  {:else}
    <div class="grid grid-cols-1 gap-6 lg:grid-cols-[1fr_400px]">
      <!-- Operation list -->
      <div class="space-y-4">
        <Input bind:value={filter} placeholder="Filter by id, path, tag, method, classification…" />

        {#if operations.isLoading}
          <p class="text-muted-foreground">Loading operations…</p>
        {:else}
          <p class="text-xs text-muted-foreground">{filtered.length} of {operations.data?.length ?? 0} operations</p>

          <!-- Sticky tag navigation -->
          {#if tagNames.length > 1}
            <nav class="sticky top-14 z-10 flex flex-wrap gap-1 rounded-md border border-border bg-background/90 p-2 backdrop-blur-sm">
              {#each tagNames as tag (tag)}
                <a href="#{tag}" class="rounded px-2 py-0.5 text-xs hover:bg-secondary">{tag}</a>
              {/each}
            </nav>
          {/if}

          {#each Object.entries(grouped) as [tag, ops] (tag)}
            <Card.Root id={tag}>
              <Card.Header class="flex-row items-center justify-between pb-2">
                <Card.Title class="text-base">{tag}</Card.Title>
                <Badge variant="secondary">{ops.length}</Badge>
              </Card.Header>
              <Card.Content class="space-y-0 p-2">
                {#each ops as op, i (op.id)}
                  <button
                    class="flex w-full items-center gap-3 rounded-md px-3 py-2 text-left text-sm transition-colors hover:bg-secondary {selectedId === op.id ? 'bg-secondary ring-1 ring-primary' : ''} {i % 2 === 1 ? 'bg-muted/30' : ''}"
                    onclick={() => (selectedId = op.id)}
                  >
                    <span class="w-14 shrink-0 font-mono text-xs font-bold {methodColor(op.method)}">{op.method}</span>
                    <span class="min-w-0 flex-1 truncate font-medium">{op.path}</span>
                    <Badge variant="outline" class="text-[0.65rem]">{op.classification}</Badge>
                    {#if op.is_destructive}
                      <Badge variant="destructive" class="text-[0.65rem]">protected</Badge>
                    {/if}
                  </button>
                {/each}
              </Card.Content>
            </Card.Root>
          {/each}
        {/if}
      </div>

      <!-- Detail panel -->
      <aside class="space-y-4 lg:sticky lg:top-20 lg:self-start" bind:this={detailPanel}>
        {#if selected}
          <Card.Root>
            <Card.Header>
              <div class="flex items-center gap-2">
                <span class="font-mono text-sm font-bold {methodColor(selected.method)}">{selected.method}</span>
                <Card.Title class="text-base">{selected.operation_id}</Card.Title>
              </div>
              <Card.Description>{selected.path}</Card.Description>
            </Card.Header>
            <Card.Content class="space-y-4">
              {#if selected.summary}
                <p class="text-sm text-muted-foreground">{selected.summary}</p>
              {/if}

              <Separator />

              <div class="space-y-2">
                <label class="text-xs font-bold uppercase tracking-wide text-muted-foreground">Classification</label>
                <select
                  class="w-full rounded-md border border-input bg-background px-3 py-2 text-sm"
                  bind:value={editClassification}
                >
                  <option value="list">list</option>
                  <option value="read">read</option>
                  <option value="create">create</option>
                  <option value="update">update</option>
                  <option value="delete">delete</option>
                  <option value="action">action</option>
                </select>
              </div>

              <label class="flex items-center gap-2 text-sm">
                <input type="checkbox" bind:checked={editDestructive} class="size-4 rounded border-input" />
                Mark as destructive (compiler protection)
              </label>

              <Button
                class="w-full"
                onclick={() => updateClassification.mutate()}
                disabled={updateClassification.isPending}
              >
                {updateClassification.isPending ? 'Saving…' : 'Save Classification'}
              </Button>

              <Separator />

              <div class="space-y-2">
                <label class="text-xs font-bold uppercase tracking-wide text-muted-foreground">Overrides JSON</label>
                <Textarea bind:value={overrideJson} class="min-h-[12rem] font-mono text-xs" />
              </div>

              <Button
                variant="outline"
                class="w-full"
                onclick={() => updateOverrides.mutate()}
                disabled={updateOverrides.isPending}
              >
                {updateOverrides.isPending ? 'Saving…' : 'Save Overrides'}
              </Button>
            </Card.Content>
          </Card.Root>
        {:else}
          <Card.Root>
            <Card.Content class="py-8 text-center text-sm text-muted-foreground">
              Select an operation to view details and configure overrides.
            </Card.Content>
          </Card.Root>
        {/if}
      </aside>
    </div>
  {/if}
</main>
