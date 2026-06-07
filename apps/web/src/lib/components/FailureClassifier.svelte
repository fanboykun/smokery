<script lang="ts">
  import * as Select from '$lib/components/ui/select';
  import { Badge } from '$lib/components/ui/badge';
  import { Wifi, FileWarning, ShieldAlert, Gauge, ServerCrash, HelpCircle } from '@lucide/svelte';

  interface Props {
    value: string;
    onchange: (classification: string) => void;
    disabled?: boolean;
  }

  let { value = $bindable(), onchange, disabled = false }: Props = $props();

  const classifications = [
    { id: 'network_timeout', label: 'Network Timeout', icon: Wifi, color: 'text-yellow-400' },
    { id: 'schema_mismatch', label: 'Schema Mismatch', icon: FileWarning, color: 'text-orange-400' },
    { id: 'auth_failure', label: 'Auth Failure', icon: ShieldAlert, color: 'text-red-400' },
    { id: 'rate_limit', label: 'Rate Limited', icon: Gauge, color: 'text-purple-400' },
    { id: 'server_error', label: 'Server Error', icon: ServerCrash, color: 'text-red-500' },
    { id: 'unknown', label: 'Unknown', icon: HelpCircle, color: 'text-muted-foreground' },
  ];

  const selected = $derived(classifications.find((c) => c.id === value));
</script>

<Select.Root type="single" {disabled} bind:value onValueChange={(v: string) => onchange(v)}>
  <Select.Trigger class="w-full">
    {#if selected}
      <span class="flex items-center gap-2">
        <selected.icon class={`size-4 ${selected.color}`} />
        {selected.label}
      </span>
    {:else}
      <span class="text-muted-foreground">Classify failure…</span>
    {/if}
  </Select.Trigger>
  <Select.Content>
    {#each classifications as cls (cls.id)}
      <Select.Item value={cls.id}>
        <span class="flex items-center gap-2">
          <cls.icon class={`size-4 ${cls.color}`} />
          {cls.label}
        </span>
      </Select.Item>
    {/each}
  </Select.Content>
</Select.Root>
