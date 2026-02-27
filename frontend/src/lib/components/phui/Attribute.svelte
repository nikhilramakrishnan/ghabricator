<script lang="ts">
  import type { Snippet } from 'svelte';

  let {
    icon = '',
    children
  }: {
    icon?: string;
    children: Snippet;
  } = $props();
</script>

<span class="attr">
  {#if icon}
    <i class="fa {icon} attr-icon mrs"></i>
  {/if}
  {@render children()}
</span>

<style>
  .attr {
    display: inline-flex;
    align-items: center;
    font-size: 12px;
    color: var(--text-muted);
  }
  .attr-icon {
    color: var(--text-muted);
  }
</style>
