<script lang="ts">
  type ActionItem = {
    label: string;
    href?: string;
    icon?: string;
    onclick?: () => void;
  };

  let { actions }: { actions: ActionItem[] } = $props();
</script>

<ul class="action-list">
  {#each actions as action}
    <li class="action-item">
      {#if action.href}
        <a href={action.href} class="action-link">
          {#if action.icon}
            <i class="fa {action.icon} action-icon"></i>
          {/if}
          <span>{action.label}</span>
        </a>
      {:else}
        <button
          type="button"
          class="action-link action-btn"
          onclick={action.onclick}
        >
          {#if action.icon}
            <i class="fa {action.icon} action-icon"></i>
          {/if}
          <span>{action.label}</span>
        </button>
      {/if}
    </li>
  {/each}
</ul>

<style>
  .action-list {
    list-style: none;
    margin: 0;
    padding: 0;
  }
  .action-item {
    padding: 6px 0;
  }
  .action-link {
    display: flex;
    align-items: center;
    gap: 8px;
    color: var(--text-link);
    font-size: 13px;
    text-decoration: none;
  }
  .action-link:hover {
    text-decoration: underline;
  }
  .action-btn {
    background: none;
    border: none;
    cursor: pointer;
    width: 100%;
    text-align: left;
    padding: 0;
    font: inherit;
  }
  .action-icon {
    width: 16px;
    text-align: center;
    color: var(--text-muted);
    flex-shrink: 0;
  }
</style>
