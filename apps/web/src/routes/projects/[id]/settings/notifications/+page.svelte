<script lang="ts">
  import { page } from '$app/stores';
  import { createQuery } from '@tanstack/svelte-query';
  import { mockGetNotificationRules } from '$lib/api/mock-phase2';
  import * as Card from '$lib/components/ui/card';
  import { Badge } from '$lib/components/ui/badge';
  import { Button } from '$lib/components/ui/button';
  import { Bell, Plus, Mail, MessageSquare, CheckCircle2, XCircle } from '@lucide/svelte';

  const projectId = $page.params.id!;

  const rules = createQuery(() => ({
    queryKey: ['notification-rules', projectId],
    queryFn: () => mockGetNotificationRules(projectId),
  }));

  function channelIcon(channel: string) {
    if (channel === 'email') return Mail;
    if (channel === 'slack') return MessageSquare;
    return Bell;
  }
</script>

<main class="mx-auto max-w-4xl px-6 py-8 space-y-6">
  <div class="flex items-center justify-between">
    <div>
      <p class="text-xs font-bold uppercase tracking-widest text-primary">Settings</p>
      <h1 class="text-2xl font-bold">Notifications</h1>
    </div>
    <Button size="sm"><Plus class="size-4" /> Add Rule</Button>
  </div>

  {#if rules.isPending}
    <p class="text-muted-foreground">Loading…</p>
  {:else if rules.data && rules.data.length > 0}
    <div class="space-y-3">
      {#each rules.data as rule (rule.id)}
        {@const ChannelIcon = channelIcon(rule.channel)}
        <Card.Root>
          <Card.Content class="flex items-center gap-4 py-4">
            <div class="flex size-10 shrink-0 items-center justify-center rounded-lg bg-muted">
              <ChannelIcon class="size-5 text-muted-foreground" />
            </div>
            <div class="flex-1 min-w-0">
              <div class="flex items-center gap-2">
                <p class="font-medium">{rule.name}</p>
                {#if rule.is_active}
                  <CheckCircle2 class="size-4 text-emerald-400" />
                {:else}
                  <XCircle class="size-4 text-muted-foreground" />
                {/if}
              </div>
              <div class="mt-1 flex flex-wrap gap-1">
                <Badge variant="outline" class="text-xs capitalize">{rule.channel}</Badge>
                {#each rule.triggers as trigger}
                  <Badge variant="secondary" class="text-xs">{trigger.event}</Badge>
                {/each}
              </div>
            </div>
            <span class="text-xs text-muted-foreground shrink-0">
              Created {new Date(rule.created_at).toLocaleDateString()}
            </span>
          </Card.Content>
        </Card.Root>
      {/each}
    </div>
  {:else}
    <Card.Root class="border-dashed">
      <Card.Content class="py-12 text-center">
        <Bell class="mx-auto size-8 text-muted-foreground" />
        <p class="mt-2 font-medium">No notification rules</p>
        <p class="text-sm text-muted-foreground">Create rules to get alerted on run failures and other events.</p>
      </Card.Content>
    </Card.Root>
  {/if}
</main>
