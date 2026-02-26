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

<div class="phui-info-view phui-info-severity-{severity}">
  {#if icon}
    <span class="phui-icon-view phui-font-fa {icon} mrs"></span>
  {/if}
  {@render children()}
</div>
