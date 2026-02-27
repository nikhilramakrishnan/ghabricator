<script lang="ts">
  import type { Snippet } from 'svelte';

  let {
    children,
    curtain
  }: {
    children: Snippet;
    curtain?: Snippet;
  } = $props();
</script>

<div class="curtain-layout">
  <div class="main-col">
    {@render children()}
  </div>
  {#if curtain}
    <div class="side-col">
      {@render curtain()}
    </div>
  {/if}
</div>

<style>
  .curtain-layout {
    display: flex;
    gap: 16px;
  }
  .main-col {
    flex: 1;
    min-width: 0;
  }
  .side-col {
    width: 280px;
    flex-shrink: 0;
  }
</style>
