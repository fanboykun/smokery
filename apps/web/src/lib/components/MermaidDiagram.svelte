<script lang="ts">
  import { onMount } from 'svelte';

  let { code }: { code: string } = $props();
  let container: HTMLDivElement;

  onMount(async () => {
    const mermaid = (await import('mermaid')).default;
    mermaid.initialize({ startOnLoad: false, theme: 'dark' });
    const { svg } = await mermaid.render('mermaid-' + Math.random().toString(36).slice(2), code);
    container.innerHTML = svg;
  });
</script>

<div bind:this={container} class="overflow-x-auto rounded-lg border border-border bg-card p-4"></div>
