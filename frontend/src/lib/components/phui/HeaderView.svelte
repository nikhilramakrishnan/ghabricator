<script lang="ts">
  let {
    title,
    icon = '',
    count
  }: {
    title: string;
    icon?: string;
    count?: number;
  } = $props();
</script>

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
</div>

<style>
  .header {
    background: var(--bg-card-header);
    padding: 6px 12px;
    border-bottom: 1px solid var(--border-subtle);
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
  .header-count {
    font-size: 11px;
    font-weight: 600;
    background: var(--bg-subtle);
    color: var(--text-muted);
    padding: 1px 6px;
    border-radius: 10px;
    margin-left: 2px;
  }
</style>
