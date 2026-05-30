<script lang="ts">
  import type { Snippet } from 'svelte';

  interface Props {
    title?: string;
    subtitle?: string;
    description?: string;
    actions?: Snippet;
    children?: Snippet;
    class?: string;
  }

  let { title, subtitle, description, actions, children, class: cls = '' }: Props = $props();
</script>

<div class={`flex flex-col gap-6 ${cls}`}>
  {#if title || actions}
    <div class="flex items-start justify-between gap-4">
      <div class="flex-1">
        {#if subtitle}
          <p class="text-xs font-bold uppercase tracking-widest text-primary">{subtitle}</p>
        {/if}
        {#if title}
          <h1 class="text-3xl font-bold tracking-tight">{title}</h1>
        {/if}
        {#if description}
          <p class="mt-2 text-sm text-muted-foreground">{description}</p>
        {/if}
      </div>
      {#if actions}
        <div class="flex items-center gap-2">
          {@render actions()}
        </div>
      {/if}
    </div>
  {/if}

  <div class="flex-1">
    {#if children}
      {@render children()}
    {/if}
  </div>
</div>
