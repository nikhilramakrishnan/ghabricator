<script lang="ts">
  import type { Snippet } from 'svelte';

  type PropertyItem = {
    label: string;
    value: string | Snippet;
  };

  let { items }: { items: PropertyItem[] } = $props();
</script>

<div class="mood-curtain-box">
  {#each items as item}
    <div class="mood-curtain-prop">
      <span class="mood-curtain-key">{item.label}</span>
      <span class="mood-curtain-val">
        {#if typeof item.value === 'string'}
          {item.value}
        {:else}
          {@render item.value()}
        {/if}
      </span>
    </div>
  {/each}
</div>
