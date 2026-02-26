<script lang="ts">
  type StatusItem = {
    name: string;
    icon: string;
    color: string;
    href?: string;
  };

  let { items }: { items: StatusItem[] } = $props();
</script>

<div>
  {#each items as item}
    <div style="display:flex;align-items:center;gap:8px;margin-bottom:6px;font-size:13px">
      <span class="phui-icon-view phui-font-fa {item.icon}" style="color:{item.color};width:16px;text-align:center"></span>
      {#if item.href}
        <a
          href={item.href}
          target="_blank"
          rel="noopener"
          style="flex:1;min-width:0;overflow:hidden;text-overflow:ellipsis;white-space:nowrap;text-decoration:none;color:inherit"
        >{item.name}</a>
      {:else}
        <span style="flex:1;min-width:0;overflow:hidden;text-overflow:ellipsis;white-space:nowrap">{item.name}</span>
      {/if}
    </div>
  {/each}
</div>
