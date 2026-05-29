<script lang="ts">
  type Crumb = { label: string; href?: string };
  let { crumbs }: { crumbs: Crumb[] } = $props();
</script>

<nav aria-label="Breadcrumb" class="mb-4 flex items-center gap-1.5 text-xs text-muted-foreground">
  <a href="/projects" class="hover:text-foreground transition-colors">Projects</a>
  {#each crumbs as crumb, i (i)}
    <span class="opacity-50">/</span>
    {#if crumb.href}
      <a href={crumb.href} class="hover:text-foreground transition-colors">{crumb.label}</a>
    {:else}
      <span class="text-foreground font-medium">{crumb.label}</span>
    {/if}
  {/each}
</nav>
