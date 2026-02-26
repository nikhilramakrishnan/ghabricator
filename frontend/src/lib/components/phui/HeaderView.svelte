<script lang="ts">
  let {
    title,
    icon = ''
  }: {
    title: string;
    icon?: string;
  } = $props();
</script>

<div class="phui-header-shell">
  <div class="phui-header-view">
    <h1 class="phui-header-header">
      {#if icon}
        <span class="phui-header-icon phui-icon-view phui-font-fa {icon}"></span>
      {/if}
      {title}
    </h1>
  </div>
</div>
