<script lang="ts">
  import { page } from '$app/stores';
  import { createQuery } from '@tanstack/svelte-query';
  import { mockGetSpecDiff } from '$lib/api/mock-phase2';
  import * as Card from '$lib/components/ui/card';
  import { Badge } from '$lib/components/ui/badge';
  import { Button } from '$lib/components/ui/button';
  import { Plus, Minus, Pencil, AlertTriangle, ArrowLeft } from '@lucide/svelte';

  const projectId = $page.params.id!;
  const fromId = $page.url.searchParams.get('from') ?? '';
  const toId = $page.url.searchParams.get('to') ?? '';

  const diff = createQuery(() => ({
    queryKey: ['spec-diff', fromId, toId],
    queryFn: () => mockGetSpecDiff(fromId, toId),
    enabled: !!fromId && !!toId,
  }));

  function changeIcon(type: string) {
    if (type === 'added') return Plus;
    if (type === 'removed') return Minus;
    return Pencil;
  }

  function changeColor(type: string) {
    if (type === 'added') return 'text-emerald-400 bg-emerald-500/10';
    if (type === 'removed') return 'text-red-400 bg-red-500/10';
    return 'text-yellow-400 bg-yellow-500/10';
  }
</script>

<main class="mx-auto max-w-4xl px-6 py-8 space-y-6">
  <div>
    <a href="/projects/{projectId}/spec/versions" class="inline-flex items-center gap-1 text-sm text-muted-foreground hover:text-foreground">
      <ArrowLeft class="size-4" />
      Back to versions
    </a>
    <h1 class="mt-2 text-2xl font-bold">Spec Diff</h1>
    <p class="text-sm text-muted-foreground">Comparing changes between versions</p>
  </div>

  {#if !fromId || !toId}
    <Card.Root>
      <Card.Content class="py-8 text-center text-muted-foreground">
        Select two versions to compare from the <a href="/projects/{projectId}/spec/versions" class="underline">version history</a>.
      </Card.Content>
    </Card.Root>
  {:else if diff.isPending}
    <p class="text-muted-foreground">Loading diff…</p>
  {:else if diff.data}
    <!-- Summary -->
    <div class="flex gap-3">
      <Badge variant="secondary" class="gap-1"><Plus class="size-3" /> {diff.data.changes.filter(c => c.type === 'added').length} added</Badge>
      <Badge variant="secondary" class="gap-1"><Pencil class="size-3" /> {diff.data.changes.filter(c => c.type === 'modified').length} modified</Badge>
      <Badge variant="secondary" class="gap-1"><Minus class="size-3" /> {diff.data.changes.filter(c => c.type === 'removed').length} removed</Badge>
    </div>

    <!-- Changes -->
    <div class="space-y-3">
      {#each diff.data.changes as change (change.operation_id + change.type)}
        {@const Icon = changeIcon(change.type)}
        <Card.Root class={change.breaking ? 'border-destructive/40' : ''}>
          <Card.Content class="flex items-start gap-3 py-4">
            <div class={`mt-0.5 flex size-7 shrink-0 items-center justify-center rounded-md ${changeColor(change.type)}`}>
              <Icon class="size-4" />
            </div>
            <div class="flex-1 min-w-0">
              <div class="flex items-center gap-2">
                <Badge variant="outline" class="font-mono text-xs">{change.method.toUpperCase()}</Badge>
                <span class="font-mono text-sm">{change.path}</span>
                {#if change.breaking}
                  <Badge variant="destructive" class="gap-1 text-xs">
                    <AlertTriangle class="size-3" />
                    Breaking
                  </Badge>
                {/if}
              </div>
              <p class="mt-1 text-sm text-muted-foreground">{change.details}</p>
              <p class="mt-0.5 text-xs text-muted-foreground">Operation: {change.operation_id}</p>
            </div>
          </Card.Content>
        </Card.Root>
      {/each}
    </div>

    {#if diff.data.changes.some(c => c.breaking)}
      <Button href="/projects/{projectId}/impact?spec-version={toId}" variant="destructive" class="w-full">
        <AlertTriangle class="size-4" />
        View Impact Analysis
      </Button>
    {/if}
  {/if}
</main>
