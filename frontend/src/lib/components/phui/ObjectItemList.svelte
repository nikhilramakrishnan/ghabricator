<script lang="ts">
  import type { Snippet } from 'svelte';

  let { children }: { children: Snippet } = $props();
</script>

<ul class="item-list">
  {@render children()}
</ul>

<style>
  .item-list {
    background: var(--bg-card);
    border-radius: 4px;
    overflow: hidden;
    list-style: none;
    margin: 0;
    padding: 0;
  }
</style>
