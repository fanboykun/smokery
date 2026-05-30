<script lang="ts">
  import { page } from '$app/stores';
  import { createQuery } from '@tanstack/svelte-query';
  import { mockGetProjectMembers } from '$lib/api/mock-phase2';
  import * as Card from '$lib/components/ui/card';
  import * as Avatar from '$lib/components/ui/avatar';
  import { Badge } from '$lib/components/ui/badge';
  import { Button } from '$lib/components/ui/button';
  import { UserPlus, Shield, Pencil, Eye } from '@lucide/svelte';

  const projectId = $page.params.id!;

  const members = createQuery(() => ({
    queryKey: ['project-members', projectId],
    queryFn: () => mockGetProjectMembers(projectId),
  }));

  function roleIcon(role: string) {
    if (role === 'admin') return Shield;
    if (role === 'editor') return Pencil;
    return Eye;
  }

  function roleColor(role: string) {
    if (role === 'admin') return 'text-primary';
    if (role === 'editor') return 'text-yellow-400';
    return 'text-muted-foreground';
  }
</script>

<main class="mx-auto max-w-4xl px-6 py-8 space-y-6">
  <div class="flex items-center justify-between">
    <div>
      <p class="text-xs font-bold uppercase tracking-widest text-primary">Settings</p>
      <h1 class="text-2xl font-bold">Members</h1>
    </div>
    <Button size="sm"><UserPlus class="size-4" /> Invite</Button>
  </div>

  {#if members.isPending}
    <p class="text-muted-foreground">Loading…</p>
  {:else if members.data}
    <Card.Root>
      <Card.Content class="divide-y divide-border p-0">
        {#each members.data as member (member.id)}
          {@const RoleIcon = roleIcon(member.role)}
          <div class="flex items-center gap-4 px-4 py-3">
            <Avatar.Root>
              <Avatar.Image src={member.avatar_url} alt={member.user_name} />
              <Avatar.Fallback>{member.user_name.slice(0, 2).toUpperCase()}</Avatar.Fallback>
            </Avatar.Root>
            <div class="flex-1 min-w-0">
              <p class="font-medium truncate">{member.user_name}</p>
              <p class="text-xs text-muted-foreground">{member.user_email}</p>
            </div>
            <Badge variant="outline" class={`gap-1 ${roleColor(member.role)}`}>
              <RoleIcon class="size-3" />
              {member.role}
            </Badge>
            <span class="text-xs text-muted-foreground">{new Date(member.added_at).toLocaleDateString()}</span>
          </div>
        {/each}
      </Card.Content>
    </Card.Root>
  {/if}
</main>
