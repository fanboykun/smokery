<script lang="ts">
  import * as Select from '$lib/components/ui/select';
  import * as Avatar from '$lib/components/ui/avatar';
  import type { TeamMember } from '$lib/api/phase2';

  interface Props {
    members: TeamMember[];
    value: string;
    onchange: (userId: string) => void;
    disabled?: boolean;
  }

  let { members, value = $bindable(), onchange, disabled = false }: Props = $props();

  const selected = $derived(members.find((m) => m.id === value));
</script>

<Select.Root type="single" {disabled} bind:value onValueChange={(v: string) => onchange(v)}>
  <Select.Trigger class="w-full">
    {#if selected}
      <span class="flex items-center gap-2">
        <Avatar.Root size="sm">
          <Avatar.Image src={selected.avatar_url} alt={selected.name} />
          <Avatar.Fallback>{selected.name.slice(0, 2).toUpperCase()}</Avatar.Fallback>
        </Avatar.Root>
        <span class="truncate">{selected.name}</span>
      </span>
    {:else}
      <span class="text-muted-foreground">Assign to…</span>
    {/if}
  </Select.Trigger>
  <Select.Content>
    {#each members as member (member.id)}
      <Select.Item value={member.id}>
        <span class="flex items-center gap-2">
          <Avatar.Root size="sm">
            <Avatar.Image src={member.avatar_url} alt={member.name} />
            <Avatar.Fallback>{member.name.slice(0, 2).toUpperCase()}</Avatar.Fallback>
          </Avatar.Root>
          <span>{member.name}</span>
          <span class="text-xs text-muted-foreground">({member.role})</span>
        </span>
      </Select.Item>
    {/each}
  </Select.Content>
</Select.Root>
