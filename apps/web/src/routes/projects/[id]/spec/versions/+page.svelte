<script lang="ts">
  import { page } from '$app/stores';
  import { createQuery } from '@tanstack/svelte-query';
  import { mockGetSpecVersions } from '$lib/api/mock-phase2';
  import * as Card from '$lib/components/ui/card';
  import { Badge } from '$lib/components/ui/badge';
  import { Button } from '$lib/components/ui/button';
  import { Upload, GitCompare, AlertTriangle } from '@lucide/svelte';

  const projectId = $page.params.id!;

  const versions = createQuery(() => ({
    queryKey: ['spec-versions', projectId],
    queryFn: () => mockGetSpecVersions(projectId),
  }));
</script>

<main class="mx-auto max-w-4xl px-6 py-8 space-y-6">
  <div class="flex items-center justify-between">
    <div>
      <p class="text-xs font-bold uppercase tracking-widest text-primary">Spec Evolution</p>
      <h1 class="text-2xl font-bold">Version History</h1>
    </div>
    <Button variant="outline" size="sm" href="/projects/{projectId}/spec">
      <Upload class="size-4" />
      Current Spec
    </Button>
  </div>

  {#if versions.isPending}
    <p class="text-muted-foreground">Loading versions…</p>
  {:else if versions.data}
    <div class="space-y-3">
      {#each versions.data as version, i (version.id)}
        <Card.Root class="transition-colors hover:border-primary/30">
          <Card.Content class="flex items-center gap-4 py-4">
            <div class="flex size-10 shrink-0 items-center justify-center rounded-full bg-primary/10">
              <span class="text-sm font-bold text-primary">{version.version}</span>
            </div>
            <div class="flex-1 min-w-0">
              <p class="font-medium">{version.summary}</p>
              <div class="mt-1 flex flex-wrap items-center gap-2 text-xs text-muted-foreground">
                <span>{version.uploaded_by}</span>
                <span>•</span>
                <span>{new Date(version.uploaded_at).toLocaleDateString()}</span>
                <span>•</span>
                <span>{version.operation_count} operations</span>
              </div>
            </div>
            <div class="flex items-center gap-2">
              {#if version.breaking_changes > 0}
                <Badge variant="destructive" class="gap-1">
                  <AlertTriangle class="size-3" />
                  {version.breaking_changes} breaking
                </Badge>
              {/if}
              {#if version.schema_changes > 0}
                <Badge variant="secondary">{version.schema_changes} changes</Badge>
              {/if}
              {#if i < versions.data.length - 1}
                <Button variant="ghost" size="sm" href="/projects/{projectId}/spec/diff?from={versions.data[i + 1].id}&to={version.id}">
                  <GitCompare class="size-4" />
                </Button>
              {/if}
            </div>
          </Card.Content>
        </Card.Root>
      {/each}
    </div>
  {/if}
</main>
