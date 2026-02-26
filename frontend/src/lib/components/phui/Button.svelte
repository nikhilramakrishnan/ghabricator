<script lang="ts">
  import type { Snippet } from 'svelte';

  let {
    color = 'default',
    icon = '',
    href = '',
    type = 'button',
    disabled = false,
    onclick,
    children
  }: {
    color?: string;
    icon?: string;
    href?: string;
    type?: 'button' | 'submit' | 'reset';
    disabled?: boolean;
    onclick?: (e: MouseEvent) => void;
    children: Snippet;
  } = $props();
</script>

{#if href}
  <a {href} class="mood-btn mood-btn-{color}" class:disabled>
    {#if icon}
      <span class="phui-icon-view phui-font-fa {icon} mrs"></span>
    {/if}
    {@render children()}
  </a>
{:else}
  <button {type} class="mood-btn mood-btn-{color}" {disabled} {onclick}>
    {#if icon}
      <span class="phui-icon-view phui-font-fa {icon} mrs"></span>
    {/if}
    {@render children()}
  </button>
{/if}
