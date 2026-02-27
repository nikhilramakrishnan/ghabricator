<script lang="ts">
  let { crumbs }: { crumbs: { name: string; href?: string }[] } = $props();
</script>

<div class="breadcrumbs">
  {#each crumbs as crumb, i}
    {#if i > 0}
      <span class="sep">&#x203a;</span>
    {/if}
    {#if crumb.href}
      <a href={crumb.href} class="crumb-link">{crumb.name}</a>
    {:else}
      <span class="crumb-current">{crumb.name}</span>
    {/if}
  {/each}
</div>

<style>
  .breadcrumbs {
    display: flex;
    align-items: center;
    gap: 6px;
    padding: 8px 16px;
    font-size: 13px;
    color: var(--text-muted);
    background: var(--bg-card);
    border-bottom: 1px solid var(--border);
  }
  .sep {
    margin: 0 2px;
  }
  .crumb-link {
    color: var(--text-link);
    text-decoration: none;
  }
  .crumb-link:hover {
    text-decoration: underline;
  }
  .crumb-current {
    color: var(--text);
  }
</style>
