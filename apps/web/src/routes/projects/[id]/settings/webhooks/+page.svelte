<script lang="ts">
  import { page } from '$app/stores';
  import { createQuery } from '@tanstack/svelte-query';
  import { mockGetWebhooks } from '$lib/api/mock-phase2';
  import * as Card from '$lib/components/ui/card';
  import { Badge } from '$lib/components/ui/badge';
  import { Button } from '$lib/components/ui/button';
  import { Webhook, Plus, CheckCircle2, XCircle } from '@lucide/svelte';

  const projectId = $page.params.id!;

  const webhooks = createQuery(() => ({
    queryKey: ['webhooks', projectId],
    queryFn: () => mockGetWebhooks(projectId),
  }));
</script>

<main class="mx-auto max-w-4xl px-6 py-8 space-y-6">
  <div class="flex items-center justify-between">
    <div>
      <p class="text-xs font-bold uppercase tracking-widest text-primary">Settings</p>
      <h1 class="text-2xl font-bold">Webhooks</h1>
    </div>
    <Button size="sm"><Plus class="size-4" /> Add Webhook</Button>
  </div>

  {#if webhooks.isPending}
    <p class="text-muted-foreground">Loading…</p>
  {:else if webhooks.data && webhooks.data.length > 0}
    <div class="space-y-3">
      {#each webhooks.data as wh (wh.id)}
        <Card.Root>
          <Card.Content class="flex items-center gap-4 py-4">
            <div class="flex size-10 shrink-0 items-center justify-center rounded-lg bg-muted">
              <Webhook class="size-5 text-muted-foreground" />
            </div>
            <div class="flex-1 min-w-0">
              <div class="flex items-center gap-2">
                <p class="font-medium">{wh.name}</p>
                {#if wh.is_active}
                  <CheckCircle2 class="size-4 text-emerald-400" />
                {:else}
                  <XCircle class="size-4 text-muted-foreground" />
                {/if}
              </div>
              <p class="text-xs text-muted-foreground font-mono truncate">{wh.url}</p>
              <div class="mt-1 flex flex-wrap gap-1">
                {#each wh.events as event}
                  <Badge variant="secondary" class="text-xs">{event}</Badge>
                {/each}
              </div>
            </div>
            <div class="text-right text-xs text-muted-foreground shrink-0">
              {#if wh.last_triggered_at}
                <p>Last triggered</p>
                <p>{new Date(wh.last_triggered_at).toLocaleDateString()}</p>
              {:else}
                <p>Never triggered</p>
              {/if}
            </div>
          </Card.Content>
        </Card.Root>
      {/each}
    </div>
  {:else}
    <Card.Root class="border-dashed">
      <Card.Content class="py-12 text-center">
        <Webhook class="mx-auto size-8 text-muted-foreground" />
        <p class="mt-2 font-medium">No webhooks configured</p>
        <p class="text-sm text-muted-foreground">Add a webhook to get notified about run events.</p>
      </Card.Content>
    </Card.Root>
  {/if}
</main>
