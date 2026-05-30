<script lang="ts">
  import { Badge } from '$lib/components/ui/badge';
  import { CheckCircle2, AlertCircle, Clock, XCircle } from '@lucide/svelte';

  type Status = 'passed' | 'failed' | 'running' | 'pending' | 'error';

  interface Props {
    status: Status;
    label?: string;
    size?: 'sm' | 'md' | 'lg';
    class?: string;
  }

  let { status, label, size = 'md', class: cls = '' }: Props = $props();

  const statusConfig: Record<Status, { variant: any; icon: any; color: string; text: string }> = {
    passed: {
      variant: 'secondary',
      icon: CheckCircle2,
      color: 'text-emerald-400',
      text: 'Passed',
    },
    failed: {
      variant: 'destructive',
      icon: XCircle,
      color: 'text-red-400',
      text: 'Failed',
    },
    running: {
      variant: 'secondary',
      icon: Clock,
      color: 'text-blue-400',
      text: 'Running',
    },
    pending: {
      variant: 'secondary',
      icon: Clock,
      color: 'text-yellow-400',
      text: 'Pending',
    },
    error: {
      variant: 'destructive',
      icon: AlertCircle,
      color: 'text-red-500',
      text: 'Error',
    },
  };

  const config = $derived(statusConfig[status]);
  const Icon = $derived(config.icon);

  const sizeMap: Record<string, string> = {
    sm: 'size-3',
    md: 'size-4',
    lg: 'size-5',
  };
</script>

<div class={`inline-flex items-center gap-1.5 ${cls}`}>
  <Icon class={`${sizeMap[size]} ${config.color} animate-pulse`} strokeWidth={2} />
  <span class="text-xs font-medium">{label || config.text}</span>
</div>
