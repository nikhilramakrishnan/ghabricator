<script lang="ts">
  import type { Snippet } from 'svelte';

  let {
    filetree,
    children
  }: {
    filetree?: Snippet;
    children: Snippet;
  } = $props();

  let collapsed = $state(true);
</script>

<div class="formation" class:collapsed>
  {#if filetree}
    <button
      class="toggle-show"
      onclick={() => collapsed = false}
      title="Show file tree"
    >
      <i class="fa fa-chevron-right"></i>
    </button>
    <div class="flank">
      <div class="flank-header">
        <span class="flank-title">File Tree</span>
        <button
          class="toggle-hide"
          onclick={() => collapsed = true}
          title="Hide file tree"
        >
          <i class="fa fa-chevron-left"></i>
        </button>
      </div>
      <div class="flank-body">
        {@render filetree()}
      </div>
    </div>
  {/if}
  <div class="center">
    <div class="main-content">
      {@render children()}
    </div>
  </div>
</div>

<style>
  .formation {
    display: flex;
    min-height: calc(100vh - 88px);
    position: relative;
  }
  .toggle-show {
    display: none;
    position: absolute;
    left: 0;
    top: 8px;
    z-index: 10;
    background: var(--bg-card);
    border: 1px solid var(--border);
    border-left: none;
    border-radius: 0 4px 4px 0;
    padding: 6px 8px;
    cursor: pointer;
    color: var(--text-muted);
    font-size: 12px;
  }
  .toggle-show:hover {
    color: var(--text);
    background: var(--bg-hover);
  }
  .collapsed .toggle-show {
    display: block;
  }
  .flank {
    width: 280px;
    min-width: 280px;
    flex-shrink: 0;
    border-right: 1px solid var(--border);
    background: var(--bg-card);
  }
  .collapsed .flank {
    display: none;
  }
  .flank-header {
    display: flex;
    align-items: center;
    justify-content: space-between;
    padding: 10px 12px;
    border-bottom: 1px solid var(--border-subtle);
    font-size: 13px;
    font-weight: 600;
    color: var(--text);
  }
  .flank-title {
    flex: 1;
  }
  .toggle-hide {
    background: none;
    border: none;
    cursor: pointer;
    color: var(--text-muted);
    padding: 2px 4px;
    font-size: 12px;
  }
  .toggle-hide:hover {
    color: var(--text);
  }
  .flank-body {
    overflow-y: auto;
    max-height: calc(100vh - 130px);
  }
  .center {
    flex: 1;
    min-width: 0;
    overflow-x: auto;
  }
  .main-content {
    padding: 0 16px;
  }
</style>
