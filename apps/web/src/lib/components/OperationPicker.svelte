<script lang="ts">
  import { createQuery } from '@tanstack/svelte-query';
  import { api } from '$lib/api/client';
  import * as Select from '$lib/components/ui/select';
  import { Badge } from '$lib/components/ui/badge';
  import { Input } from '$lib/components/ui/input';
  import { Label } from '$lib/components/ui/label';

  interface Props {
    projectId: string;
    specId: string | undefined;
    value: string;
    onchange: (value: string) => void;
    label?: string;
    placeholder?: string;
  }

  let {
    projectId,
    specId,
    value = $bindable(),
    onchange,
    label = 'Operation',
    placeholder = 'Select an operation…',
  }: Props = $props();

  let searchInput = $state('');
  let isOpen = $state(false);

  // Fetch specs for the project
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

  // Use provided specId or latest spec
  const latestSpecId = $derived(specId || specs.data?.at(-1)?.id);

  // Fetch operations from spec
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

  // Group operations by tag
  const grouped = $derived.by(() => {
    const ops = operations.data ?? [];
    const needle = searchInput.toLowerCase().trim();
    const filtered = ops.filter((op) => {
      if (!needle) return true;
      return [op.operation_id, op.path, op.method, ...(op.tags ?? [])].some((v) =>
        v?.toLowerCase().includes(needle),
      );
    });

    return filtered.reduce<Record<string, typeof filtered>>((acc, op) => {
      const tag = op.tags?.[0] ?? 'untagged';
      (acc[tag] ??= []).push(op);
      return acc;
    }, {});
  });

  const selectedOp = $derived(operations.data?.find((op) => op.operation_id === value));

  function methodColor(method: string) {
    const m = method.toUpperCase();
    if (m === 'GET') return 'bg-blue-500/20 text-blue-300';
    if (m === 'POST') return 'bg-green-500/20 text-green-300';
    if (m === 'PUT') return 'bg-yellow-500/20 text-yellow-300';
    if (m === 'DELETE') return 'bg-red-500/20 text-red-300';
    if (m === 'PATCH') return 'bg-purple-500/20 text-purple-300';
    return 'bg-gray-500/20 text-gray-300';
  }
</script>

<div class="space-y-2">
  {#if label}
    <Label>{label}</Label>
  {/if}

  <div class="relative">
    <!-- Selected value display -->
    {#if selectedOp}
      <div class="flex items-center gap-2 rounded-md border border-input bg-muted/50 px-3 py-2">
        <Badge class={`${methodColor(selectedOp.method)} font-mono text-xs`}>
          {selectedOp.method.toUpperCase()}
        </Badge>
        <span class="flex-1 text-sm font-medium">{selectedOp.operation_id}</span>
        <span class="text-xs text-muted-foreground">{selectedOp.path}</span>
        <button
          class="text-muted-foreground hover:text-foreground"
          onclick={() => {
            value = '';
            onchange('');
            searchInput = '';
            isOpen = false;
          }}
        >
          ×
        </button>
      </div>
    {:else}
      <!-- Search input -->
      <Input
        type="text"
        placeholder={placeholder}
        value={searchInput}
        oninput={(e) => {
          searchInput = e.currentTarget.value;
          isOpen = true;
        }}
        onfocus={() => (isOpen = true)}
        onblur={() => setTimeout(() => (isOpen = false), 200)}
        class="text-sm"
      />
    {/if}

    <!-- Dropdown -->
    {#if isOpen && (operations.data?.length ?? 0) > 0}
      <div class="absolute top-full z-50 mt-1 w-full max-h-64 overflow-y-auto rounded-md border border-input bg-popover shadow-md">
        {#if Object.keys(grouped).length === 0}
          <div class="px-3 py-2 text-xs text-muted-foreground">No operations found</div>
        {:else}
          {#each Object.entries(grouped) as [tag, ops] (tag)}
            <div>
              <div class="sticky top-0 bg-muted/80 px-3 py-1.5 text-xs font-semibold uppercase tracking-wider text-muted-foreground">
                {tag}
              </div>
              {#each ops as op (op.id)}
                <button
                  class="w-full flex items-center gap-2 px-3 py-2 text-left text-sm hover:bg-secondary transition-colors"
                  onclick={() => {
                    value = op.operation_id;
                    onchange(op.operation_id);
                    searchInput = '';
                    isOpen = false;
                  }}
                >
                  <Badge class={`${methodColor(op.method)} font-mono text-xs shrink-0`}>
                    {op.method.toUpperCase()}
                  </Badge>
                  <div class="flex-1 min-w-0">
                    <div class="font-medium truncate">{op.operation_id}</div>
                    <div class="text-xs text-muted-foreground truncate">{op.path}</div>
                  </div>
                  {#if op.is_destructive}
                    <span class="text-xs text-red-400 shrink-0">destr</span>
                  {/if}
                </button>
              {/each}
            </div>
          {/each}
        {/if}
      </div>
    {/if}

    <!-- Loading state -->
    {#if operations.isPending && !selectedOp}
      <div class="absolute top-full z-50 mt-1 w-full rounded-md border border-input bg-popover px-3 py-2 text-xs text-muted-foreground">
        Loading operations…
      </div>
    {/if}
  </div>

  <!-- Operation details (if selected) -->
  {#if selectedOp}
    <div class="rounded-md bg-muted/30 p-3 text-xs space-y-2">
      <div class="flex gap-2">
        {#if selectedOp.classification}
          <Badge variant="outline">{selectedOp.classification}</Badge>
        {/if}
        {#each selectedOp.tags ?? [] as tag (tag)}
          <Badge variant="secondary">{tag}</Badge>
        {/each}
      </div>
      {#if selectedOp.is_destructive}
        <div class="text-red-300">⚠ Destructive operation - handle with care</div>
      {/if}
    </div>
  {/if}
</div>
