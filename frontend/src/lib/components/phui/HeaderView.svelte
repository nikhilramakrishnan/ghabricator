<script lang="ts">
  import type { Snippet } from 'svelte';

  let {
    title,
    icon = '',
    count,
    collapsible = false,
    collapsed = false,
    onToggle,
    actions
  }: {
    title: string;
    icon?: string;
    count?: number;
    collapsible?: boolean;
    collapsed?: boolean;
    onToggle?: () => void;
    actions?: Snippet;
  } = $props();
</script>

{#if collapsible}
  <div class="header clickable" role="button" tabindex="0" onclick={onToggle} onkeydown={(e) => e.key === 'Enter' && onToggle?.()}>
    <h1 class="header-title">
      <i class="fa {collapsed ? 'fa-chevron-right' : 'fa-chevron-down'} toggle-icon"></i>
      {#if icon}
        <i class="fa {icon} header-icon"></i>
      {/if}
      {title}
      {#if count != null}
        <span class="header-count">{count}</span>
      {/if}
    </h1>
    {#if actions}
      <div class="header-actions" onclick={(e) => e.stopPropagation()}>
        {@render actions()}
      </div>
    {/if}
  </div>
{:else}
  <div class="header">
    <h1 class="header-title">
      {#if icon}
        <i class="fa {icon} header-icon"></i>
      {/if}
      {title}
      {#if count != null}
        <span class="header-count">{count}</span>
      {/if}
    </h1>
    {#if actions}
      <div class="header-actions">
        {@render actions()}
      </div>
    {/if}
  </div>
{/if}

<style>
  .header {
    background: var(--bg-card-header);
    padding: 6px 12px;
    border-bottom: 1px solid var(--border-subtle);
    display: flex;
    align-items: center;
    justify-content: space-between;
  }
  .header.clickable {
    cursor: pointer;
    user-select: none;
  }
  .header-title {
    font-size: 15px;
    font-weight: 600;
    color: var(--text);
    margin: 0;
    display: flex;
    align-items: center;
    gap: 8px;
  }
  .header-icon {
    color: var(--text-muted);
    font-size: 14px;
  }
  .toggle-icon {
    font-size: 12px;
    opacity: 0.6;
  }
  .header-count {
    font-size: 11px;
    font-weight: 600;
    background: var(--bg-subtle);
    color: var(--text-muted);
    padding: 1px 6px;
    border-radius: 10px;
    margin-left: 2px;
  }
  .header-actions {
    display: flex;
    align-items: center;
    gap: 8px;
  }
</style>
