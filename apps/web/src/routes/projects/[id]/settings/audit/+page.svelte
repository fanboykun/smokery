<script lang="ts">
  import { page } from '$app/stores';
  import { createQuery } from '@tanstack/svelte-query';
  import { mockGetAuditLog } from '$lib/api/mock-phase2';
  import * as Card from '$lib/components/ui/card';
  import { Badge } from '$lib/components/ui/badge';
  import { Pencil, Play, Trash2, Plus } from '@lucide/svelte';

  const projectId = $page.params.id!;

  const audit = createQuery(() => ({
    queryKey: ['audit-log', projectId],
    queryFn: () => mockGetAuditLog(projectId),
  }));

  function actionIcon(action: string) {
    if (action === 'update') return Pencil;
    if (action === 'run') return Play;
    if (action === 'delete') return Trash2;
    return Plus;
  }

  function actionColor(action: string) {
    if (action === 'delete') return 'text-red-400';
    if (action === 'run') return 'text-blue-400';
    return 'text-muted-foreground';
  }
</script>

<main class="mx-auto max-w-4xl px-6 py-8 space-y-6">
  <div>
    <p class="text-xs font-bold uppercase tracking-widest text-primary">Settings</p>
    <h1 class="text-2xl font-bold">Audit Log</h1>
    <p class="text-sm text-muted-foreground">All actions performed on this project</p>
  </div>

  {#if audit.isPending}
    <p class="text-muted-foreground">Loading…</p>
  {:else if audit.data}
    <Card.Root>
      <Card.Content class="divide-y divide-border p-0">
        {#each audit.data as entry (entry.id)}
          {@const Icon = actionIcon(entry.action)}
          <div class="flex items-start gap-3 px-4 py-3">
            <div class={`mt-0.5 flex size-7 shrink-0 items-center justify-center rounded-md bg-muted ${actionColor(entry.action)}`}>
              <Icon class="size-4" />
            </div>
            <div class="flex-1 min-w-0">
              <p class="text-sm">
                <span class="font-medium">{entry.actor_name}</span>
                <span class="text-muted-foreground"> {entry.action} </span>
                <Badge variant="outline" class="text-xs">{entry.resource_type}</Badge>
              </p>
              {#if entry.changes}
                <div class="mt-1 space-y-0.5">
                  {#each entry.changes as change}
                    <p class="text-xs text-muted-foreground font-mono">
                      {change.field}: <span class="text-red-400 line-through">{change.old_value}</span> → <span class="text-emerald-400">{change.new_value}</span>
                    </p>
                  {/each}
                </div>
              {/if}
            </div>
            <span class="text-xs text-muted-foreground shrink-0">{new Date(entry.timestamp).toLocaleString()}</span>
          </div>
        {/each}
      </Card.Content>
    </Card.Root>
  {/if}
</main>
