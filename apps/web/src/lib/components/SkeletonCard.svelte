<script lang="ts">
  interface Props {
    lines?: number;
    class?: string;
  }

  let { lines = 2, class: cls = '' } = $props();
</script>

<div class={`space-y-3 rounded-lg border border-border bg-card p-4 ${cls}`}>
  <div class="h-5 w-1/3 animate-pulse rounded bg-muted" />
  {#each Array(lines) as _}
    <div class="h-4 animate-pulse rounded bg-muted" style={`width: ${Math.random() * 40 + 60}%`} />
  {/each}
</div>
