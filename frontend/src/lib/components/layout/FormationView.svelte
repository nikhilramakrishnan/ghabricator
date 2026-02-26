<script lang="ts">
  import type { Snippet } from 'svelte';

  let {
    filetree,
    children,
    curtain
  }: {
    filetree?: Snippet;
    children: Snippet;
    curtain?: Snippet;
  } = $props();

  let collapsed = $state(true);
</script>

<div class="phui-formation-wrapper" class:collapsed style="display:flex;min-height:calc(100vh - 88px);position:relative">
  {#if filetree}
    <button
      class="phui-flank-toggle-show"
      onclick={() => collapsed = false}
      title="Show file tree"
    >
      <span class="phui-icon-view phui-font-fa fa-chevron-right"></span>
    </button>
    <div class="phui-flank-view" style="width:280px;min-width:280px;flex-shrink:0;border-right:1px solid rgba(55,55,55,.3)">
      <div class="phui-flank-header">
        <span class="phui-flank-header-text">File Tree</span>
        <button
          class="phui-flank-toggle"
          onclick={() => collapsed = true}
          title="Hide file tree"
        >
          <span class="phui-icon-view phui-font-fa fa-chevron-left"></span>
        </button>
      </div>
      <div class="phui-flank-view-body" style="overflow-y:auto;max-height:calc(100vh - 130px)">
        {@render filetree()}
      </div>
    </div>
  {/if}
  <div style="flex:1;min-width:0;overflow-x:auto">
    <div class="phui-two-column-view">
      <div class="phui-two-column-container">
        <div class="phui-two-column-content">
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
        </div>
      </div>
    </div>
  </div>
</div>
