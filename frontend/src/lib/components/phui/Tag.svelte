<script lang="ts">
  import type { Snippet } from 'svelte';

  let {
    shade = 'blue',
    icon = '',
    children
  }: {
    shade?: string;
    icon?: string;
    children: Snippet;
  } = $props();
</script>

<span class="phui-tag-view phui-tag-shade-{shade} phui-tag-type-shade">
  <span class="phui-tag-core">
    {#if icon}
      <span class="phui-icon-view phui-font-fa {icon} mrs"></span>
    {/if}
    {@render children()}
  </span>
</span>
