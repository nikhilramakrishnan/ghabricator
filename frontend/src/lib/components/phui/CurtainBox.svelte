<script lang="ts">
  import type { Snippet } from 'svelte';

  let {
    title,
    children
  }: {
    title: string;
    children: Snippet;
  } = $props();
</script>

<div class="curtain-box">
  <div class="curtain-title">{title}</div>
  {@render children()}
</div>

<style>
  .curtain-box {
    margin-bottom: 16px;
  }
  .curtain-title {
    font-size: 11px;
    font-weight: 700;
    text-transform: uppercase;
    letter-spacing: 1px;
    color: var(--text-muted);
    margin-bottom: 8px;
    padding-bottom: 6px;
    border-bottom: 1px solid var(--border-subtle);
  }
</style>
