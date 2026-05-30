<script lang="ts">
  import * as Card from '$lib/components/ui/card';
  import { Button } from '$lib/components/ui/button';
  import { InboxIcon } from '@lucide/svelte';

  interface Props {
    title: string;
    description?: string;
    icon?: any;
    action?: {
      label: string;
      href?: string;
      onclick?: () => void;
    };
    class?: string;
  }

  let { title, description, icon: Icon = InboxIcon, action, class: cls = '' } = $props();
</script>

<Card.Root class={`border-dashed ${cls}`}>
  <Card.Content class="py-12 text-center">
    <div class="flex justify-center mb-4">
      <Icon class="size-8 text-muted-foreground" />
    </div>
    <p class="font-semibold text-foreground">{title}</p>
    {#if description}
      <p class="mt-1 text-sm text-muted-foreground">{description}</p>
    {/if}
    {#if action}
      <div class="mt-6">
        {#if action.href}
          <Button href={action.href}>{action.label}</Button>
        {:else if action.onclick}
          <Button onclick={action.onclick}>{action.label}</Button>
        {/if}
      </div>
    {/if}
    <slot />
  </Card.Content>
</Card.Root>
