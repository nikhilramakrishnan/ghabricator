<script lang="ts">
  import type { Snippet } from 'svelte';

  let {
    severity = 'nodata',
    icon = '',
    children
  }: {
    severity?: string;
    icon?: string;
    children: Snippet;
  } = $props();
</script>

<div class="info info-{severity}">
  {#if icon}
    <i class="fa {icon} info-icon mrs"></i>
  {/if}
  {@render children()}
</div>

<style>
  .info {
    background: var(--bg-subtle);
    border: 1px solid var(--border);
    border-radius: 4px;
    padding: 12px 16px;
    color: var(--text-muted);
    font-size: 13px;
    display: flex;
    align-items: flex-start;
    gap: 8px;
  }
  .info-icon {
    color: var(--text-muted);
    flex-shrink: 0;
    margin-top: 1px;
  }
</style>
