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

<div class="phui-two-column-row grouped">
  <div class="phui-main-column">
    {@render children()}
  </div>
  {#if curtain}
    <div class="phui-side-column">
      {@render curtain()}
    </div>
  {/if}
</div>
