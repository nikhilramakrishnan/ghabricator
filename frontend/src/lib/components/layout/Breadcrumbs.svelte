<script lang="ts">
  let { crumbs }: { crumbs: { name: string; href?: string }[] } = $props();
</script>

<div class="ghabricator-crumbs" style="display:flex;align-items:center;gap:6px;padding:8px 16px;font-size:13px;color:#6B748C;background:#fff;border-bottom:1px solid #c7ccd9;">
  {#each crumbs as crumb, i}
    {#if i > 0}
      <span style="margin:0 2px;">&#x203a;</span>
    {/if}
    {#if crumb.href}
      <a href={crumb.href} style="color:#136CB2;text-decoration:none;">{crumb.name}</a>
    {:else}
      <span style="color:#464C5C;">{crumb.name}</span>
    {/if}
  {/each}
</div>
