<script lang="ts">
  import type { Snippet } from 'svelte';

  let {
    title = '',
    icon = '',
    headerRight,
    breadcrumbs,
    children
  }: {
    title?: string;
    icon?: string;
    headerRight?: Snippet;
    breadcrumbs?: Snippet;
    children: Snippet;
  } = $props();
</script>

<div class="page-shell">
  {#if breadcrumbs}
    {@render breadcrumbs()}
  {/if}
  {#if title}
    <div class="page-header">
      <h1 class="page-title">
        {#if icon}
          <i class="fa {icon} title-icon"></i>
        {/if}
        {title}
      </h1>
      {#if headerRight}
        <div class="header-right">
          {@render headerRight()}
        </div>
      {/if}
    </div>
  {/if}
  <div class="page-content">
    {@render children()}
  </div>
</div>

<style>
  .page-shell {
    padding: 0 16px;
  }
  .page-header {
    display: flex;
    align-items: center;
    justify-content: space-between;
    padding: 8px 0 6px;
  }
  .page-title {
    font-size: 18px;
    font-weight: 600;
    color: var(--text);
    margin: 0;
    display: flex;
    align-items: center;
    gap: 8px;
  }
  .title-icon {
    color: var(--text-muted);
    font-size: 16px;
  }
  .page-content {
    min-height: 0;
  }
</style>
