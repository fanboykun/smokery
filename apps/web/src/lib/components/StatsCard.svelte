<script lang="ts">
  import * as Card from '$lib/components/ui/card';
  import { Badge } from '$lib/components/ui/badge';

  interface Stat {
    label: string;
    value: string | number;
    icon?: any;
    color?: 'default' | 'success' | 'warning' | 'error';
  }

  interface Props {
    stats: Stat[];
    class?: string;
  }

  let { stats, class: cls = '' } = $props();

  const colorMap = {
    default: 'bg-muted text-muted-foreground',
    success: 'bg-emerald-500/10 text-emerald-400',
    warning: 'bg-yellow-500/10 text-yellow-400',
    error: 'bg-red-500/10 text-red-400',
  };
</script>

<Card.Root class={cls}>
  <Card.Content class="pt-6">
    <div class="grid gap-4 sm:grid-cols-2 lg:grid-cols-4">
      {#each stats as stat (stat.label)}
        <div class="space-y-2">
          <p class="text-xs font-medium text-muted-foreground uppercase tracking-widest">{stat.label}</p>
          <div class="flex items-baseline gap-2">
            <span class="text-2xl font-bold">{stat.value}</span>
            {#if stat.color}
              <span class={`inline-block size-2 rounded-full ${colorMap[stat.color as keyof typeof colorMap]}`} />
            {/if}
          </div>
        </div>
      {/each}
    </div>
  </Card.Content>
</Card.Root>
