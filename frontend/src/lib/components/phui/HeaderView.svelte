<script lang="ts">
  let {
    title,
    icon = ''
  }: {
    title: string;
    icon?: string;
  } = $props();
</script>

<div class="header">
  <h1 class="header-title">
    {#if icon}
      <i class="fa {icon} header-icon"></i>
    {/if}
    {title}
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
</style>
