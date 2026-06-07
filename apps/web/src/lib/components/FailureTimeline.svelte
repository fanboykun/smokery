<script lang="ts">
  import { Badge } from '$lib/components/ui/badge';
  import { Tag, UserPlus, ArrowRightLeft } from '@lucide/svelte';
  import type { FailureAction } from '$lib/api/phase2';

  interface Props {
    actions: FailureAction[];
  }

  let { actions }: Props = $props();

  function actionIcon(type: string) {
    if (type === 'classified') return Tag;
    if (type === 'assigned') return UserPlus;
    return ArrowRightLeft;
  }

  function actionLabel(action: FailureAction): string {
    const d = action.details ?? {};
    if (action.action_type === 'classified') return `classified as ${d.classification}`;
    if (action.action_type === 'assigned') return `assigned to ${d.assigned_to_name}`;
    if (action.action_type === 'status_changed') return `changed status ${d.from_status} → ${d.to_status}`;
    return action.action_type;
  }

  function timeAgo(iso: string): string {
    const diff = Date.now() - new Date(iso).getTime();
    const mins = Math.floor(diff / 60000);
    if (mins < 60) return `${mins}m ago`;
    const hrs = Math.floor(mins / 60);
    if (hrs < 24) return `${hrs}h ago`;
    return `${Math.floor(hrs / 24)}d ago`;
  }
</script>

{#if actions.length === 0}
  <p class="text-xs text-muted-foreground">No actions yet</p>
{:else}
  <div class="space-y-3">
    {#each actions as action (action.id)}
      {@const Icon = actionIcon(action.action_type)}
      <div class="flex items-start gap-3">
        <div class="mt-0.5 flex size-6 shrink-0 items-center justify-center rounded-full bg-muted">
          <Icon class="size-3 text-muted-foreground" />
        </div>
        <div class="flex-1 min-w-0">
          <p class="text-sm">
            <span class="font-medium">{action.actor_name}</span>
            <span class="text-muted-foreground"> {actionLabel(action)}</span>
          </p>
          <p class="text-xs text-muted-foreground">{timeAgo(action.created_at)}</p>
        </div>
      </div>
    {/each}
  </div>
{/if}
