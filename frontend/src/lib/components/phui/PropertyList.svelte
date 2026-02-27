<script lang="ts">
  import type { Snippet } from 'svelte';

  type PropertyItem = {
    label: string;
    value: string | Snippet;
  };

  let { items }: { items: PropertyItem[] } = $props();
</script>

<div class="prop-list">
  {#each items as item}
    <div class="prop-row">
      <span class="prop-key">{item.label}</span>
      <span class="prop-val">
        {#if typeof item.value === 'string'}
          {item.value}
        {:else}
          {@render item.value()}
        {/if}
      </span>
    </div>
  {/each}
</div>

<style>
  .prop-list {
    padding: 0;
  }
  .prop-row {
    display: flex;
    gap: 8px;
    padding: 4px 0;
  }
  .prop-key {
    font-size: 12px;
    color: var(--text-muted);
    width: 80px;
    flex-shrink: 0;
  }
  .prop-val {
    font-size: 13px;
    color: var(--text);
  }
</style>
