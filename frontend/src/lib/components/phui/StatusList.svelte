<script lang="ts">
  type StatusItem = {
    name: string;
    icon: string;
    color: string;
    href?: string;
  };

  let { items }: { items: StatusItem[] } = $props();
</script>

<div class="status-list">
  {#each items as item}
    <div class="status-item">
      <i class="fa {item.icon} status-icon" style="color:{item.color}"></i>
      {#if item.href}
        <a
          href={item.href}
          target="_blank"
          rel="noopener"
          class="status-link"
        >{item.name}</a>
      {:else}
        <span class="status-text">{item.name}</span>
      {/if}
    </div>
  {/each}
</div>

<style>
  .status-list {
    padding: 0;
  }
  .status-item {
    display: flex;
    align-items: center;
    gap: 8px;
    margin-bottom: 6px;
    font-size: 13px;
  }
  .status-icon {
    width: 16px;
    text-align: center;
    flex-shrink: 0;
  }
  .status-link {
    flex: 1;
    min-width: 0;
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
    text-decoration: none;
    color: inherit;
  }
  .status-link:hover {
    text-decoration: underline;
  }
  .status-text {
    flex: 1;
    min-width: 0;
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
  }
</style>
