<script lang="ts">
  import * as Card from '$lib/components/ui/card';
  import { Button } from '$lib/components/ui/button';
  import { AlertCircle } from '@lucide/svelte';

  interface Props {
    error?: Error | null;
    resetError?: () => void;
    class?: string;
  }

  let { error = null, resetError, class: cls = '' } = $props();
</script>

{#if error}
  <Card.Root class={`border-destructive ${cls}`}>
    <Card.Content class="pt-6">
      <div class="flex gap-4">
        <AlertCircle class="size-5 flex-shrink-0 text-destructive mt-0.5" />
        <div class="flex-1">
          <p class="font-semibold text-destructive">Something went wrong</p>
          <p class="mt-1 text-sm text-muted-foreground">{error.message || 'An unexpected error occurred'}</p>
          {#if resetError}
            <Button size="sm" variant="outline" class="mt-4" onclick={resetError}>
              Try again
            </Button>
          {/if}
        </div>
      </div>
    </Card.Content>
  </Card.Root>
{:else}
  <slot />
{/if}
