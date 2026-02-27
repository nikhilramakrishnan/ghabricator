<script lang="ts">
  export interface Tab {
    key: string;
    label: string;
    icon?: string;
    count?: number;
  }

  let {
    tabs,
    active = $bindable('')
  }: {
    tabs: Tab[];
    active: string;
  } = $props();
</script>

<div class="tab-bar">
  {#each tabs as tab (tab.key)}
    <button
      class="tab"
      class:active={active === tab.key}
      onclick={() => active = tab.key}
    >
      {#if tab.icon}
        <i class="fa {tab.icon}"></i>
      {/if}
      {tab.label}
      {#if tab.count != null}
        <span class="count">{tab.count}</span>
      {/if}
    </button>
  {/each}
</div>

<style>
  .tab-bar {
    display: flex;
    gap: 0;
    border-bottom: 1px solid var(--border-subtle);
    background: var(--bg-card);
    padding: 0 12px;
  }

  .tab {
    display: flex;
    align-items: center;
    gap: 6px;
    padding: 8px 14px;
    border: none;
    background: none;
    color: var(--text-muted);
    font: inherit;
    font-size: 13px;
    font-weight: 500;
    cursor: pointer;
    position: relative;
    white-space: nowrap;
  }

  .tab::after {
    content: '';
    position: absolute;
    bottom: -1px;
    left: 0;
    right: 0;
    height: 2px;
    background: transparent;
  }

  .tab:hover {
    color: var(--text);
  }

  .tab.active {
    color: var(--blue);
    font-weight: 600;
  }

  .tab.active::after {
    background: var(--blue);
  }

  .count {
    font-size: 11px;
    font-weight: 600;
    background: var(--tag-grey-bg);
    color: var(--tag-grey-text);
    padding: 0 6px;
    border-radius: 10px;
    min-width: 18px;
    text-align: center;
  }

  .tab.active .count {
    background: var(--tag-blue-bg);
    color: var(--tag-blue-text);
  }
</style>
