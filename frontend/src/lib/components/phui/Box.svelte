<script lang="ts">
  import type { Snippet } from 'svelte';

  let {
    border = false,
    children
  }: {
    border?: boolean;
    children: Snippet;
  } = $props();
</script>

<div class="box" class:bordered={border}>
  {@render children()}
</div>

<style>
  .box {
    background: var(--bg-card);
    border-radius: 4px;
    margin-bottom: 12px;
    overflow: hidden;
  }
  .bordered {
    border: 1px solid var(--border);
  }
</style>
