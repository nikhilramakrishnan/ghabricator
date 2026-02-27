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

<span class="tag tag-{shade}">
  {#if icon}
    <i class="fa {icon} mrs"></i>
  {/if}
  {@render children()}
</span>

<style>
  .tag {
    display: inline-flex;
    align-items: center;
    padding: 0 8px;
    border-radius: 3px;
    font-size: 12px;
    line-height: 22px;
    font-weight: 600;
    white-space: nowrap;
  }
  .tag-green { background: var(--tag-green-bg); color: var(--tag-green-text); }
  .tag-red { background: var(--tag-red-bg); color: var(--tag-red-text); }
  .tag-blue { background: var(--tag-blue-bg); color: var(--tag-blue-text); }
  .tag-orange { background: var(--tag-orange-bg); color: var(--tag-orange-text); }
  .tag-yellow { background: var(--tag-yellow-bg); color: var(--tag-yellow-text); }
  .tag-grey { background: var(--tag-grey-bg); color: var(--tag-grey-text); }
  .tag-violet { background: var(--tag-violet-bg); color: var(--tag-violet-text); }
  .tag-indigo { background: var(--tag-indigo-bg); color: var(--tag-indigo-text); }
</style>
