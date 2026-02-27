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
  <a {href} class="btn btn-{color}" class:disabled>
    {#if icon}
      <i class="fa {icon} mrs"></i>
    {/if}
    {@render children()}
  </a>
{:else}
  <button {type} class="btn btn-{color}" {disabled} {onclick}>
    {#if icon}
      <i class="fa {icon} mrs"></i>
    {/if}
    {@render children()}
  </button>
{/if}

<style>
  .btn {
    display: inline-flex;
    align-items: center;
    padding: 6px 14px;
    border-radius: 3px;
    font-size: 12px;
    font-weight: 600;
    cursor: pointer;
    text-decoration: none;
    border: 1px solid var(--border);
    background: var(--bg-card);
    color: var(--text);
    font-family: inherit;
    line-height: 1.4;
  }
  .btn:hover {
    text-decoration: none;
    background: var(--bg-hover);
  }
  .btn:disabled, .btn.disabled {
    opacity: 0.5;
    cursor: not-allowed;
    pointer-events: none;
  }
  .btn-green {
    background: var(--green);
    color: #fff;
    border-color: var(--green);
  }
  .btn-green:hover { background: var(--green-hover); border-color: var(--green-hover); }
  .btn-red {
    background: var(--red);
    color: #fff;
    border-color: var(--red);
  }
  .btn-red:hover { background: var(--red-hover); border-color: var(--red-hover); }
  .btn-blue {
    background: var(--blue);
    color: #fff;
    border-color: var(--blue);
  }
  .btn-blue:hover { opacity: 0.9; }
</style>
