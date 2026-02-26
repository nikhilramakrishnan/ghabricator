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

<li class="phui-oi-attribute">
  {#if icon}
    <span class="phui-icon-view phui-font-fa {icon} mrs"></span>
  {/if}
  {@render children()}
</li>
