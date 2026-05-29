<script lang="ts">
  import { GripVertical } from '@lucide/svelte';

  interface Props {
    leftLabel?: string;
    rightLabel?: string;
    class?: string;
  }

  let { leftLabel = 'Config', rightLabel = 'Preview', class: cls = '' } = $props();
  
  let isDragging = $state(false);
  let splitPos = $state(50);

  function handleMouseDown() {
    isDragging = true;
  }

  function handleMouseUp() {
    isDragging = false;
  }

  function handleMouseMove(e: MouseEvent) {
    if (!isDragging) return;

    const container = document.querySelector('[data-split-container]') as HTMLElement;
    if (!container) return;

    const rect = container.getBoundingClientRect();
    const newPos = Math.max(20, Math.min(80, ((e.clientX - rect.left) / rect.width) * 100));
    splitPos = newPos;
  }
</script>

<svelte:document onmouseup={handleMouseUp} onmousemove={handleMouseMove} />

<div 
  data-split-container
  class={`flex h-full overflow-hidden select-none ${cls}`}
>
  <!-- Left Panel -->
  <div class="flex flex-col overflow-hidden" style={`width: ${splitPos}%`}>
    {#if leftLabel}
      <div class="border-b border-border bg-card px-4 py-2">
        <p class="text-xs font-semibold uppercase tracking-widest text-muted-foreground">{leftLabel}</p>
      </div>
    {/if}
    <div class="flex-1 overflow-auto">
      <slot name="left" />
    </div>
  </div>

  <!-- Divider -->
  <div
    class="w-1 cursor-col-resize bg-muted transition-colors hover:bg-primary/50 {isDragging ? 'bg-primary/50' : ''}"
    onmousedown={handleMouseDown}
    role="separator"
    aria-label="Resize panes"
    aria-valuenow={splitPos}
    aria-valuemin={20}
    aria-valuemax={80}
    tabindex="0"
  />

  <!-- Right Panel -->
  <div class="flex flex-col overflow-hidden" style={`width: ${100 - splitPos}%`}>
    {#if rightLabel}
      <div class="border-b border-border bg-muted/30 px-4 py-2">
        <p class="text-xs font-semibold uppercase tracking-widest text-muted-foreground">{rightLabel}</p>
      </div>
    {/if}
    <div class="flex-1 overflow-auto bg-muted/10">
      <slot name="right" />
    </div>
  </div>
</div>

<style>
  :global([data-split-container]) {
    user-select: none;
  }
</style>
